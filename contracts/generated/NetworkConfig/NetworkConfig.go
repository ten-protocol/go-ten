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

// NetworkConfigContractVersion is an auto generated low-level Go binding around an user-defined struct.
type NetworkConfigContractVersion struct {
	Name           string
	Version        string
	Implementation common.Address
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
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"AdditionalContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"AdditionalContractAddressRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"forkName\",\"type\":\"string\"}],\"name\":\"HardforkUpgrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NetworkContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DATA_AVAILABILITY_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FORK_MANAGER_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MESSAGE_BUS_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NETWORK_ENCLAVE_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addAdditionalAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractName\",\"type\":\"string\"}],\"name\":\"additionalAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"addressNames\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addresses\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dataAvailabilityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1CrossChainMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2CrossChainMessenger\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.NamedAddress[]\",\"name\":\"additionalContracts\",\"type\":\"tuple[]\"}],\"internalType\":\"structNetworkConfig.Addresses\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdditionaContractNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getContractVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.ContractVersion\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dataAvailabilityRegistry\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.FixedAddresses\",\"name\":\"_addresses\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkEnclaveRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hardforkName\",\"type\":\"string\"}],\"name\":\"recordHardfork\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"removeAdditionalAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b5061001861001d565b6100fc565b5f6100266100bd565b805490915068010000000000000000900460ff16156100585760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100ba5780546001600160401b0319166001600160401b0390811782556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916100b1916100e7565b60405180910390a15b50565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6001600160401b0382168152602081016100e1565b612217806101095f395ff3fe608060405234801561000f575f5ffd5b5060043610610235575f3560e01c80638da5cb5b1161013d578063be9f8207116100b8578063ed11f78d11610088578063f5e9f2861161006e578063f5e9f28614610437578063faa5e2de1461043f578063fbfd6d9114610447575f5ffd5b8063ed11f78d14610404578063f2fde38b14610424575f5ffd5b8063be9f8207146103cc578063da0321cd146103d4578063e1825d06146103e9578063e30c3978146103fc575f5ffd5b8063adc5c2071161010d578063af454635116100f3578063af45463514610370578063b7bef9ab146103a4578063bc162c90146103b7575f5ffd5b8063adc5c20714610355578063ae61ecba14610368575f5ffd5b80638da5cb5b1461032a578063934746a71461033257806396493cc51461033a578063a1b918d61461034d575f5ffd5b80635ab2a558116101cd57806371fd11f31161019d57806379ba50971161018357806379ba509714610312578063812b1ffe1461031a57806385f427cb14610322575f5ffd5b806371fd11f31461030257806372bad9121461030a575f5ffd5b80635ab2a558146102d757806367cc852e146102df5780636c1358ac146102e7578063715018a6146102fa575f5ffd5b8063450948ad11610208578063450948ad1461029457806346a30a781461029c57806348d87239146102af578063556d89dd146102c4575f5ffd5b80630b592f45146102395780630f387b1e146102575780632fc00c761461027757806331d1464d1461027f575b5f5ffd5b61024161044f565b60405161024e91906116f4565b60405180910390f35b61026a610265366004611719565b610487565b60405161024e9190611779565b61024161052c565b61029261028d3660046117d8565b61055b565b005b610241610628565b6102926102aa366004611831565b610657565b6102b76106f2565b60405161024e9190611854565b6102926102d2366004611831565b610720565b6102416107ab565b6102b76107da565b6102926102f5366004611925565b610805565b610292610a0a565b610241610a2a565b6102b7610a59565b610292610a84565b6102b7610ac3565b6102b7610aee565b610241610b19565b610241610b4d565b610292610348366004611831565b610b7c565b610241610c07565b6102926103633660046117d8565b610c36565b6102b7610c81565b61024161037e3660046119e7565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b6102926103b2366004611a27565b610cac565b6103bf610e22565b60405161024e9190611aef565b6102b7610ef5565b6103dc610f20565b60405161024e9190611c51565b6102926103f7366004611831565b6111b8565b610241611243565b6104176104123660046117d8565b61126b565b60405161024e9190611ca8565b610292610432366004611831565b6113eb565b610241611470565b6102b761149f565b6102b76114ca565b5f61048261047e60017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611ccd565b5490565b905090565b5f8181548110610495575f80fd5b905f5260205f20015f9150905080546104ad90611cf4565b80601f01602080910402602001604051908101604052809291908181526020018280546104d990611cf4565b80156105245780601f106104fb57610100808354040283529160200191610524565b820191905f5260205f20905b81548152906001019060200180831161050757829003601f168201915b505050505081565b5f61048261047e60017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611ccd565b6105636114f5565b5f6001600160a01b03166001838360405161057f929190611d30565b908152604051908190036020019020546001600160a01b0316036105be5760405162461bcd60e51b81526004016105b590611d6d565b60405180910390fd5b600182826040516105d0929190611d30565b90815260405190819003602001812080546001600160a01b03191690557f5f9e4bc50ed4fc3c0bf14c4b518e1f8132c7a95ce5bed2e97c8675e6adf035739061061c9084908490611d9d565b60405180910390a15050565b5f61048261047e60017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611ccd565b61065f6114f5565b6001600160a01b0381166106855760405162461bcd60e51b81526004016105b590611de1565b6106b86106b360017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611ccd565b829055565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516106e79190611e23565b60405180910390a150565b61071d60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611ccd565b81565b6107286114f5565b6001600160a01b03811661074e5760405162461bcd60e51b81526004016105b590611de1565b61077c6106b360017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611ccd565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516106e79190611e74565b5f61048261047e60017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611ccd565b61071d60017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611ccd565b5f61080e611529565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f8115801561083a5750825b90505f8267ffffffffffffffff1660011480156108565750303b155b905081158015610864575080155b1561089b576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156108cf57845468ff00000000000000001916680100000000000000001785555b6108d886611551565b61090c61090660017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611ccd565b88519055565b61094361093a60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611ccd565b60208901519055565b61097a61097160017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611ccd565b60408901519055565b6109b16109a860017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611ccd565b60608901519055565b8315610a0157845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906109f890600190611ea7565b60405180910390a15b50505050505050565b610a126114f5565b60405162461bcd60e51b81526004016105b590611eb5565b5f61048261047e60017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611ccd565b61071d60017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611ccd565b3380610a8e611243565b6001600160a01b031614610ab7578060405163118cdaa760e01b81526004016105b591906116f4565b610ac08161156a565b50565b61071d60017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611ccd565b61071d60017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611ccd565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b5f61048261047e60017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611ccd565b610b846114f5565b6001600160a01b038116610baa5760405162461bcd60e51b81526004016105b590611de1565b610bd86106b360017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611ccd565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516106e79190611f48565b5f61048261047e60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611ccd565b610c3e6114f5565b8181604051610c4e929190611d30565b604051908190038120907fd8e1d6699327d650b4a450706e19aba65aef5d4aaa98f0f9b5c0c908b2f35bb9905f90a25050565b61071d60017fc9e8e7a4a583757cbcf624a50138f888cb585d449a8799952d3cc62760699622611ccd565b610cb46114f5565b6001600160a01b038116610cda5760405162461bcd60e51b81526004016105b590611de1565b5f6001600160a01b031660018484604051610cf6929190611d30565b908152604051908190036020019020546001600160a01b031614610d2c5760405162461bcd60e51b81526004016105b590611f8a565b5f6001600160a01b031660018484604051610d48929190611d30565b908152604051908190036020019020546001600160a01b031603610da1575f80546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301610d9f838583612032565b505b8060018484604051610db4929190611d30565b90815260405190819003602001812080546001600160a01b03939093166001600160a01b0319909316929092179091557f7ef997b0c9df3b39718be90c44d4d0d3d0230ac10eae31d63200210c7541ab7090610e15908590859085906120f1565b60405180910390a1505050565b60605f805480602002602001604051908101604052809291908181526020015f905b82821015610eec578382905f5260205f20018054610e6190611cf4565b80601f0160208091040260200160405190810160405280929190818152602001828054610e8d90611cf4565b8015610ed85780601f10610eaf57610100808354040283529160200191610ed8565b820191905f5260205f20905b815481529060010190602001808311610ebb57829003601f168201915b505050505081526020019060010190610e44565b50505050905090565b61071d60017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611ccd565b60408051610120810182525f8082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820152815490919067ffffffffffffffff811115610f8357610f83611862565b604051908082528060200260200182016040528015610fc857816020015b60408051808201909152606081525f6020820152815260200190600190039081610fa15790505b5090505f5b5f548110156110f25760405180604001604052805f8381548110610ff357610ff3612112565b905f5260205f2001805461100690611cf4565b80601f016020809104026020016040519081016040528092919081815260200182805461103290611cf4565b801561107d5780601f106110545761010080835404028352916020019161107d565b820191905f5260205f20905b81548152906001019060200180831161106057829003601f168201915b5050505050815260200160015f848154811061109b5761109b612112565b905f5260205f20016040516110b09190612195565b908152604051908190036020019020546001600160a01b0316905282518390839081106110df576110df612112565b6020908102919091010152600101610fcd565b50604051806101200160405280611107610c07565b6001600160a01b0316815260200161111d611470565b6001600160a01b0316815260200161113361052c565b6001600160a01b03168152602001611149610a2a565b6001600160a01b0316815260200161115f6107ab565b6001600160a01b03168152602001611175610b4d565b6001600160a01b0316815260200161118b610628565b6001600160a01b031681526020016111a161044f565b6001600160a01b0316815260200191909152919050565b6111c06114f5565b6001600160a01b0381166111e65760405162461bcd60e51b81526004016105b590611de1565b6112146106b360017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611ccd565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516106e791906121d1565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610b3d565b604080516060808201835280825260208201525f8183015290516002906112959085908590611d30565b90815260200160405180910390206040518060600160405290815f820180546112bd90611cf4565b80601f01602080910402602001604051908101604052809291908181526020018280546112e990611cf4565b80156113345780601f1061130b57610100808354040283529160200191611334565b820191905f5260205f20905b81548152906001019060200180831161131757829003601f168201915b5050505050815260200160018201805461134d90611cf4565b80601f016020809104026020016040519081016040528092919081815260200182805461137990611cf4565b80156113c45780601f1061139b576101008083540402835291602001916113c4565b820191905f5260205f20905b8154815290600101906020018083116113a757829003601f168201915b5050509183525050600291909101546001600160a01b031660209091015290505b92915050565b6113f36114f5565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b0319166001600160a01b0383169081178255611437610b19565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b5f61048261047e60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611ccd565b61071d60017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611ccd565b61071d60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611ccd565b336114fe610b19565b6001600160a01b031614611527573360405163118cdaa760e01b81526004016105b591906116f4565b565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006113e5565b6115596115a6565b611562816115e4565b610ac06115f5565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0080546001600160a01b03191681556115a2826115fd565b5050565b6115ae61166d565b611527576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6115ec6115a6565b610ac08161168b565b6115276115a6565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f611676611529565b5468010000000000000000900460ff16919050565b6116936115a6565b6001600160a01b038116610ab7575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016105b591906116f4565b5f6001600160a01b0382166113e5565b6116ee816116d5565b82525050565b602081016113e582846116e5565b805b8114610ac0575f5ffd5b80356113e581611702565b5f6020828403121561172c5761172c5f5ffd5b611736838361170e565b9392505050565b8281835e505f910152565b5f611751825190565b80845260208401935061176881856020860161173d565b601f01601f19169290920192915050565b602080825281016117368184611748565b5f5f83601f84011261179d5761179d5f5ffd5b50813567ffffffffffffffff8111156117b7576117b75f5ffd5b6020830191508360018202830111156117d1576117d15f5ffd5b9250929050565b5f5f602083850312156117ec576117ec5f5ffd5b823567ffffffffffffffff811115611805576118055f5ffd5b6118118582860161178a565b92509250509250929050565b611704816116d5565b80356113e58161181d565b5f60208284031215611844576118445f5ffd5b6117368383611826565b806116ee565b602081016113e5828461184e565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff8211171561189c5761189c611862565b6040525050565b5f6118ad60405190565b90506118b98282611876565b919050565b5f608082840312156118d1576118d15f5ffd5b6118db60806118a3565b90506118e78383611826565b81526118f68360208401611826565b60208201526119088360408401611826565b604082015261191a8360608401611826565b606082015292915050565b5f5f60a08385031215611939576119395f5ffd5b61194384846118be565b91506119528460808501611826565b90509250929050565b5f67ffffffffffffffff82111561197457611974611862565b601f19601f83011660200192915050565b82818337505f910152565b5f6119a261199d8461195b565b6118a3565b90508281528383830111156119b8576119b85f5ffd5b611736836020830184611985565b5f82601f8301126119d8576119d85f5ffd5b61173683833560208501611990565b5f602082840312156119fa576119fa5f5ffd5b813567ffffffffffffffff811115611a1357611a135f5ffd5b611a1f848285016119c6565b949350505050565b5f5f5f60408486031215611a3c57611a3c5f5ffd5b833567ffffffffffffffff811115611a5557611a555f5ffd5b611a618682870161178a565b9350935050611a738560208601611826565b90509250925092565b5f6117368383611748565b60200190565b5f611a96825190565b80845260208401935083602082028501611ab08560200190565b5f5b84811015611ae35783830388528151611acb8482611a7c565b93505060208201602098909801979150600101611ab2565b50909695505050505050565b602080825281016117368184611a8d565b805160408084525f9190840190611b178282611748565b9150506020830151611b2c60208601826116e5565b509392505050565b5f6117368383611b00565b5f611b48825190565b80845260208401935083602082028501611b628560200190565b5f5b84811015611ae35783830388528151611b7d8482611b34565b93505060208201602098909801979150600101611b64565b80515f90610120840190611ba985826116e5565b506020830151611bbc60208601826116e5565b506040830151611bcf60408601826116e5565b506060830151611be260608601826116e5565b506080830151611bf560808601826116e5565b5060a0830151611c0860a08601826116e5565b5060c0830151611c1b60c08601826116e5565b5060e0830151611c2e60e08601826116e5565b50610100830151848203610100860152611c488282611b3f565b95945050505050565b602080825281016117368184611b95565b805160608084525f9190840190611c798282611748565b91505060208301518482036020860152611c938282611748565b9150506040830151611b2c60408601826116e5565b602080825281016117368184611c62565b634e487b7160e01b5f52601160045260245ffd5b818103818111156113e5576113e5611cb9565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680611d0857607f821691505b602082108103611d1a57611d1a611ce0565b50919050565b611d2b828483611985565b500190565b611736818385611d20565b60168152602081017f4164647265737320646f6573206e6f742065786973740000000000000000000081529050611a87565b602080825281016113e581611d3b565b818352602083019250611d91828483611985565b50601f01601f19160190565b60208082528101611a1f818486611d7d565b600f8152602081017f496e76616c69642061646472657373000000000000000000000000000000000081529050611a87565b602080825281016113e581611daf565b60158152602081017f6c3143726f7373436861696e4d657373656e676572000000000000000000000081529050611a87565b60408082528101611e3381611df1565b90506113e560208301846116e5565b60088152602081017f6c3242726964676500000000000000000000000000000000000000000000000081529050611a87565b60408082528101611e3381611e42565b5f6113e582611e91565b90565b67ffffffffffffffff1690565b6116ee81611e84565b602081016113e58284611e9e565b602080825281016113e581603481527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60208201527f742072656e6f756e6365206f776e657273686970000000000000000000000000604082015260600190565b60158152602081017f6c3243726f7373436861696e4d657373656e676572000000000000000000000081529050611a87565b60408082528101611e3381611f16565b60168152602081017f4164647265737320616c7265616479206578697374730000000000000000000081529050611a87565b602080825281016113e581611f58565b5f6113e5611e8e8381565b611fae83611f9a565b81545f1960089490940293841b1916921b91909117905550565b5f611fd4818484611fa5565b505050565b818110156115a257611feb5f82611fc8565b600101611fd9565b601f821115611fd4575f818152602090206020601f850104810160208510156120195750805b61202b6020601f860104830182611fd9565b5050505050565b8267ffffffffffffffff81111561204b5761204b611862565b6120558254611cf4565b612060828285611ff3565b505f601f821160018114612092575f831561207b5750848201355b5f19600885021c19811660028502178555506120e9565b5f84815260208120601f198516915b828110156120c157878501358255602094850194600190920191016120a1565b50848210156120dd575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b60408082528101612103818587611d7d565b9050611a1f60208301846116e5565b634e487b7160e01b5f52603260045260245ffd5b5f815461213281611cf4565b600182168015612149576001811461215e5761218c565b60ff198316865281151582028601935061218c565b5f858152602090205f5b8381101561218457815488820152600190910190602001612168565b505081860193505b50505092915050565b6113e58183612126565b60088152602081017f6c3142726964676500000000000000000000000000000000000000000000000081529050611a87565b60408082528101611e338161219f56fea2646970667358221220e59e622de3e363f437175d1dfdd97781d49aaf330d2910f23512750d2435122864736f6c634300081c0033",
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

// GetAdditionaContractNames is a free data retrieval call binding the contract method 0xbc162c90.
//
// Solidity: function getAdditionaContractNames() view returns(string[])
func (_NetworkConfig *NetworkConfigCaller) GetAdditionaContractNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "getAdditionaContractNames")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetAdditionaContractNames is a free data retrieval call binding the contract method 0xbc162c90.
//
// Solidity: function getAdditionaContractNames() view returns(string[])
func (_NetworkConfig *NetworkConfigSession) GetAdditionaContractNames() ([]string, error) {
	return _NetworkConfig.Contract.GetAdditionaContractNames(&_NetworkConfig.CallOpts)
}

