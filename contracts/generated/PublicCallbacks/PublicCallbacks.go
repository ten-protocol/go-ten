// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicCallbacks

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

// PublicCallbacksMetaData contains all meta data concerning the PublicCallbacks contract.
var PublicCallbacksMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBefore\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfter\",\"type\":\"uint256\"}],\"name\":\"CallbackExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"CallbackRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"callbackBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executeNextCallbacks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"reattemptCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"callback\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"removeCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5060156019565b60c9565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161560685760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460c65780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b611435806100d65f395ff3fe60806040526004361061006d575f3560e01c8063929d34e91161004c578063929d34e9146100e05780639fb56ad2146100ff578063a67e17601461011e578063d98c616914610132575f5ffd5b8062e0d3b5146100715780638129fc1c146100aa57806382fbdc9c146100c0575b5f5ffd5b34801561007c575f5ffd5b5061009061008b366004610d9e565b61015d565b6040516100a1959493929190610e1c565b60405180910390f35b3480156100b5575f5ffd5b506100be610226565b005b6100d36100ce366004610ebd565b61035c565b6040516100a19190610f02565b3480156100eb575f5ffd5b506100be6100fa366004610d9e565b6103c3565b34801561010a575f5ffd5b506100be610119366004610d9e565b610605565b348015610129575f5ffd5b506100be61076f565b34801561013d575f5ffd5b506100d361014c366004610d9e565b60036020525f908152604090205481565b5f60208190529081526040902080546001820180546001600160a01b03909216929161018890610f24565b80601f01602080910402602001604051908101604052809291908181526020018280546101b490610f24565b80156101ff5780601f106101d6576101008083540402835291602001916101ff565b820191905f5260205f20905b8154815290600101906020018083116101e257829003601f168201915b5050505060028301546003840154600490940154929390929091506001600160a01b031685565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156102705750825b90505f8267ffffffffffffffff16600114801561028c5750303b155b90508115801561029a575080155b156102d1576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561030557845468ff00000000000000001916680100000000000000001785555b831561035557845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061034c90600190610f73565b60405180910390a15b5050505050565b5f5f34116103855760405162461bcd60e51b815260040161037c90610fb5565b60405180910390fd5b615208610391346107c0565b116103ae5760405162461bcd60e51b815260040161037c9061101f565b6103ba338484346107cb565b90505b92915050565b5f81815260036020526040902054819043116103f15760405162461bcd60e51b815260040161037c90611087565b5f82815260208181526040808320815160a0810190925280546001600160a01b03168252600181018054929391929184019161042c90610f24565b80601f016020809104026020016040519081016040528092919081815260200182805461045890610f24565b80156104a35780601f1061047a576101008083540402835291602001916104a3565b820191905f5260205f20905b81548152906001019060200180831161048657829003601f168201915b505050918352505060028201546020820152600382015460408201526004909101546001600160a01b039081166060909201919091528151919250166104fb5760405162461bcd60e51b815260040161037c906110c9565b60808101516001600160a01b031633146105275760405162461bcd60e51b815260040161037c9061110b565b805160208201516040515f926001600160a01b0316916105469161113c565b5f604051808303815f865af19150503d805f811461057f576040519150601f19603f3d011682016040523d82523d5f602084013e610584565b606091505b50509050806105a55760405162461bcd60e51b815260040161037c90611178565b5f84815260208190526040812080546001600160a01b0319168155906105ce6001830182610d4d565b505f600282018190556003808301829055600490920180546001600160a01b03191690559485526020525050604082209190915550565b5f81815260208181526040808320815160a0810190925280546001600160a01b03168252600181018054929391929184019161064090610f24565b80601f016020809104026020016040519081016040528092919081815260200182805461066c90610f24565b80156106b75780601f1061068e576101008083540402835291602001916106b7565b820191905f5260205f20905b81548152906001019060200180831161069a57829003601f168201915b505050918352505060028201546020820152600382015460408201526004909101546001600160a01b0390811660609092019190915260808201519192501633146107145760405162461bcd60e51b815260040161037c9061110b565b5f82815260208190526040812080546001600160a01b03191681559061073d6001830182610d4d565b505f600282018190556003808301829055600490920180546001600160a01b0319169055928352602052506040812055565b5f61077b60013061119c565b9050336001600160a01b038216146107a55760405162461bcd60e51b815260040161037c906111f1565b600254600154146107bd576107b861091e565b6107a5565b50565b5f6103bd4883611215565b5f60015490506040518060a00160405280866001600160a01b0316815260200185858080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92018290525093855250505060208201859052486040830152336060909201919091526001805482918261084983611228565b9091555081526020808201929092526040015f20825181546001600160a01b0319166001600160a01b0390911617815590820151600182019061088c90826112e9565b5060408281015160028301556060830151600380840191909155608090930151600490920180546001600160a01b0319166001600160a01b03909316929092179091555f83815260209290925290819020439055517f3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c49061090e908390610f02565b60405180910390a1949350505050565b6002546001540361092b57565b5f5f610935610a76565b915091505f826060015190505f8184604001516109529190611215565b90505f5a90505f855f01516001600160a01b0316838760200151604051610979919061113c565b5f604051808303815f8787f1925050503d805f81146109b3576040519150601f19603f3d011682016040523d82523d5f602084013e6109b8565b606091505b505090505f5a90505f6109cb82856113a5565b90505f818611156109ee57866109e183886113a5565b6109eb91906113b8565b90505b5f818a604001516109ff91906113a5565b8a519091508515610a1257610a12610baa565b610a1a610c0b565b610a2583828c610c21565b610a2e82610cfb565b7f79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b28a8887604051610a61939291906113cf565b60405180910390a15050505050505050505050565b610ab66040518060a001604052805f6001600160a01b03168152602001606081526020015f81526020015f81526020015f6001600160a01b031681525090565b6002545f81815260208181526040808320815160a0810190925280546001600160a01b0316825260018101805494959194919385929084019190610af990610f24565b80601f0160208091040260200160405190810160405280929190818152602001828054610b2590610f24565b8015610b705780601f10610b4757610100808354040283529160200191610b70565b820191905f5260205f20905b815481529060010190602001808311610b5357829003601f168201915b505050918352505060028201546020820152600382015460408201526004909101546001600160a01b031660609091015294909350915050565b6002545f90815260208190526040812080546001600160a01b031916815590610bd66001830182610d4d565b505f60028281018290556003808401839055600490930180546001600160a01b03191690555481526020919091526040812055565b60028054905f610c1a83611228565b9190505550565b5f81604051602401610c339190610f02565b60408051601f198184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f5ea3955800000000000000000000000000000000000000000000000000000000179052519091505f906001600160a01b0385169061afc8908790610cad90869061113c565b5f60405180830381858888f193505050503d805f8114610ce8576040519150601f19603f3d011682016040523d82523d5f602084013e610ced565b606091505b505090508061035557610355855b805f03610d055750565b604051419082905f81818185875af1925050503d805f8114610d42576040519150601f19603f3d011682016040523d82523d5f602084013e610d47565b606091505b50505050565b508054610d5990610f24565b5f825580601f10610d68575050565b601f0160209004905f5260205f20908101906107bd91905b80821115610d93575f8155600101610d80565b5090565b80356103bd565b5f60208284031215610db157610db15f5ffd5b6103ba8383610d97565b5f6001600160a01b0382166103bd565b610dd481610dbb565b82525050565b8281835e505f910152565b5f610dee825190565b808452602084019350610e05818560208601610dda565b601f01601f19169290920192915050565b80610dd4565b60a08101610e2a8288610dcb565b8181036020830152610e3c8187610de5565b9050610e4b6040830186610e16565b610e586060830185610e16565b610e656080830184610dcb565b9695505050505050565b5f5f83601f840112610e8257610e825f5ffd5b50813567ffffffffffffffff811115610e9c57610e9c5f5ffd5b602083019150836001820283011115610eb657610eb65f5ffd5b9250929050565b5f5f60208385031215610ed157610ed15f5ffd5b823567ffffffffffffffff811115610eea57610eea5f5ffd5b610ef685828601610e6f565b92509250509250929050565b602081016103bd8284610e16565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680610f3857607f821691505b602082108103610f4a57610f4a610f10565b50919050565b5f6103bd82610f5d565b90565b67ffffffffffffffff1690565b610dd481610f50565b602081016103bd8284610f6a565b600d8152602081017f4e6f2076616c75652073656e7400000000000000000000000000000000000000815290505b60200190565b602080825281016103bd81610f81565b60248152602081017f47617320746f6f206c6f7720636f6d706172656420746f20636f7374206f662081527f63616c6c00000000000000000000000000000000000000000000000000000000602082015290505b60400190565b602080825281016103bd81610fc5565b60228152602081017f43616c6c6261636b2063616e6e6f74206265207265617474656d70746564207981527f657400000000000000000000000000000000000000000000000000000000000060208201529050611019565b602080825281016103bd8161102f565b60178152602081017f43616c6c6261636b20646f6573206e6f7420657869737400000000000000000081529050610faf565b602080825281016103bd81611097565b60098152602081017f4e6f74206f776e6572000000000000000000000000000000000000000000000081529050610faf565b602080825281016103bd816110d9565b5f611124825190565b611132818560208601610dda565b9290920192915050565b6103bd818361111b565b60198152602081017f43616c6c6261636b20657865637574696f6e206661696c65640000000000000081529050610faf565b602080825281016103bd81611146565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b039182169190811690828203908111156103bd576103bd611188565b60088152602081017f4e6f742073656c6600000000000000000000000000000000000000000000000081529050610faf565b602080825281016103bd816111bf565b634e487b7160e01b5f52601260045260245ffd5b5f8261122357611223611201565b500490565b5f6001820161123957611239611188565b5060010190565b634e487b7160e01b5f52604160045260245ffd5b5f6103bd610f5a8381565b61126883611254565b81545f1960089490940293841b1916921b91909117905550565b5f61128e81848461125f565b505050565b818110156112ad576112a55f82611282565b600101611293565b5050565b601f82111561128e575f818152602090206020601f850104810160208510156112d75750805b6103556020601f860104830182611293565b815167ffffffffffffffff81111561130357611303611240565b61130d8254610f24565b6113188282856112b1565b506020601f82116001811461134b575f83156113345750848201515b5f19600885021c1981166002850217855550610355565b5f84815260208120601f198516915b8281101561137a578785015182556020948501946001909201910161135a565b508482101561139657838701515f19601f87166008021c191681555b50505050600202600101905550565b818103818111156103bd576103bd611188565b81810281158282048414176103bd576103bd611188565b606081016113dd8286610e16565b6113ea6020830185610e16565b6113f76040830184610e16565b94935050505056fea2646970667358221220f228c78b893ade19b32da2b94eabeb4529d97ba637dc0d851eedae73b50b0ad764736f6c634300081c0033",
}

