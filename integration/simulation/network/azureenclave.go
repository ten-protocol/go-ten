package network

import (
	"fmt"
	"time"

	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/obscuronet/obscuro-playground/go/ethclient"

	"github.com/obscuronet/obscuro-playground/integration/simulation/params"

	"github.com/obscuronet/obscuro-playground/integration/simulation/stats"
)

const enclavePort = 11000

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkWithAzureEnclaves struct {
	gethNetwork *gethnetwork.GethNetwork
	gethClients []ethclient.EthClient
	wallets     *params.SimWallets

	obscuroClients  []obscuroclient.Client
	azureEnclaveIps []string

	enclaveAddresses []string
}

func NewNetworkWithAzureEnclaves(enclaveIps []string, wallets *params.SimWallets) Network {
	if len(enclaveIps) == 0 {
		panic("Cannot create azure enclaves network without at least one enclave address.")
	}
	return &networkWithAzureEnclaves{
		azureEnclaveIps: enclaveIps,
		wallets:         wallets,
	}
}

func (n *networkWithAzureEnclaves) Create(params *params.SimParams, stats *stats.Stats) ([]ethclient.EthClient, []obscuroclient.Client, []string, error) {
	params.MgmtContractAddr, params.Erc20Address, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)
	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.Erc20Address)

	fmt.Printf("Please start the docker image on the azure server with with:\n")
	for i := 0; i < len(n.azureEnclaveIps); i++ {
		fmt.Printf("sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p %d:%d/tcp obscuro_enclave --hostID %d --address :11000 --managementContractAddress %s  --erc20ContractAddresses %s\n", enclavePort, enclavePort, i, params.MgmtContractAddr.Hex(), params.Erc20Address.Hex())
	}
	time.Sleep(10 * time.Second)

	// Start the rest of the enclaves
	startRemoteEnclaveServers(len(n.azureEnclaveIps), params, stats)

	n.enclaveAddresses = make([]string, params.NumberOfNodes)
	for i := 0; i < len(n.azureEnclaveIps); i++ {
		n.enclaveAddresses[i] = fmt.Sprintf("%s:%d", n.azureEnclaveIps[i], enclavePort)
	}
	for i := len(n.azureEnclaveIps); i < params.NumberOfNodes; i++ {
		n.enclaveAddresses[i] = fmt.Sprintf("%s:%d", Localhost, params.StartPort+DefaultEnclaveOffset+i)
	}

	obscuroClients, nodeP2pAddrs := startStandaloneObscuroNodes(params, stats, n.gethClients, n.enclaveAddresses)
	n.obscuroClients = obscuroClients

	return n.gethClients, n.obscuroClients, nodeP2pAddrs, nil
}

func (n *networkWithAzureEnclaves) TearDown() {
	// First stop the obscuro nodes
	StopObscuroNodes(n.obscuroClients)
	StopGethNetwork(n.gethClients, n.gethNetwork)
}
