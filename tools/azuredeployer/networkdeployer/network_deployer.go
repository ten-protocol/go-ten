package main

import (
	"github.com/obscuronet/obscuro-playground/tools/azuredeployer"
)

const (
	templateFile   = "tools/azuredeployer/networkdeployer/vm-template.json"
	parametersFile = "tools/azuredeployer/networkdeployer/vm-params.json"
	// Creates a network of two Geth nodes, two Obscuro hosts, and two Obscuro enclaves.
	// TODO - Re-enable blockchain validation for enclaves.
	setupScript = `
		sudo apt-get update
		sudo apt-get install -y docker.io
		sudo apt-get install -y make
		sudo apt-get install -y build-essential
		wget -c https://go.dev/dl/go1.18.2.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
		sudo systemctl enable --now docker
		export PATH=$PATH:/usr/local/go/bin
	
		if ! [ -d "go-obscuro" ]; then git clone https://github.com/obscuronet/go-obscuro; else :; fi
		sudo docker build -t obscuro_enclave - < go-obscuro/dockerfiles/enclave.Dockerfile

		# This step happens automatically when first running a Geth network, but it is too slow and the Obscuro hosts 
		# time out when starting up.
		cd go-obscuro/integration/gethnetwork
		source build_geth_binary.sh --version=v1.10.17

		cd ../../..
		go mod tidy
		go build -o ./go/host/main/host ./go/obscuronode/host/main/main.go
		go build -o ./go/enclave/main/enclave ./go/obscuronode/enclave/main/main.go
		go build -o ./integration/gethnetwork/main/geth ./integration/gethnetwork/main/*.go
		go build -o ./tools/obscuroscan/main/obscuroscan ./tools/obscuroscan/main/*.go
		go build -o ./tools/networkmanager/main/networkmanager ./tools/networkmanager/main/*.go
		`
)

func main() {
	azuredeployer.DeployToAzure(templateFile, parametersFile, setupScript)
}
