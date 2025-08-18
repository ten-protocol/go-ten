// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package CrossChain

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

// CrossChainMetaData contains all meta data concerning the CrossChain contract.
var CrossChainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_messageBus\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"crossChainHashes\",\"type\":\"bytes[]\"}],\"name\":\"isBundleAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"bundleHash\",\"type\":\"bytes32\"}],\"name\":\"isBundleSaved\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isBundleSaved\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"withdrawalHash\",\"type\":\"bytes32\"}],\"name\":\"isWithdrawalSpent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isWithdrawalSpent\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"contractIMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461002257610011610026565b604051610e936101968239610e9390f35b5f80fd5b61002e610038565b6100366100b9565b565b61003661002e565b61004d9060401c60ff1690565b90565b61004d9054610040565b61004d905b6001600160401b031690565b61004d905461005a565b61004d9061005f906001600160401b031682565b9061009961004d6100b592610075565b82546001600160401b0319166001600160401b03919091161790565b9055565b5f6100c261014f565b016100cc81610050565b61013e576100d98161006b565b6001600160401b03919082908116036100f0575050565b8161011f7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29361013993610089565b604051918291826001600160401b03909116815260200190565b0390a1565b63f92ee8a960e01b5f908152600490fd5b61004d61018d565b61004d61004d61004d9290565b61004d7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610157565b61004d61016456fe60806040526004361015610011575f80fd5b5f3560e01c80632f0cb9e3146100c0578063485cc955146100bb578063715018a6146100b657806379ba5097146100b157806384154826146100ac5780638da5cb5b146100a7578063a1a227fa146100a2578063a4ab2faa1461009d578063e30c397814610098578063e874eb20146100935763f2fde38b036100ce576104e6565b6104b7565b610491565b610476565b6102d4565b610234565b610219565b6101ea565b6101d5565b6101ad565b610131565b805b036100ce57565b5f80fd5b905035906100df826100c5565b565b906020828203126100ce576100f5916100d2565b90565b6100f5916008021c5b60ff1690565b906100f591546100f8565b5f6101286100f59282905f5260205260405f2090565b610107565b9052565b346100ce5761015e61014c6101473660046100e1565b610112565b60405191829182901515815260200190565b0390f35b6001600160a01b031690565b6001600160a01b0381166100c7565b905035906100df8261016e565b91906040838203126100ce576100f59060206101a6828661017d565b940161017d565b346100ce576101c66101c036600461018a565b9061086f565b604051005b5f9103126100ce57565b346100ce576101e53660046101cb565b6108ec565b346100ce576101fa3660046101cb565b6101c66108f1565b5f6101286100f5926001905f5260205260405f2090565b346100ce5761015e61014c61022f3660046100e1565b610202565b346100ce576102443660046101cb565b61015e61024f61096f565b604051918291826001600160a01b03909116815260200190565b6100f5916008021c6001600160a01b031690565b906100f59154610269565b6100f55f600261027d565b6101626100f56100f5926001600160a01b031690565b6100f590610293565b6100f5906102a9565b61012d906102b2565b6020810192916100df91906102bb565b346100ce576102e43660046101cb565b61015e6102ef610288565b604051918291826102c4565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761033157604052565b6102fb565b906100df61034360405190565b928361030f565b67ffffffffffffffff81116103315760208091020190565b67ffffffffffffffff811161033157602090601f01601f19160190565b90825f939282370152565b9092919261039f61039a82610362565b610336565b93818552818301116100ce576100df91602085019061037f565b9080601f830112156100ce578160206100f59335910161038a565b9291906103e361039a8261034a565b93818552602080860192028101918383116100ce5781905b838210610409575050505050565b813567ffffffffffffffff81116100ce5760209161042a87849387016103b9565b8152019101906103fb565b9080601f830112156100ce578160206100f5933591016103d4565b906020828203126100ce57813567ffffffffffffffff81116100ce576100f59201610435565b346100ce5761015e61014c61048c366004610450565b610a2f565b346100ce576104a13660046101cb565b61015e61024f610ad5565b6100f55f600361027d565b346100ce576104c73660046101cb565b61015e6102ef6104ac565b906020828203126100ce576100f59161017d565b346100ce576101c66104f93660046104d2565b610b79565b6100f59060401c610101565b6100f590546104fe565b6100f5905b67ffffffffffffffff1690565b6100f59054610514565b6105196100f56100f59290565b6100f56100f56100f59290565b9067ffffffffffffffff905b9181191691161790565b6105196100f56100f59267ffffffffffffffff1690565b906105876100f561058e92610560565b825461054a565b9055565b9068ff00000000000000009060401b610556565b906105b66100f561058e92151590565b8254610592565b61012d90610530565b6020810192916100df91906105bd565b6105de610b82565b9081906105fa6105f46105f08461050a565b1590565b92610526565b936106045f610530565b67ffffffffffffffff8616148061071f575b60019561063361062588610530565b9167ffffffffffffffff1690565b1490816106f7575b155b90816106ee575b506106c45761066d91836106645f61065b89610530565b97019687610577565b6106b5576107fa565b610675575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916106a45f6106b0936105a6565b604051918291826105c6565b0390a1565b6106bf86866105a6565b6107fa565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f610644565b905061063d610705306102b2565b3b6107166107125f61053d565b9190565b1491905061063b565b5082610616565b6101626100f56100f59290565b6100f590610726565b1561074357565b60405162461bcd60e51b815260206004820152601b60248201527f496e76616c6964206d65737361676520627573206164647265737300000000006044820152606490fd5b1561078f57565b60405162461bcd60e51b815260206004820152601360248201527f4f776e65722063616e6e6f7420626520307830000000000000000000000000006044820152606490fd5b906001600160a01b0390610556565b906107f36100f561058e926102b2565b82546107d4565b6100df91610848610868926108436108326108145f610733565b6101626001600160a01b0382166001600160a01b038816141561073c565b6001600160a01b0383161415610788565b610ba7565b610850610bc3565b61086361085c826102b2565b60036107e3565b6102b2565b60026107e3565b906100df916105d6565b610881610bcb565b60405162461bcd60e51b815260206004820152603460248201527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60448201527f742072656e6f756e6365206f776e6572736869700000000000000000000000006064820152608490fd5b610879565b336108fa610ad5565b6109156001600160a01b0383165b916001600160a01b031690565b03610923576100df90610c29565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b6100f590610162565b6100f5905461095c565b6100f55f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b01610965565b634e487b7160e01b5f52603260045260245ffd5b906109b8825190565b8110156109c9576020809102010190565b61099b565b6109e76109e36109dc835190565b9260200190565b5190565b90602081106109f4575090565b610a05905f19906020036008021b90565b1690565b9081526040810192916100df9160200152565b6100f590610101565b6100f59054610a1c565b610a385f61053d565b90610a425f61053d565b915b610a4f6100f5835190565b831015610ab757610a8b610a97610ab192610a73610a6d87876109af565b516109ce565b90610a7d60405190565b938492602084019283610a09565b9081038252038261030f565b610aa9610aa2825190565b9160200190565b209260010190565b91610a44565b6100f59250610ad091506001905f5260205260405f2090565b610a25565b6100f55f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610995565b6100df90610b0b610bcb565b610b35817f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c006107e3565b610b49610b4361086361096f565b916102b2565b907f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700610b7460405190565b5f90a3565b6100df90610aff565b6100f5610c85565b6100df90610b96610c8d565b610b9f90610cdd565b6100df610cee565b6100df90610b8a565b610bb8610c8d565b6100df6100df610d64565b6100df610bb0565b610bd361096f565b3390610bde82610908565b036109235750565b91906008610556910291610c006001600160a01b03841b90565b921b90565b9190610c166100f561058e936102b2565b908354610be6565b6100df915f91610c05565b6100df90610c575f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610c1e565b610d6c565b6100f57ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061053d565b6100f5610c5c565b610c986105f0610dd1565b610c9e57565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f908152600490fd5b6100df90610cd4610c8d565b6100df90610e54565b6100df90610cc8565b6100df610c8d565b6100df610ce6565b610cfe610c8d565b6100df610d30565b6100f5600161053d565b905f1990610556565b90610d296100f561058e9261053d565b8254610d10565b6100df7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f610d5d610d06565b9101610d19565b6100df610cf6565b610da6610b437f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930061086384610da083610965565b926107e3565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0610b7460405190565b6100f55f610ddd610b82565b0161050a565b6100df90610def610c8d565b610df85f610733565b6001600160a01b0381166001600160a01b03831614610e1b57506100df90610c29565b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b6100df90610de356fea26469706673582212200fa1a4ae32dd2d02553b5bf7f37203af5b16ac30e9a72855564f558202ee01bd64736f6c634300081c0033",
}

