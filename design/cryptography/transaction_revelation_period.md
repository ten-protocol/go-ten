# Transaction Revelation Period

## Scope

The data revelation feature ensures that Obscuro transactions are revealed after a certain period of time.
The transaction must be backwards compatible with EVM transactions and at the same time they must be able to indicate the revelation period.

For simplicity of design 1 block is the clock tick, and 1 block equals 12 seconds in real world time.


## Requirements
* Transactions must be backwards compatible with the EVM
  * Transactions must follow the same signature process
  * Transactions must be seamlessly created by existing web3 frameworks
* Transactions must specify the revelation period
* There are 5 revelation periods
  * XS - 12 seconds (1 block)
  * S - 1 hour (3600/12 = 300 blocks)
  * M - 1 day (24 * 3600/12 = 7200 blocks)
  * L - 1 month (30 * 24 * 3600/12 = 216000 blocks)
  * XL - 1 year (12 * 30 * 24 * 3600/12 = 2592000 blocks)

## Design

### Transactions must be backwards compatible with the EVM

Transactions are encrypted per default. 
Using any type of transaction will encrypt the transaction for a default period of 1 block (XS).
This will guarantee that the use of regular EVM compatible frameworks is assured.

### Transactions must specify the revelation period

Since [EIP-2718](https://eips.ethereum.org/EIPS/eip-2718) ethereum supports custom type transactions.
This EIP allows to specify custom transactions to be supported by the ethereum protocol.

Obscuro heavily relies on Geth to handle the ethereum protocol and EVM execution. 
Geth does not implement a generic custom type transaction. It does implement 3 types of transactions. 
`Legacy`, `AccessList` and `DynamicFee` transactions.

Both `AccessList` and `DynamicFee` transactions have an `AccessList` field that used by the transaction issuer to specify future storage/address accesses that the transaction will make.

AccessLists are available in web3 frameworks such as [ethers](https://docs.ethers.io/v5/api/providers/types/#types--access-lists).


Obscuro uses a single Address and one of multiple Storage Keys to specify the Revelation period for the transaction.


#### In depth AccessList usage 

The specification of the `AccessList` Field is as follows:

```
For the transaction to be valid, accessList must be of type [[{20 bytes}, [{32 bytes}...]]...], 
where ... means “zero or more of the thing to the left”.
For example, the following is a valid access list (all hex strings would in reality be in byte representation):

[
    [
        "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
        [
            "0x0000000000000000000000000000000000000000000000000000000000000003",
            "0x0000000000000000000000000000000000000000000000000000000000000007"
        ]
    ],
    [
        "0xbb9bc244d798123fde783fcc1c72d3bb8c189413",
        []
    ]
]
```

Obscuro uses the `AccessList` object to determine the revelation period.
A special address and special storage keys combination is checked for internally.

To simplify the protocol the Null address (0x0000000000000000000000000000000000000000) will be used and the hexadecimal representation of 1 to 5 is used to specify the revelation period. 

An example would be:
```
[
    [
        "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",                             <- Canonical AccessList usage
        [
            "0x0000000000000000000000000000000000000000000000000000000000000003",
            "0x0000000000000000000000000000000000000000000000000000000000000007"
        ]
    ],
    [
        "0x0000000000000000000000000000000000000000",                             <- Obscuro Revelation Period
         [
            "0x0000000000000000000000000000000000000000000000000000000000000003", <- 1 day revelation Period
            "0x0000000000000000000000000000000000000000000000000000000000000007"  <- ignored
        ]
    ]
]
```

