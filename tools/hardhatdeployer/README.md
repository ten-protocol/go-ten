# Contract deployer

Deploys contracts to both the TEN network (L2) and eth network (L1)

## Usage

All commands are executed by running `contractdeployer/main/main()`.

Contract is a string value to select the contract bytecode to deploy, currently ERC20, and MGMT are supported:
-  `ERC20`: a standard ERC20 contract
-  `MGMT`: the TEN L1 management contract

* Example arguments to deploy an L2 contract:

  `--nodeHost=<x> --nodePort=<x> --privateKey=<x> --contract=ERC20`

* Example arguments to deploy an L1 contract:

  `--nodeHost=<x> --nodePort=<x> --privateKey=<x> --contract=MGMT`