// GetAdditionaContractNames is a free data retrieval call binding the contract method 0xbc162c90.
//
// Solidity: function getAdditionaContractNames() view returns(string[])
func (_NetworkConfig *NetworkConfigCallerSession) GetAdditionaContractNames() ([]string, error) {
	return _NetworkConfig.Contract.GetAdditionaContractNames(&_NetworkConfig.CallOpts)
}

// GetContractVersion is a free data retrieval call binding the contract method 0xed11f78d.
//
// Solidity: function getContractVersion(string name) view returns((string,string,address))
func (_NetworkConfig *NetworkConfigCaller) GetContractVersion(opts *bind.CallOpts, name string) (NetworkConfigContractVersion, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "getContractVersion", name)

	if err != nil {
		return *new(NetworkConfigContractVersion), err
	}

	out0 := *abi.ConvertType(out[0], new(NetworkConfigContractVersion)).(*NetworkConfigContractVersion)

	return out0, err

}

// GetContractVersion is a free data retrieval call binding the contract method 0xed11f78d.
//
// Solidity: function getContractVersion(string name) view returns((string,string,address))
func (_NetworkConfig *NetworkConfigSession) GetContractVersion(name string) (NetworkConfigContractVersion, error) {
	return _NetworkConfig.Contract.GetContractVersion(&_NetworkConfig.CallOpts, name)
}

