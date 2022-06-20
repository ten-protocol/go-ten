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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_parentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_l1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_parentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"_r\",\"type\":\"tuple\"}],\"name\":\"AppendRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"element\",\"type\":\"tuple\"}],\"name\":\"GetParentRollup\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"GetRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"rollupID\",\"type\":\"uint256\"}],\"name\":\"GetRollupByID\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"ElementID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ParentID\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"rollup\",\"type\":\"tuple\"}],\"internalType\":\"structManagementContract.TreeElement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HasSecondCousinFork\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"internalType\":\"structManagementContract.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"InitializeTree\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"IsWithdrawalAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tree\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_TAIL\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_HEAD\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_nextID\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"initialized\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1b1e5402": "AddRollup(bytes32,bytes32,address,bytes32,uint256,string)",
		"1ac63930": "AppendRollup(uint256,(bytes32,bytes32,address,bytes32,uint256))",
		"31b1d255": "GetParentRollup((uint256,uint256,(bytes32,bytes32,address,bytes32,uint256)))",
		"8236a7ba": "GetRollupByHash(bytes32)",
		"92aaec79": "GetRollupByID(uint256)",
		"57b70600": "HasSecondCousinFork()",
		"c719bf50": "InitializeNetworkSecret(address,bytes)",
		"73bba846": "InitializeTree((bytes32,bytes32,address,bytes32,uint256))",
		"a52f433c": "IsWithdrawalAvailable()",
		"e34fbfc8": "RequestNetworkSecret(string)",
		"981214ba": "RespondNetworkSecret(address,address,bytes,bytes)",
		"8ef74f89": "attestationRequests(address)",
		"d4c80664": "attested(address)",
		"e0643dfc": "rollups(uint256,uint256)",
		"fd54b228": "tree()",
	},
	Bin: "0x608060405234801561001057600080fd5b506118ed806100206000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806392aaec7911610097578063d4c8066411610066578063d4c8066414610244578063e0643dfc14610267578063e34fbfc8146102b0578063fd54b228146102c357600080fd5b806392aaec79146101fb578063981214ba1461020e578063a52f433c14610221578063c719bf501461023157600080fd5b806357b70600116100d357806357b706001461019d57806373bba846146101b55780638236a7ba146101c85780638ef74f89146101db57600080fd5b80631ac63930146100fa5780631b1e54021461010f57806331b1d25514610122575b600080fd5b61010d610108366004611343565b6102fe565b005b61010d61011d3660046113b2565b610435565b61013561013036600461142d565b6105ee565b6040805192151583528151602080850191909152808301518483015291810151805160608086019190915292810151608080860191909152918101516001600160a01b031660a08501529182015160c0840152015160e0820152610100015b60405180910390f35b6101a561060e565b6040519015158152602001610194565b61010d6101c336600461148e565b6107b6565b6101356101d63660046114b1565b61098c565b6101ee6101e93660046114ca565b6109ae565b6040516101949190611515565b6101356102093660046114b1565b610a48565b61010d61021c3660046115d5565b610ac8565b600a54610100900460ff166101a5565b61010d61023f36600461165a565b610bcd565b6101a56102523660046114ca565b60026020526000908152604090205460ff1681565b61027a6102753660046116ad565b610c14565b6040805195865260208601949094526001600160a01b03909216928401929092526060830191909152608082015260a001610194565b61010d6102be3660046116cf565b610c6b565b6006546007546008546009546102dc9392919060ff1684565b6040805194855260208501939093529183015215156060820152608001610194565b60088054908190600061031083611727565b919050555060008061032185610a48565b915091508161036a5760405162461bcd60e51b815260206004820152601060248201526f1c185c995b9d081b9bdd08199bdd5b9960821b60448201526064015b60405180910390fd5b6040805160608082018352858252602080830189815283850189815260008981526003808552878220965187559251600180880191909155915180516002880155808501519387019390935582870151600480880180546001600160a01b0319166001600160a01b0390931692909217909155948301516005808801919091556080909301516006909601959095558a85529082528484208054918201815584528184200187905587810151835252208390556007548151141561042e5760078390555b5050505050565b6001600160a01b03851660009081526002602052604090205460ff1661049d5760405162461bcd60e51b815260206004820152601760248201527f61676772656761746f72206e6f742061747465737465640000000000000000006044820152606401610361565b6040805160a08101825288815260208082018981526001600160a01b03898116848601908152606085018a8152608086018a8152436000908152808752978820805460018082018355918a52969098208751600590970201958655935196850196909655516002840180546001600160a01b0319169190921617905592516003820155915160049092019190915560095460ff166105445761053e816107b6565b506105e5565b6000806105508a61098c565b91509150816105a15760405162461bcd60e51b815260206004820152601a60248201527f756e61626c6520746f2066696e6420706172656e7420686173680000000000006044820152606401610361565b600754600210156105d55760006105b661060e565b905080156105d3575050600a805461ff0019169055506105e59050565b505b80516105e190846102fe565b5050505b50505050505050565b60006105f86111bf565b6106058360200151610a48565b91509150915091565b600080610619610c8a565b9050600080610627836105ee565b91509150816106645760405162461bcd60e51b81526020600482015260096024820152681b9bc81c185c995b9d60ba1b6044820152606401610361565b600080610670836105ee565b91509150816106b35760405162461bcd60e51b815260206004820152600f60248201526e1b9bc819dc985b99081c185c995b9d608a1b6044820152606401610361565b805160009081526005602090815260408083208054825181850281018501909352808352919290919083018282801561070b57602002820191906000526020600020905b8154815260200190600101908083116106f7575b5050505050905060005b81518110156107a85760008061074384848151811061073657610736611742565b6020026020010151610a48565b915091508161075c576000995050505050505050505090565b86518151141561076d575050610796565b805160009081526005602052604090205415610793576001995050505050505050505090565b50505b806107a081611727565b915050610715565b506000965050505050505090565b60095460ff16156108095760405162461bcd60e51b815260206004820152601b60248201527f63616e6e6f7420626520696e697469616c697a656420616761696e00000000006044820152606401610361565b6009805460ff191660019081179091556040805160608082018352838252600060208084018281528486018881528784526003835294517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c55517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054d55925180517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054e55808401517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054f55808501517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3055080546001600160a01b039092166001600160a01b0319909216919091179055918201517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c30551556080909101517fa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3055255600784905560026008559381015184526004905290912055600a805461ff001916610100179055565b60006109966111bf565b60008381526004602052604090205461060590610a48565b600160205260009081526040902080546109c790611758565b80601f01602080910402602001604051908101604052809291908181526020018280546109f390611758565b8015610a405780601f10610a1557610100808354040283529160200191610a40565b820191906000526020600020905b815481529060010190602001808311610a2357829003601f168201915b505050505081565b6000610a526111bf565b505060009081526003602081815260409283902083516060808201865282548252600183015482850152855160a08101875260028401548152948301549385019390935260048201546001600160a01b031684860152600582015492840192909252600601546080830152918201528051151591565b6001600160a01b03841660009081526002602052604090205460ff1680610aee57600080fd5b6000610b1c868685604051602001610b0893929190611793565b604051602081830303815290604052610d06565b90506000610b2a8286610d41565b9050866001600160a01b0316816001600160a01b031614610ba15760405162461bcd60e51b815260206004820152602b60248201527f7265636f7665726564206164647265737320616e64206174746573746572494460448201526a040c8dedce840dac2e8c6d60ab1b6064820152608401610361565b5050506001600160a01b039092166000908152600260205260409020805460ff19166001179055505050565b600a5460ff1615610bdd57600080fd5b5050600a805460ff1990811660019081179092556001600160a01b0390921660009081526002602052604090208054909216179055565b60006020528160005260406000208181548110610c3057600080fd5b6000918252602090912060059091020180546001820154600283015460038401546004909401549295509093506001600160a01b0316919085565b336000908152600160205260409020610c85908383611201565b505050565b610c926111bf565b5060075460009081526003602081815260409283902083516060808201865282548252600183015482850152855160a08101875260028401548152948301549385019390935260048201546001600160a01b0316848601526005820154928401929092526006015460808301529182015290565b6000610d128251610d65565b82604051602001610d249291906117d9565b604051602081830303815290604052805190602001209050919050565b6000806000610d508585610e6b565b91509150610d5d81610edb565b509392505050565b606081610d895750506040805180820190915260018152600360fc1b602082015290565b8160005b8115610db35780610d9d81611727565b9150610dac9050600a8361184a565b9150610d8d565b60008167ffffffffffffffff811115610dce57610dce61129a565b6040519080825280601f01601f191660200182016040528015610df8576020820181803683370190505b5090505b8415610e6357610e0d60018361185e565b9150610e1a600a86611875565b610e25906030611889565b60f81b818381518110610e3a57610e3a611742565b60200101906001600160f81b031916908160001a905350610e5c600a8661184a565b9450610dfc565b949350505050565b600080825160411415610ea25760208301516040840151606085015160001a610e9687828585611099565b94509450505050610ed4565b825160401415610ecc5760208301516040840151610ec1868383611186565b935093505050610ed4565b506000905060025b9250929050565b6000816004811115610eef57610eef6118a1565b1415610ef85750565b6001816004811115610f0c57610f0c6118a1565b1415610f5a5760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e617475726500000000000000006044820152606401610361565b6002816004811115610f6e57610f6e6118a1565b1415610fbc5760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e677468006044820152606401610361565b6003816004811115610fd057610fd06118a1565b14156110295760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b6064820152608401610361565b600481600481111561103d5761103d6118a1565b14156110965760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b6064820152608401610361565b50565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08311156110d0575060009050600361117d565b8460ff16601b141580156110e857508460ff16601c14155b156110f9575060009050600461117d565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa15801561114d573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381166111765760006001925092505061117d565b9150600090505b94509492505050565b6000806001600160ff1b038316816111a360ff86901c601b611889565b90506111b187828885611099565b935093505050935093915050565b604080516060808201835260008083526020808401829052845160a0810186528281529081018290528085018290529182018190526080820152909182015290565b82805461120d90611758565b90600052602060002090601f01602090048101928261122f5760008555611275565b82601f106112485782800160ff19823516178555611275565b82800160010185558215611275579182015b8281111561127557823582559160200191906001019061125a565b50611281929150611285565b5090565b5b808211156112815760008155600101611286565b634e487b7160e01b600052604160045260246000fd5b80356001600160a01b03811681146112c757600080fd5b919050565b600060a082840312156112de57600080fd5b60405160a0810181811067ffffffffffffffff821117156113015761130161129a565b80604052508091508235815260208301356020820152611323604084016112b0565b604082015260608301356060820152608083013560808201525092915050565b60008060c0838503121561135657600080fd5b8235915061136784602085016112cc565b90509250929050565b60008083601f84011261138257600080fd5b50813567ffffffffffffffff81111561139a57600080fd5b602083019150836020828501011115610ed457600080fd5b600080600080600080600060c0888a0312156113cd57600080fd5b87359650602088013595506113e4604089016112b0565b9450606088013593506080880135925060a088013567ffffffffffffffff81111561140e57600080fd5b61141a8a828b01611370565b989b979a50959850939692959293505050565b600060e0828403121561143f57600080fd5b6040516060810181811067ffffffffffffffff821117156114625761146261129a565b8060405250823581526020830135602082015261148284604085016112cc565b60408201529392505050565b600060a082840312156114a057600080fd5b6114aa83836112cc565b9392505050565b6000602082840312156114c357600080fd5b5035919050565b6000602082840312156114dc57600080fd5b6114aa826112b0565b60005b838110156115005781810151838201526020016114e8565b8381111561150f576000848401525b50505050565b60208152600082518060208401526115348160408501602087016114e5565b601f01601f19169190910160400192915050565b600082601f83011261155957600080fd5b813567ffffffffffffffff808211156115745761157461129a565b604051601f8301601f19908116603f0116810190828211818310171561159c5761159c61129a565b816040528381528660208588010111156115b557600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080608085870312156115eb57600080fd5b6115f4856112b0565b9350611602602086016112b0565b9250604085013567ffffffffffffffff8082111561161f57600080fd5b61162b88838901611548565b9350606087013591508082111561164157600080fd5b5061164e87828801611548565b91505092959194509250565b60008060006040848603121561166f57600080fd5b611678846112b0565b9250602084013567ffffffffffffffff81111561169457600080fd5b6116a086828701611370565b9497909650939450505050565b600080604083850312156116c057600080fd5b50508035926020909101359150565b600080602083850312156116e257600080fd5b823567ffffffffffffffff8111156116f957600080fd5b61170585828601611370565b90969095509350505050565b634e487b7160e01b600052601160045260246000fd5b600060001982141561173b5761173b611711565b5060010190565b634e487b7160e01b600052603260045260246000fd5b600181811c9082168061176c57607f821691505b6020821081141561178d57634e487b7160e01b600052602260045260246000fd5b50919050565b60006bffffffffffffffffffffffff19808660601b168352808560601b1660148401525082516117ca8160288501602087016114e5565b91909101602801949350505050565b7f19457468657265756d205369676e6564204d6573736167653a0a00000000000081526000835161181181601a8501602088016114e5565b83519083019061182881601a8401602088016114e5565b01601a01949350505050565b634e487b7160e01b600052601260045260246000fd5b60008261185957611859611834565b500490565b60008282101561187057611870611711565b500390565b60008261188457611884611834565b500690565b6000821982111561189c5761189c611711565b500190565b634e487b7160e01b600052602160045260246000fdfea26469706673582212209de6087d5178366bf494c48603c2c170d81cb1447403bee3d1d53245acec37a564736f6c634300080c0033",
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

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_ManagementContract *ManagementContractCaller) AttestationRequests(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "attestationRequests", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_ManagementContract *ManagementContractSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _ManagementContract.Contract.AttestationRequests(&_ManagementContract.CallOpts, arg0)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_ManagementContract *ManagementContractCallerSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _ManagementContract.Contract.AttestationRequests(&_ManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_ManagementContract *ManagementContractCaller) Attested(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "attested", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_ManagementContract *ManagementContractSession) Attested(arg0 common.Address) (bool, error) {
	return _ManagementContract.Contract.Attested(&_ManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_ManagementContract *ManagementContractCallerSession) Attested(arg0 common.Address) (bool, error) {
	return _ManagementContract.Contract.Attested(&_ManagementContract.CallOpts, arg0)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, bytes32 Hash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_ManagementContract *ManagementContractCaller) Rollups(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	Hash         [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "rollups", arg0, arg1)

	outstruct := new(struct {
		ParentHash   [32]byte
		Hash         [32]byte
		AggregatorID common.Address
		L1Block      [32]byte
		Number       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ParentHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Hash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.AggregatorID = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.L1Block = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Number = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, bytes32 Hash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_ManagementContract *ManagementContractSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	Hash         [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _ManagementContract.Contract.Rollups(&_ManagementContract.CallOpts, arg0, arg1)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, bytes32 Hash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_ManagementContract *ManagementContractCallerSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	Hash         [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _ManagementContract.Contract.Rollups(&_ManagementContract.CallOpts, arg0, arg1)
}

// Tree is a free data retrieval call binding the contract method 0xfd54b228.
//
// Solidity: function tree() view returns(uint256 _TAIL, uint256 _HEAD, uint256 _nextID, bool initialized)
func (_ManagementContract *ManagementContractCaller) Tree(opts *bind.CallOpts) (struct {
	TAIL        *big.Int
	HEAD        *big.Int
	NextID      *big.Int
	Initialized bool
}, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "tree")

	outstruct := new(struct {
		TAIL        *big.Int
		HEAD        *big.Int
		NextID      *big.Int
		Initialized bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TAIL = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.HEAD = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NextID = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Initialized = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Tree is a free data retrieval call binding the contract method 0xfd54b228.
//
// Solidity: function tree() view returns(uint256 _TAIL, uint256 _HEAD, uint256 _nextID, bool initialized)
func (_ManagementContract *ManagementContractSession) Tree() (struct {
	TAIL        *big.Int
	HEAD        *big.Int
	NextID      *big.Int
	Initialized bool
}, error) {
	return _ManagementContract.Contract.Tree(&_ManagementContract.CallOpts)
}

// Tree is a free data retrieval call binding the contract method 0xfd54b228.
//
// Solidity: function tree() view returns(uint256 _TAIL, uint256 _HEAD, uint256 _nextID, bool initialized)
func (_ManagementContract *ManagementContractCallerSession) Tree() (struct {
	TAIL        *big.Int
	HEAD        *big.Int
	NextID      *big.Int
	Initialized bool
}, error) {
	return _ManagementContract.Contract.Tree(&_ManagementContract.CallOpts)
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
