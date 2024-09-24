// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ZenBase

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
	TxType               uint8
	ChainId              *big.Int
	Nonce                *big.Int
	GasPrice             *big.Int
	GasLimit             *big.Int
	To                   common.Address
	Value                *big.Int
	Data                 []byte
	V                    uint8
	R                    [32]byte
	S                    [32]byte
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
	AccessList           []common.Address
}

// ZenBaseMetaData contains all meta data concerning the ZenBase contract.
var ZenBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transactionAnalyzer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"transactionDecoder\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"accessList\",\"type\":\"address[]\"}],\"internalType\":\"structStructs.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlockEnd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523462000031576200001f62000018620000e5565b90620002c0565b6040516111e46200063d82396111e490f35b600080fd5b634e487b7160e01b600052604160045260246000fd5b90601f01601f191681019081106001600160401b038211176200006e57604052565b62000036565b906200008b6200008360405190565b92836200004c565b565b6001600160a01b031690565b90565b6001600160a01b038116036200003157565b905051906200008b826200009c565b919060408382031262000031576200009990620000db8185620000ae565b93602001620000ae565b620001086200182180380380620000fc8162000074565b928339810190620000bd565b9091565b6001600160401b0381116200006e57602090601f01601f19160190565b906200013f62000139836200010c565b62000074565b918252565b62000150600362000129565b622d32b760e91b602082015290565b6200009962000144565b62000175600362000129565b622d22a760e91b602082015290565b6200009962000169565b6200008d62000099620000999290565b62000099906200018e565b15620001b157565b60405162461bcd60e51b8152602060048201526024808201527f496e76616c6964207472616e73616374696f6e20616e616c797a6572206164646044820152637265737360e01b6064820152608490fd5b0390fd5b156200020e57565b60405162461bcd60e51b815260206004820152602360248201527f496e76616c6964207472616e73616374696f6e206465636f646572206164647260448201526265737360e81b6064820152608490fd5b906001600160a01b03905b9181191691161790565b62000099906200008d906001600160a01b031682565b620000999062000274565b62000099906200028a565b90620002b462000099620002bc9262000295565b82546200025f565b9055565b90620003386200008b92620002ea33620002d96200015f565b620002e362000184565b9162000340565b620003306200031e620002fe60006200019e565b6200008d6001600160a01b0382166001600160a01b0386161415620001a9565b6001600160a01b038516141562000206565b6006620002a0565b6007620002a0565b916200034c91620005a9565b6200035860006200019e565b6001600160a01b0381166001600160a01b038316146200037e57506200008b90620005da565b62000202906200038d60405190565b631e4fbdf760e01b8152918291600483016001600160a01b03909116815260200190565b634e487b7160e01b600052602260045260246000fd5b9060016002830492168015620003ea575b6020831014620003e457565b620003b1565b91607f1691620003d8565b9160001960089290920291821b911b6200026a565b6200009962000099620000999290565b91906200042f62000099620002bc936200040a565b908354620003f5565b6200008b916000916200041a565b81811062000452575050565b8062000462600060019362000438565b0162000446565b9190601f81116200047957505050565b6200048d6200008b93600052602060002090565b906020601f840181900483019310620004b1575b6020601f90910104019062000446565b9091508190620004a1565b90620004c6815190565b906001600160401b0382116200006e57620004ee82620004e78554620003c7565b8562000469565b602090601f83116001146200052d57620002bc92916000918362000521575b5050600019600883021c1916906002021790565b0151905038806200050d565b601f198316916200054385600052602060002090565b9260005b81811062000584575091600293918560019694106200056a575b50505002019055565b01516000196008601f8516021c1916905538808062000561565b9193602060018192878701518155019501920162000547565b906200008b91620004bc565b90620005bb6200008b9260036200059d565b60046200059d565b62000099906200008d565b620000999054620005c3565b620005e66005620005ce565b90620005f4816005620002a0565b6200062b620006247f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09362000295565b9162000295565b916200063660405190565b600090a356fe6080604052600436101561001257600080fd5b60003560e01c806306fdde03146100e2578063095ea7b3146100dd57806318160ddd146100d857806323b872dd146100d3578063313ce567146100ce578063630ac52c146100c957806370a08231146100c4578063715018a6146100bf5780638da5cb5b146100ba57806395d89b41146100b5578063a9059cbb146100b0578063dd62ed3e146100ab5763f2fde38b036100f25761043d565b610421565b6103e2565b6103c7565b610388565b610370565b610355565b610323565b61029a565b61027e565b610222565b6101f4565b610166565b60009103126100f257565b600080fd5b60005b83811061010a5750506000910152565b81810151838201526020016100fa565b61013b61014460209361014e9361012f815190565b80835293849260200190565b958691016100f7565b601f01601f191690565b0190565b60208082526101639291019061011a565b90565b346100f2576101763660046100e7565b61018d610181610781565b60405191829182610152565b0390f35b6001600160a01b031690565b6001600160a01b0381165b036100f257565b905035906101bc8261019d565b565b806101a8565b905035906101bc826101be565b91906040838203126100f257610163906101eb81856101af565b936020016101c4565b346100f25761018d61021061020a3660046101d1565b90610836565b60405191829182901515815260200190565b346100f2576102323660046100e7565b61018d61023d6107c3565b6040515b9182918290815260200190565b90916060828403126100f25761016361026784846101af565b9361027581602086016101af565b936040016101c4565b346100f25761018d61021061029436600461024e565b91610841565b346100f2576102aa3660046100e7565b61018d6102b56107a8565b6040519182918260ff909116815260200190565b909182601f830112156100f25781359167ffffffffffffffff83116100f25760200192602083028401116100f257565b906020828203126100f257813567ffffffffffffffff81116100f25761031f92016102c9565b9091565b346100f25761033c6103363660046102f9565b90611167565b604051005b906020828203126100f257610163916101af565b346100f25761018d61023d61036b366004610341565b6107e5565b346100f2576103803660046100e7565b61033c61048d565b346100f2576103983660046100e7565b61018d6103ad6005546001600160a01b031690565b604051918291826001600160a01b03909116815260200190565b346100f2576103d73660046100e7565b61018d61018161078b565b346100f25761018d6102106103f83660046101d1565b90610801565b91906040838203126100f2576101639061041881856101af565b936020016101af565b346100f25761018d61023d6104373660046103fe565b90610816565b346100f25761033c610450366004610341565b61058e565b61045d610495565b6101bc61047b565b6101916101636101639290565b61016390610465565b6101bc6104886000610472565b6105ef565b6101bc610455565b6005546001600160a01b031633906104b5825b916001600160a01b031690565b036104bd5750565b610507906104ca60405190565b9182917f118cdaa7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b0390fd5b6101bc90610517610495565b6105216000610472565b6001600160a01b0381166001600160a01b0383161461054457506101bc906105ef565b6105079061055160405190565b9182917f1e4fbdf7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6101bc9061050b565b906001600160a01b03905b9181191691161790565b610191610163610163926001600160a01b031690565b610163906105ac565b610163906105c2565b906105e46101636105eb926105cb565b8254610597565b9055565b6005546001600160a01b0316906106078160056105d4565b61063a6106347f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936105cb565b916105cb565b9161064460405190565b80805b0390a3565b634e487b7160e01b600052602260045260246000fd5b9060016002830492168015610682575b602083101461067d57565b61064c565b91607f1691610672565b805460009392916106a961069f83610662565b8085529360200190565b91600181169081156106fb57506001146106c257505050565b6106d59192939450600052602060002090565b916000925b8184106106e75750500190565b8054848401526020909301926001016106da565b92949550505060ff1916825215156020020190565b906101639161068c565b634e487b7160e01b600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff82111761075257604052565b61071a565b906101bc6107719261076860405190565b93848092610710565b0383610730565b61016390610757565b6101636003610778565b6101636004610778565b6107a26101636101639290565b60ff1690565b6101636012610795565b6101639081565b61016390546107b2565b61016360026107b9565b906107d7906105cb565b600052602052604060002090565b6107fc610163916107f4600090565b5060006107cd565b6107b9565b610811919033610852565b610852565b600190565b610163916108316107fc92610829600090565b5060016107cd565b6107cd565b610811919033610af8565b61081192919061080c833383610c22565b92919061085f6000610472565b936001600160a01b0385166001600160a01b038216146108e7576001600160a01b0385166001600160a01b0383161461089d576101bc9394506109b2565b610507856108aa60405190565b9182917fec442f05000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b610507856108f460405190565b9182917f96c6fd1e000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6001600160a01b0390911681526060810193926101bc929091604091610958906020830152565b0152565b90600019906105a2565b6101636101636101639290565b906109836101636105eb92610966565b825461095c565b634e487b7160e01b600052601160045260246000fd5b919082018092116109ad57565b61098a565b8160006109be81610472565b6001600160a01b0381166001600160a01b03851603610a83576109f8906101916109f1886109ec60026107b9565b6109a0565b6002610973565b6001600160a01b03831603610a5e575050610a1e6109f184610a1a60026107b9565b0390565b610647610a54610a4e7fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef936105cb565b936105cb565b9361024160405190565b610a7e91610a6b916107cd565b610a788561014e836107b9565b90610973565b610a1e565b909150610a936107fc84846107cd565b858110610abb57849291610191610aac886109f8940390565b610ab687866107cd565b610973565b8361050787610ac960405190565b9384937fe450d38c00000000000000000000000000000000000000000000000000000000855260048501610931565b90916101bc926001925b909192610b0f6000610472565b6001600160a01b0381166001600160a01b03841614610bd8576001600160a01b0381166001600160a01b03851614610b8e5750610b5584610ab6856108318660016107cd565b610b5e57505050565b610647610a54610a4e7f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925936105cb565b61050790610b9b60405190565b9182917f94280d62000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b61050790610be560405190565b9182917fe602df05000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b90929192610c308183610816565b936000198503610c42575b5050509050565b808510610c6857610c5690610c5f94950390565b90600092610b02565b80388080610c3b565b906105078592610c7760405190565b9384937ffb8f41b200000000000000000000000000000000000000000000000000000000855260048501610931565b15610cad57565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f43616c6c65723a2063616c6c6572206973206e6f74207468652064657369676e60448201527f61746564206164647265737300000000000000000000000000000000000000006064820152608490fd5b906101bc91610d5933610d536104a86101916006546001600160a01b031690565b14610ca6565b61108b565b60001981146109ad5760010190565b634e487b7160e01b600052603260045260246000fd5b9035906101be1936829003018212156100f2570190565b90821015610db15760206101639202810190610d83565b610d6d565b905051906101bc8261019d565b906020828203126100f25761016391610db6565b60ff81166101a8565b905035906101bc82610dd7565b50610163906020810190610de0565b506101639060208101906101c4565b506101639060208101906101af565b9035601e1936839003018112156100f257016020813591019167ffffffffffffffff82116100f2573682900383136100f257565b90826000939282370152565b919061014481610e718161014e9560209181520190565b8095610e4e565b9035601e1936839003018112156100f257016020813591019167ffffffffffffffff82116100f25760208202360383136100f257565b818352602090920191906000825b828210610eca575050505090565b90919293610ef9610ef2600192610ee18886610e0b565b6001600160a01b0316815260200190565b9560200190565b93920190610ebc565b6101639161105f610fd36101c08301610f25610f1e8680610ded565b60ff168552565b610f3c610f356020870187610dfc565b6020860152565b610f53610f4c6040870187610dfc565b6040860152565b610f6a610f636060870187610dfc565b6060860152565b610f81610f7a6080870187610dfc565b6080860152565b610fa1610f9160a0870187610e0b565b6001600160a01b031660a0860152565b610fb8610fb160c0870187610dfc565b60c0860152565b610fc560e0860186610e1a565b9085830360e0870152610e5a565b92610ff0610fe5610100830183610ded565b60ff16610100850152565b611009611001610120830183610dfc565b610120850152565b61102261101a610140830183610dfc565b610140850152565b61103b611033610160830183610dfc565b610160850152565b61105461104c610180830183610dfc565b610180850152565b6101a0810190610e78565b916101a0818503910152610eae565b602080825261016392910190610f02565b6040513d6000823e3d90fd5b91906110976000610966565b81811015611161576110f89060206110c26110bd6110bd6007546001600160a01b031690565b6105cb565b6323bc02dd906110ed6110d685888b610d9a565b926110e060405190565b9687948593849360e01b90565b83526004830161106e565b03915afa91821561115c57611129926111249160009161112e575b5061111e6001610966565b90611171565b610d5e565b611097565b61114f915060203d8111611155575b6111478183610730565b810190610dc3565b38611113565b503d61113d565b61107f565b50509050565b906101bc91610d32565b919061117d6000610472565b926001600160a01b0384166001600160a01b038216146111a1576101bc92936109b2565b610507846108aa6040519056fea2646970667358221220278d99c1097531adcce1df939920329f478a45496964d2f7bc6c675345cb9a2364736f6c63430008140033",
}

// ZenBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenBaseMetaData.ABI instead.
var ZenBaseABI = ZenBaseMetaData.ABI

// ZenBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZenBaseMetaData.Bin instead.
var ZenBaseBin = ZenBaseMetaData.Bin

// DeployZenBase deploys a new Ethereum contract, binding an instance of ZenBase to it.
func DeployZenBase(auth *bind.TransactOpts, backend bind.ContractBackend, transactionAnalyzer common.Address, transactionDecoder common.Address) (common.Address, *types.Transaction, *ZenBase, error) {
	parsed, err := ZenBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZenBaseBin), backend, transactionAnalyzer, transactionDecoder)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ZenBase{ZenBaseCaller: ZenBaseCaller{contract: contract}, ZenBaseTransactor: ZenBaseTransactor{contract: contract}, ZenBaseFilterer: ZenBaseFilterer{contract: contract}}, nil
}

// ZenBase is an auto generated Go binding around an Ethereum contract.
type ZenBase struct {
	ZenBaseCaller     // Read-only binding to the contract
	ZenBaseTransactor // Write-only binding to the contract
	ZenBaseFilterer   // Log filterer for contract events
}

// ZenBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZenBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZenBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZenBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZenBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZenBaseSession struct {
	Contract     *ZenBase          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZenBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZenBaseCallerSession struct {
	Contract *ZenBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ZenBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZenBaseTransactorSession struct {
	Contract     *ZenBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ZenBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZenBaseRaw struct {
	Contract *ZenBase // Generic contract binding to access the raw methods on
}

// ZenBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZenBaseCallerRaw struct {
	Contract *ZenBaseCaller // Generic read-only contract binding to access the raw methods on
}

// ZenBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZenBaseTransactorRaw struct {
	Contract *ZenBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZenBase creates a new instance of ZenBase, bound to a specific deployed contract.
func NewZenBase(address common.Address, backend bind.ContractBackend) (*ZenBase, error) {
	contract, err := bindZenBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZenBase{ZenBaseCaller: ZenBaseCaller{contract: contract}, ZenBaseTransactor: ZenBaseTransactor{contract: contract}, ZenBaseFilterer: ZenBaseFilterer{contract: contract}}, nil
}

// NewZenBaseCaller creates a new read-only instance of ZenBase, bound to a specific deployed contract.
func NewZenBaseCaller(address common.Address, caller bind.ContractCaller) (*ZenBaseCaller, error) {
	contract, err := bindZenBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZenBaseCaller{contract: contract}, nil
}

// NewZenBaseTransactor creates a new write-only instance of ZenBase, bound to a specific deployed contract.
func NewZenBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*ZenBaseTransactor, error) {
	contract, err := bindZenBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZenBaseTransactor{contract: contract}, nil
}

