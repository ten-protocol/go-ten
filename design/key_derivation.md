# Key Derivation design

## Scope

Obscuro blockchain is seeded with a network Key.
This network key is shared between enclaves upon successful attestation.

Transactions are encrypted in the rollups that make the Obscuro blockchain.
Each rollup uses a key to encrypt the Transactions.
This key is derived from the network Key.

## Requirements

* Each rollup is encrypted with a unique key
  * Each unique key is derived from the network key
  * Each unique key is not able to disclose other unique keys
  * Each unique key is not able to disclose the network key
  * Each unique key is deterministic give the rollup height
* Each enclave is able to determine the unique key given the enclave secret and rollup height

## Design

Use HKDF standardized in https://www.rfc-editor.org/rfc/rfc5869.

HKDF has two primary and potentially independent uses:

To "extract" (condense/blend) entropy from a larger random source to provide a more uniformly unbiased and higher entropy but smaller output (e.g., an encryption key). This is done by utilising the diffusion properties of cryptographic MACs.

To "expand" the generated output of an already reasonably random input such as an existing shared key into a larger cryptographically independent output, thereby producing multiple keys deterministically from that initial shared key, so that the same process may produce those same secret keys safely on multiple devices, as long as the same inputs are utilised.


## Known limitations

TBD

## Alternatives considered

### PBKDF2

The downside of this approach is that PBKDF2 is a Password Based key derivation function, and optimized for password based keys.
To focus on Password Based keys the PBKDF2 uses a large amount of iterations to compensate the weak nature of password based keys.

