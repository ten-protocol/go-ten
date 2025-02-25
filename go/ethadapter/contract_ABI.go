package ethadapter

import (
	"github.com/ten-protocol/go-ten/contracts/generated/CrossChain"
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

var NetworkConfigABI = NetworkConfig.NetworkConfigMetaData.ABI
var MessageBusABI = MessageBus.MessageBusMetaData.ABI
var CrossChainABI = CrossChain.CrossChainMetaData.ABI
var NetworkEnclaveRegistryABI = NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData.ABI
var RollupContractABI = RollupContract.RollupContractMetaData.ABI