// NewZenBaseFilterer creates a new log filterer instance of ZenBase, bound to a specific deployed contract.
func NewZenBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*ZenBaseFilterer, error) {
	contract, err := bindZenBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZenBaseFilterer{contract: contract}, nil
}

// bindZenBase binds a generic wrapper to an already deployed contract.
func bindZenBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZenBaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenBase *ZenBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenBase.Contract.ZenBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenBase *ZenBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBase.Contract.ZenBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenBase *ZenBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenBase.Contract.ZenBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZenBase *ZenBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZenBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZenBase *ZenBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZenBase *ZenBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZenBase.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBase *ZenBaseCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBase *ZenBaseSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenBase.Contract.Allowance(&_ZenBase.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZenBase *ZenBaseCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZenBase.Contract.Allowance(&_ZenBase.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBase *ZenBaseCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBase *ZenBaseSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenBase.Contract.BalanceOf(&_ZenBase.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZenBase *ZenBaseCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZenBase.Contract.BalanceOf(&_ZenBase.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBase *ZenBaseCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBase *ZenBaseSession) Decimals() (uint8, error) {
	return _ZenBase.Contract.Decimals(&_ZenBase.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZenBase *ZenBaseCallerSession) Decimals() (uint8, error) {
	return _ZenBase.Contract.Decimals(&_ZenBase.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBase *ZenBaseCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBase *ZenBaseSession) Name() (string, error) {
	return _ZenBase.Contract.Name(&_ZenBase.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZenBase *ZenBaseCallerSession) Name() (string, error) {
	return _ZenBase.Contract.Name(&_ZenBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenBase *ZenBaseCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenBase *ZenBaseSession) Owner() (common.Address, error) {
	return _ZenBase.Contract.Owner(&_ZenBase.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZenBase *ZenBaseCallerSession) Owner() (common.Address, error) {
	return _ZenBase.Contract.Owner(&_ZenBase.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBase *ZenBaseCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBase *ZenBaseSession) Symbol() (string, error) {
	return _ZenBase.Contract.Symbol(&_ZenBase.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZenBase *ZenBaseCallerSession) Symbol() (string, error) {
	return _ZenBase.Contract.Symbol(&_ZenBase.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBase *ZenBaseCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZenBase.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBase *ZenBaseSession) TotalSupply() (*big.Int, error) {
	return _ZenBase.Contract.TotalSupply(&_ZenBase.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZenBase *ZenBaseCallerSession) TotalSupply() (*big.Int, error) {
	return _ZenBase.Contract.TotalSupply(&_ZenBase.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBase *ZenBaseSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Approve(&_ZenBase.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Approve(&_ZenBase.TransactOpts, spender, value)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x630ac52c.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[])[] transactions) returns()
func (_ZenBase *ZenBaseTransactor) OnBlockEnd(opts *bind.TransactOpts, transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "onBlockEnd", transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x630ac52c.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[])[] transactions) returns()
func (_ZenBase *ZenBaseSession) OnBlockEnd(transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenBase.Contract.OnBlockEnd(&_ZenBase.TransactOpts, transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x630ac52c.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[])[] transactions) returns()
func (_ZenBase *ZenBaseTransactorSession) OnBlockEnd(transactions []StructsTransaction) (*types.Transaction, error) {
	return _ZenBase.Contract.OnBlockEnd(&_ZenBase.TransactOpts, transactions)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenBase *ZenBaseTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenBase *ZenBaseSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenBase.Contract.RenounceOwnership(&_ZenBase.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZenBase *ZenBaseTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZenBase.Contract.RenounceOwnership(&_ZenBase.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Transfer(&_ZenBase.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Transfer(&_ZenBase.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferFrom(&_ZenBase.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ZenBase *ZenBaseTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferFrom(&_ZenBase.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenBase *ZenBaseTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenBase *ZenBaseSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferOwnership(&_ZenBase.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZenBase *ZenBaseTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZenBase.Contract.TransferOwnership(&_ZenBase.TransactOpts, newOwner)
}

// ZenBaseApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ZenBase contract.
type ZenBaseApprovalIterator struct {
	Event *ZenBaseApproval // Event containing the contract specifics and raw log

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
func (it *ZenBaseApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBaseApproval)
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
		it.Event = new(ZenBaseApproval)
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
func (it *ZenBaseApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBaseApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBaseApproval represents a Approval event raised by the ZenBase contract.
type ZenBaseApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenBase *ZenBaseFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZenBaseApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenBase.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZenBaseApprovalIterator{contract: _ZenBase.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZenBase *ZenBaseFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZenBaseApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZenBase.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBaseApproval)
				if err := _ZenBase.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ZenBase *ZenBaseFilterer) ParseApproval(log types.Log) (*ZenBaseApproval, error) {
	event := new(ZenBaseApproval)
	if err := _ZenBase.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBaseOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZenBase contract.
type ZenBaseOwnershipTransferredIterator struct {
	Event *ZenBaseOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZenBaseOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBaseOwnershipTransferred)
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
		it.Event = new(ZenBaseOwnershipTransferred)
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
func (it *ZenBaseOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBaseOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBaseOwnershipTransferred represents a OwnershipTransferred event raised by the ZenBase contract.
type ZenBaseOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenBase *ZenBaseFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZenBaseOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenBase.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZenBaseOwnershipTransferredIterator{contract: _ZenBase.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZenBase *ZenBaseFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZenBaseOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZenBase.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBaseOwnershipTransferred)
				if err := _ZenBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ZenBase *ZenBaseFilterer) ParseOwnershipTransferred(log types.Log) (*ZenBaseOwnershipTransferred, error) {
	event := new(ZenBaseOwnershipTransferred)
	if err := _ZenBase.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZenBaseTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ZenBase contract.
type ZenBaseTransferIterator struct {
	Event *ZenBaseTransfer // Event containing the contract specifics and raw log

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
func (it *ZenBaseTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZenBaseTransfer)
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
		it.Event = new(ZenBaseTransfer)
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
func (it *ZenBaseTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZenBaseTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZenBaseTransfer represents a Transfer event raised by the ZenBase contract.
type ZenBaseTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenBase *ZenBaseFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZenBaseTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenBase.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZenBaseTransferIterator{contract: _ZenBase.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZenBase *ZenBaseFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZenBaseTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZenBase.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZenBaseTransfer)
				if err := _ZenBase.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ZenBase *ZenBaseFilterer) ParseTransfer(log types.Log) (*ZenBaseTransfer, error) {
	event := new(ZenBaseTransfer)
	if err := _ZenBase.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
