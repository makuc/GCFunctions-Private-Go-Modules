.phony: deploy

# Config
prefix = "test"

#repo = "https://source.developers.google.com/projects/testing-192515/repos/diplomska-naloga"


# END Config

# Checks if you are deploying source from Repository or Local Filesystem.
ifeq (${repo}, .)
	source = .
else
	source = --no-source
endif

# Calling this deploys all function to the Google Cloud Functions
deploy:
	gcloud builds submit --config cloudbuild.yaml --substitutions=_PREFIX=${prefix} .

# Calling this deletes all functions from the Google Cloud Functions
delete:
	gcloud builds submit --config clouddelete.yaml --substitutions=_PREFIX=${prefix} --no-source
