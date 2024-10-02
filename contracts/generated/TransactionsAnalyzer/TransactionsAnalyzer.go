// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TransactionsAnalyzer

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

// StructsTransaction is an auto generated low-level Go binding around an user-defined struct.
type StructsTransaction struct {
	TxType   uint8
	Nonce    *big.Int
	GasPrice *big.Int
	GasLimit *big.Int
	To       common.Address
	Value    *big.Int
	Data     []byte
	From     common.Address
}

// TransactionsAnalyzerMetaData contains all meta data concerning the TransactionsAnalyzer contract.
var TransactionsAnalyzerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transactionsLength\",\"type\":\"uint256\"}],\"name\":\"TransactionsConverted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EOA_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HOOK_CALLER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"}],\"name\":\"addOnBlockEndCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"eoaAdmin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"authorizedCaller\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523461001a57604051610f456100208239610f4590f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806301ffc9a7146100d2578063248a9ca3146100cd5780632f2ff15d146100c857806336568abe146100c3578063485cc955146100be5780634a44e6c1146100b9578063508a50f4146100b45780635f03a661146100af57806391d14854146100aa578063a217fddf146100a5578063d547741f146100a05763ee546fd803610102576103dd565b6103b0565b610395565b61035a565b610321565b6102e8565b6102c4565b610251565b61021e565b610200565b610189565b610131565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b0361010257565b600080fd5b90503590610114826100d7565b565b906020828203126101025761012a91610107565b90565b9052565b346101025761015e61014c610147366004610116565b6103f5565b60405191829182901515815260200190565b0390f35b806100fb565b9050359061011482610162565b906020828203126101025761012a91610168565b346101025761015e6101a461019f366004610175565b610518565b6040515b9182918290815260200190565b6001600160a01b031690565b6001600160a01b0381166100fb565b90503590610114826101c1565b91906040838203126101025761012a906101f78185610168565b936020016101d0565b34610102576102196102133660046101dd565b9061055a565b604051005b34610102576102196102313660046101dd565b9061060a565b91906040838203126101025761012a906101f781856101d0565b3461010257610219610264366004610237565b906109fb565b909182601f830112156101025781359167ffffffffffffffff831161010257602001926020830284011161010257565b9060208282031261010257813567ffffffffffffffff8111610102576102c0920161026a565b9091565b34610102576102196102d736600461029a565b90610f05565b600091031261010257565b34610102576102f83660046102dd565b61015e7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a5986101a4565b34610102576103313660046102dd565b61015e7f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db6101a4565b346101025761015e61014c6103703660046101dd565b906104e8565b61012a61012a61012a9290565b61012a6000610376565b61012a610383565b34610102576103a53660046102dd565b61015e6101a461038d565b34610102576102196103c33660046101dd565b90610600565b906020828203126101025761012a916101d0565b34610102576102196103f03660046103c9565b610ac4565b7f7965db0b000000000000000000000000000000000000000000000000000000007fffffffff00000000000000000000000000000000000000000000000000000000821614908115610445575090565b61012a91507fffffffff00000000000000000000000000000000000000000000000000000000167f01ffc9a7000000000000000000000000000000000000000000000000000000001490565b905b600052602052604060002090565b6101b561012a61012a926001600160a01b031690565b61012a906104a1565b61012a906104b7565b90610493906104c0565b61012a905b60ff1690565b61012a90546104d3565b61012a91610502916104fc60009182610491565b016104c9565b6104de565b61012a9081565b61012a9054610507565b600161053161012a92610529600090565b506000610491565b0161050e565b906101149161054d61054882610518565b610564565b9061055791610681565b50565b9061011491610537565b61011490339061058c565b6001600160a01b0390911681526040810192916101149160200152565b9061059e61059a82846104e8565b1590565b6105a6575050565b6105e16105b260405190565b9283927fe2517d3f0000000000000000000000000000000000000000000000000000000084526004840161056f565b0390fd5b90610114916105f661054882610518565b90610557916106fd565b90610114916105e5565b90610614336101b5565b6001600160a01b0382160361062c57610557916106fd565b6040517f6697b232000000000000000000000000000000000000000000000000000000008152600490fd5b9060ff905b9181191691161790565b9061067661012a61067d92151590565b8254610657565b9055565b61068e61059a83836104e8565b156106f6576001916106af836106aa8360006104fc8782610491565b610666565b33906106e56106df6106df7f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9590565b926104c0565b926106ef60405190565b600090a490565b5050600090565b9061070881836104e8565b156106f65761072160006106aa83826104fc8782610491565b33906107516106df6106df7ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9590565b9261075b60405190565b600090a4600190565b61012a9060401c6104d8565b61012a9054610764565b61012a905b67ffffffffffffffff1690565b61012a905461077a565b61077f61012a61012a9290565b9067ffffffffffffffff9061065c565b61077f61012a61012a9267ffffffffffffffff1690565b906107da61012a61067d926107b3565b82546107a3565b9068ff00000000000000009060401b61065c565b9061080561012a61067d92151590565b82546107e1565b61012d90610796565b602081019291610114919061080c565b907ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00908161085e61085861059a83610770565b9161078c565b9360009261086b84610796565b67ffffffffffffffff8716148061098d575b60019661089a61088c89610796565b9167ffffffffffffffff1690565b149081610965575b155b908161095c575b50610931576108d491836108cb866108c28a610796565b980197886107ca565b61092257610994565b6108dd57505050565b6108e6916107f5565b61091d7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29161091460405190565b91829182610815565b0390a1565b61092c87876107f5565b610994565b6040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152600490fd5b159050386108ab565b90506108a4610973306104c0565b3b61098461098088610376565b9190565b149190506108a2565b508261087d565b906109d4610557926109ad816109a8610383565b610681565b507ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a598610681565b507f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db610681565b9061011491610825565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b8054821015610a5457610a4b600191600052602060002090565b91020190600090565b610a1b565b9190600861065c910291610a736001600160a01b03841b90565b921b90565b9190610a8961012a61067d936104c0565b908354610a59565b9081549168010000000000000000831015610abf5782610ab991600161011495018155610a31565b90610a78565b610a05565b61011490610ad36001916104c0565b90610a91565b6101149190610b077f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db610564565b610d96565b0190565b634e487b7160e01b600052601160045260246000fd5b6000198114610b355760010190565b610b10565b61012a916008021c6101b5565b9061012a9154610b3a565b90601f01601f1916810190811067ffffffffffffffff821117610abf57604052565b60ff81166100fb565b9050359061011482610b74565b5061012a906020810190610b7d565b5061012a906020810190610168565b5061012a9060208101906101d0565b9035601e19368390030181121561010257016020813591019167ffffffffffffffff82116101025736829003831361010257565b90826000939282370152565b9190610c1581610c0e81610b0c9560209181520190565b8095610beb565b601f01601f191690565b9061012a9060e0610ce6610cdc6101008401610c45610c3e8880610b8a565b60ff168652565b610c5c610c556020890189610b99565b6020870152565b610c73610c6c6040890189610b99565b6040870152565b610c8a610c836060890189610b99565b6060870152565b610caa610c9a6080890189610ba8565b6001600160a01b03166080870152565b610cc1610cba60a0890189610b99565b60a0870152565b610cce60c0880188610bb7565b9086830360c0880152610bf7565b9482810190610ba8565b6001600160a01b0316910152565b9061012a91610c1f565b903560fe193683900301811215610102570190565b818352916020019081610d296020830284019490565b92836000925b848410610d3f5750505050505090565b9091929394956020610d6b610d648385600195038852610d5f8b88610cfe565b610cf4565b9860200190565b940194019294939190610d2f565b602080825261012a93910191610d13565b6040513d6000823e3d90fd5b90919060009083610da683610376565b8114610ea657610dda7f3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece41916101a860405190565b0390a1610de682610376565b6001610df361012a825490565b821015610e9e57610e10610e0a83610e1593610a31565b90610b47565b6104c0565b9063d90d786e91803b1561010257610e3b928591610e3260405190565b94859260e01b90565b8252818381610e4e8c8a60048401610d79565b03925af1918215610e9957610e6892610e6d575b50610b26565b610de6565b610e8c90853d8711610e92575b610e848183610b52565b8101906102dd565b38610e62565b503d610e7a565b610d8a565b505050509050565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f4e6f207472616e73616374696f6e7320746f20636f6e766572740000000000006044820152606490fd5b9061011491610ad956fea2646970667358221220b67ba19c18a233d92d36e36ee5b48a1bbfdebb7ee33eda44400e9229076b658d64736f6c63430008140033",
}

