package smartcontract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/obscuro-playground/contracts/compiledcontracts/generatedManagementContract"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
)

// debugMgmtContractLib is a wrapper around the MgmtContractLib
// allows the direct use of the generatedManagementContract package
type debugMgmtContractLib struct {
	mgmtcontractlib.MgmtContractLib
	genContract *generatedManagementContract.GeneratedManagementContract
}

// newDebugMgmtContractLib creates an instance of the generated contract package and allows the use of the MgmtContractLib properties
func newDebugMgmtContractLib(address common.Address, client *ethclient.Client, mgmtContractLib mgmtcontractlib.MgmtContractLib) *debugMgmtContractLib {
	genContract, err := generatedManagementContract.NewGeneratedManagementContract(address, client)
	if err != nil {
		panic(err)
	}

	return &debugMgmtContractLib{
		mgmtContractLib,
		genContract,
	}
}
