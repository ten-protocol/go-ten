// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedManagementContract

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

// ManagementContractMetaRollup is an auto generated low-level Go binding around an user-defined struct.
type ManagementContractMetaRollup struct {
	ParentHash   [32]byte
	Hash         [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}

// ManagementContractTreeElement is an auto generated low-level Go binding around an user-defined struct.
type ManagementContractTreeElement struct {
	ElementID *big.Int
	ParentID  *big.Int
	Rollup    ManagementContractMetaRollup
}

// ECDSAMetaData contains all meta data concerning the ECDSA contract.
var ECDSAMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220744e00a84dcdd38eb90023f9a414401959a8ad7256c8c6e2f82c723706df170b64736f6c634300080c0033",
}

// ECDSAABI is the input ABI used to generate the binding from.
// Deprecated: Use ECDSAMetaData.ABI instead.
var ECDSAABI = ECDSAMetaData.ABI

// ECDSABin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ECDSAMetaData.Bin instead.
var ECDSABin = ECDSAMetaData.Bin

// DeployECDSA deploys a new Ethereum contract, binding an instance of ECDSA to it.
func DeployECDSA(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ECDSA, error) {
	parsed, err := ECDSAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ECDSABin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ECDSA{ECDSACaller: ECDSACaller{contract: contract}, ECDSATransactor: ECDSATransactor{contract: contract}, ECDSAFilterer: ECDSAFilterer{contract: contract}}, nil
}

// ECDSA is an auto generated Go binding around an Ethereum contract.
type ECDSA struct {
	ECDSACaller     // Read-only binding to the contract
	ECDSATransactor // Write-only binding to the contract
	ECDSAFilterer   // Log filterer for contract events
}

// ECDSACaller is an auto generated read-only Go binding around an Ethereum contract.
type ECDSACaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSATransactor is an auto generated write-only Go binding around an Ethereum contract.
type ECDSATransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSAFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ECDSAFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSASession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ECDSASession struct {
	Contract     *ECDSA            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECDSACallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ECDSACallerSession struct {
	Contract *ECDSACaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ECDSATransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ECDSATransactorSession struct {
	Contract     *ECDSATransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECDSARaw is an auto generated low-level Go binding around an Ethereum contract.
type ECDSARaw struct {
	Contract *ECDSA // Generic contract binding to access the raw methods on
}

// ECDSACallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ECDSACallerRaw struct {
	Contract *ECDSACaller // Generic read-only contract binding to access the raw methods on
}

// ECDSATransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ECDSATransactorRaw struct {
	Contract *ECDSATransactor // Generic write-only contract binding to access the raw methods on
}

// NewECDSA creates a new instance of ECDSA, bound to a specific deployed contract.
func NewECDSA(address common.Address, backend bind.ContractBackend) (*ECDSA, error) {
	contract, err := bindECDSA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ECDSA{ECDSACaller: ECDSACaller{contract: contract}, ECDSATransactor: ECDSATransactor{contract: contract}, ECDSAFilterer: ECDSAFilterer{contract: contract}}, nil
}

// NewECDSACaller creates a new read-only instance of ECDSA, bound to a specific deployed contract.
func NewECDSACaller(address common.Address, caller bind.ContractCaller) (*ECDSACaller, error) {
	contract, err := bindECDSA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSACaller{contract: contract}, nil
}

// NewECDSATransactor creates a new write-only instance of ECDSA, bound to a specific deployed contract.
func NewECDSATransactor(address common.Address, transactor bind.ContractTransactor) (*ECDSATransactor, error) {
	contract, err := bindECDSA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSATransactor{contract: contract}, nil
}

// NewECDSAFilterer creates a new log filterer instance of ECDSA, bound to a specific deployed contract.
func NewECDSAFilterer(address common.Address, filterer bind.ContractFilterer) (*ECDSAFilterer, error) {
	contract, err := bindECDSA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ECDSAFilterer{contract: contract}, nil
}

// bindECDSA binds a generic wrapper to an already deployed contract.
func bindECDSA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ECDSAABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSA *ECDSARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSA.Contract.ECDSACaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSA *ECDSARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSA.Contract.ECDSATransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSA *ECDSARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSA.Contract.ECDSATransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSA *ECDSACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSA.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSA *ECDSATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSA.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSA *ECDSATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSA.Contract.contract.Transact(opts, method, params...)
}

