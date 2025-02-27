package network

import (
	"bufio"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	testcommon "github.com/ten-protocol/go-ten/integration/common"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/host/l1"
	"github.com/ten-protocol/go-ten/integration/noderunner"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"
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
	eth2Network eth2network.PosEth2Network
	gethClients []ethadapter.EthClient
	wallets     *params.SimWallets
	tenClients  []*obsclient.ObsClient
}

func (n networkOfSocketNodes) Create(params *params.SimParams, stats *stats.Stats) (*RPCHandles, error) {
	//TODO implement me
	panic("implement me")
}

func (n networkOfSocketNodes) TearDown() {
	//TODO implement me
	panic("implement me")
}

func NewNetworkOfSocketNodes(wallets *params.SimWallets) Network {
	return &networkOfSocketNodes{
		wallets: wallets,
	}
}

func (n *networkOfSocketNodes) Create(simParams *params.SimParams, _ *stats.Stats) (*RPCHandles, error) {
	// kickoff the network with the prefunded wallet addresses
	simParams.L1TenData, n.gethClients, n.eth2Network = SetUpGethNetwork(
		n.wallets,
		simParams.StartPort,
		simParams.NumberOfNodes,
	)
	networkContractConfig, err := contractlib.NewNetworkConfigLib(simParams.L1TenData.NetworkConfigAddress, *n.gethClients[0].EthClient())
	if err != nil {
		return nil, err
	}
	println(networkContractConfig.GetContractAddr().Hex())
	simParams.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(
		&simParams.L1TenData.NetworkConfigAddress,
		&simParams.L1TenData.ObxErc20Address,
		&simParams.L1TenData.EthErc20Address,
	)
	beaconURL := fmt.Sprintf("127.0.0.1:%d", simParams.L1BeaconPort)
	simParams.BlobResolver = l1.NewBlobResolver(ethadapter.NewL1BeaconClient(
		ethadapter.NewBeaconHTTPClient(new(http.Client), beaconURL)))

	// get the sequencer Address
	seqPrivateKey := n.wallets.NodeWallets[0].PrivateKey()
	seqPrivKey := fmt.Sprintf("%x", crypto.FromECDSA(seqPrivateKey))
	seqHostAddress := "0"

	// create the nodes
	nodes := make([]node.Node, simParams.NumberOfNodes)
	for i := 0; i < simParams.NumberOfNodes; i++ {
		privateKey := seqPrivKey
		hostAddress := seqHostAddress
		nodeTypeStr := "sequencer"
		isInboundP2PDisabled := false

		// if it's not the sequencer
		if i != 0 {
			nodeTypeStr = "validator"
			privateKey = fmt.Sprintf("%x", crypto.FromECDSA(n.wallets.NodeWallets[i].PrivateKey()))
			hostAddress = fmt.Sprintf("%d", i)
			// only the validators can have the incoming p2p disabled
			isInboundP2PDisabled = i == simParams.NodeWithInboundP2PDisabled
		}

		genesis := "{}"
		if simParams.WithPrefunding {
			genesis = testcommon.TestnetGenesisJSON()
		}
		nodeType, err := common.ToNodeType(nodeTypeStr)
		if err != nil {
			return nil, fmt.Errorf("unable to convert node type (%s): %w", nodeTypeStr, err)
		}
		hostP2PAddress := fmt.Sprintf("127.0.0.1:%d", simParams.StartPort+integration.DefaultHostP2pOffset+i)

		tenCfg, err := config.LoadTenConfig("defaults/sim/1-env-sim.yaml")
		if err != nil {
			return nil, fmt.Errorf("unable to load TEN config: %w", err)
		}
		tenCfg.Network.GenesisJSON = genesis
		tenCfg.Network.Sequencer.P2PAddress = fmt.Sprintf("127.0.0.1:%d", simParams.StartPort+integration.DefaultHostP2pOffset)
		tenCfg.Network.L1.BlockTime = simParams.AvgBlockDuration
		tenCfg.Network.L1.L1Contracts.ManagementContract = simParams.L1TenData.MgmtContractAddress
		tenCfg.Network.L1.L1Contracts.MessageBusContract = simParams.L1TenData.MessageBusAddr
		tenCfg.Network.Gas.PaymentAddress = simParams.Wallets.L2FeesWallet.Address()

		tenCfg.Node.PrivateKeyString = privateKey
		tenCfg.Node.HostAddress = hostP2PAddress
		tenCfg.Node.NodeType = nodeType
		tenCfg.Node.IsGenesis = i == 0
		tenCfg.Node.ID = hostAddress
		tenCfg.Host.P2P.IsDisabled = isInboundP2PDisabled
		tenCfg.Host.P2P.BindAddress = hostP2PAddress
		tenCfg.Host.RPC.HTTPPort = uint64(simParams.StartPort + integration.DefaultHostRPCHTTPOffset + i)
		tenCfg.Host.RPC.WSPort = uint64(simParams.StartPort + integration.DefaultHostRPCWSOffset + i)
		tenCfg.Host.Enclave.RPCAddresses = []string{fmt.Sprintf("127.0.0.1:%d", simParams.StartPort+integration.DefaultEnclaveOffset+i)}
		tenCfg.Host.L1.WebsocketURL = fmt.Sprintf("ws://127.0.0.1:%d", simParams.StartPort+100)
		tenCfg.Host.L1.L1BeaconUrl = beaconURL
		tenCfg.Host.Log.Level = 1
		tenCfg.Enclave.Log.Level = 1
		tenCfg.Enclave.RPC.BindAddress = fmt.Sprintf("127.0.0.1:%d", simParams.StartPort+integration.DefaultEnclaveOffset+i)

		// create the nodes
		nodes[i] = noderunner.NewInMemNode(tenCfg)

		// start the nodes
		err = nodes[i].Start()
		if err != nil {
			errCheck := checkProcessPort(err.Error())
			if errCheck != nil {
				testlog.Logger().Warn("no port found on error", log.ErrKey, err)
			}
			fmt.Printf("unable to start TEN node: %s\n", err)
			testlog.Logger().Error("unable to start TEN node ", log.ErrKey, err)
		}
	}

	// create the l2 and eth connections
	err = n.createConnections(simParams)
	if err != nil {
		testlog.Logger().Crit("unable to create node connections", log.ErrKey, err)
	}
	walletClients := createAuthClientsPerWallet(n.l2Clients, simParams.Wallets)

	time.Sleep(15 * simParams.AvgBlockDuration)

	// permission the sequencer enclaveID
	// we retry fetching the seqHealth until it comes back with the enclaveID as the nodes are still starting up
	startTime := time.Now()
	var seqEnclaveID *common.EnclaveID
	for ; seqEnclaveID == nil; time.Sleep(100 * time.Millisecond) {
		seqHealth, _ := n.tenClients[0].Health()
		if seqHealth.Enclaves == nil || len(seqHealth.Enclaves) == 0 {
			continue
		}
		seqEnclaveID = &seqHealth.Enclaves[0].EnclaveID
		if time.Now().After(startTime.Add(2 * time.Minute)) {
			return nil, fmt.Errorf("unable to get sequencer enclaveID after 2 minutes")
		}
	}

	// permission the sequencer enclaveID (also requires retries as the enclaveID may not be attested yet)
	err = PermissionTenSequencerEnclave(n.wallets.MCOwnerWallet, n.gethClients[0], simParams.L1TenData.MgmtContractAddress, *seqEnclaveID)
	if err != nil {
		return nil, fmt.Errorf("unable to permission sequencer enclaveID: %w", err)
	}

	// wait for nodes to be healthy now we've permissioned
	// make sure the nodes are healthy
	for _, client := range n.tenClients {
		startTime := time.Now()
		healthy := false
		for ; !healthy; time.Sleep(500 * time.Millisecond) {
			h, _ := client.Health()
			healthy = h.OverallHealth
			if !healthy {
				fmt.Println(h.Errors)
			}
			if time.Now().After(startTime.Add(3 * time.Minute)) {
				return nil, fmt.Errorf("nodes not healthy after 3 minutes")
			}
		}
	}

	return &RPCHandles{
		EthClients:     n.gethClients,
		TenClients:     n.tenClients,
		RPCClients:     n.l2Clients,
		AuthObsClients: walletClients,
	}, nil
}

