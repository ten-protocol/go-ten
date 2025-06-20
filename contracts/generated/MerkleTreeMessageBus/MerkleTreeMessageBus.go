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
	Nonce            uint64
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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STATE_ROOT_MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"activationTime\",\"type\":\"uint256\"}],\"name\":\"addStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"}],\"name\":\"addStateRootManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"withdrawalManager\",\"type\":\"address\"}],\"name\":\"addWithdrawalManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"name\":\"disableStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublishFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"withdrawalManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_fees\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"}],\"name\":\"removeStateRootManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"withdrawalManager\",\"type\":\"address\"}],\"name\":\"removeWithdrawalManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyMessageInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.ValueTransferMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyValueTransferInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506015601f565b601b601f565b60cf565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff1615606e5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460cc5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b6128e1806100dc5f395ff3fe6080604052600436106101b5575f3560e01c806391643fdd116100eb578063c0c53b8b11610089578063e30c397811610063578063e30c39781461055f578063ea4cf97014610573578063f2fde38b14610592578063fb894029146105b1576101b5565b8063c0c53b8b14610502578063ce0d7db314610521578063d547741f14610540576101b5565b8063ad7805e8116100c5578063ad7805e814610472578063b201246f146104a5578063b6aed0cb146104c4578063bb9cdd2c146104e3576101b5565b806391643fdd146103dd57806391d14854146103fc578063a217fddf1461045f576101b5565b80632ce8603d1161015857806336568abe1161013257806336568abe14610375578063715018a61461039457806379ba5097146103a85780638da5cb5b146103bc576101b5565b80632ce8603d146103185780632e1a0b8e146103375780632f2ff15d14610356576101b5565b80630d3fd67c116101945780630d3fd67c1461026b5780630fe9188e1461028b5780631050afdd146102ac578063248a9ca3146102cb576101b5565b8062a1b815146101e257806301ffc9a71461020c57806302b4df1914610238575b3480156101c0575f5ffd5b5060405162461bcd60e51b81526004016101d9906116ce565b60405180910390fd5b3480156101ed575f5ffd5b506101f66105d0565b60405161020391906116e6565b60405180910390f35b348015610217575f5ffd5b5061022b61022636600461172d565b610659565b6040516102039190611752565b348015610243575f5ffd5b506101f67f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19481565b61027e6102793660046117f4565b6106f1565b604051610203919061187c565b348015610296575f5ffd5b506102aa6102a536600461189b565b6107f9565b005b3480156102b7575f5ffd5b506102aa6102c63660046118dc565b610860565b3480156102d6575f5ffd5b506101f66102e536600461189b565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b348015610323575f5ffd5b506102aa6103323660046118dc565b610898565b348015610342575f5ffd5b506101f6610351366004611912565b6108cc565b348015610361575f5ffd5b506102aa61037036600461194a565b610929565b348015610380575f5ffd5b506102aa61038f36600461194a565b610972565b34801561039f575f5ffd5b506102aa6109c3565b3480156103b3575f5ffd5b506102aa6109e3565b3480156103c7575f5ffd5b506103d0610a22565b6040516102039190611989565b3480156103e8575f5ffd5b5061022b6103f7366004611912565b610a56565b348015610407575f5ffd5b5061022b61041636600461194a565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b34801561046a575f5ffd5b506101f65f81565b34801561047d575f5ffd5b506101f67fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9881565b3480156104b0575f5ffd5b506102aa6104bf3660046119f1565b610aa6565b3480156104cf575f5ffd5b506102aa6104de366004611a58565b610ba4565b3480156104ee575f5ffd5b506102aa6104fd3660046118dc565b610c2b565b34801561050d575f5ffd5b506102aa61051c366004611a85565b610c5f565b34801561052c575f5ffd5b506102aa61053b366004611acb565b610e32565b34801561054b575f5ffd5b506102aa61055a36600461194a565b610f79565b34801561056a575f5ffd5b506103d0610fbc565b34801561057e575f5ffd5b506102aa61058d366004611b43565b610fe4565b34801561059d575f5ffd5b506102aa6105ac3660046118dc565b611129565b3480156105bc575f5ffd5b506102aa6105cb3660046118dc565b6111bb565b600354604080517f1a90a21900000000000000000000000000000000000000000000000000000000815290515f926001600160a01b031691631a90a2199160048083019260209291908290030181865afa158015610630573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106549190611b97565b905090565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806106eb57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6003545f906001600160a01b0316156107a2575f61070d6105d0565b90508034101561072f5760405162461bcd60e51b81526004016101d990611c0e565b6003546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610779576040519150601f19603f3d011682016040523d82523d5f602084013e61077e565b606091505b505090508061079f5760405162461bcd60e51b81526004016101d990611c76565b50505b6107ab336111ef565b90507fd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061338288888888886040516107e89796959493929190611cc6565b60405180910390a195945050505050565b7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e1946108238161124c565b5f82815260046020526040812054900361084f5760405162461bcd60e51b81526004016101d990611d5b565b505f90815260046020526040812055565b5f61086a8161124c565b6108947f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19483610f79565b5050565b5f6108a28161124c565b6108947fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9883610f79565b5f5f826040516020016108df9190611ea5565b60408051601f1981840301815291815281516020928301205f81815292839052912054909150806109225760405162461bcd60e51b81526004016101d990611ef2565b9392505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546109628161124c565b61096c8383611256565b50505050565b6001600160a01b03811633146109b4576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109be8282611322565b505050565b6109cb6113c6565b60405162461bcd60e51b81526004016101d990611f5a565b33806109ed610fbc565b6001600160a01b031614610a16578060405163118cdaa760e01b81526004016101d99190611989565b610a1f816113fa565b50565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b5f5f82604051602001610a699190611ea5565b60408051601f1981840301815291815281516020928301205f818152928390529120549091508015801590610a9e5750428111155b949350505050565b5f818152600460205260408120549003610ad25760405162461bcd60e51b81526004016101d990611fc2565b5f81815260046020526040902054421015610aff5760405162461bcd60e51b81526004016101d99061200e565b5f84604051602001610b119190612092565b60405160208183030381529060405280519060200120604051602001610b3791906120d2565b604051602081830303815290604052805190602001209050610b8184848484604051602001610b6691906120fa565b6040516020818303038152906040528051906020012061143f565b610b9d5760405162461bcd60e51b81526004016101d990612164565b5050505050565b7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e194610bce8161124c565b5f8381526004602052604090205415610bf95760405162461bcd60e51b81526004016101d9906121cc565b42821015610c195760405162461bcd60e51b81526004016101d990612234565b505f9182526004602052604090912055565b5f610c358161124c565b6109be7fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9883611256565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f81158015610ca95750825b90505f8267ffffffffffffffff166001148015610cc55750303b155b905081158015610cd3575080155b15610d0a576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610d3e57845468ff00000000000000001916680100000000000000001785555b610d4788611456565b610d4f61146b565b610d595f89611256565b50610d847f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19489611256565b50610daf7fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9888611256565b506003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0388161790558315610e2857845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610e1f90600190612260565b60405180910390a15b5050505050505050565b5f818152600460205260408120549003610e5e5760405162461bcd60e51b81526004016101d990611fc2565b5f81815260046020526040902054421015610e8b5760405162461bcd60e51b81526004016101d99061200e565b5f610e9960208601866118dc565b610ea9604087016020880161226e565b610eb9606088016040890161226e565b610ec96080890160608a0161228b565b610ed660808a018a6122a8565b610ee660c08c0160a08d016122f9565b604051602001610efc9796959493929190611cc6565b6040516020818303038152906040528051906020012090505f81604051602001610f269190612348565b604051602081830303815290604052805190602001209050610f5585858584604051602001610b6691906120fa565b610f715760405162461bcd60e51b81526004016101d9906123b0565b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610fb28161124c565b61096c8383611322565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610a46565b5f610ff06001306123d4565b9050610ffa610a22565b6001600160a01b0316336001600160a01b031614806110215750336001600160a01b038216145b61103d5760405162461bcd60e51b81526004016101d990612429565b5f6110488342612439565b90505f8460405160200161105c9190611ea5565b60408051601f1981840301815291815281516020928301205f818152928390529120549091501561109f5760405162461bcd60e51b81526004016101d9906124a4565b5f818152602081815260408220849055600191906110bf908801886118dc565b6001600160a01b0316815260208101919091526040015f908120906110ea608088016060890161228b565b63ffffffff1681526020808201929092526040015f9081208054600181018255908252919020869160040201611120828261284e565b50505050505050565b6111316113c6565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081178255611182610a22565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b5f6111c58161124c565b6108947f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19483610929565b6001600160a01b0381165f908152600260205260408120805467ffffffffffffffff1691600191906112218385612858565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b610a1f8133611473565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff16611319575f848152602082815260408083206001600160a01b03871684529091529020805460ff191660011790556112cf3390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019150506106eb565b5f9150506106eb565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff1615611319575f848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a460019150506106eb565b336113cf610a22565b6001600160a01b0316146113f8573360405163118cdaa760e01b81526004016101d99190611989565b565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19168155610894826114f1565b5f8261144c86868561156e565b1495945050505050565b61145e6115af565b61146781611616565b610a1f5b6113f86115af565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff166108945780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016101d992919061287c565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f81815b848110156115a65761159c8287878481811061159057611590612897565b90506020020135611627565b9150600101611572565b50949350505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166113f8576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61161e6115af565b610a1f81611650565b5f818310611641575f828152602084905260409020610922565b505f9182526020526040902090565b6116586115af565b6001600160a01b038116610a16575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016101d99190611989565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b602080825281016106eb8161169a565b805b82525050565b602081016106eb82846116de565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b8114610a1f575f5ffd5b80356106eb816116f4565b5f60208284031215611740576117405f5ffd5b6109228383611722565b8015156116e0565b602081016106eb828461174a565b67ffffffffffffffff8116611718565b80356106eb81611760565b63ffffffff8116611718565b80356106eb8161177b565b5f5f83601f8401126117a5576117a55f5ffd5b50813567ffffffffffffffff8111156117bf576117bf5f5ffd5b6020830191508360018202830111156117d9576117d95f5ffd5b9250929050565b60ff8116611718565b80356106eb816117e0565b5f5f5f5f5f6080868803121561180b5761180b5f5ffd5b6118158787611770565b94506118248760208801611787565b9350604086013567ffffffffffffffff811115611842576118425f5ffd5b61184e88828901611792565b935093505061186087606088016117e9565b90509295509295909350565b67ffffffffffffffff81166116e0565b602081016106eb828461186c565b80611718565b80356106eb8161188a565b5f602082840312156118ae576118ae5f5ffd5b6109228383611890565b5f6001600160a01b0382166106eb565b611718816118b8565b80356106eb816118c8565b5f602082840312156118ef576118ef5f5ffd5b61092283836118d1565b5f60c0828403121561190c5761190c5f5ffd5b50919050565b5f60208284031215611925576119255f5ffd5b813567ffffffffffffffff81111561193e5761193e5f5ffd5b610a9e848285016118f9565b5f5f6040838503121561195e5761195e5f5ffd5b6119688484611890565b915061197784602085016118d1565b90509250929050565b6116e0816118b8565b602081016106eb8284611980565b5f6080828403121561190c5761190c5f5ffd5b5f5f83601f8401126119bd576119bd5f5ffd5b50813567ffffffffffffffff8111156119d7576119d75f5ffd5b6020830191508360208202830111156117d9576117d95f5ffd5b5f5f5f5f60c08587031215611a0757611a075f5ffd5b611a118686611997565b9350608085013567ffffffffffffffff811115611a2f57611a2f5f5ffd5b611a3b878288016119aa565b9350935050611a4d8660a08701611890565b905092959194509250565b5f5f60408385031215611a6c57611a6c5f5ffd5b611a768484611890565b91506119778460208501611890565b5f5f5f60608486031215611a9a57611a9a5f5ffd5b611aa485856118d1565b9250611ab385602086016118d1565b9150611ac285604086016118d1565b90509250925092565b5f5f5f5f60608587031215611ae157611ae15f5ffd5b843567ffffffffffffffff811115611afa57611afa5f5ffd5b611b06878288016118f9565b945050602085013567ffffffffffffffff811115611b2557611b255f5ffd5b611b31878288016119aa565b9350935050611a4d8660408701611890565b5f5f60408385031215611b5757611b575f5ffd5b823567ffffffffffffffff811115611b7057611b705f5ffd5b611b7c858286016118f9565b9250506119778460208501611890565b80516106eb8161188a565b5f60208284031215611baa57611baa5f5ffd5b6109228383611b8c565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f7373616765000000000000000000000000000000000000000000000000000000602082015290505b60400190565b602080825281016106eb81611bb4565b60248152602081017f4661696c656420746f2073656e64206665657320746f206665657320636f6e7481527f726163740000000000000000000000000000000000000000000000000000000060208201529050611c08565b602080825281016106eb81611c1e565b63ffffffff81166116e0565b82818337505f910152565b818352602083019250611cb1828483611c92565b50601f01601f19160190565b60ff81166116e0565b60c08101611cd4828a611980565b611ce1602083018961186c565b611cee604083018861186c565b611cfb6060830187611c86565b8181036080830152611d0e818587611c9d565b9050611d1d60a0830184611cbd565b98975050505050505050565b601a8152602081017f537461746520726f6f7420646f6573206e6f742065786973742e000000000000815290506116c8565b602080825281016106eb81611d29565b505f6106eb60208301836118d1565b505f6106eb6020830183611770565b505f6106eb6020830183611787565b5f808335601e1936859003018112611db157611db15f5ffd5b830160208101925035905067ffffffffffffffff811115611dd357611dd35f5ffd5b368190038213156117d9576117d95f5ffd5b505f6106eb60208301836117e9565b5f60c08301611e038380611d6b565b611e0d8582611980565b50611e1b6020840184611d7a565b611e28602086018261186c565b50611e366040840184611d7a565b611e43604086018261186c565b50611e516060840184611d89565b611e5e6060860182611c86565b50611e6c6080840184611d98565b8583036080870152611e7f838284611c9d565b92505050611e9060a0840184611de5565b611e9d60a0860182611cbd565b509392505050565b602080825281016109228184611df4565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d69747465648152601760f91b60208201529050611c08565b602080825281016106eb81611eb6565b60348152602081017f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f81527f742072656e6f756e6365206f776e65727368697000000000000000000000000060208201529050611c08565b602080825281016106eb81611f02565b602a8152602081017f526f6f74206973206e6f74207075626c6973686564206f6e2074686973206d6581527f7373616765206275732e0000000000000000000000000000000000000000000060208201529050611c08565b602080825281016106eb81611f6a565b60218152602081017f526f6f74206973206e6f7420636f6e736964657265642066696e616c207965748152601760f91b60208201529050611c08565b602080825281016106eb81611fd2565b505f6106eb6020830183611890565b6120378180611d6b565b6120418382611980565b5061204f6020820182611d6b565b61205c6020840182611980565b5061206a604082018261201e565b61207760408401826116de565b506120856060820182611d7a565b6109be606084018261186c565b608081016106eb828461202d565b60018152602081017f7600000000000000000000000000000000000000000000000000000000000000815290506116c8565b604080825281016120e2816120a0565b90506106eb60208301846116de565b806116e0565b90565b61210481836120f1565b602001919050565b60338152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722076616c7581527f65207472616e73666572206d6573736167652e0000000000000000000000000060208201529050611c08565b602080825281016106eb8161210c565b60258152602081017f526f6f7420616c726561647920616464656420746f20746865206d657373616781527f652062757300000000000000000000000000000000000000000000000000000060208201529050611c08565b602080825281016106eb81612174565b60258152602081017f41637469766174696f6e2074696d65206d75737420626520696e20746865206681527f757475726500000000000000000000000000000000000000000000000000000060208201529050611c08565b602080825281016106eb816121dc565b5f6106eb825b67ffffffffffffffff1690565b6116e081612244565b602081016106eb8284612257565b5f60208284031215612281576122815f5ffd5b6109228383611770565b5f6020828403121561229e5761229e5f5ffd5b6109228383611787565b5f808335601e19368590030181126122c1576122c15f5ffd5b8301915050803567ffffffffffffffff8111156122df576122df5f5ffd5b6020820191506001810236038213156117d9576117d95f5ffd5b5f6020828403121561230c5761230c5f5ffd5b61092283836117e9565b60018152602081017f6d00000000000000000000000000000000000000000000000000000000000000815290506116c8565b604080825281016120e281612316565b60308152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722063726f7381527f7320636861696e206d6573736167652e0000000000000000000000000000000060208201529050611c08565b602080825281016106eb81612358565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b039182169190811690828203908111156106eb576106eb6123c0565b60118152602081017f4e6f74206f776e6572206f722073656c66000000000000000000000000000000815290506116c8565b602080825281016106eb816123f7565b808201808211156106eb576106eb6123c0565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f210000000000000000000000000000000000000000000000000000000000000060208201529050611c08565b602080825281016106eb8161244c565b5f81356106eb816118c8565b5f6001600160a01b03835b81169019929092169190911792915050565b5f6106eb826118b8565b5f6106eb826124dd565b6124fa826124e7565b6125058183546124c0565b8255505050565b5f81356106eb81611760565b5f7bffffffffffffffff00000000000000000000000000000000000000006124cb8460a01b90565b5f6106eb67ffffffffffffffff831661224a565b61255d82612540565b612505818354612518565b5f67ffffffffffffffff836124cb565b61258182612540565b612505818354612568565b5f81356106eb8161177b565b5f6bffffffff00000000000000006124cb8460401b90565b5f63ffffffff82166106eb565b6125c6826125b0565b612505818354612598565b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061260d57607f821691505b60208210810361190c5761190c6125e5565b5f6106eb6120f78381565b6126338361261f565b81545f1960089490940293841b1916921b91909117905550565b5f6109be81848461262a565b818110156108945761266b5f8261264d565b600101612659565b601f8211156109be575f818152602090206020601f850104810160208510156126995750805b610b9d6020601f860104830182612659565b8267ffffffffffffffff8111156126c4576126c46125d1565b6126ce82546125f9565b6126d9828285612673565b505f601f82116001811461270b575f83156126f45750848201355b5f19600885021c1981166002850217855550610f71565b5f84815260208120601f198516915b8281101561273a578785013582556020948501946001909201910161271a565b5084821015612756575f196008601f8716021c19878501351681555b5050505060020260010190555050565b6109be8383836126ab565b5f81356106eb816117e0565b5f60ff82166106eb565b6127908261277d565b815460ff191660ff821617612505565b8082806127ac816124b4565b90506127b881846124f1565b505060208301806127c88261250c565b90506127d48184612554565b5050506001810160408301806127e98261250c565b90506127f58184612578565b505060608301806128058261258c565b905061281181846125bd565b50505061282160808301836122a8565b61282f818360028601612766565b505060a082018061283f82612771565b905061096c8160038501612787565b61089482826127a0565b67ffffffffffffffff9182169190811690828201908111156106eb576106eb6123c0565b6040810161288a8285611980565b61092260208301846116de565b634e487b7160e01b5f52603260045260245ffdfea2646970667358221220ac0bcb1b7a764486ac03b59a536962b81a6a57e6efbaeb1045312e87a6da0ba464736f6c634300081c0033",
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

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.DEFAULTADMINROLE(&_MerkleTreeMessageBus.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.DEFAULTADMINROLE(&_MerkleTreeMessageBus.CallOpts)
}

