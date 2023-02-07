package integration

// Tracks the start ports handed out to different tests, in a bid to minimise conflicts.
const (
	StartPortEth2NetworkTests        = 31000
	StartPortNodeRunnerTest          = 32000
	StartPortSimulationGethInMem     = 34000
	StartPortSimulationInMem         = 35000
	StartPortSimulationFullNetwork   = 37000
	StartPortSmartContractTests      = 38000
	StartPortContractDeployerTest    = 39000
	StartPortWalletExtensionUnitTest = 40000

	DefaultGethWSPortOffset      = 100
	DefaultGethAUTHPortOffset    = 200
	DefaultGethNetworkPortOffset = 300
	DefaultPrysmHTTPPortOffset   = 400
	DefaultPrysmP2PPortOffset    = 500
	DefaultHostP2pOffset         = 600 // The default offset for the host P2p
	DefaultEnclaveOffset         = 700 // The default offset between a Geth nodes port and the enclave ports. Used in Socket Simulations.
	DefaultHostRPCHTTPOffset     = 800 // The default offset for the host's RPC HTTP port
	DefaultHostRPCWSOffset       = 900 // The default offset for the host's RPC websocket port
)

const (
	EthereumChainID = 1337
	ObscuroChainID  = 777
)