// TransactionsAnalyzerABI is the input ABI used to generate the binding from.
// Deprecated: Use TransactionsAnalyzerMetaData.ABI instead.
var TransactionsAnalyzerABI = TransactionsAnalyzerMetaData.ABI

// TransactionsAnalyzerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransactionsAnalyzerMetaData.Bin instead.
var TransactionsAnalyzerBin = TransactionsAnalyzerMetaData.Bin

// DeployTransactionsAnalyzer deploys a new Ethereum contract, binding an instance of TransactionsAnalyzer to it.
func DeployTransactionsAnalyzer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TransactionsAnalyzer, error) {
	parsed, err := TransactionsAnalyzerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransactionsAnalyzerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransactionsAnalyzer{TransactionsAnalyzerCaller: TransactionsAnalyzerCaller{contract: contract}, TransactionsAnalyzerTransactor: TransactionsAnalyzerTransactor{contract: contract}, TransactionsAnalyzerFilterer: TransactionsAnalyzerFilterer{contract: contract}}, nil
}

// TransactionsAnalyzer is an auto generated Go binding around an Ethereum contract.
type TransactionsAnalyzer struct {
	TransactionsAnalyzerCaller     // Read-only binding to the contract
	TransactionsAnalyzerTransactor // Write-only binding to the contract
	TransactionsAnalyzerFilterer   // Log filterer for contract events
}

// TransactionsAnalyzerCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransactionsAnalyzerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionsAnalyzerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransactionsAnalyzerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionsAnalyzerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransactionsAnalyzerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionsAnalyzerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransactionsAnalyzerSession struct {
	Contract     *TransactionsAnalyzer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TransactionsAnalyzerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransactionsAnalyzerCallerSession struct {
	Contract *TransactionsAnalyzerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// TransactionsAnalyzerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransactionsAnalyzerTransactorSession struct {
	Contract     *TransactionsAnalyzerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// TransactionsAnalyzerRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransactionsAnalyzerRaw struct {
	Contract *TransactionsAnalyzer // Generic contract binding to access the raw methods on
}

// TransactionsAnalyzerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransactionsAnalyzerCallerRaw struct {
	Contract *TransactionsAnalyzerCaller // Generic read-only contract binding to access the raw methods on
}

// TransactionsAnalyzerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransactionsAnalyzerTransactorRaw struct {
	Contract *TransactionsAnalyzerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransactionsAnalyzer creates a new instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzer(address common.Address, backend bind.ContractBackend) (*TransactionsAnalyzer, error) {
	contract, err := bindTransactionsAnalyzer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzer{TransactionsAnalyzerCaller: TransactionsAnalyzerCaller{contract: contract}, TransactionsAnalyzerTransactor: TransactionsAnalyzerTransactor{contract: contract}, TransactionsAnalyzerFilterer: TransactionsAnalyzerFilterer{contract: contract}}, nil
}

// NewTransactionsAnalyzerCaller creates a new read-only instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzerCaller(address common.Address, caller bind.ContractCaller) (*TransactionsAnalyzerCaller, error) {
	contract, err := bindTransactionsAnalyzer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerCaller{contract: contract}, nil
}

// NewTransactionsAnalyzerTransactor creates a new write-only instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzerTransactor(address common.Address, transactor bind.ContractTransactor) (*TransactionsAnalyzerTransactor, error) {
	contract, err := bindTransactionsAnalyzer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerTransactor{contract: contract}, nil
}

// NewTransactionsAnalyzerFilterer creates a new log filterer instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzerFilterer(address common.Address, filterer bind.ContractFilterer) (*TransactionsAnalyzerFilterer, error) {
	contract, err := bindTransactionsAnalyzer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerFilterer{contract: contract}, nil
}

