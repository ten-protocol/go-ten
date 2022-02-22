package enclave

import (
	"fmt"
	common3 "github.com/otherview/obscuro-playground/go/common"
	common2 "github.com/otherview/obscuro-playground/go/obscuro-node/common"
)

type StatsCollector interface {
	// Register when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id common3.NodeId)
	RollupWithMoreRecentProof()
}

// The response sent back to the node
type SubmitBlockResponse struct {
	Hash      common3.L2RootHash
	Rollup    common2.ExtRollup
	Processed bool
}

// Enclave - The actual implementation of this interface will call an rpc service
type Enclave interface {
	// Todo - attestation, secret generation, etc

	ProduceGenesis() SubmitBlockResponse
	IngestBlocks(blocks []common3.ExtBlock)
	Start(block common3.ExtBlock)

	// SubmitBlock - When a new round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the rollup
	SubmitBlock(block common3.ExtBlock) SubmitBlockResponse

	// SubmitRollup - receive gossiped rollups
	SubmitRollup(rollup common2.ExtRollup)

	// SubmitTx - user transactions
	SubmitTx(tx common2.EncryptedTx)

	// Balance - returns the balance of an address with a block delay
	Balance(address common3.Address) uint64

	// RoundWinner - calculates and returns the winner for a round
	RoundWinner(parent common3.L2RootHash) (common2.ExtRollup, bool)
	Stop()

	// TestPeekHead - only availble for testing purposes
	TestPeekHead() BlockState

	// TestDb - only availble for testing purposes
	TestDb() Db
}

type enclaveImpl struct {
	node   common3.NodeId
	mining         bool
	db             Db
	statsCollector StatsCollector

	txCh                 chan L2Tx
	roundWinnerCh        chan *EnclaveRollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan speculativeWork
}

func (e *enclaveImpl) Start(block common3.ExtBlock) {
	s, f := e.db.FetchState(block.Header.Hash())
	if !f {
		panic("state should be calculated")
	}

	currentHead := s.Head
	currentState := newProcessedState(e.db.FetchRollupState(currentHead.Hash()))
	var currentProcessedTxs []L2Tx
	currentProcessedTxsMap := make(map[common3.TxHash]L2Tx)

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
			_, f := currentProcessedTxsMap[tx.Id]
			if !f {
				currentProcessedTxsMap[tx.Id] = tx
				currentProcessedTxs = append(currentProcessedTxs, tx)
				executeTx(&currentState, tx)
			}

		case <-e.speculativeWorkInCh:
			b := make([]L2Tx, 0)
			for _, tx := range currentProcessedTxs {
				b = append(b, tx)
			}
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
	return SubmitBlockResponse{
		Hash:      r.Header.Hash(),
		Rollup:    r.ToExtRollup(),
		Processed: true,
	}
}

func (e *enclaveImpl) IngestBlocks(blocks []common3.ExtBlock) {
	for _, block := range blocks {
		b := block.ToBlock()
		e.db.Store(b)
		updateState(b, e.db)
	}
}

