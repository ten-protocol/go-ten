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
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// MessageBusMetaData contains all meta data concerning the MessageBus contract.
var MessageBusMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"submitOutOfNetworkMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061169c806100206000396000f3fe6080604052600436106100435760003560e01c80630fcfbd11146100be57806333a88c72146100fb578063b1454caa14610138578063f238ae121461017557610083565b36610083576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161007a9061059f565b60405180910390fd5b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100b59061060b565b60405180910390fd5b3480156100ca57600080fd5b506100e560048036038101906100e09190610659565b61019e565b6040516100f291906106bb565b60405180910390f35b34801561010757600080fd5b50610122600480360381019061011d9190610659565b61022e565b60405161012f91906106f1565b60405180910390f35b34801561014457600080fd5b5061015f600480360381019061015a91906107e6565b610278565b60405161016c9190610891565b60405180910390f35b34801561018157600080fd5b5061019c600480360381019061019791906108d8565b6102d1565b005b600080826040516020016101b29190610be1565b604051602081830303815290604052805190602001209050600080600083815260200190815260200160002054905060008111610224576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161021b90610c75565b60405180910390fd5b8092505050919050565b600080826040516020016102429190610be1565b60405160208183030381529060405280519060200120905042600080838152602001908152602001600020541015915050919050565b600061028333610436565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937338288888888886040516102c09796959493929190610d00565b60405180910390a195945050505050565b600081436102df9190610d99565b90506000836040516020016102f49190610be1565b60405160208183030381529060405280519060200120905060008060008381526020019081526020016000205414610361576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035890610e3f565b60405180910390fd5b81600080838152602001908152602001600020819055506001600085600001602081019061038f9190610e5f565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008560600160208101906103de9190610e8c565b63ffffffff1663ffffffff168152602001908152602001600020849080600181540180825580915050600190039060005260206000209060040201600090919091909150818161042e919061161c565b505050505050565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900467ffffffffffffffff1690506001600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a900467ffffffffffffffff166104f1919061162a565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b600082825260208201905092915050565b7f74686520576f726d686f6c6520636f6e747261637420646f6573206e6f74206160008201527f6363657074206173736574730000000000000000000000000000000000000000602082015250565b6000610589602c8361051c565b91506105948261052d565b604082019050919050565b600060208201905081810360008301526105b88161057c565b9050919050565b7f756e737570706f72746564000000000000000000000000000000000000000000600082015250565b60006105f5600b8361051c565b9150610600826105bf565b602082019050919050565b60006020820190508181036000830152610624816105e8565b9050919050565b600080fd5b600080fd5b600080fd5b600060c082840312156106505761064f610635565b5b81905092915050565b60006020828403121561066f5761066e61062b565b5b600082013567ffffffffffffffff81111561068d5761068c610630565b5b6106998482850161063a565b91505092915050565b6000819050919050565b6106b5816106a2565b82525050565b60006020820190506106d060008301846106ac565b92915050565b60008115159050919050565b6106eb816106d6565b82525050565b600060208201905061070660008301846106e2565b92915050565b600063ffffffff82169050919050565b6107258161070c565b811461073057600080fd5b50565b6000813590506107428161071c565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261076d5761076c610748565b5b8235905067ffffffffffffffff81111561078a5761078961074d565b5b6020830191508360018202830111156107a6576107a5610752565b5b9250929050565b600060ff82169050919050565b6107c3816107ad565b81146107ce57600080fd5b50565b6000813590506107e0816107ba565b92915050565b6000806000806000608086880312156108025761080161062b565b5b600061081088828901610733565b955050602061082188828901610733565b945050604086013567ffffffffffffffff81111561084257610841610630565b5b61084e88828901610757565b93509350506060610861888289016107d1565b9150509295509295909350565b600067ffffffffffffffff82169050919050565b61088b8161086e565b82525050565b60006020820190506108a66000830184610882565b92915050565b6108b5816106a2565b81146108c057600080fd5b50565b6000813590506108d2816108ac565b92915050565b600080604083850312156108ef576108ee61062b565b5b600083013567ffffffffffffffff81111561090d5761090c610630565b5b6109198582860161063a565b925050602061092a858286016108c3565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061095f82610934565b9050919050565b61096f81610954565b811461097a57600080fd5b50565b60008135905061098c81610966565b92915050565b60006109a1602084018461097d565b905092915050565b6109b281610954565b82525050565b6109c18161086e565b81146109cc57600080fd5b50565b6000813590506109de816109b8565b92915050565b60006109f360208401846109cf565b905092915050565b610a048161086e565b82525050565b6000610a196020840184610733565b905092915050565b610a2a8161070c565b82525050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610a5c57610a5b610a3a565b5b83810192508235915060208301925067ffffffffffffffff821115610a8457610a83610a30565b5b600182023603831315610a9a57610a99610a35565b5b509250929050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b6000610adf8385610aa2565b9350610aec838584610ab3565b610af583610ac2565b840190509392505050565b6000610b0f60208401846107d1565b905092915050565b610b20816107ad565b82525050565b600060c08301610b396000840184610992565b610b4660008601826109a9565b50610b5460208401846109e4565b610b6160208601826109fb565b50610b6f6040840184610a0a565b610b7c6040860182610a21565b50610b8a6060840184610a0a565b610b976060860182610a21565b50610ba56080840184610a3f565b8583036080870152610bb8838284610ad3565b92505050610bc960a0840184610b00565b610bd660a0860182610b17565b508091505092915050565b60006020820190508181036000830152610bfb8184610b26565b905092915050565b7f54686973206d65737361676520776173206e65766572207375626d697474656460008201527f2e00000000000000000000000000000000000000000000000000000000000000602082015250565b6000610c5f60218361051c565b9150610c6a82610c03565b604082019050919050565b60006020820190508181036000830152610c8e81610c52565b9050919050565b610c9e81610954565b82525050565b610cad8161070c565b82525050565b600082825260208201905092915050565b6000610cd08385610cb3565b9350610cdd838584610ab3565b610ce683610ac2565b840190509392505050565b610cfa816107ad565b82525050565b600060c082019050610d15600083018a610c95565b610d226020830189610882565b610d2f6040830188610ca4565b610d3c6060830187610ca4565b8181036080830152610d4f818587610cc4565b9050610d5e60a0830184610cf1565b98975050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610da4826106a2565b9150610daf836106a2565b9250828201905080821115610dc757610dc6610d6a565b5b92915050565b7f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636560008201527f2100000000000000000000000000000000000000000000000000000000000000602082015250565b6000610e2960218361051c565b9150610e3482610dcd565b604082019050919050565b60006020820190508181036000830152610e5881610e1c565b9050919050565b600060208284031215610e7557610e7461062b565b5b6000610e838482850161097d565b91505092915050565b600060208284031215610ea257610ea161062b565b5b6000610eb084828501610733565b91505092915050565b60008135610ec681610966565b80915050919050565b60008160001b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff610efc84610ecf565b9350801983169250808416831791505092915050565b6000819050919050565b6000610f37610f32610f2d84610934565b610f12565b610934565b9050919050565b6000610f4982610f1c565b9050919050565b6000610f5b82610f3e565b9050919050565b6000819050919050565b610f7582610f50565b610f88610f8182610f62565b8354610edc565b8255505050565b60008135610f9c816109b8565b80915050919050565b60008160a01b9050919050565b60007bffffffffffffffff0000000000000000000000000000000000000000610fda84610fa5565b9350801983169250808416831791505092915050565b600061100b6110066110018461086e565b610f12565b61086e565b9050919050565b6000819050919050565b61102582610ff0565b61103861103182611012565b8354610fb2565b8255505050565b6000813561104c8161071c565b80915050919050565b60008160e01b9050919050565b60007fffffffff0000000000000000000000000000000000000000000000000000000061108e84611055565b9350801983169250808416831791505092915050565b60006110bf6110ba6110b58461070c565b610f12565b61070c565b9050919050565b6000819050919050565b6110d9826110a4565b6110ec6110e5826110c6565b8354611062565b8255505050565b600063ffffffff61110384610ecf565b9350801983169250808416831791505092915050565b611122826110a4565b61113561112e826110c6565b83546110f3565b8255505050565b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126111685761116761113c565b5b80840192508235915067ffffffffffffffff82111561118a57611189611141565b5b6020830192506001820236038313156111a6576111a5611146565b5b509250929050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061122f57607f821691505b602082108103611242576112416111e8565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026112aa7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261126d565b6112b4868361126d565b95508019841693508086168417925050509392505050565b60006112e76112e26112dd846106a2565b610f12565b6106a2565b9050919050565b6000819050919050565b611301836112cc565b61131561130d826112ee565b84845461127a565b825550505050565b600090565b61132a61131d565b6113358184846112f8565b505050565b5b818110156113595761134e600082611322565b60018101905061133b565b5050565b601f82111561139e5761136f81611248565b6113788461125d565b81016020851015611387578190505b61139b6113938561125d565b83018261133a565b50505b505050565b600082821c905092915050565b60006113c1600019846008026113a3565b1980831691505092915050565b60006113da83836113b0565b9150826002028217905092915050565b6113f483836111ae565b67ffffffffffffffff81111561140d5761140c6111b9565b5b6114178254611217565b61142282828561135d565b6000601f831160018114611451576000841561143f578287013590505b61144985826113ce565b8655506114b1565b601f19841661145f86611248565b60005b8281101561148757848901358255600182019150602085019450602081019050611462565b868310156114a457848901356114a0601f8916826113b0565b8355505b6001600288020188555050505b50505050505050565b6114c58383836113ea565b505050565b600081356114d7816107ba565b80915050919050565b600060ff6114ed84610ecf565b9350801983169250808416831791505092915050565b600061151e611519611514846107ad565b610f12565b6107ad565b9050919050565b6000819050919050565b61153882611503565b61154b61154482611525565b83546114e0565b8255505050565b60008101600083018061156481610eb9565b90506115708184610f6c565b50505060008101602083018061158581610f8f565b9050611591818461101c565b5050506000810160408301806115a68161103f565b90506115b281846110d0565b5050506001810160608301806115c78161103f565b90506115d38184611119565b50505060028101608083016115e8818561114b565b6115f38183866114ba565b505050506003810160a0830180611609816114ca565b9050611615818461152f565b5050505050565b6116268282611552565b5050565b60006116358261086e565b91506116408361086e565b9250828201905067ffffffffffffffff8111156116605761165f610d6a565b5b9291505056fea2646970667358221220a74205e7b42de3eb7d1d2a8b82b7b1b63268644d6f357755b42237069d4f694364736f6c63430008110033",
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

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
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

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0xf238ae12.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactor) SubmitOutOfNetworkMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "submitOutOfNetworkMessage", crossChainMessage, finalAfterTimestamp)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0xf238ae12.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusSession) SubmitOutOfNetworkMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.SubmitOutOfNetworkMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0xf238ae12.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
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
