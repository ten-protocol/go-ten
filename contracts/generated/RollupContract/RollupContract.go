// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package RollupContract

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

// StructsMetaRollup is an auto generated low-level Go binding around an user-defined struct.
type StructsMetaRollup struct {
	Hash               [32]byte
	LastSequenceNumber *big.Int
	BlockBindingHash   [32]byte
	BlockBindingNumber *big.Int
	CrossChainRoot     [32]byte
	LastBatchHash      [32]byte
	Signature          []byte
}

// RollupContractMetaData contains all meta data concerning the RollupContract contract.
var RollupContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"RollupAdded\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"addRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enclaveRegistry\",\"outputs\":[{\"internalType\":\"contractINetworkEnclaveRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"getRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_merkleMessageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_enclaveRegistry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b611536806100975f395ff3fe608060405234801561000f575f5ffd5b5060043610610085575f3560e01c80638da5cb5b116100585780638da5cb5b146100f0578063ae247b4914610128578063e874eb201461013b578063f2fde38b1461014e575f5ffd5b8063485cc955146100895780636fb6a45c1461009e578063715018a6146100c85780637c72dbd0146100d0575b5f5ffd5b61009c610097366004610c0a565b610161565b005b6100b16100ac366004610c51565b6102e3565b6040516100bf929190610d54565b60405180910390f35b61009c610416565b6004546100e3906001600160a01b031681565b6040516100bf9190610dbb565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516100bf9190610dd2565b61009c610136366004610df9565b610429565b6003546100e3906001600160a01b031681565b61009c61015c366004610e31565b6107af565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156101ab5750825b90505f8267ffffffffffffffff1660011480156101c75750303b155b9050811580156101d5575080155b1561020c576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561024057845468ff00000000000000001916680100000000000000001785555b61024933610805565b600380546001600160a01b03808a1673ffffffffffffffffffffffffffffffffffffffff199283161790925560048054928916929091169190911790555f60025583156102da57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906102d190600190610e68565b60405180910390a15b50505050505050565b5f61031e6040518060e001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081525090565b5f5f5f015f8581526020019081526020015f206040518060e00160405290815f8201548152602001600182015481526020016002820154815260200160038201548152602001600482015481526020016005820154815260200160068201805461038790610e8a565b80601f01602080910402602001604051908101604052809291908181526020018280546103b390610e8a565b80156103fe5780601f106103d5576101008083540402835291602001916103fe565b820191905f5260205f20905b8154815290600101906020018083116103e157829003601f168201915b50505091909252505081519095149590945092505050565b61041e610816565b6104275f61088a565b565b80806060013543116104565760405162461bcd60e51b815260040161044d90610eb0565b60405180910390fd5b610465606082013560ff610f25565b43106104835760405162461bcd60e51b815260040161044d90610f6c565b6060810135405f8190036104a95760405162461bcd60e51b815260040161044d90610fae565b816040013581146104cc5760405162461bcd60e51b815260040161044d90610ff0565b5f496104ea5760405162461bcd60e51b815260040161044d90611032565b5f82602001358360a001358460400135856060013586608001355f4960405160200161051b96959493929190611042565b60408051601f19818403018152919052805160209091012090505f6105808261054760c087018761109a565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061090792505050565b600480546040517f3c23afba0000000000000000000000000000000000000000000000000000000081529293506001600160a01b031691633c23afba916105c991859101610dd2565b602060405180830381865afa1580156105e4573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106089190611105565b6106245760405162461bcd60e51b815260040161044d90611154565b600480546040517f6d46e9870000000000000000000000000000000000000000000000000000000081526001600160a01b0390911691636d46e9879161066c91859101610dd2565b602060405180830381865afa158015610687573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106ab9190611105565b6106c75760405162461bcd60e51b815260040161044d90611196565b6106cf610816565b6106d885610931565b60808501355f1914610761576003546040517fb6aed0cb0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063b6aed0cb906107339060808901359042906004016111a6565b5f604051808303815f87803b15801561074a575f5ffd5b505af115801561075c573d5f5f3e3d5ffd5b505050505b7fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b12455f4961079160c088018861109a565b6040516107a0939291906111ec565b60405180910390a15050505050565b6107b7610816565b6001600160a01b0381166107f9575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161044d9190610dd2565b6108028161088a565b50565b61080d610965565b610802816109cc565b336108487f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461042757336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161044d9190610dd2565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f5f5f5f61091586866109d4565b9250925092506109258282610a1d565b50909150505b92915050565b80355f908152602081905260409020819061094c8282611496565b5050600254602082013511156108025760200135600255565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610427576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107b7610965565b5f5f5f8351604103610a0b576020840151604085015160608601515f1a6109fd88828585610b22565b955095509550505050610a16565b505081515f91506002905b9250925092565b5f826003811115610a3057610a306114a0565b03610a39575050565b6001826003811115610a4d57610a4d6114a0565b03610a84576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002826003811115610a9857610a986114a0565b03610ad1576040517ffce698f700000000000000000000000000000000000000000000000000000000815261044d9082906004016114b4565b6003826003811115610ae557610ae56114a0565b03610b1e57806040517fd78bce0c00000000000000000000000000000000000000000000000000000000815260040161044d91906114b4565b5050565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610b5b57505f91506003905082610bd2565b5f6001888888886040515f8152602001604052604051610b7e94939291906114cb565b6020604051602081039080840390855afa158015610b9e573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116610bc957505f925060019150829050610bd2565b92505f91508190505b9450945094915050565b5f6001600160a01b03821661092b565b610bf581610bdc565b8114610802575f5ffd5b803561092b81610bec565b5f5f60408385031215610c1e57610c1e5f5ffd5b610c288484610bff565b9150610c378460208501610bff565b90509250929050565b80610bf5565b803561092b81610c40565b5f60208284031215610c6457610c645f5ffd5b610c6e8383610c46565b9392505050565b8015155b82525050565b80610c79565b8281835e505f910152565b5f610c99825190565b808452602084019350610cb0818560208601610c85565b601f01601f19169290920192915050565b80515f9060e0840190610cd48582610c7f565b506020830151610ce76020860182610c7f565b506040830151610cfa6040860182610c7f565b506060830151610d0d6060860182610c7f565b506080830151610d206080860182610c7f565b5060a0830151610d3360a0860182610c7f565b5060c083015184820360c0860152610d4b8282610c90565b95945050505050565b60408101610d628285610c75565b8181036020830152610d748184610cc1565b949350505050565b5f61092b6001600160a01b038316610d92565b90565b6001600160a01b031690565b5f61092b82610d7c565b5f61092b82610d9e565b610c7981610da8565b6020810161092b8284610db2565b610c7981610bdc565b6020810161092b8284610dc9565b5f60e08284031215610df357610df35f5ffd5b50919050565b5f60208284031215610e0c57610e0c5f5ffd5b813567ffffffffffffffff811115610e2557610e255f5ffd5b610d7484828501610de0565b5f60208284031215610e4457610e445f5ffd5b610c6e8383610bff565b5f67ffffffffffffffff821661092b565b610c7981610e4e565b6020810161092b8284610e5f565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680610e9e57607f821691505b602082108103610df357610df3610e76565b6020808252810161092b81602681527f43616e6e6f742062696e6420746f20667574757265206f722063757272656e7460208201527f20626c6f636b0000000000000000000000000000000000000000000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561092b5761092b610f11565b60158152602081017f426c6f636b2062696e64696e6720746f6f206f6c640000000000000000000000815290505b60200190565b6020808252810161092b81610f38565b60128152602081017f556e6b6e6f776e20626c6f636b2068617368000000000000000000000000000081529050610f66565b6020808252810161092b81610f7c565b60168152602081017f426c6f636b2062696e64696e67206d69736d617463680000000000000000000081529050610f66565b6020808252810161092b81610fbe565b60148152602081017f426c6f622068617368206973206e6f742073657400000000000000000000000081529050610f66565b6020808252810161092b81611000565b61104c8188610c7f565b6020016110598187610c7f565b6020016110668186610c7f565b6020016110738185610c7f565b6020016110808184610c7f565b60200161108d8183610c7f565b6020019695505050505050565b5f808335601e19368590030181126110b3576110b35f5ffd5b8301915050803567ffffffffffffffff8111156110d1576110d15f5ffd5b6020820191506001810236038213156110eb576110eb5f5ffd5b9250929050565b801515610bf5565b805161092b816110f2565b5f60208284031215611118576111185f5ffd5b610c6e83836110fa565b60168152602081017f656e636c6176654944206e6f742061747465737465640000000000000000000081529050610f66565b6020808252810161092b81611122565b60198152602081017f656e636c6176654944206e6f7420612073657175656e6365720000000000000081529050610f66565b6020808252810161092b81611164565b604081016111b48285610c7f565b610c6e6020830184610c7f565b82818337505f910152565b8183526020830192506111e08284836111c1565b50601f01601f19160190565b604081016111fa8286610c7f565b8181036020830152610d4b8184866111cc565b5f813561092b81610c40565b5f8161092b565b61122982611219565b611235610d8f82611219565b8255505050565b5f61092b610d8f8381565b6112508261123c565b80611235565b634e487b7160e01b5f52604160045260245ffd5b6112738361123c565b81545f1960089490940293841b1916921b91909117905550565b5f61129981848461126a565b505050565b81811015610b1e576112b05f8261128d565b60010161129e565b601f821115611299575f818152602090206020601f850104810160208510156112de5750805b6112f06020601f86010483018261129e565b5050505050565b8267ffffffffffffffff81111561131057611310611256565b61131a8254610e8a565b6113258282856112b8565b505f601f821160018114611357575f83156113405750848201355b5f19600885021c19811660028502178555506113ae565b5f84815260208120601f198516915b828110156113865787850135825560209485019460019092019101611366565b50848210156113a2575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b6112998383836112f7565b81806113cc8161120d565b90506113d88184611220565b505060208201806113e88261120d565b90506113f78160018501611247565b505060408201806114078261120d565b90506114168160028501611220565b505060608201806114268261120d565b90506114358160038501611247565b505060808201806114458261120d565b90506114548160048501611220565b505060a08201806114648261120d565b90506114738160058501611220565b505061148260c083018361109a565b6114908183600686016113b6565b50505050565b610b1e82826113c1565b634e487b7160e01b5f52602160045260245ffd5b6020810161092b8284610c7f565b60ff8116610c79565b608081016114d98287610c7f565b6114e660208301866114c2565b6114f36040830185610c7f565b610d4b6060830184610c7f56fea26469706673582212206e18a0136a741e6051959e30c85320679de9d29df9a36babd1166b04db44880764736f6c634300081c0033",
}

// RollupContractABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupContractMetaData.ABI instead.
var RollupContractABI = RollupContractMetaData.ABI

// RollupContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RollupContractMetaData.Bin instead.
var RollupContractBin = RollupContractMetaData.Bin

// DeployRollupContract deploys a new Ethereum contract, binding an instance of RollupContract to it.
func DeployRollupContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RollupContract, error) {
	parsed, err := RollupContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RollupContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RollupContract{RollupContractCaller: RollupContractCaller{contract: contract}, RollupContractTransactor: RollupContractTransactor{contract: contract}, RollupContractFilterer: RollupContractFilterer{contract: contract}}, nil
}

// RollupContract is an auto generated Go binding around an Ethereum contract.
type RollupContract struct {
	RollupContractCaller     // Read-only binding to the contract
	RollupContractTransactor // Write-only binding to the contract
	RollupContractFilterer   // Log filterer for contract events
}

// RollupContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupContractSession struct {
	Contract     *RollupContract   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupContractCallerSession struct {
	Contract *RollupContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RollupContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupContractTransactorSession struct {
	Contract     *RollupContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RollupContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupContractRaw struct {
	Contract *RollupContract // Generic contract binding to access the raw methods on
}

// RollupContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupContractCallerRaw struct {
	Contract *RollupContractCaller // Generic read-only contract binding to access the raw methods on
}

// RollupContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupContractTransactorRaw struct {
	Contract *RollupContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollupContract creates a new instance of RollupContract, bound to a specific deployed contract.
func NewRollupContract(address common.Address, backend bind.ContractBackend) (*RollupContract, error) {
	contract, err := bindRollupContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RollupContract{RollupContractCaller: RollupContractCaller{contract: contract}, RollupContractTransactor: RollupContractTransactor{contract: contract}, RollupContractFilterer: RollupContractFilterer{contract: contract}}, nil
}

// NewRollupContractCaller creates a new read-only instance of RollupContract, bound to a specific deployed contract.
func NewRollupContractCaller(address common.Address, caller bind.ContractCaller) (*RollupContractCaller, error) {
	contract, err := bindRollupContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupContractCaller{contract: contract}, nil
}

// NewRollupContractTransactor creates a new write-only instance of RollupContract, bound to a specific deployed contract.
func NewRollupContractTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupContractTransactor, error) {
	contract, err := bindRollupContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupContractTransactor{contract: contract}, nil
}

