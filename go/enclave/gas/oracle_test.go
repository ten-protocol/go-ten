package gas

import (
	"context"
	"errors"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"

	tencommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/enclave/evm"
)

// mockBlockResolver implements storage.BlockResolver for tests
type mockBlockResolver struct {
	headByHeight map[uint64]*types.Header
	headByHash   map[tencommon.L1BlockHash]*types.Header
}

func newMockBlockResolver() *mockBlockResolver {
	return &mockBlockResolver{
		headByHeight: make(map[uint64]*types.Header),
		headByHash:   make(map[tencommon.L1BlockHash]*types.Header),
	}
}

func (m *mockBlockResolver) addHeader(h *types.Header) {
	hash := h.Hash()
	m.headByHeight[h.Number.Uint64()] = h
	m.headByHash[hash] = h
}

// interface methods
func (m *mockBlockResolver) FetchBlock(ctx context.Context, blockHash tencommon.L1BlockHash) (*types.Header, error) {
	if h, ok := m.headByHash[blockHash]; ok {
		return h, nil
	}
	return nil, errNotFound
}

func (m *mockBlockResolver) IsBlockCanonical(ctx context.Context, blockHash tencommon.L1BlockHash) (bool, error) {
	_, ok := m.headByHash[blockHash]
	return ok, nil
}

func (m *mockBlockResolver) FetchCanonicaBlockByHeight(ctx context.Context, height *big.Int) (*types.Header, error) {
	if h, ok := m.headByHeight[height.Uint64()]; ok {
		return h, nil
	}
	return nil, errNotFound
}

func (m *mockBlockResolver) FetchHeadBlock(ctx context.Context) (*types.Header, error) {
	var max uint64
	for k := range m.headByHeight {
		if k > max {
			max = k
		}
	}
	if max == 0 {
		return nil, errNotFound
	}
	return m.headByHeight[max], nil
}

func (m *mockBlockResolver) StoreBlock(ctx context.Context, block *types.Header, fork *tencommon.ChainFork) error {
	return nil
}
func (m *mockBlockResolver) UpdateProcessed(ctx context.Context, block tencommon.L1BlockHash) error {
	return nil
}
func (m *mockBlockResolver) DeleteDirtyBlocks(ctx context.Context) error { return nil }
func (m *mockBlockResolver) IsAncestor(ctx context.Context, block *types.Header, maybeAncestor *types.Header) bool {
	return false
}

var errNotFound = errors.New("not found")

// helper to construct a simple chain of headers with specified base fees
func buildChainWithBaseFees(fees []int64) (*types.Header, *mockBlockResolver) {
	resolver := newMockBlockResolver()
	var parent *types.Header
	for i, fee := range fees {
		h := &types.Header{
			Number:  big.NewInt(int64(i + 1)),
			BaseFee: big.NewInt(fee),
		}
		// Set ExcessBlobGas non-nil to avoid blob fee calculation path in MA
		zero := uint64(0)
		h.ExcessBlobGas = &zero
		if parent != nil {
			h.ParentHash = parent.Hash()
		}
		resolver.addHeader(h)
		parent = h
	}
	return parent, resolver
}

func averageInt64(values []int64) *big.Int {
	if len(values) == 0 {
		return big.NewInt(0)
	}
	sum := big.NewInt(0)
	for _, v := range values {
		sum.Add(sum, big.NewInt(v))
	}
	return sum.Div(sum, big.NewInt(int64(len(values))))
}

func roundUpToMultiple(n, multiple *big.Int) *big.Int {
	if multiple.Sign() == 0 {
		return new(big.Int).Set(n)
	}
	remainder := new(big.Int).Mod(n, multiple)
	if remainder.Sign() == 0 {
		return new(big.Int).Set(n)
	}
	return new(big.Int).Add(n, new(big.Int).Sub(multiple, remainder))
}