// ManagementContractMetaData contains all meta data concerning the ManagementContract contract.
var ManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_parentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_l1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_parentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"_r\",\"type\":\"tuple\"}],\"name\":\"AppendRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"Attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"element\",\"type\":\"tuple\"}],\"name\":\"GetParentRollup\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"GetRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rollupID\",\"type\":\"uint256\"}],\"name\":\"GetRollupByID\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HasSecondCousinFork\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"InitializeTree\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"IsWithdrawalAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1b1e5402": "AddRollup(bytes32,bytes32,address,bytes32,uint256,string)",
		"1ac63930": "AppendRollup(uint256,(bytes32,bytes32,address,bytes32,uint256))",
		"43348b2f": "Attested(address)",
		"31b1d255": "GetParentRollup((uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))",
		"8236a7ba": "GetRollupByHash(bytes32)",
		"92aaec79": "GetRollupByID(uint256)",
		"57b70600": "HasSecondCousinFork()",
		"c719bf50": "InitializeNetworkSecret(address,bytes)",
		"73bba846": "InitializeTree((bytes32,bytes32,address,bytes32,uint256))",
		"a52f433c": "IsWithdrawalAvailable()",
		"e34fbfc8": "RequestNetworkSecret(string)",
		"981214ba": "RespondNetworkSecret(address,address,bytes,bytes)",
	},
	Bin: "0x608060405234801561001057600080fd5b5061166a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c80638236a7ba116100715780638236a7ba146101b357806392aaec79146101c6578063981214ba146101d9578063a52f433c146101ec578063c719bf50146101fc578063e34fbfc81461020f57600080fd5b80631ac63930146100b95780631b1e5402146100ce57806331b1d255146100e157806343348b2f1461015c57806357b706001461019857806373bba846146101a0575b600080fd5b6100cc6100c7366004611115565b610222565b005b6100cc6100dc366004611184565b610358565b6100f46100ef3660046111ff565b6104ad565b6040805192151583528151602080850191909152808301518483015291810151805160608086019190915292810151608080860191909152918101516001600160a01b031660a08501529182015160c0840152015160e0820152610100015b60405180910390f35b61018861016a366004611260565b6001600160a01b031660009081526001602052604090205460ff1690565b6040519015158152602001610153565b6101886104cd565b6100cc6101ae366004611282565b610675565b6100f46101c136600461129e565b61084e565b6100f46101d436600461129e565b61086c565b6100cc6101e7366004611344565b6108ec565b600954610100900460ff16610188565b6100cc61020a3660046113c9565b6109f5565b6100cc61021d36600461141c565b610a3c565b60078054908190600061023483611474565b91905055506000806102458561086c565b915091508161028e5760405162461bcd60e51b815260206004820152601060248201526f1c185c995b9d081b9bdd08199bdd5b9960821b60448201526064015b60405180910390fd5b604080516060808201835285825260208083018981528385018981526000898152600280855287822096518755925160018088019190915591518051938701939093558284015160038088019190915583880151600480890180546001600160a01b0319166001600160a01b03909316929092179091559584015160058801556080909301516006968701558b81529383528584208054918201815584528284200188905588820151835290529190912084905554815114156103515760068390555b5050505050565b6001600160a01b03851660009081526001602052604090205460ff166103c05760405162461bcd60e51b815260206004820152601760248201527f61676772656761746f72206e6f742061747465737465640000000000000000006044820152606401610285565b6040805160a081018252888152602081018890526001600160a01b03871691810191909152606081018590526080810184905260085460ff1661040c5761040681610675565b506104a4565b6000806104188a61084e565b91509150816104695760405162461bcd60e51b815260206004820152601a60248201527f756e61626c6520746f2066696e6420706172656e7420686173680000000000006044820152606401610285565b6006546002101561049457600061047e6104cd565b90508015610492576009805461ff00191690555b505b80516104a09084610222565b5050505b50505050505050565b60006104b7610f91565b6104c4836020015161086c565b91509150915091565b6000806104d8610a5b565b90506000806104e6836104ad565b91509150816105235760405162461bcd60e51b81526020600482015260096024820152681b9bc81c185c995b9d60ba1b6044820152606401610285565b60008061052f836104ad565b91509150816105725760405162461bcd60e51b815260206004820152600f60248201526e1b9bc819dc985b99081c185c995b9d608a1b6044820152606401610285565b80516000908152600460209081526040808320805482518185028101850190935280835291929091908301828280156105ca57602002820191906000526020600020905b8154815260200190600101908083116105b6575b5050505050905060005b8151811015610667576000806106028484815181106105f5576105f561148f565b602002602001015161086c565b915091508161061b576000995050505050505050505090565b86518151141561062c575050610655565b805160009081526004602052604090205415610652576001995050505050505050505090565b50505b8061065f81611474565b9150506105d4565b506000965050505050505090565b60085460ff16156106c85760405162461bcd60e51b815260206004820152601b60248201527f63616e6e6f7420626520696e697469616c697a656420616761696e00000000006044820152606401610285565b6008805460ff19166001908117909155604080516060808201835283825260006020808401828152848601888152878452600280845295517fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e05590517fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e1555180517fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e255808201517fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e355808601517fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e480546001600160a01b039092166001600160a01b0319909216919091179055928301517fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e5556080909201517fe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e65560068590556007929092559384015181526003909352909120556009805461ff001916610100179055565b6000610858610f91565b6000838152600360205260409020546104c4905b6000610876610f91565b505060009081526002602081815260409283902083516060808201865282548252600183015482850152855160a08101875294830154855260038301549385019390935260048201546001600160a01b031684860152600582015492840192909252600601546080830152918201528051151591565b6001600160a01b03841660009081526001602052604090205460ff168061091257600080fd5b600061094086868560405160200161092c939291906114d5565b604051602081830303815290604052610ad8565b9050600061094e8286610b13565b9050866001600160a01b0316816001600160a01b0316146109c65760405162461bcd60e51b815260206004820152602c60248201527f63616c63756c61746564206164647265737320616e642061747465737465724960448201526b08840c8dedce840dac2e8c6d60a31b6064820152608401610285565b5050506001600160a01b039092166000908152600160208190526040909120805460ff19169091179055505050565b60095460ff1615610a0557600080fd5b50506009805460ff1990811660019081179092556001600160a01b0390921660009081526020829052604090208054909216179055565b336000908152602081905260409020610a56908383610fd3565b505050565b610a63610f91565b506006805460009081526002602081815260409283902083516060808201865282548252600183015482850152855160a08101875294830154855260038301549385019390935260048201546001600160a01b0316848601526005820154928401929092529093015460808201529082015290565b6000610ae48251610b37565b82604051602001610af692919061151b565b604051602081830303815290604052805190602001209050919050565b6000806000610b228585610c3d565b91509150610b2f81610cad565b509392505050565b606081610b5b5750506040805180820190915260018152600360fc1b602082015290565b8160005b8115610b855780610b6f81611474565b9150610b7e9050600a8361158c565b9150610b5f565b60008167ffffffffffffffff811115610ba057610ba061106c565b6040519080825280601f01601f191660200182016040528015610bca576020820181803683370190505b5090505b8415610c3557610bdf6001836115a0565b9150610bec600a866115b7565b610bf79060306115cb565b60f81b818381518110610c0c57610c0c61148f565b60200101906001600160f81b031916908160001a905350610c2e600a8661158c565b9450610bce565b949350505050565b600080825160411415610c745760208301516040840151606085015160001a610c6887828585610e6b565b94509450505050610ca6565b825160401415610c9e5760208301516040840151610c93868383610f58565b935093505050610ca6565b506000905060025b9250929050565b6000816004811115610cc157610cc16115e3565b1415610cca5750565b6001816004811115610cde57610cde6115e3565b1415610d2c5760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e617475726500000000000000006044820152606401610285565b6002816004811115610d4057610d406115e3565b1415610d8e5760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e677468006044820152606401610285565b6003816004811115610da257610da26115e3565b1415610dfb5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b6064820152608401610285565b6004816004811115610e0f57610e0f6115e3565b1415610e685760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b6064820152608401610285565b50565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0831115610ea25750600090506003610f4f565b8460ff16601b14158015610eba57508460ff16601c14155b15610ecb5750600090506004610f4f565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015610f1f573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610f4857600060019250925050610f4f565b9150600090505b94509492505050565b6000806001600160ff1b03831681610f7560ff86901c601b6115cb565b9050610f8387828885610e6b565b935093505050935093915050565b604080516060808201835260008083526020808401829052845160a0810186528281529081018290528085018290529182018190526080820152909182015290565b828054610fdf906115f9565b90600052602060002090601f0160209004810192826110015760008555611047565b82601f1061101a5782800160ff19823516178555611047565b82800160010185558215611047579182015b8281111561104757823582559160200191906001019061102c565b50611053929150611057565b5090565b5b808211156110535760008155600101611058565b634e487b7160e01b600052604160045260246000fd5b80356001600160a01b038116811461109957600080fd5b919050565b600060a082840312156110b057600080fd5b60405160a0810181811067ffffffffffffffff821117156110d3576110d361106c565b806040525080915082358152602083013560208201526110f560408401611082565b604082015260608301356060820152608083013560808201525092915050565b60008060c0838503121561112857600080fd5b82359150611139846020850161109e565b90509250929050565b60008083601f84011261115457600080fd5b50813567ffffffffffffffff81111561116c57600080fd5b602083019150836020828501011115610ca657600080fd5b600080600080600080600060c0888a03121561119f57600080fd5b87359650602088013595506111b660408901611082565b9450606088013593506080880135925060a088013567ffffffffffffffff8111156111e057600080fd5b6111ec8a828b01611142565b989b979a50959850939692959293505050565b600060e0828403121561121157600080fd5b6040516060810181811067ffffffffffffffff821117156112345761123461106c565b80604052508235815260208301356020820152611254846040850161109e565b60408201529392505050565b60006020828403121561127257600080fd5b61127b82611082565b9392505050565b600060a0828403121561129457600080fd5b61127b838361109e565b6000602082840312156112b057600080fd5b5035919050565b600082601f8301126112c857600080fd5b813567ffffffffffffffff808211156112e3576112e361106c565b604051601f8301601f19908116603f0116810190828211818310171561130b5761130b61106c565b8160405283815286602085880101111561132457600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806000806080858703121561135a57600080fd5b61136385611082565b935061137160208601611082565b9250604085013567ffffffffffffffff8082111561138e57600080fd5b61139a888389016112b7565b935060608701359150808211156113b057600080fd5b506113bd878288016112b7565b91505092959194509250565b6000806000604084860312156113de57600080fd5b6113e784611082565b9250602084013567ffffffffffffffff81111561140357600080fd5b61140f86828701611142565b9497909650939450505050565b6000806020838503121561142f57600080fd5b823567ffffffffffffffff81111561144657600080fd5b61145285828601611142565b90969095509350505050565b634e487b7160e01b600052601160045260246000fd5b60006000198214156114885761148861145e565b5060010190565b634e487b7160e01b600052603260045260246000fd5b60005b838110156114c05781810151838201526020016114a8565b838111156114cf576000848401525b50505050565b60006bffffffffffffffffffffffff19808660601b168352808560601b16601484015250825161150c8160288501602087016114a5565b91909101602801949350505050565b7f19457468657265756d205369676e6564204d6573736167653a0a00000000000081526000835161155381601a8501602088016114a5565b83519083019061156a81601a8401602088016114a5565b01601a01949350505050565b634e487b7160e01b600052601260045260246000fd5b60008261159b5761159b611576565b500490565b6000828210156115b2576115b261145e565b500390565b6000826115c6576115c6611576565b500690565b600082198211156115de576115de61145e565b500190565b634e487b7160e01b600052602160045260246000fd5b600181811c9082168061160d57607f821691505b6020821081141561162e57634e487b7160e01b600052602260045260246000fd5b5091905056fea26469706673582212204eb479abee67414e42243391afc360ca1ffead5267f5edd9fc2448169a879ca864736f6c634300080c0033",
}

// ManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ManagementContractMetaData.ABI instead.
var ManagementContractABI = ManagementContractMetaData.ABI

// Deprecated: Use ManagementContractMetaData.Sigs instead.
// ManagementContractFuncSigs maps the 4-byte function signature to its string representation.
var ManagementContractFuncSigs = ManagementContractMetaData.Sigs

// ManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ManagementContractMetaData.Bin instead.
var ManagementContractBin = ManagementContractMetaData.Bin

// DeployManagementContract deploys a new Ethereum contract, binding an instance of ManagementContract to it.
func DeployManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ManagementContract, error) {
	parsed, err := ManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ManagementContract{ManagementContractCaller: ManagementContractCaller{contract: contract}, ManagementContractTransactor: ManagementContractTransactor{contract: contract}, ManagementContractFilterer: ManagementContractFilterer{contract: contract}}, nil
}

// ManagementContract is an auto generated Go binding around an Ethereum contract.
type ManagementContract struct {
	ManagementContractCaller     // Read-only binding to the contract
	ManagementContractTransactor // Write-only binding to the contract
	ManagementContractFilterer   // Log filterer for contract events
}

// ManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ManagementContractSession struct {
	Contract     *ManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ManagementContractCallerSession struct {
	Contract *ManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ManagementContractTransactorSession struct {
	Contract     *ManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ManagementContractRaw struct {
	Contract *ManagementContract // Generic contract binding to access the raw methods on
}

// ManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ManagementContractCallerRaw struct {
	Contract *ManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// ManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ManagementContractTransactorRaw struct {
	Contract *ManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewManagementContract creates a new instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContract(address common.Address, backend bind.ContractBackend) (*ManagementContract, error) {
	contract, err := bindManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ManagementContract{ManagementContractCaller: ManagementContractCaller{contract: contract}, ManagementContractTransactor: ManagementContractTransactor{contract: contract}, ManagementContractFilterer: ManagementContractFilterer{contract: contract}}, nil
}

// NewManagementContractCaller creates a new read-only instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractCaller(address common.Address, caller bind.ContractCaller) (*ManagementContractCaller, error) {
	contract, err := bindManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ManagementContractCaller{contract: contract}, nil
}

// NewManagementContractTransactor creates a new write-only instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ManagementContractTransactor, error) {
	contract, err := bindManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ManagementContractTransactor{contract: contract}, nil
}