// CrossChainABI is the input ABI used to generate the binding from.
// Deprecated: Use CrossChainMetaData.ABI instead.
var CrossChainABI = CrossChainMetaData.ABI

// CrossChainBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrossChainMetaData.Bin instead.
var CrossChainBin = CrossChainMetaData.Bin

// DeployCrossChain deploys a new Ethereum contract, binding an instance of CrossChain to it.
func DeployCrossChain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrossChain, error) {
	parsed, err := CrossChainMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChain{CrossChainCaller: CrossChainCaller{contract: contract}, CrossChainTransactor: CrossChainTransactor{contract: contract}, CrossChainFilterer: CrossChainFilterer{contract: contract}}, nil
}

// CrossChain is an auto generated Go binding around an Ethereum contract.
type CrossChain struct {
	CrossChainCaller     // Read-only binding to the contract
	CrossChainTransactor // Write-only binding to the contract
	CrossChainFilterer   // Log filterer for contract events
}

// CrossChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrossChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrossChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrossChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrossChainSession struct {
	Contract     *CrossChain       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrossChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrossChainCallerSession struct {
	Contract *CrossChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CrossChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrossChainTransactorSession struct {
	Contract     *CrossChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CrossChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrossChainRaw struct {
	Contract *CrossChain // Generic contract binding to access the raw methods on
}

// CrossChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrossChainCallerRaw struct {
	Contract *CrossChainCaller // Generic read-only contract binding to access the raw methods on
}

// CrossChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrossChainTransactorRaw struct {
	Contract *CrossChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrossChain creates a new instance of CrossChain, bound to a specific deployed contract.
func NewCrossChain(address common.Address, backend bind.ContractBackend) (*CrossChain, error) {
	contract, err := bindCrossChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChain{CrossChainCaller: CrossChainCaller{contract: contract}, CrossChainTransactor: CrossChainTransactor{contract: contract}, CrossChainFilterer: CrossChainFilterer{contract: contract}}, nil
}

// NewCrossChainCaller creates a new read-only instance of CrossChain, bound to a specific deployed contract.
func NewCrossChainCaller(address common.Address, caller bind.ContractCaller) (*CrossChainCaller, error) {
	contract, err := bindCrossChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainCaller{contract: contract}, nil
}

// NewCrossChainTransactor creates a new write-only instance of CrossChain, bound to a specific deployed contract.
func NewCrossChainTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainTransactor, error) {
	contract, err := bindCrossChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainTransactor{contract: contract}, nil
}

