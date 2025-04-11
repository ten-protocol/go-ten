// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TenERC20

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

// TenERC20MetaData contains all meta data concerning the TenERC20 contract.
var TenERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600580546001600160a01b03191673deb34a740eca1ec42c8b8204cbec0ba34fdd27f3179055348015610035575f5ffd5b50604051610d50380380610d508339810160408190526100549161016f565b8181600361006283826102c2565b50600461006f82826102c2565b505050505061037d565b634e487b7160e01b5f52604160045260245ffd5b601f19601f83011681016001600160401b03811182821017156100b2576100b2610079565b6040525050565b5f6100c360405190565b90506100cf828261008d565b919050565b5f6001600160401b038211156100ec576100ec610079565b601f19601f83011660200192915050565b8281835e505f910152565b5f61011a610115846100d4565b6100b9565b9050828152838383011115610130576101305f5ffd5b61013e8360208301846100fd565b9392505050565b5f82601f830112610157576101575f5ffd5b8151610167848260208601610108565b949350505050565b5f5f60408385031215610183576101835f5ffd5b82516001600160401b0381111561019b5761019b5f5ffd5b6101a785828601610145565b602085015190935090506001600160401b038111156101c7576101c75f5ffd5b6101d385828601610145565b9150509250929050565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061020557607f821691505b602082108103610217576102176101dd565b50919050565b5f61022b6102288381565b90565b92915050565b61023a8361021d565b81545f1960089490940293841b1916921b91909117905550565b5f610260818484610231565b505050565b8181101561027f576102775f82610254565b600101610265565b5050565b601f821115610260575f818152602090206020601f850104810160208510156102a95750805b6102bb6020601f860104830182610265565b5050505050565b81516001600160401b038111156102db576102db610079565b6102e582546101f1565b6102f0828285610283565b506020601f821160018114610323575f831561030c5750848201515b5f19600885021c19811660028502178555506102bb565b5f84815260208120601f198516915b828110156103525787850151825560209485019460019092019101610332565b508482101561036e57838701515f19601f87166008021c191681555b50505050600202600101905550565b6109c68061038a5f395ff3fe608060405234801561000f575f5ffd5b506004361061009f575f3560e01c8063313ce5671161007257806395d89b411161005857806395d89b4114610127578063a9059cbb1461012f578063dd62ed3e14610142575f5ffd5b8063313ce5671461010557806370a0823114610114575f5ffd5b806306fdde03146100a3578063095ea7b3146100c157806318160ddd146100e157806323b872dd146100f2575b5f5ffd5b6100ab610155565b6040516100b891906106e3565b60405180910390f35b6100d46100cf36600461073d565b6101e5565b6040516100b8919061077d565b6002545b6040516100b89190610791565b6100d461010036600461079f565b6101fe565b60126040516100b891906107ee565b6100e56101223660046107fc565b610221565b6100ab61029d565b6100d461013d36600461073d565b6102ac565b6100e5610150366004610819565b6102b9565b6060600380546101649061085a565b80601f01602080910402602001604051908101604052809291908181526020018280546101909061085a565b80156101db5780601f106101b2576101008083540402835291602001916101db565b820191905f5260205f20905b8154815290600101906020018083116101be57829003601f168201915b5050505050905090565b5f336101f2818585610371565b60019150505b92915050565b5f3361020b858285610383565b6102168585856103ec565b506001949350505050565b5f6001600160a01b038216320361024f576001600160a01b0382165f908152602081905260409020546101f8565b6001600160a01b038216330361027c576001600160a01b0382165f908152602081905260409020546101f8565b60405162461bcd60e51b815260040161029490610886565b60405180910390fd5b6060600480546101649061085a565b5f336101f28185856103ec565b5f326001600160a01b03841614806102d95750326001600160a01b038316145b1561030b576001600160a01b038084165f908152600160209081526040808320938616835292905220545b90506101f8565b336001600160a01b038416148061032a5750336001600160a01b038316145b15610359576001600160a01b038084165f90815260016020908152604080832093861683529290522054610304565b60405162461bcd60e51b8152600401610294906108c1565b61037e838383600161047b565b505050565b5f61038e84846102b9565b90505f1981146103e657818110156103d8578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016102949392919061092b565b6103e684848484035f61047b565b50505050565b6001600160a01b03831661042e575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610294919061095b565b6001600160a01b038216610470575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610294919061095b565b61037e83838361057d565b6001600160a01b0384166104bd575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401610294919061095b565b6001600160a01b0383166104ff575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610294919061095b565b6001600160a01b038085165f90815260016020908152604080832093871683529290522082905580156103e657826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161056f9190610791565b60405180910390a350505050565b6001600160a01b0383166105a7578060025f82825461059c919061097d565b9091555061061d9050565b6001600160a01b0383165f90815260208190526040902054818110156105ff578381836040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016102949392919061092b565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b03821661063957600280548290039055610657565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161069a9190610791565b60405180910390a3505050565b8281835e505f910152565b5f6106bb825190565b8084526020840193506106d28185602086016106a7565b601f01601f19169290920192915050565b602080825281016106f481846106b2565b9392505050565b5f6001600160a01b0382166101f8565b610714816106fb565b811461071e575f5ffd5b50565b80356101f88161070b565b80610714565b80356101f88161072c565b5f5f60408385031215610751576107515f5ffd5b61075b8484610721565b915061076a8460208501610732565b90509250929050565b8015155b82525050565b602081016101f88284610773565b80610777565b602081016101f8828461078b565b5f5f5f606084860312156107b4576107b45f5ffd5b6107be8585610721565b92506107cd8560208601610721565b91506107dc8560408601610732565b90509250925092565b60ff8116610777565b602081016101f882846107e5565b5f6020828403121561080f5761080f5f5ffd5b6106f48383610721565b5f5f6040838503121561082d5761082d5f5ffd5b6108378484610721565b915061076a8460208501610721565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061086e57607f821691505b60208210810361088057610880610846565b50919050565b602080825281016101f881601f81527f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e636500602082015260400190565b602080825281016101f881602181527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360208201527f6500000000000000000000000000000000000000000000000000000000000000604082015260600190565b610777816106fb565b606081016109398286610922565b610946602083018561078b565b610953604083018461078b565b949350505050565b602081016101f88284610922565b634e487b7160e01b5f52601160045260245ffd5b808201808211156101f8576101f861096956fea26469706673582212206e8a2c22bb56b0c8b5e566d758ffe90721896f36d081c34bafcbe05bec3fe83364736f6c634300081c0033",
}

// TenERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use TenERC20MetaData.ABI instead.
var TenERC20ABI = TenERC20MetaData.ABI

// TenERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TenERC20MetaData.Bin instead.
var TenERC20Bin = TenERC20MetaData.Bin

// DeployTenERC20 deploys a new Ethereum contract, binding an instance of TenERC20 to it.
func DeployTenERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *TenERC20, error) {
	parsed, err := TenERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TenERC20Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TenERC20{TenERC20Caller: TenERC20Caller{contract: contract}, TenERC20Transactor: TenERC20Transactor{contract: contract}, TenERC20Filterer: TenERC20Filterer{contract: contract}}, nil
}

// TenERC20 is an auto generated Go binding around an Ethereum contract.
type TenERC20 struct {
	TenERC20Caller     // Read-only binding to the contract
	TenERC20Transactor // Write-only binding to the contract
	TenERC20Filterer   // Log filterer for contract events
}

// TenERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type TenERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TenERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TenERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TenERC20Session struct {
	Contract     *TenERC20         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TenERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TenERC20CallerSession struct {
	Contract *TenERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TenERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TenERC20TransactorSession struct {
	Contract     *TenERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TenERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type TenERC20Raw struct {
	Contract *TenERC20 // Generic contract binding to access the raw methods on
}

// TenERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TenERC20CallerRaw struct {
	Contract *TenERC20Caller // Generic read-only contract binding to access the raw methods on
}

// TenERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TenERC20TransactorRaw struct {
	Contract *TenERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTenERC20 creates a new instance of TenERC20, bound to a specific deployed contract.
func NewTenERC20(address common.Address, backend bind.ContractBackend) (*TenERC20, error) {
	contract, err := bindTenERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TenERC20{TenERC20Caller: TenERC20Caller{contract: contract}, TenERC20Transactor: TenERC20Transactor{contract: contract}, TenERC20Filterer: TenERC20Filterer{contract: contract}}, nil
}

// NewTenERC20Caller creates a new read-only instance of TenERC20, bound to a specific deployed contract.
func NewTenERC20Caller(address common.Address, caller bind.ContractCaller) (*TenERC20Caller, error) {
	contract, err := bindTenERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TenERC20Caller{contract: contract}, nil
}

// NewTenERC20Transactor creates a new write-only instance of TenERC20, bound to a specific deployed contract.
func NewTenERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*TenERC20Transactor, error) {
	contract, err := bindTenERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TenERC20Transactor{contract: contract}, nil
}

