// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EthERC20

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
)

// EthERC20MetaData contains all meta data concerning the EthERC20 contract.
var EthERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"managementContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"hm\",\"type\":\"string\"}],\"name\":\"Something\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200252438038062002524833981810160405281019062000037919062000736565b838381600390816200004a919062000a27565b5080600490816200005c919062000a27565b505050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161462000147578073ffffffffffffffffffffffffffffffffffffffff1663a1a227fa6040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000e0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000106919062000b53565b600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505b80600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506200019a3383620001a460201b60201c565b5050505062000eef565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000216576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200020d9062000be6565b60405180910390fd5b6200022a600083836200031c60201b60201c565b80600260008282546200023e919062000c37565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825462000295919062000c37565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620002fc919062000c83565b60405180910390a36200031860008383620004fe60201b60201c565b5050565b600073ffffffffffffffffffffffffffffffffffffffff16600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160315620004f957600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620004f857600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152509050600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa436001808111156200046a576200046962000ca0565b5b846040516020016200047d919062000d39565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401620004af949392919062000e24565b6020604051808303816000875af1158015620004cf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620004f5919062000ebd565b50505b5b505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200056c8262000521565b810181811067ffffffffffffffff821117156200058e576200058d62000532565b5b80604052505050565b6000620005a362000503565b9050620005b1828262000561565b919050565b600067ffffffffffffffff821115620005d457620005d362000532565b5b620005df8262000521565b9050602081019050919050565b60005b838110156200060c578082015181840152602081019050620005ef565b60008484015250505050565b60006200062f6200062984620005b6565b62000597565b9050828152602081018484840111156200064e576200064d6200051c565b5b6200065b848285620005ec565b509392505050565b600082601f8301126200067b576200067a62000517565b5b81516200068d84826020860162000618565b91505092915050565b6000819050919050565b620006ab8162000696565b8114620006b757600080fd5b50565b600081519050620006cb81620006a0565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620006fe82620006d1565b9050919050565b6200071081620006f1565b81146200071c57600080fd5b50565b600081519050620007308162000705565b92915050565b600080600080608085870312156200075357620007526200050d565b5b600085015167ffffffffffffffff81111562000774576200077362000512565b5b620007828782880162000663565b945050602085015167ffffffffffffffff811115620007a657620007a562000512565b5b620007b48782880162000663565b9350506040620007c787828801620006ba565b9250506060620007da878288016200071f565b91505092959194509250565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200083957607f821691505b6020821081036200084f576200084e620007f1565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620008b97fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826200087a565b620008c586836200087a565b95508019841693508086168417925050509392505050565b6000819050919050565b60006200090862000902620008fc8462000696565b620008dd565b62000696565b9050919050565b6000819050919050565b6200092483620008e7565b6200093c62000933826200090f565b84845462000887565b825550505050565b600090565b6200095362000944565b6200096081848462000919565b505050565b5b8181101562000988576200097c60008262000949565b60018101905062000966565b5050565b601f821115620009d757620009a18162000855565b620009ac846200086a565b81016020851015620009bc578190505b620009d4620009cb856200086a565b83018262000965565b50505b505050565b600082821c905092915050565b6000620009fc60001984600802620009dc565b1980831691505092915050565b600062000a178383620009e9565b9150826002028217905092915050565b62000a3282620007e6565b67ffffffffffffffff81111562000a4e5762000a4d62000532565b5b62000a5a825462000820565b62000a678282856200098c565b600060209050601f83116001811462000a9f576000841562000a8a578287015190505b62000a96858262000a09565b86555062000b06565b601f19841662000aaf8662000855565b60005b8281101562000ad95784890151825560018201915060208501945060208101905062000ab2565b8683101562000af9578489015162000af5601f891682620009e9565b8355505b6001600288020188555050505b505050505050565b600062000b1b82620006f1565b9050919050565b62000b2d8162000b0e565b811462000b3957600080fd5b50565b60008151905062000b4d8162000b22565b92915050565b60006020828403121562000b6c5762000b6b6200050d565b5b600062000b7c8482850162000b3c565b91505092915050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000bce601f8362000b85565b915062000bdb8262000b96565b602082019050919050565b6000602082019050818103600083015262000c018162000bbf565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600062000c448262000696565b915062000c518362000696565b925082820190508082111562000c6c5762000c6b62000c08565b5b92915050565b62000c7d8162000696565b82525050565b600060208201905062000c9a600083018462000c72565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b62000cda81620006f1565b82525050565b62000ceb8162000696565b82525050565b60608201600082015162000d09600085018262000ccf565b50602082015162000d1e602085018262000ccf565b50604082015162000d33604085018262000ce0565b50505050565b600060608201905062000d50600083018462000cf1565b92915050565b600063ffffffff82169050919050565b62000d718162000d56565b82525050565b600081519050919050565b600082825260208201905092915050565b600062000da08262000d77565b62000dac818562000d82565b935062000dbe818560208601620005ec565b62000dc98162000521565b840191505092915050565b6000819050919050565b600060ff82169050919050565b600062000e0c62000e0662000e008462000dd4565b620008dd565b62000dde565b9050919050565b62000e1e8162000deb565b82525050565b600060808201905062000e3b600083018762000d66565b62000e4a602083018662000d66565b818103604083015262000e5e818562000d93565b905062000e6f606083018462000e13565b95945050505050565b600067ffffffffffffffff82169050919050565b62000e978162000e78565b811462000ea357600080fd5b50565b60008151905062000eb78162000e8c565b92915050565b60006020828403121562000ed65762000ed56200050d565b5b600062000ee68482850162000ea6565b91505092915050565b6116258062000eff6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610ce6565b60405180910390f35b6100e660048036038101906100e19190610da1565b610308565b6040516100f39190610dfc565b60405180910390f35b61010461032b565b6040516101119190610e26565b60405180910390f35b610134600480360381019061012f9190610e41565b610335565b6040516101419190610dfc565b60405180910390f35b610152610364565b60405161015f9190610eb0565b60405180910390f35b610182600480360381019061017d9190610da1565b61036d565b60405161018f9190610dfc565b60405180910390f35b6101b260048036038101906101ad9190610ecb565b6103a4565b6040516101bf9190610e26565b60405180910390f35b6101d06103ec565b6040516101dd9190610ce6565b60405180910390f35b61020060048036038101906101fb9190610da1565b61047e565b60405161020d9190610dfc565b60405180910390f35b610230600480360381019061022b9190610da1565b6104f5565b60405161023d9190610dfc565b60405180910390f35b610260600480360381019061025b9190610ef8565b610518565b60405161026d9190610e26565b60405180910390f35b60606003805461028590610f67565b80601f01602080910402602001604051908101604052809291908181526020018280546102b190610f67565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b60008061031361059f565b90506103208185856105a7565b600191505092915050565b6000600254905090565b60008061034061059f565b905061034d858285610770565b6103588585856107fc565b60019150509392505050565b60006012905090565b60008061037861059f565b905061039981858561038a8589610518565b6103949190610fc7565b6105a7565b600191505092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6060600480546103fb90610f67565b80601f016020809104026020016040519081016040528092919081815260200182805461042790610f67565b80156104745780601f1061044957610100808354040283529160200191610474565b820191906000526020600020905b81548152906001019060200180831161045757829003601f168201915b5050505050905090565b60008061048961059f565b905060006104978286610518565b9050838110156104dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d39061106d565b60405180910390fd5b6104e982868684036105a7565b60019250505092915050565b60008061050061059f565b905061050d8185856107fc565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610616576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161060d906110ff565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610685576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067c90611191565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516107639190610e26565b60405180910390a3505050565b600061077c8484610518565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107f657818110156107e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107df906111fd565b60405180910390fd5b6107f584848484036105a7565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361086b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108629061128f565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108d190611321565b60405180910390fd5b6108e5838383610a7b565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561096b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610962906113b3565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109fe9190610fc7565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610a629190610e26565b60405180910390a3610a75848484610c51565b50505050565b600073ffffffffffffffffffffffffffffffffffffffff16600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160315610c4c57600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610c4b57600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff168152602001838152509050600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115610bc457610bc36113d3565b5b84604051602001610bd59190611462565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401610c059493929190611536565b6020604051808303816000875af1158015610c24573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c4891906115c2565b50505b5b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610c90578082015181840152602081019050610c75565b60008484015250505050565b6000601f19601f8301169050919050565b6000610cb882610c56565b610cc28185610c61565b9350610cd2818560208601610c72565b610cdb81610c9c565b840191505092915050565b60006020820190508181036000830152610d008184610cad565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610d3882610d0d565b9050919050565b610d4881610d2d565b8114610d5357600080fd5b50565b600081359050610d6581610d3f565b92915050565b6000819050919050565b610d7e81610d6b565b8114610d8957600080fd5b50565b600081359050610d9b81610d75565b92915050565b60008060408385031215610db857610db7610d08565b5b6000610dc685828601610d56565b9250506020610dd785828601610d8c565b9150509250929050565b60008115159050919050565b610df681610de1565b82525050565b6000602082019050610e116000830184610ded565b92915050565b610e2081610d6b565b82525050565b6000602082019050610e3b6000830184610e17565b92915050565b600080600060608486031215610e5a57610e59610d08565b5b6000610e6886828701610d56565b9350506020610e7986828701610d56565b9250506040610e8a86828701610d8c565b9150509250925092565b600060ff82169050919050565b610eaa81610e94565b82525050565b6000602082019050610ec56000830184610ea1565b92915050565b600060208284031215610ee157610ee0610d08565b5b6000610eef84828501610d56565b91505092915050565b60008060408385031215610f0f57610f0e610d08565b5b6000610f1d85828601610d56565b9250506020610f2e85828601610d56565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610f7f57607f821691505b602082108103610f9257610f91610f38565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610fd282610d6b565b9150610fdd83610d6b565b9250828201905080821115610ff557610ff4610f98565b5b92915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b6000611057602583610c61565b915061106282610ffb565b604082019050919050565b600060208201905081810360008301526110868161104a565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006110e9602483610c61565b91506110f48261108d565b604082019050919050565b60006020820190508181036000830152611118816110dc565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b600061117b602283610c61565b91506111868261111f565b604082019050919050565b600060208201905081810360008301526111aa8161116e565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b60006111e7601d83610c61565b91506111f2826111b1565b602082019050919050565b60006020820190508181036000830152611216816111da565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b6000611279602583610c61565b91506112848261121d565b604082019050919050565b600060208201905081810360008301526112a88161126c565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b600061130b602383610c61565b9150611316826112af565b604082019050919050565b6000602082019050818103600083015261133a816112fe565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b600061139d602683610c61565b91506113a882611341565b604082019050919050565b600060208201905081810360008301526113cc81611390565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61140b81610d2d565b82525050565b61141a81610d6b565b82525050565b6060820160008201516114366000850182611402565b5060208201516114496020850182611402565b50604082015161145c6040850182611411565b50505050565b60006060820190506114776000830184611420565b92915050565b600063ffffffff82169050919050565b6114968161147d565b82525050565b600081519050919050565b600082825260208201905092915050565b60006114c38261149c565b6114cd81856114a7565b93506114dd818560208601610c72565b6114e681610c9c565b840191505092915050565b6000819050919050565b6000819050919050565b600061152061151b611516846114f1565b6114fb565b610e94565b9050919050565b61153081611505565b82525050565b600060808201905061154b600083018761148d565b611558602083018661148d565b818103604083015261156a81856114b8565b90506115796060830184611527565b95945050505050565b600067ffffffffffffffff82169050919050565b61159f81611582565b81146115aa57600080fd5b50565b6000815190506115bc81611596565b92915050565b6000602082840312156115d8576115d7610d08565b5b60006115e6848285016115ad565b9150509291505056fea264697066735822122011b5ab79b9d0368dcf162f8f00da0c4fb291e35d95b76dccfa7c8680fc0c8a6a64736f6c63430008110033",
}

// EthERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use EthERC20MetaData.ABI instead.
var EthERC20ABI = EthERC20MetaData.ABI

// EthERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EthERC20MetaData.Bin instead.
var EthERC20Bin = EthERC20MetaData.Bin

// DeployEthERC20 deploys a new Ethereum contract, binding an instance of EthERC20 to it.
func DeployEthERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int, managementContract common.Address) (common.Address, *types.Transaction, *EthERC20, error) {
	parsed, err := EthERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EthERC20Bin), backend, name, symbol, initialSupply, managementContract)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EthERC20{EthERC20Caller: EthERC20Caller{contract: contract}, EthERC20Transactor: EthERC20Transactor{contract: contract}, EthERC20Filterer: EthERC20Filterer{contract: contract}}, nil
}

// EthERC20 is an auto generated Go binding around an Ethereum contract.
type EthERC20 struct {
	EthERC20Caller     // Read-only binding to the contract
	EthERC20Transactor // Write-only binding to the contract
	EthERC20Filterer   // Log filterer for contract events
}

// EthERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type EthERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type EthERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthERC20Session struct {
	Contract     *EthERC20         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthERC20CallerSession struct {
	Contract *EthERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// EthERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthERC20TransactorSession struct {
	Contract     *EthERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// EthERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type EthERC20Raw struct {
	Contract *EthERC20 // Generic contract binding to access the raw methods on
}

// EthERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthERC20CallerRaw struct {
	Contract *EthERC20Caller // Generic read-only contract binding to access the raw methods on
}

// EthERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthERC20TransactorRaw struct {
	Contract *EthERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewEthERC20 creates a new instance of EthERC20, bound to a specific deployed contract.
func NewEthERC20(address common.Address, backend bind.ContractBackend) (*EthERC20, error) {
	contract, err := bindEthERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthERC20{EthERC20Caller: EthERC20Caller{contract: contract}, EthERC20Transactor: EthERC20Transactor{contract: contract}, EthERC20Filterer: EthERC20Filterer{contract: contract}}, nil
}

// NewEthERC20Caller creates a new read-only instance of EthERC20, bound to a specific deployed contract.
func NewEthERC20Caller(address common.Address, caller bind.ContractCaller) (*EthERC20Caller, error) {
	contract, err := bindEthERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthERC20Caller{contract: contract}, nil
}

// NewEthERC20Transactor creates a new write-only instance of EthERC20, bound to a specific deployed contract.
func NewEthERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*EthERC20Transactor, error) {
	contract, err := bindEthERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthERC20Transactor{contract: contract}, nil
}

// NewEthERC20Filterer creates a new log filterer instance of EthERC20, bound to a specific deployed contract.
func NewEthERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*EthERC20Filterer, error) {
	contract, err := bindEthERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthERC20Filterer{contract: contract}, nil
}

