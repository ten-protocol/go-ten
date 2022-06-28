package bridge

import (
	"bytes"
	"math/big"

	crypto2 "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/crypto"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/wallet"
)

// Todo - remove all hardcoded values in the next iteration.
// The Contract addresses are the result of the deploying a smart contract from hardcoded owners.
// The "owners" are keys which are the de-facto "admins" of those erc20s and are able to transfer or mint tokens.
// The contracts and addresses cannot be random for now, because there is hardcoded logic in the core
// to generate synthetic "transfer" transactions for each erc20 deposit on ethereum
// and these transactions need to be signed. Which means the platform needs to "own" ERC20s.

// ERC20 - the supported ERC20 tokens. A list of made-up tokens used for testing.
// Todo - this will be removed together will all the keys and addresses.
type ERC20 int

const (
	BTC ERC20 = iota
	ETH
)

var WBtcOwner, _ = crypto.HexToECDSA("6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682")

// WBtcContract X- address of the deployed "btc" erc20 on the L2
var WBtcContract = common.BytesToAddress(common.Hex2Bytes("f3a8bd422097bFdd9B3519Eaeb533393a1c561aC"))

var WEthOnwer, _ = crypto.HexToECDSA("4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8")

// WEthContract - address of the deployed "eth" erc20 on the L2
var WEthContract = common.BytesToAddress(common.Hex2Bytes("9802F661d17c65527D7ABB59DAAD5439cb125a67"))

// BridgeAddress - address of the virtual bridge
var BridgeAddress = common.BytesToAddress(common.Hex2Bytes("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))

// ERC20Mapping - maps an L1 Erc20 to an L2 Erc20 address
type ERC20Mapping struct {
	Name ERC20

	// L1Owner   wallet.Wallet
	L1Address *common.Address

	Owner     wallet.Wallet // for now the wrapped L2 version is owned by a wallet, but this will change
	L2Address *common.Address
}

// Bridge encapsulates all logic around processing the interactions with an L1
type Bridge struct {
	SupportedTokens map[ERC20]*ERC20Mapping
	// BridgeAddress The address the bridge on the L2
	BridgeAddress common.Address

	MgmtContractLib  mgmtcontractlib.MgmtContractLib
	Erc20ContractLib erc20contractlib.ERC20ContractLib

	NodeID                uint64
	TransactionBlobCrypto crypto2.TransactionBlobCrypto

	ObscuroChainID  int64
	EthereumChainID int64
}

func New(
	btcAddress *common.Address,
	ethAddress *common.Address,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	nodeID uint64,
	transactionBlobCrypto crypto2.TransactionBlobCrypto,
	obscuroChainID int64,
	ethereumChainID int64,
) *Bridge {
	tokens := make(map[ERC20]*ERC20Mapping, 0)

	tokens[BTC] = &ERC20Mapping{
		Name:      BTC,
		L1Address: btcAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), WBtcOwner),
		L2Address: &WBtcContract,
	}

	tokens[ETH] = &ERC20Mapping{
		Name:      ETH,
		L1Address: ethAddress,
		Owner:     wallet.NewInMemoryWalletFromPK(big.NewInt(obscuroChainID), WEthOnwer),
		L2Address: &WEthContract,
	}

	return &Bridge{
		SupportedTokens:       tokens,
		BridgeAddress:         BridgeAddress,
		MgmtContractLib:       mgmtContractLib,
		Erc20ContractLib:      erc20ContractLib,
		NodeID:                nodeID,
		TransactionBlobCrypto: transactionBlobCrypto,
		ObscuroChainID:        obscuroChainID,
		EthereumChainID:       ethereumChainID,
	}
}

func (bridge *Bridge) IsWithdrawal(address common.Address) bool {
	return bytes.Equal(address.Bytes(), bridge.BridgeAddress.Bytes())
}

