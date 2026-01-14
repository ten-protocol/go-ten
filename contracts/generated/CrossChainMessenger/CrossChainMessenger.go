// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package CrossChainMessenger

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

// CrossChainMessengerMetaData contains all meta data concerning the CrossChainMessenger contract.
var CrossChainMessengerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"error\",\"type\":\"bytes\"}],\"name\":\"CallFailed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"dataLen\",\"type\":\"uint256\"}],\"name\":\"DEBUG_BeforeCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"name\":\"DEBUG_CallResult\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DEBUG_ConsumeMessagePassed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UNPAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainTarget\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"encodeCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantPauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantUnpauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messageBusAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContract\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"messageConsumed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"messageConsumed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"relayMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"relayMessageWithProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokePauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeUnpauserRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"multisig\",\"type\":\"address\"}],\"name\":\"transferUnpauserRoleToMultisig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50611d198061001c5f395ff3fe608060405234801561000f575f5ffd5b50600436106101a5575f3560e01c80636c11c21c116100e8578063a217fddf11610093578063d547741f1161006e578063d547741f14610409578063e63ab1e91461041c578063f865af0814610443578063fb1bb9de14610456575f5ffd5b8063a217fddf146103dc578063b859ce83146103e3578063c4d66de8146103f6575f5ffd5b80638456cb59116100c35780638456cb591461036d57806391d1485414610375578063a1a227fa146103cc575f5ffd5b80636c11c21c14610328578063772c65521461033b5780637920c9861461035a575f5ffd5b80633f4ba83a11610153578063530c1e401161012e578063530c1e401461029c5780635b76f28b146102be5780635c975abb146102de57806363012de514610308575f5ffd5b80633f4ba83a1461026e5780634c81bd20146102765780635067627214610289575f5ffd5b80632f2ff15d116101835780632f2ff15d14610235578063329687821461024857806336568abe1461025b575f5ffd5b806301ffc9a7146101a9578063248a9ca3146101d25780632540e2da14610220575b5f5ffd5b6101bc6101b7366004611381565b61047d565b6040516101c991906113a8565b60405180910390f35b6102136101e03660046113c7565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b6040516101c991906113ea565b61023361022e36600461141c565b610515565b005b610233610243366004611439565b61054e565b61023361025636600461141c565b610597565b610233610269366004611439565b6105cb565b610233610617565b610233610284366004611488565b61064c565b610233610297366004611516565b610797565b6101bc6102aa3660046113c7565b60036020525f908152604090205460ff1681565b6102d16102cc3660046115e0565b6108de565b6040516101c99190611671565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff166101bc565b60015461031b906001600160a01b031681565b6040516101c9919061168b565b61023361033636600461141c565b61095e565b5f5461034d906001600160a01b031681565b6040516101c991906116b6565b61023361036836600461141c565b610992565b610233610a17565b6101bc610383366004611439565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b5f546001600160a01b031661031b565b6102135f81565b60025461031b906001600160a01b031681565b61023361040436600461141c565b610a49565b610233610417366004611439565b610ba8565b6102137f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a81565b61023361045136600461141c565b610beb565b6102137f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a81565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061050f57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f61051f81610c1f565b6105497f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a83610c29565b505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015461058781610c1f565b6105918383610cd6565b50505050565b5f6105a181610c1f565b6105497f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a83610cd6565b6001600160a01b038116331461060d576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105498282610c29565b7f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a61064181610c1f565b610649610d99565b50565b610654610e05565b61065c610e63565b61066581610ec6565b610672602082018261141c565b600180546001600160a01b0319166001600160a01b03929092169190911790555f6106a060808301836116c4565b8101906106ad9190611887565b8051600280546001600160a01b0319166001600160a01b0390921691821790559091505f9081905a84602001516040516106e791906118e0565b5f604051808303815f8787f1925050503d805f8114610721576040519150601f19603f3d011682016040523d82523d5f602084013e610726565b606091505b509150915081610754578060405163a5fa8d2b60e01b815260040161074b9190611671565b60405180910390fd5b5050600180546001600160a01b031990811682556002805490911690557f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00555050565b61079f610e05565b6107a7610e63565b6107b384848484611003565b6107c0602085018561141c565b600180546001600160a01b0319166001600160a01b03929092169190911790555f6107ee60808601866116c4565b8101906107fb9190611887565b8051600280546001600160a01b0319166001600160a01b0390921691821790559091505f9081905a846020015160405161083591906118e0565b5f604051808303815f8787f1925050503d805f811461086f576040519150601f19603f3d011682016040523d82523d5f602084013e610874565b606091505b509150915081610899578060405163a5fa8d2b60e01b815260040161074b9190611671565b5050600180546001600160a01b031990811682556002805490911690557f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005550610591565b60606040518060600160405280856001600160a01b0316815260200184848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920182905250938552505050602091820152604051610945929101611932565b60405160208183030381529060405290505b9392505050565b5f61096881610c1f565b6105497f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a83610cd6565b5f61099c81610c1f565b6001600160a01b0382166109c25760405162461bcd60e51b815260040161074b90611977565b6109ec7f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a33610c29565b506105497f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a83610cd6565b7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a610a4181610c1f565b6106496110f4565b5f610a5261114f565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f81158015610a7e5750825b90505f8267ffffffffffffffff166001148015610a9a5750303b155b905081158015610aa8575080155b15610adf576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610b1357845468ff00000000000000001916680100000000000000001785555b610b1b611177565b6001600160a01b038616610b2d575f5ffd5b610b3633611187565b5f80546001600160a01b0319166001600160a01b0388161790558315610ba057845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610b97906001906119a1565b60405180910390a15b505050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610be181610c1f565b6105918383610c29565b5f610bf581610c1f565b6105497f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a83610c29565b6106498133611203565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff1615610ccd575f848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4600191505061050f565b5f91505061050f565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff16610ccd575f848152602082815260408083206001600160a01b03871684529091529020805460ff19166001179055610d4f3390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4600191505061050f565b610da1611281565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300805460ff191681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b604051610dfa919061168b565b60405180910390a150565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1615610e61576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00805460011901610ec0576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60029055565b5f546040517f91643fdd0000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906391643fdd90610f0e908490600401611b6c565b602060405180830381865afa158015610f29573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f4d9190611b90565b610f695760405162461bcd60e51b815260040161074b90611bdf565b5f81604051602001610f7b9190611b6c565b60408051601f1981840301815291815281516020928301205f818152600390935291205490915060ff1615610fc25760405162461bcd60e51b815260040161074b90611c21565b5f908152600360205260409020805460ff1916600117905550565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5f546040517fce0d7db30000000000000000000000000000000000000000000000000000000081526001600160a01b039091169063ce0d7db390611051908790879087908790600401611c8a565b5f6040518083038186803b158015611067575f5ffd5b505afa158015611079573d5f5f3e3d5ffd5b505050505f8460405160200161108f9190611b6c565b60408051601f1981840301815291815281516020928301205f818152600390935291205490915060ff16156110d65760405162461bcd60e51b815260040161074b90611c21565b5f908152600360205260409020805460ff1916600117905550505050565b6110fc610e05565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300805460ff191660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25833610ded565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061050f565b61117f6112dc565b610e6161131a565b61118f6112dc565b611197611322565b61119f611322565b6111a95f82610cd6565b506111d47f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a82610cd6565b506111ff7f427da25fe773164f88948d3e215c94b6554e2ed5e5f203a821c9f2f6131cf75a82610cd6565b5050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff166111ff5780826040517fe2517d3f00000000000000000000000000000000000000000000000000000000815260040161074b929190611cc8565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610e61576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6112e461132a565b610e61576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610fdd6112dc565b610e616112dc565b5f61133361114f565b5468010000000000000000900460ff16919050565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b8114610649575f5ffd5b803561050f81611348565b5f60208284031215611394576113945f5ffd5b6109578383611376565b8015155b82525050565b6020810161050f828461139e565b8061136c565b803561050f816113b6565b5f602082840312156113da576113da5f5ffd5b61095783836113bc565b806113a2565b6020810161050f82846113e4565b5f6001600160a01b03821661050f565b61136c816113f8565b803561050f81611408565b5f6020828403121561142f5761142f5f5ffd5b6109578383611411565b5f5f6040838503121561144d5761144d5f5ffd5b61145784846113bc565b91506114668460208501611411565b90509250929050565b5f60c08284031215611482576114825f5ffd5b50919050565b5f6020828403121561149b5761149b5f5ffd5b813567ffffffffffffffff8111156114b4576114b45f5ffd5b6114c08482850161146f565b949350505050565b5f5f83601f8401126114db576114db5f5ffd5b50813567ffffffffffffffff8111156114f5576114f55f5ffd5b60208301915083602082028301111561150f5761150f5f5ffd5b9250929050565b5f5f5f5f6060858703121561152c5761152c5f5ffd5b843567ffffffffffffffff811115611545576115455f5ffd5b6115518782880161146f565b945050602085013567ffffffffffffffff811115611570576115705f5ffd5b61157c878288016114c8565b935093505061158e86604087016113bc565b905092959194509250565b5f5f83601f8401126115ac576115ac5f5ffd5b50813567ffffffffffffffff8111156115c6576115c65f5ffd5b60208301915083600182028301111561150f5761150f5f5ffd5b5f5f5f604084860312156115f5576115f55f5ffd5b6115ff8585611411565b9250602084013567ffffffffffffffff81111561161d5761161d5f5ffd5b61162986828701611599565b92509250509250925092565b8281835e505f910152565b5f611649825190565b808452602084019350611660818560208601611635565b601f01601f19169290920192915050565b602080825281016109578184611640565b6113a2816113f8565b6020810161050f8284611682565b5f61050f826113f8565b5f61050f82611699565b6113a2816116a3565b6020810161050f82846116ad565b5f808335601e19368590030181126116dd576116dd5f5ffd5b8301915050803567ffffffffffffffff8111156116fb576116fb5f5ffd5b60208201915060018102360382131561150f5761150f5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff8211171561176857611768611715565b6040525050565b5f61177960405190565b90506117858282611742565b919050565b5f67ffffffffffffffff8211156117a3576117a3611715565b601f19601f83011660200192915050565b82818337505f910152565b5f6117d16117cc8461178a565b61176f565b90508281528383830111156117e7576117e75f5ffd5b6109578360208301846117b4565b5f82601f830112611807576118075f5ffd5b610957838335602085016117bf565b5f60608284031215611829576118295f5ffd5b611833606061176f565b905061183f8383611411565b8152602082013567ffffffffffffffff81111561185d5761185d5f5ffd5b611869848285016117f5565b60208301525061187c83604084016113bc565b604082015292915050565b5f6020828403121561189a5761189a5f5ffd5b813567ffffffffffffffff8111156118b3576118b35f5ffd5b6114c084828501611816565b5f6118c8825190565b6118d6818560208601611635565b9290920192915050565b61050f81836118bf565b80515f9060608401906118fd8582611682565b50602083015184820360208601526119158282611640565b915050604083015161192a60408601826113e4565b509392505050565b6020808252810161095781846118ea565b60188152602081017f496e76616c6964206d756c746973696720616464726573730000000000000000815290505b60200190565b6020808252810161050f81611943565b5f67ffffffffffffffff821661050f565b6113a281611987565b6020810161050f8284611998565b505f61050f6020830183611411565b67ffffffffffffffff811661136c565b803561050f816119be565b505f61050f60208301836119ce565b67ffffffffffffffff81166113a2565b63ffffffff811661136c565b803561050f816119f8565b505f61050f6020830183611a04565b63ffffffff81166113a2565b5f808335601e1936859003018112611a4357611a435f5ffd5b830160208101925035905067ffffffffffffffff811115611a6557611a655f5ffd5b3681900382131561150f5761150f5f5ffd5b818352602083019250611a8b8284836117b4565b50601f01601f19160190565b60ff811661136c565b803561050f81611a97565b505f61050f6020830183611aa0565b60ff81166113a2565b5f60c08301611ad283806119af565b611adc8582611682565b50611aea60208401846119d9565b611af760208601826119e8565b50611b0560408401846119d9565b611b1260408601826119e8565b50611b206060840184611a0f565b611b2d6060860182611a1e565b50611b3b6080840184611a2a565b8583036080870152611b4e838284611a77565b92505050611b5f60a0840184611aab565b61192a60a0860182611aba565b602080825281016109578184611ac3565b80151561136c565b805161050f81611b7d565b5f60208284031215611ba357611ba35f5ffd5b6109578383611b85565b601f8152602081017f4d657373616765206e6f7420666f756e64206f722066696e616c697a65642e0081529050611971565b6020808252810161050f81611bad565b60198152602081017f4d65737361676520616c726561647920636f6e73756d65642e0000000000000081529050611971565b6020808252810161050f81611bef565b82818337505050565b8183526020830192505f7f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115611c7357611c735f5ffd5b602083029250611c84838584611c31565b50500190565b60608082528101611c9b8187611ac3565b90508181036020830152611cb0818587611c3a565b9050611cbf60408301846113e4565b95945050505050565b60408101611cd68285611682565b61095760208301846113e456fea264697066735822122008725d4cab005562091bb4a48c93253e29a30e6a4271420a614f8da39116553364736f6c634300081c0033",
}

