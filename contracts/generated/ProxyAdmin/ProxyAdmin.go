// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ProxyAdmin

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

// ProxyAdminMetaData contains all meta data concerning the ProxyAdmin contract.
var ProxyAdminMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITransparentUpgradeableProxy\",\"name\":\"proxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60806040523462000030576200001e62000018620000d3565b620000f6565b60405161066362000265823961066390f35b600080fd5b634e487b7160e01b600052604160045260246000fd5b90601f01601f191681019081106001600160401b038211176200006d57604052565b62000035565b906200008a6200008260405190565b92836200004b565b565b6001600160a01b031690565b90565b6001600160a01b038116036200003057565b905051906200008a826200009b565b9060208282031262000030576200009891620000ad565b62000098620008c880380380620000ea8162000073565b928339810190620000bc565b6200008a906200012f565b6200008c62000098620000989290565b620000989062000101565b6001600160a01b03909116815260200190565b6200013b600062000111565b6001600160a01b0381166001600160a01b038316146200016157506200008a9062000202565b62000187906200017060405190565b631e4fbdf760e01b8152918291600483016200011c565b0390fd5b62000098906200008c565b6200009890546200018b565b62000098906200008c906001600160a01b031682565b6200009890620001a2565b6200009890620001b8565b90620001e262000098620001fe92620001c3565b82546001600160a01b0319166001600160a01b03919091161790565b9055565b6200020e600062000196565b906200021c816000620001ce565b620002536200024c7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093620001c3565b91620001c3565b916200025e60405190565b600090a356fe6080604052600436101561001257600080fd5b60003560e01c8063715018a6146100625780638da5cb5b1461005d5780639623609d14610058578063ad3cb1cc146100535763f2fde38b036100725761035a565b61031f565b610246565b6100a9565b610077565b600091031261007257565b600080fd5b3461007257610087366004610067565b61008f6103aa565b604051005b0390f35b6001600160a01b031690565b90565b565b34610072576100b9366004610067565b6100946100ce6000546001600160a01b031690565b604051918291826001600160a01b03909116815260200190565b6001600160a01b0381165b0361007257565b905035906100a7826100e8565b6001600160a01b0381166100f3565b905035906100a782610107565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff82111761017457604052565b610123565b906100a761018660405190565b9283610152565b67ffffffffffffffff811161017457602090601f01601f19160190565b0190565b90826000939282370152565b909291926101cf6101ca8261018d565b610179565b93818552602085019082840111610072576100a7926101ae565b9080601f83011215610072578160206100a4933591016101ba565b916060838303126100725761021982846100fa565b926102278360208301610116565b92604082013567ffffffffffffffff8111610072576100a492016101e9565b61008f610254366004610204565b91610622565b906102676101ca8361018d565b918252565b610276600561025a565b7f352e302e30000000000000000000000000000000000000000000000000000000602082015290565b6100a461026c565b6100a461029f565b6100a46102a7565b60005b8381106102ca5750506000910152565b81810151838201526020016102ba565b6102fb6103046020936101aa936102ef815190565b80835293849260200190565b958691016102b7565b601f01601f191690565b60208082526100a4929101906102da565b346100725761032f366004610067565b61009461033a6102af565b6040519182918261030e565b90602082820312610072576100a491610116565b346100725761008f61036d366004610346565b61049b565b61037a6103b2565b6100a7610398565b6100986100a46100a49290565b6100a490610382565b6100a76103a5600061038f565b610514565b6100a7610372565b60005433906001600160a01b03168190036103ca5750565b610414906103d760405190565b9182917f118cdaa7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b0390fd5b6100a7906104246103b2565b61042e600061038f565b6001600160a01b0381166001600160a01b0383161461045157506100a790610514565b6104149061045e60405190565b9182917f1e4fbdf7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6100a790610418565b6100986100a46100a4926001600160a01b031690565b6100a4906104a4565b6100a4906104ba565b906104dc6100a4610510926104c3565b82547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03919091161790565b9055565b6000546001600160a01b03169061052c8160006104cc565b61055f6105597f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936104c3565b916104c3565b9161056960405190565b600090a3565b906100a7929161057d6103b2565b6105b1565b6001600160a01b0390911681526100a491604082019160208184039101526102da565b6040513d6000823e3d90fd5b6105ba906104c3565b349190634f1ef286813b15610072576000936105e9916105f46105dc60405190565b9788968795869460e01b90565b845260048401610582565b03925af1801561061d576106055750565b6100a79060006106158183610152565b810190610067565b6105a5565b906100a7929161056f56fea26469706673582212205d5f75a536a5fd2f2e441550c1d5a4e715b6cff79e1cc29a86c7abd4b441f65364736f6c63430008140033",
}

