// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MerkleTreeMessageBus

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

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// StructsValueTransferMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsValueTransferMessage struct {
	Sender   common.Address
	Receiver common.Address
	Amount   *big.Int
	Sequence uint64
}

// MerkleTreeMessageBusMetaData contains all meta data concerning the MerkleTreeMessageBus contract.
var MerkleTreeMessageBusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"name\":\"ValueTransfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"activationTime\",\"type\":\"uint256\"}],\"name\":\"addStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"name\":\"disableStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feesAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"receiveValueFromL2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"retrieveAllFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendValueToL2\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyMessageInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.ValueTransferMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyValueTransferInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061001a33610027565b610022610098565b61014a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100e85760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146101475780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b612126806101596000396000f3fe6080604052600436106100ec5760003560e01c80638da5cb5b1161008a578063b201246f11610059578063b201246f14610311578063b6aed0cb14610331578063e138a8d214610351578063f2fde38b1461037157610160565b80638da5cb5b1461026c5780639730886d146102b157806399a3ad21146102d1578063b1454caa146102f157610160565b8063346633fb116100c6578063346633fb1461020457806336d2da9014610217578063485cc95514610237578063715018a61461025757610160565b80630fcfbd11146101815780630fe9188e146101b757806333a88c72146101d757610160565b36610160576040517f346633fb000000000000000000000000000000000000000000000000000000008152309063346633fb9034906101319033908390600401610f76565b6000604051808303818588803b15801561014a57600080fd5b505af115801561015e573d6000803e3d6000fd5b005b60405162461bcd60e51b815260040161017890610fc5565b60405180910390fd5b34801561018d57600080fd5b506101a161019c366004610ff0565b610391565b6040516101ae919061102b565b60405180910390f35b3480156101c357600080fd5b5061015e6101d2366004611051565b6103f0565b3480156101e357600080fd5b506101f76101f2366004610ff0565b610436565b6040516101ae9190611078565b61015e61021236600461109a565b610488565b34801561022357600080fd5b5061015e6102323660046110d2565b610512565b34801561024357600080fd5b5061015e6102523660046110f1565b610591565b34801561026357600080fd5b5061015e6106fc565b34801561027857600080fd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516101ae9190611120565b3480156102bd57600080fd5b5061015e6102cc36600461112e565b610710565b3480156102dd57600080fd5b5061015e6102ec36600461109a565b61087c565b6103046102ff3660046111f7565b6108fc565b6040516101ae9190611284565b34801561031d57600080fd5b5061015e61032c3660046112f2565b61098f565b34801561033d57600080fd5b5061015e61034c36600461135d565b610a90565b34801561035d57600080fd5b5061015e61036c36600461137d565b610ad6565b34801561037d57600080fd5b5061015e61038c3660046110d2565b610c21565b600080826040516020016103a59190611599565b60408051601f198184030181529181528151602092830120600081815292839052912054909150806103e95760405162461bcd60e51b8152600401610178906115e8565b9392505050565b6103f8610c78565b60008181526004602052604081205490036104255760405162461bcd60e51b81526004016101789061162a565b600090815260046020526040812055565b6000808260405160200161044a9190611599565b60408051601f19818403018152918152815160209283012060008181529283905291205490915080158015906104805750428111155b949350505050565b60003411801561049757508034145b6104b35760405162461bcd60e51b815260040161017890611692565b60006104be33610cec565b9050826001600160a01b0316336001600160a01b03167f50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b634846040516105059291906116a2565b60405180910390a3505050565b61051a610c78565b6000816001600160a01b03164760405160006040518083038185875af1925050503d8060008114610567576040519150601f19603f3d011682016040523d82523d6000602084013e61056c565b606091505b505090508061058d5760405162461bcd60e51b8152600401610178906116ef565b5050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156105dc5750825b905060008267ffffffffffffffff1660011480156105f95750303b155b905081158015610607575080155b1561063e576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561067257845468ff00000000000000001916680100000000000000001785555b61067b87610d4a565b6003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03881617905583156106f357845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906106ea90600190611723565b60405180910390a15b50505050505050565b610704610c78565b61070e6000610d5b565b565b600061071d600130611747565b90506107507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b031614806107775750336001600160a01b038216145b6107935760405162461bcd60e51b81526004016101789061179c565b600061079f83426117ac565b90506000846040516020016107b49190611599565b60408051601f198184030181529181528151602092830120600081815292839052912054909150156107f85760405162461bcd60e51b815260040161017890611817565b600081815260208181526040822084905560019190610819908801886110d2565b6001600160a01b0316815260208101919091526040016000908120906108456080880160608901611827565b63ffffffff168152602080820192909252604001600090812080546001810182559082529190208691600402016106f38282611c65565b610884610c78565b6000826001600160a01b03168260405160006040518083038185875af1925050503d80600081146108d1576040519150601f19603f3d011682016040523d82523d6000602084013e6108d6565b606091505b50509050806108f75760405162461bcd60e51b8152600401610178906116ef565b505050565b6003546000906001600160a01b0316156109385761091983610dd9565b3410156109385760405162461bcd60e51b815260040161017890611cc7565b61094133610cec565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef7759373382888888888860405161097e9796959493929190611cd7565b60405180910390a195945050505050565b60008181526004602052604081205490036109bc5760405162461bcd60e51b815260040161017890611d92565b6000818152600460205260409020544210156109ea5760405162461bcd60e51b815260040161017890611dde565b6000846040516020016109fd9190611e63565b60405160208183030381529060405280519060200120604051602001610a239190611ea3565b604051602081830303815290604052805190602001209050610a6d84848484604051602001610a529190611ec2565b60405160208183030381529060405280519060200120610e5b565b610a895760405162461bcd60e51b815260040161017890611f2c565b5050505050565b610a98610c78565b60008281526004602052604090205415610ac45760405162461bcd60e51b815260040161017890611f94565b60009182526004602052604090912055565b6000818152600460205260408120549003610b035760405162461bcd60e51b815260040161017890611d92565b600081815260046020526040902054421015610b315760405162461bcd60e51b815260040161017890611dde565b6000610b4060208601866110d2565b610b506040870160208801611fa4565b610b606060880160408901611827565b610b706080890160608a01611827565b610b7d60808a018a61197e565b610b8d60c08c0160a08d01611fc3565b604051602001610ba39796959493929190611cd7565b604051602081830303815290604052805190602001209050600081604051602001610bce9190612014565b604051602081830303815290604052805190602001209050610bfd85858584604051602001610a529190611ec2565b610c195760405162461bcd60e51b81526004016101789061207c565b505050505050565b610c29610c78565b6001600160a01b038116610c6c5760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016101789190611120565b610c7581610d5b565b50565b33610caa7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461070e57336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016101789190611120565b6001600160a01b0381166000908152600260205260408120805467ffffffffffffffff169160019190610d1f838561208c565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b610d52610e73565b610c7581610eda565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b6003546000906001600160a01b031663f1d44d51610df86011856117ac565b6040518263ffffffff1660e01b8152600401610e14919061102b565b602060405180830381865afa158015610e31573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e5591906120bb565b92915050565b600082610e69868685610ee2565b1495945050505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff1661070e576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c29610e73565b600081815b84811015610f1b57610f1182878784818110610f0557610f056120da565b90506020020135610f24565b9150600101610ee7565b50949350505050565b6000818310610f405760008281526020849052604090206103e9565b5060009182526020526040902090565b60006001600160a01b038216610e55565b610f6a81610f50565b82525050565b80610f6a565b60408101610f848285610f61565b6103e96020830184610f70565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b60208082528101610e5581610f91565b600060c08284031215610fea57610fea600080fd5b50919050565b60006020828403121561100557611005600080fd5b813567ffffffffffffffff81111561101f5761101f600080fd5b61048084828501610fd5565b60208101610e558284610f70565b805b8114610c7557600080fd5b8035610e5581611039565b60006020828403121561106657611066600080fd5b6103e98383611046565b801515610f6a565b60208101610e558284611070565b61103b81610f50565b8035610e5581611086565b600080604083850312156110b0576110b0600080fd5b6110ba848461108f565b91506110c98460208501611046565b90509250929050565b6000602082840312156110e7576110e7600080fd5b6103e9838361108f565b6000806040838503121561110757611107600080fd5b611111848461108f565b91506110c9846020850161108f565b60208101610e558284610f61565b6000806040838503121561114457611144600080fd5b823567ffffffffffffffff81111561115e5761115e600080fd5b61116a85828601610fd5565b9250506110c98460208501611046565b63ffffffff811661103b565b8035610e558161117a565b60008083601f8401126111a6576111a6600080fd5b50813567ffffffffffffffff8111156111c1576111c1600080fd5b6020830191508360018202830111156111dc576111dc600080fd5b9250929050565b60ff811661103b565b8035610e55816111e3565b60008060008060006080868803121561121257611212600080fd5b61121c8787611186565b945061122b8760208801611186565b9350604086013567ffffffffffffffff81111561124a5761124a600080fd5b61125688828901611191565b935093505061126887606088016111ec565b90509295509295909350565b67ffffffffffffffff8116610f6a565b60208101610e558284611274565b600060808284031215610fea57610fea600080fd5b60008083601f8401126112bc576112bc600080fd5b50813567ffffffffffffffff8111156112d7576112d7600080fd5b6020830191508360208202830111156111dc576111dc600080fd5b60008060008060c0858703121561130b5761130b600080fd5b6113158686611292565b9350608085013567ffffffffffffffff81111561133457611334600080fd5b611340878288016112a7565b93509350506113528660a08701611046565b905092959194509250565b6000806040838503121561137357611373600080fd5b6110ba8484611046565b6000806000806060858703121561139657611396600080fd5b843567ffffffffffffffff8111156113b0576113b0600080fd5b6113bc87828801610fd5565b945050602085013567ffffffffffffffff8111156113dc576113dc600080fd5b6113e8878288016112a7565b93509350506113528660408701611046565b506000610e55602083018361108f565b67ffffffffffffffff811661103b565b8035610e558161140a565b506000610e55602083018361141a565b506000610e556020830183611186565b63ffffffff8116610f6a565b6000808335601e193685900301811261146c5761146c600080fd5b830160208101925035905067ffffffffffffffff81111561148f5761148f600080fd5b368190038213156111dc576111dc600080fd5b82818337506000910152565b8183526020830192506114c28284836114a2565b50601f01601f19160190565b506000610e5560208301836111ec565b60ff8116610f6a565b600060c083016114f783806113fa565b6115018582610f61565b5061150f6020840184611425565b61151c6020860182611274565b5061152a6040840184611435565b6115376040860182611445565b506115456060840184611435565b6115526060860182611445565b506115606080840184611451565b85830360808701526115738382846114ae565b9250505061158460a08401846114ce565b61159160a08601826114de565b509392505050565b602080825281016103e981846114e7565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d69747465648152601760f91b602082015290505b60400190565b60208082528101610e55816115aa565b601a8152602081017f537461746520726f6f7420646f6573206e6f742065786973742e00000000000081529050610fbf565b60208082528101610e55816115f8565b60308152602081017f417474656d7074696e6720746f2073656e642076616c756520776974686f757481527f2070726f766964696e6720457468657200000000000000000000000000000000602082015290506115e2565b60208082528101610e558161163a565b604081016116b08285610f70565b6103e96020830184611274565b60148152602081017f6661696c65642073656e64696e672076616c756500000000000000000000000081529050610fbf565b60208082528101610e55816116bd565b6000610e558261170d565b90565b67ffffffffffffffff1690565b610f6a816116ff565b60208101610e55828461171a565b634e487b7160e01b600052601160045260246000fd5b6001600160a01b03918216919081169082820390811115610e5557610e55611731565b60118152602081017f4e6f74206f776e6572206f722073656c6600000000000000000000000000000081529050610fbf565b60208082528101610e558161176a565b80820180821115610e5557610e55611731565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f2100000000000000000000000000000000000000000000000000000000000000602082015290506115e2565b60208082528101610e55816117bf565b60006020828403121561183c5761183c600080fd5b6103e98383611186565b60008135610e5581611086565b60006001600160a01b03835b81169019929092169190911792915050565b6000610e5582610f50565b6000610e5582611871565b6118908261187c565b61189b818354611853565b8255505050565b60008135610e558161140a565b60007bffffffffffffffff000000000000000000000000000000000000000061185f8460a01b90565b6000610e5567ffffffffffffffff831661170d565b6118f6826118d8565b61189b8183546118af565b60008135610e558161117a565b60007fffffffff0000000000000000000000000000000000000000000000000000000061185f8460e01b90565b600063ffffffff8216610e55565b6119528261193b565b61189b81835461190e565b600063ffffffff8361185f565b6119738261193b565b61189b81835461195d565b6000808335601e193685900301811261199957611999600080fd5b8301915050803567ffffffffffffffff8111156119b8576119b8600080fd5b6020820191506001810236038213156111dc576111dc600080fd5b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052602260045260246000fd5b600281046001821680611a1357607f821691505b602082108103610fea57610fea6119e9565b6000610e5561170a8381565b611a3a83611a25565b815460001960089490940293841b1916921b91909117905550565b60006108f7818484611a31565b8181101561058d57611a75600082611a55565b600101611a62565b601f8211156108f7576000818152602090206020601f85010481016020851015611aa45750805b610a896020601f860104830182611a62565b8267ffffffffffffffff811115611acf57611acf6119d3565b611ad982546119ff565b611ae4828285611a7d565b506000601f821160018114611b195760008315611b015750848201355b600019600885021c1981166002850217855550610c19565b600084815260209020601f19841690835b82811015611b4a5787850135825560209485019460019092019101611b2a565b5084821015611b67576000196008601f8716021c19878501351681555b5050505060020260010190555050565b6108f7838383611ab6565b60008135610e55816111e3565b600060ff8216610e55565b611ba382611b8f565b815460ff191660ff82161761189b565b808280611bbf81611846565b9050611bcb8184611887565b50506020830180611bdb826118a2565b9050611be781846118ed565b50506040830180611bf782611901565b9050611c038184611949565b5050506060820180611c1482611901565b9050611c23816001850161196a565b5050611c32608083018361197e565b611c40818360028601611b77565b505060a0820180611c5082611b82565b9050611c5f8160038501611b9a565b50505050565b61058d8282611bb3565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f7373616765000000000000000000000000000000000000000000000000000000602082015290506115e2565b60208082528101610e5581611c6f565b60c08101611ce5828a610f61565b611cf26020830189611274565b611cff6040830188611445565b611d0c6060830187611445565b8181036080830152611d1f8185876114ae565b9050611d2e60a08301846114de565b98975050505050505050565b602a8152602081017f526f6f74206973206e6f74207075626c6973686564206f6e2074686973206d6581527f7373616765206275732e00000000000000000000000000000000000000000000602082015290506115e2565b60208082528101610e5581611d3a565b60218152602081017f526f6f74206973206e6f7420636f6e736964657265642066696e616c207965748152601760f91b602082015290506115e2565b60208082528101610e5581611da2565b506000610e556020830183611046565b611e0881806113fa565b611e128382610f61565b50611e2060208201826113fa565b611e2d6020840182610f61565b50611e3b6040820182611dee565b611e486040840182610f70565b50611e566060820182611425565b6108f76060840182611274565b60808101610e558284611dfe565b60018152602081017f760000000000000000000000000000000000000000000000000000000000000081529050610fbf565b60408082528101611eb381611e71565b9050610e556020830184610f70565b611ecc8183610f70565b602001919050565b60338152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722076616c7581527f65207472616e73666572206d6573736167652e00000000000000000000000000602082015290506115e2565b60208082528101610e5581611ed4565b60258152602081017f526f6f7420616c726561647920616464656420746f20746865206d657373616781527f6520627573000000000000000000000000000000000000000000000000000000602082015290506115e2565b60208082528101610e5581611f3c565b600060208284031215611fb957611fb9600080fd5b6103e9838361141a565b600060208284031215611fd857611fd8600080fd5b6103e983836111ec565b60018152602081017f6d0000000000000000000000000000000000000000000000000000000000000081529050610fbf565b60408082528101611eb381611fe2565b60308152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722063726f7381527f7320636861696e206d6573736167652e00000000000000000000000000000000602082015290506115e2565b60208082528101610e5581612024565b67ffffffffffffffff918216919081169082820190811115610e5557610e55611731565b8051610e5581611039565b6000602082840312156120d0576120d0600080fd5b6103e983836120b0565b634e487b7160e01b600052603260045260246000fdfea26469706673582212201f7f0f6c6bd3365796e0b09ad74d1b7971f84ca4f811483bd29ce1cb972c84a064736f6c634300081c0033",
}

// MerkleTreeMessageBusABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleTreeMessageBusMetaData.ABI instead.
var MerkleTreeMessageBusABI = MerkleTreeMessageBusMetaData.ABI

// MerkleTreeMessageBusBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MerkleTreeMessageBusMetaData.Bin instead.
var MerkleTreeMessageBusBin = MerkleTreeMessageBusMetaData.Bin

// DeployMerkleTreeMessageBus deploys a new Ethereum contract, binding an instance of MerkleTreeMessageBus to it.
func DeployMerkleTreeMessageBus(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MerkleTreeMessageBus, error) {
	parsed, err := MerkleTreeMessageBusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MerkleTreeMessageBusBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MerkleTreeMessageBus{MerkleTreeMessageBusCaller: MerkleTreeMessageBusCaller{contract: contract}, MerkleTreeMessageBusTransactor: MerkleTreeMessageBusTransactor{contract: contract}, MerkleTreeMessageBusFilterer: MerkleTreeMessageBusFilterer{contract: contract}}, nil
}

// MerkleTreeMessageBus is an auto generated Go binding around an Ethereum contract.
type MerkleTreeMessageBus struct {
	MerkleTreeMessageBusCaller     // Read-only binding to the contract
	MerkleTreeMessageBusTransactor // Write-only binding to the contract
	MerkleTreeMessageBusFilterer   // Log filterer for contract events
}

// MerkleTreeMessageBusCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeMessageBusTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeMessageBusFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleTreeMessageBusFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeMessageBusSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleTreeMessageBusSession struct {
	Contract     *MerkleTreeMessageBus // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MerkleTreeMessageBusCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleTreeMessageBusCallerSession struct {
	Contract *MerkleTreeMessageBusCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// MerkleTreeMessageBusTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleTreeMessageBusTransactorSession struct {
	Contract     *MerkleTreeMessageBusTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// MerkleTreeMessageBusRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleTreeMessageBusRaw struct {
	Contract *MerkleTreeMessageBus // Generic contract binding to access the raw methods on
}

// MerkleTreeMessageBusCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusCallerRaw struct {
	Contract *MerkleTreeMessageBusCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleTreeMessageBusTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusTransactorRaw struct {
	Contract *MerkleTreeMessageBusTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleTreeMessageBus creates a new instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBus(address common.Address, backend bind.ContractBackend) (*MerkleTreeMessageBus, error) {
	contract, err := bindMerkleTreeMessageBus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBus{MerkleTreeMessageBusCaller: MerkleTreeMessageBusCaller{contract: contract}, MerkleTreeMessageBusTransactor: MerkleTreeMessageBusTransactor{contract: contract}, MerkleTreeMessageBusFilterer: MerkleTreeMessageBusFilterer{contract: contract}}, nil
}

// NewMerkleTreeMessageBusCaller creates a new read-only instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBusCaller(address common.Address, caller bind.ContractCaller) (*MerkleTreeMessageBusCaller, error) {
	contract, err := bindMerkleTreeMessageBus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusCaller{contract: contract}, nil
}

