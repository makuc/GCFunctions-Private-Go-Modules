.phony:
	encrypt copy clean

# Config:
KeyRingName = "github-keyring"
CryptoKeyName = "github-token"

# END Config

sample:
	@echo OverwriteWithAccessToken> ${CryptoKeyName}

newkeyring:
	gcloud kms keyrings create ${KeyRingName} --location=global

newkey:
	gcloud kms keys create ${CryptoKeyName} --location=global --keyring=${KeyRingName} --purpose=encryption

secretenv:
	gcloud kms encrypt \
      --plaintext-file=${CryptoKeyName} \
      --ciphertext-file=- \
      --location=global \
      --keyring=${KeyRingName} \
      --key=${CryptoKeyName} | base64> encoded