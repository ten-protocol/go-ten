package bridge

import (
	"bytes"
	"fmt"
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	crypto2 "github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/common"
	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
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

// Bridge encapsulates all logic around processing the interactions with an L1
type Bridge struct {
	SupportedTokens map[ERC20]*ERC20Mapping
	// BridgeAddress The address the bridge on the L2
	BridgeAddress gethcommon.Address

	MgmtContractLib  mgmtcontractlib.MgmtContractLib
	Erc20ContractLib erc20contractlib.ERC20ContractLib

	TransactionBlobCrypto crypto2.TransactionBlobCrypto

	ObscuroChainID  int64
	EthereumChainID int64

	logger gethlog.Logger
}

func New(
	hocAddress *gethcommon.Address,
	pocAddress *gethcommon.Address,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	transactionBlobCrypto crypto2.TransactionBlobCrypto,
	obscuroChainID int64,
	ethereumChainID int64,
	logger gethlog.Logger,
) *Bridge {
	tokens := make(map[ERC20]*ERC20Mapping, 0)

	tokens[HOC] = &ERC20Mapping{
		Name:      HOC,
		L1Address: hocAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), HOCOwner, logger),
		L2Address: &HOCContract,
	}

	tokens[POC] = &ERC20Mapping{
		Name:      POC,
		L1Address: pocAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), POCOwner, logger),
		L2Address: &POCContract,
	}

	return &Bridge{
		SupportedTokens:       tokens,
		BridgeAddress:         BridgeAddress,
		MgmtContractLib:       mgmtContractLib,
		Erc20ContractLib:      erc20ContractLib,
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
		logger:                logger,
	}
}

func (bridge *Bridge) IsWithdrawal(address gethcommon.Address) bool {
	return bytes.Equal(address.Bytes(), bridge.BridgeAddress.Bytes())
}

// L1Address - returns the L1 address of a token based on the mapping
func (bridge *Bridge) L1Address(l2Address *gethcommon.Address) *gethcommon.Address {
	if l2Address == nil {
		return nil
	}
	for _, t := range bridge.SupportedTokens {
		if bytes.Equal(l2Address.Bytes(), t.L2Address.Bytes()) {
			return t.L1Address
		}
	}
	return nil
}

// GetMapping - finds the mapping based on the address that was called in an L1 transaction
func (bridge *Bridge) GetMapping(l1ContractAddress *gethcommon.Address) *ERC20Mapping {
	for _, t := range bridge.SupportedTokens {
		if bytes.Equal(t.L1Address.Bytes(), l1ContractAddress.Bytes()) {
			return t
		}
	}
	return nil
}

// ExtractRollups - returns a list of the rollups published in this block
func (bridge *Bridge) ExtractRollups(b *types.Block, blockResolver db.BlockResolver) []*obscurocore.Rollup {
	rollups := make([]*obscurocore.Rollup, 0)
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
				rollups = append(rollups, obscurocore.ToEnclaveRollup(r, bridge.TransactionBlobCrypto))
				bridge.logger.Trace(fmt.Sprintf("Extracted Rollup r_%d from block b_%d",
					common.ShortHash(r.Hash()),
					common.ShortHash(b.Hash()),
				))
			}
		}
	}
	return rollups
}

// NewDepositTx creates a synthetic Obscuro transfer transaction based on deposits into the L1 bridge.
// Todo - has to go through a few more iterations
func (bridge *Bridge) NewDepositTx(contract *gethcommon.Address, address gethcommon.Address, amount *big.Int, rollupState *state.StateDB, adjustNonce uint64) *common.L2Tx {
	transferERC20data := erc20contractlib.CreateTransferTxData(address, amount)
	signer := types.NewLondonSigner(big.NewInt(bridge.ObscuroChainID))

	token := bridge.GetMapping(contract)
	if token == nil {
		bridge.logger.Crit("This should not happen as we don't generate deposits on unsupported tokens.")
		return nil
	}

	// The nonce is adjusted with the number of deposits added to the rollup already.
	storedNonce := rollupState.GetNonce(token.Owner.Address())
	nonce := storedNonce + adjustNonce

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    gethcommon.Big0,
		Gas:      1_000_000,
		GasPrice: gethcommon.Big1,
		Data:     transferERC20data,
		To:       token.L2Address,
	})

	newTx, err := types.SignTx(tx, signer, token.Owner.PrivateKey())
	if err != nil {
		bridge.logger.Crit("could not sign synthetic deposit tx.", log.ErrKey, err)
		return nil
	}
	return newTx
}

// ExtractDeposits returns a list of L2 deposit transactions generated from the L1 deposit transactions
// starting with the proof of the parent rollup(exclusive) to the proof of the current rollup
func (bridge *Bridge) ExtractDeposits(
	fromBlock *types.Block,
	toBlock *types.Block,
	blockResolver db.BlockResolver,
	rollupState *state.StateDB,
) []*common.L2Tx {
	from := common.GenesisBlock.Hash()
	height := common.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
		if !blockResolver.IsAncestor(toBlock, fromBlock) {
			bridge.logger.Crit("Deposits can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
			return nil
		}
	}

	allDeposits := make([]*common.L2Tx, 0)
	b := toBlock
	for {
		if bytes.Equal(b.Hash().Bytes(), from.Bytes()) {
			break
		}
		for _, tx := range b.Transactions() {
			t := bridge.Erc20ContractLib.DecodeTx(tx)
			if t == nil {
				continue
			}

			if depositTx, ok := t.(*ethadapter.L1DepositTx); ok {
				// todo - the adjust has to be per token
				depL2Tx := bridge.NewDepositTx(depositTx.TokenContract, *depositTx.Sender, depositTx.Amount, rollupState, uint64(len(allDeposits)))
				allDeposits = append(allDeposits, depL2Tx)
			}
		}
		if b.NumberU64() < height {
			bridge.logger.Crit("block height is less than genesis height")
			return nil
		}
		p, err := blockResolver.ParentBlock(b)
		if err != nil {
			// todo - joel - handle error
			bridge.logger.Crit("deposits can't be processed because the rollups are not on the same Ethereum fork")
			return nil
		}
		b = p
	}

	bridge.logger.Trace(fmt.Sprintf("Extracted deposits %d ->%d: %v.", fromBlock.NumberU64(), toBlock.NumberU64(), allDeposits))
	return allDeposits
}

// Todo - this has to be implemented differently based on how we define the ObsERC20
func (bridge *Bridge) RollupPostProcessingWithdrawals(newHeadRollup *obscurocore.Rollup, state *state.StateDB, receiptsMap map[gethcommon.Hash]*types.Receipt) []common.Withdrawal {
	w := make([]common.Withdrawal, 0)
	// go through each transaction and check if the withdrawal was processed correctly
	for _, t := range newHeadRollup.Transactions {
		found, address, amount := erc20contractlib.DecodeTransferTx(t, bridge.logger)

		supportedTokenAddress := bridge.L1Address(t.To())
		if found && supportedTokenAddress != nil && bridge.IsWithdrawal(*address) {
			receipt := receiptsMap[t.Hash()]
			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				signer := types.NewLondonSigner(big.NewInt(bridge.ObscuroChainID))
				from, err := types.Sender(signer, t)
				if err != nil {
					bridge.logger.Crit("Error retrieving the sender from the signature", log.ErrKey, err)
					return nil
				}
				state.Logs()
				w = append(w, common.Withdrawal{
					Contract:  *supportedTokenAddress,
					Amount:    amount,
					Recipient: from,
				})
			}
		}
	}

	// TODO - fix the withdrawals logic
	// clearWithdrawals(state, withdrawalTxs)
	return w
}
