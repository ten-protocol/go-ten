package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

const (
	maxGasForVisibility = 200_000 // hardcode at 30k gas.
)

type contractVisibilityReader struct {
	logger gethlog.Logger
}

func NewContractVisibilityReader(logger gethlog.Logger) ContractVisibilityReader {
	return &contractVisibilityReader{logger: logger}
}

func (v *contractVisibilityReader) readVisibilityConfig(evm *vm.EVM, contractAddress gethcommon.Address, gasCap uint64) (*core.ContractVisibilityConfig, uint64, error) {
	cap := gasCap
	if cap == 0 || cap > maxGasForVisibility {
		cap = maxGasForVisibility
	}
	lcc := &localContractCaller{evm: evm, maxGasForVisibility: cap}

	cc, err := NewTransparencyConfigCaller(contractAddress, lcc)
	if err != nil {
		// unrecoverable; should not happen
		v.logger.Crit(fmt.Sprintf("could not create transparency config caller. %v", err))
	}
	visibilityRules, err := cc.VisibilityRules(nil)
	if err != nil {
		// no visibility defined => auto
		return &core.ContractVisibilityConfig{AutoConfig: true}, lcc.usedGasLast, nil
	}

	transp := visibilityRules.ContractCfg == transparent
	cfg := &core.ContractVisibilityConfig{
		AutoConfig:   false,
		Transparent:  &transp,
		EventConfigs: make(map[gethcommon.Hash]*core.EventVisibilityConfig),
	}
	if transp {
		return cfg, lcc.usedGasLast, nil
	}
	for i := range visibilityRules.EventLogConfigs {
		logConfig := visibilityRules.EventLogConfigs[i]
		eventConfig := eventCfg(logConfig)
		if valErr := eventConfig.Validate(); valErr == nil {
			cfg.EventConfigs[logConfig.EventSignature] = eventConfig
		}
	}
	return cfg, lcc.usedGasLast, nil
}

func eventCfg(logConfig ContractTransparencyConfigEventLogConfig) *core.EventVisibilityConfig {
	relevantToMap := make(map[uint8]bool)
	for _, field := range logConfig.VisibleTo {
		relevantToMap[field] = true
	}
	isPublic := relevantToMap[everyone]

	if isPublic {
		return &core.EventVisibilityConfig{AutoConfig: false, Public: true}
	}

	t1 := relevantToMap[topic1]
	t2 := relevantToMap[topic2]
	t3 := relevantToMap[topic3]
	s := relevantToMap[sender]
	return &core.EventVisibilityConfig{
		AutoConfig:    false,
		Public:        false,
		Topic1CanView: &t1,
		Topic2CanView: &t2,
		Topic3CanView: &t3,
		SenderCanView: &s,
	}
}

// used as a wrapper around the vm.EVM to allow for easier calling of smart contract view functions
type localContractCaller struct {
	evm                 *vm.EVM
	maxGasForVisibility uint64
	usedGasLast         uint64
}

// CodeAt - not implemented because it's not needed for our use case. It just has to return something non-nil
func (cc *localContractCaller) CodeAt(_ context.Context, _ gethcommon.Address, _ *big.Int) ([]byte, error) {
	return []byte{0}, nil
}

func (cc *localContractCaller) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	// Prefer a static call to enforce non-mutating behavior of view/pure methods.
	ret, left, err := cc.evm.StaticCall(call.From, *call.To, call.Data, cc.maxGasForVisibility)
	cc.usedGasLast = cc.maxGasForVisibility - left
	return ret, err
}

// ReadVisibilityConfig performs the same ABI call but returns (config, gasUsed).
// gasCap will be clamped to the package-level maxGasForVisibility.
func (v *contractVisibilityReader) ReadVisibilityConfig(ctx context.Context, evm *vm.EVM, contractAddress gethcommon.Address, gasCap uint64) (*core.ContractVisibilityConfig, uint64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
		cfg, gasUsed, err := v.readVisibilityConfig(evm, contractAddress, gasCap)
		return cfg, gasUsed, err
	}
}
