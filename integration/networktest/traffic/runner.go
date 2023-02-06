package traffic

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

// Runner is used to simulate traffic against a durationTrafficTest network
// When Start() is called then simulated traffic will begin against the provided NetworkConnector, until Stop() is called.
type Runner interface {
	// Start is only expected to be called once for the lifetime of the Runner
	Start(network networktest.NetworkConnector) error
	Stop()
	RunData() RunData
	Name() string // brief name (used in log file naming etc., avoid spaces)
}

// RunData builds up data about the simulated actions during the traffic run
// - this can be validated as part of an automated durationTrafficTest or perhaps statistics or visualisations created etc.
// Note: this data can be combined with the runner's data and the network to verify post-state
type RunData interface {
	ActionEvents() []ActionRecord // get a list of completed actions
	RunDescription() string       // few words describing the simulation that was run
}

type Action interface {
	Name() string
}

const (
	SendTransaction = "SendTransaction"
)

type ActionRecord interface {
	SimUserID() common.Address
	ActionID() int // action ID for that user
	StartTime() time.Time
	CompleteTime() time.Time
	Success() bool
	Description() string
}

// record used for simple transaction actions (sending quantity of native/erc20 from user to user)
type transferTransactionRecord struct {
	userID    common.Address
	actionID  int
	desc      string
	startTime time.Time
	endTime   time.Time
	err       error
	amount    *big.Int
	toAccount common.Address
}

func (t *transferTransactionRecord) SimUserID() common.Address {
	return t.userID
}

func (t *transferTransactionRecord) ActionID() int {
	return t.actionID
}

func (t *transferTransactionRecord) StartTime() time.Time {
	return t.startTime
}

func (t *transferTransactionRecord) CompleteTime() time.Time {
	return t.endTime
}

func (t *transferTransactionRecord) Success() bool {
	return t.err == nil
}

func (t *transferTransactionRecord) Description() string {
	return t.desc
}

//// todo: ideally all this thing needs to be given is a "create sim user" setup function
//type basicRunner struct {
//	// factory function to produce a simulated "user"
//	userFactory   func(userIdx int, actionCh <-chan Action, network networktest.NetworkConnector) (SimUser, error)
//	numSims       int
//	actionsPerSec float64
//
//	simUsers         []SimUser
//	userChans        []chan<- Action
//	data             *EventSliceRunData
//	cancelRunningSim context.CancelFunc
//}
//
//func (b *basicRunner) Start(network networktest.NetworkConnector) error {
//	defer b.cleanUp() // make sure we don't leave any goroutines dangling
//	b.simUsers = make([]SimUser, 0)
//
//	// prepare the network with any contracts etc. the runner needs
//
//	// create and prepare sim users
//	for i := 0; i < b.numSims; i++ {
//		userCh := make(chan Action)
//		user, err := b.userFactory(i, userCh, network)
//		if err != nil {
//			return fmt.Errorf("unable to create sim user idx=%d - %w", i, err)
//		}
//		b.userChans = append(b.userChans, userCh)
//		b.simUsers = append(b.simUsers, user)
//		// user will consume actions from its userCh in a separate goroutine until the channel is closed
//		go user.Start()
//	}
//
//	// run the sim. Round-robin the users at the prescribed rate, instructing them to perform their action
//	simCtx, cancelSimRun := context.WithCancel(context.Background())
//	b.cancelRunningSim = cancelSimRun
//	go b.Run(simCtx)
//
//	return nil
//}
//
//func (b *basicRunner) Stop() {
//	if b.cancelRunningSim != nil {
//		b.cancelRunningSim()
//	}
//}
//
//func (b *basicRunner) RunData() RunData {
//	return b.data
//}
//
//// Run round-robins actions to the sim users at random intervals (averaging at b.actionsPerSec rate)
//func (b *basicRunner) Run(ctx context.Context) {
//	nextUser := 0
//	for {
//		select {}
//	}
//}

type EventSliceRunData struct {
	events []ActionRecord
	desc   string
}

func (e *EventSliceRunData) ActionEvents() []ActionRecord {
	return e.events
}

func (e *EventSliceRunData) RunDescription() string {
	return e.desc
}

func (e *EventSliceRunData) AddEvent(event ActionRecord) {
	e.events = append(e.events, event)
}

type SimUser interface {
	Start()
}