// CrossChainMessengerABI is the input ABI used to generate the binding from.
// Deprecated: Use CrossChainMessengerMetaData.ABI instead.
var CrossChainMessengerABI = CrossChainMessengerMetaData.ABI

// CrossChainMessengerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrossChainMessengerMetaData.Bin instead.
var CrossChainMessengerBin = CrossChainMessengerMetaData.Bin

// DeployCrossChainMessenger deploys a new Ethereum contract, binding an instance of CrossChainMessenger to it.
func DeployCrossChainMessenger(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrossChainMessenger, error) {
	parsed, err := CrossChainMessengerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainMessengerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChainMessenger{CrossChainMessengerCaller: CrossChainMessengerCaller{contract: contract}, CrossChainMessengerTransactor: CrossChainMessengerTransactor{contract: contract}, CrossChainMessengerFilterer: CrossChainMessengerFilterer{contract: contract}}, nil
}

// CrossChainMessenger is an auto generated Go binding around an Ethereum contract.
type CrossChainMessenger struct {
	CrossChainMessengerCaller     // Read-only binding to the contract
	CrossChainMessengerTransactor // Write-only binding to the contract
	CrossChainMessengerFilterer   // Log filterer for contract events
}

// CrossChainMessengerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrossChainMessengerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrossChainMessengerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrossChainMessengerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainMessengerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrossChainMessengerSession struct {
	Contract     *CrossChainMessenger // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CrossChainMessengerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrossChainMessengerCallerSession struct {
	Contract *CrossChainMessengerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CrossChainMessengerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrossChainMessengerTransactorSession struct {
	Contract     *CrossChainMessengerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CrossChainMessengerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrossChainMessengerRaw struct {
	Contract *CrossChainMessenger // Generic contract binding to access the raw methods on
}

// CrossChainMessengerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrossChainMessengerCallerRaw struct {
	Contract *CrossChainMessengerCaller // Generic read-only contract binding to access the raw methods on
}

// CrossChainMessengerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrossChainMessengerTransactorRaw struct {
	Contract *CrossChainMessengerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrossChainMessenger creates a new instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessenger(address common.Address, backend bind.ContractBackend) (*CrossChainMessenger, error) {
	contract, err := bindCrossChainMessenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessenger{CrossChainMessengerCaller: CrossChainMessengerCaller{contract: contract}, CrossChainMessengerTransactor: CrossChainMessengerTransactor{contract: contract}, CrossChainMessengerFilterer: CrossChainMessengerFilterer{contract: contract}}, nil
}

// NewCrossChainMessengerCaller creates a new read-only instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerCaller(address common.Address, caller bind.ContractCaller) (*CrossChainMessengerCaller, error) {
	contract, err := bindCrossChainMessenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerCaller{contract: contract}, nil
}

// NewCrossChainMessengerTransactor creates a new write-only instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainMessengerTransactor, error) {
	contract, err := bindCrossChainMessenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerTransactor{contract: contract}, nil
}

// NewCrossChainMessengerFilterer creates a new log filterer instance of CrossChainMessenger, bound to a specific deployed contract.
func NewCrossChainMessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainMessengerFilterer, error) {
	contract, err := bindCrossChainMessenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerFilterer{contract: contract}, nil
}

