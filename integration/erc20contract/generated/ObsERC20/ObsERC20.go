// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObsERC20

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

// ObsERC20MetaData contains all meta data concerning the ObsERC20 contract.
var ObsERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001b5838038062001b588339818101604052810190620000379190620003d4565b828281600390816200004a9190620006af565b5080600490816200005c9190620006af565b5050506200007133826200007a60201b60201c565b505050620008d3565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620000ec576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000e390620007f7565b60405180910390fd5b6200010060008383620001f260201b60201c565b806002600082825462000114919062000848565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546200016b919062000848565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620001d29190620008b6565b60405180910390a3620001ee60008383620001f760201b60201c565b5050565b505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b62000265826200021a565b810181811067ffffffffffffffff821117156200028757620002866200022b565b5b80604052505050565b60006200029c620001fc565b9050620002aa82826200025a565b919050565b600067ffffffffffffffff821115620002cd57620002cc6200022b565b5b620002d8826200021a565b9050602081019050919050565b60005b8381101562000305578082015181840152602081019050620002e8565b8381111562000315576000848401525b50505050565b6000620003326200032c84620002af565b62000290565b90508281526020810184848401111562000351576200035062000215565b5b6200035e848285620002e5565b509392505050565b600082601f8301126200037e576200037d62000210565b5b8151620003908482602086016200031b565b91505092915050565b6000819050919050565b620003ae8162000399565b8114620003ba57600080fd5b50565b600081519050620003ce81620003a3565b92915050565b600080600060608486031215620003f057620003ef62000206565b5b600084015167ffffffffffffffff8111156200041157620004106200020b565b5b6200041f8682870162000366565b935050602084015167ffffffffffffffff8111156200044357620004426200020b565b5b620004518682870162000366565b92505060406200046486828701620003bd565b9150509250925092565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620004c157607f821691505b602082108103620004d757620004d662000479565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620005417fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000502565b6200054d868362000502565b95508019841693508086168417925050509392505050565b6000819050919050565b6000620005906200058a620005848462000399565b62000565565b62000399565b9050919050565b6000819050919050565b620005ac836200056f565b620005c4620005bb8262000597565b8484546200050f565b825550505050565b600090565b620005db620005cc565b620005e8818484620005a1565b505050565b5b81811015620006105762000604600082620005d1565b600181019050620005ee565b5050565b601f8211156200065f576200062981620004dd565b6200063484620004f2565b8101602085101562000644578190505b6200065c6200065385620004f2565b830182620005ed565b50505b505050565b600082821c905092915050565b6000620006846000198460080262000664565b1980831691505092915050565b60006200069f838362000671565b9150826002028217905092915050565b620006ba826200046e565b67ffffffffffffffff811115620006d657620006d56200022b565b5b620006e28254620004a8565b620006ef82828562000614565b600060209050601f83116001811462000727576000841562000712578287015190505b6200071e858262000691565b8655506200078e565b601f1984166200073786620004dd565b60005b8281101562000761578489015182556001820191506020850194506020810190506200073a565b868310156200078157848901516200077d601f89168262000671565b8355505b6001600288020188555050505b505050505050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b6000620007df601f8362000796565b9150620007ec82620007a7565b602082019050919050565b600060208201905081810360008301526200081281620007d0565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620008558262000399565b9150620008628362000399565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156200089a576200089962000819565b5b828201905092915050565b620008b08162000399565b82525050565b6000602082019050620008cd6000830184620008a5565b92915050565b61127580620008e36000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610b30565b60405180910390f35b6100e660048036038101906100e19190610beb565b610308565b6040516100f39190610c46565b60405180910390f35b61010461032b565b6040516101119190610c70565b60405180910390f35b610134600480360381019061012f9190610c8b565b610335565b6040516101419190610c46565b60405180910390f35b610152610364565b60405161015f9190610cfa565b60405180910390f35b610182600480360381019061017d9190610beb565b61036d565b60405161018f9190610c46565b60405180910390f35b6101b260048036038101906101ad9190610d15565b6103a4565b6040516101bf9190610c70565b60405180910390f35b6101d06103b6565b6040516101dd9190610b30565b60405180910390f35b61020060048036038101906101fb9190610beb565b610448565b60405161020d9190610c46565b60405180910390f35b610230600480360381019061022b9190610beb565b6104bf565b60405161023d9190610c46565b60405180910390f35b610260600480360381019061025b9190610d42565b6104e2565b60405161026d9190610c70565b60405180910390f35b60606003805461028590610db1565b80601f01602080910402602001604051908101604052809291908181526020018280546102b190610db1565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b600080610313610569565b9050610320818585610571565b600191505092915050565b6000600254905090565b600080610340610569565b905061034d85828561073a565b6103588585856107c6565b60019150509392505050565b60006012905090565b600080610378610569565b905061039981858561038a85896104e2565b6103949190610e11565b610571565b600191505092915050565b60006103af82610a45565b9050919050565b6060600480546103c590610db1565b80601f01602080910402602001604051908101604052809291908181526020018280546103f190610db1565b801561043e5780601f106104135761010080835404028352916020019161043e565b820191906000526020600020905b81548152906001019060200180831161042157829003601f168201915b5050505050905090565b600080610453610569565b9050600061046182866104e2565b9050838110156104a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161049d90610ed9565b60405180910390fd5b6104b38286868403610571565b60019250505092915050565b6000806104ca610569565b90506104d78185856107c6565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036105e0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105d790610f6b565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361064f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064690610ffd565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258360405161072d9190610c70565b60405180910390a3505050565b600061074684846104e2565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146107c057818110156107b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107a990611069565b60405180910390fd5b6107bf8484848403610571565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610835576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082c906110fb565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161089b9061118d565b60405180910390fd5b6108af838383610a8d565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610935576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161092c9061121f565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109c89190610e11565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610a2c9190610c70565b60405180910390a3610a3f848484610a92565b50505050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610ad1578082015181840152602081019050610ab6565b83811115610ae0576000848401525b50505050565b6000601f19601f8301169050919050565b6000610b0282610a97565b610b0c8185610aa2565b9350610b1c818560208601610ab3565b610b2581610ae6565b840191505092915050565b60006020820190508181036000830152610b4a8184610af7565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610b8282610b57565b9050919050565b610b9281610b77565b8114610b9d57600080fd5b50565b600081359050610baf81610b89565b92915050565b6000819050919050565b610bc881610bb5565b8114610bd357600080fd5b50565b600081359050610be581610bbf565b92915050565b60008060408385031215610c0257610c01610b52565b5b6000610c1085828601610ba0565b9250506020610c2185828601610bd6565b9150509250929050565b60008115159050919050565b610c4081610c2b565b82525050565b6000602082019050610c5b6000830184610c37565b92915050565b610c6a81610bb5565b82525050565b6000602082019050610c856000830184610c61565b92915050565b600080600060608486031215610ca457610ca3610b52565b5b6000610cb286828701610ba0565b9350506020610cc386828701610ba0565b9250506040610cd486828701610bd6565b9150509250925092565b600060ff82169050919050565b610cf481610cde565b82525050565b6000602082019050610d0f6000830184610ceb565b92915050565b600060208284031215610d2b57610d2a610b52565b5b6000610d3984828501610ba0565b91505092915050565b60008060408385031215610d5957610d58610b52565b5b6000610d6785828601610ba0565b9250506020610d7885828601610ba0565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610dc957607f821691505b602082108103610ddc57610ddb610d82565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610e1c82610bb5565b9150610e2783610bb5565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115610e5c57610e5b610de2565b5b828201905092915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b6000610ec3602583610aa2565b9150610ece82610e67565b604082019050919050565b60006020820190508181036000830152610ef281610eb6565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b6000610f55602483610aa2565b9150610f6082610ef9565b604082019050919050565b60006020820190508181036000830152610f8481610f48565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b6000610fe7602283610aa2565b9150610ff282610f8b565b604082019050919050565b6000602082019050818103600083015261101681610fda565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b6000611053601d83610aa2565b915061105e8261101d565b602082019050919050565b6000602082019050818103600083015261108281611046565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b60006110e5602583610aa2565b91506110f082611089565b604082019050919050565b60006020820190508181036000830152611114816110d8565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b6000611177602383610aa2565b91506111828261111b565b604082019050919050565b600060208201905081810360008301526111a68161116a565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b6000611209602683610aa2565b9150611214826111ad565b604082019050919050565b60006020820190508181036000830152611238816111fc565b905091905056fea26469706673582212207b67a90fb99649931784a7a72527659444973ac278afa6d3922fe81730be6f1264736f6c634300080f0033",
}

// ObsERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ObsERC20MetaData.ABI instead.
var ObsERC20ABI = ObsERC20MetaData.ABI

// ObsERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObsERC20MetaData.Bin instead.
var ObsERC20Bin = ObsERC20MetaData.Bin

// DeployObsERC20 deploys a new Ethereum contract, binding an instance of ObsERC20 to it.
func DeployObsERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int) (common.Address, *types.Transaction, *ObsERC20, error) {
	parsed, err := ObsERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObsERC20Bin), backend, name, symbol, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObsERC20{ObsERC20Caller: ObsERC20Caller{contract: contract}, ObsERC20Transactor: ObsERC20Transactor{contract: contract}, ObsERC20Filterer: ObsERC20Filterer{contract: contract}}, nil
}

// ObsERC20 is an auto generated Go binding around an Ethereum contract.
type ObsERC20 struct {
	ObsERC20Caller     // Read-only binding to the contract
	ObsERC20Transactor // Write-only binding to the contract
	ObsERC20Filterer   // Log filterer for contract events
}

// ObsERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ObsERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ObsERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObsERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObsERC20Session struct {
	Contract     *ObsERC20         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObsERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObsERC20CallerSession struct {
	Contract *ObsERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ObsERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObsERC20TransactorSession struct {
	Contract     *ObsERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ObsERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ObsERC20Raw struct {
	Contract *ObsERC20 // Generic contract binding to access the raw methods on
}

// ObsERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObsERC20CallerRaw struct {
	Contract *ObsERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ObsERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObsERC20TransactorRaw struct {
	Contract *ObsERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewObsERC20 creates a new instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20(address common.Address, backend bind.ContractBackend) (*ObsERC20, error) {
	contract, err := bindObsERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObsERC20{ObsERC20Caller: ObsERC20Caller{contract: contract}, ObsERC20Transactor: ObsERC20Transactor{contract: contract}, ObsERC20Filterer: ObsERC20Filterer{contract: contract}}, nil
}

// NewObsERC20Caller creates a new read-only instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Caller(address common.Address, caller bind.ContractCaller) (*ObsERC20Caller, error) {
	contract, err := bindObsERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Caller{contract: contract}, nil
}

// NewObsERC20Transactor creates a new write-only instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ObsERC20Transactor, error) {
	contract, err := bindObsERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Transactor{contract: contract}, nil
}

// NewObsERC20Filterer creates a new log filterer instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ObsERC20Filterer, error) {
	contract, err := bindObsERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Filterer{contract: contract}, nil
}

