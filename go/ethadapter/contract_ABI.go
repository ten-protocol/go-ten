package ethadapter

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkEnclaveRegistry"
	"github.com/ten-protocol/go-ten/contracts/generated/RollupContract"
)

const (
	GetCrossChainContractAddress             = "crossChainContractAddress"
	GetMessageBusContractAddress             = "messageBusContractAddress"
	GetNetworkEnclaveRegistryContractAddress = "networkEnclaveRegistryContractAddress"
	GetRollupContractAddress                 = "rollupContractAddress"
	GetContractAddresses                     = "addresses"

	ExtractNativeValueMethod = "extractNativeValue"
	PauseWithdrawals         = "pauseWithdrawals"

	RespondSecretMethod    = "respondNetworkSecret"
	RequestSecretMethod    = "requestNetworkSecret"
	InitializeSecretMethod = "initializeNetworkSecret" //#nosec

	AddRollupMethod = "addRollup"

	MethodBytesLen = 4
)

var (
	MessageBusABI, _      = abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	NetworkConfigABI, _   = abi.JSON(strings.NewReader(NetworkConfig.NetworkConfigMetaData.ABI))
	RollupContractABI, _  = abi.JSON(strings.NewReader(RollupContract.RollupContractMetaData.ABI))
	EnclaveRegistryABI, _ = abi.JSON(strings.NewReader(NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData.ABI))

	CrossChainEventName               = "LogMessagePublished"
	ValueTransferEventName            = "ValueTransfer"
	NetworkSecretInitializedEventName = "NetworkSecretInitialized"
	NetworkSecretRequestedEventName   = "NetworkSecretRequested"
	NetworkSecretRespondedEventName   = "NetworkSecretResponded"
	SequencerEnclaveGrantedEventName  = "SequencerEnclaveGranted"
	SequencerEnclaveRevokedEventName  = "SequencerEnclaveRevoked"
	RollupAddedEventName              = "RollupAdded"

	CrossChainEventID               = MessageBusABI.Events[CrossChainEventName].ID
	ValueTransferEventID            = MessageBusABI.Events[ValueTransferEventName].ID
	NetworkSecretInitializedEventID = EnclaveRegistryABI.Events[NetworkSecretInitializedEventName].ID
	SequencerEnclaveGrantedEventID  = EnclaveRegistryABI.Events[SequencerEnclaveGrantedEventName].ID
	SequencerEnclaveRevokedEventID  = EnclaveRegistryABI.Events[SequencerEnclaveRevokedEventName].ID
	NetworkSecretRequestedID        = EnclaveRegistryABI.Events[NetworkSecretRequestedEventName].ID
	NetworkSecretRespondedID        = EnclaveRegistryABI.Events[NetworkSecretRespondedEventName].ID
	RollupAddedID                   = RollupContractABI.Events[RollupAddedEventName].ID
)