// NewCrossChainFilterer creates a new log filterer instance of CrossChain, bound to a specific deployed contract.
func NewCrossChainFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainFilterer, error) {
	contract, err := bindCrossChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainFilterer{contract: contract}, nil
}

// bindCrossChain binds a generic wrapper to an already deployed contract.
func bindCrossChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrossChainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChain *CrossChainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChain.Contract.CrossChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChain *CrossChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChain.Contract.CrossChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChain *CrossChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChain.Contract.CrossChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChain *CrossChainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChain *CrossChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChain *CrossChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChain.Contract.contract.Transact(opts, method, params...)
}

// IsBundleAvailable is a free data retrieval call binding the contract method 0xa4ab2faa.
//
// Solidity: function isBundleAvailable(bytes[] crossChainHashes) view returns(bool)
func (_CrossChain *CrossChainCaller) IsBundleAvailable(opts *bind.CallOpts, crossChainHashes [][]byte) (bool, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "isBundleAvailable", crossChainHashes)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBundleAvailable is a free data retrieval call binding the contract method 0xa4ab2faa.
//
// Solidity: function isBundleAvailable(bytes[] crossChainHashes) view returns(bool)
func (_CrossChain *CrossChainSession) IsBundleAvailable(crossChainHashes [][]byte) (bool, error) {
	return _CrossChain.Contract.IsBundleAvailable(&_CrossChain.CallOpts, crossChainHashes)
}

