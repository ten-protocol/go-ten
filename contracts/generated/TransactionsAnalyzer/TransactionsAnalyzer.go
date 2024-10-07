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
	Bin: "0x608060405234801561001057600080fd5b50610dc7806100206000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c8063508a50f411610081578063a217fddf1161005b578063a217fddf14610205578063d547741f1461020d578063ee546fd81461022057600080fd5b8063508a50f4146101805780635f03a661146101a757806391d14854146101ce57600080fd5b806336568abe116100b257806336568abe14610147578063485cc9551461015a5780634a44e6c11461016d57600080fd5b806301ffc9a7146100d9578063248a9ca3146101025780632f2ff15d14610132575b600080fd5b6100ec6100e73660046108a6565b610297565b6040516100f991906108d9565b60405180910390f35b6101256101103660046108f8565b60009081526020819052604090206001015490565b6040516100f9919061091f565b610145610140366004610952565b610330565b005b610145610155366004610952565b61035b565b61014561016836600461098f565b6103ac565b61014561017b366004610a03565b610548565b6101257ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59881565b6101257f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db81565b6100ec6101dc366004610952565b6000918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b610125600081565b61014561021b366004610952565b6106ab565b61014561022e366004610a4b565b6001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061032a57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b60008281526020819052604090206001015461034b816106d0565b61035583836106dd565b50505050565b6001600160a01b038116331461039d576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103a78282610787565b505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156103f75750825b905060008267ffffffffffffffff1660011480156104145750303b155b905081158015610422575080155b15610459576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561048d57845468ff00000000000000001916680100000000000000001785555b6104986000886106dd565b506104c37ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a598886106dd565b506104ee7f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db876106dd565b50831561053f57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061053690600190610a87565b60405180910390a15b50505050505050565b7f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db610572816106d0565b60008290036105b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105ad90610a95565b60405180910390fd5b6040517f3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece41906105e690849061091f565b60405180910390a160005b6001548110156103555760006001828154811061061057610610610ad0565b6000918252602090912001546040517fd90d786e0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169150819063d90d786e906106679088908890600401610d35565b600060405180830381600087803b15801561068157600080fd5b505af1158015610695573d6000803e3d6000fd5b5050505050806106a490610d5d565b90506105f1565b6000828152602081905260409020600101546106c6816106d0565b6103558383610787565b6106da813361080a565b50565b6000828152602081815260408083206001600160a01b038516845290915281205460ff1661077f576000838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556107373390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161032a565b50600061032a565b6000828152602081815260408083206001600160a01b038516845290915281205460ff161561077f576000838152602081815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161032a565b6000828152602081815260408083206001600160a01b038516845290915290205460ff166108685780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016105ad929190610d76565b5050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b81146106da57600080fd5b803561032a8161086c565b6000602082840312156108bb576108bb600080fd5b60006108c7848461089b565b949350505050565b8015155b82525050565b6020810161032a82846108cf565b80610890565b803561032a816108e7565b60006020828403121561090d5761090d600080fd5b60006108c784846108ed565b806108d3565b6020810161032a8284610919565b60006001600160a01b03821661032a565b6108908161092d565b803561032a8161093e565b6000806040838503121561096857610968600080fd5b600061097485856108ed565b925050602061098585828601610947565b9150509250929050565b600080604083850312156109a5576109a5600080fd5b60006109748585610947565b60008083601f8401126109c6576109c6600080fd5b50813567ffffffffffffffff8111156109e1576109e1600080fd5b6020830191508360208202830111156109fc576109fc600080fd5b9250929050565b60008060208385031215610a1957610a19600080fd5b823567ffffffffffffffff811115610a3357610a33600080fd5b610a3f858286016109b1565b92509250509250929050565b600060208284031215610a6057610a60600080fd5b60006108c78484610947565b600067ffffffffffffffff821661032a565b6108d381610a6c565b6020810161032a8284610a7e565b6020808252810161032a81601a81527f4e6f207472616e73616374696f6e7320746f20636f6e76657274000000000000602082015260400190565b634e487b7160e01b600052603260045260246000fd5b60ff8116610890565b803561032a81610ae6565b50600061032a6020830183610aef565b60ff81166108d3565b50600061032a60208301836108ed565b50600061032a6020830183610947565b6108d38161092d565b6000808335601e1936859003018112610b5757610b57600080fd5b830160208101925035905067ffffffffffffffff811115610b7a57610b7a600080fd5b368190038213156109fc576109fc600080fd5b82818337506000910152565b818352602083019250610bad828483610b8d565b50601f01601f19160190565b60006101008301610bca8380610afa565b610bd48582610b0a565b50610be26020840184610b13565b610bef6020860182610919565b50610bfd6040840184610b13565b610c0a6040860182610919565b50610c186060840184610b13565b610c256060860182610919565b50610c336080840184610b23565b610c406080860182610b33565b50610c4e60a0840184610b13565b610c5b60a0860182610919565b50610c6960c0840184610b3c565b85830360c0870152610c7c838284610b99565b92505050610c8d60e0840184610b23565b610c9a60e0860182610b33565b509392505050565b6000610cae8383610bb9565b9392505050565b6000823560fe1936849003018112610ccf57610ccf600080fd5b90910192915050565b818352602083019250600083602084028101838060005b87811015610d28578484038952610d068284610cb5565b610d108582610ca2565b94505060208201602099909901989150600101610cef565b5091979650505050505050565b602080825281016108c7818486610cd8565b634e487b7160e01b600052601160045260246000fd5b600060018201610d6f57610d6f610d47565b5060010190565b60408101610d848285610b33565b610cae602083018461091956fea2646970667358221220557a7a45f4c7a153aa142c304ad3537d6721a1738e96e14650618f178649ce0764736f6c63430008140033",
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
