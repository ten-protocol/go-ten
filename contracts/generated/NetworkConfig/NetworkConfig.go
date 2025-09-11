// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package NetworkConfig

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// NetworkConfigAddresses is an auto generated low-level Go binding around an user-defined struct.
type NetworkConfigAddresses struct {
	CrossChain               common.Address
	MessageBus               common.Address
	NetworkEnclaveRegistry   common.Address
	DataAvailabilityRegistry common.Address
	L1Bridge                 common.Address
	L2Bridge                 common.Address
	L1CrossChainMessenger    common.Address
	L2CrossChainMessenger    common.Address
	AdditionalContracts      []NetworkConfigNamedAddress
}

// NetworkConfigFixedAddresses is an auto generated low-level Go binding around an user-defined struct.
type NetworkConfigFixedAddresses struct {
	CrossChain               common.Address
	MessageBus               common.Address
	NetworkEnclaveRegistry   common.Address
	DataAvailabilityRegistry common.Address
}

// NetworkConfigNamedAddress is an auto generated low-level Go binding around an user-defined struct.
type NetworkConfigNamedAddress struct {
	Name string
	Addr common.Address
}

// NetworkConfigMetaData contains all meta data concerning the NetworkConfig contract.
var NetworkConfigMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"AdditionalContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"AdditionalContractAddressRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NetworkContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"featureName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"featureData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validAtBlock\",\"type\":\"uint256\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DATA_AVAILABILITY_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FORK_MANAGER_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MESSAGE_BUS_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NETWORK_ENCLAVE_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addAdditionalAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractName\",\"type\":\"string\"}],\"name\":\"additionalAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"addressNames\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addresses\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dataAvailabilityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1CrossChainMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2CrossChainMessenger\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.NamedAddress[]\",\"name\":\"additionalContracts\",\"type\":\"tuple[]\"}],\"internalType\":\"structNetworkConfig.Addresses\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdditionalContractNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dataAvailabilityRegistry\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.FixedAddresses\",\"name\":\"_addresses\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkEnclaveRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"removeAdditionalAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"featureName\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"featureData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validAtBlock\",\"type\":\"uint256\"}],\"name\":\"upgradeFeature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b5061001861001d565b6100fc565b5f6100266100bd565b805490915068010000000000000000900460ff16156100585760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100ba5780546001600160401b0319166001600160401b0390811782556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916100b1916100e7565b60405180910390a15b50565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6001600160401b0382168152602081016100e1565b612419806101095f395ff3fe608060405234801561000f575f5ffd5b506004361061021b575f3560e01c806385f427cb11610123578063be9f8207116100b8578063e30c397811610088578063f5e9f2861161006e578063f5e9f286146103fd578063faa5e2de14610405578063fbfd6d911461040d575f5ffd5b8063e30c3978146103e2578063f2fde38b146103ea575f5ffd5b8063be9f82071461039f578063da0321cd146103a7578063dbf5d903146103bc578063e1825d06146103cf575f5ffd5b8063a1b918d6116100f3578063a1b918d614610348578063ae61ecba14610350578063af45463514610358578063b7bef9ab1461038c575f5ffd5b806385f427cb1461031d5780638da5cb5b14610325578063934746a71461032d57806396493cc514610335575f5ffd5b8063556d89dd116101b3578063715018a61161018357806372bad9121161016957806372bad9121461030557806379ba50971461030d578063812b1ffe14610315575f5ffd5b8063715018a6146102f557806371fd11f3146102fd575f5ffd5b8063556d89dd146102bf5780635ab2a558146102d257806367cc852e146102da5780636c1358ac146102e2575f5ffd5b806331d1464d116101ee57806331d1464d1461027a578063450948ad1461028f57806346a30a781461029757806348d87239146102aa575f5ffd5b80630b592f451461021f5780630f387b1e1461023d57806313eeee961461025d5780632fc00c7614610272575b5f5ffd5b610227610415565b6040516102349190611650565b60405180910390f35b61025061024b366004611675565b61044d565b60405161023491906116d5565b6102656104f2565b6040516102349190611759565b6102276105c5565b61028d6102883660046117b8565b6105f4565b005b6102276106c1565b61028d6102a5366004611811565b6106f0565b6102b261078b565b6040516102349190611834565b61028d6102cd366004611811565b6107b9565b610227610844565b6102b2610873565b61028d6102f0366004611905565b61089e565b61028d610b6e565b610227610b8e565b6102b2610bbd565b61028d610be8565b6102b2610c27565b6102b2610c52565b610227610c7d565b610227610cb1565b61028d610343366004611811565b610ce0565b610227610d6b565b6102b2610d9a565b6102276103663660046119c7565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b61028d61039a366004611a07565b610dc5565b6102b2610f58565b6103af610f83565b6040516102349190611bad565b61028d6103ca366004611bbe565b61121b565b61028d6103dd366004611811565b611292565b61022761131d565b61028d6103f8366004611811565b611345565b6102276113ca565b6102b26113f9565b6102b2611424565b5f61044861044460017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611c59565b5490565b905090565b5f818154811061045b575f80fd5b905f5260205f20015f91509050805461047390611c80565b80601f016020809104026020016040519081016040528092919081815260200182805461049f90611c80565b80156104ea5780601f106104c1576101008083540402835291602001916104ea565b820191905f5260205f20905b8154815290600101906020018083116104cd57829003601f168201915b505050505081565b60605f805480602002602001604051908101604052809291908181526020015f905b828210156105bc578382905f5260205f2001805461053190611c80565b80601f016020809104026020016040519081016040528092919081815260200182805461055d90611c80565b80156105a85780601f1061057f576101008083540402835291602001916105a8565b820191905f5260205f20905b81548152906001019060200180831161058b57829003601f168201915b505050505081526020019060010190610514565b50505050905090565b5f61044861044460017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611c59565b6105fc61144f565b5f6001600160a01b031660018383604051610618929190611cbc565b908152604051908190036020019020546001600160a01b0316036106575760405162461bcd60e51b815260040161064e90611cf9565b60405180910390fd5b60018282604051610669929190611cbc565b90815260405190819003602001812080546001600160a01b03191690557f5f9e4bc50ed4fc3c0bf14c4b518e1f8132c7a95ce5bed2e97c8675e6adf03573906106b59084908490611d29565b60405180910390a15050565b5f61044861044460017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611c59565b6106f861144f565b6001600160a01b03811661071e5760405162461bcd60e51b815260040161064e90611d6d565b61075161074c60017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611c59565b829055565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516107809190611daf565b60405180910390a150565b6107b660017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611c59565b81565b6107c161144f565b6001600160a01b0381166107e75760405162461bcd60e51b815260040161064e90611d6d565b61081561074c60017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611c59565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516107809190611e00565b5f61044861044460017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611c59565b6107b660017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611c59565b5f6108a7611483565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f811580156108d35750825b90505f8267ffffffffffffffff1660011480156108ef5750303b155b9050811580156108fd575080155b15610934576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561096857845468ff00000000000000001916680100000000000000001785555b6001600160a01b03861661098e5760405162461bcd60e51b815260040161064e90611e42565b86516001600160a01b03166109b55760405162461bcd60e51b815260040161064e90611e84565b60208701516001600160a01b03166109df5760405162461bcd60e51b815260040161064e90611ec6565b60408701516001600160a01b0316610a095760405162461bcd60e51b815260040161064e90611f30565b60608701516001600160a01b0316610a335760405162461bcd60e51b815260040161064e90611f98565b610a3c866114ad565b610a70610a6a60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611c59565b88519055565b610aa7610a9e60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611c59565b60208901519055565b610ade610ad560017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611c59565b60408901519055565b610b15610b0c60017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611c59565b60608901519055565b8315610b6557845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610b5c90600190611fcb565b60405180910390a15b50505050505050565b610b7661144f565b60405162461bcd60e51b815260040161064e90612031565b5f61044861044460017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611c59565b6107b660017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611c59565b3380610bf261131d565b6001600160a01b031614610c1b578060405163118cdaa760e01b815260040161064e9190611650565b610c24816114c6565b50565b6107b660017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611c59565b6107b660017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611c59565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b5f61044861044460017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611c59565b610ce861144f565b6001600160a01b038116610d0e5760405162461bcd60e51b815260040161064e90611d6d565b610d3c61074c60017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611c59565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516107809190612073565b5f61044861044460017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611c59565b6107b660017fc9e8e7a4a583757cbcf624a50138f888cb585d449a8799952d3cc62760699622611c59565b610dcd61144f565b6001600160a01b038116610df35760405162461bcd60e51b815260040161064e90611d6d565b5f6001600160a01b031660018484604051610e0f929190611cbc565b908152604051908190036020019020546001600160a01b031614610e455760405162461bcd60e51b815260040161064e906120b5565b81610e625760405162461bcd60e51b815260040161064e906120f7565b5f6001600160a01b031660018484604051610e7e929190611cbc565b908152604051908190036020019020546001600160a01b031603610ed7575f80546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301610ed583858361219f565b505b8060018484604051610eea929190611cbc565b90815260405190819003602001812080546001600160a01b03939093166001600160a01b0319909316929092179091557f7ef997b0c9df3b39718be90c44d4d0d3d0230ac10eae31d63200210c7541ab7090610f4b9085908590859061225e565b60405180910390a1505050565b6107b660017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611c59565b60408051610120810182525f8082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820152815490919067ffffffffffffffff811115610fe657610fe6611842565b60405190808252806020026020018201604052801561102b57816020015b60408051808201909152606081525f60208201528152602001906001900390816110045790505b5090505f5b5f548110156111555760405180604001604052805f83815481106110565761105661227f565b905f5260205f2001805461106990611c80565b80601f016020809104026020016040519081016040528092919081815260200182805461109590611c80565b80156110e05780601f106110b7576101008083540402835291602001916110e0565b820191905f5260205f20905b8154815290600101906020018083116110c357829003601f168201915b5050505050815260200160015f84815481106110fe576110fe61227f565b905f5260205f20016040516111139190612302565b908152604051908190036020019020546001600160a01b0316905282518390839081106111425761114261227f565b6020908102919091010152600101611030565b5060405180610120016040528061116a610d6b565b6001600160a01b031681526020016111806113ca565b6001600160a01b031681526020016111966105c5565b6001600160a01b031681526020016111ac610b8e565b6001600160a01b031681526020016111c2610844565b6001600160a01b031681526020016111d8610cb1565b6001600160a01b031681526020016111ee6106c1565b6001600160a01b03168152602001611204610415565b6001600160a01b0316815260200191909152919050565b61122361144f565b61122e43604061230c565b811161124c5760405162461bcd60e51b815260040161064e90612351565b7f403a465e7990270029003553a58e56a0aa4f72cdf7948740eb4d5cf4056678c18585858585604051611283959493929190612361565b60405180910390a15050505050565b61129a61144f565b6001600160a01b0381166112c05760405162461bcd60e51b815260040161064e90611d6d565b6112ee61074c60017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611c59565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b738160405161078091906123d3565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610ca1565b61134d61144f565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b0319166001600160a01b0383169081178255611391610c7d565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b5f61044861044460017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611c59565b6107b660017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611c59565b6107b660017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611c59565b33611458610c7d565b6001600160a01b031614611481573360405163118cdaa760e01b815260040161064e9190611650565b565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6114b5611502565b6114be81611540565b610c24611551565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b03191681556114fe82611559565b5050565b61150a6115c9565b611481576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611548611502565b610c24816115e7565b611481611502565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f6115d2611483565b5468010000000000000000900460ff16919050565b6115ef611502565b6001600160a01b038116610c1b575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161064e9190611650565b5f6001600160a01b0382166114a7565b61164a81611631565b82525050565b602081016114a78284611641565b805b8114610c24575f5ffd5b80356114a78161165e565b5f60208284031215611688576116885f5ffd5b611692838361166a565b9392505050565b8281835e505f910152565b5f6116ad825190565b8084526020840193506116c4818560208601611699565b601f01601f19169290920192915050565b6020808252810161169281846116a4565b5f61169283836116a4565b60200190565b5f611700825190565b8084526020840193508360208202850161171a8560200190565b5f5b8481101561174d578383038852815161173584826116e6565b9350506020820160209890980197915060010161171c565b50909695505050505050565b6020808252810161169281846116f7565b5f5f83601f84011261177d5761177d5f5ffd5b50813567ffffffffffffffff811115611797576117975f5ffd5b6020830191508360018202830111156117b1576117b15f5ffd5b9250929050565b5f5f602083850312156117cc576117cc5f5ffd5b823567ffffffffffffffff8111156117e5576117e55f5ffd5b6117f18582860161176a565b92509250509250929050565b61166081611631565b80356114a7816117fd565b5f60208284031215611824576118245f5ffd5b6116928383611806565b8061164a565b602081016114a7828461182e565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff8211171561187c5761187c611842565b6040525050565b5f61188d60405190565b90506118998282611856565b919050565b5f608082840312156118b1576118b15f5ffd5b6118bb6080611883565b90506118c78383611806565b81526118d68360208401611806565b60208201526118e88360408401611806565b60408201526118fa8360608401611806565b606082015292915050565b5f5f60a08385031215611919576119195f5ffd5b611923848461189e565b91506119328460808501611806565b90509250929050565b5f67ffffffffffffffff82111561195457611954611842565b601f19601f83011660200192915050565b82818337505f910152565b5f61198261197d8461193b565b611883565b9050828152838383011115611998576119985f5ffd5b611692836020830184611965565b5f82601f8301126119b8576119b85f5ffd5b61169283833560208501611970565b5f602082840312156119da576119da5f5ffd5b813567ffffffffffffffff8111156119f3576119f35f5ffd5b6119ff848285016119a6565b949350505050565b5f5f5f60408486031215611a1c57611a1c5f5ffd5b833567ffffffffffffffff811115611a3557611a355f5ffd5b611a418682870161176a565b9350935050611a538560208601611806565b90509250925092565b805160408084525f9190840190611a7382826116a4565b9150506020830151611a886020860182611641565b509392505050565b5f6116928383611a5c565b5f611aa4825190565b80845260208401935083602082028501611abe8560200190565b5f5b8481101561174d5783830388528151611ad98482611a90565b93505060208201602098909801979150600101611ac0565b80515f90610120840190611b058582611641565b506020830151611b186020860182611641565b506040830151611b2b6040860182611641565b506060830151611b3e6060860182611641565b506080830151611b516080860182611641565b5060a0830151611b6460a0860182611641565b5060c0830151611b7760c0860182611641565b5060e0830151611b8a60e0860182611641565b50610100830151848203610100860152611ba48282611a9b565b95945050505050565b602080825281016116928184611af1565b5f5f5f5f5f60608688031215611bd557611bd55f5ffd5b853567ffffffffffffffff811115611bee57611bee5f5ffd5b611bfa8882890161176a565b9550955050602086013567ffffffffffffffff811115611c1b57611c1b5f5ffd5b611c278882890161176a565b9350935050611c39876040880161166a565b90509295509295909350565b634e487b7160e01b5f52601160045260245ffd5b818103818111156114a7576114a7611c45565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680611c9457607f821691505b602082108103611ca657611ca6611c6c565b50919050565b611cb7828483611965565b500190565b611692818385611cac565b60168152602081017f4164647265737320646f6573206e6f7420657869737400000000000000000000815290506116f1565b602080825281016114a781611cc7565b818352602083019250611d1d828483611965565b50601f01601f19160190565b602080825281016119ff818486611d09565b600f8152602081017f496e76616c696420616464726573730000000000000000000000000000000000815290506116f1565b602080825281016114a781611d3b565b60158152602081017f6c3143726f7373436861696e4d657373656e6765720000000000000000000000815290506116f1565b60408082528101611dbf81611d7d565b90506114a76020830184611641565b60088152602081017f6c32427269646765000000000000000000000000000000000000000000000000815290506116f1565b60408082528101611dbf81611dce565b60138152602081017f4f776e65722063616e6e6f742062652030783000000000000000000000000000815290506116f1565b602080825281016114a781611e10565b60198152602081017f43726f737320636861696e2063616e6e6f742062652030783000000000000000815290506116f1565b602080825281016114a781611e52565b60198152602081017f4d657373616765206275732063616e6e6f742062652030783000000000000000815290506116f1565b602080825281016114a781611e94565b60268152602081017f4e6574776f726b20656e636c6176652072656769737472792063616e6e6f742081527f6265203078300000000000000000000000000000000000000000000000000000602082015290505b60400190565b602080825281016114a781611ed6565b60288152602081017f4461746120617661696c6162696c6974792072656769737472792063616e6e6f81527f742062652030783000000000000000000000000000000000000000000000000060208201529050611f2a565b602080825281016114a781611f40565b5f6114a782611fb5565b90565b67ffffffffffffffff1690565b61164a81611fa8565b602081016114a78284611fc2565b60348152602081017f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f81527f742072656e6f756e6365206f776e65727368697000000000000000000000000060208201529050611f2a565b602080825281016114a781611fd9565b60158152602081017f6c3243726f7373436861696e4d657373656e6765720000000000000000000000815290506116f1565b60408082528101611dbf81612041565b60168152602081017f4164647265737320616c72656164792065786973747300000000000000000000815290506116f1565b602080825281016114a781612083565b60148152602081017f4e616d652063616e6e6f7420626520656d707479000000000000000000000000815290506116f1565b602080825281016114a7816120c5565b5f6114a7611fb28381565b61211b83612107565b81545f1960089490940293841b1916921b91909117905550565b5f612141818484612112565b505050565b818110156114fe576121585f82612135565b600101612146565b601f821115612141575f818152602090206020601f850104810160208510156121865750805b6121986020601f860104830182612146565b5050505050565b8267ffffffffffffffff8111156121b8576121b8611842565b6121c28254611c80565b6121cd828285612160565b505f601f8211600181146121ff575f83156121e85750848201355b5f19600885021c1981166002850217855550612256565b5f84815260208120601f198516915b8281101561222e578785013582556020948501946001909201910161220e565b508482101561224a575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b60408082528101612270818587611d09565b90506119ff6020830184611641565b634e487b7160e01b5f52603260045260245ffd5b5f815461229f81611c80565b6001821680156122b657600181146122cb576122f9565b60ff19831686528115158202860193506122f9565b5f858152602090205f5b838110156122f1578154888201526001909101906020016122d5565b505081860193505b50505092915050565b6114a78183612293565b808201808211156114a7576114a7611c45565b60148152602081017f496e76616c69642076616c69644174426c6f636b000000000000000000000000815290506116f1565b602080825281016114a78161231f565b60608082528101612373818789611d09565b90508181036020830152612388818587611d09565b9050612397604083018461182e565b9695505050505050565b60088152602081017f6c31427269646765000000000000000000000000000000000000000000000000815290506116f1565b60408082528101611dbf816123a156fea26469706673582212207cd3c55ea986e2569a4abf40a4a04fd21c92f9b8c86d4a1b5b6e2164633c2ac064736f6c634300081c0033",
}

// NetworkConfigABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkConfigMetaData.ABI instead.
var NetworkConfigABI = NetworkConfigMetaData.ABI

// NetworkConfigBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NetworkConfigMetaData.Bin instead.
var NetworkConfigBin = NetworkConfigMetaData.Bin

// DeployNetworkConfig deploys a new Ethereum contract, binding an instance of NetworkConfig to it.
func DeployNetworkConfig(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NetworkConfig, error) {
	parsed, err := NetworkConfigMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NetworkConfigBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NetworkConfig{NetworkConfigCaller: NetworkConfigCaller{contract: contract}, NetworkConfigTransactor: NetworkConfigTransactor{contract: contract}, NetworkConfigFilterer: NetworkConfigFilterer{contract: contract}}, nil
}

// NetworkConfig is an auto generated Go binding around an Ethereum contract.
type NetworkConfig struct {
	NetworkConfigCaller     // Read-only binding to the contract
	NetworkConfigTransactor // Write-only binding to the contract
	NetworkConfigFilterer   // Log filterer for contract events
}

// NetworkConfigCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkConfigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkConfigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkConfigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkConfigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkConfigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkConfigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkConfigSession struct {
	Contract     *NetworkConfig    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NetworkConfigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkConfigCallerSession struct {
	Contract *NetworkConfigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// NetworkConfigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkConfigTransactorSession struct {
	Contract     *NetworkConfigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// NetworkConfigRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkConfigRaw struct {
	Contract *NetworkConfig // Generic contract binding to access the raw methods on
}

// NetworkConfigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkConfigCallerRaw struct {
	Contract *NetworkConfigCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkConfigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkConfigTransactorRaw struct {
	Contract *NetworkConfigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkConfig creates a new instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfig(address common.Address, backend bind.ContractBackend) (*NetworkConfig, error) {
	contract, err := bindNetworkConfig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkConfig{NetworkConfigCaller: NetworkConfigCaller{contract: contract}, NetworkConfigTransactor: NetworkConfigTransactor{contract: contract}, NetworkConfigFilterer: NetworkConfigFilterer{contract: contract}}, nil
}

// NewNetworkConfigCaller creates a new read-only instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfigCaller(address common.Address, caller bind.ContractCaller) (*NetworkConfigCaller, error) {
	contract, err := bindNetworkConfig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigCaller{contract: contract}, nil
}

// NewNetworkConfigTransactor creates a new write-only instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfigTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkConfigTransactor, error) {
	contract, err := bindNetworkConfig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigTransactor{contract: contract}, nil
}

// NewNetworkConfigFilterer creates a new log filterer instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfigFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkConfigFilterer, error) {
	contract, err := bindNetworkConfig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigFilterer{contract: contract}, nil
}

// bindNetworkConfig binds a generic wrapper to an already deployed contract.
func bindNetworkConfig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NetworkConfigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkConfig *NetworkConfigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkConfig.Contract.NetworkConfigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkConfig *NetworkConfigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkConfig.Contract.NetworkConfigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkConfig *NetworkConfigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkConfig.Contract.NetworkConfigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkConfig *NetworkConfigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkConfig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkConfig *NetworkConfigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkConfig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkConfig *NetworkConfigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkConfig.Contract.contract.Transact(opts, method, params...)
}

