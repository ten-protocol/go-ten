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
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERC20_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NATIVE_TOKEN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUSPENDED_ERC20_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"contractIMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messenger\",\"outputs\":[{\"internalType\":\"contractICrossChainMessenger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"pauseToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"promoteToAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"remoteBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridge\",\"type\":\"address\"}],\"name\":\"setRemoteBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"unpauseToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"whitelistToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052346019576040516124c161001e82396124c190f35b5f80fdfe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a7146101905780630f0a9a4b1461018b57806316ce8149146101865780631888d71214610181578063248a9ca31461017c5780632f2ff15d1461017757806336568abe146101725780633b3bff0f1461016d5780633cb747bf14610168578063485cc95514610163578063498d82ab1461015e5780635ccc9613146101595780635d8729701461015457806375b238fc1461014f5780637c41ad2c1461014a57806383bece4d1461014557806391d148541461014057806393b374421461013b578063a1a227fa14610136578063a217fddf14610131578063a381c8e21461012c578063affed0e014610127578063d547741f146101225763e4c3ebc7036101a8576106c7565b6106ae565b610693565b61067f565b610664565b61062b565b610613565b6105f7565b6105de565b61059d565b610564565b61052b565b6104f2565b6104d6565b610433565b6103f0565b61038a565b610371565b610358565b610309565b6102cf565b6102b2565b610261565b6101d6565b6001600160e01b031981165b036101a857565b5f80fd5b905035906101b982610195565b565b906020828203126101a8576101cf916101ac565b90565b9052565b346101a8576102036101f16101ec3660046101bb565b610700565b60405191829182901515815260200190565b0390f35b5f9103126101a857565b6101cf916008021c5b73ffffffffffffffffffffffffffffffffffffffff1690565b906101cf9154610211565b6101cf5f80610233565b6101d29061021a565b6020810192916101b99190610248565b346101a857610271366004610207565b61020361027c61023e565b60405191829182610251565b6101a18161021a565b905035906101b982610288565b906020828203126101a8576101cf91610291565b346101a8576102ca6102c536600461029e565b610913565b604051005b6102ca6102dd36600461029e565b610a3b565b806101a1565b905035906101b9826102e2565b906020828203126101a8576101cf916102e8565b346101a85761020361032461031f3660046102f5565b610ad9565b6040515b9182918290815260200190565b91906040838203126101a8576101cf90602061035182866102e8565b9401610291565b346101a8576102ca61036b366004610335565b90610b37565b346101a8576102ca610384366004610335565b90610b41565b346101a8576102ca61039d36600461029e565b610bea565b61021a6101cf6101cf9273ffffffffffffffffffffffffffffffffffffffff1690565b6101cf906103a2565b6101cf906103c5565b6101d2906103ce565b6020810192916101b991906103d7565b346101a857610400366004610207565b61020361040b610c6a565b604051918291826103e0565b91906040838203126101a8576101cf9060206103518286610291565b346101a8576102ca610446366004610417565b90611020565b909182601f830112156101a85781359167ffffffffffffffff83116101a85760200192600183028401116101a857565b6060818303126101a8576104908282610291565b92602082013567ffffffffffffffff81116101a857836104b191840161044c565b929093604082013567ffffffffffffffff81116101a8576104d2920161044c565b9091565b346101a8576102ca6104e936600461047c565b9392909261123e565b346101a857610502366004610207565b6102037fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e610324565b346101a85761053b366004610207565b6102037f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a610324565b346101a857610574366004610207565b6102037fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610324565b346101a8576102ca6105b036600461029e565b61132a565b90916060828403126101a8576101cf6105ce8484610291565b93604061035182602087016102e8565b346101a8576102ca6105f13660046105b5565b916115d7565b346101a8576102036101f161060d366004610335565b906115ff565b346101a8576102ca61062636600461029e565b6116ee565b346101a85761063b366004610207565b61020361040b61173a565b6101cf6101cf6101cf9290565b6101cf5f610646565b6101cf610653565b346101a857610674366004610207565b61020361032461065c565b6102ca61068d3660046105b5565b916118ac565b346101a8576106a3366004610207565b6102036103246119da565b346101a8576102ca6106c1366004610335565b90611a07565b346101a8576106d7366004610207565b6102037fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad2110578610324565b7f7965db0b000000000000000000000000000000000000000000000000000000006001600160e01b0319821614908115610738575090565b6101cf91506001600160e01b0319167f01ffc9a7000000000000000000000000000000000000000000000000000000001490565b6101b99061079e7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611a15565b611a15565b610859565b61021a6101cf6101cf9290565b6101cf906107a3565b0190565b156107c457565b60405162461bcd60e51b815260206004820152601460248201527f4272696467652063616e6e6f74206265203078300000000000000000000000006044820152606490fd5b6101cf9061021a565b6101cf9054610809565b9073ffffffffffffffffffffffffffffffffffffffff905b9181191691161790565b9061084e6101cf610855926103ce565b825461081c565b9055565b6108625f6107b0565b61087e61086e8261021a565b6108778461021a565b14156107bd565b61089961089361088d5f610812565b9261021a565b9161021a565b036108a8576101b9905f61083e565b60405162461bcd60e51b815260206004820152602260248201527f52656d6f746520627269646765206164647265737320616c726561647920736560448201527f742e0000000000000000000000000000000000000000000000000000000000006064820152608490fd5b6101b99061076c565b1561092357565b60405162461bcd60e51b815260206004820152600f60248201527f456d707479207472616e736665722e00000000000000000000000000000000006044820152606490fd5b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761099e57604052565b610968565b906101b96109b060405190565b928361097c565b6101cf60406109a3565b906101d29061021a565b805182526101b99190602090810151910190610248565b6040810192916101b991906109cb565b634e487b7160e01b5f52602160045260245ffd5b60031115610a1057565b6109f2565b906101b982610a06565b6101cf90610a15565b610a356101cf6101cf9290565b60ff1690565b6101b990610a52610a4b5f610646565b341161091c565b610a6f610a5d6109b7565b91610a66348452565b602083016109c1565b610a97610a7b60405190565b8092610a8b6020830191826109e2565b9081038252038261097c565b610aa16002610a1f565b610aaa5f610646565b90610ab45f610a28565b92611b27565b905b5f5260205260405f2090565b6101cf9081565b6101cf9054610ac8565b6001610b136101cf92610ae95f90565b505f7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268005b01610aba565b01610acf565b906101b991610b2a61079982610ad9565b90610b3491611bf8565b50565b906101b991610b19565b90610b503361021a565b61021a565b610b598261021a565b03610b6757610b3491611c9e565b7f6697b232000000000000000000000000000000000000000000000000000000005f90815260045b035ffd5b6101b990610bc07fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611a15565b610b34907fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e611c9e565b6101b990610b93565b67ffffffffffffffff811161099e57602090601f01601f19160190565b90610c22610c1d83610bf3565b6109a3565b918252565b610c31601e610c10565b7f43726f7373436861696e456e61626c656454454e2e6d657373656e6765720000602082015290565b6101cf610c27565b6101cf610c5a565b6101cf610c9c5f610c966101cf610c7f610c62565b80516020918201205f19015f9081522060ff191690565b01610812565b6103ce565b6101cf9060401c610a35565b6101cf9054610ca1565b6101cf905b67ffffffffffffffff1690565b6101cf9054610cb7565b610cbc6101cf6101cf9290565b9067ffffffffffffffff90610834565b610cbc6101cf6101cf9267ffffffffffffffff1690565b90610d176101cf61085592610cf0565b8254610ce0565b9068ff00000000000000009060401b610834565b90610d426101cf61085592151590565b8254610d1e565b6101d290610cd3565b6020810192916101b99190610d49565b610d6a611d1c565b908190610d86610d80610d7c84610cad565b1590565b92610cc9565b93610d905f610cd3565b67ffffffffffffffff86161480610eab575b600195610dbf610db188610cd3565b9167ffffffffffffffff1690565b149081610e83575b155b9081610e7a575b50610e5057610df99183610df05f610de789610cd3565b97019687610d07565b610e4157610f4a565b610e01575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291610e305f610e3c93610d32565b60405191829182610d52565b0390a1565b610e4b8686610d32565b610f4a565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f610dd0565b9050610dc9610e91306103ce565b3b610ea2610e9e5f610646565b9190565b14919050610dc7565b5082610da2565b15610eb957565b60405162461bcd60e51b815260206004820152601760248201527f4d657373656e6765722063616e6e6f74206265203078300000000000000000006044820152606490fd5b15610f0557565b60405162461bcd60e51b815260206004820152601360248201527f4f776e65722063616e6e6f7420626520307830000000000000000000000000006044820152606490fd5b610ff9610b3492610f9b610f5d5f6107b0565b93610f7a610f6a8661021a565b610f738361021a565b1415610eb2565b610f96610f868661021a565b610f8f8561021a565b1415610efe565b611d45565b610fa3611e1a565b610fab610653565b610fb58282611bf8565b50610fc8610fc2306103ce565b82611bf8565b50610ff47fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217759182611e30565b611bf8565b507fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad2110578611bf8565b906101b991610d62565b6101b9949392919061105b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611a15565b611163565b1561106757565b60405162461bcd60e51b815260206004820152601360248201527f41737365742063616e6e6f7420626520307830000000000000000000000000006044820152606490fd5b156110b357565b60405162461bcd60e51b815260206004820152601960248201527f546f6b656e20616c72656164792077686974656c6973746564000000000000006044820152606490fd5b90825f939282370152565b91906111218161111a816107b99560209181520190565b80956110f8565b601f01601f191690565b9391906101cf95936111559261114860608801935f890190610248565b8683036020880152611103565b926040818503910152611103565b90926101b9946112009161120f94611190611180610b4b5f6107b0565b6111898761021a565b1415611060565b6111ca857f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a610ff46111c5610d7c84846115ff565b6110ac565b506040519687956004602088017f458ffd630000000000000000000000000000000000000000000000000000000081520161112b565b6020820181038252038261097c565b6112185f610812565b6112226001610a1f565b9161122c5f610646565b9283916112385f610a28565b93611f4d565b906101b99493929161102a565b6101b9906112787fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775611a15565b6112c9565b1561128457565b60405162461bcd60e51b815260206004820152601860248201527f546f6b656e206973206e6f742077686974656c697374656400000000000000006044820152606490fd5b610b34906113046112ff827f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a6115ff565b6115ff565b61127d565b7fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e611bf8565b6101b99061124b565b1561133a57565b60405162461bcd60e51b815260206004820152603060248201527f436f6e74726163742063616c6c6572206973206e6f742074686520726567697360448201527f7465726564206d657373656e67657221000000000000000000000000000000006064820152608490fd5b156113ac57565b60405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f72726563742073656e646572210000000000000000000000000000006064820152608490fd5b1561141e57565b60405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f727265637420746172676574210000000000000000000000000000006064820152608490fd5b906101b992916114bf61149b5f610812565b6114ab6114a6611fce565b611333565b6114b961089361088d611fe5565b146113a5565b6114df6114ca612053565b6114d9610893610b4b306103ce565b14611417565b906114f292916114ed6120bb565b6114fa565b6101b9612106565b611524817f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a6115ff565b1561153457906101b992916121e9565b61155e907fd2fb17ceaa388942529b17e0006ffc4d559f040dd4f2157b8070f17ad21105786115ff565b1561156c576101b99161217f565b60405162461bcd60e51b815260206004820152602560248201527f417474656d7074696e6720746f20776974686472617720756e6b6e6f776e206160448201527f737365742e0000000000000000000000000000000000000000000000000000006064820152608490fd5b906101b99291611489565b90610abc906103ce565b6101cf90610a35565b6101cf90546115ec565b6101cf915f61163a611640936116125f90565b50827f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800610b0d565b016115e2565b6115f5565b6101b990611654610799610653565b6116a5565b1561166057565b60405162461bcd60e51b815260206004820152601760248201527f4e65772061646d696e2063616e6e6f74206265203078300000000000000000006044820152606490fd5b6101b9906116c86116b8610b4b5f6107b0565b6116c18361021a565b1415611659565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610b37565b6101b990611645565b611701601f610c10565b7f43726f7373436861696e456e61626c656454454e2e6d65737361676542757300602082015290565b6101cf6116f7565b6101cf61172a565b6101cf610c9c5f610c966101cf610c7f611732565b1561175657565b60405162461bcd60e51b815260206004820152601060248201527f546f6b656e206973207061757365642e000000000000000000000000000000006044820152606490fd5b156117a257565b60405162461bcd60e51b815260206004820152601a60248201527f417474656d7074696e6720656d707479207472616e736665722e0000000000006044820152606490fd5b156117ee57565b60405162461bcd60e51b815260206004820152604e60248201527f54686973206164647265737320686173206e6f74206265656e20676976656e2060448201527f61207479706520616e64206973207468757320636f6e73696465726564206e6f60648201527f742077686974656c69737465642e000000000000000000000000000000000000608482015260a490fd5b6040906118a56101b9949695939661189e60608401985f850190610248565b6020830152565b0190610248565b611985906112006101b9946118ed6118e8610d7c856112fa7fe08e2b666d5448741eeecf0d9bbc95ce21b6e73cca0d67afab4279e53c07cd3e90565b61174f565b6119006118f95f610646565b861161179b565b61193261192d847f9f225881f6e7ac8a885b63aa2269cbce78dd6a669864ccd2cd2517a8e709d73a6115ff565b6117e7565b6119508561193f856103ce565b611948306103ce565b90339061227c565b6040519485936004602086017f83bece4d0000000000000000000000000000000000000000000000000000000081520161187f565b61198e5f610812565b6112225f610a1f565b6119a1601a610c10565b7f43726f7373436861696e456e61626c656454454e2e6e6f6e6365000000000000602082015290565b6101cf611997565b6101cf6119ca565b6101cf5f610b136101cf610c7f6119d2565b906101b9916119fd61079982610ad9565b90610b3491611c9e565b906101b9916119ec565b1490565b6101b99033906122e0565b634e487b7160e01b5f52601160045260245ffd5b5f198114611a425760010190565b611a20565b905f1990610834565b90611a606101cf61085592610646565b8254611a47565b67ffffffffffffffff81166101a1565b905051906101b982611a67565b906020828203126101a8576101cf91611a77565b90825f9392825e0152565b611ac46111216020936107b993611ab8815190565b80835293849260200190565b95869101611a98565b9493916060916101b994611b05611b1293611af860808b01945f8c019067ffffffffffffffff169052565b63ffffffff1660208a0152565b8782036040890152611aa3565b94019060ff169052565b6040513d5f823e3d90fd5b9092611b8f602093611b3a610c9c61173a565b92611b9a630d3fd67c91611b785f611b566101cf610c7f6119d2565b01611b73611b6382610acf565b91611b6d83611a34565b90611a50565b610cd3565b96611b8260405190565b998a988997889660e01b90565b865260048601611acd565b03925af18015611bd357611bab5750565b610b349060203d602011611bcc575b611bc4818361097c565b810190611a84565b503d611bba565b611b1c565b9060ff90610834565b90611bf16101cf61085592151590565b8254611bd8565b7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800611c26610d7c84846115ff565b15611c97576001611c41845f61163a8682611c469701610aba565b611be1565b611c64611c5e611c58339390565b9390565b936103ce565b916103ce565b917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d611c8f60405190565b5f90a4600190565b5050505f90565b7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800611cc983836115ff565b15611c97575f611c41848261163a8682611ce39701610aba565b611cf1611c5e611c58339390565b917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b611c8f60405190565b6101cf612350565b905051906101b982610288565b906020828203126101a8576101cf91611d24565b611d5e905f611d586101cf610c7f610c62565b0161083e565b611d69610c9c610c6a565b6020611d7460405190565b7fa1a227fa00000000000000000000000000000000000000000000000000000000815291829060049082905afa8015611bd357611dc2915f91611de3575b505f611d586101cf610c7f611732565b6101b9611dce5f610646565b5f611ddd6101cf610c7f6119d2565b01611a50565b611e05915060203d602011611e0b575b611dfd818361097c565b810190611d31565b5f611db2565b503d611df3565b6101b9612358565b6101b9611e12565b90611a606101cf6108559290565b90611e7d610e9e611c547f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800946101cf856001611e77845f611e7082610ad9565b9b01610aba565b01611e22565b917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff611ea860405190565b5f90a4565b15611eb457565b60405162461bcd60e51b815260206004820152601460248201527f5461726765742063616e6e6f74206265203078300000000000000000000000006044820152606490fd5b6101cf60606109a3565b906101cf90604080611f3460608401611f225f8801515f870190610248565b60208701518582036020870152611aa3565b940151910152565b60208082526101cf92910190611f03565b9294611f9c6020959396611f95611b8f94611f7d611f6d610b4b5f6107b0565b611f768a61021a565b1415611ead565b611f8f611f88611ef9565b98896109c1565b88880152565b6040860152565b611fc3611fa860405190565b8095611fb78883019182611f3c565b9081038252038561097c565b611b3a610c9c61173a565b611fdc610b4b610c9c610c6a565b611a113361021a565b611ff0610c9c610c6a565b6020611ffb60405190565b9182907f63012de5000000000000000000000000000000000000000000000000000000005b825260049082905afa908115611bd3575f9161203a575090565b6101cf915060203d602011611e0b57611dfd818361097c565b61205e610c9c610c6a565b602061206960405190565b9182907fb859ce8300000000000000000000000000000000000000000000000000000000612020565b6101cf7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00610646565b6120c3612393565b6120dc576101b960016120d76101cf612092565b6123a5565b7f3ee5aeb5000000000000000000000000000000000000000000000000000000005f908152600490fd5b6101b95f6120d76101cf612092565b3d1561212e576121243d610c10565b903d5f602084013e565b606090565b1561213a57565b60405162461bcd60e51b815260206004820152601460248201527f4661696c656420746f2073656e642045746865720000000000000000000000006044820152606490fd5b6121a15f8061218d60405190565b5f9086865af161219b612115565b50612133565b7f2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b63986121e46121da611c586121d45f6107b0565b946103ce565b9361032860405190565b0390a3565b90916121e46121da611c587f2717ead6b9200dd235aad468c9809ea400fe33ac69b5bfaa6d3e90fc922b639893610c9c8782612224896103ce565b6123a8565b61224261223c6101cf9263ffffffff1690565b60e01b90565b6001600160e01b03191690565b6040906122786101b9949695939661226e60608401985f850190610248565b6020830190610248565b0152565b906122c1906122b26101b9956004956122986323b872dd612229565b936122a260405190565b978895602087019081520161224f565b6020820181038252038361097c565b6123dc565b9160206101b992949361227860408201965f830190610248565b906122ee610d7c82846115ff565b6122f6575050565b7fe2517d3f000000000000000000000000000000000000000000000000000000005f90815291610b8f9160046122c6565b6101cf7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610646565b6101cf612327565b612363610d7c612479565b61236957565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f908152600490fd5b6101cf6123a16101cf612092565b5c90565b5d565b6122c16004926122b26101b9956123c263a9059cbb612229565b926123cc60405190565b96879460208601908152016122c6565b905f6020916123e85f90565b50828151910182855af115611b1c573d5f5190612407610e9e5f610646565b036124655750612416816103ce565b3b612423610e9e5f610646565b145b61242c5750565b610b8f6124395f926103ce565b7f5274afe700000000000000000000000000000000000000000000000000000000835260048301610251565b612472610e9e6001610646565b1415612425565b6101cf5f612485611d1c565b01610cad56fea2646970667358221220ba96b207c630b6df8038934a98afe173875a9fa9a8d4bf1efcc1afe565cd55bb64736f6c634300081c0033",
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

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_TenBridge *TenBridgeCaller) MessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "messageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_TenBridge *TenBridgeSession) MessageBus() (common.Address, error) {
	return _TenBridge.Contract.MessageBus(&_TenBridge.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_TenBridge *TenBridgeCallerSession) MessageBus() (common.Address, error) {
	return _TenBridge.Contract.MessageBus(&_TenBridge.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_TenBridge *TenBridgeCaller) Messenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "messenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_TenBridge *TenBridgeSession) Messenger() (common.Address, error) {
	return _TenBridge.Contract.Messenger(&_TenBridge.CallOpts)
}

