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

# If the PCCS_ADDR, the host is provided
# Do not use the default PCCS_URL defined in /etc/sgx_default_qcnl.conf
# Particularly used in Alibaba cloud
if [ -n "${PCCS_ADDR}" ]; then
	PCCS_URL=https://${PCCS_ADDR}/sgx/certification/v3/
fi

# Install the libsgx-dcap-default-qpl and redefine /etc/sgx_default_qcnl.conf (Alibaba)
if [ -n "${PCCS_URL}" ]; then
	apt-get install -qq libsgx-dcap-default-qpl
	echo "PCCS_URL: ${PCCS_URL}"
	echo "PCCS_URL=${PCCS_URL}\nUSE_SECURE_CERT=FALSE" > /etc/sgx_default_qcnl.conf
else
# Otherwise use the Azure library
	apt-get install -qq az-dcap-client
fi

"$@"