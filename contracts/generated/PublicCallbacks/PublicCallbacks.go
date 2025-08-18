// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PublicCallbacks

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

// PublicCallbacksMetaData contains all meta data concerning the PublicCallbacks contract.
var PublicCallbacksMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBefore\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAfter\",\"type\":\"uint256\"}],\"name\":\"CallbackExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"CallbackRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"callbackBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"callbacks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executeNextCallbacks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"reattemptCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"callback\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackId\",\"type\":\"uint256\"}],\"name\":\"removeCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523461002257610011610026565b6040516114de61018682396114de90f35b5f80fd5b61002e6100a9565b565b61003d9060401c60ff1690565b90565b61003d9054610030565b61003d905b6001600160401b031690565b61003d905461004a565b61003d9061004f906001600160401b031682565b9061008961003d6100a592610065565b82546001600160401b0319166001600160401b03919091161790565b9055565b5f6100b261013f565b016100bc81610040565b61012e576100c98161005b565b6001600160401b03919082908116036100e0575050565b8161010f7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29361012993610079565b604051918291826001600160401b03909116815260200190565b0390a1565b63f92ee8a960e01b5f908152600490fd5b61003d61017d565b61003d61003d61003d9290565b61003d7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610147565b61003d61015456fe60806040526004361015610011575f80fd5b5f3560e01c8062e0d3b51461007f5780638129fc1c1461007a57806382fbdc9c14610075578063929d34e9146100705780639fb56ad21461006b578063a67e1760146100665763d98c61690361008457610448565b610406565b6103ee565b6103d6565b6103af565b610338565b6102f9565b5f80fd5b905035905b565b90602082820312610084576100a391610088565b90565b6100a36100a36100a39290565b906100bd906100a6565b5f5260205260405f2090565b6100a3905b6001600160a01b031690565b6100a390546100c9565b634e487b7160e01b5f52602260045260245ffd5b9060016002830492168015610118575b602083101461011357565b6100e4565b91607f1691610108565b80545f93929161013e610134836100f8565b8085529360200190565b916001811690811561018d575060011461015757505050565b61016891929394505f5260205f2090565b915f925b8184106101795750500190565b80548484015260209093019260010161016c565b92949550505060ff1916825215156020020190565b906100a391610122565b634e487b7160e01b5f52604160045260245ffd5b90601f01601f1916810190811067ffffffffffffffff8211176101e257604052565b6101ac565b9061008d610201926101f860405190565b938480926101a2565b03836101c0565b6100a39081565b6100a39054610208565b610223905f6100b3565b61022c816100da565b91610239600183016101e7565b916102466002820161020f565b916100a360046102586003850161020f565b93016100da565b610268906100ce565b9052565b90825f9392825e0152565b6102986102a16020936102ab9361028c815190565b80835293849260200190565b9586910161026c565b601f01601f191690565b0190565b9061008d946102eb6102e36102f2936080969a99979a6102d660a08801925f89019061025f565b8682036020880152610277565b986040850152565b6060830152565b019061025f565b346100845761032a61031461030f36600461008f565b610219565b9161032195939560405190565b958695866102af565b0390f35b5f91031261008457565b346100845761034836600461032e565b61035061069d565b604051005b909182601f830112156100845781359167ffffffffffffffff831161008457602001926001830284011161008457565b9060208282031261008457813567ffffffffffffffff8111610084576103ab9201610355565b9091565b61032a6103c66103c0366004610385565b90610762565b6040519182918290815260200190565b34610084576103506103e936600461008f565b610d87565b346100845761035061040136600461008f565b610d90565b346100845761041636600461032e565b610350610e8b565b6100a3916008021c81565b906100a3915461041e565b5f6104436100a39260036100b3565b610429565b346100845761032a6103c661045e36600461008f565b610434565b6100a39060401c60ff1690565b6100a39054610463565b6100a3905b67ffffffffffffffff1690565b6100a3905461047a565b61047f6100a36100a39290565b6100ce6100a36100a3926001600160a01b031690565b6100a3906104a3565b6100a3906104b9565b9067ffffffffffffffff905b9181191691161790565b61047f6100a36100a39267ffffffffffffffff1690565b906105086100a361050f926104e1565b82546104cb565b9055565b9068ff00000000000000009060401b6104d7565b906105376100a361050f92151590565b8254610513565b61026890610496565b60208101929161008d919061053e565b61055f610e93565b8061057961057361056f83610470565b1590565b9261048c565b916105835f610496565b67ffffffffffffffff84161480610696575b6001936105b26105a486610496565b9167ffffffffffffffff1690565b14908161066e575b155b9081610665575b5061063b57806105df5f6105d686610496565b940193846104f8565b61062c575b6105ec575050565b7fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29161061b5f61062793610527565b60405191829182610547565b0390a1565b6106368383610527565b6105e4565b7ff92ee8a9000000000000000000000000000000000000000000000000000000005f908152600490fd5b1590505f6105c3565b90506105bc61067c306104c2565b3b61068d6106895f6100a6565b9190565b149190506105ba565b5080610595565b61008d610557565b156106ac57565b60405162461bcd60e51b815260206004820152600d60248201527f4e6f2076616c75652073656e74000000000000000000000000000000000000006044820152606490fd5b156106f857565b60405162461bcd60e51b8152602060048201526024808201527f47617320746f6f206c6f7720636f6d706172656420746f20636f7374206f662060448201527f63616c6c000000000000000000000000000000000000000000000000000000006064820152608490fd5b6100a3916107796107725f6100a6565b34116106a5565b61079961078534610ebe565b6107936106896152086100a6565b116106f1565b3491336110f2565b156107a857565b60405162461bcd60e51b815260206004820152602260248201527f43616c6c6261636b2063616e6e6f74206265207265617474656d70746564207960448201527f65740000000000000000000000000000000000000000000000000000000000006064820152608490fd5b61008d9061083461082d6108288360036100b3565b61020f565b43116107a1565b610cd2565b90610268906100ce565b9061008d61085060405190565b92836101c0565b6100a360a0610843565b9061008d6108cd6004610872610857565b9461088561087f826100da565b87610839565b61089b610894600183016101e7565b6020880152565b6108b16108aa6002830161020f565b6040880152565b6108c76108c06003830161020f565b6060880152565b016100da565b60808401610839565b6100a390610861565b6100a390516100ce565b6100ce6100a36100a39290565b6100a3906108e9565b1561090657565b60405162461bcd60e51b815260206004820152601760248201527f43616c6c6261636b20646f6573206e6f742065786973740000000000000000006044820152606490fd5b1561095257565b60405162461bcd60e51b815260206004820152600960248201527f4e6f74206f776e657200000000000000000000000000000000000000000000006044820152606490fd5b67ffffffffffffffff81116101e257602090601f01601f19160190565b906109c66109c183610997565b610843565b918252565b3d156109e4576109da3d6109b4565b903d5f602084013e565b606090565b156109f057565b60405162461bcd60e51b815260206004820152601960248201527f43616c6c6261636b20657865637574696f6e206661696c6564000000000000006044820152606490fd5b90610a47905f19906020036008021c90565b8154169055565b915f1960089290920291821b911b6104d7565b9190610a726100a361050f936100a6565b908354610a4e565b61008d915f91610a61565b818110610a90575050565b80610a9d5f600193610a7a565b01610a85565b905f91610aca610ab6825f5260205f2090565b9283545f19600883021c1916906002021790565b905555565b9192906020821015610b3157601f8411600114610afd5761050f9293505f19600883021c1916906002021790565b5090610b2c61008d936001610b23610b18855f5260205f2090565b92601f602091010490565b82019101610a85565b610aa3565b50610b6f8293610b466001945f5260205f2090565b610b686020601f860104820192601f861680610b77575b50601f602091010490565b0190610a85565b600202179055565b610b8390888603610a35565b5f610b5d565b9290916801000000000000000082116101e25760201115610bdf576020811015610bc25761050f915f19600883021c1916906002021790565b60019160ff1916610bd6845f5260205f2090565b55600202019055565b60019150600202019055565b908154610bf7816100f8565b90818311610c20575b818310610c0e575b50505050565b610c1793610acf565b5f808080610c08565b610c2c83838387610b89565b610c00565b634e487b7160e01b5f52601160045260245ffd5b81810292918115918404141715610c5857565b610c31565b5f61008d91610beb565b634e487b7160e01b5f525f60045260245ffd5b905f03610c8a5761008d90610c5d565b610c67565b5f80825590600490610ca48360018301610c7a565b610cb18360028301610a7a565b610cbe8360038301610a7a565b0155565b905f03610c8a5761008d90610c8f565b5f610d8261008d92610d688380610cf1610cec85836100b3565b6108d6565b6020610d4c838301610d26610d05826108df565b610d1f610d19610d14896108f6565b6100ce565b916100ce565b14156108ff565b610d47610d35608086016108df565b610d41610d19336100ce565b1461094b565b6108df565b9101519082602083519301915af1610d626109cb565b506109e9565b610d7b83610d7683826100b3565b610cc2565b60036100b3565b610a7a565b61008d90610813565b5f610d8261008d92610d68610d356080610dad610cec85886100b3565b016108df565b6001600160a01b0390811691169003906001600160a01b038211610c5857565b15610dda57565b60405162461bcd60e51b815260206004820152600860248201527f4e6f742073656c660000000000000000000000000000000000000000000000006044820152606490fd5b610e5b610e4c610d14610e37610e3c610e37306104c2565b6104b9565b610e4660016108e9565b90610db3565b610e55336100ce565b14610dd3565b61008d5b610e69600161020f565b610e796106896100a3600261020f565b1461008d57610e866111e7565b610e5f565b61008d610e1f565b6100a3611341565b634e487b7160e01b5f52601260045260245ffd5b8115610eb9570490565b610e9b565b6100a3904890610eaf565b90825f939282370152565b90929192610ee46109c182610997565b93818552818301116100845761008d916020850190610ec9565b6100a3913691610ed4565b5f198114610c585760010190565b905f19906104d7565b90610f306100a361050f926100a6565b8254610f17565b906001600160a01b03906104d7565b90610f566100a361050f926104c2565b8254610f37565b9190601f8111610f6c57505050565b610f7c61008d935f5260205f2090565b906020601f840181900483019310610f9c575b6020601f90910104610b68565b9091508190610f8f565b90610faf815190565b9067ffffffffffffffff82116101e257610fd382610fcd85546100f8565b85610f5d565b602090601f831160011461100c5761050f92915f9183611001575b50505f19600883021c1916906002021790565b015190505f80610fee565b601f1983169161101f855f5260205f2090565b925f5b81811061105b57509160029391856001969410611043575b50505002019055565b01515f196008601f8516021c191690555f808061103a565b91936020600181928787015181550195019201611022565b9061008d91610fa6565b60046110e1608061008d9461109c6110965f83016108df565b86610f46565b6110b36110aa602083015190565b60018701611073565b6110ca6110c1604083015190565b60028701610f20565b610dad6110d8606083015190565b60038701610f20565b9101610f46565b9061008d9161107d565b916111a47f3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c4936111616106279497966111285f90565b5061115a611153611139600161020f565b96879b61114e611147610857565b9788610839565b610efe565b6020850152565b6040830152565b61116c486060830152565b6111793360808301610839565b61119f611186600161020f565b61119961119282610f09565b6001610f20565b5f6100b3565b6110e8565b6103c66111b28260036100b3565b4390610f20565b91908203918211610c5857565b90815260608101939261008d9290916040916111e3906020830152565b0152565b6111f1600161020f565b6112016106896100a3600261020f565b1461008d577f79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b261122f611378565b61062761123d606084015190565b6112d760408501946112d28561125a856112558a5190565b610eaf565b975a985f80611268876108df565b602088015190826020835193019186f1916112816109cb565b505a9761128e898d6111b9565b6112975f6100a6565b938181116112f0575b5050505f6112b8836112b36112bf945190565b6111b9565b96016108df565b916112e3575b6112cd6113ca565b6113e6565b611470565b604051938493846111c6565b6112eb6113a2565b6112c5565b6112b39450926113105f9361130b6112b8946112bf976111b9565b610c45565b9450926112a0565b6100a37ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a006100a6565b6100a3611318565b611351610857565b905f8252602080808080860160608152015f8152015f8152015f905250565b6100a3611349565b611380611370565b5061138e611199600261020f565b9061068961139c600261020f565b926108d6565b6113b95f610d766113b3600261020f565b826100b3565b61008d5f610d82610d7b600261020f565b61008d6113df6113da600261020f565b610f09565b6002610f20565b915f6114336114649361144283946113fd60405190565b9384916004602084017f5ea395580000000000000000000000000000000000000000000000000000000081520190815260200190565b602082018103825203836101c0565b8561144e61afc86100a6565b9160208451940192f161145f6109cb565b501590565b61146b5750565b61008d905b6114795f6100a6565b81146114a5575f809161148b416104c2565b9061149560405190565b90818003925af1506114a56109cb565b5056fea2646970667358221220487a69253bfc29679edada99cb14c526b02b6d21c502edad1d170a376a1dd86b64736f6c634300081c0033",
}

