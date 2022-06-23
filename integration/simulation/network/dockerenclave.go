package network

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/pkg/stdcopy"

	"github.com/ethereum/go-ethereum/log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

const (
	enclaveDockerImg  = "obscuro_enclave"
	enclaveAddress    = ":11000"
	enclaveDockerPort = "11000/tcp"
)

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type basicNetworkOfNodesWithDockerEnclave struct {
	obscuroClients   []obscuroclient.Client
	enclaveAddresses []string
	// Geth
	gethNetwork *gethnetwork.GethNetwork
	gethClients []ethclient.EthClient
	wallets     *params.SimWallets

	// Docker
	ctx              context.Context
	client           *client.Client
	containerIDs     map[string]string
	containerStreams map[string]*types.HijackedResponse
}

func NewBasicNetworkOfNodesWithDockerEnclave(wallets *params.SimWallets) Network {
	return &basicNetworkOfNodesWithDockerEnclave{
		wallets:          wallets,
		containerStreams: map[string]*types.HijackedResponse{},
	}
}

// Create initializes Obscuro nodes with their own Dockerised enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
// TODO - Use individual Docker containers for the Obscuro nodes and Ethereum nodes.
func (n *basicNetworkOfNodesWithDockerEnclave) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []obscuroclient.Client, error) {
	// We create Docker client, and finish early if docker or the enclave image are not available.
	if err := n.setupAndCheckDocker(); err != nil {
		return nil, nil, err
	}

	// We start a geth network with all necessary contracts deployed.
	params.MgmtContractAddr, params.Erc20Address, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)
	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.Erc20Address)

	// Start the enclave docker containers with the right addresses.
	n.startDockerEnclaves(params)

	n.enclaveAddresses = make([]string, params.NumberOfNodes)
	for i := 0; i < params.NumberOfNodes; i++ {
		n.enclaveAddresses[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
	}

	// Start the standalone obscuro nodes connected to the enclaves and to the geth nodes
	obscuroClients := startStandaloneObscuroNodes(params, stats, n.gethClients, n.enclaveAddresses)
	n.obscuroClients = obscuroClients

	return n.gethClients, obscuroClients, nil
}

func (n *basicNetworkOfNodesWithDockerEnclave) TearDown() {
	// First stop the obscuro nodes
	StopObscuroNodes(n.obscuroClients)

	StopGethNetwork(n.gethClients, n.gethNetwork)
	terminateDockerContainers(n.ctx, n.client, n.containerIDs, n.containerStreams)
}

func (n *basicNetworkOfNodesWithDockerEnclave) setupAndCheckDocker() error {
	n.ctx = context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	n.client = cli
	// We check the required Docker images are available.
	if !dockerImagesAvailable(n.ctx, cli) {
		// We don't cause the test to fail here, because we want users to be able to run all the tests in the repo
		// without having to build the Docker images.
		return fmt.Errorf("this test requires the `%s` Docker image to be built using `dockerfiles/enclave.Dockerfile`. Terminating", enclaveDockerImg)
	}
	return nil
}

func (n *basicNetworkOfNodesWithDockerEnclave) startDockerEnclaves(params *params.SimParams) {
	// We create the Docker containers and set up a hook to terminate them at the end of the test.
	n.containerIDs = createDockerContainers(n.ctx, n.client, params.NumberOfNodes, params.StartPort, params.MgmtContractAddr.Hex(), params.Erc20Address.Hex())

	// We start the Docker containers.
	for id := range n.containerIDs {
		if err := n.client.ContainerStart(n.ctx, id, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
		waiter, err := n.client.ContainerAttach(n.ctx, id, types.ContainerAttachOptions{
			Stderr: true,
			Stdout: true,
			Stdin:  false,
			Stream: true,
		})

		go func() {
			_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, waiter.Reader)
			if err != nil {
				log.Error("Could not copy output from the docker container")
			}
		}()

		if err != nil {
			panic(err)
		}
		n.containerStreams[id] = &waiter
	}
}
