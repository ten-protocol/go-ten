package integration

// Tracks the start ports handed out to different tests, in a bid to minimise conflicts.
const (
	StartPortEth2NetworkTests        = 31000
	StartPortNodeRunnerTest          = 33000
	StartPortSimulationGethInMem     = 35000
	StartPortSimulationInMem         = 37000
	StartPortSimulationFullNetwork   = 39000
	StartPortSmartContractTests      = 41000
	StartPortContractDeployerTest    = 43000
	StartPortWalletExtensionUnitTest = 45000
	StartPortFaucetUnitTest          = 47000
	StartPortFaucetHTTPUnitTest      = 49000
	StartPortTenscanUnitTest         = 51000
	StartPortTenGatewayUnitTest      = 53000

	DefaultGethWSPortOffset         = 100
	DefaultGethAUTHPortOffset       = 200
	DefaultGethNetworkPortOffset    = 300
	DefaultPrysmHTTPPortOffset      = 400
	DefaultPrysmP2PPortOffset       = 500
	DefaultHostP2pOffset            = 600 // The default offset for the host P2p
	DefaultEnclaveOffset            = 700 // The default offset between a Geth nodes port and the enclave ports. Used in Socket Simulations.
	DefaultHostRPCHTTPOffset        = 800 // The default offset for the host's RPC HTTP port
	DefaultHostRPCWSOffset          = 900 // The default offset for the host's RPC websocket port
	DefaultTenscanHTTPPortOffset    = 1000
	DefaultTenGatewayHTTPPortOffset = 1100
	DefaultTenGatewayWSPortOffset   = 1200
)

const (
	EthereumChainID = 1337
	TenChainID      = 443
)
