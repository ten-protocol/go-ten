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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProxyDeniedAdminAccess\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	Bin: "0x60a060405261001561000f610162565b916101b6565b6040516105a9610571823960805181602c01526105a990f35b634e487b7160e01b5f52604160045260245ffd5b90601f01601f191681019081106001600160401b0382111761006357604052565b61002e565b9061007c61007560405190565b9283610042565b565b6001600160a01b031690565b90565b6001600160a01b0381160361009e57565b5f80fd5b9050519061007c8261008d565b6001600160401b03811161006357602090601f01601f19160190565b90825f9392825e0152565b909291926100eb6100e6826100af565b610068565b938185528183011161009e5761007c9160208501906100cb565b9080601f8301121561009e57815161008a926020016100d6565b9160608383031261009e5761013482846100a2565b9261014283602083016100a2565b60408201519093906001600160401b03811161009e5761008a9201610105565b6101806113488038038061017581610068565b92833981019061011f565b909192565b6040513d5f823e3d90fd5b61008a9061007e906001600160a01b031682565b61008a90610190565b61008a906101a4565b916101c19192610232565b6040519061082e82016001600160401b038111838210176100635782916101fd9161082e610b1a85396001600160a01b03909116815260200190565b03905ff0801561022d57610210906101ad565b60805261007c6102286080516001600160a01b031690565b61023c565b610185565b9061007c916102a9565b61007c907f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f8161026a610353565b61029461027660405190565b928392836001600160a01b0391821681529116602082015260400190565b0390a16103ae565b61008a61008a61008a9290565b906102b382610430565b6102bc826101ad565b7fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b6102e660405190565b5f90a280516102fb6102f75f61029c565b9190565b111561030d5761030a916104c5565b50565b505061007c610473565b61008a7fb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d610361029c565b61008a9061007e565b61008a9054610340565b61008a5f61036261008a610317565b01610349565b61007e61008a61008a9290565b61008a90610368565b9061038e61008a6103aa926101ad565b82546001600160a01b0319166001600160a01b03919091161790565b9055565b6103b75f610375565b6001600160a01b0381166001600160a01b038316146103e7575061007c905f6103e161008a610317565b0161037e565b633173bdd160e11b5f9081526001600160a01b039091166004526024035ffd5b61008a7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc61029c565b803b61043e6102f75f61029c565b146104535761007c905f6103e161008a610407565b634c9c8ce360e01b5f9081526001600160a01b039091166004526024035ffd5b61047c5f61029c565b341161048457565b63b398979f60e01b5f908152600490fd5b906104a26100e6836100af565b918252565b3d156104c0576104b63d610495565b903d5f602084013e565b606090565b5f8061008a936104d3606090565b50602081519101845af46104e56104a7565b91906104f15750610543565b81516104ff6102f75f61029c565b148061052e575b61050e575090565b639996b31560e01b5f9081526001600160a01b039091166004526024035ffd5b50803b61053d6102f75f61029c565b14610506565b80516105516102f75f61029c565b111561055f57805190602001fd5b63d6bda27560e01b5f908152600490fdfe608060405261000c61000e565b005b610016610027565b565b6001600160a01b031690565b90565b6100507f0000000000000000000000000000000000000000000000000000000000000000610018565b33036100d4577f4f1ef286000000000000000000000000000000000000000000000000000000005f357fffffffff0000000000000000000000000000000000000000000000000000000016146100cc577fd2b576ec000000000000000000000000000000000000000000000000000000005f90815260045b035ffd5b610016610260565b6100dc61029a565b6102a2565b6100246100246100249290565b90939293848311610106578411610106578101920390565b5f80fd5b6001600160a01b0381160361010657565b905035906100168261010a565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761017757604052565b610128565b9061001661018960405190565b9283610155565b67ffffffffffffffff811161017757602090601f01601f19160190565b90825f939282370152565b909291926101cd6101c882610190565b61017c565b9381855281830111610106576100169160208501906101ad565b9080601f8301121561010657816020610024933591016101b8565b91909160408184031261010657610219838261011b565b92602082013567ffffffffffffffff81116101065761002492016101e7565b610018610024610024926001600160a01b031690565b61002490610238565b6100249061024e565b61001661029561028e61028661028261027960046100e1565b3690365f6100ee565b9091565b810190610202565b9190610257565b6102bf565b610024610369565b5f8091368280378136915af43d5f803e156102bb573d5ff35b3d5ffd5b906102c9826103d9565b6102d282610257565b7fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b6102fc60405190565b5f90a2805161031161030d5f6100e1565b9190565b1115610323576103209161049e565b50565b5050610016610433565b6100247f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc6100e1565b61002490610018565b6100249054610356565b6100245f61037861002461032d565b0161035f565b6001600160a01b03909116815260200190565b906103a16100246103d592610257565b82547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03919091161790565b9055565b803b6103e761030d5f6100e1565b1461040257610016905f6103fc61002461032d565b01610391565b6100c85f917f4c9c8ce30000000000000000000000000000000000000000000000000000000083526004830161037e565b61043c5f6100e1565b341161044457565b7fb398979f000000000000000000000000000000000000000000000000000000005f908152600490fd5b9061047b6101c883610190565b918252565b3d156104995761048f3d61046e565b903d5f602084013e565b606090565b5f80610024936104ac606090565b50602081519101845af46104be610480565b91906104ca575061052d565b81516104d861030d5f6100e1565b1480610518575b6104e7575090565b6100c85f917f9996b3150000000000000000000000000000000000000000000000000000000083526004830161037e565b50803b61052761030d5f6100e1565b146104df565b805161053b61030d5f6100e1565b111561054957805190602001fd5b7fd6bda275000000000000000000000000000000000000000000000000000000005f908152600490fdfea26469706673582212209d6c1feecdea0c4c8fa7c7b1c4c715a40dec2fb8af67b60cce14bc3e5a76e88964736f6c634300081c003360806040523461002a576100196100146100bf565b6100dd565b604051610619610215823961061990f35b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b90601f01601f191681019081106001600160401b0382111761006357604052565b61002e565b9061007c61007560405190565b9283610042565b565b6001600160a01b031690565b90565b6001600160a01b0381160361002a57565b9050519061007c8261008d565b9060208282031261002a5761008a9161009e565b61008a61082e803803806100d281610068565b9283398101906100ab565b61007c9061010f565b61007e61008a61008a9290565b61008a906100e6565b6001600160a01b03909116815260200190565b6101185f6100f3565b6001600160a01b0381166001600160a01b0383161461013b575061007c906101c0565b631e4fbdf760e01b5f908152906101539060046100fc565b035ffd5b61008a9061007e565b61008a9054610157565b61008a9061007e906001600160a01b031682565b61008a9061016a565b61008a9061017e565b906101a061008a6101bc92610187565b82546001600160a01b0319166001600160a01b03919091161790565b9055565b6101e46101de6101cf5f610160565b6101d9845f610190565b610187565b91610187565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e061020f60405190565b5f90a356fe60806040526004361015610011575f80fd5b5f3560e01c8063715018a6146100605780638da5cb5b1461005b5780639623609d14610056578063ad3cb1cc146100515763f2fde38b0361006f5761033b565b610300565b61023f565b6100a6565b610073565b5f91031261006f57565b5f80fd5b3461006f57610083366004610065565b61008b61038a565b60405180805b0390f35b6001600160a01b031690565b90565b565b3461006f576100b6366004610065565b6100916100ca5f546001600160a01b031690565b604051918291826001600160a01b03909116815260200190565b6001600160a01b0381165b0361006f57565b905035906100a4826100e4565b6001600160a01b0381166100ef565b905035906100a482610103565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff82111761016e57604052565b61011f565b906100a461018060405190565b928361014c565b67ffffffffffffffff811161016e57602090601f01601f19160190565b0190565b90825f939282370152565b909291926101c86101c382610187565b610173565b938185528183011161006f576100a49160208501906101a8565b9080601f8301121561006f578160206100a1933591016101b3565b9160608383031261006f5761021282846100f6565b926102208360208301610112565b92604082013567ffffffffffffffff811161006f576100a192016101e2565b61008b61024d3660046101fd565b9161046a565b906102606101c383610187565b918252565b61026f6005610253565b7f352e302e30000000000000000000000000000000000000000000000000000000602082015290565b6100a1610265565b6100a1610298565b6100a16102a0565b90825f9392825e0152565b6102dc6102e56020936101a4936102d0815190565b80835293849260200190565b958691016102b0565b601f01601f191690565b60208082526100a1929101906102bb565b3461006f57610310366004610065565b61009161031b6102a8565b604051918291826102ef565b9060208282031261006f576100a191610112565b3461006f5761008b61034e366004610327565b6104e6565b61035b6104ef565b6100a4610379565b6100956100a16100a19290565b6100a190610363565b6100a46103855f610370565b610587565b6100a4610353565b906100a492916103a06104ef565b6103fb565b6100956100a16100a1926001600160a01b031690565b6100a1906103a5565b6100a1906103bb565b6001600160a01b0390911681526100a191604082019160208184039101526102bb565b6040513d5f823e3d90fd5b610404906103c4565b90634f1ef28691803b1561006f576104305f9361043b61042360405190565b9687958694859460e01b90565b8452600484016103cd565b039134905af180156104655761044e5750565b6100a4905f61045d818361014c565b810190610065565b6103f0565b906100a49291610392565b6100a4906104816104ef565b61048a5f610370565b6001600160a01b0381166001600160a01b038316146104ad57506100a490610587565b7f1e4fbdf7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b6100a490610475565b5f5433906001600160a01b03168190036105065750565b7f118cdaa7000000000000000000000000000000000000000000000000000000005f9081526001600160a01b039091166004526024035ffd5b9061054f6100a1610583926103c4565b82547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03919091161790565b9055565b6105b36105ad61059e5f546001600160a01b031690565b6105a8845f61053f565b6103c4565b916103c4565b907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e06105de60405190565b5f90a356fea264697066735822122001967bab5efcf570bf58012072fe784bb3fac3042654f37ac70d1aa3338fbcea64736f6c634300081c0033",
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
