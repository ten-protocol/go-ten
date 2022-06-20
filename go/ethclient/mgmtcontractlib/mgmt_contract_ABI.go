package mgmtcontractlib

import "github.com/obscuronet/obscuro-playground/contracts/compiledcontracts/generatedManagementContract"

const (
	AddRollupMethod        = "AddRollup"
	RespondSecretMethod    = "RespondNetworkSecret"
	RequestSecretMethod    = "RequestNetworkSecret"
	InitializeSecretMethod = "InitializeNetworkSecret" //#nosec
)

var (
	MgmtContractByteCode = generatedManagementContract.ManagementContractMetaData.Bin[2:]
	MgmtContractABI      = generatedManagementContract.ManagementContractMetaData.ABI
)
