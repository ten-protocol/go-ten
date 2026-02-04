#!/bin/bash
set -euo pipefail

# Parameters
BINARY_PATH="${1:-main}"
SIGNER_SOURCE_PATH="${2:-/home/obscuro/go-obscuro/tools/enclavesigner/main}"
SIGNER_TOOL_PATH="/tmp/enclavesigner"

# Private key path inside the build container (file must exist in the build context)
LOCAL_SIGNING_KEY_PATH="${LOCAL_SIGNING_KEY_PATH:-/home/obscuro/go-obscuro/tools/enclavesigner/sgx_enclave_signing_key_rsa3072_e3.pem}"

TEMP_DIR="/tmp/sgx_local_signing"
SIG_B64_PATH="/tmp/signature.b64"
MOD_B64_PATH="/tmp/modulus.b64"

echo "====== LOCAL KEY SIGNATURE REPLACEMENT (FILE-BASED KEY, NO AZURE) ======"
echo "Using private key: $LOCAL_SIGNING_KEY_PATH"

echo "====== ENSURING DEPENDENCIES ======"
if ! command -v openssl >/dev/null 2>&1; then
  apt-get update
  apt-get install -y openssl ca-certificates
  apt-get clean && rm -rf /var/lib/apt/lists/*
fi

if [ ! -f "$LOCAL_SIGNING_KEY_PATH" ]; then
  echo "ERROR: Private key file not found at: $LOCAL_SIGNING_KEY_PATH"
  echo "Make sure you generated it on the host and that it is included in the Docker build context."
  exit 1
fi

echo "====== BUILDING ENCLAVE SIGNER TOOL ======"
if [ -d "$SIGNER_SOURCE_PATH" ]; then
  current_dir="$(pwd)"
  cd "$SIGNER_SOURCE_PATH"
  go build -o "$SIGNER_TOOL_PATH" .
  cd "$current_dir"
else
  echo "ERROR: Signer source path not found: $SIGNER_SOURCE_PATH"
  exit 1
fi

echo "====== STEP 1: EXTRACT HASH FROM EGO-SIGNED ENCLAVE ======"
cp "$BINARY_PATH" "$BINARY_PATH.original"
hash_b64="$("$SIGNER_TOOL_PATH" extract_hash "$BINARY_PATH")"
echo "$hash_b64" > /tmp/hash.b64
echo "Hash (Base64) extracted: $(echo "$hash_b64" | head -c 50)..."

echo "====== STEP 2: SIGN USING EXISTING PRIVATE KEY ======"
mkdir -p "$TEMP_DIR"

# Optional sanity checks (fail fast if key is wrong)
bits="$(openssl rsa -in "$LOCAL_SIGNING_KEY_PATH" -text -noout 2>/dev/null | awk '/Private-Key:/ {gsub("[()]", "", $2); print $2; exit}')"
exp="$(openssl rsa -in "$LOCAL_SIGNING_KEY_PATH" -text -noout 2>/dev/null | awk -F'[()]' '/publicExponent/ {print $2; exit}')"
if [ "${bits:-}" != "3072" ]; then
  echo "ERROR: Expected RSA-3072 key, got: ${bits:-unknown}"
  exit 1
fi
if [ "${exp:-}" != "3" ]; then
  echo "ERROR: Expected public exponent e=3, got: ${exp:-unknown}"
  exit 1
fi

# Extract modulus (n) and encode as Base64 (raw bytes)
modulus_hex="$(openssl rsa -in "$LOCAL_SIGNING_KEY_PATH" -modulus -noout 2>/dev/null | cut -d'=' -f2)"
printf '%s' "$modulus_hex" | xxd -r -p > "$TEMP_DIR/modulus.bin"
modulus_b64="$(base64 -w 0 < "$TEMP_DIR/modulus.bin" 2>/dev/null || base64 < "$TEMP_DIR/modulus.bin" | tr -d '\n')"

# Decode digest and sign using PKCS#1 v1.5 padding (raw RSA op)
printf '%s' "$hash_b64" | base64 -d > "$TEMP_DIR/digest.bin"
openssl rsautl -sign -pkcs -inkey "$LOCAL_SIGNING_KEY_PATH" -in "$TEMP_DIR/digest.bin" -out "$TEMP_DIR/signature.bin" >/dev/null 2>&1
signature_b64="$(base64 -w 0 < "$TEMP_DIR/signature.bin" 2>/dev/null || base64 < "$TEMP_DIR/signature.bin" | tr -d '\n')"

printf '%s' "$signature_b64" > "$SIG_B64_PATH"
printf '%s' "$modulus_b64" > "$MOD_B64_PATH"

if [ ! -s "$SIG_B64_PATH" ] || [ ! -s "$MOD_B64_PATH" ]; then
  echo "ERROR: Local signing failed (empty signature/modulus)."
  exit 1
fi

echo "Signature (Base64): $(echo "$signature_b64" | head -c 50)..."
echo "Modulus   (Base64): $(echo "$modulus_b64"   | head -c 50)..."

echo "====== STEP 3: REPLACE SIGNATURE IN ENCLAVE BINARY ======"
"$SIGNER_TOOL_PATH" replace "$BINARY_PATH" "$signature_b64" "$modulus_b64" "$BINARY_PATH.local_signed" 2>&1

echo "====== VERIFYING LOCALLY-SIGNED ENCLAVE ======"
"$SIGNER_TOOL_PATH" verify "$BINARY_PATH.local_signed" 2>&1 || {
  echo "ERROR: Locally-signed enclave verification FAILED!"
  exit 1
}

mv "$BINARY_PATH.local_signed" "$BINARY_PATH"
echo "Enclave signature successfully replaced with LOCAL key signature"

echo "====== FINAL VERIFICATION ======"
"$SIGNER_TOOL_PATH" verify "$BINARY_PATH" 2>&1

ls -la "$BINARY_PATH"
echo "Build completed successfully with LOCAL key signature integration"