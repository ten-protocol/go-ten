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

PCCS_URL=${PCCS_URL:-https://global.acccache.azure.net/sgx/certification/v4/}
COLLATERAL_SERVICE=${COLLATERAL_SERVICE:-https://api.trustedservices.intel.com/sgx/certification/v4/}

echo "PCCS_URL: ${PCCS_URL}"
echo "COLLATERAL_SERVICE: ${COLLATERAL_SERVICE}"

apt-get install -qq libsgx-dcap-default-qpl

echo "PCCS_URL=${PCCS_URL}\nCOLLATERAL_SERVICE=${COLLATERAL_SERVICE}\nUSE_SECURE_CERT=FALSE" > /etc/sgx_default_qcnl.conf

"$@"