// NewMerkleTreeMessageBusTransactor creates a new write-only instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBusTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleTreeMessageBusTransactor, error) {
	contract, err := bindMerkleTreeMessageBus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusTransactor{contract: contract}, nil
}

// NewMerkleTreeMessageBusFilterer creates a new log filterer instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBusFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleTreeMessageBusFilterer, error) {
	contract, err := bindMerkleTreeMessageBus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusFilterer{contract: contract}, nil
}

// bindMerkleTreeMessageBus binds a generic wrapper to an already deployed contract.
func bindMerkleTreeMessageBus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MerkleTreeMessageBusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleTreeMessageBus.Contract.MerkleTreeMessageBusCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.MerkleTreeMessageBusTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.MerkleTreeMessageBusTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleTreeMessageBus.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.contract.Transact(opts, method, params...)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetMessageTimeOfFinality(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetMessageTimeOfFinality(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Owner() (common.Address, error) {
	return _MerkleTreeMessageBus.Contract.Owner(&_MerkleTreeMessageBus.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) Owner() (common.Address, error) {
	return _MerkleTreeMessageBus.Contract.Owner(&_MerkleTreeMessageBus.CallOpts)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MerkleTreeMessageBus.Contract.VerifyMessageFinalized(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MerkleTreeMessageBus.Contract.VerifyMessageFinalized(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xe138a8d2.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint32,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyMessageInclusion(opts *bind.CallOpts, message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyMessageInclusion", message, proof, root)

	if err != nil {
		return err
	}

	return err

}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xe138a8d2.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint32,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyMessageInclusion(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyMessageInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xe138a8d2.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint32,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) VerifyMessageInclusion(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyMessageInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// VerifyValueTransferInclusion is a free data retrieval call binding the contract method 0xb201246f.
//
// Solidity: function verifyValueTransferInclusion((address,address,uint256,uint64) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyValueTransferInclusion(opts *bind.CallOpts, message StructsValueTransferMessage, proof [][32]byte, root [32]byte) error {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyValueTransferInclusion", message, proof, root)

	if err != nil {
		return err
	}

	return err

}

// VerifyValueTransferInclusion is a free data retrieval call binding the contract method 0xb201246f.
//
// Solidity: function verifyValueTransferInclusion((address,address,uint256,uint64) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyValueTransferInclusion(message StructsValueTransferMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyValueTransferInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// VerifyValueTransferInclusion is a free data retrieval call binding the contract method 0xb201246f.
//
// Solidity: function verifyValueTransferInclusion((address,address,uint256,uint64) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) VerifyValueTransferInclusion(message StructsValueTransferMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyValueTransferInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// AddStateRoot is a paid mutator transaction binding the contract method 0xb6aed0cb.
//
// Solidity: function addStateRoot(bytes32 stateRoot, uint256 activationTime) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) AddStateRoot(opts *bind.TransactOpts, stateRoot [32]byte, activationTime *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "addStateRoot", stateRoot, activationTime)
}

// AddStateRoot is a paid mutator transaction binding the contract method 0xb6aed0cb.
//
// Solidity: function addStateRoot(bytes32 stateRoot, uint256 activationTime) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) AddStateRoot(stateRoot [32]byte, activationTime *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot, activationTime)
}

// AddStateRoot is a paid mutator transaction binding the contract method 0xb6aed0cb.
//
// Solidity: function addStateRoot(bytes32 stateRoot, uint256 activationTime) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) AddStateRoot(stateRoot [32]byte, activationTime *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot, activationTime)
}

// DisableStateRoot is a paid mutator transaction binding the contract method 0x0fe9188e.
//
// Solidity: function disableStateRoot(bytes32 stateRoot) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) DisableStateRoot(opts *bind.TransactOpts, stateRoot [32]byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "disableStateRoot", stateRoot)
}

// DisableStateRoot is a paid mutator transaction binding the contract method 0x0fe9188e.
//
// Solidity: function disableStateRoot(bytes32 stateRoot) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) DisableStateRoot(stateRoot [32]byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.DisableStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot)
}

// DisableStateRoot is a paid mutator transaction binding the contract method 0x0fe9188e.
//
// Solidity: function disableStateRoot(bytes32 stateRoot) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) DisableStateRoot(stateRoot [32]byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.DisableStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address caller, address feesAddress) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Initialize(opts *bind.TransactOpts, caller common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "initialize", caller, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address caller, address feesAddress) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Initialize(caller common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Initialize(&_MerkleTreeMessageBus.TransactOpts, caller, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address caller, address feesAddress) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Initialize(caller common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Initialize(&_MerkleTreeMessageBus.TransactOpts, caller, feesAddress)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.PublishMessage(&_MerkleTreeMessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.PublishMessage(&_MerkleTreeMessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) ReceiveValueFromL2(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "receiveValueFromL2", receiver, amount)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) ReceiveValueFromL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.ReceiveValueFromL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) ReceiveValueFromL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.ReceiveValueFromL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RenounceOwnership() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RenounceOwnership(&_MerkleTreeMessageBus.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RenounceOwnership(&_MerkleTreeMessageBus.TransactOpts)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RetrieveAllFunds(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "retrieveAllFunds", receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RetrieveAllFunds(&_MerkleTreeMessageBus.TransactOpts, receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RetrieveAllFunds(&_MerkleTreeMessageBus.TransactOpts, receiver)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) SendValueToL2(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "sendValueToL2", receiver, amount)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) SendValueToL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.SendValueToL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) SendValueToL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.SendValueToL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) StoreCrossChainMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "storeCrossChainMessage", crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.StoreCrossChainMessage(&_MerkleTreeMessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.StoreCrossChainMessage(&_MerkleTreeMessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.TransferOwnership(&_MerkleTreeMessageBus.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.TransferOwnership(&_MerkleTreeMessageBus.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Fallback(&_MerkleTreeMessageBus.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Fallback(&_MerkleTreeMessageBus.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Receive() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Receive(&_MerkleTreeMessageBus.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Receive() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Receive(&_MerkleTreeMessageBus.TransactOpts)
}

// MerkleTreeMessageBusInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusInitializedIterator struct {
	Event *MerkleTreeMessageBusInitialized // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusInitialized)
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
		it.Event = new(MerkleTreeMessageBusInitialized)
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
func (it *MerkleTreeMessageBusInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusInitialized represents a Initialized event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterInitialized(opts *bind.FilterOpts) (*MerkleTreeMessageBusInitializedIterator, error) {

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusInitializedIterator{contract: _MerkleTreeMessageBus.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusInitialized) (event.Subscription, error) {

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusInitialized)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseInitialized(log types.Log) (*MerkleTreeMessageBusInitialized, error) {
	event := new(MerkleTreeMessageBusInitialized)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusLogMessagePublishedIterator is returned from FilterLogMessagePublished and is used to iterate over the raw logs and unpacked data for LogMessagePublished events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusLogMessagePublishedIterator struct {
	Event *MerkleTreeMessageBusLogMessagePublished // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusLogMessagePublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusLogMessagePublished)
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
		it.Event = new(MerkleTreeMessageBusLogMessagePublished)
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
func (it *MerkleTreeMessageBusLogMessagePublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusLogMessagePublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusLogMessagePublished represents a LogMessagePublished event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusLogMessagePublished struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogMessagePublished is a free log retrieval operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts) (*MerkleTreeMessageBusLogMessagePublishedIterator, error) {

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusLogMessagePublishedIterator{contract: _MerkleTreeMessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchLogMessagePublished(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusLogMessagePublished) (event.Subscription, error) {

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusLogMessagePublished)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
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

// ParseLogMessagePublished is a log parse operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MerkleTreeMessageBusLogMessagePublished, error) {
	event := new(MerkleTreeMessageBusLogMessagePublished)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusOwnershipTransferredIterator struct {
	Event *MerkleTreeMessageBusOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusOwnershipTransferred)
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
		it.Event = new(MerkleTreeMessageBusOwnershipTransferred)
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
func (it *MerkleTreeMessageBusOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusOwnershipTransferred represents a OwnershipTransferred event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MerkleTreeMessageBusOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusOwnershipTransferredIterator{contract: _MerkleTreeMessageBus.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusOwnershipTransferred)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseOwnershipTransferred(log types.Log) (*MerkleTreeMessageBusOwnershipTransferred, error) {
	event := new(MerkleTreeMessageBusOwnershipTransferred)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusValueTransferIterator is returned from FilterValueTransfer and is used to iterate over the raw logs and unpacked data for ValueTransfer events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusValueTransferIterator struct {
	Event *MerkleTreeMessageBusValueTransfer // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusValueTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusValueTransfer)
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
		it.Event = new(MerkleTreeMessageBusValueTransfer)
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
func (it *MerkleTreeMessageBusValueTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusValueTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusValueTransfer represents a ValueTransfer event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusValueTransfer struct {
	Sender   common.Address
	Receiver common.Address
	Amount   *big.Int
	Sequence uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValueTransfer is a free log retrieval operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterValueTransfer(opts *bind.FilterOpts, sender []common.Address, receiver []common.Address) (*MerkleTreeMessageBusValueTransferIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "ValueTransfer", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusValueTransferIterator{contract: _MerkleTreeMessageBus.contract, event: "ValueTransfer", logs: logs, sub: sub}, nil
}

// WatchValueTransfer is a free log subscription operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchValueTransfer(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusValueTransfer, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "ValueTransfer", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusValueTransfer)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "ValueTransfer", log); err != nil {
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

// ParseValueTransfer is a log parse operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseValueTransfer(log types.Log) (*MerkleTreeMessageBusValueTransfer, error) {
	event := new(MerkleTreeMessageBusValueTransfer)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "ValueTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
