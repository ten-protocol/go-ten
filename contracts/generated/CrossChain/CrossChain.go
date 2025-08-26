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
	Bin: "0x608060405234801561000f575f5ffd5b50610018610025565b610020610025565b610104565b5f61002e6100c5565b805490915068010000000000000000900460ff16156100605760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100c25780546001600160401b0319166001600160401b0390811782556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916100b9916100ef565b60405180910390a15b50565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6001600160401b0382168152602081016100e9565b610c27806101115f395ff3fe608060405234801561000f575f5ffd5b50600436106100c4575f3560e01c80638da5cb5b1161007d578063e30c397811610058578063e30c39781461018f578063e874eb2014610197578063f2fde38b146101aa575f5ffd5b80638da5cb5b14610147578063a1a227fa1461015c578063a4ab2faa1461017c575f5ffd5b8063715018a6116100ad578063715018a61461011557806379ba50971461011d5780638415482614610125575f5ffd5b80632f0cb9e3146100c8578063485cc95514610100575b5f5ffd5b6100ea6100d6366004610798565b5f6020819052908152604090205460ff1681565b6040516100f791906107c6565b60405180910390f35b61011361010e3660046107f8565b6101bd565b005b61011361037e565b61011361039e565b6100ea610133366004610798565b60016020525f908152604090205460ff1681565b61014f6103dd565b6040516100f79190610837565b60025461016f906001600160a01b031681565b6040516100f79190610862565b6100ea61018a366004610a0f565b610411565b61014f61048d565b60035461016f906001600160a01b031681565b6101136101b8366004610a4f565b6104b5565b5f6101c6610547565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f811580156101f25750825b90505f8267ffffffffffffffff16600114801561020e5750303b155b90508115801561021c575080155b15610253576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561028757845468ff00000000000000001916680100000000000000001785555b6001600160a01b0386166102b65760405162461bcd60e51b81526004016102ad90610aa0565b60405180910390fd5b6001600160a01b0387166102dc5760405162461bcd60e51b81526004016102ad90610ae2565b6102e587610571565b6102ed61058a565b600380546001600160a01b03881673ffffffffffffffffffffffffffffffffffffffff199182168117909255600280549091169091179055831561037557845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061036c90600190610b0c565b60405180910390a15b50505050505050565b61038661059c565b60405162461bcd60e51b81526004016102ad90610b1a565b33806103a861048d565b6001600160a01b0316146103d1578060405163118cdaa760e01b81526004016102ad9190610837565b6103da816105ce565b50565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b5f80805b8351811015610475578184828151811061043157610431610b7b565b602002602001015161044290610b98565b604051602001610453929190610bd6565b60408051601f1981840301815291905280516020909101209150600101610415565b505f9081526001602052604090205460ff1692915050565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610401565b6104bd61059c565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117825561050e6103dd565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b610579610617565b61058281610655565b6103da610666565b610592610617565b61059a61066e565b565b336105a56103dd565b6001600160a01b03161461059a573360405163118cdaa760e01b81526004016102ad9190610837565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff191681556106138261069c565b5050565b61061f610719565b61059a576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61065d610617565b6103da81610737565b61059a610617565b610676610617565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f610722610547565b5468010000000000000000900460ff16919050565b61073f610617565b6001600160a01b0381166103d1575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016102ad9190610837565b805b81146103da575f5ffd5b803561056b81610781565b5f602082840312156107ab576107ab5f5ffd5b6107b5838361078d565b9392505050565b8015155b82525050565b6020810161056b82846107bc565b5f6001600160a01b03821661056b565b610783816107d4565b803561056b816107e4565b5f5f6040838503121561080c5761080c5f5ffd5b61081684846107ed565b915061082584602085016107ed565b90509250929050565b6107c0816107d4565b6020810161056b828461082e565b5f61056b826107d4565b5f61056b82610845565b6107c08161084f565b6020810161056b8284610859565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff821117156108aa576108aa610870565b6040525050565b5f6108bb60405190565b90506108c78282610884565b919050565b5f67ffffffffffffffff8211156108e5576108e5610870565b5060209081020190565b5f67ffffffffffffffff82111561090857610908610870565b601f19601f83011660200192915050565b82818337505f910152565b5f610936610931846108ef565b6108b1565b905082815283838301111561094c5761094c5f5ffd5b6107b5836020830184610919565b5f82601f83011261096c5761096c5f5ffd5b6107b583833560208501610924565b5f610988610931846108cc565b838152905060208082019084028301858111156109a6576109a65f5ffd5b835b818110156109e457803567ffffffffffffffff8111156109c9576109c95f5ffd5b6109d58882880161095a565b845250602092830192016109a8565b5050509392505050565b5f82601f830112610a0057610a005f5ffd5b6107b58383356020850161097b565b5f60208284031215610a2257610a225f5ffd5b813567ffffffffffffffff811115610a3b57610a3b5f5ffd5b610a47848285016109ee565b949350505050565b5f60208284031215610a6257610a625f5ffd5b6107b583836107ed565b601b8152602081017f496e76616c6964206d6573736167652062757320616464726573730000000000815290505b60200190565b6020808252810161056b81610a6c565b60138152602081017f4f776e65722063616e6e6f74206265203078300000000000000000000000000081529050610a9a565b6020808252810161056b81610ab0565b5f67ffffffffffffffff821661056b565b6107c081610af2565b6020810161056b8284610b03565b6020808252810161056b81603481527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60208201527f742072656e6f756e6365206f776e657273686970000000000000000000000000604082015260600190565b634e487b7160e01b5f52603260045260245ffd5b5f61056b825190565b5f610ba1825190565b60208301610bae81610b8f565b9250506020811015610bca575f1960086020839003021b821691505b50919050565b806107c0565b60408101610be48285610bd0565b6107b56020830184610bd056fea264697066735822122048dbaaad494a0afb78bbdef21571fc4167001c2e551eafd45df275ab5aa16a7364736f6c634300081c0033",
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
