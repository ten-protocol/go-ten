package network

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ten-protocol/go-ten/integration/noderunner"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"github.com/ten-protocol/go-ten/go/node"
	"github.com/ten-protocol/go-ten/go/obsclient"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/eth2network"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
	"github.com/ten-protocol/go-ten/integration/simulation/stats"
)

// creates TEN nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkOfSocketNodes struct {
	l2Clients         []rpc.Client
	hostWebsocketURLs []string

	// geth
	eth2Network    eth2network.Eth2Network
	gethClients    []ethadapter.EthClient
	wallets        *params.SimWallets
	obscuroClients []*obsclient.ObsClient
}

func NewNetworkOfSocketNodes(wallets *params.SimWallets) Network {
	return &networkOfSocketNodes{
		wallets: wallets,
	}
}

func (n *networkOfSocketNodes) Create(simParams *params.SimParams, _ *stats.Stats) (*RPCHandles, error) {
	// kickoff the network with the prefunded wallet addresses
	simParams.L1SetupData, n.gethClients, n.eth2Network = SetUpGethNetwork(
		n.wallets,
		simParams.StartPort,
		simParams.NumberOfNodes,
		int(simParams.AvgBlockDuration.Seconds()),
	)

	simParams.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(&simParams.L1SetupData.MgmtContractAddress, testlog.Logger())
	simParams.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(
		&simParams.L1SetupData.MgmtContractAddress,
		&simParams.L1SetupData.ObxErc20Address,
		&simParams.L1SetupData.EthErc20Address,
	)

	// get the sequencer Address
	seqPrivateKey := n.wallets.NodeWallets[0].PrivateKey()
	seqPrivKey := fmt.Sprintf("%x", crypto.FromECDSA(seqPrivateKey))
	seqHostAddress := crypto.PubkeyToAddress(seqPrivateKey.PublicKey)

	// create the nodes
	nodes := make([]node.Node, simParams.NumberOfNodes)
	var err error
	for i := 0; i < simParams.NumberOfNodes; i++ {
		privateKey := seqPrivKey
		hostAddress := seqHostAddress
		nodeTypeStr := "sequencer"
		isInboundP2PDisabled := false

		// if it's not the sequencer
		if i != 0 {
			nodeTypeStr = "validator"
			privateKey = fmt.Sprintf("%x", crypto.FromECDSA(n.wallets.NodeWallets[i].PrivateKey()))
			hostAddress = crypto.PubkeyToAddress(n.wallets.NodeWallets[i].PrivateKey().PublicKey)

			// only the validators can have the incoming p2p disabled
			isInboundP2PDisabled = i == simParams.NodeWithInboundP2PDisabled
		}

		genesis := "{}"
		if simParams.WithPrefunding {
			genesis = ""
		}

		// create the nodes
		nodes[i] = noderunner.NewInMemNode(
			node.NewNodeConfig(
				node.WithGenesis(i == 0),
				node.WithHostID(hostAddress.String()),
				node.WithPrivateKey(privateKey),
				node.WithSequencerID(seqHostAddress.String()),
				node.WithEnclaveWSPort(simParams.StartPort+integration.DefaultEnclaveOffset+i),
				node.WithHostWSPort(simParams.StartPort+integration.DefaultHostRPCWSOffset+i),
				node.WithHostHTTPPort(simParams.StartPort+integration.DefaultHostRPCHTTPOffset+i),
				node.WithHostP2PPort(simParams.StartPort+integration.DefaultHostP2pOffset+i),
				node.WithHostPublicP2PAddr(fmt.Sprintf("127.0.0.1:%d", simParams.StartPort+integration.DefaultHostP2pOffset+i)),
				node.WithManagementContractAddress(simParams.L1SetupData.MgmtContractAddress.String()),
				node.WithMessageBusContractAddress(simParams.L1SetupData.MessageBusAddr.String()),
				node.WithNodeType(nodeTypeStr),
				node.WithCoinbase(simParams.Wallets.L2FeesWallet.Address().Hex()),
				node.WithL1WebsocketURL(fmt.Sprintf("ws://%s:%d", "127.0.0.1", simParams.StartPort+100)),
				node.WithInboundP2PDisabled(isInboundP2PDisabled),
				node.WithLogLevel(4),
				node.WithDebugNamespaceEnabled(true),
				node.WithL1BlockTime(simParams.AvgBlockDuration),
				node.WithTenGenesis(genesis),
			),
		)

		// start the nodes
		err = nodes[i].Start()
		if err != nil {
			errCheck := checkProcessPort(err.Error())
			if errCheck != nil {
				testlog.Logger().Warn("no port found on error", log.ErrKey, err)
			}
			fmt.Printf("unable to start TEN node: %s", err)
			testlog.Logger().Error("unable to start TEN node ", log.ErrKey, err)
		}
	}

	// create the l2 and eth connections
	err = n.createConnections(simParams)
	if err != nil {
		testlog.Logger().Crit("unable to create node connections", log.ErrKey, err)
	}
	walletClients := createAuthClientsPerWallet(n.l2Clients, simParams.Wallets)

	return &RPCHandles{
		EthClients:     n.gethClients,
		ObscuroClients: n.obscuroClients,
		RPCClients:     n.l2Clients,
		AuthObsClients: walletClients,
	}, nil
}

