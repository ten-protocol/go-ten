// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConstantSupplyERC20

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

// ConstantSupplyERC20MetaData contains all meta data concerning the ConstantSupplyERC20 contract.
var ConstantSupplyERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b5060405162000c3038038062000c3083398101604081905262000033916200029d565b8282600362000043838262000397565b50600462000052828262000397565b5050506200006733826200007060201b60201c565b50505062000485565b6001600160a01b0382166200009f5760405163ec442f0560e01b81525f60048201526024015b60405180910390fd5b620000ac5f8383620000b0565b5050565b6001600160a01b038316620000de578060025f828254620000d291906200045f565b90915550620001509050565b6001600160a01b0383165f9081526020819052604090205481811015620001325760405163391434e360e21b81526001600160a01b0385166004820152602481018290526044810183905260640162000096565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b0382166200016e576002805482900390556200018c565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620001d291815260200190565b60405180910390a3505050565b634e487b7160e01b5f52604160045260245ffd5b5f82601f83011262000203575f80fd5b81516001600160401b0380821115620002205762000220620001df565b604051601f8301601f19908116603f011681019082821181831017156200024b576200024b620001df565b8160405283815260209250868385880101111562000267575f80fd5b5f91505b838210156200028a57858201830151818301840152908201906200026b565b5f93810190920192909252949350505050565b5f805f60608486031215620002b0575f80fd5b83516001600160401b0380821115620002c7575f80fd5b620002d587838801620001f3565b94506020860151915080821115620002eb575f80fd5b50620002fa86828701620001f3565b925050604084015190509250925092565b600181811c908216806200032057607f821691505b6020821081036200033f57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111562000392575f81815260208120601f850160051c810160208610156200036d5750805b601f850160051c820191505b818110156200038e5782815560010162000379565b5050505b505050565b81516001600160401b03811115620003b357620003b3620001df565b620003cb81620003c484546200030b565b8462000345565b602080601f83116001811462000401575f8415620003e95750858301515b5f19600386901b1c1916600185901b1785556200038e565b5f85815260208120601f198616915b82811015620004315788860151825594840194600190910190840162000410565b50858210156200044f57878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b808201808211156200047f57634e487b7160e01b5f52601160045260245ffd5b92915050565b61079d80620004935f395ff3fe608060405234801561000f575f80fd5b506004361061009f575f3560e01c8063313ce5671161007257806395d89b411161005857806395d89b4114610140578063a9059cbb14610148578063dd62ed3e1461015b575f80fd5b8063313ce5671461010957806370a0823114610118575f80fd5b806306fdde03146100a3578063095ea7b3146100c157806318160ddd146100e457806323b872dd146100f6575b5f80fd5b6100ab610193565b6040516100b891906105f8565b60405180910390f35b6100d46100cf36600461065e565b610223565b60405190151581526020016100b8565b6002545b6040519081526020016100b8565b6100d4610104366004610686565b61023c565b604051601281526020016100b8565b6100e86101263660046106bf565b6001600160a01b03165f9081526020819052604090205490565b6100ab61025f565b6100d461015636600461065e565b61026e565b6100e86101693660046106df565b6001600160a01b039182165f90815260016020908152604080832093909416825291909152205490565b6060600380546101a290610710565b80601f01602080910402602001604051908101604052809291908181526020018280546101ce90610710565b80156102195780601f106101f057610100808354040283529160200191610219565b820191905f5260205f20905b8154815290600101906020018083116101fc57829003601f168201915b5050505050905090565b5f3361023081858561027b565b60019150505b92915050565b5f3361024985828561028d565b610254858585610326565b506001949350505050565b6060600480546101a290610710565b5f33610230818585610326565b61028883838360016103b5565b505050565b6001600160a01b038381165f908152600160209081526040808320938616835292905220545f1981146103205781811015610312576040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526001600160a01b038416600482015260248101829052604481018390526064015b60405180910390fd5b61032084848484035f6103b5565b50505050565b6001600160a01b038316610368576040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610309565b6001600160a01b0382166103aa576040517fec442f050000000000000000000000000000000000000000000000000000000081525f6004820152602401610309565b6102888383836104b9565b6001600160a01b0384166103f7576040517fe602df050000000000000000000000000000000000000000000000000000000081525f6004820152602401610309565b6001600160a01b038316610439576040517f94280d620000000000000000000000000000000000000000000000000000000081525f6004820152602401610309565b6001600160a01b038085165f908152600160209081526040808320938716835292905220829055801561032057826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516104ab91815260200190565b60405180910390a350505050565b6001600160a01b0383166104e3578060025f8282546104d89190610748565b9091555061056c9050565b6001600160a01b0383165f908152602081905260409020548181101561054e576040517fe450d38c0000000000000000000000000000000000000000000000000000000081526001600160a01b03851660048201526024810182905260448101839052606401610309565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b038216610588576002805482900390556105a6565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105eb91815260200190565b60405180910390a3505050565b5f6020808352835180828501525f5b8181101561062357858101830151858201604001528201610607565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610659575f80fd5b919050565b5f806040838503121561066f575f80fd5b61067883610643565b946020939093013593505050565b5f805f60608486031215610698575f80fd5b6106a184610643565b92506106af60208501610643565b9150604084013590509250925092565b5f602082840312156106cf575f80fd5b6106d882610643565b9392505050565b5f80604083850312156106f0575f80fd5b6106f983610643565b915061070760208401610643565b90509250929050565b600181811c9082168061072457607f821691505b60208210810361074257634e487b7160e01b5f52602260045260245ffd5b50919050565b8082018082111561023657634e487b7160e01b5f52601160045260245ffdfea2646970667358221220c2fead2f098c766cc98c4b17615ad00ab6cbb6c547c8bc6ac0ce17231eb3a9a364736f6c63430008140033",
}

// ConstantSupplyERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ConstantSupplyERC20MetaData.ABI instead.
var ConstantSupplyERC20ABI = ConstantSupplyERC20MetaData.ABI

// ConstantSupplyERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConstantSupplyERC20MetaData.Bin instead.
var ConstantSupplyERC20Bin = ConstantSupplyERC20MetaData.Bin

// DeployConstantSupplyERC20 deploys a new Ethereum contract, binding an instance of ConstantSupplyERC20 to it.
func DeployConstantSupplyERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int) (common.Address, *types.Transaction, *ConstantSupplyERC20, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConstantSupplyERC20Bin), backend, name, symbol, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// ConstantSupplyERC20 is an auto generated Go binding around an Ethereum contract.
type ConstantSupplyERC20 struct {
	ConstantSupplyERC20Caller     // Read-only binding to the contract
	ConstantSupplyERC20Transactor // Write-only binding to the contract
	ConstantSupplyERC20Filterer   // Log filterer for contract events
}

// ConstantSupplyERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConstantSupplyERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConstantSupplyERC20Session struct {
	Contract     *ConstantSupplyERC20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConstantSupplyERC20CallerSession struct {
	Contract *ConstantSupplyERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ConstantSupplyERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConstantSupplyERC20TransactorSession struct {
	Contract     *ConstantSupplyERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ConstantSupplyERC20Raw struct {
	Contract *ConstantSupplyERC20 // Generic contract binding to access the raw methods on
}

// ConstantSupplyERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20CallerRaw struct {
	Contract *ConstantSupplyERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ConstantSupplyERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20TransactorRaw struct {
	Contract *ConstantSupplyERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewConstantSupplyERC20 creates a new instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20(address common.Address, backend bind.ContractBackend) (*ConstantSupplyERC20, error) {
	contract, err := bindConstantSupplyERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// NewConstantSupplyERC20Caller creates a new read-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Caller(address common.Address, caller bind.ContractCaller) (*ConstantSupplyERC20Caller, error) {
	contract, err := bindConstantSupplyERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Caller{contract: contract}, nil
}

// NewConstantSupplyERC20Transactor creates a new write-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ConstantSupplyERC20Transactor, error) {
	contract, err := bindConstantSupplyERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Transactor{contract: contract}, nil
}

// NewConstantSupplyERC20Filterer creates a new log filterer instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ConstantSupplyERC20Filterer, error) {
	contract, err := bindConstantSupplyERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Filterer{contract: contract}, nil
}

// bindConstantSupplyERC20 binds a generic wrapper to an already deployed contract.
func bindConstantSupplyERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// ConstantSupplyERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20ApprovalIterator struct {
	Event *ConstantSupplyERC20Approval // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Approval)
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
		it.Event = new(ConstantSupplyERC20Approval)
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
func (it *ConstantSupplyERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Approval represents a Approval event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ConstantSupplyERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20ApprovalIterator{contract: _ConstantSupplyERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Approval)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseApproval(log types.Log) (*ConstantSupplyERC20Approval, error) {
	event := new(ConstantSupplyERC20Approval)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConstantSupplyERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20TransferIterator struct {
	Event *ConstantSupplyERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Transfer)
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
		it.Event = new(ConstantSupplyERC20Transfer)
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
func (it *ConstantSupplyERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Transfer represents a Transfer event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConstantSupplyERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20TransferIterator{contract: _ConstantSupplyERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Transfer)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseTransfer(log types.Log) (*ConstantSupplyERC20Transfer, error) {
	event := new(ConstantSupplyERC20Transfer)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