// bindCrossChainMessenger binds a generic wrapper to an already deployed contract.
func bindCrossChainMessenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrossChainMessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainMessenger *CrossChainMessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainMessenger.Contract.CrossChainMessengerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainMessenger *CrossChainMessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.CrossChainMessengerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainMessenger *CrossChainMessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.CrossChainMessengerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainMessenger *CrossChainMessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainMessenger.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainMessenger *CrossChainMessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainMessenger *CrossChainMessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _CrossChainMessenger.Contract.DEFAULTADMINROLE(&_CrossChainMessenger.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _CrossChainMessenger.Contract.DEFAULTADMINROLE(&_CrossChainMessenger.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerSession) PAUSERROLE() ([32]byte, error) {
	return _CrossChainMessenger.Contract.PAUSERROLE(&_CrossChainMessenger.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) PAUSERROLE() ([32]byte, error) {
	return _CrossChainMessenger.Contract.PAUSERROLE(&_CrossChainMessenger.CallOpts)
}

// UNPAUSERROLE is a free data retrieval call binding the contract method 0xfb1bb9de.
//
// Solidity: function UNPAUSER_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCaller) UNPAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "UNPAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UNPAUSERROLE is a free data retrieval call binding the contract method 0xfb1bb9de.
//
// Solidity: function UNPAUSER_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerSession) UNPAUSERROLE() ([32]byte, error) {
	return _CrossChainMessenger.Contract.UNPAUSERROLE(&_CrossChainMessenger.CallOpts)
}

// UNPAUSERROLE is a free data retrieval call binding the contract method 0xfb1bb9de.
//
// Solidity: function UNPAUSER_ROLE() view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) UNPAUSERROLE() ([32]byte, error) {
	return _CrossChainMessenger.Contract.UNPAUSERROLE(&_CrossChainMessenger.CallOpts)
}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) CrossChainSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "crossChainSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) CrossChainSender() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainSender(&_CrossChainMessenger.CallOpts)
}

