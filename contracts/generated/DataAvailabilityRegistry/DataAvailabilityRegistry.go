// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DataAvailabilityRegistry

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

// IDataAvailabilityRegistryMetaRollup is an auto generated low-level Go binding around an user-defined struct.
type IDataAvailabilityRegistryMetaRollup struct {
	Hash                [32]byte
	FirstSequenceNumber *big.Int
	LastSequenceNumber  *big.Int
	BlockBindingHash    [32]byte
	BlockBindingNumber  *big.Int
	CrossChainRoot      [32]byte
	LastBatchHash       [32]byte
	Signature           []byte
}

// DataAvailabilityRegistryMetaData contains all meta data concerning the DataAvailabilityRegistry contract.
var DataAvailabilityRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"RollupAdded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structIDataAvailabilityRegistry.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"}],\"name\":\"addRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"enclaveRegistry\",\"outputs\":[{\"internalType\":\"contractINetworkEnclaveRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChallengePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"getRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"FirstSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"BlockBindingHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"BlockBindingNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"crossChainRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"LastBatchHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"Signature\",\"type\":\"bytes\"}],\"internalType\":\"structIDataAvailabilityRegistry.MetaRollup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_merkleMessageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_enclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSeqNo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"merkleMessageBus\",\"outputs\":[{\"internalType\":\"contractIMerkleTreeMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_delay\",\"type\":\"uint256\"}],\"name\":\"setChallengePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461002257610011610026565b6040516123fe61019982396123fe90f35b5f80fd5b61002e610038565b6100366100bc565b565b61003661003661002e565b6100509060401c60ff1690565b90565b6100509054610043565b610050905b6001600160401b031690565b610050905461005d565b61005090610062906001600160401b031682565b9061009c6100506100b892610078565b82546001600160401b0319166001600160401b03919091161790565b9055565b5f6100c5610152565b016100cf81610053565b610141576100dc8161006e565b6001600160401b03919082908116036100f3575050565b816101227fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29361013c9361008c565b604051918291826001600160401b03909116815260200190565b0390a1565b63f92ee8a960e01b5f908152600490fd5b610050610190565b6100506100506100509290565b6100507ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061015a565b61005061016756fe60806040526004361015610011575f80fd5b5f3560e01c8063440c953b146100f05780635d475fdd146100eb5780635fdf31a2146100e65780636fb6a45c146100e1578063715018a6146100dc5780637864b77d146100d757806379ba5097146100d25780637c72dbd0146100cd57806384b0196e146100c85780638da5cb5b146100c3578063c0c53b8b146100be578063e30c3978146100b9578063e874eb20146100b45763f2fde38b036100ff576105f5565b6105c6565b6105a0565b610587565b61051a565b6104d6565b6103c4565b610326565b61030b565b6102f6565b6102c8565b6101d8565b610186565b61012d565b5f9103126100ff57565b5f80fd5b61010e916008021c81565b90565b9061010e9154610103565b61010e5f6001610111565b9052565b565b346100ff5761013d3660046100f5565b61015861014861011c565b6040519182918290815260200190565b0390f35b805b036100ff57565b9050359061012b8261015c565b906020828203126100ff5761010e91610165565b346100ff5761019e610199366004610172565b610660565b604051005b90816101009103126100ff5790565b906020828203126100ff57813567ffffffffffffffff81116100ff5761010e92016101a3565b346100ff5761019e6101eb3660046101b2565b610ddc565b90825f9392825e0152565b61021c61022560209361022f93610210815190565b80835293849260200190565b958691016101f0565b601f01601f191690565b0190565b8051825261010e9161010081019160e09061025360208201516020850152565b61026260408201516040850152565b61027160608201516060850152565b61028060808201516080850152565b61028f60a082015160a0850152565b61029e60c082015160c0850152565b01519060e08184039101526101fb565b901515815260406020820181905261010e92910190610233565b346100ff576102e06102db366004610172565b610fed565b906101586102ed60405190565b928392836102ae565b346100ff576103063660046100f5565b6110a1565b346100ff5761031b3660046100f5565b6101586101486110a6565b346100ff576103363660046100f5565b61019e6110b0565b61010e916008021c5b73ffffffffffffffffffffffffffffffffffffffff1690565b9061010e915461033e565b61010e5f6004610360565b61034761010e61010e9273ffffffffffffffffffffffffffffffffffffffff1690565b61010e90610376565b61010e90610399565b610127906103a2565b60208101929161012b91906103ab565b346100ff576103d43660046100f5565b6101586103df61036b565b604051918291826103b4565b61012790610347565b9061041461040d610403845190565b8084529260200190565b9260200190565b905f5b8181106104245750505090565b90919261044161043a6001928651815260200190565b9460200190565b929101610417565b939591946104b86104b06104c9956104a26104c29561010e9c9a61049560e08c01925f8d01907fff00000000000000000000000000000000000000000000000000000000000000169052565b8a820360208c01526101fb565b9088820360408a01526101fb565b976060870152565b60808501906103eb565b60a0830152565b60c08184039101526103f4565b346100ff576104e63660046100f5565b6101586104f16111b7565b9361050197959793919360405190565b97889788610449565b60208101929161012b91906103eb565b346100ff5761052a3660046100f5565b61015861053561127e565b6040519182918261050a565b61015e81610347565b9050359061012b82610541565b90916060828403126100ff5761010e610570848461054a565b936040610580826020870161054a565b940161054a565b346100ff5761019e61059a366004610557565b91611741565b346100ff576105b03660046100f5565b61015861053561174c565b61010e5f6003610360565b346100ff576105d63660046100f5565b6101586103df6105bb565b906020828203126100ff5761010e9161054a565b346100ff5761019e6106083660046105e1565b6117f0565b61012b906106196117f9565b610655565b905f19905b9181191691161790565b61010e61010e61010e9290565b9061064a61010e6106519261062d565b825461061e565b9055565b61012b90600261063a565b61012b9061060d565b3561010e8161015c565b60208082526026908201527f43616e6e6f742062696e6420746f20667574757265206f722063757272656e7460408201527f20626c6f636b0000000000000000000000000000000000000000000000000000606082015260800190565b156106d757565b60405162461bcd60e51b8152806106f060048201610673565b0390fd5b634e487b7160e01b5f52601160045260245ffd5b9190610713565b9290565b820180921161071e57565b6106f4565b60208082526015908201527f426c6f636b2062696e64696e6720746f6f206f6c640000000000000000000000604082015260600190565b1561076157565b60405162461bcd60e51b8152806106f060048201610723565b60208082526012908201527f556e6b6e6f776e20626c6f636b20686173680000000000000000000000000000604082015260600190565b156107b857565b60405162461bcd60e51b8152806106f06004820161077a565b60208082526016908201527f426c6f636b2062696e64696e67206d69736d6174636800000000000000000000604082015260600190565b1561080f57565b60405162461bcd60e51b8152806106f0600482016107d1565b60208082526014908201527f426c6f622068617368206973206e6f7420736574000000000000000000000000604082015260600190565b1561086657565b60405162461bcd60e51b8152806106f060048201610828565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff8211176108b557604052565b61087f565b903590601e1936829003018212156100ff570180359067ffffffffffffffff82116100ff57602001913682900383136100ff57565b9061012b6108fc60405190565b9283610893565b67ffffffffffffffff81116108b557602090601f01601f19160190565b90825f939282370152565b9092919261094061093b82610903565b6108ef565b93818552818301116100ff5761012b916020850190610920565b61010e91369161092b565b61010e90610347565b61010e9054610965565b80151561015e565b9050519061012b82610978565b906020828203126100ff5761010e91610980565b6040513d5f823e3d90fd5b60208082526019908201527f656e636c6176654944206e6f7420612073657175656e63657200000000000000604082015260600190565b156109ea57565b60405162461bcd60e51b8152806106f0600482016109ac565b610bf1906020610ba9610b8d60808401610a29610a2261010e83610669565b43116106d0565b610a52610a4b61010e610a3b84610669565b610a4560ff61062d565b90610708565b431061075a565b610b75610a5e82610669565b40610a73610a6b5f61062d565b8214156107b1565b610a946060880191610a8e610a8a61010e85610669565b9190565b14610808565b610a9d5f61062d565b49610ab5610aad61010e5f61062d565b82141561085f565b610b697ff1d777b6e4e8b6da895bbd02f40c91ccf99705b363740954e9a795541603ce8591610ae5898b01610669565b938a610af360408201610669565b97610b1e60a0610b17610b11610b0b60c08701610669565b96610669565b93610669565b9301610669565b92610b2860405190565b998a988e8a019889908152610100810198979690959094909390929091602087015260408601526060850152608084015260a083015260c082015260e00152565b90810382520382610893565b610b87610b80825190565b9160200190565b20611817565b610ba3610b9d60e08601866108ba565b9061095a565b9061185b565b610bbb610bb6600461096e565b6103a2565b604051948592839182917f6d46e9870000000000000000000000000000000000000000000000000000000083526004830161050a565b03915afa918215610c455761012b92610c11915f91610c16575b506109e3565b610ca8565b610c38915060203d602011610c3e575b610c308183610893565b81019061098d565b5f610c0b565b503d610c26565b6109a1565b61010e9081565b61010e9054610c4a565b90815260408101929161012b9160200152565b0152565b919061022581610c898161022f9560209181520190565b8095610920565b90815260406020820181905261010e93910191610c72565b610cb181611ac7565b60a08101610cbe81610669565b610cce610a8a61010e5f1961062d565b03610d26575b50610d217fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b124591610d12610d065f61062d565b499160e08101906108ba565b60405191939193849384610c90565b0390a1565b90610d3a610d346002610c51565b42610708565b610d4a610b11610bb6600361096e565b90833b156100ff57610d81935f9283610d6260405190565b809781958294610d7663b6aed0cb60e01b90565b845260048401610c5b565b03925af1908115610c45577fd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b124592610d2192610dbf575b509150610cd4565b610dd6905f610dce8183610893565b8101906100f5565b5f610db7565b61012b90610a03565b61010e6101006108ef565b610df8610de5565b905f825260208080808080808089015f8152015f8152015f8152015f8152015f8152015f8152016060905250565b61010e610df0565b90610e389061010e565b5f5260205260405f2090565b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610e78575b6020831014610e7357565b610e44565b91607f1691610e68565b80545f939291610e9e610e9483610e58565b8085529360200190565b9160018116908115610eed5750600114610eb757505050565b610ec891929394505f5260205f2090565b915f925b818410610ed95750500190565b805484840152602090930192600101610ecc565b92949550505060ff1916825215156020020190565b9061010e91610e82565b9061012b610f2692610f1d60405190565b93848092610f02565b0383610893565b9061012b610fd96007610f3e610de5565b94610f4f610f4b82610c51565b8752565b610f65610f5e60018301610c51565b6020880152565b610f7b610f7460028301610c51565b6040880152565b610f91610f8a60038301610c51565b6060880152565b610fa7610fa060048301610c51565b6080880152565b610fbd610fb660058301610c51565b60a0880152565b610fd3610fcc60068301610c51565b60c0880152565b01610f0c565b60e0840152565b61010e90610f2d565b5190565b90610ff6610e26565b50611009611004835f610e2e565b610fe0565b91611019610a8a61070f85610fe9565b149190565b6110266117f9565b611088565b60208082526034908201527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60408201527f742072656e6f756e6365206f776e657273686970000000000000000000000000606082015260800190565b60405162461bcd60e51b8152806106f06004820161102b565b61101e565b61010e6002610c51565b336110b961174c565b6110cb6110c583610347565b91610347565b036110d95761012b90611b76565b61110a5f917f118cdaa70000000000000000000000000000000000000000000000000000000083526004830161050a565b035ffd5b60208082526015908201527f4549503731323a20556e696e697469616c697a65640000000000000000000000604082015260600190565b1561114c57565b60405162461bcd60e51b8152806106f06004820161110e565b67ffffffffffffffff81116108b55760208091020190565b9061118a61093b83611165565b918252565b369037565b9061012b6111aa6111a48461117d565b93611165565b601f19016020840161118f565b6111fe7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1006111e481610c51565b6111f0610a8a5f61062d565b14908161125e575b50611145565b611206611c37565b9061120f611c64565b90611219306103a2565b6112225f61062d565b61123361122e5f61062d565b611194565b7f0f000000000000000000000000000000000000000000000000000000000000009594934693929190565b61126b9150600101610c51565b611277610a8a5f61062d565b145f6111f8565b61010e5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b0161096e565b61010e9060401c60ff1690565b61010e90546112aa565b61010e905b67ffffffffffffffff1690565b61010e90546112c1565b6112c661010e61010e9290565b9067ffffffffffffffff90610623565b6112c661010e61010e9267ffffffffffffffff1690565b9061132161010e610651926112fa565b82546112ea565b9068ff00000000000000009060401b610623565b151590565b9061135161010e6106519261133c565b8254611328565b610127906112dd565b60208101929161012b9190611358565b909161137b611c8f565b91829161139761139161138d856112b7565b1590565b936112d3565b946113a15f6112dd565b67ffffffffffffffff871614806114b3575b6001966113d06113c2896112dd565b9167ffffffffffffffff1690565b14908161148f575b155b9081611486575b5061145c5761140a92846114015f6113f88a6112dd565b98019788611311565b61144d57611689565b611412575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916114415f610d2193611341565b60405191829182611361565b6114578787611341565b611689565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f6113e1565b90506113da61149d306103a2565b3b6114aa610a8a5f61062d565b149190506113d8565b50836113b3565b61034761010e61010e9290565b61010e906114ba565b6020808252818101527f4d65726b6c65206d657373616765206275732063616e6e6f7420626520307830604082015260600190565b1561150c57565b60405162461bcd60e51b8152806106f0600482016114d0565b6020808252601e908201527f456e636c6176652072656769737472792063616e6e6f74206265203078300000604082015260600190565b1561156357565b60405162461bcd60e51b8152806106f060048201611525565b60208082526013908201527f4f776e65722063616e6e6f742062652030783000000000000000000000000000604082015260600190565b156115ba57565b60405162461bcd60e51b8152806106f06004820161157c565b9061118a61093b83610903565b6115ea60186115d3565b7f44617461417661696c6162696c69747952656769737472790000000000000000602082015290565b61010e6115e0565b61162560016115d3565b7f3100000000000000000000000000000000000000000000000000000000000000602082015290565b61010e61161b565b9073ffffffffffffffffffffffffffffffffffffffff90610623565b9061168261010e610651926103a2565b8254611656565b90610bb661171561171c936116fc611723966116f76116e76116aa5f6114c7565b6116c66116b682610347565b6116bf88610347565b1415611505565b6116e26116d282610347565b6116db8b610347565b141561155c565b610347565b6116f083610347565b14156115b3565b611cb4565b610bb6611707611613565b61170f61164e565b90611cd4565b6003611672565b6004611672565b61012b61172f5f61062d565b61173a81600161063a565b600261063a565b9061012b9291611371565b61010e5f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c006112a4565b61012b906117826117f9565b6117ac817f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00611672565b6117c06117ba610bb661127e565b916103a2565b907f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e227006117eb60405190565b5f90a3565b61012b90611776565b61180161127e565b339061180f6110c583610347565b036110d95750565b61010e90611823611cde565b604291604051917f19010000000000000000000000000000000000000000000000000000000000008352600283015260228201522090565b61010e9161186891611ce6565b90929192611d81565b9061064a6118816106519261010e565b61010e565b915f1960089290920291821b911b610623565b921b90565b91906118af61010e6106519361062d565b908354611886565b61012b915f9161189e565b8181106118cd575050565b806118da5f6001936118b7565b016118c2565b9190601f81116118ef57505050565b6118ff61012b935f5260205f2090565b906020601f840181900483019310611921575b6020601f9091010401906118c2565b9091508190611912565b919067ffffffffffffffff82116108b5576119508261194a8554610e58565b856118e0565b5f90601f83116001146119885761065192915f918361197d575b50505f19600883021c1916906002021790565b013590505f8061196a565b90601f1983169161199c855f5260205f2090565b925f5b8181106119d9575091600293918560019694106119c0575b50505002019055565b01355f19601f84166008021c19165b90555f80806119b7565b9293602060018192878601358155019501930161199f565b9061012b929161192b565b600790611ab461012b93611a1a611a145f8301610669565b84611871565b611a32611a2960208301610669565b6001850161063a565b611a4a611a4160408301610669565b6002850161063a565b611a62611a5960608301610669565b60038501611871565b611a7a611a7160808301610669565b6004850161063a565b611a92611a8960a08301610669565b60058501611871565b611aaa611aa160c08301610669565b60068501611871565b60e08101906108ba565b929091016119f1565b9061012b916119fc565b611ae481611adf5f611ad88161062d565b4990610e2e565b611abd565b611af060208201610669565b611b0d610a8a61010e611b036001610c51565b610a45600161062d565b14611b155750565b611b24604061012b9201610669565b600161063a565b9190600861062391029161189973ffffffffffffffffffffffffffffffffffffffff841b90565b9190611b6361010e610651936103a2565b908354611b2b565b61012b915f91611b52565b61012b90611ba45f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00611b6b565b611e7d565b80545f939291611bbb610e9483610e58565b9160018116908115610eed5750600114611bd457505050565b611be591929394505f5260205f2090565b915f925b818410611bf65750500190565b805484840152602090930192600101611be9565b9061010e91611ba9565b9061012b610f2692611c2560405190565b93848092611c0a565b61010e90611c14565b61010e60027fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1005b01611c2e565b61010e60037fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100611c5e565b61010e611f0b565b61012b90611ca3611f13565b611cac90611f63565b61012b611f74565b61012b90611c97565b9061012b91611cca611f13565b9061012b916120b8565b9061012b91611cbd565b61010e612101565b905f91611cf1825190565b611cfe610a8a604161062d565b03611d2757611d209250602082015190606060408401519301515f1a90612193565b9192909190565b509050611d43611d3e611d395f6114c7565b925190565b61062d565b909160029190565b634e487b7160e01b5f52602160045260245ffd5b60041115611d6957565b611d4b565b9061012b82611d5f565b61010e9061062d565b611d8a5f611d6e565b611d9382611d6e565b03611d9c575050565b611da66001611d6e565b611daf82611d6e565b03611dde577ff645eedf000000000000000000000000000000000000000000000000000000005f908152600490fd5b611de86002611d6e565b611df182611d6e565b03611e32575f61110a611e0384611d78565b7ffce698f700000000000000000000000000000000000000000000000000000000835260048301526024820190565b611e45611e3f6003611d6e565b91611d6e565b14611e4d5750565b7fd78bce0c000000000000000000000000000000000000000000000000000000005f908152600491909152602490fd5b611eb76117ba7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300610bb684611eb18361096e565b92611672565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e06117eb60405190565b61010e7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061062d565b61010e611ee2565b611f1e61138d612252565b611f2457565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f908152600490fd5b61012b90611f5a611f13565b61012b906122cb565b61012b90611f4e565b61012b611f13565b61012b611f6c565b9061012b91611f89611f13565b612058565b90611f97815190565b9067ffffffffffffffff82116108b557611fb58261194a8554610e58565b602090601f8311600114611fed5761065192915f9183611fe25750505f19600883021c1916906002021790565b015190505f8061196a565b601f19831691612000855f5260205f2090565b925f5b818110612036575091600293918560019694106120235750505002019055565b01515f196008601f8516021c19166119cf565b91936020600181928787015181550195019201612003565b9061012b91611f8e565b61209c61012b926120936120897fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10090565b936002850161204e565b6003830161204e565b60016120a75f61062d565b916120b28382611871565b01611871565b9061012b91611f7c565b9095949261012b946120f36120fa926120ec6080966120e560a088019c5f890152565b6020870152565b6040850152565b6060830152565b01906103eb565b7f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f61215561212d6122d4565b612135612368565b90610b69612142306103a2565b60405195869460208601944692866120c2565b612160610b80825190565b2090565b610c6e61012b9461218c606094989795612182608086019a5f870152565b60ff166020850152565b6040830152565b909161219e84611d78565b6121ca610a8a7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a061062d565b1161223f57906121eb6020945f94936121e260405190565b94859485612164565b838052039060015afa15610c45575f516122045f6114c7565b61220d81610347565b61221683610347565b1461222c57506122255f61062d565b90915f9190565b90506122375f61062d565b909160019190565b50505061224b5f6114c7565b9160039190565b61010e5f61225e611c8f565b016112b7565b61012b90612270611f13565b6122795f6114c7565b61228281610347565b61228b83610347565b1461229a575061012b90611b76565b61110a5f917f1e4fbdf70000000000000000000000000000000000000000000000000000000083526004830161050a565b61012b90612264565b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10061230061010e611c37565b90612309825190565b612315610a8a5f61062d565b11156123275750612160610b80825190565b6123319150610c51565b61233a5f61062d565b81146123435790565b507fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47090565b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10061239461010e611c64565b9061239d825190565b6123a9610a8a5f61062d565b11156123bb5750612160610b80825190565b6123319150600101610c5156fea26469706673582212208569c6a1fb6847c2f459e896a61a6eec854825d8da62eeefb9c52993e3fe402064736f6c634300081c0033",
}

// DataAvailabilityRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use DataAvailabilityRegistryMetaData.ABI instead.
var DataAvailabilityRegistryABI = DataAvailabilityRegistryMetaData.ABI

// DataAvailabilityRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DataAvailabilityRegistryMetaData.Bin instead.
var DataAvailabilityRegistryBin = DataAvailabilityRegistryMetaData.Bin

// DeployDataAvailabilityRegistry deploys a new Ethereum contract, binding an instance of DataAvailabilityRegistry to it.
func DeployDataAvailabilityRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DataAvailabilityRegistry, error) {
	parsed, err := DataAvailabilityRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DataAvailabilityRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DataAvailabilityRegistry{DataAvailabilityRegistryCaller: DataAvailabilityRegistryCaller{contract: contract}, DataAvailabilityRegistryTransactor: DataAvailabilityRegistryTransactor{contract: contract}, DataAvailabilityRegistryFilterer: DataAvailabilityRegistryFilterer{contract: contract}}, nil
}

// DataAvailabilityRegistry is an auto generated Go binding around an Ethereum contract.
type DataAvailabilityRegistry struct {
	DataAvailabilityRegistryCaller     // Read-only binding to the contract
	DataAvailabilityRegistryTransactor // Write-only binding to the contract
	DataAvailabilityRegistryFilterer   // Log filterer for contract events
}

// DataAvailabilityRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataAvailabilityRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataAvailabilityRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataAvailabilityRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataAvailabilityRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataAvailabilityRegistrySession struct {
	Contract     *DataAvailabilityRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DataAvailabilityRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataAvailabilityRegistryCallerSession struct {
	Contract *DataAvailabilityRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// DataAvailabilityRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataAvailabilityRegistryTransactorSession struct {
	Contract     *DataAvailabilityRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// DataAvailabilityRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataAvailabilityRegistryRaw struct {
	Contract *DataAvailabilityRegistry // Generic contract binding to access the raw methods on
}

// DataAvailabilityRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryCallerRaw struct {
	Contract *DataAvailabilityRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// DataAvailabilityRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataAvailabilityRegistryTransactorRaw struct {
	Contract *DataAvailabilityRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataAvailabilityRegistry creates a new instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistry(address common.Address, backend bind.ContractBackend) (*DataAvailabilityRegistry, error) {
	contract, err := bindDataAvailabilityRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistry{DataAvailabilityRegistryCaller: DataAvailabilityRegistryCaller{contract: contract}, DataAvailabilityRegistryTransactor: DataAvailabilityRegistryTransactor{contract: contract}, DataAvailabilityRegistryFilterer: DataAvailabilityRegistryFilterer{contract: contract}}, nil
}

