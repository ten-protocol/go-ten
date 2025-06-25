// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TenBridge

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

// TenBridgeMetaData contains all meta data concerning the TenBridge contract.
var TenBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERC20_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NATIVE_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUSPENDED_ERC20_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messengerAddress\",\"type\":\"address\"}],\"name\":\"configure\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"pauseToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"promoteToAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"retrieveAllFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"setRemoteBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"unpauseToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"whitelistToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50611f6c8061001c5f395ff3fe60806040526004361061016d575f3560e01c80635d872970116100c657806391d148541161007c578063a381c8e211610057578063a381c8e21461043b578063d547741f1461044e578063e4c3ebc71461046d575f5ffd5b806391d14854146103c557806393b3744214610409578063a217fddf14610428575f5ffd5b806375cb2672116100ac57806375cb2672146103685780637c41ad2c1461038757806383bece4d146103a6575f5ffd5b80635d8729701461030257806375b238fc14610335575f5ffd5b806336568abe11610126578063485cc95511610101578063485cc95514610291578063498d82ab146102b05780635ccc9613146102cf575f5ffd5b806336568abe1461023457806336d2da90146102535780633b3bff0f14610272575f5ffd5b80631888d712116101565780631888d712146101c7578063248a9ca3146101da5780632f2ff15d14610215575f5ffd5b806301ffc9a71461017157806316ce8149146101a6575b5f5ffd5b34801561017c575f5ffd5b5061019061018b3660046115f9565b6104a0565b60405161019d9190611620565b60405180910390f35b3480156101b1575f5ffd5b506101c56101c0366004611652565b610508565b005b6101c56101d5366004611652565b610594565b3480156101e5575f5ffd5b506102086101f4366004611680565b5f9081526002602052604090206001015490565b60405161019d91906116a3565b348015610220575f5ffd5b506101c561022f3660046116b1565b610605565b34801561023f575f5ffd5b506101c561024e3660046116b1565b61062f565b34801561025e575f5ffd5b506101c561026d366004611652565b610680565b34801561027d575f5ffd5b506101c561028c366004611652565b610719565b34801561029c575f5ffd5b506101c56102ab3660046116e7565b61076d565b3480156102bb575f5ffd5b506101c56102ca366004611753565b610904565b3480156102da575f5ffd5b506102087fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e81565b34801561030d575f5ffd5b506102087f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a81565b348015610340575f5ffd5b506102087fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b348015610373575f5ffd5b506101c5610382366004611652565b610a23565b348015610392575f5ffd5b506101c56103a1366004611652565b610b06565b3480156103b1575f5ffd5b506101c56103c03660046117da565b610b5a565b3480156103d0575f5ffd5b506101906103df3660046116b1565b5f9182526002602090815260408084206001600160a01b0393909316845291905290205460ff1690565b348015610414575f5ffd5b506101c5610423366004611652565b610c72565b348015610433575f5ffd5b506102085f81565b6101c56104493660046117da565b610cc6565b348015610459575f5ffd5b506101c56104683660046116b1565b610e08565b348015610478575f5ffd5b506102087fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad211057881565b5f6001600160e01b031982167f7965db0b00000000000000000000000000000000000000000000000000000000148061050257507f01ffc9a7000000000000000000000000000000000000000000000000000000006001600160e01b03198316145b92915050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561053281610e2c565b6003546001600160a01b0316156105645760405162461bcd60e51b815260040161055b9061187a565b60405180910390fd5b506003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b5f34116105b35760405162461bcd60e51b815260040161055b906118be565b5f6040518060400160405280348152602001836001600160a01b03168152506040516020016105e291906118f6565b60408051601f1981840301815291905290506106018160025f5f610e39565b5050565b5f8281526002602052604090206001015461061f81610e2c565b6106298383610ef2565b50505050565b6001600160a01b0381163314610671576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61067b8282610f9d565b505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756106aa81610e2c565b5f826001600160a01b0316476040515f6040518083038185875af1925050503d805f81146106f3576040519150601f19603f3d011682016040523d82523d5f602084013e6106f8565b606091505b505090508061067b5760405162461bcd60e51b815260040161055b90611936565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561074381610e2c565b61067b7fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e83610f9d565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156107b75750825b90505f8267ffffffffffffffff1660011480156107d35750303b155b9050811580156107e1575080155b15610818576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561084c57845468ff00000000000000001916680100000000000000001785555b61085587610a23565b61087f7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177587610ef2565b506108aa7fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad21105785f610ef2565b5083156108fb57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906108f290600190611962565b60405180910390a15b50505050505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561092e81610e2c565b6001600160a01b0386165f9081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16156109855760405162461bcd60e51b815260040161055b906119a2565b6109af7f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a87610ef2565b505f63458ffd6360e01b87878787876040516024016109d29594939291906119dd565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b0319909316929092179091526003549091506108fb906001600160a01b03168260015b5f5f5f611022565b610a2b611125565b5f805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117909155604080517fa1a227fa000000000000000000000000000000000000000000000000000000008152905163a1a227fa916004808201926020929091908290030181865afa158015610aa7573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610acb9190611a29565b600180547fffffffffffffffff000000000000000000000000000000000000000000000000166001600160a01b039290921691909117905550565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610b3081610e2c565b61067b7fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e83610ef2565b6003545f546001600160a01b0391821691163314610b8a5760405162461bcd60e51b815260040161055b90611a9e565b806001600160a01b0316610b9c61118e565b6001600160a01b031614610bc25760405162461bcd60e51b815260040161055b90611b06565b6001600160a01b0384165f9081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff1615610c1157610c0c848484611207565b610629565b6001600160a01b0384165f9081527f13ad2d85210d477fe1a6e25654c8250308cf29b050a4bf0b039d70467486712c602052604090205460ff1615610c5a57610c0c8284611262565b60405162461bcd60e51b815260040161055b90611b6e565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610c9c81610e2c565b61067b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177583610ef2565b6001600160a01b0383165f9081527f9acf0ad725a1e807fc3c987f6db0625f43fbb07e2a89626a1a226c363a0b0bf8602052604090205460ff1615610d1d5760405162461bcd60e51b815260040161055b90611bb0565b5f8211610d3c5760405162461bcd60e51b815260040161055b90611bf2565b6001600160a01b0383165f9081527f32ef73018533fa188e9e42b313c0a4048c6052342b662fb7510c0d1abcea3413602052604090205460ff16610d925760405162461bcd60e51b815260040161055b90611c02565b610d9e83333085611314565b5f6383bece4d60e01b848484604051602401610dbc93929190611c89565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152600354909150610629906001600160a01b0316825f610a1b565b5f82815260026020526040902060010154610e2281610e2c565b6106298383610f9d565b610e36813361136e565b50565b600180546001600160a01b03811691630d3fd67c918591600160a01b90910463ffffffff16906014610e6a83611ce6565b91906101000a81548163ffffffff021916908363ffffffff1602179055508688866040518663ffffffff1660e01b8152600401610eaa9493929190611d72565b60206040518083038185885af1158015610ec6573d5f5f3e3d5ffd5b50505050506040513d601f19601f82011682018060405250810190610eeb9190611dd2565b5050505050565b5f8281526002602090815260408083206001600160a01b038516845290915281205460ff16610f96575f8381526002602090815260408083206001600160a01b03861684529091529020805460ff19166001179055610f4e3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610502565b505f610502565b5f8281526002602090815260408083206001600160a01b038516845290915281205460ff1615610f96575f8381526002602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4506001610502565b5f6040518060600160405280886001600160a01b03168152602001878152602001858152506040516020016110579190611e37565b60408051808303601f19018152919052600180549192506001600160a01b03821691630d3fd67c918591600160a01b900463ffffffff1690601461109a83611ce6565b91906101000a81548163ffffffff021916908363ffffffff1602179055508885886040518663ffffffff1660e01b81526004016110da9493929190611d72565b60206040518083038185885af11580156110f6573d5f5f3e3d5ffd5b50505050506040513d601f19601f8201168201806040525081019061111b9190611dd2565b5050505050505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff1661118c576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f5f5f9054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b8152600401602060405180830381865afa1580156111de573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906112029190611a29565b905090565b6112128382846113cd565b826001600160a01b0316816001600160a01b03167f2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b63988460405161125591906116a3565b60405180910390a3505050565b5f826001600160a01b0316826040515f6040518083038185875af1925050503d805f81146112ab576040519150601f19603f3d011682016040523d82523d5f602084013e6112b0565b606091505b50509050806112d15760405162461bcd60e51b815260040161055b90611e7a565b5f6001600160a01b0316836001600160a01b03167f2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b63988460405161125591906116a3565b61062984856001600160a01b03166323b872dd86868660405160240161133c93929190611e8a565b604051602081830303815290604052915060e01b6020820180516001600160e01b0383818316178352505050506113f3565b5f8281526002602090815260408083206001600160a01b038516845290915290205460ff166106015780826040517fe2517d3f00000000000000000000000000000000000000000000000000000000815260040161055b929190611eb2565b61067b83846001600160a01b031663a9059cbb858560405160240161133c929190611eb2565b5f6114076001600160a01b03841683611464565b905080515f1415801561142b5750808060200190518101906114299190611ee0565b155b1561067b57826040517f5274afe700000000000000000000000000000000000000000000000000000000815260040161055b9190611efd565b606061147183835f611478565b9392505050565b6060814710156114b657306040517fcd78605900000000000000000000000000000000000000000000000000000000815260040161055b9190611efd565b5f5f856001600160a01b031684866040516114d19190611f2c565b5f6040518083038185875af1925050503d805f811461150b576040519150601f19603f3d011682016040523d82523d5f602084013e611510565b606091505b509150915061152086838361152a565b9695505050505050565b60608261153f5761153a82611596565b611471565b815115801561155657506001600160a01b0384163b155b1561158f57836040517f9996b31500000000000000000000000000000000000000000000000000000000815260040161055b9190611efd565b5080611471565b8051156115a65780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160e01b031981165b8114610e36575f5ffd5b8035610502816115d8565b5f6020828403121561160c5761160c5f5ffd5b61147183836115ee565b8015155b82525050565b602081016105028284611616565b5f6001600160a01b038216610502565b6115e48161162e565b80356105028161163e565b5f60208284031215611665576116655f5ffd5b6114718383611647565b806115e4565b80356105028161166f565b5f60208284031215611693576116935f5ffd5b6114718383611675565b8061161a565b60208101610502828461169d565b5f5f604083850312156116c5576116c55f5ffd5b6116cf8484611675565b91506116de8460208501611647565b90509250929050565b5f5f604083850312156116fb576116fb5f5ffd5b6116cf8484611647565b5f5f83601f840112611718576117185f5ffd5b50813567ffffffffffffffff811115611732576117325f5ffd5b60208301915083600182028301111561174c5761174c5f5ffd5b9250929050565b5f5f5f5f5f6060868803121561176a5761176a5f5ffd5b6117748787611647565b9450602086013567ffffffffffffffff811115611792576117925f5ffd5b61179e88828901611705565b9450945050604086013567ffffffffffffffff8111156117bf576117bf5f5ffd5b6117cb88828901611705565b92509250509295509295909350565b5f5f5f606084860312156117ef576117ef5f5ffd5b6117f98585611647565b92506118088560208601611675565b91506118178560408601611647565b90509250925092565b60228152602081017f52656d6f746520627269646765206164647265737320616c726561647920736581527f742e000000000000000000000000000000000000000000000000000000000000602082015290505b60400190565b6020808252810161050281611820565b600f8152602081017f456d707479207472616e736665722e0000000000000000000000000000000000815290505b60200190565b602080825281016105028161188a565b61161a8161162e565b80516118e3838261169d565b50602081015161067b60208401826118ce565b6040810161050282846118d7565b60148152602081017f6661696c65642073656e64696e672076616c7565000000000000000000000000815290506118b8565b6020808252810161050281611904565b5f610502825b67ffffffffffffffff1690565b61161a81611946565b602081016105028284611959565b60198152602081017f546f6b656e20616c72656164792077686974656c697374656400000000000000815290506118b8565b6020808252810161050281611970565b82818337505f910152565b8183526020830192506119d18284836119b2565b50601f01601f19160190565b606081016119eb82886118ce565b81810360208301526119fe8186886119bd565b90508181036040830152611a138184866119bd565b979650505050505050565b80516105028161163e565b5f60208284031215611a3c57611a3c5f5ffd5b6114718383611a1e565b60308152602081017f436f6e74726163742063616c6c6572206973206e6f742074686520726567697381527f7465726564206d657373656e676572210000000000000000000000000000000060208201529050611874565b6020808252810161050281611a46565b60318152602081017f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2081527f696e636f72726563742073656e6465722100000000000000000000000000000060208201529050611874565b6020808252810161050281611aae565b60258152602081017f417474656d7074696e6720746f20776974686472617720756e6b6e6f776e206181527f737365742e00000000000000000000000000000000000000000000000000000060208201529050611874565b6020808252810161050281611b16565b60108152602081017f546f6b656e206973207061757365642e00000000000000000000000000000000815290506118b8565b6020808252810161050281611b7e565b601a8152602081017f417474656d7074696e6720656d707479207472616e736665722e000000000000815290506118b8565b6020808252810161050281611bc0565b6020808252810161050281604e81527f54686973206164647265737320686173206e6f74206265656e20676976656e2060208201527f61207479706520616e64206973207468757320636f6e73696465726564206e6f60408201527f742077686974656c69737465642e000000000000000000000000000000000000606082015260800190565b60608101611c9782866118ce565b611ca4602083018561169d565b611cb160408301846118ce565b949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b63ffffffff165f63fffffffe198201611d0157611d01611cb9565b5060010190565b5f61050263ffffffff831661194c565b61161a81611d08565b63ffffffff811661161a565b8281835e505f910152565b5f611d41825190565b808452602084019350611d58818560208601611d2d565b601f01601f19169290920192915050565b60ff811661161a565b60808101611d808287611d18565b611d8d6020830186611d21565b8181036040830152611d9f8185611d38565b9050611dae6060830184611d69565b95945050505050565b67ffffffffffffffff81166115e4565b805161050281611db7565b5f60208284031215611de557611de55f5ffd5b6114718383611dc7565b80515f906060840190611e0285826118ce565b5060208301518482036020860152611e1a8282611d38565b9150506040830151611e2f604086018261169d565b509392505050565b602080825281016114718184611def565b60148152602081017f4661696c656420746f2073656e64204574686572000000000000000000000000815290506118b8565b6020808252810161050281611e48565b60608101611e9882866118ce565b611ea560208301856118ce565b611cb1604083018461169d565b60408101611ec082856118ce565b611471602083018461169d565b8015156115e4565b805161050281611ecd565b5f60208284031215611ef357611ef35f5ffd5b6114718383611ed5565b6020810161050282846118ce565b5f611f14825190565b611f22818560208601611d2d565b9290920192915050565b6105028183611f0b56fea264697066735822122094684564185cf819396e8f0f9e26cb622d510c487824e5bf5bf35d0445ef967a64736f6c634300081c0033",
}

// TenBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use TenBridgeMetaData.ABI instead.
var TenBridgeABI = TenBridgeMetaData.ABI

// TenBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TenBridgeMetaData.Bin instead.
var TenBridgeBin = TenBridgeMetaData.Bin

// DeployTenBridge deploys a new Ethereum contract, binding an instance of TenBridge to it.
func DeployTenBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TenBridge, error) {
	parsed, err := TenBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TenBridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TenBridge{TenBridgeCaller: TenBridgeCaller{contract: contract}, TenBridgeTransactor: TenBridgeTransactor{contract: contract}, TenBridgeFilterer: TenBridgeFilterer{contract: contract}}, nil
}

// TenBridge is an auto generated Go binding around an Ethereum contract.
type TenBridge struct {
	TenBridgeCaller     // Read-only binding to the contract
	TenBridgeTransactor // Write-only binding to the contract
	TenBridgeFilterer   // Log filterer for contract events
}

// TenBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TenBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TenBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TenBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TenBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TenBridgeSession struct {
	Contract     *TenBridge        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TenBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TenBridgeCallerSession struct {
	Contract *TenBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TenBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TenBridgeTransactorSession struct {
	Contract     *TenBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TenBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TenBridgeRaw struct {
	Contract *TenBridge // Generic contract binding to access the raw methods on
}

// TenBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TenBridgeCallerRaw struct {
	Contract *TenBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// TenBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TenBridgeTransactorRaw struct {
	Contract *TenBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTenBridge creates a new instance of TenBridge, bound to a specific deployed contract.
func NewTenBridge(address common.Address, backend bind.ContractBackend) (*TenBridge, error) {
	contract, err := bindTenBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TenBridge{TenBridgeCaller: TenBridgeCaller{contract: contract}, TenBridgeTransactor: TenBridgeTransactor{contract: contract}, TenBridgeFilterer: TenBridgeFilterer{contract: contract}}, nil
}

// NewTenBridgeCaller creates a new read-only instance of TenBridge, bound to a specific deployed contract.
func NewTenBridgeCaller(address common.Address, caller bind.ContractCaller) (*TenBridgeCaller, error) {
	contract, err := bindTenBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TenBridgeCaller{contract: contract}, nil
}

// NewTenBridgeTransactor creates a new write-only instance of TenBridge, bound to a specific deployed contract.
func NewTenBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*TenBridgeTransactor, error) {
	contract, err := bindTenBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TenBridgeTransactor{contract: contract}, nil
}

// NewTenBridgeFilterer creates a new log filterer instance of TenBridge, bound to a specific deployed contract.
func NewTenBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*TenBridgeFilterer, error) {
	contract, err := bindTenBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TenBridgeFilterer{contract: contract}, nil
}

