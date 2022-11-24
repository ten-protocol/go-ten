# Cryptography Design

## Scope

Following the Cryptography section of the Obscuro blockchain whitepaper this document describes the design implementation of the topic in Obscuro.

The Cryptography topic has been broken down into sections for iterative implementation and design.

## Areas

* Rollup Encryption design
* Key Derivation design
* Transaction Revelation Period
* Data Revelation design

## Rollup Encryption

The rollups in Obscuro have an encrypted principle that provides privacy to the transactions.
The encryption and decryption rollup approach is described in the [Rollup Encryption](cryptography/rollup_encryption.md) document.

## Key Derivation 

In order to avoid reusing the same key for all encryption, keys are deterministically derived from the Master Seed.
Using this key derivation framework, keys will also be derived to encrypt data for different revelation keys.
The derivation of new keys given a master seed is described in the [Key Derivation](cryptography/key_derivation.md) document.

## Transaction Revelation Period

Transactions are able to specify the different revelation period desired, while being EVM compatible.
The transaction changes that specify the revelation period are described in the [Transaction Revelation Period](cryptography/transaction_revelation_period.md) document.

## Data Revelation Mechanism

Generating different Encryption Keys and associating them with different time lengths allows to create multiple Data Revelation Periods.
The data revelation mechanism is described in the [Data Revelation](cryptography/data_revelation.md) document.
