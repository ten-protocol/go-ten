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
	CrossChain             common.Address
	MessageBus             common.Address
	NetworkEnclaveRegistry common.Address
	RollupContract         common.Address
}

// NetworkConfigMetaData contains all meta data concerning the NetworkConfig contract.
var NetworkConfigMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MESSAGE_BUS_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NETWORK_ENCLAVE_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROLLUP_CONTRACT_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addresses\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rollupContract\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.Addresses\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rollupContract\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.Addresses\",\"name\":\"_addresses\",\"type\":\"tuple\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkEnclaveRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506109fb8061001c5f395ff3fe608060405234801561000f575f5ffd5b50600436106100da575f3560e01c806385a7805a11610088578063da0321cd11610063578063da0321cd1461017e578063f2fde38b14610193578063f5e9f286146101a6578063fbfd6d91146101ae575f5ffd5b806385a7805a1461013e5780638da5cb5b14610146578063a1b918d614610176575f5ffd5b80634e027271116100b85780634e02727114610126578063715018a61461012e57806372bad91214610136575f5ffd5b80630415cc15146100de5780632fc00c76146100f357806348d8723914610111575b5f5ffd5b6100f16100ec3660046108ac565b6101b6565b005b6100fb6103cf565b60405161010891906108df565b60405180910390f35b610119610407565b60405161010891906108f3565b6100fb610435565b6100f1610464565b610119610477565b6101196104a2565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166100fb565b6100fb6104cd565b6101866104fc565b604051610108919061094b565b6100f16101a1366004610959565b610584565b6100fb6105e3565b610119610612565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156102005750825b90505f8267ffffffffffffffff16600114801561021c5750303b155b90508115801561022a575080155b15610261576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561029557845468ff00000000000000001916680100000000000000001785555b61029e3361063d565b6102d26102cc60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a78307195861098a565b87519055565b61030961030060017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c61098a565b60208801519055565b61034061033760017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba61098a565b60408801519055565b61037761036e60017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc61098a565b60608801519055565b83156103c757845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906103be906001906109b7565b60405180910390a15b505050505050565b5f6104026103fe60017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba61098a565b5490565b905090565b61043260017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c61098a565b81565b5f6104026103fe60017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc61098a565b61046c61064e565b6104755f6106c2565b565b61043260017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba61098a565b61043260017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc61098a565b5f6104026103fe60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a78307195861098a565b604080516080810182525f80825260208201819052918101829052606081019190915260405180608001604052806105326104cd565b6001600160a01b031681526020016105486105e3565b6001600160a01b0316815260200161055e6103cf565b6001600160a01b03168152602001610574610435565b6001600160a01b03169052919050565b61058c61064e565b6001600160a01b0381166105d7575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016105ce91906108df565b60405180910390fd5b6105e0816106c2565b50565b5f6104026103fe60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c61098a565b61043260017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a78307195861098a565b61064561074a565b6105e0816107b1565b336106807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461047557336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016105ce91906108df565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080547fffffffffffffffffffffffff000000000000000000000000000000000000000081166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610475576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61058c61074a565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff821117156107f3576107f36107b9565b6040525050565b5f61080460405190565b905061081082826107cd565b919050565b5f6001600160a01b0382165b92915050565b61083081610815565b81146105e0575f5ffd5b803561082181610827565b5f60808284031215610858576108585f5ffd5b61086260806107fa565b905061086e838361083a565b815261087d836020840161083a565b602082015261088f836040840161083a565b60408201526108a1836060840161083a565b606082015292915050565b5f608082840312156108bf576108bf5f5ffd5b6108c98383610845565b9392505050565b6108d981610815565b82525050565b6020810161082182846108d0565b806108d9565b6020810161082182846108ed565b805161090d83826108d0565b50602081015161092060208401826108d0565b50604081015161093360408401826108d0565b50606081015161094660608401826108d0565b505050565b608081016108218284610901565b5f6020828403121561096c5761096c5f5ffd5b6108c9838361083a565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561082157610821610976565b5f67ffffffffffffffff8216610821565b6108d98161099d565b6020810161082182846109ae56fea2646970667358221220ebfec208873fdc6923d769cf7314cb2e1d25d5ff24112fdd8b87915d41bdd74764736f6c634300081c0033",
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
func 	NewNetworkConfigCaller(address common.Address, caller bind.ContractCaller) (*NetworkConfigCaller, error) {
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

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address))
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
// Solidity: function addresses() view returns((address,address,address,address))
func (_NetworkConfig *NetworkConfigSession) Addresses() (NetworkConfigAddresses, error) {
	return _NetworkConfig.Contract.Addresses(&_NetworkConfig.CallOpts)
}

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address))
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

// Initialize is a paid mutator transaction binding the contract method 0x0415cc15.
//
// Solidity: function initialize((address,address,address,address) _addresses) returns()
func (_NetworkConfig *NetworkConfigTransactor) Initialize(opts *bind.TransactOpts, _addresses NetworkConfigAddresses) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "initialize", _addresses)
}

// Initialize is a paid mutator transaction binding the contract method 0x0415cc15.
//
// Solidity: function initialize((address,address,address,address) _addresses) returns()
func (_NetworkConfig *NetworkConfigSession) Initialize(_addresses NetworkConfigAddresses) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses)
}

// Initialize is a paid mutator transaction binding the contract method 0x0415cc15.
//
// Solidity: function initialize((address,address,address,address) _addresses) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) Initialize(_addresses NetworkConfigAddresses) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses)
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
