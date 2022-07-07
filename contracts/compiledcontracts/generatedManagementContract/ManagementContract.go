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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_parentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_l1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_parentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"_r\",\"type\":\"tuple\"}],\"name\":\"AppendRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"Attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"element\",\"type\":\"tuple\"}],\"name\":\"GetParentRollup\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"GetRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rollupID\",\"type\":\"uint256\"}],\"name\":\"GetRollupByID\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HasSecondCousinFork\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_hostAddress\",\"type\":\"string\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"InitializeTree\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"IsWithdrawalAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"verifyAttester\",\"type\":\"bool\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1b1e5402": "AddRollup(bytes32,bytes32,address,bytes32,uint256,string)",
		"1ac63930": "AppendRollup(uint256,(bytes32,bytes32,address,bytes32,uint256))",
		"43348b2f": "Attested(address)",
		"324ff866": "GetHostAddresses()",
		"31b1d255": "GetParentRollup((uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))",
		"8236a7ba": "GetRollupByHash(bytes32)",
		"92aaec79": "GetRollupByID(uint256)",
		"57b70600": "HasSecondCousinFork()",
		"68e10383": "InitializeNetworkSecret(address,bytes,string)",
		"73bba846": "InitializeTree((bytes32,bytes32,address,bytes32,uint256))",
		"a52f433c": "IsWithdrawalAvailable()",
		"e34fbfc8": "RequestNetworkSecret(string)",
		"bbd79e15": "RespondNetworkSecret(address,address,bytes,bytes,string,bool)",
	},
	Bin: "0x608060405234801561001057600080fd5b50611959806100206000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806368e103831161008c57806392aaec791161006657806392aaec7914610209578063a52f433c1461021c578063bbd79e151461022c578063e34fbfc81461023f57600080fd5b806368e10383146101d057806373bba846146101e35780638236a7ba146101f657600080fd5b80631ac63930146100d45780631b1e5402146100e957806331b1d255146100fc578063324ff8661461017757806343348b2f1461018c57806357b70600146101c8575b600080fd5b6100e76100e2366004611310565b610252565b005b6100e76100f736600461137f565b610389565b61010f61010a3660046113fa565b6104de565b6040805192151583528151602080850191909152808301518483015291810151805160608086019190915292810151608080860191909152918101516001600160a01b031660a08501529182015160c0840152015160e0820152610100015b60405180910390f35b61017f6104fe565b60405161016e919061148b565b6101b861019a366004611505565b6001600160a01b031660009081526001602052604090205460ff1690565b604051901515815260200161016e565b6101b86105d7565b6100e76101de3660046115b4565b61077f565b6100e76101f136600461162e565b6107fe565b61010f61020436600461164a565b6109d4565b61010f61021736600461164a565b6109f2565b600a54610100900460ff166101b8565b6100e761023a366004611663565b610a72565b6100e761024d366004611725565b610bc4565b6008805490819060006102648361177d565b9190505550600080610275856109f2565b91509150816102be5760405162461bcd60e51b815260206004820152601060248201526f1c185c995b9d081b9bdd08199bdd5b9960821b60448201526064015b60405180910390fd5b6040805160608082018352858252602080830189815283850189815260008981526003808552878220965187559251600180880191909155915180516002880155808501519387019390935582870151600480880180546001600160a01b0319166001600160a01b0390931692909217909155948301516005808801919091556080909301516006909601959095558a8552908252848420805491820181558452818420018790558781015183525220839055600754815114156103825760078390555b5050505050565b6001600160a01b03851660009081526001602052604090205460ff166103f15760405162461bcd60e51b815260206004820152601760248201527f61676772656761746f72206e6f7420617474657374656400000000000000000060448201526064016102b5565b6040805160a081018252888152602081018890526001600160a01b03871691810191909152606081018590526080810184905260095460ff1661043d57610437816107fe565b506104d5565b6000806104498a6109d4565b915091508161049a5760405162461bcd60e51b815260206004820152601a60248201527f756e61626c6520746f2066696e6420706172656e74206861736800000000000060448201526064016102b5565b600754600210156104c55760006104af6105d7565b905080156104c357600a805461ff00191690555b505b80516104d19084610252565b5050505b50505050505050565b60006104e8611118565b6104f583602001516109f2565b91509150915091565b60606002805480602002602001604051908101604052809291908181526020016000905b828210156105ce57838290600052602060002001805461054190611798565b80601f016020809104026020016040519081016040528092919081815260200182805461056d90611798565b80156105ba5780601f1061058f576101008083540402835291602001916105ba565b820191906000526020600020905b81548152906001019060200180831161059d57829003601f168201915b505050505081526020019060010190610522565b50505050905090565b6000806105e2610be3565b90506000806105f0836104de565b915091508161062d5760405162461bcd60e51b81526020600482015260096024820152681b9bc81c185c995b9d60ba1b60448201526064016102b5565b600080610639836104de565b915091508161067c5760405162461bcd60e51b815260206004820152600f60248201526e1b9bc819dc985b99081c185c995b9d608a1b60448201526064016102b5565b80516000908152600560209081526040808320805482518185028101850190935280835291929091908301828280156106d457602002820191906000526020600020905b8154815260200190600101908083116106c0575b5050505050905060005b81518110156107715760008061070c8484815181106106ff576106ff6117d3565b60200260200101516109f2565b9150915081610725576000995050505050505050505090565b86518151141561073657505061075f565b80516000908152600560205260409020541561075c576001995050505050505050505090565b50505b806107698161177d565b9150506106de565b506000965050505050505090565b600a5460ff161561078f57600080fd5b600a8054600160ff1991821681179092556001600160a01b03861660009081526020838152604082208054909316841790925560028054938401815590528251610382927f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace019184019061115a565b60095460ff16156108515760405162461bcd60e51b815260206004820152601b60248201527f63616e6e6f7420626520696e697469616c697a656420616761696e000000000060448201526064016102b5565b6009805460ff191660019081179091556040805160608082018352838252600060208084018281528486018881528784526003835294517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c55517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054d55925180517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054e55808401517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054f55808501517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3055080546001600160a01b039092166001600160a01b0319909216919091179055918201517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c30551556080909101517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3055255600784905560026008559381015184526004905290912055600a805461ff001916610100179055565b60006109de611118565b6000838152600460205260409020546104f5905b60006109fc611118565b505060009081526003602081815260409283902083516060808201865282548252600183015482850152855160a08101875260028401548152948301549385019390935260048201546001600160a01b031684860152600582015492840192909252600601546080830152918201528051151591565b6001600160a01b03861660009081526001602052604090205460ff1680610a9857600080fd5b8115610b57576000610ace88888688604051602001610aba94939291906117e9565b604051602081830303815290604052610c5f565b90506000610adc8288610c9a565b9050886001600160a01b0316816001600160a01b031614610b545760405162461bcd60e51b815260206004820152602c60248201527f63616c63756c61746564206164647265737320616e642061747465737465724960448201526b08840c8dedce840dac2e8c6d60a31b60648201526084016102b5565b50505b6001600160a01b03861660009081526001602081815260408320805460ff1916831790556002805492830181559092528451610bba927f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace9092019186019061115a565b5050505050505050565b336000908152602081905260409020610bde9083836111de565b505050565b610beb611118565b5060075460009081526003602081815260409283902083516060808201865282548252600183015482850152855160a08101875260028401548152948301549385019390935260048201546001600160a01b0316848601526005820154928401929092526006015460808301529182015290565b6000610c6b8251610cbe565b82604051602001610c7d929190611845565b604051602081830303815290604052805190602001209050919050565b6000806000610ca98585610dc4565b91509150610cb681610e34565b509392505050565b606081610ce25750506040805180820190915260018152600360fc1b602082015290565b8160005b8115610d0c5780610cf68161177d565b9150610d059050600a836118b6565b9150610ce6565b60008167ffffffffffffffff811115610d2757610d27611267565b6040519080825280601f01601f191660200182016040528015610d51576020820181803683370190505b5090505b8415610dbc57610d666001836118ca565b9150610d73600a866118e1565b610d7e9060306118f5565b60f81b818381518110610d9357610d936117d3565b60200101906001600160f81b031916908160001a905350610db5600a866118b6565b9450610d55565b949350505050565b600080825160411415610dfb5760208301516040840151606085015160001a610def87828585610ff2565b94509450505050610e2d565b825160401415610e255760208301516040840151610e1a8683836110df565b935093505050610e2d565b506000905060025b9250929050565b6000816004811115610e4857610e4861190d565b1415610e515750565b6001816004811115610e6557610e6561190d565b1415610eb35760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016102b5565b6002816004811115610ec757610ec761190d565b1415610f155760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016102b5565b6003816004811115610f2957610f2961190d565b1415610f825760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016102b5565b6004816004811115610f9657610f9661190d565b1415610fef5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b60648201526084016102b5565b50565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a083111561102957506000905060036110d6565b8460ff16601b1415801561104157508460ff16601c14155b1561105257506000905060046110d6565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa1580156110a6573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381166110cf576000600192509250506110d6565b9150600090505b94509492505050565b6000806001600160ff1b038316816110fc60ff86901c601b6118f5565b905061110a87828885610ff2565b935093505050935093915050565b604080516060808201835260008083526020808401829052845160a0810186528281529081018290528085018290529182018190526080820152909182015290565b82805461116690611798565b90600052602060002090601f01602090048101928261118857600085556111ce565b82601f106111a157805160ff19168380011785556111ce565b828001600101855582156111ce579182015b828111156111ce5782518255916020019190600101906111b3565b506111da929150611252565b5090565b8280546111ea90611798565b90600052602060002090601f01602090048101928261120c57600085556111ce565b82601f106112255782800160ff198235161785556111ce565b828001600101855582156111ce579182015b828111156111ce578235825591602001919060010190611237565b5b808211156111da5760008155600101611253565b634e487b7160e01b600052604160045260246000fd5b80356001600160a01b038116811461129457600080fd5b919050565b600060a082840312156112ab57600080fd5b60405160a0810181811067ffffffffffffffff821117156112ce576112ce611267565b806040525080915082358152602083013560208201526112f06040840161127d565b604082015260608301356060820152608083013560808201525092915050565b60008060c0838503121561132357600080fd5b823591506113348460208501611299565b90509250929050565b60008083601f84011261134f57600080fd5b50813567ffffffffffffffff81111561136757600080fd5b602083019150836020828501011115610e2d57600080fd5b600080600080600080600060c0888a03121561139a57600080fd5b87359650602088013595506113b16040890161127d565b9450606088013593506080880135925060a088013567ffffffffffffffff8111156113db57600080fd5b6113e78a828b0161133d565b989b979a50959850939692959293505050565b600060e0828403121561140c57600080fd5b6040516060810181811067ffffffffffffffff8211171561142f5761142f611267565b8060405250823581526020830135602082015261144f8460408501611299565b60408201529392505050565b60005b8381101561147657818101518382015260200161145e565b83811115611485576000848401525b50505050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156114f857878503603f19018452815180518087526114d9818989018a850161145b565b601f01601f1916959095018601945092850192908501906001016114b2565b5092979650505050505050565b60006020828403121561151757600080fd5b6115208261127d565b9392505050565b600082601f83011261153857600080fd5b813567ffffffffffffffff8082111561155357611553611267565b604051601f8301601f19908116603f0116810190828211818310171561157b5761157b611267565b8160405283815286602085880101111561159457600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080606085870312156115ca57600080fd5b6115d38561127d565b9350602085013567ffffffffffffffff808211156115f057600080fd5b6115fc8883890161133d565b9095509350604087013591508082111561161557600080fd5b5061162287828801611527565b91505092959194509250565b600060a0828403121561164057600080fd5b6115208383611299565b60006020828403121561165c57600080fd5b5035919050565b60008060008060008060c0878903121561167c57600080fd5b6116858761127d565b95506116936020880161127d565b9450604087013567ffffffffffffffff808211156116b057600080fd5b6116bc8a838b01611527565b955060608901359150808211156116d257600080fd5b6116de8a838b01611527565b945060808901359150808211156116f457600080fd5b5061170189828a01611527565b92505060a0870135801515811461171757600080fd5b809150509295509295509295565b6000806020838503121561173857600080fd5b823567ffffffffffffffff81111561174f57600080fd5b61175b8582860161133d565b90969095509350505050565b634e487b7160e01b600052601160045260246000fd5b600060001982141561179157611791611767565b5060010190565b600181811c908216806117ac57607f821691505b602082108114156117cd57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b60006bffffffffffffffffffffffff19808760601b168352808660601b16601484015250835161182081602885016020880161145b565b83519083019061183781602884016020880161145b565b016028019695505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a00000000000081526000835161187d81601a85016020880161145b565b83519083019061189481601a84016020880161145b565b01601a01949350505050565b634e487b7160e01b600052601260045260246000fd5b6000826118c5576118c56118a0565b500490565b6000828210156118dc576118dc611767565b500390565b6000826118f0576118f06118a0565b500690565b6000821982111561190857611908611767565b500190565b634e487b7160e01b600052602160045260246000fdfea2646970667358221220baaab66f3329a970b01d7a536b8fe48d7b58d508a7b63743cedea5f3bf60780464736f6c634300080c0033",
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

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_ManagementContract *ManagementContractCaller) GetHostAddresses(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "GetHostAddresses")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_ManagementContract *ManagementContractSession) GetHostAddresses() ([]string, error) {
	return _ManagementContract.Contract.GetHostAddresses(&_ManagementContract.CallOpts)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_ManagementContract *ManagementContractCallerSession) GetHostAddresses() ([]string, error) {
	return _ManagementContract.Contract.GetHostAddresses(&_ManagementContract.CallOpts)
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

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret, string _hostAddress) returns()
func (_ManagementContract *ManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, _aggregatorID common.Address, _initSecret []byte, _hostAddress string) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "InitializeNetworkSecret", _aggregatorID, _initSecret, _hostAddress)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret, string _hostAddress) returns()
func (_ManagementContract *ManagementContractSession) InitializeNetworkSecret(_aggregatorID common.Address, _initSecret []byte, _hostAddress string) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, _aggregatorID, _initSecret, _hostAddress)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret, string _hostAddress) returns()
func (_ManagementContract *ManagementContractTransactorSession) InitializeNetworkSecret(_aggregatorID common.Address, _initSecret []byte, _hostAddress string) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, _aggregatorID, _initSecret, _hostAddress)
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

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_ManagementContract *ManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_ManagementContract *ManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _ManagementContract.Contract.RespondNetworkSecret(&_ManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_ManagementContract *ManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _ManagementContract.Contract.RespondNetworkSecret(&_ManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
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
