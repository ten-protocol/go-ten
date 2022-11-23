# Cryptography Design

## Scope

Following the Cryptography section of the Obscuro blockchain whitepaper this document describes the design around cryptography in Obscuro.

The Cryptography area has been broken down into sections for interative implementation.

## Areas

* Rollup Encryption design
* Key Derivation design
* Data Revelation design

## Rollup Encryption

The rollups in Obscuro have an encrypted design.
In the [Rollup Encryption](rollup_encryption.md) document the approach to encrypt and decrypt the rollup is described.

## Key Derivation 

In order to avoid reusing the same key for all encryption, keys are deterministically derived from the Master Seed.
In the [Key Derivation](key_derivation.md) document the approach to derive new keys from the master seed is described.

## Data Revelation Mechanism

Generating different Encryption Keys and associating them with different time lengths allows to create multiple Data Revelation Periods.
In the [Data Revelation](data_revelation.md) document the approach to reveal data is described.

