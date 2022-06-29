# Network manager

A tool that performs various functions for the management of an Obscuro network:

* **Contract deployer**: Deploys the management contract and ERC20 contracts to the L1, and returns their deployed addresses
* **Transaction injector**: Injects transactions across the L1 and L2 networks (deposits from the L1 to the L2, 
  transfers on the L2, and withdrawals back to the L1), then reports on whether the injected transactions were 
  successfully incorporated into the blockchain

## Usage

All commands are executed by running `networkmanager/main/main()`.

* Arguments to deploy a management contract:

  `--l1NodeHost=<x> --l1NodePort=<x> --privateKey=<x> --chainID=<x> deployMgmtContract`

* Arguments to deploy an ERC20 contract:

  `--l1NodeHost=<x> --l1NodePort=<x> --privateKey=<x> --chainID=<x> deployERC20Contract`

* Arguments to inject transactions:

  `--l1NodeHost=<x> --l1NodePort=<x> --managementContractAddress=<x> --erc20ContractAddress=<x> --obscuroClientAddress=<x> injectTransactions <num of transactions, or 0 for unlimited>`

  