// STATEROOTMANAGERROLE is a free data retrieval call binding the contract method 0x02b4df19.
//
// Solidity: function STATE_ROOT_MANAGER_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) STATEROOTMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "STATE_ROOT_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// STATEROOTMANAGERROLE is a free data retrieval call binding the contract method 0x02b4df19.
//
// Solidity: function STATE_ROOT_MANAGER_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) STATEROOTMANAGERROLE() ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.STATEROOTMANAGERROLE(&_MerkleTreeMessageBus.CallOpts)
}

// STATEROOTMANAGERROLE is a free data retrieval call binding the contract method 0x02b4df19.
//
// Solidity: function STATE_ROOT_MANAGER_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) STATEROOTMANAGERROLE() ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.STATEROOTMANAGERROLE(&_MerkleTreeMessageBus.CallOpts)
}

// WITHDRAWALMANAGERROLE is a free data retrieval call binding the contract method 0xad7805e8.
//
// Solidity: function WITHDRAWAL_MANAGER_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) WITHDRAWALMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "WITHDRAWAL_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WITHDRAWALMANAGERROLE is a free data retrieval call binding the contract method 0xad7805e8.
//
// Solidity: function WITHDRAWAL_MANAGER_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) WITHDRAWALMANAGERROLE() ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.WITHDRAWALMANAGERROLE(&_MerkleTreeMessageBus.CallOpts)
}

