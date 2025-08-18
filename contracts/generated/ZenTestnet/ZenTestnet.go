// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ZenTestnet

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

// StructsTransaction is an auto generated low-level Go binding around an user-defined struct.
type StructsTransaction struct {
	TxType     uint8
	Nonce      *big.Int
	GasPrice   *big.Int
	GasLimit   *big.Int
	To         common.Address
	Value      *big.Int
	Data       []byte
	From       common.Address
	Successful bool
	GasUsed    uint64
}

// ZenTestnetMetaData contains all meta data concerning the ZenTestnet contract.
var ZenTestnetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransactionProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transactionPostProcessor\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"successful\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"gasUsed\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlockEnd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461002257610011610026565b6040516115b061018682396115b090f35b5f80fd5b61002e6100a9565b565b61003d9060401c60ff1690565b90565b61003d9054610030565b61003d905b6001600160401b031690565b61003d905461004a565b61003d9061004f906001600160401b031682565b9061008961003d6100a592610065565b82546001600160401b0319166001600160401b03919091161790565b9055565b5f6100b261013f565b016100bc81610040565b61012e576100c98161005b565b6001600160401b03919082908116036100e0575050565b8161010f7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29361012993610079565b604051918291826001600160401b03909116815260200190565b0390a1565b63f92ee8a960e01b5f908152600490fd5b61003d61017d565b61003d61003d61003d9290565b61003d7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610147565b61003d61015456fe60806040526004361015610011575f80fd5b5f3560e01c806306fdde03146100f0578063095ea7b3146100eb57806318160ddd146100e657806323b872dd146100e1578063313ce567146100dc57806370a08231146100d7578063715018a6146100d25780638da5cb5b146100cd57806395d89b41146100c85780639f9976af146100c3578063a9059cbb146100be578063c4d66de8146100b9578063dd62ed3e146100b45763f2fde38b036100ff5761043c565b610420565b6103e5565b6103c9565b6103b0565b61033b565b610306565b6102e9565b6102ce565b61028b565b61026f565b61021a565b6101ec565b61015a565b5f9103126100ff57565b5f80fd5b90825f9392825e0152565b61012f61013860209361014293610123815190565b80835293849260200190565b95869101610103565b601f01601f191690565b0190565b60208082526101579291019061010e565b90565b346100ff5761016a3660046100f5565b610181610175610581565b60405191829182610146565b0390f35b6001600160a01b031690565b6001600160a01b0381165b036100ff57565b905035906101b082610191565b565b8061019c565b905035906101b0826101b2565b91906040838203126100ff576101579060206101e182866101a3565b94016101b8565b9052565b346100ff576101816102086102023660046101c5565b906105ae565b60405191829182901515815260200190565b346100ff5761022a3660046100f5565b6101816102356105cf565b6040515b9182918290815260200190565b90916060828403126100ff5761015761025f84846101a3565b9360406101e182602087016101a3565b346100ff57610181610208610285366004610246565b916105f8565b346100ff5761029b3660046100f5565b6101816102a6610621565b6040519182918260ff909116815260200190565b906020828203126100ff57610157916101a3565b346100ff576101816102356102e43660046102ba565b610669565b346100ff576102f93660046100f5565b6103016106dd565b604051005b346100ff576103163660046100f5565b6101816103216106f8565b604051918291826001600160a01b03909116815260200190565b346100ff5761034b3660046100f5565b610181610175610721565b909182601f830112156100ff5781359167ffffffffffffffff83116100ff5760200192602083028401116100ff57565b906020828203126100ff57813567ffffffffffffffff81116100ff576103ac9201610356565b9091565b346100ff576103016103c3366004610386565b906108eb565b346100ff576101816102086103df3660046101c5565b906108f5565b346100ff576103016103f83660046102ba565b610cbc565b91906040838203126100ff5761015790602061041982866101a3565b94016101a3565b346100ff576101816102356104363660046103fd565b90610cc5565b346100ff5761030161044f3660046102ba565b610d76565b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610488575b602083101461048357565b610454565b91607f1691610478565b80545f9392916104ae6104a483610468565b8085529360200190565b91600181169081156104fd57506001146104c757505050565b6104d891929394505f5260205f2090565b915f925b8184106104e95750500190565b8054848401526020909301926001016104dc565b92949550505060ff1916825215156020020190565b9061015791610492565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761055257604052565b61051c565b906101b06105719261056860405190565b93848092610512565b0383610530565b61015790610557565b61015760037f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace005b01610578565b6105b9919033610d7f565b600190565b6101579081565b61015790546105be565b6101577f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace026105c5565b6105b9929190610609833383610db7565b610e28565b61061b6101576101579290565b60ff1690565b610157601261060e565b610185610157610157926001600160a01b031690565b6101579061062b565b61015790610641565b9061065d9061064a565b5f5260205260405f2090565b6106a1610157916106775f90565b505f7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace005b01610653565b6105c5565b6106ae610ec7565b6101b06106cc565b6101856101576101579290565b610157906106b6565b6101b06106d85f6106c3565b610f1d565b6101b06106a6565b61015790610185565b61015790546106e5565b6101577f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993006106ee565b61015760047f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace006105a8565b1561075357565b60405162461bcd60e51b815260206004820152602c60248201527f43616c6c65723a2063616c6c6572206973206e6f74207468652064657369676e60448201527f61746564206164647265737300000000000000000000000000000000000000006064820152608490fd5b906101b0916107d96107d26101855f6106ee565b331461074c565b61083c565b6101576101576101579290565b634e487b7160e01b5f52603260045260245ffd5b90359061013e1936829003018212156100ff570190565b9082101561082d57602061015792028101906107ff565b6107eb565b3561015781610191565b909190829161084a5f6107de565b83146108a6576108595f6107de565b835b81101561089f576108988161089261088260e061087c61085b968b8a610816565b01610832565b61088c60016107de565b90610f95565b60010190565b9050610859565b5092505050565b60405162461bcd60e51b815260206004820152601a60248201527f4e6f207472616e73616374696f6e7320746f20636f6e766572740000000000006044820152606490fd5b906101b0916107be565b6105b9919033610e28565b6101579060401c61061b565b6101579054610900565b610157905b67ffffffffffffffff1690565b6101579054610916565b61091b6101576101579290565b9067ffffffffffffffff905b9181191691161790565b61091b6101576101579267ffffffffffffffff1690565b9061097c61015761098392610955565b825461093f565b9055565b9068ff00000000000000009060401b61094b565b906109ab61015761098392151590565b8254610987565b6101e890610932565b6020810192916101b091906109b2565b6109d3610fe2565b806109ed6109e76109e38361090c565b1590565b91610928565b926109f75f610932565b67ffffffffffffffff85161480610b14575b600194610a26610a1887610932565b9167ffffffffffffffff1690565b149081610aec575b155b9081610ae3575b50610ab757610a609082610a575f610a4e88610932565b9601958661096c565b610aa857610c70565b610a68575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d291610a975f610aa39361099b565b604051918291826109bb565b0390a1565b610ab2858561099b565b610c70565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f90815260045b035ffd5b1590505f610a37565b9050610a30610afa3061064a565b3b610b0b610b075f6107de565b9190565b14919050610a2e565b5081610a09565b15610b2257565b60405162461bcd60e51b8152602060048201526024808201527f496e76616c6964207472616e73616374696f6e20616e616c797a65722061646460448201527f72657373000000000000000000000000000000000000000000000000000000006064820152608490fd5b906101b0610b9960405190565b9283610530565b67ffffffffffffffff811161055257602090601f01601f19160190565b90610bcf610bca83610ba0565b610b8c565b918252565b610bde6003610bbd565b7f5a656e0000000000000000000000000000000000000000000000000000000000602082015290565b610157610bd4565b610c196003610bbd565b7f5a454e0000000000000000000000000000000000000000000000000000000000602082015290565b610157610c0f565b906001600160a01b039061094b565b90610c696101576109839261064a565b8254610c4a565b6101b090610c94610c836101855f6106c3565b6001600160a01b0383161415610b1b565b610cad610c9f610c07565b610ca7610c42565b90611001565b610cb633611020565b5f610c59565b6101b0906109cb565b61015791610d006106a192610cd75f90565b5060017f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0061069b565b610653565b6101b090610d11610ec7565b610d1a5f6106c3565b6001600160a01b0381166001600160a01b03831614610d3d57506101b090610f1d565b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b6101b090610d05565b916001916101b093611049565b6001600160a01b0390911681526060810193926101b0929091604091610db3906020830152565b0152565b91610dc28284610cc5565b5f198110610dd1575b50505050565b818110610df65791610de7610ded94925f940390565b91611049565b5f808080610dcb565b7ffb8f41b2000000000000000000000000000000000000000000000000000000005f9081529350610adf926004610d8c565b929190610e345f6106c3565b936001600160a01b0385166001600160a01b03821614610e90576001600160a01b0385166001600160a01b03831614610e72576101b09394506111aa565b63ec442f0560e01b5f9081526001600160a01b038616600452602490fd5b7f96c6fd1e000000000000000000000000000000000000000000000000000000005f9081526001600160a01b038616600452602490fd5b610ecf6106f8565b339081906001600160a01b031603610ee45750565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b610f62610f5c7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300610f5784610f51836106ee565b92610c59565b61064a565b9161064a565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0610f8d60405190565b80805b0390a3565b9190610fa05f6106c3565b926001600160a01b0384166001600160a01b03821614610fc4576101b092936111aa565b63ec442f0560e01b5f9081526001600160a01b038516600452602490fd5b61015761133c565b906101b091610ff7611344565b906101b091611549565b906101b091610fea565b6101b090611017611344565b6101b09061155f565b6101b09061100b565b905f199061094b565b90611042610157610983926107de565b8254611029565b9091926110737f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0090565b61107c5f6106c3565b6001600160a01b0381166001600160a01b0385161461114b576001600160a01b0381166001600160a01b038616146111125750846110c485610d008660016110c99601610653565b611032565b6110d257505050565b610f906111086111027f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259361064a565b9361064a565b9361023960405190565b7f94280d62000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b7fe602df05000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b634e487b7160e01b5f52601160045260245ffd5b919082018092116111a557565b611184565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00826111d55f6106c3565b6001600160a01b0381166001600160a01b03851603611281577fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef93836111089361123e611102946101856002610f9099016112388d611233836105c5565b611198565b90611032565b6001600160a01b0382160361126657506002610f57910161123889611262836105c5565b0390565b610f57916112749190610653565b61123889610142836105c5565b9050816112916106a18583610653565b8681106112e05793610f90938661123e61110895610185611102966110c4876112da8f7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9e0390565b92610653565b610adf875f92877fe450d38c00000000000000000000000000000000000000000000000000000000855260048501610d8c565b6101577ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006107de565b610157611313565b61134f6109e3611568565b61135557565b7fd7e6bcf8000000000000000000000000000000000000000000000000000000005f908152600490fd5b906101b09161138c611344565b611508565b915f1960089290920291821b911b61094b565b91906113b5610157610983936107de565b908354611391565b6101b0915f916113a4565b8181106113d3575050565b806113e05f6001936113bd565b016113c8565b9190601f81116113f557505050565b6114056101b0935f5260205f2090565b906020601f840181900483019310611427575b6020601f9091010401906113c8565b9091508190611418565b9061143a815190565b9067ffffffffffffffff82116105525761145e826114588554610468565b856113e6565b602090601f83116001146114975761098392915f918361148c575b50505f19600883021c1916906002021790565b015190505f80611479565b601f198316916114aa855f5260205f2090565b925f5b8181106114e6575091600293918560019694106114ce575b50505002019055565b01515f196008601f8516021c191690555f80806114c5565b919360206001819287870151815501950192016114ad565b906101b091611431565b9060046101b0926115436115397f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0090565b91600383016114fe565b016114fe565b906101b09161137f565b6101b090610d11611344565b6101b090611553565b6101575f611574610fe2565b0161090c56fea26469706673582212205d9b609c2cfcf31bbe3553c6cf4d2dede355cbaf6f9721e820e6caca7791895764736f6c634300081c0033",
}

// ZenTestnetABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenTestnetMetaData.ABI instead.
var ZenTestnetABI = ZenTestnetMetaData.ABI

// ZenTestnetBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZenTestnetMetaData.Bin instead.
var ZenTestnetBin = ZenTestnetMetaData.Bin

// DeployZenTestnet deploys a new Ethereum contract, binding an instance of ZenTestnet to it.
func DeployZenTestnet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ZenTestnet, error) {
	parsed, err := ZenTestnetMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZenTestnetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ZenTestnet{ZenTestnetCaller: ZenTestnetCaller{contract: contract}, ZenTestnetTransactor: ZenTestnetTransactor{contract: contract}, ZenTestnetFilterer: ZenTestnetFilterer{contract: contract}}, nil
}

// ZenTestnet is an auto generated Go binding around an Ethereum contract.
type ZenTestnet struct {
	ZenTestnetCaller     // Read-only binding to the contract
	ZenTestnetTransactor // Write-only binding to the contract
	ZenTestnetFilterer   // Log filterer for contract events
}

// ZenTestnetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZenTestnetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenTestnetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZenTestnetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenTestnetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZenTestnetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenTestnetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZenTestnetSession struct {
	Contract     *ZenTestnet       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenTestnetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZenTestnetCallerSession struct {
	Contract *ZenTestnetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ZenTestnetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZenTestnetTransactorSession struct {
	Contract     *ZenTestnetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ZenTestnetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZenTestnetRaw struct {
	Contract *ZenTestnet // Generic contract binding to access the raw methods on
}

// ZenTestnetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZenTestnetCallerRaw struct {
	Contract *ZenTestnetCaller // Generic read-only contract binding to access the raw methods on
}

// ZenTestnetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZenTestnetTransactorRaw struct {
	Contract *ZenTestnetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZenTestnet creates a new instance of ZenTestnet, bound to a specific deployed contract.
func NewZenTestnet(address common.Address, backend bind.ContractBackend) (*ZenTestnet, error) {
	contract, err := bindZenTestnet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZenTestnet{ZenTestnetCaller: ZenTestnetCaller{contract: contract}, ZenTestnetTransactor: ZenTestnetTransactor{contract: contract}, ZenTestnetFilterer: ZenTestnetFilterer{contract: contract}}, nil
}

// NewZenTestnetCaller creates a new read-only instance of ZenTestnet, bound to a specific deployed contract.
func NewZenTestnetCaller(address common.Address, caller bind.ContractCaller) (*ZenTestnetCaller, error) {
	contract, err := bindZenTestnet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZenTestnetCaller{contract: contract}, nil
}

// NewZenTestnetTransactor creates a new write-only instance of ZenTestnet, bound to a specific deployed contract.
func NewZenTestnetTransactor(address common.Address, transactor bind.ContractTransactor) (*ZenTestnetTransactor, error) {
	contract, err := bindZenTestnet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZenTestnetTransactor{contract: contract}, nil
}

// NewZenTestnetFilterer creates a new log filterer instance of ZenTestnet, bound to a specific deployed contract.
func NewZenTestnetFilterer(address common.Address, filterer bind.ContractFilterer) (*ZenTestnetFilterer, error) {
	contract, err := bindZenTestnet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZenTestnetFilterer{contract: contract}, nil
}

// bindZenTestnet binds a generic wrapper to an already deployed contract.
func bindZenTestnet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZenTestnetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenTestnet *ZenTestnetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenTestnet.Contract.ZenTestnetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenTestnet *ZenTestnetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenTestnet.Contract.ZenTestnetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenTestnet *ZenTestnetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenTestnet.Contract.ZenTestnetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenTestnet *ZenTestnetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenTestnet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenTestnet *ZenTestnetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenTestnet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenTestnet *ZenTestnetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenTestnet.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenTestnet *ZenTestnetCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenTestnet.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenTestnet *ZenTestnetSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenTestnet.Contract.Allowance(&_ZenTestnet.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenTestnet *ZenTestnetCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenTestnet.Contract.Allowance(&_ZenTestnet.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenTestnet *ZenTestnetCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenTestnet.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenTestnet *ZenTestnetSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenTestnet.Contract.BalanceOf(&_ZenTestnet.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenTestnet *ZenTestnetCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenTestnet.Contract.BalanceOf(&_ZenTestnet.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenTestnet *ZenTestnetCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ZenTestnet.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenTestnet *ZenTestnetSession) Decimals() (uint8, error) {
	return _ZenTestnet.Contract.Decimals(&_ZenTestnet.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenTestnet *ZenTestnetCallerSession) Decimals() (uint8, error) {
	return _ZenTestnet.Contract.Decimals(&_ZenTestnet.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenTestnet *ZenTestnetCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenTestnet.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenTestnet *ZenTestnetSession) Name() (string, error) {
	return _ZenTestnet.Contract.Name(&_ZenTestnet.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenTestnet *ZenTestnetCallerSession) Name() (string, error) {
	return _ZenTestnet.Contract.Name(&_ZenTestnet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenTestnet *ZenTestnetCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZenTestnet.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenTestnet *ZenTestnetSession) Owner() (common.Address, error) {
	return _ZenTestnet.Contract.Owner(&_ZenTestnet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenTestnet *ZenTestnetCallerSession) Owner() (common.Address, error) {
	return _ZenTestnet.Contract.Owner(&_ZenTestnet.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenTestnet *ZenTestnetCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenTestnet.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenTestnet *ZenTestnetSession) Symbol() (string, error) {
	return _ZenTestnet.Contract.Symbol(&_ZenTestnet.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenTestnet *ZenTestnetCallerSession) Symbol() (string, error) {
	return _ZenTestnet.Contract.Symbol(&_ZenTestnet.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenTestnet *ZenTestnetCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZenTestnet.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenTestnet *ZenTestnetSession) TotalSupply() (*big.Int, error) {
	return _ZenTestnet.Contract.TotalSupply(&_ZenTestnet.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenTestnet *ZenTestnetCallerSession) TotalSupply() (*big.Int, error) {
	return _ZenTestnet.Contract.TotalSupply(&_ZenTestnet.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.Contract.Approve(&_ZenTestnet.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.Contract.Approve(&_ZenTestnet.TransactOpts, spender, value)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address transactionPostProcessor) returns()
func (_ZenTestnet *ZenTestnetTransactor) Initialize(opts *bind.TransactOpts, transactionPostProcessor common.Address) (*types.Transaction, error) {
	return _ZenTestnet.contract.Transact(opts, "initialize", transactionPostProcessor)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address transactionPostProcessor) returns()
func (_ZenTestnet *ZenTestnetSession) Initialize(transactionPostProcessor common.Address) (*types.Transaction, error) {
	return _ZenTestnet.Contract.Initialize(&_ZenTestnet.TransactOpts, transactionPostProcessor)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address transactionPostProcessor) returns()
func (_ZenTestnet *ZenTestnetTransactorSession) Initialize(transactionPostProcessor common.Address) (*types.Transaction, error) {
	return _ZenTestnet.Contract.Initialize(&_ZenTestnet.TransactOpts, transactionPostProcessor)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x9f9976af.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_ZenTestnet *ZenTestnetTransactor) OnBlockEnd(opts *bind.TransactOpts, transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenTestnet.contract.Transact(opts, "onBlockEnd", transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x9f9976af.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_ZenTestnet *ZenTestnetSession) OnBlockEnd(transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenTestnet.Contract.OnBlockEnd(&_ZenTestnet.TransactOpts, transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x9f9976af.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,address,uint256,bytes,address,bool,uint64)[] transactions) returns()
func (_ZenTestnet *ZenTestnetTransactorSession) OnBlockEnd(transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenTestnet.Contract.OnBlockEnd(&_ZenTestnet.TransactOpts, transactions)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenTestnet *ZenTestnetTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenTestnet.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenTestnet *ZenTestnetSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenTestnet.Contract.RenounceOwnership(&_ZenTestnet.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenTestnet *ZenTestnetTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenTestnet.Contract.RenounceOwnership(&_ZenTestnet.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.Contract.Transfer(&_ZenTestnet.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.Contract.Transfer(&_ZenTestnet.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.Contract.TransferFrom(&_ZenTestnet.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenTestnet *ZenTestnetTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenTestnet.Contract.TransferFrom(&_ZenTestnet.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenTestnet *ZenTestnetTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZenTestnet.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenTestnet *ZenTestnetSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenTestnet.Contract.TransferOwnership(&_ZenTestnet.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenTestnet *ZenTestnetTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenTestnet.Contract.TransferOwnership(&_ZenTestnet.TransactOpts, newOwner)
}

// ZenTestnetApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ZenTestnet contract.
type ZenTestnetApprovalIterator struct {
	Event *ZenTestnetApproval // Event containing the contract specifics and raw log

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
func (it *ZenTestnetApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenTestnetApproval)
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
		it.Event = new(ZenTestnetApproval)
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
func (it *ZenTestnetApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenTestnetApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenTestnetApproval represents a Approval event raised by the ZenTestnet contract.
type ZenTestnetApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenTestnet *ZenTestnetFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZenTestnetApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenTestnet.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZenTestnetApprovalIterator{contract: _ZenTestnet.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenTestnet *ZenTestnetFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZenTestnetApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenTestnet.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenTestnetApproval)
				if err := _ZenTestnet.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenTestnet *ZenTestnetFilterer) ParseApproval(log types.Log) (*ZenTestnetApproval, error) {
	event := new(ZenTestnetApproval)
	if err := _ZenTestnet.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenTestnetInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ZenTestnet contract.
type ZenTestnetInitializedIterator struct {
	Event *ZenTestnetInitialized // Event containing the contract specifics and raw log

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
func (it *ZenTestnetInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenTestnetInitialized)
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
		it.Event = new(ZenTestnetInitialized)
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
func (it *ZenTestnetInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenTestnetInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenTestnetInitialized represents a Initialized event raised by the ZenTestnet contract.
type ZenTestnetInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ZenTestnet *ZenTestnetFilterer) FilterInitialized(opts *bind.FilterOpts) (*ZenTestnetInitializedIterator, error) {

	logs, sub, err := _ZenTestnet.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ZenTestnetInitializedIterator{contract: _ZenTestnet.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ZenTestnet *ZenTestnetFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ZenTestnetInitialized) (event.Subscription, error) {

	logs, sub, err := _ZenTestnet.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenTestnetInitialized)
				if err := _ZenTestnet.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ZenTestnet *ZenTestnetFilterer) ParseInitialized(log types.Log) (*ZenTestnetInitialized, error) {
	event := new(ZenTestnetInitialized)
	if err := _ZenTestnet.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenTestnetOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZenTestnet contract.
type ZenTestnetOwnershipTransferredIterator struct {
	Event *ZenTestnetOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZenTestnetOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenTestnetOwnershipTransferred)
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
		it.Event = new(ZenTestnetOwnershipTransferred)
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
func (it *ZenTestnetOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenTestnetOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenTestnetOwnershipTransferred represents a OwnershipTransferred event raised by the ZenTestnet contract.
type ZenTestnetOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenTestnet *ZenTestnetFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZenTestnetOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenTestnet.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZenTestnetOwnershipTransferredIterator{contract: _ZenTestnet.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenTestnet *ZenTestnetFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZenTestnetOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenTestnet.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenTestnetOwnershipTransferred)
				if err := _ZenTestnet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ZenTestnet *ZenTestnetFilterer) ParseOwnershipTransferred(log types.Log) (*ZenTestnetOwnershipTransferred, error) {
	event := new(ZenTestnetOwnershipTransferred)
	if err := _ZenTestnet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenTestnetTransactionProcessedIterator is returned from FilterTransactionProcessed and is used to iterate over the raw logs and unpacked data for TransactionProcessed events raised by the ZenTestnet contract.
type ZenTestnetTransactionProcessedIterator struct {
	Event *ZenTestnetTransactionProcessed // Event containing the contract specifics and raw log

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
func (it *ZenTestnetTransactionProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenTestnetTransactionProcessed)
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
		it.Event = new(ZenTestnetTransactionProcessed)
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
func (it *ZenTestnetTransactionProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenTestnetTransactionProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenTestnetTransactionProcessed represents a TransactionProcessed event raised by the ZenTestnet contract.
type ZenTestnetTransactionProcessed struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransactionProcessed is a free log retrieval operation binding the contract event 0xe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d1.
//
// Solidity: event TransactionProcessed(address sender, uint256 amount)
func (_ZenTestnet *ZenTestnetFilterer) FilterTransactionProcessed(opts *bind.FilterOpts) (*ZenTestnetTransactionProcessedIterator, error) {

	logs, sub, err := _ZenTestnet.contract.FilterLogs(opts, "TransactionProcessed")
	if err != nil {
		return nil, err
	}
	return &ZenTestnetTransactionProcessedIterator{contract: _ZenTestnet.contract, event: "TransactionProcessed", logs: logs, sub: sub}, nil
}

// WatchTransactionProcessed is a free log subscription operation binding the contract event 0xe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d1.
//
// Solidity: event TransactionProcessed(address sender, uint256 amount)
func (_ZenTestnet *ZenTestnetFilterer) WatchTransactionProcessed(opts *bind.WatchOpts, sink chan<- *ZenTestnetTransactionProcessed) (event.Subscription, error) {

	logs, sub, err := _ZenTestnet.contract.WatchLogs(opts, "TransactionProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenTestnetTransactionProcessed)
				if err := _ZenTestnet.contract.UnpackLog(event, "TransactionProcessed", log); err != nil {
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

// ParseTransactionProcessed is a log parse operation binding the contract event 0xe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d1.
//
// Solidity: event TransactionProcessed(address sender, uint256 amount)
func (_ZenTestnet *ZenTestnetFilterer) ParseTransactionProcessed(log types.Log) (*ZenTestnetTransactionProcessed, error) {
	event := new(ZenTestnetTransactionProcessed)
	if err := _ZenTestnet.contract.UnpackLog(event, "TransactionProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenTestnetTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ZenTestnet contract.
type ZenTestnetTransferIterator struct {
	Event *ZenTestnetTransfer // Event containing the contract specifics and raw log

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
func (it *ZenTestnetTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenTestnetTransfer)
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
		it.Event = new(ZenTestnetTransfer)
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
func (it *ZenTestnetTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenTestnetTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenTestnetTransfer represents a Transfer event raised by the ZenTestnet contract.
type ZenTestnetTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenTestnet *ZenTestnetFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZenTestnetTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenTestnet.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZenTestnetTransferIterator{contract: _ZenTestnet.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenTestnet *ZenTestnetFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZenTestnetTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenTestnet.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenTestnetTransfer)
				if err := _ZenTestnet.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenTestnet *ZenTestnetFilterer) ParseTransfer(log types.Log) (*ZenTestnetTransfer, error) {
	event := new(ZenTestnetTransfer)
	if err := _ZenTestnet.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
