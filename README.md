# Private GitHub Go Modules made Easy

## Overview

This example shows how to easily deploy Cloud Functions written in Go using **private Go modules**. Achieved using custom **Cloud Build steps** to gain acces to _private_ Git repositories and _automatically_ **vendors** them in the cloud _during the deployment_.

It uses an **X-OAUTH-TOKEN** (from _GitHub_), which is securely **Encrypted** using **Cloud KMS** and after Decryption stored as an **ENV variable**.



### What this example does is:

* Helps you setup Cloud KMS token encryption
* It **vendors** all _functions_ (subdirectories) in the `/functions` dir
* Deletes `go.mod` and `go.sum` files for all _functions_ (subdirectories) in the `/functions` dir
* Uses the appropriate `gcloud` commands to deploy the specified functions as Cloud Build steps to deploy functions without `go.mod` files conflicting with `vendor`-ed dependencies.

The example uses the following products:
* Cloud Functions
* Cloud Build
* Cloud KMS



## Walkthrough

### Pre-requisites
* A Google Cloud Platform account with billing enabled



### Prepare GitHub OAuth Token and encrypt it on Cloud KMS

Prepare a file where you'll input your OAuth Token (the whole `prep` folder is in `.gitignore`, except for `Makefile`)

```console
$ cd prep
$ make sample
```

Replace content of `github-token` file (*including newline character* at the end...There can only be a single string) with your actual OAuth token.

