// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package NetworkConfig

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

// NetworkConfigAddresses is an auto generated low-level Go binding around an user-defined struct.
type NetworkConfigAddresses struct {
	CrossChain             common.Address
	MessageBus             common.Address
	NetworkEnclaveRegistry common.Address
	RollupContract         common.Address
	AdditionalContracts    []NetworkConfigNamedAddress
}

// NetworkConfigFixedAddresses is an auto generated low-level Go binding around an user-defined struct.
type NetworkConfigFixedAddresses struct {
	CrossChain             common.Address
	MessageBus             common.Address
	NetworkEnclaveRegistry common.Address
	RollupContract         common.Address
}

// NetworkConfigNamedAddress is an auto generated low-level Go binding around an user-defined struct.
type NetworkConfigNamedAddress struct {
	Name string
	Addr common.Address
}

// NetworkConfigMetaData contains all meta data concerning the NetworkConfig contract.
var NetworkConfigMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NetworkContractAddressAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MESSAGE_BUS_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"NETWORK_ENCLAVE_REGISTRY_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROLLUP_CONTRACT_SLOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"additionalAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addresses\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.NamedAddress[]\",\"name\":\"additionalContracts\",\"type\":\"tuple[]\"}],\"internalType\":\"structNetworkConfig.Addresses\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crossChainContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"crossChain\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"messageBus\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"networkEnclaveRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rollupContract\",\"type\":\"address\"}],\"internalType\":\"structNetworkConfig.FixedAddresses\",\"name\":\"_addresses\",\"type\":\"tuple\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBusContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkEnclaveRegistryContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506112b48061001c5f395ff3fe608060405234801561000f575f5ffd5b50600436106100f0575f3560e01c80638da5cb5b11610093578063da0321cd11610063578063da0321cd146101db578063f2fde38b146101f0578063f5e9f28614610203578063fbfd6d911461020b575f5ffd5b80638da5cb5b1461015c578063a1b918d61461018c578063af45463514610194578063ceb35b0f146101c8575f5ffd5b80634e027271116100ce5780634e0272711461013c578063715018a61461014457806372bad9121461014c57806385a7805a14610154575f5ffd5b80630415cc15146100f45780632fc00c761461010957806348d8723914610127575b5f5ffd5b610107610102366004610be0565b610213565b005b61011161042c565b60405161011e9190610c13565b60405180910390f35b61012f610464565b60405161011e9190610c27565b610111610492565b6101076104c1565b61012f6104d4565b61012f6104ff565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b0316610111565b61011161052a565b6101116101a2366004610cc1565b80516020818301810180516001825292820191909301209152546001600160a01b031681565b6101076101d6366004610d4f565b610559565b6101e36106ab565b60405161011e9190610ef4565b6101076101fe366004610f05565b6108cc565b610111610922565b61012f610951565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f8115801561025d5750825b90505f8267ffffffffffffffff1660011480156102795750303b155b905081158015610287575080155b156102be576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156102f257845468ff00000000000000001916680100000000000000001785555b6102fb3361097c565b61032f61032960017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958610f36565b87519055565b61036661035d60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c610f36565b60208801519055565b61039d61039460017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba610f36565b60408801519055565b6103d46103cb60017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc610f36565b60608801519055565b831561042457845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061041b90600190610f6c565b60405180910390a15b505050505050565b5f61045f61045b60017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba610f36565b5490565b905090565b61048f60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c610f36565b81565b5f61045f61045b60017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc610f36565b6104c961098d565b6104d25f610a01565b565b61048f60017fa8dc982740f2c3c626e5e571dc05dd1658ff80318c0fb06acc8b264b5ed7ebba610f36565b61048f60017f11991024cba012d71648dba5812a933585fbe131af1382f6830696bff498c6fc610f36565b5f61045f61045b60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958610f36565b61056161098d565b6001600160a01b0381166105905760405162461bcd60e51b815260040161058790610fac565b60405180910390fd5b5f6001600160a01b0316600184846040516105ac929190610fcc565b908152604051908190036020019020546001600160a01b0316146105e25760405162461bcd60e51b815260040161058790611009565b80600184846040516105f5929190610fcc565b90815260405190819003602001902080546001600160a01b039290921673ffffffffffffffffffffffffffffffffffffffff199092169190911790555f80546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5630161066a8385836110f5565b507f8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b7383838360405161069e939291906111d0565b60405180910390a1505050565b6040805160a0810182525f8082526020820181905291810182905260608082018390526080820152815490919067ffffffffffffffff8111156106f0576106f0610aed565b60405190808252806020026020018201604052801561073557816020015b60408051808201909152606081525f602082015281526020019060019003908161070e5790505b5090505f5b5f5481101561085f5760405180604001604052805f8381548110610760576107606111f1565b905f5260205f200180546107739061102d565b80601f016020809104026020016040519081016040528092919081815260200182805461079f9061102d565b80156107ea5780601f106107c1576101008083540402835291602001916107ea565b820191905f5260205f20905b8154815290600101906020018083116107cd57829003601f168201915b5050505050815260200160015f8481548110610808576108086111f1565b905f5260205f200160405161081d9190611274565b908152604051908190036020019020546001600160a01b03169052825183908390811061084c5761084c6111f1565b602090810291909101015260010161073a565b506040518060a0016040528061087361052a565b6001600160a01b03168152602001610889610922565b6001600160a01b0316815260200161089f61042c565b6001600160a01b031681526020016108b5610492565b6001600160a01b0316815260200191909152919050565b6108d461098d565b6001600160a01b038116610916575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016105879190610c13565b61091f81610a01565b50565b5f61045f61045b60017f83a6e12707c3cce2dda8a0b6be6d727d0c7e3f872360a29f026e5f6fb65eff2c610f36565b61048f60017fa508d09e1d1c531b763d64886006a2907e36a4e174a478e71c5c12a783071958610f36565b610984610a7e565b61091f81610ae5565b336109bf7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146104d257336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016105879190610c13565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166104d2576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108d4610a7e565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff82111715610b2757610b27610aed565b6040525050565b5f610b3860405190565b9050610b448282610b01565b919050565b5f6001600160a01b0382165b92915050565b610b6481610b49565b811461091f575f5ffd5b8035610b5581610b5b565b5f60808284031215610b8c57610b8c5f5ffd5b610b966080610b2e565b9050610ba28383610b6e565b8152610bb18360208401610b6e565b6020820152610bc38360408401610b6e565b6040820152610bd58360608401610b6e565b606082015292915050565b5f60808284031215610bf357610bf35f5ffd5b610bfd8383610b79565b9392505050565b610c0d81610b49565b82525050565b60208101610b558284610c04565b80610c0d565b60208101610b558284610c21565b5f67ffffffffffffffff821115610c4e57610c4e610aed565b601f19601f83011660200192915050565b82818337505f910152565b5f610c7c610c7784610c35565b610b2e565b9050828152838383011115610c9257610c925f5ffd5b610bfd836020830184610c5f565b5f82601f830112610cb257610cb25f5ffd5b610bfd83833560208501610c6a565b5f60208284031215610cd457610cd45f5ffd5b813567ffffffffffffffff811115610ced57610ced5f5ffd5b610cf984828501610ca0565b949350505050565b5f5f83601f840112610d1457610d145f5ffd5b50813567ffffffffffffffff811115610d2e57610d2e5f5ffd5b602083019150836001820283011115610d4857610d485f5ffd5b9250929050565b5f5f5f60408486031215610d6457610d645f5ffd5b833567ffffffffffffffff811115610d7d57610d7d5f5ffd5b610d8986828701610d01565b9350935050610d9b8560208601610b6e565b90509250925092565b8281835e505f910152565b5f610db8825190565b808452602084019350610dcf818560208601610da4565b601f01601f19169290920192915050565b805160408084525f9190840190610df78282610daf565b9150506020830151610e0c6020860182610c04565b509392505050565b5f610bfd8383610de0565b60200190565b5f610e2e825190565b80845260208401935083602082028501610e488560200190565b5f5b84811015610e7b5783830388528151610e638482610e14565b93505060208201602098909801979150600101610e4a565b50909695505050505050565b80515f9060a0840190610e9a8582610c04565b506020830151610ead6020860182610c04565b506040830151610ec06040860182610c04565b506060830151610ed36060860182610c04565b5060808301518482036080860152610eeb8282610e25565b95945050505050565b60208082528101610bfd8184610e87565b5f60208284031215610f1857610f185f5ffd5b610bfd8383610b6e565b634e487b7160e01b5f52601160045260245ffd5b81810381811115610b5557610b55610f22565b5f610b5582610f56565b90565b67ffffffffffffffff1690565b610c0d81610f49565b60208101610b558284610f63565b600f8152602081017f496e76616c69642061646472657373000000000000000000000000000000000081529050610e1f565b60208082528101610b5581610f7a565b610fc7828483610c5f565b500190565b610bfd818385610fbc565b601b8152602081017f41646472657373206e616d6520616c726561647920657869737473000000000081529050610e1f565b60208082528101610b5581610fd7565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061104157607f821691505b60208210810361105357611053611019565b50919050565b5f610b55610f538381565b61106d83611059565b81545f1960089490940293841b1916921b91909117905550565b5f611093818484611064565b505050565b818110156110b2576110aa5f82611087565b600101611098565b5050565b601f821115611093575f818152602090206020601f850104810160208510156110dc5750805b6110ee6020601f860104830182611098565b5050505050565b8267ffffffffffffffff81111561110e5761110e610aed565b611118825461102d565b6111238282856110b6565b505f601f821160018114611155575f831561113e5750848201355b5f19600885021c1981166002850217855550610424565b5f84815260208120601f198516915b828110156111845787850135825560209485019460019092019101611164565b50848210156111a0575f196008601f8716021c19878501351681555b5050505060020260010190555050565b8183526020830192506111c4828483610c5f565b50601f01601f19160190565b604080825281016111e28185876111b0565b9050610cf96020830184610c04565b634e487b7160e01b5f52603260045260245ffd5b5f81546112118161102d565b600182168015611228576001811461123d5761126b565b60ff198316865281151582028601935061126b565b5f858152602090205f5b8381101561126357815488820152600190910190602001611247565b505081860193505b50505092915050565b610b55818361120556fea2646970667358221220519220c36274a083998bbb863c09d6dfa5edff9a3862a894ff9da65c2bd3435664736f6c634300081c0033",
}

// NetworkConfigABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkConfigMetaData.ABI instead.
var NetworkConfigABI = NetworkConfigMetaData.ABI

// NetworkConfigBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NetworkConfigMetaData.Bin instead.
var NetworkConfigBin = NetworkConfigMetaData.Bin

// DeployNetworkConfig deploys a new Ethereum contract, binding an instance of NetworkConfig to it.
func DeployNetworkConfig(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NetworkConfig, error) {
	parsed, err := NetworkConfigMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NetworkConfigBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NetworkConfig{NetworkConfigCaller: NetworkConfigCaller{contract: contract}, NetworkConfigTransactor: NetworkConfigTransactor{contract: contract}, NetworkConfigFilterer: NetworkConfigFilterer{contract: contract}}, nil
}

// NetworkConfig is an auto generated Go binding around an Ethereum contract.
type NetworkConfig struct {
	NetworkConfigCaller     // Read-only binding to the contract
	NetworkConfigTransactor // Write-only binding to the contract
	NetworkConfigFilterer   // Log filterer for contract events
}

// NetworkConfigCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkConfigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkConfigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkConfigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkConfigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkConfigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkConfigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkConfigSession struct {
	Contract     *NetworkConfig    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NetworkConfigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkConfigCallerSession struct {
	Contract *NetworkConfigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// NetworkConfigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkConfigTransactorSession struct {
	Contract     *NetworkConfigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// NetworkConfigRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkConfigRaw struct {
	Contract *NetworkConfig // Generic contract binding to access the raw methods on
}

// NetworkConfigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkConfigCallerRaw struct {
	Contract *NetworkConfigCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkConfigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkConfigTransactorRaw struct {
	Contract *NetworkConfigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkConfig creates a new instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfig(address common.Address, backend bind.ContractBackend) (*NetworkConfig, error) {
	contract, err := bindNetworkConfig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkConfig{NetworkConfigCaller: NetworkConfigCaller{contract: contract}, NetworkConfigTransactor: NetworkConfigTransactor{contract: contract}, NetworkConfigFilterer: NetworkConfigFilterer{contract: contract}}, nil
}

// NewNetworkConfigCaller creates a new read-only instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfigCaller(address common.Address, caller bind.ContractCaller) (*NetworkConfigCaller, error) {
	contract, err := bindNetworkConfig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigCaller{contract: contract}, nil
}

// NewNetworkConfigTransactor creates a new write-only instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfigTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkConfigTransactor, error) {
	contract, err := bindNetworkConfig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigTransactor{contract: contract}, nil
}

// NewNetworkConfigFilterer creates a new log filterer instance of NetworkConfig, bound to a specific deployed contract.
func NewNetworkConfigFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkConfigFilterer, error) {
	contract, err := bindNetworkConfig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigFilterer{contract: contract}, nil
}

// bindNetworkConfig binds a generic wrapper to an already deployed contract.
func bindNetworkConfig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NetworkConfigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkConfig *NetworkConfigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkConfig.Contract.NetworkConfigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkConfig *NetworkConfigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkConfig.Contract.NetworkConfigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkConfig *NetworkConfigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkConfig.Contract.NetworkConfigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkConfig *NetworkConfigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkConfig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkConfig *NetworkConfigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkConfig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkConfig *NetworkConfigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkConfig.Contract.contract.Transact(opts, method, params...)
}

