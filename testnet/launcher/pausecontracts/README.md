# Pause/Unpause All Contracts Launcher

This launcher pauses or unpauses all contracts that implement `PausableWithRoles` and then verifies they are successfully in the desired state.

## Prerequisites

You need to have a local testnet running which can be done using:

```bash
./testnet/testnet-local-build_images.sh                                         
go run ./testnet/launcher/cmd
```

## Usage

Once the testnet is running, you can run the pause/unpause all contracts script using:

```bash
go run ./testnet/launcher/pausecontracts/cmd \
    -l1_http_url="http://eth2network:8025" \
    -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
    -network_config_addr="0x..." \
    -merkle_message_bus_addr="0x..." \
    -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest" \
    -action="PAUSE"
```

## Parameters

- `l1_http_url`: Layer 1 network HTTP RPC address (default: http://eth2network:8025)
- `private_key`: L1 private key used for deployment (default: f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb)
- `network_config_addr`: L1 network config contract address (required)
- `merkle_message_bus_addr`: L1 merkle message bus contract address (required)
- `docker_image`: Docker image to run (default: testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest)
- `action`: Action to perform - either "PAUSE" or "UNPAUSE" (default: PAUSE)

The script expects these environment variables to be set in the Docker container:

- `NETWORK_JSON`: Network configuration for hardhat
- `NETWORK_CONFIG_ADDR`: Network config contract address
- `MERKLE_MESSAGE_BUS_ADDR`: Merkle message bus contract address
- `ACTION`: Action to perform (PAUSE or UNPAUSE)

## What it does

1. Connects to all contracts that implement `PausableWithRoles`
2. Pauses or unpauses each contract using the deployer account (which should have `PAUSER_ROLE`)
3. Verifies that each contract is successfully in the desired state
4. Reports the status of each contract

## Contracts affected

- MessageBus
- NetworkEnclaveRegistry
- DataAvailabilityRegistry
- TenBridge
- CrossChainMessenger
- MerkleTreeMessageBus
