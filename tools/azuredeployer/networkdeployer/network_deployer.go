package main

import (
	"github.com/obscuronet/obscuro-playground/tools/azuredeployer"
)

const (
	templateFile   = "tools/azuredeployer/networkdeployer/vm-template.json"
	parametersFile = "tools/azuredeployer/networkdeployer/vm-params.json"
	// Creates a network of two Geth nodes, two Obscuro hosts, and two Obscuro enclaves.
	// TODO - Start Obscuro enclave services in Docker. Requires some changes to the Docker file.
	// TODO - Re-enable blockchain validation for enclaves.
	// TODO - Detach at end of script to show final message.
	setupScript = `
		sudo apt-get update
		sudo apt-get install -y docker.io
		sudo apt-get install -y make
		sudo apt-get install -y build-essential
		wget -c https://go.dev/dl/go1.18.2.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
		sudo systemctl enable --now docker
		export PATH=$PATH:/usr/local/go/bin
	
		if ! [ -d "obscuro-playground" ]; then git clone https://github.com/obscuronet/obscuro-playground; else :; fi
		sudo docker build -t obscuro_enclave - < obscuro-playground/dockerfiles/enclave.Dockerfile

		# While this step is performed automatically as part of running the Geth network, it takes too long and the 
		# Obscuro hosts will time out when starting up.
		cd obscuro-playground/integration/gethnetwork
		source build_geth_binary.sh --version=v1.10.17

		cd ../../..
		go mod tidy
		go build -o ./go/obscuronode/host/main/host ./go/obscuronode/host/main/main.go
		go build -o ./go/obscuronode/enclave/main/enclave ./go/obscuronode/enclave/main/main.go
		go build -o ./integration/gethnetwork/main/geth ./integration/gethnetwork/main/*.go

		integration/gethnetwork/main/geth --numNodes=2 --prefundedAddrs=0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf,0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF > /dev/null &
		go/obscuronode/enclave/main/enclave --nodeID=1 --address=localhost:11000 > /dev/null &
		go/obscuronode/enclave/main/enclave --nodeID=2 --address=localhost:11001 > /dev/null &
		go/obscuronode/host/main/host --nodeID=1 --isGenesis=true --enclaveAddress=localhost:11000 --clientServerAddress=localhost:13000 --ethClientPort=12100 --privateKey=0000000000000000000000000000000000000000000000000000000000000001 > /dev/null &
		go/obscuronode/host/main/host --nodeID=2 --isGenesis=false --enclaveAddress=localhost:11001 --clientServerAddress=localhost:13001 --ethClientPort=12101 --privateKey=0000000000000000000000000000000000000000000000000000000000000002 > /dev/null &
		`
)

func main() {
	azuredeployer.DeployToAzure(templateFile, parametersFile, setupScript)
}
