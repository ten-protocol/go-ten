// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MessageBus

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
)

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender   common.Address
	Sequence uint64
	Nonce    uint32
	Topic    uint32
	Payload  []byte
}

// MessageBusMetaData contains all meta data concerning the MessageBus contract.
var MessageBusMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"submitOutOfNetworkMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b506114c4806100206000396000f3fe6080604052600436106100435760003560e01c80638d8f8cde146100be578063b1454caa146100e7578063b9d229ca14610124578063c5b219871461016157610083565b36610083576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161007a906104e4565b60405180910390fd5b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100b590610550565b60405180910390fd5b3480156100ca57600080fd5b506100e560048036038101906100e091906105d4565b61019e565b005b3480156100f357600080fd5b5061010e6004803603810190610109919061070a565b610335565b60405161011b91906107b5565b60405180910390f35b34801561013057600080fd5b5061014b600480360381019061014691906107d0565b610387565b6040516101589190610828565b60405180910390f35b34801561016d57600080fd5b50610188600480360381019061018391906107d0565b610417565b604051610195919061085e565b60405180910390f35b600081116101e1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101d8906108c5565b60405180910390fd5b6000826040516020016101f49190610b51565b60405160208183030381529060405280519060200120905060008060008381526020019081526020016000205414610261576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161025890610be5565b60405180910390fd5b81600080838152602001908152602001600020819055506001600084600001602081019061028f9190610c05565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008460600160208101906102de9190610c32565b63ffffffff1663ffffffff168152602001908152602001600020839080600181540180825580915050600190039060005260206000209060030201600090919091909150818161032e9190611319565b5050505050565b60007fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937338288888888886040516103729796959493929190611392565b60405180910390a16001905095945050505050565b6000808260405160200161039b9190610b51565b60405160208183030381529060405280519060200120905060008060008381526020019081526020016000205490506000811161040d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104049061146e565b60405180910390fd5b8092505050919050565b6000808260405160200161042b9190610b51565b60405160208183030381529060405280519060200120905042600080838152602001908152602001600020541015915050919050565b600082825260208201905092915050565b7f74686520576f726d686f6c6520636f6e747261637420646f6573206e6f74206160008201527f6363657074206173736574730000000000000000000000000000000000000000602082015250565b60006104ce602c83610461565b91506104d982610472565b604082019050919050565b600060208201905081810360008301526104fd816104c1565b9050919050565b7f756e737570706f72746564000000000000000000000000000000000000000000600082015250565b600061053a600b83610461565b915061054582610504565b602082019050919050565b600060208201905081810360008301526105698161052d565b9050919050565b600080fd5b600080fd5b600080fd5b600060a082840312156105955761059461057a565b5b81905092915050565b6000819050919050565b6105b18161059e565b81146105bc57600080fd5b50565b6000813590506105ce816105a8565b92915050565b600080604083850312156105eb576105ea610570565b5b600083013567ffffffffffffffff81111561060957610608610575565b5b6106158582860161057f565b9250506020610626858286016105bf565b9150509250929050565b600063ffffffff82169050919050565b61064981610630565b811461065457600080fd5b50565b60008135905061066681610640565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126106915761069061066c565b5b8235905067ffffffffffffffff8111156106ae576106ad610671565b5b6020830191508360018202830111156106ca576106c9610676565b5b9250929050565b600060ff82169050919050565b6106e7816106d1565b81146106f257600080fd5b50565b600081359050610704816106de565b92915050565b60008060008060006080868803121561072657610725610570565b5b600061073488828901610657565b955050602061074588828901610657565b945050604086013567ffffffffffffffff81111561076657610765610575565b5b6107728882890161067b565b93509350506060610785888289016106f5565b9150509295509295909350565b600067ffffffffffffffff82169050919050565b6107af81610792565b82525050565b60006020820190506107ca60008301846107a6565b92915050565b6000602082840312156107e6576107e5610570565b5b600082013567ffffffffffffffff81111561080457610803610575565b5b6108108482850161057f565b91505092915050565b6108228161059e565b82525050565b600060208201905061083d6000830184610819565b92915050565b60008115159050919050565b61085881610843565b82525050565b6000602082019050610873600083018461084f565b92915050565b7f4e6f2e0000000000000000000000000000000000000000000000000000000000600082015250565b60006108af600383610461565b91506108ba82610879565b602082019050919050565b600060208201905081810360008301526108de816108a2565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610910826108e5565b9050919050565b61092081610905565b811461092b57600080fd5b50565b60008135905061093d81610917565b92915050565b6000610952602084018461092e565b905092915050565b61096381610905565b82525050565b61097281610792565b811461097d57600080fd5b50565b60008135905061098f81610969565b92915050565b60006109a46020840184610980565b905092915050565b6109b581610792565b82525050565b60006109ca6020840184610657565b905092915050565b6109db81610630565b82525050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610a0d57610a0c6109eb565b5b83810192508235915060208301925067ffffffffffffffff821115610a3557610a346109e1565b5b600182023603831315610a4b57610a4a6109e6565b5b509250929050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b6000610a908385610a53565b9350610a9d838584610a64565b610aa683610a73565b840190509392505050565b600060a08301610ac46000840184610943565b610ad1600086018261095a565b50610adf6020840184610995565b610aec60208601826109ac565b50610afa60408401846109bb565b610b0760408601826109d2565b50610b1560608401846109bb565b610b2260608601826109d2565b50610b3060808401846109f0565b8583036080870152610b43838284610a84565b925050508091505092915050565b60006020820190508181036000830152610b6b8184610ab1565b905092915050565b7f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636560008201527f2100000000000000000000000000000000000000000000000000000000000000602082015250565b6000610bcf602183610461565b9150610bda82610b73565b604082019050919050565b60006020820190508181036000830152610bfe81610bc2565b9050919050565b600060208284031215610c1b57610c1a610570565b5b6000610c298482850161092e565b91505092915050565b600060208284031215610c4857610c47610570565b5b6000610c5684828501610657565b91505092915050565b60008135610c6c81610917565b80915050919050565b60008160001b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff610ca284610c75565b9350801983169250808416831791505092915050565b6000819050919050565b6000610cdd610cd8610cd3846108e5565b610cb8565b6108e5565b9050919050565b6000610cef82610cc2565b9050919050565b6000610d0182610ce4565b9050919050565b6000819050919050565b610d1b82610cf6565b610d2e610d2782610d08565b8354610c82565b8255505050565b60008135610d4281610969565b80915050919050565b60008160a01b9050919050565b60007bffffffffffffffff0000000000000000000000000000000000000000610d8084610d4b565b9350801983169250808416831791505092915050565b6000610db1610dac610da784610792565b610cb8565b610792565b9050919050565b6000819050919050565b610dcb82610d96565b610dde610dd782610db8565b8354610d58565b8255505050565b60008135610df281610640565b80915050919050565b60008160e01b9050919050565b60007fffffffff00000000000000000000000000000000000000000000000000000000610e3484610dfb565b9350801983169250808416831791505092915050565b6000610e65610e60610e5b84610630565b610cb8565b610630565b9050919050565b6000819050919050565b610e7f82610e4a565b610e92610e8b82610e6c565b8354610e08565b8255505050565b600063ffffffff610ea984610c75565b9350801983169250808416831791505092915050565b610ec882610e4a565b610edb610ed482610e6c565b8354610e99565b8255505050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610f0e57610f0d610ee2565b5b80840192508235915067ffffffffffffffff821115610f3057610f2f610ee7565b5b602083019250600182023603831315610f4c57610f4b610eec565b5b509250929050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610fd557607f821691505b602082108103610fe857610fe7610f8e565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026110507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611013565b61105a8683611013565b95508019841693508086168417925050509392505050565b600061108d6110886110838461059e565b610cb8565b61059e565b9050919050565b6000819050919050565b6110a783611072565b6110bb6110b382611094565b848454611020565b825550505050565b600090565b6110d06110c3565b6110db81848461109e565b505050565b5b818110156110ff576110f46000826110c8565b6001810190506110e1565b5050565b601f8211156111445761111581610fee565b61111e84611003565b8101602085101561112d578190505b61114161113985611003565b8301826110e0565b50505b505050565b600082821c905092915050565b600061116760001984600802611149565b1980831691505092915050565b60006111808383611156565b9150826002028217905092915050565b61119a8383610f54565b67ffffffffffffffff8111156111b3576111b2610f5f565b5b6111bd8254610fbd565b6111c8828285611103565b6000601f8311600181146111f757600084156111e5578287013590505b6111ef8582611174565b865550611257565b601f19841661120586610fee565b60005b8281101561122d57848901358255600182019150602085019450602081019050611208565b8683101561124a5784890135611246601f891682611156565b8355505b6001600288020188555050505b50505050505050565b61126b838383611190565b505050565b60008101600083018061128281610c5f565b905061128e8184610d12565b5050506000810160208301806112a381610d35565b90506112af8184610dc2565b5050506000810160408301806112c481610de5565b90506112d08184610e76565b5050506001810160608301806112e581610de5565b90506112f18184610ebf565b50505060028101608083016113068185610ef1565b611311818386611260565b505050505050565b6113238282611270565b5050565b61133081610905565b82525050565b61133f81610630565b82525050565b600082825260208201905092915050565b60006113628385611345565b935061136f838584610a64565b61137883610a73565b840190509392505050565b61138c816106d1565b82525050565b600060c0820190506113a7600083018a611327565b6113b460208301896107a6565b6113c16040830188611336565b6113ce6060830187611336565b81810360808301526113e1818587611356565b90506113f060a0830184611383565b98975050505050505050565b7f54686973206d65737361676520776173206e65766572207375626d697474656460008201527f2e00000000000000000000000000000000000000000000000000000000000000602082015250565b6000611458602183610461565b9150611463826113fc565b604082019050919050565b600060208201905081810360008301526114878161144b565b905091905056fea2646970667358221220ee5e898f9b3da4fc346cd8f137a60352825f4847d1f5b67f5c521f9bb7d7543364736f6c63430008110033",
}

