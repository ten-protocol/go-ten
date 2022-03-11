package enclave

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

const ChainID = 777 // The unique ID for the Obscuro chain. Required for Geth signing.

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id gethcommon.Address)
	RollupWithMoreRecentProof()
}

// The response sent back to the node
type SubmitBlockResponse struct {
	Hash      common.L2RootHash
	Rollup    nodecommon.ExtRollup
	Processed bool
}

// Enclave - The actual implementation of this interface will call an rpc service
type Enclave interface {
	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation() common.AttestationReport

	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret() common.EncryptedSharedEnclaveSecret

	// FetchSecret - return the shared secret encrypted with the key from the attestation
	FetchSecret(report common.AttestationReport) common.EncryptedSharedEnclaveSecret

	// Init - initialise an enclave with a seed received by another enclave
	Init(secret common.EncryptedSharedEnclaveSecret)

	// IsInitialised - true if the shared secret is avaible
	IsInitialised() bool

	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis() SubmitBlockResponse

	// IngestBlocks - feed L1 blocks into the enclave to catch up
	IngestBlocks(blocks []*types.Block)

	// Start - start speculative execution
	Start(block types.Block)

	// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving a parent of it, will result in it being ignored
	SubmitBlock(block types.Block) SubmitBlockResponse

	// SubmitRollup - receive gossiped rollups
	SubmitRollup(rollup nodecommon.ExtRollup)

	// SubmitTx - user transactions
	SubmitTx(tx nodecommon.EncryptedTx) error

	// Balance - returns the balance of an address with a block delay
	Balance(address gethcommon.Address) uint64

	// RoundWinner - calculates and returns the winner for a round
	RoundWinner(parent common.L2RootHash) (nodecommon.ExtRollup, bool)
	Stop()

	// TestPeekHead - only available for testing purposes
	TestPeekHead() BlockState

	// TestDb - only available for testing purposes
	TestDB() DB
}

type enclaveImpl struct {
	node           gethcommon.Address
	mining         bool
	db             DB
	blockResolver  common.BlockResolver
	statsCollector StatsCollector

	txCh                 chan L2Tx
	roundWinnerCh        chan *Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork
}

