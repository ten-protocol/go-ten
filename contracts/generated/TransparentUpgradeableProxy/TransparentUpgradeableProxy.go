// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TransparentUpgradeableProxy

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

// TransparentUpgradeableProxyMetaData contains all meta data concerning the TransparentUpgradeableProxy contract.
var TransparentUpgradeableProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProxyDeniedAdminAccess\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	Bin: "0x60a060405262000019620000126200019f565b91620001ff565b6040516105e762000681823960805181603001526105e790f35b634e487b7160e01b600052604160045260246000fd5b90601f01601f191681019081106001600160401b038211176200006b57604052565b62000033565b90620000886200008060405190565b928362000049565b565b6001600160a01b031690565b90565b6001600160a01b03811603620000ab57565b600080fd5b90505190620000888262000099565b6001600160401b0381116200006b57602090601f01601f19160190565b60005b838110620000f05750506000910152565b8181015183820152602001620000df565b909291926200011a6200011482620000bf565b62000071565b93818552602085019082840111620000ab576200008892620000dc565b9080601f83011215620000ab578151620000969260200162000101565b91606083830312620000ab576200016c8284620000b0565b926200017c8360208301620000b0565b60408201519093906001600160401b038111620000ab5762000096920162000137565b620001c26200153080380380620001b68162000071565b92833981019062000154565b909192565b6040513d6000823e3d90fd5b62000096906200008a906001600160a01b031682565b6200009690620001d3565b6200009690620001e9565b916200020c919262000288565b604051906108c882016001600160401b038111838210176200006b5782916200024b916108c862000c6885396001600160a01b03909116815260200190565b03906000f0801562000282576200026290620001f4565b608052620000886200027c6080516001600160a01b031690565b62000406565b620001c7565b906200008891620002a4565b6200009662000096620000969290565b90620002b08262000381565b7fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b620002dc83620001f4565b90620002e760405190565b600090a2805162000301620002fd600062000294565b9190565b111562000316576200031391620005a5565b50565b50506200008862000543565b620000967f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc62000294565b9062000361620000966200037d92620001f4565b82546001600160a01b0319166001600160a01b03919091161790565b9055565b6000813b62000394620002fd8362000294565b14620003b557906200008891620003ae6200009662000322565b016200034d565b620003e882620003c460405190565b634c9c8ce360e01b8152918291600483016001600160a01b03909116815260200190565b0390fd5b6001600160a01b0391821681529116602082015260400190565b6200008890620004156200049b565b817f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f91620004506200044660405190565b92839283620003ec565b0390a1620004d1565b620000967fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d610362000294565b62000096906200008a565b62000096905462000484565b620000966000620004af6200009662000459565b016200048f565b6200008a62000096620000969290565b6200009690620004b6565b90600091620004e083620004c6565b926001600160a01b0384166001600160a01b03831614620005105762000088929350620003ae6200009662000459565b620003e8846200051f60405190565b633173bdd160e11b8152918291600483016001600160a01b03909116815260200190565b3462000554620002fd600062000294565b116200055c57565b60405163b398979f60e01b8152600490fd5b906200057e6200011483620000bf565b918252565b3d15620005a057620005953d6200056e565b903d6000602084013e565b606090565b6000806200009693620005b6606090565b50805190602001845af4620005ca62000583565b9190620005d857506200064d565b8151600090620005ec620002fd8362000294565b14908162000632575b50620005ff575090565b620003e8906200060e60405190565b639996b31560e01b8152918291600483016001600160a01b03909116815260200190565b905062000645620002fd833b9262000294565b1438620005f5565b80516200065f620002fd600062000294565b11156200066e57805190602001fd5b604051630a12f52160e11b8152600490fdfe608060405261000c61000e565b005b610016610027565b565b6001600160a01b031690565b90565b336100616100547f0000000000000000000000000000000000000000000000000000000000000000610018565b916001600160a01b031690565b036100e8577f4f1ef286000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000060003516146100e0576040517fd2b576ec000000000000000000000000000000000000000000000000000000008152600490fd5b0390fd5b6100166102ed565b6100f06100f5565b610156565b610024610140565b6100246100246100249290565b6100247f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6100fd565b546001600160a01b031690565b610024600061015061002461010a565b01610133565b60008091368280378136915af43d6000803e15610172573d6000f35b3d6000fd5b9093929384831161018f57841161018f578101920390565b600080fd5b6001600160a01b0381160361018f57565b9050359061001682610194565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff82111761020357604052565b6101b2565b9061001661021560405190565b92836101e1565b67ffffffffffffffff811161020357602090601f01601f19160190565b90826000939282370152565b9092919261025a6102558261021c565b610208565b9381855260208501908284011161018f5761001692610239565b9080601f8301121561018f5781602061002493359101610245565b91909160408184031261018f576102a683826101a5565b92602082013567ffffffffffffffff811161018f576100249201610274565b610018610024610024926001600160a01b031690565b610024906102c5565b610024906102db565b61001661032361031c6103146103103660008161030a60046100fd565b91610177565b9091565b81019061028f565b91906102e4565b9061032d826103ef565b7fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b610357836102e4565b9061036160405190565b600090a2805161037861037460006100fd565b9190565b111561038a57610387916104c7565b50565b5050610016610456565b6001600160a01b03909116815260200190565b906103b76100246103eb926102e4565b82547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03919091161790565b9055565b6000813b6103ff610374836100fd565b1461041a57906100169161041461002461010a565b016103a7565b6100dc8261042760405190565b9182917f4c9c8ce300000000000000000000000000000000000000000000000000000000835260048301610394565b3461046461037460006100fd565b1161046b57565b6040517fb398979f000000000000000000000000000000000000000000000000000000008152600490fd5b906104a36102558361021c565b918252565b3d156104c2576104b73d610496565b903d6000602084013e565b606090565b600080610024936104d6606090565b50805190602001845af46104e86104a8565b91906104f45750610569565b8151600090610505610374836100fd565b149081610552575b50610516575090565b6100dc9061052360405190565b9182917f9996b31500000000000000000000000000000000000000000000000000000000835260048301610394565b9050610562610374833b926100fd565b143861050d565b805161057861037460006100fd565b111561058657805190602001fd5b6040517f1425ea42000000000000000000000000000000000000000000000000000000008152600490fdfea2646970667358221220782f2801dba6b92d163f770087c095f4a3bcf48f3a41dbeafc28edfa481d940f64736f6c6343000814003360806040523462000030576200001e62000018620000d3565b620000f6565b60405161066362000265823961066390f35b600080fd5b634e487b7160e01b600052604160045260246000fd5b90601f01601f191681019081106001600160401b038211176200006d57604052565b62000035565b906200008a6200008260405190565b92836200004b565b565b6001600160a01b031690565b90565b6001600160a01b038116036200003057565b905051906200008a826200009b565b9060208282031262000030576200009891620000ad565b62000098620008c880380380620000ea8162000073565b928339810190620000bc565b6200008a906200012f565b6200008c62000098620000989290565b620000989062000101565b6001600160a01b03909116815260200190565b6200013b600062000111565b6001600160a01b0381166001600160a01b038316146200016157506200008a9062000202565b62000187906200017060405190565b631e4fbdf760e01b8152918291600483016200011c565b0390fd5b62000098906200008c565b6200009890546200018b565b62000098906200008c906001600160a01b031682565b6200009890620001a2565b6200009890620001b8565b90620001e262000098620001fe92620001c3565b82546001600160a01b0319166001600160a01b03919091161790565b9055565b6200020e600062000196565b906200021c816000620001ce565b620002536200024c7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e093620001c3565b91620001c3565b916200025e60405190565b600090a356fe6080604052600436101561001257600080fd5b60003560e01c8063715018a6146100625780638da5cb5b1461005d5780639623609d14610058578063ad3cb1cc146100535763f2fde38b036100725761035a565b61031f565b610246565b6100a9565b610077565b600091031261007257565b600080fd5b3461007257610087366004610067565b61008f6103aa565b604051005b0390f35b6001600160a01b031690565b90565b565b34610072576100b9366004610067565b6100946100ce6000546001600160a01b031690565b604051918291826001600160a01b03909116815260200190565b6001600160a01b0381165b0361007257565b905035906100a7826100e8565b6001600160a01b0381166100f3565b905035906100a782610107565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff82111761017457604052565b610123565b906100a761018660405190565b9283610152565b67ffffffffffffffff811161017457602090601f01601f19160190565b0190565b90826000939282370152565b909291926101cf6101ca8261018d565b610179565b93818552602085019082840111610072576100a7926101ae565b9080601f83011215610072578160206100a4933591016101ba565b916060838303126100725761021982846100fa565b926102278360208301610116565b92604082013567ffffffffffffffff8111610072576100a492016101e9565b61008f610254366004610204565b91610622565b906102676101ca8361018d565b918252565b610276600561025a565b7f352e302e30000000000000000000000000000000000000000000000000000000602082015290565b6100a461026c565b6100a461029f565b6100a46102a7565b60005b8381106102ca5750506000910152565b81810151838201526020016102ba565b6102fb6103046020936101aa936102ef815190565b80835293849260200190565b958691016102b7565b601f01601f191690565b60208082526100a4929101906102da565b346100725761032f366004610067565b61009461033a6102af565b6040519182918261030e565b90602082820312610072576100a491610116565b346100725761008f61036d366004610346565b61049b565b61037a6103b2565b6100a7610398565b6100986100a46100a49290565b6100a490610382565b6100a76103a5600061038f565b610514565b6100a7610372565b60005433906001600160a01b03168190036103ca5750565b610414906103d760405190565b9182917f118cdaa7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b0390fd5b6100a7906104246103b2565b61042e600061038f565b6001600160a01b0381166001600160a01b0383161461045157506100a790610514565b6104149061045e60405190565b9182917f1e4fbdf7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b6100a790610418565b6100986100a46100a4926001600160a01b031690565b6100a4906104a4565b6100a4906104ba565b906104dc6100a4610510926104c3565b82547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03919091161790565b9055565b6000546001600160a01b03169061052c8160006104cc565b61055f6105597f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936104c3565b916104c3565b9161056960405190565b600090a3565b906100a7929161057d6103b2565b6105b1565b6001600160a01b0390911681526100a491604082019160208184039101526102da565b6040513d6000823e3d90fd5b6105ba906104c3565b349190634f1ef286813b15610072576000936105e9916105f46105dc60405190565b9788968795869460e01b90565b845260048401610582565b03925af1801561061d576106055750565b6100a79060006106158183610152565b810190610067565b6105a5565b906100a7929161056f56fea26469706673582212205d5f75a536a5fd2f2e441550c1d5a4e715b6cff79e1cc29a86c7abd4b441f65364736f6c63430008140033",
}

// TransparentUpgradeableProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use TransparentUpgradeableProxyMetaData.ABI instead.
var TransparentUpgradeableProxyABI = TransparentUpgradeableProxyMetaData.ABI

// TransparentUpgradeableProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransparentUpgradeableProxyMetaData.Bin instead.
var TransparentUpgradeableProxyBin = TransparentUpgradeableProxyMetaData.Bin

// DeployTransparentUpgradeableProxy deploys a new Ethereum contract, binding an instance of TransparentUpgradeableProxy to it.
func DeployTransparentUpgradeableProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, initialOwner common.Address, _data []byte) (common.Address, *types.Transaction, *TransparentUpgradeableProxy, error) {
	parsed, err := TransparentUpgradeableProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransparentUpgradeableProxyBin), backend, _logic, initialOwner, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransparentUpgradeableProxy{TransparentUpgradeableProxyCaller: TransparentUpgradeableProxyCaller{contract: contract}, TransparentUpgradeableProxyTransactor: TransparentUpgradeableProxyTransactor{contract: contract}, TransparentUpgradeableProxyFilterer: TransparentUpgradeableProxyFilterer{contract: contract}}, nil
}

// TransparentUpgradeableProxy is an auto generated Go binding around an Ethereum contract.
type TransparentUpgradeableProxy struct {
	TransparentUpgradeableProxyCaller     // Read-only binding to the contract
	TransparentUpgradeableProxyTransactor // Write-only binding to the contract
	TransparentUpgradeableProxyFilterer   // Log filterer for contract events
}

// TransparentUpgradeableProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransparentUpgradeableProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransparentUpgradeableProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransparentUpgradeableProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransparentUpgradeableProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransparentUpgradeableProxySession struct {
	Contract     *TransparentUpgradeableProxy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// TransparentUpgradeableProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransparentUpgradeableProxyCallerSession struct {
	Contract *TransparentUpgradeableProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// TransparentUpgradeableProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransparentUpgradeableProxyTransactorSession struct {
	Contract     *TransparentUpgradeableProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// TransparentUpgradeableProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransparentUpgradeableProxyRaw struct {
	Contract *TransparentUpgradeableProxy // Generic contract binding to access the raw methods on
}

// TransparentUpgradeableProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyCallerRaw struct {
	Contract *TransparentUpgradeableProxyCaller // Generic read-only contract binding to access the raw methods on
}

// TransparentUpgradeableProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransparentUpgradeableProxyTransactorRaw struct {
	Contract *TransparentUpgradeableProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransparentUpgradeableProxy creates a new instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxy(address common.Address, backend bind.ContractBackend) (*TransparentUpgradeableProxy, error) {
	contract, err := bindTransparentUpgradeableProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxy{TransparentUpgradeableProxyCaller: TransparentUpgradeableProxyCaller{contract: contract}, TransparentUpgradeableProxyTransactor: TransparentUpgradeableProxyTransactor{contract: contract}, TransparentUpgradeableProxyFilterer: TransparentUpgradeableProxyFilterer{contract: contract}}, nil
}