// CROSSCHAINSLOT is a free data retrieval call binding the contract method 0xfbfd6d91.
//
// Solidity: function CROSS_CHAIN_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) CROSSCHAINSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "CROSS_CHAIN_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CROSSCHAINSLOT is a free data retrieval call binding the contract method 0xfbfd6d91.
//
// Solidity: function CROSS_CHAIN_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) CROSSCHAINSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.CROSSCHAINSLOT(&_NetworkConfig.CallOpts)
}

// CROSSCHAINSLOT is a free data retrieval call binding the contract method 0xfbfd6d91.
//
// Solidity: function CROSS_CHAIN_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) CROSSCHAINSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.CROSSCHAINSLOT(&_NetworkConfig.CallOpts)
}

// DATAAVAILABILITYREGISTRYSLOT is a free data retrieval call binding the contract method 0xfaa5e2de.
//
// Solidity: function DATA_AVAILABILITY_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) DATAAVAILABILITYREGISTRYSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "DATA_AVAILABILITY_REGISTRY_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DATAAVAILABILITYREGISTRYSLOT is a free data retrieval call binding the contract method 0xfaa5e2de.
//
// Solidity: function DATA_AVAILABILITY_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) DATAAVAILABILITYREGISTRYSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.DATAAVAILABILITYREGISTRYSLOT(&_NetworkConfig.CallOpts)
}

// DATAAVAILABILITYREGISTRYSLOT is a free data retrieval call binding the contract method 0xfaa5e2de.
//
// Solidity: function DATA_AVAILABILITY_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) DATAAVAILABILITYREGISTRYSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.DATAAVAILABILITYREGISTRYSLOT(&_NetworkConfig.CallOpts)
}

// FORKMANAGERSLOT is a free data retrieval call binding the contract method 0xae61ecba.
//
// Solidity: function FORK_MANAGER_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) FORKMANAGERSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "FORK_MANAGER_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FORKMANAGERSLOT is a free data retrieval call binding the contract method 0xae61ecba.
//
// Solidity: function FORK_MANAGER_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) FORKMANAGERSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.FORKMANAGERSLOT(&_NetworkConfig.CallOpts)
}

// FORKMANAGERSLOT is a free data retrieval call binding the contract method 0xae61ecba.
//
// Solidity: function FORK_MANAGER_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) FORKMANAGERSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.FORKMANAGERSLOT(&_NetworkConfig.CallOpts)
}

// L1BRIDGESLOT is a free data retrieval call binding the contract method 0xbe9f8207.
//
// Solidity: function L1_BRIDGE_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) L1BRIDGESLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "L1_BRIDGE_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// L1BRIDGESLOT is a free data retrieval call binding the contract method 0xbe9f8207.
//
// Solidity: function L1_BRIDGE_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) L1BRIDGESLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.L1BRIDGESLOT(&_NetworkConfig.CallOpts)
}

// L1BRIDGESLOT is a free data retrieval call binding the contract method 0xbe9f8207.
//
// Solidity: function L1_BRIDGE_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) L1BRIDGESLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.L1BRIDGESLOT(&_NetworkConfig.CallOpts)
}

// L1CROSSCHAINMESSENGER is a free data retrieval call binding the contract method 0x67cc852e.
//
// Solidity: function L1_CROSS_CHAIN_MESSENGER() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) L1CROSSCHAINMESSENGER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "L1_CROSS_CHAIN_MESSENGER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// L1CROSSCHAINMESSENGER is a free data retrieval call binding the contract method 0x67cc852e.
//
// Solidity: function L1_CROSS_CHAIN_MESSENGER() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) L1CROSSCHAINMESSENGER() ([32]byte, error) {
	return _NetworkConfig.Contract.L1CROSSCHAINMESSENGER(&_NetworkConfig.CallOpts)
}

// L1CROSSCHAINMESSENGER is a free data retrieval call binding the contract method 0x67cc852e.
//
// Solidity: function L1_CROSS_CHAIN_MESSENGER() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) L1CROSSCHAINMESSENGER() ([32]byte, error) {
	return _NetworkConfig.Contract.L1CROSSCHAINMESSENGER(&_NetworkConfig.CallOpts)
}

// L2BRIDGESLOT is a free data retrieval call binding the contract method 0x812b1ffe.
//
// Solidity: function L2_BRIDGE_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) L2BRIDGESLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "L2_BRIDGE_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// L2BRIDGESLOT is a free data retrieval call binding the contract method 0x812b1ffe.
//
// Solidity: function L2_BRIDGE_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) L2BRIDGESLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.L2BRIDGESLOT(&_NetworkConfig.CallOpts)
}

// L2BRIDGESLOT is a free data retrieval call binding the contract method 0x812b1ffe.
//
// Solidity: function L2_BRIDGE_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) L2BRIDGESLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.L2BRIDGESLOT(&_NetworkConfig.CallOpts)
}

// L2CROSSCHAINMESSENGER is a free data retrieval call binding the contract method 0x85f427cb.
//
// Solidity: function L2_CROSS_CHAIN_MESSENGER() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) L2CROSSCHAINMESSENGER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "L2_CROSS_CHAIN_MESSENGER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// L2CROSSCHAINMESSENGER is a free data retrieval call binding the contract method 0x85f427cb.
//
// Solidity: function L2_CROSS_CHAIN_MESSENGER() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) L2CROSSCHAINMESSENGER() ([32]byte, error) {
	return _NetworkConfig.Contract.L2CROSSCHAINMESSENGER(&_NetworkConfig.CallOpts)
}

// L2CROSSCHAINMESSENGER is a free data retrieval call binding the contract method 0x85f427cb.
//
// Solidity: function L2_CROSS_CHAIN_MESSENGER() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) L2CROSSCHAINMESSENGER() ([32]byte, error) {
	return _NetworkConfig.Contract.L2CROSSCHAINMESSENGER(&_NetworkConfig.CallOpts)
}

// MESSAGEBUSSLOT is a free data retrieval call binding the contract method 0x48d87239.
//
// Solidity: function MESSAGE_BUS_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) MESSAGEBUSSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "MESSAGE_BUS_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MESSAGEBUSSLOT is a free data retrieval call binding the contract method 0x48d87239.
//
// Solidity: function MESSAGE_BUS_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) MESSAGEBUSSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.MESSAGEBUSSLOT(&_NetworkConfig.CallOpts)
}

// MESSAGEBUSSLOT is a free data retrieval call binding the contract method 0x48d87239.
//
// Solidity: function MESSAGE_BUS_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) MESSAGEBUSSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.MESSAGEBUSSLOT(&_NetworkConfig.CallOpts)
}

// NETWORKENCLAVEREGISTRYSLOT is a free data retrieval call binding the contract method 0x72bad912.
//
// Solidity: function NETWORK_ENCLAVE_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) NETWORKENCLAVEREGISTRYSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "NETWORK_ENCLAVE_REGISTRY_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NETWORKENCLAVEREGISTRYSLOT is a free data retrieval call binding the contract method 0x72bad912.
//
// Solidity: function NETWORK_ENCLAVE_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) NETWORKENCLAVEREGISTRYSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.NETWORKENCLAVEREGISTRYSLOT(&_NetworkConfig.CallOpts)
}

