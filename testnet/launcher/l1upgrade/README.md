# Upgrading contracts locally

You will need to have a local tesntet running which can be done using:

```bash
./testnet/testnet-local-build_images.sh                                         
go run ./testnet/launcher/cmd
```

Once that is done you can run the l1 upgrade using the following:

```bash
go run ./testnet/launcher/l1upgrade/cmd \
    -l1_http_url="http://eth2network:8025" \
    -private_key="f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb" \
    -network_config_addr="0x..." \
    -docker_image="testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
```

The `network-config-addr` can be found in the following endpoint `http://localhost:3000/v1/network-config/`


