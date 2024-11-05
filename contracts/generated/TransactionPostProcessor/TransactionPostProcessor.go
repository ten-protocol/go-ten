// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TransactionPostProcessor

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

// TransactionPostProcessorMetaData contains all meta data concerning the TransactionPostProcessor contract.
var TransactionPostProcessorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transactionsLength\",\"type\":\"uint256\"}],\"name\":\"TransactionsConverted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EOA_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"}],\"name\":\"addOnBlockEndCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"eoaAdmin\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"successful\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"gasUsed\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610e41806100206000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c806364c55a9d11610081578063c4d66de81161005b578063c4d66de8146101c8578063d547741f146101db578063ee546fd8146101ee57600080fd5b806364c55a9d1461017657806391d1485414610189578063a217fddf146101c057600080fd5b80632f2ff15d116100b25780632f2ff15d1461012757806336568abe1461013c578063508a50f41461014f57600080fd5b806301ffc9a7146100ce578063248a9ca3146100f7575b600080fd5b6100e16100dc36600461083b565b610265565b6040516100ee919061086e565b60405180910390f35b61011a61010536600461088d565b60009081526020819052604090206001015490565b6040516100ee91906108b4565b61013a6101353660046108e7565b6102fe565b005b61013a61014a3660046108e7565b610329565b61011a7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59881565b61013a610184366004610976565b61037a565b6100e16101973660046108e7565b6000918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b61011a600081565b61013a6101d63660046109be565b6104d0565b61013a6101e93660046108e7565b610640565b61013a6101fc3660046109be565b6001805480820182556000919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806102f857507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b60008281526020819052604090206001015461031981610665565b6103238383610672565b50505050565b6001600160a01b038116331461036b576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610375828261071c565b505050565b60006103876001306109f5565b9050336001600160a01b038216146103ba5760405162461bcd60e51b81526004016103b190610a4c565b60405180910390fd5b60008290036103db5760405162461bcd60e51b81526004016103b190610a8e565b6040517f3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece419061040b9084906108b4565b60405180910390a160005b6001548110156103235760006001828154811061043557610435610a9e565b6000918252602090912001546040517f9f9976af0000000000000000000000000000000000000000000000000000000081526001600160a01b0390911691508190639f9976af9061048c9088908890600401610d9c565b600060405180830381600087803b1580156104a657600080fd5b505af11580156104ba573d6000803e3d6000fd5b5050505050806104c990610dae565b9050610416565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff1660008115801561051b5750825b905060008267ffffffffffffffff1660011480156105385750303b155b905081158015610546575080155b1561057d576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156105b157845468ff00000000000000001916680100000000000000001785555b6105bc600087610672565b506105e77ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59887610672565b50831561063857845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061062f90600190610de2565b60405180910390a15b505050505050565b60008281526020819052604090206001015461065b81610665565b610323838361071c565b61066f813361079f565b50565b6000828152602081815260408083206001600160a01b038516845290915281205460ff16610714576000838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556106cc3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016102f8565b5060006102f8565b6000828152602081815260408083206001600160a01b038516845290915281205460ff1615610714576000838152602081815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45060016102f8565b6000828152602081815260408083206001600160a01b038516845290915290205460ff166107fd5780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016103b1929190610df0565b5050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b811461066f57600080fd5b80356102f881610801565b60006020828403121561085057610850600080fd5b600061085c8484610830565b949350505050565b8015155b82525050565b602081016102f88284610864565b80610825565b80356102f88161087c565b6000602082840312156108a2576108a2600080fd5b600061085c8484610882565b80610868565b602081016102f882846108ae565b60006001600160a01b0382166102f8565b610825816108c2565b80356102f8816108d3565b600080604083850312156108fd576108fd600080fd5b60006109098585610882565b925050602061091a858286016108dc565b9150509250929050565b60008083601f84011261093957610939600080fd5b50813567ffffffffffffffff81111561095457610954600080fd5b60208301915083602082028301111561096f5761096f600080fd5b9250929050565b6000806020838503121561098c5761098c600080fd5b823567ffffffffffffffff8111156109a6576109a6600080fd5b6109b285828601610924565b92509250509250929050565b6000602082840312156109d3576109d3600080fd5b600061085c84846108dc565b634e487b7160e01b600052601160045260246000fd5b6001600160a01b039182169190811690828203908111156102f8576102f86109df565b60088152602081017f4e6f742073656c66000000000000000000000000000000000000000000000000815290505b60200190565b602080825281016102f881610a18565b601a8152602081017f4e6f207472616e73616374696f6e7320746f20636f6e7665727400000000000081529050610a46565b602080825281016102f881610a5c565b634e487b7160e01b600052603260045260246000fd5b60ff8116610825565b80356102f881610ab4565b5060006102f86020830183610abd565b60ff8116610868565b5060006102f86020830183610882565b5060006102f860208301836108dc565b610868816108c2565b6000808335601e1936859003018112610b2557610b25600080fd5b830160208101925035905067ffffffffffffffff811115610b4857610b48600080fd5b3681900382131561096f5761096f600080fd5b82818337506000910152565b818352602083019250610b7b828483610b5b565b50601f01601f19160190565b801515610825565b80356102f881610b87565b5060006102f86020830183610b8f565b67ffffffffffffffff8116610825565b80356102f881610baa565b5060006102f86020830183610bba565b67ffffffffffffffff8116610868565b60006101408301610bf68380610ac8565b610c008582610ad8565b50610c0e6020840184610ae1565b610c1b60208601826108ae565b50610c296040840184610ae1565b610c3660408601826108ae565b50610c446060840184610ae1565b610c5160608601826108ae565b50610c5f6080840184610af1565b610c6c6080860182610b01565b50610c7a60a0840184610ae1565b610c8760a08601826108ae565b50610c9560c0840184610b0a565b85830360c0870152610ca8838284610b67565b92505050610cb960e0840184610af1565b610cc660e0860182610b01565b50610cd5610100840184610b9a565b610ce3610100860182610864565b50610cf2610120840184610bc5565b610d00610120860182610bd5565b509392505050565b6000610d148383610be5565b9392505050565b6000823561013e1936849003018112610d3657610d36600080fd5b90910192915050565b818352602083019250600083602084028101838060005b87811015610d8f578484038952610d6d8284610d1b565b610d778582610d08565b94505060208201602099909901989150600101610d56565b5091979650505050505050565b6020808252810161085c818486610d3f565b600060018201610dc057610dc06109df565b5060010190565b600067ffffffffffffffff82166102f8565b61086881610dc7565b602081016102f88284610dd9565b60408101610dfe8285610b01565b610d1460208301846108ae56fea2646970667358221220a2829afc1b81bbd1feb9c6e28aa15934fd8a00d5bce0d721a5aef549094dc19364736f6c63430008150033",
}

