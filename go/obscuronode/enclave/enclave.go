package enclave

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const ChainID = 777 // The unique ID for the Obscuro chain. Required for Geth signing.

// todo - this should become an elaborate data structure
type SharedEnclaveSecret []byte

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common.Address)
	RollupWithMoreRecentProof()
}

// BlockSubmissionResponse is the response sent from the enclave back to the node after ingesting a block
type BlockSubmissionResponse struct {
	L1Hash      obscurocommon.L1RootHash // The Header Hash of the ingested Block
	L1Height    uint64                   // The L1 Height of the ingested Block
	L1Parent    obscurocommon.L2RootHash // The L1 Parent of the ingested Block
	L2Hash      obscurocommon.L2RootHash // The Rollup Hash in the ingested Block
	L2Height    uint64                   // The Rollup Height in the ingested Block
	L2Parent    obscurocommon.L2RootHash // The Rollup Hash Parent inside the ingested Block
	Withdrawals []nodecommon.Withdrawal  // The Withdrawals available in Rollup of the ingested Block

	ProducedRollup    nodecommon.ExtRollup // The new Rollup when ingesting the block produces a new Rollup
	IngestedBlock     bool                 // Whether the Block was ingested or discarded
	IngestedNewRollup bool                 // Whether the Block had a new Rollup and the enclave has ingested it
}

// Enclave - The actual implementation of this interface will call an rpc service
type Enclave interface {
	// IsReady checks whether the enclave is ready to process requests
	IsReady() error

	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation() obscurocommon.AttestationReport

	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret() obscurocommon.EncryptedSharedEnclaveSecret

	// FetchSecret - return the shared secret encrypted with the key from the attestation
	FetchSecret(report obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret

	// InitEnclave - initialise an enclave with a seed received by another enclave
	InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret)

	// IsInitialised - true if the shared secret is avaible
	IsInitialised() bool

	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis() BlockSubmissionResponse

	// IngestBlocks - feed L1 blocks into the enclave to catch up
	IngestBlocks(blocks []*types.Block) []BlockSubmissionResponse

	// Start - start speculative execution
	Start(block types.Block)

	// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving a parent of it, will result in it being ignored
	SubmitBlock(block types.Block) BlockSubmissionResponse

	// SubmitRollup - receive gossiped rollups
	SubmitRollup(rollup nodecommon.ExtRollup)

	// SubmitTx - user transactions
	SubmitTx(tx nodecommon.EncryptedTx) error

	// Balance - returns the balance of an address with a block delay
	Balance(address common.Address) uint64

	// RoundWinner - calculates and returns the winner for a round, and whether this node is the winner
	RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool)

	// Stop gracefully stops the enclave
	Stop()

	// GetTransaction returns a transaction given its signed hash, or nil if the transaction is unknown
	GetTransaction(txHash common.Hash) *L2Tx

	// StopClient stops the enclave client if one exists
	StopClient()
}

type enclaveImpl struct {
	node           common.Address
	mining         bool
	db             DB
	blockResolver  obscurocommon.BlockResolver
	statsCollector StatsCollector

	txCh                 chan L2Tx
	roundWinnerCh        chan *Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork
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
	var currentHead *Rollup
	var currentState RollupState

	var currentProcessedTxs []L2Tx
	currentProcessedTxsMap := make(map[common.Hash]L2Tx)

	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.db.FetchState(block.Hash())
	if f {
		currentHead = blockState.Head
		if currentHead != nil {
			currentState = newProcessedState(e.db.FetchRollupState(currentHead.Hash()))
		}
	}

	// start the speculative rollup execution loop
	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:

			currentHead = winnerRollup
			currentState = newProcessedState(e.db.FetchRollupState(winnerRollup.Hash()))

			// determine the transactions that were not yet included
			currentProcessedTxs = currentTxs(winnerRollup, e.db.FetchTxs(), e.db)
			currentProcessedTxsMap = makeMap(currentProcessedTxs)

			// calculate the State after executing them
			currentState = executeTransactions(currentProcessedTxs, currentState)

		case tx := <-e.txCh:
			// only process transactions if there is already a rollup to use as parent
			if currentHead != nil {
				_, found := currentProcessedTxsMap[tx.Hash()]
				if !found {
					currentProcessedTxsMap[tx.Hash()] = tx
					currentProcessedTxs = append(currentProcessedTxs, tx)
					executeTx(&currentState, tx)
				}
			}

		case <-e.speculativeWorkInCh:
			b := make([]L2Tx, 0, len(currentProcessedTxs))
			b = append(b, currentProcessedTxs...)
			state := copyProcessedState(currentState)
			e.speculativeWorkOutCh <- speculativeWork{
				r:   currentHead,
				s:   &state,
				txs: b,
			}

		case <-e.exitCh:
			return
		}
	}
}

