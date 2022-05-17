package main

import (
	"github.com/obscuronet/obscuro-playground/tools/azuredeployer"
)

const (
	templateFile   = "tools/azuredeployer/networkdeployer/vm-template.json"
	parametersFile = "tools/azuredeployer/networkdeployer/vm-params.json"
	setupScript    = `
		sudo apt-get update
		sudo apt-get install -y docker.io
		sudo apt-get install -y make
		sudo apt-get install -y build-essential
		wget -c https://go.dev/dl/go1.18.2.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
	
		if ! [ -d "obscuro-playground" ]; then git clone https://github.com/obscuronet/obscuro-playground; else :; fi
		sudo systemctl enable --now docker
		sudo docker build -t obscuro_enclave - < obscuro-playground/dockerfiles/enclave.Dockerfile
	
		cd obscuro-playground
		export PATH=$PATH:/usr/local/go/bin
		go mod tidy
		go build -o ./go/obscuronode/host/main/host ./go/obscuronode/host/main/main.go
		go build -o ./go/obscuronode/enclave/main/enclave ./go/obscuronode/enclave/main/main.go
		integration/gethnetwork/build_geth_binary.sh --version=v1.10.17
	`

	// todo - joel - provide script to start all the components, incl. Geth (start with one of each component)
)

func main() {
	azuredeployer.DeployToAzure(templateFile, parametersFile, setupScript)
}