// NewDataAvailabilityRegistryCaller creates a new read-only instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistryCaller(address common.Address, caller bind.ContractCaller) (*DataAvailabilityRegistryCaller, error) {
	contract, err := bindDataAvailabilityRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryCaller{contract: contract}, nil
}

// NewDataAvailabilityRegistryTransactor creates a new write-only instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*DataAvailabilityRegistryTransactor, error) {
	contract, err := bindDataAvailabilityRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryTransactor{contract: contract}, nil
}

// NewDataAvailabilityRegistryFilterer creates a new log filterer instance of DataAvailabilityRegistry, bound to a specific deployed contract.
func NewDataAvailabilityRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*DataAvailabilityRegistryFilterer, error) {
	contract, err := bindDataAvailabilityRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryFilterer{contract: contract}, nil
}

// bindDataAvailabilityRegistry binds a generic wrapper to an already deployed contract.
func bindDataAvailabilityRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DataAvailabilityRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataAvailabilityRegistry.Contract.DataAvailabilityRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.DataAvailabilityRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.DataAvailabilityRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataAvailabilityRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _DataAvailabilityRegistry.Contract.Eip712Domain(&_DataAvailabilityRegistry.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _DataAvailabilityRegistry.Contract.Eip712Domain(&_DataAvailabilityRegistry.CallOpts)
}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) EnclaveRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "enclaveRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) EnclaveRegistry() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.EnclaveRegistry(&_DataAvailabilityRegistry.CallOpts)
}

