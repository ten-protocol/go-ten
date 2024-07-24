package integration

// Tracks the start ports handed out to different tests, in a bid to minimise conflicts.
// Note: the max should not exceed 30000 because the OS can use those ports and we'll get conflicts
const (
	StartPortEth2NetworkTests        = 10000
	StartPortTenscanUnitTest         = 11000
	StartPortNodeRunnerTest          = 12000
	StartPortSimulationGethInMem     = 14000
	StartPortSimulationInMem         = 15000
	StartPortSimulationFullNetwork   = 16000
	StartPortNetworkTests            = 17000
	StartPortSmartContractTests      = 18000
	StartPortContractDeployerTest1   = 19000
	StartPortContractDeployerTest2   = 21000
	StartPortFaucetUnitTest          = 22000
	StartPortFaucetHTTPUnitTest      = 23000
	StartPortTenGatewayUnitTest      = 24000
	StartPortWalletExtensionUnitTest = 25000

	DefaultGethWSPortOffset         = 100
	DefaultGethAUTHPortOffset       = 200
	DefaultGethNetworkPortOffset    = 300
	DefaultGethHTTPPortOffset       = 400
	DefaultPrysmP2PPortOffset       = 500
	DefaultPrysmRPCPortOffset       = 550
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

const (
	GethNodeAddress = "0x123463a4b065722e99115d6c222f267d9cabb524"
	GethNodePK      = "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622"
)
