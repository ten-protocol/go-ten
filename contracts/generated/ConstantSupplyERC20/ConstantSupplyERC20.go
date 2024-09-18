// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConstantSupplyERC20

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

// ConstantSupplyERC20MetaData contains all meta data concerning the ConstantSupplyERC20 contract.
var ConstantSupplyERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523462000031576200001f6200001862000183565b91620001ab565b6040516109f8620006de82396109f890f35b600080fd5b634e487b7160e01b600052604160045260246000fd5b90601f01601f191681019081106001600160401b038211176200006e57604052565b62000036565b906200008b6200008360405190565b92836200004c565b565b6001600160401b0381116200006e57602090601f01601f19160190565b60005b838110620000be5750506000910152565b8181015183820152602001620000ad565b90929192620000e8620000e2826200008d565b62000074565b9381855260208501908284011162000031576200008b92620000aa565b9080601f83011215620000315781516200012292602001620000cf565b90565b9091606082840312620000315781516001600160401b0381116200003157836200015191840162000105565b602083015190936001600160401b038211620000315762000179816200012293860162000105565b9360400190505190565b620001a6620010d6803803806200019a8162000074565b92833981019062000125565b909192565b6200008b9291620001bc91620003c4565b3362000426565b634e487b7160e01b600052602260045260246000fd5b9060016002830492168015620001fc575b6020831014620001f657565b620001c3565b91607f1691620001ea565b9160001960089290920291821b911b5b9181191691161790565b6200012262000122620001229290565b919062000246620001226200024f9362000221565b90835462000207565b9055565b6200008b9160009162000231565b8181106200026d575050565b806200027d600060019362000253565b0162000261565b9190601f81116200029457505050565b620002a86200008b93600052602060002090565b906020601f840181900483019310620002cc575b6020601f90910104019062000261565b9091508190620002bc565b90620002e1815190565b906001600160401b0382116200006e576200030982620003028554620001d9565b8562000284565b602090601f831160011462000348576200024f9291600091836200033c575b5050600019600883021c1916906002021790565b01519050388062000328565b601f198316916200035e85600052602060002090565b9260005b8181106200039f5750916002939185600196941062000385575b50505002019055565b01516000196008601f8516021c191690553880806200037c565b9193602060018192878701518155019501920162000362565b906200008b91620002d7565b90620003d66200008b926003620003b8565b6004620003b8565b620003ee62000122620001229290565b6001600160a01b031690565b6200012290620003de565b6200041090620003ee565b9052565b6020810192916200008b919062000405565b9190620004346000620003fa565b926200044084620003ee565b6200044b82620003ee565b146200045d576200008b929362000561565b62000483846200046c60405190565b63ec442f0560e01b81529182916004830162000414565b0390fd5b6200012290620003ee906001600160a01b031682565b620001229062000487565b62000122906200049d565b90620004bf90620004a8565b600052602052604060002090565b620001229081565b620001229054620004cd565b6040906200050c6200008b9496959396620005056060840198600085019062000405565b6020830152565b0152565b906000199062000217565b906200052f620001226200024f9262000221565b825462000510565b634e487b7160e01b600052601160045260246000fd5b919082018092116200055b57565b62000537565b8160006200056f81620003fa565b6200057a81620003ee565b6200058585620003ee565b036200066e57620005b990620005b3620005ab88620005a56002620004d5565b6200054d565b60026200051b565b620003ee565b620005c483620003ee565b036200063d575050620005e5620005ab84620005e16002620004d5565b0390565b6200063862000620620006197fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef93620004a8565b93620004a8565b936200062b60405190565b9182918290815260200190565b0390a3565b62000668916200064d91620004b3565b62000661856200065d83620004d5565b0190565b906200051b565b620005e5565b90915062000687620006818484620004b3565b620004d5565b858110620006b657849291620005b3620006a488620005b9940390565b620006b08786620004b3565b6200051b565b836200048387620006c660405190565b63391434e360e21b815293849360048501620004e156fe6080604052600436101561001257600080fd5b60003560e01c806306fdde03146100a2578063095ea7b31461009d57806318160ddd1461009857806323b872dd14610093578063313ce5671461008e57806370a082311461008957806395d89b4114610084578063a9059cbb1461007f5763dd62ed3e036100b257610312565b6102d3565b6102b8565b61029d565b61025a565b61023e565b6101e2565b6101b4565b610126565b60009103126100b257565b600080fd5b60005b8381106100ca5750506000910152565b81810151838201526020016100ba565b6100fb61010460209361010e936100ef815190565b80835293849260200190565b958691016100b7565b601f01601f191690565b0190565b6020808252610123929101906100da565b90565b346100b2576101363660046100a7565b61014d610141610463565b60405191829182610112565b0390f35b6001600160a01b031690565b6001600160a01b0381165b036100b257565b9050359061017c8261015d565b565b80610168565b9050359061017c8261017e565b91906040838203126100b257610123906101ab818561016f565b93602001610184565b346100b25761014d6101d06101ca366004610191565b90610540565b60405191829182901515815260200190565b346100b2576101f23660046100a7565b61014d6101fd6104a5565b6040515b9182918290815260200190565b90916060828403126100b257610123610227848461016f565b93610235816020860161016f565b93604001610184565b346100b25761014d6101d061025436600461020e565b9161054b565b346100b25761026a3660046100a7565b61014d61027561048a565b6040519182918260ff909116815260200190565b906020828203126100b2576101239161016f565b346100b25761014d6101fd6102b3366004610289565b6104ef565b346100b2576102c83660046100a7565b61014d61014161046d565b346100b25761014d6101d06102e9366004610191565b9061050b565b91906040838203126100b25761012390610309818561016f565b9360200161016f565b346100b25761014d6101fd6103283660046102ef565b90610520565b634e487b7160e01b600052602260045260246000fd5b9060016002830492168015610364575b602083101461035f57565b61032e565b91607f1691610354565b8054600093929161038b61038183610344565b8085529360200190565b91600181169081156103dd57506001146103a457505050565b6103b79192939450600052602060002090565b916000925b8184106103c95750500190565b8054848401526020909301926001016103bc565b92949550505060ff1916825215156020020190565b906101239161036e565b634e487b7160e01b600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff82111761043457604052565b6103fc565b9061017c6104539261044a60405190565b938480926103f2565b0383610412565b61012390610439565b610123600361045a565b610123600461045a565b6104846101236101239290565b60ff1690565b6101236012610477565b6101239081565b6101239054610494565b610123600261049b565b610151610123610123926001600160a01b031690565b610123906104af565b610123906104c5565b906104e1906104ce565b600052602052604060002090565b610506610123916104fe600090565b5060006104d7565b61049b565b61051b919033610572565b610572565b600190565b6101239161053b61050692610533600090565b5060016104d7565b6104d7565b61051b919033610814565b61051b92919061051683338361093e565b6101516101236101239290565b6101239061055c565b92919061057f6000610569565b936001600160a01b0385166001600160a01b0382161461060b576001600160a01b0385166001600160a01b038316146105bd5761017c9394506106c9565b610607856105ca60405190565b9182917fec442f05000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b0390fd5b6106078561061860405190565b9182917f96c6fd1e000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6001600160a01b03909116815260608101939261017c92909160409161067c906020830152565b0152565b6101236101236101239290565b9061012361012361069d92610680565b9055565b634e487b7160e01b600052601160045260246000fd5b919082018092116106c457565b6106a1565b8160006106d581610569565b6001600160a01b0381166001600160a01b0385160361079f5761070f9061015161070888610703600261049b565b6106b7565b600261068d565b6001600160a01b0383160361077a57505061073561070884610731600261049b565b0390565b61077561076b6107657fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef936104ce565b936104ce565b9361020160405190565b0390a3565b61079a91610787916104d7565b6107948561010e8361049b565b9061068d565b610735565b9091506107af61050684846104d7565b8581106107d7578492916101516107c88861070f940390565b6107d287866104d7565b61068d565b83610607876107e560405190565b9384937fe450d38c00000000000000000000000000000000000000000000000000000000855260048501610655565b909161017c926001925b90919261082b6000610569565b6001600160a01b0381166001600160a01b038416146108f4576001600160a01b0381166001600160a01b038516146108aa5750610871846107d28561053b8660016104d7565b61087a57505050565b61077561076b6107657f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925936104ce565b610607906108b760405190565b9182917f94280d62000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6106079061090160405190565b9182917fe602df05000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b9092919261094c8183610520565b93600019850361095e575b5050509050565b808510610984576109729061097b94950390565b9060009261081e565b80388080610957565b90610607859261099360405190565b9384937ffb8f41b20000000000000000000000000000000000000000000000000000000085526004850161065556fea26469706673582212206c5fbc4c274c1d9b437e400eeb4b616e6aa1bcff9af56e0a9b15e71aada6082664736f6c63430008140033",
}

// ConstantSupplyERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ConstantSupplyERC20MetaData.ABI instead.
var ConstantSupplyERC20ABI = ConstantSupplyERC20MetaData.ABI

// ConstantSupplyERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConstantSupplyERC20MetaData.Bin instead.
var ConstantSupplyERC20Bin = ConstantSupplyERC20MetaData.Bin

// DeployConstantSupplyERC20 deploys a new Ethereum contract, binding an instance of ConstantSupplyERC20 to it.
func DeployConstantSupplyERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int) (common.Address, *types.Transaction, *ConstantSupplyERC20, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConstantSupplyERC20Bin), backend, name, symbol, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// ConstantSupplyERC20 is an auto generated Go binding around an Ethereum contract.
type ConstantSupplyERC20 struct {
	ConstantSupplyERC20Caller     // Read-only binding to the contract
	ConstantSupplyERC20Transactor // Write-only binding to the contract
	ConstantSupplyERC20Filterer   // Log filterer for contract events
}

// ConstantSupplyERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConstantSupplyERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstantSupplyERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConstantSupplyERC20Session struct {
	Contract     *ConstantSupplyERC20 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConstantSupplyERC20CallerSession struct {
	Contract *ConstantSupplyERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ConstantSupplyERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConstantSupplyERC20TransactorSession struct {
	Contract     *ConstantSupplyERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ConstantSupplyERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ConstantSupplyERC20Raw struct {
	Contract *ConstantSupplyERC20 // Generic contract binding to access the raw methods on
}

// ConstantSupplyERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConstantSupplyERC20CallerRaw struct {
	Contract *ConstantSupplyERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ConstantSupplyERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConstantSupplyERC20TransactorRaw struct {
	Contract *ConstantSupplyERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewConstantSupplyERC20 creates a new instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20(address common.Address, backend bind.ContractBackend) (*ConstantSupplyERC20, error) {
	contract, err := bindConstantSupplyERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20{ConstantSupplyERC20Caller: ConstantSupplyERC20Caller{contract: contract}, ConstantSupplyERC20Transactor: ConstantSupplyERC20Transactor{contract: contract}, ConstantSupplyERC20Filterer: ConstantSupplyERC20Filterer{contract: contract}}, nil
}

// NewConstantSupplyERC20Caller creates a new read-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Caller(address common.Address, caller bind.ContractCaller) (*ConstantSupplyERC20Caller, error) {
	contract, err := bindConstantSupplyERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Caller{contract: contract}, nil
}

// NewConstantSupplyERC20Transactor creates a new write-only instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ConstantSupplyERC20Transactor, error) {
	contract, err := bindConstantSupplyERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Transactor{contract: contract}, nil
}

// NewConstantSupplyERC20Filterer creates a new log filterer instance of ConstantSupplyERC20, bound to a specific deployed contract.
func NewConstantSupplyERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ConstantSupplyERC20Filterer, error) {
	contract, err := bindConstantSupplyERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20Filterer{contract: contract}, nil
}

