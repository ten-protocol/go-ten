package ethadapter

import (
	"strings"

	"github.com/ten-protocol/go-ten/contracts/generated/DataAvailabilityRegistry"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkEnclaveRegistry"
)

const (
	RespondSecretMethod    = "respondNetworkSecret"
	RequestSecretMethod    = "requestNetworkSecret"
	InitializeSecretMethod = "initializeNetworkSecret" //#nosec

	AddRollupMethod = "addRollup"

	AddAdditionalAddressMethod = "addAdditionalAddress"
	MethodBytesLen             = 4
)

var (
	MessageBusABI, _               = abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	NetworkConfigABI, _            = abi.JSON(strings.NewReader(NetworkConfig.NetworkConfigMetaData.ABI))
	DataAvailabilityRegistryABI, _ = abi.JSON(strings.NewReader(DataAvailabilityRegistry.DataAvailabilityRegistryMetaData.ABI))
	EnclaveRegistryABI, _          = abi.JSON(strings.NewReader(NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData.ABI))

	CrossChainEventName                = "LogMessagePublished"
	ValueTransferEventName             = "ValueTransfer"
	NetworkSecretInitializedEventName  = "NetworkSecretInitialized"
	NetworkSecretRequestedEventName    = "NetworkSecretRequested"
	NetworkSecretRespondedEventName    = "NetworkSecretResponded"
	SequencerEnclaveGrantedEventName   = "SequencerEnclaveGranted"
	SequencerEnclaveRevokedEventName   = "SequencerEnclaveRevoked"
	RollupAddedEventName               = "RollupAdded"
	NetworkContractAddressAddedName    = "NetworkContractAddressAdded"
	AdditionalContractAddressAddedName = "AdditionalContractAddressAdded"

	CrossChainEventID                = MessageBusABI.Events[CrossChainEventName].ID
	ValueTransferEventID             = MessageBusABI.Events[ValueTransferEventName].ID
	NetworkSecretInitializedEventID  = EnclaveRegistryABI.Events[NetworkSecretInitializedEventName].ID
	SequencerEnclaveGrantedEventID   = EnclaveRegistryABI.Events[SequencerEnclaveGrantedEventName].ID
	SequencerEnclaveRevokedEventID   = EnclaveRegistryABI.Events[SequencerEnclaveRevokedEventName].ID
	NetworkSecretRequestedID         = EnclaveRegistryABI.Events[NetworkSecretRequestedEventName].ID
	NetworkSecretRespondedID         = EnclaveRegistryABI.Events[NetworkSecretRespondedEventName].ID
	RollupAddedID                    = DataAvailabilityRegistryABI.Events[RollupAddedEventName].ID
	NetworkContractAddressAddedID    = NetworkConfigABI.Events[NetworkContractAddressAddedName].ID
	AdditionalContractAddressAddedID = NetworkConfigABI.Events[AdditionalContractAddressAddedName].ID
)
