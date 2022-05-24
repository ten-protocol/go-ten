# Go Obscuro
This is the reference implementation of the Obscuro Protocol.

## Structure
```
|
| - contracts: This folder contains the source code for the solidity Management contract, which will be deployed on Ethereum.  
|               For now it's a very basic version that just stores rollups.
| - go: The obscuro protocol code.
|.   |
|.   | - obscuronode: source code for the obscuro node. Note that the node is composed of two executables: "enclave" and "host".
|.   |      | 
|.   |      | - enclave: This is the component that is loaded up inside SGX.
|.   |      | - host: The component that performs the networking duties.
|.   |      | - obscuroclient: basic RPC library client code.
|.   |      | - nodecommon:
|
| 
| - integration: source code for end to end integration tests. 
|

```

## Usage
Todo

## High level overview
Todo


