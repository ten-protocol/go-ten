// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package NetworkEnclaveRegistry

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

// NetworkEnclaveRegistryMetaData contains all meta data concerning the NetworkEnclaveRegistry contract.
var NetworkEnclaveRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"NetworkSecretInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"NetworkSecretRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"attester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"NetworkSecretResponded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"SequencerEnclaveGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"SequencerEnclaveRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"grantSequencerEnclave\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_sequencerHost\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_genesisAttestation\",\"type\":\"string\"}],\"name\":\"initializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"isAttested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"isSequencer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"requestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"respondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"revokeSequencerEnclave\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461002257610011610026565b604051611fd76101948239611fd790f35b5f80fd5b61002e610030565b565b61002e61002e61002e6100b7565b61004b9060401c60ff1690565b90565b61004b905461003e565b61004b905b6001600160401b031690565b61004b9054610058565b61004b9061005d906001600160401b031682565b9061009761004b6100b392610073565b82546001600160401b0319166001600160401b03919091161790565b9055565b5f6100c061014d565b016100ca8161004e565b61013c576100d781610069565b6001600160401b03919082908116036100ee575050565b8161011d7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29361013793610087565b604051918291826001600160401b03909116815260200190565b0390a1565b63f92ee8a960e01b5f908152600490fd5b61004b61018b565b61004b61004b61004b9290565b61004b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610155565b61004b61016256fe60806040526004361015610011575f80fd5b5f3560e01c80633c23afba146100f0578063485cc955146100eb578063534ddc7a146100e65780635ac49c4e146100e15780635ad124ef146100dc5780636d46e987146100d7578063715018a6146100d257806379ba5097146100cd57806384b0196e146100c85780638da5cb5b146100c3578063a3411155146100be578063e30c3978146100b9578063f2fde38b146100b45763f3cbc5f80361011557610604565b610596565b61057b565b610563565b61052e565b6104fa565b6103bb565b6103a6565b610381565b610368565b6102f2565b6101b2565b610194565b610140565b6001600160a01b031690565b90565b6001600160a01b0381160361011557565b5f80fd5b9050359061012682610104565b565b906020828203126101155761010191610119565b9052565b346101155761016d61015b610156366004610128565b610673565b60405191829182901515815260200190565b0390f35b91906040838203126101155761010190602061018d8286610119565b9401610119565b34610115576101ad6101a7366004610171565b90610aa9565b604051005b34610115576101ad6101c5366004610128565b610b9c565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761020057604052565b6101ca565b9061012661021260405190565b92836101de565b67ffffffffffffffff811161020057602090601f01601f19160190565b0190565b90825f939282370152565b9092919261025a61025582610219565b610205565b93818552818301116101155761012691602085019061023a565b9080601f830112156101155781602061010193359101610245565b90608082820312610115576102a48183610119565b926102b28260208501610119565b92604081013567ffffffffffffffff811161011557836102d3918301610274565b92606082013567ffffffffffffffff8111610115576101019201610274565b34610115576101ad61030536600461028f565b92919091610d6e565b909182601f830112156101155781359167ffffffffffffffff831161011557602001926001830284011161011557565b9060208282031261011557813567ffffffffffffffff811161011557610364920161030e565b9091565b34610115576101ad61037b36600461033e565b90610f4e565b346101155761016d61015b610397366004610128565b610fae565b5f91031261011557565b34610115576103b636600461039c565b611037565b34610115576103cb36600461039c565b6101ad61103c565b90825f9392825e0152565b6103ff610408602093610236936103f3815190565b80835293849260200190565b958691016103d3565b601f01601f191690565b9061043261042b610421845190565b8084529260200190565b9260200190565b905f5b8181106104425750505090565b90919261045f6104586001928651815260200190565b9460200190565b929101610435565b939591946104d66104ce6104ed956104c06104e6956101019c9a6104b360e08c01925f8d01907fff00000000000000000000000000000000000000000000000000000000000000169052565b8a820360208c01526103de565b9088820360408a01526103de565b976060870152565b6001600160a01b03166080850152565b60a0830152565b60c0818403910152610412565b346101155761050a36600461039c565b61016d610515611148565b9361052597959793919360405190565b97889788610467565b346101155761053e36600461039c565b61016d610549611222565b604051918291826001600160a01b03909116815260200190565b34610115576101ad610576366004610128565b6112f6565b346101155761058b36600461039c565b61016d6105496112ff565b34610115576101ad6105a9366004610128565b611398565b606081830312610115576105c28282610119565b92602082013567ffffffffffffffff811161011557836105e391840161030e565b929093604082013567ffffffffffffffff811161011557610364920161030e565b34610115576101ad6106173660046105ae565b939290926114ab565b6100f5610101610101926001600160a01b031690565b61010190610620565b61010190610636565b906106529061063f565b5f5260205260405f2090565b610101905b60ff1690565b610101905461065e565b610689610101916106815f90565b506001610648565b610669565b6101019060401c610663565b610101905461068e565b610101905b67ffffffffffffffff1690565b61010190546106a4565b6106a96101016101019290565b6101016101016101019290565b9067ffffffffffffffff905b9181191691161790565b6106a96101016101019267ffffffffffffffff1690565b9061071761010161071e926106f0565b82546106da565b9055565b9068ff00000000000000009060401b6106e6565b9061074661010161071e92151590565b8254610722565b61013c906106c0565b602081019291610126919061074d565b61076e611549565b90819061078a6107846107808461069a565b1590565b926106b6565b936107945f6106c0565b67ffffffffffffffff861614806108b1575b6001956107c36107b5886106c0565b9167ffffffffffffffff1690565b149081610889575b155b9081610880575b50610854576107fd91836107f45f6107eb896106c0565b97019687610707565b61084557610a34565b610805575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916108345f61084093610736565b60405191829182610756565b0390a1565b61084f8686610736565b610a34565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f90815260045b035ffd5b1590505f6107d4565b90506107cd6108973061063f565b3b6108a86108a45f6106cd565b9190565b149190506107cb565b50826107a6565b6100f56101016101019290565b610101906108b8565b156108d557565b60405162461bcd60e51b815260206004820152601360248201527f4f776e65722063616e6e6f7420626520307830000000000000000000000000006044820152606490fd5b1561092157565b60405162461bcd60e51b815260206004820152601c60248201527f53657175656e63657220686f73742063616e6e6f7420626520307830000000006044820152606490fd5b9061097361025583610219565b918252565b6109826016610966565b7f4e6574776f726b456e636c617665526567697374727900000000000000000000602082015290565b610101610978565b6109bd6001610966565b7f3100000000000000000000000000000000000000000000000000000000000000602082015290565b6101016109b3565b9060ff906106e6565b90610a0761010161071e92151590565b82546109ee565b906001600160a01b03906106e6565b90610a2d61010161071e9261063f565b8254610a0e565b90610a7f61012692610a7a610a69610a4b5f6108c5565b6100f56001600160a01b0382166001600160a01b03861614156108ce565b6001600160a01b038516141561091a565b61156e565b610a98610a8a6109ab565b610a926109e6565b9061158e565b610aa25f806109f7565b6003610a1d565b9061012691610766565b61012690610abf611598565b610b4d565b15610acb57565b60405162461bcd60e51b815260206004820152601960248201527f656e636c6176654944206e6f7420612073657175656e636572000000000000006044820152606490fd5b919060086106e6910291610b2460ff841b90565b921b90565b9190610b3a61010161071e93151590565b908354610b10565b610126915f91610b29565b6108407f0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b4791610b88610b83610689836002610648565b610ac4565b6105495f610b97836002610648565b610b42565b61012690610ab3565b15610bac57565b60405162461bcd60e51b815260206004820152602660248201527f726573706f6e64696e67206174746573746572206973206e6f7420612073657160448201527f75656e63657200000000000000000000000000000000000000000000000000006064820152608490fd5b15610c1e57565b60405162461bcd60e51b815260206004820152601a60248201527f72657175657374657220616c72656164792061747465737465640000000000006044820152606490fd5b15610c6a57565b60405162461bcd60e51b815260206004820152601960248201527f696e76616c6964207265717565737465722061646472657373000000000000006044820152606490fd5b15610cb657565b60405162461bcd60e51b815260206004820152601e60248201527f696e76616c69642073656372657420726573706f6e7365206c656e67687400006044820152606490fd5b9081526001600160a01b0390911660208201526060810192916101269160400152565b0152565b15610d2957565b60405162461bcd60e51b815260206004820152601160248201527f696e76616c6964207369676e61747572650000000000000000000000000000006044820152606490fd5b610e9882610e89610e68610ea396610e63610e9d9799610d9a610d956106898a6002610648565b610ba5565b610db4610daf6107806106896001998a610648565b610c17565b610dd4610dc36100f55f6108c5565b6001600160a01b038d161415610c63565b610df2610ddf825190565b610dec6108a460916106cd565b14610caf565b610e528b7f6cf6b11fb1cd2d1b02d7e8188664fa0cef474b883701fe1020f57b4677837f9e92610e2a610e23825190565b9160200190565b2092610e46610e3860405190565b948593602085019384610cfb565b908103825203826101de565b610e5d610e23825190565b206115b3565b6115f7565b610e836001600160a01b0386165b916001600160a01b031690565b14610d22565b610e938682610648565b6109f7565b61063f565b9161063f565b907fb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be1610ece60405190565b5f90a3565b15610eda57565b60405162461bcd60e51b815260206004820152601060248201527f616c7265616479206174746573746564000000000000000000000000000000006044820152606490fd5b919061040881610f36816102369560209181520190565b809561023a565b602080825261010193910191610f1f565b907f0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d430190610f8a610f85610780610689336001610648565b610ed3565b610f933361063f565b92610fa9610fa060405190565b92839283610f3d565b0390a2565b61068961010191610fbc5f90565b506002610648565b610fcc611598565b60405162461bcd60e51b815260206004820152603460248201527f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f60448201527f742072656e6f756e6365206f776e6572736869700000000000000000000000006064820152608490fd5b610fc4565b336110456112ff565b6110576001600160a01b038316610e76565b03611065576101269061164b565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b6101019081565b610101905461109e565b156110b657565b60405162461bcd60e51b815260206004820152601560248201527f4549503731323a20556e696e697469616c697a656400000000000000000000006044820152606490fd5b67ffffffffffffffff81116102005760208091020190565b90610973610255836110fb565b369037565b9061012661113b61113584611113565b936110fb565b601f190160208401611120565b61118f7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100611175816110a5565b6111816108a45f6106cd565b1490816111ef575b506110af565b611197611770565b906111a061179d565b906111aa3061063f565b6111b35f6106cd565b6111c46111bf5f6106cd565b611125565b7f0f000000000000000000000000000000000000000000000000000000000000009594934693929190565b6111fc91506001016110a5565b6112086108a45f6106cd565b145f611189565b610101906100f5565b610101905461120f565b6101015f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b01611218565b6101269061125a611598565b6112ab565b1561126657565b60405162461bcd60e51b815260206004820152601660248201527f656e636c6176654944206e6f74206174746573746564000000000000000000006044820152606490fd5b6108407ffe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e760936916112e66112e1610689836001610648565b61125f565b6105496001610e93836002610648565b6101269061124e565b6101015f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00611248565b61012690611335611598565b61135f817f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00610a1d565b61136d610e9d610e98611222565b907f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700610ece60405190565b61012690611329565b156113a857565b60405162461bcd60e51b815260206004820152600e60248201527f6e6f7420617574686f72697a65640000000000000000000000000000000000006044820152606490fd5b156113f457565b60405162461bcd60e51b815260206004820152602260248201527f6e6574776f726b2073656372657420616c726561647920696e697469616c697a60448201527f65640000000000000000000000000000000000000000000000000000000000006064820152608490fd5b1561146657565b60405162461bcd60e51b815260206004820152601760248201527f696e76616c696420656e636c61766520616464726573730000000000000000006044820152606490fd5b9150507fd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23925061084091506114ed6114e66100f56003611218565b33146113a1565b6115016114fc6107805f610669565b6113ed565b6115216115106100f55f6108c5565b6001600160a01b038316141561145f565b6105496001611530815f6109f7565b61153e81610e938582610648565b610e93836002610648565b6101016117f1565b6101269061155d6117f9565b61156690611849565b61012661185a565b61012690611551565b90610126916115846117f9565b9061012691611a69565b9061012691611577565b6115a0611222565b33906115ab82610e76565b036110655750565b610101906115bf611a73565b604291604051917f19010000000000000000000000000000000000000000000000000000000000008352600283015260228201522090565b6101019161160491611a7b565b90929192611b16565b919060086106e6910291610b246001600160a01b03841b90565b919061163861010161071e9361063f565b90835461160d565b610126915f91611627565b610126906116795f7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00611640565b611c12565b634e487b7160e01b5f52602260045260245ffd5b90600160028304921680156116b2575b60208310146116ad57565b61167e565b91607f16916116a2565b80545f9392916116d86116ce83611692565b8085529360200190565b916001811690811561172757506001146116f157505050565b61170291929394505f5260205f2090565b915f925b8184106117135750500190565b805484840152602090930192600101611706565b92949550505060ff1916825215156020020190565b90610101916116bc565b906101266117609261175760405190565b9384809261173c565b03836101de565b61010190611746565b61010160027fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d1005b01611767565b61010160037fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100611797565b6101017ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006106cd565b6101016117c8565b611804610780611c77565b61180a57565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f908152600490fd5b610126906118406117f9565b61012690611cfa565b61012690611834565b6101266117f9565b610126611852565b906101269161186f6117f9565b611a09565b915f1960089290920291821b911b6106e6565b919061189861010161071e936106cd565b908354611874565b610126915f91611887565b8181106118b6575050565b806118c35f6001936118a0565b016118ab565b9190601f81116118d857505050565b6118e8610126935f5260205f2090565b906020601f84018190048301931061190a575b6020601f9091010401906118ab565b90915081906118fb565b9061191d815190565b9067ffffffffffffffff8211610200576119418261193b8554611692565b856118c9565b602090601f831160011461197a5761071e92915f918361196f575b50505f19600883021c1916906002021790565b015190505f8061195c565b601f1983169161198d855f5260205f2090565b925f5b8181106119c9575091600293918560019694106119b1575b50505002019055565b01515f196008601f8516021c191690555f80806119a8565b91936020600181928787015181550195019201611990565b9061012691611914565b905f19906106e6565b90611a0261010161071e9290565b82546119eb565b611a4d61012692611a44611a3a7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d10090565b93600285016119e1565b600383016119e1565b6001611a585f6106cd565b91611a6383826119f4565b016119f4565b9061012691611862565b610101611d4a565b905f91611a86825190565b611a936108a460416106cd565b03611abc57611ab59250602082015190606060408401519301515f1a90611de7565b9192909190565b509050611ad8611ad3611ace5f6108c5565b925190565b6106cd565b909160029190565b634e487b7160e01b5f52602160045260245ffd5b60041115611afe57565b611ae0565b9061012682611af4565b610101906106cd565b611b1f5f611b03565b611b2882611b03565b03611b31575050565b611b3b6001611b03565b611b4482611b03565b03611b73577ff645eedf000000000000000000000000000000000000000000000000000000005f908152600490fd5b611b7d6002611b03565b611b8682611b03565b03611bc7575f61087c611b9884611b0d565b7ffce698f700000000000000000000000000000000000000000000000000000000835260048301526024820190565b611bda611bd46003611b03565b91611b03565b14611be25750565b7fd78bce0c000000000000000000000000000000000000000000000000000000005f908152600491909152602490fd5b611c4c610e9d7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300610e9884611c4683611218565b92610a1d565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0610ece60405190565b6101015f611c83611549565b0161069a565b61012690611c956117f9565b611c9e5f6108c5565b6001600160a01b0381166001600160a01b03831614611cc157506101269061164b565b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b61012690611c89565b9095949261012694611d34611d3b92611d2d608096611d2660a088019c5f890152565b6020870152565b6040850152565b6060830152565b01906001600160a01b03169052565b7f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f611d9e611d76611ead565b611d7e611f41565b90610e46611d8b3061063f565b6040519586946020860194469286611d03565b611da9610e23825190565b2090565b610d1e61012694611dd5606094989795611dcb608086019a5f870152565b60ff166020850152565b6040830152565b6040513d5f823e3d90fd5b9091611df284611b0d565b611e1e6108a47f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a06106cd565b11611e9a5790611e3f6020945f9493611e3660405190565b94859485611dad565b838052039060015afa15611e95575f51611e585f6108c5565b6001600160a01b0381166001600160a01b03831614611e825750611e7b5f6106cd565b90915f9190565b9050611e8d5f6106cd565b909160019190565b611ddc565b505050611ea65f6108c5565b9160039190565b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100611ed9610101611770565b90611ee2825190565b611eee6108a45f6106cd565b1115611f005750611da9610e23825190565b611f0a91506110a5565b611f135f6106cd565b8114611f1c5790565b507fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47090565b7fa16a46d94261c7517cc8ff89f61c0ce93598e3c849801011dee649a6a557d100611f6d61010161179d565b90611f76825190565b611f826108a45f6106cd565b1115611f945750611da9610e23825190565b611f0a91506001016110a556fea26469706673582212206d77df66ec5d78a530bc2f3899892ce99dc9f451abefdc69ab0159f95c710ebb64736f6c634300081c0033",
}

// NetworkEnclaveRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkEnclaveRegistryMetaData.ABI instead.
var NetworkEnclaveRegistryABI = NetworkEnclaveRegistryMetaData.ABI

// NetworkEnclaveRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NetworkEnclaveRegistryMetaData.Bin instead.
var NetworkEnclaveRegistryBin = NetworkEnclaveRegistryMetaData.Bin

// DeployNetworkEnclaveRegistry deploys a new Ethereum contract, binding an instance of NetworkEnclaveRegistry to it.
func DeployNetworkEnclaveRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NetworkEnclaveRegistry, error) {
	parsed, err := NetworkEnclaveRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NetworkEnclaveRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NetworkEnclaveRegistry{NetworkEnclaveRegistryCaller: NetworkEnclaveRegistryCaller{contract: contract}, NetworkEnclaveRegistryTransactor: NetworkEnclaveRegistryTransactor{contract: contract}, NetworkEnclaveRegistryFilterer: NetworkEnclaveRegistryFilterer{contract: contract}}, nil
}

// NetworkEnclaveRegistry is an auto generated Go binding around an Ethereum contract.
type NetworkEnclaveRegistry struct {
	NetworkEnclaveRegistryCaller     // Read-only binding to the contract
	NetworkEnclaveRegistryTransactor // Write-only binding to the contract
	NetworkEnclaveRegistryFilterer   // Log filterer for contract events
}

// NetworkEnclaveRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkEnclaveRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkEnclaveRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkEnclaveRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkEnclaveRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkEnclaveRegistrySession struct {
	Contract     *NetworkEnclaveRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NetworkEnclaveRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkEnclaveRegistryCallerSession struct {
	Contract *NetworkEnclaveRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// NetworkEnclaveRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkEnclaveRegistryTransactorSession struct {
	Contract     *NetworkEnclaveRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// NetworkEnclaveRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkEnclaveRegistryRaw struct {
	Contract *NetworkEnclaveRegistry // Generic contract binding to access the raw methods on
}

// NetworkEnclaveRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryCallerRaw struct {
	Contract *NetworkEnclaveRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkEnclaveRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryTransactorRaw struct {
	Contract *NetworkEnclaveRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkEnclaveRegistry creates a new instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistry(address common.Address, backend bind.ContractBackend) (*NetworkEnclaveRegistry, error) {
	contract, err := bindNetworkEnclaveRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistry{NetworkEnclaveRegistryCaller: NetworkEnclaveRegistryCaller{contract: contract}, NetworkEnclaveRegistryTransactor: NetworkEnclaveRegistryTransactor{contract: contract}, NetworkEnclaveRegistryFilterer: NetworkEnclaveRegistryFilterer{contract: contract}}, nil
}

// NewNetworkEnclaveRegistryCaller creates a new read-only instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistryCaller(address common.Address, caller bind.ContractCaller) (*NetworkEnclaveRegistryCaller, error) {
	contract, err := bindNetworkEnclaveRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryCaller{contract: contract}, nil
}

