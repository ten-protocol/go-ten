package smartcontract

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/contracts/compiledcontracts/generatedManagementContract"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	ethereumclient "github.com/ethereum/go-ethereum/ethclient"
)

// debugMgmtContractLib is a wrapper around the MgmtContractLib
// allows the direct use of the generatedManagementContract package
type debugMgmtContractLib struct {
	mgmtcontractlib.MgmtContractLib
	GenContract *generatedManagementContract.ManagementContract
}

// newDebugMgmtContractLib creates an instance of the generated contract package and allows the use of the MgmtContractLib properties
func newDebugMgmtContractLib(address common.Address, client *ethereumclient.Client, mgmtContractLib mgmtcontractlib.MgmtContractLib) *debugMgmtContractLib {
	genContract, err := generatedManagementContract.NewManagementContract(address, client)
	if err != nil {
		panic(err)
	}

	return &debugMgmtContractLib{
		mgmtContractLib,
		genContract,
	}
}

// AwaitedIssueRollup speeds ups the issuance of rollup, await of tx to be minted and makes sure the values are correctly stored
func (d *debugMgmtContractLib) AwaitedIssueRollup(rollup nodecommon.Rollup, client ethclient.EthClient, w *debugWallet) error {
	txData := d.CreateRollup(
		&obscurocommon.L1RollupTx{Rollup: nodecommon.EncodeRollup(&rollup)},
		w.GetNonceAndIncrement(),
	)

	issuedTx, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		_, err := w.debugTransaction(client, issuedTx)
		if err != nil {
			return fmt.Errorf("transaction should have succeeded, expected %d got %d - reason: %w", types.ReceiptStatusSuccessful, receipt.Status, err)
		}
	}

	// rollup meta data is actually stored
	found, rollupElement, err := d.GenContract.GetRollupByHash(nil, rollup.Hash())
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("rollup not stored in tree")
	}

	if rollupElement.Rollup.Number.Int64() != rollup.Header.Number.Int64() ||
		!bytes.Equal(rollupElement.Rollup.ParentHash[:], rollup.Header.ParentHash.Bytes()) ||
		!bytes.Equal(rollupElement.Rollup.AggregatorID[:], rollup.Header.Agg.Bytes()) ||
		!bytes.Equal(rollupElement.Rollup.L1Block[:], rollup.Header.L1Proof.Bytes()) {
		return fmt.Errorf("stored rollup does not match the generated rollup")
	}

	return nil
}