// NewTransparentUpgradeableProxyCaller creates a new read-only instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxyCaller(address common.Address, caller bind.ContractCaller) (*TransparentUpgradeableProxyCaller, error) {
	contract, err := bindTransparentUpgradeableProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyCaller{contract: contract}, nil
}

// NewTransparentUpgradeableProxyTransactor creates a new write-only instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*TransparentUpgradeableProxyTransactor, error) {
	contract, err := bindTransparentUpgradeableProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyTransactor{contract: contract}, nil
}

// NewTransparentUpgradeableProxyFilterer creates a new log filterer instance of TransparentUpgradeableProxy, bound to a specific deployed contract.
func NewTransparentUpgradeableProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*TransparentUpgradeableProxyFilterer, error) {
	contract, err := bindTransparentUpgradeableProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyFilterer{contract: contract}, nil
}

// bindTransparentUpgradeableProxy binds a generic wrapper to an already deployed contract.
func bindTransparentUpgradeableProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransparentUpgradeableProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransparentUpgradeableProxy.Contract.TransparentUpgradeableProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.TransparentUpgradeableProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.TransparentUpgradeableProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransparentUpgradeableProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.Fallback(&_TransparentUpgradeableProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TransparentUpgradeableProxy.Contract.Fallback(&_TransparentUpgradeableProxy.TransactOpts, calldata)
}

// TransparentUpgradeableProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyAdminChangedIterator struct {
	Event *TransparentUpgradeableProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *TransparentUpgradeableProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransparentUpgradeableProxyAdminChanged)
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
		it.Event = new(TransparentUpgradeableProxyAdminChanged)
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
func (it *TransparentUpgradeableProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransparentUpgradeableProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransparentUpgradeableProxyAdminChanged represents a AdminChanged event raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*TransparentUpgradeableProxyAdminChangedIterator, error) {

	logs, sub, err := _TransparentUpgradeableProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyAdminChangedIterator{contract: _TransparentUpgradeableProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *TransparentUpgradeableProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _TransparentUpgradeableProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransparentUpgradeableProxyAdminChanged)
				if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) ParseAdminChanged(log types.Log) (*TransparentUpgradeableProxyAdminChanged, error) {
	event := new(TransparentUpgradeableProxyAdminChanged)
	if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransparentUpgradeableProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyUpgradedIterator struct {
	Event *TransparentUpgradeableProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *TransparentUpgradeableProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransparentUpgradeableProxyUpgraded)
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
		it.Event = new(TransparentUpgradeableProxyUpgraded)
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
func (it *TransparentUpgradeableProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransparentUpgradeableProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransparentUpgradeableProxyUpgraded represents a Upgraded event raised by the TransparentUpgradeableProxy contract.
type TransparentUpgradeableProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TransparentUpgradeableProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TransparentUpgradeableProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TransparentUpgradeableProxyUpgradedIterator{contract: _TransparentUpgradeableProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TransparentUpgradeableProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TransparentUpgradeableProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransparentUpgradeableProxyUpgraded)
				if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TransparentUpgradeableProxy *TransparentUpgradeableProxyFilterer) ParseUpgraded(log types.Log) (*TransparentUpgradeableProxyUpgraded, error) {
	event := new(TransparentUpgradeableProxyUpgraded)
	if err := _TransparentUpgradeableProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
