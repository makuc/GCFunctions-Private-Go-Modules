.phony:
	encrypt secretenv

# BEGIN Config:

KeyRingName = "github-keyring"
CryptoKeyName = "github-token"

# END Config




sample:
	@echo OverwriteWithAccessToken> ${CryptoKeyName}

newkeyring:
	gcloud kms keyrings create ${KeyRingName} --location=global

newkey:
	gcloud kms keys create ${CryptoKeyName} --location=global --keyring=${KeyRingName} --purpose=encryption

# Standalone - The only one needed for Secret Env !!ss
secretenv:
	gcloud kms encrypt \
      --plaintext-file=${CryptoKeyName} \
      --ciphertext-file=- \
      --location=global \
      --keyring=${KeyRingName} \
      --key=${CryptoKeyName} | base64> encoded