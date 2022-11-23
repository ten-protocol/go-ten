# Rollup Encryption Design

## Scope

The rollup is a public available object that's written as a bytecode in the Layer 1 chain.
Obscuro blockchain is a chain of rollups that represent the state changes.
Given the private nature of Obscuro, the rollups must be confidential.
The confidentiality model is described in this section.


## Confidentiality Model

The rollups are composed of public and private segments.
Public segments are metadata related information, not covered in this document.
Private segments are the state changes, namely the transactions.
These transactions are encrypted.


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

Use AES-256-GCM to encrypt and decrypt the transaction payload.
