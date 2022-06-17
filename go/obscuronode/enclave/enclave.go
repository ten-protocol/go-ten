package enclave

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/rpcencryptionmanager"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/sql"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/mempool"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
)

const (
	msgNoRollup  = "could not fetch rollup"
	DummyBalance = "0x0"
	// EnclavePrivateKeyHex is the private key used for sensitive communication with the enclave.
	// TODO - Replace this fixed key with a derived key.
	enclavePrivateKeyHex = "81acce9620f0adf1728cb8df7f6b8b8df857955eb9e8b7aed6ef8390c09fc207"

	// The relevant fields in a Call request's params.
	CallFieldTo   = "to"
	CallFieldFrom = "from"
	CallFieldData = "data"
)

type StatsCollector interface {
	// L2Recalc registers when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	RollupWithMoreRecentProof()
}

type enclaveImpl struct {
	config               config.EnclaveConfig
	nodeShortID          uint64
	storage              db.Storage
	blockResolver        db.BlockResolver
	mempool              mempool.Manager
	statsCollector       StatsCollector
	l1Blockchain         *core.BlockChain
	rpcEncryptionManager rpcencryptionmanager.RPCEncryptionManager

	txCh                 chan *nodecommon.L2Tx
	roundWinnerCh        chan *obscurocore.Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork

	mgmtContractLib       mgmtcontractlib.MgmtContractLib
	erc20ContractLib      erc20contractlib.ERC20ContractLib
	attestationProvider   AttestationProvider // interface for producing attestation reports and verifying them
	publicKeySerialized   []byte
	privateKey            *ecdsa.PrivateKey
	transactionBlobCrypto obscurocore.TransactionBlobCrypto

	blockProcessingMutex sync.Mutex
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(
	config config.EnclaveConfig,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	erc20ContractLib erc20contractlib.ERC20ContractLib,
	collector StatsCollector,
) nodecommon.Enclave {
	nodeShortID := obscurocommon.ShortAddress(config.HostID)

	backingDB, err := getDB(nodeShortID, config)
	if err != nil {
		log.Panic("Failed to connect to backing database - %s", err)
	}
	storage := db.NewStorage(backingDB, nodeShortID)

	var l1Blockchain *core.BlockChain
	if config.ValidateL1Blocks {
		if config.GenesisJSON == nil {
			log.Panic("enclave is configured to validate blocks, but genesis JSON is nil")
		}
		l1Blockchain = NewL1Blockchain(config.GenesisJSON)
	} else {
		nodecommon.LogWithID(obscurocommon.ShortAddress(config.HostID), "validateBlocks is set to false. L1 blocks will not be validated.")
	}

	var attestationProvider AttestationProvider
	if config.WillAttest {
		attestationProvider = &EgoAttestationProvider{}
	} else {
		nodecommon.LogWithID(nodeShortID, "WARNING - Attestation is not enabled, enclave will not create a verified attestation report.")
		attestationProvider = &DummyAttestationProvider{}
	}

	nodecommon.LogWithID(nodeShortID, "Generating public key")
	privKey := generateKeyPair()
	serializedPubKey := crypto.CompressPubkey(&privKey.PublicKey)
	nodecommon.LogWithID(nodeShortID, "Generated public key %s", common.Bytes2Hex(serializedPubKey))

	return &enclaveImpl{
		config:                config,
		nodeShortID:           nodeShortID,
		storage:               storage,
		blockResolver:         storage,
		mempool:               mempool.New(),
		statsCollector:        collector,
		l1Blockchain:          l1Blockchain,
		rpcEncryptionManager:  rpcencryptionmanager.NewRPCEncryptionManager(config.ViewingKeysEnabled, ecies.ImportECDSA(privKey)),
		txCh:                  make(chan *nodecommon.L2Tx),
		roundWinnerCh:         make(chan *obscurocore.Rollup),
		exitCh:                make(chan bool),
		speculativeWorkInCh:   make(chan bool),
		speculativeWorkOutCh:  make(chan speculativeWork),
		mgmtContractLib:       mgmtContractLib,
		erc20ContractLib:      erc20ContractLib,
		attestationProvider:   attestationProvider,
		privateKey:            privKey,
		publicKeySerialized:   serializedPubKey,
		transactionBlobCrypto: obscurocore.NewTransactionBlobCryptoImpl(),
	}
}

func (e *enclaveImpl) IsReady() error {
	return nil // The enclave is local so it is always ready
}

func (e *enclaveImpl) StopClient() error {
	return nil // The enclave is local so there is no client to stop
}

func (e *enclaveImpl) Start(block types.Block) {
	if e.config.SpeculativeExecution {
		// start the speculative rollup execution loop on its own go routine
		go e.start(block)
	}
}

func (e *enclaveImpl) start(block types.Block) {
	env := processingEnvironment{processedTxsMap: make(map[common.Hash]*nodecommon.L2Tx)}
	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.storage.FetchBlockState(block.Hash())
	if f {
		env.headRollup, _ = e.storage.FetchRollup(blockState.HeadRollup)
		if env.headRollup != nil {
			env.state = e.storage.CreateStateDB(env.headRollup.Hash())
		}
	}

	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:
			hash := winnerRollup.Hash()
			env.header = obscurocore.NewHeader(&hash, winnerRollup.Header.Number.Uint64()+1, e.config.HostID)
			env.headRollup = winnerRollup
			env.state = e.storage.CreateStateDB(winnerRollup.Hash())
			log.Trace(fmt.Sprintf(">   Agg%d: Create new speculative env  r_%d(%d).",
				e.nodeShortID,
				obscurocommon.ShortHash(winnerRollup.Header.Hash()),
				winnerRollup.Header.Number,
			))

			// determine the transactions that were not yet included
			env.processedTxs = currentTxs(winnerRollup, e.mempool.FetchMempoolTxs(), e.storage)
			env.processedTxsMap = makeMap(env.processedTxs)

			// calculate the State after executing them
			evm.ExecuteTransactions(env.processedTxs, env.state, env.headRollup.Header, e.storage, e.config.ObscuroChainID, 0)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if env.headRollup != nil {
				_, found := env.processedTxsMap[tx.Hash()]
				if !found {
					env.processedTxsMap[tx.Hash()] = tx
					env.processedTxs = append(env.processedTxs, tx)
					evm.ExecuteTransactions([]*nodecommon.L2Tx{tx}, env.state, env.header, e.storage, e.config.ObscuroChainID, 0)
				}
			}

		case <-e.speculativeWorkInCh:
			if env.header == nil {
				e.speculativeWorkOutCh <- speculativeWork{found: false}
			} else {
				b := make([]*nodecommon.L2Tx, 0, len(env.processedTxs))
				b = append(b, env.processedTxs...)
				e.speculativeWorkOutCh <- speculativeWork{
					found: true,
					r:     env.headRollup,
					s:     env.state,
					h:     env.header,
					txs:   b,
				}
			}

		case <-e.exitCh:
			return
		}
	}
}