// Messenger is a free data retrieval call binding the contract method 0x3cb747bf.
//
// Solidity: function messenger() view returns(address)
func (_TenBridge *TenBridgeCallerSession) Messenger() (common.Address, error) {
	return _TenBridge.Contract.Messenger(&_TenBridge.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_TenBridge *TenBridgeCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_TenBridge *TenBridgeSession) Nonce() (*big.Int, error) {
	return _TenBridge.Contract.Nonce(&_TenBridge.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_TenBridge *TenBridgeCallerSession) Nonce() (*big.Int, error) {
	return _TenBridge.Contract.Nonce(&_TenBridge.CallOpts)
}

// RemoteBridgeAddress is a free data retrieval call binding the contract method 0x0f0a9a4b.
//
// Solidity: function remoteBridgeAddress() view returns(address)
func (_TenBridge *TenBridgeCaller) RemoteBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TenBridge.contract.Call(opts, &out, "remoteBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RemoteBridgeAddress is a free data retrieval call binding the contract method 0x0f0a9a4b.
//
// Solidity: function remoteBridgeAddress() view returns(address)
func (_TenBridge *TenBridgeSession) RemoteBridgeAddress() (common.Address, error) {
	return _TenBridge.Contract.RemoteBridgeAddress(&_TenBridge.CallOpts)
}

// RemoteBridgeAddress is a free data retrieval call binding the contract method 0x0f0a9a4b.
//
// Solidity: function remoteBridgeAddress() view returns(address)
func (_TenBridge *TenBridgeCallerSession) RemoteBridgeAddress() (common.Address, error) {
	return _TenBridge.Contract.RemoteBridgeAddress(&_TenBridge.CallOpts)
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
