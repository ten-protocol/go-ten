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

// TransactionDecoderTransaction is an auto generated low-level Go binding around an user-defined struct.
type TransactionDecoderTransaction struct {
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transactionsAnalyzer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"accessList\",\"type\":\"address[]\"}],\"internalType\":\"structTransactionDecoder.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"name\":\"onBlockEnd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523462000031576200001f62000018620001b2565b916200023b565b6040516112256200057e823961122590f35b600080fd5b634e487b7160e01b600052604160045260246000fd5b90601f01601f191681019081106001600160401b038211176200006e57604052565b62000036565b906200008b6200008360405190565b92836200004c565b565b6001600160a01b031690565b90565b6001600160a01b038116036200003157565b905051906200008b826200009c565b6001600160401b0381116200006e57602090601f01601f19160190565b60005b838110620000ee5750506000910152565b8181015183820152602001620000dd565b90929192620001186200011282620000bd565b62000074565b9381855260208501908284011162000031576200008b92620000da565b9080601f83011215620000315781516200009992602001620000ff565b9160608383031262000031576200016a8284620000ae565b60208401519093906001600160401b0381116200003157836200018f91830162000135565b60408201519093906001600160401b038111620000315762000099920162000135565b620001d5620017a380380380620001c98162000074565b92833981019062000152565b909192565b906001600160a01b03905b9181191691161790565b62000099906200008d906001600160a01b031682565b6200009990620001ef565b620000999062000205565b906200022f62000099620002379262000210565b8254620001da565b9055565b906200008b926200024d913362000255565b60066200021b565b906200008b92916200027d565b6200008d62000099620000999290565b620000999062000262565b916200028991620004ea565b62000295600062000272565b6001600160a01b0381166001600160a01b03831614620002bb57506200008b906200051b565b620002ee90620002ca60405190565b631e4fbdf760e01b8152918291600483016001600160a01b03909116815260200190565b0390fd5b634e487b7160e01b600052602260045260246000fd5b90600160028304921680156200032b575b60208310146200032557565b620002f2565b91607f169162000319565b9160001960089290920291821b911b620001e5565b6200009962000099620000999290565b9190620003706200009962000237936200034b565b90835462000336565b6200008b916000916200035b565b81811062000393575050565b80620003a3600060019362000379565b0162000387565b9190601f8111620003ba57505050565b620003ce6200008b93600052602060002090565b906020601f840181900483019310620003f2575b6020601f90910104019062000387565b9091508190620003e2565b9062000407815190565b906001600160401b0382116200006e576200042f8262000428855462000308565b85620003aa565b602090601f83116001146200046e576200023792916000918362000462575b5050600019600883021c1916906002021790565b0151905038806200044e565b601f198316916200048485600052602060002090565b9260005b818110620004c557509160029391856001969410620004ab575b50505002019055565b01516000196008601f8516021c19169055388080620004a2565b9193602060018192878701518155019501920162000488565b906200008b91620003fd565b90620004fc6200008b926003620004de565b6004620004de565b62000099906200008d565b62000099905462000504565b6200052760056200050f565b90620005358160056200021b565b6200056c620005657f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09362000210565b9162000210565b916200057760405190565b600090a356fe6080604052600436101561001257600080fd5b60003560e01c806306fdde03146100f2578063095ea7b3146100ed57806318160ddd146100e857806323b872dd146100e3578063313ce567146100de57806340c10f19146100d9578063630ac52c146100d457806370a08231146100cf578063715018a6146100ca5780638da5cb5b146100c557806395d89b41146100c0578063a9059cbb146100bb578063dd62ed3e146100b65763f2fde38b0361010257610466565b61044a565b61040b565b6103f0565b6103b1565b610399565b61037e565b610351565b6102d9565b6102aa565b61028e565b610232565b610204565b610176565b600091031261010257565b600080fd5b60005b83811061011a5750506000910152565b818101518382015260200161010a565b61014b61015460209361015e9361013f815190565b80835293849260200190565b95869101610107565b601f01601f191690565b0190565b60208082526101739291019061012a565b90565b34610102576101863660046100f7565b61019d6101916107aa565b60405191829182610162565b0390f35b6001600160a01b031690565b6001600160a01b0381165b0361010257565b905035906101cc826101ad565b565b806101b8565b905035906101cc826101ce565b919060408382031261010257610173906101fb81856101bf565b936020016101d4565b346101025761019d61022061021a3660046101e1565b9061085f565b60405191829182901515815260200190565b34610102576102423660046100f7565b61019d61024d6107ec565b6040515b9182918290815260200190565b90916060828403126101025761017361027784846101bf565b9361028581602086016101bf565b936040016101d4565b346101025761019d6102206102a436600461025e565b9161086a565b34610102576102ba3660046100f7565b61019d6102c56107d1565b6040519182918260ff909116815260200190565b34610102576102f26102ec3660046101e1565b906111e5565b604051005b909182601f830112156101025781359167ffffffffffffffff831161010257602001926020830284011161010257565b9060208282031261010257813567ffffffffffffffff81116101025761034d92016102f7565b9091565b34610102576102f2610364366004610327565b90611187565b9060208282031261010257610173916101bf565b346101025761019d61024d61039436600461036a565b61080e565b34610102576103a93660046100f7565b6102f26104b6565b34610102576103c13660046100f7565b61019d6103d66005546001600160a01b031690565b604051918291826001600160a01b03909116815260200190565b34610102576104003660046100f7565b61019d6101916107b4565b346101025761019d6102206104213660046101e1565b9061082a565b9190604083820312610102576101739061044181856101bf565b936020016101bf565b346101025761019d61024d610460366004610427565b9061083f565b34610102576102f261047936600461036a565b6105b7565b6104866104be565b6101cc6104a4565b6101a16101736101739290565b6101739061048e565b6101cc6104b1600061049b565b610618565b6101cc61047e565b6005546001600160a01b031633906104de825b916001600160a01b031690565b036104e65750565b610530906104f360405190565b9182917f118cdaa7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b0390fd5b6101cc906105406104be565b61054a600061049b565b6001600160a01b0381166001600160a01b0383161461056d57506101cc90610618565b6105309061057a60405190565b9182917f1e4fbdf7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6101cc90610534565b906001600160a01b03905b9181191691161790565b6101a1610173610173926001600160a01b031690565b610173906105d5565b610173906105eb565b9061060d610173610614926105f4565b82546105c0565b9055565b6005546001600160a01b0316906106308160056105fd565b61066361065d7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936105f4565b916105f4565b9161066d60405190565b80805b0390a3565b634e487b7160e01b600052602260045260246000fd5b90600160028304921680156106ab575b60208310146106a657565b610675565b91607f169161069b565b805460009392916106d26106c88361068b565b8085529360200190565b916001811690811561072457506001146106eb57505050565b6106fe9192939450600052602060002090565b916000925b8184106107105750500190565b805484840152602090930192600101610703565b92949550505060ff1916825215156020020190565b90610173916106b5565b634e487b7160e01b600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff82111761077b57604052565b610743565b906101cc61079a9261079160405190565b93848092610739565b0383610759565b61017390610780565b61017360036107a1565b61017360046107a1565b6107cb6101736101739290565b60ff1690565b61017360126107be565b6101739081565b61017390546107db565b61017360026107e2565b90610800906105f4565b600052602052604060002090565b6108256101739161081d600090565b5060006107f6565b6107e2565b61083a91903361087b565b61087b565b600190565b6101739161085a61082592610852600090565b5060016107f6565b6107f6565b61083a919033610b21565b61083a929190610835833383610c4b565b929190610888600061049b565b936001600160a01b0385166001600160a01b03821614610910576001600160a01b0385166001600160a01b038316146108c6576101cc9394506109db565b610530856108d360405190565b9182917fec442f05000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6105308561091d60405190565b9182917f96c6fd1e000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6001600160a01b0390911681526060810193926101cc929091604091610981906020830152565b0152565b90600019906105cb565b6101736101736101739290565b906109ac6101736106149261098f565b8254610985565b634e487b7160e01b600052601160045260246000fd5b919082018092116109d657565b6109b3565b8160006109e78161049b565b6001600160a01b0381166001600160a01b03851603610aac57610a21906101a1610a1a88610a1560026107e2565b6109c9565b600261099c565b6001600160a01b03831603610a87575050610a47610a1a84610a4360026107e2565b0390565b610670610a7d610a777fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef936105f4565b936105f4565b9361025160405190565b610aa791610a94916107f6565b610aa18561015e836107e2565b9061099c565b610a47565b909150610abc61082584846107f6565b858110610ae4578492916101a1610ad588610a21940390565b610adf87866107f6565b61099c565b8361053087610af260405190565b9384937fe450d38c0000000000000000000000000000000000000000000000000000000085526004850161095a565b90916101cc926001925b909192610b38600061049b565b6001600160a01b0381166001600160a01b03841614610c01576001600160a01b0381166001600160a01b03851614610bb75750610b7e84610adf8561085a8660016107f6565b610b8757505050565b610670610a7d610a777f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925936105f4565b61053090610bc460405190565b9182917f94280d62000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b61053090610c0e60405190565b9182917fe602df05000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b90929192610c59818361083f565b936000198503610c6b575b5050509050565b808510610c9157610c7f90610c8894950390565b90600092610b2b565b80388080610c64565b906105308592610ca060405190565b9384937ffb8f41b20000000000000000000000000000000000000000000000000000000085526004850161095a565b15610cd657565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f43616c6c65723a2063616c6c6572206973206e6f74207468652064657369676e60448201527f61746564206164647265737300000000000000000000000000000000000000006064820152608490fd5b906101cc91610d8233610d7c6104d16101a16006546001600160a01b031690565b14610ccf565b6110b4565b60001981146109d65760010190565b634e487b7160e01b600052603260045260246000fd5b9035906101be193682900301821215610102570190565b90821015610dda5760206101739202810190610dac565b610d96565b905051906101cc826101ad565b906020828203126101025761017391610ddf565b60ff81166101b8565b905035906101cc82610e00565b50610173906020810190610e09565b506101739060208101906101d4565b506101739060208101906101bf565b9035601e19368390030181121561010257016020813591019167ffffffffffffffff82116101025736829003831361010257565b90826000939282370152565b919061015481610e9a8161015e9560209181520190565b8095610e77565b9035601e19368390030181121561010257016020813591019167ffffffffffffffff821161010257602082023603831361010257565b818352602090920191906000825b828210610ef3575050505090565b90919293610f22610f1b600192610f0a8886610e34565b6001600160a01b0316815260200190565b9560200190565b93920190610ee5565b61017391611088610ffc6101c08301610f4e610f478680610e16565b60ff168552565b610f65610f5e6020870187610e25565b6020860152565b610f7c610f756040870187610e25565b6040860152565b610f93610f8c6060870187610e25565b6060860152565b610faa610fa36080870187610e25565b6080860152565b610fca610fba60a0870187610e34565b6001600160a01b031660a0860152565b610fe1610fda60c0870187610e25565b60c0860152565b610fee60e0860186610e43565b9085830360e0870152610e83565b9261101961100e610100830183610e16565b60ff16610100850152565b61103261102a610120830183610e25565b610120850152565b61104b611043610140830183610e25565b610140850152565b61106461105c610160830183610e25565b610160850152565b61107d611075610180830183610e25565b610180850152565b6101a0810190610ea1565b916101a0818503910152610ed7565b602080825261017392910190610f2b565b6040513d6000823e3d90fd5b91906110c0600061098f565b818110156111815761111890602063fe7fbd1873__$3a457adf3f2d33c60ddff735cdd91d6a07$__61110d6110f685888b610dc3565b9261110060405190565b9687948593849360e01b90565b835260048301611097565b03915af491821561117c57611149926111449160009161114e575b5061113e600161098f565b90611191565b610d87565b6110c0565b61116f915060203d8111611175575b6111678183610759565b810190610dec565b38611133565b503d61115d565b6110a8565b50509050565b906101cc91610d5b565b919061119d600061049b565b926001600160a01b0384166001600160a01b038216146111c1576101cc92936109db565b610530846108d360405190565b906101cc916111db6104be565b906101cc91611191565b906101cc916111ce56fea2646970667358221220319d7df78beef1437b2f4a0a18444da02fd93e6b175e55d493734819ea0f6b7664736f6c63430008140033",
}

// ZenBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use ZenBaseMetaData.ABI instead.
var ZenBaseABI = ZenBaseMetaData.ABI

// ZenBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZenBaseMetaData.Bin instead.
var ZenBaseBin = ZenBaseMetaData.Bin

// DeployZenBase deploys a new Ethereum contract, binding an instance of ZenBase to it.
func DeployZenBase(auth *bind.TransactOpts, backend bind.ContractBackend, transactionsAnalyzer common.Address, name string, symbol string) (common.Address, *types.Transaction, *ZenBase, error) {
	parsed, err := ZenBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZenBaseBin), backend, transactionsAnalyzer, name, symbol)
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

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ZenBase *ZenBaseTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ZenBase *ZenBaseSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Mint(&_ZenBase.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ZenBase *ZenBaseTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZenBase.Contract.Mint(&_ZenBase.TransactOpts, to, amount)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x630ac52c.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[])[] transactions) returns()
func (_ZenBase *ZenBaseTransactor) OnBlockEnd(opts *bind.TransactOpts, transactions []TransactionDecoderTransaction) (*types.Transaction, error) {
	return _ZenBase.contract.Transact(opts, "onBlockEnd", transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x630ac52c.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[])[] transactions) returns()
func (_ZenBase *ZenBaseSession) OnBlockEnd(transactions []TransactionDecoderTransaction) (*types.Transaction, error) {
	return _ZenBase.Contract.OnBlockEnd(&_ZenBase.TransactOpts, transactions)
}

// OnBlockEnd is a paid mutator transaction binding the contract method 0x630ac52c.
//
// Solidity: function onBlockEnd((uint8,uint256,uint256,uint256,uint256,address,uint256,bytes,uint8,bytes32,bytes32,uint256,uint256,address[])[] transactions) returns()
func (_ZenBase *ZenBaseTransactorSession) OnBlockEnd(transactions []TransactionDecoderTransaction) (*types.Transaction, error) {
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
