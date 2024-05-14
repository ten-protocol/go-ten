package subscription

import (
	"reflect"
	"sync/atomic"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// ForwardFromChannels - reads messages from all input channels, and calls the `onMessage` callback.
// Exits when the "stopped" flag is true or when the connection times out.
// Must be called as a go routine!
func ForwardFromChannels[R any](inputChannels []chan R, onMessage func(R) error, onBackendDisconnect func(), backendDisconnected *atomic.Bool, stopped *atomic.Bool, timeoutInterval time.Duration, logger gethlog.Logger) {
	inputCases := make([]reflect.SelectCase, len(inputChannels)+1)

	// create a ticker to handle cleanup, check the "stopped" flag and exit the goroutine
	inputCases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(time.NewTicker(2 * time.Second).C),
	}

	// create a select "case" for each input channel
	for i, ch := range inputChannels {
		inputCases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	lastMessageTime := time.Now()
loop:
	for {
		// this mechanism removes closed input channels. When there is none left, the subscription is considered "disconnected".
		_, value, ok := reflect.Select(inputCases)
		if !ok {
			logger.Error("Failed to read from the channel")
			break loop
		}

		// flag that the service needs to stop
		if stopped != nil && stopped.Load() {
			return
		}

		// flag that the backend channels were disconnected
		if backendDisconnected != nil && backendDisconnected.Load() {
			break loop
		}

		switch v := value.Interface().(type) {
		case time.Time:
			// no message was received longer than the timeout. Exiting.
			if time.Since(lastMessageTime) > timeoutInterval {
				break loop
			}
		case R:
			lastMessageTime = time.Now()
			err := onMessage(v)
			if err != nil {
				logger.Error("Failed to process message", log.ErrKey, err)
				break loop
			}
		default:
			// ignore unexpected element
			logger.Warn("Received unexpected message type.", "type", reflect.TypeOf(v), "value", value)
			break loop
		}
	}

	if onBackendDisconnect != nil {
		onBackendDisconnect()
	}
}

// HandleUnsubscribe - when the client calls "unsubscribe" or the subscription times out, it calls `onSub`
// Must be called as a go routine!
func HandleUnsubscribe(connectionSub *rpc.Subscription, onUnsub func()) {
	<-connectionSub.Err()
	onUnsub()
}

// HandleUnsubscribeErrChan - when the client calls "unsubscribe" or the subscription times out, it calls `onSub`
// Must be called as a go routine!
func HandleUnsubscribeErrChan(errChan []<-chan error, onUnsub func()) {
	inputCases := make([]reflect.SelectCase, len(errChan))
	for i, ch := range errChan {
		inputCases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	reflect.Select(inputCases)
	onUnsub()
}