// CrossChainSender is a free data retrieval call binding the contract method 0x63012de5.
//
// Solidity: function crossChainSender() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) CrossChainSender() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainSender(&_CrossChainMessenger.CallOpts)
}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) CrossChainTarget(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "crossChainTarget")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) CrossChainTarget() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainTarget(&_CrossChainMessenger.CallOpts)
}

// CrossChainTarget is a free data retrieval call binding the contract method 0xb859ce83.
//
// Solidity: function crossChainTarget() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) CrossChainTarget() (common.Address, error) {
	return _CrossChainMessenger.Contract.CrossChainTarget(&_CrossChainMessenger.CallOpts)
}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerCaller) EncodeCall(opts *bind.CallOpts, target common.Address, payload []byte) ([]byte, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "encodeCall", target, payload)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerSession) EncodeCall(target common.Address, payload []byte) ([]byte, error) {
	return _CrossChainMessenger.Contract.EncodeCall(&_CrossChainMessenger.CallOpts, target, payload)
}

// EncodeCall is a free data retrieval call binding the contract method 0x5b76f28b.
//
// Solidity: function encodeCall(address target, bytes payload) pure returns(bytes)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) EncodeCall(target common.Address, payload []byte) ([]byte, error) {
	return _CrossChainMessenger.Contract.EncodeCall(&_CrossChainMessenger.CallOpts, target, payload)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _CrossChainMessenger.Contract.GetRoleAdmin(&_CrossChainMessenger.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _CrossChainMessenger.Contract.GetRoleAdmin(&_CrossChainMessenger.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _CrossChainMessenger.Contract.HasRole(&_CrossChainMessenger.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _CrossChainMessenger.Contract.HasRole(&_CrossChainMessenger.CallOpts, role, account)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageBus() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBus(&_CrossChainMessenger.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageBus() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBus(&_CrossChainMessenger.CallOpts)
}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageBusContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageBusContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageBusContract() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBusContract(&_CrossChainMessenger.CallOpts)
}

// MessageBusContract is a free data retrieval call binding the contract method 0x772c6552.
//
// Solidity: function messageBusContract() view returns(address)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageBusContract() (common.Address, error) {
	return _CrossChainMessenger.Contract.MessageBusContract(&_CrossChainMessenger.CallOpts)
}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerCaller) MessageConsumed(opts *bind.CallOpts, messageHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "messageConsumed", messageHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerSession) MessageConsumed(messageHash [32]byte) (bool, error) {
	return _CrossChainMessenger.Contract.MessageConsumed(&_CrossChainMessenger.CallOpts, messageHash)
}

