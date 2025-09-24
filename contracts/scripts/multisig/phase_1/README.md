# Executing Upgrade TX in Gnosis Multisig

1. Go to https://app.safe.global
2. Connect your wallet and select your Safe
3. Click "New Transaction" and then "Transaction Builder"
4. For each contract to upgrade:
    - Enter the Proxy Admin Address, ABI will be filled and Contract Method Selector will appear
    - Select "upgrade" function
    - Enter specific contract proxy address
    - Enter new implementation address
    - Click "Add new transaction"
5. Repeat for all contracts and click "Create Batch"
6. Click "Send Batch"
7. Click "Continue"
8. Click "Sign" and sign with your wallet
9. Wait for all multisig members to sign
10. Execute the transaction

## Walkthrough

### 1. Grant Multisig Proxy Admin Ownership of Contracts

This step only needs to be done once. There are three env vars that need to be checked before running this:
- `NETWORK_CONFIG_ADDR`
- `MULTISIG_ADDR`
- `PROXY_ADMIN_ADDR`

The contract addresses can be found in the output of the k8s prepare workflow and downloading the 
[deploy-l1-artifacts](https://github.com/ten-protocol/go-ten/actions/runs/17130412785).

1. Run the [Multisig Setup](https://github.com/ten-protocol/go-ten/actions/workflows/manual-l1-contracts-multisig-setup.yml) workflow
2. Download the `multisig-setup` artifact after it has completed


