package subscription

import (
	"reflect"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// ForwardFromChannels - reads messages from the input channels, and calls the `onMessage` callback.
// Exits when the unsubscribed flag is true.
// Must be called as a go routine!
func ForwardFromChannels[R any](inputChannels []chan R, unsubscribed *atomic.Bool, onMessage func(R) error) {
	inputCases := make([]reflect.SelectCase, len(inputChannels)+1)

	// create a ticker to handle cleanup, check the "unsubscribed" flag and exit the goroutine
	inputCases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(time.NewTicker(2 * time.Second).C),
	}

	// create a select "case" for each input channel
	for i, ch := range inputChannels {
		inputCases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	unclosedInputChannels := len(inputCases)
	for unclosedInputChannels > 0 {
		chosen, value, ok := reflect.Select(inputCases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			inputCases[chosen].Chan = reflect.ValueOf(nil)
			unclosedInputChannels--
			continue
		}

		if unsubscribed != nil && unsubscribed.Load() {
			return
		}

		switch v := value.Interface().(type) {
		case time.Time:
			// exit the loop to avoid a goroutine leak
			if unsubscribed != nil && unsubscribed.Load() {
				return
			}
		case R:
			err := onMessage(v)
			if err != nil {
				// todo - log
				return
			}
		default:
			// ignore unexpected element
			continue
		}
	}
}

// HandleUnsubscribe - when the client calls "unsubscribe" or the subscription times out, it calls `onSub`
// Must be called as a go routine!
func HandleUnsubscribe(connectionSub *rpc.Subscription, unsubscribed *atomic.Bool, onUnsub func()) {
	<-connectionSub.Err()
	onUnsub()
}
