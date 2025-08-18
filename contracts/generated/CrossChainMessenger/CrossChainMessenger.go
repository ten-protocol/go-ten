// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package CrossChainMessenger

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

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint64
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// CrossChainMessengerMetaData contains all meta data concerning the CrossChainMessenger contract.
var CrossChainMessengerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"error\",\"type\":\"bytes\"}],\"name\":\"CallFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"crossChainSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainTarget\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"encodeCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messageBusAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContract\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"messageConsumed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"messageConsumed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"relayMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"relayMessageWithProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523460195760405161102061001e823961102090f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c80634c81bd20146100a0578063506762721461009b578063530c1e40146100965780635b76f28b1461009157806363012de51461008c578063772c655214610087578063a1a227fa14610082578063b859ce831461007d5763c4d66de8036100b357610477565b610448565b610422565b6103fb565b61037b565b61031f565b61020c565b6101a2565b6100e0565b908160c09103126100b35790565b5f80fd5b906020828203126100b357813567ffffffffffffffff81116100b3576100dd92016100a5565b90565b346100b3576100f86100f33660046100b7565b6106a0565b60405180805b0390f35b909182601f830112156100b35781359167ffffffffffffffff83116100b35760200192602083028401116100b357565b805b036100b357565b9050359061014882610132565b565b906060828203126100b357813567ffffffffffffffff81116100b357816101729184016100a5565b92602083013567ffffffffffffffff81116100b357826101996040946100dd938701610102565b9490950161013b565b346100b3576100f86101b536600461014a565b92919091610784565b906020828203126100b3576100dd9161013b565b6100dd916008021c5b60ff1690565b906100dd91546101d2565b5f6102036100dd926003905f5260205260405f2090565b6101e1565b9052565b346100b3576100fe6102276102223660046101be565b6101ec565b60405191829182901515815260200190565b6001600160a01b031690565b6001600160a01b038116610134565b9050359061014882610245565b909182601f830112156100b35781359167ffffffffffffffff83116100b35760200192600183028401116100b357565b9190916040818403126100b3576102a88382610254565b92602082013567ffffffffffffffff81116100b3576102c79201610261565b9091565b90825f9392825e0152565b6102f761030060209361030a936102eb815190565b80835293849260200190565b958691016102cb565b601f01601f191690565b0190565b60208082526100dd929101906102d6565b346100b3576100fe61033b610335366004610291565b91610802565b6040519182918261030e565b5f9103126100b357565b6100dd916008021c6001600160a01b031690565b906100dd9154610351565b6100dd5f6001610365565b346100b35761038b366004610347565b6100fe610396610370565b604051918291826001600160a01b03909116815260200190565b6100dd5f80610365565b6102396100dd6100dd926001600160a01b031690565b6100dd906103ba565b6100dd906103d0565b610208906103d9565b60208101929161014891906103e2565b346100b35761040b366004610347565b6100fe6104166103b0565b604051918291826103eb565b346100b357610432366004610347565b6100fe610396610881565b6100dd5f6002610365565b346100b357610458366004610347565b6100fe61039661043d565b906020828203126100b3576100dd91610254565b346100b3576100f861048a366004610463565b610ab3565b356100dd81610245565b906001600160a01b03905b9181191691161790565b906104be6100dd6104c5926103d9565b8254610499565b9055565b903590601e1936829003018212156100b3570180359067ffffffffffffffff82116100b357602001913682900383136100b357565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761054d57604052565b6104fe565b9061014861055f60405190565b928361052b565b67ffffffffffffffff811161054d57602090601f01601f19160190565b90825f939282370152565b909291926105a361059e82610566565b610552565b93818552818301116100b357610148916020850190610583565b9080601f830112156100b3578160206100dd9335910161058e565b9190916060818403126100b3576105ef6060610552565b926105fa8183610254565b8452602082013567ffffffffffffffff81116100b357826106228360409361062d96016105bd565b60208701520161013b565b6040830152565b906020828203126100b357813567ffffffffffffffff81116100b3576100dd92016105d8565b9061066761059e83610566565b918252565b3d156106855761067b3d61065a565b903d5f602084013e565b606090565b6102396100dd6100dd9290565b6100dd9061068a565b5f6106d96106d182936106b281610d5d565b6106c76106c085830161048f565b60016104ae565b60808101906104c9565b810190610634565b6107078282016106fa6106f382516001600160a01b031690565b60026104ae565b516001600160a01b031690565b908260205a9201519260208451940192f161072861072361066c565b911590565b61074457506101486107395f610697565b6106f38160016104ae565b6107809061075160405190565b9182917fa5fa8d2b0000000000000000000000000000000000000000000000000000000083526004830161030e565b0390fd5b6106d1906106b25f9586956106d99584610eeb565b6100dd6060610552565b6100dd91369161058e565b6100dd6100dd6100dd9290565b80516001600160a01b03168252906100dd906040806107e960608401602087015185820360208701526102d6565b940151910152565b60208082526100dd929101906107bb565b916108339061083a92610813606090565b5061082e61081f610799565b6001600160a01b039096168652565b6107a3565b6020830152565b61084661062d5f6107ae565b6100dd61085260405190565b80926108626020830191826107f1565b9081038252038261052b565b6100dd90610239565b6100dd905461086e565b6100dd61088d5f610877565b6103d9565b6100dd9060401c6101db565b6100dd9054610892565b6100dd905b67ffffffffffffffff1690565b6100dd90546108a8565b6108ad6100dd6100dd9290565b9067ffffffffffffffff906104a4565b6108ad6100dd6100dd9267ffffffffffffffff1690565b906109086100dd6104c5926108e1565b82546108d1565b9068ff00000000000000009060401b6104a4565b906109336100dd6104c592151590565b825461090f565b610208906108c4565b602081019291610148919061093a565b61095b610fb1565b8061097561096f61096b8361089e565b1590565b916108ba565b9261097f5f6108c4565b67ffffffffffffffff85161480610a9a575b6001946109ae6109a0876108c4565b9167ffffffffffffffff1690565b149081610a72575b155b9081610a69575b50610a3f576109e890826109df5f6109d6886108c4565b960195866108f8565b610a3057610aa1565b6109f0575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291610a1f5f610a2b93610923565b60405191829182610943565b0390a1565b610a3a8585610923565b610aa1565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f6109bf565b90506109b8610a80306103d9565b3b610a91610a8d5f6107ae565b9190565b149190506109b6565b5081610991565b610aad610148916103d9565b5f6104ae565b61014890610953565b801515610134565b9050519061014882610abc565b906020828203126100b3576100dd91610ac4565b506100dd906020810190610254565b67ffffffffffffffff8116610134565b9050359061014882610af4565b506100dd906020810190610b04565b63ffffffff8116610134565b9050359061014882610b20565b506100dd906020810190610b2c565b9035601e1936839003018112156100b357016020813591019167ffffffffffffffff82116100b3573682900383136100b357565b919061030081610b938161030a9560209181520190565b8095610583565b60ff8116610134565b9050359061014882610b9a565b506100dd906020810190610ba3565b906100dd9060a0610c6e610c6460c08401610bea610bdd8880610ae5565b6001600160a01b03168652565b610c0b610bfa6020890189610b11565b67ffffffffffffffff166020870152565b610c2c610c1b6040890189610b11565b67ffffffffffffffff166040870152565b610c49610c3c6060890189610b39565b63ffffffff166060870152565b610c566080880188610b48565b908683036080880152610b7c565b9482810190610bb0565b60ff16910152565b60208082526100dd92910190610bbf565b6040513d5f823e3d90fd5b15610c9957565b60405162461bcd60e51b815260206004820152601f60248201527f4d657373616765206e6f7420666f756e64206f722066696e616c697a65642e006044820152606490fd5b6100dd906101db565b6100dd9054610cde565b15610cf857565b60405162461bcd60e51b815260206004820152601960248201527f4d65737361676520616c726561647920636f6e73756d65642e000000000000006044820152606490fd5b9060ff906104a4565b90610d566100dd6104c592151590565b8254610d3d565b610d6f61088d61088d61088d5f610877565b6020610d7a60405190565b9182907f91643fdd0000000000000000000000000000000000000000000000000000000082528180610daf8760048301610c76565b03915afa8015610e6c5761014892610dd4610e38926001945f91610e3d575b50610c92565b610df0610de060405190565b8092610862602083019182610c76565b610e02610dfb825190565b9160200190565b20610e29610e2461096b610e1f846003905f5260205260405f2090565b610ce7565b610cf1565b6003905f5260205260405f2090565b610d46565b610e5f915060203d602011610e65575b610e57818361052b565b810190610ad1565b5f610dce565b503d610e4d565b610c87565b9037565b8183529091602001917f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81116100b3578291602061030a9202938491610e71565b949391610ee69061014894610ed860409460608a01908a82035f8c0152610bbf565b9188830360208a0152610e75565b940152565b92610ef861088d5f610877565b803b156100b357610f42935f93610f0e60405190565b958694859384937fce0d7db30000000000000000000000000000000000000000000000000000000085528a60048601610eb6565b03915afa8015610e6c5761014892600192610e3892610f94575b50610f69610de060405190565b610f74610dfb825190565b20610e29610f8e610e1f836003905f5260205260405f2090565b15610cf1565b610fab905f610fa3818361052b565b810190610347565b5f610f5c565b6100dd610fe2565b6100dd7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006107ae565b6100dd610fb956fea2646970667358221220cdbd1ae8e378c5c75825b338b69a5a9a9cf03d7a1580af7808d505520a9a747864736f6c634300081c0033",
}

