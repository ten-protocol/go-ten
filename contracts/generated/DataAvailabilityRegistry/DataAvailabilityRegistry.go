// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DataAvailabilityRegistry

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

// DataAvailabilityRegistryMetaData contains all meta data concerning the DataAvailabilityRegistry contract.
var DataAvailabilityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"RollupAdded\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"addRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enclaveRegistry\",\"outputs\":[{\"internalType\":\"contractINetworkEnclaveRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChallengePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"getRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_merkleMessageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_enclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSeqNo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_delay\",\"type\":\"uint256\"}],\"name\":\"setChallengePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b61161b806100975f395ff3fe608060405234801561000f575f5ffd5b50600436106100c4575f3560e01c80637864b77d1161007d578063c0c53b8b11610058578063c0c53b8b14610198578063e874eb20146101ab578063f2fde38b146101be575f5ffd5b80637864b77d146101385780637c72dbd0146101405780638da5cb5b14610160575f5ffd5b80635fdf31a2116100ad5780635fdf31a2146100fc5780636fb6a45c1461010f578063715018a614610130575f5ffd5b8063440c953b146100c85780635d475fdd146100e7575b5f5ffd5b6100d160025481565b6040516100de9190610c83565b60405180910390f35b6100fa6100f5366004610ca8565b6101d1565b005b6100fa61010a366004610ce6565b6101de565b61012261011d366004610ca8565b61056e565b6040516100de929190610e11565b6100fa6106b3565b6003546100d1565b600554610153906001600160a01b031681565b6040516100de9190610e70565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516100de9190610e97565b6100fa6101a6366004610eb9565b6106c6565b600454610153906001600160a01b031681565b6100fa6101cc366004610eff565b61084e565b6101d96108a4565b600355565b808060800135431161020b5760405162461bcd60e51b815260040161020290610f1c565b60405180910390fd5b61021a608082013560ff610f91565b43106102385760405162461bcd60e51b815260040161020290610fd8565b6080810135405f81900361025e5760405162461bcd60e51b81526004016102029061101a565b816060013581146102815760405162461bcd60e51b81526004016102029061105c565b5f4961029f5760405162461bcd60e51b81526004016102029061109e565b5f82604001358360c00135846060013585608001358660a001355f496040516020016102d0969594939291906110b4565b60408051601f19818403018152919052805160209091012090505f610335826102fc60e087018761110c565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061091892505050565b6005546040517f3c23afba0000000000000000000000000000000000000000000000000000000081529192506001600160a01b031690633c23afba9061037f908490600401610e97565b602060405180830381865afa15801561039a573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103be9190611177565b6103da5760405162461bcd60e51b8152600401610202906111c6565b6005546040517f6d46e9870000000000000000000000000000000000000000000000000000000081526001600160a01b0390911690636d46e98790610423908490600401610e97565b602060405180830381865afa15801561043e573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104629190611177565b61047e5760405162461bcd60e51b815260040161020290611208565b61048785610942565b60a08501355f1914610520575f600354426104a29190610f91565b600480546040517fb6aed0cb0000000000000000000000000000000000000000000000000000000081529293506001600160a01b03169163b6aed0cb916104f19160a08b013591869101611218565b5f604051808303815f87803b158015610508575f5ffd5b505af115801561051a573d5f5f3e3d5ffd5b50505050505b7fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b12455f4961055060e088018861110c565b60405161055f9392919061125e565b60405180910390a15050505050565b5f6105b06040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081525090565b5f5f5f015f8581526020019081526020015f20604051806101000160405290815f820154815260200160018201548152602001600282015481526020016003820154815260200160048201548152602001600582015481526020016006820154815260200160078201805461062490611293565b80601f016020809104026020016040519081016040528092919081815260200182805461065090611293565b801561069b5780601f106106725761010080835404028352916020019161069b565b820191905f5260205f20905b81548152906001019060200180831161067e57829003601f168201915b50505091909252505081519095149590945092505050565b6106bb6108a4565b6106c45f610976565b565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156107105750825b90505f8267ffffffffffffffff16600114801561072c5750303b155b90508115801561073a575080155b15610771576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156107a557845468ff00000000000000001916680100000000000000001785555b6107ae866109f3565b600480546001600160a01b03808b1673ffffffffffffffffffffffffffffffffffffffff199283161790925560058054928a16929091169190911790555f6002819055600355831561084457845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061083b906001906112d3565b60405180910390a15b5050505050505050565b6108566108a4565b6001600160a01b038116610898575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016102029190610e97565b6108a181610976565b50565b336108d67f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146106c457336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016102029190610e97565b5f5f5f5f6109268686610a04565b9250925092506109368282610a4d565b50909150505b92915050565b80355f908152602081905260409020819061095d8282611589565b5050600254604082013511156108a15760400135600255565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b6109fb610b52565b6108a181610bb9565b5f5f5f8351604103610a3b576020840151604085015160608601515f1a610a2d88828585610bc1565b955095509550505050610a46565b505081515f91506002905b9250925092565b5f826003811115610a6057610a60611593565b03610a69575050565b6001826003811115610a7d57610a7d611593565b03610ab4576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002826003811115610ac857610ac8611593565b03610b01576040517ffce698f7000000000000000000000000000000000000000000000000000000008152610202908290600401610c83565b6003826003811115610b1557610b15611593565b03610b4e57806040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004016102029190610c83565b5050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166106c4576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610856610b52565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610bfa57505f91506003905082610c71565b5f6001888888886040515f8152602001604052604051610c1d94939291906115b0565b6020604051602081039080840390855afa158015610c3d573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116610c6857505f925060019150829050610c71565b92505f91508190505b9450945094915050565b805b82525050565b6020810161093c8284610c7b565b805b81146108a1575f5ffd5b803561093c81610c91565b5f60208284031215610cbb57610cbb5f5ffd5b610cc58383610c9d565b9392505050565b5f6101008284031215610ce057610ce05f5ffd5b50919050565b5f60208284031215610cf957610cf95f5ffd5b813567ffffffffffffffff811115610d1257610d125f5ffd5b610d1e84828501610ccc565b949350505050565b801515610c7d565b8281835e505f910152565b5f610d42825190565b808452602084019350610d59818560208601610d2e565b601f01601f19169290920192915050565b80515f90610100840190610d7e8582610c7b565b506020830151610d916020860182610c7b565b506040830151610da46040860182610c7b565b506060830151610db76060860182610c7b565b506080830151610dca6080860182610c7b565b5060a0830151610ddd60a0860182610c7b565b5060c0830151610df060c0860182610c7b565b5060e083015184820360e0860152610e088282610d39565b95945050505050565b60408101610e1f8285610d26565b8181036020830152610d1e8184610d6a565b5f61093c6001600160a01b038316610e47565b90565b6001600160a01b031690565b5f61093c82610e31565b5f61093c82610e53565b610c7d81610e5d565b6020810161093c8284610e67565b5f6001600160a01b03821661093c565b610c7d81610e7e565b6020810161093c8284610e8e565b610c9381610e7e565b803561093c81610ea5565b5f5f5f60608486031215610ece57610ece5f5ffd5b610ed88585610eae565b9250610ee78560208601610eae565b9150610ef68560408601610eae565b90509250925092565b5f60208284031215610f1257610f125f5ffd5b610cc58383610eae565b6020808252810161093c81602681527f43616e6e6f742062696e6420746f20667574757265206f722063757272656e7460208201527f20626c6f636b0000000000000000000000000000000000000000000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561093c5761093c610f7d565b60158152602081017f426c6f636b2062696e64696e6720746f6f206f6c640000000000000000000000815290505b60200190565b6020808252810161093c81610fa4565b60128152602081017f556e6b6e6f776e20626c6f636b2068617368000000000000000000000000000081529050610fd2565b6020808252810161093c81610fe8565b60168152602081017f426c6f636b2062696e64696e67206d69736d617463680000000000000000000081529050610fd2565b6020808252810161093c8161102a565b60148152602081017f426c6f622068617368206973206e6f742073657400000000000000000000000081529050610fd2565b6020808252810161093c8161106c565b80610c7d565b6110be81886110ae565b6020016110cb81876110ae565b6020016110d881866110ae565b6020016110e581856110ae565b6020016110f281846110ae565b6020016110ff81836110ae565b6020019695505050505050565b5f808335601e1936859003018112611125576111255f5ffd5b8301915050803567ffffffffffffffff811115611143576111435f5ffd5b60208201915060018102360382131561115d5761115d5f5ffd5b9250929050565b801515610c93565b805161093c81611164565b5f6020828403121561118a5761118a5f5ffd5b610cc5838361116c565b60168152602081017f656e636c6176654944206e6f742061747465737465640000000000000000000081529050610fd2565b6020808252810161093c81611194565b60198152602081017f656e636c6176654944206e6f7420612073657175656e6365720000000000000081529050610fd2565b6020808252810161093c816111d6565b604081016112268285610c7b565b610cc56020830184610c7b565b82818337505f910152565b818352602083019250611252828483611233565b50601f01601f19160190565b6040810161126c8286610c7b565b8181036020830152610e0881848661123e565b634e487b7160e01b5f52602260045260245ffd5b6002810460018216806112a757607f821691505b602082108103610ce057610ce061127f565b5f67ffffffffffffffff821661093c565b610c7d816112b9565b6020810161093c82846112ca565b5f813561093c81610c91565b5f8161093c565b6112fd826112ed565b611309610e44826112ed565b8255505050565b5f61093c610e448381565b61132482611310565b80611309565b634e487b7160e01b5f52604160045260245ffd5b61134783611310565b81545f1960089490940293841b1916921b91909117905550565b5f61136d81848461133e565b505050565b81811015610b4e576113845f82611361565b600101611372565b601f82111561136d575f818152602090206020601f850104810160208510156113b25750805b6113c46020601f860104830182611372565b5050505050565b8267ffffffffffffffff8111156113e4576113e461132a565b6113ee8254611293565b6113f982828561138c565b505f601f82116001811461142b575f83156114145750848201355b5f19600885021c1981166002850217855550611482565b5f84815260208120601f198516915b8281101561145a578785013582556020948501946001909201910161143a565b5084821015611476575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b61136d8383836113cb565b81806114a0816112e1565b90506114ac81846112f4565b505060208201806114bc826112e1565b90506114cb816001850161131b565b505060408201806114db826112e1565b90506114ea816002850161131b565b505060608201806114fa826112e1565b905061150981600385016112f4565b50506080820180611519826112e1565b9050611528816004850161131b565b505060a0820180611538826112e1565b905061154781600585016112f4565b505060c0820180611557826112e1565b905061156681600685016112f4565b505061157560e083018361110c565b61158381836007860161148a565b50505050565b610b4e8282611495565b634e487b7160e01b5f52602160045260245ffd5b60ff8116610c7d565b608081016115be8287610c7b565b6115cb60208301866115a7565b6115d86040830185610c7b565b610e086060830184610c7b56fea26469706673582212204327d6ae43f584c2eb60b259ca1e8c1f20b2fa72fd1c765bace509040e81aa9664736f6c634300081c0033",
}

// DataAvailabilityRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use DataAvailabilityRegistryMetaData.ABI instead.
var DataAvailabilityRegistryABI = DataAvailabilityRegistryMetaData.ABI

// DataAvailabilityRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DataAvailabilityRegistryMetaData.Bin instead.
var DataAvailabilityRegistryBin = DataAvailabilityRegistryMetaData.Bin

// DeployDataAvailabilityRegistry deploys a new Ethereum contract, binding an instance of DataAvailabilityRegistry to it.
func DeployDataAvailabilityRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DataAvailabilityRegistry, error) {
	parsed, err := DataAvailabilityRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DataAvailabilityRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DataAvailabilityRegistry{DataAvailabilityRegistryCaller: DataAvailabilityRegistryCaller{contract: contract}, DataAvailabilityRegistryTransactor: DataAvailabilityRegistryTransactor{contract: contract}, DataAvailabilityRegistryFilterer: DataAvailabilityRegistryFilterer{contract: contract}}, nil
}

// DataAvailabilityRegistry is an auto generated Go binding around an Ethereum contract.
type DataAvailabilityRegistry struct {
	DataAvailabilityRegistryCaller     // Read-only binding to the contract
	DataAvailabilityRegistryTransactor // Write-only binding to the contract
	DataAvailabilityRegistryFilterer   // Log filterer for contract events
}

// DataAvailabilityRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataAvailabilityRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataAvailabilityRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataAvailabilityRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataAvailabilityRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataAvailabilityRegistrySession struct {
	Contract     *DataAvailabilityRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DataAvailabilityRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataAvailabilityRegistryCallerSession struct {
	Contract *DataAvailabilityRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// DataAvailabilityRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataAvailabilityRegistryTransactorSession struct {
	Contract     *DataAvailabilityRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// DataAvailabilityRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataAvailabilityRegistryRaw struct {
	Contract *DataAvailabilityRegistry // Generic contract binding to access the raw methods on
}

// DataAvailabilityRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryCallerRaw struct {
	Contract *DataAvailabilityRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// DataAvailabilityRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryTransactorRaw struct {
	Contract *DataAvailabilityRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataAvailabilityRegistry creates a new instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistry(address common.Address, backend bind.ContractBackend) (*DataAvailabilityRegistry, error) {
	contract, err := bindDataAvailabilityRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistry{DataAvailabilityRegistryCaller: DataAvailabilityRegistryCaller{contract: contract}, DataAvailabilityRegistryTransactor: DataAvailabilityRegistryTransactor{contract: contract}, DataAvailabilityRegistryFilterer: DataAvailabilityRegistryFilterer{contract: contract}}, nil
}

