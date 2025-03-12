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
	Hash                [32]byte
	FirstSequenceNumber *big.Int
	LastSequenceNumber  *big.Int
	BlockBindingHash    [32]byte
	BlockBindingNumber  *big.Int
	CrossChainRoot      [32]byte
	LastBatchHash       [32]byte
	Signature           []byte
}

// RollupContractMetaData contains all meta data concerning the RollupContract contract.
var RollupContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"RollupAdded\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"addRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enclaveRegistry\",\"outputs\":[{\"internalType\":\"contractINetworkEnclaveRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"getRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_merkleMessageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_enclaveRegistry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSeqNo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b6115aa806100975f395ff3fe608060405234801561000f575f5ffd5b506004361061009f575f3560e01c8063715018a6116100725780638da5cb5b116100585780638da5cb5b14610133578063e874eb201461016b578063f2fde38b1461017e575f5ffd5b8063715018a61461010b5780637c72dbd014610113575f5ffd5b8063440c953b146100a3578063485cc955146100c25780635fdf31a2146100d75780636fb6a45c146100ea575b5f5ffd5b6100ac60025481565b6040516100b99190610c1e565b60405180910390f35b6100d56100d0366004610c5a565b610191565b005b6100d56100e5366004610caa565b610313565b6100fd6100f8366004610cfb565b610691565b6040516100b9929190610e0a565b6100d56107d6565b600454610126906001600160a01b031681565b6040516100b99190610e69565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516100b99190610e80565b600354610126906001600160a01b031681565b6100d561018c366004610e8e565b6107e9565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156101db5750825b90505f8267ffffffffffffffff1660011480156101f75750303b155b905081158015610205575080155b1561023c576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561027057845468ff00000000000000001916680100000000000000001785555b6102793361083f565b600380546001600160a01b03808a1673ffffffffffffffffffffffffffffffffffffffff199283161790925560048054928916929091169190911790555f600255831561030a57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061030190600190610ec5565b60405180910390a15b50505050505050565b80806080013543116103405760405162461bcd60e51b815260040161033790610ed3565b60405180910390fd5b61034f608082013560ff610f48565b431061036d5760405162461bcd60e51b815260040161033790610f8f565b6080810135405f8190036103935760405162461bcd60e51b815260040161033790610fd1565b816060013581146103b65760405162461bcd60e51b815260040161033790611013565b5f496103d45760405162461bcd60e51b815260040161033790611055565b5f82604001358360c00135846060013585608001358660a001355f496040516020016104059695949392919061106b565b60408051601f19818403018152919052805160209091012090505f61046a8261043160e08701876110c3565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061085092505050565b600480546040517f3c23afba0000000000000000000000000000000000000000000000000000000081529293506001600160a01b031691633c23afba916104b391859101610e80565b602060405180830381865afa1580156104ce573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104f2919061112e565b61050e5760405162461bcd60e51b81526004016103379061117d565b600480546040517f6d46e9870000000000000000000000000000000000000000000000000000000081526001600160a01b0390911691636d46e9879161055691859101610e80565b602060405180830381865afa158015610571573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610595919061112e565b6105b15760405162461bcd60e51b8152600401610337906111bf565b6105ba8561087a565b60a08501355f1914610643576003546040517fb6aed0cb0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063b6aed0cb906106159060a08901359042906004016111cf565b5f604051808303815f87803b15801561062c575f5ffd5b505af115801561063e573d5f5f3e3d5ffd5b505050505b7fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b12455f4961067360e08801886110c3565b60405161068293929190611215565b60405180910390a15050505050565b5f6106d36040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081525090565b5f5f5f015f8581526020019081526020015f20604051806101000160405290815f82015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820180546107479061124a565b80601f01602080910402602001604051908101604052809291908181526020018280546107739061124a565b80156107be5780601f10610795576101008083540402835291602001916107be565b820191905f5260205f20905b8154815290600101906020018083116107a157829003601f168201915b50505091909252505081519095149590945092505050565b6107de6108ae565b6107e75f610922565b565b6107f16108ae565b6001600160a01b038116610833575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016103379190610e80565b61083c81610922565b50565b61084761099f565b61083c81610a06565b5f5f5f5f61085e8686610a0e565b92509250925061086e8282610a57565b50909150505b92915050565b80355f90815260208190526040902081906108958282611518565b50506002546040820135111561083c5760400135600255565b336108e07f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146107e757336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016103379190610e80565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166107e7576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107f161099f565b5f5f5f8351604103610a45576020840151604085015160608601515f1a610a3788828585610b5c565b955095509550505050610a50565b505081515f91506002905b9250925092565b5f826003811115610a6a57610a6a611522565b03610a73575050565b6001826003811115610a8757610a87611522565b03610abe576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002826003811115610ad257610ad2611522565b03610b0b576040517ffce698f7000000000000000000000000000000000000000000000000000000008152610337908290600401610c1e565b6003826003811115610b1f57610b1f611522565b03610b5857806040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004016103379190610c1e565b5050565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610b9557505f91506003905082610c0c565b5f6001888888886040515f8152602001604052604051610bb8949392919061153f565b6020604051602081039080840390855afa158015610bd8573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116610c0357505f925060019150829050610c0c565b92505f91508190505b9450945094915050565b805b82525050565b602081016108748284610c16565b5f6001600160a01b038216610874565b610c4581610c2c565b811461083c575f5ffd5b803561087481610c3c565b5f5f60408385031215610c6e57610c6e5f5ffd5b610c788484610c4f565b9150610c878460208501610c4f565b90509250929050565b5f6101008284031215610ca457610ca45f5ffd5b50919050565b5f60208284031215610cbd57610cbd5f5ffd5b813567ffffffffffffffff811115610cd657610cd65f5ffd5b610ce284828501610c90565b949350505050565b80610c45565b803561087481610cea565b5f60208284031215610d0e57610d0e5f5ffd5b610d188383610cf0565b9392505050565b801515610c18565b8281835e505f910152565b5f610d3b825190565b808452602084019350610d52818560208601610d27565b601f01601f19169290920192915050565b80515f90610100840190610d778582610c16565b506020830151610d8a6020860182610c16565b506040830151610d9d6040860182610c16565b506060830151610db06060860182610c16565b506080830151610dc36080860182610c16565b5060a0830151610dd660a0860182610c16565b5060c0830151610de960c0860182610c16565b5060e083015184820360e0860152610e018282610d32565b95945050505050565b60408101610e188285610d1f565b8181036020830152610ce28184610d63565b5f6108746001600160a01b038316610e40565b90565b6001600160a01b031690565b5f61087482610e2a565b5f61087482610e4c565b610c1881610e56565b602081016108748284610e60565b610c1881610c2c565b602081016108748284610e77565b5f60208284031215610ea157610ea15f5ffd5b610d188383610c4f565b5f67ffffffffffffffff8216610874565b610c1881610eab565b602081016108748284610ebc565b6020808252810161087481602681527f43616e6e6f742062696e6420746f20667574757265206f722063757272656e7460208201527f20626c6f636b0000000000000000000000000000000000000000000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561087457610874610f34565b60158152602081017f426c6f636b2062696e64696e6720746f6f206f6c640000000000000000000000815290505b60200190565b6020808252810161087481610f5b565b60128152602081017f556e6b6e6f776e20626c6f636b2068617368000000000000000000000000000081529050610f89565b6020808252810161087481610f9f565b60168152602081017f426c6f636b2062696e64696e67206d69736d617463680000000000000000000081529050610f89565b6020808252810161087481610fe1565b60148152602081017f426c6f622068617368206973206e6f742073657400000000000000000000000081529050610f89565b6020808252810161087481611023565b80610c18565b6110758188611065565b6020016110828187611065565b60200161108f8186611065565b60200161109c8185611065565b6020016110a98184611065565b6020016110b68183611065565b6020019695505050505050565b5f808335601e19368590030181126110dc576110dc5f5ffd5b8301915050803567ffffffffffffffff8111156110fa576110fa5f5ffd5b602082019150600181023603821315611114576111145f5ffd5b9250929050565b801515610c45565b80516108748161111b565b5f60208284031215611141576111415f5ffd5b610d188383611123565b60168152602081017f656e636c6176654944206e6f742061747465737465640000000000000000000081529050610f89565b602080825281016108748161114b565b60198152602081017f656e636c6176654944206e6f7420612073657175656e6365720000000000000081529050610f89565b602080825281016108748161118d565b604081016111dd8285610c16565b610d186020830184610c16565b82818337505f910152565b8183526020830192506112098284836111ea565b50601f01601f19160190565b604081016112238286610c16565b8181036020830152610e018184866111f5565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061125e57607f821691505b602082108103610ca457610ca4611236565b5f813561087481610cea565b5f81610874565b61128c8261127c565b611298610e3d8261127c565b8255505050565b5f610874610e3d8381565b6112b38261129f565b80611298565b634e487b7160e01b5f52604160045260245ffd5b6112d68361129f565b81545f1960089490940293841b1916921b91909117905550565b5f6112fc8184846112cd565b505050565b81811015610b58576113135f826112f0565b600101611301565b601f8211156112fc575f818152602090206020601f850104810160208510156113415750805b6113536020601f860104830182611301565b5050505050565b8267ffffffffffffffff811115611373576113736112b9565b61137d825461124a565b61138882828561131b565b505f601f8211600181146113ba575f83156113a35750848201355b5f19600885021c1981166002850217855550611411565b5f84815260208120601f198516915b828110156113e957878501358255602094850194600190920191016113c9565b5084821015611405575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b6112fc83838361135a565b818061142f81611270565b905061143b8184611283565b5050602082018061144b82611270565b905061145a81600185016112aa565b5050604082018061146a82611270565b905061147981600285016112aa565b5050606082018061148982611270565b90506114988160038501611283565b505060808201806114a882611270565b90506114b781600485016112aa565b505060a08201806114c782611270565b90506114d68160058501611283565b505060c08201806114e682611270565b90506114f58160068501611283565b505061150460e08301836110c3565b611512818360078601611419565b50505050565b610b588282611424565b634e487b7160e01b5f52602160045260245ffd5b60ff8116610c18565b6080810161154d8287610c16565b61155a6020830186611536565b6115676040830185610c16565b610e016060830184610c1656fea26469706673582212205a55b243b561de22558554de851d375f3b28395edfabe6b784433c0fc9e60cbb64736f6c634300081c0033",
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
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
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
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_RollupContract *RollupContractSession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _RollupContract.Contract.GetRollupByHash(&_RollupContract.CallOpts, rollupHash)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_RollupContract *RollupContractCallerSession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _RollupContract.Contract.GetRollupByHash(&_RollupContract.CallOpts, rollupHash)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_RollupContract *RollupContractCaller) LastBatchSeqNo(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RollupContract.contract.Call(opts, &out, "lastBatchSeqNo")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_RollupContract *RollupContractSession) LastBatchSeqNo() (*big.Int, error) {
	return _RollupContract.Contract.LastBatchSeqNo(&_RollupContract.CallOpts)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_RollupContract *RollupContractCallerSession) LastBatchSeqNo() (*big.Int, error) {
	return _RollupContract.Contract.LastBatchSeqNo(&_RollupContract.CallOpts)
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

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_RollupContract *RollupContractTransactor) AddRollup(opts *bind.TransactOpts, r StructsMetaRollup) (*types.Transaction, error) {
	return _RollupContract.contract.Transact(opts, "addRollup", r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_RollupContract *RollupContractSession) AddRollup(r StructsMetaRollup) (*types.Transaction, error) {
	return _RollupContract.Contract.AddRollup(&_RollupContract.TransactOpts, r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
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