func (e *enclaveImpl) SubmitBlock(block common3.ExtBlock) SubmitBlockResponse {
	b := block.ToBlock()
	e.db.Store(b)
	// this is where much more will actually happen.
	// the "blockchain" logic from geth has to be executed here, to determine the total proof of work, to verify some key aspects, etc

	_, f := e.db.Resolve(b.Header.ParentHash)
	if !f && b.Height(e.db) > common3.L1GenesisHeight {
		return SubmitBlockResponse{Processed: false}
	}
	blockState := updateState(b, e.db)

	if e.mining {
		e.db.PruneTxs(historicTxs(blockState.Head, e.db))

		r := e.produceRollup(b, blockState)
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

func (e *enclaveImpl) SubmitRollup(rollup common2.ExtRollup) {
	r := EnclaveRollup{
		Header:       rollup.Header,
		Transactions: decryptTransactions(rollup.Txs),
	}
	e.db.StoreRollup(&r)
}

func (e *enclaveImpl) SubmitTx(tx common2.EncryptedTx) {
	t := DecryptTx(tx)
	e.db.StoreTx(t)
	e.txCh <- t
}

func (e *enclaveImpl) RoundWinner(parent common3.L2RootHash) (common2.ExtRollup, bool) {
	head := e.db.FetchRollup(parent)

	rollupsReceivedFromPeers := e.db.FetchRollups(e.db.Height(head) + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*EnclaveRollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.db.Parent(rol)
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.db.FetchRollupState(head.Hash())
	// determine the winner of the round
	winnerRollup, s := findRoundWinner(usefulRollups, head, parentState, e.db)
	// common.Log(fmt.Sprintf(">   Agg%d: Round=r_%d Winner=r_%d(%d)[r_%d]{proof=b_%d}.", e.node, parent.ID(), winnerRollup.L2RootHash.ID(), winnerRollup.Height(), winnerRollup.Parent().L2RootHash.ID(), winnerRollup.Proof().L2RootHash.ID()))

	e.db.SetRollupState(winnerRollup.Hash(), s)
	go e.notifySpeculative(winnerRollup)

	// we are the winner
	if winnerRollup.Header.Agg == e.node {
		v := winnerRollup.Proof(e.db)
		w := e.db.Parent(winnerRollup)
		common3.Log(fmt.Sprintf(">   Agg%d: create rollup=r_%s(%d)[r_%s]{proof=b_%s}. Txs: %v. State=%v.", e.node, common3.Str(winnerRollup.Hash()), e.db.Height(winnerRollup), common3.Str(w.Hash()), common3.Str(v.Hash()), printTxs(winnerRollup.Transactions), winnerRollup.Header.State))
		return winnerRollup.ToExtRollup(), true
	}
	return common2.ExtRollup{}, false
}

func (e *enclaveImpl) notifySpeculative(winnerRollup *EnclaveRollup) {
	//if atomic.LoadInt32(e.interrupt) == 1 {
	//	return
	//}
	e.roundWinnerCh <- winnerRollup
}

func (e *enclaveImpl) Balance(address common3.Address) uint64 {
	// todo
	return 0
}

func (e *enclaveImpl) produceRollup(b *common3.Block, bs BlockState) *EnclaveRollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := *speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
	// if true {
	if (speculativeRollup.r == nil) || (speculativeRollup.r.Hash() != bs.Head.Hash()) {
		if speculativeRollup.r != nil {
			common3.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%s(%d), published=r_%s(%d)", e.node, common3.Str(speculativeRollup.r.Hash()), e.db.Height(speculativeRollup.r), common3.Str(bs.Head.Hash()), e.db.Height(bs.Head)))
			e.statsCollector.L2Recalc(e.node)
		}

		// determine transactions to include in new rollup and process them
		newRollupTxs = currentTxs(bs.Head, e.db.FetchTxs(), e.db)
		newRollupState = executeTransactions(newRollupTxs, newProcessedState(bs.State))
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := bs.Head.Proof(e.db)
	newRollupState = processDeposits(proof, b, copyProcessedState(newRollupState), e.db)

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := NewRollup(b, bs.Head, e.node, newRollupTxs, newRollupState.w, common3.GenerateNonce(), serialize(newRollupState.s))
	// h := r.Height(e.db)
	// fmt.Printf("h:=%d\n", h)
	return &r
}

func (e *enclaveImpl) TestPeekHead() BlockState {
	return e.db.Head()
}

func (e *enclaveImpl) TestDb() Db {
	return e.db
}

func (e *enclaveImpl) Stop() {
	e.exitCh <- true
}

// internal structure to pass information.
type speculativeWork struct {
	r   *EnclaveRollup
	s   *RollupState
	txs []L2Tx
}

func NewEnclave(id common3.NodeId, mining bool, collector StatsCollector) Enclave {
	return &enclaveImpl{
		node:                 id,
		db:                   NewInMemoryDb(),
		mining:               mining,
		txCh:                 make(chan L2Tx),
		roundWinnerCh:        make(chan *EnclaveRollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan speculativeWork),
		statsCollector:       collector,
	}
}
