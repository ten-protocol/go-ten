# Contract deployer

Deploys the management contract and/or ERC20 contract to the L1, and returns its deployed address.

## Usage

* Get the client server address for an Obscuro host on a running Obscuro network

* Run `contractdeployer/main/main()` with the following flags to deploy the management contract:

  ```tools/contractdeployer/main/contractdeployer --l1NodeHost=<L1 node host> --l1NodePort=<L1 node port> --privateKey=<l1 node private key> --chainID=<L1 chain ID> management```

  This will return the address of the management contract

* Switch `management` to `erc20` in the command above to deploy an ERC20 contract and get its address
