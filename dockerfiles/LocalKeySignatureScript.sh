#!/bin/bash

set -e

TEMP_DIR="/tmp/sgx_local_signing"
DOCKER_SIGNATURE_PATH="/tmp/signature.b64"
DOCKER_MODULUS_PATH="/tmp/modulus.b64"

main() {
    DIGEST_B64="$1"

    mkdir -p "$TEMP_DIR"

    # Generate RSA key
    openssl genrsa -3 -out "$TEMP_DIR/private_key.pem" 3072 2>/dev/null

    # Extract modulus
    modulus_hex=$(openssl rsa -in "$TEMP_DIR/private_key.pem" -modulus -noout 2>/dev/null | cut -d'=' -f2)
    printf '%s' "$modulus_hex" | xxd -r -p > "$TEMP_DIR/modulus.bin"
    modulus_b64=$(base64 -w 0 < "$TEMP_DIR/modulus.bin" 2>/dev/null || base64 < "$TEMP_DIR/modulus.bin" | tr -d '\n')
    echo "modulus: $modulus_b64"

    # Sign
    printf '%s' "$DIGEST_B64" | base64 -d > "$TEMP_DIR/digest.bin"
    openssl rsautl -sign -in "$TEMP_DIR/digest.bin" -inkey "$TEMP_DIR/private_key.pem" -pkcs -out "$TEMP_DIR/signature.bin" 2>/dev/null
    signature_b64=$(base64 -w 0 < "$TEMP_DIR/signature.bin" 2>/dev/null || base64 < "$TEMP_DIR/signature.bin" | tr -d '\n')
    echo "sig: signature_b64"

    # Save outputs
    printf '%s' "$signature_b64" > "$DOCKER_SIGNATURE_PATH"
    printf '%s' "$modulus_b64" > "$DOCKER_MODULUS_PATH"
}

main "$@"