// NewRollupContractFilterer creates a new log filterer instance of RollupContract, bound to a specific deployed contract.
func NewRollupContractFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupContractFilterer, error) {
	contract, err := bindRollupContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupContractFilterer{contract: contract}, nil
}

// bindRollupContract binds a generic wrapper to an already deployed contract.
func bindRollupContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RollupContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupContract *RollupContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupContract.Contract.RollupContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupContract *RollupContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupContract.Contract.RollupContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupContract *RollupContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupContract.Contract.RollupContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RollupContract *RollupContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RollupContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RollupContract *RollupContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RollupContract *RollupContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RollupContract.Contract.contract.Transact(opts, method, params...)
}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_RollupContract *RollupContractCaller) EnclaveRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollupContract.contract.Call(opts, &out, "enclaveRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_RollupContract *RollupContractSession) EnclaveRegistry() (common.Address, error) {
	return _RollupContract.Contract.EnclaveRegistry(&_RollupContract.CallOpts)
}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_RollupContract *RollupContractCallerSession) EnclaveRegistry() (common.Address, error) {
	return _RollupContract.Contract.EnclaveRegistry(&_RollupContract.CallOpts)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_RollupContract *RollupContractCaller) GetRollupByHash(opts *bind.CallOpts, rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	var out []interface{}
	err := _RollupContract.contract.Call(opts, &out, "getRollupByHash", rollupHash)

	if err != nil {
		return *new(bool), *new(StructsMetaRollup), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(StructsMetaRollup)).(*StructsMetaRollup)

	return out0, out1, err

}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_RollupContract *RollupContractSession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _RollupContract.Contract.GetRollupByHash(&_RollupContract.CallOpts, rollupHash)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_RollupContract *RollupContractCallerSession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _RollupContract.Contract.GetRollupByHash(&_RollupContract.CallOpts, rollupHash)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_RollupContract *RollupContractCaller) MerkleMessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollupContract.contract.Call(opts, &out, "merkleMessageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_RollupContract *RollupContractSession) MerkleMessageBus() (common.Address, error) {
	return _RollupContract.Contract.MerkleMessageBus(&_RollupContract.CallOpts)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_RollupContract *RollupContractCallerSession) MerkleMessageBus() (common.Address, error) {
	return _RollupContract.Contract.MerkleMessageBus(&_RollupContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollupContract *RollupContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RollupContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollupContract *RollupContractSession) Owner() (common.Address, error) {
	return _RollupContract.Contract.Owner(&_RollupContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RollupContract *RollupContractCallerSession) Owner() (common.Address, error) {
	return _RollupContract.Contract.Owner(&_RollupContract.CallOpts)
}

// AddRollup is a paid mutator transaction binding the contract method 0xae247b49.
//
// Solidity: function addRollup((bytes32,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_RollupContract *RollupContractTransactor) AddRollup(opts *bind.TransactOpts, r StructsMetaRollup) (*types.Transaction, error) {
	return _RollupContract.contract.Transact(opts, "addRollup", r)
}

// AddRollup is a paid mutator transaction binding the contract method 0xae247b49.
//
// Solidity: function addRollup((bytes32,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_RollupContract *RollupContractSession) AddRollup(r StructsMetaRollup) (*types.Transaction, error) {
	return _RollupContract.Contract.AddRollup(&_RollupContract.TransactOpts, r)
}

