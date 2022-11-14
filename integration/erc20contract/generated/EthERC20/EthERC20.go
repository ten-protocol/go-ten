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
	Bin: "0x60806040523480156200001157600080fd5b506040516200256638038062002566833981810160405281019062000037919062000716565b848481600390816200004a919062000a1d565b5080600490816200005c919062000a1d565b50505081600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa4360006001811115620000f657620000f562000b04565b5b8660405160200162000109919062000b58565b60405160208183030381529060405260006040518563ffffffff1660e01b81526004016200013b949392919062000c43565b6020604051808303816000875af11580156200015b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000181919062000cdc565b50620001943384620001e060201b60201c565b81600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050505062000f22565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000252576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620002499062000d6f565b60405180910390fd5b62000266600083836200035860201b60201c565b80600260008282546200027a919062000dc0565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254620002d1919062000dc0565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405162000338919062000e0c565b60405180910390a36200035460008383620004de60201b60201c565b5050565b600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115620003fa57620003f962000b04565b5b856040516020016200040d919062000e93565b60405160208183030381529060405260006040518563ffffffff1660e01b81526004016200043f949392919062000c43565b6020604051808303816000875af11580156200045f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000485919062000cdc565b905060018167ffffffffffffffff1614620004d7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620004ce9062000f00565b60405180910390fd5b5050505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200054c8262000501565b810181811067ffffffffffffffff821117156200056e576200056d62000512565b5b80604052505050565b600062000583620004e3565b905062000591828262000541565b919050565b600067ffffffffffffffff821115620005b457620005b362000512565b5b620005bf8262000501565b9050602081019050919050565b60005b83811015620005ec578082015181840152602081019050620005cf565b60008484015250505050565b60006200060f620006098462000596565b62000577565b9050828152602081018484840111156200062e576200062d620004fc565b5b6200063b848285620005cc565b509392505050565b600082601f8301126200065b576200065a620004f7565b5b81516200066d848260208601620005f8565b91505092915050565b6000819050919050565b6200068b8162000676565b81146200069757600080fd5b50565b600081519050620006ab8162000680565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620006de82620006b1565b9050919050565b620006f081620006d1565b8114620006fc57600080fd5b50565b6000815190506200071081620006e5565b92915050565b600080600080600060a08688031215620007355762000734620004ed565b5b600086015167ffffffffffffffff811115620007565762000755620004f2565b5b620007648882890162000643565b955050602086015167ffffffffffffffff811115620007885762000787620004f2565b5b620007968882890162000643565b9450506040620007a9888289016200069a565b9350506060620007bc88828901620006ff565b9250506080620007cf88828901620006ff565b9150509295509295909350565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200082f57607f821691505b602082108103620008455762000844620007e7565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620008af7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000870565b620008bb868362000870565b95508019841693508086168417925050509392505050565b6000819050919050565b6000620008fe620008f8620008f28462000676565b620008d3565b62000676565b9050919050565b6000819050919050565b6200091a83620008dd565b62000932620009298262000905565b8484546200087d565b825550505050565b600090565b620009496200093a565b620009568184846200090f565b505050565b5b818110156200097e57620009726000826200093f565b6001810190506200095c565b5050565b601f821115620009cd5762000997816200084b565b620009a28462000860565b81016020851015620009b2578190505b620009ca620009c18562000860565b8301826200095b565b50505b505050565b600082821c905092915050565b6000620009f260001984600802620009d2565b1980831691505092915050565b600062000a0d8383620009df565b9150826002028217905092915050565b62000a2882620007dc565b67ffffffffffffffff81111562000a445762000a4362000512565b5b62000a50825462000816565b62000a5d82828562000982565b600060209050601f83116001811462000a95576000841562000a80578287015190505b62000a8c8582620009ff565b86555062000afc565b601f19841662000aa5866200084b565b60005b8281101562000acf5784890151825560018201915060208501945060208101905062000aa8565b8683101562000aef578489015162000aeb601f891682620009df565b8355505b6001600288020188555050505b505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6000819050919050565b62000b5262000b4c8262000676565b62000b33565b82525050565b600062000b66828462000b3d565b60208201915081905092915050565b600063ffffffff82169050919050565b62000b908162000b75565b82525050565b600081519050919050565b600082825260208201905092915050565b600062000bbf8262000b96565b62000bcb818562000ba1565b935062000bdd818560208601620005cc565b62000be88162000501565b840191505092915050565b6000819050919050565b600060ff82169050919050565b600062000c2b62000c2562000c1f8462000bf3565b620008d3565b62000bfd565b9050919050565b62000c3d8162000c0a565b82525050565b600060808201905062000c5a600083018762000b85565b62000c69602083018662000b85565b818103604083015262000c7d818562000bb2565b905062000c8e606083018462000c32565b95945050505050565b600067ffffffffffffffff82169050919050565b62000cb68162000c97565b811462000cc257600080fd5b50565b60008151905062000cd68162000cab565b92915050565b60006020828403121562000cf55762000cf4620004ed565b5b600062000d058482850162000cc5565b91505092915050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000d57601f8362000d0e565b915062000d648262000d1f565b602082019050919050565b6000602082019050818103600083015262000d8a8162000d48565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600062000dcd8262000676565b915062000dda8362000676565b925082820190508082111562000df55762000df462000d91565b5b92915050565b62000e068162000676565b82525050565b600060208201905062000e23600083018462000dfb565b92915050565b62000e3481620006d1565b82525050565b62000e458162000676565b82525050565b60608201600082015162000e63600085018262000e29565b50602082015162000e78602085018262000e29565b50604082015162000e8d604085018262000e3a565b50505050565b600060608201905062000eaa600083018462000e4b565b92915050565b7f53616e69747920636865636b206661696c000000000000000000000000000000600082015250565b600062000ee860118362000d0e565b915062000ef58262000eb0565b602082019050919050565b6000602082019050818103600083015262000f1b8162000ed9565b9050919050565b6116348062000f326000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610c89565b60405180910390f35b6100e660048036038101906100e19190610d44565b610308565b6040516100f39190610d9f565b60405180910390f35b61010461032b565b6040516101119190610dc9565b60405180910390f35b610134600480360381019061012f9190610de4565b610335565b6040516101419190610d9f565b60405180910390f35b610152610364565b60405161015f9190610e53565b60405180910390f35b610182600480360381019061017d9190610d44565b61036d565b60405161018f9190610d9f565b60405180910390f35b6101b260048036038101906101ad9190610e6e565b6103a4565b6040516101bf9190610dc9565b60405180910390f35b6101d06103ec565b6040516101dd9190610c89565b60405180910390f35b61020060048036038101906101fb9190610d44565b61047e565b60405161020d9190610d9f565b60405180910390f35b610230600480360381019061022b9190610d44565b6104f5565b60405161023d9190610d9f565b60405180910390f35b610260600480360381019061025b9190610e9b565b610518565b60405161026d9190610dc9565b60405180910390f35b60606003805461028590610f0a565b80601f01602080910402602001604051908101604052809291908181526020018280546102b190610f0a565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b60008061031361059f565b90506103208185856105a7565b600191505092915050565b6000600254905090565b60008061034061059f565b905061034d858285610770565b6103588585856107fc565b60019150509392505050565b60006012905090565b60008061037861059f565b905061039981858561038a8589610518565b6103949190610f6a565b6105a7565b600191505092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6060600480546103fb90610f0a565b80601f016020809104026020016040519081016040528092919081815260200182805461042790610f0a565b80156104745780601f1061044957610100808354040283529160200191610474565b820191906000526020600020905b81548152906001019060200180831161045757829003601f168201915b5050505050905090565b60008061048961059f565b905060006104978286610518565b9050838110156104dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d390611010565b60405180910390fd5b6104e982868684036105a7565b60019250505092915050565b60008061050061059f565b905061050d8185856107fc565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610616576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161060d906110a2565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610685576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067c90611134565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516107639190610dc9565b60405180910390a3505050565b600061077c8484610518565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107f657818110156107e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107df906111a0565b60405180910390fd5b6107f584848484036105a7565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361086b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161086290611232565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108d1906112c4565b60405180910390fd5b6108e5838383610a7b565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561096b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161096290611356565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109fe9190610f6a565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610a629190610dc9565b60405180910390a3610a75848484610bf4565b50505050565b600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115610b1a57610b19611376565b5b85604051602001610b2b9190611405565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401610b5b94939291906114d9565b6020604051808303816000875af1158015610b7a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b9e9190611565565b905060018167ffffffffffffffff1614610bed576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610be4906115de565b60405180910390fd5b5050505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610c33578082015181840152602081019050610c18565b60008484015250505050565b6000601f19601f8301169050919050565b6000610c5b82610bf9565b610c658185610c04565b9350610c75818560208601610c15565b610c7e81610c3f565b840191505092915050565b60006020820190508181036000830152610ca38184610c50565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610cdb82610cb0565b9050919050565b610ceb81610cd0565b8114610cf657600080fd5b50565b600081359050610d0881610ce2565b92915050565b6000819050919050565b610d2181610d0e565b8114610d2c57600080fd5b50565b600081359050610d3e81610d18565b92915050565b60008060408385031215610d5b57610d5a610cab565b5b6000610d6985828601610cf9565b9250506020610d7a85828601610d2f565b9150509250929050565b60008115159050919050565b610d9981610d84565b82525050565b6000602082019050610db46000830184610d90565b92915050565b610dc381610d0e565b82525050565b6000602082019050610dde6000830184610dba565b92915050565b600080600060608486031215610dfd57610dfc610cab565b5b6000610e0b86828701610cf9565b9350506020610e1c86828701610cf9565b9250506040610e2d86828701610d2f565b9150509250925092565b600060ff82169050919050565b610e4d81610e37565b82525050565b6000602082019050610e686000830184610e44565b92915050565b600060208284031215610e8457610e83610cab565b5b6000610e9284828501610cf9565b91505092915050565b60008060408385031215610eb257610eb1610cab565b5b6000610ec085828601610cf9565b9250506020610ed185828601610cf9565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610f2257607f821691505b602082108103610f3557610f34610edb565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610f7582610d0e565b9150610f8083610d0e565b9250828201905080821115610f9857610f97610f3b565b5b92915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b6000610ffa602583610c04565b915061100582610f9e565b604082019050919050565b6000602082019050818103600083015261102981610fed565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b600061108c602483610c04565b915061109782611030565b604082019050919050565b600060208201905081810360008301526110bb8161107f565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b600061111e602283610c04565b9150611129826110c2565b604082019050919050565b6000602082019050818103600083015261114d81611111565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b600061118a601d83610c04565b915061119582611154565b602082019050919050565b600060208201905081810360008301526111b98161117d565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b600061121c602583610c04565b9150611227826111c0565b604082019050919050565b6000602082019050818103600083015261124b8161120f565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b60006112ae602383610c04565b91506112b982611252565b604082019050919050565b600060208201905081810360008301526112dd816112a1565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b6000611340602683610c04565b915061134b826112e4565b604082019050919050565b6000602082019050818103600083015261136f81611333565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6113ae81610cd0565b82525050565b6113bd81610d0e565b82525050565b6060820160008201516113d960008501826113a5565b5060208201516113ec60208501826113a5565b5060408201516113ff60408501826113b4565b50505050565b600060608201905061141a60008301846113c3565b92915050565b600063ffffffff82169050919050565b61143981611420565b82525050565b600081519050919050565b600082825260208201905092915050565b60006114668261143f565b611470818561144a565b9350611480818560208601610c15565b61148981610c3f565b840191505092915050565b6000819050919050565b6000819050919050565b60006114c36114be6114b984611494565b61149e565b610e37565b9050919050565b6114d3816114a8565b82525050565b60006080820190506114ee6000830187611430565b6114fb6020830186611430565b818103604083015261150d818561145b565b905061151c60608301846114ca565b95945050505050565b600067ffffffffffffffff82169050919050565b61154281611525565b811461154d57600080fd5b50565b60008151905061155f81611539565b92915050565b60006020828403121561157b5761157a610cab565b5b600061158984828501611550565b91505092915050565b7f53616e69747920636865636b206661696c000000000000000000000000000000600082015250565b60006115c8601183610c04565b91506115d382611592565b602082019050919050565b600060208201905081810360008301526115f7816115bb565b905091905056fea2646970667358221220b01926836177607aef305b3a98a2a16370ee338e6145a6ec1947d1c54bcdadb264736f6c63430008110033",
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
