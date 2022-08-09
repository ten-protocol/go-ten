// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Guess

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

// GuessMetaData contains all meta data concerning the Guess contract.
var GuessMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_size\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"guess\",\"type\":\"uint8\"}],\"name\":\"attempt\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"close\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"guesses\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"size\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000f9438038062000f948339818101604052810190620000379190620001fc565b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600060176101000a81548160ff021916908360ff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000e3620000eb60201b60201c565b505062000309565b600060179054906101000a900460ff1660ff1642446040516020016200011392919062000272565b6040516020818303038152906040528051906020012060001c620001389190620002d1565b600060146101000a81548160ff021916908360ff160217905550565b600080fd5b600060ff82169050919050565b620001718162000159565b81146200017d57600080fd5b50565b600081519050620001918162000166565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620001c48262000197565b9050919050565b620001d681620001b7565b8114620001e257600080fd5b50565b600081519050620001f681620001cb565b92915050565b6000806040838503121562000216576200021562000154565b5b6000620002268582860162000180565b92505060206200023985828601620001e5565b9150509250929050565b6000819050919050565b6000819050919050565b6200026c620002668262000243565b6200024d565b82525050565b600062000280828562000257565b60208201915062000292828462000257565b6020820191508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000620002de8262000243565b9150620002eb8362000243565b925082620002fe57620002fd620002a2565b5b828206905092915050565b610c7b80620003196000396000f3fe6080604052600436106100705760003560e01c8063785e9e861161004e578063785e9e86146100d5578063949d225d146101005780639d76ea581461012b578063c5c884211461015657610070565b806312065fe01461007557806343d726d6146100a0578063449bf569146100aa575b600080fd5b34801561008157600080fd5b5061008a610172565b60405161009791906106f0565b60405180910390f35b6100a8610215565b005b3480156100b657600080fd5b506100bf6102dc565b6040516100cc9190610728565b60405180910390f35b3480156100e157600080fd5b506100ea6102f0565b6040516100f791906107c2565b60405180910390f35b34801561010c57600080fd5b50610115610316565b60405161012291906107f9565b60405180910390f35b34801561013757600080fd5b50610140610329565b60405161014d9190610835565b60405180910390f35b610170600480360381019061016b9190610881565b61034f565b005b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016101cf9190610835565b602060405180830381865afa1580156101ec573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021091906108da565b905090565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029a9061098a565b60405180910390fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b600060159054906101000a900461ffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600060179054906101000a900460ff1681565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60018060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dd62ed3e33306040518363ffffffff1660e01b81526004016103ad9291906109aa565b602060405180830381865afa1580156103ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ee91906108da565b101561042f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161042690610a1f565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd333060016040518463ffffffff1660e01b815260040161048f93929190610a7a565b6020604051808303816000875af11580156104ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d29190610ae9565b506000601581819054906101000a900461ffff16809291906104f390610b45565b91906101000a81548161ffff021916908361ffff16021790555050600060149054906101000a900460ff1660ff168160ff160361066f57600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016105c49190610835565b602060405180830381865afa1580156105e1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061060591906108da565b6040518363ffffffff1660e01b8152600401610622929190610b6f565b6020604051808303816000875af1158015610641573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106659190610ae9565b5061066e610672565b5b50565b600060179054906101000a900460ff1660ff164244604051602001610698929190610bb9565b6040516020818303038152906040528051906020012060001c6106bb9190610c14565b600060146101000a81548160ff021916908360ff160217905550565b6000819050919050565b6106ea816106d7565b82525050565b600060208201905061070560008301846106e1565b92915050565b600061ffff82169050919050565b6107228161070b565b82525050565b600060208201905061073d6000830184610719565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600061078861078361077e84610743565b610763565b610743565b9050919050565b600061079a8261076d565b9050919050565b60006107ac8261078f565b9050919050565b6107bc816107a1565b82525050565b60006020820190506107d760008301846107b3565b92915050565b600060ff82169050919050565b6107f3816107dd565b82525050565b600060208201905061080e60008301846107ea565b92915050565b600061081f82610743565b9050919050565b61082f81610814565b82525050565b600060208201905061084a6000830184610826565b92915050565b600080fd5b61085e816107dd565b811461086957600080fd5b50565b60008135905061087b81610855565b92915050565b60006020828403121561089757610896610850565b5b60006108a58482850161086c565b91505092915050565b6108b7816106d7565b81146108c257600080fd5b50565b6000815190506108d4816108ae565b92915050565b6000602082840312156108f0576108ef610850565b5b60006108fe848285016108c5565b91505092915050565b600082825260208201905092915050565b7f4f6e6c79206f776e65722063616e2063616c6c20746869732066756e6374696f60008201527f6e2e000000000000000000000000000000000000000000000000000000000000602082015250565b6000610974602283610907565b915061097f82610918565b604082019050919050565b600060208201905081810360008301526109a381610967565b9050919050565b60006040820190506109bf6000830185610826565b6109cc6020830184610826565b9392505050565b7f436865636b2074686520746f6b656e20616c6c6f77616e63652e000000000000600082015250565b6000610a09601a83610907565b9150610a14826109d3565b602082019050919050565b60006020820190508181036000830152610a38816109fc565b9050919050565b6000819050919050565b6000610a64610a5f610a5a84610a3f565b610763565b6106d7565b9050919050565b610a7481610a49565b82525050565b6000606082019050610a8f6000830186610826565b610a9c6020830185610826565b610aa96040830184610a6b565b949350505050565b60008115159050919050565b610ac681610ab1565b8114610ad157600080fd5b50565b600081519050610ae381610abd565b92915050565b600060208284031215610aff57610afe610850565b5b6000610b0d84828501610ad4565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610b508261070b565b915061ffff8203610b6457610b63610b16565b5b600182019050919050565b6000604082019050610b846000830185610826565b610b9160208301846106e1565b9392505050565b6000819050919050565b610bb3610bae826106d7565b610b98565b82525050565b6000610bc58285610ba2565b602082019150610bd58284610ba2565b6020820191508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000610c1f826106d7565b9150610c2a836106d7565b925082610c3a57610c39610be5565b5b82820690509291505056fea264697066735822122008eb75450fe17cc6c5614094920d47cf43ad17046171800728c1659b5c67abd064736f6c634300080f0033",
}

