// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TenBridgeTestnet

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

// TenBridgeTestnetMetaData contains all meta data concerning the TenBridgeTestnet contract.
var TenBridgeTestnetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERC20_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NATIVE_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUSPENDED_ERC20_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"contractIMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messenger\",\"outputs\":[{\"internalType\":\"contractICrossChainMessenger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"pauseToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"promoteToAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"recoverTestnetFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"remoteBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"setRemoteBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"unpauseToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"whitelistToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506129448061001c5f395ff3fe60806040526004361061018e575f3560e01c80635d872970116100dc578063a1a227fa11610087578063affed0e011610062578063affed0e0146104cf578063c091a58b1461055a578063d547741f14610579578063e4c3ebc714610598575f5ffd5b8063a1a227fa14610495578063a217fddf146104a9578063a381c8e2146104bc575f5ffd5b806383bece4d116100b757806383bece4d146103f457806391d148541461041357806393b3744214610476575f5ffd5b80635d8729701461036f57806375b238fc146103a25780637c41ad2c146103d5575f5ffd5b806336568abe1161013c578063485cc95511610117578063485cc955146102fe578063498d82ab1461031d5780635ccc96131461033c575f5ffd5b806336568abe1461029f5780633b3bff0f146102be5780633cb747bf146102dd575f5ffd5b80631888d7121161016c5780631888d71214610213578063248a9ca3146102265780632f2ff15d14610280575f5ffd5b806301ffc9a7146101925780630f0a9a4b146101c757806316ce8149146101f2575b5f5ffd5b34801561019d575f5ffd5b506101b16101ac366004611d2f565b6105cb565b6040516101be9190611d5d565b60405180910390f35b3480156101d2575f5ffd5b505f546101e5906001600160a01b031681565b6040516101be9190611d84565b3480156101fd575f5ffd5b5061021161020c366004611da6565b610633565b005b610211610221366004611da6565b6106e3565b348015610231575f5ffd5b50610273610240366004611dd4565b5f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b6040516101be9190611df7565b34801561028b575f5ffd5b5061021161029a366004611e05565b610754565b3480156102aa575f5ffd5b506102116102b9366004611e05565b61079d565b3480156102c9575f5ffd5b506102116102d8366004611da6565b6107ee565b3480156102e8575f5ffd5b506102f1610842565b6040516101be9190611e58565b348015610309575f5ffd5b50610211610318366004611e66565b6108cd565b348015610328575f5ffd5b50610211610337366004611ed2565b610ae3565b348015610347575f5ffd5b506102737fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e81565b34801561037a575f5ffd5b506102737f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a81565b3480156103ad575f5ffd5b506102737fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b3480156103e0575f5ffd5b506102116103ef366004611da6565b610c27565b3480156103ff575f5ffd5b5061021161040e366004611f59565b610cd1565b34801561041e575f5ffd5b506101b161042d366004611e05565b5f9182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b348015610481575f5ffd5b50610211610490366004611da6565b610e26565b3480156104a0575f5ffd5b506102f1610e80565b3480156104b4575f5ffd5b506102735f81565b6102116104ca366004611f59565b610f00565b3480156104da575f5ffd5b5060408051808201909152601a81527f43726f7373436861696e456e61626c656454454e2e6e6f6e63650000000000006020909101527f896d106647b57d520a34062c7c0dde769b7551e327629f69d5e9844e20c864625f527fe7fbfe9855ab39eb2e984ddc0938c4134151cf94d8a54d4770b35b584ad4bd0054610273565b348015610565575f5ffd5b50610211610574366004611da6565b611043565b348015610584575f5ffd5b50610211610593366004611e05565b611188565b3480156105a3575f5ffd5b506102737fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad211057881565b5f6001600160e01b031982167f7965db0b00000000000000000000000000000000000000000000000000000000148061062d57507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561065d816111cb565b6001600160a01b03821661068c5760405162461bcd60e51b815260040161068390611fd3565b60405180910390fd5b5f546001600160a01b0316156106b45760405162461bcd60e51b81526004016106839061203d565b505f805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b5f34116107025760405162461bcd60e51b81526004016106839061207f565b5f6040518060400160405280348152602001836001600160a01b031681525060405160200161073191906120ae565b60408051601f1981840301815291905290506107508160025f5f6111d8565b5050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015461078d816111cb565b61079783836112e3565b50505050565b6001600160a01b03811633146107df576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107e982826113af565b505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610818816111cb565b6107e97fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e836113af565b60408051808201909152601e81527f43726f7373436861696e456e61626c656454454e2e6d657373656e67657200006020909101527fa8b5aada5c72138bb5566a3940e8fe06f59ef8af1e490446ba6ea7fa80395d525f9081527f3b49b3a570909bb4d324cb0ca029c61a2f4f7251edd27af783a6ad02851382005b546001600160a01b0316919050565b5f6108d6611453565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f811580156109025750825b90505f8267ffffffffffffffff16600114801561091e5750303b155b90508115801561092c575080155b15610963576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561099757845468ff00000000000000001916680100000000000000001785555b6001600160a01b0387166109bd5760405162461bcd60e51b8152600401610683906120ee565b6001600160a01b0386166109e35760405162461bcd60e51b815260040161068390612130565b6109ec8761147b565b6109f461167c565b6109fe5f876112e3565b50610a095f306112e3565b50610a347fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217755f611686565b610a5e7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775876112e3565b50610a897fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad21105785f6112e3565b508315610ada57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610ad19060019061215a565b60405180910390a15b50505050505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610b0d816111cb565b6001600160a01b038616610b335760405162461bcd60e51b81526004016106839061219a565b6001600160a01b0386165f9081527fe0305390dd8de2e924991bcde8c43652df4370e71f9558170e600f8cd2fe1d57602052604090205460ff1615610b8a5760405162461bcd60e51b8152600401610683906121dc565b610bb47f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a876112e3565b505f63458ffd6360e01b8787878787604051602401610bd7959493929190612217565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091525f54909150610ada906001600160a01b03168260015b5f5f5f611727565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610c51816111cb565b6001600160a01b0382165f9081527fe0305390dd8de2e924991bcde8c43652df4370e71f9558170e600f8cd2fe1d57602052604090205460ff16610ca75760405162461bcd60e51b81526004016106839061228a565b6107e97fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e836112e3565b5f546001600160a01b0316610ce46118a1565b610d005760405162461bcd60e51b8152600401610683906122f2565b806001600160a01b0316610d126118c3565b6001600160a01b031614610d385760405162461bcd60e51b81526004016106839061235a565b30610d41611930565b6001600160a01b031614610d675760405162461bcd60e51b8152600401610683906123c2565b610d6f611974565b6001600160a01b0384165f9081527fe0305390dd8de2e924991bcde8c43652df4370e71f9558170e600f8cd2fe1d57602052604090205460ff1615610dbe57610db98484846119fa565b610e1e565b6001600160a01b0384165f9081527e5fd0bb0e17815069821c0eac859eb66a4da90a93d511d999ef71402d667e27602052604090205460ff1615610e0657610db98284611a55565b60405162461bcd60e51b81526004016106839061242a565b610797611b07565b5f610e30816111cb565b6001600160a01b038216610e565760405162461bcd60e51b81526004016106839061246c565b6107507fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177583610754565b60408051808201909152601f81527f43726f7373436861696e456e61626c656454454e2e6d657373616765427573006020909101527f3e1bb302f668bd876eab4a48b3759a1d614a1ecbcc67ee27a10c9a116878004e5f9081527f6c6664e79adefe2c614a8e3c94fc27135b7678c3722965a80d01e330dd948d006108be565b6001600160a01b0383165f9081527ff8f9f0c07f8f13fae35355825022a620ece4ae820bcc59c97dd7358124668dc9602052604090205460ff1615610f575760405162461bcd60e51b8152600401610683906124ae565b5f8211610f765760405162461bcd60e51b8152600401610683906124f0565b6001600160a01b0383165f9081527fe0305390dd8de2e924991bcde8c43652df4370e71f9558170e600f8cd2fe1d57602052604090205460ff16610fcc5760405162461bcd60e51b815260040161068390612500565b610fd883333085611b31565b5f6383bece4d60e01b848484604051602401610ff693929190612587565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091525f8054919250610797916001600160a01b0316908390610c1f565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561106d816111cb565b6001600160a01b0382166110935760405162461bcd60e51b8152600401610683906125e9565b4662aa36a714806110a5575046610539145b6110c15760405162461bcd60e51b815260040161068390612651565b4780156107e9575f836001600160a01b0316826040515f6040518083038185875af1925050503d805f8114611111576040519150601f19603f3d011682016040523d82523d5f602084013e611116565b606091505b50509050806111375760405162461bcd60e51b815260040161068390612693565b5f6001600160a01b0316846001600160a01b03167f2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b63988460405161117a9190611df7565b60405180910390a350505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b62680060205260409020600101546111c1816111cb565b61079783836113af565b6111d58133611b8b565b50565b6111e0610e80565b60408051808201909152601a81527f43726f7373436861696e456e61626c656454454e2e6e6f6e63650000000000006020909101527f896d106647b57d520a34062c7c0dde769b7551e327629f69d5e9844e20c864625f526001600160a01b0316630d3fd67c837fe7fbfe9855ab39eb2e984ddc0938c4134151cf94d8a54d4770b35b584ad4bd008054905f611275836126d0565b919050558688866040518663ffffffff1660e01b815260040161129b9493929190612749565b60206040518083038185885af11580156112b7573d5f5f3e3d5ffd5b50505050506040513d601f19601f820116820180604052508101906112dc91906127a9565b5050505050565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff166113a6575f848152602082815260408083206001600160a01b03871684529091529020805460ff1916600117905561135c3390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4600191505061062d565b5f91505061062d565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff16156113a6575f848152602082815260408083206001600160a01b0387168085529252808320805460ff1916905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4600191505061062d565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061062d565b60408051808201909152601e81527f43726f7373436861696e456e61626c656454454e2e6d657373656e67657200006020909101527fa8b5aada5c72138bb5566a3940e8fe06f59ef8af1e490446ba6ea7fa80395d525f527f3b49b3a570909bb4d324cb0ca029c61a2f4f7251edd27af783a6ad0285138200805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316179055611522610842565b6001600160a01b031663a1a227fa6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561155d573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061158191906127d1565b604080518082018252601f81527f43726f7373436861696e456e61626c656454454e2e6d657373616765427573006020918201527f6c6664e79adefe2c614a8e3c94fc27135b7678c3722965a80d01e330dd948d00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0394909416939093179092558051808201909152601a81527f43726f7373436861696e456e61626c656454454e2e6e6f6e6365000000000000910152507f896d106647b57d520a34062c7c0dde769b7551e327629f69d5e9844e20c864625f9081527fe7fbfe9855ab39eb2e984ddc0938c4134151cf94d8a54d4770b35b584ad4bd0055565b611684611c09565b565b7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268005f6116df845f9081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b5f85815260208490526040808220600101869055519192508491839187917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a450505050565b6001600160a01b03861661174d5760405162461bcd60e51b815260040161068390612820565b5f6040518060600160405280886001600160a01b03168152602001878152602001858152506040516020016117829190612878565b604051602081830303815290604052905061179b610e80565b60408051808201909152601a81527f43726f7373436861696e456e61626c656454454e2e6e6f6e63650000000000006020909101527f896d106647b57d520a34062c7c0dde769b7551e327629f69d5e9844e20c864625f526001600160a01b0316630d3fd67c837fe7fbfe9855ab39eb2e984ddc0938c4134151cf94d8a54d4770b35b584ad4bd008054905f611830836126d0565b919050558885886040518663ffffffff1660e01b81526004016118569493929190612749565b60206040518083038185885af1158015611872573d5f5f3e3d5ffd5b50505050506040513d601f19601f8201168201806040525081019061189791906127a9565b5050505050505050565b5f6118aa610842565b6001600160a01b0316336001600160a01b031614905090565b5f6118cc610842565b6001600160a01b03166363012de56040518163ffffffff1660e01b8152600401602060405180830381865afa158015611907573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061192b91906127d1565b905090565b5f611939610842565b6001600160a01b031663b859ce836040518163ffffffff1660e01b8152600401602060405180830381865afa158015611907573d5f5f3e3d5ffd5b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c156119cd576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61168460017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005b90611c47565b611a05838284611c4e565b826001600160a01b0316816001600160a01b03167f2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b639884604051611a489190611df7565b60405180910390a3505050565b5f826001600160a01b0316826040515f6040518083038185875af1925050503d805f8114611a9e576040519150601f19603f3d011682016040523d82523d5f602084013e611aa3565b606091505b5050905080611ac45760405162461bcd60e51b8152600401610683906128bb565b5f6001600160a01b0316836001600160a01b03167f2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b639884604051611a489190611df7565b6116845f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f006119f4565b61079784856001600160a01b03166323b872dd868686604051602401611b59939291906128cb565b604051602081830303815290604052915060e01b6020820180516001600160e01b038381831617835250505050611c74565b5f8281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff166107505780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016106839291906128f3565b611c11611cf0565b611684576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80825d5050565b6107e983846001600160a01b031663a9059cbb8585604051602401611b599291906128f3565b5f5f60205f8451602086015f885af180611c93576040513d5f823e3d81fd5b50505f513d91508115611caa578060011415611cb7565b6001600160a01b0384163b155b1561079757836040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016106839190611d84565b5f611cf9611453565b5468010000000000000000900460ff16919050565b6001600160e01b031981165b81146111d5575f5ffd5b803561062d81611d0e565b5f60208284031215611d4257611d425f5ffd5b611d4c8383611d24565b9392505050565b8015155b82525050565b6020810161062d8284611d53565b5f6001600160a01b03821661062d565b611d5781611d6b565b6020810161062d8284611d7b565b611d1a81611d6b565b803561062d81611d92565b5f60208284031215611db957611db95f5ffd5b611d4c8383611d9b565b80611d1a565b803561062d81611dc3565b5f60208284031215611de757611de75f5ffd5b611d4c8383611dc9565b80611d57565b6020810161062d8284611df1565b5f5f60408385031215611e1957611e195f5ffd5b611e238484611dc9565b9150611e328460208501611d9b565b90509250929050565b5f61062d82611d6b565b5f61062d82611e3b565b611d5781611e45565b6020810161062d8284611e4f565b5f5f60408385031215611e7a57611e7a5f5ffd5b611e238484611d9b565b5f5f83601f840112611e9757611e975f5ffd5b50813567ffffffffffffffff811115611eb157611eb15f5ffd5b602083019150836001820283011115611ecb57611ecb5f5ffd5b9250929050565b5f5f5f5f5f60608688031215611ee957611ee95f5ffd5b611ef38787611d9b565b9450602086013567ffffffffffffffff811115611f1157611f115f5ffd5b611f1d88828901611e84565b9450945050604086013567ffffffffffffffff811115611f3e57611f3e5f5ffd5b611f4a88828901611e84565b92509250509295509295909350565b5f5f5f60608486031215611f6e57611f6e5f5ffd5b611f788585611d9b565b9250611f878560208601611dc9565b9150611f968560408601611d9b565b90509250925092565b60148152602081017f4272696467652063616e6e6f7420626520307830000000000000000000000000815290505b60200190565b6020808252810161062d81611f9f565b60228152602081017f52656d6f746520627269646765206164647265737320616c726561647920736581527f742e000000000000000000000000000000000000000000000000000000000000602082015290505b60400190565b6020808252810161062d81611fe3565b600f8152602081017f456d707479207472616e736665722e000000000000000000000000000000000081529050611fcd565b6020808252810161062d8161204d565b805161209b8382611df1565b5060208101516107e96020840182611d7b565b6040810161062d828461208f565b60178152602081017f4d657373656e6765722063616e6e6f742062652030783000000000000000000081529050611fcd565b6020808252810161062d816120bc565b60138152602081017f4f776e65722063616e6e6f74206265203078300000000000000000000000000081529050611fcd565b6020808252810161062d816120fe565b5f67ffffffffffffffff821661062d565b611d5781612140565b6020810161062d8284612151565b60138152602081017f41737365742063616e6e6f74206265203078300000000000000000000000000081529050611fcd565b6020808252810161062d81612168565b60198152602081017f546f6b656e20616c72656164792077686974656c69737465640000000000000081529050611fcd565b6020808252810161062d816121aa565b82818337505f910152565b81835260208301925061220b8284836121ec565b50601f01601f19160190565b606081016122258288611d7b565b81810360208301526122388186886121f7565b9050818103604083015261224d8184866121f7565b979650505050505050565b60188152602081017f546f6b656e206973206e6f742077686974656c6973746564000000000000000081529050611fcd565b6020808252810161062d81612258565b60308152602081017f436f6e74726163742063616c6c6572206973206e6f742074686520726567697381527f7465726564206d657373656e676572210000000000000000000000000000000060208201529050612037565b6020808252810161062d8161229a565b60318152602081017f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2081527f696e636f72726563742073656e6465722100000000000000000000000000000060208201529050612037565b6020808252810161062d81612302565b60318152602081017f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2081527f696e636f7272656374207461726765742100000000000000000000000000000060208201529050612037565b6020808252810161062d8161236a565b60258152602081017f417474656d7074696e6720746f20776974686472617720756e6b6e6f776e206181527f737365742e00000000000000000000000000000000000000000000000000000060208201529050612037565b6020808252810161062d816123d2565b60178152602081017f4e65772061646d696e2063616e6e6f742062652030783000000000000000000081529050611fcd565b6020808252810161062d8161243a565b60108152602081017f546f6b656e206973207061757365642e0000000000000000000000000000000081529050611fcd565b6020808252810161062d8161247c565b601a8152602081017f417474656d7074696e6720656d707479207472616e736665722e00000000000081529050611fcd565b6020808252810161062d816124be565b6020808252810161062d81604e81527f54686973206164647265737320686173206e6f74206265656e20676976656e2060208201527f61207479706520616e64206973207468757320636f6e73696465726564206e6f60408201527f742077686974656c69737465642e000000000000000000000000000000000000606082015260800190565b606081016125958286611d7b565b6125a26020830185611df1565b6125af6040830184611d7b565b949350505050565b60168152602081017f52656365697665722063616e6e6f74206265203078300000000000000000000081529050611fcd565b6020808252810161062d816125b7565b602a8152602081017f5265636f76657279206f6e6c7920616c6c6f776564206f6e20617070726f766581527f6420746573746e6574730000000000000000000000000000000000000000000060208201529050612037565b6020808252810161062d816125f9565b60178152602081017f4661696c656420746f207265636f76657220457468657200000000000000000081529050611fcd565b6020808252810161062d81612661565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f600182016126e1576126e16126a3565b5060010190565b67ffffffffffffffff8116611d57565b63ffffffff8116611d57565b8281835e505f910152565b5f612718825190565b80845260208401935061272f818560208601612704565b601f01601f19169290920192915050565b60ff8116611d57565b6080810161275782876126e8565b61276460208301866126f8565b8181036040830152612776818561270f565b90506127856060830184612740565b95945050505050565b67ffffffffffffffff8116611d1a565b805161062d8161278e565b5f602082840312156127bc576127bc5f5ffd5b611d4c838361279e565b805161062d81611d92565b5f602082840312156127e4576127e45f5ffd5b611d4c83836127c6565b60148152602081017f5461726765742063616e6e6f742062652030783000000000000000000000000081529050611fcd565b6020808252810161062d816127ee565b80515f9060608401906128438582611d7b565b506020830151848203602086015261285b828261270f565b91505060408301516128706040860182611df1565b509392505050565b60208082528101611d4c8184612830565b60148152602081017f4661696c656420746f2073656e6420457468657200000000000000000000000081529050611fcd565b6020808252810161062d81612889565b606081016128d98286611d7b565b6128e66020830185611d7b565b6125af6040830184611df1565b604081016129018285611d7b565b611d4c6020830184611df156fea26469706673582212204ef01fc3730b5c8e22138f4d5b71bf8f5b6b5bcf772bdd6d4ae9d149ec7f543c64736f6c634300081c0033",
}

