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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"name\":\"ValueTransfer\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STATE_ROOT_MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAWAL_MANAGER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"activationTime\",\"type\":\"uint256\"}],\"name\":\"addStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"}],\"name\":\"addStateRootManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"name\":\"disableStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublishFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"withdrawalManager\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"notifyDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"}],\"name\":\"removeStateRootManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"retrieveAllFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyMessageInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.ValueTransferMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyValueTransferInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b612759806100975f395ff3fe60806040526004361061019f575f3560e01c80638da5cb5b116100e0578063b1454caa11610089578063d547741f11610063578063d547741f14610539578063e138a8d214610558578063f2fde38b14610577578063fb894029146105965761019f565b8063b1454caa146104db578063b201246f146104fb578063b6aed0cb1461051a5761019f565b8063a217fddf116100ba578063a217fddf14610476578063ab53bddc14610489578063ad7805e8146104a85761019f565b80638da5cb5b146103b057806391d14854146103f45780639730886d146104575761019f565b8063248a9ca31161014d57806336568abe1161012757806336568abe1461033f57806336d2da901461035e578063485cc9551461037d578063715018a61461039c5761019f565b8063248a9ca3146102b45780632f2ff15d1461030157806333a88c72146103205761019f565b80630fcfbd111161017e5780630fcfbd11146102555780630fe9188e146102745780631050afdd146102955761019f565b8062a1b815146101cc57806301ffc9a7146101f657806302b4df1914610222575b3480156101aa575f5ffd5b5060405162461bcd60e51b81526004016101c3906115fc565b60405180910390fd5b3480156101d7575f5ffd5b506101e06105b5565b6040516101ed9190611614565b60405180910390f35b348015610201575f5ffd5b50610215610210366004611643565b61063e565b6040516101ed9190611668565b34801561022d575f5ffd5b506101e07f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19481565b348015610260575f5ffd5b506101e061026f36600461168f565b6106a6565b34801561027f575f5ffd5b5061029361028e3660046116d8565b610703565b005b3480156102a0575f5ffd5b506102936102af366004611719565b61076a565b3480156102bf575f5ffd5b506101e06102ce3660046116d8565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b34801561030c575f5ffd5b5061029361031b366004611736565b6107a2565b34801561032b575f5ffd5b5061021561033a36600461168f565b6107eb565b34801561034a575f5ffd5b50610293610359366004611736565b61083b565b348015610369575f5ffd5b50610293610378366004611719565b61088c565b348015610388575f5ffd5b5061029361039736600461176c565b610903565b3480156103a7575f5ffd5b50610293610aad565b3480156103bb575f5ffd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516101ed9190611793565b3480156103ff575f5ffd5b5061021561040e366004611736565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b348015610462575f5ffd5b506102936104713660046117a1565b610ac0565b348015610481575f5ffd5b506101e05f81565b348015610494575f5ffd5b506102936104a33660046117ea565b610c25565b3480156104b3575f5ffd5b506101e07fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9881565b6104ee6104e9366004611890565b610ced565b6040516101ed9190611918565b348015610506575f5ffd5b50610293610515366004611980565b610df5565b348015610525575f5ffd5b506102936105343660046119e7565b610ef3565b348015610544575f5ffd5b50610293610553366004611736565b610f5a565b348015610563575f5ffd5b50610293610572366004611a05565b610f9d565b348015610582575f5ffd5b50610293610591366004611719565b6110e4565b3480156105a1575f5ffd5b506102936105b0366004611719565b61113a565b600354604080517f1a90a21900000000000000000000000000000000000000000000000000000000815290515f926001600160a01b031691631a90a2199160048083019260209291908290030181865afa158015610615573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106399190611a88565b905090565b5f6001600160e01b031982167f7965db0b0000000000000000000000000000000000000000000000000000000014806106a057507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b5f5f826040516020016106b99190611c3a565b60408051601f1981840301815291815281516020928301205f81815292839052912054909150806106fc5760405162461bcd60e51b81526004016101c390611c89565b9392505050565b7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19461072d8161116e565b5f8281526004602052604081205490036107595760405162461bcd60e51b81526004016101c390611ccb565b505f90815260046020526040812055565b5f6107748161116e565b61079e7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19483610f5a565b5050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546107db8161116e565b6107e58383611178565b50505050565b5f5f826040516020016107fe9190611c3a565b60408051601f1981840301815291815281516020928301205f8181529283905291205490915080158015906108335750428111155b949350505050565b6001600160a01b038116331461087d576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108878282611244565b505050565b6108946112e8565b5f816001600160a01b0316476040515f6040518083038185875af1925050503d805f81146108dd576040519150601f19603f3d011682016040523d82523d5f602084013e6108e2565b606091505b505090508061079e5760405162461bcd60e51b81526004016101c390611d0d565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f8115801561094d5750825b90505f8267ffffffffffffffff1660011480156109695750303b155b905081158015610977575080155b156109ae576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156109e257845468ff00000000000000001916680100000000000000001785555b6109eb8761135c565b6109f361136d565b6109fd5f88611178565b50610a287f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e19488611178565b50610a537fe0d563514842a8c29151c49cd2698127f54dd344a9b2c74a42fe9be3e305fe9887611178565b508315610aa457845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610a9b90600190611d40565b60405180910390a15b50505050505050565b610ab56112e8565b610abe5f611375565b565b5f610acc600130611d62565b9050610aff7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b03161480610b265750336001600160a01b038216145b610b425760405162461bcd60e51b81526004016101c390611db7565b5f610b4d8342611dc7565b90505f84604051602001610b619190611c3a565b60408051601f1981840301815291815281516020928301205f8181529283905291205490915015610ba45760405162461bcd60e51b81526004016101c390611e32565b5f81815260208181526040822084905560019190610bc490880188611719565b6001600160a01b0316815260208101919091526040015f90812090610bef6080880160608901611e42565b63ffffffff1681526020808201929092526040015f9081208054600181018255908252919020869160040201610aa48282612240565b5f610c31600130611d62565b9050610c647f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b03161480610c8b5750336001600160a01b038216145b610ca75760405162461bcd60e51b81526004016101c390611db7565b826001600160a01b03167fcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d383604051610ce09190611614565b60405180910390a2505050565b6003545f906001600160a01b031615610d9e575f610d096105b5565b905080341015610d2b5760405162461bcd60e51b81526004016101c3906122a2565b6003546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610d75576040519150601f19603f3d011682016040523d82523d5f602084013e610d7a565b606091505b5050905080610d9b5760405162461bcd60e51b81526004016101c39061230a565b50505b610da7336113fd565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef77593733828888888888604051610de4979695949392919061231a565b60405180910390a195945050505050565b5f818152600460205260408120549003610e215760405162461bcd60e51b81526004016101c3906123d5565b5f81815260046020526040902054421015610e4e5760405162461bcd60e51b81526004016101c390612421565b5f84604051602001610e6091906124a5565b60405160208183030381529060405280519060200120604051602001610e8691906124e5565b604051602081830303815290604052805190602001209050610ed084848484604051602001610eb5919061250a565b6040516020818303038152906040528051906020012061145a565b610eec5760405162461bcd60e51b81526004016101c390612574565b5050505050565b7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e194610f1d8161116e565b5f8381526004602052604090205415610f485760405162461bcd60e51b81526004016101c3906125dc565b505f9182526004602052604090912055565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610f938161116e565b6107e58383611244565b5f818152600460205260408120549003610fc95760405162461bcd60e51b81526004016101c3906123d5565b5f81815260046020526040902054421015610ff65760405162461bcd60e51b81526004016101c390612421565b5f6110046020860186611719565b61101460408701602088016125ec565b6110246060880160408901611e42565b6110346080890160608a01611e42565b61104160808a018a611f74565b61105160c08c0160a08d01612609565b604051602001611067979695949392919061231a565b6040516020818303038152906040528051906020012090505f816040516020016110919190612658565b6040516020818303038152906040528051906020012090506110c085858584604051602001610eb5919061250a565b6110dc5760405162461bcd60e51b81526004016101c3906126c0565b505050505050565b6110ec6112e8565b6001600160a01b03811661112e575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016101c39190611793565b61113781611375565b50565b5f6111448161116e565b61079e7f65c4b771cce18ff228842b3883a73079ee3e76bd08965f6dadb7cb56dbf6e194836107a2565b6111378133611471565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff1661123b575f848152602082815260408083206001600160a01b03871684529091529020805460ff191660011790556111f13390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019150506106a0565b5f9150506106a0565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff161561123b575f848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a460019150506106a0565b3361131a7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b031614610abe57336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016101c39190611793565b6113646114ef565b61113781611556565b610abe6114ef565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080547fffffffffffffffffffffffff000000000000000000000000000000000000000081166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b6001600160a01b0381165f908152600260205260408120805467ffffffffffffffff16916001919061142f83856126d0565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b5f8261146786868561155e565b1495945050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff1661079e5780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016101c39291906126f4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610abe576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6110ec6114ef565b5f81815b848110156115965761158c828787848181106115805761158061270f565b9050602002013561159f565b9150600101611562565b50949350505050565b5f8183106115b9575f8281526020849052604090206106fc565b505f9182526020526040902090565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b602080825281016106a0816115c8565b805b82525050565b602081016106a0828461160c565b6001600160e01b031981165b8114611137575f5ffd5b80356106a081611622565b5f60208284031215611656576116565f5ffd5b6106fc8383611638565b80151561160e565b602081016106a08284611660565b5f60c08284031215611689576116895f5ffd5b50919050565b5f602082840312156116a2576116a25f5ffd5b813567ffffffffffffffff8111156116bb576116bb5f5ffd5b61083384828501611676565b8061162e565b80356106a0816116c7565b5f602082840312156116eb576116eb5f5ffd5b6106fc83836116cd565b5f6001600160a01b0382166106a0565b61162e816116f5565b80356106a081611705565b5f6020828403121561172c5761172c5f5ffd5b6106fc838361170e565b5f5f6040838503121561174a5761174a5f5ffd5b61175484846116cd565b9150611763846020850161170e565b90509250929050565b5f5f60408385031215611780576117805f5ffd5b611754848461170e565b61160e816116f5565b602081016106a0828461178a565b5f5f604083850312156117b5576117b55f5ffd5b823567ffffffffffffffff8111156117ce576117ce5f5ffd5b6117da85828601611676565b92505061176384602085016116cd565b5f5f604083850312156117fe576117fe5f5ffd5b611808848461170e565b915061176384602085016116cd565b63ffffffff811661162e565b80356106a081611817565b5f5f83601f840112611841576118415f5ffd5b50813567ffffffffffffffff81111561185b5761185b5f5ffd5b602083019150836001820283011115611875576118755f5ffd5b9250929050565b60ff811661162e565b80356106a08161187c565b5f5f5f5f5f608086880312156118a7576118a75f5ffd5b6118b18787611823565b94506118c08760208801611823565b9350604086013567ffffffffffffffff8111156118de576118de5f5ffd5b6118ea8882890161182e565b93509350506118fc8760608801611885565b90509295509295909350565b67ffffffffffffffff811661160e565b602081016106a08284611908565b5f60808284031215611689576116895f5ffd5b5f5f83601f84011261194c5761194c5f5ffd5b50813567ffffffffffffffff811115611966576119665f5ffd5b602083019150836020820283011115611875576118755f5ffd5b5f5f5f5f60c08587031215611996576119965f5ffd5b6119a08686611926565b9350608085013567ffffffffffffffff8111156119be576119be5f5ffd5b6119ca87828801611939565b93509350506119dc8660a087016116cd565b905092959194509250565b5f5f604083850312156119fb576119fb5f5ffd5b61180884846116cd565b5f5f5f5f60608587031215611a1b57611a1b5f5ffd5b843567ffffffffffffffff811115611a3457611a345f5ffd5b611a4087828801611676565b945050602085013567ffffffffffffffff811115611a5f57611a5f5f5ffd5b611a6b87828801611939565b93509350506119dc86604087016116cd565b80516106a0816116c7565b5f60208284031215611a9b57611a9b5f5ffd5b6106fc8383611a7d565b505f6106a0602083018361170e565b67ffffffffffffffff811661162e565b80356106a081611ab4565b505f6106a06020830183611ac4565b505f6106a06020830183611823565b63ffffffff811661160e565b5f808335601e1936859003018112611b1257611b125f5ffd5b830160208101925035905067ffffffffffffffff811115611b3457611b345f5ffd5b36819003821315611875576118755f5ffd5b82818337505f910152565b818352602083019250611b65828483611b46565b50601f01601f19160190565b505f6106a06020830183611885565b60ff811661160e565b5f60c08301611b988380611aa5565b611ba2858261178a565b50611bb06020840184611acf565b611bbd6020860182611908565b50611bcb6040840184611ade565b611bd86040860182611aed565b50611be66060840184611ade565b611bf36060860182611aed565b50611c016080840184611af9565b8583036080870152611c14838284611b51565b92505050611c2560a0840184611b71565b611c3260a0860182611b80565b509392505050565b602080825281016106fc8184611b89565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d69747465648152601760f91b602082015290505b60400190565b602080825281016106a081611c4b565b601a8152602081017f537461746520726f6f7420646f6573206e6f742065786973742e000000000000815290506115f6565b602080825281016106a081611c99565b60148152602081017f6661696c65642073656e64696e672076616c7565000000000000000000000000815290506115f6565b602080825281016106a081611cdb565b5f6106a082611d2a565b90565b67ffffffffffffffff1690565b61160e81611d1d565b602081016106a08284611d37565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b039182169190811690828203908111156106a0576106a0611d4e565b60118152602081017f4e6f74206f776e6572206f722073656c66000000000000000000000000000000815290506115f6565b602080825281016106a081611d85565b808201808211156106a0576106a0611d4e565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f210000000000000000000000000000000000000000000000000000000000000060208201529050611c83565b602080825281016106a081611dda565b5f60208284031215611e5557611e555f5ffd5b6106fc8383611823565b5f81356106a081611705565b5f6001600160a01b03835b81169019929092169190911792915050565b5f6106a0826116f5565b5f6106a082611e88565b611ea582611e92565b611eb0818354611e6b565b8255505050565b5f81356106a081611ab4565b5f7bffffffffffffffff0000000000000000000000000000000000000000611e768460a01b90565b5f6106a067ffffffffffffffff8316611d2a565b611f0882611eeb565b611eb0818354611ec3565b5f81356106a081611817565b5f6001600160e01b0319611e768460e01b90565b5f63ffffffff82166106a0565b611f4982611f33565b611eb0818354611f1f565b5f63ffffffff83611e76565b611f6982611f33565b611eb0818354611f54565b5f808335601e1936859003018112611f8d57611f8d5f5ffd5b8301915050803567ffffffffffffffff811115611fab57611fab5f5ffd5b602082019150600181023603821315611875576118755f5ffd5b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061200157607f821691505b60208210810361168957611689611fd9565b5f6106a0611d278381565b61202783612013565b81545f1960089490940293841b1916921b91909117905550565b5f61088781848461201e565b8181101561079e5761205f5f82612041565b60010161204d565b601f821115610887575f818152602090206020601f8501048101602085101561208d5750805b610eec6020601f86010483018261204d565b8267ffffffffffffffff8111156120b8576120b8611fc5565b6120c28254611fed565b6120cd828285612067565b505f601f8211600181146120ff575f83156120e85750848201355b5f19600885021c19811660028502178555506110dc565b5f84815260208120601f198516915b8281101561212e578785013582556020948501946001909201910161210e565b508482101561214a575f196008601f8716021c19878501351681555b5050505060020260010190555050565b61088783838361209f565b5f81356106a08161187c565b5f60ff82166106a0565b61218482612171565b815460ff191660ff821617611eb0565b8082806121a081611e5f565b90506121ac8184611e9c565b505060208301806121bc82611eb7565b90506121c88184611eff565b505060408301806121d882611f13565b90506121e48184611f40565b50505060608201806121f582611f13565b90506122048160018501611f60565b50506122136080830183611f74565b61222181836002860161215a565b505060a082018061223182612165565b90506107e5816003850161217b565b61079e8282612194565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f737361676500000000000000000000000000000000000000000000000000000060208201529050611c83565b602080825281016106a08161224a565b60248152602081017f4661696c656420746f2073656e64206665657320746f206665657320636f6e7481527f726163740000000000000000000000000000000000000000000000000000000060208201529050611c83565b602080825281016106a0816122b2565b60c08101612328828a61178a565b6123356020830189611908565b6123426040830188611aed565b61234f6060830187611aed565b8181036080830152612362818587611b51565b905061237160a0830184611b80565b98975050505050505050565b602a8152602081017f526f6f74206973206e6f74207075626c6973686564206f6e2074686973206d6581527f7373616765206275732e0000000000000000000000000000000000000000000060208201529050611c83565b602080825281016106a08161237d565b60218152602081017f526f6f74206973206e6f7420636f6e736964657265642066696e616c207965748152601760f91b60208201529050611c83565b602080825281016106a0816123e5565b505f6106a060208301836116cd565b61244a8180611aa5565b612454838261178a565b506124626020820182611aa5565b61246f602084018261178a565b5061247d6040820182612431565b61248a604084018261160c565b506124986060820182611acf565b6108876060840182611908565b608081016106a08284612440565b60018152602081017f7600000000000000000000000000000000000000000000000000000000000000815290506115f6565b604080825281016124f5816124b3565b90506106a0602083018461160c565b8061160e565b6125148183612504565b602001919050565b60338152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722076616c7581527f65207472616e73666572206d6573736167652e0000000000000000000000000060208201529050611c83565b602080825281016106a08161251c565b60258152602081017f526f6f7420616c726561647920616464656420746f20746865206d657373616781527f652062757300000000000000000000000000000000000000000000000000000060208201529050611c83565b602080825281016106a081612584565b5f602082840312156125ff576125ff5f5ffd5b6106fc8383611ac4565b5f6020828403121561261c5761261c5f5ffd5b6106fc8383611885565b60018152602081017f6d00000000000000000000000000000000000000000000000000000000000000815290506115f6565b604080825281016124f581612626565b60308152602081017f496e76616c696420696e636c7573696f6e2070726f6f6620666f722063726f7381527f7320636861696e206d6573736167652e0000000000000000000000000000000060208201529050611c83565b602080825281016106a081612668565b67ffffffffffffffff9182169190811690828201908111156106a0576106a0611d4e565b60408101612702828561178a565b6106fc602083018461160c565b634e487b7160e01b5f52603260045260245ffdfea2646970667358221220578b0a41c24ca2d361d6d6084b393db3afeae9e3dbe03afe78c53cc4f87371b964736f6c634300081c0033",
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

