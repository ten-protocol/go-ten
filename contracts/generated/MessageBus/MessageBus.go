// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MessageBus

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

// MessageBusMetaData contains all meta data concerning the MessageBus contract.
var MessageBusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UNPAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublishFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantPauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantUnpauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"withdrawal\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feesAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokePauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeUnpauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"multisig\",\"type\":\"address\"}],\"name\":\"transferUnpauserRoleToMultisig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b5061001861002d565b61002061002d565b61002861002d565b61010c565b5f6100366100cd565b805490915068010000000000000000900460ff16156100685760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100ca5780546001600160401b0319166001600160401b0390811782556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916100c1916100f7565b60405180910390a15b50565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6001600160401b0382168152602081016100f1565b612479806101195f395ff3fe6080604052600436106101b5575f3560e01c806379ba5097116100eb578063d547741f11610089578063ea4cf97011610063578063ea4cf97014610541578063f2fde38b14610560578063f865af081461057f578063fb1bb9de1461059e576101b5565b8063d547741f146104db578063e30c3978146104fa578063e63ab1e91461050e576101b5565b806391643fdd116100c557806391643fdd1461042757806391d1485414610446578063a217fddf146104a9578063c0c53b8b146104bc576101b5565b806379ba5097146103de5780638456cb59146103f25780638da5cb5b14610406576101b5565b806332968782116101585780635c975abb116101325780635c975abb146103565780636c11c21c1461038c578063715018a6146103ab5780637920c986146103bf576101b5565b8063329687821461030457806336568abe146103235780633f4ba83a14610342576101b5565b8063248a9ca311610194578063248a9ca3146102585780632540e2da146102a55780632e1a0b8e146102c65780632f2ff15d146102e5576101b5565b8062a1b815146101e257806301ffc9a71461020c5780630d3fd67c14610238575b3480156101c0575f5ffd5b5060405162461bcd60e51b81526004016101d9906116ec565b60405180910390fd5b3480156101ed575f5ffd5b506101f66105d1565b6040516102039190611704565b60405180910390f35b348015610217575f5ffd5b5061022b61022636600461174b565b61065a565b6040516102039190611770565b61024b610246366004611812565b6106f2565b604051610203919061189a565b348015610263575f5ffd5b506101f66102723660046118b9565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b3480156102b0575f5ffd5b506102c46102bf3660046118fa565b610801565b005b3480156102d1575f5ffd5b506101f66102e0366004611930565b61083a565b3480156102f0575f5ffd5b506102c46102ff366004611968565b610897565b34801561030f575f5ffd5b506102c461031e3660046118fa565b6108e0565b34801561032e575f5ffd5b506102c461033d366004611968565b610914565b34801561034d575f5ffd5b506102c4610960565b348015610361575f5ffd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1661022b565b348015610397575f5ffd5b506102c46103a63660046118fa565b610995565b3480156103b6575f5ffd5b506102c46109c9565b3480156103ca575f5ffd5b506102c46103d93660046118fa565b6109e9565b3480156103e9575f5ffd5b506102c4610a6e565b3480156103fd575f5ffd5b506102c4610aaa565b348015610411575f5ffd5b5061041a610adc565b60405161020391906119a7565b348015610432575f5ffd5b5061022b610441366004611930565b610b10565b348015610451575f5ffd5b5061022b610460366004611968565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b3480156104b4575f5ffd5b506101f65f81565b3480156104c7575f5ffd5b506102c46104d63660046119b5565b610b60565b3480156104e6575f5ffd5b506102c46104f5366004611968565b610d0a565b348015610505575f5ffd5b5061041a610d4d565b348015610519575f5ffd5b506101f67f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a81565b34801561054c575f5ffd5b506102c461055b3660046119fb565b610d75565b34801561056b575f5ffd5b506102c461057a3660046118fa565b610ec2565b34801561058a575f5ffd5b506102c46105993660046118fa565b610f54565b3480156105a9575f5ffd5b506101f67f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a81565b600354604080517f1a90a21900000000000000000000000000000000000000000000000000000000815290515f926001600160a01b031691631a90a2199160048083019260209291908290030181865afa158015610631573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106559190611a4f565b905090565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806106ec57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f6106fb610f88565b6003546001600160a01b0316156107aa575f6107156105d1565b9050803410156107375760405162461bcd60e51b81526004016101d990611ac6565b6003546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610781576040519150601f19603f3d011682016040523d82523d5f602084013e610786565b606091505b50509050806107a75760405162461bcd60e51b81526004016101d990611b2e565b50505b6107b333610fe6565b90507fd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061338288888888886040516107f09796959493929190611b7e565b60405180910390a195945050505050565b5f61080b81611043565b6108357f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a8361104d565b505050565b5f5f8260405160200161084d9190611d1b565b60408051601f1981840301815291815281516020928301205f81815292839052912054909150806108905760405162461bcd60e51b81526004016101d990611d84565b9392505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546108d081611043565b6108da83836110fa565b50505050565b5f6108ea81611043565b6108357f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a836110fa565b6001600160a01b0381163314610956576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610835828261104d565b7f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a61098a81611043565b6109926111bd565b50565b5f61099f81611043565b6108357f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a836110fa565b6109d1611229565b60405162461bcd60e51b81526004016101d990611dec565b5f6109f381611043565b6001600160a01b038216610a195760405162461bcd60e51b81526004016101d990611e2e565b610a437f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a3361104d565b506108357f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a836110fa565b3380610a78610d4d565b6001600160a01b031614610aa1578060405163118cdaa760e01b81526004016101d991906119a7565b6109928161125b565b7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a610ad481611043565b6109926112a4565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b5f5f82604051602001610b239190611d1b565b60408051601f1981840301815291815281516020928301205f818152928390529120549091508015801590610b585750428111155b949350505050565b5f610b696112ff565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f81158015610b955750825b90505f8267ffffffffffffffff166001148015610bb15750303b155b905081158015610bbf575080155b15610bf6576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610c2a57845468ff00000000000000001916680100000000000000001785555b6001600160a01b038616610c505760405162461bcd60e51b81526004016101d990611e70565b6001600160a01b038816610c765760405162461bcd60e51b81526004016101d990611eb2565b610c7f88611327565b610c8888611340565b6003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0388161790558315610d0057845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610cf790600190611ee5565b60405180910390a15b5050505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610d4381611043565b6108da838361104d565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610b00565b5f610d81600130611f07565b9050610d8b610adc565b6001600160a01b0316336001600160a01b03161480610db25750336001600160a01b038216145b610dce5760405162461bcd60e51b81526004016101d990611f5c565b610dd6610f88565b5f610de18342611f6c565b90505f84604051602001610df59190611d1b565b60408051601f1981840301815291815281516020928301205f8181529283905291205490915015610e385760405162461bcd60e51b81526004016101d990611fd7565b5f81815260208181526040822084905560019190610e58908801886118fa565b6001600160a01b0316815260208101919091526040015f90812090610e836080880160608901611fe7565b63ffffffff1681526020808201929092526040015f9081208054600181018255908252919020869160040201610eb982826123fa565b50505050505050565b610eca611229565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081178255610f1b610adc565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b5f610f5e81611043565b6108357f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a8361104d565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1615610fe4576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6001600160a01b0381165f908152600260205260408120805467ffffffffffffffff1691600191906110188385612404565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b6109928133611402565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff16156110f1575f848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a460019150506106ec565b5f9150506106ec565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff166110f1575f848152602082815260408083206001600160a01b03871684529091529020805460ff191660011790556111733390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019150506106ec565b6111c5611480565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300805460ff191681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405161121e91906119a7565b60405180910390a150565b33611232610adc565b6001600160a01b031614610fe4573360405163118cdaa760e01b81526004016101d991906119a7565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff191681556112a0826114db565b5050565b6112ac610f88565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300805460ff191660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25833611211565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006106ec565b61132f611558565b61133881611596565b6109926115a7565b611348611558565b6113506115a7565b6113586115a7565b6113827f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a826110fa565b506113ad7f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a826110fa565b506113d87f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a5f6115af565b6109927f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a5f6115af565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff166112a05780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016101d9929190612428565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610fe4576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b611560611650565b610fe4576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61159e611558565b6109928161166e565b610fe4611558565b7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268005f611608845f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b5f85815260208490526040808220600101869055519192508491839187917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a450505050565b5f6116596112ff565b5468010000000000000000900460ff16919050565b611676611558565b6001600160a01b038116610aa1575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016101d991906119a7565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b602080825281016106ec816116b8565b805b82525050565b602081016106ec82846116fc565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b8114610992575f5ffd5b80356106ec81611712565b5f6020828403121561175e5761175e5f5ffd5b6108908383611740565b8015156116fe565b602081016106ec8284611768565b67ffffffffffffffff8116611736565b80356106ec8161177e565b63ffffffff8116611736565b80356106ec81611799565b5f5f83601f8401126117c3576117c35f5ffd5b50813567ffffffffffffffff8111156117dd576117dd5f5ffd5b6020830191508360018202830111156117f7576117f75f5ffd5b9250929050565b60ff8116611736565b80356106ec816117fe565b5f5f5f5f5f60808688031215611829576118295f5ffd5b611833878761178e565b945061184287602088016117a5565b9350604086013567ffffffffffffffff811115611860576118605f5ffd5b61186c888289016117b0565b935093505061187e8760608801611807565b90509295509295909350565b67ffffffffffffffff81166116fe565b602081016106ec828461188a565b80611736565b80356106ec816118a8565b5f602082840312156118cc576118cc5f5ffd5b61089083836118ae565b5f6001600160a01b0382166106ec565b611736816118d6565b80356106ec816118e6565b5f6020828403121561190d5761190d5f5ffd5b61089083836118ef565b5f60c0828403121561192a5761192a5f5ffd5b50919050565b5f60208284031215611943576119435f5ffd5b813567ffffffffffffffff81111561195c5761195c5f5ffd5b610b5884828501611917565b5f5f6040838503121561197c5761197c5f5ffd5b61198684846118ae565b915061199584602085016118ef565b90509250929050565b6116fe816118d6565b602081016106ec828461199e565b5f5f5f606084860312156119ca576119ca5f5ffd5b6119d485856118ef565b92506119e385602086016118ef565b91506119f285604086016118ef565b90509250925092565b5f5f60408385031215611a0f57611a0f5f5ffd5b823567ffffffffffffffff811115611a2857611a285f5ffd5b611a3485828601611917565b92505061199584602085016118ae565b80516106ec816118a8565b5f60208284031215611a6257611a625f5ffd5b6108908383611a44565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f7373616765000000000000000000000000000000000000000000000000000000602082015290505b60400190565b602080825281016106ec81611a6c565b60248152602081017f4661696c656420746f2073656e64206665657320746f206665657320636f6e7481527f726163740000000000000000000000000000000000000000000000000000000060208201529050611ac0565b602080825281016106ec81611ad6565b63ffffffff81166116fe565b82818337505f910152565b818352602083019250611b69828483611b4a565b50601f01601f19160190565b60ff81166116fe565b60c08101611b8c828a61199e565b611b99602083018961188a565b611ba6604083018861188a565b611bb36060830187611b3e565b8181036080830152611bc6818587611b55565b9050611bd560a0830184611b75565b98975050505050505050565b505f6106ec60208301836118ef565b505f6106ec602083018361178e565b505f6106ec60208301836117a5565b5f808335601e1936859003018112611c2757611c275f5ffd5b830160208101925035905067ffffffffffffffff811115611c4957611c495f5ffd5b368190038213156117f7576117f75f5ffd5b505f6106ec6020830183611807565b5f60c08301611c798380611be1565b611c83858261199e565b50611c916020840184611bf0565b611c9e602086018261188a565b50611cac6040840184611bf0565b611cb9604086018261188a565b50611cc76060840184611bff565b611cd46060860182611b3e565b50611ce26080840184611c0e565b8583036080870152611cf5838284611b55565b92505050611d0660a0840184611c5b565b611d1360a0860182611b75565b509392505050565b602080825281016108908184611c6a565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d697474656481527f2e0000000000000000000000000000000000000000000000000000000000000060208201529050611ac0565b602080825281016106ec81611d2c565b60348152602081017f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f81527f742072656e6f756e6365206f776e65727368697000000000000000000000000060208201529050611ac0565b602080825281016106ec81611d94565b60188152602081017f496e76616c6964206d756c746973696720616464726573730000000000000000815290506116e6565b602080825281016106ec81611dfc565b601a8152602081017f4665657320616464726573732063616e6e6f7420626520307830000000000000815290506116e6565b602080825281016106ec81611e3e565b60148152602081017f43616c6c65722063616e6e6f7420626520307830000000000000000000000000815290506116e6565b602080825281016106ec81611e80565b5f6106ec82611ecf565b90565b67ffffffffffffffff1690565b6116fe81611ec2565b602081016106ec8284611edc565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b039182169190811690828203908111156106ec576106ec611ef3565b60118152602081017f4e6f74206f776e6572206f722073656c66000000000000000000000000000000815290506116e6565b602080825281016106ec81611f2a565b808201808211156106ec576106ec611ef3565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f210000000000000000000000000000000000000000000000000000000000000060208201529050611ac0565b602080825281016106ec81611f7f565b5f60208284031215611ffa57611ffa5f5ffd5b61089083836117a5565b5f81356106ec816118e6565b5f6001600160a01b03835b81169019929092169190911792915050565b5f6106ec826118d6565b5f6106ec8261202d565b61204a82612037565b612055818354612010565b8255505050565b5f81356106ec8161177e565b5f7bffffffffffffffff000000000000000000000000000000000000000061201b8460a01b90565b5f6106ec67ffffffffffffffff8316611ecf565b6120ad82612090565b612055818354612068565b5f67ffffffffffffffff8361201b565b6120d182612090565b6120558183546120b8565b5f81356106ec81611799565b5f6bffffffff000000000000000061201b8460401b90565b5f63ffffffff82166106ec565b61211682612100565b6120558183546120e8565b5f808335601e193685900301811261213a5761213a5f5ffd5b8301915050803567ffffffffffffffff811115612158576121585f5ffd5b6020820191506001810236038213156117f7576117f75f5ffd5b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b6002810460018216806121ae57607f821691505b60208210810361192a5761192a612186565b5f6106ec611ecc8381565b6121d4836121c0565b81545f1960089490940293841b1916921b91909117905550565b5f6108358184846121cb565b818110156112a05761220c5f826121ee565b6001016121fa565b601f821115610835575f818152602090206020601f8501048101602085101561223a5750805b61224c6020601f8601048301826121fa565b5050505050565b8267ffffffffffffffff81111561226c5761226c612172565b612276825461219a565b612281828285612214565b505f601f8211600181146122b3575f831561229c5750848201355b5f19600885021c198116600285021785555061230a565b5f84815260208120601f198516915b828110156122e257878501358255602094850194600190920191016122c2565b50848210156122fe575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b610835838383612253565b5f81356106ec816117fe565b5f60ff82166106ec565b61233c82612329565b815460ff191660ff821617612055565b80828061235881612004565b90506123648184612041565b505060208301806123748261205c565b905061238081846120a4565b5050506001810160408301806123958261205c565b90506123a181846120c8565b505060608301806123b1826120dc565b90506123bd818461210d565b5050506123cd6080830183612121565b6123db818360028601612312565b505060a08201806123eb8261231d565b90506108da8160038501612333565b6112a0828261234c565b67ffffffffffffffff9182169190811690828201908111156106ec576106ec611ef3565b60408101612436828561199e565b61089060208301846116fc56fea2646970667358221220e2cd1ad96b3d2d1e9e8137cff5dead4bc236c97adedb81433d12b9659093eda164736f6c634300081c0033",
}

