// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package WrappedERC20

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

// WrappedERC20MetaData contains all meta data concerning the WrappedERC20 contract.
var WrappedERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"giver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issueFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461002b5761001a61001461013b565b9061015d565b60405161115e61049d823961115e90f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b90601f01601f191681019081106001600160401b0382111761006457604052565b61002f565b9061007d61007660405190565b9283610043565b565b6001600160401b03811161006457602090601f01601f19160190565b90825f9392825e0152565b909291926100bb6100b68261007f565b610069565b938185528183011161002b5761007d91602085019061009b565b9080601f8301121561002b5781516100ef926020016100a6565b90565b91909160408184031261002b5780516001600160401b03811161002b578361011b9183016100d5565b60208201519093906001600160401b03811161002b576100ef92016100d5565b6101596115fb8038038061014e81610069565b9283398101906100f2565b9091565b9061016791610194565b610191337fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756103f1565b50565b61007d91829182918291829190610371565b634e487b7160e01b5f52602260045260245ffd5b90600160028304921680156101da575b60208310146101d557565b6101a6565b91607f16916101ca565b915f1960089290920291821b911b5b9181191691161790565b6100ef6100ef6100ef9290565b919061021b6100ef610223936101fd565b9083546101e4565b9055565b61007d915f9161020a565b81811061023d575050565b8061024a5f600193610227565b01610232565b9190601f811161025f57505050565b61026f61007d935f5260205f2090565b906020601f840181900483019310610291575b6020601f909101040190610232565b9091508190610282565b906102a4815190565b906001600160401b038211610064576102c7826102c185546101ba565b85610250565b602090601f83116001146103005761022392915f91836102f5575b50505f19600883021c1916906002021790565b015190505f806102e2565b601f19831691610313855f5260205f2090565b925f5b81811061034f57509160029391856001969410610337575b50505002019055565b01515f196008601f8516021c191690555f808061032e565b91936020600181928787015181550195019201610316565b9061007d9161029b565b9061038061007d926003610367565b6004610367565b905b5f5260205260405f2090565b6100ef906103a9906001600160a01b031682565b6001600160a01b031690565b6100ef90610395565b6100ef906103b5565b90610389906103be565b9060ff906101f3565b906103ea6100ef61022392151590565b82546103d1565b6104026103fe838361047a565b1590565b15610474576104276001610422845f61041c866005610387565b016103c7565b6103da565b61044161043b610435339390565b936103be565b916103be565b917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d61046c60405190565b5f90a4600190565b50505f90565b6100ef915f61041c6104959361048d5f90565b506005610387565b5460ff169056fe60806040526004361015610011575f80fd5b5f3560e01c806301ffc9a71461014057806306fdde031461013b578063095ea7b31461013657806318160ddd146101315780631dd319cb1461012c57806323b872dd14610127578063248a9ca3146101225780632f2ff15d1461011d578063313ce5671461011857806336568abe1461011357806370a082311461010e57806375b238fc1461010957806391d148541461010457806395d89b41146100ff578063979005ad146100fa578063a217fddf146100f5578063a9059cbb146100f0578063d547741f146100eb5763dd62ed3e036101705761054e565b610519565b6104fd565b6104e2565b6104ab565b610490565b610474565b61043b565b610420565b6103f3565b6103c4565b6103ab565b61036d565b61033d565b6102f6565b6102ca565b6102ae565b610229565b61019a565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b0361017057565b5f80fd5b9050359061018182610145565b565b906020828203126101705761019791610174565b90565b34610170576101c76101b56101b0366004610183565b61056a565b60405191829182901515815260200190565b0390f35b5f91031261017057565b90825f9392825e0152565b61020161020a602093610214936101f5815190565b80835293849260200190565b958691016101d5565b601f01601f191690565b0190565b6020808252610197929101906101e0565b34610170576102393660046101cb565b6101c7610244610733565b60405191829182610218565b6001600160a01b031690565b6001600160a01b038116610169565b905035906101818261025c565b80610169565b9050359061018182610278565b9190604083820312610170576101979060206102a7828661026b565b940161027e565b34610170576101c76101b56102c436600461028b565b9061073d565b34610170576102da3660046101cb565b6101c76102e561075e565b6040515b9182918290815260200190565b346101705761030f61030936600461028b565b9061080a565b604051005b90916060828403126101705761019761032d848461026b565b9360406102a7826020870161026b565b34610170576101c76101b5610353366004610314565b91610814565b90602082820312610170576101979161027e565b34610170576101c76102e5610383366004610359565b610838565b9190604083820312610170576101979060206103a4828661027e565b940161026b565b346101705761030f6103be366004610388565b90610874565b34610170576103d43660046101cb565b6101c76103df610891565b6040519182918260ff909116815260200190565b346101705761030f610406366004610388565b9061089b565b90602082820312610170576101979161026b565b34610170576101c76102e561043636600461040c565b6108e9565b346101705761044b3660046101cb565b6101c77fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756102e5565b34610170576101c76101b561048a366004610388565b9061099c565b34610170576104a03660046101cb565b6101c76102446109ba565b346101705761030f6104be36600461028b565b906109fc565b6101976101976101979290565b6101975f6104c4565b6101976104d1565b34610170576104f23660046101cb565b6101c76102e56104da565b34610170576101c76101b561051336600461028b565b90610a06565b346101705761030f61052c366004610388565b90610a2c565b9190604083820312610170576101979060206103a4828661026b565b34610170576101c76102e5610564366004610532565b90610a36565b7f7965db0b000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008216149081156105ba575090565b61019791507fffffffff00000000000000000000000000000000000000000000000000000000167f01ffc9a7000000000000000000000000000000000000000000000000000000001490565b634e487b7160e01b5f52602260045260245ffd5b906001600283049216801561063a575b602083101461063557565b610606565b91607f169161062a565b80545f9392916106606106568361061a565b8085529360200190565b91600181169081156106af575060011461067957505050565b61068a91929394505f5260205f2090565b915f925b81841061069b5750500190565b80548484015260209093019260010161068e565b92949550505060ff1916825215156020020190565b9061019791610644565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761070457604052565b6106ce565b906101816107239261071a60405190565b938480926106c4565b03836106e2565b61019790610709565b610197600361072a565b610748919033610afd565b600190565b6101979081565b610197905461074d565b6101976002610754565b610181919061079b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610b0a565b610b0a565b6107ec565b156107a757565b60405162461bcd60e51b815260206004820152601560248201527f496e73756666696369656e742062616c616e63652e00000000000000000000006044820152606490fd5b90610181916108056107fd826108e9565b8311156107a0565b610b2b565b9061018191610768565b610748929190610825833383610ba4565b610c15565b905b5f5260205260405f2090565b6001610850610197926108485f90565b50600561082a565b01610754565b906101819161086761079682610838565b9061087191610cc5565b50565b9061018191610856565b61088b6101976101979290565b60ff1690565b610197601261087e565b906108a533610250565b6001600160a01b038216036108bd5761087191610d48565b7f6697b232000000000000000000000000000000000000000000000000000000005f90815260045b035ffd5b6001600160a01b038116321461094e576001600160a01b038116331461094e5760405162461bcd60e51b815260206004820152601f60248201527f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e6365006044820152606490fd5b61019790610da4565b610250610197610197926001600160a01b031690565b61019790610957565b6101979061096d565b9061082c90610976565b6101979061088b565b6101979054610989565b610197915f6109af6109b5936108485f90565b0161097f565b610992565b610197600461072a565b61018191906109f27fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610b0a565b9061018191610dbe565b90610181916109c4565b610748919033610c15565b9061018191610a2261079682610838565b9061087191610d48565b9061018191610a11565b906001600160a01b03821632148015610aeb575b610ad0576001600160a01b03821633148015610ad9575b610ad05760405162461bcd60e51b815260206004820152602160248201527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360448201527f65000000000000000000000000000000000000000000000000000000000000006064820152608490fd5b61019791610e0b565b506001600160a01b0381163314610a61565b506001600160a01b0381163214610a4a565b9160019161018193610e4a565b610181903390610f78565b6102506101976101979290565b61019790610b15565b9190610b365f610b22565b926001600160a01b0384166001600160a01b03821614610b5b57926101819293610fe5565b634b637e8f60e11b5f9081526001600160a01b038516600452602490fd5b6001600160a01b039091168152606081019392610181929091604091610ba0906020830152565b0152565b91610baf8284610a36565b5f198110610bbe575b50505050565b818110610be35791610bd4610bda94925f940390565b91610e4a565b5f808080610bb8565b7ffb8f41b2000000000000000000000000000000000000000000000000000000005f90815293506108e5926004610b79565b929190610c215f610b22565b936001600160a01b0385166001600160a01b03821614610c7d576001600160a01b0385166001600160a01b03831614610c5f57610181939450610fe5565b63ec442f0560e01b5f9081526001600160a01b038616600452602490fd5b634b637e8f60e11b5f9081526001600160a01b038616600452602490fd5b9060ff905b9181191691161790565b90610cba610197610cc192151590565b8254610c9b565b9055565b610cd6610cd2838361099c565b1590565b15610d4257610cf56001610cf0845f6109af86600561082a565b610caa565b610d0f610d09610d03339390565b93610976565b91610976565b917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d610d3a60405190565b5f90a4600190565b50505f90565b610d52828261099c565b15610d4257610d6b5f610cf084826109af86600561082a565b610d79610d09610d03339390565b917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b610d3a60405190565b610db961019791610db25f90565b505f61097f565b610754565b9190610dc95f610b22565b926001600160a01b0384166001600160a01b03821614610ded576101819293610fe5565b63ec442f0560e01b5f9081526001600160a01b038516600452602490fd5b61019791610e25610db992610e1d5f90565b50600161097f565b61097f565b905f1990610ca0565b90610e43610197610cc1926104c4565b8254610e2a565b909192610e565f610b22565b6001600160a01b0381166001600160a01b03841614610f22576001600160a01b0381166001600160a01b03851614610ee95750610ea184610e9c85610e2586600161097f565b610e33565b610eaa57505050565b610ee4610eda610d037f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92593610976565b936102e960405190565b0390a3565b7f94280d62000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b7fe602df05000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b6001600160a01b0390911681526040810192916101819160200152565b90610f86610cd2828461099c565b610f8e575050565b7fe2517d3f000000000000000000000000000000000000000000000000000000005f908152916108e5916004610f5b565b634e487b7160e01b5f52601160045260245ffd5b91908201809211610fe057565b610fbf565b610fee5f610b22565b6001600160a01b0381166001600160a01b0383160361109c57610eda610d037fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef93611053610ee49461025061104c8a6110476002610754565b610fd3565b6002610e33565b6001600160a01b0387160361107c5761107761104c886110736002610754565b0390565b610976565b611077611089875f61097f565b6110968961021483610754565b90610e33565b6110a9610db9835f61097f565b8481106110f557610d037fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef93611053610ee4946102506110eb8a610eda970390565b610e9c855f61097f565b6108e5855f92857fe450d38c00000000000000000000000000000000000000000000000000000000855260048501610b7956fea2646970667358221220ed36be82a8df75b008b97873a55ebd9856720cae81aa177f283e591a11a7ce1b64736f6c634300081c0033",
}

// WrappedERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use WrappedERC20MetaData.ABI instead.
var WrappedERC20ABI = WrappedERC20MetaData.ABI

// WrappedERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WrappedERC20MetaData.Bin instead.
var WrappedERC20Bin = WrappedERC20MetaData.Bin

// DeployWrappedERC20 deploys a new Ethereum contract, binding an instance of WrappedERC20 to it.
func DeployWrappedERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *WrappedERC20, error) {
	parsed, err := WrappedERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WrappedERC20Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WrappedERC20{WrappedERC20Caller: WrappedERC20Caller{contract: contract}, WrappedERC20Transactor: WrappedERC20Transactor{contract: contract}, WrappedERC20Filterer: WrappedERC20Filterer{contract: contract}}, nil
}

// WrappedERC20 is an auto generated Go binding around an Ethereum contract.
type WrappedERC20 struct {
	WrappedERC20Caller     // Read-only binding to the contract
	WrappedERC20Transactor // Write-only binding to the contract
	WrappedERC20Filterer   // Log filterer for contract events
}

// WrappedERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type WrappedERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type WrappedERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WrappedERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappedERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WrappedERC20Session struct {
	Contract     *WrappedERC20     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WrappedERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WrappedERC20CallerSession struct {
	Contract *WrappedERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// WrappedERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WrappedERC20TransactorSession struct {
	Contract     *WrappedERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// WrappedERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type WrappedERC20Raw struct {
	Contract *WrappedERC20 // Generic contract binding to access the raw methods on
}

// WrappedERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WrappedERC20CallerRaw struct {
	Contract *WrappedERC20Caller // Generic read-only contract binding to access the raw methods on
}

// WrappedERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WrappedERC20TransactorRaw struct {
	Contract *WrappedERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewWrappedERC20 creates a new instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20(address common.Address, backend bind.ContractBackend) (*WrappedERC20, error) {
	contract, err := bindWrappedERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20{WrappedERC20Caller: WrappedERC20Caller{contract: contract}, WrappedERC20Transactor: WrappedERC20Transactor{contract: contract}, WrappedERC20Filterer: WrappedERC20Filterer{contract: contract}}, nil
}

// NewWrappedERC20Caller creates a new read-only instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20Caller(address common.Address, caller bind.ContractCaller) (*WrappedERC20Caller, error) {
	contract, err := bindWrappedERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20Caller{contract: contract}, nil
}

// NewWrappedERC20Transactor creates a new write-only instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*WrappedERC20Transactor, error) {
	contract, err := bindWrappedERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20Transactor{contract: contract}, nil
}

// NewWrappedERC20Filterer creates a new log filterer instance of WrappedERC20, bound to a specific deployed contract.
func NewWrappedERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*WrappedERC20Filterer, error) {
	contract, err := bindWrappedERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20Filterer{contract: contract}, nil
}