// TransactionPostProcessorABI is the input ABI used to generate the binding from.
// Deprecated: Use TransactionPostProcessorMetaData.ABI instead.
var TransactionPostProcessorABI = TransactionPostProcessorMetaData.ABI

// TransactionPostProcessorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransactionPostProcessorMetaData.Bin instead.
var TransactionPostProcessorBin = TransactionPostProcessorMetaData.Bin

// DeployTransactionPostProcessor deploys a new Ethereum contract, binding an instance of TransactionPostProcessor to it.
func DeployTransactionPostProcessor(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TransactionPostProcessor, error) {
	parsed, err := TransactionPostProcessorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransactionPostProcessorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransactionPostProcessor{TransactionPostProcessorCaller: TransactionPostProcessorCaller{contract: contract}, TransactionPostProcessorTransactor: TransactionPostProcessorTransactor{contract: contract}, TransactionPostProcessorFilterer: TransactionPostProcessorFilterer{contract: contract}}, nil
}

// TransactionPostProcessor is an auto generated Go binding around an Ethereum contract.
type TransactionPostProcessor struct {
	TransactionPostProcessorCaller     // Read-only binding to the contract
	TransactionPostProcessorTransactor // Write-only binding to the contract
	TransactionPostProcessorFilterer   // Log filterer for contract events
}

// TransactionPostProcessorCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransactionPostProcessorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionPostProcessorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransactionPostProcessorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionPostProcessorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransactionPostProcessorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionPostProcessorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransactionPostProcessorSession struct {
	Contract     *TransactionPostProcessor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TransactionPostProcessorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransactionPostProcessorCallerSession struct {
	Contract *TransactionPostProcessorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// TransactionPostProcessorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransactionPostProcessorTransactorSession struct {
	Contract     *TransactionPostProcessorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// TransactionPostProcessorRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransactionPostProcessorRaw struct {
	Contract *TransactionPostProcessor // Generic contract binding to access the raw methods on
}

// TransactionPostProcessorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransactionPostProcessorCallerRaw struct {
	Contract *TransactionPostProcessorCaller // Generic read-only contract binding to access the raw methods on
}

// TransactionPostProcessorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransactionPostProcessorTransactorRaw struct {
	Contract *TransactionPostProcessorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransactionPostProcessor creates a new instance of TransactionPostProcessor, bound to a specific deployed contract.
func NewTransactionPostProcessor(address common.Address, backend bind.ContractBackend) (*TransactionPostProcessor, error) {
	contract, err := bindTransactionPostProcessor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessor{TransactionPostProcessorCaller: TransactionPostProcessorCaller{contract: contract}, TransactionPostProcessorTransactor: TransactionPostProcessorTransactor{contract: contract}, TransactionPostProcessorFilterer: TransactionPostProcessorFilterer{contract: contract}}, nil
}

// NewTransactionPostProcessorCaller creates a new read-only instance of TransactionPostProcessor, bound to a specific deployed contract.
func NewTransactionPostProcessorCaller(address common.Address, caller bind.ContractCaller) (*TransactionPostProcessorCaller, error) {
	contract, err := bindTransactionPostProcessor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorCaller{contract: contract}, nil
}

// NewTransactionPostProcessorTransactor creates a new write-only instance of TransactionPostProcessor, bound to a specific deployed contract.
func NewTransactionPostProcessorTransactor(address common.Address, transactor bind.ContractTransactor) (*TransactionPostProcessorTransactor, error) {
	contract, err := bindTransactionPostProcessor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorTransactor{contract: contract}, nil
}

// NewTransactionPostProcessorFilterer creates a new log filterer instance of TransactionPostProcessor, bound to a specific deployed contract.
func NewTransactionPostProcessorFilterer(address common.Address, filterer bind.ContractFilterer) (*TransactionPostProcessorFilterer, error) {
	contract, err := bindTransactionPostProcessor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorFilterer{contract: contract}, nil
}

