# Private GitHub Go Modules made Easy

## Overview

This example shows how to easily deploy Cloud Functions written in Go which are using *private Go modules*.
It achieves this using custom Cloud Build steps to gain acces to private Git repositories.
It uses an X-OAUTH-TOKEN (from GitHub), which is securely Encrypted using Cloud KMS and after Decryption stored as an ENV.

What this example does is:
* Helps you setup Cloud KMS token encryption
* Sets Cloud Build Git *credentials.helper* to *cache* with a timeout of 5 min
* During this time it *vendor's* all functions (subdirectories) in the `functions` dir
* Deletes `go.mod` and `go.sum` files for all functions (subdirectories) in the `functions` dir
* Continues the build process in the following Cloud Build steps which can freely use *gcloud alpha functions deploy* commands without conflicts.

The example uses the following products:
* Cloud Functions
* Cloud Build

## Walkthrough

### Pre-requisites
* A Google Cloud Platform account with billing enabled

### Prepare GitHub OAuth Token and encrypt it on Cloud KMS

Prepare a file where you'll input your OAuth Token (don't worry, this whole folder is in .gitignore)
```console
$ cd prep
$ make sample
```

Replace content of the file (*including newline character* at the end...There can only be a single string) `github-token`
with your actual token.

To generate this token, go here: [Personal access tokens](https://github.com/settings/tokens)
and click on "Generate new token". Check the checkbox next to *"repo"* to give it full access to your repositories
and click on "Generate token" at the bottom of the page
(you can also give it a description so you'll know what you'll be using it for).

After the sample file contains only your access token, execute:

```console
$ make newkeyring
$ make newkey
$ make secretenv
```

Open a newly created file `encoded` and copy its content to `cloudbuild.yaml` inside the root directory of your repository.
Look at the 4th line and replace <EncodedAccessToken> with your encoded/encrypted string.

Now also replace <PROJECT-ID> on the 2nd line with your actual Project-ID from your Google Cloud.

And make sure you are *cloning* a correct repository (to save credentials) by adjusting URL on line 17
(replace <USERNAME> and <REPO-NAME> with you actual values)

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

Note:
Notice the flag: `--substitutions=` - this is how we tell the Cloud Build what out _PREFIX is going to be during this build-time.
By using the same PREFIX, one can easily delete the whole group of testing functions from the same stage (deployed using the same _PREFIX) by executing this command:
If you want to remove this option from this command, you'd also have to remove all its occurrences within `cloudbuild.yaml`.

```console
gcloud builds submit --config cloudbuilddelete.yaml --substitutions=_PREFIX="myprefix" --no-source
```

Note:
Notice the `--no-source` flag.
This flag is used because we aren't compiling/building the app, only deleting the already existing Cloud Functions.
This simply means we have no need for sending any source code,
which saves us the time it normally takes to upload the source code archive and verify its content.

Your functions are now online!

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
---------------------------------------------------------------------------- REMOTE BUILD OUTPUT -----------------------------------------------------------------------------
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
Step #0: Cloning into '<REPO-NAME>'...
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
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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
