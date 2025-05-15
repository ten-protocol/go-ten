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
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transactionsLength\",\"type\":\"uint256\"}],\"name\":\"TransactionsConverted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EOA_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"}],\"name\":\"addOnBlockEndCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"eoaAdmin\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"successful\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"gasUsed\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"onBlockEndListeners\",\"outputs\":[{\"internalType\":\"contractIOnBlockEndCallback\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50610f398061001c5f395ff3fe608060405234801561000f575f5ffd5b50600436106100cf575f3560e01c806364c55a9d1161007d578063c4d66de811610058578063c4d66de8146101ea578063d547741f146101fd578063ee546fd814610210575f5ffd5b806364c55a9d1461019a57806391d14854146101ad578063a217fddf146101e3575f5ffd5b806336568abe116100ad57806336568abe146101405780634d4a73c414610153578063508a50f414610173575f5ffd5b806301ffc9a7146100d3578063248a9ca3146100fc5780632f2ff15d1461012b575b5f5ffd5b6100e66100e13660046108ae565b610223565b6040516100f391906108dc565b60405180910390f35b61011e61010a3660046108fb565b5f9081526020819052604090206001015490565b6040516100f3919061091e565b61013e610139366004610950565b6102bb565b005b61013e61014e366004610950565b6102e5565b6101666101613660046108fb565b610336565b6040516100f391906109a3565b61011e7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59881565b61013e6101a83660046109ff565b61035e565b6100e66101bb366004610950565b5f918252602082815260408084206001600160a01b0393909316845291905290205460ff1690565b61011e5f81565b61013e6101f8366004610a44565b61046c565b61013e61020b366004610950565b6105d9565b61013e61021e366004610a44565b6105fd565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806102b557507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f828152602081905260409020600101546102d5816106df565b6102df83836106ec565b50505050565b6001600160a01b0381163314610327576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103318282610793565b505050565b60018181548110610345575f80fd5b5f918252602090912001546001600160a01b0316905081565b5f61036a600130610a75565b9050336001600160a01b0382161461039d5760405162461bcd60e51b815260040161039490610acc565b60405180910390fd5b5f8290036103bd5760405162461bcd60e51b815260040161039490610b0e565b5f5b6001548110156102df575f600182815481106103dd576103dd610b1e565b5f918252602090912001546040517f9f9976af0000000000000000000000000000000000000000000000000000000081526001600160a01b0390911691508190639f9976af906104339088908890600401610e03565b5f604051808303815f87803b15801561044a575f5ffd5b505af115801561045c573d5f5f3e3d5ffd5b50505050508060010190506103bf565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156104b65750825b90505f8267ffffffffffffffff1660011480156104d25750303b155b9050811580156104e0575080155b15610517576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561054b57845468ff00000000000000001916680100000000000000001785555b6105555f876106ec565b506105807ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a598876106ec565b5083156105d157845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906105c890600190610e37565b60405180910390a15b505050505050565b5f828152602081905260409020600101546105f3816106df565b6102df8383610793565b7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a598610627816106df565b6001600160a01b03821661064d5760405162461bcd60e51b815260040161039490610e77565b5f826001600160a01b03163b116106765760405162461bcd60e51b815260040161039490610e87565b506001805480820182555f919091527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b6106e98133610814565b50565b5f828152602081815260408083206001600160a01b038516845290915281205460ff1661078c575f838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556107443390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016102b5565b505f6102b5565b5f828152602081815260408083206001600160a01b038516845290915281205460ff161561078c575f838152602081815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45060016102b5565b5f828152602081815260408083206001600160a01b038516845290915290205460ff166108715780826040517fe2517d3f000000000000000000000000000000000000000000000000000000008152600401610394929190610ee8565b5050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b81146106e9575f5ffd5b80356102b581610875565b5f602082840312156108c1576108c15f5ffd5b6108cb83836108a3565b9392505050565b8015155b82525050565b602081016102b582846108d2565b80610899565b80356102b5816108ea565b5f6020828403121561090e5761090e5f5ffd5b6108cb83836108f0565b806108d6565b602081016102b58284610918565b5f6001600160a01b0382166102b5565b6108998161092c565b80356102b58161093c565b5f5f60408385031215610964576109645f5ffd5b61096e84846108f0565b915061097d8460208501610945565b90509250929050565b5f6102b58261092c565b5f6102b582610986565b6108d681610990565b602081016102b5828461099a565b5f5f83601f8401126109c4576109c45f5ffd5b50813567ffffffffffffffff8111156109de576109de5f5ffd5b6020830191508360208202830111156109f8576109f85f5ffd5b9250929050565b5f5f60208385031215610a1357610a135f5ffd5b823567ffffffffffffffff811115610a2c57610a2c5f5ffd5b610a38858286016109b1565b92509250509250929050565b5f60208284031215610a5757610a575f5ffd5b6108cb8383610945565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b039182169190811690828203908111156102b5576102b5610a61565b60088152602081017f4e6f742073656c66000000000000000000000000000000000000000000000000815290505b60200190565b602080825281016102b581610a98565b601a8152602081017f4e6f207472616e73616374696f6e7320746f20636f6e7665727400000000000081529050610ac6565b602080825281016102b581610adc565b634e487b7160e01b5f52603260045260245ffd5b60ff8116610899565b80356102b581610b32565b505f6102b56020830183610b3b565b60ff81166108d6565b505f6102b560208301836108f0565b505f6102b56020830183610945565b6108d68161092c565b5f808335601e1936859003018112610b9e57610b9e5f5ffd5b830160208101925035905067ffffffffffffffff811115610bc057610bc05f5ffd5b368190038213156109f8576109f85f5ffd5b82818337505f910152565b818352602083019250610bf1828483610bd2565b50601f01601f19160190565b801515610899565b80356102b581610bfd565b505f6102b56020830183610c05565b67ffffffffffffffff8116610899565b80356102b581610c1f565b505f6102b56020830183610c2f565b67ffffffffffffffff81166108d6565b5f6101408301610c698380610b46565b610c738582610b55565b50610c816020840184610b5e565b610c8e6020860182610918565b50610c9c6040840184610b5e565b610ca96040860182610918565b50610cb76060840184610b5e565b610cc46060860182610918565b50610cd26080840184610b6d565b610cdf6080860182610b7c565b50610ced60a0840184610b5e565b610cfa60a0860182610918565b50610d0860c0840184610b85565b85830360c0870152610d1b838284610bdd565b92505050610d2c60e0840184610b6d565b610d3960e0860182610b7c565b50610d48610100840184610c10565b610d566101008601826108d2565b50610d65610120840184610c3a565b610d73610120860182610c49565b509392505050565b5f6108cb8383610c59565b5f823561013e1936849003018112610d9f57610d9f5f5ffd5b90910192915050565b8183526020830192505f8360208402810183805f5b87811015610df6578484038952610dd48284610d86565b610dde8582610d7b565b94505060208201602099909901989150600101610dbd565b5091979650505050505050565b60208082528101610e15818486610da8565b949350505050565b5f67ffffffffffffffff82166102b5565b6108d681610e1d565b602081016102b58284610e2e565b60188152602081017f496e76616c69642063616c6c6261636b2061646472657373000000000000000081529050610ac6565b602080825281016102b581610e45565b602080825281016102b581602381527f43616c6c6261636b2061646472657373206d757374206265206120636f6e747260208201527f6163740000000000000000000000000000000000000000000000000000000000604082015260600190565b60408101610ef68285610b7c565b6108cb602083018461091856fea26469706673582212200f92a464ae312815ba2cfa03e48cfaa923c4a445135b5f275e1969974137068064736f6c634300081c0033",
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