// MessageBusABI is the input ABI used to generate the binding from.
// Deprecated: Use MessageBusMetaData.ABI instead.
var MessageBusABI = MessageBusMetaData.ABI

// MessageBusBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MessageBusMetaData.Bin instead.
var MessageBusBin = MessageBusMetaData.Bin

// DeployMessageBus deploys a new Ethereum contract, binding an instance of MessageBus to it.
func DeployMessageBus(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MessageBus, error) {
	parsed, err := MessageBusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageBusBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MessageBus{MessageBusCaller: MessageBusCaller{contract: contract}, MessageBusTransactor: MessageBusTransactor{contract: contract}, MessageBusFilterer: MessageBusFilterer{contract: contract}}, nil
}

// MessageBus is an auto generated Go binding around an Ethereum contract.
type MessageBus struct {
	MessageBusCaller     // Read-only binding to the contract
	MessageBusTransactor // Write-only binding to the contract
	MessageBusFilterer   // Log filterer for contract events
}

// MessageBusCaller is an auto generated read-only Go binding around an Ethereum contract.
type MessageBusCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MessageBusTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MessageBusFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MessageBusSession struct {
	Contract     *MessageBus       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MessageBusCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MessageBusCallerSession struct {
	Contract *MessageBusCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// MessageBusTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MessageBusTransactorSession struct {
	Contract     *MessageBusTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MessageBusRaw is an auto generated low-level Go binding around an Ethereum contract.
type MessageBusRaw struct {
	Contract *MessageBus // Generic contract binding to access the raw methods on
}

// MessageBusCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MessageBusCallerRaw struct {
	Contract *MessageBusCaller // Generic read-only contract binding to access the raw methods on
}

// MessageBusTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MessageBusTransactorRaw struct {
	Contract *MessageBusTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMessageBus creates a new instance of MessageBus, bound to a specific deployed contract.
func NewMessageBus(address common.Address, backend bind.ContractBackend) (*MessageBus, error) {
	contract, err := bindMessageBus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageBus{MessageBusCaller: MessageBusCaller{contract: contract}, MessageBusTransactor: MessageBusTransactor{contract: contract}, MessageBusFilterer: MessageBusFilterer{contract: contract}}, nil
}

// NewMessageBusCaller creates a new read-only instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusCaller(address common.Address, caller bind.ContractCaller) (*MessageBusCaller, error) {
	contract, err := bindMessageBus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageBusCaller{contract: contract}, nil
}