// CROSSCHAINSLOT is a free data retrieval call binding the contract method 0xfbfd6d91.
//
// Solidity: function CROSS_CHAIN_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) CROSSCHAINSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "CROSS_CHAIN_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CROSSCHAINSLOT is a free data retrieval call binding the contract method 0xfbfd6d91.
//
// Solidity: function CROSS_CHAIN_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) CROSSCHAINSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.CROSSCHAINSLOT(&_NetworkConfig.CallOpts)
}

// CROSSCHAINSLOT is a free data retrieval call binding the contract method 0xfbfd6d91.
//
// Solidity: function CROSS_CHAIN_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) CROSSCHAINSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.CROSSCHAINSLOT(&_NetworkConfig.CallOpts)
}

// MESSAGEBUSSLOT is a free data retrieval call binding the contract method 0x48d87239.
//
// Solidity: function MESSAGE_BUS_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) MESSAGEBUSSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "MESSAGE_BUS_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MESSAGEBUSSLOT is a free data retrieval call binding the contract method 0x48d87239.
//
// Solidity: function MESSAGE_BUS_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) MESSAGEBUSSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.MESSAGEBUSSLOT(&_NetworkConfig.CallOpts)
}

// MESSAGEBUSSLOT is a free data retrieval call binding the contract method 0x48d87239.
//
// Solidity: function MESSAGE_BUS_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) MESSAGEBUSSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.MESSAGEBUSSLOT(&_NetworkConfig.CallOpts)
}

// NETWORKENCLAVEREGISTRYSLOT is a free data retrieval call binding the contract method 0x72bad912.
//
// Solidity: function NETWORK_ENCLAVE_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) NETWORKENCLAVEREGISTRYSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "NETWORK_ENCLAVE_REGISTRY_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NETWORKENCLAVEREGISTRYSLOT is a free data retrieval call binding the contract method 0x72bad912.
//
// Solidity: function NETWORK_ENCLAVE_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) NETWORKENCLAVEREGISTRYSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.NETWORKENCLAVEREGISTRYSLOT(&_NetworkConfig.CallOpts)
}

// NETWORKENCLAVEREGISTRYSLOT is a free data retrieval call binding the contract method 0x72bad912.
//
// Solidity: function NETWORK_ENCLAVE_REGISTRY_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) NETWORKENCLAVEREGISTRYSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.NETWORKENCLAVEREGISTRYSLOT(&_NetworkConfig.CallOpts)
}