// MessageConsumed is a free data retrieval call binding the contract method 0x530c1e40.
//
// Solidity: function messageConsumed(bytes32 messageHash) view returns(bool messageConsumed)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) MessageConsumed(messageHash [32]byte) (bool, error) {
	return _CrossChainMessenger.Contract.MessageConsumed(&_CrossChainMessenger.CallOpts, messageHash)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerSession) Paused() (bool, error) {
	return _CrossChainMessenger.Contract.Paused(&_CrossChainMessenger.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) Paused() (bool, error) {
	return _CrossChainMessenger.Contract.Paused(&_CrossChainMessenger.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CrossChainMessenger.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CrossChainMessenger.Contract.SupportsInterface(&_CrossChainMessenger.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CrossChainMessenger *CrossChainMessengerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CrossChainMessenger.Contract.SupportsInterface(&_CrossChainMessenger.CallOpts, interfaceId)
}

// GrantPauserRole is a paid mutator transaction binding the contract method 0x6c11c21c.
//
// Solidity: function grantPauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) GrantPauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "grantPauserRole", account)
}

// GrantPauserRole is a paid mutator transaction binding the contract method 0x6c11c21c.
//
// Solidity: function grantPauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) GrantPauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.GrantPauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// GrantPauserRole is a paid mutator transaction binding the contract method 0x6c11c21c.
//
// Solidity: function grantPauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) GrantPauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.GrantPauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.GrantRole(&_CrossChainMessenger.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.GrantRole(&_CrossChainMessenger.TransactOpts, role, account)
}

// GrantUnpauserRole is a paid mutator transaction binding the contract method 0x32968782.
//
// Solidity: function grantUnpauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) GrantUnpauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "grantUnpauserRole", account)
}

// GrantUnpauserRole is a paid mutator transaction binding the contract method 0x32968782.
//
// Solidity: function grantUnpauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) GrantUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.GrantUnpauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// GrantUnpauserRole is a paid mutator transaction binding the contract method 0x32968782.
//
// Solidity: function grantUnpauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) GrantUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.GrantUnpauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) Initialize(opts *bind.TransactOpts, messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "initialize", messageBusAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) Initialize(messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Initialize(&_CrossChainMessenger.TransactOpts, messageBusAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messageBusAddr) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) Initialize(messageBusAddr common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Initialize(&_CrossChainMessenger.TransactOpts, messageBusAddr)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CrossChainMessenger *CrossChainMessengerSession) Pause() (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Pause(&_CrossChainMessenger.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) Pause() (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Pause(&_CrossChainMessenger.TransactOpts)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RelayMessage(opts *bind.TransactOpts, message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "relayMessage", message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}

// RelayMessage is a paid mutator transaction binding the contract method 0x4c81bd20.
//
// Solidity: function relayMessage((address,uint64,uint64,uint32,bytes,uint8) message) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RelayMessage(message StructsCrossChainMessage) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessage(&_CrossChainMessenger.TransactOpts, message)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RelayMessageWithProof(opts *bind.TransactOpts, message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "relayMessageWithProof", message, proof, root)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RelayMessageWithProof(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessageWithProof(&_CrossChainMessenger.TransactOpts, message, proof, root)
}

