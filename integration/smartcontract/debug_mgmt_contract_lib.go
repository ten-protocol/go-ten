package smartcontract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/obscuronet/obscuro-playground/contracts/compiledcontracts/generatedManagementContract"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
)

type debugMgmtContractLib struct {
	mgmtcontractlib.MgmtContractLib
	genContract *generatedManagementContract.GeneratedManagementContract
}

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
