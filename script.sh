#!/usr/bin/env bash

go run /Users/WillHester/Desktop/Dev/_ten/go-ten/go/node/cmd  \
               -is_genesis=false \
               -node_type=validator \
               -is_sgx_enabled=false \
               -host_id="0x8c1B575155d61f39aBdf4DBb766E2c67384a34E1" \
               -l1_ws_url="ws://eth2network:9000" \
               -management_contract_addr="0x51D43a3Ca257584E770B6188232b199E76B022A2" \
               -message_bus_contract_addr="0xDaBD89EEA0f08B602Ec509c3C608Cb8ED095249C" \
               -l1_start="0xd1113e62dcb78d5c7c71ea98198bb589fca04d51099ec032c19e1cc899b32a51" \
               -private_key="ebca545772d6438bbbe1a16afbed455733eccf96157b52384f1722ea65ccfa89" \
               -sequencer_id="0x0654D8B60033144D567f25bF41baC1FB0D60F23B" \
               -host_public_p2p_addr="127.0.0.1:10000" \
               -host_p2p_port=10000 \
               -enclave_docker_image="testnetobscuronet.azurecr.io/obscuronet/uat_enclave:latest" \
               -host_docker_image="testnetobscuronet.azurecr.io/obscuronet/uat_host:latest" \
               -is_debug_namespace_enabled=true \
               -log_level=3 \
               -batch_interval=1s \
               -max_batch_interval=1m \
               -rollup_interval=10s \
               -l1_chain_id=1337 \
               -host_http_port=90 \
               -host_ws_port=91 \
               -node_name=postgres-node \
		           -postgres_db_host="postgres://tenuser:PB!@Z!jx564eFZb@postgres-ten-dev-testnet.postgres.database.azure.com:5432/postrgres?sslmode=disable" \
               start