// bindTransactionPostProcessor binds a generic wrapper to an already deployed contract.
func bindTransactionPostProcessor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransactionPostProcessorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransactionPostProcessor *TransactionPostProcessorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransactionPostProcessor.Contract.TransactionPostProcessorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransactionPostProcessor *TransactionPostProcessorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.TransactionPostProcessorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransactionPostProcessor *TransactionPostProcessorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.TransactionPostProcessorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransactionPostProcessor *TransactionPostProcessorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransactionPostProcessor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransactionPostProcessor *TransactionPostProcessorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransactionPostProcessor *TransactionPostProcessorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionPostProcessor.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TransactionPostProcessor.Contract.DEFAULTADMINROLE(&_TransactionPostProcessor.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TransactionPostProcessor.Contract.DEFAULTADMINROLE(&_TransactionPostProcessor.CallOpts)
}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorCaller) EOAADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionPostProcessor.contract.Call(opts, &out, "EOA_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorSession) EOAADMINROLE() ([32]byte, error) {
	return _TransactionPostProcessor.Contract.EOAADMINROLE(&_TransactionPostProcessor.CallOpts)
}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorCallerSession) EOAADMINROLE() ([32]byte, error) {
	return _TransactionPostProcessor.Contract.EOAADMINROLE(&_TransactionPostProcessor.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TransactionPostProcessor.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TransactionPostProcessor.Contract.GetRoleAdmin(&_TransactionPostProcessor.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionPostProcessor *TransactionPostProcessorCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TransactionPostProcessor.Contract.GetRoleAdmin(&_TransactionPostProcessor.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionPostProcessor *TransactionPostProcessorCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TransactionPostProcessor.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionPostProcessor *TransactionPostProcessorSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TransactionPostProcessor.Contract.HasRole(&_TransactionPostProcessor.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionPostProcessor *TransactionPostProcessorCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TransactionPostProcessor.Contract.HasRole(&_TransactionPostProcessor.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionPostProcessor *TransactionPostProcessorCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TransactionPostProcessor.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionPostProcessor *TransactionPostProcessorSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TransactionPostProcessor.Contract.SupportsInterface(&_TransactionPostProcessor.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionPostProcessor *TransactionPostProcessorCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TransactionPostProcessor.Contract.SupportsInterface(&_TransactionPostProcessor.CallOpts, interfaceId)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactor) AddOnBlockEndCallback(opts *bind.TransactOpts, callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.contract.Transact(opts, "addOnBlockEndCallback", callbackAddress)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionPostProcessor *TransactionPostProcessorSession) AddOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.AddOnBlockEndCallback(&_TransactionPostProcessor.TransactOpts, callbackAddress)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactorSession) AddOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.AddOnBlockEndCallback(&_TransactionPostProcessor.TransactOpts, callbackAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionPostProcessor *TransactionPostProcessorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.GrantRole(&_TransactionPostProcessor.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.GrantRole(&_TransactionPostProcessor.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address eoaAdmin) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactor) Initialize(opts *bind.TransactOpts, eoaAdmin common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.contract.Transact(opts, "initialize", eoaAdmin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address eoaAdmin) returns()
func (_TransactionPostProcessor *TransactionPostProcessorSession) Initialize(eoaAdmin common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.Initialize(&_TransactionPostProcessor.TransactOpts, eoaAdmin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address eoaAdmin) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactorSession) Initialize(eoaAdmin common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.Initialize(&_TransactionPostProcessor.TransactOpts, eoaAdmin)
}

// OnBlock is a paid mutator transaction binding the contract method 0x64c55a9d.
//
// Solidity: function onBlock((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactor) OnBlock(opts *bind.TransactOpts, transactions []StructsTransaction) (*types.Transaction, error) {
	return _TransactionPostProcessor.contract.Transact(opts, "onBlock", transactions)
}

// OnBlock is a paid mutator transaction binding the contract method 0x64c55a9d.
//
// Solidity: function onBlock((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_TransactionPostProcessor *TransactionPostProcessorSession) OnBlock(transactions []StructsTransaction) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.OnBlock(&_TransactionPostProcessor.TransactOpts, transactions)
}

// OnBlock is a paid mutator transaction binding the contract method 0x64c55a9d.
//
// Solidity: function onBlock((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactorSession) OnBlock(transactions []StructsTransaction) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.OnBlock(&_TransactionPostProcessor.TransactOpts, transactions)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionPostProcessor *TransactionPostProcessorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.RenounceRole(&_TransactionPostProcessor.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.RenounceRole(&_TransactionPostProcessor.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionPostProcessor *TransactionPostProcessorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.RevokeRole(&_TransactionPostProcessor.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.RevokeRole(&_TransactionPostProcessor.TransactOpts, role, account)
}

// TransactionPostProcessorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TransactionPostProcessor contract.
type TransactionPostProcessorInitializedIterator struct {
	Event *TransactionPostProcessorInitialized // Event containing the contract specifics and raw log

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
func (it *TransactionPostProcessorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionPostProcessorInitialized)
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
		it.Event = new(TransactionPostProcessorInitialized)
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
func (it *TransactionPostProcessorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionPostProcessorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionPostProcessorInitialized represents a Initialized event raised by the TransactionPostProcessor contract.
type TransactionPostProcessorInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) FilterInitialized(opts *bind.FilterOpts) (*TransactionPostProcessorInitializedIterator, error) {

	logs, sub, err := _TransactionPostProcessor.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorInitializedIterator{contract: _TransactionPostProcessor.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TransactionPostProcessorInitialized) (event.Subscription, error) {

	logs, sub, err := _TransactionPostProcessor.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionPostProcessorInitialized)
				if err := _TransactionPostProcessor.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) ParseInitialized(log types.Log) (*TransactionPostProcessorInitialized, error) {
	event := new(TransactionPostProcessorInitialized)
	if err := _TransactionPostProcessor.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionPostProcessorRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TransactionPostProcessor contract.
type TransactionPostProcessorRoleAdminChangedIterator struct {
	Event *TransactionPostProcessorRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TransactionPostProcessorRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionPostProcessorRoleAdminChanged)
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
		it.Event = new(TransactionPostProcessorRoleAdminChanged)
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
func (it *TransactionPostProcessorRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionPostProcessorRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionPostProcessorRoleAdminChanged represents a RoleAdminChanged event raised by the TransactionPostProcessor contract.
type TransactionPostProcessorRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TransactionPostProcessorRoleAdminChangedIterator, error) {

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

	logs, sub, err := _TransactionPostProcessor.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorRoleAdminChangedIterator{contract: _TransactionPostProcessor.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TransactionPostProcessorRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _TransactionPostProcessor.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionPostProcessorRoleAdminChanged)
				if err := _TransactionPostProcessor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) ParseRoleAdminChanged(log types.Log) (*TransactionPostProcessorRoleAdminChanged, error) {
	event := new(TransactionPostProcessorRoleAdminChanged)
	if err := _TransactionPostProcessor.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionPostProcessorRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TransactionPostProcessor contract.
type TransactionPostProcessorRoleGrantedIterator struct {
	Event *TransactionPostProcessorRoleGranted // Event containing the contract specifics and raw log

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
func (it *TransactionPostProcessorRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionPostProcessorRoleGranted)
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
		it.Event = new(TransactionPostProcessorRoleGranted)
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
func (it *TransactionPostProcessorRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionPostProcessorRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionPostProcessorRoleGranted represents a RoleGranted event raised by the TransactionPostProcessor contract.
type TransactionPostProcessorRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TransactionPostProcessorRoleGrantedIterator, error) {

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

	logs, sub, err := _TransactionPostProcessor.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorRoleGrantedIterator{contract: _TransactionPostProcessor.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TransactionPostProcessorRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TransactionPostProcessor.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionPostProcessorRoleGranted)
				if err := _TransactionPostProcessor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) ParseRoleGranted(log types.Log) (*TransactionPostProcessorRoleGranted, error) {
	event := new(TransactionPostProcessorRoleGranted)
	if err := _TransactionPostProcessor.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionPostProcessorRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TransactionPostProcessor contract.
type TransactionPostProcessorRoleRevokedIterator struct {
	Event *TransactionPostProcessorRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TransactionPostProcessorRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionPostProcessorRoleRevoked)
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
		it.Event = new(TransactionPostProcessorRoleRevoked)
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
func (it *TransactionPostProcessorRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionPostProcessorRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionPostProcessorRoleRevoked represents a RoleRevoked event raised by the TransactionPostProcessor contract.
type TransactionPostProcessorRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TransactionPostProcessorRoleRevokedIterator, error) {

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

	logs, sub, err := _TransactionPostProcessor.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorRoleRevokedIterator{contract: _TransactionPostProcessor.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TransactionPostProcessorRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TransactionPostProcessor.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionPostProcessorRoleRevoked)
				if err := _TransactionPostProcessor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) ParseRoleRevoked(log types.Log) (*TransactionPostProcessorRoleRevoked, error) {
	event := new(TransactionPostProcessorRoleRevoked)
	if err := _TransactionPostProcessor.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionPostProcessorTransactionsConvertedIterator is returned from FilterTransactionsConverted and is used to iterate over the raw logs and unpacked data for TransactionsConverted events raised by the TransactionPostProcessor contract.
type TransactionPostProcessorTransactionsConvertedIterator struct {
	Event *TransactionPostProcessorTransactionsConverted // Event containing the contract specifics and raw log

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
func (it *TransactionPostProcessorTransactionsConvertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionPostProcessorTransactionsConverted)
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
		it.Event = new(TransactionPostProcessorTransactionsConverted)
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
func (it *TransactionPostProcessorTransactionsConvertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionPostProcessorTransactionsConvertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionPostProcessorTransactionsConverted represents a TransactionsConverted event raised by the TransactionPostProcessor contract.
type TransactionPostProcessorTransactionsConverted struct {
	TransactionsLength *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterTransactionsConverted is a free log retrieval operation binding the contract event 0x3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece41.
//
// Solidity: event TransactionsConverted(uint256 transactionsLength)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) FilterTransactionsConverted(opts *bind.FilterOpts) (*TransactionPostProcessorTransactionsConvertedIterator, error) {

	logs, sub, err := _TransactionPostProcessor.contract.FilterLogs(opts, "TransactionsConverted")
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorTransactionsConvertedIterator{contract: _TransactionPostProcessor.contract, event: "TransactionsConverted", logs: logs, sub: sub}, nil
}

// WatchTransactionsConverted is a free log subscription operation binding the contract event 0x3357352afe45ddda257f56623a512152c527b6f11555ec2fb2fdbbe72ddece41.
//
// Solidity: event TransactionsConverted(uint256 transactionsLength)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) WatchTransactionsConverted(opts *bind.WatchOpts, sink chan<- *TransactionPostProcessorTransactionsConverted) (event.Subscription, error) {

	logs, sub, err := _TransactionPostProcessor.contract.WatchLogs(opts, "TransactionsConverted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionPostProcessorTransactionsConverted)
				if err := _TransactionPostProcessor.contract.UnpackLog(event, "TransactionsConverted", log); err != nil {
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
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) ParseTransactionsConverted(log types.Log) (*TransactionPostProcessorTransactionsConverted, error) {
	event := new(TransactionPostProcessorTransactionsConverted)
	if err := _TransactionPostProcessor.contract.UnpackLog(event, "TransactionsConverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
