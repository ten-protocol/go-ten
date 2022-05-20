package enclave

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"

	"github.com/ethereum/go-ethereum/core/state"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/mempool"

	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/core"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type StatsCollector interface {
	// L2Recalc registers when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	RollupWithMoreRecentProof()
}

type enclaveImpl struct {
	nodeID         common.Address
	nodeShortID    uint64
	mining         bool
	storage        db.Storage
	blockResolver  db.BlockResolver
	mempool        mempool.Manager
	statsCollector StatsCollector
	l1Blockchain   *core.BlockChain

	txCh                 chan nodecommon.L2Tx
	roundWinnerCh        chan *obscurocore.Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork

	txHandler mgmtcontractlib.TxHandler

	// Toggles the speculative execution background process
	speculativeExecutionEnabled bool
}

func (e *enclaveImpl) IsReady() error {
	return nil // The enclave is local so it is always ready
}

func (e *enclaveImpl) StopClient() {
	// The enclave is local so there is no client to stop
}

func (e *enclaveImpl) Start(block types.Block) {
	if e.speculativeExecutionEnabled {
		// start the speculative rollup execution loop on its own go routine
		go e.start(block)
	}
}

func (e *enclaveImpl) start(block types.Block) {
	env := processingEnvironment{processedTxsMap: make(map[common.Hash]nodecommon.L2Tx)}
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
			env.header = obscurocore.NewHeader(&hash, winnerRollup.Header.Number+1, e.nodeID)
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
			evm.ExecuteTransactions(env.processedTxs, env.state, env.headRollup.Header, e.storage)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if env.headRollup != nil {
				_, found := env.processedTxsMap[tx.Hash()]
				if !found {
					env.processedTxsMap[tx.Hash()] = tx
					env.processedTxs = append(env.processedTxs, tx)
					evm.ExecuteTransactions([]nodecommon.L2Tx{tx}, env.state, env.header, e.storage)
				}
			}

		case <-e.speculativeWorkInCh:
			if env.header == nil {
				e.speculativeWorkOutCh <- speculativeWork{found: false}
			} else {
				b := make([]nodecommon.L2Tx, 0, len(env.processedTxs))
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
	rolGenesis := obscurocore.NewRollup(blkHash, nil, obscurocommon.L2GenesisHeight, common.HexToAddress("0x0"), []nodecommon.L2Tx{}, []nodecommon.Withdrawal{}, obscurocommon.GenerateNonce(), common.BigToHash(big.NewInt(0)))
	b, f := e.storage.FetchBlock(blkHash)
	if !f {
		panic("Could not find the block used as proof for the genesis rollup.")
	}
	return nodecommon.BlockSubmissionResponse{
		ProducedRollup: rolGenesis.ToExtRollup(),
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
		bs := updateState(block, e.blockResolver, e.txHandler, e.storage, e.storage, e.nodeShortID)
		if bs == nil {
			result[i] = e.noBlockStateBlockSubmissionResponse(block)
		} else {
			var rollup nodecommon.ExtRollup
			if bs.FoundNewRollup {
				hr, f := e.storage.FetchRollup(bs.HeadRollup)
				if !f {
					panic("Should not happen")
				}

				rollup = hr.ToExtRollup()
			}
			result[i] = e.blockStateBlockSubmissionResponse(bs, rollup)
		}
	}

	return result
}

// SubmitBlock is used to update the enclave with an additional block.
func (e *enclaveImpl) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
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

	blockState := updateState(&block, e.blockResolver, e.txHandler, e.storage, e.storage, e.nodeShortID)
	if blockState == nil {
		return e.noBlockStateBlockSubmissionResponse(&block)
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	hr, f := e.storage.FetchRollup(blockState.HeadRollup)
	if !f {
		panic("Failed to fetch rollup. Should not happen")
	}
	e.mempool.RemoveMempoolTxs(historicTxs(hr, e.storage))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.storage.StoreRollup(r)

	nodecommon.LogWithID(e.nodeShortID, "Processed block: b_%d(%d)", obscurocommon.ShortHash(block.Hash()), block.NumberU64())

	return e.blockStateBlockSubmissionResponse(blockState, r.ToExtRollup())
}

func (e *enclaveImpl) SubmitRollup(rollup nodecommon.ExtRollup) {
	r := obscurocore.Rollup{
		Header:       rollup.Header,
		Transactions: obscurocore.DecryptTransactions(rollup.Txs),
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
	err := verifySignature(&decryptedTx)
	if err != nil {
		return err
	}
	e.mempool.AddMempoolTx(decryptedTx)
	if e.speculativeExecutionEnabled {
		e.txCh <- decryptedTx
	}
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(decryptedTx *nodecommon.L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(evm.ChainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	head, found := e.storage.FetchRollup(parent)
	if !found {
		return nodecommon.ExtRollup{}, false, fmt.Errorf("rollup not found: r_%s", parent)
	}

	nodecommon.LogWithID(e.nodeShortID, "Round winner height: %d", head.Header.Number)
	rollupsReceivedFromPeers := e.storage.FetchRollups(head.Header.Number + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*obscurocore.Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.storage.ParentRollup(rol)
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.storage.CreateStateDB(head.Hash())
	// determine the winner of the round
	winnerRollup, _ := e.findRoundWinner(usefulRollups, head, parentState, e.blockResolver, e.storage)
	if e.speculativeExecutionEnabled {
		go e.notifySpeculative(winnerRollup)
	}

	// we are the winner
	if winnerRollup.Header.Agg == e.nodeID {
		v := e.blockResolver.Proof(winnerRollup)
		w := e.storage.ParentRollup(winnerRollup)
		nodecommon.LogWithID(e.nodeShortID, "Publish rollup=r_%d(%d)[r_%d]{proof=b_%d(%d)}. Num Txs: %d. Txs: %v.  State=%v. ",
			obscurocommon.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Number,
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			v.NumberU64(),
			len(winnerRollup.Transactions),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.State,
		)
		return winnerRollup.ToExtRollup(), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *obscurocore.Rollup) {
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address common.Address) uint64 {
	// todo user encryption
	s := e.storage.CreateStateDB(e.storage.FetchHeadState().HeadRollup)
	r, f := e.storage.FetchRollup(e.storage.FetchHeadState().HeadRollup)
	if !f {
		panic("not found")
	}
	return evm.BalanceOfErc20(s, address, r.Header, e.storage)
}

func (e *enclaveImpl) Nonce(address common.Address) uint64 {
	// todo user encryption
	s := e.storage.CreateStateDB(e.storage.FetchHeadState().HeadRollup)
	return s.GetNonce(address)
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs *obscurocore.BlockState) *obscurocore.Rollup {
	headRollup, f := e.storage.FetchRollup(bs.HeadRollup)
	if !f {
		panic("Should not happen")
	}

	// These variables will be used to create the new rollup
	var newRollupTxs []nodecommon.L2Tx
	var newRollupState *state.StateDB
	var newRollupHeader *nodecommon.Header

	speculativeExecutionSucceeded := false

	if e.speculativeExecutionEnabled {
		// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
		e.speculativeWorkInCh <- true
		speculativeRollup := <-e.speculativeWorkOutCh

		// newRollupTxs = speculativeRollup.txs
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

	successfulTransactions := make([]nodecommon.L2Tx, 0)
	receipts := map[common.Hash]*types.Receipt{}
	if !speculativeExecutionSucceeded {
		// In case the speculative execution thread has not succeeded in producing a valid rollup
		// we have to create a new one from the mempool transactions
		newRollupHeader = obscurocore.NewHeader(&bs.HeadRollup, headRollup.Header.Number+1, e.nodeID)
		newRollupTxs = currentTxs(headRollup, e.mempool.FetchMempoolTxs(), e.storage)

		newRollupState = e.storage.CreateStateDB(bs.HeadRollup)
		receipts = evm.ExecuteTransactions(newRollupTxs, newRollupState, newRollupHeader, e.storage)
		// todo - only transactions that fail because of the nonce should be excluded
		for _, tx := range newRollupTxs {
			_, f := receipts[tx.Hash()]
			if f {
				successfulTransactions = append(successfulTransactions, tx)
			} else {
				fmt.Printf("Excluding transaction %d\n", obscurocommon.ShortHash(tx.Hash()))
			}
		}
	}

	// always process deposits last, either on top of the rollup produced speculatively or the newly created rollup
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := e.blockResolver.Proof(headRollup)
	depositTxs := extractDeposits(proof, b, e.blockResolver, e.txHandler, newRollupState)
	depositReceipts := evm.ExecuteTransactions(depositTxs, newRollupState, newRollupHeader, e.storage)
	for _, tx := range depositTxs {
		if depositReceipts[tx.Hash()] == nil {
			panic("Should not happen")
		}
	}

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	rootHash, err := newRollupState.Commit(true)
	if err != nil {
		return nil
	}
	// dump := newRollupState.Dump(&state.DumpConfig{})
	// log.Info(fmt.Sprintf(">   Agg%d: State:%s", obscurocommon.ShortAddress(e.nodeID), dump))
	r := obscurocore.NewRollupFromHeader(newRollupHeader, b.Hash(), successfulTransactions, obscurocommon.GenerateNonce(), rootHash)

	// Postprocessing - withdrawals
	r.Header.Withdrawals = rollupPostProcessingWithdrawals(&r, newRollupState, receipts)

	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	// todo add some sort of cache
	rollup, found := e.storage.FetchRollup(e.storage.FetchHeadState().HeadRollup)
	if !found {
		panic("should not happen")
	}

	for {
		txs := rollup.Transactions
		for _, tx := range txs {
			if tx.Hash() == txHash {
				return &tx
			}
		}
		rollup = e.storage.ParentRollup(rollup)
		if rollup == nil || rollup.Header.Number == obscurocommon.L2GenesisHeight {
			return nil
		}
	}
}

func (e *enclaveImpl) Stop() error {
	if e.speculativeExecutionEnabled {
		e.exitCh <- true
	}
	return nil
}

func (e *enclaveImpl) Attestation() obscurocommon.AttestationReport {
	// Todo
	return obscurocommon.AttestationReport{Owner: e.nodeID}
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret {
	secret := make([]byte, 32)
	n, err := rand.Read(secret)
	if n != 32 || err != nil {
		panic(fmt.Sprintf("Could not generate secret: %s", err))
	}
	e.storage.StoreSecret(secret)
	return encryptSecret(secret)
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	e.storage.StoreSecret(decryptSecret(secret))
}

func (e *enclaveImpl) FetchSecret(obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	return encryptSecret(e.storage.FetchSecret())
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
		panic("Should not happen")
	}

	headBlock, f := e.storage.FetchBlock(bs.Block)
	if !f {
		panic("Should not happen")
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

// Todo - implement with crypto
func decryptSecret(secret obscurocommon.EncryptedSharedEnclaveSecret) obscurocore.SharedEnclaveSecret {
	return obscurocore.SharedEnclaveSecret(secret)
}

// Todo - implement with crypto
func encryptSecret(secret obscurocore.SharedEnclaveSecret) obscurocommon.EncryptedSharedEnclaveSecret {
	return obscurocommon.EncryptedSharedEnclaveSecret(secret)
}

// internal structure to pass information.
type speculativeWork struct {
	found bool
	r     *obscurocore.Rollup
	s     *state.StateDB
	h     *nodecommon.Header
	txs   []nodecommon.L2Tx
}

// internal structure used for the speculative execution.
type processingEnvironment struct {
	headRollup      *obscurocore.Rollup             // the current head rollup, which will be the parent of the new rollup
	header          *nodecommon.Header              // the header of the new rollup
	processedTxs    []nodecommon.L2Tx               // txs that were already processed
	processedTxsMap map[common.Hash]nodecommon.L2Tx // structure used to prevent duplicates
	state           *state.StateDB                  // the state as calculated from the previous rollup and the processed transactions
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(nodeID common.Address, mining bool, txHandler mgmtcontractlib.TxHandler, validateBlocks bool, genesisJSON []byte, collector StatsCollector) nodecommon.Enclave {
	backingDB := db.NewInMemoryDB()
	nodeShortID := obscurocommon.ShortAddress(nodeID)
	storage := db.NewStorage(backingDB, nodeShortID)

	var l1Blockchain *core.BlockChain
	if validateBlocks {
		if genesisJSON == nil {
			panic("enclave was configured to validate blocks, but genesis JSON was nil")
		}
		l1Blockchain = NewL1Blockchain(genesisJSON)
	} else {
		nodecommon.LogWithID(obscurocommon.ShortAddress(nodeID), "validateBlocks is set to false. L1 blocks will not be validated.")
	}

	return &enclaveImpl{
		nodeID:                      nodeID,
		nodeShortID:                 nodeShortID,
		mining:                      mining,
		storage:                     storage,
		blockResolver:               storage,
		mempool:                     mempool.New(),
		statsCollector:              collector,
		l1Blockchain:                l1Blockchain,
		txCh:                        make(chan nodecommon.L2Tx),
		roundWinnerCh:               make(chan *obscurocore.Rollup),
		exitCh:                      make(chan bool),
		speculativeWorkInCh:         make(chan bool),
		speculativeWorkOutCh:        make(chan speculativeWork),
		txHandler:                   txHandler,
		speculativeExecutionEnabled: false, // TODO - reenable
	}
}