// TenBridgeTestnetABI is the input ABI used to generate the binding from.
// Deprecated: Use TenBridgeTestnetMetaData.ABI instead.
var TenBridgeTestnetABI = TenBridgeTestnetMetaData.ABI

// TenBridgeTestnetBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TenBridgeTestnetMetaData.Bin instead.
var TenBridgeTestnetBin = TenBridgeTestnetMetaData.Bin

// DeployTenBridgeTestnet deploys a new Ethereum contract, binding an instance of TenBridgeTestnet to it.
func DeployTenBridgeTestnet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TenBridgeTestnet, error) {
	parsed, err := TenBridgeTestnetMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TenBridgeTestnetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TenBridgeTestnet{TenBridgeTestnetCaller: TenBridgeTestnetCaller{contract: contract}, TenBridgeTestnetTransactor: TenBridgeTestnetTransactor{contract: contract}, TenBridgeTestnetFilterer: TenBridgeTestnetFilterer{contract: contract}}, nil
}

// TenBridgeTestnet is an auto generated Go binding around an Ethereum contract.
type TenBridgeTestnet struct {
	TenBridgeTestnetCaller     // Read-only binding to the contract
	TenBridgeTestnetTransactor // Write-only binding to the contract
	TenBridgeTestnetFilterer   // Log filterer for contract events
}

