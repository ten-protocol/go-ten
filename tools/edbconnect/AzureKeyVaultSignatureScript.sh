#!/bin/bash

# Azure Key Vault Signature Script
# This script connects to Azure and makes a signature request to Azure Key Vault

set -e  # Exit on any error

# Configuration variables
TENANT_ID="${AZURE_TENANT_ID}"
SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID}"

RESOURCE_GROUP="Testnet"
KEY_VAULT_NAME="EnclaveSigning"
KEY_NAME="EnclaveSigningKey3072"

CLIENT_ID=""
CLIENT_SECRET=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# Function to check if Azure CLI is installed
check_azure_cli() {
    if ! command -v az &> /dev/null; then
        print_error "Azure CLI is not installed. Please install it first."
        print_status "Install instructions: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
        exit 1
    fi
    print_status "Azure CLI is installed"
}

# Function to check if jq is installed
check_jq() {
    if ! command -v jq &> /dev/null; then
        print_error "jq is not installed. Please install it first."
        print_status "Install with: brew install jq (macOS) or apt-get install jq (Ubuntu)"
        exit 1
    fi
    print_status "jq is installed"
}

# Function to authenticate with Azure
authenticate_azure() {
    print_status "Authenticating with Azure..."
    
    if [ -n "$CLIENT_ID" ] && [ -n "$CLIENT_SECRET" ] && [ -n "$TENANT_ID" ]; then
        print_status "Using service principal authentication..."
        az login --service-principal \
            --username "$CLIENT_ID" \
            --password "$CLIENT_SECRET" \
            --tenant "$TENANT_ID" > /dev/null
    else
        print_status "Using interactive authentication..."
        az login > /dev/null
    fi
    
    # Set the subscription
    az account set --subscription "$SUBSCRIPTION_ID"
    print_status "Set subscription to: $SUBSCRIPTION_ID"
}

# Function to verify Key Vault access
verify_key_vault_access() {
    print_status "Verifying Key Vault access..."
    
    # Check if Key Vault exists and is accessible
    if ! az keyvault show --name "$KEY_VAULT_NAME" --resource-group "$RESOURCE_GROUP" > /dev/null 2>&1; then
        print_error "Cannot access Key Vault '$KEY_VAULT_NAME' in resource group '$RESOURCE_GROUP'"
        print_warning "Please check:"
        print_warning "1. Key Vault name and resource group are correct"
        print_warning "2. You have appropriate permissions"
        print_warning "3. Key Vault exists"
        exit 1
    fi
    
    print_status "Key Vault '$KEY_VAULT_NAME' is accessible"
}

# Function to verify key exists
verify_key_exists() {
    print_status "Verifying key exists..."
    
    if ! az keyvault key show --vault-name "$KEY_VAULT_NAME" --name "$KEY_NAME" > /dev/null 2>&1; then
        print_error "Key '$KEY_NAME' not found in Key Vault '$KEY_VAULT_NAME'"
        print_warning "Available keys:"
        az keyvault key list --vault-name "$KEY_VAULT_NAME" --query "[].name" -o table
        exit 1
    fi
    
    print_status "Key '$KEY_NAME' exists in Key Vault"
}

# Function to create sample data to sign
check_input() {
    # If DIGEST_HEX is provided as a parameter, use it
    if [ -n "$1" ]; then
        DIGEST_HEX="$1"
        print_status "Using provided digest: $DIGEST_HEX"
    else
        print_error "Please provide a digest to sign"
        exit 1
    fi

    print_status "SHA-256 hash (base64): $DIGEST_HEX"
}

