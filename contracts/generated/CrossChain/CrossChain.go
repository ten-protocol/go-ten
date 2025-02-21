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

// StructsValueTransferMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsValueTransferMessage struct {
	Sender   common.Address
	Receiver common.Address
	Amount   *big.Int
	Sequence uint64
}

// CrossChainMetaData contains all meta data concerning the CrossChain contract.
var CrossChainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"messageBusAddress\",\"type\":\"address\"}],\"name\":\"LogManagementContractCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"WithdrawalsPaused\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.ValueTransferMessage\",\"name\":\"_msg\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"extractNativeValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChallengePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_messageBus\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"crossChainHashes\",\"type\":\"bytes[]\"}],\"name\":\"isBundleAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isBundleSaved\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isWithdrawalSpent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"contractIMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_pause\",\"type\":\"bool\"}],\"name\":\"pauseWithdrawals\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_delay\",\"type\":\"uint256\"}],\"name\":\"setChallengePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b610f85806100975f395ff3fe608060405234801561000f575f5ffd5b50600436106100da575f3560e01c80638da5cb5b11610088578063c4d66de811610063578063c4d66de8146101e3578063e874eb20146101f6578063f2fde38b14610209578063f4cc87ba1461021c575f5ffd5b80638da5cb5b14610178578063a1a227fa146101b0578063a4ab2faa146101d0575f5ffd5b8063715018a6116100b8578063715018a61461013e5780637864b77d146101465780638415482614610156575f5ffd5b80632f0cb9e3146100de5780635d475fdd146101165780636af526621461012b575b5f5ffd5b6101006100ec366004610915565b60046020525f908152604090205460ff1681565b60405161010d9190610943565b60405180910390f35b610129610124366004610915565b61022f565b005b6101296101393660046109b8565b61023c565b6101296103ff565b60015460405161010d9190610a25565b610100610164366004610915565b60056020525f908152604090205460ff1681565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031660405161010d9190610a4c565b6002546101c3906001600160a01b031681565b60405161010d9190610a77565b6101006101de366004610c24565b610412565b6101296101f1366004610c78565b61048e565b6003546101c3906001600160a01b031681565b610129610217366004610c78565b610640565b61012961022a366004610ca8565b610696565b6102376106e6565b600155565b61024461075a565b5f5460ff161561026f5760405162461bcd60e51b815260040161026690610cf9565b60405180910390fd5b6003546040517fb201246f0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063b201246f906102be908790879087908790600401610e24565b5f6040518083038186803b1580156102d4575f5ffd5b505afa1580156102e6573d5f5f3e3d5ffd5b505050505f846040516020016102fc9190610e5d565b60408051601f1981840301815291815281516020928301205f818152600490935291205490915060ff16156103435760405162461bcd60e51b815260040161026690610e9d565b5f81815260046020908152604091829020805460ff191660011790556002546001600160a01b0316916399a3ad2191610380918901908901610c78565b87604001356040518363ffffffff1660e01b81526004016103a2929190610ead565b5f604051808303815f87803b1580156103b9575f5ffd5b505af11580156103cb573d5f5f3e3d5ffd5b50505050506103f960017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b50505050565b6104076106e6565b6104105f6107e3565b565b5f80805b8351811015610476578184828151811061043257610432610ec8565b602002602001015161044390610ee5565b604051602001610454929190610f19565b60408051601f1981840301815291905280516020909101209150600101610416565b505f9081526005602052604090205460ff1692915050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156104d85750825b90505f8267ffffffffffffffff1660011480156104f45750303b155b905081158015610502575080155b15610539576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561056d57845468ff00000000000000001916680100000000000000001785555b61057633610860565b61057e610871565b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0388169081179091555f805460ff191690556040517fbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf916105e091610a4c565b60405180910390a1831561063857845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061062f90600190610f41565b60405180910390a15b505050505050565b6106486106e6565b6001600160a01b03811661068a575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016102669190610a4c565b610693816107e3565b50565b61069e6106e6565b5f805460ff19168215151790556040517f129d33f7856617012aed60524381cfff7233cfc57df58d9f6613a5593d3dc218906106db908390610943565b60405180910390a150565b336107187f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461041057336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016102669190610a4c565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f008054600119016107b7576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60029055565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b610868610881565b610693816108e8565b610879610881565b6104106108f0565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610410576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610648610881565b6107bd610881565b805b8114610693575f5ffd5b803561090f816108f8565b92915050565b5f60208284031215610928576109285f5ffd5b6109328383610904565b9392505050565b8015155b82525050565b6020810161090f8284610939565b5f60808284031215610964576109645f5ffd5b50919050565b5f5f83601f84011261097d5761097d5f5ffd5b50813567ffffffffffffffff811115610997576109975f5ffd5b6020830191508360208202830111156109b1576109b15f5ffd5b9250929050565b5f5f5f5f60c085870312156109ce576109ce5f5ffd5b6109d88686610951565b9350608085013567ffffffffffffffff8111156109f6576109f65f5ffd5b610a028782880161096a565b9350935050610a148660a08701610904565b905092959194509250565b8061093d565b6020810161090f8284610a1f565b5f6001600160a01b03821661090f565b61093d81610a33565b6020810161090f8284610a43565b5f61090f82610a33565b5f61090f82610a5a565b61093d81610a64565b6020810161090f8284610a6e565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff82111715610abf57610abf610a85565b6040525050565b5f610ad060405190565b9050610adc8282610a99565b919050565b5f67ffffffffffffffff821115610afa57610afa610a85565b5060209081020190565b5f67ffffffffffffffff821115610b1d57610b1d610a85565b601f19601f83011660200192915050565b82818337505f910152565b5f610b4b610b4684610b04565b610ac6565b9050828152838383011115610b6157610b615f5ffd5b610932836020830184610b2e565b5f82601f830112610b8157610b815f5ffd5b61093283833560208501610b39565b5f610b9d610b4684610ae1565b83815290506020808201908402830185811115610bbb57610bbb5f5ffd5b835b81811015610bf957803567ffffffffffffffff811115610bde57610bde5f5ffd5b610bea88828801610b6f565b84525060209283019201610bbd565b5050509392505050565b5f82601f830112610c1557610c155f5ffd5b61093283833560208501610b90565b5f60208284031215610c3757610c375f5ffd5b813567ffffffffffffffff811115610c5057610c505f5ffd5b610c5c84828501610c03565b949350505050565b6108fa81610a33565b803561090f81610c64565b5f60208284031215610c8b57610c8b5f5ffd5b6109328383610c6d565b8015156108fa565b803561090f81610c95565b5f60208284031215610cbb57610cbb5f5ffd5b6109328383610c9d565b60168152602081017f7769746864726177616c73206172652070617573656400000000000000000000815290505b60200190565b6020808252810161090f81610cc5565b505f61090f6020830183610c6d565b505f61090f6020830183610904565b67ffffffffffffffff81166108fa565b803561090f81610d27565b505f61090f6020830183610d37565b67ffffffffffffffff811661093d565b610d6b8180610d09565b610d758382610a43565b50610d836020820182610d09565b610d906020840182610a43565b50610d9e6040820182610d18565b610dab6040840182610a1f565b50610db96060820182610d42565b610dc66060840182610d51565b505050565b82818337505050565b8183526020830192505f7f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115610e0d57610e0d5f5ffd5b602083029250610e1e838584610dcb565b50500190565b60c08101610e328287610d61565b8181036080830152610e45818587610dd4565b9050610e5460a0830184610a1f565b95945050505050565b6080810161090f8284610d61565b60188152602081017f7769746864726177616c20616c7265616479207370656e74000000000000000081529050610cf3565b6020808252810161090f81610e6b565b60408101610ebb8285610a43565b6109326020830184610a1f565b634e487b7160e01b5f52603260045260245ffd5b5f61090f825190565b5f610eee825190565b60208301610efb81610edc565b9250506020811015610964575f196020919091036008021b16919050565b60408101610ebb8285610a1f565b5f67ffffffffffffffff821661090f565b61093d81610f27565b6020810161090f8284610f3856fea2646970667358221220b400986b0cceb989195ab181101a7fcf861d508f34577a02d9220f91374b6f2a64736f6c634300081c0033",
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

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_CrossChain *CrossChainCaller) GetChallengePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "getChallengePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_CrossChain *CrossChainSession) GetChallengePeriod() (*big.Int, error) {
	return _CrossChain.Contract.GetChallengePeriod(&_CrossChain.CallOpts)
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_CrossChain *CrossChainCallerSession) GetChallengePeriod() (*big.Int, error) {
	return _CrossChain.Contract.GetChallengePeriod(&_CrossChain.CallOpts)
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
// Solidity: function isBundleSaved(bytes32 ) view returns(bool)
func (_CrossChain *CrossChainCaller) IsBundleSaved(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "isBundleSaved", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBundleSaved is a free data retrieval call binding the contract method 0x84154826.
//
// Solidity: function isBundleSaved(bytes32 ) view returns(bool)
func (_CrossChain *CrossChainSession) IsBundleSaved(arg0 [32]byte) (bool, error) {
	return _CrossChain.Contract.IsBundleSaved(&_CrossChain.CallOpts, arg0)
}

// IsBundleSaved is a free data retrieval call binding the contract method 0x84154826.
//
// Solidity: function isBundleSaved(bytes32 ) view returns(bool)
func (_CrossChain *CrossChainCallerSession) IsBundleSaved(arg0 [32]byte) (bool, error) {
	return _CrossChain.Contract.IsBundleSaved(&_CrossChain.CallOpts, arg0)
}

// IsWithdrawalSpent is a free data retrieval call binding the contract method 0x2f0cb9e3.
//
// Solidity: function isWithdrawalSpent(bytes32 ) view returns(bool)
func (_CrossChain *CrossChainCaller) IsWithdrawalSpent(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _CrossChain.contract.Call(opts, &out, "isWithdrawalSpent", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWithdrawalSpent is a free data retrieval call binding the contract method 0x2f0cb9e3.
//
// Solidity: function isWithdrawalSpent(bytes32 ) view returns(bool)
func (_CrossChain *CrossChainSession) IsWithdrawalSpent(arg0 [32]byte) (bool, error) {
	return _CrossChain.Contract.IsWithdrawalSpent(&_CrossChain.CallOpts, arg0)
}

// IsWithdrawalSpent is a free data retrieval call binding the contract method 0x2f0cb9e3.
//
// Solidity: function isWithdrawalSpent(bytes32 ) view returns(bool)
func (_CrossChain *CrossChainCallerSession) IsWithdrawalSpent(arg0 [32]byte) (bool, error) {
	return _CrossChain.Contract.IsWithdrawalSpent(&_CrossChain.CallOpts, arg0)
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

// ExtractNativeValue is a paid mutator transaction binding the contract method 0x6af52662.
//
// Solidity: function extractNativeValue((address,address,uint256,uint64) _msg, bytes32[] proof, bytes32 root) returns()
func (_CrossChain *CrossChainTransactor) ExtractNativeValue(opts *bind.TransactOpts, _msg StructsValueTransferMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "extractNativeValue", _msg, proof, root)
}

// ExtractNativeValue is a paid mutator transaction binding the contract method 0x6af52662.
//
// Solidity: function extractNativeValue((address,address,uint256,uint64) _msg, bytes32[] proof, bytes32 root) returns()
func (_CrossChain *CrossChainSession) ExtractNativeValue(_msg StructsValueTransferMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChain.Contract.ExtractNativeValue(&_CrossChain.TransactOpts, _msg, proof, root)
}

// ExtractNativeValue is a paid mutator transaction binding the contract method 0x6af52662.
//
// Solidity: function extractNativeValue((address,address,uint256,uint64) _msg, bytes32[] proof, bytes32 root) returns()
func (_CrossChain *CrossChainTransactorSession) ExtractNativeValue(_msg StructsValueTransferMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChain.Contract.ExtractNativeValue(&_CrossChain.TransactOpts, _msg, proof, root)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _messageBus) returns()
func (_CrossChain *CrossChainTransactor) Initialize(opts *bind.TransactOpts, _messageBus common.Address) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "initialize", _messageBus)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _messageBus) returns()
func (_CrossChain *CrossChainSession) Initialize(_messageBus common.Address) (*types.Transaction, error) {
	return _CrossChain.Contract.Initialize(&_CrossChain.TransactOpts, _messageBus)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _messageBus) returns()
func (_CrossChain *CrossChainTransactorSession) Initialize(_messageBus common.Address) (*types.Transaction, error) {
	return _CrossChain.Contract.Initialize(&_CrossChain.TransactOpts, _messageBus)
}

// PauseWithdrawals is a paid mutator transaction binding the contract method 0xf4cc87ba.
//
// Solidity: function pauseWithdrawals(bool _pause) returns()
func (_CrossChain *CrossChainTransactor) PauseWithdrawals(opts *bind.TransactOpts, _pause bool) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "pauseWithdrawals", _pause)
}