func TestEstimateL1StorageGasCost_UsesMAAndRounds_IncludingBlobShare(t *testing.T) {
	ctx := context.Background()

	fees := []int64{100, 200, 300}
	head, resolver := buildChainWithBaseFees(fees)

	oracleIface := NewGasOracle(params.MainnetChainConfig, resolver, gethlog.New())
	impl := oracleIface.(*oracle)

	to := gethcommon.HexToAddress("0x0000000000000000000000000000000000000001")
	tx := types.NewTx(&types.LegacyTx{To: &to, Gas: 21000, GasPrice: big.NewInt(1), Value: big.NewInt(0), Data: []byte{0x01, 0x02}})
	batch := &tencommon.BatchHeader{BaseFee: big.NewInt(1), GasLimit: 10_000_000}

	got, err := oracleIface.EstimateL1StorageGasCost(ctx, tx, head, batch)
	if err != nil {
		t.Fatalf("EstimateL1StorageGasCost returned error: %v", err)
	}

	// Compute expected using the same MA and encoded size
	baseMA, blobMA, err := impl.calculateMA(ctx, head.Number.Uint64())
	if err != nil {
		t.Fatalf("calculateMA error: %v", err)
	}
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		t.Fatalf("rlp encode: %v", err)
	}
	txL1Size := CalculateL1Size(encodedTx)

	shareOfBlob := big.NewInt(0)
	if blobMA.Sign() > 0 {
		shareOfBlob = new(big.Int).Mul(txL1Size, blobMA)
	}
	shareOfL1TxGas := big.NewInt(L1TxGas / TxsPerRollup)
	shareOfL1 := new(big.Int).Mul(shareOfL1TxGas, baseMA)
	expected := new(big.Int).Add(shareOfBlob, shareOfL1)
	expected = roundUpToMultiple(expected, evm.FIXED_L2_GAS_COST_FOR_L1_PUBLISHING)

	if got.Cmp(expected) != 0 {
		t.Fatalf("unexpected l1 storage gas cost: got %s want %s", got.String(), expected.String())
	}
}

func TestEstimateL1CostForMsg_UsesHeadBlockAndRounds(t *testing.T) {
	ctx := context.Background()

	fees := []int64{500, 700, 900}
	head, resolver := buildChainWithBaseFees(fees)
	oracleIface := NewGasOracle(params.MainnetChainConfig, resolver, gethlog.New())
	impl := oracleIface.(*oracle)

	if err := oracleIface.SubmitL1Block(ctx, head); err != nil {
		t.Fatalf("SubmitL1Block error: %v", err)
	}

	from := gethcommon.HexToAddress("0x0000000000000000000000000000000000000002")
	args := &gethapi.TransactionArgs{From: &from}
	batch := &tencommon.BatchHeader{BaseFee: big.NewInt(1), GasLimit: 10_000_000}

	got, err := oracleIface.EstimateL1CostForMsg(ctx, args, batch)
	if err != nil {
		t.Fatalf("EstimateL1CostForMsg returned error: %v", err)
	}

	baseMA, blobMA, err := impl.calculateMA(ctx, head.Number.Uint64())
	if err != nil {
		t.Fatalf("calculateMA error: %v", err)
	}
	encodedArgs, err := rlp.EncodeToBytes(args)
	if err != nil {
		t.Fatalf("rlp encode: %v", err)
	}
	msgL1Size := CalculateL1Size(encodedArgs)

	shareOfBlob := big.NewInt(0)
	if blobMA.Sign() > 0 {
		shareOfBlob = new(big.Int).Mul(msgL1Size, blobMA)
	}
	shareOfL1TxGas := big.NewInt(L1TxGas / TxsPerRollup)
	shareOfL1 := new(big.Int).Mul(shareOfL1TxGas, baseMA)
	expected := new(big.Int).Add(shareOfBlob, shareOfL1)
	expected = roundUpToMultiple(expected, evm.FIXED_L2_GAS_COST_FOR_L1_PUBLISHING)

	if got.Cmp(expected) != 0 {
		t.Fatalf("unexpected l1 cost for msg: got %s want %s", got.String(), expected.String())
	}
}

func TestCalculateL1Size_CompressionFactor(t *testing.T) {
	data := make([]byte, 100)
	size := CalculateL1Size(data)
	expected := uint64((100 * compressionFactor) / 100)
	if size.Uint64() != expected {
		t.Fatalf("unexpected compressed size: got %d want %d", size.Uint64(), expected)
	}
}