// NETWORKENCLAVEREGISTRYSLOT is a free data retrieval call binding the contract method 0x72bad912.
//
// Solidity: function NETWORK_ENCLAVE_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) NETWORKENCLAVEREGISTRYSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.NETWORKENCLAVEREGISTRYSLOT(&_NetworkConfig.CallOpts)
}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string contractName) view returns(address contractAddress)
func (_NetworkConfig *NetworkConfigCaller) AdditionalAddresses(opts *bind.CallOpts, contractName string) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "additionalAddresses", contractName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string contractName) view returns(address contractAddress)
func (_NetworkConfig *NetworkConfigSession) AdditionalAddresses(contractName string) (common.Address, error) {
	return _NetworkConfig.Contract.AdditionalAddresses(&_NetworkConfig.CallOpts, contractName)
}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string contractName) view returns(address contractAddress)
func (_NetworkConfig *NetworkConfigCallerSession) AdditionalAddresses(contractName string) (common.Address, error) {
	return _NetworkConfig.Contract.AdditionalAddresses(&_NetworkConfig.CallOpts, contractName)
}

// AddressNames is a free data retrieval call binding the contract method 0x0f387b1e.
//
// Solidity: function addressNames(uint256 ) view returns(string)
func (_NetworkConfig *NetworkConfigCaller) AddressNames(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "addressNames", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AddressNames is a free data retrieval call binding the contract method 0x0f387b1e.
//
// Solidity: function addressNames(uint256 ) view returns(string)
func (_NetworkConfig *NetworkConfigSession) AddressNames(arg0 *big.Int) (string, error) {
	return _NetworkConfig.Contract.AddressNames(&_NetworkConfig.CallOpts, arg0)
}

// AddressNames is a free data retrieval call binding the contract method 0x0f387b1e.
//
// Solidity: function addressNames(uint256 ) view returns(string)
func (_NetworkConfig *NetworkConfigCallerSession) AddressNames(arg0 *big.Int) (string, error) {
	return _NetworkConfig.Contract.AddressNames(&_NetworkConfig.CallOpts, arg0)
}

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address,address,address,address,address,(string,address)[]))
func (_NetworkConfig *NetworkConfigCaller) Addresses(opts *bind.CallOpts) (NetworkConfigAddresses, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "addresses")

	if err != nil {
		return *new(NetworkConfigAddresses), err
	}

	out0 := *abi.ConvertType(out[0], new(NetworkConfigAddresses)).(*NetworkConfigAddresses)

	return out0, err

}

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address,address,address,address,address,(string,address)[]))
func (_NetworkConfig *NetworkConfigSession) Addresses() (NetworkConfigAddresses, error) {
	return _NetworkConfig.Contract.Addresses(&_NetworkConfig.CallOpts)
}

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address,address,address,address,address,(string,address)[]))
func (_NetworkConfig *NetworkConfigCallerSession) Addresses() (NetworkConfigAddresses, error) {
	return _NetworkConfig.Contract.Addresses(&_NetworkConfig.CallOpts)
}

// CrossChainContractAddress is a free data retrieval call binding the contract method 0xa1b918d6.
//
// Solidity: function crossChainContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) CrossChainContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "crossChainContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainContractAddress is a free data retrieval call binding the contract method 0xa1b918d6.
//
// Solidity: function crossChainContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) CrossChainContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.CrossChainContractAddress(&_NetworkConfig.CallOpts)
}

// CrossChainContractAddress is a free data retrieval call binding the contract method 0xa1b918d6.
//
// Solidity: function crossChainContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) CrossChainContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.CrossChainContractAddress(&_NetworkConfig.CallOpts)
}

// DaRegistryContractAddress is a free data retrieval call binding the contract method 0x71fd11f3.
//
// Solidity: function daRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) DaRegistryContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "daRegistryContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DaRegistryContractAddress is a free data retrieval call binding the contract method 0x71fd11f3.
//
// Solidity: function daRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) DaRegistryContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.DaRegistryContractAddress(&_NetworkConfig.CallOpts)
}

// DaRegistryContractAddress is a free data retrieval call binding the contract method 0x71fd11f3.
//
// Solidity: function daRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) DaRegistryContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.DaRegistryContractAddress(&_NetworkConfig.CallOpts)
}

// GetAdditionalContractNames is a free data retrieval call binding the contract method 0x13eeee96.
//
// Solidity: function getAdditionalContractNames() view returns(string[])
func (_NetworkConfig *NetworkConfigCaller) GetAdditionalContractNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "getAdditionalContractNames")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetAdditionalContractNames is a free data retrieval call binding the contract method 0x13eeee96.
//
// Solidity: function getAdditionalContractNames() view returns(string[])
func (_NetworkConfig *NetworkConfigSession) GetAdditionalContractNames() ([]string, error) {
	return _NetworkConfig.Contract.GetAdditionalContractNames(&_NetworkConfig.CallOpts)
}

// GetAdditionalContractNames is a free data retrieval call binding the contract method 0x13eeee96.
//
// Solidity: function getAdditionalContractNames() view returns(string[])
func (_NetworkConfig *NetworkConfigCallerSession) GetAdditionalContractNames() ([]string, error) {
	return _NetworkConfig.Contract.GetAdditionalContractNames(&_NetworkConfig.CallOpts)
}

// L1BridgeAddress is a free data retrieval call binding the contract method 0x5ab2a558.
//
// Solidity: function l1BridgeAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) L1BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "l1BridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1BridgeAddress is a free data retrieval call binding the contract method 0x5ab2a558.
//
// Solidity: function l1BridgeAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) L1BridgeAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L1BridgeAddress(&_NetworkConfig.CallOpts)
}

// L1BridgeAddress is a free data retrieval call binding the contract method 0x5ab2a558.
//
// Solidity: function l1BridgeAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) L1BridgeAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L1BridgeAddress(&_NetworkConfig.CallOpts)
}

// L1CrossChainMessengerAddress is a free data retrieval call binding the contract method 0x450948ad.
//
// Solidity: function l1CrossChainMessengerAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) L1CrossChainMessengerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "l1CrossChainMessengerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L1CrossChainMessengerAddress is a free data retrieval call binding the contract method 0x450948ad.
//
// Solidity: function l1CrossChainMessengerAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) L1CrossChainMessengerAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L1CrossChainMessengerAddress(&_NetworkConfig.CallOpts)
}

// L1CrossChainMessengerAddress is a free data retrieval call binding the contract method 0x450948ad.
//
// Solidity: function l1CrossChainMessengerAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) L1CrossChainMessengerAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L1CrossChainMessengerAddress(&_NetworkConfig.CallOpts)
}

// L2BridgeAddress is a free data retrieval call binding the contract method 0x934746a7.
//
// Solidity: function l2BridgeAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) L2BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "l2BridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2BridgeAddress is a free data retrieval call binding the contract method 0x934746a7.
//
// Solidity: function l2BridgeAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) L2BridgeAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L2BridgeAddress(&_NetworkConfig.CallOpts)
}

// L2BridgeAddress is a free data retrieval call binding the contract method 0x934746a7.
//
// Solidity: function l2BridgeAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) L2BridgeAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L2BridgeAddress(&_NetworkConfig.CallOpts)
}

// L2CrossChainMessengerAddress is a free data retrieval call binding the contract method 0x0b592f45.
//
// Solidity: function l2CrossChainMessengerAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) L2CrossChainMessengerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "l2CrossChainMessengerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2CrossChainMessengerAddress is a free data retrieval call binding the contract method 0x0b592f45.
//
// Solidity: function l2CrossChainMessengerAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) L2CrossChainMessengerAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L2CrossChainMessengerAddress(&_NetworkConfig.CallOpts)
}