// NotifyDeposit is a paid mutator transaction binding the contract method 0xab53bddc.
//
// Solidity: function notifyDeposit(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) NotifyDeposit(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "notifyDeposit", receiver, amount)
}

// NotifyDeposit is a paid mutator transaction binding the contract method 0xab53bddc.
//
// Solidity: function notifyDeposit(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) NotifyDeposit(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.NotifyDeposit(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// NotifyDeposit is a paid mutator transaction binding the contract method 0xab53bddc.
//
// Solidity: function notifyDeposit(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) NotifyDeposit(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.NotifyDeposit(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
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

// MerkleTreeMessageBusNativeDepositIterator is returned from FilterNativeDeposit and is used to iterate over the raw logs and unpacked data for NativeDeposit events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusNativeDepositIterator struct {
	Event *MerkleTreeMessageBusNativeDeposit // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusNativeDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusNativeDeposit)
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
		it.Event = new(MerkleTreeMessageBusNativeDeposit)
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
func (it *MerkleTreeMessageBusNativeDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusNativeDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusNativeDeposit represents a NativeDeposit event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusNativeDeposit struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNativeDeposit is a free log retrieval operation binding the contract event 0xcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d3.
//
// Solidity: event NativeDeposit(address indexed receiver, uint256 amount)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterNativeDeposit(opts *bind.FilterOpts, receiver []common.Address) (*MerkleTreeMessageBusNativeDepositIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "NativeDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusNativeDepositIterator{contract: _MerkleTreeMessageBus.contract, event: "NativeDeposit", logs: logs, sub: sub}, nil
}

// WatchNativeDeposit is a free log subscription operation binding the contract event 0xcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d3.
//
// Solidity: event NativeDeposit(address indexed receiver, uint256 amount)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchNativeDeposit(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusNativeDeposit, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "NativeDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusNativeDeposit)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "NativeDeposit", log); err != nil {
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

// ParseNativeDeposit is a log parse operation binding the contract event 0xcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d3.
//
// Solidity: event NativeDeposit(address indexed receiver, uint256 amount)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseNativeDeposit(log types.Log) (*MerkleTreeMessageBusNativeDeposit, error) {
	event := new(MerkleTreeMessageBusNativeDeposit)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "NativeDeposit", log); err != nil {
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
