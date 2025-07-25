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
	Bin: "0x608060405234801561000f575f5ffd5b5061001861001d565b6100fc565b5f6100266100bd565b805490915068010000000000000000900460ff16156100585760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100ba5780546001600160401b0319166001600160401b0390811782556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916100b1916100e7565b60405180910390a15b50565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6001600160401b0382168152602081016100e1565b6112d2806101095f395ff3fe608060405234801561000f575f5ffd5b50600436106100e5575f3560e01c80638da5cb5b11610088578063a9059cbb11610063578063a9059cbb1461021d578063c4d66de814610230578063dd62ed3e14610243578063f2fde38b1461029a575f5ffd5b80638da5cb5b146101ca57806395d89b41146102025780639f9976af1461020a575f5ffd5b806323b872dd116100c357806323b872dd14610157578063313ce5671461016a57806370a0823114610179578063715018a6146101c0575f5ffd5b806306fdde03146100e9578063095ea7b31461010757806318160ddd14610127575b5f5ffd5b6100f16102ad565b6040516100fe9190610d2a565b60405180910390f35b61011a610115366004610d81565b610380565b6040516100fe9190610dc1565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace02545b6040516100fe9190610dd5565b61011a610165366004610de3565b610399565b60126040516100fe9190610e32565b61014a610187366004610e40565b6001600160a01b03165f9081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00602052604090205490565b6101c86103bc565b005b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516100fe9190610e66565b6100f16103cf565b6101c8610218366004610ec2565b610420565b61011a61022b366004610d81565b6104c9565b6101c861023e366004610e40565b6104d6565b61014a610251366004610f07565b6001600160a01b039182165f9081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace016020908152604080832093909416825291909152205490565b6101c86102a8366004610e40565b6106c2565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0380546060917f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00916102fe90610f48565b80601f016020809104026020016040519081016040528092919081815260200182805461032a90610f48565b80156103755780601f1061034c57610100808354040283529160200191610375565b820191905f5260205f20905b81548152906001019060200180831161035857829003601f168201915b505050505091505090565b5f3361038d818585610718565b60019150505b92915050565b5f336103a6858285610725565b6103b18585856107c6565b506001949350505050565b6103c461083c565b6103cd5f6108b0565b565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0480546060917f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00916102fe90610f48565b5f546001600160a01b031633146104525760405162461bcd60e51b815260040161044990610fce565b60405180910390fd5b5f8190036104725760405162461bcd60e51b815260040161044990610fde565b5f5b818110156104c4576104bc83838381811061049157610491611019565b90506020028101906104a3919061102d565b6104b59061010081019060e001610e40565b600161092d565b600101610474565b505050565b5f3361038d8185856107c6565b5f6104df610965565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f8115801561050b5750825b90505f8267ffffffffffffffff1660011480156105275750303b155b905081158015610535575080155b1561056c576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156105a057845468ff00000000000000001916680100000000000000001785555b6001600160a01b0386166105c65760405162461bcd60e51b8152600401610449906110a8565b61063a6040518060400160405280600381526020017f5a656e00000000000000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f5a454e000000000000000000000000000000000000000000000000000000000081525061098d565b6106433361099f565b5f805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03881617905583156106ba57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906106b1906001906110db565b60405180910390a15b505050505050565b6106ca61083c565b6001600160a01b03811661070c575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016104499190610e66565b610715816108b0565b50565b6104c483838360016109b0565b6001600160a01b038381165f9081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0160209081526040808320938616835292905220545f198110156107c057818110156107b2578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401610449939291906110e9565b6107c084848484035f6109b0565b50505050565b6001600160a01b038316610808575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016104499190610e66565b6001600160a01b038216610831575f60405163ec442f0560e01b81526004016104499190610e66565b6104c4838383610ad7565b3361086e7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146103cd57336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016104499190610e66565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b6001600160a01b038216610956575f60405163ec442f0560e01b81526004016104499190610e66565b6109615f8383610ad7565b5050565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610393565b610995610c27565b6109618282610c65565b6109a7610c27565b61071581610cc8565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace006001600160a01b038516610a13575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016104499190610e66565b6001600160a01b038416610a55575f6040517f94280d620000000000000000000000000000000000000000000000000000000081526004016104499190610e66565b6001600160a01b038086165f90815260018301602090815260408083209388168352929052208390558115610ad057836001600160a01b0316856001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92585604051610ac79190610dd5565b60405180910390a35b5050505050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace006001600160a01b038416610b245781816002015f828254610b19919061112d565b90915550610b9a9050565b6001600160a01b0384165f9081526020829052604090205482811015610b7c578481846040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401610449939291906110e9565b6001600160a01b0385165f9081526020839052604090209083900390555b6001600160a01b038316610bb8576002810180548390039055610bd6565b6001600160a01b0383165f9081526020829052604090208054830190555b826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610c199190610dd5565b60405180910390a350505050565b610c2f610cd0565b6103cd576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c6d610c27565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace007f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace03610cb984826111e0565b50600481016107c083826111e0565b6106ca610c27565b5f610cd9610965565b5468010000000000000000900460ff16919050565b8281835e505f910152565b5f610d02825190565b808452602084019350610d19818560208601610cee565b601f01601f19169290920192915050565b60208082528101610d3b8184610cf9565b9392505050565b5f6001600160a01b038216610393565b610d5b81610d42565b8114610715575f5ffd5b803561039381610d52565b80610d5b565b803561039381610d70565b5f5f60408385031215610d9557610d955f5ffd5b610d9f8484610d65565b9150610dae8460208501610d76565b90509250929050565b8015155b82525050565b602081016103938284610db7565b80610dbb565b602081016103938284610dcf565b5f5f5f60608486031215610df857610df85f5ffd5b610e028585610d65565b9250610e118560208601610d65565b9150610e208560408601610d76565b90509250925092565b60ff8116610dbb565b602081016103938284610e29565b5f60208284031215610e5357610e535f5ffd5b610d3b8383610d65565b610dbb81610d42565b602081016103938284610e5d565b5f5f83601f840112610e8757610e875f5ffd5b50813567ffffffffffffffff811115610ea157610ea15f5ffd5b602083019150836020820283011115610ebb57610ebb5f5ffd5b9250929050565b5f5f60208385031215610ed657610ed65f5ffd5b823567ffffffffffffffff811115610eef57610eef5f5ffd5b610efb85828601610e74565b92509250509250929050565b5f5f60408385031215610f1b57610f1b5f5ffd5b610f258484610d65565b9150610dae8460208501610d65565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680610f5c57607f821691505b602082108103610f6e57610f6e610f34565b50919050565b602c8152602081017f43616c6c65723a2063616c6c6572206973206e6f74207468652064657369676e81527f6174656420616464726573730000000000000000000000000000000000000000602082015290505b60400190565b6020808252810161039381610f74565b6020808252810161039381601a81527f4e6f207472616e73616374696f6e7320746f20636f6e76657274000000000000602082015260400190565b634e487b7160e01b5f52603260045260245ffd5b5f823561013e1936849003018112611046576110465f5ffd5b9190910192915050565b60248152602081017f496e76616c6964207472616e73616374696f6e20616e616c797a65722061646481527f726573730000000000000000000000000000000000000000000000000000000060208201529050610fc8565b6020808252810161039381611050565b5f610393826110c5565b90565b67ffffffffffffffff1690565b610dbb816110b8565b6020810161039382846110d2565b606081016110f78286610e5d565b6111046020830185610dcf565b6111116040830184610dcf565b949350505050565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561039357610393611119565b634e487b7160e01b5f52604160045260245ffd5b5f6103936110c28381565b61116883611154565b81545f1960089490940293841b1916921b91909117905550565b5f6104c481848461115f565b81811015610961576111a05f82611182565b60010161118e565b601f8211156104c4575f818152602090206020601f850104810160208510156111ce5750805b610ad06020601f86010483018261118e565b815167ffffffffffffffff8111156111fa576111fa611140565b6112048254610f48565b61120f8282856111a8565b506020601f821160018114611242575f831561122b5750848201515b5f19600885021c1981166002850217855550610ad0565b5f84815260208120601f198516915b828110156112715787850151825560209485019460019092019101611251565b508482101561128d57838701515f19601f87166008021c191681555b5050505060020260010190555056fea2646970667358221220ff3070db823350c324d6985771d24cc42fc36731ec67c4cb3c22484f3c4cf26f64736f6c634300081c0033",
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
