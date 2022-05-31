# Network manager

A tool that performs various functions for the management of an Obscuro network.

* Deploys the management contract to the L1, and returns its deployed address
* Deploys the ERC20 contract to the L1, and returns its deployed address

## Usage

All commands are executed by running `networkmanager/main/main()`.

* Arguments to deploy a management contract:

  `--l1NodeHost=<L1 node host> --l1NodePort=<L1 node port> --privateKey=<l1 node private key> --chainID=<L1 chain ID> deployMgmtContract`

* Arguments to deploy an ERC20 contract:

  `--l1NodeHost=<L1 node host> --l1NodePort=<L1 node port> --privateKey=<l1 node private key> --chainID=<L1 chain ID> deployERC20Contract`
