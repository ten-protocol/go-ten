package main

import (
	"github.com/obscuronet/go-obscuro/tools/azuredeployer"
)

const (
	templateFile   = "tools/azuredeployer/enclavedeployer/vm-template.json"
	parametersFile = "tools/azuredeployer/enclavedeployer/vm-params.json"
	setupScript    = `
		sudo apt-get update
		sudo apt-get install -y docker.io
		sudo systemctl enable --now docker
		if ! [ -d "go-obscuro" ]; then git clone https://github.com/obscuronet/go-obscuro; else :; fi
		sudo docker build -t enclave - < go-obscuro/dockerfiles/enclave.Dockerfile`
)

func main() {
	azuredeployer.DeployToAzure(templateFile, parametersFile, setupScript)
}
