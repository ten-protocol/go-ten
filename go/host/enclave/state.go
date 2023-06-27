package enclave

import (
	"fmt"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/gethutil"
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

func (es Status) String() string {
	return [...]string{"Live", "Disconnected", "Unavailable", "AwaitingSecret", "L1Catchup", "L2Catchup"}[es]
}

// StateTracker is the state machine for the enclave
type StateTracker struct {
	// status is the cached status of the enclave
	// It is a function of the properties below and recalculated when any of them change
	status Status

	// enclave states (updated when enclave returns Status and optimistically after successful actions)
	enclaveStatusCode common.StatusCode
	enclaveL1Head     gethcommon.Hash
	enclaveL2Head     gethcommon.Hash

	// latest seen heads of L1 and L2 chains from external sources
	hostL1Head gethcommon.Hash
	hostL2Head gethcommon.Hash

	m      *sync.RWMutex
	logger gethlog.Logger
}

func NewStateTracker(logger gethlog.Logger) *StateTracker {
	return &StateTracker{status: Disconnected, m: &sync.RWMutex{}, logger: logger}
}

func (s *StateTracker) String() string {
	return fmt.Sprintf("StateTracker: [%s] enclave(StatusCode=%d, L1Head=%s, L2Head=%s), Host(L1Head=%s, L2Head=%s)",
		s.status, s.enclaveStatusCode, s.enclaveL1Head, s.enclaveL2Head, s.hostL1Head, s.hostL2Head)
}

func (s *StateTracker) GetStatus() Status {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.status
}

func (s *StateTracker) OnProcessedBlock(enclL1Head gethcommon.Hash) {
	s.m.Lock()
	s.enclaveL1Head = enclL1Head
	s.m.Unlock()
	s.recalculateStatus()
}

func (s *StateTracker) OnReceivedBlock(l1Head gethcommon.Hash) {
	s.m.Lock()
	defer s.m.Unlock()
	s.hostL1Head = l1Head
}

func (s *StateTracker) OnProcessedBatch(enclL2Head gethcommon.Hash) {
	s.m.Lock()
	s.enclaveL2Head = enclL2Head
	s.m.Unlock()
	s.recalculateStatus()
}

func (s *StateTracker) OnReceivedBatch(l2Head gethcommon.Hash) {
	s.m.Lock()
	defer s.m.Unlock()
	s.hostL2Head = l2Head
}

func (s *StateTracker) OnSecretProvided() {
	s.m.Lock()
	if s.enclaveStatusCode == common.AwaitingSecret {
		s.enclaveStatusCode = common.Running
	}
	s.m.Unlock()
	s.recalculateStatus()
}

func (s *StateTracker) OnEnclaveStatus(es common.Status) {
	s.m.Lock()
	s.enclaveStatusCode = es.StatusCode
	s.enclaveL1Head = es.L1Head
	s.enclaveL2Head = es.L2Head
	s.m.Unlock()

	s.recalculateStatus()
}

// OnDisconnected is called if the enclave is unreachable/not returning a valid Status
func (s *StateTracker) OnDisconnected() {
	s.m.Lock()
	defer s.m.Unlock()
	s.setStatus(Disconnected)
}

// when enclave is operational, this method will update the status based on comparison of current chain heads with enclave heads
func (s *StateTracker) recalculateStatus() {
	s.m.Lock()
	defer s.m.Unlock()
	switch s.enclaveStatusCode {
	case common.AwaitingSecret:
		s.setStatus(AwaitingSecret)
	case common.Unavailable:
		s.setStatus(Unavailable)
	case common.Running:
		if s.hostL1Head != s.enclaveL1Head || s.enclaveL1Head == gethutil.EmptyHash {
			s.setStatus(L1Catchup)
			return
		}
		if s.hostL2Head != s.enclaveL2Head || s.enclaveL2Head == gethutil.EmptyHash {
			s.setStatus(L2Catchup)
			return
		}
		s.setStatus(Live)
	}
}

// InSyncWithL1 returns true if the enclave is up-to-date with L1 data so guardian can process L1 blocks as they arrive
func (s *StateTracker) InSyncWithL1() bool {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.status == Live || s.status == L2Catchup
}

func (s *StateTracker) IsUpToDate() bool {
	return s.status == Live
}

func (s *StateTracker) GetEnclaveL1Head() gethcommon.Hash {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.enclaveL1Head
}

func (s *StateTracker) GetEnclaveL2Head() gethcommon.Hash {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.enclaveL2Head
}

// this must be called from within write-lock
func (s *StateTracker) setStatus(newStatus Status) {
	if s.status == newStatus {
		return
	}
	s.logger.Debug(fmt.Sprintf("Updating enclave status from [%s] to [%s]", s.status, newStatus))
	s.status = newStatus
}
