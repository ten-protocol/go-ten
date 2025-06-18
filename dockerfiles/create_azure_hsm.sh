#!/bin/bash
set -e

AZURE_TENANT_ID="6ecc106d-6ece-4212-930d-289825fa8379"
AZURE_SUBSCRIPTION_ID="6e522d3f-f104-4fdf-bf88-76d2fe400cca"

# Configuration
TENANT_ID="$AZURE_TENANT_ID"
SUBSCRIPTION_ID="$AZURE_SUBSCRIPTION_ID"
RESOURCE_GROUP="Testnet"
HSM_NAME="EnclaveSigning"
KEY_NAME="EnclaveSigningKey3072"
LOCATION="East US"

echo "=== Setting up Azure Managed HSM with RSA 3072-bit key (exponent 3) ==="

# Login and set subscription
az login --tenant "$TENANT_ID"
az account set --subscription "$SUBSCRIPTION_ID"

# Get current user for HSM admin
USER_OBJECT_ID=$(az ad signed-in-user show --query id -o tsv)

# Create Managed HSM
echo "Creating Managed HSM: $HSM_NAME"
az keyvault create \
  --hsm-name "$HSM_NAME" \
  --resource-group "$RESOURCE_GROUP" \
  --location "$LOCATION" \
  --administrators "$USER_OBJECT_ID" \
  --retention-days 7

# Wait for HSM to be ready
echo "Waiting for HSM provisioning..."
az keyvault wait --hsm-name "$HSM_NAME" --created

# Create RSA key with exponent 3
echo "Creating RSA key with exponent 3..."
az rest \
  --method POST \
  --url "https://$HSM_NAME.managedhsm.azure.net/keys/$KEY_NAME/create?api-version=7.4" \
  --headers "Content-Type=application/json" \
  --body '{
    "kty": "RSA-HSM",
    "key_size": 3072,
    "public_exponent": 3,
    "key_ops": ["sign", "verify"]
  }'

echo "=== Setup completed successfully! ==="
echo "HSM Name: $HSM_NAME"
echo "Key Name: $KEY_NAME"
echo "HSM URL: https://$HSM_NAME.managedhsm.azure.net/"