func (e *enclaveImpl) ProduceGenesis(blkHash common.Hash) nodecommon.BlockSubmissionResponse {
	rolGenesis := obscurocore.NewRollup(blkHash, nil, obscurocommon.L2GenesisHeight, common.HexToAddress("0x0"), []*nodecommon.L2Tx{}, []nodecommon.Withdrawal{}, obscurocommon.GenerateNonce(), common.BigToHash(big.NewInt(0)))
	b, f := e.storage.FetchBlock(blkHash)
	if !f {
		log.Panic("Could not find the block used as proof for the genesis rollup.")
	}
	return nodecommon.BlockSubmissionResponse{
		ProducedRollup: rolGenesis.ToExtRollup(e.transactionBlobCrypto),
		BlockHeader:    b.Header(),
		IngestedBlock:  true,
	}
}

// IngestBlocks is used to update the enclave with the full history of the L1 chain to date.
func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) []nodecommon.BlockSubmissionResponse {
	result := make([]nodecommon.BlockSubmissionResponse, len(blocks))
	for i, block := range blocks {
		// We ignore a failure on the genesis block, since insertion of the genesis also produces a failure in Geth
		// (at least with Clique, where it fails with a `vote nonce not 0x00..0 or 0xff..f`).
		if ingestionFailedResponse := e.insertBlockIntoL1Chain(block); !e.isGenesisBlock(block) && ingestionFailedResponse != nil {
			result[i] = *ingestionFailedResponse
			return result // We return early, as all descendant blocks will also fail verification.
		}

		e.storage.StoreBlock(block)
		bs := updateState(block, e.blockResolver, e.mgmtContractLib, e.erc20ContractLib, e.storage, e.storage, e.nodeShortID, e.config.ObscuroChainID, e.transactionBlobCrypto)
		if bs == nil {
			result[i] = e.noBlockStateBlockSubmissionResponse(block)
		} else {
			var rollup nodecommon.ExtRollup
			if bs.FoundNewRollup {
				hr, f := e.storage.FetchRollup(bs.HeadRollup)
				if !f {
					log.Panic(msgNoRollup)
				}

				rollup = hr.ToExtRollup(e.transactionBlobCrypto)
			}
			result[i] = e.blockStateBlockSubmissionResponse(bs, rollup)
		}
	}

	return result
}

