#!/bin/sh
set -e

mkdir -p /dev/sgx
if [ ! -L /dev/sgx/enclave ]; then
	ln -s /dev/sgx_enclave /dev/sgx/enclave
fi

if [ -n "${PCCS_ADDR}" ]; then
	PCCS_URL=https://${PCCS_ADDR}/sgx/certification/v3/
fi

if [ -n "${PCCS_URL}" ]; then
	apt-get install -qq libsgx-dcap-default-qpl
	echo "PCCS_URL: ${PCCS_URL}"
	echo "PCCS_URL=${PCCS_URL}\nUSE_SECURE_CERT=FALSE" > /etc/sgx_default_qcnl.conf
else
	apt-get install -qq az-dcap-client
fi

"$@"