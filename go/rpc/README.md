This package contains library code to allow client applications to connect to TEN nodes via RPC.

### Viewing keys

Viewing keys are generated inside the wallet extension (or other users of the TEN rpc client), and then signed by the wallet (e.g. MetaMask)
to which the keys relate.
The keys are then are sent to the enclave via RPC and processed by:
- checking the validity of the signature over the viewing key
- finding the account to which this viewing key corresponds

We can do that by retrieving the signing public key from the signature.
By hashing the public key, we can then determine the address of the account.
- finally the enclave will save the viewing key (which is a public key) against the account, and use it to encrypt any
sensitive requests (e.g. "eth_call" and "eth_getBalance") permitted to be viewed by that account

Client requests to the enclave are encrypted by the client with the enclave's public key and the response will be encrypted
with the relevant viewing key (pre-added using the method above) so only the intended recipient can read it.