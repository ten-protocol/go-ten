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
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"l1MessageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"managementContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"hm\",\"type\":\"string\"}],\"name\":\"Something\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620024ef380380620024ef83398181016040528101906200003791906200068b565b848481600390816200004a919062000992565b5080600490816200005c919062000992565b50505081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000b23384620000fe60201b60201c565b80600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050505062000e55565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000170576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620001679062000ada565b60405180910390fd5b62000184600083836200027660201b60201c565b806002600082825462000198919062000b2b565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254620001ef919062000b2b565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405162000256919062000b77565b60405180910390a362000272600083836200045360201b60201c565b5050565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036200044e57600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa436001808111156200036e576200036d62000b94565b5b8560405160200162000381919062000c2d565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401620003b3949392919062000d18565b6020604051808303816000875af1158015620003d3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620003f9919062000db1565b905060018167ffffffffffffffff16146200044b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620004429062000e33565b60405180910390fd5b50505b505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620004c18262000476565b810181811067ffffffffffffffff82111715620004e357620004e262000487565b5b80604052505050565b6000620004f862000458565b9050620005068282620004b6565b919050565b600067ffffffffffffffff82111562000529576200052862000487565b5b620005348262000476565b9050602081019050919050565b60005b838110156200056157808201518184015260208101905062000544565b60008484015250505050565b6000620005846200057e846200050b565b620004ec565b905082815260208101848484011115620005a357620005a262000471565b5b620005b084828562000541565b509392505050565b600082601f830112620005d057620005cf6200046c565b5b8151620005e28482602086016200056d565b91505092915050565b6000819050919050565b6200060081620005eb565b81146200060c57600080fd5b50565b6000815190506200062081620005f5565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620006538262000626565b9050919050565b620006658162000646565b81146200067157600080fd5b50565b60008151905062000685816200065a565b92915050565b600080600080600060a08688031215620006aa57620006a962000462565b5b600086015167ffffffffffffffff811115620006cb57620006ca62000467565b5b620006d988828901620005b8565b955050602086015167ffffffffffffffff811115620006fd57620006fc62000467565b5b6200070b88828901620005b8565b94505060406200071e888289016200060f565b9350506060620007318882890162000674565b9250506080620007448882890162000674565b9150509295509295909350565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620007a457607f821691505b602082108103620007ba57620007b96200075c565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620008247fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620007e5565b620008308683620007e5565b95508019841693508086168417925050509392505050565b6000819050919050565b6000620008736200086d6200086784620005eb565b62000848565b620005eb565b9050919050565b6000819050919050565b6200088f8362000852565b620008a76200089e826200087a565b848454620007f2565b825550505050565b600090565b620008be620008af565b620008cb81848462000884565b505050565b5b81811015620008f357620008e7600082620008b4565b600181019050620008d1565b5050565b601f82111562000942576200090c81620007c0565b6200091784620007d5565b8101602085101562000927578190505b6200093f6200093685620007d5565b830182620008d0565b50505b505050565b600082821c905092915050565b6000620009676000198460080262000947565b1980831691505092915050565b600062000982838362000954565b9150826002028217905092915050565b6200099d8262000751565b67ffffffffffffffff811115620009b957620009b862000487565b5b620009c582546200078b565b620009d2828285620008f7565b600060209050601f83116001811462000a0a5760008415620009f5578287015190505b62000a01858262000974565b86555062000a71565b601f19841662000a1a86620007c0565b60005b8281101562000a445784890151825560018201915060208501945060208101905062000a1d565b8683101562000a64578489015162000a60601f89168262000954565b8355505b6001600288020188555050505b505050505050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000ac2601f8362000a79565b915062000acf8262000a8a565b602082019050919050565b6000602082019050818103600083015262000af58162000ab3565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600062000b3882620005eb565b915062000b4583620005eb565b925082820190508082111562000b605762000b5f62000afc565b5b92915050565b62000b7181620005eb565b82525050565b600060208201905062000b8e600083018462000b66565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b62000bce8162000646565b82525050565b62000bdf81620005eb565b82525050565b60608201600082015162000bfd600085018262000bc3565b50602082015162000c12602085018262000bc3565b50604082015162000c27604085018262000bd4565b50505050565b600060608201905062000c44600083018462000be5565b92915050565b600063ffffffff82169050919050565b62000c658162000c4a565b82525050565b600081519050919050565b600082825260208201905092915050565b600062000c948262000c6b565b62000ca0818562000c76565b935062000cb281856020860162000541565b62000cbd8162000476565b840191505092915050565b6000819050919050565b600060ff82169050919050565b600062000d0062000cfa62000cf48462000cc8565b62000848565b62000cd2565b9050919050565b62000d128162000cdf565b82525050565b600060808201905062000d2f600083018762000c5a565b62000d3e602083018662000c5a565b818103604083015262000d52818562000c87565b905062000d63606083018462000d07565b95945050505050565b600067ffffffffffffffff82169050919050565b62000d8b8162000d6c565b811462000d9757600080fd5b50565b60008151905062000dab8162000d80565b92915050565b60006020828403121562000dca5762000dc962000462565b5b600062000dda8482850162000d9a565b91505092915050565b7f53616e69747920636865636b206661696c000000000000000000000000000000600082015250565b600062000e1b60118362000a79565b915062000e288262000de3565b602082019050919050565b6000602082019050818103600083015262000e4e8162000e0c565b9050919050565b61168a8062000e656000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610cdf565b60405180910390f35b6100e660048036038101906100e19190610d9a565b610308565b6040516100f39190610df5565b60405180910390f35b61010461032b565b6040516101119190610e1f565b60405180910390f35b610134600480360381019061012f9190610e3a565b610335565b6040516101419190610df5565b60405180910390f35b610152610364565b60405161015f9190610ea9565b60405180910390f35b610182600480360381019061017d9190610d9a565b61036d565b60405161018f9190610df5565b60405180910390f35b6101b260048036038101906101ad9190610ec4565b6103a4565b6040516101bf9190610e1f565b60405180910390f35b6101d06103ec565b6040516101dd9190610cdf565b60405180910390f35b61020060048036038101906101fb9190610d9a565b61047e565b60405161020d9190610df5565b60405180910390f35b610230600480360381019061022b9190610d9a565b6104f5565b60405161023d9190610df5565b60405180910390f35b610260600480360381019061025b9190610ef1565b610518565b60405161026d9190610e1f565b60405180910390f35b60606003805461028590610f60565b80601f01602080910402602001604051908101604052809291908181526020018280546102b190610f60565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b60008061031361059f565b90506103208185856105a7565b600191505092915050565b6000600254905090565b60008061034061059f565b905061034d858285610770565b6103588585856107fc565b60019150509392505050565b60006012905090565b60008061037861059f565b905061039981858561038a8589610518565b6103949190610fc0565b6105a7565b600191505092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6060600480546103fb90610f60565b80601f016020809104026020016040519081016040528092919081815260200182805461042790610f60565b80156104745780601f1061044957610100808354040283529160200191610474565b820191906000526020600020905b81548152906001019060200180831161045757829003601f168201915b5050505050905090565b60008061048961059f565b905060006104978286610518565b9050838110156104dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d390611066565b60405180910390fd5b6104e982868684036105a7565b60019250505092915050565b60008061050061059f565b905061050d8185856107fc565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610616576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161060d906110f8565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610685576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067c9061118a565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516107639190610e1f565b60405180910390a3505050565b600061077c8484610518565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107f657818110156107e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107df906111f6565b60405180910390fd5b6107f584848484036105a7565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361086b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161086290611288565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108d19061131a565b60405180910390fd5b6108e5838383610a7b565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561096b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610962906113ac565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109fe9190610fc0565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610a629190610e1f565b60405180910390a3610a75848484610c4a565b50505050565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610c4557600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115610b6f57610b6e6113cc565b5b85604051602001610b80919061145b565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401610bb0949392919061152f565b6020604051808303816000875af1158015610bcf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bf391906115bb565b905060018167ffffffffffffffff1614610c42576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c3990611634565b60405180910390fd5b50505b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610c89578082015181840152602081019050610c6e565b60008484015250505050565b6000601f19601f8301169050919050565b6000610cb182610c4f565b610cbb8185610c5a565b9350610ccb818560208601610c6b565b610cd481610c95565b840191505092915050565b60006020820190508181036000830152610cf98184610ca6565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610d3182610d06565b9050919050565b610d4181610d26565b8114610d4c57600080fd5b50565b600081359050610d5e81610d38565b92915050565b6000819050919050565b610d7781610d64565b8114610d8257600080fd5b50565b600081359050610d9481610d6e565b92915050565b60008060408385031215610db157610db0610d01565b5b6000610dbf85828601610d4f565b9250506020610dd085828601610d85565b9150509250929050565b60008115159050919050565b610def81610dda565b82525050565b6000602082019050610e0a6000830184610de6565b92915050565b610e1981610d64565b82525050565b6000602082019050610e346000830184610e10565b92915050565b600080600060608486031215610e5357610e52610d01565b5b6000610e6186828701610d4f565b9350506020610e7286828701610d4f565b9250506040610e8386828701610d85565b9150509250925092565b600060ff82169050919050565b610ea381610e8d565b82525050565b6000602082019050610ebe6000830184610e9a565b92915050565b600060208284031215610eda57610ed9610d01565b5b6000610ee884828501610d4f565b91505092915050565b60008060408385031215610f0857610f07610d01565b5b6000610f1685828601610d4f565b9250506020610f2785828601610d4f565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610f7857607f821691505b602082108103610f8b57610f8a610f31565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610fcb82610d64565b9150610fd683610d64565b9250828201905080821115610fee57610fed610f91565b5b92915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b6000611050602583610c5a565b915061105b82610ff4565b604082019050919050565b6000602082019050818103600083015261107f81611043565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006110e2602483610c5a565b91506110ed82611086565b604082019050919050565b60006020820190508181036000830152611111816110d5565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b6000611174602283610c5a565b915061117f82611118565b604082019050919050565b600060208201905081810360008301526111a381611167565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b60006111e0601d83610c5a565b91506111eb826111aa565b602082019050919050565b6000602082019050818103600083015261120f816111d3565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b6000611272602583610c5a565b915061127d82611216565b604082019050919050565b600060208201905081810360008301526112a181611265565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b6000611304602383610c5a565b915061130f826112a8565b604082019050919050565b60006020820190508181036000830152611333816112f7565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b6000611396602683610c5a565b91506113a18261133a565b604082019050919050565b600060208201905081810360008301526113c581611389565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61140481610d26565b82525050565b61141381610d64565b82525050565b60608201600082015161142f60008501826113fb565b50602082015161144260208501826113fb565b506040820151611455604085018261140a565b50505050565b60006060820190506114706000830184611419565b92915050565b600063ffffffff82169050919050565b61148f81611476565b82525050565b600081519050919050565b600082825260208201905092915050565b60006114bc82611495565b6114c681856114a0565b93506114d6818560208601610c6b565b6114df81610c95565b840191505092915050565b6000819050919050565b6000819050919050565b600061151961151461150f846114ea565b6114f4565b610e8d565b9050919050565b611529816114fe565b82525050565b60006080820190506115446000830187611486565b6115516020830186611486565b818103604083015261156381856114b1565b90506115726060830184611520565b95945050505050565b600067ffffffffffffffff82169050919050565b6115988161157b565b81146115a357600080fd5b50565b6000815190506115b58161158f565b92915050565b6000602082840312156115d1576115d0610d01565b5b60006115df848285016115a6565b91505092915050565b7f53616e69747920636865636b206661696c000000000000000000000000000000600082015250565b600061161e601183610c5a565b9150611629826115e8565b602082019050919050565b6000602082019050818103600083015261164d81611611565b905091905056fea26469706673582212208602c487139c77b32992ab54d506dc2e2737fbb0fd02974fe529345a4ede20ee64736f6c63430008110033",
}

// EthERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use EthERC20MetaData.ABI instead.
var EthERC20ABI = EthERC20MetaData.ABI

// EthERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EthERC20MetaData.Bin instead.
var EthERC20Bin = EthERC20MetaData.Bin

// DeployEthERC20 deploys a new Ethereum contract, binding an instance of EthERC20 to it.
func DeployEthERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int, l1MessageBus common.Address, managementContract common.Address) (common.Address, *types.Transaction, *EthERC20, error) {
	parsed, err := EthERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EthERC20Bin), backend, name, symbol, initialSupply, l1MessageBus, managementContract)
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