func (e *enclaveImpl) ProduceGenesis() BlockSubmissionResponse {
	return BlockSubmissionResponse{
		L2Hash:         GenesisRollup.Header.Hash(),
		L1Hash:         obscurocommon.GenesisHash,
		ProducedRollup: GenesisRollup.ToExtRollup(),
		IngestedBlock:  true,
	}
}

func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) []BlockSubmissionResponse {
	result := make([]BlockSubmissionResponse, len(blocks))
	for i, block := range blocks {
		e.db.StoreBlock(block)
		bs := updateState(block, e.db, e.blockResolver)
		if bs == nil {
			result[i] = BlockSubmissionResponse{
				L1Hash:            block.Hash(),
				L1Height:          e.blockResolver.HeightBlock(block),
				L1Parent:          block.ParentHash(),
				IngestedBlock:     true,
				IngestedNewRollup: false,
			}
		} else {
			var rollup nodecommon.ExtRollup
			if bs.foundNewRollup {
				rollup = bs.Head.ToExtRollup()
			}
			result[i] = BlockSubmissionResponse{
				L1Hash:            bs.Block.Hash(),
				L1Height:          e.blockResolver.HeightBlock(bs.Block),
				L1Parent:          bs.Block.ParentHash(),
				L2Hash:            bs.Head.Hash(),
				L2Height:          e.db.HeightRollup(bs.Head),
				L2Parent:          bs.Head.Header.ParentHash,
				ProducedRollup:    rollup,
				IngestedBlock:     true,
				IngestedNewRollup: bs.foundNewRollup,
			}
		}
	}

	return result
}

