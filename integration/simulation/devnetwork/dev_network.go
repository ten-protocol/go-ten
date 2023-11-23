package devnetwork

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/simulation/network"

	"github.com/ten-protocol/go-ten/go/ethadapter"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/integration/networktest/userwallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/simulation/params"
)

var _defaultFaucetAmount = big.NewInt(750_000_000_000_000)

// InMemDevNetwork is a local dev network (L1 and L2) - the obscuro nodes are in-memory in a single go process, the L1 nodes are a docker geth network
//
// It can play the role of node operators and network admins to reproduce complex scenarios around nodes joining/leaving/failing.
//
// It also implements networktest.NetworkConnector to allow us to run the same NetworkTests against it that we can run against Testnets.
type InMemDevNetwork struct {
	logger gethlog.Logger

	// todo (@matt) - replace this with a struct for accs/contracts that are controlled by network admins
	// 	(don't pollute with "user" wallets, they will be controlled by the individual network test runners)
	networkWallets *params.SimWallets

	l1Network L1Network

	// When Obscuro network has been initialised on the L1 network, this will be populated
	// - if reconnecting to an existing network it needs to be populated when initialising this object
	// - if it is nil when `Start()` is called then Obscuro contracts will be deployed on the L1
	l1SetupData *params.L1SetupData

	obscuroConfig     ObscuroConfig
	obscuroSequencer  *InMemNodeOperator
	obscuroValidators []*InMemNodeOperator

	faucet     *userwallet.UserWallet
	faucetLock sync.Mutex
}

func (s *InMemDevNetwork) GetMCOwnerWallet() (wallet.Wallet, error) {
	return s.networkWallets.MCOwnerWallet, nil
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

	txHash, err := s.faucet.SendFunds(ctx, account, _defaultFaucetAmount, 1_000_000)
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
// todo (@matt) - this allows tests some basic L1 verification but in future this will need support more manipulation of L1 nodes,
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
	if s.logger == nil {
		s.logger = testlog.Logger()
	}
	fmt.Println("Starting L1 network")
	s.l1Network.Prepare()
	if s.l1SetupData == nil {
		// this is a new network, deploy the contracts to the L1
		fmt.Println("Deploying obscuro contracts to L1")
		s.deployObscuroNetworkContracts()
	}
	fmt.Println("Starting obscuro nodes")
	s.startNodes()
}

func (s *InMemDevNetwork) DeployL1StandardContracts() {
	// todo (@matt) - separate out L1 contract deployment from the geth network setup to give better sim control
}

func (s *InMemDevNetwork) startNodes() {
	if s.obscuroSequencer == nil {
		// initialise node operators
		s.obscuroSequencer = NewInMemNodeOperator(0, s.obscuroConfig, common.Sequencer, s.l1SetupData, s.l1Network.GetClient(0), s.networkWallets.NodeWallets[0], s.logger)
		for i := 1; i <= s.obscuroConfig.InitNumValidators; i++ {
			l1Client := s.l1Network.GetClient(i % s.l1Network.NumNodes())
			s.obscuroValidators = append(s.obscuroValidators, NewInMemNodeOperator(i, s.obscuroConfig, common.Validator, s.l1SetupData, l1Client, s.networkWallets.NodeWallets[i], s.logger))
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
			fmt.Println("failed to stop sequencer", err.Error())
		}
	}()
	go s.l1Network.CleanUp()

	s.logger.Info("Waiting for servers to stop.")
	time.Sleep(3 * time.Second)
}

func (s *InMemDevNetwork) deployObscuroNetworkContracts() {
	client := s.l1Network.GetClient(0)
	// note: we don't currently deploy ERC20s here, don't want to waste gas on sepolia
	l1SetupData, err := network.DeployObscuroNetworkContracts(client, s.networkWallets, false)
	if err != nil {
		panic(err)
	}
	s.l1SetupData = l1SetupData
}