// MessageBusABI is the input ABI used to generate the binding from.
// Deprecated: Use MessageBusMetaData.ABI instead.
var MessageBusABI = MessageBusMetaData.ABI

// MessageBusBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MessageBusMetaData.Bin instead.
var MessageBusBin = MessageBusMetaData.Bin

// DeployMessageBus deploys a new Ethereum contract, binding an instance of MessageBus to it.
func DeployMessageBus(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MessageBus, error) {
	parsed, err := MessageBusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageBusBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MessageBus{MessageBusCaller: MessageBusCaller{contract: contract}, MessageBusTransactor: MessageBusTransactor{contract: contract}, MessageBusFilterer: MessageBusFilterer{contract: contract}}, nil
}

// MessageBus is an auto generated Go binding around an Ethereum contract.
type MessageBus struct {
	MessageBusCaller     // Read-only binding to the contract
	MessageBusTransactor // Write-only binding to the contract
	MessageBusFilterer   // Log filterer for contract events
}

// MessageBusCaller is an auto generated read-only Go binding around an Ethereum contract.
type MessageBusCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MessageBusTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MessageBusFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MessageBusSession struct {
	Contract     *MessageBus       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MessageBusCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MessageBusCallerSession struct {
	Contract *MessageBusCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// MessageBusTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MessageBusTransactorSession struct {
	Contract     *MessageBusTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MessageBusRaw is an auto generated low-level Go binding around an Ethereum contract.
type MessageBusRaw struct {
	Contract *MessageBus // Generic contract binding to access the raw methods on
}

// MessageBusCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MessageBusCallerRaw struct {
	Contract *MessageBusCaller // Generic read-only contract binding to access the raw methods on
}

// MessageBusTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MessageBusTransactorRaw struct {
	Contract *MessageBusTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMessageBus creates a new instance of MessageBus, bound to a specific deployed contract.
func NewMessageBus(address common.Address, backend bind.ContractBackend) (*MessageBus, error) {
	contract, err := bindMessageBus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageBus{MessageBusCaller: MessageBusCaller{contract: contract}, MessageBusTransactor: MessageBusTransactor{contract: contract}, MessageBusFilterer: MessageBusFilterer{contract: contract}}, nil
}

// NewMessageBusCaller creates a new read-only instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusCaller(address common.Address, caller bind.ContractCaller) (*MessageBusCaller, error) {
	contract, err := bindMessageBus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageBusCaller{contract: contract}, nil
}

