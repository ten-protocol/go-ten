package enclave

import (
	"crypto/rand"
	"fmt"
	"math/big"

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

const ChainID = 777 // The unique ID for the Obscuro chain. Required for Geth signing.

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	RollupWithMoreRecentProof()
}

type enclaveImpl struct {
	node           common.Address
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
	txHandler            mgmtcontractlib.TxHandler
}

func (e *enclaveImpl) IsReady() error {
	return nil // The enclave is local so it is always ready
}

func (e *enclaveImpl) StopClient() {
	// The enclave is local so there is no client to stop
}

func (e *enclaveImpl) Start(block types.Block) {
	// start the speculative rollup execution loop on its own go routine
	go e.start(block)
}

func (e *enclaveImpl) start(block types.Block) {
	env := processingEnvironment{processedTxsMap: make(map[common.Hash]nodecommon.L2Tx)}
	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.storage.FetchBlockState(block.Hash())
	if f {
		env.headRollup = blockState.Head
		if env.headRollup != nil {
			env.state = e.storage.CreateStateDB(env.headRollup.Hash())
		}
	}

	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:
			env.header = obscurocore.NewHeader(winnerRollup, winnerRollup.Header.Number+1, e.node)
			env.headRollup = winnerRollup
			env.state = e.storage.CreateStateDB(winnerRollup.Hash())

			// determine the transactions that were not yet included
			env.processedTxs = currentTxs(winnerRollup, e.mempool.FetchMempoolTxs(), e.storage)
			env.processedTxsMap = makeMap(env.processedTxs)

			// calculate the State after executing them
			executeTransactions(env.processedTxs, env.state, env.headRollup.Header)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if env.headRollup != nil {
				_, found := env.processedTxsMap[tx.Hash()]
				if !found {
					env.processedTxsMap[tx.Hash()] = tx
					env.processedTxs = append(env.processedTxs, tx)
					executeTx(env.state, tx)
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
					s:     env.state.Copy(),
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
		bs := updateState(block, e.blockResolver, e.txHandler, e.storage, e.storage)
		if bs == nil {
			result[i] = e.noBlockStateBlockSubmissionResponse(block)
		} else {
			var rollup nodecommon.ExtRollup
			if bs.FoundNewRollup {
				rollup = bs.Head.ToExtRollup()
			}
			result[i] = e.blockStateBlockSubmissionResponse(bs, rollup)
		}
	}

	return result
}

// SubmitBlock is used to update the enclave with an additional block.
func (e *enclaveImpl) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	// The genesis block will always be ingested, not submitted.
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

	blockState := updateState(&block, e.blockResolver, e.txHandler, e.storage, e.storage)
	if blockState == nil {
		return e.noBlockStateBlockSubmissionResponse(&block)
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	e.mempool.RemoveMempoolTxs(historicTxs(blockState.Head, e.storage))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.storage.StoreRollup(r)

	log.Log(fmt.Sprintf("Agg%d:> Processed block: b_%d", obscurocommon.ShortAddress(e.node), obscurocommon.ShortHash(block.Hash())))

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
		log.Log(fmt.Sprintf("Agg%d:> Received rollup with no parent: r_%d", obscurocommon.ShortAddress(e.node), obscurocommon.ShortHash(r.Hash())))
	}
}

