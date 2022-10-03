package integration

// Tracks the start ports handed out to different tests, in a bid to minimise conflicts.
const (
	StartPortGethNetworkTest         uint64 = 30000
	StartPortNodeRunnerTest          uint64 = 32000
	StartPortSimulationGethInMem            = 34000
	StartPortSimulationInMem                = 35000
	StartPortSimulationFullNetwork          = 37000
	StartPortSmartContractTests             = 38000
	StartPortContractDeployerTest           = 39000
	StartPortWalletExtensionUnitTest        = 40000
)

const (
	EthereumChainID = 1337
	ObscuroChainID  = 777
)