// bindTenBridge binds a generic wrapper to an already deployed contract.
func bindTenBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TenBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenBridge *TenBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenBridge.Contract.TenBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenBridge *TenBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenBridge.Contract.TenBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenBridge *TenBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenBridge.Contract.TenBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TenBridge *TenBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TenBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TenBridge *TenBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TenBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TenBridge *TenBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TenBridge.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) ADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.ADMINROLE(&_TenBridge.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) ADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.ADMINROLE(&_TenBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.DEFAULTADMINROLE(&_TenBridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TenBridge.Contract.DEFAULTADMINROLE(&_TenBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) ERC20TOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "ERC20_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) ERC20TOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.ERC20TOKENROLE(&_TenBridge.CallOpts)
}

// ERC20TOKENROLE is a free data retrieval call binding the contract method 0x5d872970.
//
// Solidity: function ERC20_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) ERC20TOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.ERC20TOKENROLE(&_TenBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) NATIVETOKENROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "NATIVE_TOKEN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) NATIVETOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.NATIVETOKENROLE(&_TenBridge.CallOpts)
}

// NATIVETOKENROLE is a free data retrieval call binding the contract method 0xe4c3ebc7.
//
// Solidity: function NATIVE_TOKEN_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) NATIVETOKENROLE() ([32]byte, error) {
	return _TenBridge.Contract.NATIVETOKENROLE(&_TenBridge.CallOpts)
}

