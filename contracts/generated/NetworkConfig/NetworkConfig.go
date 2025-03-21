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
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"AdditionalContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NetworkContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L1_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_BRIDGE_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_CROSS_CHAIN_MESSENGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MESSAGE_BUS_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NETWORK_ENCLAVE_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROLLUP_CONTRACT_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addAdditionalAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"additionalAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"addressNames\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addresses\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2Bridge\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l1CrossChainMessenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"l2CrossChainMessenger\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.NamedAddress[]\",\"name\":\"additionalContracts\",\"type\":\"tuple[]\"}],\"internalType\":\"structNetworkConfig.Addresses\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAdditionaContractNames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rollupContract\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.FixedAddresses\",\"name\":\"_addresses\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2BridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2CrossChainMessengerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkEnclaveRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL1CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2BridgeAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"setL2CrossChainMessengerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50611bb58061001c5f395ff3fe608060405234801561000f575f5ffd5b50600436106101c6575f3560e01c806385a7805a116100fe578063b7bef9ab1161009e578063e1825d061161006e578063e1825d0614610374578063f2fde38b14610387578063f5e9f2861461039a578063fbfd6d91146103a2575f5ffd5b8063b7bef9ab1461032f578063bc162c9014610342578063be9f820714610357578063da0321cd1461035f575f5ffd5b8063934746a7116100d9578063934746a7146102d857806396493cc5146102e0578063a1b918d6146102f3578063af454635146102fb575f5ffd5b806385a7805a1461029857806385f427cb146102a05780638da5cb5b146102a8575f5ffd5b8063556d89dd116101695780636c1358ac116101445780636c1358ac1461026d578063715018a61461028057806372bad91214610288578063812b1ffe14610290575f5ffd5b8063556d89dd1461024a5780635ab2a5581461025d57806367cc852e14610265575f5ffd5b8063450948ad116101a4578063450948ad1461021057806346a30a781461021857806348d872391461022d5780634e02727114610242575f5ffd5b80630b592f45146101ca5780630f387b1e146101e85780632fc00c7614610208575b5f5ffd5b6101d26103aa565b6040516101df9190611221565b60405180910390f35b6101fb6101f6366004611246565b6103e2565b6040516101df91906112a6565b6101d2610487565b6101d26104b6565b61022b6102263660046112cb565b6104e5565b005b610235610589565b6040516101df91906112ee565b6101d26105b7565b61022b6102583660046112cb565b6105e6565b6101d2610671565b6102356106a0565b61022b61027b3660046113bf565b6106cb565b61022b6108e5565b6102356108f8565b610235610923565b61023561094e565b610235610979565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166101d2565b6101d26109a4565b61022b6102ee3660046112cb565b6109d3565b6101d2610a5e565b6101d2610309366004611481565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b61022b61033d36600461150f565b610a8d565b61034a610bbe565b6040516101df91906115d7565b610235610c91565b610367610cbc565b6040516101df9190611739565b61022b6103823660046112cb565b610f54565b61022b6103953660046112cb565b610fdf565b6101d2611035565b610235611064565b5f6103dd6103d960017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf61175e565b5490565b905090565b5f81815481106103f0575f80fd5b905f5260205f20015f91509050805461040890611785565b80601f016020809104026020016040519081016040528092919081815260200182805461043490611785565b801561047f5780601f106104565761010080835404028352916020019161047f565b820191905f5260205f20905b81548152906001019060200180831161046257829003601f168201915b505050505081565b5f6103dd6103d960017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba61175e565b5f6103dd6103d960017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f8261175e565b6104ed61108f565b6001600160a01b03811661051c5760405162461bcd60e51b8152600401610513906117e3565b60405180910390fd5b61054f61054a60017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f8261175e565b829055565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b738160405161057e9190611825565b60405180910390a150565b6105b460017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c61175e565b81565b5f6103dd6103d960017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc61175e565b6105ee61108f565b6001600160a01b0381166106145760405162461bcd60e51b8152600401610513906117e3565b61064261054a60017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c9161175e565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b738160405161057e9190611876565b5f6103dd6103d960017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a61175e565b6105b460017f9c4bf36639b03148aa45703f540290d6a1c2225945d7196e6b3ef866efdf4f8261175e565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156107155750825b90505f8267ffffffffffffffff1660011480156107315750303b155b90508115801561073f575080155b15610776576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156107aa57845468ff00000000000000001916680100000000000000001785555b6107b386611103565b6107e76107e160017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a78307195861175e565b88519055565b61081e61081560017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c61175e565b60208901519055565b61085561084c60017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba61175e565b60408901519055565b61088c61088360017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc61175e565b60608901519055565b83156108dc57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906108d3906001906118a9565b60405180910390a15b50505050505050565b6108ed61108f565b6108f65f611114565b565b6105b460017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba61175e565b6105b460017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c9161175e565b6105b460017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc61175e565b6105b460017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf61175e565b5f6103dd6103d960017f01487b9e499fe6b85dcf5493ba4e0a725bd52278ead06c0d370c5b5e3d513c9161175e565b6109db61108f565b6001600160a01b038116610a015760405162461bcd60e51b8152600401610513906117e3565b610a2f61054a60017f73731984b9847fe0f6bd6840905e6dc77e4bfa84e759f8238ab14c4e91ca3cbf61175e565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b738160405161057e91906118e9565b5f6103dd6103d960017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a78307195861175e565b610a9561108f565b6001600160a01b038116610abb5760405162461bcd60e51b8152600401610513906117e3565b5f6001600160a01b031660018484604051610ad7929190611909565b908152604051908190036020019020546001600160a01b031603610b30575f80546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301610b2e8385836119b0565b505b8060018484604051610b43929190611909565b90815260405190819003602001812080546001600160a01b039390931673ffffffffffffffffffffffffffffffffffffffff19909316929092179091557f7ef997b0c9df3b39718be90c44d4d0d3d0230ac10eae31d63200210c7541ab7090610bb190859085908590611a8f565b60405180910390a1505050565b60605f805480602002602001604051908101604052809291908181526020015f905b82821015610c88578382905f5260205f20018054610bfd90611785565b80601f0160208091040260200160405190810160405280929190818152602001828054610c2990611785565b8015610c745780601f10610c4b57610100808354040283529160200191610c74565b820191905f5260205f20905b815481529060010190602001808311610c5757829003601f168201915b505050505081526020019060010190610be0565b50505050905090565b6105b460017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a61175e565b60408051610120810182525f8082526020820181905291810182905260608082018390526080820183905260a0820183905260c0820183905260e08201839052610100820152815490919067ffffffffffffffff811115610d1f57610d1f6112fc565b604051908082528060200260200182016040528015610d6457816020015b60408051808201909152606081525f6020820152815260200190600190039081610d3d5790505b5090505f5b5f54811015610e8e5760405180604001604052805f8381548110610d8f57610d8f611ab0565b905f5260205f20018054610da290611785565b80601f0160208091040260200160405190810160405280929190818152602001828054610dce90611785565b8015610e195780601f10610df057610100808354040283529160200191610e19565b820191905f5260205f20905b815481529060010190602001808311610dfc57829003601f168201915b5050505050815260200160015f8481548110610e3757610e37611ab0565b905f5260205f2001604051610e4c9190611b33565b908152604051908190036020019020546001600160a01b031690528251839083908110610e7b57610e7b611ab0565b6020908102919091010152600101610d69565b50604051806101200160405280610ea3610a5e565b6001600160a01b03168152602001610eb9611035565b6001600160a01b03168152602001610ecf610487565b6001600160a01b03168152602001610ee56105b7565b6001600160a01b03168152602001610efb610671565b6001600160a01b03168152602001610f116109a4565b6001600160a01b03168152602001610f276104b6565b6001600160a01b03168152602001610f3d6103aa565b6001600160a01b0316815260200191909152919050565b610f5c61108f565b6001600160a01b038116610f825760405162461bcd60e51b8152600401610513906117e3565b610fb061054a60017f5c3a696a2e63ec310c7ce2ec9686153b437260263b2fd38923a13e3adc7a8d8a61175e565b7f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b738160405161057e9190611b6f565b610fe761108f565b6001600160a01b038116611029575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016105139190611221565b61103281611114565b50565b5f6103dd6103d960017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c61175e565b6105b460017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a78307195861175e565b336110c17f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146108f657336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016105139190611221565b61110b611191565b611032816111f8565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166108f6576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610fe7611191565b5f6001600160a01b0382165b92915050565b61121b81611200565b82525050565b6020810161120c8284611212565b805b8114611032575f5ffd5b803561120c8161122f565b5f60208284031215611259576112595f5ffd5b611263838361123b565b9392505050565b8281835e505f910152565b5f61127e825190565b80845260208401935061129581856020860161126a565b601f01601f19169290920192915050565b602080825281016112638184611275565b61123181611200565b803561120c816112b7565b5f602082840312156112de576112de5f5ffd5b61126383836112c0565b8061121b565b6020810161120c82846112e8565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff82111715611336576113366112fc565b6040525050565b5f61134760405190565b90506113538282611310565b919050565b5f6080828403121561136b5761136b5f5ffd5b611375608061133d565b905061138183836112c0565b815261139083602084016112c0565b60208201526113a283604084016112c0565b60408201526113b483606084016112c0565b606082015292915050565b5f5f60a083850312156113d3576113d35f5ffd5b6113dd8484611358565b91506113ec84608085016112c0565b90509250929050565b5f67ffffffffffffffff82111561140e5761140e6112fc565b601f19601f83011660200192915050565b82818337505f910152565b5f61143c611437846113f5565b61133d565b9050828152838383011115611452576114525f5ffd5b61126383602083018461141f565b5f82601f830112611472576114725f5ffd5b6112638383356020850161142a565b5f60208284031215611494576114945f5ffd5b813567ffffffffffffffff8111156114ad576114ad5f5ffd5b6114b984828501611460565b949350505050565b5f5f83601f8401126114d4576114d45f5ffd5b50813567ffffffffffffffff8111156114ee576114ee5f5ffd5b602083019150836001820283011115611508576115085f5ffd5b9250929050565b5f5f5f60408486031215611524576115245f5ffd5b833567ffffffffffffffff81111561153d5761153d5f5ffd5b611549868287016114c1565b935093505061155b85602086016112c0565b90509250925092565b5f6112638383611275565b60200190565b5f61157e825190565b808452602084019350836020820285016115988560200190565b5f5b848110156115cb57838303885281516115b38482611564565b9350506020820160209890980197915060010161159a565b50909695505050505050565b602080825281016112638184611575565b805160408084525f91908401906115ff8282611275565b91505060208301516116146020860182611212565b509392505050565b5f61126383836115e8565b5f611630825190565b8084526020840193508360208202850161164a8560200190565b5f5b848110156115cb5783830388528151611665848261161c565b9350506020820160209890980197915060010161164c565b80515f906101208401906116918582611212565b5060208301516116a46020860182611212565b5060408301516116b76040860182611212565b5060608301516116ca6060860182611212565b5060808301516116dd6080860182611212565b5060a08301516116f060a0860182611212565b5060c083015161170360c0860182611212565b5060e083015161171660e0860182611212565b506101008301518482036101008601526117308282611627565b95945050505050565b60208082528101611263818461167d565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561120c5761120c61174a565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061179957607f821691505b6020821081036117ab576117ab611771565b50919050565b600f8152602081017f496e76616c6964206164647265737300000000000000000000000000000000008152905061156f565b6020808252810161120c816117b1565b60158152602081017f6c3143726f7373436861696e4d657373656e67657200000000000000000000008152905061156f565b60408082528101611835816117f3565b905061120c6020830184611212565b60088152602081017f6c324272696467650000000000000000000000000000000000000000000000008152905061156f565b6040808252810161183581611844565b5f61120c82611893565b90565b67ffffffffffffffff1690565b61121b81611886565b6020810161120c82846118a0565b60158152602081017f6c3243726f7373436861696e4d657373656e67657200000000000000000000008152905061156f565b60408082528101611835816118b7565b61190482848361141f565b500190565b6112638183856118f9565b5f61120c6118908381565b61192883611914565b81545f1960089490940293841b1916921b91909117905550565b5f61194e81848461191f565b505050565b8181101561196d576119655f82611942565b600101611953565b5050565b601f82111561194e575f818152602090206020601f850104810160208510156119975750805b6119a96020601f860104830182611953565b5050505050565b8267ffffffffffffffff8111156119c9576119c96112fc565b6119d38254611785565b6119de828285611971565b505f601f821160018114611a10575f83156119f95750848201355b5f19600885021c1981166002850217855550611a67565b5f84815260208120601f198516915b82811015611a3f5787850135825560209485019460019092019101611a1f565b5084821015611a5b575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b818352602083019250611a8382848361141f565b50601f01601f19160190565b60408082528101611aa1818587611a6f565b90506114b96020830184611212565b634e487b7160e01b5f52603260045260245ffd5b5f8154611ad081611785565b600182168015611ae75760018114611afc57611b2a565b60ff1983168652811515820286019350611b2a565b5f858152602090205f5b83811015611b2257815488820152600190910190602001611b06565b505081860193505b50505092915050565b61120c8183611ac4565b60088152602081017f6c314272696467650000000000000000000000000000000000000000000000008152905061156f565b6040808252810161183581611b3d56fea26469706673582212207587d0f8ad5cadd9e10d6ff4bbbad9f38f9ef1d0b68f1c52922578d9cad7cf2864736f6c634300081c0033",
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

