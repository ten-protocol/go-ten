# Multisig Setup for L1 Contracts

This launcher sets up multisig control for L1 contracts by transferring proxy admin ownership to a specified multisig address.

## Prerequisites

You need to have a local testnet running which can be done using:

```bash
./testnet/testnet-local-build_images.sh                                         
go run ./testnet/launcher/cmd
```

## Usage

Once the testnet is running, you can run the multisig setup using:

```bash
go run ./testnet/launcher/multisigsetup/cmd \
    -l1_http_url="http://eth2network:8025" \
    -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
    -network_config_addr="0x..." \
    -multisig_address="0x..." \
    -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
```

## Parameters

- `l1_http_url`: Layer 1 network HTTP RPC address (default: http://eth2network:8025)
- `private_key`: L1 private key used for deployment (default: f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb)
- `network_config_addr`: L1 network config contract address (required)
- `multisig_address`: Multisig address to transfer proxy admin ownership to (required)
- `docker_image`: Docker image to run (default: testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest)

## What it does

1. Runs the `001_direct_multisig_setup.ts` script in a Docker container
2. Transfers proxy admin ownership for all contracts to the specified multisig address
3. Verifies that the transfer was successful
4. Provides detailed logging of the process

## Finding the network-config-addr

The `network-config-addr` can be found in the following endpoint: `http://localhost:3000/v1/network-config/`