// SubmitBlock is used to update the enclave with an additional L1 block.
func (e *enclaveImpl) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	e.blockProcessingMutex.Lock()
	defer e.blockProcessingMutex.Unlock()

	// The genesis block should always be ingested, not submitted, so we ignore it if it's passed in here.
	if e.isGenesisBlock(&block) {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block was genesis block."}
	}

	_, foundBlock := e.storage.FetchBlock(block.Hash())
	if foundBlock {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block already ingested."}
	}

	if ingestionFailedResponse := e.insertBlockIntoL1Chain(&block); ingestionFailedResponse != nil {
		return *ingestionFailedResponse
	}

	_, f := e.storage.FetchBlock(block.Header().ParentHash)
	if !f && block.NumberU64() > obscurocommon.L1GenesisHeight {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: "Block parent not stored."}
	}

	// Only store the block if the parent is available.
	stored := e.storage.StoreBlock(&block)
	if !stored {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false}
	}

	nodecommon.LogWithID(e.nodeShortID, "Update state: %d", obscurocommon.ShortHash(block.Hash()))
	blockState := updateState(&block, e.blockResolver, e.mgmtContractLib, e.erc20ContractLib, e.storage, e.storage, e.nodeShortID, e.config.ObscuroChainID, e.transactionBlobCrypto)
	if blockState == nil {
		return e.noBlockStateBlockSubmissionResponse(&block)
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	hr, f := e.storage.FetchRollup(blockState.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}
	e.mempool.RemoveMempoolTxs(historicTxs(hr, e.storage))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.storage.StoreRollup(r)

	nodecommon.LogWithID(e.nodeShortID, "Processed block: b_%d(%d)", obscurocommon.ShortHash(block.Hash()), block.NumberU64())

	return e.blockStateBlockSubmissionResponse(blockState, r.ToExtRollup(e.transactionBlobCrypto))
}

func (e *enclaveImpl) SubmitRollup(rollup nodecommon.ExtRollup) {
	r := obscurocore.Rollup{
		Header:       rollup.Header,
		Transactions: e.transactionBlobCrypto.Decrypt(rollup.EncryptedTxBlob),
	}

	// only store if the parent exists
	_, found := e.storage.FetchRollup(r.Header.ParentHash)
	if found {
		e.storage.StoreRollup(&r)
	} else {
		nodecommon.LogWithID(e.nodeShortID, "Received rollup with no parent: r_%d", obscurocommon.ShortHash(r.Hash()))
	}
}

