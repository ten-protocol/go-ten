package mgmtcontractlib

import "github.com/obscuronet/obscuro-playground/contracts/compiledcontracts/generatedManagementContract"

const (
	AddRollupMethod        = "AddRollup"
	RespondSecretMethod    = "RespondNetworkSecret"
	RequestSecretMethod    = "RequestNetworkSecret"
	InitializeSecretMethod = "InitializeNetworkSecret" //#nosec
	GetHostAddressesMethod = "GetHostAddresses"
)

var (
	MgmtContractByteCode = generatedManagementContract.GeneratedManagementContractMetaData.Bin[2:]
	MgmtContractABI      = generatedManagementContract.GeneratedManagementContractMetaData.ABI
)