func (n *networkOfSocketNodes) TearDown() {
	// Stop the TEN nodes first (each host will attempt to shut down its enclave as part of shutdown).
	StopTenNodes(n.l2Clients)
	CheckHostRPCServersStopped(n.hostWebsocketURLs)
	StopEth2Network(n.gethClients, n.eth2Network)
}

func (n *networkOfSocketNodes) createConnections(simParams *params.SimParams) error {
	// create the clients in the structs
	n.l2Clients = make([]rpc.Client, simParams.NumberOfNodes)
	n.hostWebsocketURLs = make([]string, simParams.NumberOfNodes)
	n.tenClients = make([]*obsclient.ObsClient, simParams.NumberOfNodes)

	for i := 0; i < simParams.NumberOfNodes; i++ {
		var client rpc.Client
		var err error

		// create a connection to the newly created nodes - panic if no connection is made after some time
		startTime := time.Now()
		for connected := false; !connected; time.Sleep(500 * time.Millisecond) {
			port := simParams.StartPort + integration.DefaultHostRPCWSOffset + i
			client, err = rpc.NewNetworkClient(fmt.Sprintf("ws://127.0.0.1:%d", port))
			connected = err == nil // The client cannot be created until the node has started.
			if time.Now().After(startTime.Add(2 * time.Minute)) {
				return fmt.Errorf("failed to create a connect to node after 2 minute - %w", err)
			}

			testlog.Logger().Info(fmt.Sprintf("Could not create client %d at port %d. Retrying...", i, port), log.ErrKey, err)
		}

		n.l2Clients[i] = client
		n.hostWebsocketURLs[i] = fmt.Sprintf("ws://%s:%d", Localhost, simParams.StartPort+integration.DefaultHostRPCWSOffset+i)
	}

	for idx, l2Client := range n.l2Clients {
		n.tenClients[idx] = obsclient.NewObsClient(l2Client)
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