func (e *enclaveImpl) SubmitTx(tx nodecommon.EncryptedTx) error {
	decryptedTx := obscurocore.DecryptTx(tx)
	err := verifySignature(e.config.ObscuroChainID, decryptedTx)
	if err != nil {
		return err
	}
	e.mempool.AddMempoolTx(decryptedTx)
	if e.config.SpeculativeExecution {
		e.txCh <- decryptedTx
	}
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(chainID int64, decryptedTx *nodecommon.L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(chainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	head, found := e.storage.FetchRollup(parent)
	if !found {
		return nodecommon.ExtRollup{}, false, fmt.Errorf("rollup not found: r_%s", parent)
	}

	nodecommon.LogWithID(e.nodeShortID, "Round winner height: %d", head.Header.Number)
	rollupsReceivedFromPeers := e.storage.FetchRollups(head.Header.Number.Uint64() + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*obscurocore.Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.storage.ParentRollup(rol)
		if p == nil {
			nodecommon.LogWithID(e.nodeShortID, "Received rollup from peer but don't have parent rollup - discarding...")
			continue
		}
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.storage.CreateStateDB(head.Hash())
	// determine the winner of the round
	winnerRollup, _ := e.findRoundWinner(usefulRollups, head, parentState, e.blockResolver, e.storage)
	if e.config.SpeculativeExecution {
		go e.notifySpeculative(winnerRollup)
	}

	// we are the winner
	if winnerRollup.Header.Agg == e.config.HostID {
		v := e.blockResolver.Proof(winnerRollup)
		w := e.storage.ParentRollup(winnerRollup)
		nodecommon.LogWithID(e.nodeShortID, "Publish rollup=r_%d(%d)[r_%d]{proof=b_%d(%d)}. Num Txs: %d. Txs: %v.  Root=%v. ",
			obscurocommon.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Number,
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			v.NumberU64(),
			len(winnerRollup.Transactions),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.Root,
		)
		return winnerRollup.ToExtRollup(e.transactionBlobCrypto), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *obscurocore.Rollup) {
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) ExecuteOffChainTransaction(encryptedParams nodecommon.EncryptedParams) (nodecommon.EncryptedResponse, error) {
	paramBytes, err := e.rpcEncryptionManager.DecryptWithEnclaveKey(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in Call request. Cause: %w", err)
	}

	contractAddress, from, data, err := extractCallParams(paramBytes)
	if err != nil {
		return nil, err
	}

	hs := e.storage.FetchHeadState()
	if hs == nil {
		panic("Not initialised")
	}
	// todo - get the parent
	r, f := e.storage.FetchRollup(hs.HeadRollup)
	if !f {
		panic("not found")
	}
	s := e.storage.CreateStateDB(hs.HeadRollup)
	result, err := evm.ExecuteOffChainCall(from, contractAddress, data, s, r.Header, e.storage, e.config.ObscuroChainID)
	if err != nil {
		return nil, err
	}
	if result.Failed() {
		log.Info("Failed to execute contract %s: %s\n", contractAddress.Hex(), result.Err)
		return nil, result.Err
	}

	encryptedResult, err := e.rpcEncryptionManager.EncryptWithViewingKey(from, result.ReturnData)
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_call request. Cause: %w", err)
	}

	return encryptedResult, nil
}

func (e *enclaveImpl) Nonce(address common.Address) uint64 {
	// todo user encryption
	hs := e.storage.FetchHeadState()
	if hs == nil {
		return 0
	}
	s := e.storage.CreateStateDB(hs.HeadRollup)
	return s.GetNonce(address)
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs *obscurocore.BlockState) *obscurocore.Rollup {
	headRollup, f := e.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	// These variables will be used to create the new rollup
	var newRollupTxs obscurocore.L2Txs
	var newRollupState *state.StateDB
	var newRollupHeader *nodecommon.Header

	/*
			speculativeExecutionSucceeded := false
		   todo - reenable
			if e.speculativeExecutionEnabled {
				// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
				e.speculativeWorkInCh <- true
				speculativeRollup := <-e.speculativeWorkOutCh

				newRollupTxs = speculativeRollup.txs
				newRollupState = speculativeRollup.s
				newRollupHeader = speculativeRollup.h

				// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
				// or speculative execution is disabled
				speculativeExecutionSucceeded = speculativeRollup.found && (speculativeRollup.r.Hash() == bs.HeadRollup)

				if !speculativeExecutionSucceeded && speculativeRollup.r != nil {
					nodecommon.LogWithID(e.nodeShortID, "Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
						obscurocommon.ShortHash(speculativeRollup.r.Hash()),
						speculativeRollup.r.Header.Number,
						obscurocommon.ShortHash(bs.HeadRollup),
						headRollup.Header.Number)
					if e.statsCollector != nil {
						e.statsCollector.L2Recalc(e.nodeID)
					}
				}
			}
	*/

	successfulTransactions := make([]*nodecommon.L2Tx, 0)
	// if !speculativeExecutionSucceeded {
	// In case the speculative execution thread has not succeeded in producing a valid rollup
	// we have to create a new one from the mempool transactions
	newRollupHeader = obscurocore.NewHeader(&bs.HeadRollup, headRollup.Header.Number.Uint64()+1, e.config.HostID)
	newRollupTxs = currentTxs(headRollup, e.mempool.FetchMempoolTxs(), e.storage)
	newRollupState = e.storage.CreateStateDB(bs.HeadRollup)
	txReceipts := evm.ExecuteTransactions(newRollupTxs, newRollupState, newRollupHeader, e.storage, e.config.ObscuroChainID, 0)
	txReceiptsMap := toReceiptMap(txReceipts)
	// todo - only transactions that fail because of the nonce should be excluded
	for i, tx := range newRollupTxs {
		_, f := txReceiptsMap[tx.Hash()]
		if f {
			successfulTransactions = append(successfulTransactions, newRollupTxs[i])
		} else {
			log.Info(">   Agg%d: Excluding transaction %d", obscurocommon.ShortAddress(e.config.HostID), obscurocommon.ShortHash(tx.Hash()))
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := e.blockResolver.Proof(headRollup)
	depositTxs := extractDeposits(proof, b, e.blockResolver, e.erc20ContractLib, newRollupState, e.config.ObscuroChainID)
	depositReceipts := evm.ExecuteTransactions(depositTxs, newRollupState, newRollupHeader, e.storage, e.config.ObscuroChainID, len(newRollupTxs))
	depositReceiptsMap := toReceiptMap(depositReceipts)
	for _, tx := range depositTxs {
		if depositReceiptsMap[tx.Hash()] == nil {
			panic("Should not happen")
		}
	}

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	rootHash, err := newRollupState.Commit(true)
	if err != nil {
		return nil
	}
	r := obscurocore.NewRollupFromHeader(newRollupHeader, b.Hash(), successfulTransactions, obscurocommon.GenerateNonce(), rootHash)

	// Postprocessing - withdrawals
	r.Header.Withdrawals = e.rollupPostProcessingWithdrawals(&r, newRollupState, txReceipts)
	receipts := getReceipts(txReceipts, depositReceipts)

	if len(receipts) == 0 {
		r.Header.ReceiptHash = types.EmptyRootHash
	} else {
		r.Header.ReceiptHash = types.DeriveSha(receipts, trie.NewStackTrie(nil))
		r.Header.Bloom = types.CreateBloom(receipts)
	}

	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	// todo add some sort of cache
	hs := e.storage.FetchHeadState()
	if hs == nil {
		panic("should not happen")
	}
	rollup, found := e.storage.FetchRollup(hs.HeadRollup)
	if !found {
		log.Panic("could not fetch block's head rollup")
	}

	for {
		txs := rollup.Transactions
		for _, tx := range txs {
			if tx.Hash() == txHash {
				return tx
			}
		}
		rollup = e.storage.ParentRollup(rollup)
		if rollup == nil || rollup.Header.Number.Uint64() == obscurocommon.L2GenesisHeight {
			return nil
		}
	}
}

func (e *enclaveImpl) GetTransactionReceipt(common.Hash) *types.Receipt {
	return nil // todo - joel - implement
}

func (e *enclaveImpl) GetRollup(rollupHash obscurocommon.L2RootHash) *nodecommon.ExtRollup {
	rollup, found := e.storage.FetchRollup(rollupHash)
	if found {
		extRollup := rollup.ToExtRollup(e.transactionBlobCrypto)
		return &extRollup
	}
	return nil
}

func (e *enclaveImpl) Stop() error {
	if e.config.SpeculativeExecution {
		e.exitCh <- true
	}
	return nil
}

func (e *enclaveImpl) Attestation() *obscurocommon.AttestationReport {
	if e.publicKeySerialized == nil {
		panic("public key not initialized, we can't produce the attestation report")
	}
	report, err := e.attestationProvider.GetReport(e.publicKeySerialized, e.config.HostID)
	if err != nil {
		panic("Failed to produce remote report.")
	}
	return report
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	secret := make([]byte, 32)
	n, err := rand.Read(secret)
	if n != 32 || err != nil {
		log.Panic("could not generate secret. Cause: %s", err)
	}
	e.storage.StoreSecret(secret)
	encSec, err := e.encryptSecret(e.publicKeySerialized, secret)
	if err != nil {
		log.Panic("failed to encrypt secret. Cause: %s", err)
	}
	return encSec
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(s obscurocommon.EncryptedSharedEnclaveSecret) error {
	secret, err := e.decryptSecret(s)
	if err != nil {
		return err
	}
	e.storage.StoreSecret(secret)
	log.Trace(">   Agg%d: Secret decrypted and stored. Secret: %v", e.nodeShortID, secret)
	return nil
}

// ShareSecret verifies the request and if it trusts the report and the public key it will return the secret encrypted with that public key.
func (e *enclaveImpl) ShareSecret(att *obscurocommon.AttestationReport) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	// First we verify the attestation report has come from a valid obscuro enclave running in a verified TEE.
	data, err := e.attestationProvider.VerifyReport(att)
	if err != nil {
		return nil, err
	}
	// Then we verify the public key provided has come from the same enclave as that attestation report
	if err = verifyIdentity(data, att); err != nil {
		return nil, err
	}
	nodecommon.LogWithID(e.nodeShortID, "Successfully verified attestation and identity. Owner: %s", att.Owner)

	secret := e.storage.FetchSecret()
	if secret == nil {
		return nil, errors.New("secret was nil, no secret to share - this shouldn't happen")
	}
	return e.encryptSecret(att.PubKey, secret)
}

func (e *enclaveImpl) AddViewingKey(encryptedViewingKeyBytes []byte, signature []byte) error {
	viewingKeyBytes, err := ecies.ImportECDSA(e.privateKey).Decrypt(encryptedViewingKeyBytes, nil, nil)
	if err != nil {
		return fmt.Errorf("could not decrypt viewing key when adding it to enclave. Cause: %w", err)
	}
	return e.rpcEncryptionManager.AddViewingKey(viewingKeyBytes, signature)
}

func (e *enclaveImpl) GetBalance(encryptedParams nodecommon.EncryptedParams) (nodecommon.EncryptedResponse, error) {
	paramBytes, err := e.rpcEncryptionManager.DecryptWithEnclaveKey(encryptedParams)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt params in GetBalance request. Cause: %w", err)
	}

	var paramsJSONMap []string
	err = json.Unmarshal(paramBytes, &paramsJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON params in GetBalance request. Cause: %w", err)
	}
	address := common.HexToAddress(paramsJSONMap[0]) // The first argument is the address, the second the block.

	// TODO - Calculate balance correctly, rather than returning this dummy value.
	balance := DummyBalance // The Ethereum API is to return the balance in hex.

	encryptedBalance, err := e.rpcEncryptionManager.EncryptWithViewingKey(address, []byte(balance))
	if err != nil {
		return nil, fmt.Errorf("enclave could not respond securely to eth_getBalance request. Cause: %w", err)
	}

	return encryptedBalance, nil
}

func verifyIdentity(data []byte, att *obscurocommon.AttestationReport) error {
	expectedIDHash := getIDHash(att.Owner, att.PubKey)
	// we trim the actual data because data extracted from the verified attestation is always 64 bytes long (padded with zeroes at the end)
	if !bytes.Equal(expectedIDHash, data[:len(expectedIDHash)]) {
		return fmt.Errorf("failed to verify hash for attestation report with owner: %s", att.Owner)
	}
	return nil
}

func (e *enclaveImpl) IsInitialised() bool {
	return e.storage.FetchSecret() != nil
}

func (e *enclaveImpl) isGenesisBlock(block *types.Block) bool {
	return e.l1Blockchain != nil && block.Hash() == e.l1Blockchain.Genesis().Hash()
}

// Inserts the block into the L1 chain if it exists and the block is not the genesis block. Returns a non-nil
// BlockSubmissionResponse if the insertion failed.
func (e *enclaveImpl) insertBlockIntoL1Chain(block *types.Block) *nodecommon.BlockSubmissionResponse {
	if e.l1Blockchain != nil {
		_, err := e.l1Blockchain.InsertChain(types.Blocks{block})
		if err != nil {
			causeMsg := fmt.Sprintf("Block was invalid: %v", err)
			return &nodecommon.BlockSubmissionResponse{IngestedBlock: false, BlockNotIngestedCause: causeMsg}
		}
	}
	return nil
}

func (e *enclaveImpl) noBlockStateBlockSubmissionResponse(block *types.Block) nodecommon.BlockSubmissionResponse {
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:   block.Header(),
		IngestedBlock: true,
		FoundNewHead:  false,
	}
}

func (e *enclaveImpl) blockStateBlockSubmissionResponse(bs *obscurocore.BlockState, rollup nodecommon.ExtRollup) nodecommon.BlockSubmissionResponse {
	headRollup, f := e.storage.FetchRollup(bs.HeadRollup)
	if !f {
		log.Panic(msgNoRollup)
	}

	headBlock, f := e.storage.FetchBlock(bs.Block)
	if !f {
		log.Panic("could not fetch block")
	}

	var head *nodecommon.Header
	if bs.FoundNewRollup {
		head = headRollup.Header
	}
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:    headBlock.Header(),
		ProducedRollup: rollup,
		IngestedBlock:  true,
		FoundNewHead:   bs.FoundNewRollup,
		RollupHead:     head,
	}
}

func generateKeyPair() *ecdsa.PrivateKey {
	// todo: This should be generated deterministically based on some enclave attributes if possible
	key, err := crypto.HexToECDSA(enclavePrivateKeyHex)
	if err != nil {
		panic("Failed to create enclave private key")
	}
	return key
}

// Todo - implement with better crypto
func (e *enclaveImpl) decryptSecret(secret obscurocommon.EncryptedSharedEnclaveSecret) ([]byte, error) {
	if e.privateKey == nil {
		return nil, errors.New("private key not found - shouldn't happen")
	}
	return decryptWithPrivateKey(secret, e.privateKey)
}

// Todo - implement with better crypto
func (e *enclaveImpl) encryptSecret(pubKeyEncoded []byte, secret obscurocore.SharedEnclaveSecret) (obscurocommon.EncryptedSharedEnclaveSecret, error) {
	nodecommon.LogWithID(e.nodeShortID, "Encrypting secret with public key %s", common.Bytes2Hex(pubKeyEncoded))
	key, err := crypto.DecompressPubkey(pubKeyEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key %w", err)
	}

	encKey, err := encryptWithPublicKey(secret, key)
	if err != nil {
		nodecommon.LogWithID(e.nodeShortID, "Failed to encrypt key, err: %s\nsecret: %v\npubkey: %v\nencKey:%v", err, secret, pubKeyEncoded, encKey)
	}
	return encKey, err
}

// Extracts and validates the relevant parameters in a Call request.
func extractCallParams(decryptedParams []byte) (common.Address, common.Address, []byte, error) {
	var paramsJSONMap []interface{}
	err := json.Unmarshal(decryptedParams, &paramsJSONMap)
	if err != nil {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("could not parse JSON params in Call request. Cause: %w", err)
	}

	txArgs := paramsJSONMap[0] // The first argument is the transaction arguments, the second the block, the third the state overrides.
	contractAddressString, ok := txArgs.(map[string]interface{})[CallFieldTo].(string)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("to field in Call request params was not of expected type string")
	}
	fromString, ok := txArgs.(map[string]interface{})[CallFieldFrom].(string)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("from field in Call request params was not of expected type string")
	}
	dataString, ok := txArgs.(map[string]interface{})[CallFieldData].(string)
	if !ok {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("data field in Call request params was not of expected type string")
	}

	contractAddress := common.HexToAddress(contractAddressString)
	from := common.HexToAddress(fromString)
	data, err := hexutil.Decode(dataString)
	if err != nil {
		return common.Address{}, common.Address{}, nil, fmt.Errorf("could not decode data in Call request. Cause: %w", err)
	}
	return contractAddress, from, data, nil
}

