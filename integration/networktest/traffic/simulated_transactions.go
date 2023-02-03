package traffic

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/obscuronet/go-obscuro/go/obsclient"

	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration/networktest/userwallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

const (
	Low              = 2  // tx per sec
	High             = 20 // tx per sec
	_defaultNumUsers = 10
)

var _defaultTransferAmt = big.NewInt(10000000000)

func NativeFundsTransfers() Runner {
	return &randomTransactingSimRunner{
		numSimUsers: _defaultNumUsers,
		txPerSecond: Low,
		prepareSim:  nil, // no special prep needed, sim users will already have neem allocated native funds
		randomTransaction: func(ctx context.Context, user *randomTransactingSimUser) {
			startTime := time.Now()
			actionID := atomic.AddInt32(&user.actionID, 1)
			target := user.randomPeer()

			err := user.userWallet.SendFunds(ctx, target.Address(), _defaultTransferAmt)
			user.recordTxComplete(int(actionID), startTime, target.Address(), err)
		},
	}
}

func ERC20Transfers() Runner {
	return &randomTransactingSimRunner{
		numSimUsers: _defaultNumUsers,
		txPerSecond: Low,
		prepareSim: func(users []*randomTransactingSimUser, network networktest.NetworkConnector) error {
			// todo: deploy erc20 contract - how to get the contract back to the validators for checking?
			return nil
		},
		randomTransaction: func(ctx context.Context, user *randomTransactingSimUser) {
			// startTime := time.Now()
			// actionID := atomic.AddInt32(&user.actionID, 1)

			// target := user.randomPeer()

			// todo: execute erc20 transfers here
			// err := user.userWallet.SendFunds(ctx, target.Address(), _defaultTransferAmt)
			// user.recordTxComplete(int(actionID), startTime, target.Address(), err)
		},
	}
}

type randomTransactingSimRunner struct {
	// config
	txPerSecond float64
	numSimUsers int
	// method configured during setup that allows network and sim user preparation (e.g. deploy an erc20 and allocate some to each user)
	prepareSim func([]*randomTransactingSimUser, networktest.NetworkConnector) error
	// method configured during setup that dictates what the sim user will do when asked to make a random transaction
	// (this allows us to use this code for sending native funds, erc20 tx or anything else)
	randomTransaction func(ctx context.Context, user *randomTransactingSimUser)

	// state
	users             []*randomTransactingSimUser
	cancelUserContext context.CancelFunc
	eventsRecord      []ActionRecord
}

func (c *randomTransactingSimRunner) RunData() RunData {
	return &EventSliceRunData{
		events: c.eventsRecord,
		desc:   fmt.Sprintf("Simulated making native fund transfers - %f tx per sec across %d users", c.txPerSecond, c.numSimUsers),
	}
}

func (c *randomTransactingSimRunner) Start(network networktest.NetworkConnector) error {
	logger := testlog.Logger()
	userCtx, cancelUserCtx := context.WithCancel(context.Background())
	c.cancelUserContext = cancelUserCtx
	wallets := make([]*userwallet.UserWallet, c.numSimUsers)
	// create wallets
	for i := 0; i < c.numSimUsers; i++ {
		wal := datagenerator.RandomWallet(integration.ObscuroChainID)
		// traffic sim users are round robin-ed onto the validators for now
		wallets[i] = userwallet.NewUserWallet(wal.PrivateKey(), network.ValidatorRPCAddress(i%network.NumValidators()), logger)
		// register viewing key on sequencer todo: this is a workaround for bug, can remove this after #1396 is complete
		_, err := obsclient.DialWithAuth(network.SequencerRPCAddress(), wal, logger)
		if err != nil {
			return fmt.Errorf("temporary hack of setting up VK on seq failed - %w", err)
		}
	}
	// create users
	userTxRate := c.txPerSecond / float64(c.numSimUsers)
	for i := 0; i < c.numSimUsers; i++ {
		peerWallets := append(wallets[:i], wallets[i+1:]...)
		user := newNativeFundsSimRunner(i, wallets[i], c.randomTransaction, userTxRate, peerWallets, network)
		err := network.AllocateFaucetFunds(user.Account())
		if err != nil {
			cancelUserCtx()
			return fmt.Errorf("unable to allocate faucet funds - %w", err)
		}
		c.users = append(c.users, user)
	}
	if c.prepareSim != nil {
		err := c.prepareSim(c.users, network)
		if err != nil {
			cancelUserCtx()
			return fmt.Errorf("unable to complete custom `prepareSim` method - %w", err)
		}
	}

	// sim users send events for completed actions on this channel
	eventsCh := make(chan ActionRecord, 100)
	// start users (stoppable with ctx)
	for _, u := range c.users {
		go u.Start(userCtx, eventsCh)
	}

	// for now we just monitor the eventsCh and just store the events to process them all afterwards
	go func() {
		for {
			select {
			case ev := <-eventsCh:
				c.eventsRecord = append(c.eventsRecord, ev)
			case <-userCtx.Done():
				// todo: make sure we process any left in the queue before we exit?
				return
			}
		}
	}()
	return nil
}