func (e *enclaveImpl) SubmitTx(tx nodecommon.EncryptedTx) error {
	decryptedTx := obscurocore.DecryptTx(tx)
	err := verifySignature(&decryptedTx)
	if err != nil {
		return err
	}
	e.mempool.AddMempoolTx(decryptedTx)
	e.txCh <- decryptedTx
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(decryptedTx *nodecommon.L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(ChainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool, error) {
	head, found := e.storage.FetchRollup(parent)
	if !found {
		return nodecommon.ExtRollup{}, false, fmt.Errorf("rollup not found: r_%s", parent)
	}

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
	winnerRollup, s := e.findRoundWinner(usefulRollups, head, parentState, e.blockResolver, e.storage)
	s.Commit(winnerRollup.Hash())
	// e.storage.SetRollupState(winnerRollup.Hash(), s)
	go e.notifySpeculative(winnerRollup)

	// we are the winner
	if winnerRollup.Header.Agg == e.node {
		v := e.blockResolver.Proof(winnerRollup)
		w := e.storage.ParentRollup(winnerRollup)
		log.Log(fmt.Sprintf(">   Agg%d: publish rollup=r_%d(%d)[r_%d]{proof=b_%d}. Num Txs: %d. Txs: %v.  State=%v. ",
			obscurocommon.ShortAddress(e.node),
			obscurocommon.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Number,
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			len(winnerRollup.Transactions),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.State,
		))
		return winnerRollup.ToExtRollup(), true, nil
	}
	return nodecommon.ExtRollup{}, false, nil
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *obscurocore.Rollup) {
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address common.Address) uint64 {
	// todo user encryption
	return e.storage.CreateStateDB(e.storage.FetchHeadState().Head.Hash()).GetBalance(address)
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs *db.BlockState) *obscurocore.Rollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := speculativeRollup.s
	newRollupHeader := speculativeRollup.h

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
	if !speculativeRollup.found || (speculativeRollup.r.Hash() != bs.Head.Hash()) {
		if speculativeRollup.r != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
				obscurocommon.ShortAddress(e.node),
				obscurocommon.ShortHash(speculativeRollup.r.Hash()),
				speculativeRollup.r.Header.Number,
				obscurocommon.ShortHash(bs.Head.Hash()),
				bs.Head.Header.Number),
			)
			if e.statsCollector != nil {
				e.statsCollector.L2Recalc(e.node)
			}
		}

		newRollupHeader = obscurocore.NewHeader(bs.Head, bs.Head.Header.Number+1, e.node)
		// determine transactions to include in new rollup and process them
		newRollupTxs = currentTxs(bs.Head, e.mempool.FetchMempoolTxs(), e.storage)

		newRollupState = e.storage.CreateStateDB(bs.Head.Hash())
		executeTransactions(newRollupTxs, newRollupState, newRollupHeader)
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := e.blockResolver.Proof(bs.Head)
	depositTxs := processDeposits(proof, b, e.blockResolver, e.txHandler)
	executeTransactions(depositTxs, newRollupState, newRollupHeader)

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := obscurocore.NewRollupFromHeader(newRollupHeader, b.Hash(), newRollupTxs, obscurocommon.GenerateNonce(), newRollupState.StateRoot())

	// Postprocessing - withdrawals
	r.Header.Withdrawals = rollupPostProcessingWithdrawals(&r, newRollupState)

	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	// todo add some sort of cache
	rollup := e.storage.FetchHeadState().Head

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
	e.exitCh <- true
	return nil
}

func (e *enclaveImpl) Attestation() obscurocommon.AttestationReport {
	// Todo
	return obscurocommon.AttestationReport{Owner: e.node}
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

func (e *enclaveImpl) blockStateBlockSubmissionResponse(bs *db.BlockState, rollup nodecommon.ExtRollup) nodecommon.BlockSubmissionResponse {
	var head *nodecommon.Header
	if bs.FoundNewRollup {
		head = bs.Head.Header
	}
	return nodecommon.BlockSubmissionResponse{
		BlockHeader:    bs.Block.Header(),
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
	s     db.StateDB
	h     *nodecommon.Header
	txs   []nodecommon.L2Tx
}

// internal structure used for the speculative execution.
type processingEnvironment struct {
	headRollup      *obscurocore.Rollup             // the current head rollup, which will be the parent of the new rollup
	header          *nodecommon.Header              // the header of the new rollup
	processedTxs    []nodecommon.L2Tx               // txs that were already processed
	processedTxsMap map[common.Hash]nodecommon.L2Tx // structure used to prevent duplicates
	state           db.StateDB                      // the state as calculated from the previous rollup and the processed transactions
}

// NewEnclave creates a new enclave.
// `genesisJSON` is the configuration for the corresponding L1's genesis block. This is used to validate the blocks
// received from the L1 node if `validateBlocks` is set to true.
func NewEnclave(id common.Address, mining bool, txHandler mgmtcontractlib.TxHandler, validateBlocks bool, genesisJSON []byte, collector StatsCollector) nodecommon.Enclave {
	backingDB := db.NewInMemoryDB()
	storage := db.NewStorage(backingDB)

	var l1Blockchain *core.BlockChain
	if validateBlocks {
		if genesisJSON == nil {
			panic("enclave was configured to validate blocks, but genesis JSON was nil")
		}
		l1Blockchain = NewL1Blockchain(genesisJSON)
	} else {
		log.Log(fmt.Sprintf("Enclave-%d: validateBlocks is set to false. L1 blocks will not be validated.", obscurocommon.ShortAddress(id)))
	}

	return &enclaveImpl{
		node:                 id,
		mining:               mining,
		storage:              storage,
		blockResolver:        storage,
		mempool:              mempool.New(),
		statsCollector:       collector,
		l1Blockchain:         l1Blockchain,
		txCh:                 make(chan nodecommon.L2Tx),
		roundWinnerCh:        make(chan *obscurocore.Rollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan speculativeWork),
		txHandler:            txHandler,
	}
}
