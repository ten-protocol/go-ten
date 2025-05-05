/usr/local/go/bin/go run ./go/node/cmd \
-node_name=new_node \
-node_type=validator \
-is_genesis=false \
-num_enclaves=1 \
-is_sgx_enabled=False \
-enclave_docker_image=testnetobscuronet.azurecr.io/obscuronet/enclave:latest \
-host_docker_image=testnetobscuronet.azurecr.io/obscuronet/host:latest \
-l1_ws_url=ws://eth2network:9000 \
-host_http_port=40592 \
-host_ws_port=16716 \
-host_p2p_port=28439 \
-host_p2p_host=http://127.0.0.1 \
-enclave_http_port=11000 \
-enclave_WS_port=11001 \
-private_key=e80da309bbfae5f18088701a7059ab5dc82288b627d5d53295ea6230e391c06b \
-sequencer_addr= \
-enclave_registry_addr=0x26c62148Cf06C9742b8506A2BCEcd7d72E51A206 \
-cross_chain_addr=0x19e98b050662b49D6AbDFBe2467016430197BA90 \
-da_registry_addr=0x5b8b9160C4C2084cd8dDA7B4E2428C231cf29E7d \
-network_config_addr=0x0a0b7fdB9B79D7c838675Aca65ec7293b6Cb0846 \
-message_bus_contract_addr=0xDaBD89EEA0f08B602Ec509c3C608Cb8ED095249C \
-bridge_contract_addr=0xeDFD5157b075f423795F5406527014A80B72C40F \
-l1_start=None \
-edgeless_db_image=ghcr.io/edgelesssys/edgelessdb-sgx-1gb:v0.3.2 \
-is_debug_namespace_enabled=False \
-log_level=3 \
-is_inbound_p2p_disabled=True \
-batch_interval=1s \
-host_id=0x71B0cd82d3444d27d2Cb0889540D194993f875d1 \
-max_batch_interval=1s \
-rollup_interval=3s \
-l1_chain_id=1337 \
-host_public_p2p_addr=new_node:28439 \
-l1_beacon_url=eth2network:126000 \
-system_contracts_upgrader=0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77 \
start




