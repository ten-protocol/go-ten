#!/bin/bash

set -e

TENANT_ID="${AZURE_TENANT_ID}"
SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID}"
RESOURCE_GROUP="Testnet"
HSM_NAME="EnclaveSigningV2"
KEY_NAME="sgx-enclave-signing-key"
DIGEST_B64="$1"

# Function to convert Base64 to Base64URL
base64_to_base64url() {
    echo "$1" | tr '+/' '-_' | tr -d '='
}

# Function to convert Base64URL to Base64
base64url_to_base64() {
    local input="$1"
    # Replace URL-safe characters
    input=$(echo "$input" | tr -- '-_' '+/')

    # Add padding
    local padding=$((4 - ${#input} % 4))
    if [ $padding -ne 4 ]; then
        local pad_chars=$(printf "%*s" $padding | tr ' ' '=')
        input="${input}${pad_chars}"
    fi

    echo "$input"
}

echo "Authenticating with Azure..."
az login --tenant "$TENANT_ID" > /dev/null
az account set --subscription "$SUBSCRIPTION_ID"

echo "Getting key modulus from Managed HSM..."
az rest \
  --method GET \
  --url "https://$HSM_NAME.managedhsm.azure.net/keys/$KEY_NAME?api-version=7.4" \
  --resource "https://managedhsm.azure.net" \
  --output json > /tmp/key.json

modulus_b64url=$(jq -r '.key.n' /tmp/key.json)

# Convert input digest to Base64URL for Azure API
DIGEST_B64URL=$(base64_to_base64url "$DIGEST_B64")

echo "Original Base64: $DIGEST_B64"
echo "Converted Base64URL: $DIGEST_B64URL"

echo "Signing with Managed HSM..."
az rest \
  --method POST \
  --url "https://$HSM_NAME.managedhsm.azure.net/keys/$KEY_NAME/sign?api-version=7.4" \
  --resource "https://managedhsm.azure.net" \
  --headers "Content-Type=application/json" \
  --body "{\"alg\": \"RS256\", \"value\": \"$DIGEST_B64URL\"}" \
  --output json > /tmp/sign.json

signature_b64url=$(jq -r '.value' /tmp/sign.json)

# Convert Azure responses back to Base64 for Go application
signature_b64=$(base64url_to_base64 "$signature_b64url")
modulus_b64=$(base64url_to_base64 "$modulus_b64url")

# Save in Base64 format for Go application
echo "$signature_b64" > /tmp/signature.b64
echo "$modulus_b64" > /tmp/modulus.b64

echo "Signing completed successfully!"
echo "Signature (Base64): $signature_b64"
echo "Modulus (Base64): $modulus_b64"

rm /tmp/key.json /tmp/sign.json