// bindTransactionsAnalyzer binds a generic wrapper to an already deployed contract.
func bindTransactionsAnalyzer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransactionsAnalyzerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransactionsAnalyzer *TransactionsAnalyzerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransactionsAnalyzer.Contract.TransactionsAnalyzerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransactionsAnalyzer *TransactionsAnalyzerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.TransactionsAnalyzerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransactionsAnalyzer *TransactionsAnalyzerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.TransactionsAnalyzerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransactionsAnalyzer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.DEFAULTADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.DEFAULTADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) EOAADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "EOA_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) EOAADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.EOAADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) EOAADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.EOAADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// HOOKCALLERROLE is a free data retrieval call binding the contract method 0x5f03a661.
//
// Solidity: function HOOK_CALLER_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) HOOKCALLERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "HOOK_CALLER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HOOKCALLERROLE is a free data retrieval call binding the contract method 0x5f03a661.
//
// Solidity: function HOOK_CALLER_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) HOOKCALLERROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.HOOKCALLERROLE(&_TransactionsAnalyzer.CallOpts)
}

// HOOKCALLERROLE is a free data retrieval call binding the contract method 0x5f03a661.
//
// Solidity: function HOOK_CALLER_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) HOOKCALLERROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.HOOKCALLERROLE(&_TransactionsAnalyzer.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.GetRoleAdmin(&_TransactionsAnalyzer.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.GetRoleAdmin(&_TransactionsAnalyzer.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TransactionsAnalyzer.Contract.HasRole(&_TransactionsAnalyzer.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TransactionsAnalyzer.Contract.HasRole(&_TransactionsAnalyzer.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TransactionsAnalyzer.Contract.SupportsInterface(&_TransactionsAnalyzer.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TransactionsAnalyzer.Contract.SupportsInterface(&_TransactionsAnalyzer.CallOpts, interfaceId)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) AddOnBlockEndCallback(opts *bind.TransactOpts, callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "addOnBlockEndCallback", callbackAddress)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) AddOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.AddOnBlockEndCallback(&_TransactionsAnalyzer.TransactOpts, callbackAddress)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) AddOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.AddOnBlockEndCallback(&_TransactionsAnalyzer.TransactOpts, callbackAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.GrantRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.GrantRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address eoaAdmin, address authorizedCaller) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) Initialize(opts *bind.TransactOpts, eoaAdmin common.Address, authorizedCaller common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "initialize", eoaAdmin, authorizedCaller)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address eoaAdmin, address authorizedCaller) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) Initialize(eoaAdmin common.Address, authorizedCaller common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.Initialize(&_TransactionsAnalyzer.TransactOpts, eoaAdmin, authorizedCaller)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address eoaAdmin, address authorizedCaller) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) Initialize(eoaAdmin common.Address, authorizedCaller common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.Initialize(&_TransactionsAnalyzer.TransactOpts, eoaAdmin, authorizedCaller)
}

// OnBlock is a paid mutator transaction binding the contract method 0x4a44e6c1.
//
// Solidity: function onBlock((uint8,uint256,uint256,uint256,address,uint256,bytes,address)[] transactions) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) OnBlock(opts *bind.TransactOpts, transactions []StructsTransaction) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "onBlock", transactions)
}

// OnBlock is a paid mutator transaction binding the contract method 0x4a44e6c1.
//
// Solidity: function onBlock((uint8,uint256,uint256,uint256,address,uint256,bytes,address)[] transactions) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) OnBlock(transactions []StructsTransaction) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.OnBlock(&_TransactionsAnalyzer.TransactOpts, transactions)
}