// NewDataAvailabilityRegistryCaller creates a new read-only instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistryCaller(address common.Address, caller bind.ContractCaller) (*DataAvailabilityRegistryCaller, error) {
	contract, err := bindDataAvailabilityRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryCaller{contract: contract}, nil
}

// NewDataAvailabilityRegistryTransactor creates a new write-only instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*DataAvailabilityRegistryTransactor, error) {
	contract, err := bindDataAvailabilityRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryTransactor{contract: contract}, nil
}

// NewDataAvailabilityRegistryFilterer creates a new log filterer instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*DataAvailabilityRegistryFilterer, error) {
	contract, err := bindDataAvailabilityRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryFilterer{contract: contract}, nil
}

// bindDataAvailabilityRegistry binds a generic wrapper to an already deployed contract.
func bindDataAvailabilityRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DataAvailabilityRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataAvailabilityRegistry.Contract.DataAvailabilityRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.DataAvailabilityRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.DataAvailabilityRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataAvailabilityRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.contract.Transact(opts, method, params...)
}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) EnclaveRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "enclaveRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) EnclaveRegistry() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.EnclaveRegistry(&_DataAvailabilityRegistry.CallOpts)
}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) EnclaveRegistry() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.EnclaveRegistry(&_DataAvailabilityRegistry.CallOpts)
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) GetChallengePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "getChallengePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) GetChallengePeriod() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.GetChallengePeriod(&_DataAvailabilityRegistry.CallOpts)
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) GetChallengePeriod() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.GetChallengePeriod(&_DataAvailabilityRegistry.CallOpts)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) GetRollupByHash(opts *bind.CallOpts, rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "getRollupByHash", rollupHash)

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
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _DataAvailabilityRegistry.Contract.GetRollupByHash(&_DataAvailabilityRegistry.CallOpts, rollupHash)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _DataAvailabilityRegistry.Contract.GetRollupByHash(&_DataAvailabilityRegistry.CallOpts, rollupHash)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) LastBatchSeqNo(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "lastBatchSeqNo")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) LastBatchSeqNo() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.LastBatchSeqNo(&_DataAvailabilityRegistry.CallOpts)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) LastBatchSeqNo() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.LastBatchSeqNo(&_DataAvailabilityRegistry.CallOpts)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) MerkleMessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "merkleMessageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) MerkleMessageBus() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.MerkleMessageBus(&_DataAvailabilityRegistry.CallOpts)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) MerkleMessageBus() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.MerkleMessageBus(&_DataAvailabilityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) Owner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.Owner(&_DataAvailabilityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) Owner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.Owner(&_DataAvailabilityRegistry.CallOpts)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) AddRollup(opts *bind.TransactOpts, r StructsMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "addRollup", r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) AddRollup(r StructsMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AddRollup(&_DataAvailabilityRegistry.TransactOpts, r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) AddRollup(r StructsMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AddRollup(&_DataAvailabilityRegistry.TransactOpts, r)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry, address _owner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) Initialize(opts *bind.TransactOpts, _merkleMessageBus common.Address, _enclaveRegistry common.Address, _owner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "initialize", _merkleMessageBus, _enclaveRegistry, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry, address _owner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) Initialize(_merkleMessageBus common.Address, _enclaveRegistry common.Address, _owner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.Initialize(&_DataAvailabilityRegistry.TransactOpts, _merkleMessageBus, _enclaveRegistry, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry, address _owner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) Initialize(_merkleMessageBus common.Address, _enclaveRegistry common.Address, _owner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.Initialize(&_DataAvailabilityRegistry.TransactOpts, _merkleMessageBus, _enclaveRegistry, _owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.RenounceOwnership(&_DataAvailabilityRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.RenounceOwnership(&_DataAvailabilityRegistry.TransactOpts)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) SetChallengePeriod(opts *bind.TransactOpts, _delay *big.Int) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "setChallengePeriod", _delay)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) SetChallengePeriod(_delay *big.Int) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.SetChallengePeriod(&_DataAvailabilityRegistry.TransactOpts, _delay)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) SetChallengePeriod(_delay *big.Int) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.SetChallengePeriod(&_DataAvailabilityRegistry.TransactOpts, _delay)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.TransferOwnership(&_DataAvailabilityRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.TransferOwnership(&_DataAvailabilityRegistry.TransactOpts, newOwner)
}

// DataAvailabilityRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryInitializedIterator struct {
	Event *DataAvailabilityRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryInitialized)
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
		it.Event = new(DataAvailabilityRegistryInitialized)
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
func (it *DataAvailabilityRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryInitialized represents a Initialized event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*DataAvailabilityRegistryInitializedIterator, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryInitializedIterator{contract: _DataAvailabilityRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryInitialized)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseInitialized(log types.Log) (*DataAvailabilityRegistryInitialized, error) {
	event := new(DataAvailabilityRegistryInitialized)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataAvailabilityRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferredIterator struct {
	Event *DataAvailabilityRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryOwnershipTransferred)
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
		it.Event = new(DataAvailabilityRegistryOwnershipTransferred)
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
func (it *DataAvailabilityRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DataAvailabilityRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryOwnershipTransferredIterator{contract: _DataAvailabilityRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryOwnershipTransferred)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*DataAvailabilityRegistryOwnershipTransferred, error) {
	event := new(DataAvailabilityRegistryOwnershipTransferred)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataAvailabilityRegistryRollupAddedIterator is returned from FilterRollupAdded and is used to iterate over the raw logs and unpacked data for RollupAdded events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryRollupAddedIterator struct {
	Event *DataAvailabilityRegistryRollupAdded // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryRollupAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryRollupAdded)
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
		it.Event = new(DataAvailabilityRegistryRollupAdded)
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
func (it *DataAvailabilityRegistryRollupAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryRollupAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryRollupAdded represents a RollupAdded event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryRollupAdded struct {
	RollupHash [32]byte
	Signature  []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRollupAdded is a free log retrieval operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterRollupAdded(opts *bind.FilterOpts) (*DataAvailabilityRegistryRollupAddedIterator, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "RollupAdded")
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryRollupAddedIterator{contract: _DataAvailabilityRegistry.contract, event: "RollupAdded", logs: logs, sub: sub}, nil
}

// WatchRollupAdded is a free log subscription operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchRollupAdded(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryRollupAdded) (event.Subscription, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "RollupAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryRollupAdded)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "RollupAdded", log); err != nil {
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
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseRollupAdded(log types.Log) (*DataAvailabilityRegistryRollupAdded, error) {
	event := new(DataAvailabilityRegistryRollupAdded)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "RollupAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
