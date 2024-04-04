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

# Todo - pass this in as a parameter
PCCS_URL=https://global.acccache.azure.net/sgx/certification/v3/

# Install the libsgx-dcap-default-qpl and redefine /etc/sgx_default_qcnl.conf (Alibaba)
apt-get install -qq libsgx-dcap-default-qpl
echo "PCCS_URL: ${PCCS_URL}"

echo '{
  "pccs_url": "https://global.acccache.azure.net/sgx/certification/v3/",
  "use_secure_cert": false,
  "collateral_service": "https://global.acccache.azure.net/sgx/certification/v3/",
  "pccs_api_version": "3.1",
  "retry_times": 6,
  "retry_delay": 5,
  "local_pck_url": "http://169.254.169.254/metadata/THIM/sgx/certification/v3/",
  "pck_cache_expire_hours": 24,
  "verify_collateral_cache_expire_hours": 24,
  "custom_request_options": {
      "get_cert": {
          "headers": {
              "metadata": "true"
          },
          "params": {
              "api-version": "2021-07-22-preview"
          }
      }
  }
}'  > /etc/sgx_default_qcnl.conf

"$@"