// TenBridgeTestnetCaller is an auto generated read-only Go binding around an Ethereum contract.
type TenBridgeTestnetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeTestnetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TenBridgeTestnetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeTestnetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TenBridgeTestnetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeTestnetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TenBridgeTestnetSession struct {
	Contract     *TenBridgeTestnet // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TenBridgeTestnetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TenBridgeTestnetCallerSession struct {
	Contract *TenBridgeTestnetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TenBridgeTestnetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TenBridgeTestnetTransactorSession struct {
	Contract     *TenBridgeTestnetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TenBridgeTestnetRaw is an auto generated low-level Go binding around an Ethereum contract.
type TenBridgeTestnetRaw struct {
	Contract *TenBridgeTestnet // Generic contract binding to access the raw methods on
}

// TenBridgeTestnetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TenBridgeTestnetCallerRaw struct {
	Contract *TenBridgeTestnetCaller // Generic read-only contract binding to access the raw methods on
}

// TenBridgeTestnetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TenBridgeTestnetTransactorRaw struct {
	Contract *TenBridgeTestnetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTenBridgeTestnet creates a new instance of TenBridgeTestnet, bound to a specific deployed contract.
func NewTenBridgeTestnet(address common.Address, backend bind.ContractBackend) (*TenBridgeTestnet, error) {
	contract, err := bindTenBridgeTestnet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnet{TenBridgeTestnetCaller: TenBridgeTestnetCaller{contract: contract}, TenBridgeTestnetTransactor: TenBridgeTestnetTransactor{contract: contract}, TenBridgeTestnetFilterer: TenBridgeTestnetFilterer{contract: contract}}, nil
}

// NewTenBridgeTestnetCaller creates a new read-only instance of TenBridgeTestnet, bound to a specific deployed contract.
func NewTenBridgeTestnetCaller(address common.Address, caller bind.ContractCaller) (*TenBridgeTestnetCaller, error) {
	contract, err := bindTenBridgeTestnet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetCaller{contract: contract}, nil
}

// NewTenBridgeTestnetTransactor creates a new write-only instance of TenBridgeTestnet, bound to a specific deployed contract.
func NewTenBridgeTestnetTransactor(address common.Address, transactor bind.ContractTransactor) (*TenBridgeTestnetTransactor, error) {
	contract, err := bindTenBridgeTestnet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetTransactor{contract: contract}, nil
}

// NewTenBridgeTestnetFilterer creates a new log filterer instance of TenBridgeTestnet, bound to a specific deployed contract.
func NewTenBridgeTestnetFilterer(address common.Address, filterer bind.ContractFilterer) (*TenBridgeTestnetFilterer, error) {
	contract, err := bindTenBridgeTestnet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetFilterer{contract: contract}, nil
}

// bindTenBridgeTestnet binds a generic wrapper to an already deployed contract.
func bindTenBridgeTestnet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TenBridgeTestnetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenBridgeTestnet *TenBridgeTestnetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenBridgeTestnet.Contract.TenBridgeTestnetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenBridgeTestnet *TenBridgeTestnetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.TenBridgeTestnetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenBridgeTestnet *TenBridgeTestnetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.TenBridgeTestnetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenBridgeTestnet *TenBridgeTestnetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenBridgeTestnet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenBridgeTestnet *TenBridgeTestnetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenBridgeTestnet *TenBridgeTestnetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetSession) ADMINROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.ADMINROLE(&_TenBridgeTestnet.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) ADMINROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.ADMINROLE(&_TenBridgeTestnet.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.DEFAULTADMINROLE(&_TenBridgeTestnet.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.DEFAULTADMINROLE(&_TenBridgeTestnet.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) ERC20TOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "ERC20_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetSession) ERC20TOKENROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.ERC20TOKENROLE(&_TenBridgeTestnet.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) ERC20TOKENROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.ERC20TOKENROLE(&_TenBridgeTestnet.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) NATIVETOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "NATIVE_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetSession) NATIVETOKENROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.NATIVETOKENROLE(&_TenBridgeTestnet.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) NATIVETOKENROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.NATIVETOKENROLE(&_TenBridgeTestnet.CallOpts)
}