// AddRollup is a paid mutator transaction binding the contract method 0xae247b49.
//
// Solidity: function addRollup((bytes32,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_RollupContract *RollupContractTransactorSession) AddRollup(r StructsMetaRollup) (*types.Transaction, error) {
	return _RollupContract.Contract.AddRollup(&_RollupContract.TransactOpts, r)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry) returns()
func (_RollupContract *RollupContractTransactor) Initialize(opts *bind.TransactOpts, _merkleMessageBus common.Address, _enclaveRegistry common.Address) (*types.Transaction, error) {
	return _RollupContract.contract.Transact(opts, "initialize", _merkleMessageBus, _enclaveRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry) returns()
func (_RollupContract *RollupContractSession) Initialize(_merkleMessageBus common.Address, _enclaveRegistry common.Address) (*types.Transaction, error) {
	return _RollupContract.Contract.Initialize(&_RollupContract.TransactOpts, _merkleMessageBus, _enclaveRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry) returns()
func (_RollupContract *RollupContractTransactorSession) Initialize(_merkleMessageBus common.Address, _enclaveRegistry common.Address) (*types.Transaction, error) {
	return _RollupContract.Contract.Initialize(&_RollupContract.TransactOpts, _merkleMessageBus, _enclaveRegistry)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollupContract *RollupContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RollupContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollupContract *RollupContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _RollupContract.Contract.RenounceOwnership(&_RollupContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RollupContract *RollupContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RollupContract.Contract.RenounceOwnership(&_RollupContract.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollupContract *RollupContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RollupContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollupContract *RollupContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RollupContract.Contract.TransferOwnership(&_RollupContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RollupContract *RollupContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RollupContract.Contract.TransferOwnership(&_RollupContract.TransactOpts, newOwner)
}

// RollupContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the RollupContract contract.
type RollupContractInitializedIterator struct {
	Event *RollupContractInitialized // Event containing the contract specifics and raw log

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
func (it *RollupContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupContractInitialized)
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
		it.Event = new(RollupContractInitialized)
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
func (it *RollupContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupContractInitialized represents a Initialized event raised by the RollupContract contract.
type RollupContractInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RollupContract *RollupContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*RollupContractInitializedIterator, error) {

	logs, sub, err := _RollupContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RollupContractInitializedIterator{contract: _RollupContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_RollupContract *RollupContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RollupContractInitialized) (event.Subscription, error) {

	logs, sub, err := _RollupContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupContractInitialized)
				if err := _RollupContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_RollupContract *RollupContractFilterer) ParseInitialized(log types.Log) (*RollupContractInitialized, error) {
	event := new(RollupContractInitialized)
	if err := _RollupContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RollupContract contract.
type RollupContractOwnershipTransferredIterator struct {
	Event *RollupContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RollupContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupContractOwnershipTransferred)
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
		it.Event = new(RollupContractOwnershipTransferred)
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
func (it *RollupContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupContractOwnershipTransferred represents a OwnershipTransferred event raised by the RollupContract contract.
type RollupContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RollupContract *RollupContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RollupContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollupContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RollupContractOwnershipTransferredIterator{contract: _RollupContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RollupContract *RollupContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RollupContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RollupContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupContractOwnershipTransferred)
				if err := _RollupContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_RollupContract *RollupContractFilterer) ParseOwnershipTransferred(log types.Log) (*RollupContractOwnershipTransferred, error) {
	event := new(RollupContractOwnershipTransferred)
	if err := _RollupContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupContractRollupAddedIterator is returned from FilterRollupAdded and is used to iterate over the raw logs and unpacked data for RollupAdded events raised by the RollupContract contract.
type RollupContractRollupAddedIterator struct {
	Event *RollupContractRollupAdded // Event containing the contract specifics and raw log

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
func (it *RollupContractRollupAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupContractRollupAdded)
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
		it.Event = new(RollupContractRollupAdded)
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
func (it *RollupContractRollupAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupContractRollupAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupContractRollupAdded represents a RollupAdded event raised by the RollupContract contract.
type RollupContractRollupAdded struct {
	RollupHash [32]byte
	Signature  []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRollupAdded is a free log retrieval operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_RollupContract *RollupContractFilterer) FilterRollupAdded(opts *bind.FilterOpts) (*RollupContractRollupAddedIterator, error) {

	logs, sub, err := _RollupContract.contract.FilterLogs(opts, "RollupAdded")
	if err != nil {
		return nil, err
	}
	return &RollupContractRollupAddedIterator{contract: _RollupContract.contract, event: "RollupAdded", logs: logs, sub: sub}, nil
}

// WatchRollupAdded is a free log subscription operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_RollupContract *RollupContractFilterer) WatchRollupAdded(opts *bind.WatchOpts, sink chan<- *RollupContractRollupAdded) (event.Subscription, error) {

	logs, sub, err := _RollupContract.contract.WatchLogs(opts, "RollupAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupContractRollupAdded)
				if err := _RollupContract.contract.UnpackLog(event, "RollupAdded", log); err != nil {
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

// ParseRollupAdded is a log parse operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_RollupContract *RollupContractFilterer) ParseRollupAdded(log types.Log) (*RollupContractRollupAdded, error) {
	event := new(RollupContractRollupAdded)
	if err := _RollupContract.contract.UnpackLog(event, "RollupAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
