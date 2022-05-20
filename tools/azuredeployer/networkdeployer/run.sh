#!/usr/bin/env bash

export PATH=$PATH:/usr/local/go/bin
sudo systemctl enable --now docker

rm ./run_logs.txt
touch ./run_logs.txt

obscuro-playground/integration/gethnetwork/main/geth --numNodes=2 --prefundedAddrs=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf,0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF &>> ./run_logs.txt &
sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 127.0.0.1:11000:11000/tcp obscuro_enclave --nodeID 1 --address :11000 &>> ./run_logs.txt &
sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 127.0.0.1:11001:11000/tcp obscuro_enclave --nodeID 2 --address :11000 &>> ./run_logs.txt &
obscuro-playground/go/obscuronode/host/main/host --nodeID=1 --isGenesis=true --enclaveAddress=localhost:11000 --clientServerAddress=0.0.0.0:13000 --ethClientPort=12100 --privateKey=0000000000000000000000000000000000000000000000000000000000000001 &>> ./run_logs.txt &
obscuro-playground/go/obscuronode/host/main/host --nodeID=2 --isGenesis=false --enclaveAddress=localhost:11001 --clientServerAddress=localhost:13001 --ethClientPort=12101 --privateKey=0000000000000000000000000000000000000000000000000000000000000002 &>> ./run_logs.txt &
cd obscuro-playground
sudo ./tools/obscuroscan/main/obscuroscan --clientServerAddress=127.0.0.1:13000 --address=0.0.0.0:80 &>> ../run_logs.txt &