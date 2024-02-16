package integration

// Tracks the start ports handed out to different tests, in a bid to minimise conflicts.
const (
	StartPortEth2NetworkTests        = 10000
	StartPortNodeRunnerTest          = 14000
	StartPortSimulationGethInMem     = 18000
	StartPortSimulationInMem         = 22000
	StartPortSimulationFullNetwork   = 26000
	StartPortSmartContractTests      = 30000
	StartPortContractDeployerTest1   = 34000
	StartPortContractDeployerTest2   = 35000
	StartPortWalletExtensionUnitTest = 38000
	StartPortFaucetUnitTest          = 42000
	StartPortFaucetHTTPUnitTest      = 48000
	StartPortTenscanUnitTest         = 52000
	StartPortTenGatewayUnitTest      = 56000

	DefaultGethWSPortOffset         = 100
	DefaultGethAUTHPortOffset       = 200
	DefaultGethNetworkPortOffset    = 300
	DefaultPrysmHTTPPortOffset      = 400
	DefaultPrysmP2PPortOffset       = 500
	DefaultHostP2pOffset            = 600 // The default offset for the host P2p
	DefaultEnclaveOffset            = 700 // The default offset between a Geth nodes port and the enclave ports. Used in Socket Simulations.
	DefaultHostRPCHTTPOffset        = 800 // The default offset for the host's RPC HTTP port
	DefaultHostRPCWSOffset          = 900 // The default offset for the host's RPC websocket port
	DefaultTenscanHTTPPortOffset    = 950
	DefaultTenGatewayHTTPPortOffset = 951
	DefaultTenGatewayWSPortOffset   = 952
)

const (
	EthereumChainID = 1337
	TenChainID      = 443
)
