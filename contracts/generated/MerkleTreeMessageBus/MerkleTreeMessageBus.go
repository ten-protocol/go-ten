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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STATE_ROOT_MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"activationTime\",\"type\":\"uint256\"}],\"name\":\"addStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"}],\"name\":\"addStateRootManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"name\":\"disableStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublishFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"withdrawalManager\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"}],\"name\":\"removeStateRootManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"retrieveAllFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyMessageInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.ValueTransferMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyValueTransferInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b612649806100975f395ff3fe608060405260043610610194575f3560e01c8063715018a6116100e0578063b1454caa11610089578063d547741f11610063578063d547741f1461050f578063e138a8d21461052e578063f2fde38b1461054d578063fb8940291461056c57610194565b8063b1454caa146104b1578063b201246f146104d1578063b6aed0cb146104f057610194565b80639730886d116100ba5780639730886d1461044c578063a217fddf1461046b578063ad7805e81461047e57610194565b8063715018a6146103915780638da5cb5b146103a557806391d14854146103e957610194565b8063248a9ca31161014257806336568abe1161011c57806336568abe1461033457806336d2da9014610353578063485cc9551461037257610194565b8063248a9ca3146102a95780632f2ff15d146102f657806333a88c721461031557610194565b80630fcfbd11116101735780630fcfbd111461024a5780630fe9188e146102695780631050afdd1461028a57610194565b8062a1b815146101c157806301ffc9a7146101eb57806302b4df1914610217575b34801561019f575f5ffd5b5060405162461bcd60e51b81526004016101b89061150a565b60405180910390fd5b3480156101cc575f5ffd5b506101d561058b565b6040516101e29190611522565b60405180910390f35b3480156101f6575f5ffd5b5061020a610205366004611551565b610614565b6040516101e29190611576565b348015610222575f5ffd5b506101d57f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19481565b348015610255575f5ffd5b506101d561026436600461159d565b61067c565b348015610274575f5ffd5b506102886102833660046115e6565b6106d9565b005b348015610295575f5ffd5b506102886102a4366004611627565b610740565b3480156102b4575f5ffd5b506101d56102c33660046115e6565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b348015610301575f5ffd5b50610288610310366004611644565b610778565b348015610320575f5ffd5b5061020a61032f36600461159d565b6107c1565b34801561033f575f5ffd5b5061028861034e366004611644565b610811565b34801561035e575f5ffd5b5061028861036d366004611627565b610862565b34801561037d575f5ffd5b5061028861038c36600461167a565b6108d9565b34801561039c575f5ffd5b50610288610a83565b3480156103b0575f5ffd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516101e291906116a1565b3480156103f4575f5ffd5b5061020a610403366004611644565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b348015610457575f5ffd5b506102886104663660046116af565b610a96565b348015610476575f5ffd5b506101d55f81565b348015610489575f5ffd5b506101d57fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9881565b6104c46104bf366004611771565b610bfb565b6040516101e291906117f9565b3480156104dc575f5ffd5b506102886104eb366004611861565b610d03565b3480156104fb575f5ffd5b5061028861050a3660046118c8565b610e01565b34801561051a575f5ffd5b50610288610529366004611644565b610e68565b348015610539575f5ffd5b506102886105483660046118f5565b610eab565b348015610558575f5ffd5b50610288610567366004611627565b610ff2565b348015610577575f5ffd5b50610288610586366004611627565b611048565b600354604080517f1a90a21900000000000000000000000000000000000000000000000000000000815290515f926001600160a01b031691631a90a2199160048083019260209291908290030181865afa1580156105eb573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061060f9190611978565b905090565b5f6001600160e01b031982167f7965db0b00000000000000000000000000000000000000000000000000000000148061067657507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b5f5f8260405160200161068f9190611b2a565b60408051601f1981840301815291815281516020928301205f81815292839052912054909150806106d25760405162461bcd60e51b81526004016101b890611b79565b9392505050565b7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e1946107038161107c565b5f82815260046020526040812054900361072f5760405162461bcd60e51b81526004016101b890611bbb565b505f90815260046020526040812055565b5f61074a8161107c565b6107747f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19483610e68565b5050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546107b18161107c565b6107bb8383611086565b50505050565b5f5f826040516020016107d49190611b2a565b60408051601f1981840301815291815281516020928301205f8181529283905291205490915080158015906108095750428111155b949350505050565b6001600160a01b0381163314610853576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61085d8282611152565b505050565b61086a6111f6565b5f816001600160a01b0316476040515f6040518083038185875af1925050503d805f81146108b3576040519150601f19603f3d011682016040523d82523d5f602084013e6108b8565b606091505b50509050806107745760405162461bcd60e51b81526004016101b890611bfd565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156109235750825b90505f8267ffffffffffffffff16600114801561093f5750303b155b90508115801561094d575080155b15610984576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156109b857845468ff00000000000000001916680100000000000000001785555b6109c18761126a565b6109c961127b565b6109d35f88611086565b506109fe7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19488611086565b50610a297fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9887611086565b508315610a7a57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610a7190600190611c30565b60405180910390a15b50505050505050565b610a8b6111f6565b610a945f611283565b565b5f610aa2600130611c52565b9050610ad57f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b03161480610afc5750336001600160a01b038216145b610b185760405162461bcd60e51b81526004016101b890611ca7565b5f610b238342611cb7565b90505f84604051602001610b379190611b2a565b60408051601f1981840301815291815281516020928301205f8181529283905291205490915015610b7a5760405162461bcd60e51b81526004016101b890611d22565b5f81815260208181526040822084905560019190610b9a90880188611627565b6001600160a01b0316815260208101919091526040015f90812090610bc56080880160608901611d32565b63ffffffff1681526020808201929092526040015f9081208054600181018255908252919020869160040201610a7a8282612130565b6003545f906001600160a01b031615610cac575f610c1761058b565b905080341015610c395760405162461bcd60e51b81526004016101b890612192565b6003546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610c83576040519150601f19603f3d011682016040523d82523d5f602084013e610c88565b606091505b5050905080610ca95760405162461bcd60e51b81526004016101b8906121fa565b50505b610cb53361130b565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef77593733828888888888604051610cf2979695949392919061220a565b60405180910390a195945050505050565b5f818152600460205260408120549003610d2f5760405162461bcd60e51b81526004016101b8906122c5565b5f81815260046020526040902054421015610d5c5760405162461bcd60e51b81526004016101b890612311565b5f84604051602001610d6e9190612395565b60405160208183030381529060405280519060200120604051602001610d9491906123d5565b604051602081830303815290604052805190602001209050610dde84848484604051602001610dc391906123fa565b60405160208183030381529060405280519060200120611368565b610dfa5760405162461bcd60e51b81526004016101b890612464565b5050505050565b7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e194610e2b8161107c565b5f8381526004602052604090205415610e565760405162461bcd60e51b81526004016101b8906124cc565b505f9182526004602052604090912055565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610ea18161107c565b6107bb8383611152565b5f818152600460205260408120549003610ed75760405162461bcd60e51b81526004016101b8906122c5565b5f81815260046020526040902054421015610f045760405162461bcd60e51b81526004016101b890612311565b5f610f126020860186611627565b610f2260408701602088016124dc565b610f326060880160408901611d32565b610f426080890160608a01611d32565b610f4f60808a018a611e64565b610f5f60c08c0160a08d016124f9565b604051602001610f75979695949392919061220a565b6040516020818303038152906040528051906020012090505f81604051602001610f9f9190612548565b604051602081830303815290604052805190602001209050610fce85858584604051602001610dc391906123fa565b610fea5760405162461bcd60e51b81526004016101b8906125b0565b505050505050565b610ffa6111f6565b6001600160a01b03811661103c575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016101b891906116a1565b61104581611283565b50565b5f6110528161107c565b6107747f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19483610778565b611045813361137f565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff16611149575f848152602082815260408083206001600160a01b03871684529091529020805460ff191660011790556110ff3390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a46001915050610676565b5f915050610676565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff1615611149575f848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a46001915050610676565b336112287f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b031614610a9457336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016101b891906116a1565b6112726113fd565b61104581611464565b610a946113fd565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080547fffffffffffffffffffffffff000000000000000000000000000000000000000081166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b6001600160a01b0381165f908152600260205260408120805467ffffffffffffffff16916001919061133d83856125c0565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b5f8261137586868561146c565b1495945050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff166107745780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016101b89291906125e4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610a94576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ffa6113fd565b5f81815b848110156114a45761149a8287878481811061148e5761148e6125ff565b905060200201356114ad565b9150600101611470565b50949350505050565b5f8183106114c7575f8281526020849052604090206106d2565b505f9182526020526040902090565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b60208082528101610676816114d6565b805b82525050565b60208101610676828461151a565b6001600160e01b031981165b8114611045575f5ffd5b803561067681611530565b5f60208284031215611564576115645f5ffd5b6106d28383611546565b80151561151c565b60208101610676828461156e565b5f60c08284031215611597576115975f5ffd5b50919050565b5f602082840312156115b0576115b05f5ffd5b813567ffffffffffffffff8111156115c9576115c95f5ffd5b61080984828501611584565b8061153c565b8035610676816115d5565b5f602082840312156115f9576115f95f5ffd5b6106d283836115db565b5f6001600160a01b038216610676565b61153c81611603565b803561067681611613565b5f6020828403121561163a5761163a5f5ffd5b6106d2838361161c565b5f5f60408385031215611658576116585f5ffd5b61166284846115db565b9150611671846020850161161c565b90509250929050565b5f5f6040838503121561168e5761168e5f5ffd5b611662848461161c565b61151c81611603565b602081016106768284611698565b5f5f604083850312156116c3576116c35f5ffd5b823567ffffffffffffffff8111156116dc576116dc5f5ffd5b6116e885828601611584565b92505061167184602085016115db565b63ffffffff811661153c565b8035610676816116f8565b5f5f83601f840112611722576117225f5ffd5b50813567ffffffffffffffff81111561173c5761173c5f5ffd5b602083019150836001820283011115611756576117565f5ffd5b9250929050565b60ff811661153c565b80356106768161175d565b5f5f5f5f5f60808688031215611788576117885f5ffd5b6117928787611704565b94506117a18760208801611704565b9350604086013567ffffffffffffffff8111156117bf576117bf5f5ffd5b6117cb8882890161170f565b93509350506117dd8760608801611766565b90509295509295909350565b67ffffffffffffffff811661151c565b6020810161067682846117e9565b5f60808284031215611597576115975f5ffd5b5f5f83601f84011261182d5761182d5f5ffd5b50813567ffffffffffffffff811115611847576118475f5ffd5b602083019150836020820283011115611756576117565f5ffd5b5f5f5f5f60c08587031215611877576118775f5ffd5b6118818686611807565b9350608085013567ffffffffffffffff81111561189f5761189f5f5ffd5b6118ab8782880161181a565b93509350506118bd8660a087016115db565b905092959194509250565b5f5f604083850312156118dc576118dc5f5ffd5b6118e684846115db565b915061167184602085016115db565b5f5f5f5f6060858703121561190b5761190b5f5ffd5b843567ffffffffffffffff811115611924576119245f5ffd5b61193087828801611584565b945050602085013567ffffffffffffffff81111561194f5761194f5f5ffd5b61195b8782880161181a565b93509350506118bd86604087016115db565b8051610676816115d5565b5f6020828403121561198b5761198b5f5ffd5b6106d2838361196d565b505f610676602083018361161c565b67ffffffffffffffff811661153c565b8035610676816119a4565b505f61067660208301836119b4565b505f6106766020830183611704565b63ffffffff811661151c565b5f808335601e1936859003018112611a0257611a025f5ffd5b830160208101925035905067ffffffffffffffff811115611a2457611a245f5ffd5b36819003821315611756576117565f5ffd5b82818337505f910152565b818352602083019250611a55828483611a36565b50601f01601f19160190565b505f6106766020830183611766565b60ff811661151c565b5f60c08301611a888380611995565b611a928582611698565b50611aa060208401846119bf565b611aad60208601826117e9565b50611abb60408401846119ce565b611ac860408601826119dd565b50611ad660608401846119ce565b611ae360608601826119dd565b50611af160808401846119e9565b8583036080870152611b04838284611a41565b92505050611b1560a0840184611a61565b611b2260a0860182611a70565b509392505050565b602080825281016106d28184611a79565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d69747465648152601760f91b602082015290505b60400190565b6020808252810161067681611b3b565b601a8152602081017f537461746520726f6f7420646f6573206e6f742065786973742e00000000000081529050611504565b6020808252810161067681611b89565b60148152602081017f6661696c65642073656e64696e672076616c756500000000000000000000000081529050611504565b6020808252810161067681611bcb565b5f61067682611c1a565b90565b67ffffffffffffffff1690565b61151c81611c0d565b602081016106768284611c27565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b0391821691908116908282039081111561067657610676611c3e565b60118152602081017f4e6f74206f776e6572206f722073656c6600000000000000000000000000000081529050611504565b6020808252810161067681611c75565b8082018082111561067657610676611c3e565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f210000000000000000000000000000000000000000000000000000000000000060208201529050611b73565b6020808252810161067681611cca565b5f60208284031215611d4557611d455f5ffd5b6106d28383611704565b5f813561067681611613565b5f6001600160a01b03835b81169019929092169190911792915050565b5f61067682611603565b5f61067682611d78565b611d9582611d82565b611da0818354611d5b565b8255505050565b5f8135610676816119a4565b5f7bffffffffffffffff0000000000000000000000000000000000000000611d668460a01b90565b5f61067667ffffffffffffffff8316611c1a565b611df882611ddb565b611da0818354611db3565b5f8135610676816116f8565b5f6001600160e01b0319611d668460e01b90565b5f63ffffffff8216610676565b611e3982611e23565b611da0818354611e0f565b5f63ffffffff83611d66565b611e5982611e23565b611da0818354611e44565b5f808335601e1936859003018112611e7d57611e7d5f5ffd5b8301915050803567ffffffffffffffff811115611e9b57611e9b5f5ffd5b602082019150600181023603821315611756576117565f5ffd5b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b600281046001821680611ef157607f821691505b60208210810361159757611597611ec9565b5f610676611c178381565b611f1783611f03565b81545f1960089490940293841b1916921b91909117905550565b5f61085d818484611f0e565b8181101561077457611f4f5f82611f31565b600101611f3d565b601f82111561085d575f818152602090206020601f85010481016020851015611f7d5750805b610dfa6020601f860104830182611f3d565b8267ffffffffffffffff811115611fa857611fa8611eb5565b611fb28254611edd565b611fbd828285611f57565b505f601f821160018114611fef575f8315611fd85750848201355b5f19600885021c1981166002850217855550610fea565b5f84815260208120601f198516915b8281101561201e5787850135825560209485019460019092019101611ffe565b508482101561203a575f196008601f8716021c19878501351681555b5050505060020260010190555050565b61085d838383611f8f565b5f81356106768161175d565b5f60ff8216610676565b61207482612061565b815460ff191660ff821617611da0565b80828061209081611d4f565b905061209c8184611d8c565b505060208301806120ac82611da7565b90506120b88184611def565b505060408301806120c882611e03565b90506120d48184611e30565b50505060608201806120e582611e03565b90506120f48160018501611e50565b50506121036080830183611e64565b61211181836002860161204a565b505060a082018061212182612055565b90506107bb816003850161206b565b6107748282612084565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f737361676500000000000000000000000000000000000000000000000000000060208201529050611b73565b602080825281016106768161213a565b60248152602081017f4661696c656420746f2073656e64206665657320746f206665657320636f6e7481527f726163740000000000000000000000000000000000000000000000000000000060208201529050611b73565b60208082528101610676816121a2565b60c08101612218828a611698565b61222560208301896117e9565b61223260408301886119dd565b61223f60608301876119dd565b8181036080830152612252818587611a41565b905061226160a0830184611a70565b98975050505050505050565b602a8152602081017f526f6f74206973206e6f74207075626c6973686564206f6e2074686973206d6581527f7373616765206275732e0000000000000000000000000000000000000000000060208201529050611b73565b602080825281016106768161226d565b60218152602081017f526f6f74206973206e6f7420636f6e736964657265642066696e616c207965748152601760f91b60208201529050611b73565b60208082528101610676816122d5565b505f61067660208301836115db565b61233a8180611995565b6123448382611698565b506123526020820182611995565b61235f6020840182611698565b5061236d6040820182612321565b61237a604084018261151a565b5061238860608201826119bf565b61085d60608401826117e9565b608081016106768284612330565b60018152602081017f760000000000000000000000000000000000000000000000000000000000000081529050611504565b604080825281016123e5816123a3565b9050610676602083018461151a565b8061151c565b61240481836123f4565b602001919050565b60338152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722076616c7581527f65207472616e73666572206d6573736167652e0000000000000000000000000060208201529050611b73565b602080825281016106768161240c565b60258152602081017f526f6f7420616c726561647920616464656420746f20746865206d657373616781527f652062757300000000000000000000000000000000000000000000000000000060208201529050611b73565b6020808252810161067681612474565b5f602082840312156124ef576124ef5f5ffd5b6106d283836119b4565b5f6020828403121561250c5761250c5f5ffd5b6106d28383611766565b60018152602081017f6d0000000000000000000000000000000000000000000000000000000000000081529050611504565b604080825281016123e581612516565b60308152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722063726f7381527f7320636861696e206d6573736167652e0000000000000000000000000000000060208201529050611b73565b6020808252810161067681612558565b67ffffffffffffffff91821691908116908282019081111561067657610676611c3e565b604081016125f28285611698565b6106d2602083018461151a565b634e487b7160e01b5f52603260045260245ffdfea2646970667358221220336963fe788b9daf25603880fee3b2b36f869df4fa2914a67b951e91e692be3e64736f6c634300081c0033",
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

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address initialOwner, address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "initialize", initialOwner, withdrawalManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address initialOwner, address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Initialize(initialOwner common.Address, withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Initialize(&_MerkleTreeMessageBus.TransactOpts, initialOwner, withdrawalManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address initialOwner, address withdrawalManager) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Initialize(initialOwner common.Address, withdrawalManager common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Initialize(&_MerkleTreeMessageBus.TransactOpts, initialOwner, withdrawalManager)
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
