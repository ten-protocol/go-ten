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

// IDataAvailabilityRegistryMetaRollup is an auto generated low-level Go binding around an user-defined struct.
type IDataAvailabilityRegistryMetaRollup struct {
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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"RollupAdded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structIDataAvailabilityRegistry.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"addRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enclaveRegistry\",\"outputs\":[{\"internalType\":\"contractINetworkEnclaveRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChallengePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"getRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structIDataAvailabilityRegistry.MetaRollup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_merkleMessageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_enclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSeqNo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_delay\",\"type\":\"uint256\"}],\"name\":\"setChallengePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506015601f565b601b601f565b60cf565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff1615606e5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460cc5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b6117f8806100dc5f395ff3fe608060405234801561000f575f5ffd5b50600436106100da575f3560e01c806379ba509711610088578063c0c53b8b11610063578063c0c53b8b14610193578063e30c3978146101a6578063e874eb20146101ae578063f2fde38b146101c1575f5ffd5b806379ba5097146101565780637c72dbd01461015e5780638da5cb5b1461017e575f5ffd5b80636fb6a45c116100b85780636fb6a45c14610125578063715018a6146101465780637864b77d1461014e575f5ffd5b8063440c953b146100de5780635d475fdd146100fd5780635fdf31a214610112575b5f5ffd5b6100e760025481565b6040516100f49190610de1565b60405180910390f35b61011061010b366004610e06565b6101d4565b005b610110610120366004610e44565b6101e1565b610138610133366004610e06565b610577565b6040516100f4929190610f6f565b6101106106bc565b6003546100e7565b6101106106dc565b600554610171906001600160a01b031681565b6040516100f49190610fce565b61018661071b565b6040516100f49190610ff5565b6101106101a1366004611017565b61074f565b6101866108d7565b600454610171906001600160a01b031681565b6101106101cf36600461105d565b6108ff565b6101dc610991565b600355565b808060800135431161020e5760405162461bcd60e51b8152600401610205906110d4565b60405180910390fd5b61021d608082013560ff6110f8565b431061023b5760405162461bcd60e51b81526004016102059061113f565b6080810135405f8190036102615760405162461bcd60e51b815260040161020590611181565b816060013581146102845760405162461bcd60e51b8152600401610205906111c3565b5f496102a25760405162461bcd60e51b815260040161020590611205565b5f826020013583604001358460c00135856060013586608001358760a001355f496040516020016102d9979695949392919061121b565b60408051601f19818403018152919052805160209091012090505f61033e8261030560e0870187611281565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506109c592505050565b6005546040517f3c23afba0000000000000000000000000000000000000000000000000000000081529192506001600160a01b031690633c23afba90610388908490600401610ff5565b602060405180830381865afa1580156103a3573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103c791906112ec565b6103e35760405162461bcd60e51b81526004016102059061133b565b6005546040517f6d46e9870000000000000000000000000000000000000000000000000000000081526001600160a01b0390911690636d46e9879061042c908490600401610ff5565b602060405180830381865afa158015610447573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061046b91906112ec565b6104875760405162461bcd60e51b81526004016102059061137d565b610490856109ef565b60a08501355f1914610529575f600354426104ab91906110f8565b600480546040517fb6aed0cb0000000000000000000000000000000000000000000000000000000081529293506001600160a01b03169163b6aed0cb916104fa9160a08b01359186910161138d565b5f604051808303815f87803b158015610511575f5ffd5b505af1158015610523573d5f5f3e3d5ffd5b50505050505b7fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b12455f4961055960e0880188611281565b604051610568939291906113d3565b60405180910390a15050505050565b5f6105b96040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081525090565b5f5f5f015f8581526020019081526020015f20604051806101000160405290815f820154815260200160018201548152602001600282015481526020016003820154815260200160048201548152602001600582015481526020016006820154815260200160078201805461062d90611408565b80601f016020809104026020016040519081016040528092919081815260200182805461065990611408565b80156106a45780601f1061067b576101008083540402835291602001916106a4565b820191905f5260205f20905b81548152906001019060200180831161068757829003601f168201915b50505091909252505081519095149590945092505050565b6106c4610991565b60405162461bcd60e51b815260040161020590611486565b33806106e66108d7565b6001600160a01b03161461070f578060405163118cdaa760e01b81526004016102059190610ff5565b61071881610a2c565b50565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156107995750825b90505f8267ffffffffffffffff1660011480156107b55750303b155b9050811580156107c3575080155b156107fa576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561082e57845468ff00000000000000001916680100000000000000001785555b61083786610a75565b600480546001600160a01b03808b1673ffffffffffffffffffffffffffffffffffffffff199283161790925560058054928a16929091169190911790555f600281905560035583156108cd57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906108c4906001906114b0565b60405180910390a15b5050505050505050565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c0061073f565b610907610991565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117825561095861071b565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b3361099a61071b565b6001600160a01b0316146109c3573360405163118cdaa760e01b81526004016102059190610ff5565b565b5f5f5f5f6109d38686610a8e565b9250925092506109e38282610ad7565b50909150505b92915050565b5f804981526020819052604090208190610a098282611766565b5050600254610a199060016110f8565b8160200135036107185760400135600255565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19168155610a7182610bd8565b5050565b610a7d610c55565b610a8681610cbc565b610718610ccd565b5f5f5f8351604103610ac5576020840151604085015160608601515f1a610ab788828585610cd5565b955095509550505050610ad0565b505081515f91506002905b9250925092565b5f826003811115610aea57610aea611770565b03610af3575050565b6001826003811115610b0757610b07611770565b03610b3e576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6002826003811115610b5257610b52611770565b03610b8b576040517ffce698f7000000000000000000000000000000000000000000000000000000008152610205908290600401610de1565b6003826003811115610b9f57610b9f611770565b03610a7157806040517fd78bce0c0000000000000000000000000000000000000000000000000000000081526004016102059190610de1565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166109c3576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610cc4610c55565b61071881610d8f565b6109c3610c55565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610d0e57505f91506003905082610d85565b5f6001888888886040515f8152602001604052604051610d31949392919061178d565b6020604051602081039080840390855afa158015610d51573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116610d7c57505f925060019150829050610d85565b92505f91508190505b9450945094915050565b610d97610c55565b6001600160a01b03811661070f575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016102059190610ff5565b805b82525050565b602081016109e98284610dd9565b805b8114610718575f5ffd5b80356109e981610def565b5f60208284031215610e1957610e195f5ffd5b610e238383610dfb565b9392505050565b5f6101008284031215610e3e57610e3e5f5ffd5b50919050565b5f60208284031215610e5757610e575f5ffd5b813567ffffffffffffffff811115610e7057610e705f5ffd5b610e7c84828501610e2a565b949350505050565b801515610ddb565b8281835e505f910152565b5f610ea0825190565b808452602084019350610eb7818560208601610e8c565b601f01601f19169290920192915050565b80515f90610100840190610edc8582610dd9565b506020830151610eef6020860182610dd9565b506040830151610f026040860182610dd9565b506060830151610f156060860182610dd9565b506080830151610f286080860182610dd9565b5060a0830151610f3b60a0860182610dd9565b5060c0830151610f4e60c0860182610dd9565b5060e083015184820360e0860152610f668282610e97565b95945050505050565b60408101610f7d8285610e84565b8181036020830152610e7c8184610ec8565b5f6109e96001600160a01b038316610fa5565b90565b6001600160a01b031690565b5f6109e982610f8f565b5f6109e982610fb1565b610ddb81610fbb565b602081016109e98284610fc5565b5f6001600160a01b0382166109e9565b610ddb81610fdc565b602081016109e98284610fec565b610df181610fdc565b80356109e981611003565b5f5f5f6060848603121561102c5761102c5f5ffd5b611036858561100c565b9250611045856020860161100c565b9150611054856040860161100c565b90509250925092565b5f60208284031215611070576110705f5ffd5b610e23838361100c565b60268152602081017f43616e6e6f742062696e6420746f20667574757265206f722063757272656e7481527f20626c6f636b0000000000000000000000000000000000000000000000000000602082015290505b60400190565b602080825281016109e98161107a565b634e487b7160e01b5f52601160045260245ffd5b808201808211156109e9576109e96110e4565b60158152602081017f426c6f636b2062696e64696e6720746f6f206f6c640000000000000000000000815290505b60200190565b602080825281016109e98161110b565b60128152602081017f556e6b6e6f776e20626c6f636b2068617368000000000000000000000000000081529050611139565b602080825281016109e98161114f565b60168152602081017f426c6f636b2062696e64696e67206d69736d617463680000000000000000000081529050611139565b602080825281016109e981611191565b60148152602081017f426c6f622068617368206973206e6f742073657400000000000000000000000081529050611139565b602080825281016109e9816111d3565b80610ddb565b6112258189611215565b6020016112328188611215565b60200161123f8187611215565b60200161124c8186611215565b6020016112598185611215565b6020016112668184611215565b6020016112738183611215565b602001979650505050505050565b5f808335601e193685900301811261129a5761129a5f5ffd5b8301915050803567ffffffffffffffff8111156112b8576112b85f5ffd5b6020820191506001810236038213156112d2576112d25f5ffd5b9250929050565b801515610df1565b80516109e9816112d9565b5f602082840312156112ff576112ff5f5ffd5b610e2383836112e1565b60168152602081017f656e636c6176654944206e6f742061747465737465640000000000000000000081529050611139565b602080825281016109e981611309565b60198152602081017f656e636c6176654944206e6f7420612073657175656e6365720000000000000081529050611139565b602080825281016109e98161134b565b6040810161139b8285610dd9565b610e236020830184610dd9565b82818337505f910152565b8183526020830192506113c78284836113a8565b50601f01601f19160190565b604081016113e18286610dd9565b8181036020830152610f668184866113b3565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061141c57607f821691505b602082108103610e3e57610e3e6113f4565b60348152602081017f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f81527f742072656e6f756e6365206f776e657273686970000000000000000000000000602082015290506110ce565b602080825281016109e98161142e565b5f67ffffffffffffffff82166109e9565b610ddb81611496565b602081016109e982846114a7565b5f81356109e981610def565b5f816109e9565b6114da826114ca565b6114e6610fa2826114ca565b8255505050565b5f6109e9610fa28381565b611501826114ed565b806114e6565b634e487b7160e01b5f52604160045260245ffd5b611524836114ed565b81545f1960089490940293841b1916921b91909117905550565b5f61154a81848461151b565b505050565b81811015610a71576115615f8261153e565b60010161154f565b601f82111561154a575f818152602090206020601f8501048101602085101561158f5750805b6115a16020601f86010483018261154f565b5050505050565b8267ffffffffffffffff8111156115c1576115c1611507565b6115cb8254611408565b6115d6828285611569565b505f601f821160018114611608575f83156115f15750848201355b5f19600885021c198116600285021785555061165f565b5f84815260208120601f198516915b828110156116375787850135825560209485019460019092019101611617565b5084821015611653575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b61154a8383836115a8565b818061167d816114be565b905061168981846114d1565b50506020820180611699826114be565b90506116a881600185016114f8565b505060408201806116b8826114be565b90506116c781600285016114f8565b505060608201806116d7826114be565b90506116e681600385016114d1565b505060808201806116f6826114be565b905061170581600485016114f8565b505060a0820180611715826114be565b905061172481600585016114d1565b505060c0820180611734826114be565b905061174381600685016114d1565b505061175260e0830183611281565b611760818360078601611667565b50505050565b610a718282611672565b634e487b7160e01b5f52602160045260245ffd5b60ff8116610ddb565b6080810161179b8287610dd9565b6117a86020830186611784565b6117b56040830185610dd9565b610f666060830184610dd956fea2646970667358221220c7ca6ebbbe95262d123fad635652d742d29de4d4c4120142b3405cbe19d3028064736f6c634300081c0033",
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
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) GetRollupByHash(opts *bind.CallOpts, rollupHash [32]byte) (bool, IDataAvailabilityRegistryMetaRollup, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "getRollupByHash", rollupHash)

	if err != nil {
		return *new(bool), *new(IDataAvailabilityRegistryMetaRollup), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(IDataAvailabilityRegistryMetaRollup)).(*IDataAvailabilityRegistryMetaRollup)

	return out0, out1, err

}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) GetRollupByHash(rollupHash [32]byte) (bool, IDataAvailabilityRegistryMetaRollup, error) {
	return _DataAvailabilityRegistry.Contract.GetRollupByHash(&_DataAvailabilityRegistry.CallOpts, rollupHash)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) GetRollupByHash(rollupHash [32]byte) (bool, IDataAvailabilityRegistryMetaRollup, error) {
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

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) PendingOwner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.PendingOwner(&_DataAvailabilityRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) PendingOwner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.PendingOwner(&_DataAvailabilityRegistry.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) RenounceOwnership() error {
	return _DataAvailabilityRegistry.Contract.RenounceOwnership(&_DataAvailabilityRegistry.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) RenounceOwnership() error {
	return _DataAvailabilityRegistry.Contract.RenounceOwnership(&_DataAvailabilityRegistry.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AcceptOwnership(&_DataAvailabilityRegistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AcceptOwnership(&_DataAvailabilityRegistry.TransactOpts)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) AddRollup(opts *bind.TransactOpts, r IDataAvailabilityRegistryMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "addRollup", r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) AddRollup(r IDataAvailabilityRegistryMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AddRollup(&_DataAvailabilityRegistry.TransactOpts, r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) AddRollup(r IDataAvailabilityRegistryMetaRollup) (*types.Transaction, error) {
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

// DataAvailabilityRegistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferStartedIterator struct {
	Event *DataAvailabilityRegistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryOwnershipTransferStarted)
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
		it.Event = new(DataAvailabilityRegistryOwnershipTransferStarted)
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
func (it *DataAvailabilityRegistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DataAvailabilityRegistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryOwnershipTransferStartedIterator{contract: _DataAvailabilityRegistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryOwnershipTransferStarted)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*DataAvailabilityRegistryOwnershipTransferStarted, error) {
	event := new(DataAvailabilityRegistryOwnershipTransferStarted)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