// SUSPENDEDERC20ROLE is a free data retrieval call binding the contract method 0x5ccc9613.
//
// Solidity: function SUSPENDED_ERC20_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) SUSPENDEDERC20ROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "SUSPENDED_ERC20_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SUSPENDEDERC20ROLE is a free data retrieval call binding the contract method 0x5ccc9613.
//
// Solidity: function SUSPENDED_ERC20_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetSession) SUSPENDEDERC20ROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.SUSPENDEDERC20ROLE(&_TenBridgeTestnet.CallOpts)
}

// SUSPENDEDERC20ROLE is a free data retrieval call binding the contract method 0x5ccc9613.
//
// Solidity: function SUSPENDED_ERC20_ROLE() view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) SUSPENDEDERC20ROLE() ([32]byte, error) {
	return _TenBridgeTestnet.Contract.SUSPENDEDERC20ROLE(&_TenBridgeTestnet.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TenBridgeTestnet.Contract.GetRoleAdmin(&_TenBridgeTestnet.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TenBridgeTestnet.Contract.GetRoleAdmin(&_TenBridgeTestnet.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridgeTestnet *TenBridgeTestnetSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TenBridgeTestnet.Contract.HasRole(&_TenBridgeTestnet.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TenBridgeTestnet.Contract.HasRole(&_TenBridgeTestnet.CallOpts, role, account)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) MessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "messageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetSession) MessageBus() (common.Address, error) {
	return _TenBridgeTestnet.Contract.MessageBus(&_TenBridgeTestnet.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) MessageBus() (common.Address, error) {
	return _TenBridgeTestnet.Contract.MessageBus(&_TenBridgeTestnet.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetSession) Messenger() (common.Address, error) {
	return _TenBridgeTestnet.Contract.Messenger(&_TenBridgeTestnet.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) Messenger() (common.Address, error) {
	return _TenBridgeTestnet.Contract.Messenger(&_TenBridgeTestnet.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_TenBridgeTestnet *TenBridgeTestnetSession) Nonce() (*big.Int, error) {
	return _TenBridgeTestnet.Contract.Nonce(&_TenBridgeTestnet.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) Nonce() (*big.Int, error) {
	return _TenBridgeTestnet.Contract.Nonce(&_TenBridgeTestnet.CallOpts)
}

// RemoteBridgeAddress is a free data retrieval call binding the contract method 0x0f0a9a4b.
//
// Solidity: function remoteBridgeAddress() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) RemoteBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "remoteBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RemoteBridgeAddress is a free data retrieval call binding the contract method 0x0f0a9a4b.
//
// Solidity: function remoteBridgeAddress() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetSession) RemoteBridgeAddress() (common.Address, error) {
	return _TenBridgeTestnet.Contract.RemoteBridgeAddress(&_TenBridgeTestnet.CallOpts)
}

// RemoteBridgeAddress is a free data retrieval call binding the contract method 0x0f0a9a4b.
//
// Solidity: function remoteBridgeAddress() view returns(address)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) RemoteBridgeAddress() (common.Address, error) {
	return _TenBridgeTestnet.Contract.RemoteBridgeAddress(&_TenBridgeTestnet.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridgeTestnet *TenBridgeTestnetCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TenBridgeTestnet.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridgeTestnet *TenBridgeTestnetSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TenBridgeTestnet.Contract.SupportsInterface(&_TenBridgeTestnet.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridgeTestnet *TenBridgeTestnetCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TenBridgeTestnet.Contract.SupportsInterface(&_TenBridgeTestnet.CallOpts, interfaceId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.GrantRole(&_TenBridgeTestnet.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.GrantRole(&_TenBridgeTestnet.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address owner) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) Initialize(opts *bind.TransactOpts, messenger common.Address, owner common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "initialize", messenger, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address owner) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) Initialize(messenger common.Address, owner common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.Initialize(&_TenBridgeTestnet.TransactOpts, messenger, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address owner) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) Initialize(messenger common.Address, owner common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.Initialize(&_TenBridgeTestnet.TransactOpts, messenger, owner)
}

// PauseToken is a paid mutator transaction binding the contract method 0x7c41ad2c.
//
// Solidity: function pauseToken(address asset) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) PauseToken(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "pauseToken", asset)
}

// PauseToken is a paid mutator transaction binding the contract method 0x7c41ad2c.
//
// Solidity: function pauseToken(address asset) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) PauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.PauseToken(&_TenBridgeTestnet.TransactOpts, asset)
}

// PauseToken is a paid mutator transaction binding the contract method 0x7c41ad2c.
//
// Solidity: function pauseToken(address asset) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) PauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.PauseToken(&_TenBridgeTestnet.TransactOpts, asset)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) PromoteToAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "promoteToAdmin", newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.PromoteToAdmin(&_TenBridgeTestnet.TransactOpts, newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.PromoteToAdmin(&_TenBridgeTestnet.TransactOpts, newAdmin)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.ReceiveAssets(&_TenBridgeTestnet.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.ReceiveAssets(&_TenBridgeTestnet.TransactOpts, asset, amount, receiver)
}

// RecoverTestnetFunds is a paid mutator transaction binding the contract method 0xc091a58b.
//
// Solidity: function recoverTestnetFunds(address receiver) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) RecoverTestnetFunds(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "recoverTestnetFunds", receiver)
}

// RecoverTestnetFunds is a paid mutator transaction binding the contract method 0xc091a58b.
//
// Solidity: function recoverTestnetFunds(address receiver) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) RecoverTestnetFunds(receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.RecoverTestnetFunds(&_TenBridgeTestnet.TransactOpts, receiver)
}

// RecoverTestnetFunds is a paid mutator transaction binding the contract method 0xc091a58b.
//
// Solidity: function recoverTestnetFunds(address receiver) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) RecoverTestnetFunds(receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.RecoverTestnetFunds(&_TenBridgeTestnet.TransactOpts, receiver)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.RenounceRole(&_TenBridgeTestnet.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.RenounceRole(&_TenBridgeTestnet.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.RevokeRole(&_TenBridgeTestnet.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.RevokeRole(&_TenBridgeTestnet.TransactOpts, role, account)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) SendERC20(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "sendERC20", asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.SendERC20(&_TenBridgeTestnet.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.SendERC20(&_TenBridgeTestnet.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) SendNative(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "sendNative", receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.SendNative(&_TenBridgeTestnet.TransactOpts, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.SendNative(&_TenBridgeTestnet.TransactOpts, receiver)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) SetRemoteBridge(opts *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "setRemoteBridge", bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.SetRemoteBridge(&_TenBridgeTestnet.TransactOpts, bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.SetRemoteBridge(&_TenBridgeTestnet.TransactOpts, bridge)
}

// UnpauseToken is a paid mutator transaction binding the contract method 0x3b3bff0f.
//
// Solidity: function unpauseToken(address asset) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) UnpauseToken(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "unpauseToken", asset)
}

// UnpauseToken is a paid mutator transaction binding the contract method 0x3b3bff0f.
//
// Solidity: function unpauseToken(address asset) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) UnpauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.UnpauseToken(&_TenBridgeTestnet.TransactOpts, asset)
}

