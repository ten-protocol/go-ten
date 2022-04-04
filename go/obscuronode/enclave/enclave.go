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

type enclaveImpl struct {
	node           common.Address
	mining         bool
	storage        Storage
	blockResolver  BlockResolver
	statsCollector StatsCollector

	txCh                 chan nodecommon.L2Tx
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
	var currentProcessedTxs []nodecommon.L2Tx
	currentProcessedTxsMap := make(map[common.Hash]nodecommon.L2Tx)
	// determine whether the block where the speculative execution will start already contains Obscuro state
	blockState, f := e.storage.FetchBlockState(block.Hash())
	if f {
		currentHead = blockState.head
		if currentHead != nil {
			currentState = newProcessedState(e.storage.FetchRollupState(currentHead.Hash()))
		}
	}

	for {
		select {
		// A new winner was found after gossiping. Start speculatively executing incoming transactions to already have a rollup ready when the next round starts.
		case winnerRollup := <-e.roundWinnerCh:

			currentHead = winnerRollup
			currentState = newProcessedState(e.storage.FetchRollupState(winnerRollup.Hash()))

			// determine the transactions that were not yet included
			currentProcessedTxs = currentTxs(winnerRollup, e.storage.FetchMempoolTxs(), e.storage)
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
			b := make([]nodecommon.L2Tx, 0, len(currentProcessedTxs))
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

func (e *enclaveImpl) ProduceGenesis() nodecommon.BlockSubmissionResponse {
	return nodecommon.BlockSubmissionResponse{
		L2Hash:         GenesisRollup.Header.Hash(),
		L1Hash:         obscurocommon.GenesisHash,
		ProducedRollup: GenesisRollup.ToExtRollup(),
		IngestedBlock:  true,
	}
}

func (e *enclaveImpl) IngestBlocks(blocks []*types.Block) []nodecommon.BlockSubmissionResponse {
	result := make([]nodecommon.BlockSubmissionResponse, len(blocks))
	for i, block := range blocks {
		e.storage.StoreBlock(block)
		bs := updateState(block, e.storage, e.blockResolver)
		if bs == nil {
			result[i] = nodecommon.BlockSubmissionResponse{
				L1Hash:            block.Hash(),
				L1Height:          e.blockResolver.HeightBlock(block),
				L1Parent:          block.ParentHash(),
				IngestedBlock:     true,
				IngestedNewRollup: false,
			}
		} else {
			var rollup nodecommon.ExtRollup
			if bs.foundNewRollup {
				rollup = bs.head.ToExtRollup()
			}
			result[i] = nodecommon.BlockSubmissionResponse{
				L1Hash:            bs.block.Hash(),
				L1Height:          e.blockResolver.HeightBlock(bs.block),
				L1Parent:          bs.block.ParentHash(),
				L2Hash:            bs.head.Hash(),
				L2Height:          bs.head.Header.Height,
				L2Parent:          bs.head.Header.ParentHash,
				ProducedRollup:    rollup,
				IngestedBlock:     true,
				IngestedNewRollup: bs.foundNewRollup,
			}
		}
	}

	return result
}

func (e *enclaveImpl) SubmitBlock(block types.Block) nodecommon.BlockSubmissionResponse {
	_, foundBlock := e.storage.FetchBlock(block.Hash())
	if foundBlock {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false}
	}

	e.storage.StoreBlock(&block)
	// this is where much more will actually happen.
	// the "blockchain" logic from geth has to be executed here,
	// to determine the total proof of work, to verify some key aspects, etc

	_, f := e.storage.FetchBlock(block.Header().ParentHash)
	if !f && e.storage.HeightBlock(&block) > obscurocommon.L1GenesisHeight {
		return nodecommon.BlockSubmissionResponse{IngestedBlock: false}
	}
	blockState := updateState(&block, e.storage, e.blockResolver)

	if blockState == nil {
		return nodecommon.BlockSubmissionResponse{
			L1Hash:            block.Hash(),
			L1Height:          e.blockResolver.HeightBlock(&block),
			L1Parent:          block.ParentHash(),
			IngestedBlock:     true,
			IngestedNewRollup: false,
		}
	}

	// todo - A verifier node will not produce rollups, we can check the e.mining to get the node behaviour
	e.storage.RemoveMempoolTxs(historicTxs(blockState.head, e.storage))
	r := e.produceRollup(&block, blockState)
	// todo - should store proposal rollups in a different storage as they are ephemeral (round based)
	e.storage.StoreRollup(r)

	return nodecommon.BlockSubmissionResponse{
		L1Hash:      block.Hash(),
		L1Height:    e.blockResolver.HeightBlock(&block),
		L1Parent:    blockState.block.Header().ParentHash,
		L2Hash:      blockState.head.Hash(),
		L2Height:    blockState.head.Header.Height,
		L2Parent:    blockState.head.Header.ParentHash,
		Withdrawals: blockState.head.Header.Withdrawals,

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
	_, found := e.storage.FetchRollup(r.Header.ParentHash)
	if found {
		e.storage.StoreRollup(&r)
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
	e.storage.AddMempoolTx(decryptedTx)
	e.txCh <- decryptedTx
	return nil
}

// Checks that the L2Tx has a valid signature.
func verifySignature(decryptedTx *nodecommon.L2Tx) error {
	signer := types.NewLondonSigner(big.NewInt(ChainID))
	_, err := types.Sender(signer, decryptedTx)
	return err
}

func (e *enclaveImpl) RoundWinner(parent obscurocommon.L2RootHash) (nodecommon.ExtRollup, bool) {
	head, found := e.storage.FetchRollup(parent)
	if !found {
		panic(fmt.Sprintf("Could not find rollup: r_%s", parent))
	}

	rollupsReceivedFromPeers := e.storage.FetchRollups(head.Header.Height + 1)
	// filter out rollups with a different Parent
	var usefulRollups []*Rollup
	for _, rol := range rollupsReceivedFromPeers {
		p := e.storage.ParentRollup(rol)
		if p.Hash() == head.Hash() {
			usefulRollups = append(usefulRollups, rol)
		}
	}

	parentState := e.storage.FetchRollupState(head.Hash())
	// determine the winner of the round
	winnerRollup, s := findRoundWinner(usefulRollups, head, parentState, e.storage, e.blockResolver)

	e.storage.SetRollupState(winnerRollup.Hash(), s)
	go e.notifySpeculative(winnerRollup)

	// we are the winner
	if winnerRollup.Header.Agg == e.node {
		v := winnerRollup.Proof(e.blockResolver)
		w := e.storage.ParentRollup(winnerRollup)
		log.Log(fmt.Sprintf(">   Agg%d: create rollup=r_%d(%d)[r_%d]{proof=b_%d}. Txs: %v. State=%v.",
			obscurocommon.ShortAddress(e.node),
			obscurocommon.ShortHash(winnerRollup.Hash()), winnerRollup.Header.Height,
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
	return e.storage.FetchHeadState().state[address]
}

func (e *enclaveImpl) produceRollup(b *types.Block, bs *blockState) *Rollup {
	// retrieve the speculatively calculated State based on the previous winner and the incoming transactions
	e.speculativeWorkInCh <- true
	speculativeRollup := <-e.speculativeWorkOutCh

	newRollupTxs := speculativeRollup.txs
	newRollupState := *speculativeRollup.s

	// the speculative execution has been processing on top of the wrong parent - due to failure in gossip or publishing to L1
	if (speculativeRollup.r == nil) || (speculativeRollup.r.Hash() != bs.head.Hash()) {
		if speculativeRollup.r != nil {
			log.Log(fmt.Sprintf(">   Agg%d: Recalculate. speculative=r_%d(%d), published=r_%d(%d)",
				obscurocommon.ShortAddress(e.node),
				obscurocommon.ShortHash(speculativeRollup.r.Hash()),
				speculativeRollup.r.Header.Height,
				obscurocommon.ShortHash(bs.head.Hash()),
				bs.head.Header.Height),
			)
			if e.statsCollector != nil {
				e.statsCollector.L2Recalc(e.node)
			}
		}

		// determine transactions to include in new rollup and process them
		newRollupTxs = currentTxs(bs.head, e.storage.FetchMempoolTxs(), e.storage)
		newRollupState = executeTransactions(newRollupTxs, newProcessedState(bs.state))
	}

	// always process deposits last
	// process deposits from the proof of the parent to the current block (which is the proof of the new rollup)
	proof := bs.head.Proof(e.blockResolver)
	newRollupState = processDeposits(proof, b, copyProcessedState(newRollupState), e.blockResolver)

	// Create a new rollup based on the proof of inclusion of the previous, including all new transactions
	r := NewRollup(b, bs.head, bs.head.Header.Height+1, e.node, newRollupTxs, newRollupState.w, obscurocommon.GenerateNonce(), serialize(newRollupState.s))
	return &r
}

func (e *enclaveImpl) GetTransaction(txHash common.Hash) *nodecommon.L2Tx {
	// todo add some sort of cache
	rollup := e.storage.FetchHeadState().head

	var found bool
	for {
		txs := rollup.Transactions
		for _, tx := range txs {
			if tx.Hash() == txHash {
				return &tx
			}
		}
		rollup = e.storage.ParentRollup(rollup)
		rollup, found = e.storage.FetchRollup(rollup.Hash())
		if !found {
			panic(fmt.Sprintf("Could not find rollup: r_%s", rollup.Hash()))
		}
		if rollup.Header.Height == obscurocommon.L2GenesisHeight {
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
	txs []nodecommon.L2Tx
}

func NewEnclave(id common.Address, mining bool, collector StatsCollector) nodecommon.Enclave {
	storage := NewStorage()
	return &enclaveImpl{
		node:                 id,
		storage:              storage,
		blockResolver:        storage,
		mining:               mining,
		txCh:                 make(chan nodecommon.L2Tx),
		roundWinnerCh:        make(chan *Rollup),
		exitCh:               make(chan bool),
		speculativeWorkInCh:  make(chan bool),
		speculativeWorkOutCh: make(chan speculativeWork),
		statsCollector:       collector,
	}
}
