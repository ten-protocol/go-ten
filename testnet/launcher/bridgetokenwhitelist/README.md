# Bridge Token Whitelister

This tool whitelists existing ERC20 tokens (like USDC and USDT) on the TenBridge and registers them in the NetworkConfig contract, making them bridgeable and discoverable via the `ten_config()` RPC method.

## Overview

The Bridge Token Whitelister:
1. Whitelists an existing ERC20 token on the TenBridge contract (grants `ERC20_TOKEN_ROLE`)
2. Registers the token in the NetworkConfig contract via `addAdditionalAddress()`
3. L2 nodes automatically detect the event and sync the address within seconds

## Use Cases

- **Whitelisting USDC/USDT**: For production networks using real stablecoins
- **Verified Tokens**: For allowing specific verified ERC20s to be bridged
- **Curated Token List**: Maintaining a list of approved tokens

## Known Token Addresses

### Sepolia Testnet
- **USDC**: `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238`
- **USDT**: `0x7169D38820dfd117C3FA1f22a697dBA58d90BA06`

### Ethereum Mainnet
- **USDC**: `0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`
- **USDT**: `0xdAC17F958D2ee523a2206206994597C13D831ec7`

## Usage

### GitHub Actions

1. Go to **Actions** → **[M] Whitelist Bridge Token**
2. Click **Run workflow**
3. Fill in the parameters:
   - **Testnet Type**: `uat-testnet`, `sepolia-testnet`, or `mainnet`
   - **Token Address**: e.g., `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238` (USDC on Sepolia)
   - **Token Name**: e.g., `"USD Coin"`
   - **Token Symbol**: e.g., `"USDC"`
   - **Confirmation**: Type `"confirm"` if whitelisting on sepolia or mainnet
4. Click **Run workflow**

### CLI

```bash
go run ./testnet/launcher/bridgetokenwhitelist/cmd \
  -token_address="0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238" \
  -token_name="USD Coin" \
  -token_symbol="USDC" \
  -l1_http_url="https://sepolia.infura.io/v3/YOUR_KEY" \
  -l2_rpc_url="erpc.testnet.ten.xyz"
  -l2_nonce="0"
  -private_key="0x..." \
  -docker_image="your-docker-image" \
  -network_config_addr="0x..."
```

go run ./testnet/launcher/erc20deployer/cmd \
  -token_name="USD Coin" \
  -token_symbol="USDC" \
  -token_decimals="6" \
  -token_supply="100000" \
  -l1_http_url="http://eth2network:8025" \
  -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
  -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest" \
  -network_config_addr="0x2a8b83Fd5EB49A7a620F27f34D52DFA86Dabf393"

  go run ./testnet/launcher/bridgetokenwhitelist/cmd \
  -token_address="0xc9f23cA31AF54cFa22591e430F6557598fe8911F" \
  -token_name="USD Coin" \
  -token_symbol="USDC" \
  -l1_http_url="http://eth2network:8025" \
  -l2_rpc_url="http://sequencer-host:80" \
  -l2_nonce="70" \
  -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
  -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest" \
  -network_config_addr="0x2a8b83Fd5EB49A7a620F27f34D52DFA86Dabf393"