// GuessABI is the input ABI used to generate the binding from.
// Deprecated: Use GuessMetaData.ABI instead.
var GuessABI = GuessMetaData.ABI

// GuessBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GuessMetaData.Bin instead.
var GuessBin = GuessMetaData.Bin

// DeployGuess deploys a new Ethereum contract, binding an instance of Guess to it.
func DeployGuess(auth *bind.TransactOpts, backend bind.ContractBackend, _size uint8, _tokenAddress common.Address) (common.Address, *types.Transaction, *Guess, error) {
	parsed, err := GuessMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GuessBin), backend, _size, _tokenAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Guess{GuessCaller: GuessCaller{contract: contract}, GuessTransactor: GuessTransactor{contract: contract}, GuessFilterer: GuessFilterer{contract: contract}}, nil
}

// Guess is an auto generated Go binding around an Ethereum contract.
type Guess struct {
	GuessCaller     // Read-only binding to the contract
	GuessTransactor // Write-only binding to the contract
	GuessFilterer   // Log filterer for contract events
}

// GuessCaller is an auto generated read-only Go binding around an Ethereum contract.
type GuessCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GuessTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GuessTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GuessFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GuessFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GuessSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GuessSession struct {
	Contract     *Guess            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GuessCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GuessCallerSession struct {
	Contract *GuessCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// GuessTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GuessTransactorSession struct {
	Contract     *GuessTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GuessRaw is an auto generated low-level Go binding around an Ethereum contract.
type GuessRaw struct {
	Contract *Guess // Generic contract binding to access the raw methods on
}

// GuessCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GuessCallerRaw struct {
	Contract *GuessCaller // Generic read-only contract binding to access the raw methods on
}

// GuessTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GuessTransactorRaw struct {
	Contract *GuessTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGuess creates a new instance of Guess, bound to a specific deployed contract.
func NewGuess(address common.Address, backend bind.ContractBackend) (*Guess, error) {
	contract, err := bindGuess(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Guess{GuessCaller: GuessCaller{contract: contract}, GuessTransactor: GuessTransactor{contract: contract}, GuessFilterer: GuessFilterer{contract: contract}}, nil
}

// NewGuessCaller creates a new read-only instance of Guess, bound to a specific deployed contract.
func NewGuessCaller(address common.Address, caller bind.ContractCaller) (*GuessCaller, error) {
	contract, err := bindGuess(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GuessCaller{contract: contract}, nil
}

// NewGuessTransactor creates a new write-only instance of Guess, bound to a specific deployed contract.
func NewGuessTransactor(address common.Address, transactor bind.ContractTransactor) (*GuessTransactor, error) {
	contract, err := bindGuess(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GuessTransactor{contract: contract}, nil
}

// NewGuessFilterer creates a new log filterer instance of Guess, bound to a specific deployed contract.
func NewGuessFilterer(address common.Address, filterer bind.ContractFilterer) (*GuessFilterer, error) {
	contract, err := bindGuess(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GuessFilterer{contract: contract}, nil
}

// bindGuess binds a generic wrapper to an already deployed contract.
func bindGuess(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GuessABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Guess *GuessRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Guess.Contract.GuessCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Guess *GuessRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Guess.Contract.GuessTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Guess *GuessRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Guess.Contract.GuessTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Guess *GuessCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Guess.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Guess *GuessTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Guess.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Guess *GuessTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Guess.Contract.contract.Transact(opts, method, params...)
}

// Erc20 is a free data retrieval call binding the contract method 0x785e9e86.
//
// Solidity: function erc20() view returns(address)
func (_Guess *GuessCaller) Erc20(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Guess.contract.Call(opts, &out, "erc20")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20 is a free data retrieval call binding the contract method 0x785e9e86.
//
// Solidity: function erc20() view returns(address)
func (_Guess *GuessSession) Erc20() (common.Address, error) {
	return _Guess.Contract.Erc20(&_Guess.CallOpts)
}

// Erc20 is a free data retrieval call binding the contract method 0x785e9e86.
//
// Solidity: function erc20() view returns(address)
func (_Guess *GuessCallerSession) Erc20() (common.Address, error) {
	return _Guess.Contract.Erc20(&_Guess.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_Guess *GuessCaller) GetBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Guess.contract.Call(opts, &out, "getBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_Guess *GuessSession) GetBalance() (*big.Int, error) {
	return _Guess.Contract.GetBalance(&_Guess.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_Guess *GuessCallerSession) GetBalance() (*big.Int, error) {
	return _Guess.Contract.GetBalance(&_Guess.CallOpts)
}

// Guesses is a free data retrieval call binding the contract method 0x449bf569.
//
// Solidity: function guesses() view returns(uint16)
func (_Guess *GuessCaller) Guesses(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Guess.contract.Call(opts, &out, "guesses")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// Guesses is a free data retrieval call binding the contract method 0x449bf569.
//
// Solidity: function guesses() view returns(uint16)
func (_Guess *GuessSession) Guesses() (uint16, error) {
	return _Guess.Contract.Guesses(&_Guess.CallOpts)
}

// Guesses is a free data retrieval call binding the contract method 0x449bf569.
//
// Solidity: function guesses() view returns(uint16)
func (_Guess *GuessCallerSession) Guesses() (uint16, error) {
	return _Guess.Contract.Guesses(&_Guess.CallOpts)
}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() view returns(uint8)
func (_Guess *GuessCaller) Size(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Guess.contract.Call(opts, &out, "size")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() view returns(uint8)
func (_Guess *GuessSession) Size() (uint8, error) {
	return _Guess.Contract.Size(&_Guess.CallOpts)
}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() view returns(uint8)
func (_Guess *GuessCallerSession) Size() (uint8, error) {
	return _Guess.Contract.Size(&_Guess.CallOpts)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() view returns(address)
func (_Guess *GuessCaller) TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Guess.contract.Call(opts, &out, "tokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() view returns(address)
func (_Guess *GuessSession) TokenAddress() (common.Address, error) {
	return _Guess.Contract.TokenAddress(&_Guess.CallOpts)
}

// TokenAddress is a free data retrieval call binding the contract method 0x9d76ea58.
//
// Solidity: function tokenAddress() view returns(address)
func (_Guess *GuessCallerSession) TokenAddress() (common.Address, error) {
	return _Guess.Contract.TokenAddress(&_Guess.CallOpts)
}

// Attempt is a paid mutator transaction binding the contract method 0xc5c88421.
//
// Solidity: function attempt(uint8 guess) payable returns()
func (_Guess *GuessTransactor) Attempt(opts *bind.TransactOpts, guess uint8) (*types.Transaction, error) {
	return _Guess.contract.Transact(opts, "attempt", guess)
}

// Attempt is a paid mutator transaction binding the contract method 0xc5c88421.
//
// Solidity: function attempt(uint8 guess) payable returns()
func (_Guess *GuessSession) Attempt(guess uint8) (*types.Transaction, error) {
	return _Guess.Contract.Attempt(&_Guess.TransactOpts, guess)
}

// Attempt is a paid mutator transaction binding the contract method 0xc5c88421.
//
// Solidity: function attempt(uint8 guess) payable returns()
func (_Guess *GuessTransactorSession) Attempt(guess uint8) (*types.Transaction, error) {
	return _Guess.Contract.Attempt(&_Guess.TransactOpts, guess)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() payable returns()
func (_Guess *GuessTransactor) Close(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Guess.contract.Transact(opts, "close")
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() payable returns()
func (_Guess *GuessSession) Close() (*types.Transaction, error) {
	return _Guess.Contract.Close(&_Guess.TransactOpts)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() payable returns()
func (_Guess *GuessTransactorSession) Close() (*types.Transaction, error) {
	return _Guess.Contract.Close(&_Guess.TransactOpts)
}