// ROLLUPCONTRACTSLOT is a free data retrieval call binding the contract method 0x85a7805a.
//
// Solidity: function ROLLUP_CONTRACT_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCaller) ROLLUPCONTRACTSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "ROLLUP_CONTRACT_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ROLLUPCONTRACTSLOT is a free data retrieval call binding the contract method 0x85a7805a.
//
// Solidity: function ROLLUP_CONTRACT_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigSession) ROLLUPCONTRACTSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.ROLLUPCONTRACTSLOT(&_NetworkConfig.CallOpts)
}

// ROLLUPCONTRACTSLOT is a free data retrieval call binding the contract method 0x85a7805a.
//
// Solidity: function ROLLUP_CONTRACT_SLOT() view returns(bytes32)
func (_NetworkConfig *NetworkConfigCallerSession) ROLLUPCONTRACTSLOT() ([32]byte, error) {
	return _NetworkConfig.Contract.ROLLUPCONTRACTSLOT(&_NetworkConfig.CallOpts)
}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string ) view returns(address)
func (_NetworkConfig *NetworkConfigCaller) AdditionalAddresses(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "additionalAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string ) view returns(address)
func (_NetworkConfig *NetworkConfigSession) AdditionalAddresses(arg0 string) (common.Address, error) {
	return _NetworkConfig.Contract.AdditionalAddresses(&_NetworkConfig.CallOpts, arg0)
}

// AdditionalAddresses is a free data retrieval call binding the contract method 0xaf454635.
//
// Solidity: function additionalAddresses(string ) view returns(address)
func (_NetworkConfig *NetworkConfigCallerSession) AdditionalAddresses(arg0 string) (common.Address, error) {
	return _NetworkConfig.Contract.AdditionalAddresses(&_NetworkConfig.CallOpts, arg0)
}

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address,(string,address)[]))
func (_NetworkConfig *NetworkConfigCaller) Addresses(opts *bind.CallOpts) (NetworkConfigAddresses, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "addresses")

	if err != nil {
		return *new(NetworkConfigAddresses), err
	}

	out0 := *abi.ConvertType(out[0], new(NetworkConfigAddresses)).(*NetworkConfigAddresses)

	return out0, err

}

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address,(string,address)[]))
func (_NetworkConfig *NetworkConfigSession) Addresses() (NetworkConfigAddresses, error) {
	return _NetworkConfig.Contract.Addresses(&_NetworkConfig.CallOpts)
}

// Addresses is a free data retrieval call binding the contract method 0xda0321cd.
//
// Solidity: function addresses() view returns((address,address,address,address,(string,address)[]))
func (_NetworkConfig *NetworkConfigCallerSession) Addresses() (NetworkConfigAddresses, error) {
	return _NetworkConfig.Contract.Addresses(&_NetworkConfig.CallOpts)
}

// CrossChainContractAddress is a free data retrieval call binding the contract method 0xa1b918d6.
//
// Solidity: function crossChainContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) CrossChainContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "crossChainContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrossChainContractAddress is a free data retrieval call binding the contract method 0xa1b918d6.
//
// Solidity: function crossChainContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) CrossChainContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.CrossChainContractAddress(&_NetworkConfig.CallOpts)
}

// CrossChainContractAddress is a free data retrieval call binding the contract method 0xa1b918d6.
//
// Solidity: function crossChainContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) CrossChainContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.CrossChainContractAddress(&_NetworkConfig.CallOpts)
}

// MessageBusContractAddress is a free data retrieval call binding the contract method 0xf5e9f286.
//
// Solidity: function messageBusContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) MessageBusContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "messageBusContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBusContractAddress is a free data retrieval call binding the contract method 0xf5e9f286.
//
// Solidity: function messageBusContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) MessageBusContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.MessageBusContractAddress(&_NetworkConfig.CallOpts)
}

// MessageBusContractAddress is a free data retrieval call binding the contract method 0xf5e9f286.
//
// Solidity: function messageBusContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) MessageBusContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.MessageBusContractAddress(&_NetworkConfig.CallOpts)
}

