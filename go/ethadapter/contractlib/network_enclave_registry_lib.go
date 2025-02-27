package contractlib

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

type NetworkEnclaveRegistryLib interface {
	ContractLib
	CreateInitializeSecret(tx *common.L1InitializeSecretTx) (types.TxData, error)
	CreateRequestSecret(tx *common.L1RequestSecretTx) (types.TxData, error)
	CreateRespondSecret(tx *common.L1RespondSecretTx, verifyAttester bool) (types.TxData, error)
}

type networkEnclaveRegistryLibImpl struct {
	addr        *gethcommon.Address
	contractABI abi.ABI
	logger      gethlog.Logger
}

func NewNetworkEnclaveRegistryLib(addr *gethcommon.Address, logger gethlog.Logger) NetworkEnclaveRegistryLib {
	return &networkEnclaveRegistryLibImpl{
		addr:        addr,
		contractABI: ethadapter.EnclaveRegistryABI,
		logger:      logger,
	}
}

func (n *networkEnclaveRegistryLibImpl) IsMock() bool {
	return false
}

func (n *networkEnclaveRegistryLibImpl) GetContractAddr() *gethcommon.Address {
	return n.addr
}

func (n *networkEnclaveRegistryLibImpl) CreateInitializeSecret(tx *common.L1InitializeSecretTx) (types.TxData, error) {
	data, err := n.contractABI.Pack(
		ethadapter.InitializeSecretMethod,
		tx.EnclaveID,
		tx.InitialSecret,
		ethadapter.Base64EncodeToString(tx.Attestation),
	)

	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return &types.LegacyTx{
		To:   n.addr,
		Data: data,
	}, nil
}

func (n *networkEnclaveRegistryLibImpl) CreateRequestSecret(tx *common.L1RequestSecretTx) (types.TxData, error) {
	data, err := n.contractABI.Pack(ethadapter.RequestSecretMethod, ethadapter.Base64EncodeToString(tx.Attestation))
	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}

	return &types.LegacyTx{
		To:   n.addr,
		Data: data,
	}, nil
}

func (n *networkEnclaveRegistryLibImpl) CreateRespondSecret(tx *common.L1RespondSecretTx, verifyAttester bool) (types.TxData, error) {
	data, err := n.contractABI.Pack(
		ethadapter.RespondSecretMethod,
		tx.AttesterID,
		tx.RequesterID,
		tx.AttesterSig,
		tx.Secret,
		verifyAttester,
	)

	if err != nil {
		return nil, fmt.Errorf("could not pack the call data. Cause: %w", err)
	}
	return &types.LegacyTx{
		To:   n.addr,
		Data: data,
	}, nil
}

func (n *networkEnclaveRegistryLibImpl) DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error) {
	if tx.To() == nil || tx.To().Hex() != n.addr.Hex() || len(tx.Data()) == 0 {
		return nil, nil
	}
	method, err := n.contractABI.MethodById(tx.Data()[:ethadapter.MethodBytesLen])
	if err != nil {
		return nil, fmt.Errorf("could not decode tx. Cause: %w", err)
	}

	contractCallData := map[string]interface{}{}
	switch method.Name {
	case ethadapter.RespondSecretMethod:
		return n.unpackRespondSecretTx(tx, method, contractCallData)

	case ethadapter.RequestSecretMethod:
		return n.unpackRequestSecretTx(tx, method, contractCallData)

	case ethadapter.InitializeSecretMethod:
		return n.unpackInitSecretTx(tx, method, contractCallData)
	}
	return nil, nil
}

func (n *networkEnclaveRegistryLibImpl) unpackInitSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*common.L1InitializeSecretTx, error) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[ethadapter.MethodBytesLen:])
	if err != nil {
		return nil, fmt.Errorf("could not unpack transaction. Cause: %w", err)
	}
	callData, found := contractCallData["_genesisAttestation"]
	if !found {
		return nil, fmt.Errorf("call data not found for _genesisAttestation")
	}

	att, err := ethadapter.Base64DecodeFromString(callData.(string))
	if err != nil {
		return nil, err
	}

	return &common.L1InitializeSecretTx{
		Attestation: att,
	}, nil
}

func (n *networkEnclaveRegistryLibImpl) unpackRequestSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*common.L1RequestSecretTx, error) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[ethadapter.MethodBytesLen:])
	if err != nil {
		return nil, fmt.Errorf("could not unpack transaction. Cause: %w", err)
	}
	callData, found := contractCallData["requestReport"]
	if !found {
		return nil, fmt.Errorf("call data not found for requestReport")
	}

	att, err := ethadapter.Base64DecodeFromString(callData.(string))
	if err != nil {
		return nil, fmt.Errorf("could not decode attestation request. Cause: %w", err)
	}
	return &common.L1RequestSecretTx{
		Attestation: att,
	}, nil
}

func (n *networkEnclaveRegistryLibImpl) unpackRespondSecretTx(tx *types.Transaction, method *abi.Method, contractCallData map[string]interface{}) (*common.L1RespondSecretTx, error) {
	err := method.Inputs.UnpackIntoMap(contractCallData, tx.Data()[ethadapter.MethodBytesLen:])
	if err != nil {
		return nil, fmt.Errorf("could not unpack transaction. Cause: %w", err)
	}

	requesterData, found := contractCallData["requesterID"]
	if !found {
		return nil, fmt.Errorf("call data not found for requesterID")
	}
	requesterAddr, ok := requesterData.(gethcommon.Address)
	if !ok {
		return nil, fmt.Errorf("could not decode requester data")
	}

	attesterData, found := contractCallData["attesterID"]
	if !found {
		return nil, fmt.Errorf("call data not found for attesterID")
	}
	attesterAddr, ok := attesterData.(gethcommon.Address)
	if !ok {
		return nil, fmt.Errorf("could not decode attester data")
	}

	responseSecretData, found := contractCallData["responseSecret"]
	if !found {
		return nil, fmt.Errorf("call data not found for responseSecret")
	}
	responseSecretBytes, ok := responseSecretData.([]uint8)
	if !ok {
		return nil, fmt.Errorf("could not decode responseSecret data")
	}

	return &common.L1RespondSecretTx{
		AttesterID:  attesterAddr,
		RequesterID: requesterAddr,
		Secret:      responseSecretBytes[:],
	}, nil
}