// WITHDRAWALMANAGERROLE is a free data retrieval call binding the contract method 0xad7805e8.
//
// Solidity: function WITHDRAWAL_MANAGER_ROLE() view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) WITHDRAWALMANAGERROLE() ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.WITHDRAWALMANAGERROLE(&_MerkleTreeMessageBus.CallOpts)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetMessageTimeOfFinality(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetMessageTimeOfFinality(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) GetPublishFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "getPublishFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) GetPublishFee() (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetPublishFee(&_MerkleTreeMessageBus.CallOpts)
}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) GetPublishFee() (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetPublishFee(&_MerkleTreeMessageBus.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.GetRoleAdmin(&_MerkleTreeMessageBus.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _MerkleTreeMessageBus.Contract.GetRoleAdmin(&_MerkleTreeMessageBus.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _MerkleTreeMessageBus.Contract.HasRole(&_MerkleTreeMessageBus.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _MerkleTreeMessageBus.Contract.HasRole(&_MerkleTreeMessageBus.CallOpts, role, account)
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

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) PendingOwner() (common.Address, error) {
	return _MerkleTreeMessageBus.Contract.PendingOwner(&_MerkleTreeMessageBus.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) PendingOwner() (common.Address, error) {
	return _MerkleTreeMessageBus.Contract.PendingOwner(&_MerkleTreeMessageBus.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RenounceOwnership() error {
	return _MerkleTreeMessageBus.Contract.RenounceOwnership(&_MerkleTreeMessageBus.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) RenounceOwnership() error {
	return _MerkleTreeMessageBus.Contract.RenounceOwnership(&_MerkleTreeMessageBus.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MerkleTreeMessageBus.Contract.SupportsInterface(&_MerkleTreeMessageBus.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MerkleTreeMessageBus.Contract.SupportsInterface(&_MerkleTreeMessageBus.CallOpts, interfaceId)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MerkleTreeMessageBus.Contract.VerifyMessageFinalized(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MerkleTreeMessageBus.Contract.VerifyMessageFinalized(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xce0d7db3.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyMessageInclusion(opts *bind.CallOpts, message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyMessageInclusion", message, proof, root)

	if err != nil {
		return err
	}

	return err

}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xce0d7db3.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyMessageInclusion(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyMessageInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xce0d7db3.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
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

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) AcceptOwnership() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AcceptOwnership(&_MerkleTreeMessageBus.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AcceptOwnership(&_MerkleTreeMessageBus.TransactOpts)
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

// AddStateRootManager is a paid mutator transaction binding the contract method 0xfb894029.
//
// Solidity: function addStateRootManager(address manager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) AddStateRootManager(opts *bind.TransactOpts, manager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "addStateRootManager", manager)
}

// AddStateRootManager is a paid mutator transaction binding the contract method 0xfb894029.
//
// Solidity: function addStateRootManager(address manager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) AddStateRootManager(manager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddStateRootManager(&_MerkleTreeMessageBus.TransactOpts, manager)
}

// AddStateRootManager is a paid mutator transaction binding the contract method 0xfb894029.
//
// Solidity: function addStateRootManager(address manager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) AddStateRootManager(manager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddStateRootManager(&_MerkleTreeMessageBus.TransactOpts, manager)
}

// AddWithdrawalManager is a paid mutator transaction binding the contract method 0xbb9cdd2c.
//
// Solidity: function addWithdrawalManager(address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) AddWithdrawalManager(opts *bind.TransactOpts, withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "addWithdrawalManager", withdrawalManager)
}

// AddWithdrawalManager is a paid mutator transaction binding the contract method 0xbb9cdd2c.
//
// Solidity: function addWithdrawalManager(address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) AddWithdrawalManager(withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddWithdrawalManager(&_MerkleTreeMessageBus.TransactOpts, withdrawalManager)
}

// AddWithdrawalManager is a paid mutator transaction binding the contract method 0xbb9cdd2c.
//
// Solidity: function addWithdrawalManager(address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) AddWithdrawalManager(withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddWithdrawalManager(&_MerkleTreeMessageBus.TransactOpts, withdrawalManager)
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

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.GrantRole(&_MerkleTreeMessageBus.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.GrantRole(&_MerkleTreeMessageBus.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address withdrawalManager, address _fees) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, withdrawalManager common.Address, _fees common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "initialize", initialOwner, withdrawalManager, _fees)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address withdrawalManager, address _fees) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Initialize(initialOwner common.Address, withdrawalManager common.Address, _fees common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Initialize(&_MerkleTreeMessageBus.TransactOpts, initialOwner, withdrawalManager, _fees)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address initialOwner, address withdrawalManager, address _fees) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Initialize(initialOwner common.Address, withdrawalManager common.Address, _fees common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Initialize(&_MerkleTreeMessageBus.TransactOpts, initialOwner, withdrawalManager, _fees)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) PublishMessage(nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.PublishMessage(&_MerkleTreeMessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) PublishMessage(nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.PublishMessage(&_MerkleTreeMessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// RemoveStateRootManager is a paid mutator transaction binding the contract method 0x1050afdd.
//
// Solidity: function removeStateRootManager(address manager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RemoveStateRootManager(opts *bind.TransactOpts, manager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "removeStateRootManager", manager)
}

// RemoveStateRootManager is a paid mutator transaction binding the contract method 0x1050afdd.
//
// Solidity: function removeStateRootManager(address manager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RemoveStateRootManager(manager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RemoveStateRootManager(&_MerkleTreeMessageBus.TransactOpts, manager)
}

// RemoveStateRootManager is a paid mutator transaction binding the contract method 0x1050afdd.
//
// Solidity: function removeStateRootManager(address manager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RemoveStateRootManager(manager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RemoveStateRootManager(&_MerkleTreeMessageBus.TransactOpts, manager)
}

// RemoveWithdrawalManager is a paid mutator transaction binding the contract method 0x2ce8603d.
//
// Solidity: function removeWithdrawalManager(address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RemoveWithdrawalManager(opts *bind.TransactOpts, withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "removeWithdrawalManager", withdrawalManager)
}

// RemoveWithdrawalManager is a paid mutator transaction binding the contract method 0x2ce8603d.
//
// Solidity: function removeWithdrawalManager(address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RemoveWithdrawalManager(withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RemoveWithdrawalManager(&_MerkleTreeMessageBus.TransactOpts, withdrawalManager)
}

// RemoveWithdrawalManager is a paid mutator transaction binding the contract method 0x2ce8603d.
//
// Solidity: function removeWithdrawalManager(address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RemoveWithdrawalManager(withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RemoveWithdrawalManager(&_MerkleTreeMessageBus.TransactOpts, withdrawalManager)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RenounceRole(&_MerkleTreeMessageBus.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RenounceRole(&_MerkleTreeMessageBus.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RevokeRole(&_MerkleTreeMessageBus.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RevokeRole(&_MerkleTreeMessageBus.TransactOpts, role, account)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) StoreCrossChainMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "storeCrossChainMessage", crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.StoreCrossChainMessage(&_MerkleTreeMessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
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
// Solidity: fallback() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Fallback(&_MerkleTreeMessageBus.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Fallback(&_MerkleTreeMessageBus.TransactOpts, calldata)
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
	Nonce            uint64
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogMessagePublished is a free log retrieval operation binding the contract event 0xd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts) (*MerkleTreeMessageBusLogMessagePublishedIterator, error) {

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusLogMessagePublishedIterator{contract: _MerkleTreeMessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
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

// ParseLogMessagePublished is a log parse operation binding the contract event 0xd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MerkleTreeMessageBusLogMessagePublished, error) {
	event := new(MerkleTreeMessageBusLogMessagePublished)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusOwnershipTransferStartedIterator struct {
	Event *MerkleTreeMessageBusOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusOwnershipTransferStarted)
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
		it.Event = new(MerkleTreeMessageBusOwnershipTransferStarted)
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
func (it *MerkleTreeMessageBusOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MerkleTreeMessageBusOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusOwnershipTransferStartedIterator{contract: _MerkleTreeMessageBus.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusOwnershipTransferStarted)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseOwnershipTransferStarted(log types.Log) (*MerkleTreeMessageBusOwnershipTransferStarted, error) {
	event := new(MerkleTreeMessageBusOwnershipTransferStarted)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// MerkleTreeMessageBusRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusRoleAdminChangedIterator struct {
	Event *MerkleTreeMessageBusRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusRoleAdminChanged)
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
		it.Event = new(MerkleTreeMessageBusRoleAdminChanged)
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
func (it *MerkleTreeMessageBusRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusRoleAdminChanged represents a RoleAdminChanged event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*MerkleTreeMessageBusRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusRoleAdminChangedIterator{contract: _MerkleTreeMessageBus.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusRoleAdminChanged)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseRoleAdminChanged(log types.Log) (*MerkleTreeMessageBusRoleAdminChanged, error) {
	event := new(MerkleTreeMessageBusRoleAdminChanged)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusRoleGrantedIterator struct {
	Event *MerkleTreeMessageBusRoleGranted // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusRoleGranted)
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
		it.Event = new(MerkleTreeMessageBusRoleGranted)
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
func (it *MerkleTreeMessageBusRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusRoleGranted represents a RoleGranted event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MerkleTreeMessageBusRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusRoleGrantedIterator{contract: _MerkleTreeMessageBus.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusRoleGranted)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseRoleGranted(log types.Log) (*MerkleTreeMessageBusRoleGranted, error) {
	event := new(MerkleTreeMessageBusRoleGranted)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusRoleRevokedIterator struct {
	Event *MerkleTreeMessageBusRoleRevoked // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusRoleRevoked)
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
		it.Event = new(MerkleTreeMessageBusRoleRevoked)
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
func (it *MerkleTreeMessageBusRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusRoleRevoked represents a RoleRevoked event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MerkleTreeMessageBusRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusRoleRevokedIterator{contract: _MerkleTreeMessageBus.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusRoleRevoked)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseRoleRevoked(log types.Log) (*MerkleTreeMessageBusRoleRevoked, error) {
	event := new(MerkleTreeMessageBusRoleRevoked)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
