package devnetwork

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/integration/networktest/userwallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/wallet"
	"github.com/obscuronet/go-obscuro/integration"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/simulation/params"
)

var _defaultFaucetAmount = big.NewInt(750_000_000_000_000)

// InMemDevNetwork is a local dev network (L1 and L2) - the obscuro nodes are in-memory in a single go process, the L1 nodes are a docker geth network
//
// It can play the role of node operators and network admins to reproduce complex scenarios around nodes joining/leaving/failing.
//
// It also implements networktest.NetworkConnector to allow us to run the same NetworkTests against it that we can run against Testnets.
type InMemDevNetwork struct {
	logger gethlog.Logger

	// todo: replace this with a struct for accs/contracts that are controlled by network admins
	// 	(don't pollute with "user" wallets, they will be controlled by the individual network test runners)
	networkWallets *params.SimWallets

	l1Network L1Network

	obscuroConfig     ObscuroConfig
	obscuroSequencer  *InMemNodeOperator
	obscuroValidators []*InMemNodeOperator

	faucet     *userwallet.UserWallet
	faucetLock sync.Mutex
}

func (s *InMemDevNetwork) ChainID() int64 {
	return integration.ObscuroChainID
}

func (s *InMemDevNetwork) FaucetWallet() wallet.Wallet {
	return s.networkWallets.L2FaucetWallet
}

func (s *InMemDevNetwork) AllocateFaucetFunds(ctx context.Context, account gethcommon.Address) error {
	// ensure only one test account is getting faucet funds at a time, faucet client isn't thread-safe
	s.faucetLock.Lock()
	defer s.faucetLock.Unlock()

	txHash, err := s.faucet.SendFunds(ctx, account, _defaultFaucetAmount)
	if err != nil {
		return err
	}

	receipt, err := s.faucet.AwaitReceipt(ctx, txHash)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("faucet transaction receipt status not successful - %v", receipt.Status)
	}
	return nil
}

func (s *InMemDevNetwork) SequencerRPCAddress() string {
	seq := s.GetSequencerNode()
	return seq.HostRPCAddress()
}

func (s *InMemDevNetwork) ValidatorRPCAddress(idx int) string {
	val := s.GetValidatorNode(idx)
	return val.HostRPCAddress()
}

// GetL1Client returns the first client we have for our local L1 network
// todo: this allows tests some basic L1 verification but in future this will need support more manipulation of L1 nodes,
//
//	(to allow us to simulate various scenarios where L1 is unavailable, under attack, etc.)
func (s *InMemDevNetwork) GetL1Client() (ethadapter.EthClient, error) {
	return s.l1Network.GetClient(0), nil
}

func (s *InMemDevNetwork) GetSequencerNode() networktest.NodeOperator {
	return s.obscuroSequencer
}

func (s *InMemDevNetwork) GetValidatorNode(i int) networktest.NodeOperator {
	return s.obscuroValidators[i]
}

func (s *InMemDevNetwork) NumValidators() int {
	return len(s.obscuroValidators)
}

func (s *InMemDevNetwork) Start() {
	s.l1Network.Start()
	s.startNodes()
}

func (s *InMemDevNetwork) DeployL1StandardContracts() {
	// todo: separate out L1 contract deployment from the geth network setup to give better sim control
}

func (s *InMemDevNetwork) startNodes() {
	if s.obscuroSequencer == nil {
		// initialise node operators
		s.obscuroSequencer = NewInMemNodeOperator(0, s.obscuroConfig, common.Sequencer, s.l1Network.ObscuroSetupData(), s.l1Network.GetClient(0), s.networkWallets.NodeWallets[0], s.logger)
		for i := 1; i <= s.obscuroConfig.InitNumValidators; i++ {
			l1Client := s.l1Network.GetClient(i % s.l1Network.NumNodes())
			s.obscuroValidators = append(s.obscuroValidators, NewInMemNodeOperator(i, s.obscuroConfig, common.Validator, s.l1Network.ObscuroSetupData(), l1Client, s.networkWallets.NodeWallets[i], s.logger))
		}
	}

	go func() {
		err := s.obscuroSequencer.Start()
		if err != nil {
			panic(err)
		}
	}()
	for _, v := range s.obscuroValidators {
		go func(v networktest.NodeOperator) {
			err := v.Start()
			if err != nil {
				panic(err)
			}
		}(v)
	}
	s.faucet = userwallet.NewUserWallet(s.networkWallets.L2FaucetWallet.PrivateKey(), s.SequencerRPCAddress(), s.logger)
}

func (s *InMemDevNetwork) CleanUp() {
	for _, v := range s.obscuroValidators {
		go func(v networktest.NodeOperator) {
			err := v.Stop()
			if err != nil {
				fmt.Println("failed to stop validator", err.Error())
			}
		}(v)
	}
	go func() {
		err := s.obscuroSequencer.Stop()
		if err != nil {
			panic(err)
		}
	}()
	go s.l1Network.Stop()

	s.logger.Info("Waiting for servers to stop.")
	time.Sleep(3 * time.Second)
}
