// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObsERC20

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

// ObsERC20MetaData contains all meta data concerning the ObsERC20 contract.
var ObsERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600580546001600160a01b03191673deb34a740eca1ec42c8b8204cbec0ba34fdd27f31790553480156200003757600080fd5b5060405162000c4738038062000c478339810160408190526200005a9162000149565b818160036200006a838262000242565b50600462000079828262000242565b50505050506200030e565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620000ac57600080fd5b81516001600160401b0380821115620000c957620000c962000084565b604051601f8301601f19908116603f01168101908282118183101715620000f457620000f462000084565b816040528381526020925086838588010111156200011157600080fd5b600091505b8382101562000135578582018301518183018401529082019062000116565b600093810190920192909252949350505050565b600080604083850312156200015d57600080fd5b82516001600160401b03808211156200017557600080fd5b62000183868387016200009a565b935060208501519150808211156200019a57600080fd5b50620001a9858286016200009a565b9150509250929050565b600181811c90821680620001c857607f821691505b602082108103620001e957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200023d57600081815260208120601f850160051c81016020861015620002185750805b601f850160051c820191505b81811015620002395782815560010162000224565b5050505b505050565b81516001600160401b038111156200025e576200025e62000084565b62000276816200026f8454620001b3565b84620001ef565b602080601f831160018114620002ae5760008415620002955750858301515b600019600386901b1c1916600185901b17855562000239565b600085815260208120601f198616915b82811015620002df57888601518255948401946001909101908401620002be565b5085821015620002fe5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610929806200031e6000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063313ce5671161007657806395d89b411161005b57806395d89b4114610130578063a9059cbb14610138578063dd62ed3e1461014b57600080fd5b8063313ce5671461010e57806370a082311461011d57600080fd5b806306fdde03146100a8578063095ea7b3146100c657806318160ddd146100e957806323b872dd146100fb575b600080fd5b6100b061015e565b6040516100bd9190610773565b60405180910390f35b6100d96100d43660046107dd565b6101f0565b60405190151581526020016100bd565b6002545b6040519081526020016100bd565b6100d9610109366004610807565b61020a565b604051601281526020016100bd565b6100ed61012b366004610843565b61022e565b6100b06102d9565b6100d96101463660046107dd565b6102e8565b6100ed610159366004610865565b6102f6565b60606003805461016d90610898565b80601f016020809104026020016040519081016040528092919081815260200182805461019990610898565b80156101e65780601f106101bb576101008083540402835291602001916101e6565b820191906000526020600020905b8154815290600101906020018083116101c957829003601f168201915b5050505050905090565b6000336101fe818585610407565b60019150505b92915050565b600033610218858285610419565b610223858585610498565b506001949350505050565b60006001600160a01b038216320361025e576001600160a01b038216600090815260208190526040902054610204565b6001600160a01b038216330361028c576001600160a01b038216600090815260208190526040902054610204565b60405162461bcd60e51b815260206004820152601f60248201527f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e63650060448201526064015b60405180910390fd5b60606004805461016d90610898565b6000336101fe818585610498565b6000326001600160a01b03841614806103175750326001600160a01b038316145b1561034a576001600160a01b038084166000908152600160209081526040808320938616835292905220545b9050610204565b336001600160a01b03841614806103695750336001600160a01b038316145b15610399576001600160a01b03808416600090815260016020908152604080832093861683529290522054610343565b60405162461bcd60e51b815260206004820152602160248201527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360448201527f650000000000000000000000000000000000000000000000000000000000000060648201526084016102d0565b6104148383836001610529565b505050565b600061042584846102f6565b905060001981146104925781811015610483576040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526001600160a01b038416600482015260248101829052604481018390526064016102d0565b61049284848484036000610529565b50505050565b6001600160a01b0383166104db576040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600060048201526024016102d0565b6001600160a01b03821661051e576040517fec442f05000000000000000000000000000000000000000000000000000000008152600060048201526024016102d0565b610414838383610630565b6001600160a01b03841661056c576040517fe602df05000000000000000000000000000000000000000000000000000000008152600060048201526024016102d0565b6001600160a01b0383166105af576040517f94280d62000000000000000000000000000000000000000000000000000000008152600060048201526024016102d0565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561049257826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161062291815260200190565b60405180910390a350505050565b6001600160a01b03831661065b57806002600082825461065091906108d2565b909155506106e69050565b6001600160a01b038316600090815260208190526040902054818110156106c7576040517fe450d38c0000000000000000000000000000000000000000000000000000000081526001600160a01b038516600482015260248101829052604481018390526064016102d0565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b03821661070257600280548290039055610721565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161076691815260200190565b60405180910390a3505050565b600060208083528351808285015260005b818110156107a057858101830151858201604001528201610784565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146107d857600080fd5b919050565b600080604083850312156107f057600080fd5b6107f9836107c1565b946020939093013593505050565b60008060006060848603121561081c57600080fd5b610825846107c1565b9250610833602085016107c1565b9150604084013590509250925092565b60006020828403121561085557600080fd5b61085e826107c1565b9392505050565b6000806040838503121561087857600080fd5b610881836107c1565b915061088f602084016107c1565b90509250929050565b600181811c908216806108ac57607f821691505b6020821081036108cc57634e487b7160e01b600052602260045260246000fd5b50919050565b8082018082111561020457634e487b7160e01b600052601160045260246000fdfea2646970667358221220747492c6ad6919950d95fb045569a252ef4e67afdd9357bceeee7893599576cd64736f6c63430008140033",
}

// ObsERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ObsERC20MetaData.ABI instead.
var ObsERC20ABI = ObsERC20MetaData.ABI

// ObsERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObsERC20MetaData.Bin instead.
var ObsERC20Bin = ObsERC20MetaData.Bin

// DeployObsERC20 deploys a new Ethereum contract, binding an instance of ObsERC20 to it.
func DeployObsERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *ObsERC20, error) {
	parsed, err := ObsERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObsERC20Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObsERC20{ObsERC20Caller: ObsERC20Caller{contract: contract}, ObsERC20Transactor: ObsERC20Transactor{contract: contract}, ObsERC20Filterer: ObsERC20Filterer{contract: contract}}, nil
}

// ObsERC20 is an auto generated Go binding around an Ethereum contract.
type ObsERC20 struct {
	ObsERC20Caller     // Read-only binding to the contract
	ObsERC20Transactor // Write-only binding to the contract
	ObsERC20Filterer   // Log filterer for contract events
}

// ObsERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ObsERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ObsERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObsERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObsERC20Session struct {
	Contract     *ObsERC20         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObsERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObsERC20CallerSession struct {
	Contract *ObsERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ObsERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObsERC20TransactorSession struct {
	Contract     *ObsERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ObsERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ObsERC20Raw struct {
	Contract *ObsERC20 // Generic contract binding to access the raw methods on
}

// ObsERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObsERC20CallerRaw struct {
	Contract *ObsERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ObsERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObsERC20TransactorRaw struct {
	Contract *ObsERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewObsERC20 creates a new instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20(address common.Address, backend bind.ContractBackend) (*ObsERC20, error) {
	contract, err := bindObsERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObsERC20{ObsERC20Caller: ObsERC20Caller{contract: contract}, ObsERC20Transactor: ObsERC20Transactor{contract: contract}, ObsERC20Filterer: ObsERC20Filterer{contract: contract}}, nil
}

// NewObsERC20Caller creates a new read-only instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Caller(address common.Address, caller bind.ContractCaller) (*ObsERC20Caller, error) {
	contract, err := bindObsERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Caller{contract: contract}, nil
}

// NewObsERC20Transactor creates a new write-only instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ObsERC20Transactor, error) {
	contract, err := bindObsERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Transactor{contract: contract}, nil
}

// NewObsERC20Filterer creates a new log filterer instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ObsERC20Filterer, error) {
	contract, err := bindObsERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Filterer{contract: contract}, nil
}

// bindObsERC20 binds a generic wrapper to an already deployed contract.
func bindObsERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ObsERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObsERC20 *ObsERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObsERC20.Contract.ObsERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObsERC20 *ObsERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObsERC20.Contract.ObsERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObsERC20 *ObsERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObsERC20.Contract.ObsERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObsERC20 *ObsERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObsERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObsERC20 *ObsERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObsERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObsERC20 *ObsERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObsERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.Allowance(&_ObsERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.Allowance(&_ObsERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.BalanceOf(&_ObsERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.BalanceOf(&_ObsERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20Session) Decimals() (uint8, error) {
	return _ObsERC20.Contract.Decimals(&_ObsERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20CallerSession) Decimals() (uint8, error) {
	return _ObsERC20.Contract.Decimals(&_ObsERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20Session) Name() (string, error) {
	return _ObsERC20.Contract.Name(&_ObsERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20CallerSession) Name() (string, error) {
	return _ObsERC20.Contract.Name(&_ObsERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20Session) Symbol() (string, error) {
	return _ObsERC20.Contract.Symbol(&_ObsERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20CallerSession) Symbol() (string, error) {
	return _ObsERC20.Contract.Symbol(&_ObsERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20Session) TotalSupply() (*big.Int, error) {
	return _ObsERC20.Contract.TotalSupply(&_ObsERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ObsERC20.Contract.TotalSupply(&_ObsERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Approve(&_ObsERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Approve(&_ObsERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Transfer(&_ObsERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Transfer(&_ObsERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.TransferFrom(&_ObsERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.TransferFrom(&_ObsERC20.TransactOpts, from, to, value)
}

// ObsERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ObsERC20 contract.
type ObsERC20ApprovalIterator struct {
	Event *ObsERC20Approval // Event containing the contract specifics and raw log

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
func (it *ObsERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObsERC20Approval)
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
		it.Event = new(ObsERC20Approval)
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
func (it *ObsERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObsERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObsERC20Approval represents a Approval event raised by the ObsERC20 contract.
type ObsERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ObsERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObsERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ObsERC20ApprovalIterator{contract: _ObsERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ObsERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObsERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObsERC20Approval)
				if err := _ObsERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ObsERC20 *ObsERC20Filterer) ParseApproval(log types.Log) (*ObsERC20Approval, error) {
	event := new(ObsERC20Approval)
	if err := _ObsERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObsERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ObsERC20 contract.
type ObsERC20TransferIterator struct {
	Event *ObsERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ObsERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObsERC20Transfer)
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
		it.Event = new(ObsERC20Transfer)
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
func (it *ObsERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObsERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObsERC20Transfer represents a Transfer event raised by the ObsERC20 contract.
type ObsERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ObsERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObsERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ObsERC20TransferIterator{contract: _ObsERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ObsERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObsERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObsERC20Transfer)
				if err := _ObsERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ObsERC20 *ObsERC20Filterer) ParseTransfer(log types.Log) (*ObsERC20Transfer, error) {
	event := new(ObsERC20Transfer)
	if err := _ObsERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
