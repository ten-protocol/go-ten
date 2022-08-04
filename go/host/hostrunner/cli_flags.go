package hostrunner

// Flag names and usages.
const (
	configName  = "config"
	configUsage = "The path to the node's config file. Overrides all other flags"

	nodeIDName  = "id"
	nodeIDUsage = "The 20 bytes of the node's address"

	isGenesisName  = "isGenesis"
	isGenesisUsage = "Whether the node is the first node to join the network"

	gossipRoundNanosName  = "gossipRoundNanos"
	gossipRoundNanosUsage = "The duration of the gossip round"

	clientRPCPortHTTPName  = "clientRPCPortHttp"
	clientRPCPortHTTPUsage = "The port on which to listen for client application RPC requests over HTTP"

	clientRPCPortWSName  = "clientRPCPortWs"
	clientRPCPortWSUsage = "The port on which to listen for client application RPC requests over websockets"

	clientRPCHostName  = "clientRPCHost"
	clientRPCHostUsage = "The host on which to handle client application RPC requests"

	clientRPCTimeoutSecsName  = "clientRPCTimeoutSecs"
	clientRPCTimeoutSecsUsage = "The timeout for client <-> host RPC communication"

	enclaveRPCAddressName  = "enclaveRPCAddress"
	enclaveRPCAddressUsage = "The address to use to connect to the Obscuro enclave service"

	enclaveRPCTimeoutSecsName  = "enclaveRPCTimeoutSecs"
	enclaveRPCTimeoutSecsUsage = "The timeout for host <-> enclave RPC communication"

	p2pBindAddressName  = "p2pBindAddress"
	p2pBindAddressUsage = "The address where the p2p server is bound to. Defaults to 0.0.0.0:10000"

	p2pPublicAddressName  = "p2pPublicAddress"
	p2pPublicAddressUsage = "The P2P address where the other servers should connect to. Defaults to 127.0.0.1:10000"

	l1NodeHostName  = "l1NodeHost"
	l1NodeHostUsage = "The host on which to connect to the Ethereum client"

	l1NodePortName  = "l1NodePort"
	l1NodePortUsage = "The port on which to connect to the Ethereum client"

	l1ConnectionTimeoutSecsName  = "l1ConnectionTimeoutSecs"
	l1ConnectionTimeoutSecsUsage = "The timeout for connecting to the Ethereum client"

	rollupContractAddrName  = "rollupContractAddress"
	rollupContractAddrUsage = "The management contract address on the L1"

	logLevelName  = "logLevel"
	logLevelUsage = "The verbosity level of logs. (Defaults to Info)"

	logPathName  = "logPath"
	logPathUsage = "The path to use for the host's log file"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the L1 node account"

	l1ChainIDName  = "l1ChainID"
	l1ChainIDUsage = "An integer representing the unique chain id of the Ethereum chain used as an L1 (default 1337)"

	obscuroChainIDName  = "obscuroChainID"
	obscuroChainIDUsage = "An integer representing the unique chain id of the Obscuro chain (default 777)"

	profilerEnabledName  = "profilerEnabled"
	profilerEnabledUsage = "Runs a profiler instance (Defaults to false)"
)
