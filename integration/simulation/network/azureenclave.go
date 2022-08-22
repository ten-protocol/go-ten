package network

import (
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/integration/gethnetwork"

	"github.com/obscuronet/go-obscuro/go/rpcclientlib"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/obscuronet/go-obscuro/integration/simulation/params"

	"github.com/obscuronet/go-obscuro/integration/simulation/stats"
)

const enclavePort = 11000

// creates Obscuro nodes with their own enclave servers that communicate with peers via sockets, wires them up, and populates the network objects
type networkWithAzureEnclaves struct {
	gethNetwork *gethnetwork.GethNetwork
	gethClients []ethadapter.EthClient
	wallets     *params.SimWallets

	obscuroClients  []rpcclientlib.Client
	azureEnclaveIps []string

	hostRPCAddresses []string
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

func (n *networkWithAzureEnclaves) Create(params *params.SimParams, stats *stats.Stats) (*RPCHandles, error) {
	params.MgmtContractAddr, params.ObxErc20Address, params.EthErc20Address, n.gethClients, n.gethNetwork = SetUpGethNetwork(
		n.wallets,
		params.StartPort,
		params.NumberOfNodes,
		int(params.AvgBlockDuration.Seconds()),
	)
	params.MgmtContractLib = mgmtcontractlib.NewMgmtContractLib(params.MgmtContractAddr)
	params.ERC20ContractLib = erc20contractlib.NewERC20ContractLib(params.MgmtContractAddr, params.ObxErc20Address, params.EthErc20Address)

	fmt.Printf("Please start the edgeless DB instances. Then start the docker image on the azure server with below cmds:\n")
	for i := 0; i < len(n.azureEnclaveIps); i++ {
		fmt.Printf("sudo docker run --net enclavenet --name enclave -h enclave -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p %d:%d/tcp obscuro_enclave --willAttest --useInMemoryDB=false "+
			"--edgelessDBHost obscuroedb --hostID %d --address :11000 --managementContractAddress %s  --erc20ContractAddresses %s,%s\n",
			enclavePort, enclavePort, i, params.MgmtContractAddr.Hex(), params.ObxErc20Address.Hex(), params.EthErc20Address.Hex())
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

	obscuroClients, walletClients, hostRPCAddresses := startStandaloneObscuroNodes(params, stats, n.gethClients, n.enclaveAddresses)
	n.obscuroClients = obscuroClients
	n.hostRPCAddresses = hostRPCAddresses

	return &RPCHandles{
		EthClients:                    n.gethClients,
		ObscuroClients:                n.obscuroClients,
		VirtualWalletExtensionClients: walletClients,
	}, nil
}

func (n *networkWithAzureEnclaves) TearDown() {
	// First stop the obscuro nodes
	StopObscuroNodes(n.obscuroClients)
	StopGethNetwork(n.gethClients, n.gethNetwork)
	CheckHostRPCServersStopped(n.hostRPCAddresses)
}