// NewNetworkEnclaveRegistryTransactor creates a new write-only instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkEnclaveRegistryTransactor, error) {
	contract, err := bindNetworkEnclaveRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryTransactor{contract: contract}, nil
}

// NewNetworkEnclaveRegistryFilterer creates a new log filterer instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkEnclaveRegistryFilterer, error) {
	contract, err := bindNetworkEnclaveRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryFilterer{contract: contract}, nil
}

// bindNetworkEnclaveRegistry binds a generic wrapper to an already deployed contract.
func bindNetworkEnclaveRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NetworkEnclaveRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkEnclaveRegistry.Contract.NetworkEnclaveRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.NetworkEnclaveRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.NetworkEnclaveRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkEnclaveRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "eip712Domain")

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
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _NetworkEnclaveRegistry.Contract.Eip712Domain(&_NetworkEnclaveRegistry.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _NetworkEnclaveRegistry.Contract.Eip712Domain(&_NetworkEnclaveRegistry.CallOpts)
}

// IsAttested is a free data retrieval call binding the contract method 0x3c23afba.
//
// Solidity: function isAttested(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) IsAttested(opts *bind.CallOpts, enclaveID common.Address) (bool, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "isAttested", enclaveID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAttested is a free data retrieval call binding the contract method 0x3c23afba.
//
// Solidity: function isAttested(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) IsAttested(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsAttested(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// IsAttested is a free data retrieval call binding the contract method 0x3c23afba.
//
// Solidity: function isAttested(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) IsAttested(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsAttested(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// IsSequencer is a free data retrieval call binding the contract method 0x6d46e987.
//
// Solidity: function isSequencer(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) IsSequencer(opts *bind.CallOpts, enclaveID common.Address) (bool, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "isSequencer", enclaveID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSequencer is a free data retrieval call binding the contract method 0x6d46e987.
//
// Solidity: function isSequencer(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) IsSequencer(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsSequencer(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// IsSequencer is a free data retrieval call binding the contract method 0x6d46e987.
//
// Solidity: function isSequencer(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) IsSequencer(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsSequencer(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) Owner() (common.Address, error) {
	return _NetworkEnclaveRegistry.Contract.Owner(&_NetworkEnclaveRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) Owner() (common.Address, error) {
	return _NetworkEnclaveRegistry.Contract.Owner(&_NetworkEnclaveRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) PendingOwner() (common.Address, error) {
	return _NetworkEnclaveRegistry.Contract.PendingOwner(&_NetworkEnclaveRegistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) PendingOwner() (common.Address, error) {
	return _NetworkEnclaveRegistry.Contract.PendingOwner(&_NetworkEnclaveRegistry.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RenounceOwnership() error {
	return _NetworkEnclaveRegistry.Contract.RenounceOwnership(&_NetworkEnclaveRegistry.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) RenounceOwnership() error {
	return _NetworkEnclaveRegistry.Contract.RenounceOwnership(&_NetworkEnclaveRegistry.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.AcceptOwnership(&_NetworkEnclaveRegistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.AcceptOwnership(&_NetworkEnclaveRegistry.TransactOpts)
}

// GrantSequencerEnclave is a paid mutator transaction binding the contract method 0xa3411155.
//
// Solidity: function grantSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) GrantSequencerEnclave(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "grantSequencerEnclave", _addr)
}

// GrantSequencerEnclave is a paid mutator transaction binding the contract method 0xa3411155.
//
// Solidity: function grantSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) GrantSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.GrantSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// GrantSequencerEnclave is a paid mutator transaction binding the contract method 0xa3411155.
//
// Solidity: function grantSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) GrantSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.GrantSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _owner, address _sequencerHost) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _sequencerHost common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "initialize", _owner, _sequencerHost)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _owner, address _sequencerHost) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) Initialize(_owner common.Address, _sequencerHost common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.Initialize(&_NetworkEnclaveRegistry.TransactOpts, _owner, _sequencerHost)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _owner, address _sequencerHost) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) Initialize(_owner common.Address, _sequencerHost common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.Initialize(&_NetworkEnclaveRegistry.TransactOpts, _owner, _sequencerHost)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xf3cbc5f8.
//
// Solidity: function initializeNetworkSecret(address enclaveID, bytes _initSecret, string _genesisAttestation) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, enclaveID common.Address, _initSecret []byte, _genesisAttestation string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "initializeNetworkSecret", enclaveID, _initSecret, _genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xf3cbc5f8.
//
// Solidity: function initializeNetworkSecret(address enclaveID, bytes _initSecret, string _genesisAttestation) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) InitializeNetworkSecret(enclaveID common.Address, _initSecret []byte, _genesisAttestation string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.InitializeNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, enclaveID, _initSecret, _genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xf3cbc5f8.
//
// Solidity: function initializeNetworkSecret(address enclaveID, bytes _initSecret, string _genesisAttestation) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) InitializeNetworkSecret(enclaveID common.Address, _initSecret []byte, _genesisAttestation string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.InitializeNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, enclaveID, _initSecret, _genesisAttestation)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0x5ad124ef.
//
// Solidity: function requestNetworkSecret(string requestReport) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "requestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0x5ad124ef.
//
// Solidity: function requestNetworkSecret(string requestReport) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RequestNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0x5ad124ef.
//
// Solidity: function requestNetworkSecret(string requestReport) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RequestNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x5ac49c4e.
//
// Solidity: function respondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "respondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x5ac49c4e.
//
// Solidity: function respondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RespondNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x5ac49c4e.
//
// Solidity: function respondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RespondNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}

// RevokeSequencerEnclave is a paid mutator transaction binding the contract method 0x534ddc7a.
//
// Solidity: function revokeSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) RevokeSequencerEnclave(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "revokeSequencerEnclave", _addr)
}

// RevokeSequencerEnclave is a paid mutator transaction binding the contract method 0x534ddc7a.
//
// Solidity: function revokeSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RevokeSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RevokeSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// RevokeSequencerEnclave is a paid mutator transaction binding the contract method 0x534ddc7a.
//
// Solidity: function revokeSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) RevokeSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RevokeSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.TransferOwnership(&_NetworkEnclaveRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.TransferOwnership(&_NetworkEnclaveRegistry.TransactOpts, newOwner)
}

// NetworkEnclaveRegistryEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryEIP712DomainChangedIterator struct {
	Event *NetworkEnclaveRegistryEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryEIP712DomainChanged)
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
		it.Event = new(NetworkEnclaveRegistryEIP712DomainChanged)
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
func (it *NetworkEnclaveRegistryEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryEIP712DomainChanged represents a EIP712DomainChanged event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*NetworkEnclaveRegistryEIP712DomainChangedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryEIP712DomainChangedIterator{contract: _NetworkEnclaveRegistry.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryEIP712DomainChanged)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseEIP712DomainChanged(log types.Log) (*NetworkEnclaveRegistryEIP712DomainChanged, error) {
	event := new(NetworkEnclaveRegistryEIP712DomainChanged)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryInitializedIterator struct {
	Event *NetworkEnclaveRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryInitialized)
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
		it.Event = new(NetworkEnclaveRegistryInitialized)
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
func (it *NetworkEnclaveRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryInitialized represents a Initialized event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*NetworkEnclaveRegistryInitializedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryInitializedIterator{contract: _NetworkEnclaveRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryInitialized)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseInitialized(log types.Log) (*NetworkEnclaveRegistryInitialized, error) {
	event := new(NetworkEnclaveRegistryInitialized)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryNetworkSecretInitializedIterator is returned from FilterNetworkSecretInitialized and is used to iterate over the raw logs and unpacked data for NetworkSecretInitialized events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretInitializedIterator struct {
	Event *NetworkEnclaveRegistryNetworkSecretInitialized // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryNetworkSecretInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryNetworkSecretInitialized)
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
		it.Event = new(NetworkEnclaveRegistryNetworkSecretInitialized)
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
func (it *NetworkEnclaveRegistryNetworkSecretInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryNetworkSecretInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryNetworkSecretInitialized represents a NetworkSecretInitialized event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretInitialized struct {
	EnclaveID common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNetworkSecretInitialized is a free log retrieval operation binding the contract event 0xd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23.
//
// Solidity: event NetworkSecretInitialized(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterNetworkSecretInitialized(opts *bind.FilterOpts) (*NetworkEnclaveRegistryNetworkSecretInitializedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "NetworkSecretInitialized")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryNetworkSecretInitializedIterator{contract: _NetworkEnclaveRegistry.contract, event: "NetworkSecretInitialized", logs: logs, sub: sub}, nil
}

// WatchNetworkSecretInitialized is a free log subscription operation binding the contract event 0xd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23.
//
// Solidity: event NetworkSecretInitialized(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchNetworkSecretInitialized(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryNetworkSecretInitialized) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "NetworkSecretInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryNetworkSecretInitialized)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretInitialized", log); err != nil {
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

// ParseNetworkSecretInitialized is a log parse operation binding the contract event 0xd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23.
//
// Solidity: event NetworkSecretInitialized(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseNetworkSecretInitialized(log types.Log) (*NetworkEnclaveRegistryNetworkSecretInitialized, error) {
	event := new(NetworkEnclaveRegistryNetworkSecretInitialized)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryNetworkSecretRequestedIterator is returned from FilterNetworkSecretRequested and is used to iterate over the raw logs and unpacked data for NetworkSecretRequested events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretRequestedIterator struct {
	Event *NetworkEnclaveRegistryNetworkSecretRequested // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryNetworkSecretRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryNetworkSecretRequested)
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
		it.Event = new(NetworkEnclaveRegistryNetworkSecretRequested)
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
func (it *NetworkEnclaveRegistryNetworkSecretRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryNetworkSecretRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryNetworkSecretRequested represents a NetworkSecretRequested event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretRequested struct {
	Requester     common.Address
	RequestReport string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNetworkSecretRequested is a free log retrieval operation binding the contract event 0x0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d4301.
//
// Solidity: event NetworkSecretRequested(address indexed requester, string requestReport)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterNetworkSecretRequested(opts *bind.FilterOpts, requester []common.Address) (*NetworkEnclaveRegistryNetworkSecretRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "NetworkSecretRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryNetworkSecretRequestedIterator{contract: _NetworkEnclaveRegistry.contract, event: "NetworkSecretRequested", logs: logs, sub: sub}, nil
}

// WatchNetworkSecretRequested is a free log subscription operation binding the contract event 0x0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d4301.
//
// Solidity: event NetworkSecretRequested(address indexed requester, string requestReport)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchNetworkSecretRequested(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryNetworkSecretRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "NetworkSecretRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryNetworkSecretRequested)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretRequested", log); err != nil {
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

// ParseNetworkSecretRequested is a log parse operation binding the contract event 0x0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d4301.
//
// Solidity: event NetworkSecretRequested(address indexed requester, string requestReport)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseNetworkSecretRequested(log types.Log) (*NetworkEnclaveRegistryNetworkSecretRequested, error) {
	event := new(NetworkEnclaveRegistryNetworkSecretRequested)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryNetworkSecretRespondedIterator is returned from FilterNetworkSecretResponded and is used to iterate over the raw logs and unpacked data for NetworkSecretResponded events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretRespondedIterator struct {
	Event *NetworkEnclaveRegistryNetworkSecretResponded // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryNetworkSecretRespondedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryNetworkSecretResponded)
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
		it.Event = new(NetworkEnclaveRegistryNetworkSecretResponded)
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
func (it *NetworkEnclaveRegistryNetworkSecretRespondedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryNetworkSecretRespondedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryNetworkSecretResponded represents a NetworkSecretResponded event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretResponded struct {
	Attester  common.Address
	Requester common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNetworkSecretResponded is a free log retrieval operation binding the contract event 0xb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be1.
//
// Solidity: event NetworkSecretResponded(address indexed attester, address indexed requester)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterNetworkSecretResponded(opts *bind.FilterOpts, attester []common.Address, requester []common.Address) (*NetworkEnclaveRegistryNetworkSecretRespondedIterator, error) {

	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "NetworkSecretResponded", attesterRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryNetworkSecretRespondedIterator{contract: _NetworkEnclaveRegistry.contract, event: "NetworkSecretResponded", logs: logs, sub: sub}, nil
}

// WatchNetworkSecretResponded is a free log subscription operation binding the contract event 0xb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be1.
//
// Solidity: event NetworkSecretResponded(address indexed attester, address indexed requester)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchNetworkSecretResponded(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryNetworkSecretResponded, attester []common.Address, requester []common.Address) (event.Subscription, error) {

	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "NetworkSecretResponded", attesterRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryNetworkSecretResponded)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretResponded", log); err != nil {
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

// ParseNetworkSecretResponded is a log parse operation binding the contract event 0xb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be1.
//
// Solidity: event NetworkSecretResponded(address indexed attester, address indexed requester)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseNetworkSecretResponded(log types.Log) (*NetworkEnclaveRegistryNetworkSecretResponded, error) {
	event := new(NetworkEnclaveRegistryNetworkSecretResponded)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretResponded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryOwnershipTransferStartedIterator struct {
	Event *NetworkEnclaveRegistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryOwnershipTransferStarted)
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
		it.Event = new(NetworkEnclaveRegistryOwnershipTransferStarted)
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
func (it *NetworkEnclaveRegistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NetworkEnclaveRegistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryOwnershipTransferStartedIterator{contract: _NetworkEnclaveRegistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryOwnershipTransferStarted)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*NetworkEnclaveRegistryOwnershipTransferStarted, error) {
	event := new(NetworkEnclaveRegistryOwnershipTransferStarted)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryOwnershipTransferredIterator struct {
	Event *NetworkEnclaveRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryOwnershipTransferred)
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
		it.Event = new(NetworkEnclaveRegistryOwnershipTransferred)
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
func (it *NetworkEnclaveRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NetworkEnclaveRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryOwnershipTransferredIterator{contract: _NetworkEnclaveRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryOwnershipTransferred)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*NetworkEnclaveRegistryOwnershipTransferred, error) {
	event := new(NetworkEnclaveRegistryOwnershipTransferred)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistrySequencerEnclaveGrantedIterator is returned from FilterSequencerEnclaveGranted and is used to iterate over the raw logs and unpacked data for SequencerEnclaveGranted events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveGrantedIterator struct {
	Event *NetworkEnclaveRegistrySequencerEnclaveGranted // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistrySequencerEnclaveGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistrySequencerEnclaveGranted)
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
		it.Event = new(NetworkEnclaveRegistrySequencerEnclaveGranted)
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
func (it *NetworkEnclaveRegistrySequencerEnclaveGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistrySequencerEnclaveGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistrySequencerEnclaveGranted represents a SequencerEnclaveGranted event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveGranted struct {
	EnclaveID common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSequencerEnclaveGranted is a free log retrieval operation binding the contract event 0xfe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e760936.
//
// Solidity: event SequencerEnclaveGranted(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterSequencerEnclaveGranted(opts *bind.FilterOpts) (*NetworkEnclaveRegistrySequencerEnclaveGrantedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "SequencerEnclaveGranted")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistrySequencerEnclaveGrantedIterator{contract: _NetworkEnclaveRegistry.contract, event: "SequencerEnclaveGranted", logs: logs, sub: sub}, nil
}

// WatchSequencerEnclaveGranted is a free log subscription operation binding the contract event 0xfe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e760936.
//
// Solidity: event SequencerEnclaveGranted(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchSequencerEnclaveGranted(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistrySequencerEnclaveGranted) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "SequencerEnclaveGranted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistrySequencerEnclaveGranted)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveGranted", log); err != nil {
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

// ParseSequencerEnclaveGranted is a log parse operation binding the contract event 0xfe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e760936.
//
// Solidity: event SequencerEnclaveGranted(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseSequencerEnclaveGranted(log types.Log) (*NetworkEnclaveRegistrySequencerEnclaveGranted, error) {
	event := new(NetworkEnclaveRegistrySequencerEnclaveGranted)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistrySequencerEnclaveRevokedIterator is returned from FilterSequencerEnclaveRevoked and is used to iterate over the raw logs and unpacked data for SequencerEnclaveRevoked events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveRevokedIterator struct {
	Event *NetworkEnclaveRegistrySequencerEnclaveRevoked // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistrySequencerEnclaveRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
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
		it.Event = new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
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
func (it *NetworkEnclaveRegistrySequencerEnclaveRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistrySequencerEnclaveRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistrySequencerEnclaveRevoked represents a SequencerEnclaveRevoked event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveRevoked struct {
	EnclaveID common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSequencerEnclaveRevoked is a free log retrieval operation binding the contract event 0x0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b47.
//
// Solidity: event SequencerEnclaveRevoked(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterSequencerEnclaveRevoked(opts *bind.FilterOpts) (*NetworkEnclaveRegistrySequencerEnclaveRevokedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "SequencerEnclaveRevoked")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistrySequencerEnclaveRevokedIterator{contract: _NetworkEnclaveRegistry.contract, event: "SequencerEnclaveRevoked", logs: logs, sub: sub}, nil
}

// WatchSequencerEnclaveRevoked is a free log subscription operation binding the contract event 0x0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b47.
//
// Solidity: event SequencerEnclaveRevoked(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchSequencerEnclaveRevoked(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistrySequencerEnclaveRevoked) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "SequencerEnclaveRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveRevoked", log); err != nil {
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

// ParseSequencerEnclaveRevoked is a log parse operation binding the contract event 0x0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b47.
//
// Solidity: event SequencerEnclaveRevoked(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseSequencerEnclaveRevoked(log types.Log) (*NetworkEnclaveRegistrySequencerEnclaveRevoked, error) {
	event := new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
