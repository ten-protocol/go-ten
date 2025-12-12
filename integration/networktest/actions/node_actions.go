package actions

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/wallet"
	"github.com/ten-protocol/go-ten/integration"
	"github.com/ten-protocol/go-ten/integration/common/testlog"
	"github.com/ten-protocol/go-ten/integration/networktest"
	"github.com/ten-protocol/go-ten/integration/simulation/devnetwork"
)

func StartValidatorEnclave(validatorIdx int) networktest.Action {
	return &startValidatorEnclaveAction{validatorIdx: validatorIdx}
}

type startValidatorEnclaveAction struct {
	validatorIdx int
}

func (s *startValidatorEnclaveAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	fmt.Printf("Validator %d: starting enclave\n", s.validatorIdx)
	validator := network.GetValidatorNode(s.validatorIdx)
	// note: these actions are assuming single-enclave setups
	err := validator.StartEnclave(0)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (s *startValidatorEnclaveAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

// ConfigApplier is a function that can read from context and apply config changes
// This allows config to be determined dynamically at runtime based on previous test steps
type ConfigApplier func(ctx context.Context, cfg *devnetwork.TenConfig)

// ApplySharedSecretFromContext creates a ConfigApplier that reads SharedSecret from context
// and applies it to the config. This is useful for tests that need to retrieve the shared
// secret dynamically during test execution (e.g., from a backup) and pass it to a new node.
func ApplySharedSecretFromContext() ConfigApplier {
	return func(ctx context.Context, cfg *devnetwork.TenConfig) {
		if secret := ctx.Value("SharedSecret"); secret != nil {
			cfg.SharedSecret = secret.(string)
		}
	}
}

func StartNewValidatorNode(configAppliers ...ConfigApplier) networktest.Action {
	return &startNewValidatorNodeAction{configAppliers: configAppliers}
}

type startNewValidatorNodeAction struct {
	configAppliers []ConfigApplier
}

func (s *startNewValidatorNodeAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	// Generate a new wallet for the validator node
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return ctx, fmt.Errorf("failed to generate private key for new validator: %w", err)
	}

	nodeWallet := wallet.NewInMemoryWalletFromPK(big.NewInt(integration.EthereumChainID), privateKey, testlog.Logger())

	/*	// Fund the new wallet from the contract owner wallet
		contractOwner, err := network.GetContractOwnerWallet()
		if err != nil {
			return ctx, fmt.Errorf("failed to get contract owner wallet: %w", err)
		}

		// Transfer some ETH to the new validator wallet for gas fees
		fundingAmount := big.NewInt(1000000000000000000) // 1 ETH
		tx, err := contractOwner.SendFunds(ctx, nodeWallet.Address(), fundingAmount)
		if err != nil {
			return ctx, fmt.Errorf("failed to fund new validator wallet: %w", err)
		}

		fmt.Printf("Funding new validator wallet %s with transaction %s\n", nodeWallet.Address().Hex(), tx.Hex())

		// Wait for the funding transaction to be mined
		time.Sleep(2 * time.Second)
	*/
	// Get the current TenConfig from the network
	devNetwork, ok := network.(*devnetwork.InMemDevNetwork)
	if !ok {
		return ctx, fmt.Errorf("network does not support creating new validator nodes")
	}

	// Create a copy of the config to avoid modifying the shared config
	currentConfig := devNetwork.TenConfig()
	newConfig := *currentConfig

	// Apply each config applier function - they can read from context and modify config
	for _, applier := range s.configAppliers {
		applier(ctx, &newConfig)
	}

	// Create the new validator node
	newValidator := network.NewValidatorNode(&newConfig, nodeWallet)

	fmt.Printf("Starting new validator node (index: %d)\n", network.NumValidators()-1)

	// Start the new validator node
	err = newValidator.Start()
	if err != nil {
		return ctx, fmt.Errorf("failed to start new validator node: %w", err)
	}

	fmt.Printf("New validator node started successfully at %s\n", newValidator.HostRPCWSAddress())

	return ctx, nil
}

func (s *startNewValidatorNodeAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

func StopValidatorEnclave(validatorIdx int) networktest.Action {
	return &stopValidatorEnclaveAction{validatorIdx: validatorIdx}
}

type stopValidatorEnclaveAction struct {
	validatorIdx int
}

func (s *stopValidatorEnclaveAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	fmt.Printf("Validator %d: stopping enclave\n", s.validatorIdx)
	validator := network.GetValidatorNode(s.validatorIdx)
	err := validator.StopEnclave(0)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (s *stopValidatorEnclaveAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

func StopValidatorHost(validatorIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Validator %d: stopping host\n", validatorIdx)
		validator := network.GetValidatorNode(validatorIdx)
		err := validator.StopHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StartValidatorHost(validatorIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Validator %d: starting host\n", validatorIdx)
		validator := network.GetValidatorNode(validatorIdx)
		err := validator.StartHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StopSequencerHost() networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Println("Sequencer: stopping host")
		sequencer := network.GetSequencerNode()
		err := sequencer.StopHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StartSequencerHost() networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Println("Sequencer: starting host")
		sequencer := network.GetSequencerNode()
		err := sequencer.StartHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StopSequencerEnclave(enclaveIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Sequencer: stopping enclave %d\n", enclaveIdx)
		sequencer := network.GetSequencerNode()
		err := sequencer.StopEnclave(enclaveIdx)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StartSequencerEnclave(enclaveIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Sequencer: starting enclave %d\n", enclaveIdx)
		sequencer := network.GetSequencerNode()
		err := sequencer.StartEnclave(enclaveIdx)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func WaitForValidatorHealthCheck(validatorIdx int, maxWait time.Duration) networktest.Action {
	return &waitForValidatorHealthCheckAction{
		validatorIdx: validatorIdx,
		maxWait:      maxWait,
	}
}

type waitForValidatorHealthCheckAction struct {
	validatorIdx int
	maxWait      time.Duration
}

func (w *waitForValidatorHealthCheckAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	validator := network.GetValidatorNode(w.validatorIdx)
	// poll the health check until success or timeout
	err := retry.Do(func() error {
		return networktest.NodeHealthCheck(validator.HostRPCWSAddress())
	}, retry.NewTimeoutStrategy(w.maxWait, 1*time.Second))
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (w *waitForValidatorHealthCheckAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

func WaitForSequencerHealthCheck(maxWait time.Duration) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		sequencer := network.GetSequencerNode()
		// poll the health check until success or timeout
		err := retry.Do(func() error {
			return networktest.NodeHealthCheck(sequencer.HostRPCWSAddress())
		}, retry.NewTimeoutStrategy(maxWait, 1*time.Second))
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}
