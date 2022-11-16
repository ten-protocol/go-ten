package smartcontract

import (
	"bytes"
	"fmt"

	generatedManagementContract "github.com/obscuronet/go-obscuro/contracts/managementcontract/generated/ManagementContract"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	gethcommon "github.com/ethereum/go-ethereum/common"
	ethereumclient "github.com/ethereum/go-ethereum/ethclient"
)

// debugMgmtContractLib is a wrapper around the MgmtContractLib
// allows the direct use of the generatedManagementContract package
type debugMgmtContractLib struct {
	mgmtcontractlib.MgmtContractLib
	GenContract *generatedManagementContract.ManagementContract
}

// newDebugMgmtContractLib creates an instance of the generated contract package and allows the use of the MgmtContractLib properties
func newDebugMgmtContractLib(address gethcommon.Address, client *ethereumclient.Client, mgmtContractLib mgmtcontractlib.MgmtContractLib) *debugMgmtContractLib {
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
func (d *debugMgmtContractLib) AwaitedIssueRollup(rollup common.ExtRollup, client ethadapter.EthClient, w *debugWallet) error {
	encodedRollup, err := common.EncodeRollup(&rollup)
	if err != nil {
		return err
	}
	txData := d.CreateRollup(
		&ethadapter.L1RollupTx{Rollup: encodedRollup},
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