// EnclaveRegistry is a free data retrieval call binding the contract method 0x7c72dbd0.
//
// Solidity: function enclaveRegistry() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) EnclaveRegistry() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.EnclaveRegistry(&_DataAvailabilityRegistry.CallOpts)
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) GetChallengePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "getChallengePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) GetChallengePeriod() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.GetChallengePeriod(&_DataAvailabilityRegistry.CallOpts)
}

// GetChallengePeriod is a free data retrieval call binding the contract method 0x7864b77d.
//
// Solidity: function getChallengePeriod() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) GetChallengePeriod() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.GetChallengePeriod(&_DataAvailabilityRegistry.CallOpts)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) GetRollupByHash(opts *bind.CallOpts, rollupHash [32]byte) (bool, IDataAvailabilityRegistryMetaRollup, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "getRollupByHash", rollupHash)

	if err != nil {
		return *new(bool), *new(IDataAvailabilityRegistryMetaRollup), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(IDataAvailabilityRegistryMetaRollup)).(*IDataAvailabilityRegistryMetaRollup)

	return out0, out1, err

}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) GetRollupByHash(rollupHash [32]byte) (bool, IDataAvailabilityRegistryMetaRollup, error) {
	return _DataAvailabilityRegistry.Contract.GetRollupByHash(&_DataAvailabilityRegistry.CallOpts, rollupHash)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x6fb6a45c.
//
// Solidity: function getRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes))
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) GetRollupByHash(rollupHash [32]byte) (bool, IDataAvailabilityRegistryMetaRollup, error) {
	return _DataAvailabilityRegistry.Contract.GetRollupByHash(&_DataAvailabilityRegistry.CallOpts, rollupHash)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) LastBatchSeqNo(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "lastBatchSeqNo")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) LastBatchSeqNo() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.LastBatchSeqNo(&_DataAvailabilityRegistry.CallOpts)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) LastBatchSeqNo() (*big.Int, error) {
	return _DataAvailabilityRegistry.Contract.LastBatchSeqNo(&_DataAvailabilityRegistry.CallOpts)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) MerkleMessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "merkleMessageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) MerkleMessageBus() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.MerkleMessageBus(&_DataAvailabilityRegistry.CallOpts)
}

