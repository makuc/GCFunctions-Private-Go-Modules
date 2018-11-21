.phony:
	deploy

# Config
prefix = "test"

# END Config



# Calling this deploys all function to the Google Cloud Functions
deploy:
	gcloud builds submit --config cloudbuild.yaml --substitutions=_PREFIX=${prefix} .

# Calling this deletes all functions from the Google Cloud Functions
delete:
	gcloud builds submit --config clouddelete.yaml --substitutions=_PREFIX=${prefix} --no-source