// PauseWithdrawals is a paid mutator transaction binding the contract method 0xf4cc87ba.
//
// Solidity: function pauseWithdrawals(bool _pause) returns()
func (_CrossChain *CrossChainSession) PauseWithdrawals(_pause bool) (*types.Transaction, error) {
	return _CrossChain.Contract.PauseWithdrawals(&_CrossChain.TransactOpts, _pause)
}

// PauseWithdrawals is a paid mutator transaction binding the contract method 0xf4cc87ba.
//
// Solidity: function pauseWithdrawals(bool _pause) returns()
func (_CrossChain *CrossChainTransactorSession) PauseWithdrawals(_pause bool) (*types.Transaction, error) {
	return _CrossChain.Contract.PauseWithdrawals(&_CrossChain.TransactOpts, _pause)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CrossChain *CrossChainTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CrossChain *CrossChainSession) RenounceOwnership() (*types.Transaction, error) {
	return _CrossChain.Contract.RenounceOwnership(&_CrossChain.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CrossChain *CrossChainTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CrossChain.Contract.RenounceOwnership(&_CrossChain.TransactOpts)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_CrossChain *CrossChainTransactor) SetChallengePeriod(opts *bind.TransactOpts, _delay *big.Int) (*types.Transaction, error) {
	return _CrossChain.contract.Transact(opts, "setChallengePeriod", _delay)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_CrossChain *CrossChainSession) SetChallengePeriod(_delay *big.Int) (*types.Transaction, error) {
	return _CrossChain.Contract.SetChallengePeriod(&_CrossChain.TransactOpts, _delay)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_CrossChain *CrossChainTransactorSession) SetChallengePeriod(_delay *big.Int) (*types.Transaction, error) {
	return _CrossChain.Contract.SetChallengePeriod(&_CrossChain.TransactOpts, _delay)
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

// CrossChainLogManagementContractCreatedIterator is returned from FilterLogManagementContractCreated and is used to iterate over the raw logs and unpacked data for LogManagementContractCreated events raised by the CrossChain contract.
type CrossChainLogManagementContractCreatedIterator struct {
	Event *CrossChainLogManagementContractCreated // Event containing the contract specifics and raw log

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
func (it *CrossChainLogManagementContractCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainLogManagementContractCreated)
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
		it.Event = new(CrossChainLogManagementContractCreated)
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
func (it *CrossChainLogManagementContractCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainLogManagementContractCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainLogManagementContractCreated represents a LogManagementContractCreated event raised by the CrossChain contract.
type CrossChainLogManagementContractCreated struct {
	MessageBusAddress common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterLogManagementContractCreated is a free log retrieval operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_CrossChain *CrossChainFilterer) FilterLogManagementContractCreated(opts *bind.FilterOpts) (*CrossChainLogManagementContractCreatedIterator, error) {

	logs, sub, err := _CrossChain.contract.FilterLogs(opts, "LogManagementContractCreated")
	if err != nil {
		return nil, err
	}
	return &CrossChainLogManagementContractCreatedIterator{contract: _CrossChain.contract, event: "LogManagementContractCreated", logs: logs, sub: sub}, nil
}

// WatchLogManagementContractCreated is a free log subscription operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_CrossChain *CrossChainFilterer) WatchLogManagementContractCreated(opts *bind.WatchOpts, sink chan<- *CrossChainLogManagementContractCreated) (event.Subscription, error) {

	logs, sub, err := _CrossChain.contract.WatchLogs(opts, "LogManagementContractCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainLogManagementContractCreated)
				if err := _CrossChain.contract.UnpackLog(event, "LogManagementContractCreated", log); err != nil {
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

// ParseLogManagementContractCreated is a log parse operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_CrossChain *CrossChainFilterer) ParseLogManagementContractCreated(log types.Log) (*CrossChainLogManagementContractCreated, error) {
	event := new(CrossChainLogManagementContractCreated)
	if err := _CrossChain.contract.UnpackLog(event, "LogManagementContractCreated", log); err != nil {
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

// CrossChainWithdrawalsPausedIterator is returned from FilterWithdrawalsPaused and is used to iterate over the raw logs and unpacked data for WithdrawalsPaused events raised by the CrossChain contract.
type CrossChainWithdrawalsPausedIterator struct {
	Event *CrossChainWithdrawalsPaused // Event containing the contract specifics and raw log

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
func (it *CrossChainWithdrawalsPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainWithdrawalsPaused)
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
		it.Event = new(CrossChainWithdrawalsPaused)
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
func (it *CrossChainWithdrawalsPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainWithdrawalsPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainWithdrawalsPaused represents a WithdrawalsPaused event raised by the CrossChain contract.
type CrossChainWithdrawalsPaused struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalsPaused is a free log retrieval operation binding the contract event 0x129d33f7856617012aed60524381cfff7233cfc57df58d9f6613a5593d3dc218.
//
// Solidity: event WithdrawalsPaused(bool paused)
func (_CrossChain *CrossChainFilterer) FilterWithdrawalsPaused(opts *bind.FilterOpts) (*CrossChainWithdrawalsPausedIterator, error) {

	logs, sub, err := _CrossChain.contract.FilterLogs(opts, "WithdrawalsPaused")
	if err != nil {
		return nil, err
	}
	return &CrossChainWithdrawalsPausedIterator{contract: _CrossChain.contract, event: "WithdrawalsPaused", logs: logs, sub: sub}, nil
}

// WatchWithdrawalsPaused is a free log subscription operation binding the contract event 0x129d33f7856617012aed60524381cfff7233cfc57df58d9f6613a5593d3dc218.
//
// Solidity: event WithdrawalsPaused(bool paused)
func (_CrossChain *CrossChainFilterer) WatchWithdrawalsPaused(opts *bind.WatchOpts, sink chan<- *CrossChainWithdrawalsPaused) (event.Subscription, error) {

	logs, sub, err := _CrossChain.contract.WatchLogs(opts, "WithdrawalsPaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainWithdrawalsPaused)
				if err := _CrossChain.contract.UnpackLog(event, "WithdrawalsPaused", log); err != nil {
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

// ParseWithdrawalsPaused is a log parse operation binding the contract event 0x129d33f7856617012aed60524381cfff7233cfc57df58d9f6613a5593d3dc218.
//
// Solidity: event WithdrawalsPaused(bool paused)
func (_CrossChain *CrossChainFilterer) ParseWithdrawalsPaused(log types.Log) (*CrossChainWithdrawalsPaused, error) {
	event := new(CrossChainWithdrawalsPaused)
	if err := _CrossChain.contract.UnpackLog(event, "WithdrawalsPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
