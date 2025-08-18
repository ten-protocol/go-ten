// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ZenBase

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
	TxType     uint8
	Nonce      *big.Int
	GasPrice   *big.Int
	GasLimit   *big.Int
	To         common.Address
	Value      *big.Int
	Data       []byte
	From       common.Address
	Successful bool
	GasUsed    uint64
}

// ZenBaseMetaData contains all meta data concerning the ZenBase contract.
var ZenBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transactionPostProcessor\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransactionProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"successful\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"gasUsed\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlockEnd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461002a576100196100146100bf565b610216565b604051610e466104f28239610e4690f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b90601f01601f191681019081106001600160401b0382111761006357604052565b61002e565b9061007c61007560405190565b9283610042565b565b6001600160a01b031690565b90565b6001600160a01b0381160361002a57565b9050519061007c8261008d565b9060208282031261002a5761008a9161009e565b61008a611338803803806100d281610068565b9283398101906100ab565b6001600160401b03811161006357602090601f01601f19160190565b9061010b610106836100dd565b610068565b918252565b61011a60036100f9565b622d32b760e91b602082015290565b61008a610110565b61013b60036100f9565b622d22a760e91b602082015290565b61008a610131565b61007e61008a61008a9290565b61008a90610152565b1561016f57565b60405162461bcd60e51b8152602060048201526024808201527f496e76616c6964207472616e73616374696f6e20616e616c797a6572206164646044820152637265737360e01b6064820152608490fd5b906001600160a01b03905b9181191691161790565b61008a9061007e906001600160a01b031682565b61008a906101d5565b61008a906101e9565b9061020b61008a610212926101f2565b82546101c0565b9055565b61007c90610234610225610129565b61022d61014a565b903361025b565b61025461024361007e5f61015f565b6001600160a01b0383161415610168565b60066101fb565b9161026591610472565b61026e5f61015f565b6001600160a01b0381166001600160a01b03831614610291575061007c9061049b565b631e4fbdf760e01b5f9081526001600160a01b039091166004526024035ffd5b634e487b7160e01b5f52602260045260245ffd5b90600160028304921680156102e5575b60208310146102e057565b6102b1565b91607f16916102d5565b915f1960089290920291821b911b6101cb565b61008a61008a61008a9290565b919061032061008a61021293610302565b9083546102ef565b61007c915f9161030f565b81811061033e575050565b8061034b5f600193610328565b01610333565b9190601f811161036057505050565b61037061007c935f5260205f2090565b906020601f840181900483019310610392575b6020601f909101040190610333565b9091508190610383565b906103a5815190565b906001600160401b038211610063576103c8826103c285546102c5565b85610351565b602090601f83116001146104015761021292915f91836103f6575b50505f19600883021c1916906002021790565b015190505f806103e3565b601f19831691610414855f5260205f2090565b925f5b81811061045057509160029391856001969410610438575b50505002019055565b01515f196008601f8516021c191690555f808061042f565b91936020600181928787015181550195019201610417565b9061007c9161039c565b9061048161007c926003610468565b6004610468565b61008a9061007e565b61008a9054610488565b6104c16104bb6104ab6005610491565b6104b68460056101fb565b6101f2565b916101f2565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e06104ec60405190565b5f90a356fe60806040526004361015610011575f80fd5b5f3560e01c806306fdde03146100e0578063095ea7b3146100db57806318160ddd146100d657806323b872dd146100d1578063313ce567146100cc57806370a08231146100c7578063715018a6146100c25780638da5cb5b146100bd57806395d89b41146100b85780639f9976af146100b3578063a9059cbb146100ae578063dd62ed3e146100a95763f2fde38b036100ef57610410565b6103f4565b6103b5565b61039c565b610327565b6102f2565b6102d5565b6102ba565b610277565b61025b565b610206565b6101d8565b61014a565b5f9103126100ef57565b5f80fd5b90825f9392825e0152565b61011f61012860209361013293610113815190565b80835293849260200190565b958691016100f3565b601f01601f191690565b0190565b6020808252610147929101906100fe565b90565b346100ef5761015a3660046100e5565b610171610165610555565b60405191829182610136565b0390f35b6001600160a01b031690565b6001600160a01b0381165b036100ef57565b905035906101a082610181565b565b8061018c565b905035906101a0826101a2565b91906040838203126100ef576101479060206101d18286610193565b94016101a8565b346100ef576101716101f46101ee3660046101b5565b9061055f565b60405191829182901515815260200190565b346100ef576102163660046100e5565b610171610221610580565b6040515b9182918290815260200190565b90916060828403126100ef5761014761024b8484610193565b9360406101d18260208701610193565b346100ef576101716101f4610271366004610232565b9161058a565b346100ef576102873660046100e5565b6101716102926105b3565b6040519182918260ff909116815260200190565b906020828203126100ef5761014791610193565b346100ef576101716102216102d03660046102a6565b6105fb565b346100ef576102e53660046100e5565b6102ed61064c565b604051005b346100ef576103023660046100e5565b61017161030d610667565b604051918291826001600160a01b03909116815260200190565b346100ef576103373660046100e5565b610171610165610671565b909182601f830112156100ef5781359167ffffffffffffffff83116100ef5760200192602083028401116100ef57565b906020828203126100ef57813567ffffffffffffffff81116100ef576103989201610342565b9091565b346100ef576102ed6103af366004610372565b90610870565b346100ef576101716101f46103cb3660046101b5565b9061087a565b91906040838203126100ef576101479060206103ed8286610193565b9401610193565b346100ef5761017161022161040a3660046103d1565b90610885565b346100ef576102ed6104233660046102a6565b610916565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561045c575b602083101461045757565b610428565b91607f169161044c565b80545f9392916104826104788361043c565b8085529360200190565b91600181169081156104d1575060011461049b57505050565b6104ac91929394505f5260205f2090565b915f925b8184106104bd5750500190565b8054848401526020909301926001016104b0565b92949550505060ff1916825215156020020190565b9061014791610466565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761052657604052565b6104f0565b906101a06105459261053c60405190565b938480926104e6565b0383610504565b6101479061052b565b610147600361054c565b61056a91903361091f565b600190565b6101479081565b610147905461056f565b6101476002610576565b61056a92919061059b833383610957565b6109c8565b6105ad6101476101479290565b60ff1690565b61014760126105a0565b610175610147610147926001600160a01b031690565b610147906105bd565b610147906105d3565b906105ef906105dc565b5f5260205260405f2090565b610610610147916106095f90565b505f6105e5565b610576565b61061d610a67565b6101a061063b565b6101756101476101479290565b61014790610625565b6101a06106475f610632565b610aed565b6101a0610615565b61014790610175565b6101479054610654565b610147600561065d565b610147600461054c565b1561068257565b60405162461bcd60e51b815260206004820152602c60248201527f43616c6c65723a2063616c6c6572206973206e6f74207468652064657369676e60448201527f61746564206164647265737300000000000000000000000000000000000000006064820152608490fd5b906101a091610709610702610175600661065d565b331461067b565b6107a9565b6101476101476101479290565b634e487b7160e01b5f52603260045260245ffd5b90359061013e1936829003018212156100ef570190565b9082101561075d576020610147920281019061072f565b61071b565b3561014781610181565b634e487b7160e01b5f52601160045260245ffd5b60ff16604d811161079157600a0a90565b61076c565b8181029291811591840414171561079157565b90919082916107b75f61070e565b831461082b576107c65f61070e565b835b8110156108245761081d816108176107ef60e06107e96107c8968b8a610746565b01610762565b6108116108026107fd6105b3565b610780565b61080c600161070e565b610796565b90610b46565b60010190565b90506107c6565b5092505050565b60405162461bcd60e51b815260206004820152601a60248201527f4e6f207472616e73616374696f6e7320746f20636f6e766572740000000000006044820152606490fd5b906101a0916106ed565b61056a9190336109c8565b6101479161089f610610926108975f90565b5060016105e5565b6105e5565b6101a0906108b0610a67565b6108b95f610632565b6001600160a01b0381166001600160a01b038316146108dc57506101a090610aed565b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b0390911660045260245b035ffd5b6101a0906108a4565b916001916101a093610bb3565b6001600160a01b0390911681526060810193926101a0929091604091610953906020830152565b0152565b916109628284610885565b5f198110610971575b50505050565b818110610996579161098761098d94925f940390565b91610bb3565b5f80808061096b565b7ffb8f41b2000000000000000000000000000000000000000000000000000000005f908152935061091292600461092c565b9291906109d45f610632565b936001600160a01b0385166001600160a01b03821614610a30576001600160a01b0385166001600160a01b03831614610a12576101a0939450610cd2565b63ec442f0560e01b5f9081526001600160a01b038616600452602490fd5b7f96c6fd1e000000000000000000000000000000000000000000000000000000005f9081526001600160a01b038616600452602490fd5b610a6f610667565b339081906001600160a01b031603610a845750565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b906001600160a01b03905b9181191691161790565b90610ae2610147610ae9926105dc565b8254610abd565b9055565b610b13610b0d610afd600561065d565b610b08846005610ad2565b6105dc565b916105dc565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0610b3e60405190565b80805b0390a3565b9190610b515f610632565b926001600160a01b0384166001600160a01b03821614610b75576101a09293610cd2565b63ec442f0560e01b5f9081526001600160a01b038516600452602490fd5b905f1990610ac8565b90610bac610147610ae99261070e565b8254610b93565b909192610bbf5f610632565b6001600160a01b0381166001600160a01b03841614610c8c576001600160a01b0381166001600160a01b03851614610c535750610c0a84610c058561089f8660016105e5565b610b9c565b610c1357505050565b610b41610c49610c437f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925936105dc565b936105dc565b9361022560405190565b7f94280d62000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b7fe602df05000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b9190820180921161079157565b610cdb5f610632565b6001600160a01b0381166001600160a01b03831603610d8457610c49610c437fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef93610d40610b4194610175610d398a610d346002610576565b610cc5565b6002610b9c565b6001600160a01b03871603610d6457610b08610d3988610d606002610576565b0390565b610b08610d71875f6105e5565b610d7e8961013283610576565b90610b9c565b610d91610610835f6105e5565b848110610ddd57610c437fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef93610d40610b4194610175610dd38a610c49970390565b610c05855f6105e5565b610912855f92857fe450d38c0000000000000000000000000000000000000000000000000000000085526004850161092c56fea26469706673582212203fa28b72eda753d01773fc2d0972a245b6b25cefcfbc41ce5f27402f5c14320364736f6c634300081c0033",
}

// ZenBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenBaseMetaData.ABI instead.
var ZenBaseABI = ZenBaseMetaData.ABI

// ZenBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZenBaseMetaData.Bin instead.
var ZenBaseBin = ZenBaseMetaData.Bin

// DeployZenBase deploys a new Ethereum contract, binding an instance of ZenBase to it.
func DeployZenBase(auth *bind.TransactOpts, backend bind.ContractBackend, transactionPostProcessor common.Address) (common.Address, *types.Transaction, *ZenBase, error) {
	parsed, err := ZenBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZenBaseBin), backend, transactionPostProcessor)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ZenBase{ZenBaseCaller: ZenBaseCaller{contract: contract}, ZenBaseTransactor: ZenBaseTransactor{contract: contract}, ZenBaseFilterer: ZenBaseFilterer{contract: contract}}, nil
}

// ZenBase is an auto generated Go binding around an Ethereum contract.
type ZenBase struct {
	ZenBaseCaller     // Read-only binding to the contract
	ZenBaseTransactor // Write-only binding to the contract
	ZenBaseFilterer   // Log filterer for contract events
}

// ZenBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZenBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZenBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZenBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZenBaseSession struct {
	Contract     *ZenBase          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZenBaseCallerSession struct {
	Contract *ZenBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ZenBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZenBaseTransactorSession struct {
	Contract     *ZenBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ZenBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZenBaseRaw struct {
	Contract *ZenBase // Generic contract binding to access the raw methods on
}

// ZenBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZenBaseCallerRaw struct {
	Contract *ZenBaseCaller // Generic read-only contract binding to access the raw methods on
}

// ZenBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZenBaseTransactorRaw struct {
	Contract *ZenBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZenBase creates a new instance of ZenBase, bound to a specific deployed contract.
func NewZenBase(address common.Address, backend bind.ContractBackend) (*ZenBase, error) {
	contract, err := bindZenBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZenBase{ZenBaseCaller: ZenBaseCaller{contract: contract}, ZenBaseTransactor: ZenBaseTransactor{contract: contract}, ZenBaseFilterer: ZenBaseFilterer{contract: contract}}, nil
}

// NewZenBaseCaller creates a new read-only instance of ZenBase, bound to a specific deployed contract.
func NewZenBaseCaller(address common.Address, caller bind.ContractCaller) (*ZenBaseCaller, error) {
	contract, err := bindZenBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZenBaseCaller{contract: contract}, nil
}

// NewZenBaseTransactor creates a new write-only instance of ZenBase, bound to a specific deployed contract.
func NewZenBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*ZenBaseTransactor, error) {
	contract, err := bindZenBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZenBaseTransactor{contract: contract}, nil
}

