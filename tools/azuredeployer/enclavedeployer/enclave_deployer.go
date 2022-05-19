package main

import (
	"github.com/obscuronet/obscuro-playground/tools/azuredeployer"
)

const (
	templateFile   = "tools/azuredeployer/enclavedeployer/vm-template.json"
	parametersFile = "tools/azuredeployer/enclavedeployer/vm-params.json"
	setupScript    = `
		sudo apt-get update
		sudo apt-get install -y docker.io
		sudo systemctl enable --now docker
		if ! [ -d "obscuro-playground" ]; then git clone https://github.com/obscuronet/obscuro-playground; else :; fi
		sudo docker build -t obscuro_enclave - < obscuro-playground/dockerfiles/enclave.Dockerfile`
)

func main() {
	azuredeployer.DeployToAzure(templateFile, parametersFile, setupScript)
}
