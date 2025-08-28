package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

const (
	maxGasForVisibility = 30_000 // hardcode at 30k gas.
)

type contractVisibilityReader struct {
	logger gethlog.Logger
}

func NewContractVisibilityReader(logger gethlog.Logger) ContractVisibilityReader {
	return &contractVisibilityReader{logger: logger}
}

func (v *contractVisibilityReader) ReadVisibilityConfig(ctx context.Context, evm *vm.EVM, contractAddress gethcommon.Address) (*core.ContractVisibilityConfig, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return v.readVisibilityConfig(evm, contractAddress)
	}
}

func (v *contractVisibilityReader) readVisibilityConfig(evm *vm.EVM, contractAddress gethcommon.Address) (*core.ContractVisibilityConfig, error) {
	cc, err := NewTransparencyConfigCaller(contractAddress, &localContractCaller{evm: evm, maxGasForVisibility: maxGasForVisibility})
	if err != nil {
		// unrecoverable error. should not happen
		v.logger.Crit(fmt.Sprintf("could not create transparency config caller. %v", err))
	}
	visibilityRules, err := cc.VisibilityRules(nil)
	if err != nil {
		// there is no visibility defined, so we return auto
		return &core.ContractVisibilityConfig{AutoConfig: true}, nil
	}

	transp := false
	if visibilityRules.ContractCfg == transparent {
		transp = true
	}

	cfg := &core.ContractVisibilityConfig{
		AutoConfig:   false,
		Transparent:  &transp,
		EventConfigs: make(map[gethcommon.Hash]*core.EventVisibilityConfig),
	}

	if transp {
		return cfg, nil
	}

	// only check the config for non-transparent contracts
	for i := range visibilityRules.EventLogConfigs {
		logConfig := visibilityRules.EventLogConfigs[i]
		eventConfig := eventCfg(logConfig)
		valErr := eventConfig.Validate()
		if valErr == nil {
			// ignore invalid configs
			cfg.EventConfigs[logConfig.EventSignature] = eventConfig
		}
	}

	return cfg, nil
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
}

// CodeAt - not implemented because it's not needed for our use case. It just has to return something non-nil
func (cc *localContractCaller) CodeAt(_ context.Context, _ gethcommon.Address, _ *big.Int) ([]byte, error) {
	return []byte{0}, nil
}

func (cc *localContractCaller) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	ret, _, err := cc.evm.Call(call.From, *call.To, call.Data, cc.maxGasForVisibility, uint256.NewInt(0))
	return ret, err
}
