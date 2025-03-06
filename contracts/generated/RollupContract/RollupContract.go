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
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b6115b2806100975f395ff3fe608060405234801561000f575f5ffd5b506004361061009f575f3560e01c8063715018a6116100725780638da5cb5b116100585780638da5cb5b14610133578063e874eb201461016b578063f2fde38b1461017e575f5ffd5b8063715018a61461010b5780637c72dbd014610113575f5ffd5b8063440c953b146100a3578063485cc955146100c25780635fdf31a2146100d75780636fb6a45c146100ea575b5f5ffd5b6100ac60025481565b6040516100b99190610c26565b60405180910390f35b6100d56100d0366004610c62565b610191565b005b6100d56100e5366004610cb2565b610313565b6100fd6100f8366004610d03565b610699565b6040516100b9929190610e12565b6100d56107de565b600454610126906001600160a01b031681565b6040516100b99190610e71565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516100b99190610e88565b600354610126906001600160a01b031681565b6100d561018c366004610e96565b6107f1565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156101db5750825b90505f8267ffffffffffffffff1660011480156101f75750303b155b905081158015610205575080155b1561023c576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561027057845468ff00000000000000001916680100000000000000001785555b61027933610847565b600380546001600160a01b03808a1673ffffffffffffffffffffffffffffffffffffffff199283161790925560048054928916929091169190911790555f600255831561030a57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061030190600190610ecd565b60405180910390a15b50505050505050565b80806080013543116103405760405162461bcd60e51b815260040161033790610edb565b60405180910390fd5b61034f608082013560ff610f50565b431061036d5760405162461bcd60e51b815260040161033790610f97565b6080810135405f8190036103935760405162461bcd60e51b815260040161033790610fd9565b816060013581146103b65760405162461bcd60e51b81526004016103379061101b565b5f496103d45760405162461bcd60e51b81526004016103379061105d565b5f82604001358360c00135846060013585608001358660a001355f4960405160200161040596959493929190611073565b60408051601f19818403018152919052805160209091012090505f61046a8261043160e08701876110cb565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061085892505050565b600480546040517f3c23afba0000000000000000000000000000000000000000000000000000000081529293506001600160a01b031691633c23afba916104b391859101610e88565b602060405180830381865afa1580156104ce573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104f29190611136565b61050e5760405162461bcd60e51b815260040161033790611185565b600480546040517f6d46e9870000000000000000000000000000000000000000000000000000000081526001600160a01b0390911691636d46e9879161055691859101610e88565b602060405180830381865afa158015610571573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105959190611136565b6105b15760405162461bcd60e51b8152600401610337906111c7565b6105b9610882565b6105c2856108f6565b60a08501355f191461064b576003546040517fb6aed0cb0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063b6aed0cb9061061d9060a08901359042906004016111d7565b5f604051808303815f87803b158015610634575f5ffd5b505af1158015610646573d5f5f3e3d5ffd5b505050505b7fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b12455f4961067b60e08801886110cb565b60405161068a9392919061121d565b60405180910390a15050505050565b5f6106db6040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081525090565b5f5f5f015f8581526020019081526020015f20604051806101000160405290815f820154815260200160018201548152602001600282015481526020016003820154815260200160048201548152602001600582015481526020016006820154815260200160078201805461074f90611252565b80601f016020809104026020016040519081016040528092919081815260200182805461077b90611252565b80156107c65780601f1061079d576101008083540402835291602001916107c6565b820191905f5260205f20905b8154815290600101906020018083116107a957829003601f168201915b50505091909252505081519095149590945092505050565b6107e6610882565b6107ef5f61092a565b565b6107f9610882565b6001600160a01b03811661083b575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016103379190610e88565b6108448161092a565b50565b61084f6109a7565b61084481610a0e565b5f5f5f5f6108668686610a16565b9250925092506108768282610a5f565b50909150505b92915050565b336108b47f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146107ef57336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016103379190610e88565b80355f90815260208190526040902081906109118282611520565b5050600254604082013511156108445760400135600255565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166107ef576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107f96109a7565b5f5f5f8351604103610a4d576020840151604085015160608601515f1a610a3f88828585610b64565b955095509550505050610a58565b505081515f91506002905b9250925092565b5f826003811115610a7257610a7261152a565b03610a7b575050565b6001826003811115610a8f57610a8f61152a565b03610ac6576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002826003811115610ada57610ada61152a565b03610b13576040517ffce698f7000000000000000000000000000000000000000000000000000000008152610337908290600401610c26565b6003826003811115610b2757610b2761152a565b03610b6057806040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004016103379190610c26565b5050565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610b9d57505f91506003905082610c14565b5f6001888888886040515f8152602001604052604051610bc09493929190611547565b6020604051602081039080840390855afa158015610be0573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116610c0b57505f925060019150829050610c14565b92505f91508190505b9450945094915050565b805b82525050565b6020810161087c8284610c1e565b5f6001600160a01b03821661087c565b610c4d81610c34565b8114610844575f5ffd5b803561087c81610c44565b5f5f60408385031215610c7657610c765f5ffd5b610c808484610c57565b9150610c8f8460208501610c57565b90509250929050565b5f6101008284031215610cac57610cac5f5ffd5b50919050565b5f60208284031215610cc557610cc55f5ffd5b813567ffffffffffffffff811115610cde57610cde5f5ffd5b610cea84828501610c98565b949350505050565b80610c4d565b803561087c81610cf2565b5f60208284031215610d1657610d165f5ffd5b610d208383610cf8565b9392505050565b801515610c20565b8281835e505f910152565b5f610d43825190565b808452602084019350610d5a818560208601610d2f565b601f01601f19169290920192915050565b80515f90610100840190610d7f8582610c1e565b506020830151610d926020860182610c1e565b506040830151610da56040860182610c1e565b506060830151610db86060860182610c1e565b506080830151610dcb6080860182610c1e565b5060a0830151610dde60a0860182610c1e565b5060c0830151610df160c0860182610c1e565b5060e083015184820360e0860152610e098282610d3a565b95945050505050565b60408101610e208285610d27565b8181036020830152610cea8184610d6b565b5f61087c6001600160a01b038316610e48565b90565b6001600160a01b031690565b5f61087c82610e32565b5f61087c82610e54565b610c2081610e5e565b6020810161087c8284610e68565b610c2081610c34565b6020810161087c8284610e7f565b5f60208284031215610ea957610ea95f5ffd5b610d208383610c57565b5f67ffffffffffffffff821661087c565b610c2081610eb3565b6020810161087c8284610ec4565b6020808252810161087c81602681527f43616e6e6f742062696e6420746f20667574757265206f722063757272656e7460208201527f20626c6f636b0000000000000000000000000000000000000000000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561087c5761087c610f3c565b60158152602081017f426c6f636b2062696e64696e6720746f6f206f6c640000000000000000000000815290505b60200190565b6020808252810161087c81610f63565b60128152602081017f556e6b6e6f776e20626c6f636b2068617368000000000000000000000000000081529050610f91565b6020808252810161087c81610fa7565b60168152602081017f426c6f636b2062696e64696e67206d69736d617463680000000000000000000081529050610f91565b6020808252810161087c81610fe9565b60148152602081017f426c6f622068617368206973206e6f742073657400000000000000000000000081529050610f91565b6020808252810161087c8161102b565b80610c20565b61107d818861106d565b60200161108a818761106d565b602001611097818661106d565b6020016110a4818561106d565b6020016110b1818461106d565b6020016110be818361106d565b6020019695505050505050565b5f808335601e19368590030181126110e4576110e45f5ffd5b8301915050803567ffffffffffffffff811115611102576111025f5ffd5b60208201915060018102360382131561111c5761111c5f5ffd5b9250929050565b801515610c4d565b805161087c81611123565b5f60208284031215611149576111495f5ffd5b610d20838361112b565b60168152602081017f656e636c6176654944206e6f742061747465737465640000000000000000000081529050610f91565b6020808252810161087c81611153565b60198152602081017f656e636c6176654944206e6f7420612073657175656e6365720000000000000081529050610f91565b6020808252810161087c81611195565b604081016111e58285610c1e565b610d206020830184610c1e565b82818337505f910152565b8183526020830192506112118284836111f2565b50601f01601f19160190565b6040810161122b8286610c1e565b8181036020830152610e098184866111fd565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061126657607f821691505b602082108103610cac57610cac61123e565b5f813561087c81610cf2565b5f8161087c565b61129482611284565b6112a0610e4582611284565b8255505050565b5f61087c610e458381565b6112bb826112a7565b806112a0565b634e487b7160e01b5f52604160045260245ffd5b6112de836112a7565b81545f1960089490940293841b1916921b91909117905550565b5f6113048184846112d5565b505050565b81811015610b605761131b5f826112f8565b600101611309565b601f821115611304575f818152602090206020601f850104810160208510156113495750805b61135b6020601f860104830182611309565b5050505050565b8267ffffffffffffffff81111561137b5761137b6112c1565b6113858254611252565b611390828285611323565b505f601f8211600181146113c2575f83156113ab5750848201355b5f19600885021c1981166002850217855550611419565b5f84815260208120601f198516915b828110156113f157878501358255602094850194600190920191016113d1565b508482101561140d575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b611304838383611362565b818061143781611278565b9050611443818461128b565b5050602082018061145382611278565b905061146281600185016112b2565b5050604082018061147282611278565b905061148181600285016112b2565b5050606082018061149182611278565b90506114a0816003850161128b565b505060808201806114b082611278565b90506114bf81600485016112b2565b505060a08201806114cf82611278565b90506114de816005850161128b565b505060c08201806114ee82611278565b90506114fd816006850161128b565b505061150c60e08301836110cb565b61151a818360078601611421565b50505050565b610b60828261142c565b634e487b7160e01b5f52602160045260245ffd5b60ff8116610c20565b608081016115558287610c1e565b611562602083018661153e565b61156f6040830185610c1e565b610e096060830184610c1e56fea2646970667358221220be64ff1afa425091ac9a03ad10d929b5425fdc670e93682c18323e873d34181264736f6c634300081c0033",
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
