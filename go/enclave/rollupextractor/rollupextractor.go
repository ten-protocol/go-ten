package rollupextractor

// TODO: once the cross chain messages based bridge is implemented remove this completely.

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	crypto2 "github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// Todo - remove all hardcoded values in the next iteration.
// The Contract addresses are the result of the deploying a smart contract from hardcoded owners.
// The "owners" are keys which are the de-facto "admins" of those erc20s and are able to transfer or mint tokens.
// The contracts and addresses cannot be random for now, because there is hardcoded logic in the core
// to generate synthetic "transfer" transactions for each erc20 deposit on ethereum
// and these transactions need to be signed. Which means the platform needs to "own" ERC20s.

// ERC20 - the supported ERC20 tokens. A list of made-up tokens used for testing.
// Todo - this will be removed together will all the keys and addresses.
type ERC20 string

const (
	HOC            ERC20 = "HOC"
	POC            ERC20 = "POC"
	HOCAddr              = "f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"
	pocAddr              = "9802F661d17c65527D7ABB59DAAD5439cb125a67"
	bridgeAddr           = "deB34A740ECa1eC42C8b8204CBEC0bA34FDD27f3"
	hocOwnerKeyHex       = "6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
	pocOwnerKeyHex       = "4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"
)

var HOCOwner, _ = crypto.HexToECDSA(hocOwnerKeyHex)

// HOCContract - address of the deployed "hocus" erc20 on the L2
var HOCContract = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(HOCAddr))

var POCOwner, _ = crypto.HexToECDSA(pocOwnerKeyHex)

// POCContract - address of the deployed "pocus" erc20 on the L2
var POCContract = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(pocAddr))

// BridgeAddress - address of the virtual bridge
var BridgeAddress = gethcommon.BytesToAddress(gethcommon.Hex2Bytes(bridgeAddr))

// ERC20Mapping - maps an L1 Erc20 to an L2 Erc20 address
type ERC20Mapping struct {
	Name ERC20

	// L1Owner   wallet.Wallet
	L1Address *gethcommon.Address

	Owner     wallet.Wallet // for now the wrapped L2 version is owned by a wallet, but this will change
	L2Address *gethcommon.Address
}

// RollupExtractor encapsulates the logic of decoding rollup transactions submitted to the L1 and resolving them
// to rollups that the enclave can process.
type RollupExtractor struct {
	MgmtContractLib mgmtcontractlib.MgmtContractLib

	TransactionBlobCrypto crypto2.TransactionBlobCrypto

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger
}

func New(
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	transactionBlobCrypto crypto2.TransactionBlobCrypto,
	obscuroChainID int64,
	ethereumChainID int64,
	logger gethlog.Logger,
) *RollupExtractor {
	return &RollupExtractor{
		MgmtContractLib:       mgmtContractLib,
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
		logger:                logger,
	}
}

// ExtractRollups - returns a list of the rollups published in this block
func (bridge *RollupExtractor) ExtractRollups(b *types.Block, blockResolver db.BlockResolver) []*core.Rollup {
	rollups := make([]*core.Rollup, 0)
	for _, tx := range b.Transactions() {
		// go through all rollup transactions
		t := bridge.MgmtContractLib.DecodeTx(tx)
		if t == nil {
			continue
		}

		if rolTx, ok := t.(*ethadapter.L1RollupTx); ok {
			r, err := common.DecodeRollup(rolTx.Rollup)
			if err != nil {
				bridge.logger.Crit("could not decode rollup.", log.ErrKey, err)
				return nil
			}

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
				rollups = append(rollups, core.ToRollup(r, bridge.TransactionBlobCrypto))
				bridge.logger.Trace(fmt.Sprintf("Extracted Rollup r_%d from block b_%d",
					common.ShortHash(r.Hash()),
					common.ShortHash(b.Hash()),
				))
			}
		}
	}
	return rollups
}
