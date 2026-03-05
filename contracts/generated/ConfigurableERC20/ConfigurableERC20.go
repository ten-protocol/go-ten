// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConfigurableERC20

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

// ConfigurableERC20MetaData contains all meta data concerning the ConfigurableERC20 contract.
var ConfigurableERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b50604051610e79380380610e7983398101604081905261002e91610314565b8484600361003c8382610490565b5060046100498282610490565b50505060128360ff1611156100795760405162461bcd60e51b81526004016100709061054b565b60405180910390fd5b6005805460ff191660ff8516179055600681905561009733836100a1565b5050505050610616565b6001600160a01b0382166100ca575f60405163ec442f0560e01b815260040161007091906105a5565b6100d55f83836100d9565b5050565b6001600160a01b038316610103578060025f8282546100f891906105c7565b909155506101609050565b6001600160a01b0383165f90815260208190526040902054818110156101425783818360405163391434e360e21b8152600401610070939291906105e0565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b03821661017c5760028054829003905561019a565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516101dd9190610608565b60405180910390a3505050565b634e487b7160e01b5f52604160045260245ffd5b601f19601f83011681016001600160401b0381118282101715610223576102236101ea565b6040525050565b5f61023460405190565b905061024082826101fe565b919050565b5f6001600160401b0382111561025d5761025d6101ea565b601f19601f83011660200192915050565b8281835e505f910152565b5f61028b61028684610245565b61022a565b90508281528383830111156102a1576102a15f5ffd5b6102af83602083018461026e565b9392505050565b5f82601f8301126102c8576102c85f5ffd5b81516102d8848260208601610279565b949350505050565b60ff81165b81146102ef575f5ffd5b50565b80516102fd816102e0565b92915050565b806102e5565b80516102fd81610303565b5f5f5f5f5f60a0868803121561032b5761032b5f5ffd5b85516001600160401b03811115610343576103435f5ffd5b61034f888289016102b6565b602088015190965090506001600160401b0381111561036f5761036f5f5ffd5b61037b888289016102b6565b94505061038b87604088016102f2565b925061039a8760608801610309565b91506103a98760808801610309565b90509295509295909350565b634e487b7160e01b5f52602260045260245ffd5b6002810460018216806103dd57607f821691505b6020821081036103ef576103ef6103b5565b50919050565b5f6102fd6104008381565b90565b61040c836103f5565b81545f1960089490940293841b1916921b91909117905550565b5f610432818484610403565b505050565b818110156100d5576104495f82610426565b600101610437565b601f821115610432575f818152602090206020601f850104810160208510156104775750805b6104896020601f860104830182610437565b5050505050565b81516001600160401b038111156104a9576104a96101ea565b6104b382546103c9565b6104be828285610451565b506020601f8211600181146104f1575f83156104da5750848201515b5f19600885021c1981166002850217855550610489565b5f84815260208120601f198516915b828110156105205787850151825560209485019460019092019101610500565b508482101561053c57838701515f19601f87166008021c191681555b50505050600202600101905550565b602080825281016102fd81601681527f446563696d616c73206d757374206265203c3d20313800000000000000000000602082015260400190565b5f6001600160a01b0382166102fd565b61059f81610586565b82525050565b602081016102fd8284610596565b634e487b7160e01b5f52601160045260245ffd5b808201808211156102fd576102fd6105b3565b8061059f565b606081016105ee8286610596565b6105fb60208301856105da565b6102d860408301846105da565b602081016102fd82846105da565b610856806106235f395ff3fe608060405234801561000f575f5ffd5b506004361061009f575f3560e01c8063313ce5671161007257806395d89b411161005857806395d89b4114610140578063a9059cbb14610148578063dd62ed3e1461015b575f5ffd5b8063313ce5671461010557806370a0823114610118575f5ffd5b806306fdde03146100a3578063095ea7b3146100c157806318160ddd146100e157806323b872dd146100f2575b5f5ffd5b6100ab610193565b6040516100b8919061060f565b60405180910390f35b6100d46100cf366004610669565b610223565b6040516100b891906106a9565b6002545b6040516100b891906106bd565b6100d46101003660046106cb565b61023c565b60055460ff166040516100b8919061071a565b6100e5610126366004610728565b6001600160a01b03165f9081526020819052604090205490565b6100ab61025f565b6100d4610156366004610669565b61026e565b6100e5610169366004610745565b6001600160a01b039182165f90815260016020908152604080832093909416825291909152205490565b6060600380546101a290610786565b80601f01602080910402602001604051908101604052809291908181526020018280546101ce90610786565b80156102195780601f106101f057610100808354040283529160200191610219565b820191905f5260205f20905b8154815290600101906020018083116101fc57829003601f168201915b5050505050905090565b5f3361023081858561027b565b60019150505b92915050565b5f3361024985828561028d565b610254858585610318565b506001949350505050565b6060600480546101a290610786565b5f33610230818585610318565b61028883838360016103a7565b505050565b6001600160a01b038381165f908152600160209081526040808320938616835292905220545f198110156103125781811015610304578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016102fb939291906107bb565b60405180910390fd5b61031284848484035f6103a7565b50505050565b6001600160a01b03831661035a575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016102fb91906107eb565b6001600160a01b03821661039c575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016102fb91906107eb565b6102888383836104a9565b6001600160a01b0384166103e9575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016102fb91906107eb565b6001600160a01b03831661042b575f6040517f94280d620000000000000000000000000000000000000000000000000000000081526004016102fb91906107eb565b6001600160a01b038085165f908152600160209081526040808320938716835292905220829055801561031257826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161049b91906106bd565b60405180910390a350505050565b6001600160a01b0383166104d3578060025f8282546104c8919061080d565b909155506105499050565b6001600160a01b0383165f908152602081905260409020548181101561052b578381836040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016102fb939291906107bb565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b03821661056557600280548290039055610583565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105c691906106bd565b60405180910390a3505050565b8281835e505f910152565b5f6105e7825190565b8084526020840193506105fe8185602086016105d3565b601f01601f19169290920192915050565b6020808252810161062081846105de565b9392505050565b5f6001600160a01b038216610236565b61064081610627565b811461064a575f5ffd5b50565b803561023681610637565b80610640565b803561023681610658565b5f5f6040838503121561067d5761067d5f5ffd5b610687848461064d565b9150610696846020850161065e565b90509250929050565b8015155b82525050565b60208101610236828461069f565b806106a3565b6020810161023682846106b7565b5f5f5f606084860312156106e0576106e05f5ffd5b6106ea858561064d565b92506106f9856020860161064d565b9150610708856040860161065e565b90509250925092565b60ff81166106a3565b602081016102368284610711565b5f6020828403121561073b5761073b5f5ffd5b610620838361064d565b5f5f60408385031215610759576107595f5ffd5b610763848461064d565b9150610696846020850161064d565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061079a57607f821691505b6020821081036107ac576107ac610772565b50919050565b6106a381610627565b606081016107c982866107b2565b6107d660208301856106b7565b6107e360408301846106b7565b949350505050565b6020810161023682846107b2565b634e487b7160e01b5f52601160045260245ffd5b80820180821115610236576102366107f956fea26469706673582212204766103870d87a41d587c552041f7a898b86be40eb2f7db500e155591945d97664736f6c634300081c0033",
}

// ConfigurableERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfigurableERC20MetaData.ABI instead.
var ConfigurableERC20ABI = ConfigurableERC20MetaData.ABI

// ConfigurableERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConfigurableERC20MetaData.Bin instead.
var ConfigurableERC20Bin = ConfigurableERC20MetaData.Bin

// DeployConfigurableERC20 deploys a new Ethereum contract, binding an instance of ConfigurableERC20 to it.
func DeployConfigurableERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals_ uint8, initialSupply *big.Int, salt *big.Int) (common.Address, *types.Transaction, *ConfigurableERC20, error) {
	parsed, err := ConfigurableERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConfigurableERC20Bin), backend, name, symbol, decimals_, initialSupply, salt)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConfigurableERC20{ConfigurableERC20Caller: ConfigurableERC20Caller{contract: contract}, ConfigurableERC20Transactor: ConfigurableERC20Transactor{contract: contract}, ConfigurableERC20Filterer: ConfigurableERC20Filterer{contract: contract}}, nil
}

// ConfigurableERC20 is an auto generated Go binding around an Ethereum contract.
type ConfigurableERC20 struct {
	ConfigurableERC20Caller     // Read-only binding to the contract
	ConfigurableERC20Transactor // Write-only binding to the contract
	ConfigurableERC20Filterer   // Log filterer for contract events
}

// ConfigurableERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ConfigurableERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigurableERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfigurableERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigurableERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfigurableERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfigurableERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfigurableERC20Session struct {
	Contract     *ConfigurableERC20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ConfigurableERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfigurableERC20CallerSession struct {
	Contract *ConfigurableERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ConfigurableERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfigurableERC20TransactorSession struct {
	Contract     *ConfigurableERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ConfigurableERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ConfigurableERC20Raw struct {
	Contract *ConfigurableERC20 // Generic contract binding to access the raw methods on
}

// ConfigurableERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfigurableERC20CallerRaw struct {
	Contract *ConfigurableERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ConfigurableERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfigurableERC20TransactorRaw struct {
	Contract *ConfigurableERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewConfigurableERC20 creates a new instance of ConfigurableERC20, bound to a specific deployed contract.
func NewConfigurableERC20(address common.Address, backend bind.ContractBackend) (*ConfigurableERC20, error) {
	contract, err := bindConfigurableERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfigurableERC20{ConfigurableERC20Caller: ConfigurableERC20Caller{contract: contract}, ConfigurableERC20Transactor: ConfigurableERC20Transactor{contract: contract}, ConfigurableERC20Filterer: ConfigurableERC20Filterer{contract: contract}}, nil
}

// NewConfigurableERC20Caller creates a new read-only instance of ConfigurableERC20, bound to a specific deployed contract.
func NewConfigurableERC20Caller(address common.Address, caller bind.ContractCaller) (*ConfigurableERC20Caller, error) {
	contract, err := bindConfigurableERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigurableERC20Caller{contract: contract}, nil
}

// NewConfigurableERC20Transactor creates a new write-only instance of ConfigurableERC20, bound to a specific deployed contract.
func NewConfigurableERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ConfigurableERC20Transactor, error) {
	contract, err := bindConfigurableERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfigurableERC20Transactor{contract: contract}, nil
}

// NewConfigurableERC20Filterer creates a new log filterer instance of ConfigurableERC20, bound to a specific deployed contract.
func NewConfigurableERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ConfigurableERC20Filterer, error) {
	contract, err := bindConfigurableERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfigurableERC20Filterer{contract: contract}, nil
}

// bindConfigurableERC20 binds a generic wrapper to an already deployed contract.
func bindConfigurableERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfigurableERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfigurableERC20 *ConfigurableERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfigurableERC20.Contract.ConfigurableERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfigurableERC20 *ConfigurableERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.ConfigurableERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfigurableERC20 *ConfigurableERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.ConfigurableERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfigurableERC20 *ConfigurableERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfigurableERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfigurableERC20 *ConfigurableERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfigurableERC20 *ConfigurableERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConfigurableERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConfigurableERC20.Contract.Allowance(&_ConfigurableERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConfigurableERC20.Contract.Allowance(&_ConfigurableERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConfigurableERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConfigurableERC20.Contract.BalanceOf(&_ConfigurableERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConfigurableERC20.Contract.BalanceOf(&_ConfigurableERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConfigurableERC20 *ConfigurableERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ConfigurableERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConfigurableERC20 *ConfigurableERC20Session) Decimals() (uint8, error) {
	return _ConfigurableERC20.Contract.Decimals(&_ConfigurableERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConfigurableERC20 *ConfigurableERC20CallerSession) Decimals() (uint8, error) {
	return _ConfigurableERC20.Contract.Decimals(&_ConfigurableERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConfigurableERC20 *ConfigurableERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConfigurableERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConfigurableERC20 *ConfigurableERC20Session) Name() (string, error) {
	return _ConfigurableERC20.Contract.Name(&_ConfigurableERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConfigurableERC20 *ConfigurableERC20CallerSession) Name() (string, error) {
	return _ConfigurableERC20.Contract.Name(&_ConfigurableERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConfigurableERC20 *ConfigurableERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConfigurableERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConfigurableERC20 *ConfigurableERC20Session) Symbol() (string, error) {
	return _ConfigurableERC20.Contract.Symbol(&_ConfigurableERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConfigurableERC20 *ConfigurableERC20CallerSession) Symbol() (string, error) {
	return _ConfigurableERC20.Contract.Symbol(&_ConfigurableERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConfigurableERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20Session) TotalSupply() (*big.Int, error) {
	return _ConfigurableERC20.Contract.TotalSupply(&_ConfigurableERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConfigurableERC20 *ConfigurableERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ConfigurableERC20.Contract.TotalSupply(&_ConfigurableERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.Approve(&_ConfigurableERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.Approve(&_ConfigurableERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.Transfer(&_ConfigurableERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.Transfer(&_ConfigurableERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.TransferFrom(&_ConfigurableERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConfigurableERC20 *ConfigurableERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConfigurableERC20.Contract.TransferFrom(&_ConfigurableERC20.TransactOpts, from, to, value)
}

// ConfigurableERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ConfigurableERC20 contract.
type ConfigurableERC20ApprovalIterator struct {
	Event *ConfigurableERC20Approval // Event containing the contract specifics and raw log

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
func (it *ConfigurableERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigurableERC20Approval)
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
		it.Event = new(ConfigurableERC20Approval)
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
func (it *ConfigurableERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigurableERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigurableERC20Approval represents a Approval event raised by the ConfigurableERC20 contract.
type ConfigurableERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConfigurableERC20 *ConfigurableERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ConfigurableERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConfigurableERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ConfigurableERC20ApprovalIterator{contract: _ConfigurableERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConfigurableERC20 *ConfigurableERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ConfigurableERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConfigurableERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigurableERC20Approval)
				if err := _ConfigurableERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ConfigurableERC20 *ConfigurableERC20Filterer) ParseApproval(log types.Log) (*ConfigurableERC20Approval, error) {
	event := new(ConfigurableERC20Approval)
	if err := _ConfigurableERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfigurableERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ConfigurableERC20 contract.
type ConfigurableERC20TransferIterator struct {
	Event *ConfigurableERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ConfigurableERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfigurableERC20Transfer)
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
		it.Event = new(ConfigurableERC20Transfer)
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
func (it *ConfigurableERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfigurableERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfigurableERC20Transfer represents a Transfer event raised by the ConfigurableERC20 contract.
type ConfigurableERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConfigurableERC20 *ConfigurableERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConfigurableERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfigurableERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConfigurableERC20TransferIterator{contract: _ConfigurableERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConfigurableERC20 *ConfigurableERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ConfigurableERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConfigurableERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfigurableERC20Transfer)
				if err := _ConfigurableERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ConfigurableERC20 *ConfigurableERC20Filterer) ParseTransfer(log types.Log) (*ConfigurableERC20Transfer, error) {
	event := new(ConfigurableERC20Transfer)
	if err := _ConfigurableERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
