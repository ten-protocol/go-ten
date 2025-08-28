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
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"}],\"name\":\"CallbackAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EOA_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"}],\"name\":\"addOnBlockEndCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"eoaAdmin\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"successful\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"gasUsed\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"onBlockEndListeners\",\"outputs\":[{\"internalType\":\"contractIOnBlockEndCallback\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"}],\"name\":\"removeOnBlockEndCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506112cb8061001c5f395ff3fe608060405234801561000f575f5ffd5b50600436106100da575f3560e01c80635100f2ad11610088578063a217fddf11610063578063a217fddf14610241578063c4d66de814610248578063d547741f1461025b578063ee546fd81461026e575f5ffd5b80635100f2ad146101c457806364c55a9d146101d757806391d14854146101ea575f5ffd5b806336568abe116100b857806336568abe1461016a5780634d4a73c41461017d578063508a50f41461019d575f5ffd5b806301ffc9a7146100de578063248a9ca3146101075780632f2ff15d14610155575b5f5ffd5b6100f16100ec366004610bc9565b610281565b6040516100fe9190610bf7565b60405180910390f35b610148610115366004610c16565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b6040516100fe9190610c39565b610168610163366004610c6b565b610319565b005b610168610178366004610c6b565b610362565b61019061018b366004610c16565b6103b3565b6040516100fe9190610cbe565b6101487ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59881565b6101686101d2366004610ccc565b6103da565b6101686101e5366004610d37565b610558565b6100f16101f8366004610c6b565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b6101485f81565b610168610256366004610ccc565b61065b565b610168610269366004610c6b565b6107bb565b61016861027c366004610ccc565b6107fe565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061031357507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015461035281610907565b61035c8383610914565b50505050565b6001600160a01b03811633146103a4576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103ae82826109e0565b505050565b5f81815481106103c1575f80fd5b5f918252602090912001546001600160a01b0316905081565b7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59861040481610907565b6001600160a01b0382166104335760405162461bcd60e51b815260040161042a90610db0565b60405180910390fd5b5f8054905b8181101561053b57836001600160a01b03165f828154811061045c5761045c610dc0565b5f918252602090912001546001600160a01b031603610533575f610481600184610de8565b8154811061049157610491610dc0565b5f91825260208220015481546001600160a01b039091169190839081106104ba576104ba610dc0565b5f9182526020822001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0393909316929092179091558054806104fe576104fe610dfb565b5f8281526020902081015f19908101805473ffffffffffffffffffffffffffffffffffffffff19169055019055506105549050565b600101610438565b5060405162461bcd60e51b815260040161042a90610e41565b5050565b5f610564600130610e51565b9050336001600160a01b0382161461058e5760405162461bcd60e51b815260040161042a90610ea6565b5f8290036105ae5760405162461bcd60e51b815260040161042a90610ee8565b5f5b5f5481101561035c575f5f82815481106105cc576105cc610dc0565b5f918252602090912001546040517f9f9976af0000000000000000000000000000000000000000000000000000000081526001600160a01b0390911691508190639f9976af9061062290889088906004016111c9565b5f604051808303815f87803b158015610639575f5ffd5b505af115801561064b573d5f5f3e3d5ffd5b50505050508060010190506105b0565b5f610664610a84565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f811580156106905750825b90505f8267ffffffffffffffff1660011480156106ac5750303b155b9050811580156106ba575080155b156106f1576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561072557845468ff00000000000000001916680100000000000000001785555b61072d610aac565b6107375f87610914565b506107627ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59887610914565b5083156107b357845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906107aa906001906111fd565b60405180910390a15b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546107f481610907565b61035c83836109e0565b7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59861082881610907565b6001600160a01b03821661084e5760405162461bcd60e51b815260040161042a90610db0565b5f826001600160a01b03163b116108775760405162461bcd60e51b815260040161042a9061120b565b5f80546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0384161790556040517f3206984d30c94bcf064cb1df53d334a1fe97a7931023e3c1ea98fa76a973cc80906108fb90849061126c565b60405180910390a15050565b6109118133610ab6565b50565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff166109d7575f848152602082815260408083206001600160a01b03871684529091529020805460ff1916600117905561098d3390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a46001915050610313565b5f915050610313565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff16156109d7575f848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a46001915050610313565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610313565b610ab4610b34565b565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff166105545780826040517fe2517d3f00000000000000000000000000000000000000000000000000000000815260040161042a92919061127a565b610b3c610b72565b610ab4576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f610b7b610a84565b5468010000000000000000900460ff16919050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b8114610911575f5ffd5b803561031381610b90565b5f60208284031215610bdc57610bdc5f5ffd5b610be68383610bbe565b9392505050565b8015155b82525050565b602081016103138284610bed565b80610bb4565b803561031381610c05565b5f60208284031215610c2957610c295f5ffd5b610be68383610c0b565b80610bf1565b602081016103138284610c33565b5f6001600160a01b038216610313565b610bb481610c47565b803561031381610c57565b5f5f60408385031215610c7f57610c7f5f5ffd5b610c898484610c0b565b9150610c988460208501610c60565b90509250929050565b5f61031382610c47565b5f61031382610ca1565b610bf181610cab565b602081016103138284610cb5565b5f60208284031215610cdf57610cdf5f5ffd5b610be68383610c60565b5f5f83601f840112610cfc57610cfc5f5ffd5b50813567ffffffffffffffff811115610d1657610d165f5ffd5b602083019150836020820283011115610d3057610d305f5ffd5b9250929050565b5f5f60208385031215610d4b57610d4b5f5ffd5b823567ffffffffffffffff811115610d6457610d645f5ffd5b610d7085828601610ce9565b92509250509250929050565b60188152602081017f496e76616c69642063616c6c6261636b20616464726573730000000000000000815290505b60200190565b6020808252810161031381610d7c565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b8181038181111561031357610313610dd4565b634e487b7160e01b5f52603160045260245ffd5b60128152602081017f43616c6c6261636b206e6f7420666f756e64000000000000000000000000000081529050610daa565b6020808252810161031381610e0f565b6001600160a01b0391821691908116908282039081111561031357610313610dd4565b60088152602081017f4e6f742073656c6600000000000000000000000000000000000000000000000081529050610daa565b6020808252810161031381610e74565b601a8152602081017f4e6f207472616e73616374696f6e7320746f20636f6e7665727400000000000081529050610daa565b6020808252810161031381610eb6565b60ff8116610bb4565b803561031381610ef8565b505f6103136020830183610f01565b60ff8116610bf1565b505f6103136020830183610c0b565b505f6103136020830183610c60565b610bf181610c47565b5f808335601e1936859003018112610f6457610f645f5ffd5b830160208101925035905067ffffffffffffffff811115610f8657610f865f5ffd5b36819003821315610d3057610d305f5ffd5b82818337505f910152565b818352602083019250610fb7828483610f98565b50601f01601f19160190565b801515610bb4565b803561031381610fc3565b505f6103136020830183610fcb565b67ffffffffffffffff8116610bb4565b803561031381610fe5565b505f6103136020830183610ff5565b67ffffffffffffffff8116610bf1565b5f610140830161102f8380610f0c565b6110398582610f1b565b506110476020840184610f24565b6110546020860182610c33565b506110626040840184610f24565b61106f6040860182610c33565b5061107d6060840184610f24565b61108a6060860182610c33565b506110986080840184610f33565b6110a56080860182610f42565b506110b360a0840184610f24565b6110c060a0860182610c33565b506110ce60c0840184610f4b565b85830360c08701526110e1838284610fa3565b925050506110f260e0840184610f33565b6110ff60e0860182610f42565b5061110e610100840184610fd6565b61111c610100860182610bed565b5061112b610120840184611000565b61113961012086018261100f565b509392505050565b5f610be6838361101f565b5f823561013e1936849003018112611165576111655f5ffd5b90910192915050565b8183526020830192505f8360208402810183805f5b878110156111bc57848403895261119a828461114c565b6111a48582611141565b94505060208201602099909901989150600101611183565b5091979650505050505050565b602080825281016111db81848661116e565b949350505050565b5f67ffffffffffffffff8216610313565b610bf1816111e3565b6020810161031382846111f4565b6020808252810161031381602381527f43616c6c6261636b2061646472657373206d757374206265206120636f6e747260208201527f6163740000000000000000000000000000000000000000000000000000000000604082015260600190565b602081016103138284610f42565b604081016112888285610f42565b610be66020830184610c3356fea264697066735822122064379b19e0bc504747b3d3e6fbd250d5d3c7ae9a8d6cc5d96d3591f30e13a86a64736f6c634300081c0033",
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