// L2CrossChainMessengerAddress is a free data retrieval call binding the contract method 0x0b592f45.
//
// Solidity: function l2CrossChainMessengerAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) L2CrossChainMessengerAddress() (common.Address, error) {
	return _NetworkConfig.Contract.L2CrossChainMessengerAddress(&_NetworkConfig.CallOpts)
}

// MessageBusContractAddress is a free data retrieval call binding the contract method 0xf5e9f286.
//
// Solidity: function messageBusContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) MessageBusContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "messageBusContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBusContractAddress is a free data retrieval call binding the contract method 0xf5e9f286.
//
// Solidity: function messageBusContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) MessageBusContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.MessageBusContractAddress(&_NetworkConfig.CallOpts)
}

// MessageBusContractAddress is a free data retrieval call binding the contract method 0xf5e9f286.
//
// Solidity: function messageBusContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) MessageBusContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.MessageBusContractAddress(&_NetworkConfig.CallOpts)
}

// NetworkEnclaveRegistryContractAddress is a free data retrieval call binding the contract method 0x2fc00c76.
//
// Solidity: function networkEnclaveRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) NetworkEnclaveRegistryContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "networkEnclaveRegistryContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NetworkEnclaveRegistryContractAddress is a free data retrieval call binding the contract method 0x2fc00c76.
//
// Solidity: function networkEnclaveRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) NetworkEnclaveRegistryContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.NetworkEnclaveRegistryContractAddress(&_NetworkConfig.CallOpts)
}

// NetworkEnclaveRegistryContractAddress is a free data retrieval call binding the contract method 0x2fc00c76.
//
// Solidity: function networkEnclaveRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) NetworkEnclaveRegistryContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.NetworkEnclaveRegistryContractAddress(&_NetworkConfig.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkConfig *NetworkConfigCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkConfig *NetworkConfigSession) Owner() (common.Address, error) {
	return _NetworkConfig.Contract.Owner(&_NetworkConfig.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkConfig *NetworkConfigCallerSession) Owner() (common.Address, error) {
	return _NetworkConfig.Contract.Owner(&_NetworkConfig.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_NetworkConfig *NetworkConfigCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_NetworkConfig *NetworkConfigSession) PendingOwner() (common.Address, error) {
	return _NetworkConfig.Contract.PendingOwner(&_NetworkConfig.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_NetworkConfig *NetworkConfigCallerSession) PendingOwner() (common.Address, error) {
	return _NetworkConfig.Contract.PendingOwner(&_NetworkConfig.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_NetworkConfig *NetworkConfigCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_NetworkConfig *NetworkConfigSession) RenounceOwnership() error {
	return _NetworkConfig.Contract.RenounceOwnership(&_NetworkConfig.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_NetworkConfig *NetworkConfigCallerSession) RenounceOwnership() error {
	return _NetworkConfig.Contract.RenounceOwnership(&_NetworkConfig.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_NetworkConfig *NetworkConfigTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_NetworkConfig *NetworkConfigSession) AcceptOwnership() (*types.Transaction, error) {
	return _NetworkConfig.Contract.AcceptOwnership(&_NetworkConfig.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_NetworkConfig *NetworkConfigTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _NetworkConfig.Contract.AcceptOwnership(&_NetworkConfig.TransactOpts)
}

// AddAdditionalAddress is a paid mutator transaction binding the contract method 0xb7bef9ab.
//
// Solidity: function addAdditionalAddress(string name, address addr) returns()
func (_NetworkConfig *NetworkConfigTransactor) AddAdditionalAddress(opts *bind.TransactOpts, name string, addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "addAdditionalAddress", name, addr)
}

// AddAdditionalAddress is a paid mutator transaction binding the contract method 0xb7bef9ab.
//
// Solidity: function addAdditionalAddress(string name, address addr) returns()
func (_NetworkConfig *NetworkConfigSession) AddAdditionalAddress(name string, addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.AddAdditionalAddress(&_NetworkConfig.TransactOpts, name, addr)
}

// AddAdditionalAddress is a paid mutator transaction binding the contract method 0xb7bef9ab.
//
// Solidity: function addAdditionalAddress(string name, address addr) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) AddAdditionalAddress(name string, addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.AddAdditionalAddress(&_NetworkConfig.TransactOpts, name, addr)
}

// Initialize is a paid mutator transaction binding the contract method 0x6c1358ac.
//
// Solidity: function initialize((address,address,address,address) _addresses, address _owner) returns()
func (_NetworkConfig *NetworkConfigTransactor) Initialize(opts *bind.TransactOpts, _addresses NetworkConfigFixedAddresses, _owner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "initialize", _addresses, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x6c1358ac.
//
// Solidity: function initialize((address,address,address,address) _addresses, address _owner) returns()
func (_NetworkConfig *NetworkConfigSession) Initialize(_addresses NetworkConfigFixedAddresses, _owner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x6c1358ac.
//
// Solidity: function initialize((address,address,address,address) _addresses, address _owner) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) Initialize(_addresses NetworkConfigFixedAddresses, _owner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses, _owner)
}

// RemoveAdditionalAddress is a paid mutator transaction binding the contract method 0x31d1464d.
//
// Solidity: function removeAdditionalAddress(string name) returns()
func (_NetworkConfig *NetworkConfigTransactor) RemoveAdditionalAddress(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "removeAdditionalAddress", name)
}

// RemoveAdditionalAddress is a paid mutator transaction binding the contract method 0x31d1464d.
//
// Solidity: function removeAdditionalAddress(string name) returns()
func (_NetworkConfig *NetworkConfigSession) RemoveAdditionalAddress(name string) (*types.Transaction, error) {
	return _NetworkConfig.Contract.RemoveAdditionalAddress(&_NetworkConfig.TransactOpts, name)
}

// RemoveAdditionalAddress is a paid mutator transaction binding the contract method 0x31d1464d.
//
// Solidity: function removeAdditionalAddress(string name) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) RemoveAdditionalAddress(name string) (*types.Transaction, error) {
	return _NetworkConfig.Contract.RemoveAdditionalAddress(&_NetworkConfig.TransactOpts, name)
}

// SetL1BridgeAddress is a paid mutator transaction binding the contract method 0xe1825d06.
//
// Solidity: function setL1BridgeAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactor) SetL1BridgeAddress(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "setL1BridgeAddress", _addr)
}

// SetL1BridgeAddress is a paid mutator transaction binding the contract method 0xe1825d06.
//
// Solidity: function setL1BridgeAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigSession) SetL1BridgeAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL1BridgeAddress(&_NetworkConfig.TransactOpts, _addr)
}

// SetL1BridgeAddress is a paid mutator transaction binding the contract method 0xe1825d06.
//
// Solidity: function setL1BridgeAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) SetL1BridgeAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL1BridgeAddress(&_NetworkConfig.TransactOpts, _addr)
}

// SetL1CrossChainMessengerAddress is a paid mutator transaction binding the contract method 0x46a30a78.
//
// Solidity: function setL1CrossChainMessengerAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactor) SetL1CrossChainMessengerAddress(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "setL1CrossChainMessengerAddress", _addr)
}

// SetL1CrossChainMessengerAddress is a paid mutator transaction binding the contract method 0x46a30a78.
//
// Solidity: function setL1CrossChainMessengerAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigSession) SetL1CrossChainMessengerAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL1CrossChainMessengerAddress(&_NetworkConfig.TransactOpts, _addr)
}