// bindObsERC20 binds a generic wrapper to an already deployed contract.
func bindObsERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ObsERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObsERC20 *ObsERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObsERC20.Contract.ObsERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObsERC20 *ObsERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObsERC20.Contract.ObsERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObsERC20 *ObsERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObsERC20.Contract.ObsERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObsERC20 *ObsERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObsERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObsERC20 *ObsERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObsERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObsERC20 *ObsERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObsERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.Allowance(&_ObsERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.Allowance(&_ObsERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.BalanceOf(&_ObsERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.BalanceOf(&_ObsERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20Session) Decimals() (uint8, error) {
	return _ObsERC20.Contract.Decimals(&_ObsERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20CallerSession) Decimals() (uint8, error) {
	return _ObsERC20.Contract.Decimals(&_ObsERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20Session) Name() (string, error) {
	return _ObsERC20.Contract.Name(&_ObsERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20CallerSession) Name() (string, error) {
	return _ObsERC20.Contract.Name(&_ObsERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20Session) Symbol() (string, error) {
	return _ObsERC20.Contract.Symbol(&_ObsERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20CallerSession) Symbol() (string, error) {
	return _ObsERC20.Contract.Symbol(&_ObsERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20Session) TotalSupply() (*big.Int, error) {
	return _ObsERC20.Contract.TotalSupply(&_ObsERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ObsERC20.Contract.TotalSupply(&_ObsERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Approve(&_ObsERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Approve(&_ObsERC20.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObsERC20 *ObsERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.DecreaseAllowance(&_ObsERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.DecreaseAllowance(&_ObsERC20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObsERC20 *ObsERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.IncreaseAllowance(&_ObsERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.IncreaseAllowance(&_ObsERC20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Transfer(&_ObsERC20.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Transfer(&_ObsERC20.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.TransferFrom(&_ObsERC20.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.TransferFrom(&_ObsERC20.TransactOpts, from, to, amount)
}

// ObsERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ObsERC20 contract.
type ObsERC20ApprovalIterator struct {
	Event *ObsERC20Approval // Event containing the contract specifics and raw log

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
func (it *ObsERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObsERC20Approval)
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
		it.Event = new(ObsERC20Approval)
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
func (it *ObsERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObsERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObsERC20Approval represents a Approval event raised by the ObsERC20 contract.
type ObsERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ObsERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObsERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ObsERC20ApprovalIterator{contract: _ObsERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ObsERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObsERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObsERC20Approval)
				if err := _ObsERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ObsERC20 *ObsERC20Filterer) ParseApproval(log types.Log) (*ObsERC20Approval, error) {
	event := new(ObsERC20Approval)
	if err := _ObsERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObsERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ObsERC20 contract.
type ObsERC20TransferIterator struct {
	Event *ObsERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ObsERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObsERC20Transfer)
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
		it.Event = new(ObsERC20Transfer)
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
func (it *ObsERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObsERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObsERC20Transfer represents a Transfer event raised by the ObsERC20 contract.
type ObsERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ObsERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObsERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ObsERC20TransferIterator{contract: _ObsERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ObsERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObsERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObsERC20Transfer)
				if err := _ObsERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ObsERC20 *ObsERC20Filterer) ParseTransfer(log types.Log) (*ObsERC20Transfer, error) {
	event := new(ObsERC20Transfer)
	if err := _ObsERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