// CrossChainMessengerABI is the input ABI used to generate the binding from.
// Deprecated: Use CrossChainMessengerMetaData.ABI instead.
var CrossChainMessengerABI = CrossChainMessengerMetaData.ABI

// CrossChainMessengerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrossChainMessengerMetaData.Bin instead.
var CrossChainMessengerBin = CrossChainMessengerMetaData.Bin

// DeployCrossChainMessenger deploys a new Ethereum contract, binding an instance of CrossChainMessenger to it.
func DeployCrossChainMessenger(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrossChainMessenger, error) {
	parsed, err := CrossChainMessengerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainMessengerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChainMessenger{CrossChainMessengerCaller: CrossChainMessengerCaller{contract: contract}, CrossChainMessengerTransactor: CrossChainMessengerTransactor{contract: contract}, CrossChainMessengerFilterer: CrossChainMessengerFilterer{contract: contract}}, nil
}

// CrossChainMessenger is an auto generated Go binding around an Ethereum contract.
type CrossChainMessenger struct {
	CrossChainMessengerCaller     // Read-only binding to the contract
	CrossChainMessengerTransactor // Write-only binding to the contract
	CrossChainMessengerFilterer   // Log filterer for contract events
}

// CrossChainMessengerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrossChainMessengerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrossChainMessengerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrossChainMessengerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrossChainMessengerSession struct {
	Contract     *CrossChainMessenger // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CrossChainMessengerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrossChainMessengerCallerSession struct {
	Contract *CrossChainMessengerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CrossChainMessengerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrossChainMessengerTransactorSession struct {
	Contract     *CrossChainMessengerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CrossChainMessengerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrossChainMessengerRaw struct {
	Contract *CrossChainMessenger // Generic contract binding to access the raw methods on
}

// CrossChainMessengerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrossChainMessengerCallerRaw struct {
	Contract *CrossChainMessengerCaller // Generic read-only contract binding to access the raw methods on
}

// CrossChainMessengerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrossChainMessengerTransactorRaw struct {
	Contract *CrossChainMessengerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrossChainMessenger creates a new instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessenger(address common.Address, backend bind.ContractBackend) (*CrossChainMessenger, error) {
	contract, err := bindCrossChainMessenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessenger{CrossChainMessengerCaller: CrossChainMessengerCaller{contract: contract}, CrossChainMessengerTransactor: CrossChainMessengerTransactor{contract: contract}, CrossChainMessengerFilterer: CrossChainMessengerFilterer{contract: contract}}, nil
}

// NewCrossChainMessengerCaller creates a new read-only instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerCaller(address common.Address, caller bind.ContractCaller) (*CrossChainMessengerCaller, error) {
	contract, err := bindCrossChainMessenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerCaller{contract: contract}, nil
}

// NewCrossChainMessengerTransactor creates a new write-only instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainMessengerTransactor, error) {
	contract, err := bindCrossChainMessenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerTransactor{contract: contract}, nil
}

// NewCrossChainMessengerFilterer creates a new log filterer instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainMessengerFilterer, error) {
	contract, err := bindCrossChainMessenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerFilterer{contract: contract}, nil
}

// bindCrossChainMessenger binds a generic wrapper to an already deployed contract.
func bindCrossChainMessenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrossChainMessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainMessenger *CrossChainMessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainMessenger.Contract.CrossChainMessengerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainMessenger *CrossChainMessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.CrossChainMessengerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainMessenger *CrossChainMessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.CrossChainMessengerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainMessenger *CrossChainMessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainMessenger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainMessenger *CrossChainMessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainMessenger *CrossChainMessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.contract.Transact(opts, method, params...)
}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) CrossChainSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "crossChainSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) CrossChainSender() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainSender(&_CrossChainMessenger.CallOpts)
}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) CrossChainSender() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainSender(&_CrossChainMessenger.CallOpts)
}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) CrossChainTarget(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "crossChainTarget")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) CrossChainTarget() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainTarget(&_CrossChainMessenger.CallOpts)
}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) CrossChainTarget() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainTarget(&_CrossChainMessenger.CallOpts)
}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerCaller) EncodeCall(opts *bind.CallOpts, target common.Address, payload []byte) ([]byte, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "encodeCall", target, payload)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerSession) EncodeCall(target common.Address, payload []byte) ([]byte, error) {
	return _CrossChainMessenger.Contract.EncodeCall(&_CrossChainMessenger.CallOpts, target, payload)
}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) EncodeCall(target common.Address, payload []byte) ([]byte, error) {
	return _CrossChainMessenger.Contract.EncodeCall(&_CrossChainMessenger.CallOpts, target, payload)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageBus() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBus(&_CrossChainMessenger.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageBus() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBus(&_CrossChainMessenger.CallOpts)
}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageBusContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageBusContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageBusContract() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBusContract(&_CrossChainMessenger.CallOpts)
}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageBusContract() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBusContract(&_CrossChainMessenger.CallOpts)
}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageConsumed(opts *bind.CallOpts, messageHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageConsumed", messageHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageConsumed(messageHash [32]byte) (bool, error) {
	return _CrossChainMessenger.Contract.MessageConsumed(&_CrossChainMessenger.CallOpts, messageHash)
}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageConsumed(messageHash [32]byte) (bool, error) {
	return _CrossChainMessenger.Contract.MessageConsumed(&_CrossChainMessenger.CallOpts, messageHash)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) Initialize(opts *bind.TransactOpts, messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "initialize", messageBusAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) Initialize(messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Initialize(&_CrossChainMessenger.TransactOpts, messageBusAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) Initialize(messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Initialize(&_CrossChainMessenger.TransactOpts, messageBusAddr)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RelayMessage(opts *bind.TransactOpts, message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "relayMessage", message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RelayMessageWithProof(opts *bind.TransactOpts, message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "relayMessageWithProof", message, proof, root)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RelayMessageWithProof(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessageWithProof(&_CrossChainMessenger.TransactOpts, message, proof, root)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RelayMessageWithProof(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessageWithProof(&_CrossChainMessenger.TransactOpts, message, proof, root)
}

// CrossChainMessengerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CrossChainMessenger contract.
type CrossChainMessengerInitializedIterator struct {
	Event *CrossChainMessengerInitialized // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerInitialized)
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
		it.Event = new(CrossChainMessengerInitialized)
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
func (it *CrossChainMessengerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerInitialized represents a Initialized event raised by the CrossChainMessenger contract.
type CrossChainMessengerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterInitialized(opts *bind.FilterOpts) (*CrossChainMessengerInitializedIterator, error) {

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerInitializedIterator{contract: _CrossChainMessenger.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerInitialized) (event.Subscription, error) {

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerInitialized)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParseInitialized(log types.Log) (*CrossChainMessengerInitialized, error) {
	event := new(CrossChainMessengerInitialized)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
