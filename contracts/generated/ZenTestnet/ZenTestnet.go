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
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100d0565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161561006e5760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100cd5780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b6113ea806100df6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80638da5cb5b1161008c578063a9059cbb11610066578063a9059cbb14610224578063c4d66de814610237578063dd62ed3e1461024a578063f2fde38b146102a257600080fd5b80638da5cb5b146101d157806395d89b41146102095780639f9976af1461021157600080fd5b806323b872dd116100c857806323b872dd1461015d578063313ce5671461017057806370a082311461017f578063715018a6146101c757600080fd5b806306fdde03146100ef578063095ea7b31461010d57806318160ddd1461012d575b600080fd5b6100f76102b5565b6040516101049190610dda565b60405180910390f35b61012061011b366004610e2c565b61038a565b6040516101049190610e73565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace02545b6040516101049190610e87565b61012061016b366004610e95565b6103a4565b60126040516101049190610eee565b61015061018d366004610efc565b6001600160a01b031660009081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00602052604090205490565b6101cf6103ca565b005b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516101049190610f2e565b6100f76103de565b6101cf61021f366004610f8e565b61042f565b610120610232366004610e2c565b610554565b6101cf610245366004610efc565b610562565b610150610258366004610fd6565b6001600160a01b0391821660009081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace016020908152604080832093909416825291909152205490565b6101cf6102b0366004610efc565b610766565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0380546060917f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00916103069061101f565b80601f01602080910402602001604051908101604052809291908181526020018280546103329061101f565b801561037f5780601f106103545761010080835404028352916020019161037f565b820191906000526020600020905b81548152906001019060200180831161036257829003601f168201915b505050505091505090565b6000336103988185856107bd565b60019150505b92915050565b6000336103b28582856107ca565b6103bd85858561086d565b60019150505b9392505050565b6103d26108e5565b6103dc6000610959565b565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0480546060917f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00916103069061101f565b6000546001600160a01b031633146104625760405162461bcd60e51b8152600401610459906110a5565b60405180910390fd5b60008190036104835760405162461bcd60e51b8152600401610459906110b5565b60005b8181101561054f576104ce8383838181106104a3576104a36110f0565b90506020028101906104b59190611106565b6104c79061010081019060e001610efc565b60016109d7565b7fe848a9a1096c6a1986f56a70fb7fb3250e28b1f56d41fa97ac643492c6c853d1838383818110610501576105016110f0565b90506020028101906105139190611106565b6105259061010081019060e001610efc565b6001604051610535929190611143565b60405180910390a18061054781611174565b915050610486565b505050565b60003361039881858561086d565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156105ad5750825b905060008267ffffffffffffffff1660011480156105ca5750303b155b9050811580156105d8575080155b1561060f576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561064357845468ff00000000000000001916680100000000000000001785555b6001600160a01b0386166106695760405162461bcd60e51b8152600401610459906111e5565b6106dd6040518060400160405280600381526020017f5a656e00000000000000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f5a454e0000000000000000000000000000000000000000000000000000000000815250610a11565b6106e633610a23565b6000805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038816179055831561075e57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061075590600190611210565b60405180910390a15b505050505050565b61076e6108e5565b6001600160a01b0381166107b15760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016104599190610f2e565b6107ba81610959565b50565b61054f8383836001610a34565b6001600160a01b0383811660009081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace01602090815260408083209386168352929052205460001981146108675781811015610858578281836040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526004016104599392919061121e565b61086784848484036000610a34565b50505050565b6001600160a01b0383166108b05760006040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016104599190610f2e565b6001600160a01b0382166108da57600060405163ec442f0560e01b81526004016104599190610f2e565b61054f838383610b5e565b336109177f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146103dc57336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016104599190610f2e565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b6001600160a01b038216610a0157600060405163ec442f0560e01b81526004016104599190610f2e565b610a0d60008383610b5e565b5050565b610a19610cb2565b610a0d8282610d19565b610a2b610cb2565b6107ba81610d7c565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace006001600160a01b038516610a985760006040517fe602df050000000000000000000000000000000000000000000000000000000081526004016104599190610f2e565b6001600160a01b038416610adb5760006040517f94280d620000000000000000000000000000000000000000000000000000000081526004016104599190610f2e565b6001600160a01b03808616600090815260018301602090815260408083209388168352929052208390558115610b5757836001600160a01b0316856001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92585604051610b4e9190610e87565b60405180910390a35b5050505050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace006001600160a01b038416610bac5781816002016000828254610ba19190611246565b90915550610c249050565b6001600160a01b03841660009081526020829052604090205482811015610c05578481846040517fe450d38c0000000000000000000000000000000000000000000000000000000081526004016104599392919061121e565b6001600160a01b03851660009081526020839052604090209083900390555b6001600160a01b038316610c42576002810180548390039055610c61565b6001600160a01b03831660009081526020829052604090208054830190555b826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610ca49190610e87565b60405180910390a350505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166103dc576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610d21610cb2565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace007f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace03610d6d84826112f4565b506004810161086783826112f4565b61076e610cb2565b60005b83811015610d9f578181015183820152602001610d87565b50506000910152565b6000610db2825190565b808452602084019350610dc9818560208601610d84565b601f01601f19169290920192915050565b602080825281016103c38184610da8565b60006001600160a01b03821661039e565b610e0581610deb565b81146107ba57600080fd5b803561039e81610dfc565b80610e05565b803561039e81610e1b565b60008060408385031215610e4257610e42600080fd5b6000610e4e8585610e10565b9250506020610e5f85828601610e21565b9150509250929050565b8015155b82525050565b6020810161039e8284610e69565b80610e6d565b6020810161039e8284610e81565b600080600060608486031215610ead57610ead600080fd5b6000610eb98686610e10565b9350506020610eca86828701610e10565b9250506040610edb86828701610e21565b9150509250925092565b60ff8116610e6d565b6020810161039e8284610ee5565b600060208284031215610f1157610f11600080fd5b6000610f1d8484610e10565b949350505050565b610e6d81610deb565b6020810161039e8284610f25565b60008083601f840112610f5157610f51600080fd5b50813567ffffffffffffffff811115610f6c57610f6c600080fd5b602083019150836020820283011115610f8757610f87600080fd5b9250929050565b60008060208385031215610fa457610fa4600080fd5b823567ffffffffffffffff811115610fbe57610fbe600080fd5b610fca85828601610f3c565b92509250509250929050565b60008060408385031215610fec57610fec600080fd5b6000610ff88585610e10565b9250506020610e5f85828601610e10565b634e487b7160e01b600052602260045260246000fd5b60028104600182168061103357607f821691505b60208210810361104557611045611009565b50919050565b602c8152602081017f43616c6c65723a2063616c6c6572206973206e6f74207468652064657369676e81527f6174656420616464726573730000000000000000000000000000000000000000602082015290505b60400190565b6020808252810161039e8161104b565b6020808252810161039e81601a81527f4e6f207472616e73616374696f6e7320746f20636f6e76657274000000000000602082015260400190565b634e487b7160e01b600052603260045260246000fd5b6000823561013e193684900301811261112157611121600080fd5b9190910192915050565b600061039e6111378381565b90565b610e6d8161112b565b604081016111518285610f25565b6103c3602083018461113a565b634e487b7160e01b600052601160045260246000fd5b6000600182016111865761118661115e565b5060010190565b60248152602081017f496e76616c6964207472616e73616374696f6e20616e616c797a65722061646481527f72657373000000000000000000000000000000000000000000000000000000006020820152905061109f565b6020808252810161039e8161118d565b600067ffffffffffffffff821661039e565b610e6d816111f5565b6020810161039e8284611207565b6060810161122c8286610f25565b6112396020830185610e81565b610f1d6040830184610e81565b8082018082111561039e5761039e61115e565b634e487b7160e01b600052604160045260246000fd5b6112788361112b565b815460001960089490940293841b1916921b91909117905550565b600061054f81848461126f565b81811015610a0d576112b3600082611293565b6001016112a0565b601f82111561054f576000818152602090206020601f850104810160208510156112e25750805b610b576020601f8601048301826112a0565b815167ffffffffffffffff81111561130e5761130e611259565b611318825461101f565b6113238282856112bb565b506020601f82116001811461135857600083156113405750848201515b600019600885021c1981166002850217855550610b57565b600084815260208120601f198516915b828110156113885787850151825560209485019460019092019101611368565b50848210156113a55783870151600019601f87166008021c191681555b5050505060020260010190555056fea2646970667358221220c92d36796669f09a165d85f5ad6ccaec7365a19d0cba54a86db26b29f35b60dd64736f6c63430008150033",
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