// MerkleMessageBus is a free data retrieval call binding the contract method 0xe874eb20.
//
// Solidity: function merkleMessageBus() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) MerkleMessageBus() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.MerkleMessageBus(&_DataAvailabilityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) Owner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.Owner(&_DataAvailabilityRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) Owner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.Owner(&_DataAvailabilityRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) PendingOwner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.PendingOwner(&_DataAvailabilityRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) PendingOwner() (common.Address, error) {
	return _DataAvailabilityRegistry.Contract.PendingOwner(&_DataAvailabilityRegistry.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _DataAvailabilityRegistry.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) RenounceOwnership() error {
	return _DataAvailabilityRegistry.Contract.RenounceOwnership(&_DataAvailabilityRegistry.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryCallerSession) RenounceOwnership() error {
	return _DataAvailabilityRegistry.Contract.RenounceOwnership(&_DataAvailabilityRegistry.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AcceptOwnership(&_DataAvailabilityRegistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AcceptOwnership(&_DataAvailabilityRegistry.TransactOpts)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) AddRollup(opts *bind.TransactOpts, r IDataAvailabilityRegistryMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "addRollup", r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) AddRollup(r IDataAvailabilityRegistryMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AddRollup(&_DataAvailabilityRegistry.TransactOpts, r)
}

// AddRollup is a paid mutator transaction binding the contract method 0x5fdf31a2.
//
// Solidity: function addRollup((bytes32,uint256,uint256,bytes32,uint256,bytes32,bytes32,bytes) r) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) AddRollup(r IDataAvailabilityRegistryMetaRollup) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.AddRollup(&_DataAvailabilityRegistry.TransactOpts, r)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry, address _owner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) Initialize(opts *bind.TransactOpts, _merkleMessageBus common.Address, _enclaveRegistry common.Address, _owner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "initialize", _merkleMessageBus, _enclaveRegistry, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry, address _owner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) Initialize(_merkleMessageBus common.Address, _enclaveRegistry common.Address, _owner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.Initialize(&_DataAvailabilityRegistry.TransactOpts, _merkleMessageBus, _enclaveRegistry, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _merkleMessageBus, address _enclaveRegistry, address _owner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) Initialize(_merkleMessageBus common.Address, _enclaveRegistry common.Address, _owner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.Initialize(&_DataAvailabilityRegistry.TransactOpts, _merkleMessageBus, _enclaveRegistry, _owner)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) SetChallengePeriod(opts *bind.TransactOpts, _delay *big.Int) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "setChallengePeriod", _delay)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) SetChallengePeriod(_delay *big.Int) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.SetChallengePeriod(&_DataAvailabilityRegistry.TransactOpts, _delay)
}

// SetChallengePeriod is a paid mutator transaction binding the contract method 0x5d475fdd.
//
// Solidity: function setChallengePeriod(uint256 _delay) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) SetChallengePeriod(_delay *big.Int) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.SetChallengePeriod(&_DataAvailabilityRegistry.TransactOpts, _delay)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.TransferOwnership(&_DataAvailabilityRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DataAvailabilityRegistry.Contract.TransferOwnership(&_DataAvailabilityRegistry.TransactOpts, newOwner)
}

// DataAvailabilityRegistryEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryEIP712DomainChangedIterator struct {
	Event *DataAvailabilityRegistryEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryEIP712DomainChanged)
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
		it.Event = new(DataAvailabilityRegistryEIP712DomainChanged)
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
func (it *DataAvailabilityRegistryEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryEIP712DomainChanged represents a EIP712DomainChanged event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*DataAvailabilityRegistryEIP712DomainChangedIterator, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryEIP712DomainChangedIterator{contract: _DataAvailabilityRegistry.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryEIP712DomainChanged)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseEIP712DomainChanged(log types.Log) (*DataAvailabilityRegistryEIP712DomainChanged, error) {
	event := new(DataAvailabilityRegistryEIP712DomainChanged)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataAvailabilityRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryInitializedIterator struct {
	Event *DataAvailabilityRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryInitialized)
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
		it.Event = new(DataAvailabilityRegistryInitialized)
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
func (it *DataAvailabilityRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryInitialized represents a Initialized event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*DataAvailabilityRegistryInitializedIterator, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryInitializedIterator{contract: _DataAvailabilityRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryInitialized)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseInitialized(log types.Log) (*DataAvailabilityRegistryInitialized, error) {
	event := new(DataAvailabilityRegistryInitialized)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataAvailabilityRegistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferStartedIterator struct {
	Event *DataAvailabilityRegistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryOwnershipTransferStarted)
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
		it.Event = new(DataAvailabilityRegistryOwnershipTransferStarted)
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
func (it *DataAvailabilityRegistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DataAvailabilityRegistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryOwnershipTransferStartedIterator{contract: _DataAvailabilityRegistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryOwnershipTransferStarted)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*DataAvailabilityRegistryOwnershipTransferStarted, error) {
	event := new(DataAvailabilityRegistryOwnershipTransferStarted)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataAvailabilityRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferredIterator struct {
	Event *DataAvailabilityRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryOwnershipTransferred)
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
		it.Event = new(DataAvailabilityRegistryOwnershipTransferred)
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
func (it *DataAvailabilityRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DataAvailabilityRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryOwnershipTransferredIterator{contract: _DataAvailabilityRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryOwnershipTransferred)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*DataAvailabilityRegistryOwnershipTransferred, error) {
	event := new(DataAvailabilityRegistryOwnershipTransferred)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataAvailabilityRegistryRollupAddedIterator is returned from FilterRollupAdded and is used to iterate over the raw logs and unpacked data for RollupAdded events raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryRollupAddedIterator struct {
	Event *DataAvailabilityRegistryRollupAdded // Event containing the contract specifics and raw log

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
func (it *DataAvailabilityRegistryRollupAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataAvailabilityRegistryRollupAdded)
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
		it.Event = new(DataAvailabilityRegistryRollupAdded)
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
func (it *DataAvailabilityRegistryRollupAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataAvailabilityRegistryRollupAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataAvailabilityRegistryRollupAdded represents a RollupAdded event raised by the DataAvailabilityRegistry contract.
type DataAvailabilityRegistryRollupAdded struct {
	RollupHash [32]byte
	Signature  []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRollupAdded is a free log retrieval operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) FilterRollupAdded(opts *bind.FilterOpts) (*DataAvailabilityRegistryRollupAddedIterator, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.FilterLogs(opts, "RollupAdded")
	if err != nil {
		return nil, err
	}
	return &DataAvailabilityRegistryRollupAddedIterator{contract: _DataAvailabilityRegistry.contract, event: "RollupAdded", logs: logs, sub: sub}, nil
}

// WatchRollupAdded is a free log subscription operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) WatchRollupAdded(opts *bind.WatchOpts, sink chan<- *DataAvailabilityRegistryRollupAdded) (event.Subscription, error) {

	logs, sub, err := _DataAvailabilityRegistry.contract.WatchLogs(opts, "RollupAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataAvailabilityRegistryRollupAdded)
				if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "RollupAdded", log); err != nil {
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

// ParseRollupAdded is a log parse operation binding the contract event 0xd0fa8825abde6f2b225ec23c2fb943dd8b2414208adf66eb6da3f24be40b1245.
//
// Solidity: event RollupAdded(bytes32 rollupHash, bytes signature)
func (_DataAvailabilityRegistry *DataAvailabilityRegistryFilterer) ParseRollupAdded(log types.Log) (*DataAvailabilityRegistryRollupAdded, error) {
	event := new(DataAvailabilityRegistryRollupAdded)
	if err := _DataAvailabilityRegistry.contract.UnpackLog(event, "RollupAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
