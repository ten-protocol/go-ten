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
	Bin: "0x6080604052346019576040516112e561001e82396112e590f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146100e0578063248a9ca3146100db5780632f2ff15d146100d657806336568abe146100d15780634d4a73c4146100cc578063508a50f4146100c75780635100f2ad146100c257806364c55a9d146100bd57806391d14854146100b8578063a217fddf146100b3578063c4d66de8146100ae578063d547741f146100a95763ee546fd80361011057610489565b610470565b610458565b61043d565b610403565b6103ea565b610378565b61032b565b6102fa565b61022a565b61020c565b610196565b61013e565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b0361011057565b5f80fd5b90503590610121826100e5565b565b906020828203126101105761013791610114565b90565b9052565b346101105761016b610159610154366004610123565b6104a1565b60405191829182901515815260200190565b0390f35b80610109565b905035906101218261016f565b906020828203126101105761013791610175565b346101105761016b6101b16101ac366004610182565b61055c565b6040519182918290815260200190565b6001600160a01b031690565b6001600160a01b038116610109565b90503590610121826101cd565b9190604083820312610110576101379060206102058286610175565b94016101dc565b346101105761022561021f3660046101e9565b9061059c565b604051005b346101105761022561023d3660046101e9565b906105a6565b634e487b7160e01b5f52603260045260245ffd5b80548210156102775761026f6001915f5260205f2090565b910201905f90565b610243565b610137916008021c6001600160a01b031690565b90610137915461027c565b600154811015610110576102b3610137916001610257565b90610290565b6101c1610137610137926001600160a01b031690565b610137906102b9565b610137906102cf565b61013a906102d8565b60208101929161012191906102e1565b346101105761016b610315610310366004610182565b61029b565b604051918291826102ea565b5f91031261011057565b346101105761033b366004610321565b61016b7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a5986101b1565b9060208282031261011057610137916101dc565b346101105761022561038b366004610364565b610845565b909182601f830112156101105781359167ffffffffffffffff831161011057602001926020830284011161011057565b9060208282031261011057813567ffffffffffffffff8111610110576103e69201610390565b9091565b34610110576102256103fd3660046103c0565b90610cfd565b346101105761016b6101596104193660046101e9565b90610d26565b6101376101376101379290565b6101375f61041f565b61013761042c565b346101105761044d366004610321565b61016b6101b1610435565b346101105761022561046b366004610364565b610f92565b34610110576102256104833660046101e9565b90610fb6565b346101105761022561049c366004610364565b611105565b7f7965db0b000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008216149081156104f1575090565b61013791507fffffffff00000000000000000000000000000000000000000000000000000000167f01ffc9a7000000000000000000000000000000000000000000000000000000001490565b905b5f5260205260405f2090565b6101379081565b610137905461054b565b60016105736101379261056c5f90565b505f61053d565b01610552565b906101219161058f61058a8261055c565b61110e565b9061059991611139565b50565b9061012191610579565b906105b0336101c1565b6001600160a01b038216036105c857610599916111b7565b7f6697b232000000000000000000000000000000000000000000000000000000005f90815260045b035ffd5b610121906106217ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59861110e565b61073e565b6101c16101376101379290565b61013790610626565b0190565b1561064757565b60405162461bcd60e51b815260206004820152601860248201527f496e76616c69642063616c6c6261636b206164647265737300000000000000006044820152606490fd5b634e487b7160e01b5f52601160045260245ffd5b919082039182116106ad57565b61068c565b919060086106d19102916106cc6001600160a01b03841b90565b921b90565b9181191691161790565b91906106ec6101376106f4936102d8565b9083546106b2565b9055565b634e487b7160e01b5f52603160045260245ffd5b610121915f916106db565b80548015610739575f1901906107366107308383610257565b9061070c565b55565b6106f8565b61075e61074d6101c15f610633565b6001600160a01b0383161415610640565b6001809161076a825490565b906107745f61041f565b935b6107c0575b60405162461bcd60e51b815260206004820152601260248201527f43616c6c6261636b206e6f7420666f756e6400000000000000000000000000006044820152606490fd5b81841015610840576107dd6107d86102b38686610257565b6102d8565b6001600160a01b03828116911614610800576107fa839460010190565b93610776565b50906101376101219361083561082e6102b361082861083b976108228861041f565b906106a0565b86610257565b9184610257565b906106db565b610717565b61077b565b610121906105f4565b6001600160a01b0390811691169003906001600160a01b0382116106ad57565b1561087557565b60405162461bcd60e51b815260206004820152600860248201527f4e6f742073656c660000000000000000000000000000000000000000000000006044820152606490fd5b90610121916108f36108ec6101c16108d76108dc6108d7306102d8565b6102cf565b6108e66001610626565b9061084e565b331461086e565b610be0565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761092e57604052565b6108f8565b60ff8116610109565b9050359061012182610933565b5061013790602081019061093c565b50610137906020810190610175565b506101379060208101906101dc565b9035601e19368390030181121561011057016020813591019167ffffffffffffffff82116101105736829003831361011057565b90825f939282370152565b91906109d3816109cc8161063c9560209181520190565b80956109aa565b601f01601f191690565b801515610109565b90503590610121826109dd565b506101379060208101906109e5565b67ffffffffffffffff8116610109565b9050359061012182610a01565b50610137906020810190610a11565b9061013790610120610b30610aeb6101408401610a54610a4d8880610949565b60ff168652565b610a6b610a646020890189610958565b6020870152565b610a82610a7b6040890189610958565b6040870152565b610a99610a926060890189610958565b6060870152565b610ab9610aa96080890189610967565b6001600160a01b03166080870152565b610ad0610ac960a0890189610958565b60a0870152565b610add60c0880188610976565b9086830360c08801526109b5565b94610b0c610afc60e0830183610967565b6001600160a01b031660e0860152565b610b27610b1d6101008301836109f2565b1515610100860152565b82810190610a1e565b67ffffffffffffffff16910152565b9061013791610a2d565b903561013e193683900301811215610110570190565b818352916020019081610b756020830284019490565b92835f925b848410610b8a5750505050505090565b9091929394956020610bb6610baf8385600195038852610baa8b88610b49565b610b3f565b9860200190565b940194019294939190610b7a565b602080825261013793910191610b5f565b6040513d5f823e3d90fd5b90919082610bf4610bf05f61041f565b9190565b14610cb857610c025f61041f565b6001610c0f610137825490565b821015610cb1576107d86102b383610c2693610257565b90813b15610110575f610c3860405190565b9283907f9f9976af000000000000000000000000000000000000000000000000000000008252818381610c6f8b8a60048401610bc4565b03925af1918215610cac57610c8a92610c8f575b5060010190565b610c02565b610ca6905f610c9e818361090c565b810190610321565b5f610c83565b610bd5565b5050509050565b60405162461bcd60e51b815260206004820152601a60248201527f4e6f207472616e73616374696f6e7320746f20636f6e766572740000000000006044820152606490fd5b90610121916108ba565b9061053f906102d8565b610137905b60ff1690565b6101379054610d11565b610137915f610d40610d4693610d395f90565b508261053d565b01610d07565b610d1c565b6101379060401c610d16565b6101379054610d4b565b610137905b67ffffffffffffffff1690565b6101379054610d61565b610d666101376101379290565b9067ffffffffffffffff906106d1565b610d666101376101379267ffffffffffffffff1690565b90610dc16101376106f492610d9a565b8254610d8a565b9068ff00000000000000009060401b6106d1565b90610dec6101376106f492151590565b8254610dc8565b61013a90610d7d565b6020810192916101219190610df3565b610e14611212565b80610e2e610e28610e2483610d57565b1590565b91610d73565b92610e385f610d7d565b67ffffffffffffffff85161480610f4f575b600194610e67610e5987610d7d565b9167ffffffffffffffff1690565b149081610f2b575b155b9081610f22575b50610ef857610ea19082610e985f610e8f88610d7d565b96019586610db1565b610ee957610f56565b610ea9575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291610ed85f610ee493610ddc565b60405191829182610dfc565b0390a1565b610ef38585610ddc565b610f56565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f610e78565b9050610e71610f39306102d8565b3b610f46610bf05f61041f565b14919050610e6f565b5081610e4a565b61059990610f6b81610f6661042c565b611139565b507ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a598611139565b61012190610e0c565b9061012191610fac61058a8261055c565b90610599916111b7565b9061012191610f9b565b61012190610fed7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59861110e565b61108c565b15610ff957565b60405162461bcd60e51b815260206004820152602360248201527f43616c6c6261636b2061646472657373206d757374206265206120636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608490fd5b908154916801000000000000000083101561092e578261083591600161012195018155610257565b610ee47f3206984d30c94bcf064cb1df53d334a1fe97a7931023e3c1ea98fa76a973cc80916110c061074d6101c15f610633565b6110d7813b6110d1610bf05f61041f565b11610ff2565b6110eb60016110e5836102d8565b90611064565b604051918291826001600160a01b03909116815260200190565b61012190610fc0565b610121903390611237565b9060ff906106d1565b906111326101376106f492151590565b8254611119565b611146610e248383610d26565b156111b157611164600161115f845f610d40868261053d565b611122565b61117e611178611172339390565b936102d8565b916102d8565b917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d6111a960405190565b5f90a4600190565b50505f90565b6111c18282610d26565b156111b1576111d95f61115f8482610d40868261053d565b6111e7611178611172339390565b917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b6111a960405190565b6101376112a7565b6001600160a01b0390911681526040810192916101219160200152565b90611245610e248284610d26565b61124d575050565b7fe2517d3f000000000000000000000000000000000000000000000000000000005f908152916105f091600461121a565b6101377ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061041f565b61013761127e56fea264697066735822122017fb968df2bd7b22cfc02b7a1cbe489ee1ed8af3d67d007a76e972686d591fb964736f6c634300081c0033",
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
