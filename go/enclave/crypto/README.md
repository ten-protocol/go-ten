This package contains logic which implements the cryptographic requirements of TEN.

1. Manage the shared secret of the network.(SS) - shared_secret_service
2. Manage the "Ten RPC" encryption - which is the key used by all clients to communicate with the TEN network (key derived from SS) - rpc_key_service
3. Manage the Data availability(DA) (Rollup and Batches) Encryption/Decryption ( key derived from SS). - da_enc_service
4. Manage the enclave key signature/encryption/decryption/ id derivation. - enclave_key_service
5. Manage entropy per batch and tx - evm_entropy_service
