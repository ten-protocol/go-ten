package ethadapter

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ten-protocol/go-ten/contracts/generated/MessageBus"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkEnclaveRegistry"
	"github.com/ten-protocol/go-ten/contracts/generated/RollupContract"
	"strings"
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

//var NetworkConfigABI = NetworkConfig.NetworkConfigMetaData.ABI
//var MessageBusABI = MessageBus.MessageBusMetaData.ABI
//var CrossChainABI = CrossChain.CrossChainMetaData.ABI
//var NetworkEnclaveRegistryABI = NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData.ABI
//var RollupContractABI = RollupContract.RollupContractMetaData.ABI

var (
	MessageBusABI, _               = abi.JSON(strings.NewReader(MessageBus.MessageBusMetaData.ABI))
	NetworkConfigABI, _            = abi.JSON(strings.NewReader(NetworkConfig.NetworkConfigMetaData.ABI))
	RollupContractABI, _           = abi.JSON(strings.NewReader(RollupContract.RollupContractMetaData.ABI))
	EnclaveRegistryABI, _          = abi.JSON(strings.NewReader(NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData.ABI))
	CrossChainEventName            = "LogMessagePublished"
	CrossChainEventID              = MessageBusABI.Events[CrossChainEventName].ID
	ValueTransferEventName         = "ValueTransfer"
	ValueTransferEventID           = MessageBusABI.Events["ValueTransfer"].ID
	SequencerEnclaveGrantedEventID = EnclaveRegistryABI.Events["SequencerEnclaveGranted"].ID
	SequencerEnclaveRevokedEventID = EnclaveRegistryABI.Events["SequencerEnclaveRevoked"].ID
	NetworkSecretRequestedID       = EnclaveRegistryABI.Events["NetworkSecretRequested"].ID
	NetworkSecretRespondedID       = EnclaveRegistryABI.Events["NetworkSecretResponded"].ID
	RollupAddedID                  = RollupContractABI.Events["RollupAdded"].ID
)
