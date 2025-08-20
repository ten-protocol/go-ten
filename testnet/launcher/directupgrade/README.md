# Direct Upgrade for L1 Contracts

This launcher performs direct upgrades of L1 contracts using the multisig that has been set up previously.

## Prerequisites

1. You need to have a local testnet running which can be done using:

```bash
./testnet/testnet-local-build_images.sh                                         
go run ./testnet/launcher/cmd
```

2. The multisig setup must have been completed first using the multisig setup launcher.

## Usage

Once the testnet is running and multisig setup is complete, you can run the direct upgrade using:

```bash
go run ./testnet/launcher/directupgrade/cmd \
    -l1_http_url="http://eth2network:8025" \
    -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
    -network_config_addr="0x..." \
    -multisig_addr="0x..." \
    -proxy_admin_addr="0x..." \
    -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
```

## Parameters

- `l1_http_url`: Layer 1 network HTTP RPC address (default: http://eth2network:8025)
- `private_key`: L1 private key used for deployment (default: f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb)
- `network_config_addr`: L1 network config contract address (required)
- `multisig_addr`: Multisig address that controls the contracts (required)
- `proxy_admin_addr`: Proxy admin contract address (required)
- `docker_image`: Docker image to run (default: testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest)

## What it does

1. Runs the `002_direct_upgrade.ts` script in a Docker container
2. Deploys new implementations for all contracts
3. Generates Safe transaction bundle for manual execution
4. Provides implementation addresses for verification
5. Bypasses all timelock delays (immediate execution possible)

## Workflow

1. **First**: Run the multisig setup launcher to transfer proxy admin ownership
2. **Then**: Run this direct upgrade launcher to deploy new implementations and generate upgrade transactions
3. **Finally**: Execute the generated transactions through your Safe multisig

## Finding the network-config-addr

The `network-config-addr` can be found in the following endpoint: `http://localhost:3000/v1/network-config/`

## Important Notes

- This script bypasses all timelock delays
- Only use during initial mainnet phase for rapid iteration
- The multisig must have control of the proxy admin before running this
- Generated transactions need to be executed manually through the Safe UI