func (e *enclaveImpl) SubmitBlock(block types.Block) BlockSubmissionResponse {
	_, foundBlock := e.db.ResolveBlock(block.Hash())
	if foundBlock {
		return BlockSubmissionResponse{IngestedBlock: false}
	}

	e.db.StoreBlock(&block)
	// this is where much more will actually happen.
	// the "blockchain" logic from geth has to be executed here,
	// to determine the total proof of work, to verify some key aspects, etc

	_, f := e.db.ResolveBlock(block.Header().ParentHash)
	if !f && e.db.HeightBlock(&block) > obscurocommon.L1GenesisHeight {
		return BlockSubmissionResponse{IngestedBlock: false}
	}
	blockState := updateState(&block, e.db, e.blockResolver)

	if blockState == nil {
		return BlockSubmissionResponse{
			L1Hash:            block.Hash(),
			L1Height:          e.blockResolver.HeightBlock(&block),
			L1Parent:          block.ParentHash(),
			IngestedBlock:     true,
			IngestedNewRollup: false,
		}
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	e.db.PruneTxs(historicTxs(blockState.Head, e.db))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.db.StoreRollup(r)

	return BlockSubmissionResponse{
		L1Hash:      block.Hash(),
		L1Height:    e.blockResolver.HeightBlock(&block),
		L1Parent:    blockState.Block.Header().ParentHash,
		L2Hash:      blockState.Head.Hash(),
		L2Height:    blockState.Head.Height.Load().(uint64),
		L2Parent:    blockState.Head.Header.ParentHash,
		Withdrawals: blockState.Head.Header.Withdrawals,

		ProducedRollup:    r.ToExtRollup(),
		IngestedBlock:     true,
		IngestedNewRollup: blockState.foundNewRollup,
	}
}

func (e *enclaveImpl) SubmitRollup(rollup nodecommon.ExtRollup) {
	r := Rollup{
		Header:       rollup.Header,
		Transactions: decryptTransactions(rollup.Txs),
	}

	// only store if the parent exists
	if e.db.ExistRollup(r.Header.ParentHash) {
		e.db.StoreRollup(&r)
	} else {
		log.Log(fmt.Sprintf("Agg%d:> Received rollup with no parent: r_%d", obscurocommon.ShortAddress(e.node), obscurocommon.ShortHash(r.Hash())))
	}
}

func (e *enclaveImpl) SubmitTx(tx nodecommon.EncryptedTx) error {
	decryptedTx := DecryptTx(tx)
	err := verifySignature(&decryptedTx)
	if err != nil {
		return err
	}
	e.db.StoreTx(decryptedTx)
	e.txCh <- decryptedTx
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(decryptedTx *L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(ChainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	head := e.db.FetchRollup(parent)

	rollupsReceivedFromPeers := e.db.FetchGossipedRollups(e.db.HeightRollup(head) + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.db.ParentRollup(rol)
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.db.FetchRollupState(head.Hash())
	// determine the winner of the round
	winnerRollup, s := findRoundWinner(usefulRollups, head, parentState, e.db, e.blockResolver)
	// nodecommon.Log(fmt.Sprintf(">   Agg%d: Round=r_%d Winner=r_%d(%d)[r_%d]{proof=b_%d}.", e.node, parent.ID(),
	// winnerRollup.L2RootHash.ID(), winnerRollup.Height(), winnerRollup.Parent().L2RootHash.ID(),
	// winnerRollup.Proof().L2RootHash.ID()))

	e.db.SetRollupState(winnerRollup.Hash(), s)
	go e.notifySpeculative(winnerRollup)

	// we are the winner
	if winnerRollup.Header.Agg == e.node {
		v := winnerRollup.Proof(e.blockResolver)
		w := e.db.ParentRollup(winnerRollup)
		log.Log(fmt.Sprintf(">   Agg%d: create rollup=r_%d(%d)[r_%d]{proof=b_%d}. Txs: %v. State=%v.",
			obscurocommon.ShortAddress(e.node),
			obscurocommon.ShortHash(winnerRollup.Hash()), e.db.HeightRollup(winnerRollup),
			obscurocommon.ShortHash(w.Hash()),
			obscurocommon.ShortHash(v.Hash()),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.State),
		)
		return winnerRollup.ToExtRollup(), true
	}
	return nodecommon.ExtRollup{}, false
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *Rollup) {
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address common.Address) uint64 {
	// todo user encryption
	return e.db.Balance(address)
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs *BlockState) *Rollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := *speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
	// if true {
	if (speculativeRollup.r == nil) || (speculativeRollup.r.Hash() != bs.Head.Hash()) {
		if speculativeRollup.r != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
				obscurocommon.ShortAddress(e.node),
				obscurocommon.ShortHash(speculativeRollup.r.Hash()),
				e.db.HeightRollup(speculativeRollup.r),
				obscurocommon.ShortHash(bs.Head.Hash()),
				e.db.HeightRollup(bs.Head)),
			)
			e.statsCollector.L2Recalc(e.node)
		}

		// determine transactions to include in new rollup and process them
		newRollupTxs = currentTxs(bs.Head, e.db.FetchTxs(), e.db)
		newRollupState = executeTransactions(newRollupTxs, newProcessedState(bs.State))
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := bs.Head.Proof(e.blockResolver)
	newRollupState = processDeposits(proof, b, copyProcessedState(newRollupState), e.blockResolver)

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := NewRollup(b, bs.Head, e.node, newRollupTxs, newRollupState.w, obscurocommon.GenerateNonce(), serialize(newRollupState.s))
	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) *L2Tx {
	// todo add some sort of cache
	rollup := e.db.Head().Head
	for {
		txs := rollup.Transactions
		for _, tx := range txs {
			if tx.Hash() == txHash {
				return &tx
			}
		}
		rollup = e.db.FetchRollup(rollup.Header.ParentHash)
		if rollup.Height.Load() == obscurocommon.L2GenesisHeight {
			return nil
		}
	}
}

func (e *enclaveImpl) Stop() {
	e.exitCh <- true
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
	e.db.StoreSecret(secret)
	return encryptSecret(secret)
}

// InitEnclave - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) InitEnclave(secret obscurocommon.EncryptedSharedEnclaveSecret) {
	e.db.StoreSecret(decryptSecret(secret))
}

func (e *enclaveImpl) FetchSecret(obscurocommon.AttestationReport) obscurocommon.EncryptedSharedEnclaveSecret {
	return encryptSecret(e.db.FetchSecret())
}

func (e *enclaveImpl) IsInitialised() bool {
	return e.db.FetchSecret() != nil
}

// Todo - implement with crypto
func decryptSecret(secret obscurocommon.EncryptedSharedEnclaveSecret) SharedEnclaveSecret {
	return SharedEnclaveSecret(secret)
}

// Todo - implement with crypto
func encryptSecret(secret SharedEnclaveSecret) obscurocommon.EncryptedSharedEnclaveSecret {
	return obscurocommon.EncryptedSharedEnclaveSecret(secret)
}

// internal structure to pass information.
type speculativeWork struct {
	r   *Rollup
	s   *RollupState
	txs []L2Tx
}

func NewEnclave(id common.Address, mining bool, collector StatsCollector) Enclave {
	db := NewInMemoryDB()
	return &enclaveImpl{
		node:                 id,
		db:                   db,
		blockResolver:        db,
		mining:               mining,
		txCh:                 make(chan L2Tx),
		roundWinnerCh:        make(chan *Rollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan speculativeWork),
		statsCollector:       collector,
	}
}
