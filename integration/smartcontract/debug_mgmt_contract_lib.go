package smartcontract

import (
	"bytes"
	"fmt"

	generatedManagementContract "github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

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
	txData, err := d.CreateBlobRollup(&common.L1RollupTx{Rollup: encodedRollup})
	if err != nil {
		return fmt.Errorf("failed to create blob rollup: %w", err)
	}

	issuedTx, receipt, err := w.AwaitedSignAndSendTransaction(client, txData)
	if err != nil {
		return fmt.Errorf("failed to send and await transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		debugOutput, debugErr := w.debugTransaction(client, issuedTx)
		if debugErr != nil {
			return fmt.Errorf("transaction failed with status %d and debug failed: %v", receipt.Status, debugErr)
		}
		return fmt.Errorf("transaction failed with status %d: %s", receipt.Status, string(debugOutput))
	}

	// rollup meta data is actually stored
	found, rollupElement, err := d.GenContract.GetRollupByHash(nil, rollup.Hash())
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("rollup not stored in tree")
	}

	if !bytes.Equal(rollupElement.Signature, rollup.Header.Signature) {
		return fmt.Errorf("stored rollup does not match the generated rollup")
	}

	return nil
}