func (c *randomTransactingSimRunner) Stop() {
	c.cancelUserContext()
}

func (c *randomTransactingSimRunner) Name() string {
	return "native-funds"
}

// randomTransactingSimUser wraps a sim user, adds functionality to run for a period of time, simulating interactions and recording results
type randomTransactingSimUser struct {
	idx            int
	avgDelayMillis int

	// method configured during setup that dictates what the sim user will do when asked to make a random transaction
	// (this allows us to use this code for sending native funds, erc20 tx or anything else)
	randomTransaction func(ctx context.Context, user *randomTransactingSimUser)

	peers      []*userwallet.UserWallet
	userWallet *userwallet.UserWallet
	network    networktest.NetworkConnector
	logger     gethlog.Logger

	eventsCh chan<- ActionRecord

	actionID int32
}

func (u *randomTransactingSimUser) Start(ctx context.Context, eventsCh chan<- ActionRecord) {
	u.eventsCh = eventsCh
	// loop until context closes
	for ctx.Err() == nil {
		// delay randomly between 0 and 2*avg delay
		delay := rand.Intn(2 * u.avgDelayMillis) //nolint:gosec
		time.Sleep(time.Duration(delay) * time.Millisecond)

		// separate go-routine to avoid skewing the tx rate waiting for the receipt
		// todo: is this safe? what if we start the next tx while the prev one is just getting started, will we hit nonce issues?
		go u.randomTransaction(ctx, u)
	}
}

func (u *randomTransactingSimUser) Account() gethcommon.Address {
	return u.userWallet.Address()
}

func (u *randomTransactingSimUser) randomPeer() wallet.Wallet {
	if len(u.peers) <= 1 {
		panic("no peers")
	}
	rndIdx := rand.Intn(len(u.peers)) //nolint:gosec
	for rndIdx == u.idx {
		// keep retrying if we get our own index
		rndIdx = rand.Intn(len(u.peers)) //nolint:gosec
	}
	return u.peers[rndIdx]
}

func (u *randomTransactingSimUser) recordTxComplete(id int, startTime time.Time, target gethcommon.Address, err error) {
	txEvent := &transferTransactionRecord{
		userID:    u.userWallet.Address(),
		actionID:  id,
		startTime: startTime,
		err:       err,
		amount:    _defaultTransferAmt,
		toAccount: target,
	}
	u.eventsCh <- txEvent
}

func newNativeFundsSimRunner(idx int,
	wallet *userwallet.UserWallet,
	randomTransaction func(ctx context.Context, user *randomTransactingSimUser),
	txPerSec float64,
	peers []*userwallet.UserWallet,
	network networktest.NetworkConnector,
) *randomTransactingSimUser {
	return &randomTransactingSimUser{
		idx: idx,
		// convert txPerSec to avg delay (e.g. 0.5 tx per sec -> 2000ms avg sleep)
		avgDelayMillis:    int(1000 * float64(1) / txPerSec),
		peers:             peers,
		userWallet:        wallet,
		randomTransaction: randomTransaction,
		network:           network,
		logger:            testlog.Logger(),
	}
}