// UnpauseToken is a paid mutator transaction binding the contract method 0x3b3bff0f.
//
// Solidity: function unpauseToken(address asset) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) UnpauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.UnpauseToken(&_TenBridgeTestnet.TransactOpts, asset)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactor) WhitelistToken(opts *bind.TransactOpts, asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridgeTestnet.contract.Transact(opts, "whitelistToken", asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridgeTestnet *TenBridgeTestnetSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.WhitelistToken(&_TenBridgeTestnet.TransactOpts, asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridgeTestnet *TenBridgeTestnetTransactorSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridgeTestnet.Contract.WhitelistToken(&_TenBridgeTestnet.TransactOpts, asset, name, symbol)
}

// TenBridgeTestnetInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TenBridgeTestnet contract.
type TenBridgeTestnetInitializedIterator struct {
	Event *TenBridgeTestnetInitialized // Event containing the contract specifics and raw log

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
func (it *TenBridgeTestnetInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeTestnetInitialized)
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
		it.Event = new(TenBridgeTestnetInitialized)
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
func (it *TenBridgeTestnetInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeTestnetInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeTestnetInitialized represents a Initialized event raised by the TenBridgeTestnet contract.
type TenBridgeTestnetInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) FilterInitialized(opts *bind.FilterOpts) (*TenBridgeTestnetInitializedIterator, error) {

	logs, sub, err := _TenBridgeTestnet.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetInitializedIterator{contract: _TenBridgeTestnet.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TenBridgeTestnetInitialized) (event.Subscription, error) {

	logs, sub, err := _TenBridgeTestnet.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeTestnetInitialized)
				if err := _TenBridgeTestnet.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) ParseInitialized(log types.Log) (*TenBridgeTestnetInitialized, error) {
	event := new(TenBridgeTestnetInitialized)
	if err := _TenBridgeTestnet.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeTestnetRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TenBridgeTestnet contract.
type TenBridgeTestnetRoleAdminChangedIterator struct {
	Event *TenBridgeTestnetRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TenBridgeTestnetRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeTestnetRoleAdminChanged)
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
		it.Event = new(TenBridgeTestnetRoleAdminChanged)
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
func (it *TenBridgeTestnetRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeTestnetRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeTestnetRoleAdminChanged represents a RoleAdminChanged event raised by the TenBridgeTestnet contract.
type TenBridgeTestnetRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TenBridgeTestnetRoleAdminChangedIterator, error) {

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

	logs, sub, err := _TenBridgeTestnet.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetRoleAdminChangedIterator{contract: _TenBridgeTestnet.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TenBridgeTestnetRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _TenBridgeTestnet.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeTestnetRoleAdminChanged)
				if err := _TenBridgeTestnet.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) ParseRoleAdminChanged(log types.Log) (*TenBridgeTestnetRoleAdminChanged, error) {
	event := new(TenBridgeTestnetRoleAdminChanged)
	if err := _TenBridgeTestnet.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeTestnetRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TenBridgeTestnet contract.
type TenBridgeTestnetRoleGrantedIterator struct {
	Event *TenBridgeTestnetRoleGranted // Event containing the contract specifics and raw log

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
func (it *TenBridgeTestnetRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeTestnetRoleGranted)
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
		it.Event = new(TenBridgeTestnetRoleGranted)
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
func (it *TenBridgeTestnetRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeTestnetRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeTestnetRoleGranted represents a RoleGranted event raised by the TenBridgeTestnet contract.
type TenBridgeTestnetRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TenBridgeTestnetRoleGrantedIterator, error) {

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

	logs, sub, err := _TenBridgeTestnet.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetRoleGrantedIterator{contract: _TenBridgeTestnet.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TenBridgeTestnetRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TenBridgeTestnet.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeTestnetRoleGranted)
				if err := _TenBridgeTestnet.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) ParseRoleGranted(log types.Log) (*TenBridgeTestnetRoleGranted, error) {
	event := new(TenBridgeTestnetRoleGranted)
	if err := _TenBridgeTestnet.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeTestnetRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TenBridgeTestnet contract.
type TenBridgeTestnetRoleRevokedIterator struct {
	Event *TenBridgeTestnetRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TenBridgeTestnetRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeTestnetRoleRevoked)
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
		it.Event = new(TenBridgeTestnetRoleRevoked)
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
func (it *TenBridgeTestnetRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeTestnetRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeTestnetRoleRevoked represents a RoleRevoked event raised by the TenBridgeTestnet contract.
type TenBridgeTestnetRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TenBridgeTestnetRoleRevokedIterator, error) {

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

	logs, sub, err := _TenBridgeTestnet.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetRoleRevokedIterator{contract: _TenBridgeTestnet.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TenBridgeTestnetRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TenBridgeTestnet.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeTestnetRoleRevoked)
				if err := _TenBridgeTestnet.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) ParseRoleRevoked(log types.Log) (*TenBridgeTestnetRoleRevoked, error) {
	event := new(TenBridgeTestnetRoleRevoked)
	if err := _TenBridgeTestnet.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeTestnetWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the TenBridgeTestnet contract.
type TenBridgeTestnetWithdrawalIterator struct {
	Event *TenBridgeTestnetWithdrawal // Event containing the contract specifics and raw log

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
func (it *TenBridgeTestnetWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeTestnetWithdrawal)
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
		it.Event = new(TenBridgeTestnetWithdrawal)
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
func (it *TenBridgeTestnetWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeTestnetWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeTestnetWithdrawal represents a Withdrawal event raised by the TenBridgeTestnet contract.
type TenBridgeTestnetWithdrawal struct {
	Receiver common.Address
	Asset    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b6398.
//
// Solidity: event Withdrawal(address indexed receiver, address indexed asset, uint256 amount)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) FilterWithdrawal(opts *bind.FilterOpts, receiver []common.Address, asset []common.Address) (*TenBridgeTestnetWithdrawalIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _TenBridgeTestnet.contract.FilterLogs(opts, "Withdrawal", receiverRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTestnetWithdrawalIterator{contract: _TenBridgeTestnet.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b6398.
//
// Solidity: event Withdrawal(address indexed receiver, address indexed asset, uint256 amount)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *TenBridgeTestnetWithdrawal, receiver []common.Address, asset []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _TenBridgeTestnet.contract.WatchLogs(opts, "Withdrawal", receiverRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeTestnetWithdrawal)
				if err := _TenBridgeTestnet.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b6398.
//
// Solidity: event Withdrawal(address indexed receiver, address indexed asset, uint256 amount)
func (_TenBridgeTestnet *TenBridgeTestnetFilterer) ParseWithdrawal(log types.Log) (*TenBridgeTestnetWithdrawal, error) {
	event := new(TenBridgeTestnetWithdrawal)
	if err := _TenBridgeTestnet.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
