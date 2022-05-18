#!/usr/bin/env bash

export PATH=$PATH:/usr/local/go/bin

obscuro-playground/integration/gethnetwork/main/geth --numNodes=2 --prefundedAddrs=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf,0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF > /dev/null &
obscuro-playground/go/obscuronode/enclave/main/enclave --nodeID=1 --address=localhost:11000 --writeToLogs=true > /dev/null &
obscuro-playground/go/obscuronode/enclave/main/enclave --nodeID=2 --address=localhost:11001 --writeToLogs=true > /dev/null &
obscuro-playground/go/obscuronode/host/main/host --nodeID=1 --isGenesis=true --enclaveAddress=localhost:11000 --clientServerAddress=0.0.0.0:13000 --ethClientPort=12100 --privateKey=0000000000000000000000000000000000000000000000000000000000000001 > /dev/null &
obscuro-playground/go/obscuronode/host/main/host --nodeID=2 --isGenesis=false --enclaveAddress=localhost:11001 --clientServerAddress=localhost:13001 --ethClientPort=12101 --privateKey=0000000000000000000000000000000000000000000000000000000000000002 > /dev/null &