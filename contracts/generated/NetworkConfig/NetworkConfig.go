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
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"AdditionalContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"forkName\",\"type\":\"string\"}],\"name\":\"HardforkUpgrade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NetworkContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DATA_AVAILABILITY_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FORK_MANAGER_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MESSAGE_BUS_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NETWORK_ENCLAVE_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addAdditionalAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"additionalAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"addressNames\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addresses\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dataAvailabilityRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1CrossChainMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2CrossChainMessenger\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.NamedAddress[]\",\"name\":\"additionalContracts\",\"type\":\"tuple[]\"}],\"internalType\":\"structNetworkConfig.Addresses\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"daRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdditionaContractNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getContractVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.ContractVersion\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dataAvailabilityRegistry\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.FixedAddresses\",\"name\":\"_addresses\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkEnclaveRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hardforkName\",\"type\":\"string\"}],\"name\":\"recordHardfork\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50611ea18061001c5f395ff3fe608060405234801561000f575f5ffd5b50600436106101e7575f3560e01c8063934746a711610109578063be9f82071161009e578063f2fde38b1161006e578063f2fde38b146103db578063f5e9f286146103ee578063faa5e2de146103f6578063fbfd6d91146103fe575f5ffd5b8063be9f82071461038b578063da0321cd14610393578063e1825d06146103a8578063ed11f78d146103bb575f5ffd5b8063ae61ecba116100d9578063ae61ecba14610327578063af4546351461032f578063b7bef9ab14610363578063bc162c9014610376575f5ffd5b8063934746a7146102f157806396493cc5146102f9578063a1b918d61461030c578063adc5c20714610314575f5ffd5b806367cc852e1161017f57806372bad9121161014f57806372bad912146102a9578063812b1ffe146102b157806385f427cb146102b95780638da5cb5b146102c1575f5ffd5b806367cc852e1461027e5780636c1358ac14610286578063715018a61461029957806371fd11f3146102a1575f5ffd5b806346a30a78116101ba57806346a30a781461023957806348d872391461024e578063556d89dd146102635780635ab2a55814610276575f5ffd5b80630b592f45146101eb5780630f387b1e146102095780632fc00c7614610229578063450948ad14610231575b5f5ffd5b6101f3610406565b6040516102009190611471565b60405180910390f35b61021c610217366004611496565b61043e565b60405161020091906114f6565b6101f36104e3565b6101f3610512565b61024c61024736600461151b565b610541565b005b6102566105e5565b604051610200919061153e565b61024c61027136600461151b565b610613565b6101f361069e565b6102566106cd565b61024c61029436600461160f565b6106f8565b61024c610912565b6101f3610925565b610256610954565b61025661097f565b6102566109aa565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166101f3565b6101f36109d5565b61024c61030736600461151b565b610a04565b6101f3610a8f565b61024c610322366004611693565b610abe565b610256610b09565b6101f361033d366004611764565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b61024c6103713660046117a4565b610b34565b61037e610c65565b604051610200919061186c565b610256610d38565b61039b610d63565b60405161020091906119ce565b61024c6103b636600461151b565b610ffb565b6103ce6103c9366004611693565b611086565b6040516102009190611a25565b61024c6103e936600461151b565b611206565b6101f361125c565b61025661128b565b6102566112b6565b5f61043961043560017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611a4a565b5490565b905090565b5f818154811061044c575f80fd5b905f5260205f20015f91509050805461046490611a71565b80601f016020809104026020016040519081016040528092919081815260200182805461049090611a71565b80156104db5780601f106104b2576101008083540402835291602001916104db565b820191905f5260205f20905b8154815290600101906020018083116104be57829003601f168201915b505050505081565b5f61043961043560017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611a4a565b5f61043961043560017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611a4a565b6105496112e1565b6001600160a01b0381166105785760405162461bcd60e51b815260040161056f90611acf565b60405180910390fd5b6105ab6105a660017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611a4a565b829055565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516105da9190611b11565b60405180910390a150565b61061060017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611a4a565b81565b61061b6112e1565b6001600160a01b0381166106415760405162461bcd60e51b815260040161056f90611acf565b61066f6105a660017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611a4a565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516105da9190611b62565b5f61043961043560017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611a4a565b61061060017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f82611a4a565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156107425750825b90505f8267ffffffffffffffff16600114801561075e5750303b155b90508115801561076c575080155b156107a3576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156107d757845468ff00000000000000001916680100000000000000001785555b6107e086611355565b61081461080e60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611a4a565b88519055565b61084b61084260017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611a4a565b60208901519055565b61088261087960017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611a4a565b60408901519055565b6108b96108b060017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611a4a565b60608901519055565b831561090957845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061090090600190611b95565b60405180910390a15b50505050505050565b61091a6112e1565b6109235f611366565b565b5f61043961043560017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611a4a565b61061060017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba611a4a565b61061060017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611a4a565b61061060017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611a4a565b5f61043961043560017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c91611a4a565b610a0c6112e1565b6001600160a01b038116610a325760405162461bcd60e51b815260040161056f90611acf565b610a606105a660017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf611a4a565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516105da9190611bd5565b5f61043961043560017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611a4a565b610ac66112e1565b8181604051610ad6929190611bf5565b604051908190038120907fd8e1d6699327d650b4a450706e19aba65aef5d4aaa98f0f9b5c0c908b2f35bb9905f90a25050565b61061060017fc9e8e7a4a583757cbcf624a50138f888cb585d449a8799952d3cc62760699622611a4a565b610b3c6112e1565b6001600160a01b038116610b625760405162461bcd60e51b815260040161056f90611acf565b5f6001600160a01b031660018484604051610b7e929190611bf5565b908152604051908190036020019020546001600160a01b031603610bd7575f80546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301610bd5838583611c9c565b505b8060018484604051610bea929190611bf5565b90815260405190819003602001812080546001600160a01b039390931673ffffffffffffffffffffffffffffffffffffffff19909316929092179091557f7ef997b0c9df3b39718be90c44d4d0d3d0230ac10eae31d63200210c7541ab7090610c5890859085908590611d7b565b60405180910390a1505050565b60605f805480602002602001604051908101604052809291908181526020015f905b82821015610d2f578382905f5260205f20018054610ca490611a71565b80601f0160208091040260200160405190810160405280929190818152602001828054610cd090611a71565b8015610d1b5780601f10610cf257610100808354040283529160200191610d1b565b820191905f5260205f20905b815481529060010190602001808311610cfe57829003601f168201915b505050505081526020019060010190610c87565b50505050905090565b61061060017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611a4a565b60408051610120810182525f8082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820152815490919067ffffffffffffffff811115610dc657610dc661154c565b604051908082528060200260200182016040528015610e0b57816020015b60408051808201909152606081525f6020820152815260200190600190039081610de45790505b5090505f5b5f54811015610f355760405180604001604052805f8381548110610e3657610e36611d9c565b905f5260205f20018054610e4990611a71565b80601f0160208091040260200160405190810160405280929190818152602001828054610e7590611a71565b8015610ec05780601f10610e9757610100808354040283529160200191610ec0565b820191905f5260205f20905b815481529060010190602001808311610ea357829003601f168201915b5050505050815260200160015f8481548110610ede57610ede611d9c565b905f5260205f2001604051610ef39190611e1f565b908152604051908190036020019020546001600160a01b031690528251839083908110610f2257610f22611d9c565b6020908102919091010152600101610e10565b50604051806101200160405280610f4a610a8f565b6001600160a01b03168152602001610f6061125c565b6001600160a01b03168152602001610f766104e3565b6001600160a01b03168152602001610f8c610925565b6001600160a01b03168152602001610fa261069e565b6001600160a01b03168152602001610fb86109d5565b6001600160a01b03168152602001610fce610512565b6001600160a01b03168152602001610fe4610406565b6001600160a01b0316815260200191909152919050565b6110036112e1565b6001600160a01b0381166110295760405162461bcd60e51b815260040161056f90611acf565b6110576105a660017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a611a4a565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73816040516105da9190611e5b565b604080516060808201835280825260208201525f8183015290516002906110b09085908590611bf5565b90815260200160405180910390206040518060600160405290815f820180546110d890611a71565b80601f016020809104026020016040519081016040528092919081815260200182805461110490611a71565b801561114f5780601f106111265761010080835404028352916020019161114f565b820191905f5260205f20905b81548152906001019060200180831161113257829003601f168201915b5050505050815260200160018201805461116890611a71565b80601f016020809104026020016040519081016040528092919081815260200182805461119490611a71565b80156111df5780601f106111b6576101008083540402835291602001916111df565b820191905f5260205f20905b8154815290600101906020018083116111c257829003601f168201915b5050509183525050600291909101546001600160a01b031660209091015290505b92915050565b61120e6112e1565b6001600160a01b038116611250575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161056f9190611471565b61125981611366565b50565b5f61043961043560017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c611a4a565b61061060017f5c104a5cbc447428f263418ae884b79d6ce229d3ad858ef1544ef89e88adc15f611a4a565b61061060017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958611a4a565b336113137f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461092357336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161056f9190611471565b61135d6113e3565b6112598161144a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610923576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61120e6113e3565b5f6001600160a01b038216611200565b61146b81611452565b82525050565b602081016112008284611462565b805b8114611259575f5ffd5b80356112008161147f565b5f602082840312156114a9576114a95f5ffd5b6114b3838361148b565b9392505050565b8281835e505f910152565b5f6114ce825190565b8084526020840193506114e58185602086016114ba565b601f01601f19169290920192915050565b602080825281016114b381846114c5565b61148181611452565b803561120081611507565b5f6020828403121561152e5761152e5f5ffd5b6114b38383611510565b8061146b565b602081016112008284611538565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff821117156115865761158661154c565b6040525050565b5f61159760405190565b90506115a38282611560565b919050565b5f608082840312156115bb576115bb5f5ffd5b6115c5608061158d565b90506115d18383611510565b81526115e08360208401611510565b60208201526115f28360408401611510565b60408201526116048360608401611510565b606082015292915050565b5f5f60a08385031215611623576116235f5ffd5b61162d84846115a8565b915061163c8460808501611510565b90509250929050565b5f5f83601f840112611658576116585f5ffd5b50813567ffffffffffffffff811115611672576116725f5ffd5b60208301915083600182028301111561168c5761168c5f5ffd5b9250929050565b5f5f602083850312156116a7576116a75f5ffd5b823567ffffffffffffffff8111156116c0576116c05f5ffd5b6116cc85828601611645565b92509250509250929050565b5f67ffffffffffffffff8211156116f1576116f161154c565b601f19601f83011660200192915050565b82818337505f910152565b5f61171f61171a846116d8565b61158d565b9050828152838383011115611735576117355f5ffd5b6114b3836020830184611702565b5f82601f830112611755576117555f5ffd5b6114b38383356020850161170d565b5f60208284031215611777576117775f5ffd5b813567ffffffffffffffff811115611790576117905f5ffd5b61179c84828501611743565b949350505050565b5f5f5f604084860312156117b9576117b95f5ffd5b833567ffffffffffffffff8111156117d2576117d25f5ffd5b6117de86828701611645565b93509350506117f08560208601611510565b90509250925092565b5f6114b383836114c5565b60200190565b5f611813825190565b8084526020840193508360208202850161182d8560200190565b5f5b84811015611860578383038852815161184884826117f9565b9350506020820160209890980197915060010161182f565b50909695505050505050565b602080825281016114b3818461180a565b805160408084525f919084019061189482826114c5565b91505060208301516118a96020860182611462565b509392505050565b5f6114b3838361187d565b5f6118c5825190565b808452602084019350836020820285016118df8560200190565b5f5b8481101561186057838303885281516118fa84826118b1565b935050602082016020989098019791506001016118e1565b80515f906101208401906119268582611462565b5060208301516119396020860182611462565b50604083015161194c6040860182611462565b50606083015161195f6060860182611462565b5060808301516119726080860182611462565b5060a083015161198560a0860182611462565b5060c083015161199860c0860182611462565b5060e08301516119ab60e0860182611462565b506101008301518482036101008601526119c582826118bc565b95945050505050565b602080825281016114b38184611912565b805160608084525f91908401906119f682826114c5565b91505060208301518482036020860152611a1082826114c5565b91505060408301516118a96040860182611462565b602080825281016114b381846119df565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561120057611200611a36565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680611a8557607f821691505b602082108103611a9757611a97611a5d565b50919050565b600f8152602081017f496e76616c69642061646472657373000000000000000000000000000000000081529050611804565b6020808252810161120081611a9d565b60158152602081017f6c3143726f7373436861696e4d657373656e676572000000000000000000000081529050611804565b60408082528101611b2181611adf565b90506112006020830184611462565b60088152602081017f6c3242726964676500000000000000000000000000000000000000000000000081529050611804565b60408082528101611b2181611b30565b5f61120082611b7f565b90565b67ffffffffffffffff1690565b61146b81611b72565b602081016112008284611b8c565b60158152602081017f6c3243726f7373436861696e4d657373656e676572000000000000000000000081529050611804565b60408082528101611b2181611ba3565b611bf0828483611702565b500190565b6114b3818385611be5565b5f611200611b7c8381565b611c1483611c00565b81545f1960089490940293841b1916921b91909117905550565b5f611c3a818484611c0b565b505050565b81811015611c5957611c515f82611c2e565b600101611c3f565b5050565b601f821115611c3a575f818152602090206020601f85010481016020851015611c835750805b611c956020601f860104830182611c3f565b5050505050565b8267ffffffffffffffff811115611cb557611cb561154c565b611cbf8254611a71565b611cca828285611c5d565b505f601f821160018114611cfc575f8315611ce55750848201355b5f19600885021c1981166002850217855550611d53565b5f84815260208120601f198516915b82811015611d2b5787850135825560209485019460019092019101611d0b565b5084821015611d47575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b818352602083019250611d6f828483611702565b50601f01601f19160190565b60408082528101611d8d818587611d5b565b905061179c6020830184611462565b634e487b7160e01b5f52603260045260245ffd5b5f8154611dbc81611a71565b600182168015611dd35760018114611de857611e16565b60ff1983168652811515820286019350611e16565b5f858152602090205f5b83811015611e0e57815488820152600190910190602001611df2565b505081860193505b50505092915050565b6112008183611db0565b60088152602081017f6c3142726964676500000000000000000000000000000000000000000000000081529050611804565b60408082528101611b2181611e2956fea2646970667358221220e94b03665ca98ab96272260eba53fc212ed40b0de3cd2c439a6be109970ae71e64736f6c634300081c0033",
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
// Solidity: function additionalAddresses(string ) view returns(address)
func (_NetworkConfig *NetworkConfigCaller) AdditionalAddresses(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "additionalAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string ) view returns(address)
func (_NetworkConfig *NetworkConfigSession) AdditionalAddresses(arg0 string) (common.Address, error) {
	return _NetworkConfig.Contract.AdditionalAddresses(&_NetworkConfig.CallOpts, arg0)
}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string ) view returns(address)
func (_NetworkConfig *NetworkConfigCallerSession) AdditionalAddresses(arg0 string) (common.Address, error) {
	return _NetworkConfig.Contract.AdditionalAddresses(&_NetworkConfig.CallOpts, arg0)
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
// Solidity: function initialize((address,address,address,address) _addresses, address owner) returns()
func (_NetworkConfig *NetworkConfigTransactor) Initialize(opts *bind.TransactOpts, _addresses NetworkConfigFixedAddresses, owner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "initialize", _addresses, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x6c1358ac.
//
// Solidity: function initialize((address,address,address,address) _addresses, address owner) returns()
func (_NetworkConfig *NetworkConfigSession) Initialize(_addresses NetworkConfigFixedAddresses, owner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x6c1358ac.
//
// Solidity: function initialize((address,address,address,address) _addresses, address owner) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) Initialize(_addresses NetworkConfigFixedAddresses, owner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses, owner)
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkConfig *NetworkConfigTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkConfig *NetworkConfigSession) RenounceOwnership() (*types.Transaction, error) {
	return _NetworkConfig.Contract.RenounceOwnership(&_NetworkConfig.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkConfig *NetworkConfigTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NetworkConfig.Contract.RenounceOwnership(&_NetworkConfig.TransactOpts)
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
