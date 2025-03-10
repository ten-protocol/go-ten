package enclave

import (
	"fmt"
	"math/big"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethutil"
)

// This state machine compares the state of the enclave to the state of the world and is used to determine what actions can be taken with the enclave.
// It records the last known status code of the enclave. It also records the l1 head and the l2 head that it believes the
// enclave has processed, optimistically updating these after successful actions and verifying the status when errors occur.

// Usage notes:
// - The status is updated by the host when the enclave successfully processed blocks and batches
// - The status is updated when we receive a status from the enclave
// - The status is **not** updated immediately when it receives blocks/batches from the outside world (this is to avoid flickering between catch-up and live when new blocks/batches arrive)
// - The state should be notified of a live block/batch arrival before notifying if it is successfully processed
// - If unexpected error occurs when interacting with the enclave, then status should be requested and this state updated with the result

// Status is the status of the enclave from the host's perspective (including what it knows of the outside world)
type Status int

const (
	// Live - enclave is up-to-date with known external data. It can process L1 and L2 blocks as they arrive and respond to requests.
	Live Status = iota
	// Disconnected - enclave is unreachable or not returning a valid status (this overrides state calculations)
	Disconnected
	// Unavailable - enclave responding with 'Unavailable' status code
	Unavailable
	// AwaitingSecret - enclave is waiting for host to request and provide secret
	AwaitingSecret
	// L1Catchup - enclave is behind on L1 data, host should submit L1 blocks to catch up
	L1Catchup
	// L2Catchup - enclave is behind on L2 data, host should request and submit L2 batches to catch up
	L2Catchup
)

// when the L2 head is 0 then it means no batch has been seen or processed (first seq number is always 1)
var _noBatch = big.NewInt(0)

func (es Status) String() string {
	return [...]string{"Live", "Disconnected", "Unavailable", "AwaitingSecret", "L1Catchup", "L2Catchup"}[es]
}

// StateTracker is the state machine for the enclave
type StateTracker struct {
	// status is the status according to this enclave tracker
	// It is a function of the properties below and recalculated when any of them change
	status Status

	// enclave states (updated when enclave returns Status and optimistically after successful actions)
	enclaveStatusCode common.StatusCode // this is the status code reported by the enclave (Running/AwaitingSecret/Unavailable)
	enclaveL1Head     gethcommon.Hash
	enclaveL2Head     *big.Int
	isActiveSequencer bool

	// latest seen heads of L1 and L2 chains from external sources
	hostL1Head gethcommon.Hash
	hostL2Head *big.Int

	m      *sync.RWMutex
	logger gethlog.Logger
}

// StateSnapshot is a snapshot of the state of the enclave, used to ensure a consistent, thread-safe view of the state
type StateSnapshot struct {
	Status  Status
	Enclave struct {
		StatusCode        common.StatusCode
		L1Head            gethcommon.Hash
		L2Head            *big.Int
		IsActiveSequencer bool
	}
	Host struct {
		L1Head gethcommon.Hash
		L2Head *big.Int
	}
}

// InSyncWithL1 returns true if the enclave is up-to-date with L1 data so guardian can process L1 blocks as they arrive
func (s StateSnapshot) InSyncWithL1() bool {
	return s.Status == Live || s.Status == L2Catchup
}

func (s StateSnapshot) IsLive() bool {
	return s.Status == Live
}

func (s StateSnapshot) String() string {
	return fmt.Sprintf("StateSnapshot: [%s] enclave(StatusCode=%d, L1Head=%s, L2Head=%s IsActive=%v), Host(L1Head=%s, L2Head=%s)",
		s.Status, s.Enclave.StatusCode, s.Enclave.L1Head, s.Enclave.L2Head, s.Enclave.IsActiveSequencer, s.Host.L1Head, s.Host.L2Head)
}

func NewStateTracker(logger gethlog.Logger) *StateTracker {
	return &StateTracker{status: Disconnected, m: &sync.RWMutex{}, logger: logger}
}

func (s *StateTracker) GetStatus() Status {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.status
}