// PublicCallbacksABI is the input ABI used to generate the binding from.
// Deprecated: Use PublicCallbacksMetaData.ABI instead.
var PublicCallbacksABI = PublicCallbacksMetaData.ABI

// PublicCallbacksBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PublicCallbacksMetaData.Bin instead.
var PublicCallbacksBin = PublicCallbacksMetaData.Bin

// DeployPublicCallbacks deploys a new Ethereum contract, binding an instance of PublicCallbacks to it.
func DeployPublicCallbacks(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PublicCallbacks, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PublicCallbacksBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// PublicCallbacks is an auto generated Go binding around an Ethereum contract.
type PublicCallbacks struct {
	PublicCallbacksCaller     // Read-only binding to the contract
	PublicCallbacksTransactor // Write-only binding to the contract
	PublicCallbacksFilterer   // Log filterer for contract events
}

// PublicCallbacksCaller is an auto generated read-only Go binding around an Ethereum contract.
type PublicCallbacksCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PublicCallbacksFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicCallbacksSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PublicCallbacksSession struct {
	Contract     *PublicCallbacks  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PublicCallbacksCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PublicCallbacksCallerSession struct {
	Contract *PublicCallbacksCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// PublicCallbacksTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PublicCallbacksTransactorSession struct {
	Contract     *PublicCallbacksTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// PublicCallbacksRaw is an auto generated low-level Go binding around an Ethereum contract.
type PublicCallbacksRaw struct {
	Contract *PublicCallbacks // Generic contract binding to access the raw methods on
}

// PublicCallbacksCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PublicCallbacksCallerRaw struct {
	Contract *PublicCallbacksCaller // Generic read-only contract binding to access the raw methods on
}

// PublicCallbacksTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PublicCallbacksTransactorRaw struct {
	Contract *PublicCallbacksTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicCallbacks creates a new instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacks(address common.Address, backend bind.ContractBackend) (*PublicCallbacks, error) {
	contract, err := bindPublicCallbacks(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacks{PublicCallbacksCaller: PublicCallbacksCaller{contract: contract}, PublicCallbacksTransactor: PublicCallbacksTransactor{contract: contract}, PublicCallbacksFilterer: PublicCallbacksFilterer{contract: contract}}, nil
}

// NewPublicCallbacksCaller creates a new read-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksCaller(address common.Address, caller bind.ContractCaller) (*PublicCallbacksCaller, error) {
	contract, err := bindPublicCallbacks(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCaller{contract: contract}, nil
}

// NewPublicCallbacksTransactor creates a new write-only instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicCallbacksTransactor, error) {
	contract, err := bindPublicCallbacks(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksTransactor{contract: contract}, nil
}

// NewPublicCallbacksFilterer creates a new log filterer instance of PublicCallbacks, bound to a specific deployed contract.
func NewPublicCallbacksFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicCallbacksFilterer, error) {
	contract, err := bindPublicCallbacks(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksFilterer{contract: contract}, nil
}

// bindPublicCallbacks binds a generic wrapper to an already deployed contract.
func bindPublicCallbacks(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PublicCallbacksMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.PublicCallbacksCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.PublicCallbacksTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicCallbacks *PublicCallbacksCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicCallbacks.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicCallbacks *PublicCallbacksTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.contract.Transact(opts, method, params...)
}

// CallbackBlockNumber is a free data retrieval call binding the contract method 0xd98c6169.
//
// Solidity: function callbackBlockNumber(uint256 callbackId) view returns(uint256 blockNumber)
func (_PublicCallbacks *PublicCallbacksCaller) CallbackBlockNumber(opts *bind.CallOpts, callbackId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicCallbacks.contract.Call(opts, &out, "callbackBlockNumber", callbackId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CallbackBlockNumber is a free data retrieval call binding the contract method 0xd98c6169.
//
// Solidity: function callbackBlockNumber(uint256 callbackId) view returns(uint256 blockNumber)
func (_PublicCallbacks *PublicCallbacksSession) CallbackBlockNumber(callbackId *big.Int) (*big.Int, error) {
	return _PublicCallbacks.Contract.CallbackBlockNumber(&_PublicCallbacks.CallOpts, callbackId)
}

// CallbackBlockNumber is a free data retrieval call binding the contract method 0xd98c6169.
//
// Solidity: function callbackBlockNumber(uint256 callbackId) view returns(uint256 blockNumber)
func (_PublicCallbacks *PublicCallbacksCallerSession) CallbackBlockNumber(callbackId *big.Int) (*big.Int, error) {
	return _PublicCallbacks.Contract.CallbackBlockNumber(&_PublicCallbacks.CallOpts, callbackId)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 callbackId) view returns(address target, bytes data, uint256 value, uint256 baseFee, address owner)
func (_PublicCallbacks *PublicCallbacksCaller) Callbacks(opts *bind.CallOpts, callbackId *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
	Owner   common.Address
}, error) {
	var out []interface{}
	err := _PublicCallbacks.contract.Call(opts, &out, "callbacks", callbackId)

	outstruct := new(struct {
		Target  common.Address
		Data    []byte
		Value   *big.Int
		BaseFee *big.Int
		Owner   common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Target = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Data = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Value = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.BaseFee = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 callbackId) view returns(address target, bytes data, uint256 value, uint256 baseFee, address owner)
func (_PublicCallbacks *PublicCallbacksSession) Callbacks(callbackId *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
	Owner   common.Address
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, callbackId)
}

// Callbacks is a free data retrieval call binding the contract method 0x00e0d3b5.
//
// Solidity: function callbacks(uint256 callbackId) view returns(address target, bytes data, uint256 value, uint256 baseFee, address owner)
func (_PublicCallbacks *PublicCallbacksCallerSession) Callbacks(callbackId *big.Int) (struct {
	Target  common.Address
	Data    []byte
	Value   *big.Int
	BaseFee *big.Int
	Owner   common.Address
}, error) {
	return _PublicCallbacks.Contract.Callbacks(&_PublicCallbacks.CallOpts, callbackId)
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ExecuteNextCallbacks(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "executeNextCallbacks")
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksSession) ExecuteNextCallbacks() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallbacks(&_PublicCallbacks.TransactOpts)
}

// ExecuteNextCallbacks is a paid mutator transaction binding the contract method 0xa67e1760.
//
// Solidity: function executeNextCallbacks() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ExecuteNextCallbacks() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ExecuteNextCallbacks(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) Initialize() (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Initialize(&_PublicCallbacks.TransactOpts)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactor) ReattemptCallback(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "reattemptCallback", callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// ReattemptCallback is a paid mutator transaction binding the contract method 0x929d34e9.
//
// Solidity: function reattemptCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) ReattemptCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.ReattemptCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksTransactor) Register(opts *bind.TransactOpts, callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "register", callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// Register is a paid mutator transaction binding the contract method 0x82fbdc9c.
//
// Solidity: function register(bytes callback) payable returns(uint256)
func (_PublicCallbacks *PublicCallbacksTransactorSession) Register(callback []byte) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.Register(&_PublicCallbacks.TransactOpts, callback)
}

// RemoveCallback is a paid mutator transaction binding the contract method 0x9fb56ad2.
//
// Solidity: function removeCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactor) RemoveCallback(opts *bind.TransactOpts, callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.contract.Transact(opts, "removeCallback", callbackId)
}

// RemoveCallback is a paid mutator transaction binding the contract method 0x9fb56ad2.
//
// Solidity: function removeCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksSession) RemoveCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.RemoveCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// RemoveCallback is a paid mutator transaction binding the contract method 0x9fb56ad2.
//
// Solidity: function removeCallback(uint256 callbackId) returns()
func (_PublicCallbacks *PublicCallbacksTransactorSession) RemoveCallback(callbackId *big.Int) (*types.Transaction, error) {
	return _PublicCallbacks.Contract.RemoveCallback(&_PublicCallbacks.TransactOpts, callbackId)
}

// PublicCallbacksCallbackExecutedIterator is returned from FilterCallbackExecuted and is used to iterate over the raw logs and unpacked data for CallbackExecuted events raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecutedIterator struct {
	Event *PublicCallbacksCallbackExecuted // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksCallbackExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksCallbackExecuted)
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
		it.Event = new(PublicCallbacksCallbackExecuted)
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
func (it *PublicCallbacksCallbackExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksCallbackExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksCallbackExecuted represents a CallbackExecuted event raised by the PublicCallbacks contract.
type PublicCallbacksCallbackExecuted struct {
	CallbackId *big.Int
	GasBefore  *big.Int
	GasAfter   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCallbackExecuted is a free log retrieval operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterCallbackExecuted(opts *bind.FilterOpts) (*PublicCallbacksCallbackExecutedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCallbackExecutedIterator{contract: _PublicCallbacks.contract, event: "CallbackExecuted", logs: logs, sub: sub}, nil
}

// WatchCallbackExecuted is a free log subscription operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchCallbackExecuted(opts *bind.WatchOpts, sink chan<- *PublicCallbacksCallbackExecuted) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "CallbackExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksCallbackExecuted)
				if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
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

// ParseCallbackExecuted is a log parse operation binding the contract event 0x79867de645e468e8c09d74e8be7ed5d3ffcb800407d63d145988787eb329c9b2.
//
// Solidity: event CallbackExecuted(uint256 callbackId, uint256 gasBefore, uint256 gasAfter)
func (_PublicCallbacks *PublicCallbacksFilterer) ParseCallbackExecuted(log types.Log) (*PublicCallbacksCallbackExecuted, error) {
	event := new(PublicCallbacksCallbackExecuted)
	if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicCallbacksCallbackRegisteredIterator is returned from FilterCallbackRegistered and is used to iterate over the raw logs and unpacked data for CallbackRegistered events raised by the PublicCallbacks contract.
type PublicCallbacksCallbackRegisteredIterator struct {
	Event *PublicCallbacksCallbackRegistered // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksCallbackRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksCallbackRegistered)
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
		it.Event = new(PublicCallbacksCallbackRegistered)
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
func (it *PublicCallbacksCallbackRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksCallbackRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksCallbackRegistered represents a CallbackRegistered event raised by the PublicCallbacks contract.
type PublicCallbacksCallbackRegistered struct {
	CallbackId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCallbackRegistered is a free log retrieval operation binding the contract event 0x3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c4.
//
// Solidity: event CallbackRegistered(uint256 callbackId)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterCallbackRegistered(opts *bind.FilterOpts) (*PublicCallbacksCallbackRegisteredIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "CallbackRegistered")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksCallbackRegisteredIterator{contract: _PublicCallbacks.contract, event: "CallbackRegistered", logs: logs, sub: sub}, nil
}

// WatchCallbackRegistered is a free log subscription operation binding the contract event 0x3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c4.
//
// Solidity: event CallbackRegistered(uint256 callbackId)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchCallbackRegistered(opts *bind.WatchOpts, sink chan<- *PublicCallbacksCallbackRegistered) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "CallbackRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksCallbackRegistered)
				if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackRegistered", log); err != nil {
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

// ParseCallbackRegistered is a log parse operation binding the contract event 0x3dcab49ed22ae75e876805945c8ebd2a1b325e5e8c4a1f0bc0959e62e18cc9c4.
//
// Solidity: event CallbackRegistered(uint256 callbackId)
func (_PublicCallbacks *PublicCallbacksFilterer) ParseCallbackRegistered(log types.Log) (*PublicCallbacksCallbackRegistered, error) {
	event := new(PublicCallbacksCallbackRegistered)
	if err := _PublicCallbacks.contract.UnpackLog(event, "CallbackRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicCallbacksInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PublicCallbacks contract.
type PublicCallbacksInitializedIterator struct {
	Event *PublicCallbacksInitialized // Event containing the contract specifics and raw log

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
func (it *PublicCallbacksInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicCallbacksInitialized)
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
		it.Event = new(PublicCallbacksInitialized)
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
func (it *PublicCallbacksInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicCallbacksInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicCallbacksInitialized represents a Initialized event raised by the PublicCallbacks contract.
type PublicCallbacksInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) FilterInitialized(opts *bind.FilterOpts) (*PublicCallbacksInitializedIterator, error) {

	logs, sub, err := _PublicCallbacks.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PublicCallbacksInitializedIterator{contract: _PublicCallbacks.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicCallbacks *PublicCallbacksFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PublicCallbacksInitialized) (event.Subscription, error) {

	logs, sub, err := _PublicCallbacks.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicCallbacksInitialized)
				if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_PublicCallbacks *PublicCallbacksFilterer) ParseInitialized(log types.Log) (*PublicCallbacksInitialized, error) {
	event := new(PublicCallbacksInitialized)
	if err := _PublicCallbacks.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