// NetworkEnclaveRegistryContractAddress is a free data retrieval call binding the contract method 0x2fc00c76.
//
// Solidity: function networkEnclaveRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) NetworkEnclaveRegistryContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "networkEnclaveRegistryContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NetworkEnclaveRegistryContractAddress is a free data retrieval call binding the contract method 0x2fc00c76.
//
// Solidity: function networkEnclaveRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) NetworkEnclaveRegistryContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.NetworkEnclaveRegistryContractAddress(&_NetworkConfig.CallOpts)
}

// NetworkEnclaveRegistryContractAddress is a free data retrieval call binding the contract method 0x2fc00c76.
//
// Solidity: function networkEnclaveRegistryContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) NetworkEnclaveRegistryContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.NetworkEnclaveRegistryContractAddress(&_NetworkConfig.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkConfig *NetworkConfigCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkConfig *NetworkConfigSession) Owner() (common.Address, error) {
	return _NetworkConfig.Contract.Owner(&_NetworkConfig.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkConfig *NetworkConfigCallerSession) Owner() (common.Address, error) {
	return _NetworkConfig.Contract.Owner(&_NetworkConfig.CallOpts)
}

// RollupContractAddress is a free data retrieval call binding the contract method 0x4e027271.
//
// Solidity: function rollupContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCaller) RollupContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkConfig.contract.Call(opts, &out, "rollupContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupContractAddress is a free data retrieval call binding the contract method 0x4e027271.
//
// Solidity: function rollupContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigSession) RollupContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.RollupContractAddress(&_NetworkConfig.CallOpts)
}

// RollupContractAddress is a free data retrieval call binding the contract method 0x4e027271.
//
// Solidity: function rollupContractAddress() view returns(address addr_)
func (_NetworkConfig *NetworkConfigCallerSession) RollupContractAddress() (common.Address, error) {
	return _NetworkConfig.Contract.RollupContractAddress(&_NetworkConfig.CallOpts)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address addr) returns()
func (_NetworkConfig *NetworkConfigTransactor) AddAddress(opts *bind.TransactOpts, name string, addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "addAddress", name, addr)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address addr) returns()
func (_NetworkConfig *NetworkConfigSession) AddAddress(name string, addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.AddAddress(&_NetworkConfig.TransactOpts, name, addr)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address addr) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) AddAddress(name string, addr common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.AddAddress(&_NetworkConfig.TransactOpts, name, addr)
}

// Initialize is a paid mutator transaction binding the contract method 0x0415cc15.
//
// Solidity: function initialize((address,address,address,address) _addresses) returns()
func (_NetworkConfig *NetworkConfigTransactor) Initialize(opts *bind.TransactOpts, _addresses NetworkConfigFixedAddresses) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "initialize", _addresses)
}

// Initialize is a paid mutator transaction binding the contract method 0x0415cc15.
//
// Solidity: function initialize((address,address,address,address) _addresses) returns()
func (_NetworkConfig *NetworkConfigSession) Initialize(_addresses NetworkConfigFixedAddresses) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses)
}

