# Emergency Pause/Unpause for L1 Contracts

This launcher provides emergency pause and unpause functionality for all L1 contracts that implement the pause mechanism. This is designed for emergency situations where immediate action is required to stop malicious code execution.

## Prerequisites

1. You need to have a local testnet running which can be done using:

```bash
./testnet/testnet-local-build_images.sh                                         
go run ./testnet/launcher/cmd
```

2. The contracts must have been deployed and the multisig setup must have been completed.

## Usage

Once the testnet is running, you can run the emergency pause/unpause using:

```bash
go run ./testnet/launcher/emergencypause/cmd \
    -l1_http_url="http://eth2network:8025" \
    -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
    -network_config_addr="0x..." \
    -action="pause" \
    -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
```

## Parameters

- `l1_http_url`: Layer 1 network HTTP RPC address (default: http://eth2network:8025)
- `private_key`: L1 private key used for deployment (default: f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb)
- `network_config_addr`: L1 network config contract address (required)
- `action`: Action to perform - either "pause" or "unpause" (required)
- `docker_image`: Docker image to run (default: testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest)

## What it does

1. Runs the `001_emergency_pause.ts` script in a Docker container
2. Automatically detects all contracts that have pause/unpause functionality
3. Generates transaction data for manual execution through multisig wallet
4. Provides detailed logging of the process
5. Bypasses all timelock delays (immediate execution possible)

## Supported Contracts

The following contracts support pause/unpause functionality:
- `CrossChain`
- `NetworkEnclaveRegistry`
- `DataAvailabilityRegistry`
- `TenBridge`
- `EthereumBridge`
- `MessageBus`
- `MerkleTreeMessageBus`

## Workflow

1. **Emergency Pause**: Use `-action="pause"` to generate pause transaction data for all contracts
2. **Emergency Unpause**: Use `-action="unpause"` to generate unpause transaction data for all contracts
3. **Manual Execution**: Execute the generated transactions through your multisig wallet

## Finding the network-config-addr

The `network-config-addr` can be found in the following endpoint: `http://localhost:3000/v1/network-config/`

## Important Notes

- This script generates transaction data for manual execution
- Only use during genuine emergency situations
- Transactions must be executed manually through your multisig wallet
- This is designed for rapid response to security incidents
- No direct contract interaction - only transaction data generation
