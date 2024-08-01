// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObscuroBridge

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

// ObscuroBridgeMetaData contains all meta data concerning the ObscuroBridge contract.
var ObscuroBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERC20_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NATIVE_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messengerAddress\",\"type\":\"address\"}],\"name\":\"configure\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"promoteToAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"removeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"setRemoteBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"whitelistToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805463ffffffff60a01b1916905534801561002057600080fd5b50611944806100306000396000f3fe6080604052600436106101445760003560e01c806375b238fc116100c0578063a217fddf11610074578063c4d66de811610059578063c4d66de8146103b4578063d547741f146103d4578063e4c3ebc7146103f457600080fd5b8063a217fddf1461037f578063a381c8e21461039457600080fd5b806383bece4d116100a557806383bece4d146102f957806391d148541461031957806393b374421461035f57600080fd5b806375b238fc146102a557806375cb2672146102d957600080fd5b80632f2ff15d11610117578063498d82ab116100fc578063498d82ab146102315780635d872970146102515780635fa7b5841461028557600080fd5b80632f2ff15d146101f157806336568abe1461021157600080fd5b806301ffc9a71461014957806316ce81491461017e5780631888d712146101a0578063248a9ca3146101b3575b600080fd5b34801561015557600080fd5b5061016961016436600461154f565b610428565b60405190151581526020015b60405180910390f35b34801561018a57600080fd5b5061019e61019936600461158e565b610491565b005b61019e6101ae36600461158e565b6104eb565b3480156101bf57600080fd5b506101e36101ce3660046115ab565b60009081526002602052604090206001015490565b604051908152602001610175565b3480156101fd57600080fd5b5061019e61020c3660046115c4565b610624565b34801561021d57600080fd5b5061019e61022c3660046115c4565b61064f565b34801561023d57600080fd5b5061019e61024c36600461163d565b6106a0565b34801561025d57600080fd5b506101e37f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a81565b34801561029157600080fd5b5061019e6102a036600461158e565b610765565b3480156102b157600080fd5b506101e37fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b3480156102e557600080fd5b5061019e6102f436600461158e565b6107b9565b34801561030557600080fd5b5061019e6103143660046116c0565b610894565b34801561032557600080fd5b506101696103343660046115c4565b60009182526002602090815260408084206001600160a01b0393909316845291905290205460ff1690565b34801561036b57600080fd5b5061019e61037a36600461158e565b610ab0565b34801561038b57600080fd5b506101e3600081565b3480156103a057600080fd5b5061019e6103af3660046116c0565b610b04565b3480156103c057600080fd5b5061019e6103cf36600461158e565b610cb1565b3480156103e057600080fd5b5061019e6103ef3660046115c4565b610e45565b34801561040057600080fd5b506101e37fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad211057881565b60006001600160e01b031982167f7965db0b00000000000000000000000000000000000000000000000000000000148061048b57507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756104bb81610e6a565b506003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b600034116105405760405162461bcd60e51b815260206004820152600f60248201527f456d707479207472616e736665722e000000000000000000000000000000000060448201526064015b60405180910390fd5b604080518082018252348082526001600160a01b03848116602093840190815284519384019290925290518116828401528251808303840181526060909201909252600354909161059891168260025b600080610e77565b6001546001600160a01b03166040517f346633fb0000000000000000000000000000000000000000000000000000000081526001600160a01b038481166004830152346024830181905292169163346633fb916044016000604051808303818588803b15801561060757600080fd5b505af115801561061b573d6000803e3d6000fd5b50505050505050565b60008281526002602052604090206001015461063f81610e6a565b6106498383610f83565b50505050565b6001600160a01b0381163314610691576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61069b8282611031565b505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756106ca81610e6a565b6106f47f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a87610f83565b50600063458ffd6360e01b878787878760405160240161071895949392919061172b565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b03199093169290921790915260035490915061061b906001600160a01b0316826001610590565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561078f81610e6a565b61069b7f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a83611031565b6107c16110b8565b6000805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117909155604080517fa1a227fa000000000000000000000000000000000000000000000000000000008152905163a1a227fa916004808201926020929091908290030181865afa158015610840573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610864919061176d565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039290921691909117905550565b6003546000546001600160a01b039182169116331461091b5760405162461bcd60e51b815260206004820152603060248201527f436f6e74726163742063616c6c6572206973206e6f742074686520726567697360448201527f7465726564206d657373656e67657221000000000000000000000000000000006064820152608401610537565b806001600160a01b031661092d611121565b6001600160a01b0316146109a95760405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f72726563742073656e646572210000000000000000000000000000006064820152608401610537565b6001600160a01b03841660009081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16156109f9576109f484848461119e565b610649565b6001600160a01b03841660009081527f13ad2d85210d477fe1a6e25654c8250308cf29b050a4bf0b039d70467486712c602052604090205460ff1615610a42576109f4826111a9565b60405162461bcd60e51b815260206004820152602560248201527f417474656d7074696e6720746f20776974686472617720756e6b6e6f776e206160448201527f737365742e0000000000000000000000000000000000000000000000000000006064820152608401610537565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610ada81610e6a565b61069b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177583610f83565b60008211610b545760405162461bcd60e51b815260206004820152601a60248201527f417474656d7074696e6720656d707479207472616e736665722e0000000000006044820152606401610537565b6001600160a01b03831660009081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16610c275760405162461bcd60e51b815260206004820152604e60248201527f54686973206164647265737320686173206e6f74206265656e20676976656e2060448201527f61207479706520616e64206973207468757320636f6e73696465726564206e6f60648201527f742077686974656c69737465642e000000000000000000000000000000000000608482015260a401610537565b610c338333308561124b565b604080516001600160a01b038581166024830152604482018590528381166064808401919091528351808403909101815260849092019092526020810180516001600160e01b03167f83bece4d0000000000000000000000000000000000000000000000000000000017905260035490916106499116826000610590565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff16600081158015610cfc5750825b905060008267ffffffffffffffff166001148015610d195750303b155b905081158015610d27575080155b15610d5e576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff191660011785558315610d9257845468ff00000000000000001916680100000000000000001785555b610d9b866107b9565b610dc57fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177533610f83565b50610df17fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad21105786000610f83565b508315610e3d57845468ff000000000000000019168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b600082815260026020526040902060010154610e6081610e6a565b6106498383611031565b610e7481336112b2565b50565b60006040518060600160405280876001600160a01b0316815260200186815260200184815250604051602001610ead91906117da565b60408051808303601f19018152919052600180549192506001600160a01b0382169163b1454caa917401000000000000000000000000000000000000000090910463ffffffff16906014610f008361181f565b91906101000a81548163ffffffff021916908363ffffffff1602179055508684866040518563ffffffff1660e01b8152600401610f409493929190611869565b6020604051808303816000875af1158015610f5f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061061b91906118a6565b60008281526002602090815260408083206001600160a01b038516845290915281205460ff166110295760008381526002602090815260408083206001600160a01b03861684529091529020805460ff19166001179055610fe13390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161048b565b50600061048b565b60008281526002602090815260408083206001600160a01b038516845290915281205460ff16156110295760008381526002602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161048b565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff1661111f576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60008060009054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b8152600401602060405180830381865afa158015611175573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611199919061176d565b905090565b61069b838284611320565b6040516000906001600160a01b038316908281818181865af19150503d80600081146111f1576040519150601f19603f3d011682016040523d82523d6000602084013e6111f6565b606091505b50509050806112475760405162461bcd60e51b815260206004820152601460248201527f4661696c656420746f2073656e642045746865720000000000000000000000006044820152606401610537565b5050565b6040516001600160a01b0384811660248301528381166044830152606482018390526106499186918216906323b872dd906084015b604051602081830303815290604052915060e01b6020820180516001600160e01b038381831617835250505050611351565b60008281526002602090815260408083206001600160a01b038516845290915290205460ff16611247576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b038216600482015260248101839052604401610537565b6040516001600160a01b0383811660248301526044820183905261069b91859182169063a9059cbb90606401611280565b60006113666001600160a01b038416836113cd565b9050805160001415801561138b57508080602001905181019061138991906118d0565b155b1561069b576040517f5274afe70000000000000000000000000000000000000000000000000000000081526001600160a01b0384166004820152602401610537565b60606113db838360006113e2565b9392505050565b606081471015611420576040517fcd786059000000000000000000000000000000000000000000000000000000008152306004820152602401610537565b600080856001600160a01b0316848660405161143c91906118f2565b60006040518083038185875af1925050503d8060008114611479576040519150601f19603f3d011682016040523d82523d6000602084013e61147e565b606091505b509150915061148e868383611498565b9695505050505050565b6060826114ad576114a88261150d565b6113db565b81511580156114c457506001600160a01b0384163b155b15611506576040517f9996b3150000000000000000000000000000000000000000000000000000000081526001600160a01b0385166004820152602401610537565b50806113db565b80511561151d5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006020828403121561156157600080fd5b81356001600160e01b0319811681146113db57600080fd5b6001600160a01b0381168114610e7457600080fd5b6000602082840312156115a057600080fd5b81356113db81611579565b6000602082840312156115bd57600080fd5b5035919050565b600080604083850312156115d757600080fd5b8235915060208301356115e981611579565b809150509250929050565b60008083601f84011261160657600080fd5b50813567ffffffffffffffff81111561161e57600080fd5b60208301915083602082850101111561163657600080fd5b9250929050565b60008060008060006060868803121561165557600080fd5b853561166081611579565b9450602086013567ffffffffffffffff8082111561167d57600080fd5b61168989838a016115f4565b909650945060408801359150808211156116a257600080fd5b506116af888289016115f4565b969995985093965092949392505050565b6000806000606084860312156116d557600080fd5b83356116e081611579565b92506020840135915060408401356116f781611579565b809150509250925092565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b038616815260606020820152600061174e606083018688611702565b8281036040840152611761818587611702565b98975050505050505050565b60006020828403121561177f57600080fd5b81516113db81611579565b60005b838110156117a557818101518382015260200161178d565b50506000910152565b600081518084526117c681602086016020860161178a565b601f01601f19169290920160200192915050565b602081526001600160a01b038251166020820152600060208301516060604084015261180960808401826117ae565b9050604084015160608401528091505092915050565b600063ffffffff80831681810361185f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b600063ffffffff80871683528086166020840152506080604083015261189260808301856117ae565b905060ff8316606083015295945050505050565b6000602082840312156118b857600080fd5b815167ffffffffffffffff811681146113db57600080fd5b6000602082840312156118e257600080fd5b815180151581146113db57600080fd5b6000825161190481846020870161178a565b919091019291505056fea26469706673582212209080b4341dacf82f9a993515a09081c2f61eed4e485d15f5fba6af9c0305b82c64736f6c63430008140033",
}

// ObscuroBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ObscuroBridgeMetaData.ABI instead.
var ObscuroBridgeABI = ObscuroBridgeMetaData.ABI

// ObscuroBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObscuroBridgeMetaData.Bin instead.
var ObscuroBridgeBin = ObscuroBridgeMetaData.Bin

// DeployObscuroBridge deploys a new Ethereum contract, binding an instance of ObscuroBridge to it.
func DeployObscuroBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ObscuroBridge, error) {
	parsed, err := ObscuroBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObscuroBridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObscuroBridge{ObscuroBridgeCaller: ObscuroBridgeCaller{contract: contract}, ObscuroBridgeTransactor: ObscuroBridgeTransactor{contract: contract}, ObscuroBridgeFilterer: ObscuroBridgeFilterer{contract: contract}}, nil
}

// ObscuroBridge is an auto generated Go binding around an Ethereum contract.
type ObscuroBridge struct {
	ObscuroBridgeCaller     // Read-only binding to the contract
	ObscuroBridgeTransactor // Write-only binding to the contract
	ObscuroBridgeFilterer   // Log filterer for contract events
}

// ObscuroBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ObscuroBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ObscuroBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObscuroBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObscuroBridgeSession struct {
	Contract     *ObscuroBridge    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObscuroBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObscuroBridgeCallerSession struct {
	Contract *ObscuroBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ObscuroBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObscuroBridgeTransactorSession struct {
	Contract     *ObscuroBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ObscuroBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ObscuroBridgeRaw struct {
	Contract *ObscuroBridge // Generic contract binding to access the raw methods on
}

// ObscuroBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObscuroBridgeCallerRaw struct {
	Contract *ObscuroBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ObscuroBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObscuroBridgeTransactorRaw struct {
	Contract *ObscuroBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewObscuroBridge creates a new instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridge(address common.Address, backend bind.ContractBackend) (*ObscuroBridge, error) {
	contract, err := bindObscuroBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridge{ObscuroBridgeCaller: ObscuroBridgeCaller{contract: contract}, ObscuroBridgeTransactor: ObscuroBridgeTransactor{contract: contract}, ObscuroBridgeFilterer: ObscuroBridgeFilterer{contract: contract}}, nil
}

// NewObscuroBridgeCaller creates a new read-only instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeCaller(address common.Address, caller bind.ContractCaller) (*ObscuroBridgeCaller, error) {
	contract, err := bindObscuroBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeCaller{contract: contract}, nil
}

// NewObscuroBridgeTransactor creates a new write-only instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ObscuroBridgeTransactor, error) {
	contract, err := bindObscuroBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeTransactor{contract: contract}, nil
}

// NewObscuroBridgeFilterer creates a new log filterer instance of ObscuroBridge, bound to a specific deployed contract.
func NewObscuroBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ObscuroBridgeFilterer, error) {
	contract, err := bindObscuroBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeFilterer{contract: contract}, nil
}

