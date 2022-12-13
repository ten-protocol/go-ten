package p2p

import (
	"sync"
	"time"

	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	_failedMessageRead        = "Failed Message Reads"
	_failedMessageDecode      = "Failed Message Decodes"
	_failedConnectSendMessage = "Failed Peer Connects"
	_failedWriteSendMessage   = "Failed Socket Writes"
	_receivedMessage          = "Received Messages"
)

var _rollingPeriod = 5 * time.Minute

type status struct {
	lock          sync.RWMutex
	currentStatus map[string]*ttlMap
}

func newStatus() *status {
	return &status{
		lock: sync.RWMutex{},
		currentStatus: map[string]*ttlMap{
			_failedMessageDecode:      newTTLMap(_rollingPeriod),
			_failedMessageRead:        newTTLMap(_rollingPeriod),
			_failedConnectSendMessage: newTTLMap(_rollingPeriod),
			_failedWriteSendMessage:   newTTLMap(_rollingPeriod),
			_receivedMessage:          newTTLMap(_rollingPeriod),
		},
	}
}

func (s *status) increment(eventType string, host string) {
	// don't wait for locks - fire and forget
	go func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.currentStatus[eventType].increment(host)
	}()
}

// make sure it's returning a deep copy of the object to avoid multiple thread issues
func (s *status) status() *hostcommon.P2PStatus {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return &hostcommon.P2PStatus{
		FailedMessageReads:               s.currentStatus[_failedMessageRead].toMap(),
		FailedMessageDecodes:             s.currentStatus[_failedMessageDecode].toMap(),
		FailedSendMessagesPeerConnection: s.currentStatus[_failedConnectSendMessage].toMap(),
		FailedSendMessageWrites:          s.currentStatus[_failedWriteSendMessage].toMap(),
		ReceivedMessages:                 s.currentStatus[_receivedMessage].toMap(),
	}
}

func sumFailures(failures map[string]int64) int64 {
	total := int64(0)
	for _, v := range failures {
		total += v
	}
	return total
}

func peerNoMessage(receivedMsgs map[string]int64, knownPeers []string) []string {
	var disconnectedPeers []string
	for _, peerAddr := range knownPeers {
		// if the peer list was just updated then there's a bit chance the validator has not sent any message
		if _, ok := receivedMsgs[peerAddr]; !ok {
			disconnectedPeers = append(disconnectedPeers, peerAddr)
		}
	}

	return disconnectedPeers
}