// NewMessageBusTransactor creates a new write-only instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageBusTransactor, error) {
	contract, err := bindMessageBus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageBusTransactor{contract: contract}, nil
}

// NewMessageBusFilterer creates a new log filterer instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageBusFilterer, error) {
	contract, err := bindMessageBus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageBusFilterer{contract: contract}, nil
}

// bindMessageBus binds a generic wrapper to an already deployed contract.
func bindMessageBus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MessageBusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageBus *MessageBusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageBus.Contract.MessageBusCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageBus *MessageBusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.Contract.MessageBusTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageBus *MessageBusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageBus.Contract.MessageBusTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageBus *MessageBusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageBus.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageBus *MessageBusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageBus *MessageBusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageBus.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MessageBus.Contract.DEFAULTADMINROLE(&_MessageBus.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _MessageBus.Contract.DEFAULTADMINROLE(&_MessageBus.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusSession) PAUSERROLE() ([32]byte, error) {
	return _MessageBus.Contract.PAUSERROLE(&_MessageBus.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusCallerSession) PAUSERROLE() ([32]byte, error) {
	return _MessageBus.Contract.PAUSERROLE(&_MessageBus.CallOpts)
}

// UNPAUSERROLE is a free data retrieval call binding the contract method 0xfb1bb9de.
//
// Solidity: function UNPAUSER_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusCaller) UNPAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "UNPAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UNPAUSERROLE is a free data retrieval call binding the contract method 0xfb1bb9de.
//
// Solidity: function UNPAUSER_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusSession) UNPAUSERROLE() ([32]byte, error) {
	return _MessageBus.Contract.UNPAUSERROLE(&_MessageBus.CallOpts)
}

// UNPAUSERROLE is a free data retrieval call binding the contract method 0xfb1bb9de.
//
// Solidity: function UNPAUSER_ROLE() view returns(bytes32)
func (_MessageBus *MessageBusCallerSession) UNPAUSERROLE() ([32]byte, error) {
	return _MessageBus.Contract.UNPAUSERROLE(&_MessageBus.CallOpts)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MessageBus *MessageBusCaller) GetPublishFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getPublishFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MessageBus *MessageBusSession) GetPublishFee() (*big.Int, error) {
	return _MessageBus.Contract.GetPublishFee(&_MessageBus.CallOpts)
}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetPublishFee() (*big.Int, error) {
	return _MessageBus.Contract.GetPublishFee(&_MessageBus.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MessageBus *MessageBusCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MessageBus *MessageBusSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _MessageBus.Contract.GetRoleAdmin(&_MessageBus.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_MessageBus *MessageBusCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _MessageBus.Contract.GetRoleAdmin(&_MessageBus.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MessageBus *MessageBusCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MessageBus *MessageBusSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _MessageBus.Contract.HasRole(&_MessageBus.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_MessageBus *MessageBusCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _MessageBus.Contract.HasRole(&_MessageBus.CallOpts, role, account)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MessageBus *MessageBusCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MessageBus *MessageBusSession) Owner() (common.Address, error) {
	return _MessageBus.Contract.Owner(&_MessageBus.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MessageBus *MessageBusCallerSession) Owner() (common.Address, error) {
	return _MessageBus.Contract.Owner(&_MessageBus.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MessageBus *MessageBusCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MessageBus *MessageBusSession) Paused() (bool, error) {
	return _MessageBus.Contract.Paused(&_MessageBus.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MessageBus *MessageBusCallerSession) Paused() (bool, error) {
	return _MessageBus.Contract.Paused(&_MessageBus.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MessageBus *MessageBusCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MessageBus *MessageBusSession) PendingOwner() (common.Address, error) {
	return _MessageBus.Contract.PendingOwner(&_MessageBus.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MessageBus *MessageBusCallerSession) PendingOwner() (common.Address, error) {
	return _MessageBus.Contract.PendingOwner(&_MessageBus.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MessageBus *MessageBusCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MessageBus *MessageBusSession) RenounceOwnership() error {
	return _MessageBus.Contract.RenounceOwnership(&_MessageBus.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MessageBus *MessageBusCallerSession) RenounceOwnership() error {
	return _MessageBus.Contract.RenounceOwnership(&_MessageBus.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MessageBus *MessageBusCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MessageBus *MessageBusSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MessageBus.Contract.SupportsInterface(&_MessageBus.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_MessageBus *MessageBusCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MessageBus.Contract.SupportsInterface(&_MessageBus.CallOpts, interfaceId)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MessageBus *MessageBusTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MessageBus *MessageBusSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageBus.Contract.AcceptOwnership(&_MessageBus.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MessageBus *MessageBusTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageBus.Contract.AcceptOwnership(&_MessageBus.TransactOpts)
}

// GrantPauserRole is a paid mutator transaction binding the contract method 0x6c11c21c.
//
// Solidity: function grantPauserRole(address account) returns()
func (_MessageBus *MessageBusTransactor) GrantPauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "grantPauserRole", account)
}

// GrantPauserRole is a paid mutator transaction binding the contract method 0x6c11c21c.
//
// Solidity: function grantPauserRole(address account) returns()
func (_MessageBus *MessageBusSession) GrantPauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.GrantPauserRole(&_MessageBus.TransactOpts, account)
}

// GrantPauserRole is a paid mutator transaction binding the contract method 0x6c11c21c.
//
// Solidity: function grantPauserRole(address account) returns()
func (_MessageBus *MessageBusTransactorSession) GrantPauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.GrantPauserRole(&_MessageBus.TransactOpts, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MessageBus *MessageBusTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MessageBus *MessageBusSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.GrantRole(&_MessageBus.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_MessageBus *MessageBusTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.GrantRole(&_MessageBus.TransactOpts, role, account)
}

// GrantUnpauserRole is a paid mutator transaction binding the contract method 0x32968782.
//
// Solidity: function grantUnpauserRole(address account) returns()
func (_MessageBus *MessageBusTransactor) GrantUnpauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "grantUnpauserRole", account)
}

// GrantUnpauserRole is a paid mutator transaction binding the contract method 0x32968782.
//
// Solidity: function grantUnpauserRole(address account) returns()
func (_MessageBus *MessageBusSession) GrantUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.GrantUnpauserRole(&_MessageBus.TransactOpts, account)
}

// GrantUnpauserRole is a paid mutator transaction binding the contract method 0x32968782.
//
// Solidity: function grantUnpauserRole(address account) returns()
func (_MessageBus *MessageBusTransactorSession) GrantUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.GrantUnpauserRole(&_MessageBus.TransactOpts, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address owner, address withdrawal, address feesAddress) returns()
func (_MessageBus *MessageBusTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, withdrawal common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "initialize", owner, withdrawal, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address owner, address withdrawal, address feesAddress) returns()
func (_MessageBus *MessageBusSession) Initialize(owner common.Address, withdrawal common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.Initialize(&_MessageBus.TransactOpts, owner, withdrawal, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address owner, address withdrawal, address feesAddress) returns()
func (_MessageBus *MessageBusTransactorSession) Initialize(owner common.Address, withdrawal common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.Initialize(&_MessageBus.TransactOpts, owner, withdrawal, feesAddress)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MessageBus *MessageBusTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MessageBus *MessageBusSession) Pause() (*types.Transaction, error) {
	return _MessageBus.Contract.Pause(&_MessageBus.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MessageBus *MessageBusTransactorSession) Pause() (*types.Transaction, error) {
	return _MessageBus.Contract.Pause(&_MessageBus.TransactOpts)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusSession) PublishMessage(nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactorSession) PublishMessage(nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_MessageBus *MessageBusTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_MessageBus *MessageBusSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RenounceRole(&_MessageBus.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_MessageBus *MessageBusTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RenounceRole(&_MessageBus.TransactOpts, role, callerConfirmation)
}

// RevokePauserRole is a paid mutator transaction binding the contract method 0xf865af08.
//
// Solidity: function revokePauserRole(address account) returns()
func (_MessageBus *MessageBusTransactor) RevokePauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "revokePauserRole", account)
}

// RevokePauserRole is a paid mutator transaction binding the contract method 0xf865af08.
//
// Solidity: function revokePauserRole(address account) returns()
func (_MessageBus *MessageBusSession) RevokePauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RevokePauserRole(&_MessageBus.TransactOpts, account)
}

// RevokePauserRole is a paid mutator transaction binding the contract method 0xf865af08.
//
// Solidity: function revokePauserRole(address account) returns()
func (_MessageBus *MessageBusTransactorSession) RevokePauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RevokePauserRole(&_MessageBus.TransactOpts, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MessageBus *MessageBusTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MessageBus *MessageBusSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RevokeRole(&_MessageBus.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_MessageBus *MessageBusTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RevokeRole(&_MessageBus.TransactOpts, role, account)
}

// RevokeUnpauserRole is a paid mutator transaction binding the contract method 0x2540e2da.
//
// Solidity: function revokeUnpauserRole(address account) returns()
func (_MessageBus *MessageBusTransactor) RevokeUnpauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "revokeUnpauserRole", account)
}

// RevokeUnpauserRole is a paid mutator transaction binding the contract method 0x2540e2da.
//
// Solidity: function revokeUnpauserRole(address account) returns()
func (_MessageBus *MessageBusSession) RevokeUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RevokeUnpauserRole(&_MessageBus.TransactOpts, account)
}

// RevokeUnpauserRole is a paid mutator transaction binding the contract method 0x2540e2da.
//
// Solidity: function revokeUnpauserRole(address account) returns()
func (_MessageBus *MessageBusTransactorSession) RevokeUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RevokeUnpauserRole(&_MessageBus.TransactOpts, account)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactor) StoreCrossChainMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "storeCrossChainMessage", crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.StoreCrossChainMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactorSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.StoreCrossChainMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MessageBus *MessageBusTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MessageBus *MessageBusSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.TransferOwnership(&_MessageBus.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MessageBus *MessageBusTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.TransferOwnership(&_MessageBus.TransactOpts, newOwner)
}

// TransferUnpauserRoleToMultisig is a paid mutator transaction binding the contract method 0x7920c986.
//
// Solidity: function transferUnpauserRoleToMultisig(address multisig) returns()
func (_MessageBus *MessageBusTransactor) TransferUnpauserRoleToMultisig(opts *bind.TransactOpts, multisig common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "transferUnpauserRoleToMultisig", multisig)
}

// TransferUnpauserRoleToMultisig is a paid mutator transaction binding the contract method 0x7920c986.
//
// Solidity: function transferUnpauserRoleToMultisig(address multisig) returns()
func (_MessageBus *MessageBusSession) TransferUnpauserRoleToMultisig(multisig common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.TransferUnpauserRoleToMultisig(&_MessageBus.TransactOpts, multisig)
}

// TransferUnpauserRoleToMultisig is a paid mutator transaction binding the contract method 0x7920c986.
//
// Solidity: function transferUnpauserRoleToMultisig(address multisig) returns()
func (_MessageBus *MessageBusTransactorSession) TransferUnpauserRoleToMultisig(multisig common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.TransferUnpauserRoleToMultisig(&_MessageBus.TransactOpts, multisig)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MessageBus *MessageBusTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MessageBus *MessageBusSession) Unpause() (*types.Transaction, error) {
	return _MessageBus.Contract.Unpause(&_MessageBus.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MessageBus *MessageBusTransactorSession) Unpause() (*types.Transaction, error) {
	return _MessageBus.Contract.Unpause(&_MessageBus.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_MessageBus *MessageBusTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _MessageBus.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_MessageBus *MessageBusSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_MessageBus *MessageBusTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// MessageBusInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MessageBus contract.
type MessageBusInitializedIterator struct {
	Event *MessageBusInitialized // Event containing the contract specifics and raw log

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
func (it *MessageBusInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusInitialized)
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
		it.Event = new(MessageBusInitialized)
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
func (it *MessageBusInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusInitialized represents a Initialized event raised by the MessageBus contract.
type MessageBusInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MessageBus *MessageBusFilterer) FilterInitialized(opts *bind.FilterOpts) (*MessageBusInitializedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MessageBusInitializedIterator{contract: _MessageBus.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MessageBus *MessageBusFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MessageBusInitialized) (event.Subscription, error) {

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusInitialized)
				if err := _MessageBus.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseInitialized(log types.Log) (*MessageBusInitialized, error) {
	event := new(MessageBusInitialized)
	if err := _MessageBus.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusLogMessagePublishedIterator is returned from FilterLogMessagePublished and is used to iterate over the raw logs and unpacked data for LogMessagePublished events raised by the MessageBus contract.
type MessageBusLogMessagePublishedIterator struct {
	Event *MessageBusLogMessagePublished // Event containing the contract specifics and raw log

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
func (it *MessageBusLogMessagePublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusLogMessagePublished)
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
		it.Event = new(MessageBusLogMessagePublished)
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
func (it *MessageBusLogMessagePublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusLogMessagePublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusLogMessagePublished represents a LogMessagePublished event raised by the MessageBus contract.
type MessageBusLogMessagePublished struct {
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
func (_MessageBus *MessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts) (*MessageBusLogMessagePublishedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return &MessageBusLogMessagePublishedIterator{contract: _MessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) WatchLogMessagePublished(opts *bind.WatchOpts, sink chan<- *MessageBusLogMessagePublished) (event.Subscription, error) {

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusLogMessagePublished)
				if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MessageBusLogMessagePublished, error) {
	event := new(MessageBusLogMessagePublished)
	if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the MessageBus contract.
type MessageBusOwnershipTransferStartedIterator struct {
	Event *MessageBusOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *MessageBusOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusOwnershipTransferStarted)
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
		it.Event = new(MessageBusOwnershipTransferStarted)
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
func (it *MessageBusOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MessageBus contract.
type MessageBusOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MessageBusOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusOwnershipTransferStartedIterator{contract: _MessageBus.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *MessageBusOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusOwnershipTransferStarted)
				if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseOwnershipTransferStarted(log types.Log) (*MessageBusOwnershipTransferStarted, error) {
	event := new(MessageBusOwnershipTransferStarted)
	if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MessageBus contract.
type MessageBusOwnershipTransferredIterator struct {
	Event *MessageBusOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MessageBusOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusOwnershipTransferred)
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
		it.Event = new(MessageBusOwnershipTransferred)
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
func (it *MessageBusOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusOwnershipTransferred represents a OwnershipTransferred event raised by the MessageBus contract.
type MessageBusOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MessageBusOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusOwnershipTransferredIterator{contract: _MessageBus.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MessageBusOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusOwnershipTransferred)
				if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseOwnershipTransferred(log types.Log) (*MessageBusOwnershipTransferred, error) {
	event := new(MessageBusOwnershipTransferred)
	if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the MessageBus contract.
type MessageBusPausedIterator struct {
	Event *MessageBusPaused // Event containing the contract specifics and raw log

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
func (it *MessageBusPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusPaused)
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
		it.Event = new(MessageBusPaused)
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
func (it *MessageBusPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusPaused represents a Paused event raised by the MessageBus contract.
type MessageBusPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MessageBus *MessageBusFilterer) FilterPaused(opts *bind.FilterOpts) (*MessageBusPausedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MessageBusPausedIterator{contract: _MessageBus.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MessageBus *MessageBusFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MessageBusPaused) (event.Subscription, error) {

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusPaused)
				if err := _MessageBus.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MessageBus *MessageBusFilterer) ParsePaused(log types.Log) (*MessageBusPaused, error) {
	event := new(MessageBusPaused)
	if err := _MessageBus.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the MessageBus contract.
type MessageBusRoleAdminChangedIterator struct {
	Event *MessageBusRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *MessageBusRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusRoleAdminChanged)
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
		it.Event = new(MessageBusRoleAdminChanged)
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
func (it *MessageBusRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusRoleAdminChanged represents a RoleAdminChanged event raised by the MessageBus contract.
type MessageBusRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_MessageBus *MessageBusFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*MessageBusRoleAdminChangedIterator, error) {

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

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusRoleAdminChangedIterator{contract: _MessageBus.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_MessageBus *MessageBusFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *MessageBusRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusRoleAdminChanged)
				if err := _MessageBus.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseRoleAdminChanged(log types.Log) (*MessageBusRoleAdminChanged, error) {
	event := new(MessageBusRoleAdminChanged)
	if err := _MessageBus.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the MessageBus contract.
type MessageBusRoleGrantedIterator struct {
	Event *MessageBusRoleGranted // Event containing the contract specifics and raw log

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
func (it *MessageBusRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusRoleGranted)
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
		it.Event = new(MessageBusRoleGranted)
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
func (it *MessageBusRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusRoleGranted represents a RoleGranted event raised by the MessageBus contract.
type MessageBusRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_MessageBus *MessageBusFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MessageBusRoleGrantedIterator, error) {

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

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusRoleGrantedIterator{contract: _MessageBus.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_MessageBus *MessageBusFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *MessageBusRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusRoleGranted)
				if err := _MessageBus.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseRoleGranted(log types.Log) (*MessageBusRoleGranted, error) {
	event := new(MessageBusRoleGranted)
	if err := _MessageBus.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the MessageBus contract.
type MessageBusRoleRevokedIterator struct {
	Event *MessageBusRoleRevoked // Event containing the contract specifics and raw log

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
func (it *MessageBusRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusRoleRevoked)
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
		it.Event = new(MessageBusRoleRevoked)
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
func (it *MessageBusRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusRoleRevoked represents a RoleRevoked event raised by the MessageBus contract.
type MessageBusRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_MessageBus *MessageBusFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MessageBusRoleRevokedIterator, error) {

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

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusRoleRevokedIterator{contract: _MessageBus.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_MessageBus *MessageBusFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *MessageBusRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusRoleRevoked)
				if err := _MessageBus.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseRoleRevoked(log types.Log) (*MessageBusRoleRevoked, error) {
	event := new(MessageBusRoleRevoked)
	if err := _MessageBus.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the MessageBus contract.
type MessageBusUnpausedIterator struct {
	Event *MessageBusUnpaused // Event containing the contract specifics and raw log

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
func (it *MessageBusUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusUnpaused)
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
		it.Event = new(MessageBusUnpaused)
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
func (it *MessageBusUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusUnpaused represents a Unpaused event raised by the MessageBus contract.
type MessageBusUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MessageBus *MessageBusFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MessageBusUnpausedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MessageBusUnpausedIterator{contract: _MessageBus.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MessageBus *MessageBusFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MessageBusUnpaused) (event.Subscription, error) {

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusUnpaused)
				if err := _MessageBus.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MessageBus *MessageBusFilterer) ParseUnpaused(log types.Log) (*MessageBusUnpaused, error) {
	event := new(MessageBusUnpaused)
	if err := _MessageBus.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
