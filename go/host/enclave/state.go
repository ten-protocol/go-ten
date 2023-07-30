package enclave

import (
	"fmt"
	"math/big"
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

// when the L2 head is 0 then it means no batch has been seen or processed (first seq number is always 1)
var _noBatch = big.NewInt(0)

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
	enclaveL2Head     *big.Int

	// latest seen heads of L1 and L2 chains from external sources
	hostL1Head gethcommon.Hash
	hostL2Head *big.Int

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
	defer s.m.Unlock()
	s.enclaveL1Head = enclL1Head
	s.setStatus(s.calculateStatus())
}

func (s *StateTracker) OnReceivedBlock(l1Head gethcommon.Hash) {
	s.m.Lock()
	defer s.m.Unlock()
	s.hostL1Head = l1Head
}

func (s *StateTracker) OnProcessedBatch(enclL2HeadSeqNo *big.Int) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.hostL2Head == nil || s.hostL2Head.Cmp(enclL2HeadSeqNo) < 0 {
		// we've successfully processed this batch, so the host's head should be at least as high as the enclave's (this shouldn't happen, we want it to be visible if it happens)
		s.logger.Trace("host head behind enclave head - updating to match", "hostHead", s.hostL2Head, "enclaveHead", enclL2HeadSeqNo)
		s.hostL2Head = enclL2HeadSeqNo
	}

	if s.enclaveL2Head == nil || enclL2HeadSeqNo == nil {
		fmt.Println("±±±±± head increased from nil", s.enclaveL2Head, "to", enclL2HeadSeqNo)
	} else if big.NewInt(0).Add(s.enclaveL2Head, big.NewInt(1)).Cmp(enclL2HeadSeqNo) < 0 {
		fmt.Println("±±± onProcessedBatch head increased by more than one from", s.enclaveL2Head, "to", enclL2HeadSeqNo)
	}
	s.enclaveL2Head = enclL2HeadSeqNo
	s.setStatus(s.calculateStatus())
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
	s.setStatus(s.calculateStatus())
}

func (s *StateTracker) OnEnclaveStatus(es common.Status) {
	s.m.Lock()
	defer s.m.Unlock()
	s.enclaveStatusCode = es.StatusCode
	s.enclaveL1Head = es.L1Head
	if s.enclaveL2Head == nil || es.L2Head == nil {
		fmt.Println("±±±±± head increased from nil", s.enclaveL2Head, "to", es.L2Head)
	} else if big.NewInt(0).Add(s.enclaveL2Head, big.NewInt(1)).Cmp(es.L2Head) < 0 {
		fmt.Println("±±±±± head increased by more than one from", s.enclaveL2Head, "to", es.L2Head)
	}
	s.enclaveL2Head = es.L2Head

	s.setStatus(s.calculateStatus())
}

// OnDisconnected is called if the enclave is unreachable/not returning a valid Status
func (s *StateTracker) OnDisconnected() {
	s.m.Lock()
	defer s.m.Unlock()
	s.setStatus(Disconnected)
}

// when enclave is operational, this method will calculate the status based on comparison of current chain heads with enclave heads
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
		s.logger.Error("unknown enclave status code - this should not happen", "code", s.enclaveStatusCode)
		return Unavailable
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

func (s *StateTracker) GetEnclaveL2Head() *big.Int {
	s.m.RLock()
	defer s.m.RUnlock()
	if s.enclaveL2Head == nil {
		return nil
	}
	return big.NewInt(0).SetBytes(s.enclaveL2Head.Bytes())
}

// this must be called from within write-lock
func (s *StateTracker) setStatus(newStatus Status) {
	if s.status == newStatus {
		return
	}
	s.logger.Info(fmt.Sprintf("Updating enclave status from [%s] to [%s]", s.status, newStatus), "STATE", s)
	s.status = newStatus
}