func (n *networkOfSocketNodes) TearDown() {
	// Stop the TEN nodes first (each host will attempt to shut down its enclave as part of shutdown).
	StopObscuroNodes(n.l2Clients)
	StopEth2Network(n.gethClients, n.eth2Network)
	CheckHostRPCServersStopped(n.hostWebsocketURLs)
}

func (n *networkOfSocketNodes) createConnections(simParams *params.SimParams) error {
	// create the clients in the structs
	n.l2Clients = make([]rpc.Client, simParams.NumberOfNodes)
	n.hostWebsocketURLs = make([]string, simParams.NumberOfNodes)
	n.obscuroClients = make([]*obsclient.ObsClient, simParams.NumberOfNodes)

	for i := 0; i < simParams.NumberOfNodes; i++ {
		var client rpc.Client
		var err error

		// create a connection to the newly created nodes - panic if no connection is made after some time
		startTime := time.Now()
		for connected := false; !connected; time.Sleep(500 * time.Millisecond) {
			client, err = rpc.NewNetworkClient(fmt.Sprintf("ws://127.0.0.1:%d", simParams.StartPort+integration.DefaultHostRPCWSOffset+i))
			connected = err == nil // The client cannot be created until the node has started.
			if time.Now().After(startTime.Add(2 * time.Minute)) {
				return fmt.Errorf("failed to create a connect to node after 2 minute - %w", err)
			}

			testlog.Logger().Info(fmt.Sprintf("Could not create client %d. Retrying...", i), log.ErrKey, err)
		}

		n.l2Clients[i] = client
		n.hostWebsocketURLs[i] = fmt.Sprintf("ws://%s:%d", Localhost, simParams.StartPort+integration.DefaultHostRPCWSOffset+i)
	}

	for idx, l2Client := range n.l2Clients {
		n.obscuroClients[idx] = obsclient.NewObsClient(l2Client)
	}

	// make sure the nodes are healthy
	for _, client := range n.obscuroClients {
		startTime := time.Now()
		healthy := false
		for ; !healthy; time.Sleep(500 * time.Millisecond) {
			healthy, _ = client.Health()
			if time.Now().After(startTime.Add(3 * time.Minute)) {
				return fmt.Errorf("nodes not healthy after 3 minutes")
			}
		}
	}
	return nil
}

// getProcessesUsingPort returns a slice of process details using the specified port.
func checkProcessPort(errPort string) error {
	re := regexp.MustCompile(`:(\d+):`)
	matches := re.FindStringSubmatch(errPort)

	if len(matches) < 2 {
		return fmt.Errorf("no port found in string")
	}

	port := matches[1]

	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%s", port)) //nolint:gosec

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var processes []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "LISTEN") || strings.Contains(line, "ESTABLISHED") {
			processes = append(processes, line)
		}
	}

	fmt.Printf("Found processes still opened on port %s - %+v\n", port, processes)

	return nil
}
