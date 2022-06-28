# Network manager

A tool that performs various functions for the management of an Obscuro network.

* Deploys the management contract to the L1, and returns its deployed address
* Deploys the ERC20 contract to the L1, and returns its deployed address
* Injects transactions (deposits from the L1 to the L2, transfers on the L2, and withdrawals back to the L1)

## Usage

All commands are executed by running `networkmanager/main/main()`.

* Arguments to deploy a management contract:

  `--l1NodeHost=<x> --l1NodePort=<x> --privateKey=<x> --chainID=<x> deployMgmtContract`

* Arguments to deploy an ERC20 contract:

  `--l1NodeHost=<x> --l1NodePort=<x> --privateKey=<x> --chainID=<x> deployERC20Contract`

* Arguments to inject transactions (runs until user interrupt):

  `--l1NodeHost=<x> --l1NodePort=<x> --managementContractAddress=<x> --erc20ContractAddress=<x> --obscuroClientAddress=<x> injectTransactions`