// L1Address - returns the L1 address of a token based on the mapping
func (bridge *Bridge) L1Address(l2Address *common.Address) *common.Address {
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
func (bridge *Bridge) GetMapping(l1ContractAddress *common.Address) *ERC20Mapping {
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

		if rolTx, ok := t.(*obscurocommon.L1RollupTx); ok {
			r := nodecommon.DecodeRollupOrPanic(rolTx.Rollup)

			// Ignore rollups created with proofs from different L1 blocks
			// In case of L1 reorgs, rollups may end published on a fork
			if blockResolver.IsBlockAncestor(b, r.Header.L1Proof) {
				rollups = append(rollups, bridge.TransactionBlobCrypto.ToEnclaveRollup(r))
				nodecommon.LogWithID(bridge.NodeID, "Extracted Rollup r_%d from block b_%d",
					obscurocommon.ShortHash(r.Hash()),
					obscurocommon.ShortHash(b.Hash()),
				)
			}
		}
	}
	return rollups
}

// this function creates a synthetic Obscuro transfer transaction based on deposits into the L1 bridge.
// Todo - has to go through a few more iterations
func (bridge *Bridge) NewDepositTx(contract *common.Address, address common.Address, amount uint64, rollupState *state.StateDB, adjustNonce uint64) *nodecommon.L2Tx {
	transferERC20data := erc20contractlib.CreateTransferTxData(address, amount)
	signer := types.NewLondonSigner(big.NewInt(bridge.ObscuroChainID))

	token := bridge.GetMapping(contract)
	if token == nil {
		panic("This should not happen as we don't generate deposits on unsupported tokens.")
	}

	// The nonce is adjusted with the number of deposits added to the rollup already.
	storedNonce := rollupState.GetNonce(token.Owner.Address())
	nonce := storedNonce + adjustNonce

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    common.Big0,
		Gas:      1_000_000,
		GasPrice: common.Big0,
		Data:     transferERC20data,
		To:       token.L2Address,
	})

	newTx, err := types.SignTx(tx, signer, token.Owner.PrivateKey())
	if err != nil {
		log.Panic("could not sign synthetic deposit tx. Cause: %s", err)
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
) []*nodecommon.L2Tx {
	from := obscurocommon.GenesisBlock.Hash()
	height := obscurocommon.L1GenesisHeight
	if fromBlock != nil {
		from = fromBlock.Hash()
		height = fromBlock.NumberU64()
		if !blockResolver.IsAncestor(toBlock, fromBlock) {
			log.Panic("Deposits can't be processed because the rollups are not on the same Ethereum fork. This should not happen.")
		}
	}

	allDeposits := make([]*nodecommon.L2Tx, 0)
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

			if depositTx, ok := t.(*obscurocommon.L1DepositTx); ok {
				// todo - the adjust has to be per token
				depL2Tx := bridge.NewDepositTx(depositTx.TokenContract, *depositTx.Sender, depositTx.Amount, rollupState, uint64(len(allDeposits)))
				allDeposits = append(allDeposits, depL2Tx)
			}
		}
		if b.NumberU64() < height {
			log.Panic("block height is less than genesis height")
		}
		p, f := blockResolver.ParentBlock(b)
		if !f {
			log.Panic("deposits can't be processed because the rollups are not on the same Ethereum fork")
		}
		b = p
	}

	log.Info("Extracted deposits %d ->%d: %v.", fromBlock.NumberU64(), toBlock.NumberU64(), allDeposits)
	return allDeposits
}

// Todo - this has to be implemented differently based on how we define the ObsERC20
// this belongs in the bridge
func (bridge *Bridge) RollupPostProcessingWithdrawals(newHeadRollup *obscurocore.Rollup, state *state.StateDB, receiptsMap map[common.Hash]*types.Receipt) []nodecommon.Withdrawal {
	w := make([]nodecommon.Withdrawal, 0)
	// go through each transaction and check if the withdrawal was processed correctly
	for _, t := range newHeadRollup.Transactions {
		found, address, amount := erc20contractlib.DecodeTransferTx(t)

		supportedTokenAddress := bridge.L1Address(t.To())
		if found && supportedTokenAddress != nil && bridge.IsWithdrawal(*address) {
			receipt := receiptsMap[t.Hash()]
			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				signer := types.NewLondonSigner(big.NewInt(bridge.ObscuroChainID))
				from, err := types.Sender(signer, t)
				if err != nil {
					panic(err)
				}
				state.Logs()
				w = append(w, nodecommon.Withdrawal{
					Contract:  *supportedTokenAddress,
					Amount:    amount.Uint64(),
					Recipient: from,
				})
			}
		}
	}

	// TODO - fix the withdrawals logic
	// clearWithdrawals(state, withdrawalTxs)
	return w
}