// OnBlockEndListeners is a free data retrieval call binding the contract method 0x4d4a73c4.
//
// Solidity: function onBlockEndListeners(uint256 ) view returns(address)
func (_TransactionPostProcessor *TransactionPostProcessorCaller) OnBlockEndListeners(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TransactionPostProcessor.contract.Call(opts, &out, "onBlockEndListeners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OnBlockEndListeners is a free data retrieval call binding the contract method 0x4d4a73c4.
//
// Solidity: function onBlockEndListeners(uint256 ) view returns(address)
func (_TransactionPostProcessor *TransactionPostProcessorSession) OnBlockEndListeners(arg0 *big.Int) (common.Address, error) {
	return _TransactionPostProcessor.Contract.OnBlockEndListeners(&_TransactionPostProcessor.CallOpts, arg0)
}

// OnBlockEndListeners is a free data retrieval call binding the contract method 0x4d4a73c4.
//
// Solidity: function onBlockEndListeners(uint256 ) view returns(address)
func (_TransactionPostProcessor *TransactionPostProcessorCallerSession) OnBlockEndListeners(arg0 *big.Int) (common.Address, error) {
	return _TransactionPostProcessor.Contract.OnBlockEndListeners(&_TransactionPostProcessor.CallOpts, arg0)
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