// GetContractVersion is a free data retrieval call binding the contract method 0xed11f78d.
//
// Solidity: function getContractVersion(string name) view returns((string,string,address))
func (_NetworkConfig *NetworkConfigCallerSession) GetContractVersion(name string) (NetworkConfigContractVersion, error) {
	return _NetworkConfig.Contract.GetContractVersion(&_NetworkConfig.CallOpts, name)
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

// RecordHardfork is a paid mutator transaction binding the contract method 0xadc5c207.
//
// Solidity: function recordHardfork(string hardforkName) returns()
func (_NetworkConfig *NetworkConfigTransactor) RecordHardfork(opts *bind.TransactOpts, hardforkName string) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "recordHardfork", hardforkName)
}

// RecordHardfork is a paid mutator transaction binding the contract method 0xadc5c207.
//
// Solidity: function recordHardfork(string hardforkName) returns()
func (_NetworkConfig *NetworkConfigSession) RecordHardfork(hardforkName string) (*types.Transaction, error) {
	return _NetworkConfig.Contract.RecordHardfork(&_NetworkConfig.TransactOpts, hardforkName)
}

// RecordHardfork is a paid mutator transaction binding the contract method 0xadc5c207.
//
// Solidity: function recordHardfork(string hardforkName) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) RecordHardfork(hardforkName string) (*types.Transaction, error) {
	return _NetworkConfig.Contract.RecordHardfork(&_NetworkConfig.TransactOpts, hardforkName)
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

