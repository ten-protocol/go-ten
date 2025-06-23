#!/bin/bash
set -e

AZURE_TENANT_ID="6ecc106d-6ece-4212-930d-289825fa8379"
AZURE_SUBSCRIPTION_ID="6e522d3f-f104-4fdf-bf88-76d2fe400cca"

# Configuration variables
TENANT_ID="${AZURE_TENANT_ID}"
SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID}"
RESOURCE_GROUP="Testnet_HSM"
LOCATION="East US"  # Change to your preferred region
HSM_NAME="EnclaveSigningV2"
KEY_NAME="EnclaveSigningKey3072"
ADMIN_OBJECT_ID="${AZURE_ADMIN_OBJECT_ID}"  # Your Azure AD user/service principal object ID

echo "Authenticating with Azure..."
az login --tenant "$TENANT_ID" > /dev/null
az account set --subscription "$SUBSCRIPTION_ID"

echo "Getting current user's object ID..."
ADMIN_OBJECT_ID=$(az ad signed-in-user show --query id --output tsv)
echo "Using Object ID: $ADMIN_OBJECT_ID"

echo "Creating resource group if it doesn't exist..."
az group create --name "$RESOURCE_GROUP" --location "$LOCATION" --output none

echo "Creating Managed HSM (this may take 10-15 minutes)..."
az keyvault create \
  --hsm-name "$HSM_NAME" \
  --resource-group "$RESOURCE_GROUP" \
  --location "$LOCATION" \
  --administrators "$ADMIN_OBJECT_ID" \
  --retention-days 7

echo "Waiting for HSM to be fully provisioned..."
sleep 60

echo "Activating the HSM security domain..."
az keyvault security-domain download \
  --hsm-name "$HSM_NAME" \
  --sd-wrapping-keys ./sd-wrapping-key.json \
  --sd-quorum 1 \
  --security-domain-file ./security-domain.json

echo "Creating RSA key with exponent 3 for SGX..."
az rest \
  --method POST \
  --url "https://$HSM_NAME.managedhsm.azure.net/keys/$KEY_NAME/create?api-version=7.4" \
  --headers "Content-Type=application/json" \
  --body '{
    "kty": "RSA-HSM",
    "key_size": 3072,
    "public_exponent": "Aw",
    "key_ops": ["sign", "verify"],
    "attributes": {
      "enabled": true,
      "exportable": false
    }
  }' \
  --output json > /tmp/key_creation.json

echo "Key creation response:"
cat /tmp/key_creation.json | jq '.'

echo "Verifying key was created with exponent 3..."
az rest \
  --method GET \
  --url "https://$HSM_NAME.managedhsm.azure.net/keys/$KEY_NAME?api-version=7.4" \
  --output json > /tmp/key_info.json

echo "Key information (should show exponent 3):"
cat /tmp/key_info.json | jq '.key | {kty, key_size, e, n}'

# Clean up temporary files
rm -f /tmp/key_creation.json /tmp/key_info.json

echo "Azure Managed HSM setup completed successfully!"
echo "HSM Name: $HSM_NAME"
echo "Key Name: $KEY_NAME"
echo "Exponent: 3 (for SGX compatibility)"
echo "HSM URL: https://$HSM_NAME.managedhsm.azure.net/"