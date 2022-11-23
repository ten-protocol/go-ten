# Key Derivation design

## Scope

The Obscuro blockchain has a master seed that is at the heart of the chain's encryption.
This master seed is shared between enclaves upon successful attestation.
Keys are derived from the master seed, that read/write to the transaction in the blockchain.

As per defined in the whitepaper, the key used per rollup will actually be multiple keys (one for each revelation period). However, that was de-scoped in this document.


## Requirements

* The same unique key is deterministically generated from the same inputs
  * Each unique key is derived from the master seed
  * Each unique key is not able to disclose other unique keys
  * Each unique key is not able to disclose the master seed
* Each enclave is able to determine the unique key given the enclave master seed and rollup height

## Design

KDF's or Key Derivation Functions are functions or schemes to derive key or output keying material (OKM) from other secret information, the input keying material (IKM). That information may be another key or, for instance a password. It is important that the secret contains enough randomness to generate keys, without an attacker to be able to perform attacks using information about the input.

There are basically two families of KDFs. PBKDF's such as PBKDF1, PBKDF2, bcrypt, scrypt, Argon2 take passwords that need to be strengthened as input keying material and then perform strengthening. KBKDF's - Key Based Key Derivation Functions - take key material that already contains enough entropy as input.


The base construction of a KDF is:

- input:

  a binary encoded secret or key;  
  other information to derive a specific key (optional);
  output size (if configurable).


- output:
 
  a derived key or derived keying material.

Furthermore, there are many parameters possible:

  - a salt;
  
  - work factor (for PBKDF's);
  
  - memory usage (for PBKDF's);
  
  - parallelism (for PBKDF's).


### Comparison

#### PBKDF Family
- Variable cost algorithm (given the params might cost more or less to compute, specially iterations)
- Typically used for deriving a cryptographic key from a password.
- Specially focused on Key-stretching or password-stretching used protecting (as best it can) the low entropy source material.

Pros:
 - Adds entropy to weak keys (Key-stretching)
 - Allows a configurable set of interations to protect from brute-force attacks

Cons:
  - If the key already has sufficient entropy then it's slow
  - Resource hungry

#### HKDF Family
- Fixed cost algorithm (no iterations)
- Typically used for deriving a cryptographic key from a strong key.

Pros:
- Fast
- Not resource intensive

Cons:
- Does not do Key-stretching


## Decision

Use HKDF standardized in https://www.rfc-editor.org/rfc/rfc5869.

Use rollup height ( or some other shared field ) in the Info component of HKDF and use a fixed size salt as per https://soatok.blog/2021/11/17/understanding-hkdf/