// NetworkConfigHardforkUpgradeIterator is returned from FilterHardforkUpgrade and is used to iterate over the raw logs and unpacked data for HardforkUpgrade events raised by the NetworkConfig contract.
type NetworkConfigHardforkUpgradeIterator struct {
	Event *NetworkConfigHardforkUpgrade // Event containing the contract specifics and raw log

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
func (it *NetworkConfigHardforkUpgradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigHardforkUpgrade)
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
		it.Event = new(NetworkConfigHardforkUpgrade)
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
func (it *NetworkConfigHardforkUpgradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigHardforkUpgradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigHardforkUpgrade represents a HardforkUpgrade event raised by the NetworkConfig contract.
type NetworkConfigHardforkUpgrade struct {
	ForkName common.Hash
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterHardforkUpgrade is a free log retrieval operation binding the contract event 0xd8e1d6699327d650b4a450706e19aba65aef5d4aaa98f0f9b5c0c908b2f35bb9.
//
// Solidity: event HardforkUpgrade(string indexed forkName)
func (_NetworkConfig *NetworkConfigFilterer) FilterHardforkUpgrade(opts *bind.FilterOpts, forkName []string) (*NetworkConfigHardforkUpgradeIterator, error) {

	var forkNameRule []interface{}
	for _, forkNameItem := range forkName {
		forkNameRule = append(forkNameRule, forkNameItem)
	}

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "HardforkUpgrade", forkNameRule)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigHardforkUpgradeIterator{contract: _NetworkConfig.contract, event: "HardforkUpgrade", logs: logs, sub: sub}, nil
}

// WatchHardforkUpgrade is a free log subscription operation binding the contract event 0xd8e1d6699327d650b4a450706e19aba65aef5d4aaa98f0f9b5c0c908b2f35bb9.
//
// Solidity: event HardforkUpgrade(string indexed forkName)
func (_NetworkConfig *NetworkConfigFilterer) WatchHardforkUpgrade(opts *bind.WatchOpts, sink chan<- *NetworkConfigHardforkUpgrade, forkName []string) (event.Subscription, error) {

	var forkNameRule []interface{}
	for _, forkNameItem := range forkName {
		forkNameRule = append(forkNameRule, forkNameItem)
	}

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "HardforkUpgrade", forkNameRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigHardforkUpgrade)
				if err := _NetworkConfig.contract.UnpackLog(event, "HardforkUpgrade", log); err != nil {
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

// ParseHardforkUpgrade is a log parse operation binding the contract event 0xd8e1d6699327d650b4a450706e19aba65aef5d4aaa98f0f9b5c0c908b2f35bb9.
//
// Solidity: event HardforkUpgrade(string indexed forkName)
func (_NetworkConfig *NetworkConfigFilterer) ParseHardforkUpgrade(log types.Log) (*NetworkConfigHardforkUpgrade, error) {
	event := new(NetworkConfigHardforkUpgrade)
	if err := _NetworkConfig.contract.UnpackLog(event, "HardforkUpgrade", log); err != nil {
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
