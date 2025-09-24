#!/bin/bash

set -e

# Check if Azure integration is needed
if [ -z "$AZURE_TENANT_ID" ]; then
    echo "Skipping Azure setup - AZURE_TENANT_ID not set"
    exit 0
fi

# Parameters
BINARY_PATH="${1:-main}"
SIGNER_SOURCE_PATH="${2:-/home/obscuro/go-obscuro/tools/enclavesigner/main}"
SIGNER_TOOL_PATH="/tmp/enclavesigner"

# Azure configuration
TENANT_ID="${AZURE_TENANT_ID}"
SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID}"
HSM_NAME="EnclaveSigningV2"
KEY_NAME="sgx-enclave-signing-key"

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

echo "====== INSTALLING AZURE CLI AND DEPENDENCIES ======"
apt-get update && apt-get install -y jq \
    ca-certificates \
    curl \
    apt-transport-https \
    lsb-release \
    gnupg

mkdir -p /etc/apt/keyrings
curl -sLS https://packages.microsoft.com/keys/microsoft.asc | \
    gpg --dearmor | \
    tee /etc/apt/keyrings/microsoft.gpg > /dev/null
chmod go+r /etc/apt/keyrings/microsoft.gpg

AZ_REPO=$(lsb_release -cs)
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ $AZ_REPO main" | \
    tee /etc/apt/sources.list.d/azure-cli.list

apt-get update && apt-get install -y azure-cli
az --version
echo "Azure CLI installation completed successfully"

# Clean up
apt-get clean && rm -rf /var/lib/apt/lists/*

echo "====== BUILDING ENCLAVE SIGNER TOOL ======"
if [ -d "$SIGNER_SOURCE_PATH" ]; then
    echo "Building enclavesigner from source: $SIGNER_SOURCE_PATH"
    current_dir=$(pwd)
    cd "$SIGNER_SOURCE_PATH"
    go build -o "$SIGNER_TOOL_PATH" .
    cd "$current_dir"
    echo "Enclavesigner built successfully at $SIGNER_TOOL_PATH"
else
    echo "ERROR: Signer source path not found: $SIGNER_SOURCE_PATH"
    exit 1
fi

echo "====== STARTING AZURE KEY VAULT SIGNATURE REPLACEMENT ======"

echo "====== STEP 1: EXTRACT HASH FROM EGO-SIGNED ENCLAVE ======"
cp "$BINARY_PATH" "$BINARY_PATH.original"
"$SIGNER_TOOL_PATH" extract_hash "$BINARY_PATH" > /tmp/hash.hex
hash_to_sign=$(cat /tmp/hash.hex)
echo "Hash to sign extracted: $hash_to_sign"

echo "====== STEP 2: AZURE KEY VAULT SIGNING ======"
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
DIGEST_B64URL=$(base64_to_base64url "$hash_to_sign")

echo "Original Base64: $hash_to_sign"
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

# Clean up Azure temporary files
rm /tmp/key.json /tmp/sign.json

# Validate signature files
if [ ! -f /tmp/signature.b64 ] || [ ! -s /tmp/signature.b64 ]; then
    echo "ERROR: Azure Key Vault signing failed - no signature file generated"
    exit 1
fi

if [ ! -f /tmp/modulus.b64 ] || [ ! -s /tmp/modulus.b64 ]; then
    echo "ERROR: Azure Key Vault signing failed - no modulus file generated"
    exit 1
fi

echo "Using signature: $(echo "$signature_b64" | head -c 50)..."
echo "Using modulus: $(echo "$modulus_b64" | head -c 50)..."

echo "====== STEP 3: REPLACE SIGNATURE IN ENCLAVE BINARY ======"
"$SIGNER_TOOL_PATH" replace "$BINARY_PATH" "$signature_b64" "$modulus_b64" "$BINARY_PATH.azure_signed" 2>&1

echo "====== VERIFYING AZURE-SIGNED ENCLAVE ======"
"$SIGNER_TOOL_PATH" verify "$BINARY_PATH.azure_signed" 2>&1 || {
    echo "ERROR: Azure-signed enclave verification FAILED!"
    exit 1
}

mv "$BINARY_PATH.azure_signed" "$BINARY_PATH"
echo "Enclave signature successfully replaced with Azure Key Vault signature"

echo "====== SIGNATURE REPLACEMENT COMPLETED ======"
echo "====== FINAL VERIFICATION ======"
"$SIGNER_TOOL_PATH" verify "$BINARY_PATH" 2>&1 || {
    echo "ERROR: Final enclave verification FAILED!"
    exit 1
}

ls -la "$BINARY_PATH"
echo "Build completed successfully with Azure Key Vault signature integration"