// bindObscuroBridge binds a generic wrapper to an already deployed contract.
func bindObscuroBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ObscuroBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroBridge *ObscuroBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroBridge.Contract.ObscuroBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroBridge *ObscuroBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ObscuroBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroBridge *ObscuroBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ObscuroBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroBridge *ObscuroBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroBridge *ObscuroBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroBridge *ObscuroBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) ADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ADMINROLE(&_ObscuroBridge.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) ADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ADMINROLE(&_ObscuroBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.DEFAULTADMINROLE(&_ObscuroBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.DEFAULTADMINROLE(&_ObscuroBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) ERC20TOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "ERC20_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) ERC20TOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ERC20TOKENROLE(&_ObscuroBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) ERC20TOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.ERC20TOKENROLE(&_ObscuroBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) NATIVETOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "NATIVE_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) NATIVETOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.NATIVETOKENROLE(&_ObscuroBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) NATIVETOKENROLE() ([32]byte, error) {
	return _ObscuroBridge.Contract.NATIVETOKENROLE(&_ObscuroBridge.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ObscuroBridge.Contract.GetRoleAdmin(&_ObscuroBridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroBridge *ObscuroBridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ObscuroBridge.Contract.GetRoleAdmin(&_ObscuroBridge.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ObscuroBridge.Contract.HasRole(&_ObscuroBridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ObscuroBridge.Contract.HasRole(&_ObscuroBridge.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ObscuroBridge.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ObscuroBridge.Contract.SupportsInterface(&_ObscuroBridge.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroBridge *ObscuroBridgeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ObscuroBridge.Contract.SupportsInterface(&_ObscuroBridge.CallOpts, interfaceId)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) Configure(opts *bind.TransactOpts, messengerAddress common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "configure", messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_ObscuroBridge *ObscuroBridgeSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Configure(&_ObscuroBridge.TransactOpts, messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Configure(&_ObscuroBridge.TransactOpts, messengerAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.GrantRole(&_ObscuroBridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.GrantRole(&_ObscuroBridge.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) Initialize(opts *bind.TransactOpts, messenger common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "initialize", messenger)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_ObscuroBridge *ObscuroBridgeSession) Initialize(messenger common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Initialize(&_ObscuroBridge.TransactOpts, messenger)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address messenger) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) Initialize(messenger common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.Initialize(&_ObscuroBridge.TransactOpts, messenger)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) PromoteToAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "promoteToAdmin", newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_ObscuroBridge *ObscuroBridgeSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.PromoteToAdmin(&_ObscuroBridge.TransactOpts, newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.PromoteToAdmin(&_ObscuroBridge.TransactOpts, newAdmin)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ReceiveAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.ReceiveAssets(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) RemoveToken(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "removeToken", asset)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_ObscuroBridge *ObscuroBridgeSession) RemoveToken(asset common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RemoveToken(&_ObscuroBridge.TransactOpts, asset)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address asset) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) RemoveToken(asset common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RemoveToken(&_ObscuroBridge.TransactOpts, asset)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ObscuroBridge *ObscuroBridgeSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RenounceRole(&_ObscuroBridge.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RenounceRole(&_ObscuroBridge.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RevokeRole(&_ObscuroBridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.RevokeRole(&_ObscuroBridge.TransactOpts, role, account)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SendERC20(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "sendERC20", asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendERC20(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendERC20(&_ObscuroBridge.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SendNative(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "sendNative", receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendNative(&_ObscuroBridge.TransactOpts, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SendNative(&_ObscuroBridge.TransactOpts, receiver)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) SetRemoteBridge(opts *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "setRemoteBridge", bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_ObscuroBridge *ObscuroBridgeSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SetRemoteBridge(&_ObscuroBridge.TransactOpts, bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.SetRemoteBridge(&_ObscuroBridge.TransactOpts, bridge)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_ObscuroBridge *ObscuroBridgeTransactor) WhitelistToken(opts *bind.TransactOpts, asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroBridge.contract.Transact(opts, "whitelistToken", asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_ObscuroBridge *ObscuroBridgeSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.WhitelistToken(&_ObscuroBridge.TransactOpts, asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_ObscuroBridge *ObscuroBridgeTransactorSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroBridge.Contract.WhitelistToken(&_ObscuroBridge.TransactOpts, asset, name, symbol)
}

// ObscuroBridgeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ObscuroBridge contract.
type ObscuroBridgeInitializedIterator struct {
	Event *ObscuroBridgeInitialized // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeInitialized)
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
		it.Event = new(ObscuroBridgeInitialized)
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
func (it *ObscuroBridgeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeInitialized represents a Initialized event raised by the ObscuroBridge contract.
type ObscuroBridgeInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterInitialized(opts *bind.FilterOpts) (*ObscuroBridgeInitializedIterator, error) {

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeInitializedIterator{contract: _ObscuroBridge.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeInitialized) (event.Subscription, error) {

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeInitialized)
				if err := _ObscuroBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseInitialized(log types.Log) (*ObscuroBridgeInitialized, error) {
	event := new(ObscuroBridgeInitialized)
	if err := _ObscuroBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroBridgeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ObscuroBridge contract.
type ObscuroBridgeRoleAdminChangedIterator struct {
	Event *ObscuroBridgeRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeRoleAdminChanged)
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
		it.Event = new(ObscuroBridgeRoleAdminChanged)
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
func (it *ObscuroBridgeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeRoleAdminChanged represents a RoleAdminChanged event raised by the ObscuroBridge contract.
type ObscuroBridgeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ObscuroBridgeRoleAdminChangedIterator, error) {

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

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeRoleAdminChangedIterator{contract: _ObscuroBridge.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeRoleAdminChanged)
				if err := _ObscuroBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseRoleAdminChanged(log types.Log) (*ObscuroBridgeRoleAdminChanged, error) {
	event := new(ObscuroBridgeRoleAdminChanged)
	if err := _ObscuroBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroBridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ObscuroBridge contract.
type ObscuroBridgeRoleGrantedIterator struct {
	Event *ObscuroBridgeRoleGranted // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeRoleGranted)
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
		it.Event = new(ObscuroBridgeRoleGranted)
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
func (it *ObscuroBridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeRoleGranted represents a RoleGranted event raised by the ObscuroBridge contract.
type ObscuroBridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ObscuroBridgeRoleGrantedIterator, error) {

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

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeRoleGrantedIterator{contract: _ObscuroBridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeRoleGranted)
				if err := _ObscuroBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseRoleGranted(log types.Log) (*ObscuroBridgeRoleGranted, error) {
	event := new(ObscuroBridgeRoleGranted)
	if err := _ObscuroBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroBridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ObscuroBridge contract.
type ObscuroBridgeRoleRevokedIterator struct {
	Event *ObscuroBridgeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ObscuroBridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroBridgeRoleRevoked)
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
		it.Event = new(ObscuroBridgeRoleRevoked)
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
func (it *ObscuroBridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroBridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroBridgeRoleRevoked represents a RoleRevoked event raised by the ObscuroBridge contract.
type ObscuroBridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ObscuroBridgeRoleRevokedIterator, error) {

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

	logs, sub, err := _ObscuroBridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroBridgeRoleRevokedIterator{contract: _ObscuroBridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroBridge *ObscuroBridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ObscuroBridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ObscuroBridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroBridgeRoleRevoked)
				if err := _ObscuroBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_ObscuroBridge *ObscuroBridgeFilterer) ParseRoleRevoked(log types.Log) (*ObscuroBridgeRoleRevoked, error) {
	event := new(ObscuroBridgeRoleRevoked)
	if err := _ObscuroBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