// NewTenERC20Filterer creates a new log filterer instance of TenERC20, bound to a specific deployed contract.
func NewTenERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*TenERC20Filterer, error) {
	contract, err := bindTenERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TenERC20Filterer{contract: contract}, nil
}

// bindTenERC20 binds a generic wrapper to an already deployed contract.
func bindTenERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TenERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenERC20 *TenERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenERC20.Contract.TenERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenERC20 *TenERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenERC20.Contract.TenERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenERC20 *TenERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenERC20.Contract.TenERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenERC20 *TenERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenERC20 *TenERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenERC20 *TenERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TenERC20 *TenERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TenERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TenERC20 *TenERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TenERC20.Contract.Allowance(&_TenERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TenERC20 *TenERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TenERC20.Contract.Allowance(&_TenERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TenERC20 *TenERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TenERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TenERC20 *TenERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _TenERC20.Contract.BalanceOf(&_TenERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TenERC20 *TenERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TenERC20.Contract.BalanceOf(&_TenERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TenERC20 *TenERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TenERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TenERC20 *TenERC20Session) Decimals() (uint8, error) {
	return _TenERC20.Contract.Decimals(&_TenERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TenERC20 *TenERC20CallerSession) Decimals() (uint8, error) {
	return _TenERC20.Contract.Decimals(&_TenERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TenERC20 *TenERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TenERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TenERC20 *TenERC20Session) Name() (string, error) {
	return _TenERC20.Contract.Name(&_TenERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TenERC20 *TenERC20CallerSession) Name() (string, error) {
	return _TenERC20.Contract.Name(&_TenERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TenERC20 *TenERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TenERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TenERC20 *TenERC20Session) Symbol() (string, error) {
	return _TenERC20.Contract.Symbol(&_TenERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TenERC20 *TenERC20CallerSession) Symbol() (string, error) {
	return _TenERC20.Contract.Symbol(&_TenERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TenERC20 *TenERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TenERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TenERC20 *TenERC20Session) TotalSupply() (*big.Int, error) {
	return _TenERC20.Contract.TotalSupply(&_TenERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TenERC20 *TenERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _TenERC20.Contract.TotalSupply(&_TenERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TenERC20 *TenERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TenERC20 *TenERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.Contract.Approve(&_TenERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TenERC20 *TenERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.Contract.Approve(&_TenERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TenERC20 *TenERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TenERC20 *TenERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.Contract.Transfer(&_TenERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TenERC20 *TenERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.Contract.Transfer(&_TenERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TenERC20 *TenERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TenERC20 *TenERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.Contract.TransferFrom(&_TenERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TenERC20 *TenERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TenERC20.Contract.TransferFrom(&_TenERC20.TransactOpts, from, to, value)
}

// TenERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TenERC20 contract.
type TenERC20ApprovalIterator struct {
	Event *TenERC20Approval // Event containing the contract specifics and raw log

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
func (it *TenERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenERC20Approval)
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
		it.Event = new(TenERC20Approval)
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
func (it *TenERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenERC20Approval represents a Approval event raised by the TenERC20 contract.
type TenERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TenERC20 *TenERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TenERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TenERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TenERC20ApprovalIterator{contract: _TenERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TenERC20 *TenERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TenERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TenERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenERC20Approval)
				if err := _TenERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TenERC20 *TenERC20Filterer) ParseApproval(log types.Log) (*TenERC20Approval, error) {
	event := new(TenERC20Approval)
	if err := _TenERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TenERC20 contract.
type TenERC20TransferIterator struct {
	Event *TenERC20Transfer // Event containing the contract specifics and raw log

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
func (it *TenERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenERC20Transfer)
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
		it.Event = new(TenERC20Transfer)
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
func (it *TenERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenERC20Transfer represents a Transfer event raised by the TenERC20 contract.
type TenERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TenERC20 *TenERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TenERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TenERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TenERC20TransferIterator{contract: _TenERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TenERC20 *TenERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TenERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TenERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenERC20Transfer)
				if err := _TenERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TenERC20 *TenERC20Filterer) ParseTransfer(log types.Log) (*TenERC20Transfer, error) {
	event := new(TenERC20Transfer)
	if err := _TenERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
