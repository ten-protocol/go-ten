package ethadapter

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/stretchr/testify/mock"
)

func TestBlockProviderHappyPath_LiveStream(t *testing.T) {
	mockEthClient := mockEthClient(3)
	blockProvider, ctxCancel := setupBlockProvider(mockEthClient)
	defer ctxCancel()

	blkStream, err := blockProvider.StartStreamingFromHeight(big.NewInt(3))
	if err != nil {
		t.Error(err)
	}
	blkCount := 0

	for blkCount < 3 {
		select {
		case blk := <-blkStream:
			if blk != nil {
				blkCount++
			}

		case <-time.After(3 * time.Second): // shouldn't have >1sec delay between blocks in this test
			t.Errorf("expected 3 blocks from stream but got %d", blkCount)
		}
	}
}

func TestBlockProviderHappyPath_HistoricThenStream(t *testing.T) {
	mockEthClient := mockEthClient(3)
	blockProvider, ctxCancel := setupBlockProvider(mockEthClient)
	defer ctxCancel()

	blkStream, err := blockProvider.StartStreamingFromHeight(big.NewInt(1))
	if err != nil {
		t.Error(err)
	}
	blkCount := 0

	for blkCount < 3 {
		select {
		case blk := <-blkStream:
			if blk != nil {
				blkCount++
			}

		case <-time.After(3 * time.Second): // shouldn't have >1sec delay between blocks in this test
			t.Errorf("expected 3 blocks from stream but got %d", blkCount)
		}
	}
}

func setupBlockProvider(mockEthClient EthClient) (EthBlockProvider, context.CancelFunc) {
	ctx, cancelCtx := context.WithCancel(context.Background())

	logger := log.New(log.HostCmp, int(gethlog.LvlInfo), log.SysOut, log.NodeIDKey, "test")
	blockProvider := EthBlockProvider{
		ethClient:     mockEthClient,
		ctx:           ctx,
		runningStatus: new(int32),
		streamCh:      make(chan *types.Block),
		logger:        logger,
	}
	return blockProvider, cancelCtx
}

func mockEthClient(liveStreamingStart int) EthClient {
	mockClient := &ethClientMock{
		ctx:               context.TODO(),
		blks:              map[gethcommon.Hash]*types.Block{},
		blksByNum:         map[int]*types.Block{},
		liveStreamingNext: liveStreamingStart,
		lastBlockCreation: time.Now(),
	}

	// create the blocks before the streaming portion
	for i := 0; i < liveStreamingStart; i++ {
		mockClient.createHeader(i)
	}

	return mockClient
}

type ethClientMock struct {
	mock.Mock
	ctx               context.Context
	blks              map[gethcommon.Hash]*types.Block
	blksByNum         map[int]*types.Block
	liveStreamingNext int // BlockListener() will stream from this height
	lastBlockCreation time.Time
}

func (e *ethClientMock) createHeader(i int) *types.Block {
	blkHead := &types.Header{
		ParentHash: getHash(i - 1),
		Root:       getHash(i),
		TxHash:     getHash(i),
		Number:     big.NewInt(int64(i)),
	}
	block := types.NewBlock(blkHead, nil, nil, nil, nil)
	e.blks[block.Hash()] = block
	e.blksByNum[i] = block
	e.lastBlockCreation = time.Now()
	return block
}

func (e *ethClientMock) BlockListener() (chan *types.Header, ethereum.Subscription) {
	headChan := make(chan *types.Header)
	subCtx, cancel := context.WithCancel(e.ctx)
	sub := &ethSubscriptionMock{cancel: cancel}
	go func() {
		for {
			select {
			case <-subCtx.Done():
				return
			case <-time.After(800 * time.Millisecond):
				if time.Since(e.lastBlockCreation) < 500*time.Millisecond {
					// not been long enough since last stream, do it next time
					continue
				}
				blk := e.createHeader(e.liveStreamingNext)
				e.liveStreamingNext++
				header := blk.Header()
				headChan <- header
			}
		}
	}()
	return headChan, sub
}

func getHash(i int) gethcommon.Hash {
	return gethcommon.HexToHash(fmt.Sprintf("%d", i))
}

func (e *ethClientMock) BlockByHash(id gethcommon.Hash) (*types.Block, error) {
	if time.Since(e.lastBlockCreation) > 500*time.Millisecond {
		e.createHeader(e.liveStreamingNext)
		e.liveStreamingNext++
	}
	block, f := e.blks[id]
	if !f {
		return nil, fmt.Errorf("block not found")
	}
	return block, nil
}

func (e *ethClientMock) BlockByNumber(num *big.Int) (*types.Block, error) {
	if time.Since(e.lastBlockCreation) > 500*time.Millisecond {
		e.createHeader(e.liveStreamingNext)
		e.liveStreamingNext++
	}
	block, f := e.blksByNum[int(num.Int64())]
	if !f {
		return nil, fmt.Errorf("block not found")
	}
	return block, nil
}

func (e *ethClientMock) SendTransaction(signedTx *types.Transaction) error {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error) {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) Nonce(address gethcommon.Address) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) BalanceAt(account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) Info() Info {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) FetchHeadBlock() (*types.Block, bool) {
	if time.Since(e.lastBlockCreation) > 500*time.Millisecond {
		e.createHeader(e.liveStreamingNext)
		e.liveStreamingNext++
	}
	return e.blksByNum[e.liveStreamingNext-1], true
}

func (e *ethClientMock) BlocksBetween(block *types.Block, head *types.Block) []*types.Block {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) IsBlockAncestor(block *types.Block, proof common.L1RootHash) bool {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) CallContract(msg ethereum.CallMsg) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) Stop() {
	// TODO implement me
	panic("implement me")
}

func (e *ethClientMock) EthClient() *ethclient.Client {
	// TODO implement me
	panic("implement me")
}

type ethSubscriptionMock struct {
	mock.Mock
	cancel context.CancelFunc
}

func (e *ethSubscriptionMock) Unsubscribe() {
	cancel := e.cancel
	cancel()
}

func (e *ethSubscriptionMock) Err() <-chan error {
	return make(<-chan error)
}
