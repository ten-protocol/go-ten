package enclave

import (
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/ten-protocol/go-ten/go/common/log"
)

var (
	_l1Block123        = gethcommon.BytesToHash([]byte{1, 2, 3})
	_l1Block124        = gethcommon.BytesToHash([]byte{1, 2, 4})
	_l2Batch456        = big.NewInt(456)
	_l2Batch456Hash    = gethcommon.BytesToHash([]byte{4, 5, 6})
	_l2Batch457        = big.NewInt(457)
	_l2Batch457Hash    = gethcommon.BytesToHash([]byte{4, 5, 7})
	stateTrackerLogger = log.New("stateTracker", int(gethlog.LevelWarn), log.SysOut)
)

func TestStateTracker_InSyncWithL1(t *testing.T) {
	s := NewStateTracker(stateTrackerLogger)
	// state tracker is up-to-date with L1
	s.OnReceivedBlock(_l1Block123)
	s.OnProcessedBlock(_l1Block123)
	// but not processed L2 yet
	s.OnReceivedBatch(_l2Batch456, _l2Batch456Hash)
	assert.Equal(t, L2Catchup, s.GetStatus())
}

func TestStateTracker_InSyncWithL2(t *testing.T) {
	s := NewStateTracker(stateTrackerLogger)
	// state tracker is up-to-date with L1
	s.OnReceivedBlock(_l1Block123)
	s.OnProcessedBlock(_l1Block123)
	// state tracker is up-to-date with L2
	s.OnReceivedBatch(_l2Batch456, _l2Batch456Hash)
	s.OnProcessedBatch(_l2Batch456, _l2Batch456Hash)
	assert.Equal(t, Live, s.GetStatus())
}

func TestStateTracker_InSyncL2ButBehindL1(t *testing.T) {
	s := NewStateTracker(stateTrackerLogger)
	// block 124 is received before block 123 is processed
	s.OnReceivedBlock(_l1Block124)
	// state tracker becomes aware that it is behind (it works this way to avoid flickering to catch-up every time a block arrives)
	s.OnProcessedBlock(_l1Block123)
	// batches are up-to-date
	s.OnReceivedBatch(_l2Batch456, _l2Batch456Hash)
	s.OnProcessedBatch(_l2Batch456, _l2Batch456Hash)
	assert.Equal(t, L1Catchup, s.GetStatus())
}

func TestStateTracker_InSyncWithL1ButBehindL2(t *testing.T) {
	s := NewStateTracker(stateTrackerLogger)
	s.OnReceivedBlock(_l1Block123)
	s.OnProcessedBlock(_l1Block123)
	// batch 457 is received before batch 456 is processed
	s.OnReceivedBatch(_l2Batch457, _l2Batch457Hash)
	// state tracker becomes aware that it is behind (it works this way to avoid flickering to catch-up every time a batch arrives)
	s.OnProcessedBatch(_l2Batch456, _l2Batch456Hash)
	assert.Equal(t, L2Catchup, s.GetStatus())
}

func TestStateTracker_Disconnected(t *testing.T) {
	s := NewStateTracker(stateTrackerLogger)
	// state tracker is up-to-date with L1 and L2
	s.OnReceivedBlock(_l1Block123)
	s.OnProcessedBlock(_l1Block123)
	s.OnReceivedBatch(_l2Batch456, _l2Batch456Hash)
	s.OnProcessedBatch(_l2Batch456, _l2Batch456Hash)
	// but it gets disconnected
	s.OnDisconnected()
	assert.Equal(t, Disconnected, s.GetStatus())
}