// SetL1CrossChainMessengerAddress is a paid mutator transaction binding the contract method 0x46a30a78.
//
// Solidity: function setL1CrossChainMessengerAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) SetL1CrossChainMessengerAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL1CrossChainMessengerAddress(&_NetworkConfig.TransactOpts, _addr)
}

// SetL2BridgeAddress is a paid mutator transaction binding the contract method 0x556d89dd.
//
// Solidity: function setL2BridgeAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactor) SetL2BridgeAddress(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "setL2BridgeAddress", _addr)
}

// SetL2BridgeAddress is a paid mutator transaction binding the contract method 0x556d89dd.
//
// Solidity: function setL2BridgeAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigSession) SetL2BridgeAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL2BridgeAddress(&_NetworkConfig.TransactOpts, _addr)
}

// SetL2BridgeAddress is a paid mutator transaction binding the contract method 0x556d89dd.
//
// Solidity: function setL2BridgeAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) SetL2BridgeAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL2BridgeAddress(&_NetworkConfig.TransactOpts, _addr)
}

// SetL2CrossChainMessengerAddress is a paid mutator transaction binding the contract method 0x96493cc5.
//
// Solidity: function setL2CrossChainMessengerAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactor) SetL2CrossChainMessengerAddress(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "setL2CrossChainMessengerAddress", _addr)
}

// SetL2CrossChainMessengerAddress is a paid mutator transaction binding the contract method 0x96493cc5.
//
// Solidity: function setL2CrossChainMessengerAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigSession) SetL2CrossChainMessengerAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL2CrossChainMessengerAddress(&_NetworkConfig.TransactOpts, _addr)
}

// SetL2CrossChainMessengerAddress is a paid mutator transaction binding the contract method 0x96493cc5.
//
// Solidity: function setL2CrossChainMessengerAddress(address _addr) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) SetL2CrossChainMessengerAddress(_addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.SetL2CrossChainMessengerAddress(&_NetworkConfig.TransactOpts, _addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkConfig *NetworkConfigTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkConfig *NetworkConfigSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.TransferOwnership(&_NetworkConfig.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.TransferOwnership(&_NetworkConfig.TransactOpts, newOwner)
}

// UpgradeFeature is a paid mutator transaction binding the contract method 0xdbf5d903.
//
// Solidity: function upgradeFeature(string featureName, bytes featureData, uint256 validAtBlock) returns()
func (_NetworkConfig *NetworkConfigTransactor) UpgradeFeature(opts *bind.TransactOpts, featureName string, featureData []byte, validAtBlock *big.Int) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "upgradeFeature", featureName, featureData, validAtBlock)
}

// UpgradeFeature is a paid mutator transaction binding the contract method 0xdbf5d903.
//
// Solidity: function upgradeFeature(string featureName, bytes featureData, uint256 validAtBlock) returns()
func (_NetworkConfig *NetworkConfigSession) UpgradeFeature(featureName string, featureData []byte, validAtBlock *big.Int) (*types.Transaction, error) {
	return _NetworkConfig.Contract.UpgradeFeature(&_NetworkConfig.TransactOpts, featureName, featureData, validAtBlock)
}

// UpgradeFeature is a paid mutator transaction binding the contract method 0xdbf5d903.
//
// Solidity: function upgradeFeature(string featureName, bytes featureData, uint256 validAtBlock) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) UpgradeFeature(featureName string, featureData []byte, validAtBlock *big.Int) (*types.Transaction, error) {
	return _NetworkConfig.Contract.UpgradeFeature(&_NetworkConfig.TransactOpts, featureName, featureData, validAtBlock)
}

// NetworkConfigAdditionalContractAddressAddedIterator is returned from FilterAdditionalContractAddressAdded and is used to iterate over the raw logs and unpacked data for AdditionalContractAddressAdded events raised by the NetworkConfig contract.
type NetworkConfigAdditionalContractAddressAddedIterator struct {
	Event *NetworkConfigAdditionalContractAddressAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkConfigAdditionalContractAddressAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigAdditionalContractAddressAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkConfigAdditionalContractAddressAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkConfigAdditionalContractAddressAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigAdditionalContractAddressAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigAdditionalContractAddressAdded represents a AdditionalContractAddressAdded event raised by the NetworkConfig contract.
type NetworkConfigAdditionalContractAddressAdded struct {
	Name string
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAdditionalContractAddressAdded is a free log retrieval operation binding the contract event 0x7ef997b0c9df3b39718be90c44d4d0d3d0230ac10eae31d63200210c7541ab70.
//
// Solidity: event AdditionalContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) FilterAdditionalContractAddressAdded(opts *bind.FilterOpts) (*NetworkConfigAdditionalContractAddressAddedIterator, error) {

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "AdditionalContractAddressAdded")
	if err != nil {
		return nil, err
	}
	return &NetworkConfigAdditionalContractAddressAddedIterator{contract: _NetworkConfig.contract, event: "AdditionalContractAddressAdded", logs: logs, sub: sub}, nil
}