// NewManagementContractFilterer creates a new log filterer instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ManagementContractFilterer, error) {
	contract, err := bindManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ManagementContractFilterer{contract: contract}, nil
}

// bindManagementContract binds a generic wrapper to an already deployed contract.
func bindManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ManagementContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagementContract *ManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ManagementContract.Contract.ManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagementContract *ManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.Contract.ManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagementContract *ManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagementContract.Contract.ManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagementContract *ManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagementContract *ManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagementContract *ManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagementContract.Contract.contract.Transact(opts, method, params...)
}

// Attested is a free data retrieval call binding the contract method 0x43348b2f.
//
// Solidity: function Attested(address _addr) view returns(bool)
func (_ManagementContract *ManagementContractCaller) Attested(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "Attested", _addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Attested is a free data retrieval call binding the contract method 0x43348b2f.
//
// Solidity: function Attested(address _addr) view returns(bool)
func (_ManagementContract *ManagementContractSession) Attested(_addr common.Address) (bool, error) {
	return _ManagementContract.Contract.Attested(&_ManagementContract.CallOpts, _addr)
}

// Attested is a free data retrieval call binding the contract method 0x43348b2f.
//
// Solidity: function Attested(address _addr) view returns(bool)
func (_ManagementContract *ManagementContractCallerSession) Attested(_addr common.Address) (bool, error) {
	return _ManagementContract.Contract.Attested(&_ManagementContract.CallOpts, _addr)
}

// GetParentRollup is a free data retrieval call binding the contract method 0x31b1d255.
//
// Solidity: function GetParentRollup((uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)) element) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractCaller) GetParentRollup(opts *bind.CallOpts, element ManagementContractTreeElement) (bool, ManagementContractTreeElement, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "GetParentRollup", element)

	if err != nil {
		return *new(bool), *new(ManagementContractTreeElement), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(ManagementContractTreeElement)).(*ManagementContractTreeElement)

	return out0, out1, err

}

