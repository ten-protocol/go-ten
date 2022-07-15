package mgmtcontractlib

import ManagementContract "github.com/obscuronet/go-obscuro/contracts/compiledcontracts/generatedManagementContract"

const (
	AddRollupMethod        = "AddRollup"
	RespondSecretMethod    = "RespondNetworkSecret"
	RequestSecretMethod    = "RequestNetworkSecret"
	InitializeSecretMethod = "InitializeNetworkSecret" //#nosec
	GetHostAddressesMethod = "GetHostAddresses"
)

var (
	MgmtContractByteCode = ManagementContract.ManagementContractMetaData.Bin[2:]
	MgmtContractABI      = ManagementContract.ManagementContractMetaData.ABI
)
