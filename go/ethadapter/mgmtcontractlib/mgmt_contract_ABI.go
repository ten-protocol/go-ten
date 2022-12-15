package mgmtcontractlib

import "github.com/obscuronet/go-obscuro/contracts/generated/ManagementContract"

const (
	AddRollupMethod        = "AddRollup"
	RespondSecretMethod    = "RespondNetworkSecret"
	RequestSecretMethod    = "RequestNetworkSecret"
	InitializeSecretMethod = "InitializeNetworkSecret" //#nosec
	GetHostAddressesMethod = "GetHostAddresses"
)

var MgmtContractABI = ManagementContract.ManagementContractMetaData.ABI