// GetParentRollup is a free data retrieval call binding the contract method 0x31b1d255.
//
// Solidity: function GetParentRollup((uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)) element) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractSession) GetParentRollup(element ManagementContractTreeElement) (bool, ManagementContractTreeElement, error) {
	return _ManagementContract.Contract.GetParentRollup(&_ManagementContract.CallOpts, element)
}

// GetParentRollup is a free data retrieval call binding the contract method 0x31b1d255.
//
// Solidity: function GetParentRollup((uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)) element) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractCallerSession) GetParentRollup(element ManagementContractTreeElement) (bool, ManagementContractTreeElement, error) {
	return _ManagementContract.Contract.GetParentRollup(&_ManagementContract.CallOpts, element)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x8236a7ba.
//
// Solidity: function GetRollupByHash(bytes32 rollupHash) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractCaller) GetRollupByHash(opts *bind.CallOpts, rollupHash [32]byte) (bool, ManagementContractTreeElement, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "GetRollupByHash", rollupHash)

	if err != nil {
		return *new(bool), *new(ManagementContractTreeElement), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(ManagementContractTreeElement)).(*ManagementContractTreeElement)

	return out0, out1, err

}

// GetRollupByHash is a free data retrieval call binding the contract method 0x8236a7ba.
//
// Solidity: function GetRollupByHash(bytes32 rollupHash) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractSession) GetRollupByHash(rollupHash [32]byte) (bool, ManagementContractTreeElement, error) {
	return _ManagementContract.Contract.GetRollupByHash(&_ManagementContract.CallOpts, rollupHash)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x8236a7ba.