func (s *StateTracker) Snapshot() StateSnapshot {
	s.m.RLock()
	defer s.m.RUnlock()
	var snap StateSnapshot
	snap.Status = s.status
	snap.Enclave.StatusCode = s.enclaveStatusCode
	snap.Enclave.L1Head = s.enclaveL1Head
	snap.Enclave.L2Head = s.enclaveL2Head
	snap.Enclave.IsActiveSequencer = s.isActiveSequencer
	snap.Host.L1Head = s.hostL1Head
	snap.Host.L2Head = s.hostL2Head
	return snap
}

func (s *StateTracker) OnProcessedBlock(enclL1Head gethcommon.Hash) {
	s.m.Lock()
	defer s.m.Unlock()
	s.enclaveL1Head = enclL1Head
	s.setStatus("onProcessedBlock", s.calculateStatus())
}

func (s *StateTracker) OnReceivedBlock(l1Head gethcommon.Hash) {
	s.m.Lock()
	defer s.m.Unlock()
	s.hostL1Head = l1Head
}

func (s *StateTracker) OnProcessedBatch(enclL2HeadSeqNo *big.Int) {
	s.m.Lock()
	defer s.m.Unlock()
	s.enclaveL2Head = enclL2HeadSeqNo
	s.setStatus("onProcessedBatch", s.calculateStatus())
}

func (s *StateTracker) OnReceivedBatch(l2HeadSeqNo *big.Int) {
	s.m.Lock()
	defer s.m.Unlock()
	s.hostL2Head = l2HeadSeqNo
}

func (s *StateTracker) OnSecretProvided() {
	s.m.Lock()
	defer s.m.Unlock()
	if s.enclaveStatusCode == common.AwaitingSecret {
		s.enclaveStatusCode = common.Running
	}
	s.setStatus("onSecretProvided", s.calculateStatus())
}

func (s *StateTracker) OnEnclaveStatus(es common.Status) {
	s.m.Lock()
	defer s.m.Unlock()
	s.enclaveStatusCode = es.StatusCode
	// only update L1 head if non-empty head reported
	if es.L1Head != gethutil.EmptyHash {
		s.enclaveL1Head = es.L1Head
	}
	s.enclaveL2Head = es.L2Head
	if s.isActiveSequencer != es.IsActiveSequencer {
		if es.IsActiveSequencer {
			s.logger.Info("Enclave is now active sequencer")
		} else {
			s.logger.Info("Enclave is no longer active sequencer")
		}
	}
	s.isActiveSequencer = es.IsActiveSequencer
	s.setStatus("onEnclaveStatus", s.calculateStatus())
}

// OnDisconnected is called if the enclave is unreachable/not returning a valid Status
func (s *StateTracker) OnDisconnected() {
	s.m.Lock()
	defer s.m.Unlock()
	s.setStatus("onDisconnect", Disconnected)
}

// when enclave is operational, this method will calculate the status based on comparison of current chain heads with enclave heads
// for consistency, this should be called within a lock
func (s *StateTracker) calculateStatus() Status {
	switch s.enclaveStatusCode {
	case common.AwaitingSecret:
		return AwaitingSecret
	case common.Unavailable:
		return Unavailable
	case common.Running:
		if s.hostL1Head != s.enclaveL1Head || s.enclaveL1Head == gethutil.EmptyHash {
			return L1Catchup
		}
		if s.hostL2Head == nil || s.enclaveL2Head == nil || s.enclaveL2Head.Cmp(_noBatch) == 0 || s.hostL2Head.Cmp(s.enclaveL2Head) > 0 {
			return L2Catchup
		}
		return Live
	default:
		// this shouldn't happen
		s.logger.Error("Unknown enclave status code - this should not happen", "code", s.enclaveStatusCode)
		return Unavailable
	}
}

// this must be called from within write-lock
func (s *StateTracker) setStatus(event string, newStatus Status) {
	if s.status == newStatus {
		return
	}
	s.logger.Info(fmt.Sprintf("Updating enclave status from [%s] to [%s] following %s", s.status, newStatus, event), "state", s)
	s.status = newStatus
}

func (s *StateTracker) SequencerPromoted() {
	s.m.Lock()
	defer s.m.Unlock()
	s.isActiveSequencer = true
}
