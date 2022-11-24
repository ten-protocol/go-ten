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
The encryption and decryption rollup approach is described in the [Rollup Encryption](rollup_encryption.md) document.

## Key Derivation 

In order to avoid reusing the same key for all encryption, keys are deterministically derived from the Master Seed.
The derivation of new keys given a master seed is described in the [Key Derivation](key_derivation.md) document.

## Transaction Revelation Period

Transactions are able to specify the different revelation period desired, while being EVM compatible.
The transaction changes that specify the revelation period are described in the [Transaction Revelation Period](transaction_revelation_period.md) document.

## Data Revelation Mechanism

Generating different Encryption Keys and associating them with different time lengths allows to create multiple Data Revelation Periods.
The data revelation mechanism is described in the [Data Revelation](data_revelation.md) document.