//
// Solidity: function GetRollupByHash(bytes32 rollupHash) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractCallerSession) GetRollupByHash(rollupHash [32]byte) (bool, ManagementContractTreeElement, error) {
	return _ManagementContract.Contract.GetRollupByHash(&_ManagementContract.CallOpts, rollupHash)
}

// GetRollupByID is a free data retrieval call binding the contract method 0x92aaec79.
//
// Solidity: function GetRollupByID(uint256 rollupID) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractCaller) GetRollupByID(opts *bind.CallOpts, rollupID *big.Int) (bool, ManagementContractTreeElement, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "GetRollupByID", rollupID)

	if err != nil {
		return *new(bool), *new(ManagementContractTreeElement), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(ManagementContractTreeElement)).(*ManagementContractTreeElement)

	return out0, out1, err

}

// GetRollupByID is a free data retrieval call binding the contract method 0x92aaec79.
//
// Solidity: function GetRollupByID(uint256 rollupID) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractSession) GetRollupByID(rollupID *big.Int) (bool, ManagementContractTreeElement, error) {
	return _ManagementContract.Contract.GetRollupByID(&_ManagementContract.CallOpts, rollupID)
}

// GetRollupByID is a free data retrieval call binding the contract method 0x92aaec79.
//
// Solidity: function GetRollupByID(uint256 rollupID) view returns(bool, (uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))
func (_ManagementContract *ManagementContractCallerSession) GetRollupByID(rollupID *big.Int) (bool, ManagementContractTreeElement, error) {
	return _ManagementContract.Contract.GetRollupByID(&_ManagementContract.CallOpts, rollupID)
}

// HasSecondCousinFork is a free data retrieval call binding the contract method 0x57b70600.
//
// Solidity: function HasSecondCousinFork() view returns(bool)
func (_ManagementContract *ManagementContractCaller) HasSecondCousinFork(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "HasSecondCousinFork")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasSecondCousinFork is a free data retrieval call binding the contract method 0x57b70600.
//
// Solidity: function HasSecondCousinFork() view returns(bool)
func (_ManagementContract *ManagementContractSession) HasSecondCousinFork() (bool, error) {
	return _ManagementContract.Contract.HasSecondCousinFork(&_ManagementContract.CallOpts)
}

// HasSecondCousinFork is a free data retrieval call binding the contract method 0x57b70600.
//
// Solidity: function HasSecondCousinFork() view returns(bool)
func (_ManagementContract *ManagementContractCallerSession) HasSecondCousinFork() (bool, error) {
	return _ManagementContract.Contract.HasSecondCousinFork(&_ManagementContract.CallOpts)
}

// IsWithdrawalAvailable is a free data retrieval call binding the contract method 0xa52f433c.
//
// Solidity: function IsWithdrawalAvailable() view returns(bool)
func (_ManagementContract *ManagementContractCaller) IsWithdrawalAvailable(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "IsWithdrawalAvailable")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWithdrawalAvailable is a free data retrieval call binding the contract method 0xa52f433c.
//
// Solidity: function IsWithdrawalAvailable() view returns(bool)
func (_ManagementContract *ManagementContractSession) IsWithdrawalAvailable() (bool, error) {
	return _ManagementContract.Contract.IsWithdrawalAvailable(&_ManagementContract.CallOpts)
}

// IsWithdrawalAvailable is a free data retrieval call binding the contract method 0xa52f433c.
//
// Solidity: function IsWithdrawalAvailable() view returns(bool)
func (_ManagementContract *ManagementContractCallerSession) IsWithdrawalAvailable() (bool, error) {
	return _ManagementContract.Contract.IsWithdrawalAvailable(&_ManagementContract.CallOpts)
}

