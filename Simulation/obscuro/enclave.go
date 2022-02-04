package obscuro

import (
	"fmt"
	"simulation/common"
	wallet_mock "simulation/wallet-mock"
)

// Enclave - The actual implementation of this interface will call an rpc service
type Enclave interface {
	// Todo - attestation, secret generation, etc

	// SubmitBlock - When a new round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the rollup
	SubmitBlock(block common.EncodedBlock) (common.EncodedRollup, common.RootHash)

	Stop()
	Start()

	// SubmitRollup - receive gossiped rollups
	SubmitRollup(rollup common.EncodedRollup)

	// SubmitTx - user transactions
	SubmitTx(tx common.EncodedL2Tx)

	// Balance
	Balance(address wallet_mock.Address) uint64

	// RoundWinner - calculates and returns the winner for a round
	RoundWinner(parent common.RootHash) (common.EncodedRollup, bool)

	// PeekHead - only availble for testing purposes
	PeekHead() BlockState
}

type enclaveImpl struct {
	node           common.NodeId
	mining         bool
	db             Db
	statsCollector common.StatsCollector

	txCh                 chan common.L2Tx
	roundWinnerCh        chan common.Rollup
	exitCh               chan bool
	speculativeWorkInCh  chan bool
	speculativeWorkOutCh chan currentWork
}

func NewEnclave(id common.NodeId, mining bool, collector common.StatsCollector) Enclave {
	return &enclaveImpl{
		node:                 id,
		db:                   NewInMemoryDb(),
		mining:               mining,
		txCh:                 make(chan common.L2Tx),
		roundWinnerCh:        make(chan common.Rollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan currentWork),
		statsCollector:       collector,
	}
}

func (e *enclaveImpl) Start() {
	var currentHead common.Rollup
	var currentState ProcessedState
	var currentProcessedTxs []common.L2Tx

	//start the speculative rollup execution loop
	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:

			currentHead = winnerRollup
			currentState = newProcessedState(e.db.FetchRollupState(winnerRollup.RootHash))

			// determine the transactions that were not yet included
			currentProcessedTxs = currentTxs(winnerRollup, e.db.FetchTxs())

			// calculate the State after executing them
			currentState = executeTransactions(currentProcessedTxs, currentState)

		case tx := <-e.txCh:
			currentProcessedTxs = append(currentProcessedTxs, tx)
			executeTx(&currentState, tx)

		case <-e.speculativeWorkInCh:
			b := make([]common.L2Tx, len(currentProcessedTxs))
			copy(b, currentProcessedTxs)
			e.speculativeWorkOutCh <- currentWork{
				r:   currentHead,
				s:   copyProcessedState(currentState),
				txs: b,
			}
		case <-e.exitCh:
			return
		}
	}
}

func (e *enclaveImpl) SubmitBlock(block common.EncodedBlock) (common.EncodedRollup, common.RootHash) {
	b := common.DecodeBlock(block)

	blockState := updateState(b, e.db)

	if e.mining {
		e.db.PruneTxs(historicTxs(blockState.Head))

		r := e.produceRollup(b, blockState)
		e.db.StoreRollup(r.Height(), r)

		return common.EncodeRollup(r), blockState.Head.RootHash
	}

	return nil, blockState.Head.RootHash
}

func (e *enclaveImpl) SubmitRollup(rollup common.EncodedRollup) {
	r := common.DecodeRollup(rollup)
	e.db.StoreRollup(r.Height(), r)
}

func (e *enclaveImpl) SubmitTx(tx common.EncodedL2Tx) {
	t := common.DecodeTx(tx)
	e.db.StoreTx(t)
	e.txCh <- t
}

func (e *enclaveImpl) RoundWinner(parent common.RootHash) (common.EncodedRollup, bool) {

	head := e.db.FetchRollup(parent)
	rollupsReceivedFromPeers := e.db.FetchRollups(head.H + 1)
	// filter out rollups with a different Parent
	var usefulRollups []common.Rollup

	for _, rol := range rollupsReceivedFromPeers {
		if rol.Parent().Root() == head.Root() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.db.FetchRollupState(head.RootHash)
	// determine the winner of the round
	winnerRollup, s := findRoundWinner(usefulRollups, head, parentState)
	//common.Log(fmt.Sprintf(">   Agg%d: Round=r_%d Winner=r_%d(%d)[r_%d]{proof=b_%d}.", e.node, parent.ID(), winnerRollup.Root().ID(), winnerRollup.Height(), winnerRollup.Parent().Root().ID(), winnerRollup.Proof().RootHash.ID()))

	e.db.SetRollupState(winnerRollup.RootHash, s)

	// we are the winner
	if winnerRollup.Agg == e.node {
		common.Log(fmt.Sprintf(">   Agg%d: create rollup=r_%d(%d)[r_%d]{proof=b_%d}. Txs: %v. State=%v.", e.node, winnerRollup.Root().ID(), winnerRollup.Height(), winnerRollup.Parent().Root().ID(), winnerRollup.Proof().RootHash.ID(), printTxs(winnerRollup.Transactions), winnerRollup.State))
		result := common.EncodeRollup(winnerRollup)
		return result, true
	}
	e.roundWinnerCh <- winnerRollup
	return nil, false
}

func (e *enclaveImpl) Balance(address wallet_mock.Address) uint64 {
	//todo
	return 0
}

func (e *enclaveImpl) produceRollup(b common.Block, bs BlockState) common.Rollup {

	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip
	if speculativeRollup.r.Root() != bs.Head.Root() {
		common.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)", e.node, speculativeRollup.r.Root().ID(), speculativeRollup.r.Height(), bs.Head.Root().ID(), bs.Head.Height()))
		e.statsCollector.L2Recalc(e.node)

		// determine transactions to include in new rollup and process them
		newRollupTxs = currentTxs(bs.Head, e.db.FetchTxs())
		newRollupState = executeTransactions(newRollupTxs, newProcessedState(bs.State))
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := bs.Head.Proof()
	newRollupState = processDeposits(&proof, b, copyProcessedState(newRollupState))

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	return common.NewRollup(&b, &bs.Head, e.node, newRollupTxs, newRollupState.w, serialize(newRollupState.s))
}

func (e *enclaveImpl) PeekHead() BlockState {
	return e.db.Head()
}

func (e *enclaveImpl) Stop() {
	e.exitCh <- true
}

func printTxs(txs []common.L2Tx) (txsString []string) {
	for _, t1 := range txs {
		switch t1.TxType {
		case common.TransferTx:
			txsString = append(txsString, fmt.Sprintf("%v->%v(%d)", t1.From, t1.To, t1.Amount))
		case common.WithdrawalTx:
			txsString = append(txsString, fmt.Sprintf("%v->*(%d)", t1.From, t1.Amount))
		}
	}
	return txsString
}