// bindConstantSupplyERC20 binds a generic wrapper to an already deployed contract.
func bindConstantSupplyERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConstantSupplyERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.ConstantSupplyERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstantSupplyERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.Allowance(&_ConstantSupplyERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.BalanceOf(&_ConstantSupplyERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Decimals() (uint8, error) {
	return _ConstantSupplyERC20.Contract.Decimals(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Name() (string, error) {
	return _ConstantSupplyERC20.Contract.Name(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) Symbol() (string, error) {
	return _ConstantSupplyERC20.Contract.Symbol(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConstantSupplyERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ConstantSupplyERC20 *ConstantSupplyERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ConstantSupplyERC20.Contract.TotalSupply(&_ConstantSupplyERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Approve(&_ConstantSupplyERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.Transfer(&_ConstantSupplyERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ConstantSupplyERC20 *ConstantSupplyERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ConstantSupplyERC20.Contract.TransferFrom(&_ConstantSupplyERC20.TransactOpts, from, to, value)
}

// ConstantSupplyERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20ApprovalIterator struct {
	Event *ConstantSupplyERC20Approval // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Approval)
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
		it.Event = new(ConstantSupplyERC20Approval)
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
func (it *ConstantSupplyERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Approval represents a Approval event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ConstantSupplyERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20ApprovalIterator{contract: _ConstantSupplyERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Approval)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseApproval(log types.Log) (*ConstantSupplyERC20Approval, error) {
	event := new(ConstantSupplyERC20Approval)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConstantSupplyERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20TransferIterator struct {
	Event *ConstantSupplyERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ConstantSupplyERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConstantSupplyERC20Transfer)
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
		it.Event = new(ConstantSupplyERC20Transfer)
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
func (it *ConstantSupplyERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConstantSupplyERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConstantSupplyERC20Transfer represents a Transfer event raised by the ConstantSupplyERC20 contract.
type ConstantSupplyERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ConstantSupplyERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConstantSupplyERC20TransferIterator{contract: _ConstantSupplyERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ConstantSupplyERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConstantSupplyERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConstantSupplyERC20Transfer)
				if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ConstantSupplyERC20 *ConstantSupplyERC20Filterer) ParseTransfer(log types.Log) (*ConstantSupplyERC20Transfer, error) {
	event := new(ConstantSupplyERC20Transfer)
	if err := _ConstantSupplyERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
