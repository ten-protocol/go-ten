./testnet/testnet-local-build_images.sh 
go run ./testnet/launcher/cmd 

go run ./testnet/launcher/erc20deployer/cmd \
  -token_name="USD Coin" \
  -token_symbol="USDC" \
  -token_decimals="6" \
  -token_supply="1000000000" \
  -l1_http_url="http://eth2network:8025" \
  -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
  -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest" \
  -network_config_addr="0x2a8b83Fd5EB49A7a620F27f34D52DFA86Dabf393"

  go run ./testnet/launcher/bridgetokenwhitelist/cmd \
  -token_address="0xC8fA185016648022AEd9563928cfe98f5A8a5aa2" \
  -token_name="USD Coin" \
  -token_symbol="USDC" \
  -l1_http_url="http://eth2network:8025" \
  -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
  -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest" \
  -network_config_addr="0x2a8b83Fd5EB49A7a620F27f34D52DFA86Dabf393"