// OnBlock is a paid mutator transaction binding the contract method 0x4a44e6c1.
//
// Solidity: function onBlock((uint8,uint256,uint256,uint256,address,uint256,bytes,address)[] transactions) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) OnBlock(transactions []StructsTransaction) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.OnBlock(&_TransactionsAnalyzer.TransactOpts, transactions)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RenounceRole(&_TransactionsAnalyzer.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RenounceRole(&_TransactionsAnalyzer.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RevokeRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RevokeRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// TransactionsAnalyzerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerInitializedIterator struct {
	Event *TransactionsAnalyzerInitialized // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerInitialized)
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
		it.Event = new(TransactionsAnalyzerInitialized)
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
func (it *TransactionsAnalyzerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerInitialized represents a Initialized event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterInitialized(opts *bind.FilterOpts) (*TransactionsAnalyzerInitializedIterator, error) {

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerInitializedIterator{contract: _TransactionsAnalyzer.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerInitialized) (event.Subscription, error) {

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerInitialized)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseInitialized(log types.Log) (*TransactionsAnalyzerInitialized, error) {
	event := new(TransactionsAnalyzerInitialized)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionsAnalyzerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleAdminChangedIterator struct {
	Event *TransactionsAnalyzerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerRoleAdminChanged)
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
		it.Event = new(TransactionsAnalyzerRoleAdminChanged)
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
func (it *TransactionsAnalyzerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerRoleAdminChanged represents a RoleAdminChanged event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TransactionsAnalyzerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerRoleAdminChangedIterator{contract: _TransactionsAnalyzer.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerRoleAdminChanged)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseRoleAdminChanged(log types.Log) (*TransactionsAnalyzerRoleAdminChanged, error) {
	event := new(TransactionsAnalyzerRoleAdminChanged)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionsAnalyzerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleGrantedIterator struct {
	Event *TransactionsAnalyzerRoleGranted // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerRoleGranted)
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
		it.Event = new(TransactionsAnalyzerRoleGranted)
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
func (it *TransactionsAnalyzerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerRoleGranted represents a RoleGranted event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TransactionsAnalyzerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerRoleGrantedIterator{contract: _TransactionsAnalyzer.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerRoleGranted)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseRoleGranted(log types.Log) (*TransactionsAnalyzerRoleGranted, error) {
	event := new(TransactionsAnalyzerRoleGranted)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionsAnalyzerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleRevokedIterator struct {
	Event *TransactionsAnalyzerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerRoleRevoked)
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
		it.Event = new(TransactionsAnalyzerRoleRevoked)
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
func (it *TransactionsAnalyzerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerRoleRevoked represents a RoleRevoked event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TransactionsAnalyzerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerRoleRevokedIterator{contract: _TransactionsAnalyzer.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerRoleRevoked)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseRoleRevoked(log types.Log) (*TransactionsAnalyzerRoleRevoked, error) {
	event := new(TransactionsAnalyzerRoleRevoked)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionsAnalyzerTransactionsConvertedIterator is returned from FilterTransactionsConverted and is used to iterate over the raw logs and unpacked data for TransactionsConverted events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerTransactionsConvertedIterator struct {
	Event *TransactionsAnalyzerTransactionsConverted // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerTransactionsConvertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerTransactionsConverted)
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
		it.Event = new(TransactionsAnalyzerTransactionsConverted)
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
func (it *TransactionsAnalyzerTransactionsConvertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerTransactionsConvertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerTransactionsConverted represents a TransactionsConverted event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerTransactionsConverted struct {
	TransactionsLength *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterTransactionsConverted is a free log retrieval operation binding the contract event 0x3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece41.
//
// Solidity: event TransactionsConverted(uint256 transactionsLength)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterTransactionsConverted(opts *bind.FilterOpts) (*TransactionsAnalyzerTransactionsConvertedIterator, error) {

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "TransactionsConverted")
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerTransactionsConvertedIterator{contract: _TransactionsAnalyzer.contract, event: "TransactionsConverted", logs: logs, sub: sub}, nil
}

// WatchTransactionsConverted is a free log subscription operation binding the contract event 0x3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece41.
//
// Solidity: event TransactionsConverted(uint256 transactionsLength)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchTransactionsConverted(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerTransactionsConverted) (event.Subscription, error) {

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "TransactionsConverted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerTransactionsConverted)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "TransactionsConverted", log); err != nil {
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

// ParseTransactionsConverted is a log parse operation binding the contract event 0x3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece41.
//
// Solidity: event TransactionsConverted(uint256 transactionsLength)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseTransactionsConverted(log types.Log) (*TransactionsAnalyzerTransactionsConverted, error) {
	event := new(TransactionsAnalyzerTransactionsConverted)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "TransactionsConverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
