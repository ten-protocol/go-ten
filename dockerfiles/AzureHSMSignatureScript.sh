#!/bin/bash

set -e

TENANT_ID="${AZURE_TENANT_ID}"
SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID}"
RESOURCE_GROUP="Testnet"
HSM_NAME="EnclaveSigning"
KEY_NAME="EnclaveSigningKey3072"
DIGEST_B64="$1"

echo "Authenticating with Azure..."
az login --tenant "$TENANT_ID" > /dev/null
az account set --subscription "$SUBSCRIPTION_ID"

echo "Getting key modulus from Managed HSM..."
az rest \
  --method GET \
  --url "https://$HSM_NAME.managedhsm.azure.net/keys/$KEY_NAME?api-version=7.4" \
  --output json > /tmp/key.json

modulus_b64url=$(jq -r '.key.n' /tmp/key.json)
echo "$modulus_b64url" > /tmp/modulus.b64

echo "Signing with Managed HSM..."
az rest \
  --method POST \
  --url "https://$HSM_NAME.managedhsm.azure.net/keys/$KEY_NAME/sign?api-version=7.4" \
  --headers "Content-Type=application/json" \
  --body "{\"alg\": \"RS256\", \"value\": \"$DIGEST_B64\"}" \
  --output json > /tmp/sign.json

signature_b64url=$(jq -r '.value' /tmp/sign.json)
echo "$signature_b64url" > /tmp/signature.b64

echo "Signing completed successfully!"
rm /tmp/key.json /tmp/sign.json