// RemoveOnBlockEndCallback is a paid mutator transaction binding the contract method 0x5100f2ad.
//
// Solidity: function removeOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactor) RemoveOnBlockEndCallback(opts *bind.TransactOpts, callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.contract.Transact(opts, "removeOnBlockEndCallback", callbackAddress)
}

// RemoveOnBlockEndCallback is a paid mutator transaction binding the contract method 0x5100f2ad.
//
// Solidity: function removeOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionPostProcessor *TransactionPostProcessorSession) RemoveOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.RemoveOnBlockEndCallback(&_TransactionPostProcessor.TransactOpts, callbackAddress)
}

// RemoveOnBlockEndCallback is a paid mutator transaction binding the contract method 0x5100f2ad.
//
// Solidity: function removeOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionPostProcessor *TransactionPostProcessorTransactorSession) RemoveOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionPostProcessor.Contract.RemoveOnBlockEndCallback(&_TransactionPostProcessor.TransactOpts, callbackAddress)
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

// TransactionPostProcessorCallbackAddedIterator is returned from FilterCallbackAdded and is used to iterate over the raw logs and unpacked data for CallbackAdded events raised by the TransactionPostProcessor contract.
type TransactionPostProcessorCallbackAddedIterator struct {
	Event *TransactionPostProcessorCallbackAdded // Event containing the contract specifics and raw log

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
func (it *TransactionPostProcessorCallbackAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionPostProcessorCallbackAdded)
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
		it.Event = new(TransactionPostProcessorCallbackAdded)
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
func (it *TransactionPostProcessorCallbackAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionPostProcessorCallbackAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionPostProcessorCallbackAdded represents a CallbackAdded event raised by the TransactionPostProcessor contract.
type TransactionPostProcessorCallbackAdded struct {
	CallbackAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCallbackAdded is a free log retrieval operation binding the contract event 0x3206984d30c94bcf064cb1df53d334a1fe97a7931023e3c1ea98fa76a973cc80.
//
// Solidity: event CallbackAdded(address callbackAddress)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) FilterCallbackAdded(opts *bind.FilterOpts) (*TransactionPostProcessorCallbackAddedIterator, error) {

	logs, sub, err := _TransactionPostProcessor.contract.FilterLogs(opts, "CallbackAdded")
	if err != nil {
		return nil, err
	}
	return &TransactionPostProcessorCallbackAddedIterator{contract: _TransactionPostProcessor.contract, event: "CallbackAdded", logs: logs, sub: sub}, nil
}

// WatchCallbackAdded is a free log subscription operation binding the contract event 0x3206984d30c94bcf064cb1df53d334a1fe97a7931023e3c1ea98fa76a973cc80.
//
// Solidity: event CallbackAdded(address callbackAddress)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) WatchCallbackAdded(opts *bind.WatchOpts, sink chan<- *TransactionPostProcessorCallbackAdded) (event.Subscription, error) {

	logs, sub, err := _TransactionPostProcessor.contract.WatchLogs(opts, "CallbackAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionPostProcessorCallbackAdded)
				if err := _TransactionPostProcessor.contract.UnpackLog(event, "CallbackAdded", log); err != nil {
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

// ParseCallbackAdded is a log parse operation binding the contract event 0x3206984d30c94bcf064cb1df53d334a1fe97a7931023e3c1ea98fa76a973cc80.
//
// Solidity: event CallbackAdded(address callbackAddress)
func (_TransactionPostProcessor *TransactionPostProcessorFilterer) ParseCallbackAdded(log types.Log) (*TransactionPostProcessorCallbackAdded, error) {
	event := new(TransactionPostProcessorCallbackAdded)
	if err := _TransactionPostProcessor.contract.UnpackLog(event, "CallbackAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