// internal structure to pass information.
type speculativeWork struct {
	found bool
	r     *obscurocore.Rollup
	s     *state.StateDB
	h     *nodecommon.Header
	txs   []*nodecommon.L2Tx
}

// internal structure used for the speculative execution.
type processingEnvironment struct {
	headRollup      *obscurocore.Rollup              // the current head rollup, which will be the parent of the new rollup
	header          *nodecommon.Header               // the header of the new rollup
	processedTxs    []*nodecommon.L2Tx               // txs that were already processed
	processedTxsMap map[common.Hash]*nodecommon.L2Tx // structure used to prevent duplicates
	state           *state.StateDB                   // the state as calculated from the previous rollup and the processed transactions
}

// encryptWithPublicKey encrypts data with public key
func encryptWithPublicKey(msg []byte, pub *ecdsa.PublicKey) ([]byte, error) {
	ciphertext, err := ecies.Encrypt(rand.Reader, ecies.ImportECDSAPublic(pub), msg, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with public key. %w", err)
	}
	return ciphertext, nil
}

// decryptWithPrivateKey decrypts data with private key
func decryptWithPrivateKey(ciphertext []byte, priv *ecdsa.PrivateKey) ([]byte, error) {
	plaintext, err := ecies.ImportECDSA(priv).Decrypt(ciphertext, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt with private key. %w", err)
	}
	return plaintext, nil
}

// getDB creates an appropriate ethdb.Database instance based on your config
func getDB(nodeID uint64, cfg config.EnclaveConfig) (ethdb.Database, error) {
	if cfg.UseInMemoryDB {
		nodecommon.LogWithID(nodeID, "UseInMemoryDB flag is true, data will not be persisted. Creating in-memory database...")
		return getInMemDB()
	}

	if !cfg.WillAttest {
		// persistent but not secure in an enclave, we'll connect to a throwaway sqlite DB and test out persistence/sql implementations
		nodecommon.LogWithID(nodeID, "Attestation is disabled, using a basic sqlite DB for persistence")
		// todo: for now we pass in an empty dbPath which will provide a throwaway temp file,
		// 		when we want to test persistence after node restart we should change this path to be config driven
		return sql.CreateTemporarySQLiteDB("")
	}

	// persistent and with attestation means connecting to edgeless DB in a trusted enclave from a secure enclave
	panic("Haven't implemented edgeless DB enclave connection yet")
}

func getInMemDB() (ethdb.Database, error) {
	return rawdb.NewMemoryDatabase(), nil
}