# Function to perform signature operation with detailed debugging
perform_signature() {
    print_status "Performing signature operation..."

    # Create temporary files for debugging
    local temp_output="/tmp/az_sign_output.json"
    local temp_error="/tmp/az_sign_error.log"

    print_debug "Signing digest: $DIGEST_HEX"
    print_debug "Using algorithm: RS256"

    # Execute the command with full output capture
    print_status "Executing Azure CLI sign command..."

    # First, let's try the command and capture all output
    local exit_code=0
    az keyvault key sign \
        --vault-name "$KEY_VAULT_NAME" \
        --name "$KEY_NAME" \
        --algorithm "RS256" \
        --digest "$DIGEST_HEX" \
        --output json > "$temp_output" 2> "$temp_error" || exit_code=$?

    print_debug "Command exit code: $exit_code"

    # Check if there were any errors
    if [ -s "$temp_error" ]; then
        print_warning "Command stderr output:"
        cat "$temp_error"
    fi

    # Check if output file exists and has content
    if [ ! -s "$temp_output" ]; then
        print_error "No output received from Azure CLI command"
        print_error "This could indicate:"
        print_error "1. Permission denied for sign operation"
        print_error "2. Key Vault access policy issues"
        print_error "3. Network connectivity problems"
        print_error "4. Azure CLI authentication issues"

        # Try to get more information
        print_status "Checking current Azure context..."
        az account show --query "{name:name, id:id, tenantId:tenantId}" -o table 2>/dev/null || print_warning "Could not get account info"

        # Check Key Vault permissions
        print_status "Checking Key Vault access policies..."
        local policies=$(az keyvault show --name "$KEY_VAULT_NAME" --resource-group "$RESOURCE_GROUP" --query "properties.accessPolicies" -o json 2>/dev/null || echo "[]")
        if [ "$policies" != "[]" ]; then
            echo "$policies" | jq -r '.[] | "Object ID: \(.objectId), Permissions: \(.permissions.keys // []), \(.permissions.secrets // []), \(.permissions.certificates // [])"' 2>/dev/null || print_warning "Could not parse access policies"
        else
            print_warning "No access policies found or unable to retrieve them"
        fi

        return 1
    fi

    print_debug "Raw Azure CLI output:"
    cat "$temp_output"

    # Parse the output
    local signature_result=""

    # Try different ways to extract the signature
    if command -v jq &> /dev/null; then
        signature_result=$(jq -r '.result // empty' "$temp_output" 2>/dev/null)

        if [ -z "$signature_result" ]; then
            # Try alternative field names
            signature_result=$(jq -r '.signature // .value // empty' "$temp_output" 2>/dev/null)
        fi
    fi

    # If jq failed or signature is still empty, try alternative parsing
    if [ -z "$signature_result" ]; then
        print_warning "jq parsing failed, trying alternative methods..."

        # Try grep-based extraction
        signature_result=$(grep -o '"result"[[:space:]]*:[[:space:]]*"[^"]*"' "$temp_output" 2>/dev/null | cut -d'"' -f4)

        if [ -z "$signature_result" ]; then
            signature_result=$(grep -o '"signature"[[:space:]]*:[[:space:]]*"[^"]*"' "$temp_output" 2>/dev/null | cut -d'"' -f4)
        fi
    fi

    # Check if we got a signature
    if [ -n "$signature_result" ] && [ "$signature_result" != "null" ]; then
        print_status "Signature operation completed successfully!"
        print_status "Signature (base64): $signature_result"

        # Validate signature format (should be base64)
        if echo "$signature_result" | base64 -d >/dev/null 2>&1; then
            print_status "Signature format validation: PASSED"

            # Get signature length
            local sig_length=$(echo -n "$signature_result" | wc -c)
            print_debug "Signature length: $sig_length characters"

        else
            print_warning "Signature format validation: FAILED (not valid base64)"
        fi

        # Save signature to file
        echo "$signature_result" > /tmp/signature.txt
        print_status "Signature saved to: /tmp/signature.txt"

        # Save all operation data for reference
        cat > /tmp/signature_info.json << EOF
{
    "originalMessage": "$SAMPLE_MESSAGE",
    "digestHex": "$DIGEST_HEX",
    "hashBase64": "$HASH_BASE64",
    "signature": "$signature_result",
    "algorithm": "RS256",
    "keyVault": "$KEY_VAULT_NAME",
    "keyName": "$KEY_NAME",
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF
        print_status "Operation details saved to: /tmp/signature_info.json"

        # Cleanup temp files
        rm -f "$temp_output" "$temp_error"
        return 0
    else
        print_error "Signature operation failed - empty or null result"
        print_error "Raw output content:"
        cat "$temp_output" 2>/dev/null || print_error "No output file"

        print_warning "Possible causes:"
        print_warning "1. Insufficient Key Vault permissions for 'sign' operation"
        print_warning "2. Key type incompatible with RS256 algorithm"
        print_warning "3. Azure Key Vault access policy restrictions"
        print_warning "4. Network or authentication issues"

        # Keep temp files for debugging
        print_status "Debug files preserved:"
        print_status "- Output: $temp_output"
        print_status "- Errors: $temp_error"

        return 1
    fi
}


# Function to verify signature (optional)
verify_signature() {
    print_status "Verifying signature..."
    
    # Get the public key for verification
    PUBLIC_KEY=$(az keyvault key show \
        --vault-name "$KEY_VAULT_NAME" \
        --name "$KEY_NAME" \
        --query "key" \
        --output json)
    
    print_status "Retrieved public key for verification"
    print_status "Note: Full signature verification requires additional processing of the public key"
}

# Function to cleanup temporary files
cleanup() {
    print_status "Cleaning up temporary files..."
    rm -f /tmp/sample_data.txt /tmp/signature.txt
}

# Function to display summary
display_summary() {
    print_status "=== Operation Summary ==="
    print_status "Key Vault: $KEY_VAULT_NAME"
    print_status "Key Name: $KEY_NAME"
    print_status "Resource Group: $RESOURCE_GROUP"
    print_status "Subscription: $SUBSCRIPTION_ID"
    print_status "Operation: Digital Signature"
    print_status "Algorithm: RS256"
    print_status "Status: SUCCESS"
}

# Main execution function
main() {
    local digest_param="$1"
    print_status "Starting Azure Key Vault Signature Script..."

    # Pre-flight checks
    check_azure_cli
    check_jq

    # Authenticate with Azure
    authenticate_azure

    # Verify access and resources
    verify_key_vault_access
    verify_key_exists

    # Perform signature operation
    check_input "$digest_param"
    if perform_signature; then
        verify_signature
        display_summary
    else
        print_error "Script execution failed"
        exit 1
    fi

    print_status "Script completed successfully!"
}

# Error handling
trap 'print_error "Script interrupted"; cleanup; exit 1' INT TERM

# Help function
show_help() {
    echo "Azure Key Vault Signature Script"
    echo ""
    echo "Usage: $0 [OPTIONS] [DIGEST_HEX]"
    echo ""
    echo "Environment Variables:"
    echo "  AZURE_SUBSCRIPTION_ID    - Azure subscription ID"
    echo "  AZURE_RESOURCE_GROUP     - Resource group name"
    echo "  AZURE_KEY_VAULT_NAME     - Key Vault name"
    echo "  AZURE_KEY_NAME           - Key name for signing"
    echo "  AZURE_TENANT_ID          - Tenant ID (for service principal auth)"
    echo "  AZURE_CLIENT_ID          - Client ID (for service principal auth)"
    echo "  AZURE_CLIENT_SECRET      - Client secret (for service principal auth)"
    echo ""
    echo "Arguments:"
    echo "  DIGEST_HEX               - Optional: Base64-encoded digest to sign"
    echo "                             If not provided, a sample message will be created"
    echo ""
    echo "Options:"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  # Using default sample message:"
    echo "  export AZURE_SUBSCRIPTION_ID='your-subscription-id'"
    echo "  export AZURE_KEY_VAULT_NAME='your-key-vault'"
    echo "  export AZURE_KEY_NAME='your-key-name'"
    echo "  export AZURE_RESOURCE_GROUP='your-resource-group'"
    echo "  $0"
    echo ""
    echo "  # Providing a specific digest to sign:"
    echo "  $0 'base64-encoded-digest-here'"
}

# Parse command line arguments
case "$1" in
    -h|--help)
        show_help
        exit 0
        ;;
    "")
        # No arguments, run with default sample data
        main
        ;;
    *)
        # First argument is the digest
        main "$1"
        ;;
esac