// AddRollup is a paid mutator transaction binding the contract method 0x1b1e5402.
//
// Solidity: function AddRollup(bytes32 _parentHash, bytes32 _hash, address _aggregatorID, bytes32 _l1Block, uint256 _number, string _rollupData) returns()
func (_ManagementContract *ManagementContractTransactor) AddRollup(opts *bind.TransactOpts, _parentHash [32]byte, _hash [32]byte, _aggregatorID common.Address, _l1Block [32]byte, _number *big.Int, _rollupData string) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "AddRollup", _parentHash, _hash, _aggregatorID, _l1Block, _number, _rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0x1b1e5402.
//
// Solidity: function AddRollup(bytes32 _parentHash, bytes32 _hash, address _aggregatorID, bytes32 _l1Block, uint256 _number, string _rollupData) returns()
func (_ManagementContract *ManagementContractSession) AddRollup(_parentHash [32]byte, _hash [32]byte, _aggregatorID common.Address, _l1Block [32]byte, _number *big.Int, _rollupData string) (*types.Transaction, error) {
	return _ManagementContract.Contract.AddRollup(&_ManagementContract.TransactOpts, _parentHash, _hash, _aggregatorID, _l1Block, _number, _rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0x1b1e5402.
//
// Solidity: function AddRollup(bytes32 _parentHash, bytes32 _hash, address _aggregatorID, bytes32 _l1Block, uint256 _number, string _rollupData) returns()
func (_ManagementContract *ManagementContractTransactorSession) AddRollup(_parentHash [32]byte, _hash [32]byte, _aggregatorID common.Address, _l1Block [32]byte, _number *big.Int, _rollupData string) (*types.Transaction, error) {
	return _ManagementContract.Contract.AddRollup(&_ManagementContract.TransactOpts, _parentHash, _hash, _aggregatorID, _l1Block, _number, _rollupData)
}

// AppendRollup is a paid mutator transaction binding the contract method 0x1ac63930.
//
// Solidity: function AppendRollup(uint256 _parentID, (bytes32,bytes32,address,bytes32,uint256) _r) returns()
func (_ManagementContract *ManagementContractTransactor) AppendRollup(opts *bind.TransactOpts, _parentID *big.Int, _r ManagementContractMetaRollup) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "AppendRollup", _parentID, _r)
}

// AppendRollup is a paid mutator transaction binding the contract method 0x1ac63930.
//
// Solidity: function AppendRollup(uint256 _parentID, (bytes32,bytes32,address,bytes32,uint256) _r) returns()
func (_ManagementContract *ManagementContractSession) AppendRollup(_parentID *big.Int, _r ManagementContractMetaRollup) (*types.Transaction, error) {
	return _ManagementContract.Contract.AppendRollup(&_ManagementContract.TransactOpts, _parentID, _r)
}

// AppendRollup is a paid mutator transaction binding the contract method 0x1ac63930.
//
// Solidity: function AppendRollup(uint256 _parentID, (bytes32,bytes32,address,bytes32,uint256) _r) returns()
func (_ManagementContract *ManagementContractTransactorSession) AppendRollup(_parentID *big.Int, _r ManagementContractMetaRollup) (*types.Transaction, error) {
	return _ManagementContract.Contract.AppendRollup(&_ManagementContract.TransactOpts, _parentID, _r)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret) returns()
func (_ManagementContract *ManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, _aggregatorID common.Address, _initSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "InitializeNetworkSecret", _aggregatorID, _initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret) returns()
func (_ManagementContract *ManagementContractSession) InitializeNetworkSecret(_aggregatorID common.Address, _initSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, _aggregatorID, _initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret) returns()
func (_ManagementContract *ManagementContractTransactorSession) InitializeNetworkSecret(_aggregatorID common.Address, _initSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, _aggregatorID, _initSecret)
}

// InitializeTree is a paid mutator transaction binding the contract method 0x73bba846.
//
// Solidity: function InitializeTree((bytes32,bytes32,address,bytes32,uint256) r) returns()
func (_ManagementContract *ManagementContractTransactor) InitializeTree(opts *bind.TransactOpts, r ManagementContractMetaRollup) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "InitializeTree", r)
}

// InitializeTree is a paid mutator transaction binding the contract method 0x73bba846.
//
// Solidity: function InitializeTree((bytes32,bytes32,address,bytes32,uint256) r) returns()
func (_ManagementContract *ManagementContractSession) InitializeTree(r ManagementContractMetaRollup) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeTree(&_ManagementContract.TransactOpts, r)
}

