package mgmtcontractlib

import "github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"

const (
	AddRollupMethod                = "AddRollup"
	RespondSecretMethod            = "RespondNetworkSecret"
	RequestSecretMethod            = "RequestNetworkSecret"
	InitializeSecretMethod         = "InitializeNetworkSecret" //#nosec
	GetHostAddressesMethod         = "GetHostAddresses"
	GetImportantContractKeysMethod = "GetImportantContractKeys"
	SetImportantContractsMethod    = "SetImportantContractAddress"
	GetImportantAddressMethod      = "importantContractAddresses"
)

var MgmtContractABI = ManagementContract.ManagementContractMetaData.ABI
