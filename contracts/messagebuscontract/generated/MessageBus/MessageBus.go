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
	Topic    []byte
	Payload  []byte
}

// MessageBusMetaData contains all meta data concerning the MessageBus contract.
var MessageBusMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"topic\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"topic\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"topic\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"topic\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"submitOutOfNetworkMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"topic\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611597806100206000396000f3fe6080604052600436106100435760003560e01c80633bb454e5146100be5780639d8acc5b146100fb578063a5804d851461012b578063f2a64f5a1461015457610083565b36610083576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161007a906104ff565b60405180910390fd5b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100b59061056b565b60405180910390fd5b3480156100ca57600080fd5b506100e560048036038101906100e091906105c3565b610191565b6040516100f29190610627565b60405180910390f35b610115600480360381019061011091906107fd565b6101db565b60405161012291906108bf565b60405180910390f35b34801561013757600080fd5b50610152600480360381019061014d9190610910565b610253565b005b34801561016057600080fd5b5061017b600480360381019061017691906105c3565b6103e7565b604051610188919061097b565b60405180910390f35b600080826040516020016101a59190610beb565b60405160208183030381529060405280519060200120905042600080838152602001908152602001600020541015915050919050565b60006101e5610477565b3410156101f157600080fd5b600190503373ffffffffffffffffffffffffffffffffffffffff167fb8f8e2e252f184a45933e7d1c6d419d1408270cb65bfce9c7ec5f0f425af500f8287878787604051610243959493929190610caa565b60405180910390a2949350505050565b60008111610296576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161028d90610d57565b60405180910390fd5b6000826040516020016102a99190610beb565b60405160208183030381529060405280519060200120905060008060008381526020019081526020016000205414610316576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161030d90610de9565b60405180910390fd5b8160008083815260200190815260200160002081905550600160008460000160208101906103449190610e09565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002083806060019061038e9190610e45565b60405161039c929190610ed8565b908152602001604051809103902083908060018154018082558091505060019003906000526020600020906003020160009091909190915081816103e091906114c1565b5050505050565b600080826040516020016103fb9190610beb565b60405160208183030381529060405280519060200120905060008060008381526020019081526020016000205490506000811161046d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161046490611541565b60405180910390fd5b8092505050919050565b600090565b600082825260208201905092915050565b7f74686520576f726d686f6c6520636f6e747261637420646f6573206e6f74206160008201527f6363657074206173736574730000000000000000000000000000000000000000602082015250565b60006104e9602c8361047c565b91506104f48261048d565b604082019050919050565b60006020820190508181036000830152610518816104dc565b9050919050565b7f756e737570706f72746564000000000000000000000000000000000000000000600082015250565b6000610555600b8361047c565b91506105608261051f565b602082019050919050565b6000602082019050818103600083015261058481610548565b9050919050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600060a082840312156105ba576105b961059f565b5b81905092915050565b6000602082840312156105d9576105d8610595565b5b600082013567ffffffffffffffff8111156105f7576105f661059a565b5b610603848285016105a4565b91505092915050565b60008115159050919050565b6106218161060c565b82525050565b600060208201905061063c6000830184610618565b92915050565b600063ffffffff82169050919050565b61065b81610642565b811461066657600080fd5b50565b60008135905061067881610652565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6106d182610688565b810181811067ffffffffffffffff821117156106f0576106ef610699565b5b80604052505050565b600061070361058b565b905061070f82826106c8565b919050565b600067ffffffffffffffff82111561072f5761072e610699565b5b61073882610688565b9050602081019050919050565b82818337600083830152505050565b600061076761076284610714565b6106f9565b90508281526020810184848401111561078357610782610683565b5b61078e848285610745565b509392505050565b600082601f8301126107ab576107aa61067e565b5b81356107bb848260208601610754565b91505092915050565b600060ff82169050919050565b6107da816107c4565b81146107e557600080fd5b50565b6000813590506107f7816107d1565b92915050565b6000806000806080858703121561081757610816610595565b5b600061082587828801610669565b945050602085013567ffffffffffffffff8111156108465761084561059a565b5b61085287828801610796565b935050604085013567ffffffffffffffff8111156108735761087261059a565b5b61087f87828801610796565b9250506060610890878288016107e8565b91505092959194509250565b600067ffffffffffffffff82169050919050565b6108b98161089c565b82525050565b60006020820190506108d460008301846108b0565b92915050565b6000819050919050565b6108ed816108da565b81146108f857600080fd5b50565b60008135905061090a816108e4565b92915050565b6000806040838503121561092757610926610595565b5b600083013567ffffffffffffffff8111156109455761094461059a565b5b610951858286016105a4565b9250506020610962858286016108fb565b9150509250929050565b610975816108da565b82525050565b6000602082019050610990600083018461096c565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006109c182610996565b9050919050565b6109d1816109b6565b81146109dc57600080fd5b50565b6000813590506109ee816109c8565b92915050565b6000610a0360208401846109df565b905092915050565b610a14816109b6565b82525050565b610a238161089c565b8114610a2e57600080fd5b50565b600081359050610a4081610a1a565b92915050565b6000610a556020840184610a31565b905092915050565b610a668161089c565b82525050565b6000610a7b6020840184610669565b905092915050565b610a8c81610642565b82525050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610abe57610abd610a9c565b5b83810192508235915060208301925067ffffffffffffffff821115610ae657610ae5610a92565b5b600182023603831315610afc57610afb610a97565b5b509250929050565b600082825260208201905092915050565b6000610b218385610b04565b9350610b2e838584610745565b610b3783610688565b840190509392505050565b600060a08301610b5560008401846109f4565b610b626000860182610a0b565b50610b706020840184610a46565b610b7d6020860182610a5d565b50610b8b6040840184610a6c565b610b986040860182610a83565b50610ba66060840184610aa1565b8583036060870152610bb9838284610b15565b92505050610bca6080840184610aa1565b8583036080870152610bdd838284610b15565b925050508091505092915050565b60006020820190508181036000830152610c058184610b42565b905092915050565b610c1681610642565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610c56578082015181840152602081019050610c3b565b60008484015250505050565b6000610c6d82610c1c565b610c778185610c27565b9350610c87818560208601610c38565b610c9081610688565b840191505092915050565b610ca4816107c4565b82525050565b600060a082019050610cbf60008301886108b0565b610ccc6020830187610c0d565b8181036040830152610cde8186610c62565b90508181036060830152610cf28185610c62565b9050610d016080830184610c9b565b9695505050505050565b7f4e6f2e0000000000000000000000000000000000000000000000000000000000600082015250565b6000610d4160038361047c565b9150610d4c82610d0b565b602082019050919050565b60006020820190508181036000830152610d7081610d34565b9050919050565b7f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636560008201527f2100000000000000000000000000000000000000000000000000000000000000602082015250565b6000610dd360218361047c565b9150610dde82610d77565b604082019050919050565b60006020820190508181036000830152610e0281610dc6565b9050919050565b600060208284031215610e1f57610e1e610595565b5b6000610e2d848285016109df565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610e6257610e61610e36565b5b80840192508235915067ffffffffffffffff821115610e8457610e83610e3b565b5b602083019250600182023603831315610ea057610e9f610e40565b5b509250929050565b600081905092915050565b6000610ebf8385610ea8565b9350610ecc838584610745565b82840190509392505050565b6000610ee5828486610eb3565b91508190509392505050565b60008135610efe816109c8565b80915050919050565b60008160001b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff610f3484610f07565b9350801983169250808416831791505092915050565b6000819050919050565b6000610f6f610f6a610f6584610996565b610f4a565b610996565b9050919050565b6000610f8182610f54565b9050919050565b6000610f9382610f76565b9050919050565b6000819050919050565b610fad82610f88565b610fc0610fb982610f9a565b8354610f14565b8255505050565b60008135610fd481610a1a565b80915050919050565b60008160a01b9050919050565b60007bffffffffffffffff000000000000000000000000000000000000000061101284610fdd565b9350801983169250808416831791505092915050565b600061104361103e6110398461089c565b610f4a565b61089c565b9050919050565b6000819050919050565b61105d82611028565b6110706110698261104a565b8354610fea565b8255505050565b6000813561108481610652565b80915050919050565b60008160e01b9050919050565b60007fffffffff000000000000000000000000000000000000000000000000000000006110c68461108d565b9350801983169250808416831791505092915050565b60006110f76110f26110ed84610642565b610f4a565b610642565b9050919050565b6000819050919050565b611111826110dc565b61112461111d826110fe565b835461109a565b8255505050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061117d57607f821691505b6020821081036111905761118f611136565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026111f87fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826111bb565b61120286836111bb565b95508019841693508086168417925050509392505050565b600061123561123061122b846108da565b610f4a565b6108da565b9050919050565b6000819050919050565b61124f8361121a565b61126361125b8261123c565b8484546111c8565b825550505050565b600090565b61127861126b565b611283818484611246565b505050565b5b818110156112a75761129c600082611270565b600181019050611289565b5050565b601f8211156112ec576112bd81611196565b6112c6846111ab565b810160208510156112d5578190505b6112e96112e1856111ab565b830182611288565b50505b505050565b600082821c905092915050565b600061130f600019846008026112f1565b1980831691505092915050565b600061132883836112fe565b9150826002028217905092915050565b611342838361112b565b67ffffffffffffffff81111561135b5761135a610699565b5b6113658254611165565b6113708282856112ab565b6000601f83116001811461139f576000841561138d578287013590505b611397858261131c565b8655506113ff565b601f1984166113ad86611196565b60005b828110156113d5578489013582556001820191506020850194506020810190506113b0565b868310156113f257848901356113ee601f8916826112fe565b8355505b6001600288020188555050505b50505050505050565b611413838383611338565b505050565b60008101600083018061142a81610ef1565b90506114368184610fa4565b50505060008101602083018061144b81610fc7565b90506114578184611054565b50505060008101604083018061146c81611077565b90506114788184611108565b505050600181016060830161148d8185610e45565b611498818386611408565b5050505060028101608083016114ae8185610e45565b6114b9818386611408565b505050505050565b6114cb8282611418565b5050565b7f54686973206d65737361676520776173206e65766572207375626d697474656460008201527f2e00000000000000000000000000000000000000000000000000000000000000602082015250565b600061152b60218361047c565b9150611536826114cf565b604082019050919050565b6000602082019050818103600083015261155a8161151e565b905091905056fea2646970667358221220bfa69f7069f47e3b0ed1bd376edfb2e18d8a6ffb519ae2c255db32d07ee3df7c64736f6c63430008110033",
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

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0xf2a64f5a.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,bytes,bytes) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0xf2a64f5a.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,bytes,bytes) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0xf2a64f5a.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,bytes,bytes) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x3bb454e5.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,bytes,bytes) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x3bb454e5.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,bytes,bytes) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x3bb454e5.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,bytes,bytes) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x9d8acc5b.
//
// Solidity: function publishMessage(uint32 nonce, bytes topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint32, topic []byte, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x9d8acc5b.
//
// Solidity: function publishMessage(uint32 nonce, bytes topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusSession) PublishMessage(nonce uint32, topic []byte, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x9d8acc5b.
//
// Solidity: function publishMessage(uint32 nonce, bytes topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactorSession) PublishMessage(nonce uint32, topic []byte, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0xa5804d85.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,bytes,bytes) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactor) SubmitOutOfNetworkMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "submitOutOfNetworkMessage", crossChainMessage, finalAfterTimestamp)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0xa5804d85.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,bytes,bytes) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusSession) SubmitOutOfNetworkMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.SubmitOutOfNetworkMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// SubmitOutOfNetworkMessage is a paid mutator transaction binding the contract method 0xa5804d85.
//
// Solidity: function submitOutOfNetworkMessage((address,uint64,uint32,bytes,bytes) crossChainMessage, uint256 finalAfterTimestamp) returns()
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
	Topic            []byte
	Payload          []byte
	ConsistencyLevel uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogMessagePublished is a free log retrieval operation binding the contract event 0xb8f8e2e252f184a45933e7d1c6d419d1408270cb65bfce9c7ec5f0f425af500f.
//
// Solidity: event LogMessagePublished(address indexed sender, uint64 sequence, uint32 nonce, bytes topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts, sender []common.Address) (*MessageBusLogMessagePublishedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "LogMessagePublished", senderRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusLogMessagePublishedIterator{contract: _MessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xb8f8e2e252f184a45933e7d1c6d419d1408270cb65bfce9c7ec5f0f425af500f.
//
// Solidity: event LogMessagePublished(address indexed sender, uint64 sequence, uint32 nonce, bytes topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) WatchLogMessagePublished(opts *bind.WatchOpts, sink chan<- *MessageBusLogMessagePublished, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "LogMessagePublished", senderRule)
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

// ParseLogMessagePublished is a log parse operation binding the contract event 0xb8f8e2e252f184a45933e7d1c6d419d1408270cb65bfce9c7ec5f0f425af500f.
//
// Solidity: event LogMessagePublished(address indexed sender, uint64 sequence, uint32 nonce, bytes topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MessageBusLogMessagePublished, error) {
	event := new(MessageBusLogMessagePublished)
	if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
