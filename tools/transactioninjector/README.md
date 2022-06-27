# Transaction injector

A tool that injects transactions (deposits from the L1 to the L2, transfers on the L2, and withdrawals back to the L1) 
into a running Obscuro network.

## Usage

Inject transactions by running `transactioninjector/main/main()` with the following arguments (runs until user 
interrupt):

  `--l1NodeHost=<x> --l1NodePort=<x> --privateKey=<x> --chainID=<x> --managementContractAddress=<x> --erc20ContractAddress=<x> --obscuroClientAddress=<x>`