// RelayMessageWithProof is a paid mutator transaction binding the contract method 0x50676272.
//
// Solidity: function relayMessageWithProof((address,uint64,uint64,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RelayMessageWithProof(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RelayMessageWithProof(&_CrossChainMessenger.TransactOpts, message, proof, root)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RenounceRole(&_CrossChainMessenger.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RenounceRole(&_CrossChainMessenger.TransactOpts, role, callerConfirmation)
}

// RevokePauserRole is a paid mutator transaction binding the contract method 0xf865af08.
//
// Solidity: function revokePauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RevokePauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "revokePauserRole", account)
}

// RevokePauserRole is a paid mutator transaction binding the contract method 0xf865af08.
//
// Solidity: function revokePauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RevokePauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RevokePauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// RevokePauserRole is a paid mutator transaction binding the contract method 0xf865af08.
//
// Solidity: function revokePauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RevokePauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RevokePauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RevokeRole(&_CrossChainMessenger.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RevokeRole(&_CrossChainMessenger.TransactOpts, role, account)
}

// RevokeUnpauserRole is a paid mutator transaction binding the contract method 0x2540e2da.
//
// Solidity: function revokeUnpauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) RevokeUnpauserRole(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "revokeUnpauserRole", account)
}

// RevokeUnpauserRole is a paid mutator transaction binding the contract method 0x2540e2da.
//
// Solidity: function revokeUnpauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) RevokeUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RevokeUnpauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// RevokeUnpauserRole is a paid mutator transaction binding the contract method 0x2540e2da.
//
// Solidity: function revokeUnpauserRole(address account) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) RevokeUnpauserRole(account common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.RevokeUnpauserRole(&_CrossChainMessenger.TransactOpts, account)
}

// TransferUnpauserRoleToMultisig is a paid mutator transaction binding the contract method 0x7920c986.
//
// Solidity: function transferUnpauserRoleToMultisig(address multisig) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) TransferUnpauserRoleToMultisig(opts *bind.TransactOpts, multisig common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "transferUnpauserRoleToMultisig", multisig)
}

// TransferUnpauserRoleToMultisig is a paid mutator transaction binding the contract method 0x7920c986.
//
// Solidity: function transferUnpauserRoleToMultisig(address multisig) returns()
func (_CrossChainMessenger *CrossChainMessengerSession) TransferUnpauserRoleToMultisig(multisig common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.TransferUnpauserRoleToMultisig(&_CrossChainMessenger.TransactOpts, multisig)
}

// TransferUnpauserRoleToMultisig is a paid mutator transaction binding the contract method 0x7920c986.
//
// Solidity: function transferUnpauserRoleToMultisig(address multisig) returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) TransferUnpauserRoleToMultisig(multisig common.Address) (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.TransferUnpauserRoleToMultisig(&_CrossChainMessenger.TransactOpts, multisig)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CrossChainMessenger *CrossChainMessengerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainMessenger.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CrossChainMessenger *CrossChainMessengerSession) Unpause() (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Unpause(&_CrossChainMessenger.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CrossChainMessenger *CrossChainMessengerTransactorSession) Unpause() (*types.Transaction, error) {
	return _CrossChainMessenger.Contract.Unpause(&_CrossChainMessenger.TransactOpts)
}

// CrossChainMessengerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CrossChainMessenger contract.
type CrossChainMessengerInitializedIterator struct {
	Event *CrossChainMessengerInitialized // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerInitialized)
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
		it.Event = new(CrossChainMessengerInitialized)
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
func (it *CrossChainMessengerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerInitialized represents a Initialized event raised by the CrossChainMessenger contract.
type CrossChainMessengerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterInitialized(opts *bind.FilterOpts) (*CrossChainMessengerInitializedIterator, error) {

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerInitializedIterator{contract: _CrossChainMessenger.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerInitialized) (event.Subscription, error) {

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerInitialized)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParseInitialized(log types.Log) (*CrossChainMessengerInitialized, error) {
	event := new(CrossChainMessengerInitialized)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainMessengerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the CrossChainMessenger contract.
type CrossChainMessengerPausedIterator struct {
	Event *CrossChainMessengerPaused // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerPaused)
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
		it.Event = new(CrossChainMessengerPaused)
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
func (it *CrossChainMessengerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerPaused represents a Paused event raised by the CrossChainMessenger contract.
type CrossChainMessengerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterPaused(opts *bind.FilterOpts) (*CrossChainMessengerPausedIterator, error) {

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerPausedIterator{contract: _CrossChainMessenger.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerPaused) (event.Subscription, error) {

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerPaused)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParsePaused(log types.Log) (*CrossChainMessengerPaused, error) {
	event := new(CrossChainMessengerPaused)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainMessengerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the CrossChainMessenger contract.
type CrossChainMessengerRoleAdminChangedIterator struct {
	Event *CrossChainMessengerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerRoleAdminChanged)
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
		it.Event = new(CrossChainMessengerRoleAdminChanged)
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
func (it *CrossChainMessengerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerRoleAdminChanged represents a RoleAdminChanged event raised by the CrossChainMessenger contract.
type CrossChainMessengerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*CrossChainMessengerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerRoleAdminChangedIterator{contract: _CrossChainMessenger.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerRoleAdminChanged)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParseRoleAdminChanged(log types.Log) (*CrossChainMessengerRoleAdminChanged, error) {
	event := new(CrossChainMessengerRoleAdminChanged)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainMessengerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the CrossChainMessenger contract.
type CrossChainMessengerRoleGrantedIterator struct {
	Event *CrossChainMessengerRoleGranted // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerRoleGranted)
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
		it.Event = new(CrossChainMessengerRoleGranted)
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
func (it *CrossChainMessengerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerRoleGranted represents a RoleGranted event raised by the CrossChainMessenger contract.
type CrossChainMessengerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CrossChainMessengerRoleGrantedIterator, error) {

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

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerRoleGrantedIterator{contract: _CrossChainMessenger.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerRoleGranted)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParseRoleGranted(log types.Log) (*CrossChainMessengerRoleGranted, error) {
	event := new(CrossChainMessengerRoleGranted)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainMessengerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the CrossChainMessenger contract.
type CrossChainMessengerRoleRevokedIterator struct {
	Event *CrossChainMessengerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerRoleRevoked)
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
		it.Event = new(CrossChainMessengerRoleRevoked)
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
func (it *CrossChainMessengerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerRoleRevoked represents a RoleRevoked event raised by the CrossChainMessenger contract.
type CrossChainMessengerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CrossChainMessengerRoleRevokedIterator, error) {

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

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerRoleRevokedIterator{contract: _CrossChainMessenger.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerRoleRevoked)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParseRoleRevoked(log types.Log) (*CrossChainMessengerRoleRevoked, error) {
	event := new(CrossChainMessengerRoleRevoked)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrossChainMessengerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the CrossChainMessenger contract.
type CrossChainMessengerUnpausedIterator struct {
	Event *CrossChainMessengerUnpaused // Event containing the contract specifics and raw log

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
func (it *CrossChainMessengerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrossChainMessengerUnpaused)
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
		it.Event = new(CrossChainMessengerUnpaused)
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
func (it *CrossChainMessengerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrossChainMessengerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrossChainMessengerUnpaused represents a Unpaused event raised by the CrossChainMessenger contract.
type CrossChainMessengerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CrossChainMessenger *CrossChainMessengerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*CrossChainMessengerUnpausedIterator, error) {

	logs, sub, err := _CrossChainMessenger.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &CrossChainMessengerUnpausedIterator{contract: _CrossChainMessenger.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_CrossChainMessenger *CrossChainMessengerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CrossChainMessengerUnpaused) (event.Subscription, error) {

	logs, sub, err := _CrossChainMessenger.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrossChainMessengerUnpaused)
				if err := _CrossChainMessenger.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_CrossChainMessenger *CrossChainMessengerFilterer) ParseUnpaused(log types.Log) (*CrossChainMessengerUnpaused, error) {
	event := new(CrossChainMessengerUnpaused)
	if err := _CrossChainMessenger.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