// NewZenBaseFilterer creates a new log filterer instance of ZenBase, bound to a specific deployed contract.
func NewZenBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*ZenBaseFilterer, error) {
	contract, err := bindZenBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZenBaseFilterer{contract: contract}, nil
}

// bindZenBase binds a generic wrapper to an already deployed contract.
func bindZenBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZenBaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenBase *ZenBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenBase.Contract.ZenBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenBase *ZenBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBase.Contract.ZenBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenBase *ZenBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenBase.Contract.ZenBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenBase *ZenBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenBase *ZenBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenBase *ZenBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenBase.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBase *ZenBaseCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBase *ZenBaseSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenBase.Contract.Allowance(&_ZenBase.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBase *ZenBaseCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenBase.Contract.Allowance(&_ZenBase.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBase *ZenBaseCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBase *ZenBaseSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenBase.Contract.BalanceOf(&_ZenBase.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBase *ZenBaseCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenBase.Contract.BalanceOf(&_ZenBase.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBase *ZenBaseCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBase *ZenBaseSession) Decimals() (uint8, error) {
	return _ZenBase.Contract.Decimals(&_ZenBase.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBase *ZenBaseCallerSession) Decimals() (uint8, error) {
	return _ZenBase.Contract.Decimals(&_ZenBase.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBase *ZenBaseCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBase *ZenBaseSession) Name() (string, error) {
	return _ZenBase.Contract.Name(&_ZenBase.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBase *ZenBaseCallerSession) Name() (string, error) {
	return _ZenBase.Contract.Name(&_ZenBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenBase *ZenBaseCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenBase *ZenBaseSession) Owner() (common.Address, error) {
	return _ZenBase.Contract.Owner(&_ZenBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenBase *ZenBaseCallerSession) Owner() (common.Address, error) {
	return _ZenBase.Contract.Owner(&_ZenBase.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBase *ZenBaseCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBase *ZenBaseSession) Symbol() (string, error) {
	return _ZenBase.Contract.Symbol(&_ZenBase.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBase *ZenBaseCallerSession) Symbol() (string, error) {
	return _ZenBase.Contract.Symbol(&_ZenBase.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBase *ZenBaseCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBase *ZenBaseSession) TotalSupply() (*big.Int, error) {
	return _ZenBase.Contract.TotalSupply(&_ZenBase.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBase *ZenBaseCallerSession) TotalSupply() (*big.Int, error) {
	return _ZenBase.Contract.TotalSupply(&_ZenBase.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBase *ZenBaseSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Approve(&_ZenBase.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Approve(&_ZenBase.TransactOpts, spender, value)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x9f9976af.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_ZenBase *ZenBaseTransactor) OnBlockEnd(opts *bind.TransactOpts, transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "onBlockEnd", transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x9f9976af.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_ZenBase *ZenBaseSession) OnBlockEnd(transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenBase.Contract.OnBlockEnd(&_ZenBase.TransactOpts, transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x9f9976af.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_ZenBase *ZenBaseTransactorSession) OnBlockEnd(transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenBase.Contract.OnBlockEnd(&_ZenBase.TransactOpts, transactions)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenBase *ZenBaseTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenBase *ZenBaseSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenBase.Contract.RenounceOwnership(&_ZenBase.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenBase *ZenBaseTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenBase.Contract.RenounceOwnership(&_ZenBase.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Transfer(&_ZenBase.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Transfer(&_ZenBase.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferFrom(&_ZenBase.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferFrom(&_ZenBase.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenBase *ZenBaseTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenBase *ZenBaseSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferOwnership(&_ZenBase.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenBase *ZenBaseTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferOwnership(&_ZenBase.TransactOpts, newOwner)
}

// ZenBaseApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ZenBase contract.
type ZenBaseApprovalIterator struct {
	Event *ZenBaseApproval // Event containing the contract specifics and raw log

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
func (it *ZenBaseApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBaseApproval)
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
		it.Event = new(ZenBaseApproval)
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
func (it *ZenBaseApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBaseApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBaseApproval represents a Approval event raised by the ZenBase contract.
type ZenBaseApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenBase *ZenBaseFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZenBaseApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenBase.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZenBaseApprovalIterator{contract: _ZenBase.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenBase *ZenBaseFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZenBaseApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenBase.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBaseApproval)
				if err := _ZenBase.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ZenBase *ZenBaseFilterer) ParseApproval(log types.Log) (*ZenBaseApproval, error) {
	event := new(ZenBaseApproval)
	if err := _ZenBase.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBaseOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZenBase contract.
type ZenBaseOwnershipTransferredIterator struct {
	Event *ZenBaseOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZenBaseOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBaseOwnershipTransferred)
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
		it.Event = new(ZenBaseOwnershipTransferred)
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
func (it *ZenBaseOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBaseOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBaseOwnershipTransferred represents a OwnershipTransferred event raised by the ZenBase contract.
type ZenBaseOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenBase *ZenBaseFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZenBaseOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenBase.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZenBaseOwnershipTransferredIterator{contract: _ZenBase.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenBase *ZenBaseFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZenBaseOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenBase.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBaseOwnershipTransferred)
				if err := _ZenBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ZenBase *ZenBaseFilterer) ParseOwnershipTransferred(log types.Log) (*ZenBaseOwnershipTransferred, error) {
	event := new(ZenBaseOwnershipTransferred)
	if err := _ZenBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBaseTransactionProcessedIterator is returned from FilterTransactionProcessed and is used to iterate over the raw logs and unpacked data for TransactionProcessed events raised by the ZenBase contract.
type ZenBaseTransactionProcessedIterator struct {
	Event *ZenBaseTransactionProcessed // Event containing the contract specifics and raw log

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
func (it *ZenBaseTransactionProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBaseTransactionProcessed)
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
		it.Event = new(ZenBaseTransactionProcessed)
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
func (it *ZenBaseTransactionProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBaseTransactionProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBaseTransactionProcessed represents a TransactionProcessed event raised by the ZenBase contract.
type ZenBaseTransactionProcessed struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransactionProcessed is a free log retrieval operation binding the contract event 0xe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d1.
//
// Solidity: event TransactionProcessed(address sender, uint256 amount)
func (_ZenBase *ZenBaseFilterer) FilterTransactionProcessed(opts *bind.FilterOpts) (*ZenBaseTransactionProcessedIterator, error) {

	logs, sub, err := _ZenBase.contract.FilterLogs(opts, "TransactionProcessed")
	if err != nil {
		return nil, err
	}
	return &ZenBaseTransactionProcessedIterator{contract: _ZenBase.contract, event: "TransactionProcessed", logs: logs, sub: sub}, nil
}

// WatchTransactionProcessed is a free log subscription operation binding the contract event 0xe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d1.
//
// Solidity: event TransactionProcessed(address sender, uint256 amount)
func (_ZenBase *ZenBaseFilterer) WatchTransactionProcessed(opts *bind.WatchOpts, sink chan<- *ZenBaseTransactionProcessed) (event.Subscription, error) {

	logs, sub, err := _ZenBase.contract.WatchLogs(opts, "TransactionProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBaseTransactionProcessed)
				if err := _ZenBase.contract.UnpackLog(event, "TransactionProcessed", log); err != nil {
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

// ParseTransactionProcessed is a log parse operation binding the contract event 0xe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d1.
//
// Solidity: event TransactionProcessed(address sender, uint256 amount)
func (_ZenBase *ZenBaseFilterer) ParseTransactionProcessed(log types.Log) (*ZenBaseTransactionProcessed, error) {
	event := new(ZenBaseTransactionProcessed)
	if err := _ZenBase.contract.UnpackLog(event, "TransactionProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBaseTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ZenBase contract.
type ZenBaseTransferIterator struct {
	Event *ZenBaseTransfer // Event containing the contract specifics and raw log

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
func (it *ZenBaseTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBaseTransfer)
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
		it.Event = new(ZenBaseTransfer)
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
func (it *ZenBaseTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBaseTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBaseTransfer represents a Transfer event raised by the ZenBase contract.
type ZenBaseTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenBase *ZenBaseFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZenBaseTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenBase.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZenBaseTransferIterator{contract: _ZenBase.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenBase *ZenBaseFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZenBaseTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenBase.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBaseTransfer)
				if err := _ZenBase.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ZenBase *ZenBaseFilterer) ParseTransfer(log types.Log) (*ZenBaseTransfer, error) {
	event := new(ZenBaseTransfer)
	if err := _ZenBase.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
