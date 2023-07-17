package main

import "github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/container"

func main() {

	cliConfig := parseCLIArgs()
	obsScanContainer, err := container.NewObscuroScanContainer(cliConfig)
	if err != nil {
		panic(err)
	}

	err = obsScanContainer.Start()
	if err != nil {
		panic(err)
	}

	select {}
}
