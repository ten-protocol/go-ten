package integration

// This file tracks the start ports handed out to different tests, in a bid to minimise conflicts.

var (
	StartPortGethNetworkTest        uint64 = 30000
	StartPortWalletExtensionTest    uint64 = 31000
	StartPortNodeRunnerTest         uint64 = 32000
	StartPortSimulationDocker              = 33000
	StartPortSimulationGethInMem           = 34000
	StartPortSimulationInMem               = 35000
	StartPortSimulationAzureEnclave        = 36000
	StartPortSimulationSocket              = 37000
)