// ROLLUPCONTRACTSLOT is a free data retrieval call binding the contract method 0x85a7805a.
//
// Solidity: function ROLLUP_CONTRACT_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) ROLLUPCONTRACTSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "ROLLUP_CONTRACT_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ROLLUPCONTRACTSLOT is a free data retrieval call binding the contract method 0x85a7805a.
//
// Solidity: function ROLLUP_CONTRACT_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) ROLLUPCONTRACTSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.ROLLUPCONTRACTSLOT(&_NetworkConfig.CallOpts)
}

// ROLLUPCONTRACTSLOT is a free data retrieval call binding the contract method 0x85a7805a.
//
// Solidity: function ROLLUP_CONTRACT_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) ROLLUPCONTRACTSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.ROLLUPCONTRACTSLOT(&_NetworkConfig.CallOpts)
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

// RollupContractAddress is a free data retrieval call binding the contract method 0x4e027271.
//
// Solidity: function rollupContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) RollupContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "rollupContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupContractAddress is a free data retrieval call binding the contract method 0x4e027271.
//
// Solidity: function rollupContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) RollupContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.RollupContractAddress(&_NetworkConfig.CallOpts)
}

// RollupContractAddress is a free data retrieval call binding the contract method 0x4e027271.
//
// Solidity: function rollupContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) RollupContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.RollupContractAddress(&_NetworkConfig.CallOpts)
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
