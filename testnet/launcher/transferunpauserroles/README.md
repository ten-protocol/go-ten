# Transfer Unpauser Roles Launcher

This launcher transfers the `UNPAUSER_ROLE` from the deployer to a specified multisig address for contracts using `PausableWithRoles`.

## Prerequisites

You need to have a local testnet running which can be done using:

```bash
./testnet/testnet-local-build_images.sh                                         
go run ./testnet/launcher/cmd
```

## Usage

Once the testnet is running, you can run the unpauser role transfer using:

```bash
go run ./testnet/launcher/transferunpauserroles/cmd \
    -l1_http_url="http://eth2network:8025" \
    -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
    -network_config_addr="0x..." \
    -multisig_addr="0x..." \
    -merkle_message_bus_addr="0x..." \
    -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
```

## Parameters

- `l1_http_url`: Layer 1 network HTTP RPC address (default: http://eth2network:8025)
- `private_key`: L1 private key used for deployment (default: f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb)
- `network_config_addr`: L1 network config contract address (required)
- `multisig_addr`: Multisig address to transfer unpauser role to (required)
- `docker_image`: Docker image to run (default: testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest)

The script expects these environment variables to be set in the Docker container:

- `NETWORK_JSON`: Network configuration for hardhat
- `NETWORK_CONFIG_ADDR`: Network config contract address
- `MULTISIG_ADDR`: Multisig wallet address
- `MERKLE_MESSAGE_BUS_ADDR`: Merkle message bus contract address