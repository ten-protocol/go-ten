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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e0fd84bd": "AddRollup(bytes32,address,bytes32,uint256,string)",
		"c719bf50": "InitializeNetworkSecret(address,bytes)",
		"e34fbfc8": "RequestNetworkSecret(string)",
		"981214ba": "RespondNetworkSecret(address,address,bytes,bytes)",
		"8ef74f89": "attestationRequests(address)",
		"d4c80664": "attested(address)",
		"e0643dfc": "rollups(uint256,uint256)",
	},
	Bin: "0x608060405234801561001057600080fd5b506112a9806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063d4c806641161005b578063d4c80664146100d3578063e0643dfc14610106578063e0fd84bd14610141578063e34fbfc81461015457600080fd5b80638ef74f8914610082578063981214ba146100ab578063c719bf50146100c0575b600080fd5b610095610090366004610b46565b610167565b6040516100a29190610b98565b60405180910390f35b6100be6100b9366004610c6e565b610201565b005b6100be6100ce366004610d35565b610300565b6100f66100e1366004610b46565b60026020526000908152604090205460ff1681565b60405190151581526020016100a2565b610119610114366004610d88565b610347565b604080519485526001600160a01b0390931660208501529183015260608201526080016100a2565b6100be61014f366004610daa565b610396565b6100be610162366004610e1b565b61043a565b6001602052600090815260409020805461018090610e5d565b80601f01602080910402602001604051908101604052809291908181526020018280546101ac90610e5d565b80156101f95780601f106101ce576101008083540402835291602001916101f9565b820191906000526020600020905b8154815290600101906020018083116101dc57829003601f168201915b505050505081565b6001600160a01b03841660009081526002602052604090205460ff168061022757600080fd5b600061025586868560405160200161024193929190610e98565b604051602081830303815290604052610459565b905060006102638286610494565b9050866001600160a01b0316816001600160a01b031614610283886104ba565b61028c836104ba565b60405160200161029d929190610ede565b604051602081830303815290604052906102d35760405162461bcd60e51b81526004016102ca9190610b98565b60405180910390fd5b5050506001600160a01b039093166000908152600260205260409020805460ff1916600117905550505050565b60035460ff161561031057600080fd5b50506003805460ff1990811660019081179092556001600160a01b0390921660009081526002602052604090208054909216179055565b6000602052816000526040600020818154811061036357600080fd5b600091825260209091206004909102018054600182015460028301546003909301549194506001600160a01b0316925084565b6001600160a01b03851660009081526002602052604090205460ff166103bb57600080fd5b5050604080516080810182529485526001600160a01b0393841660208087019182528683019485526060870193845243600090815280825292832080546001808201835591855291909320965160049091029096019586555190850180546001600160a01b031916919094161790925551600283015551600390910155565b336000908152600160205260409020610454908383610a96565b505050565b60006104658251610601565b82604051602001610477929190610fb4565b604051602081830303815290604052805190602001209050919050565b60008060006104a38585610707565b915091506104b081610777565b5090505b92915050565b60408051602880825260608281019093526000919060208201818036833701905050905060005b60148110156105fa5760006104f7826013611025565b61050290600861103c565b61050d90600261113f565b610520906001600160a01b038716611161565b60f81b9050600060108260f81c6105379190611175565b60f81b905060008160f81c601061054e9190611197565b8360f81c61055c91906111b8565b60f81b905061056a82610935565b8561057686600261103c565b81518110610586576105866111db565b60200101906001600160f81b031916908160001a9053506105a681610935565b856105b286600261103c565b6105bd9060016111f1565b815181106105cd576105cd6111db565b60200101906001600160f81b031916908160001a90535050505080806105f290611209565b9150506104e1565b5092915050565b6060816106255750506040805180820190915260018152600360fc1b602082015290565b8160005b811561064f578061063981611209565b91506106489050600a83611161565b9150610629565b60008167ffffffffffffffff81111561066a5761066a610bcb565b6040519080825280601f01601f191660200182016040528015610694576020820181803683370190505b5090505b84156106ff576106a9600183611025565b91506106b6600a86611224565b6106c19060306111f1565b60f81b8183815181106106d6576106d66111db565b60200101906001600160f81b031916908160001a9053506106f8600a86611161565b9450610698565b949350505050565b60008082516041141561073e5760208301516040840151606085015160001a61073287828585610970565b94509450505050610770565b825160401415610768576020830151604084015161075d868383610a5d565b935093505050610770565b506000905060025b9250929050565b600081600481111561078b5761078b611238565b14156107945750565b60018160048111156107a8576107a8611238565b14156107f65760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016102ca565b600281600481111561080a5761080a611238565b14156108585760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016102ca565b600381600481111561086c5761086c611238565b14156108c55760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016102ca565b60048160048111156108d9576108d9611238565b14156109325760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b60648201526084016102ca565b50565b6000600a60f883901c101561095c5761095360f883901c603061124e565b60f81b92915050565b61095360f883901c605761124e565b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08311156109a75750600090506003610a54565b8460ff16601b141580156109bf57508460ff16601c14155b156109d05750600090506004610a54565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015610a24573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610a4d57600060019250925050610a54565b9150600090505b94509492505050565b6000806001600160ff1b03831681610a7a60ff86901c601b6111f1565b9050610a8887828885610970565b935093505050935093915050565b828054610aa290610e5d565b90600052602060002090601f016020900481019282610ac45760008555610b0a565b82601f10610add5782800160ff19823516178555610b0a565b82800160010185558215610b0a579182015b82811115610b0a578235825591602001919060010190610aef565b50610b16929150610b1a565b5090565b5b80821115610b165760008155600101610b1b565b80356001600160a01b038116811461096b57600080fd5b600060208284031215610b5857600080fd5b610b6182610b2f565b9392505050565b60005b83811015610b83578181015183820152602001610b6b565b83811115610b92576000848401525b50505050565b6020815260008251806020840152610bb7816040850160208701610b68565b601f01601f19169190910160400192915050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112610bf257600080fd5b813567ffffffffffffffff80821115610c0d57610c0d610bcb565b604051601f8301601f19908116603f01168101908282118183101715610c3557610c35610bcb565b81604052838152866020858801011115610c4e57600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008060008060808587031215610c8457600080fd5b610c8d85610b2f565b9350610c9b60208601610b2f565b9250604085013567ffffffffffffffff80821115610cb857600080fd5b610cc488838901610be1565b93506060870135915080821115610cda57600080fd5b50610ce787828801610be1565b91505092959194509250565b60008083601f840112610d0557600080fd5b50813567ffffffffffffffff811115610d1d57600080fd5b60208301915083602082850101111561077057600080fd5b600080600060408486031215610d4a57600080fd5b610d5384610b2f565b9250602084013567ffffffffffffffff811115610d6f57600080fd5b610d7b86828701610cf3565b9497909650939450505050565b60008060408385031215610d9b57600080fd5b50508035926020909101359150565b60008060008060008060a08789031215610dc357600080fd5b86359550610dd360208801610b2f565b94506040870135935060608701359250608087013567ffffffffffffffff811115610dfd57600080fd5b610e0989828a01610cf3565b979a9699509497509295939492505050565b60008060208385031215610e2e57600080fd5b823567ffffffffffffffff811115610e4557600080fd5b610e5185828601610cf3565b90969095509350505050565b600181811c90821680610e7157607f821691505b60208210811415610e9257634e487b7160e01b600052602260045260246000fd5b50919050565b60006bffffffffffffffffffffffff19808660601b168352808560601b166014840152508251610ecf816028850160208701610b68565b91909101602801949350505050565b7f7265636f7665726564206164647265737320616e64206174746573746572494481526b0103237b73a1036b0ba31b4160a51b60208201527f0a2045787065637465643a202020202020202020202020202020202020202020602c820152630101010160e51b604c82015260008351610f5e816050850160208801610b68565b7f0a202f207265636f7665726564416464725369676e656443616c63756c617465605091840191820152630321d10160e51b60708201528351610fa8816074840160208801610b68565b01607401949350505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000815260008351610fec81601a850160208801610b68565b83519083019061100381601a840160208801610b68565b01601a01949350505050565b634e487b7160e01b600052601160045260246000fd5b6000828210156110375761103761100f565b500390565b60008160001904831182151516156110565761105661100f565b500290565b600181815b8085111561109657816000190482111561107c5761107c61100f565b8085161561108957918102915b93841c9390800290611060565b509250929050565b6000826110ad575060016104b4565b816110ba575060006104b4565b81600181146110d057600281146110da576110f6565b60019150506104b4565b60ff8411156110eb576110eb61100f565b50506001821b6104b4565b5060208310610133831016604e8410600b8410161715611119575081810a6104b4565b611123838361105b565b80600019048211156111375761113761100f565b029392505050565b6000610b61838361109e565b634e487b7160e01b600052601260045260246000fd5b6000826111705761117061114b565b500490565b600060ff8316806111885761118861114b565b8060ff84160491505092915050565b600060ff821660ff84168160ff04811182151516156111375761113761100f565b600060ff821660ff8416808210156111d2576111d261100f565b90039392505050565b634e487b7160e01b600052603260045260246000fd5b600082198211156112045761120461100f565b500190565b600060001982141561121d5761121d61100f565b5060010190565b6000826112335761123361114b565b500690565b634e487b7160e01b600052602160045260246000fd5b600060ff821660ff84168060ff0382111561126b5761126b61100f565b01939250505056fea26469706673582212203a32da7221ac2cdc6baeacf19980dbfa9f9f81d643bc6c76875aebcde57bd66264736f6c634300080c0033",
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
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_ManagementContract *ManagementContractCaller) Rollups(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "rollups", arg0, arg1)

	outstruct := new(struct {
		ParentHash   [32]byte
		AggregatorID common.Address
		L1Block      [32]byte
		Number       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ParentHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.AggregatorID = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.L1Block = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Number = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_ManagementContract *ManagementContractSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _ManagementContract.Contract.Rollups(&_ManagementContract.CallOpts, arg0, arg1)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_ManagementContract *ManagementContractCallerSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _ManagementContract.Contract.Rollups(&_ManagementContract.CallOpts, arg0, arg1)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_ManagementContract *ManagementContractTransactor) AddRollup(opts *bind.TransactOpts, ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "AddRollup", ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_ManagementContract *ManagementContractSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _ManagementContract.Contract.AddRollup(&_ManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_ManagementContract *ManagementContractTransactorSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _ManagementContract.Contract.AddRollup(&_ManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_ManagementContract *ManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_ManagementContract *ManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_ManagementContract *ManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, aggregatorID, initSecret)
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
