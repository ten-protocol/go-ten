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
	mockEthClient := mockEthClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(log.HostCmp, int(gethlog.LvlInfo), log.SysOut, log.NodeIDKey, "test")

	blockProvider := EthBlockProvider{
		ethClient: mockEthClient,
		ctx:       ctx,
		streamCh:  make(chan *types.Block),
		logger:    logger,
	}

	blockProvider.Start()
	// time.Sleep(20 * time.Millisecond)

	blkStream, err := blockProvider.StartStreamingFromHeight(big.NewInt(0))
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

		case err := <-blockProvider.Err():
			t.Errorf("unexpected error: %s", err)

		case <-time.After(3 * time.Second): // shouldn't have >1sec delay between blocks in this test
			t.Errorf("expected 3 blocks from stream but got %d", blkCount)
		}
	}
}

func TestBlockProviderHappyPath_HistoricThenStream(t *testing.T) {
	mockEthClient := mockEthClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(log.HostCmp, int(gethlog.LvlInfo), log.SysOut, log.NodeIDKey, "test")
	blockProvider := EthBlockProvider{
		ethClient: mockEthClient,
		ctx:       ctx,
		streamCh:  make(chan *types.Block),
		logger:    logger,
	}

	blockProvider.Start()
	time.Sleep(20 * time.Millisecond)

	blkStream, err := blockProvider.StartStreamingFromHeight(big.NewInt(2))
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

		case err := <-blockProvider.Err():
			t.Errorf("unexpected error: %s", err)

		case <-time.After(3 * time.Second): // shouldn't have >1sec delay between blocks in this test
			t.Errorf("expected 3 blocks from stream but got %d", blkCount)
		}
	}
}

func mockEthClient() EthClient {
	return &ethClientMock{
		ctx:       context.TODO(),
		blks:      map[gethcommon.Hash]*types.Block{},
		blksByNum: map[int]*types.Block{},
	}
}

type ethClientMock struct {
	mock.Mock
	ctx       context.Context
	blks      map[gethcommon.Hash]*types.Block
	blksByNum map[int]*types.Block
}

func (r *ethClientMock) BlockListener() (chan *types.Header, ethereum.Subscription) {
	headChan := make(chan *types.Header)
	sub := &ethSubscriptionMock{}
	go func() {
		blkNum := 0
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-time.After(800 * time.Millisecond):
				blkNum++
				blkHead := &types.Header{
					ParentHash: getHash(blkNum - 1),
					Root:       getHash(blkNum),
					TxHash:     getHash(blkNum),
					Number:     big.NewInt(int64(blkNum)),
				}
				block := types.NewBlock(blkHead, nil, nil, nil, nil)
				r.blks[block.Hash()] = block
				r.blksByNum[blkNum] = block
				headChan <- blkHead
			}
		}
	}()
	return headChan, sub
}

func getHash(i int) gethcommon.Hash {
	return gethcommon.HexToHash(fmt.Sprintf("%d", i))
}

func (r *ethClientMock) BlockByHash(id gethcommon.Hash) (*types.Block, error) {
	block, f := r.blks[id]
	if !f {
		return nil, fmt.Errorf("block not found")
	}
	return block, nil
}

func (r *ethClientMock) BlockByNumber(num *big.Int) (*types.Block, error) {
	block, f := r.blksByNum[int(num.Int64())]
	if !f {
		return nil, fmt.Errorf("block not found")
	}
	return block, nil
}

func (r *ethClientMock) SendTransaction(signedTx *types.Transaction) error {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) TransactionReceipt(hash gethcommon.Hash) (*types.Receipt, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) Nonce(address gethcommon.Address) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) BalanceAt(account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) Info() Info {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) FetchHeadBlock() *types.Block {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) BlocksBetween(block *types.Block, head *types.Block) []*types.Block {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) IsBlockAncestor(block *types.Block, proof common.L1RootHash) bool {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) CallContract(msg ethereum.CallMsg) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) Stop() {
	// TODO implement me
	panic("implement me")
}

func (r *ethClientMock) EthClient() *ethclient.Client {
	// TODO implement me
	panic("implement me")
}

type ethSubscriptionMock struct {
	mock.Mock
}

func (e *ethSubscriptionMock) Unsubscribe() {
	// TODO implement me
	panic("implement me")
}

func (e *ethSubscriptionMock) Err() <-chan error {
	return make(<-chan error)
}
