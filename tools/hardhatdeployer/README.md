# Contract deployer

Deploys contracts to both the TEN network (L2) and eth network (L1)

## Usage

All commands are executed by running `contractdeployer/main/main()`.

Contract is a string value to select the contract bytecode to deploy, currently ERC20, ENCLAVE_REGISTRY, ROLLUP, CROSS_CHAIN
and NETWORK_CONFIG are supported:
-  `ERC20`: a standard ERC20 contract
-  `ENCLAVE_REGISTRY`: the TEN L1 enclave registry contract
-  `CROSS_CHAIN`: the TEN L1 cross chain contract
-  `ROLLUP`: the TEN L1 rollup contract
-  `NETWORK_CONFIG`: the TEN L1 network config contract

* Example arguments to deploy an L2 contract:

  `--nodeHost=<x> --nodePort=<x> --privateKey=<x> --contract=ERC20`

* Example arguments to deploy an L1 contract:

  `--nodeHost=<x> --nodePort=<x> --privateKey=<x> --contract=ENCLAVE_REGISTRY`