// Initialize is a paid mutator transaction binding the contract method 0x0415cc15.
//
// Solidity: function initialize((address,address,address,address) _addresses) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) Initialize(_addresses NetworkConfigFixedAddresses) (*types.Transaction, error) {
	return _NetworkConfig.Contract.Initialize(&_NetworkConfig.TransactOpts, _addresses)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkConfig *NetworkConfigTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkConfig *NetworkConfigSession) RenounceOwnership() (*types.Transaction, error) {
	return _NetworkConfig.Contract.RenounceOwnership(&_NetworkConfig.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkConfig *NetworkConfigTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NetworkConfig.Contract.RenounceOwnership(&_NetworkConfig.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkConfig *NetworkConfigTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkConfig *NetworkConfigSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.TransferOwnership(&_NetworkConfig.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkConfig *NetworkConfigTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkConfig.Contract.TransferOwnership(&_NetworkConfig.TransactOpts, newOwner)
}

// NetworkConfigInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NetworkConfig contract.
type NetworkConfigInitializedIterator struct {
	Event *NetworkConfigInitialized // Event containing the contract specifics and raw log

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
func (it *NetworkConfigInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigInitialized)
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
		it.Event = new(NetworkConfigInitialized)
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
func (it *NetworkConfigInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigInitialized represents a Initialized event raised by the NetworkConfig contract.
type NetworkConfigInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkConfig *NetworkConfigFilterer) FilterInitialized(opts *bind.FilterOpts) (*NetworkConfigInitializedIterator, error) {

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NetworkConfigInitializedIterator{contract: _NetworkConfig.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkConfig *NetworkConfigFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NetworkConfigInitialized) (event.Subscription, error) {

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigInitialized)
				if err := _NetworkConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NetworkConfig *NetworkConfigFilterer) ParseInitialized(log types.Log) (*NetworkConfigInitialized, error) {
	event := new(NetworkConfigInitialized)
	if err := _NetworkConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigNetworkContractAddressAddedIterator is returned from FilterNetworkContractAddressAdded and is used to iterate over the raw logs and unpacked data for NetworkContractAddressAdded events raised by the NetworkConfig contract.
type NetworkConfigNetworkContractAddressAddedIterator struct {
	Event *NetworkConfigNetworkContractAddressAdded // Event containing the contract specifics and raw log

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
func (it *NetworkConfigNetworkContractAddressAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigNetworkContractAddressAdded)
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
		it.Event = new(NetworkConfigNetworkContractAddressAdded)
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
func (it *NetworkConfigNetworkContractAddressAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigNetworkContractAddressAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigNetworkContractAddressAdded represents a NetworkContractAddressAdded event raised by the NetworkConfig contract.
type NetworkConfigNetworkContractAddressAdded struct {
	Name string
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNetworkContractAddressAdded is a free log retrieval operation binding the contract event 0x8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73.
//
// Solidity: event NetworkContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) FilterNetworkContractAddressAdded(opts *bind.FilterOpts) (*NetworkConfigNetworkContractAddressAddedIterator, error) {

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "NetworkContractAddressAdded")
	if err != nil {
		return nil, err
	}
	return &NetworkConfigNetworkContractAddressAddedIterator{contract: _NetworkConfig.contract, event: "NetworkContractAddressAdded", logs: logs, sub: sub}, nil
}

// WatchNetworkContractAddressAdded is a free log subscription operation binding the contract event 0x8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73.
//
// Solidity: event NetworkContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) WatchNetworkContractAddressAdded(opts *bind.WatchOpts, sink chan<- *NetworkConfigNetworkContractAddressAdded) (event.Subscription, error) {

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "NetworkContractAddressAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigNetworkContractAddressAdded)
				if err := _NetworkConfig.contract.UnpackLog(event, "NetworkContractAddressAdded", log); err != nil {
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

// ParseNetworkContractAddressAdded is a log parse operation binding the contract event 0x8fda284de6722991b87a5152e250cda7a6342080d9895b760e10ce0fa5050b73.
//
// Solidity: event NetworkContractAddressAdded(string name, address addr)
func (_NetworkConfig *NetworkConfigFilterer) ParseNetworkContractAddressAdded(log types.Log) (*NetworkConfigNetworkContractAddressAdded, error) {
	event := new(NetworkConfigNetworkContractAddressAdded)
	if err := _NetworkConfig.contract.UnpackLog(event, "NetworkContractAddressAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkConfigOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NetworkConfig contract.
type NetworkConfigOwnershipTransferredIterator struct {
	Event *NetworkConfigOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NetworkConfigOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkConfigOwnershipTransferred)
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
		it.Event = new(NetworkConfigOwnershipTransferred)
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
func (it *NetworkConfigOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkConfigOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkConfigOwnershipTransferred represents a OwnershipTransferred event raised by the NetworkConfig contract.
type NetworkConfigOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NetworkConfigOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkConfig.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NetworkConfigOwnershipTransferredIterator{contract: _NetworkConfig.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkConfig *NetworkConfigFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NetworkConfigOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkConfig.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkConfigOwnershipTransferred)
				if err := _NetworkConfig.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_NetworkConfig *NetworkConfigFilterer) ParseOwnershipTransferred(log types.Log) (*NetworkConfigOwnershipTransferred, error) {
	event := new(NetworkConfigOwnershipTransferred)
	if err := _NetworkConfig.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
