#!/usr/bin/env bash

export PATH=$PATH:/usr/local/go/bin
sudo systemctl enable --now docker

rm -f ./run_logs.txt
touch ./run_logs.txt

PRIV_KEY_ONE="0000000000000000000000000000000000000000000000000000000000000001"
PUB_KEY_ONE="0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
PRIV_KEY_TWO="0000000000000000000000000000000000000000000000000000000000000002"
PUB_KEY_TWO="0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF"

obscuro-playground/integration/gethnetwork/main/geth --numNodes=2 --startPort=12000 --websocketStartPort=12100 --prefundedAddrs=$PUB_KEY_ONE,$PUB_KEY_TWO > ./run_logs.txt 2>&1 &
sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 127.0.0.1:11000:11000/tcp obscuro_enclave --hostID 1 --address :11000 > ./run_logs.txt 2>&1 &
sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 127.0.0.1:11001:11000/tcp obscuro_enclave --hostID 2 --address :11000 > ./run_logs.txt 2>&1 &
obscuro-playground/go/obscuronode/host/main/host --id=1 --isGenesis=true --enclaveRPCAddress=localhost:11000 --clientRPCAddress=0.0.0.0:13000 --l1NodePort=12100 --privateKey=$PRIV_KEY_ONE > ./run_logs.txt 2>&1 &
obscuro-playground/go/obscuronode/host/main/host --id=2 --isGenesis=false --enclaveRPCAddress=localhost:11001 --clientRPCAddress=localhost:13001 --l1NodePort=12101 --privateKey=$PRIV_KEY_TWO > ./run_logs.txt 2>&1 &
cd obscuro-playground
sudo ./tools/obscuroscan/main/obscuroscan --clientServerAddress=127.0.0.1:13000 --address=0.0.0.0:80 > ../run_logs.txt 2>&1 &