func TestEstimateL1StorageGasCost_IncreasesWithBaseFeeAndDecreasesOnDrop(t *testing.T) {
	ctx := context.Background()

	// Build minimal headers with distinct numbers
	h1 := &types.Header{Number: big.NewInt(1)}
	h2 := &types.Header{Number: big.NewInt(2)}
	h3 := &types.Header{Number: big.NewInt(3)}

	resolver := newMockBlockResolver()
	resolver.addHeader(h1)
	resolver.addHeader(h2)
	resolver.addHeader(h3)

	oracleIface := NewGasOracle(params.MainnetChainConfig, resolver, gethlog.New())
	impl := oracleIface.(*oracle)

	// Keep blob fee at 0 to isolate base fee effect; ensure blob cache hit
	impl.blobFeeMA[1] = big.NewInt(0)
	impl.blobFeeMA[2] = big.NewInt(0)
	impl.blobFeeMA[3] = big.NewInt(0)

	// Set base fee MA values (all well above the fixed cost scale)
	bfLow := big.NewInt(1_000)
	bfHigh := big.NewInt(2_000)
	bfMid := big.NewInt(1_500)
	impl.baseFeeMA[1] = bfLow
	impl.baseFeeMA[2] = bfHigh
	impl.baseFeeMA[3] = bfMid

	to := gethcommon.HexToAddress("0x0000000000000000000000000000000000000001")
	tx := types.NewTx(&types.LegacyTx{To: &to, Gas: 21000, GasPrice: big.NewInt(1), Value: big.NewInt(0), Data: []byte{0x01, 0x02}})
	batch := &tencommon.BatchHeader{BaseFee: big.NewInt(1), GasLimit: 10_000_000}

	costLow, err := oracleIface.EstimateL1StorageGasCost(ctx, tx, h1, batch)
	if err != nil {
		t.Fatalf("low base fee estimate error: %v", err)
	}
	costHigh, err := oracleIface.EstimateL1StorageGasCost(ctx, tx, h2, batch)
	if err != nil {
		t.Fatalf("high base fee estimate error: %v", err)
	}
	costMid, err := oracleIface.EstimateL1StorageGasCost(ctx, tx, h3, batch)
	if err != nil {
		t.Fatalf("mid base fee estimate error: %v", err)
	}

	if costHigh.Cmp(costLow) <= 0 {
		t.Fatalf("expected cost to increase with higher base fee: low=%s high=%s", costLow, costHigh)
	}
	if costMid.Cmp(costHigh) >= 0 {
		t.Fatalf("expected cost to decrease when base fee drops: mid=%s high=%s", costMid, costHigh)
	}
}

func TestEstimateL1StorageGasCost_IncreasesWithBlobFeeAndDecreasesOnDrop(t *testing.T) {
	ctx := context.Background()

	// Build minimal headers with distinct numbers
	h1 := &types.Header{Number: big.NewInt(10)}
	h2 := &types.Header{Number: big.NewInt(11)}
	h3 := &types.Header{Number: big.NewInt(12)}

	resolver := newMockBlockResolver()
	resolver.addHeader(h1)
	resolver.addHeader(h2)
	resolver.addHeader(h3)

	oracleIface := NewGasOracle(params.MainnetChainConfig, resolver, gethlog.New())
	impl := oracleIface.(*oracle)

	// Fix base fee to isolate blob fee effect
	impl.baseFeeMA[10] = big.NewInt(1_000)
	impl.baseFeeMA[11] = big.NewInt(1_000)
	impl.baseFeeMA[12] = big.NewInt(1_000)

	// Vary blob fee MA; well above fixed cost when multiplied by size
	blobLow := big.NewInt(1_000)
	blobHigh := big.NewInt(2_000)
	blobMid := big.NewInt(1_500)
	impl.blobFeeMA[10] = blobLow
	impl.blobFeeMA[11] = blobHigh
	impl.blobFeeMA[12] = blobMid

	to := gethcommon.HexToAddress("0x0000000000000000000000000000000000000001")
	tx := types.NewTx(&types.LegacyTx{To: &to, Gas: 21000, GasPrice: big.NewInt(1), Value: big.NewInt(0), Data: []byte{0x01, 0x02, 0x03, 0x04}})
	batch := &tencommon.BatchHeader{BaseFee: big.NewInt(1), GasLimit: 10_000_000}

	costLow, err := oracleIface.EstimateL1StorageGasCost(ctx, tx, h1, batch)
	if err != nil {
		t.Fatalf("low blob fee estimate error: %v", err)
	}
	costHigh, err := oracleIface.EstimateL1StorageGasCost(ctx, tx, h2, batch)
	if err != nil {
		t.Fatalf("high blob fee estimate error: %v", err)
	}
	costMid, err := oracleIface.EstimateL1StorageGasCost(ctx, tx, h3, batch)
	if err != nil {
		t.Fatalf("mid blob fee estimate error: %v", err)
	}

	if costHigh.Cmp(costLow) <= 0 {
		t.Fatalf("expected cost to increase with higher blob fee: low=%s high=%s", costLow, costHigh)
	}
	if costMid.Cmp(costHigh) >= 0 {
		t.Fatalf("expected cost to decrease when blob fee drops: mid=%s high=%s", costMid, costHigh)
	}
}