// SUSPENDEDERC20ROLE is a free data retrieval call binding the contract method 0x5ccc9613.
//
// Solidity: function SUSPENDED_ERC20_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCaller) SUSPENDEDERC20ROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "SUSPENDED_ERC20_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SUSPENDEDERC20ROLE is a free data retrieval call binding the contract method 0x5ccc9613.
//
// Solidity: function SUSPENDED_ERC20_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeSession) SUSPENDEDERC20ROLE() ([32]byte, error) {
	return _TenBridge.Contract.SUSPENDEDERC20ROLE(&_TenBridge.CallOpts)
}

// SUSPENDEDERC20ROLE is a free data retrieval call binding the contract method 0x5ccc9613.
//
// Solidity: function SUSPENDED_ERC20_ROLE() view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) SUSPENDEDERC20ROLE() ([32]byte, error) {
	return _TenBridge.Contract.SUSPENDEDERC20ROLE(&_TenBridge.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridge *TenBridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridge *TenBridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TenBridge.Contract.GetRoleAdmin(&_TenBridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TenBridge *TenBridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TenBridge.Contract.GetRoleAdmin(&_TenBridge.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridge *TenBridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridge *TenBridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TenBridge.Contract.HasRole(&_TenBridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TenBridge *TenBridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TenBridge.Contract.HasRole(&_TenBridge.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridge *TenBridgeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridge *TenBridgeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TenBridge.Contract.SupportsInterface(&_TenBridge.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TenBridge *TenBridgeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TenBridge.Contract.SupportsInterface(&_TenBridge.CallOpts, interfaceId)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_TenBridge *TenBridgeTransactor) Configure(opts *bind.TransactOpts, messengerAddress common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "configure", messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_TenBridge *TenBridgeSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Configure(&_TenBridge.TransactOpts, messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_TenBridge *TenBridgeTransactorSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Configure(&_TenBridge.TransactOpts, messengerAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.GrantRole(&_TenBridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.GrantRole(&_TenBridge.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address owner) returns()
func (_TenBridge *TenBridgeTransactor) Initialize(opts *bind.TransactOpts, messenger common.Address, owner common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "initialize", messenger, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address owner) returns()
func (_TenBridge *TenBridgeSession) Initialize(messenger common.Address, owner common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Initialize(&_TenBridge.TransactOpts, messenger, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address owner) returns()
func (_TenBridge *TenBridgeTransactorSession) Initialize(messenger common.Address, owner common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.Initialize(&_TenBridge.TransactOpts, messenger, owner)
}

// PauseToken is a paid mutator transaction binding the contract method 0x7c41ad2c.
//
// Solidity: function pauseToken(address asset) returns()
func (_TenBridge *TenBridgeTransactor) PauseToken(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "pauseToken", asset)
}

// PauseToken is a paid mutator transaction binding the contract method 0x7c41ad2c.
//
// Solidity: function pauseToken(address asset) returns()
func (_TenBridge *TenBridgeSession) PauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.PauseToken(&_TenBridge.TransactOpts, asset)
}

// PauseToken is a paid mutator transaction binding the contract method 0x7c41ad2c.
//
// Solidity: function pauseToken(address asset) returns()
func (_TenBridge *TenBridgeTransactorSession) PauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.PauseToken(&_TenBridge.TransactOpts, asset)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridge *TenBridgeTransactor) PromoteToAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "promoteToAdmin", newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridge *TenBridgeSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.PromoteToAdmin(&_TenBridge.TransactOpts, newAdmin)
}

// PromoteToAdmin is a paid mutator transaction binding the contract method 0x93b37442.
//
// Solidity: function promoteToAdmin(address newAdmin) returns()
func (_TenBridge *TenBridgeTransactorSession) PromoteToAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.PromoteToAdmin(&_TenBridge.TransactOpts, newAdmin)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridge *TenBridgeTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridge *TenBridgeSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.ReceiveAssets(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_TenBridge *TenBridgeTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.ReceiveAssets(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridge *TenBridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridge *TenBridgeSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RenounceRole(&_TenBridge.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TenBridge *TenBridgeTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RenounceRole(&_TenBridge.TransactOpts, role, callerConfirmation)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_TenBridge *TenBridgeTransactor) RetrieveAllFunds(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "retrieveAllFunds", receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_TenBridge *TenBridgeSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RetrieveAllFunds(&_TenBridge.TransactOpts, receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_TenBridge *TenBridgeTransactorSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RetrieveAllFunds(&_TenBridge.TransactOpts, receiver)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RevokeRole(&_TenBridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TenBridge *TenBridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.RevokeRole(&_TenBridge.TransactOpts, role, account)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridge *TenBridgeTransactor) SendERC20(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "sendERC20", asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridge *TenBridgeSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendERC20(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) payable returns()
func (_TenBridge *TenBridgeTransactorSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendERC20(&_TenBridge.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridge *TenBridgeTransactor) SendNative(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "sendNative", receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridge *TenBridgeSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendNative(&_TenBridge.TransactOpts, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_TenBridge *TenBridgeTransactorSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SendNative(&_TenBridge.TransactOpts, receiver)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridge *TenBridgeTransactor) SetRemoteBridge(opts *bind.TransactOpts, bridge common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "setRemoteBridge", bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridge *TenBridgeSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SetRemoteBridge(&_TenBridge.TransactOpts, bridge)
}

// SetRemoteBridge is a paid mutator transaction binding the contract method 0x16ce8149.
//
// Solidity: function setRemoteBridge(address bridge) returns()
func (_TenBridge *TenBridgeTransactorSession) SetRemoteBridge(bridge common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.SetRemoteBridge(&_TenBridge.TransactOpts, bridge)
}

// UnpauseToken is a paid mutator transaction binding the contract method 0x3b3bff0f.
//
// Solidity: function unpauseToken(address asset) returns()
func (_TenBridge *TenBridgeTransactor) UnpauseToken(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "unpauseToken", asset)
}

// UnpauseToken is a paid mutator transaction binding the contract method 0x3b3bff0f.
//
// Solidity: function unpauseToken(address asset) returns()
func (_TenBridge *TenBridgeSession) UnpauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.UnpauseToken(&_TenBridge.TransactOpts, asset)
}

// UnpauseToken is a paid mutator transaction binding the contract method 0x3b3bff0f.
//
// Solidity: function unpauseToken(address asset) returns()
func (_TenBridge *TenBridgeTransactorSession) UnpauseToken(asset common.Address) (*types.Transaction, error) {
	return _TenBridge.Contract.UnpauseToken(&_TenBridge.TransactOpts, asset)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridge *TenBridgeTransactor) WhitelistToken(opts *bind.TransactOpts, asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridge.contract.Transact(opts, "whitelistToken", asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridge *TenBridgeSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridge.Contract.WhitelistToken(&_TenBridge.TransactOpts, asset, name, symbol)
}

// WhitelistToken is a paid mutator transaction binding the contract method 0x498d82ab.
//
// Solidity: function whitelistToken(address asset, string name, string symbol) returns()
func (_TenBridge *TenBridgeTransactorSession) WhitelistToken(asset common.Address, name string, symbol string) (*types.Transaction, error) {
	return _TenBridge.Contract.WhitelistToken(&_TenBridge.TransactOpts, asset, name, symbol)
}

// TenBridgeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TenBridge contract.
type TenBridgeInitializedIterator struct {
	Event *TenBridgeInitialized // Event containing the contract specifics and raw log

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
func (it *TenBridgeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeInitialized)
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
		it.Event = new(TenBridgeInitialized)
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
func (it *TenBridgeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeInitialized represents a Initialized event raised by the TenBridge contract.
type TenBridgeInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenBridge *TenBridgeFilterer) FilterInitialized(opts *bind.FilterOpts) (*TenBridgeInitializedIterator, error) {

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TenBridgeInitializedIterator{contract: _TenBridge.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TenBridge *TenBridgeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TenBridgeInitialized) (event.Subscription, error) {

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeInitialized)
				if err := _TenBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseInitialized(log types.Log) (*TenBridgeInitialized, error) {
	event := new(TenBridgeInitialized)
	if err := _TenBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TenBridge contract.
type TenBridgeRoleAdminChangedIterator struct {
	Event *TenBridgeRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TenBridgeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeRoleAdminChanged)
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
		it.Event = new(TenBridgeRoleAdminChanged)
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
func (it *TenBridgeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeRoleAdminChanged represents a RoleAdminChanged event raised by the TenBridge contract.
type TenBridgeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TenBridge *TenBridgeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TenBridgeRoleAdminChangedIterator, error) {

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

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeRoleAdminChangedIterator{contract: _TenBridge.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TenBridge *TenBridgeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TenBridgeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeRoleAdminChanged)
				if err := _TenBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseRoleAdminChanged(log types.Log) (*TenBridgeRoleAdminChanged, error) {
	event := new(TenBridgeRoleAdminChanged)
	if err := _TenBridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TenBridge contract.
type TenBridgeRoleGrantedIterator struct {
	Event *TenBridgeRoleGranted // Event containing the contract specifics and raw log

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
func (it *TenBridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeRoleGranted)
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
		it.Event = new(TenBridgeRoleGranted)
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
func (it *TenBridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeRoleGranted represents a RoleGranted event raised by the TenBridge contract.
type TenBridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TenBridgeRoleGrantedIterator, error) {

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

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeRoleGrantedIterator{contract: _TenBridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TenBridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeRoleGranted)
				if err := _TenBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseRoleGranted(log types.Log) (*TenBridgeRoleGranted, error) {
	event := new(TenBridgeRoleGranted)
	if err := _TenBridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TenBridge contract.
type TenBridgeRoleRevokedIterator struct {
	Event *TenBridgeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TenBridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeRoleRevoked)
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
		it.Event = new(TenBridgeRoleRevoked)
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
func (it *TenBridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeRoleRevoked represents a RoleRevoked event raised by the TenBridge contract.
type TenBridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TenBridgeRoleRevokedIterator, error) {

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

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeRoleRevokedIterator{contract: _TenBridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TenBridge *TenBridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TenBridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeRoleRevoked)
				if err := _TenBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseRoleRevoked(log types.Log) (*TenBridgeRoleRevoked, error) {
	event := new(TenBridgeRoleRevoked)
	if err := _TenBridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TenBridgeWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the TenBridge contract.
type TenBridgeWithdrawalIterator struct {
	Event *TenBridgeWithdrawal // Event containing the contract specifics and raw log

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
func (it *TenBridgeWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TenBridgeWithdrawal)
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
		it.Event = new(TenBridgeWithdrawal)
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
func (it *TenBridgeWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TenBridgeWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TenBridgeWithdrawal represents a Withdrawal event raised by the TenBridge contract.
type TenBridgeWithdrawal struct {
	Receiver common.Address
	Asset    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b6398.
//
// Solidity: event Withdrawal(address indexed receiver, address indexed asset, uint256 amount)
func (_TenBridge *TenBridgeFilterer) FilterWithdrawal(opts *bind.FilterOpts, receiver []common.Address, asset []common.Address) (*TenBridgeWithdrawalIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _TenBridge.contract.FilterLogs(opts, "Withdrawal", receiverRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &TenBridgeWithdrawalIterator{contract: _TenBridge.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b6398.
//
// Solidity: event Withdrawal(address indexed receiver, address indexed asset, uint256 amount)
func (_TenBridge *TenBridgeFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *TenBridgeWithdrawal, receiver []common.Address, asset []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _TenBridge.contract.WatchLogs(opts, "Withdrawal", receiverRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TenBridgeWithdrawal)
				if err := _TenBridge.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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
func (_TenBridge *TenBridgeFilterer) ParseWithdrawal(log types.Log) (*TenBridgeWithdrawal, error) {
	event := new(TenBridgeWithdrawal)
	if err := _TenBridge.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
