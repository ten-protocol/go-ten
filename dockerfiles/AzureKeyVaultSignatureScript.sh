#!/bin/bash

# Azure Key Vault Signature Script
# This script connects to Azure and makes a signature request to Azure Key Vault
# Returns signature and modulus in base64 format for SGX enclave signing

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
        exit 1
    fi
    print_status "Azure CLI is installed"
}

# Function to check if jq is installed
check_jq() {
    if ! command -v jq &> /dev/null; then
        print_error "jq is not installed. Please install it first."
        exit 1
    fi
    print_status "jq is installed"
}

# Function to base64url decode (Azure Key Vault format) and convert to regular base64
base64url_to_base64() {
    local input="$1"
    # Replace URL-safe characters first
    input=$(echo "$input" | tr '_-' '/+')
    # Add padding if needed
    local mod=$((${#input} % 4))
    if [ $mod -ne 0 ]; then
        input="${input}$(printf '%*s' $((4 - mod)) '' | tr ' ' '=')"
    fi
    echo "$input"
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

    if ! az keyvault show --name "$KEY_VAULT_NAME" --resource-group "$RESOURCE_GROUP" > /dev/null 2>&1; then
        print_error "Cannot access Key Vault '$KEY_VAULT_NAME' in resource group '$RESOURCE_GROUP'"
        exit 1
    fi

    print_status "Key Vault '$KEY_VAULT_NAME' is accessible"
}

# Function to get modulus from public key
get_modulus_base64() {
    print_status "Retrieving RSA modulus..."

    local temp_key_output="/tmp/azure_key_info.json"

    if ! az keyvault key show --vault-name "$KEY_VAULT_NAME" --name "$KEY_NAME" --output json > "$temp_key_output" 2>/dev/null; then
        print_error "Key '$KEY_NAME' not found in Key Vault '$KEY_VAULT_NAME'"
        exit 1
    fi

    # Extract modulus from public key
    local modulus_b64url=$(jq -r '.key.n' "$temp_key_output" 2>/dev/null)

    if [ -z "$modulus_b64url" ] || [ "$modulus_b64url" = "null" ]; then
        print_error "Failed to extract RSA modulus from public key"
        exit 1
    fi

    # Convert base64url to standard base64
    local modulus_b64=$(base64url_to_base64 "$modulus_b64url")

    # Save modulus
    echo "$modulus_b64" > /tmp/modulus.b64
    print_status "Modulus (base64) saved to: /tmp/modulus.b64"
    print_debug "Modulus: $modulus_b64"

    rm -f "$temp_key_output"
}
# Function to check input
check_input() {
    if [ -n "$1" ]; then
        DIGEST_B64="$1"
        print_status "Using provided digest (base64): $DIGEST_B64"

        # Validate base64 format
        if ! echo "$DIGEST_B64" | base64 -d > /dev/null 2>&1; then
            print_error "Invalid base64 digest format"
            exit 1
        fi

        # Convert to hex for debugging
        local digest_hex=$(echo "$DIGEST_B64" | base64 -d | xxd -p -c 32)
        print_debug "Digest in hex: $digest_hex"
    else
        print_error "Please provide a digest to sign"
        exit 1
    fi
}

# Function to perform signature operation
perform_signature() {
    print_status "Performing signature operation..."

    local temp_output="/tmp/az_sign_output.json"
    local temp_error="/tmp/az_sign_error.log"

    print_debug "Signing digest: $DIGEST_B64"

    # Execute Azure CLI sign command with base64 digest
    print_debug "Executing Azure CLI sign command..."
    local exit_code=0
    az keyvault key sign \
        --vault-name "$KEY_VAULT_NAME" \
        --name "$KEY_NAME" \
        --algorithm "RS256" \
        --digest "$DIGEST_B64" \
        --output json > "$temp_output" 2> "$temp_error" || exit_code=$?

    print_debug "Azure CLI exit code: $exit_code"

    # Check if we have any output
    if [ -s "$temp_error" ]; then
        print_error "Azure CLI error output:"
        cat "$temp_error"
    fi

    if [ ! -s "$temp_output" ]; then
        print_error "No output received from Azure CLI command"
        return 1
    fi

    print_debug "Raw Azure response:"
    cat "$temp_output"

    # Extract signature - Azure returns it in the 'signature' field
    local signature_b64url=""
    if command -v jq &> /dev/null; then
        signature_b64url=$(jq -r '.signature' "$temp_output" 2>/dev/null)
    fi

    if [ -z "$signature_b64url" ] || [ "$signature_b64url" = "null" ]; then
        print_error "Failed to extract signature from Azure response"
        print_debug "Available fields in response:"
        jq 'keys[]' "$temp_output" 2>/dev/null || echo "No valid JSON found"
        return 1
    fi

    print_debug "Raw signature from Azure (base64url): $signature_b64url"

    # Validate we can decode it (just for verification)
    local test_signature_b64=$(base64url_to_base64 "$signature_b64url")
    local signature_bytes=""
    if signature_bytes=$(echo "$test_signature_b64" | base64 -d 2>/dev/null); then
        local sig_length=${#signature_bytes}
        print_debug "Signature length: $sig_length bytes (will be padded to 384 by Go if needed)"

        # Remove the length check - let Go handle padding
        if [ $sig_length -lt 380 ] || [ $sig_length -gt 384 ]; then
            print_error "Signature length seems invalid: got $sig_length bytes"
            return 1
        fi
    else
        print_error "Invalid signature format"
        return 1
    fi

    print_status "✓ Signature operation completed successfully!"

    # Save signature in base64url format (as received from Azure)
    echo "$signature_b64url" > /tmp/signature.b64
    print_status "Signature (base64url) saved to: /tmp/signature.b64"

    rm -f "$temp_output" "$temp_error"
    return 0
}

# Function to get modulus from public key
get_modulus_base64() {
    print_status "Retrieving RSA modulus..."

    local temp_key_output="/tmp/azure_key_info.json"

    if ! az keyvault key show --vault-name "$KEY_VAULT_NAME" --name "$KEY_NAME" --output json > "$temp_key_output" 2>/dev/null; then
        print_error "Key '$KEY_NAME' not found in Key Vault '$KEY_VAULT_NAME'"
        exit 1
    fi

    # Extract modulus from public key (in base64url format)
    local modulus_b64url=$(jq -r '.key.n' "$temp_key_output" 2>/dev/null)

    if [ -z "$modulus_b64url" ] || [ "$modulus_b64url" = "null" ]; then
        print_error "Failed to extract RSA modulus from public key"
        exit 1
    fi

    # No conversion needed - pass base64url directly to Go
    # Save modulus in base64url format (as received from Azure)
    echo "$modulus_b64url" > /tmp/modulus.b64
    print_status "Modulus (base64url) saved to: /tmp/modulus.b64"
    print_debug "Modulus: $modulus_b64url"

    rm -f "$temp_key_output"
}

# Function to display summary
display_summary() {
    print_status "=== Operation Summary ==="
    print_status "Key Vault: $KEY_VAULT_NAME"
    print_status "Key Name: $KEY_NAME"
    print_status "Files Created:"
    print_status "  - Signature: /tmp/signature.b64 (base64)"
    print_status "  - Modulus: /tmp/modulus.b64 (base64)"
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

    # Verify access and get modulus
    verify_key_vault_access
    get_modulus_base64

    # Perform signature operation
    check_input "$digest_param"
    if perform_signature; then
        display_summary
        print_status "Script completed successfully!"
        return 0
    else
        print_error "Script execution failed"
        return 1
    fi
}

# Help function
show_help() {
    echo "Azure Key Vault Signature Script"
    echo ""
    echo "Usage: $0 [DIGEST_HEX]"
    echo ""
    echo "This script connects to Azure Key Vault to:"
    echo "1. Retrieve the RSA modulus in base64 format"
    echo "2. Sign the provided digest"
    echo "3. Output signature in base64 format"
    echo ""
    echo "Output Files:"
    echo "  /tmp/signature.b64  - Signature in base64"
    echo "  /tmp/modulus.b64    - RSA modulus in base64"
    echo ""
    echo "Arguments:"
    echo "  DIGEST_HEX         - Base64-encoded SHA-256 digest to sign"
}

# Parse command line arguments
case "$1" in
    -h|--help)
        show_help
        exit 0
        ;;
    "")
        print_error "Digest parameter is required"
        show_help
        exit 1
        ;;
    *)
        main "$1"
        ;;
esac