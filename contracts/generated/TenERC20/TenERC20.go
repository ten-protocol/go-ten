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
	Bin: "0x60806040523461002b5761001a61001461013b565b9061015d565b604051610ab361033f8239610ab390f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b90601f01601f191681019081106001600160401b0382111761006457604052565b61002f565b9061007d61007660405190565b9283610043565b565b6001600160401b03811161006457602090601f01601f19160190565b90825f9392825e0152565b909291926100bb6100b68261007f565b610069565b938185528183011161002b5761007d91602085019061009b565b9080601f8301121561002b5781516100ef926020016100a6565b90565b91909160408184031261002b5780516001600160401b03811161002b578361011b9183016100d5565b60208201519093906001600160401b03811161002b576100ef92016100d5565b610159610df28038038061014e81610069565b9283398101906100f2565b9091565b9061007d91610328565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561019b575b602083101461019657565b610167565b91607f169161018b565b6100ef6100ef6100ef9290565b91906101c36100ef6101da936101a5565b9083545f1960089290920291821b191691901b1790565b9055565b61007d915f916101b2565b8181106101f4575050565b806102015f6001936101de565b016101e9565b9190601f811161021657505050565b61022661007d935f5260205f2090565b906020601f840181900483019310610248575b6020601f9091010401906101e9565b9091508190610239565b9061025b815190565b906001600160401b0382116100645761027e82610278855461017b565b85610207565b602090601f83116001146102b7576101da92915f91836102ac575b50505f19600883021c1916906002021790565b015190505f80610299565b601f198316916102ca855f5260205f2090565b925f5b818110610306575091600293918560019694106102ee575b50505002019055565b01515f196008601f8516021c191690555f80806102e5565b919360206001819287870151815501950192016102cd565b9061007d91610252565b9061033761007d92600361031e565b600461031e56fe60806040526004361015610011575f80fd5b5f3560e01c806306fdde03146100a0578063095ea7b31461009b57806318160ddd1461009657806323b872dd14610091578063313ce5671461008c57806370a082311461008757806395d89b4114610082578063a9059cbb1461007d5763dd62ed3e036100af576102ef565b6102b0565b610295565b61027a565b610237565b61021b565b6101c6565b610198565b61010a565b5f9103126100af57565b5f80fd5b90825f9392825e0152565b6100df6100e86020936100f2936100d3815190565b80835293849260200190565b958691016100b3565b601f01601f191690565b0190565b6020808252610107929101906100be565b90565b346100af5761011a3660046100a5565b610131610125610438565b604051918291826100f6565b0390f35b6001600160a01b031690565b6001600160a01b0381165b036100af57565b9050359061016082610141565b565b8061014c565b9050359061016082610162565b91906040838203126100af576101079060206101918286610153565b9401610168565b346100af576101316101b46101ae366004610175565b90610442565b60405191829182901515815260200190565b346100af576101d63660046100a5565b6101316101e1610463565b6040515b9182918290815260200190565b90916060828403126100af5761010761020b8484610153565b9360406101918260208701610153565b346100af576101316101b46102313660046101f2565b9161046d565b346100af576102473660046100a5565b610131610252610496565b6040519182918260ff909116815260200190565b906020828203126100af5761010791610153565b346100af576101316101e1610290366004610266565b6104a0565b346100af576102a53660046100a5565b61013161012561050e565b346100af576101316101b46102c6366004610175565b90610518565b91906040838203126100af576101079060206102e88286610153565b9401610153565b346100af576101316101e16103053660046102cc565b90610523565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561033f575b602083101461033a57565b61030b565b91607f169161032f565b80545f93929161036561035b8361031f565b8085529360200190565b91600181169081156103b4575060011461037e57505050565b61038f91929394505f5260205f2090565b915f925b8184106103a05750500190565b805484840152602090930192600101610393565b92949550505060ff1916825215156020020190565b9061010791610349565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761040957604052565b6103d3565b906101606104289261041f60405190565b938480926103c9565b03836103e7565b6101079061040e565b610107600361042f565b61044d9190336105ea565b600190565b6101079081565b6101079054610452565b6101076002610459565b61044d92919061047e833383610622565b6106ad565b6104906101076101079290565b60ff1690565b6101076012610483565b6001600160a01b0381163214610505576001600160a01b03811633146105055760405162461bcd60e51b815260206004820152601f60248201527f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e6365006044820152606490fd5b610107906107a3565b610107600461042f565b61044d9190336106ad565b906001600160a01b038216321480156105d8575b6105bd576001600160a01b038216331480156105c6575b6105bd5760405162461bcd60e51b815260206004820152602160248201527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360448201527f65000000000000000000000000000000000000000000000000000000000000006064820152608490fd5b610107916107bd565b506001600160a01b038116331461054e565b506001600160a01b0381163214610537565b91600191610160936107fd565b6001600160a01b03909116815260608101939261016092909160409161061e906020830152565b0152565b9161062d8284610523565b5f19811061063c575b50505050565b818110610661579161065261065894925f940390565b916107fd565b5f808080610636565b7ffb8f41b2000000000000000000000000000000000000000000000000000000005f90815293506106939260046105f7565b035ffd5b6101356101076101079290565b61010790610697565b9291906106b95f6106a4565b936001600160a01b0385166001600160a01b0382161461072e576001600160a01b0385166001600160a01b038316146106f75761016093945061093a565b7fec442f05000000000000000000000000000000000000000000000000000000005f9081526001600160a01b038616600452602490fd5b7f96c6fd1e000000000000000000000000000000000000000000000000000000005f9081526001600160a01b038616600452602490fd5b610135610107610107926001600160a01b031690565b61010790610765565b6101079061077b565b9061079790610784565b5f5260205260405f2090565b6107b8610107916107b15f90565b505f61078d565b610459565b610107916107d76107b8926107cf5f90565b50600161078d565b61078d565b6101076101076101079290565b906101076101076107f9926107dc565b9055565b9091926108095f6106a4565b6001600160a01b0381166001600160a01b038416146108db576001600160a01b0381166001600160a01b038516146108a257506108548461084f856107d786600161078d565b6107e9565b61085d57505050565b61089d61089361088d7f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92593610784565b93610784565b936101e560405190565b0390a3565b7f94280d62000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b7fe602df05000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b634e487b7160e01b5f52601160045260245ffd5b9190820180921161093557565b610914565b6109435f6106a4565b6001600160a01b0381166001600160a01b038316036109f15761089361088d7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef936109a861089d946101356109a18a61099c6002610459565b610928565b60026107e9565b6001600160a01b038716036109d1576109cc6109a1886109c86002610459565b0390565b610784565b6109cc6109de875f61078d565b6109eb896100f283610459565b906107e9565b6109fe6107b8835f61078d565b848110610a4a5761088d7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef936109a861089d94610135610a408a610893970390565b61084f855f61078d565b610693855f92857fe450d38c000000000000000000000000000000000000000000000000000000008552600485016105f756fea2646970667358221220e96a2012161727f7ae966352331a79f79cfe60a00dde1f9190b684d0e6a30a8f64736f6c634300081c0033",
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