// WatchAdditionalContractAddressAdded is a free log subscription operation binding the contract event 0x7ef997b0c9df3b39718be90c44d4d0d3d0230ac10eae31d63200210c7541ab70.
//
// Solidity: event AdditionalContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) WatchAdditionalContractAddressAdded(opts *bind.WatchOpts, sink chan<- *NetworkConfigAdditionalContractAddressAdded) (event.Subscription, error) {

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "AdditionalContractAddressAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigAdditionalContractAddressAdded)
				if err := _NetworkConfig.contract.UnpackLog(event, "AdditionalContractAddressAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdditionalContractAddressAdded is a log parse operation binding the contract event 0x7ef997b0c9df3b39718be90c44d4d0d3d0230ac10eae31d63200210c7541ab70.
//
// Solidity: event AdditionalContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) ParseAdditionalContractAddressAdded(log types.Log) (*NetworkConfigAdditionalContractAddressAdded, error) {
	event := new(NetworkConfigAdditionalContractAddressAdded)
	if err := _NetworkConfig.contract.UnpackLog(event, "AdditionalContractAddressAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigAdditionalContractAddressRemovedIterator is returned from FilterAdditionalContractAddressRemoved and is used to iterate over the raw logs and unpacked data for AdditionalContractAddressRemoved events raised by the NetworkConfig contract.
type NetworkConfigAdditionalContractAddressRemovedIterator struct {
	Event *NetworkConfigAdditionalContractAddressRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkConfigAdditionalContractAddressRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigAdditionalContractAddressRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkConfigAdditionalContractAddressRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkConfigAdditionalContractAddressRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigAdditionalContractAddressRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigAdditionalContractAddressRemoved represents a AdditionalContractAddressRemoved event raised by the NetworkConfig contract.
type NetworkConfigAdditionalContractAddressRemoved struct {
	Name string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAdditionalContractAddressRemoved is a free log retrieval operation binding the contract event 0x5f9e4bc50ed4fc3c0bf14c4b518e1f8132c7a95ce5bed2e97c8675e6adf03573.
//
// Solidity: event AdditionalContractAddressRemoved(string name)
func (_NetworkConfig *NetworkConfigFilterer) FilterAdditionalContractAddressRemoved(opts *bind.FilterOpts) (*NetworkConfigAdditionalContractAddressRemovedIterator, error) {

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "AdditionalContractAddressRemoved")
	if err != nil {
		return nil, err
	}
	return &NetworkConfigAdditionalContractAddressRemovedIterator{contract: _NetworkConfig.contract, event: "AdditionalContractAddressRemoved", logs: logs, sub: sub}, nil
}

// WatchAdditionalContractAddressRemoved is a free log subscription operation binding the contract event 0x5f9e4bc50ed4fc3c0bf14c4b518e1f8132c7a95ce5bed2e97c8675e6adf03573.
//
// Solidity: event AdditionalContractAddressRemoved(string name)
func (_NetworkConfig *NetworkConfigFilterer) WatchAdditionalContractAddressRemoved(opts *bind.WatchOpts, sink chan<- *NetworkConfigAdditionalContractAddressRemoved) (event.Subscription, error) {

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "AdditionalContractAddressRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigAdditionalContractAddressRemoved)
				if err := _NetworkConfig.contract.UnpackLog(event, "AdditionalContractAddressRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdditionalContractAddressRemoved is a log parse operation binding the contract event 0x5f9e4bc50ed4fc3c0bf14c4b518e1f8132c7a95ce5bed2e97c8675e6adf03573.
//
// Solidity: event AdditionalContractAddressRemoved(string name)
func (_NetworkConfig *NetworkConfigFilterer) ParseAdditionalContractAddressRemoved(log types.Log) (*NetworkConfigAdditionalContractAddressRemoved, error) {
	event := new(NetworkConfigAdditionalContractAddressRemoved)
	if err := _NetworkConfig.contract.UnpackLog(event, "AdditionalContractAddressRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NetworkConfig contract.
type NetworkConfigInitializedIterator struct {
	Event *NetworkConfigInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkConfigInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkConfigInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkConfigInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigInitialized represents a Initialized event raised by the NetworkConfig contract.
type NetworkConfigInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkConfig *NetworkConfigFilterer) FilterInitialized(opts *bind.FilterOpts) (*NetworkConfigInitializedIterator, error) {

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NetworkConfigInitializedIterator{contract: _NetworkConfig.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkConfig *NetworkConfigFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NetworkConfigInitialized) (event.Subscription, error) {

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigInitialized)
				if err := _NetworkConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkConfig *NetworkConfigFilterer) ParseInitialized(log types.Log) (*NetworkConfigInitialized, error) {
	event := new(NetworkConfigInitialized)
	if err := _NetworkConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigNetworkContractAddressAddedIterator is returned from FilterNetworkContractAddressAdded and is used to iterate over the raw logs and unpacked data for NetworkContractAddressAdded events raised by the NetworkConfig contract.
type NetworkConfigNetworkContractAddressAddedIterator struct {
	Event *NetworkConfigNetworkContractAddressAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkConfigNetworkContractAddressAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigNetworkContractAddressAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkConfigNetworkContractAddressAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkConfigNetworkContractAddressAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigNetworkContractAddressAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigNetworkContractAddressAdded represents a NetworkContractAddressAdded event raised by the NetworkConfig contract.
type NetworkConfigNetworkContractAddressAdded struct {
	Name string
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNetworkContractAddressAdded is a free log retrieval operation binding the contract event 0x8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73.
//
// Solidity: event NetworkContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) FilterNetworkContractAddressAdded(opts *bind.FilterOpts) (*NetworkConfigNetworkContractAddressAddedIterator, error) {

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "NetworkContractAddressAdded")
	if err != nil {
		return nil, err
	}
	return &NetworkConfigNetworkContractAddressAddedIterator{contract: _NetworkConfig.contract, event: "NetworkContractAddressAdded", logs: logs, sub: sub}, nil
}

// WatchNetworkContractAddressAdded is a free log subscription operation binding the contract event 0x8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73.
//
// Solidity: event NetworkContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) WatchNetworkContractAddressAdded(opts *bind.WatchOpts, sink chan<- *NetworkConfigNetworkContractAddressAdded) (event.Subscription, error) {

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "NetworkContractAddressAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigNetworkContractAddressAdded)
				if err := _NetworkConfig.contract.UnpackLog(event, "NetworkContractAddressAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNetworkContractAddressAdded is a log parse operation binding the contract event 0x8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73.
//
// Solidity: event NetworkContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) ParseNetworkContractAddressAdded(log types.Log) (*NetworkConfigNetworkContractAddressAdded, error) {
	event := new(NetworkConfigNetworkContractAddressAdded)
	if err := _NetworkConfig.contract.UnpackLog(event, "NetworkContractAddressAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the NetworkConfig contract.
type NetworkConfigOwnershipTransferStartedIterator struct {
	Event *NetworkConfigOwnershipTransferStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkConfigOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigOwnershipTransferStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkConfigOwnershipTransferStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkConfigOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the NetworkConfig contract.
type NetworkConfigOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NetworkConfigOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigOwnershipTransferStartedIterator{contract: _NetworkConfig.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *NetworkConfigOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigOwnershipTransferStarted)
				if err := _NetworkConfig.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) ParseOwnershipTransferStarted(log types.Log) (*NetworkConfigOwnershipTransferStarted, error) {
	event := new(NetworkConfigOwnershipTransferStarted)
	if err := _NetworkConfig.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NetworkConfig contract.
type NetworkConfigOwnershipTransferredIterator struct {
	Event *NetworkConfigOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkConfigOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkConfigOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkConfigOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigOwnershipTransferred represents a OwnershipTransferred event raised by the NetworkConfig contract.
type NetworkConfigOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NetworkConfigOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigOwnershipTransferredIterator{contract: _NetworkConfig.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NetworkConfigOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigOwnershipTransferred)
				if err := _NetworkConfig.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) ParseOwnershipTransferred(log types.Log) (*NetworkConfigOwnershipTransferred, error) {
	event := new(NetworkConfigOwnershipTransferred)
	if err := _NetworkConfig.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the NetworkConfig contract.
type NetworkConfigUpgradedIterator struct {
	Event *NetworkConfigUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NetworkConfigUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NetworkConfigUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NetworkConfigUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigUpgraded represents a Upgraded event raised by the NetworkConfig contract.
type NetworkConfigUpgraded struct {
	FeatureName  string
	FeatureData  []byte
	ValidAtBlock *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0x403a465e7990270029003553a58e56a0aa4f72cdf7948740eb4d5cf4056678c1.
//
// Solidity: event Upgraded(string featureName, bytes featureData, uint256 validAtBlock)
func (_NetworkConfig *NetworkConfigFilterer) FilterUpgraded(opts *bind.FilterOpts) (*NetworkConfigUpgradedIterator, error) {

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "Upgraded")
	if err != nil {
		return nil, err
	}
	return &NetworkConfigUpgradedIterator{contract: _NetworkConfig.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0x403a465e7990270029003553a58e56a0aa4f72cdf7948740eb4d5cf4056678c1.
//
// Solidity: event Upgraded(string featureName, bytes featureData, uint256 validAtBlock)
func (_NetworkConfig *NetworkConfigFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *NetworkConfigUpgraded) (event.Subscription, error) {

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "Upgraded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigUpgraded)
				if err := _NetworkConfig.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0x403a465e7990270029003553a58e56a0aa4f72cdf7948740eb4d5cf4056678c1.
//
// Solidity: event Upgraded(string featureName, bytes featureData, uint256 validAtBlock)
func (_NetworkConfig *NetworkConfigFilterer) ParseUpgraded(log types.Log) (*NetworkConfigUpgraded, error) {
	event := new(NetworkConfigUpgraded)
	if err := _NetworkConfig.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