To generate this token for **Github**, go here: [Personal access tokens](https://github.com/settings/tokens)
and click on "Generate new token". Check the checkbox next to *"repo"* to give it access to your private repositories and click on "Generate token" at the bottom of the page (you can also give it a description so you'll know what you'll be using it for).

After the sample file contains only your access token, execute:

```console
$ make newkeyring
$ make newkey
$ make secretenv
```

Open a newly created file `encoded` and copy its content to `cloudbuild.yaml` inside the root directory of your repository. Find the 4th line and replace the value with your encrypted string.

Now also replace `<PROJECT-ID>` in the 2nd line with your actual Project-ID from Google Cloud.

**Note:**
Default names for _Keyring_ and _Key_ are specified (and configurable) directly inside `Makefile`.


### Steps for Gaining Access to Private Repositories

To gain access and use the OAuth Token, you have to make proper `cloudbuild.yaml` steps:

```yaml
secrets:
- kmsKeyName: projects/<PROJECT-ID>/locations/global/keyRings/<KeyRingName>/cryptoKeys/<CryptoKeyName>
# Default `KeyRingName` and `CryptoKeyName`, specified in Makefile,
# are `github-keyring` and `github-token` respectively.
  secretEnv:
    GITHUB: "<EncodedAccessToken>"
```

This must be at the beginning of the `cloudbuild.yaml` file. It is used to Decrypt the earlier prepared OAuth Token needed to gain access to your Private Repository.

The decrypted value is stored in the ENV variable called `GITHUB`.

To configure access to our private repository, we have to tell _Git_ how to access it. For this purpose we have an _ash_ script (we are using alpine), where we configure _Git_ to use proper credentials when accessing private repositories. At the end, `privateRepoAccess` should look something like this:

```bash
#!/bin/ash

git config --global url."https://$GITHUB:x-oauth-basic@github.com/<username>/".insteadOf "https://github.com/<username>"
```

We configure _Git_ to modify all URLs trying to access Github repositories from `<username>` to use proper credentials to access them. This way, we can configure proper Github access at a _per user_ basis, in case multiple people are working on our project and they have their own private repositories hosting certain packages they are working on.

And now we must execute this script to actually use these configurations:

```yaml
steps:
# Load dependencies
- name: 'gcr.io/cloud-builders/go'
  entrypoint: 'ash' # because this is alpine
  args:
  - -c
  - |
    ash privateRepoAccess  # Exec script to gain Private Access
    cd functions           # Our functions are here
    for d in */ ; do       # Now vendor all of them
      cd $d
      pwd
      go mod tidy
      go mod vendor
      rm -f go.mod
      rm -f go.sum
      cd ..
    done
  secretEnv: ['GITHUB']
```

We use the Cloud Builder for Go with a *custom* entrypoint, since we want to execute some *ash* (bash alternative) commands.

First we execute the script to make proper _Git_ configurations, and then we navigate to the folder where all our functions (one function per subfolder) are located - `functions`.

Now all we have to do is go into each subfolder (which is an independent Cloud Function) and execute *Go* commands to *vendor* our modules ... First we **tidy** our `go.mod` file (you can remove this step) and then we **vendor** our dependencies.

Once our dependencies are vendored, we delete `go.mod` and `go.sum` files, since we don't want them to conflict with our vendored dependencies (`gcloud alpha functions deploy` gives preference to using *modules* to build when possible).

All the following steps can be used to freely deploy as many functions as you want without having to pay any more attention to whether you've **manually** *vendored* them or not!

Each Cloud Build step for Cloud Functions can be defined like this:

```yaml
- name: 'gcr.io/cloud-builders/gcloud'
  args:
  - alpha
  - functions
  - deploy
  - ${_PREFIX}-<function-name>
  - --trigger-http
  - --entry-point=BrezBaze
  - --runtime=go111
  - --memory=128MB
  dir: 'functions/noDBExp'
```

**Note:**
* ${_PREFIX} - used to give custom prefix to all functions when deploying them as a custom (probably temporary) group and you don't want it to interfere with their another, separate instance.
* <function-name> - Replace with your own name for the function
* dir: 'functions/<function-subdirectory>' - This is where you specify the ***root*** folder for this **individual** Cloud Function.


### Correct modules Import paths in files:

Since you'll be testing this in a Private Repository, don't forget to properly adjust import paths and/or module URLs in the following files:
* `go.mod` - module URL
* `modules/mytest/go.mod` - module URL
* `functions/anotherFunc/go.mod` - module URL
* `functions/anotherFunc/function.go` - import path
* `functions/noDBExp/go.mod` - module URL
* `functions/noDBExp/function.go` - import path

### Manually trigger a Cloud Build

To submit a command to build and deploy your group of functions can be easily done by executing the following command
(from the *root* dir of the repository):

```console
gcloud builds submit --config cloudbuild.yaml --substitutions=_PREFIX="myprefix" .
```

**Note:**

Notice the flag: `--substitutions=` - this is how we tell the Cloud Build what out _PREFIX is going to be during this build-time. By using the same PREFIX, one can easily delete the whole group of testing functions from the same stage (deployed using the same _PREFIX) by executing this command:

```console
gcloud builds submit --config cloudbuilddelete.yaml --substitutions=_PREFIX="myprefix" --no-source
```

**Note:**

If you want to remove _substitutions_ from this command, you'd also have to remove all its occurrences within `cloudbuild.yaml`.

Notice the `--no-source` flag.
This flag is used because we aren't compiling/building the app, only deleting the already existing Cloud Functions. This simply means we have no need for sending any source code, which saves us the time it normally takes to upload the source code archive and verify its contents.

Your functions are now online! Or removed ...

To simplify the process, you can also use:

```console
make deploy
make delete
```

Or make a BASH alternative, since it is pretty much the same...



### Set up the required permissions

Cloud Build doesn't, by default, have access to the Cloud Functions API within
your project. Before attempting a deployment, follow the
[Cloud Functions Deploying artifacts](https://cloud.google.com/cloud-build/docs/configuring-builds/build-test-deploy-artifacts#deploying_artifacts) instructions (steps 2 and 3, in particular).

Note that you will need your project number (not your project name/id). When looking
at the IAM page, there will typically only be one entry that matches
`[YOUR-PROJECT-NUMBER]@cloudbuild.gserviceaccount.com`.

Cloud Build also doesn't, by default, have access to the Cloud KMS API within your project.
Before attempting a deployment, follow the
[Grant the Cloud Build service account access to the CryptoKey](https://cloud.google.com/cloud-build/docs/securing-builds/use-encrypted-secrets-credentials#grant_the_product_name_short_service_account_access_to_the_cryptokey)

### Submit a build request

If the submitted build succeeded, you should now see a similar LOG to this one:

```console
gcloud builds submit --config cloudbuild.yaml --substitutions=_PREFIX="test" .
Creating temporary tarball archive of 20 file(s) totalling 6.4 KiB before compression.
Some files were not included in the source upload.

Check the gcloud log [<OMITTED>.log] to see which files and the contents of the
default gcloudignore file used (see `$ gcloud topic gcloudignore` to learn
more).

Uploading tarball of [.] to [gs://<PROJECT-ID>_cloudbuild/source/<OMITTED>.tgz]
Created [https://cloudbuild.googleapis.com/v1/projects/<PROJECT-ID>/builds/<OMITTED>].
Logs are available at [https://console.cloud.google.com/gcr/builds/<OMITTED>].
--------------------------------------------- REMOTE BUILD OUTPUT ------------------------------------------------
starting build "<OMITTED>"

FETCHSOURCE
Fetching storage object: gs://<PROJECT-ID>_cloudbuild/source/<OMITTED>.tgz#<OMITTED>
Copying gs://<PROJECT-ID>_cloudbuild/source/<OMITTED>.tgz#<OMITTED>...
/ [1 files][  3.1 KiB/  3.1 KiB]
Operation completed over 1 objects/3.1 KiB.
BUILD
Starting Step #0
Step #0: Already have image (with digest): gcr.io/cloud-builders/go
Step #0: <OMITTED>
Step #0: go: finding github.com/<USERNAME>/<REPO-NAME>/modules/mytest v0.0.0-<OMITTED>-<OMITTED>
Step #0: go: downloading github.com/<USERNAME>/<REPO-NAME>/modules/mytest v0.0.0-<OMITTED>-<OMITTED>
Finished Step #0
Starting Step #1
Step #1: Already have image (with digest): gcr.io/cloud-builders/gcloud
Step #1: Deploying function (may take a while - up to 2 minutes)...
Step #1: ..................done.
Step #1: availableMemoryMb: 128
Step #1: entryPoint: BrezBaze
Step #1: httpsTrigger:
Step #1:   url: https://<REGION>-<PROJECT-ID>.cloudfunctions.net/test-func1
Step #1: labels:
Step #1:   deployment-tool: cli-gcloud
Step #1: name: projects/<PROJECT-ID>/locations/<REGION>/functions/test-func1
Step #1: runtime: go111
Step #1: serviceAccountEmail: <PROJECT-ID>@appspot.gserviceaccount.com
Step #1: sourceUploadUrl: https://storage.googleapis.com/gcf-upload-<REGION>-<OMITTED>
Step #1: status: ACTIVE
Step #1: timeout: 60s
Step #1: updateTime: '2018-11-21T13:41:28Z'
Step #1: versionId: '29'
Finished Step #1
PUSH
DONE
------------------------------------------------------------------------------------------------------------------

ID         CREATE_TIME                DURATION  SOURCE                                             IMAGES  STATUS
<OMITTED>  2018-11-21T13:40:49+00:00  41S       gs://<PROJECT-ID>_cloudbuild/source/<OMITTED>.tgz  -       SUCCESS
```

### Test that the deployment worked

Use `curl` to send a request to your function. You can find the function's
endpoint in the Cloud Build logs.

```console
$ curl https://<REGION>-<PROJECT_NAME>.cloudfunctions.net/test-func1
Yo, World!
```