// InitializeTree is a paid mutator transaction binding the contract method 0x73bba846.
//
// Solidity: function InitializeTree((bytes32,bytes32,address,bytes32,uint256) r) returns()
func (_ManagementContract *ManagementContractTransactorSession) InitializeTree(r ManagementContractMetaRollup) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeTree(&_ManagementContract.TransactOpts, r)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_ManagementContract *ManagementContractTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "RequestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_ManagementContract *ManagementContractSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _ManagementContract.Contract.RequestNetworkSecret(&_ManagementContract.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_ManagementContract *ManagementContractTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _ManagementContract.Contract.RequestNetworkSecret(&_ManagementContract.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_ManagementContract *ManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_ManagementContract *ManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.Contract.RespondNetworkSecret(&_ManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_ManagementContract *ManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.Contract.RespondNetworkSecret(&_ManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}

// StringsMetaData contains all meta data concerning the Strings contract.
var StringsMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122096e640a4442fdf7ef4a072b5526cb009808d85b9a856a0f78712415574c622e864736f6c634300080c0033",
}

// StringsABI is the input ABI used to generate the binding from.
// Deprecated: Use StringsMetaData.ABI instead.
var StringsABI = StringsMetaData.ABI

// StringsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StringsMetaData.Bin instead.
var StringsBin = StringsMetaData.Bin

// DeployStrings deploys a new Ethereum contract, binding an instance of Strings to it.
func DeployStrings(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Strings, error) {
	parsed, err := StringsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StringsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Strings{StringsCaller: StringsCaller{contract: contract}, StringsTransactor: StringsTransactor{contract: contract}, StringsFilterer: StringsFilterer{contract: contract}}, nil
}

// Strings is an auto generated Go binding around an Ethereum contract.
type Strings struct {
	StringsCaller     // Read-only binding to the contract
	StringsTransactor // Write-only binding to the contract
	StringsFilterer   // Log filterer for contract events
}

// StringsCaller is an auto generated read-only Go binding around an Ethereum contract.
type StringsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StringsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StringsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StringsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StringsSession struct {
	Contract     *Strings          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StringsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StringsCallerSession struct {
	Contract *StringsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StringsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StringsTransactorSession struct {
	Contract     *StringsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StringsRaw is an auto generated low-level Go binding around an Ethereum contract.
type StringsRaw struct {
	Contract *Strings // Generic contract binding to access the raw methods on
}

// StringsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StringsCallerRaw struct {
	Contract *StringsCaller // Generic read-only contract binding to access the raw methods on
}

// StringsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StringsTransactorRaw struct {
	Contract *StringsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStrings creates a new instance of Strings, bound to a specific deployed contract.
func NewStrings(address common.Address, backend bind.ContractBackend) (*Strings, error) {
	contract, err := bindStrings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Strings{StringsCaller: StringsCaller{contract: contract}, StringsTransactor: StringsTransactor{contract: contract}, StringsFilterer: StringsFilterer{contract: contract}}, nil
}

// NewStringsCaller creates a new read-only instance of Strings, bound to a specific deployed contract.
func NewStringsCaller(address common.Address, caller bind.ContractCaller) (*StringsCaller, error) {
	contract, err := bindStrings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StringsCaller{contract: contract}, nil
}

// NewStringsTransactor creates a new write-only instance of Strings, bound to a specific deployed contract.
func NewStringsTransactor(address common.Address, transactor bind.ContractTransactor) (*StringsTransactor, error) {
	contract, err := bindStrings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StringsTransactor{contract: contract}, nil
}

// NewStringsFilterer creates a new log filterer instance of Strings, bound to a specific deployed contract.
func NewStringsFilterer(address common.Address, filterer bind.ContractFilterer) (*StringsFilterer, error) {
	contract, err := bindStrings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StringsFilterer{contract: contract}, nil
}

// bindStrings binds a generic wrapper to an already deployed contract.
func bindStrings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StringsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Strings *StringsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Strings.Contract.StringsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Strings *StringsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Strings.Contract.StringsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Strings *StringsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Strings.Contract.StringsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Strings *StringsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Strings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Strings *StringsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Strings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Strings *StringsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Strings.Contract.contract.Transact(opts, method, params...)
}