// NewMessageBusTransactor creates a new write-only instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageBusTransactor, error) {
	contract, err := bindMessageBus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageBusTransactor{contract: contract}, nil
}

// NewMessageBusFilterer creates a new log filterer instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageBusFilterer, error) {
	contract, err := bindMessageBus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageBusFilterer{contract: contract}, nil
}

// bindMessageBus binds a generic wrapper to an already deployed contract.
func bindMessageBus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MessageBusABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageBus *MessageBusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageBus.Contract.MessageBusCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageBus *MessageBusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.Contract.MessageBusTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageBus *MessageBusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageBus.Contract.MessageBusTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageBus *MessageBusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageBus.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageBus *MessageBusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageBus *MessageBusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageBus.Contract.contract.Transact(opts, method, params...)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0xb9d229ca.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0xb9d229ca.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0xb9d229ca.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0xc5b21987.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0xc5b21987.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0xc5b21987.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MessageBus *MessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MessageBus *MessageBusSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MessageBus *MessageBusTransactorSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0x8d8f8cde.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,uint32,bytes) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactor) SubmitOutOfNetworkMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "submitOutOfNetworkMessage", crossChainMessage, finalAfterTimestamp)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0x8d8f8cde.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,uint32,bytes) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusSession) SubmitOutOfNetworkMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.SubmitOutOfNetworkMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0x8d8f8cde.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,uint32,bytes) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactorSession) SubmitOutOfNetworkMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.SubmitOutOfNetworkMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MessageBus *MessageBusTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _MessageBus.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MessageBus *MessageBusSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MessageBus *MessageBusTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MessageBus *MessageBusTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MessageBus *MessageBusSession) Receive() (*types.Transaction, error) {
	return _MessageBus.Contract.Receive(&_MessageBus.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MessageBus *MessageBusTransactorSession) Receive() (*types.Transaction, error) {
	return _MessageBus.Contract.Receive(&_MessageBus.TransactOpts)
}

// MessageBusLogMessagePublishedIterator is returned from FilterLogMessagePublished and is used to iterate over the raw logs and unpacked data for LogMessagePublished events raised by the MessageBus contract.
type MessageBusLogMessagePublishedIterator struct {
	Event *MessageBusLogMessagePublished // Event containing the contract specifics and raw log

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
func (it *MessageBusLogMessagePublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusLogMessagePublished)
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
		it.Event = new(MessageBusLogMessagePublished)
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
func (it *MessageBusLogMessagePublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusLogMessagePublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusLogMessagePublished represents a LogMessagePublished event raised by the MessageBus contract.
type MessageBusLogMessagePublished struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogMessagePublished is a free log retrieval operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts) (*MessageBusLogMessagePublishedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return &MessageBusLogMessagePublishedIterator{contract: _MessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) WatchLogMessagePublished(opts *bind.WatchOpts, sink chan<- *MessageBusLogMessagePublished) (event.Subscription, error) {

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusLogMessagePublished)
				if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
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

// ParseLogMessagePublished is a log parse operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MessageBusLogMessagePublished, error) {
	event := new(MessageBusLogMessagePublished)
	if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
