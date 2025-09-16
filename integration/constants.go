package integration

import (
	"reflect"
)

const (
	DefaultGethWSPortOffset         = 100
	DefaultGethAUTHPortOffset       = 200
	DefaultGethNetworkPortOffset    = 300
	DefaultGethHTTPPortOffset       = 400
	DefaultPrysmP2PPortOffset       = 500
	DefaultPrysmRPCPortOffset       = 550
	DefaultPrysmGatewayPortOffset   = 560
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
	TenChainID      = 5443
)

const (
	GethNodeAddress = "0x123463a4b065722e99115d6c222f267d9cabb524"
	GethNodePK      = "2e0834786285daccd064ca17f1654f67b4aef298acbb82cef9ec422fb4975622"
)

type Ports struct {
	TestStartPosEth2NetworkPort                 int
	TestTenscanPort                             int
	TestCanStartStandaloneTenHostAndEnclavePort int
	TestGethSimulationPort                      int
	TestInMemoryMonteCarloSimulationPort        int
	TestFullNetworkMonteCarloSimulationPort     int
	DoNotUSePort                                int
	TestManagementContractPort                  int
	TestCanDeployLayer2ERC20ContractPort        int
	TestFaucetSendsFundsOnlyIfNeededPort        int
	TestFaucetPort                              int
	TestFaucetHTTPPort                          int
	TestTenGatewayPort                          int
	NetworkTestsPort                            int
}

var TestPorts = Ports{
	TestStartPosEth2NetworkPort:                 10000,
	TestTenscanPort:                             11000,
	TestCanStartStandaloneTenHostAndEnclavePort: 12000,
	TestGethSimulationPort:                      14000,
	TestInMemoryMonteCarloSimulationPort:        15000,
	TestFullNetworkMonteCarloSimulationPort:     16000,
	DoNotUSePort:                                17000,
	TestManagementContractPort:                  18000,
	TestCanDeployLayer2ERC20ContractPort:        19000,
	TestFaucetSendsFundsOnlyIfNeededPort:        21000,
	TestFaucetPort:                              22000,
	TestFaucetHTTPPort:                          23000,
	TestTenGatewayPort:                          24000,
	NetworkTestsPort:                            25000,
}

// GetTestName looks up the test name from the port number using reflection
func GetTestName(port int) string {
	port = port - DefaultGethNetworkPortOffset
	val := reflect.ValueOf(TestPorts)
	typ := reflect.TypeOf(TestPorts)

	for i := 0; i < val.NumField(); i++ {
		fieldValue, ok := val.Field(i).Interface().(int)
		if ok && fieldValue == port {
			return typ.Field(i).Name
		}
	}
	return "UnknownTest"
}
