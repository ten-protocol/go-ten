# Key Derivation design

## Scope

The Obscuro blockchain has a network key that allows read/write to the transaction in the block chain.
This network key is shared between enclaves upon successful attestation.

Transactions are encrypted in the rollups that make the Obscuro blockchain.
Each rollup uses a key to encrypt the Transactions.
Each key is derived from the initial network Key.

As per defined in the whitepaper, the key used per rollup will actually be multiple keys (one for each revelation period). However, that was de-scoped in this document.


## Requirements

* Each rollup has a transaction payload that is encrypted with a unique key
  * Each unique key is derived from the network key
  * Each unique key is not able to disclose other unique keys
  * Each unique key is not able to disclose the network key
  * Each unique key is deterministic give the rollup height
* Each enclave is able to determine the unique key given the enclave secret and rollup height

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


## Encryption Algorithm

Each transaction payload will be encrypted with a specific cypher which will take as one of the inputs the derived key.

AES-256 is the industry standard. It's based on a design principle known as a Substitution permutation network. It is fast in both software and hardware.

### AES Cypher mode

There are many cypher mode that AES can operate on.

#### Modes that require padding (block encryption modes: ECB, CBC)
Padding can generally be dangerous because it opens up the possibility of padding oracle attacks

#### Stream cipher modes (CTR, OFB, CFB's)
These modes generate a pseudo random stream of data that may or may not depend on the plaintext. 
Similarly to stream ciphers generally, the generated pseudo random stream is XORed with the plaintext to generate the ciphertext. 
As you can use as many bits of the random stream as you like you don't need padding at all.


#### Authenticated encryption (CCM, OCB, GCM)
To prevent padding oracle attacks and changes to the ciphertext, one can compute a message authentication code (MAC) on the ciphertext and only decrypt it if it has not been tampered with. This is called encrypt-then-mac.


#### Conclusion

While there is *a lot* of literature around the different types of cypher mode and its implementation, the clear winner is GCM.

AES-GCM provides the following features:
- Stream cypher ( no padding attacks )
- Authentication based (if Message Authentication Code is tampered then decryption the remaining ciphertext is avoided)
- Faster than CCM and free to use when compared with OCB
- Intel has implemented special instructions that speed up the underlying GHASH calculation


### Other symmetric encryption algorithms

Other symmetric algorithms exist like XChaCha20, but they lack the battle tested scrutiny that AES-256 as been submitted to.

Nonetheless, the implementation of the encryption algorithm should be interfaced in such a way that swapping algorithms is easy.
This will allow testing different algorithms behaviour and performance in-loco.



## Decision

Use HKDF standardized in https://www.rfc-editor.org/rfc/rfc5869.

Use rollup height ( or some other shared field ) in the Info component of HKDF and use a fixed size salt as per https://soatok.blog/2021/11/17/understanding-hkdf/

Use AES-256-GCM to encrypt and decrypt the transaction payload with the derived key.