// bindEthERC20 binds a generic wrapper to an already deployed contract.
func bindEthERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthERC20 *EthERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthERC20.Contract.EthERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthERC20 *EthERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthERC20.Contract.EthERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthERC20 *EthERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthERC20.Contract.EthERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthERC20 *EthERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthERC20 *EthERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthERC20 *EthERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_EthERC20 *EthERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_EthERC20 *EthERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _EthERC20.Contract.Allowance(&_EthERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_EthERC20 *EthERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _EthERC20.Contract.Allowance(&_EthERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_EthERC20 *EthERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EthERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_EthERC20 *EthERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _EthERC20.Contract.BalanceOf(&_EthERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_EthERC20 *EthERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _EthERC20.Contract.BalanceOf(&_EthERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EthERC20 *EthERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _EthERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EthERC20 *EthERC20Session) Decimals() (uint8, error) {
	return _EthERC20.Contract.Decimals(&_EthERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EthERC20 *EthERC20CallerSession) Decimals() (uint8, error) {
	return _EthERC20.Contract.Decimals(&_EthERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EthERC20 *EthERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EthERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EthERC20 *EthERC20Session) Name() (string, error) {
	return _EthERC20.Contract.Name(&_EthERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EthERC20 *EthERC20CallerSession) Name() (string, error) {
	return _EthERC20.Contract.Name(&_EthERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EthERC20 *EthERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EthERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EthERC20 *EthERC20Session) Symbol() (string, error) {
	return _EthERC20.Contract.Symbol(&_EthERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EthERC20 *EthERC20CallerSession) Symbol() (string, error) {
	return _EthERC20.Contract.Symbol(&_EthERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_EthERC20 *EthERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EthERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_EthERC20 *EthERC20Session) TotalSupply() (*big.Int, error) {
	return _EthERC20.Contract.TotalSupply(&_EthERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_EthERC20 *EthERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _EthERC20.Contract.TotalSupply(&_EthERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.Approve(&_EthERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.Approve(&_EthERC20.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_EthERC20 *EthERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _EthERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_EthERC20 *EthERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.DecreaseAllowance(&_EthERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_EthERC20 *EthERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.DecreaseAllowance(&_EthERC20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_EthERC20 *EthERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _EthERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_EthERC20 *EthERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.IncreaseAllowance(&_EthERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_EthERC20 *EthERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.IncreaseAllowance(&_EthERC20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.Transfer(&_EthERC20.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.Transfer(&_EthERC20.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.TransferFrom(&_EthERC20.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_EthERC20 *EthERC20TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _EthERC20.Contract.TransferFrom(&_EthERC20.TransactOpts, from, to, amount)
}

// EthERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the EthERC20 contract.
type EthERC20ApprovalIterator struct {
	Event *EthERC20Approval // Event containing the contract specifics and raw log

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
func (it *EthERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthERC20Approval)
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
		it.Event = new(EthERC20Approval)
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
func (it *EthERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthERC20Approval represents a Approval event raised by the EthERC20 contract.
type EthERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_EthERC20 *EthERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*EthERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _EthERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &EthERC20ApprovalIterator{contract: _EthERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_EthERC20 *EthERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *EthERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _EthERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthERC20Approval)
				if err := _EthERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_EthERC20 *EthERC20Filterer) ParseApproval(log types.Log) (*EthERC20Approval, error) {
	event := new(EthERC20Approval)
	if err := _EthERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthERC20SomethingIterator is returned from FilterSomething and is used to iterate over the raw logs and unpacked data for Something events raised by the EthERC20 contract.
type EthERC20SomethingIterator struct {
	Event *EthERC20Something // Event containing the contract specifics and raw log

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
func (it *EthERC20SomethingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthERC20Something)
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
		it.Event = new(EthERC20Something)
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
func (it *EthERC20SomethingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthERC20SomethingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthERC20Something represents a Something event raised by the EthERC20 contract.
type EthERC20Something struct {
	Hm  string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterSomething is a free log retrieval operation binding the contract event 0xb3526317000a7a05cd5babafed7f7e3d5f0beb28668dc6ae1eb34911312f704a.
//
// Solidity: event Something(string hm)
func (_EthERC20 *EthERC20Filterer) FilterSomething(opts *bind.FilterOpts) (*EthERC20SomethingIterator, error) {

	logs, sub, err := _EthERC20.contract.FilterLogs(opts, "Something")
	if err != nil {
		return nil, err
	}
	return &EthERC20SomethingIterator{contract: _EthERC20.contract, event: "Something", logs: logs, sub: sub}, nil
}

// WatchSomething is a free log subscription operation binding the contract event 0xb3526317000a7a05cd5babafed7f7e3d5f0beb28668dc6ae1eb34911312f704a.
//
// Solidity: event Something(string hm)
func (_EthERC20 *EthERC20Filterer) WatchSomething(opts *bind.WatchOpts, sink chan<- *EthERC20Something) (event.Subscription, error) {

	logs, sub, err := _EthERC20.contract.WatchLogs(opts, "Something")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthERC20Something)
				if err := _EthERC20.contract.UnpackLog(event, "Something", log); err != nil {
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

// ParseSomething is a log parse operation binding the contract event 0xb3526317000a7a05cd5babafed7f7e3d5f0beb28668dc6ae1eb34911312f704a.
//
// Solidity: event Something(string hm)
func (_EthERC20 *EthERC20Filterer) ParseSomething(log types.Log) (*EthERC20Something, error) {
	event := new(EthERC20Something)
	if err := _EthERC20.contract.UnpackLog(event, "Something", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the EthERC20 contract.
type EthERC20TransferIterator struct {
	Event *EthERC20Transfer // Event containing the contract specifics and raw log

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
func (it *EthERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthERC20Transfer)
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
		it.Event = new(EthERC20Transfer)
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
func (it *EthERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthERC20Transfer represents a Transfer event raised by the EthERC20 contract.
type EthERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_EthERC20 *EthERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EthERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EthERC20TransferIterator{contract: _EthERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_EthERC20 *EthERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EthERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _EthERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthERC20Transfer)
				if err := _EthERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_EthERC20 *EthERC20Filterer) ParseTransfer(log types.Log) (*EthERC20Transfer, error) {
	event := new(EthERC20Transfer)
	if err := _EthERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
