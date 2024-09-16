package ethereummock

//import (
//	"encoding/binary"
//	"github.com/ten-protocol/go-ten/go/ethadapter"
//	"github.com/ten-protocol/go-ten/integration/datagenerator"
//	"k8s.io/utils/clock"
//	"math/big"
//	"time"
//
//	"github.com/ethereum/go-ethereum"
//	"github.com/ethereum/go-ethereum/beacon/engine"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/ethereum/go-ethereum/eth"
//	"github.com/ethereum/go-ethereum/eth/catalyst"
//	"github.com/ethereum/go-ethereum/event"
//	"github.com/ethereum/go-ethereum/log"
//)
//
//type Beacon interface {
//	StoreBlobsBundle(slot uint64, bundle *engine.BlobsBundleV1) error
//}
//
//// fakePoS is a testing-only utility to attach to Geth,
//// to build a fake proof-of-stake L1 chain with fixed block time and basic lagging safe/finalized blocks.
//type mockPoS struct {
//	clock     clock.Clock
//	eth       *eth.Ethereum
//	log       log.Logger
//	blockTime uint64
//
//	withdrawalsIndex uint64
//
//	finalizedDistance uint64
//	safeDistance      uint64
//
//	engineAPI *catalyst.ConsensusAPI
//	sub       ethereum.Subscription
//
//	beacon Beacon
//}
//
//func (f *mockPoS) FakeBeaconBlockRoot(time uint64) common.Hash {
//	var dat [8]byte
//	binary.LittleEndian.PutUint64(dat[:], time)
//	return crypto.Keccak256Hash(dat[:])
//}
//
//func (f *mockPoS) Start() error {
//	withdrawalsRNG := datagenerator.RandomUInt64()
//	f.sub = event.NewSubscription(func(quit <-chan struct{}) error {
//		ticker := time.NewTicker(time.Second / 2)
//		defer ticker.Stop()
//
//		for {
//			select {
//			case now := <-ticker.C:
//				chain := f.eth.BlockChain()
//				head := chain.CurrentBlock()
//				finalized := chain.CurrentFinalBlock()
//				if finalized == nil {
//					finalized = chain.Genesis().Header()
//				}
//				safe := chain.CurrentSafeBlock()
//				if safe == nil {
//					safe = finalized
//				}
//				if head.Number.Uint64() > f.finalizedDistance {
//					finalized = f.eth.BlockChain().GetHeaderByNumber(head.Number.Uint64() - f.finalizedDistance)
//				}
//				if head.Number.Uint64() > f.safeDistance {
//					safe = f.eth.BlockChain().GetHeaderByNumber(head.Number.Uint64() - f.safeDistance)
//				}
//				if head.Time >= uint64(now.Unix()) {
//					continue
//				}
//				newBlockTime := head.Time + f.blockTime
//				if time.Unix(int64(newBlockTime), 0).Add(5 * time.Minute).Before(time.Now()) {
//					newBlockTime = uint64(time.Now().Unix())
//				}
//				// create some random withdrawals
//				withdrawals := make([]*types.Withdrawal, withdrawalsRNG)
//				for i := 0; i < len(withdrawals); i++ {
//					withdrawals[i] = &types.Withdrawal{
//						Index:     f.withdrawalsIndex + uint64(i),
//						Validator: withdrawalsRNG % 100_000_000,
//						Address:   datagenerator.RandomAddress(),
//						Amount:    uint64(50_000_000_000),
//					}
//				}
//				attrs := &engine.PayloadAttributes{
//					Timestamp:             newBlockTime,
//					Random:                common.Hash{},
//					SuggestedFeeRecipient: head.Coinbase,
//					Withdrawals:           withdrawals,
//				}
//				parentBeaconBlockRoot := f.FakeBeaconBlockRoot(head.Time)
//				isCancun := f.eth.BlockChain().Config().IsCancun(new(big.Int).SetUint64(head.Number.Uint64()+1), newBlockTime)
//				if isCancun {
//					attrs.BeaconRoot = &parentBeaconBlockRoot
//				}
//				fcState := engine.ForkchoiceStateV1{
//					HeadBlockHash:      head.Hash(),
//					SafeBlockHash:      safe.Hash(),
//					FinalizedBlockHash: finalized.Hash(),
//				}
//				var err error
//				var res engine.ForkChoiceResponse
//				if isCancun {
//					res, err = f.engineAPI.ForkchoiceUpdatedV3(fcState, attrs)
//				} else {
//					res, err = f.engineAPI.ForkchoiceUpdatedV2(fcState, attrs)
//				}
//				if err != nil {
//					f.log.Error("failed to start building L1 block", "err", err)
//					continue
//				}
//				if res.PayloadID == nil {
//					f.log.Error("failed to start block building", "res", res)
//					continue
//				}
//				// wait with sealing, if we are not behind already
//				delay := time.Until(time.Unix(int64(newBlockTime), 0))
//				if delay > 0 {
//					timer := time.NewTimer(delay)
//					select {
//					case <-timer.C:
//						// no-op
//					case <-quit:
//						timer.Stop()
//						return nil
//					}
//				}
//				envelope, err := f.engineAPI.GetPayloadV3(*res.PayloadID)
//				if err != nil {
//					f.log.Error("failed to finish building L1 block", "err", err)
//					continue
//				}
//
//				blobHashes := make([]common.Hash, 0)
//				for _, commitment := range envelope.BlobsBundle.Commitments {
//					if len(commitment) != 48 {
//						f.log.Error("got malformed kzg commitment from engine", "commitment", commitment)
//						break
//					}
//					blobHashes = append(blobHashes, ethadapter.KZGToVersionedHash(*(*[48]byte)(commitment)))
//				}
//				if len(blobHashes) != len(envelope.BlobsBundle.Commitments) {
//					f.log.Error("invalid or incomplete blob data", "collected", len(blobHashes), "engine", len(envelope.BlobsBundle.Commitments))
//					continue
//				}
//				if isCancun {
//					if _, err := f.engineAPI.NewPayloadV3(*envelope.ExecutionPayload, blobHashes, &parentBeaconBlockRoot); err != nil {
//						f.log.Error("failed to insert built L1 block", "err", err)
//						continue
//					}
//				} else {
//					if _, err := f.engineAPI.NewPayloadV2(*envelope.ExecutionPayload); err != nil {
//						f.log.Error("failed to insert built L1 block", "err", err)
//						continue
//					}
//				}
//				if envelope.BlobsBundle != nil {
//					slot := (envelope.ExecutionPayload.Timestamp - f.eth.BlockChain().Genesis().Time()) / f.blockTime
//					if f.beacon == nil {
//						f.log.Error("no blobs storage available")
//						continue
//					}
//					if err := f.beacon.StoreBlobsBundle(slot, envelope.BlobsBundle); err != nil {
//						f.log.Error("failed to persist blobs-bundle of block, not making block canonical now", "err", err)
//						continue
//					}
//				}
//				if _, err := f.engineAPI.ForkchoiceUpdatedV3(engine.ForkchoiceStateV1{
//					HeadBlockHash:      envelope.ExecutionPayload.BlockHash,
//					SafeBlockHash:      safe.Hash(),
//					FinalizedBlockHash: finalized.Hash(),
//				}, nil); err != nil {
//					f.log.Error("failed to make built L1 block canonical", "err", err)
//					continue
//				}
//				f.withdrawalsIndex += uint64(len(withdrawals))
//			case <-quit:
//				return nil
//			}
//		}
//	})
//	return nil
//}
//
//func (f *mockPoS) Stop() error {
//	f.sub.Unsubscribe()
//	return nil
//}