// bindWrappedERC20 binds a generic wrapper to an already deployed contract.
func bindWrappedERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WrappedERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappedERC20 *WrappedERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedERC20.Contract.WrappedERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappedERC20 *WrappedERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedERC20.Contract.WrappedERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappedERC20 *WrappedERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedERC20.Contract.WrappedERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappedERC20 *WrappedERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappedERC20 *WrappedERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappedERC20 *WrappedERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedERC20.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Caller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Session) ADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.ADMINROLE(&_WrappedERC20.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20CallerSession) ADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.ADMINROLE(&_WrappedERC20.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.DEFAULTADMINROLE(&_WrappedERC20.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_WrappedERC20 *WrappedERC20CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _WrappedERC20.Contract.DEFAULTADMINROLE(&_WrappedERC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.Allowance(&_WrappedERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WrappedERC20 *WrappedERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.Allowance(&_WrappedERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedERC20 *WrappedERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.BalanceOf(&_WrappedERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WrappedERC20 *WrappedERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedERC20.Contract.BalanceOf(&_WrappedERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedERC20 *WrappedERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedERC20 *WrappedERC20Session) Decimals() (uint8, error) {
	return _WrappedERC20.Contract.Decimals(&_WrappedERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WrappedERC20 *WrappedERC20CallerSession) Decimals() (uint8, error) {
	return _WrappedERC20.Contract.Decimals(&_WrappedERC20.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedERC20 *WrappedERC20Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _WrappedERC20.Contract.GetRoleAdmin(&_WrappedERC20.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_WrappedERC20 *WrappedERC20CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _WrappedERC20.Contract.GetRoleAdmin(&_WrappedERC20.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedERC20 *WrappedERC20Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedERC20 *WrappedERC20Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _WrappedERC20.Contract.HasRole(&_WrappedERC20.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_WrappedERC20 *WrappedERC20CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _WrappedERC20.Contract.HasRole(&_WrappedERC20.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedERC20 *WrappedERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedERC20 *WrappedERC20Session) Name() (string, error) {
	return _WrappedERC20.Contract.Name(&_WrappedERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WrappedERC20 *WrappedERC20CallerSession) Name() (string, error) {
	return _WrappedERC20.Contract.Name(&_WrappedERC20.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedERC20 *WrappedERC20Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedERC20 *WrappedERC20Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WrappedERC20.Contract.SupportsInterface(&_WrappedERC20.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WrappedERC20 *WrappedERC20CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WrappedERC20.Contract.SupportsInterface(&_WrappedERC20.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedERC20 *WrappedERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedERC20 *WrappedERC20Session) Symbol() (string, error) {
	return _WrappedERC20.Contract.Symbol(&_WrappedERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WrappedERC20 *WrappedERC20CallerSession) Symbol() (string, error) {
	return _WrappedERC20.Contract.Symbol(&_WrappedERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedERC20 *WrappedERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WrappedERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedERC20 *WrappedERC20Session) TotalSupply() (*big.Int, error) {
	return _WrappedERC20.Contract.TotalSupply(&_WrappedERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WrappedERC20 *WrappedERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _WrappedERC20.Contract.TotalSupply(&_WrappedERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Approve(&_WrappedERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Approve(&_WrappedERC20.TransactOpts, spender, value)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Transactor) BurnFor(opts *bind.TransactOpts, giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "burnFor", giver, amount)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Session) BurnFor(giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.BurnFor(&_WrappedERC20.TransactOpts, giver, amount)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) BurnFor(giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.BurnFor(&_WrappedERC20.TransactOpts, giver, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.GrantRole(&_WrappedERC20.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.GrantRole(&_WrappedERC20.TransactOpts, role, account)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Transactor) IssueFor(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "issueFor", receiver, amount)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20Session) IssueFor(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.IssueFor(&_WrappedERC20.TransactOpts, receiver, amount)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) IssueFor(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.IssueFor(&_WrappedERC20.TransactOpts, receiver, amount)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_WrappedERC20 *WrappedERC20Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_WrappedERC20 *WrappedERC20Session) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RenounceRole(&_WrappedERC20.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RenounceRole(&_WrappedERC20.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RevokeRole(&_WrappedERC20.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_WrappedERC20 *WrappedERC20TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _WrappedERC20.Contract.RevokeRole(&_WrappedERC20.TransactOpts, role, account)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Transfer(&_WrappedERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.Transfer(&_WrappedERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.TransferFrom(&_WrappedERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WrappedERC20 *WrappedERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WrappedERC20.Contract.TransferFrom(&_WrappedERC20.TransactOpts, from, to, value)
}

// WrappedERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WrappedERC20 contract.
type WrappedERC20ApprovalIterator struct {
	Event *WrappedERC20Approval // Event containing the contract specifics and raw log

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
func (it *WrappedERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20Approval)
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
		it.Event = new(WrappedERC20Approval)
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
func (it *WrappedERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20Approval represents a Approval event raised by the WrappedERC20 contract.
type WrappedERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WrappedERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20ApprovalIterator{contract: _WrappedERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WrappedERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20Approval)
				if err := _WrappedERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_WrappedERC20 *WrappedERC20Filterer) ParseApproval(log types.Log) (*WrappedERC20Approval, error) {
	event := new(WrappedERC20Approval)
	if err := _WrappedERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the WrappedERC20 contract.
type WrappedERC20RoleAdminChangedIterator struct {
	Event *WrappedERC20RoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *WrappedERC20RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20RoleAdminChanged)
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
		it.Event = new(WrappedERC20RoleAdminChanged)
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
func (it *WrappedERC20RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20RoleAdminChanged represents a RoleAdminChanged event raised by the WrappedERC20 contract.
type WrappedERC20RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedERC20 *WrappedERC20Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*WrappedERC20RoleAdminChangedIterator, error) {

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

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20RoleAdminChangedIterator{contract: _WrappedERC20.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_WrappedERC20 *WrappedERC20Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *WrappedERC20RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20RoleAdminChanged)
				if err := _WrappedERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_WrappedERC20 *WrappedERC20Filterer) ParseRoleAdminChanged(log types.Log) (*WrappedERC20RoleAdminChanged, error) {
	event := new(WrappedERC20RoleAdminChanged)
	if err := _WrappedERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the WrappedERC20 contract.
type WrappedERC20RoleGrantedIterator struct {
	Event *WrappedERC20RoleGranted // Event containing the contract specifics and raw log

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
func (it *WrappedERC20RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20RoleGranted)
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
		it.Event = new(WrappedERC20RoleGranted)
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
func (it *WrappedERC20RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20RoleGranted represents a RoleGranted event raised by the WrappedERC20 contract.
type WrappedERC20RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*WrappedERC20RoleGrantedIterator, error) {

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

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20RoleGrantedIterator{contract: _WrappedERC20.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *WrappedERC20RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20RoleGranted)
				if err := _WrappedERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_WrappedERC20 *WrappedERC20Filterer) ParseRoleGranted(log types.Log) (*WrappedERC20RoleGranted, error) {
	event := new(WrappedERC20RoleGranted)
	if err := _WrappedERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the WrappedERC20 contract.
type WrappedERC20RoleRevokedIterator struct {
	Event *WrappedERC20RoleRevoked // Event containing the contract specifics and raw log

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
func (it *WrappedERC20RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20RoleRevoked)
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
		it.Event = new(WrappedERC20RoleRevoked)
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
func (it *WrappedERC20RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20RoleRevoked represents a RoleRevoked event raised by the WrappedERC20 contract.
type WrappedERC20RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*WrappedERC20RoleRevokedIterator, error) {

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

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20RoleRevokedIterator{contract: _WrappedERC20.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_WrappedERC20 *WrappedERC20Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *WrappedERC20RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20RoleRevoked)
				if err := _WrappedERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_WrappedERC20 *WrappedERC20Filterer) ParseRoleRevoked(log types.Log) (*WrappedERC20RoleRevoked, error) {
	event := new(WrappedERC20RoleRevoked)
	if err := _WrappedERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WrappedERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WrappedERC20 contract.
type WrappedERC20TransferIterator struct {
	Event *WrappedERC20Transfer // Event containing the contract specifics and raw log

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
func (it *WrappedERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedERC20Transfer)
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
		it.Event = new(WrappedERC20Transfer)
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
func (it *WrappedERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappedERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappedERC20Transfer represents a Transfer event raised by the WrappedERC20 contract.
type WrappedERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WrappedERC20TransferIterator{contract: _WrappedERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WrappedERC20 *WrappedERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WrappedERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappedERC20Transfer)
				if err := _WrappedERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_WrappedERC20 *WrappedERC20Filterer) ParseTransfer(log types.Log) (*WrappedERC20Transfer, error) {
	event := new(WrappedERC20Transfer)
	if err := _WrappedERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
