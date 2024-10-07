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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transactionPostProcessor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"transactionDecoder\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransactionProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"successful\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"gasUsed\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlockEnd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001205380380620012058339810160408190526200003491620001e2565b33604051806040016040528060038152602001622d32b760e91b815250604051806040016040528060038152602001622d22a760e91b81525081600390816200007e919062000333565b5060046200008d828262000333565b5050506001600160a01b038116620000c6576000604051631e4fbdf760e01b8152600401620000bd919062000410565b60405180910390fd5b620000d18162000156565b506001600160a01b038216620000fb5760405162461bcd60e51b8152600401620000bd9062000461565b6001600160a01b038116620001245760405162461bcd60e51b8152600401620000bd90620004b2565b600680546001600160a01b039384166001600160a01b03199182161790915560078054929093169116179055620004c4565b600580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60006001600160a01b0382165b92915050565b620001c681620001a8565b8114620001d257600080fd5b50565b8051620001b581620001bb565b60008060408385031215620001fa57620001fa600080fd5b6000620002088585620001d5565b92505060206200021b85828601620001d5565b9150509250929050565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052602260045260246000fd5b6002810460018216806200026657607f821691505b6020821081036200027b576200027b6200023b565b50919050565b6000620001b56200028f8381565b90565b6200029d8362000281565b815460001960089490940293841b1916921b91909117905550565b6000620002c781848462000292565b505050565b81811015620002eb57620002e2600082620002b8565b600101620002cc565b5050565b601f821115620002c7576000818152602090206020601f85010481016020851015620003185750805b6200032c6020601f860104830182620002cc565b5050505050565b81516001600160401b038111156200034f576200034f62000225565b6200035b825462000251565b62000368828285620002ef565b506020601f821160018114620003a05760008315620003875750848201515b600019600885021c19811660028502178555506200032c565b600084815260208120601f198516915b82811015620003d25787850151825560209485019460019092019101620003b0565b5084821015620003f05783870151600019601f87166008021c191681555b50505050600202600101905550565b6200040a81620001a8565b82525050565b60208101620001b58284620003ff565b60248152602081017f496e76616c6964207472616e73616374696f6e20616e616c797a6572206164648152637265737360e01b602082015290505b60400190565b60208082528101620001b58162000420565b60238152602081017f496e76616c6964207472616e73616374696f6e206465636f646572206164647281526265737360e81b602082015290506200045b565b60208082528101620001b58162000473565b610d3180620004d46000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c8063715018a61161008c5780639f9976af116100665780639f9976af146101a9578063a9059cbb146101bc578063dd62ed3e146101cf578063f2fde38b1461020857600080fd5b8063715018a61461017e5780638da5cb5b1461018857806395d89b41146101a157600080fd5b806323b872dd116100bd57806323b872dd14610133578063313ce5671461014657806370a082311461015557600080fd5b806306fdde03146100e4578063095ea7b31461010257806318160ddd14610122575b600080fd5b6100ec61021b565b6040516100f99190610916565b60405180910390f35b610115610110366004610968565b6102ad565b6040516100f991906109af565b6002545b6040516100f991906109c3565b6101156101413660046109d1565b6102c7565b60126040516100f99190610a2a565b610126610163366004610a38565b6001600160a01b031660009081526020819052604090205490565b6101866102ed565b005b6005546001600160a01b03166040516100f99190610a6a565b6100ec610301565b6101866101b7366004610aca565b610310565b6101156101ca366004610968565b610435565b6101266101dd366004610b12565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b610186610216366004610a38565b610443565b60606003805461022a90610b5b565b80601f016020809104026020016040519081016040528092919081815260200182805461025690610b5b565b80156102a35780601f10610278576101008083540402835291602001916102a3565b820191906000526020600020905b81548152906001019060200180831161028657829003601f168201915b5050505050905090565b6000336102bb81858561049a565b60019150505b92915050565b6000336102d58582856104a7565b6102e085858561052b565b60019150505b9392505050565b6102f56105a3565b6102ff60006105e9565b565b60606004805461022a90610b5b565b6006546001600160a01b031633146103435760405162461bcd60e51b815260040161033a90610b87565b60405180910390fd5b60008190036103645760405162461bcd60e51b815260040161033a90610be8565b60005b81811015610430576103af83838381811061038457610384610c23565b90506020028101906103969190610c39565b6103a89061010081019060e001610a38565b6001610653565b7fe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d18383838181106103e2576103e2610c23565b90506020028101906103f49190610c39565b6104069061010081019060e001610a38565b6001604051610416929190610c76565b60405180910390a18061042881610ca7565b915050610367565b505050565b6000336102bb81858561052b565b61044b6105a3565b6001600160a01b03811661048e5760006040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161033a9190610a6a565b610497816105e9565b50565b610430838383600161068d565b6001600160a01b0383811660009081526001602090815260408083209386168352929052205460001981146105255781811015610516578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161033a93929190610cc0565b6105258484848403600061068d565b50505050565b6001600160a01b03831661056e5760006040517f96c6fd1e00000000000000000000000000000000000000000000000000000000815260040161033a9190610a6a565b6001600160a01b03821661059857600060405163ec442f0560e01b815260040161033a9190610a6a565b610430838383610792565b6005546001600160a01b031633146102ff57336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161033a9190610a6a565b600580546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b03821661067d57600060405163ec442f0560e01b815260040161033a9190610a6a565b61068960008383610792565b5050565b6001600160a01b0384166106d05760006040517fe602df0500000000000000000000000000000000000000000000000000000000815260040161033a9190610a6a565b6001600160a01b0383166107135760006040517f94280d6200000000000000000000000000000000000000000000000000000000815260040161033a9190610a6a565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561052557826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161078491906109c3565b60405180910390a350505050565b6001600160a01b0383166107bd5780600260008282546107b29190610ce8565b909155506108359050565b6001600160a01b03831660009081526020819052604090205481811015610816578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161033a93929190610cc0565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b03821661085157600280548290039055610870565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516108b391906109c3565b60405180910390a3505050565b60005b838110156108db5781810151838201526020016108c3565b50506000910152565b60006108ee825190565b8084526020840193506109058185602086016108c0565b601f01601f19169290920192915050565b602080825281016102e681846108e4565b60006001600160a01b0382166102c1565b61094181610927565b811461049757600080fd5b80356102c181610938565b80610941565b80356102c181610957565b6000806040838503121561097e5761097e600080fd5b600061098a858561094c565b925050602061099b8582860161095d565b9150509250929050565b8015155b82525050565b602081016102c182846109a5565b806109a9565b602081016102c182846109bd565b6000806000606084860312156109e9576109e9600080fd5b60006109f5868661094c565b9350506020610a068682870161094c565b9250506040610a178682870161095d565b9150509250925092565b60ff81166109a9565b602081016102c18284610a21565b600060208284031215610a4d57610a4d600080fd5b6000610a59848461094c565b949350505050565b6109a981610927565b602081016102c18284610a61565b60008083601f840112610a8d57610a8d600080fd5b50813567ffffffffffffffff811115610aa857610aa8600080fd5b602083019150836020820283011115610ac357610ac3600080fd5b9250929050565b60008060208385031215610ae057610ae0600080fd5b823567ffffffffffffffff811115610afa57610afa600080fd5b610b0685828601610a78565b92509250509250929050565b60008060408385031215610b2857610b28600080fd5b6000610b34858561094c565b925050602061099b8582860161094c565b634e487b7160e01b600052602260045260246000fd5b600281046001821680610b6f57607f821691505b602082108103610b8157610b81610b45565b50919050565b602080825281016102c181602c81527f43616c6c65723a2063616c6c6572206973206e6f74207468652064657369676e60208201527f6174656420616464726573730000000000000000000000000000000000000000604082015260600190565b602080825281016102c181601a81527f4e6f207472616e73616374696f6e7320746f20636f6e76657274000000000000602082015260400190565b634e487b7160e01b600052603260045260246000fd5b6000823561013e1936849003018112610c5457610c54600080fd5b9190910192915050565b60006102c1610c6a8381565b90565b6109a981610c5e565b60408101610c848285610a61565b6102e66020830184610c6d565b634e487b7160e01b600052601160045260246000fd5b600060018201610cb957610cb9610c91565b5060010190565b60608101610cce8286610a61565b610cdb60208301856109bd565b610a5960408301846109bd565b808201808211156102c1576102c1610c9156fea2646970667358221220d63c10873d835a1e945a1a99a4407630eb7252a1288c2492ce8c551ac08a1bfb64736f6c63430008140033",
}

// ZenBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenBaseMetaData.ABI instead.
var ZenBaseABI = ZenBaseMetaData.ABI

// ZenBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZenBaseMetaData.Bin instead.
var ZenBaseBin = ZenBaseMetaData.Bin

// DeployZenBase deploys a new Ethereum contract, binding an instance of ZenBase to it.
func DeployZenBase(auth *bind.TransactOpts, backend bind.ContractBackend, transactionPostProcessor common.Address, transactionDecoder common.Address) (common.Address, *types.Transaction, *ZenBase, error) {
	parsed, err := ZenBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZenBaseBin), backend, transactionPostProcessor, transactionDecoder)
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