// ProxyAdminABI is the input ABI used to generate the binding from.
// Deprecated: Use ProxyAdminMetaData.ABI instead.
var ProxyAdminABI = ProxyAdminMetaData.ABI

// ProxyAdminBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProxyAdminMetaData.Bin instead.
var ProxyAdminBin = ProxyAdminMetaData.Bin

// DeployProxyAdmin deploys a new Ethereum contract, binding an instance of ProxyAdmin to it.
func DeployProxyAdmin(auth *bind.TransactOpts, backend bind.ContractBackend, initialOwner common.Address) (common.Address, *types.Transaction, *ProxyAdmin, error) {
	parsed, err := ProxyAdminMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProxyAdminBin), backend, initialOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ProxyAdmin{ProxyAdminCaller: ProxyAdminCaller{contract: contract}, ProxyAdminTransactor: ProxyAdminTransactor{contract: contract}, ProxyAdminFilterer: ProxyAdminFilterer{contract: contract}}, nil
}

// ProxyAdmin is an auto generated Go binding around an Ethereum contract.
type ProxyAdmin struct {
	ProxyAdminCaller     // Read-only binding to the contract
	ProxyAdminTransactor // Write-only binding to the contract
	ProxyAdminFilterer   // Log filterer for contract events
}

// ProxyAdminCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProxyAdminCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyAdminTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProxyAdminTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyAdminFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProxyAdminFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyAdminSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProxyAdminSession struct {
	Contract     *ProxyAdmin       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyAdminCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProxyAdminCallerSession struct {
	Contract *ProxyAdminCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ProxyAdminTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProxyAdminTransactorSession struct {
	Contract     *ProxyAdminTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ProxyAdminRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProxyAdminRaw struct {
	Contract *ProxyAdmin // Generic contract binding to access the raw methods on
}

// ProxyAdminCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProxyAdminCallerRaw struct {
	Contract *ProxyAdminCaller // Generic read-only contract binding to access the raw methods on
}

// ProxyAdminTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProxyAdminTransactorRaw struct {
	Contract *ProxyAdminTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProxyAdmin creates a new instance of ProxyAdmin, bound to a specific deployed contract.
func NewProxyAdmin(address common.Address, backend bind.ContractBackend) (*ProxyAdmin, error) {
	contract, err := bindProxyAdmin(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProxyAdmin{ProxyAdminCaller: ProxyAdminCaller{contract: contract}, ProxyAdminTransactor: ProxyAdminTransactor{contract: contract}, ProxyAdminFilterer: ProxyAdminFilterer{contract: contract}}, nil
}

// NewProxyAdminCaller creates a new read-only instance of ProxyAdmin, bound to a specific deployed contract.
func NewProxyAdminCaller(address common.Address, caller bind.ContractCaller) (*ProxyAdminCaller, error) {
	contract, err := bindProxyAdmin(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyAdminCaller{contract: contract}, nil
}

// NewProxyAdminTransactor creates a new write-only instance of ProxyAdmin, bound to a specific deployed contract.
func NewProxyAdminTransactor(address common.Address, transactor bind.ContractTransactor) (*ProxyAdminTransactor, error) {
	contract, err := bindProxyAdmin(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyAdminTransactor{contract: contract}, nil
}

// NewProxyAdminFilterer creates a new log filterer instance of ProxyAdmin, bound to a specific deployed contract.
func NewProxyAdminFilterer(address common.Address, filterer bind.ContractFilterer) (*ProxyAdminFilterer, error) {
	contract, err := bindProxyAdmin(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProxyAdminFilterer{contract: contract}, nil
}

// bindProxyAdmin binds a generic wrapper to an already deployed contract.
func bindProxyAdmin(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProxyAdminMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxyAdmin *ProxyAdminRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProxyAdmin.Contract.ProxyAdminCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxyAdmin *ProxyAdminRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.ProxyAdminTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxyAdmin *ProxyAdminRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.ProxyAdminTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProxyAdmin *ProxyAdminCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProxyAdmin.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProxyAdmin *ProxyAdminTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProxyAdmin *ProxyAdminTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ProxyAdmin *ProxyAdminCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ProxyAdmin.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ProxyAdmin *ProxyAdminSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ProxyAdmin.Contract.UPGRADEINTERFACEVERSION(&_ProxyAdmin.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ProxyAdmin *ProxyAdminCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ProxyAdmin.Contract.UPGRADEINTERFACEVERSION(&_ProxyAdmin.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProxyAdmin *ProxyAdminCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProxyAdmin.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProxyAdmin *ProxyAdminSession) Owner() (common.Address, error) {
	return _ProxyAdmin.Contract.Owner(&_ProxyAdmin.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProxyAdmin *ProxyAdminCallerSession) Owner() (common.Address, error) {
	return _ProxyAdmin.Contract.Owner(&_ProxyAdmin.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ProxyAdmin *ProxyAdminTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProxyAdmin.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ProxyAdmin *ProxyAdminSession) RenounceOwnership() (*types.Transaction, error) {
	return _ProxyAdmin.Contract.RenounceOwnership(&_ProxyAdmin.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ProxyAdmin *ProxyAdminTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ProxyAdmin.Contract.RenounceOwnership(&_ProxyAdmin.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ProxyAdmin *ProxyAdminTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ProxyAdmin.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ProxyAdmin *ProxyAdminSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.TransferOwnership(&_ProxyAdmin.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ProxyAdmin *ProxyAdminTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.TransferOwnership(&_ProxyAdmin.TransactOpts, newOwner)
}

// UpgradeAndCall is a paid mutator transaction binding the contract method 0x9623609d.
//
// Solidity: function upgradeAndCall(address proxy, address implementation, bytes data) payable returns()
func (_ProxyAdmin *ProxyAdminTransactor) UpgradeAndCall(opts *bind.TransactOpts, proxy common.Address, implementation common.Address, data []byte) (*types.Transaction, error) {
	return _ProxyAdmin.contract.Transact(opts, "upgradeAndCall", proxy, implementation, data)
}

// UpgradeAndCall is a paid mutator transaction binding the contract method 0x9623609d.
//
// Solidity: function upgradeAndCall(address proxy, address implementation, bytes data) payable returns()
func (_ProxyAdmin *ProxyAdminSession) UpgradeAndCall(proxy common.Address, implementation common.Address, data []byte) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.UpgradeAndCall(&_ProxyAdmin.TransactOpts, proxy, implementation, data)
}

// UpgradeAndCall is a paid mutator transaction binding the contract method 0x9623609d.
//
// Solidity: function upgradeAndCall(address proxy, address implementation, bytes data) payable returns()
func (_ProxyAdmin *ProxyAdminTransactorSession) UpgradeAndCall(proxy common.Address, implementation common.Address, data []byte) (*types.Transaction, error) {
	return _ProxyAdmin.Contract.UpgradeAndCall(&_ProxyAdmin.TransactOpts, proxy, implementation, data)
}

// ProxyAdminOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ProxyAdmin contract.
type ProxyAdminOwnershipTransferredIterator struct {
	Event *ProxyAdminOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ProxyAdminOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProxyAdminOwnershipTransferred)
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
		it.Event = new(ProxyAdminOwnershipTransferred)
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
func (it *ProxyAdminOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProxyAdminOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProxyAdminOwnershipTransferred represents a OwnershipTransferred event raised by the ProxyAdmin contract.
type ProxyAdminOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ProxyAdmin *ProxyAdminFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ProxyAdminOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ProxyAdmin.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ProxyAdminOwnershipTransferredIterator{contract: _ProxyAdmin.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ProxyAdmin *ProxyAdminFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ProxyAdminOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ProxyAdmin.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProxyAdminOwnershipTransferred)
				if err := _ProxyAdmin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ProxyAdmin *ProxyAdminFilterer) ParseOwnershipTransferred(log types.Log) (*ProxyAdminOwnershipTransferred, error) {
	event := new(ProxyAdminOwnershipTransferred)
	if err := _ProxyAdmin.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