func (e *enclaveImpl) Start(block types.Block) {
	headerHash := block.Hash()
	s, f := e.db.FetchState(headerHash)
	if !f {
		panic("state should be calculated")
	}

	currentHead := s.Head
	currentState := newProcessedState(e.db.FetchRollupState(currentHead.Hash()))
	var currentProcessedTxs []L2Tx
	currentProcessedTxsMap := make(map[gethcommon.Hash]L2Tx)

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
			_, found := currentProcessedTxsMap[tx.Hash()]
			if !found {
				currentProcessedTxsMap[tx.Hash()] = tx
				currentProcessedTxs = append(currentProcessedTxs, tx)
				executeTx(&currentState, tx)
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

func (e *enclaveImpl) ProduceGenesis() SubmitBlockResponse {
	r := GenesisRollup
	hash := r.Header.Hash()
	return SubmitBlockResponse{
		Hash:      hash,
		Rollup:    r.ToExtRollup(),
		Processed: true,
	}
}

func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) {
	for _, block := range blocks {
		e.db.StoreBlock(block)
		updateState(block, e.db, e.blockResolver)
	}
}

func (e *enclaveImpl) SubmitBlock(block types.Block) SubmitBlockResponse {
	// Todo - investigate further why this is needed.
	// So far this seems to recover correctly
	defer func() {
		if r := recover(); r != nil {
			common.Log(fmt.Sprintf("Agg%d Panic %s", common.ShortAddress(e.node), r))
		}
	}()

	_, foundBlock := e.db.ResolveBlock(block.Hash())
	if foundBlock {
		return SubmitBlockResponse{Processed: false}
	}

	e.db.StoreBlock(&block)
	// this is where much more will actually happen.
	// the "blockchain" logic from geth has to be executed here,
	// to determine the total proof of work, to verify some key aspects, etc

	_, f := e.db.ResolveBlock(block.Header().ParentHash)
	if !f && e.db.HeightBlock(&block) > common.L1GenesisHeight {
		return SubmitBlockResponse{Processed: false}
	}
	blockState := updateState(&block, e.db, e.blockResolver)

	if e.mining {
		e.db.PruneTxs(historicTxs(blockState.Head, e.db))

		r := e.produceRollup(&block, blockState)
		e.db.StoreRollup(r)

		return SubmitBlockResponse{
			Hash:      blockState.Head.Hash(),
			Rollup:    r.ToExtRollup(),
			Processed: true,
		}
	}

	return SubmitBlockResponse{
		Hash:      blockState.Head.Hash(),
		Processed: true,
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
		common.Log(fmt.Sprintf("Agg%d:> Received rollup with no parent: r_%d", common.ShortAddress(e.node), common.ShortHash(r.Hash())))
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

func (e *enclaveImpl) RoundWinner(parent common.L2RootHash) (nodecommon.ExtRollup, bool) {
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
		common.Log(fmt.Sprintf(">   Agg%d: create rollup=r_%d(%d)[r_%d]{proof=b_%d}. Txs: %v. State=%v.",
			common.ShortAddress(e.node),
			common.ShortHash(winnerRollup.Hash()), e.db.HeightRollup(winnerRollup),
			common.ShortHash(w.Hash()),
			common.ShortHash(v.Hash()),
			printTxs(winnerRollup.Transactions),
			winnerRollup.Header.State),
		)
		return winnerRollup.ToExtRollup(), true
	}
	return nodecommon.ExtRollup{}, false
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *Rollup) {
	//if atomic.LoadInt32(e.interrupt) == 1 {
	//	return
	//}
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address gethcommon.Address) uint64 {
	// todo
	return 0
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs BlockState) *Rollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := *speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
	// if true {
	if (speculativeRollup.r == nil) || (speculativeRollup.r.Hash() != bs.Head.Hash()) {
		if speculativeRollup.r != nil {
			common.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
				common.ShortAddress(e.node),
				common.ShortHash(speculativeRollup.r.Hash()),
				e.db.HeightRollup(speculativeRollup.r),
				common.ShortHash(bs.Head.Hash()),
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
	r := NewRollup(b, bs.Head, e.node, newRollupTxs, newRollupState.w, common.GenerateNonce(), serialize(newRollupState.s))
	// h := r.Height(e.db)
	// fmt.Printf("h:=%d\n", h)
	return &r
}

func (e *enclaveImpl) TestPeekHead() BlockState {
	return e.db.Head()
}

func (e *enclaveImpl) TestDB() DB {
	return e.db
}

func (e *enclaveImpl) Stop() {
	e.exitCh <- true
}

func (e *enclaveImpl) Attestation() common.AttestationReport {
	// Todo
	return common.AttestationReport{Owner: e.node}
}

// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
func (e *enclaveImpl) GenerateSecret() common.EncryptedSharedEnclaveSecret {
	secret := make([]byte, 32)
	n, err := rand.Read(secret)
	if n != 32 || err != nil {
		panic(fmt.Sprintf("Could not generate secret: %s", err))
	}
	e.db.StoreSecret(secret)
	return encryptSecret(secret)
}

// Init - initialise an enclave with a seed received by another enclave
func (e *enclaveImpl) Init(secret common.EncryptedSharedEnclaveSecret) {
	e.db.StoreSecret(decryptSecret(secret))
}

func (e *enclaveImpl) FetchSecret(report common.AttestationReport) common.EncryptedSharedEnclaveSecret {
	return encryptSecret(e.db.FetchSecret())
}

func (e *enclaveImpl) IsInitialised() bool {
	return e.db.FetchSecret() != nil
}

// Todo - implement with crypto
func decryptSecret(secret common.EncryptedSharedEnclaveSecret) SharedEnclaveSecret {
	return SharedEnclaveSecret(secret)
}

// Todo - implement with crypto
func encryptSecret(secret SharedEnclaveSecret) common.EncryptedSharedEnclaveSecret {
	return common.EncryptedSharedEnclaveSecret(secret)
}

// internal structure to pass information.
type speculativeWork struct {
	r   *Rollup
	s   *RollupState
	txs []L2Tx
}

func NewEnclave(id gethcommon.Address, mining bool, collector StatsCollector) Enclave {
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
