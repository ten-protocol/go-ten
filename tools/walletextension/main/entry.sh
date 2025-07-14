#!/bin/sh
set -e
#
# This script is the entry point for starting the enclave under a Docker container.
# It allows running SGX sdk using different parameters.
#

# It's expected to be a link between the /dev/sgx_enclave Docker device and the container /dev/sgx/enclave
mkdir -p /dev/sgx
if [ ! -L /dev/sgx/enclave ]; then
	ln -s /dev/sgx_enclave /dev/sgx/enclave
fi

# Set and export PCCS_URL environment variable
export PCCS_URL=https://global.acccache.azure.net/sgx/certification/v4/
echo "PCCS_URL: ${PCCS_URL}"

# Export other important SGX environment variables if they're not already set
if [ -z "$AESM_PATH" ]; then
    export AESM_PATH=/var/run/aesmd/aesm.socket
fi

# Debug: Show key environment variables
echo "Environment variables:"
echo "OE_SIMULATION: ${OE_SIMULATION:-not set}"
echo "PCCS_URL: ${PCCS_URL:-not set}"
echo "AESM_PATH: ${AESM_PATH:-not set}"
echo "SGX_MODE: ${SGX_MODE:-not set}"

apt-get install -qq libsgx-dcap-default-qpl

echo "PCCS_URL=${PCCS_URL}\nUSE_SECURE_CERT=FALSE" > /etc/sgx_default_qcnl.conf

"$@"