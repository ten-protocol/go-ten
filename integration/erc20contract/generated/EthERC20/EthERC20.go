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
	Bin: "0x60806040523480156200001157600080fd5b50604051620025a0380380620025a08339818101604052810190620000379190620006e4565b848481600390816200004a9190620009eb565b5080600490816200005c9190620009eb565b50505081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000f33384620000fe60201b60201c565b505050505062000eae565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000170576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620001679062000b33565b60405180910390fd5b62000184600083836200027660201b60201c565b806002600082825462000198919062000b84565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254620001ef919062000b84565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405162000256919062000bd0565b60405180910390a36200027260008383620004ac60201b60201c565b5050565b600073ffffffffffffffffffffffffffffffffffffffff16600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160315620004a757600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620004a657600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115620003c657620003c562000bed565b5b85604051602001620003d9919062000c86565b60405160208183030381529060405260006040518563ffffffff1660e01b81526004016200040b949392919062000d71565b6020604051808303816000875af11580156200042b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000451919062000e0a565b905060018167ffffffffffffffff1614620004a3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200049a9062000e8c565b60405180910390fd5b50505b5b505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200051a82620004cf565b810181811067ffffffffffffffff821117156200053c576200053b620004e0565b5b80604052505050565b600062000551620004b1565b90506200055f82826200050f565b919050565b600067ffffffffffffffff821115620005825762000581620004e0565b5b6200058d82620004cf565b9050602081019050919050565b60005b83811015620005ba5780820151818401526020810190506200059d565b60008484015250505050565b6000620005dd620005d78462000564565b62000545565b905082815260208101848484011115620005fc57620005fb620004ca565b5b620006098482856200059a565b509392505050565b600082601f830112620006295762000628620004c5565b5b81516200063b848260208601620005c6565b91505092915050565b6000819050919050565b620006598162000644565b81146200066557600080fd5b50565b60008151905062000679816200064e565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620006ac826200067f565b9050919050565b620006be816200069f565b8114620006ca57600080fd5b50565b600081519050620006de81620006b3565b92915050565b600080600080600060a08688031215620007035762000702620004bb565b5b600086015167ffffffffffffffff811115620007245762000723620004c0565b5b620007328882890162000611565b955050602086015167ffffffffffffffff811115620007565762000755620004c0565b5b620007648882890162000611565b9450506040620007778882890162000668565b93505060606200078a88828901620006cd565b92505060806200079d88828901620006cd565b9150509295509295909350565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620007fd57607f821691505b602082108103620008135762000812620007b5565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026200087d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826200083e565b6200088986836200083e565b95508019841693508086168417925050509392505050565b6000819050919050565b6000620008cc620008c6620008c08462000644565b620008a1565b62000644565b9050919050565b6000819050919050565b620008e883620008ab565b62000900620008f782620008d3565b8484546200084b565b825550505050565b600090565b6200091762000908565b62000924818484620008dd565b505050565b5b818110156200094c57620009406000826200090d565b6001810190506200092a565b5050565b601f8211156200099b57620009658162000819565b62000970846200082e565b8101602085101562000980578190505b620009986200098f856200082e565b83018262000929565b50505b505050565b600082821c905092915050565b6000620009c060001984600802620009a0565b1980831691505092915050565b6000620009db8383620009ad565b9150826002028217905092915050565b620009f682620007aa565b67ffffffffffffffff81111562000a125762000a11620004e0565b5b62000a1e8254620007e4565b62000a2b82828562000950565b600060209050601f83116001811462000a63576000841562000a4e578287015190505b62000a5a8582620009cd565b86555062000aca565b601f19841662000a738662000819565b60005b8281101562000a9d5784890151825560018201915060208501945060208101905062000a76565b8683101562000abd578489015162000ab9601f891682620009ad565b8355505b6001600288020188555050505b505050505050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000b1b601f8362000ad2565b915062000b288262000ae3565b602082019050919050565b6000602082019050818103600083015262000b4e8162000b0c565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600062000b918262000644565b915062000b9e8362000644565b925082820190508082111562000bb95762000bb862000b55565b5b92915050565b62000bca8162000644565b82525050565b600060208201905062000be7600083018462000bbf565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b62000c27816200069f565b82525050565b62000c388162000644565b82525050565b60608201600082015162000c56600085018262000c1c565b50602082015162000c6b602085018262000c1c565b50604082015162000c80604085018262000c2d565b50505050565b600060608201905062000c9d600083018462000c3e565b92915050565b600063ffffffff82169050919050565b62000cbe8162000ca3565b82525050565b600081519050919050565b600082825260208201905092915050565b600062000ced8262000cc4565b62000cf9818562000ccf565b935062000d0b8185602086016200059a565b62000d1681620004cf565b840191505092915050565b6000819050919050565b600060ff82169050919050565b600062000d5962000d5362000d4d8462000d21565b620008a1565b62000d2b565b9050919050565b62000d6b8162000d38565b82525050565b600060808201905062000d88600083018762000cb3565b62000d97602083018662000cb3565b818103604083015262000dab818562000ce0565b905062000dbc606083018462000d60565b95945050505050565b600067ffffffffffffffff82169050919050565b62000de48162000dc5565b811462000df057600080fd5b50565b60008151905062000e048162000dd9565b92915050565b60006020828403121562000e235762000e22620004bb565b5b600062000e338482850162000df3565b91505092915050565b7f53616e69747920636865636b206661696c000000000000000000000000000000600082015250565b600062000e7460118362000ad2565b915062000e818262000e3c565b602082019050919050565b6000602082019050818103600083015262000ea78162000e65565b9050919050565b6116e28062000ebe6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610d37565b60405180910390f35b6100e660048036038101906100e19190610df2565b610308565b6040516100f39190610e4d565b60405180910390f35b61010461032b565b6040516101119190610e77565b60405180910390f35b610134600480360381019061012f9190610e92565b610335565b6040516101419190610e4d565b60405180910390f35b610152610364565b60405161015f9190610f01565b60405180910390f35b610182600480360381019061017d9190610df2565b61036d565b60405161018f9190610e4d565b60405180910390f35b6101b260048036038101906101ad9190610f1c565b6103a4565b6040516101bf9190610e77565b60405180910390f35b6101d06103ec565b6040516101dd9190610d37565b60405180910390f35b61020060048036038101906101fb9190610df2565b61047e565b60405161020d9190610e4d565b60405180910390f35b610230600480360381019061022b9190610df2565b6104f5565b60405161023d9190610e4d565b60405180910390f35b610260600480360381019061025b9190610f49565b610518565b60405161026d9190610e77565b60405180910390f35b60606003805461028590610fb8565b80601f01602080910402602001604051908101604052809291908181526020018280546102b190610fb8565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b60008061031361059f565b90506103208185856105a7565b600191505092915050565b6000600254905090565b60008061034061059f565b905061034d858285610770565b6103588585856107fc565b60019150509392505050565b60006012905090565b60008061037861059f565b905061039981858561038a8589610518565b6103949190611018565b6105a7565b600191505092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6060600480546103fb90610fb8565b80601f016020809104026020016040519081016040528092919081815260200182805461042790610fb8565b80156104745780601f1061044957610100808354040283529160200191610474565b820191906000526020600020905b81548152906001019060200180831161045757829003601f168201915b5050505050905090565b60008061048961059f565b905060006104978286610518565b9050838110156104dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d3906110be565b60405180910390fd5b6104e982868684036105a7565b60019250505092915050565b60008061050061059f565b905061050d8185856107fc565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610616576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161060d90611150565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610685576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067c906111e2565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516107639190610e77565b60405180910390a3505050565b600061077c8484610518565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107f657818110156107e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107df9061124e565b60405180910390fd5b6107f584848484036105a7565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361086b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610862906112e0565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108d190611372565b60405180910390fd5b6108e5838383610a7b565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561096b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161096290611404565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109fe9190611018565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610a629190610e77565b60405180910390a3610a75848484610ca2565b50505050565b600073ffffffffffffffffffffffffffffffffffffffff16600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160315610c9d57600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610c9c57600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115610bc657610bc5611424565b5b85604051602001610bd791906114b3565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401610c079493929190611587565b6020604051808303816000875af1158015610c26573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c4a9190611613565b905060018167ffffffffffffffff1614610c99576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c909061168c565b60405180910390fd5b50505b5b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610ce1578082015181840152602081019050610cc6565b60008484015250505050565b6000601f19601f8301169050919050565b6000610d0982610ca7565b610d138185610cb2565b9350610d23818560208601610cc3565b610d2c81610ced565b840191505092915050565b60006020820190508181036000830152610d518184610cfe565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610d8982610d5e565b9050919050565b610d9981610d7e565b8114610da457600080fd5b50565b600081359050610db681610d90565b92915050565b6000819050919050565b610dcf81610dbc565b8114610dda57600080fd5b50565b600081359050610dec81610dc6565b92915050565b60008060408385031215610e0957610e08610d59565b5b6000610e1785828601610da7565b9250506020610e2885828601610ddd565b9150509250929050565b60008115159050919050565b610e4781610e32565b82525050565b6000602082019050610e626000830184610e3e565b92915050565b610e7181610dbc565b82525050565b6000602082019050610e8c6000830184610e68565b92915050565b600080600060608486031215610eab57610eaa610d59565b5b6000610eb986828701610da7565b9350506020610eca86828701610da7565b9250506040610edb86828701610ddd565b9150509250925092565b600060ff82169050919050565b610efb81610ee5565b82525050565b6000602082019050610f166000830184610ef2565b92915050565b600060208284031215610f3257610f31610d59565b5b6000610f4084828501610da7565b91505092915050565b60008060408385031215610f6057610f5f610d59565b5b6000610f6e85828601610da7565b9250506020610f7f85828601610da7565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610fd057607f821691505b602082108103610fe357610fe2610f89565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061102382610dbc565b915061102e83610dbc565b925082820190508082111561104657611045610fe9565b5b92915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b60006110a8602583610cb2565b91506110b38261104c565b604082019050919050565b600060208201905081810360008301526110d78161109b565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b600061113a602483610cb2565b9150611145826110de565b604082019050919050565b600060208201905081810360008301526111698161112d565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b60006111cc602283610cb2565b91506111d782611170565b604082019050919050565b600060208201905081810360008301526111fb816111bf565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b6000611238601d83610cb2565b915061124382611202565b602082019050919050565b600060208201905081810360008301526112678161122b565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b60006112ca602583610cb2565b91506112d58261126e565b604082019050919050565b600060208201905081810360008301526112f9816112bd565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b600061135c602383610cb2565b915061136782611300565b604082019050919050565b6000602082019050818103600083015261138b8161134f565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b60006113ee602683610cb2565b91506113f982611392565b604082019050919050565b6000602082019050818103600083015261141d816113e1565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61145c81610d7e565b82525050565b61146b81610dbc565b82525050565b6060820160008201516114876000850182611453565b50602082015161149a6020850182611453565b5060408201516114ad6040850182611462565b50505050565b60006060820190506114c86000830184611471565b92915050565b600063ffffffff82169050919050565b6114e7816114ce565b82525050565b600081519050919050565b600082825260208201905092915050565b6000611514826114ed565b61151e81856114f8565b935061152e818560208601610cc3565b61153781610ced565b840191505092915050565b6000819050919050565b6000819050919050565b600061157161156c61156784611542565b61154c565b610ee5565b9050919050565b61158181611556565b82525050565b600060808201905061159c60008301876114de565b6115a960208301866114de565b81810360408301526115bb8185611509565b90506115ca6060830184611578565b95945050505050565b600067ffffffffffffffff82169050919050565b6115f0816115d3565b81146115fb57600080fd5b50565b60008151905061160d816115e7565b92915050565b60006020828403121561162957611628610d59565b5b6000611637848285016115fe565b91505092915050565b7f53616e69747920636865636b206661696c000000000000000000000000000000600082015250565b6000611676601183610cb2565b915061168182611640565b602082019050919050565b600060208201905081810360008301526116a581611669565b905091905056fea2646970667358221220988104b8f9e3dbeb0a24b87f9557812688c1e2f1ee33845b9d3e0e5a819a0a5e64736f6c63430008110033",
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