// IsBundleAvailable is a free data retrieval call binding the contract method 0xa4ab2faa.
//
// Solidity: function isBundleAvailable(bytes[] crossChainHashes) view returns(bool)
func (_CrossChain *CrossChainCallerSession) IsBundleAvailable(crossChainHashes [][]byte) (bool, error) {
	return _CrossChain.Contract.IsBundleAvailable(&_CrossChain.CallOpts, crossChainHashes)
}

// IsBundleSaved is a free data retrieval call binding the contract method 0x84154826.
//
// Solidity: function isBundleSaved(bytes32 bundleHash) view returns(bool isBundleSaved)
func (_CrossChain *CrossChainCaller) IsBundleSaved(opts *bind.CallOpts, bundleHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "isBundleSaved", bundleHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBundleSaved is a free data retrieval call binding the contract method 0x84154826.
//
// Solidity: function isBundleSaved(bytes32 bundleHash) view returns(bool isBundleSaved)
func (_CrossChain *CrossChainSession) IsBundleSaved(bundleHash [32]byte) (bool, error) {
	return _CrossChain.Contract.IsBundleSaved(&_CrossChain.CallOpts, bundleHash)
}

// IsBundleSaved is a free data retrieval call binding the contract method 0x84154826.
//
// Solidity: function isBundleSaved(bytes32 bundleHash) view returns(bool isBundleSaved)
func (_CrossChain *CrossChainCallerSession) IsBundleSaved(bundleHash [32]byte) (bool, error) {
	return _CrossChain.Contract.IsBundleSaved(&_CrossChain.CallOpts, bundleHash)
}

// IsWithdrawalSpent is a free data retrieval call binding the contract method 0x2f0cb9e3.
//
// Solidity: function isWithdrawalSpent(bytes32 withdrawalHash) view returns(bool isWithdrawalSpent)
func (_CrossChain *CrossChainCaller) IsWithdrawalSpent(opts *bind.CallOpts, withdrawalHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "isWithdrawalSpent", withdrawalHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWithdrawalSpent is a free data retrieval call binding the contract method 0x2f0cb9e3.
//
// Solidity: function isWithdrawalSpent(bytes32 withdrawalHash) view returns(bool isWithdrawalSpent)
func (_CrossChain *CrossChainSession) IsWithdrawalSpent(withdrawalHash [32]byte) (bool, error) {
	return _CrossChain.Contract.IsWithdrawalSpent(&_CrossChain.CallOpts, withdrawalHash)
}

// IsWithdrawalSpent is a free data retrieval call binding the contract method 0x2f0cb9e3.
//
// Solidity: function isWithdrawalSpent(bytes32 withdrawalHash) view returns(bool isWithdrawalSpent)
func (_CrossChain *CrossChainCallerSession) IsWithdrawalSpent(withdrawalHash [32]byte) (bool, error) {
	return _CrossChain.Contract.IsWithdrawalSpent(&_CrossChain.CallOpts, withdrawalHash)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_CrossChain *CrossChainCaller) MerkleMessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "merkleMessageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_CrossChain *CrossChainSession) MerkleMessageBus() (common.Address, error) {
	return _CrossChain.Contract.MerkleMessageBus(&_CrossChain.CallOpts)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_CrossChain *CrossChainCallerSession) MerkleMessageBus() (common.Address, error) {
	return _CrossChain.Contract.MerkleMessageBus(&_CrossChain.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChain *CrossChainCaller) MessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "messageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChain *CrossChainSession) MessageBus() (common.Address, error) {
	return _CrossChain.Contract.MessageBus(&_CrossChain.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChain *CrossChainCallerSession) MessageBus() (common.Address, error) {
	return _CrossChain.Contract.MessageBus(&_CrossChain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CrossChain *CrossChainCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CrossChain *CrossChainSession) Owner() (common.Address, error) {
	return _CrossChain.Contract.Owner(&_CrossChain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CrossChain *CrossChainCallerSession) Owner() (common.Address, error) {
	return _CrossChain.Contract.Owner(&_CrossChain.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CrossChain *CrossChainCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CrossChain *CrossChainSession) PendingOwner() (common.Address, error) {
	return _CrossChain.Contract.PendingOwner(&_CrossChain.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_CrossChain *CrossChainCallerSession) PendingOwner() (common.Address, error) {
	return _CrossChain.Contract.PendingOwner(&_CrossChain.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_CrossChain *CrossChainCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_CrossChain *CrossChainSession) RenounceOwnership() error {
	return _CrossChain.Contract.RenounceOwnership(&_CrossChain.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_CrossChain *CrossChainCallerSession) RenounceOwnership() error {
	return _CrossChain.Contract.RenounceOwnership(&_CrossChain.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CrossChain *CrossChainTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CrossChain *CrossChainSession) AcceptOwnership() (*types.Transaction, error) {
	return _CrossChain.Contract.AcceptOwnership(&_CrossChain.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_CrossChain *CrossChainTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _CrossChain.Contract.AcceptOwnership(&_CrossChain.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address owner, address _messageBus) returns()
func (_CrossChain *CrossChainTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, _messageBus common.Address) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "initialize", owner, _messageBus)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address owner, address _messageBus) returns()
func (_CrossChain *CrossChainSession) Initialize(owner common.Address, _messageBus common.Address) (*types.Transaction, error) {
	return _CrossChain.Contract.Initialize(&_CrossChain.TransactOpts, owner, _messageBus)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address owner, address _messageBus) returns()
func (_CrossChain *CrossChainTransactorSession) Initialize(owner common.Address, _messageBus common.Address) (*types.Transaction, error) {
	return _CrossChain.Contract.Initialize(&_CrossChain.TransactOpts, owner, _messageBus)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CrossChain *CrossChainTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CrossChain *CrossChainSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CrossChain.Contract.TransferOwnership(&_CrossChain.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CrossChain *CrossChainTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CrossChain.Contract.TransferOwnership(&_CrossChain.TransactOpts, newOwner)
}

// CrossChainInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CrossChain contract.
type CrossChainInitializedIterator struct {
	Event *CrossChainInitialized // Event containing the contract specifics and raw log

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
func (it *CrossChainInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainInitialized)
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
		it.Event = new(CrossChainInitialized)
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
func (it *CrossChainInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainInitialized represents a Initialized event raised by the CrossChain contract.
type CrossChainInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChain *CrossChainFilterer) FilterInitialized(opts *bind.FilterOpts) (*CrossChainInitializedIterator, error) {

	logs, sub, err := _CrossChain.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CrossChainInitializedIterator{contract: _CrossChain.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChain *CrossChainFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CrossChainInitialized) (event.Subscription, error) {

	logs, sub, err := _CrossChain.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainInitialized)
				if err := _CrossChain.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CrossChain *CrossChainFilterer) ParseInitialized(log types.Log) (*CrossChainInitialized, error) {
	event := new(CrossChainInitialized)
	if err := _CrossChain.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the CrossChain contract.
type CrossChainOwnershipTransferStartedIterator struct {
	Event *CrossChainOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *CrossChainOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainOwnershipTransferStarted)
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
		it.Event = new(CrossChainOwnershipTransferStarted)
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
func (it *CrossChainOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the CrossChain contract.
type CrossChainOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CrossChain *CrossChainFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CrossChainOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CrossChain.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainOwnershipTransferStartedIterator{contract: _CrossChain.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_CrossChain *CrossChainFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *CrossChainOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CrossChain.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainOwnershipTransferStarted)
				if err := _CrossChain.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_CrossChain *CrossChainFilterer) ParseOwnershipTransferStarted(log types.Log) (*CrossChainOwnershipTransferStarted, error) {
	event := new(CrossChainOwnershipTransferStarted)
	if err := _CrossChain.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CrossChain contract.
type CrossChainOwnershipTransferredIterator struct {
	Event *CrossChainOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CrossChainOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainOwnershipTransferred)
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
		it.Event = new(CrossChainOwnershipTransferred)
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
func (it *CrossChainOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainOwnershipTransferred represents a OwnershipTransferred event raised by the CrossChain contract.
type CrossChainOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CrossChain *CrossChainFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CrossChainOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CrossChain.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainOwnershipTransferredIterator{contract: _CrossChain.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CrossChain *CrossChainFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CrossChainOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CrossChain.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainOwnershipTransferred)
				if err := _CrossChain.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CrossChain *CrossChainFilterer) ParseOwnershipTransferred(log types.Log) (*CrossChainOwnershipTransferred, error) {
	event := new(CrossChainOwnershipTransferred)
	if err := _CrossChain.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
