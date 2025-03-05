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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"RollupAdded\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"addRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enclaveRegistry\",\"outputs\":[{\"internalType\":\"contractINetworkEnclaveRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"getRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_merkleMessageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_enclaveRegistry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSeqNo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b61156c806100975f395ff3fe608060405234801561000f575f5ffd5b506004361061009f575f3560e01c80637c72dbd011610072578063ae247b4911610058578063ae247b4914610158578063e874eb201461016b578063f2fde38b1461017e575f5ffd5b80637c72dbd0146101005780638da5cb5b14610120575f5ffd5b8063440c953b146100a3578063485cc955146100c25780636fb6a45c146100d7578063715018a6146100f8575b5f5ffd5b6100ac60025481565b6040516100b99190610c14565b60405180910390f35b6100d56100d0366004610c50565b610191565b005b6100ea6100e5366004610c97565b610313565b6040516100b9929190610d92565b6100d5610446565b600454610113906001600160a01b031681565b6040516100b99190610df9565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516100b99190610e10565b6100d5610166366004610e37565b610459565b600354610113906001600160a01b031681565b6100d561018c366004610e6f565b6107df565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156101db5750825b90505f8267ffffffffffffffff1660011480156101f75750303b155b905081158015610205575080155b1561023c576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561027057845468ff00000000000000001916680100000000000000001785555b61027933610835565b600380546001600160a01b03808a1673ffffffffffffffffffffffffffffffffffffffff199283161790925560048054928916929091169190911790555f600255831561030a57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061030190600190610ea6565b60405180910390a15b50505050505050565b5f61034e6040518060e001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081525090565b5f5f5f015f8581526020019081526020015f206040518060e00160405290815f820154815260200160018201548152602001600282015481526020016003820154815260200160048201548152602001600582015481526020016006820180546103b790610ec8565b80601f01602080910402602001604051908101604052809291908181526020018280546103e390610ec8565b801561042e5780601f106104055761010080835404028352916020019161042e565b820191905f5260205f20905b81548152906001019060200180831161041157829003601f168201915b50505091909252505081519095149590945092505050565b61044e610846565b6104575f6108ba565b565b80806060013543116104865760405162461bcd60e51b815260040161047d90610eee565b60405180910390fd5b610495606082013560ff610f63565b43106104b35760405162461bcd60e51b815260040161047d90610faa565b6060810135405f8190036104d95760405162461bcd60e51b815260040161047d90610fec565b816040013581146104fc5760405162461bcd60e51b815260040161047d9061102e565b5f4961051a5760405162461bcd60e51b815260040161047d90611070565b5f82602001358360a001358460400135856060013586608001355f4960405160200161054b96959493929190611086565b60408051601f19818403018152919052805160209091012090505f6105b08261057760c08701876110de565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061093792505050565b600480546040517f3c23afba0000000000000000000000000000000000000000000000000000000081529293506001600160a01b031691633c23afba916105f991859101610e10565b602060405180830381865afa158015610614573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106389190611149565b6106545760405162461bcd60e51b815260040161047d90611198565b600480546040517f6d46e9870000000000000000000000000000000000000000000000000000000081526001600160a01b0390911691636d46e9879161069c91859101610e10565b602060405180830381865afa1580156106b7573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106db9190611149565b6106f75760405162461bcd60e51b815260040161047d906111da565b6106ff610846565b61070885610961565b60808501355f1914610791576003546040517fb6aed0cb0000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063b6aed0cb906107639060808901359042906004016111ea565b5f604051808303815f87803b15801561077a575f5ffd5b505af115801561078c573d5f5f3e3d5ffd5b505050505b7fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b12455f496107c160c08801886110de565b6040516107d093929190611230565b60405180910390a15050505050565b6107e7610846565b6001600160a01b038116610829575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161047d9190610e10565b610832816108ba565b50565b61083d610995565b610832816109fc565b336108787f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461045757336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161047d9190610e10565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f5f5f5f6109458686610a04565b9250925092506109558282610a4d565b50909150505b92915050565b80355f908152602081905260409020819061097c82826114da565b5050600254602082013511156108325760200135600255565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610457576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107e7610995565b5f5f5f8351604103610a3b576020840151604085015160608601515f1a610a2d88828585610b52565b955095509550505050610a46565b505081515f91506002905b9250925092565b5f826003811115610a6057610a606114e4565b03610a69575050565b6001826003811115610a7d57610a7d6114e4565b03610ab4576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002826003811115610ac857610ac86114e4565b03610b01576040517ffce698f700000000000000000000000000000000000000000000000000000000815261047d908290600401610c14565b6003826003811115610b1557610b156114e4565b03610b4e57806040517fd78bce0c00000000000000000000000000000000000000000000000000000000815260040161047d9190610c14565b5050565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610b8b57505f91506003905082610c02565b5f6001888888886040515f8152602001604052604051610bae9493929190611501565b6020604051602081039080840390855afa158015610bce573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116610bf957505f925060019150829050610c02565b92505f91508190505b9450945094915050565b805b82525050565b6020810161095b8284610c0c565b5f6001600160a01b03821661095b565b610c3b81610c22565b8114610832575f5ffd5b803561095b81610c32565b5f5f60408385031215610c6457610c645f5ffd5b610c6e8484610c45565b9150610c7d8460208501610c45565b90509250929050565b80610c3b565b803561095b81610c86565b5f60208284031215610caa57610caa5f5ffd5b610cb48383610c8c565b9392505050565b801515610c0e565b8281835e505f910152565b5f610cd7825190565b808452602084019350610cee818560208601610cc3565b601f01601f19169290920192915050565b80515f9060e0840190610d128582610c0c565b506020830151610d256020860182610c0c565b506040830151610d386040860182610c0c565b506060830151610d4b6060860182610c0c565b506080830151610d5e6080860182610c0c565b5060a0830151610d7160a0860182610c0c565b5060c083015184820360c0860152610d898282610cce565b95945050505050565b60408101610da08285610cbb565b8181036020830152610db28184610cff565b949350505050565b5f61095b6001600160a01b038316610dd0565b90565b6001600160a01b031690565b5f61095b82610dba565b5f61095b82610ddc565b610c0e81610de6565b6020810161095b8284610df0565b610c0e81610c22565b6020810161095b8284610e07565b5f60e08284031215610e3157610e315f5ffd5b50919050565b5f60208284031215610e4a57610e4a5f5ffd5b813567ffffffffffffffff811115610e6357610e635f5ffd5b610db284828501610e1e565b5f60208284031215610e8257610e825f5ffd5b610cb48383610c45565b5f67ffffffffffffffff821661095b565b610c0e81610e8c565b6020810161095b8284610e9d565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680610edc57607f821691505b602082108103610e3157610e31610eb4565b6020808252810161095b81602681527f43616e6e6f742062696e6420746f20667574757265206f722063757272656e7460208201527f20626c6f636b0000000000000000000000000000000000000000000000000000604082015260600190565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561095b5761095b610f4f565b60158152602081017f426c6f636b2062696e64696e6720746f6f206f6c640000000000000000000000815290505b60200190565b6020808252810161095b81610f76565b60128152602081017f556e6b6e6f776e20626c6f636b2068617368000000000000000000000000000081529050610fa4565b6020808252810161095b81610fba565b60168152602081017f426c6f636b2062696e64696e67206d69736d617463680000000000000000000081529050610fa4565b6020808252810161095b81610ffc565b60148152602081017f426c6f622068617368206973206e6f742073657400000000000000000000000081529050610fa4565b6020808252810161095b8161103e565b80610c0e565b6110908188611080565b60200161109d8187611080565b6020016110aa8186611080565b6020016110b78185611080565b6020016110c48184611080565b6020016110d18183611080565b6020019695505050505050565b5f808335601e19368590030181126110f7576110f75f5ffd5b8301915050803567ffffffffffffffff811115611115576111155f5ffd5b60208201915060018102360382131561112f5761112f5f5ffd5b9250929050565b801515610c3b565b805161095b81611136565b5f6020828403121561115c5761115c5f5ffd5b610cb4838361113e565b60168152602081017f656e636c6176654944206e6f742061747465737465640000000000000000000081529050610fa4565b6020808252810161095b81611166565b60198152602081017f656e636c6176654944206e6f7420612073657175656e6365720000000000000081529050610fa4565b6020808252810161095b816111a8565b604081016111f88285610c0c565b610cb46020830184610c0c565b82818337505f910152565b818352602083019250611224828483611205565b50601f01601f19160190565b6040810161123e8286610c0c565b8181036020830152610d89818486611210565b5f813561095b81610c86565b5f8161095b565b61126d8261125d565b611279610dcd8261125d565b8255505050565b5f61095b610dcd8381565b61129482611280565b80611279565b634e487b7160e01b5f52604160045260245ffd5b6112b783611280565b81545f1960089490940293841b1916921b91909117905550565b5f6112dd8184846112ae565b505050565b81811015610b4e576112f45f826112d1565b6001016112e2565b601f8211156112dd575f818152602090206020601f850104810160208510156113225750805b6113346020601f8601048301826112e2565b5050505050565b8267ffffffffffffffff8111156113545761135461129a565b61135e8254610ec8565b6113698282856112fc565b505f601f82116001811461139b575f83156113845750848201355b5f19600885021c19811660028502178555506113f2565b5f84815260208120601f198516915b828110156113ca57878501358255602094850194600190920191016113aa565b50848210156113e6575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b6112dd83838361133b565b818061141081611251565b905061141c8184611264565b5050602082018061142c82611251565b905061143b816001850161128b565b5050604082018061144b82611251565b905061145a8160028501611264565b5050606082018061146a82611251565b9050611479816003850161128b565b5050608082018061148982611251565b90506114988160048501611264565b505060a08201806114a882611251565b90506114b78160058501611264565b50506114c660c08301836110de565b6114d48183600686016113fa565b50505050565b610b4e8282611405565b634e487b7160e01b5f52602160045260245ffd5b60ff8116610c0e565b6080810161150f8287610c0c565b61151c60208301866114f8565b6115296040830185610c0c565b610d896060830184610c0c56fea2646970667358221220eaa96b2da78977df8d5ce8afb369b392003cd0155b53138af813152daaf2dee564736f6c634300081c0033",
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
