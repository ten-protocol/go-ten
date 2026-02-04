Utility to verify and replace the signature of enclave binaries.

We use it during our build process to replace the signature on the enclave binaries with one generated securely.


## Signing using local key

1. First generate the key using `GenerateImageSigningKey.sh`.
This will create a file ``sgx_enclave_signing_key_rsa3072_e3.pem`` in the current directory.
2. If you already have the key , just copy it to the local file system.
3. Build the enclave binaries using 
```DOCKER_BUILDKIT=1 docker build \
-f dockerfiles/enclave.Dockerfile \
--secret id=sgx_signing_key,src=/absolute/path/to/sgx_enclave_signing_key_rsa3072_e3.pem \
-t enclave:local-signed .```