// PublicCallbacksABI is the input ABI used to generate the binding from.
// Deprecated: Use PublicCallbacksMetaData.ABI instead.
var PublicCallbacksABI = PublicCallbacksMetaData.ABI

// PublicCallbacksBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PublicCallbacksMetaData.Bin instead.
var PublicCallbacksBin = PublicCallbacksMetaData.Bin

// DeployPublicCallbacks deploys a new Ethereum contract, binding an instance of PublicCallbacks to it.
func DeployPublicCallbacks(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PublicCallbacks, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PublicCallbacksBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// PublicCallbacks is an auto generated Go binding around an Ethereum contract.
type PublicCallbacks struct {
	PublicCallbacksCaller     // Read-only binding to the contract
	PublicCallbacksTransactor // Write-only binding to the contract
	PublicCallbacksFilterer   // Log filterer for contract events
}

// PublicCallbacksCaller is an auto generated read-only Go binding around an Ethereum contract.
type PublicCallbacksCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PublicCallbacksFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PublicCallbacksSession struct {
	Contract     *PublicCallbacks  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PublicCallbacksCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PublicCallbacksCallerSession struct {
	Contract *PublicCallbacksCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// PublicCallbacksTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PublicCallbacksTransactorSession struct {
	Contract     *PublicCallbacksTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// PublicCallbacksRaw is an auto generated low-level Go binding around an Ethereum contract.
type PublicCallbacksRaw struct {
	Contract *PublicCallbacks // Generic contract binding to access the raw methods on
}

// PublicCallbacksCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PublicCallbacksCallerRaw struct {
	Contract *PublicCallbacksCaller // Generic read-only contract binding to access the raw methods on
}

// PublicCallbacksTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactorRaw struct {
	Contract *PublicCallbacksTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicCallbacks creates a new instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacks(address common.Address, backend bind.ContractBackend) (*PublicCallbacks, error) {
	contract, err := bindPublicCallbacks(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// NewPublicCallbacksCaller creates a new read-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksCaller(address common.Address, caller bind.ContractCaller) (*PublicCallbacksCaller, error) {
	contract, err := bindPublicCallbacks(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCaller{contract: contract}, nil
}

// NewPublicCallbacksTransactor creates a new write-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicCallbacksTransactor, error) {
	contract, err := bindPublicCallbacks(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTransactor{contract: contract}, nil
}

// NewPublicCallbacksFilterer creates a new log filterer instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicCallbacksFilterer, error) {
	contract, err := bindPublicCallbacks(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksFilterer{contract: contract}, nil
}

// bindPublicCallbacks binds a generic wrapper to an already deployed contract.
func bindPublicCallbacks(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.PublicCallbacksCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transact(opts, method, params...)
}

// CallbackBlockNumber is a free data retrieval call binding the contract method 0xd98c6169.
//
// Solidity: function callbackBlockNumber(uint256 callbackId) view returns(uint256 blockNumber)
func (_PublicCallbacks *PublicCallbacksCaller) CallbackBlockNumber(opts *bind.CallOpts, callbackId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicCallbacks.contract.Call(opts, &out, "callbackBlockNumber", callbackId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CallbackBlockNumber is a free data retrieval call binding the contract method 0xd98c6169.
//
// Solidity: function callbackBlockNumber(uint256 callbackId) view returns(uint256 blockNumber)
func (_PublicCallbacks *PublicCallbacksSession) CallbackBlockNumber(callbackId *big.Int) (*big.Int, error) {
	return _PublicCallbacks.Contract.CallbackBlockNumber(&_PublicCallbacks.CallOpts, callbackId)
}

// CallbackBlockNumber is a free data retrieval call binding the contract method 0xd98c6169.
//
// Solidity: function callbackBlockNumber(uint256 callbackId) view returns(uint256 blockNumber)
func (_PublicCallbacks *PublicCallbacksCallerSession) CallbackBlockNumber(callbackId *big.Int) (*big.Int, error) {
	return _PublicCallbacks.Contract.CallbackBlockNumber(&_PublicCallbacks.CallOpts, callbackId)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 callbackId) view returns(address target, bytes data, uint256 value, uint256 baseFee, address owner)
func (_PublicCallbacks *PublicCallbacksCaller) Callbacks(opts *bind.CallOpts, callbackId *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
	Owner   common.Address
}, error) {
	var out []interface{}
	err := _PublicCallbacks.contract.Call(opts, &out, "callbacks", callbackId)

	outstruct := new(struct {
		Target  common.Address
		Data    []byte
		Value   *big.Int
		BaseFee *big.Int
		Owner   common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Data = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Value = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.BaseFee = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 callbackId) view returns(address target, bytes data, uint256 value, uint256 baseFee, address owner)
func (_PublicCallbacks *PublicCallbacksSession) Callbacks(callbackId *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
	Owner   common.Address
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, callbackId)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 callbackId) view returns(address target, bytes data, uint256 value, uint256 baseFee, address owner)
func (_PublicCallbacks *PublicCallbacksCallerSession) Callbacks(callbackId *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
	Owner   common.Address
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, callbackId)
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ExecuteNextCallbacks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "executeNextCallbacks")
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksSession) ExecuteNextCallbacks() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallbacks(&_PublicCallbacks.TransactOpts)
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ExecuteNextCallbacks() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallbacks(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ReattemptCallback(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "reattemptCallback", callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksTransactor) Register(opts *bind.TransactOpts, callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "register", callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksTransactorSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// RemoveCallback is a paid mutator transaction binding the contract method 0x9fb56ad2.
//
// Solidity: function removeCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactor) RemoveCallback(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "removeCallback", callbackId)
}

// RemoveCallback is a paid mutator transaction binding the contract method 0x9fb56ad2.
//
// Solidity: function removeCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksSession) RemoveCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.RemoveCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// RemoveCallback is a paid mutator transaction binding the contract method 0x9fb56ad2.
//
// Solidity: function removeCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) RemoveCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.RemoveCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// PublicCallbacksCallbackExecutedIterator is returned from FilterCallbackExecuted and is used to iterate over the raw logs and unpacked data for CallbackExecuted events raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecutedIterator struct {
	Event *PublicCallbacksCallbackExecuted // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksCallbackExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksCallbackExecuted)
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
		it.Event = new(PublicCallbacksCallbackExecuted)
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
func (it *PublicCallbacksCallbackExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksCallbackExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksCallbackExecuted represents a CallbackExecuted event raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecuted struct {
	CallbackId *big.Int
	GasBefore  *big.Int
	GasAfter   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCallbackExecuted is a free log retrieval operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterCallbackExecuted(opts *bind.FilterOpts) (*PublicCallbacksCallbackExecutedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCallbackExecutedIterator{contract: _PublicCallbacks.contract, event: "CallbackExecuted", logs: logs, sub: sub}, nil
}

// WatchCallbackExecuted is a free log subscription operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchCallbackExecuted(opts *bind.WatchOpts, sink chan<- *PublicCallbacksCallbackExecuted) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksCallbackExecuted)
				if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
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

// ParseCallbackExecuted is a log parse operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) ParseCallbackExecuted(log types.Log) (*PublicCallbacksCallbackExecuted, error) {
	event := new(PublicCallbacksCallbackExecuted)
	if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicCallbacksCallbackRegisteredIterator is returned from FilterCallbackRegistered and is used to iterate over the raw logs and unpacked data for CallbackRegistered events raised by the PublicCallbacks contract.
type PublicCallbacksCallbackRegisteredIterator struct {
	Event *PublicCallbacksCallbackRegistered // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksCallbackRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksCallbackRegistered)
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
		it.Event = new(PublicCallbacksCallbackRegistered)
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
func (it *PublicCallbacksCallbackRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksCallbackRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksCallbackRegistered represents a CallbackRegistered event raised by the PublicCallbacks contract.
type PublicCallbacksCallbackRegistered struct {
	CallbackId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCallbackRegistered is a free log retrieval operation binding the contract event 0x3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c4.
//
// Solidity: event CallbackRegistered(uint256 callbackId)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterCallbackRegistered(opts *bind.FilterOpts) (*PublicCallbacksCallbackRegisteredIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "CallbackRegistered")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCallbackRegisteredIterator{contract: _PublicCallbacks.contract, event: "CallbackRegistered", logs: logs, sub: sub}, nil
}

// WatchCallbackRegistered is a free log subscription operation binding the contract event 0x3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c4.
//
// Solidity: event CallbackRegistered(uint256 callbackId)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchCallbackRegistered(opts *bind.WatchOpts, sink chan<- *PublicCallbacksCallbackRegistered) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "CallbackRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksCallbackRegistered)
				if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackRegistered", log); err != nil {
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

// ParseCallbackRegistered is a log parse operation binding the contract event 0x3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c4.
//
// Solidity: event CallbackRegistered(uint256 callbackId)
func (_PublicCallbacks *PublicCallbacksFilterer) ParseCallbackRegistered(log types.Log) (*PublicCallbacksCallbackRegistered, error) {
	event := new(PublicCallbacksCallbackRegistered)
	if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicCallbacksInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PublicCallbacks contract.
type PublicCallbacksInitializedIterator struct {
	Event *PublicCallbacksInitialized // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksInitialized)
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
		it.Event = new(PublicCallbacksInitialized)
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
func (it *PublicCallbacksInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksInitialized represents a Initialized event raised by the PublicCallbacks contract.
type PublicCallbacksInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterInitialized(opts *bind.FilterOpts) (*PublicCallbacksInitializedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksInitializedIterator{contract: _PublicCallbacks.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PublicCallbacksInitialized) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksInitialized)
				if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_PublicCallbacks *PublicCallbacksFilterer) ParseInitialized(log types.Log) (*PublicCallbacksInitialized, error) {
	event := new(PublicCallbacksInitialized)
	if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
