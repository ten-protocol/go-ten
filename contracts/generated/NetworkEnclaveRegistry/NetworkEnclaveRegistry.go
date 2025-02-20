// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package NetworkEnclaveRegistry

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

// NetworkEnclaveRegistryMetaData contains all meta data concerning the NetworkEnclaveRegistry contract.
var NetworkEnclaveRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"initializer\",\"type\":\"address\"}],\"name\":\"NetworkSecretInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"NetworkSecretRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"attester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"}],\"name\":\"NetworkSecretResponded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"SequencerEnclaveGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"SequencerEnclaveRevoked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"grantSequencerEnclave\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_genesisAttestation\",\"type\":\"string\"}],\"name\":\"initializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"isAttested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"enclaveID\",\"type\":\"address\"}],\"name\":\"isSequencer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"requestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"verifyAttester\",\"type\":\"bool\"}],\"name\":\"respondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"revokeSequencerEnclave\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50601633601a565b608a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b61128a806100975f395ff3fe608060405234801561000f575f5ffd5b50600436106100c4575f3560e01c8063715018a61161007d578063a341115511610058578063a3411155146101b7578063f2fde38b146101ca578063f3cbc5f8146101dd575f5ffd5b8063715018a61461016f5780638129fc1c146101775780638da5cb5b1461017f575f5ffd5b80635ad124ef116100ad5780635ad124ef1461011e5780635b719ceb146101315780636d46e98714610144575f5ffd5b80633c23afba146100c8578063534ddc7a14610109575b5f5ffd5b6100f36100d6366004610b6f565b6001600160a01b03165f9081526001602052604090205460ff1690565b6040516101009190610b9d565b60405180910390f35b61011c610117366004610b6f565b6101f0565b005b61011c61012c366004610bf9565b610292565b61011c61013f366004610d39565b610308565b6100f3610152366004610b6f565b6001600160a01b03165f9081526002602052604090205460ff1690565b61011c61049d565b61011c6104b0565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516101009190610de3565b61011c6101c5366004610b6f565b6105f8565b61011c6101d8366004610b6f565b610689565b61011c6101eb366004610df1565b6106df565b6101f8610793565b6001600160a01b0381165f9081526002602052604090205460ff166102385760405162461bcd60e51b815260040161022f90610eac565b60405180910390fd5b6001600160a01b0381165f9081526002602052604090819020805460ff19169055517f0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b4790610287908390610de3565b60405180910390a150565b335f9081526001602052604090205460ff16156102c15760405162461bcd60e51b815260040161022f90610eee565b336001600160a01b03167f0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d430183836040516102fc929190610f1e565b60405180910390a25050565b6001600160a01b0385165f9081526001602052604090205460ff1661033f5760405162461bcd60e51b815260040161022f90610f92565b6001600160a01b0384165f9081526001602052604090205460ff16156103775760405162461bcd60e51b815260040161022f90610fd4565b6001600160a01b03841661039d5760405162461bcd60e51b815260040161022f90611016565b8015610445575f61040285846040516020016103ba929190611078565b604051602081830303815290604052805190602001207f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f908152601c91909152603c902090565b90505f61040f8286610807565b9050866001600160a01b0316816001600160a01b0316146104425760405162461bcd60e51b815260040161022f906110c1565b50505b6001600160a01b038085165f818152600160208190526040808320805460ff19169092179091555191928816917fb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be19190a35050505050565b6104a5610793565b6104ae5f610831565b565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156104fa5750825b90505f8267ffffffffffffffff1660011480156105165750303b155b905081158015610524575080155b1561055b576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561058f57845468ff00000000000000001916680100000000000000001785555b610598336108b9565b5f805460ff1916905583156105f157845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2906105e8906001906110eb565b60405180910390a15b5050505050565b610600610793565b6001600160a01b0381165f9081526001602052604090205460ff166106375760405162461bcd60e51b815260040161022f9061112b565b6001600160a01b0381165f9081526002602052604090819020805460ff19166001179055517ffe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e76093690610287908390610de3565b610691610793565b6001600160a01b0381166106d3575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161022f9190610de3565b6106dc81610831565b50565b5f5460ff16156107015760405162461bcd60e51b815260040161022f90611193565b6001600160a01b0385166107275760405162461bcd60e51b815260040161022f906111d5565b5f8054600160ff19918216811783556001600160a01b03881683526020818152604080852080548516841790556002909152928390208054909216179055517fd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23906105e8908790610de3565b336107c57f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146104ae57336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161022f9190610de3565b5f5f5f5f61081586866108ca565b9250925092506108258282610913565b50909150505b92915050565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080547fffffffffffffffffffffffff000000000000000000000000000000000000000081166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b6108c1610a18565b6106dc81610a7f565b5f5f5f8351604103610901576020840151604085015160608601515f1a6108f388828585610a87565b95509550955050505061090c565b505081515f91506002905b9250925092565b5f826003811115610926576109266111e5565b0361092f575050565b6001826003811115610943576109436111e5565b0361097a576040517ff645eedf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600282600381111561098e5761098e6111e5565b036109c7576040517ffce698f700000000000000000000000000000000000000000000000000000000815261022f9082906004016111ff565b60038260038111156109db576109db6111e5565b03610a1457806040517fd78bce0c00000000000000000000000000000000000000000000000000000000815260040161022f91906111ff565b5050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166104ae576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610691610a18565b5f80807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841115610ac057505f91506003905082610b37565b5f6001888888886040515f8152602001604052604051610ae39493929190611216565b6020604051602081039080840390855afa158015610b03573d5f5f3e3d5ffd5b5050604051601f1901519150506001600160a01b038116610b2e57505f925060019150829050610b37565b92505f91508190505b9450945094915050565b5f6001600160a01b03821661082b565b610b5a81610b41565b81146106dc575f5ffd5b803561082b81610b51565b5f60208284031215610b8257610b825f5ffd5b610b8c8383610b64565b9392505050565b8015155b82525050565b6020810161082b8284610b93565b5f5f83601f840112610bbe57610bbe5f5ffd5b50813567ffffffffffffffff811115610bd857610bd85f5ffd5b602083019150836001820283011115610bf257610bf25f5ffd5b9250929050565b5f5f60208385031215610c0d57610c0d5f5ffd5b823567ffffffffffffffff811115610c2657610c265f5ffd5b610c3285828601610bab565b92509250509250929050565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff82111715610c7857610c78610c3e565b6040525050565b5f610c8960405190565b9050610c958282610c52565b919050565b5f67ffffffffffffffff821115610cb357610cb3610c3e565b601f19601f83011660200192915050565b82818337505f910152565b5f610ce1610cdc84610c9a565b610c7f565b9050828152838383011115610cf757610cf75f5ffd5b610b8c836020830184610cc4565b5f82601f830112610d1757610d175f5ffd5b610b8c83833560208501610ccf565b801515610b5a565b803561082b81610d26565b5f5f5f5f5f60a08688031215610d5057610d505f5ffd5b610d5a8787610b64565b9450610d698760208801610b64565b9350604086013567ffffffffffffffff811115610d8757610d875f5ffd5b610d9388828901610d05565b935050606086013567ffffffffffffffff811115610db257610db25f5ffd5b610dbe88828901610d05565b925050610dce8760808801610d2e565b90509295509295909350565b610b9781610b41565b6020810161082b8284610dda565b5f5f5f5f5f60608688031215610e0857610e085f5ffd5b610e128787610b64565b9450602086013567ffffffffffffffff811115610e3057610e305f5ffd5b610e3c88828901610bab565b9450945050604086013567ffffffffffffffff811115610e5d57610e5d5f5ffd5b610e6988828901610bab565b92509250509295509295909350565b60198152602081017f656e636c6176654944206e6f7420612073657175656e63657200000000000000815290505b60200190565b6020808252810161082b81610e78565b60108152602081017f616c72656164792061747465737465640000000000000000000000000000000081529050610ea6565b6020808252810161082b81610ebc565b818352602083019250610f12828483610cc4565b50601f01601f19160190565b60208082528101610f30818486610efe565b949350505050565b60238152602081017f726573706f6e64696e67206174746573746572206973206e6f7420617474657381527f7465640000000000000000000000000000000000000000000000000000000000602082015290505b60400190565b6020808252810161082b81610f38565b601a8152602081017f72657175657374657220616c726561647920617474657374656400000000000081529050610ea6565b6020808252810161082b81610fa2565b60198152602081017f696e76616c69642072657175657374657220616464726573730000000000000081529050610ea6565b6020808252810161082b81610fe4565b5f61082b8260601b90565b5f61082b82611026565b610b9761104782610b41565b611031565b8281835e505f910152565b5f611060825190565b61106e81856020860161104c565b9290920192915050565b611082818461103b565b601401610b8c8183611057565b60118152602081017f696e76616c6964207369676e617475726500000000000000000000000000000081529050610ea6565b6020808252810161082b8161108f565b5f67ffffffffffffffff821661082b565b610b97816110d1565b6020810161082b82846110e2565b60168152602081017f656e636c6176654944206e6f742061747465737465640000000000000000000081529050610ea6565b6020808252810161082b816110f9565b60228152602081017f6e6574776f726b2073656372657420616c726561647920696e697469616c697a81527f656400000000000000000000000000000000000000000000000000000000000060208201529050610f8c565b6020808252810161082b8161113b565b60178152602081017f696e76616c696420656e636c617665206164647265737300000000000000000081529050610ea6565b6020808252810161082b816111a3565b634e487b7160e01b5f52602160045260245ffd5b80610b97565b6020810161082b82846111f9565b60ff8116610b97565b6080810161122482876111f9565b611231602083018661120d565b61123e60408301856111f9565b61124b60608301846111f9565b9594505050505056fea2646970667358221220f0b0da3c82c96483cc6937775b592f9f0028a69f0decc50d1737fb64930cc95864736f6c634300081c0033",
}

// NetworkEnclaveRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use NetworkEnclaveRegistryMetaData.ABI instead.
var NetworkEnclaveRegistryABI = NetworkEnclaveRegistryMetaData.ABI

// NetworkEnclaveRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NetworkEnclaveRegistryMetaData.Bin instead.
var NetworkEnclaveRegistryBin = NetworkEnclaveRegistryMetaData.Bin

// DeployNetworkEnclaveRegistry deploys a new Ethereum contract, binding an instance of NetworkEnclaveRegistry to it.
func DeployNetworkEnclaveRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NetworkEnclaveRegistry, error) {
	parsed, err := NetworkEnclaveRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NetworkEnclaveRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NetworkEnclaveRegistry{NetworkEnclaveRegistryCaller: NetworkEnclaveRegistryCaller{contract: contract}, NetworkEnclaveRegistryTransactor: NetworkEnclaveRegistryTransactor{contract: contract}, NetworkEnclaveRegistryFilterer: NetworkEnclaveRegistryFilterer{contract: contract}}, nil
}

// NetworkEnclaveRegistry is an auto generated Go binding around an Ethereum contract.
type NetworkEnclaveRegistry struct {
	NetworkEnclaveRegistryCaller     // Read-only binding to the contract
	NetworkEnclaveRegistryTransactor // Write-only binding to the contract
	NetworkEnclaveRegistryFilterer   // Log filterer for contract events
}

// NetworkEnclaveRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkEnclaveRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkEnclaveRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkEnclaveRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkEnclaveRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkEnclaveRegistrySession struct {
	Contract     *NetworkEnclaveRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NetworkEnclaveRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkEnclaveRegistryCallerSession struct {
	Contract *NetworkEnclaveRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// NetworkEnclaveRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkEnclaveRegistryTransactorSession struct {
	Contract     *NetworkEnclaveRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// NetworkEnclaveRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkEnclaveRegistryRaw struct {
	Contract *NetworkEnclaveRegistry // Generic contract binding to access the raw methods on
}

// NetworkEnclaveRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryCallerRaw struct {
	Contract *NetworkEnclaveRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkEnclaveRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkEnclaveRegistryTransactorRaw struct {
	Contract *NetworkEnclaveRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkEnclaveRegistry creates a new instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistry(address common.Address, backend bind.ContractBackend) (*NetworkEnclaveRegistry, error) {
	contract, err := bindNetworkEnclaveRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistry{NetworkEnclaveRegistryCaller: NetworkEnclaveRegistryCaller{contract: contract}, NetworkEnclaveRegistryTransactor: NetworkEnclaveRegistryTransactor{contract: contract}, NetworkEnclaveRegistryFilterer: NetworkEnclaveRegistryFilterer{contract: contract}}, nil
}

// NewNetworkEnclaveRegistryCaller creates a new read-only instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistryCaller(address common.Address, caller bind.ContractCaller) (*NetworkEnclaveRegistryCaller, error) {
	contract, err := bindNetworkEnclaveRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryCaller{contract: contract}, nil
}

// NewNetworkEnclaveRegistryTransactor creates a new write-only instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkEnclaveRegistryTransactor, error) {
	contract, err := bindNetworkEnclaveRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryTransactor{contract: contract}, nil
}

// NewNetworkEnclaveRegistryFilterer creates a new log filterer instance of NetworkEnclaveRegistry, bound to a specific deployed contract.
func NewNetworkEnclaveRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkEnclaveRegistryFilterer, error) {
	contract, err := bindNetworkEnclaveRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryFilterer{contract: contract}, nil
}

// bindNetworkEnclaveRegistry binds a generic wrapper to an already deployed contract.
func bindNetworkEnclaveRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NetworkEnclaveRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkEnclaveRegistry.Contract.NetworkEnclaveRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.NetworkEnclaveRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.NetworkEnclaveRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkEnclaveRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.contract.Transact(opts, method, params...)
}

// IsAttested is a free data retrieval call binding the contract method 0x3c23afba.
//
// Solidity: function isAttested(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) IsAttested(opts *bind.CallOpts, enclaveID common.Address) (bool, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "isAttested", enclaveID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAttested is a free data retrieval call binding the contract method 0x3c23afba.
//
// Solidity: function isAttested(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) IsAttested(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsAttested(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// IsAttested is a free data retrieval call binding the contract method 0x3c23afba.
//
// Solidity: function isAttested(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) IsAttested(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsAttested(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// IsSequencer is a free data retrieval call binding the contract method 0x6d46e987.
//
// Solidity: function isSequencer(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) IsSequencer(opts *bind.CallOpts, enclaveID common.Address) (bool, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "isSequencer", enclaveID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSequencer is a free data retrieval call binding the contract method 0x6d46e987.
//
// Solidity: function isSequencer(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) IsSequencer(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsSequencer(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// IsSequencer is a free data retrieval call binding the contract method 0x6d46e987.
//
// Solidity: function isSequencer(address enclaveID) view returns(bool)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) IsSequencer(enclaveID common.Address) (bool, error) {
	return _NetworkEnclaveRegistry.Contract.IsSequencer(&_NetworkEnclaveRegistry.CallOpts, enclaveID)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NetworkEnclaveRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) Owner() (common.Address, error) {
	return _NetworkEnclaveRegistry.Contract.Owner(&_NetworkEnclaveRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryCallerSession) Owner() (common.Address, error) {
	return _NetworkEnclaveRegistry.Contract.Owner(&_NetworkEnclaveRegistry.CallOpts)
}

// GrantSequencerEnclave is a paid mutator transaction binding the contract method 0xa3411155.
//
// Solidity: function grantSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) GrantSequencerEnclave(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "grantSequencerEnclave", _addr)
}

// GrantSequencerEnclave is a paid mutator transaction binding the contract method 0xa3411155.
//
// Solidity: function grantSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) GrantSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.GrantSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// GrantSequencerEnclave is a paid mutator transaction binding the contract method 0xa3411155.
//
// Solidity: function grantSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) GrantSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.GrantSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) Initialize() (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.Initialize(&_NetworkEnclaveRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.Initialize(&_NetworkEnclaveRegistry.TransactOpts)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xf3cbc5f8.
//
// Solidity: function initializeNetworkSecret(address enclaveID, bytes _initSecret, string _genesisAttestation) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, enclaveID common.Address, _initSecret []byte, _genesisAttestation string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "initializeNetworkSecret", enclaveID, _initSecret, _genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xf3cbc5f8.
//
// Solidity: function initializeNetworkSecret(address enclaveID, bytes _initSecret, string _genesisAttestation) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) InitializeNetworkSecret(enclaveID common.Address, _initSecret []byte, _genesisAttestation string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.InitializeNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, enclaveID, _initSecret, _genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xf3cbc5f8.
//
// Solidity: function initializeNetworkSecret(address enclaveID, bytes _initSecret, string _genesisAttestation) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) InitializeNetworkSecret(enclaveID common.Address, _initSecret []byte, _genesisAttestation string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.InitializeNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, enclaveID, _initSecret, _genesisAttestation)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RenounceOwnership(&_NetworkEnclaveRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RenounceOwnership(&_NetworkEnclaveRegistry.TransactOpts)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0x5ad124ef.
//
// Solidity: function requestNetworkSecret(string requestReport) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "requestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0x5ad124ef.
//
// Solidity: function requestNetworkSecret(string requestReport) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RequestNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0x5ad124ef.
//
// Solidity: function requestNetworkSecret(string requestReport) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RequestNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x5b719ceb.
//
// Solidity: function respondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, bool verifyAttester) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, verifyAttester bool) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "respondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x5b719ceb.
//
// Solidity: function respondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, bool verifyAttester) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, verifyAttester bool) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RespondNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x5b719ceb.
//
// Solidity: function respondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, bool verifyAttester) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, verifyAttester bool) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RespondNetworkSecret(&_NetworkEnclaveRegistry.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, verifyAttester)
}

// RevokeSequencerEnclave is a paid mutator transaction binding the contract method 0x534ddc7a.
//
// Solidity: function revokeSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) RevokeSequencerEnclave(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "revokeSequencerEnclave", _addr)
}

// RevokeSequencerEnclave is a paid mutator transaction binding the contract method 0x534ddc7a.
//
// Solidity: function revokeSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) RevokeSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RevokeSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// RevokeSequencerEnclave is a paid mutator transaction binding the contract method 0x534ddc7a.
//
// Solidity: function revokeSequencerEnclave(address _addr) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) RevokeSequencerEnclave(_addr common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.RevokeSequencerEnclave(&_NetworkEnclaveRegistry.TransactOpts, _addr)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.TransferOwnership(&_NetworkEnclaveRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NetworkEnclaveRegistry.Contract.TransferOwnership(&_NetworkEnclaveRegistry.TransactOpts, newOwner)
}

// NetworkEnclaveRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryInitializedIterator struct {
	Event *NetworkEnclaveRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryInitialized)
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
		it.Event = new(NetworkEnclaveRegistryInitialized)
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
func (it *NetworkEnclaveRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryInitialized represents a Initialized event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*NetworkEnclaveRegistryInitializedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryInitializedIterator{contract: _NetworkEnclaveRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryInitialized)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseInitialized(log types.Log) (*NetworkEnclaveRegistryInitialized, error) {
	event := new(NetworkEnclaveRegistryInitialized)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryNetworkSecretInitializedIterator is returned from FilterNetworkSecretInitialized and is used to iterate over the raw logs and unpacked data for NetworkSecretInitialized events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretInitializedIterator struct {
	Event *NetworkEnclaveRegistryNetworkSecretInitialized // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryNetworkSecretInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryNetworkSecretInitialized)
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
		it.Event = new(NetworkEnclaveRegistryNetworkSecretInitialized)
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
func (it *NetworkEnclaveRegistryNetworkSecretInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryNetworkSecretInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryNetworkSecretInitialized represents a NetworkSecretInitialized event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretInitialized struct {
	Initializer common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNetworkSecretInitialized is a free log retrieval operation binding the contract event 0xd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23.
//
// Solidity: event NetworkSecretInitialized(address initializer)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterNetworkSecretInitialized(opts *bind.FilterOpts) (*NetworkEnclaveRegistryNetworkSecretInitializedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "NetworkSecretInitialized")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryNetworkSecretInitializedIterator{contract: _NetworkEnclaveRegistry.contract, event: "NetworkSecretInitialized", logs: logs, sub: sub}, nil
}

// WatchNetworkSecretInitialized is a free log subscription operation binding the contract event 0xd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23.
//
// Solidity: event NetworkSecretInitialized(address initializer)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchNetworkSecretInitialized(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryNetworkSecretInitialized) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "NetworkSecretInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryNetworkSecretInitialized)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretInitialized", log); err != nil {
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

// ParseNetworkSecretInitialized is a log parse operation binding the contract event 0xd1d44220b7bc8275d2a3a1a307706da99997c90e84e42e5d50670da649fcab23.
//
// Solidity: event NetworkSecretInitialized(address initializer)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseNetworkSecretInitialized(log types.Log) (*NetworkEnclaveRegistryNetworkSecretInitialized, error) {
	event := new(NetworkEnclaveRegistryNetworkSecretInitialized)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryNetworkSecretRequestedIterator is returned from FilterNetworkSecretRequested and is used to iterate over the raw logs and unpacked data for NetworkSecretRequested events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretRequestedIterator struct {
	Event *NetworkEnclaveRegistryNetworkSecretRequested // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryNetworkSecretRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryNetworkSecretRequested)
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
		it.Event = new(NetworkEnclaveRegistryNetworkSecretRequested)
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
func (it *NetworkEnclaveRegistryNetworkSecretRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryNetworkSecretRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryNetworkSecretRequested represents a NetworkSecretRequested event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretRequested struct {
	Requester     common.Address
	RequestReport string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNetworkSecretRequested is a free log retrieval operation binding the contract event 0x0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d4301.
//
// Solidity: event NetworkSecretRequested(address indexed requester, string requestReport)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterNetworkSecretRequested(opts *bind.FilterOpts, requester []common.Address) (*NetworkEnclaveRegistryNetworkSecretRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "NetworkSecretRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryNetworkSecretRequestedIterator{contract: _NetworkEnclaveRegistry.contract, event: "NetworkSecretRequested", logs: logs, sub: sub}, nil
}

// WatchNetworkSecretRequested is a free log subscription operation binding the contract event 0x0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d4301.
//
// Solidity: event NetworkSecretRequested(address indexed requester, string requestReport)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchNetworkSecretRequested(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryNetworkSecretRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "NetworkSecretRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryNetworkSecretRequested)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretRequested", log); err != nil {
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

// ParseNetworkSecretRequested is a log parse operation binding the contract event 0x0b0ecdedd12079aa2d6c5e0186026c711cb0c8d04f1b724ba5880fb6328d4301.
//
// Solidity: event NetworkSecretRequested(address indexed requester, string requestReport)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseNetworkSecretRequested(log types.Log) (*NetworkEnclaveRegistryNetworkSecretRequested, error) {
	event := new(NetworkEnclaveRegistryNetworkSecretRequested)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryNetworkSecretRespondedIterator is returned from FilterNetworkSecretResponded and is used to iterate over the raw logs and unpacked data for NetworkSecretResponded events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretRespondedIterator struct {
	Event *NetworkEnclaveRegistryNetworkSecretResponded // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryNetworkSecretRespondedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryNetworkSecretResponded)
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
		it.Event = new(NetworkEnclaveRegistryNetworkSecretResponded)
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
func (it *NetworkEnclaveRegistryNetworkSecretRespondedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryNetworkSecretRespondedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryNetworkSecretResponded represents a NetworkSecretResponded event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryNetworkSecretResponded struct {
	Attester  common.Address
	Requester common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNetworkSecretResponded is a free log retrieval operation binding the contract event 0xb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be1.
//
// Solidity: event NetworkSecretResponded(address indexed attester, address indexed requester)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterNetworkSecretResponded(opts *bind.FilterOpts, attester []common.Address, requester []common.Address) (*NetworkEnclaveRegistryNetworkSecretRespondedIterator, error) {

	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "NetworkSecretResponded", attesterRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryNetworkSecretRespondedIterator{contract: _NetworkEnclaveRegistry.contract, event: "NetworkSecretResponded", logs: logs, sub: sub}, nil
}

// WatchNetworkSecretResponded is a free log subscription operation binding the contract event 0xb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be1.
//
// Solidity: event NetworkSecretResponded(address indexed attester, address indexed requester)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchNetworkSecretResponded(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryNetworkSecretResponded, attester []common.Address, requester []common.Address) (event.Subscription, error) {

	var attesterRule []interface{}
	for _, attesterItem := range attester {
		attesterRule = append(attesterRule, attesterItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "NetworkSecretResponded", attesterRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryNetworkSecretResponded)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretResponded", log); err != nil {
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

// ParseNetworkSecretResponded is a log parse operation binding the contract event 0xb869e23ebc7c717d76e345eee8ec282612603e45c44f7ae5494b197c8d9d1be1.
//
// Solidity: event NetworkSecretResponded(address indexed attester, address indexed requester)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseNetworkSecretResponded(log types.Log) (*NetworkEnclaveRegistryNetworkSecretResponded, error) {
	event := new(NetworkEnclaveRegistryNetworkSecretResponded)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "NetworkSecretResponded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryOwnershipTransferredIterator struct {
	Event *NetworkEnclaveRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistryOwnershipTransferred)
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
		it.Event = new(NetworkEnclaveRegistryOwnershipTransferred)
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
func (it *NetworkEnclaveRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NetworkEnclaveRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistryOwnershipTransferredIterator{contract: _NetworkEnclaveRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistryOwnershipTransferred)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*NetworkEnclaveRegistryOwnershipTransferred, error) {
	event := new(NetworkEnclaveRegistryOwnershipTransferred)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistrySequencerEnclaveGrantedIterator is returned from FilterSequencerEnclaveGranted and is used to iterate over the raw logs and unpacked data for SequencerEnclaveGranted events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveGrantedIterator struct {
	Event *NetworkEnclaveRegistrySequencerEnclaveGranted // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistrySequencerEnclaveGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistrySequencerEnclaveGranted)
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
		it.Event = new(NetworkEnclaveRegistrySequencerEnclaveGranted)
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
func (it *NetworkEnclaveRegistrySequencerEnclaveGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistrySequencerEnclaveGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistrySequencerEnclaveGranted represents a SequencerEnclaveGranted event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveGranted struct {
	EnclaveID common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSequencerEnclaveGranted is a free log retrieval operation binding the contract event 0xfe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e760936.
//
// Solidity: event SequencerEnclaveGranted(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterSequencerEnclaveGranted(opts *bind.FilterOpts) (*NetworkEnclaveRegistrySequencerEnclaveGrantedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "SequencerEnclaveGranted")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistrySequencerEnclaveGrantedIterator{contract: _NetworkEnclaveRegistry.contract, event: "SequencerEnclaveGranted", logs: logs, sub: sub}, nil
}

// WatchSequencerEnclaveGranted is a free log subscription operation binding the contract event 0xfe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e760936.
//
// Solidity: event SequencerEnclaveGranted(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchSequencerEnclaveGranted(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistrySequencerEnclaveGranted) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "SequencerEnclaveGranted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistrySequencerEnclaveGranted)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveGranted", log); err != nil {
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

// ParseSequencerEnclaveGranted is a log parse operation binding the contract event 0xfe64c7181f0fc60e300dc02cca368cdfa94d7ca45902de3b9a9d80070e760936.
//
// Solidity: event SequencerEnclaveGranted(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseSequencerEnclaveGranted(log types.Log) (*NetworkEnclaveRegistrySequencerEnclaveGranted, error) {
	event := new(NetworkEnclaveRegistrySequencerEnclaveGranted)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NetworkEnclaveRegistrySequencerEnclaveRevokedIterator is returned from FilterSequencerEnclaveRevoked and is used to iterate over the raw logs and unpacked data for SequencerEnclaveRevoked events raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveRevokedIterator struct {
	Event *NetworkEnclaveRegistrySequencerEnclaveRevoked // Event containing the contract specifics and raw log

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
func (it *NetworkEnclaveRegistrySequencerEnclaveRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
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
		it.Event = new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
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
func (it *NetworkEnclaveRegistrySequencerEnclaveRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkEnclaveRegistrySequencerEnclaveRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkEnclaveRegistrySequencerEnclaveRevoked represents a SequencerEnclaveRevoked event raised by the NetworkEnclaveRegistry contract.
type NetworkEnclaveRegistrySequencerEnclaveRevoked struct {
	EnclaveID common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSequencerEnclaveRevoked is a free log retrieval operation binding the contract event 0x0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b47.
//
// Solidity: event SequencerEnclaveRevoked(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) FilterSequencerEnclaveRevoked(opts *bind.FilterOpts) (*NetworkEnclaveRegistrySequencerEnclaveRevokedIterator, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.FilterLogs(opts, "SequencerEnclaveRevoked")
	if err != nil {
		return nil, err
	}
	return &NetworkEnclaveRegistrySequencerEnclaveRevokedIterator{contract: _NetworkEnclaveRegistry.contract, event: "SequencerEnclaveRevoked", logs: logs, sub: sub}, nil
}

// WatchSequencerEnclaveRevoked is a free log subscription operation binding the contract event 0x0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b47.
//
// Solidity: event SequencerEnclaveRevoked(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) WatchSequencerEnclaveRevoked(opts *bind.WatchOpts, sink chan<- *NetworkEnclaveRegistrySequencerEnclaveRevoked) (event.Subscription, error) {

	logs, sub, err := _NetworkEnclaveRegistry.contract.WatchLogs(opts, "SequencerEnclaveRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
				if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveRevoked", log); err != nil {
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

// ParseSequencerEnclaveRevoked is a log parse operation binding the contract event 0x0f279980343c7ca542fde9fa5396555068efb5cd560d9cf9c191aa2911079b47.
//
// Solidity: event SequencerEnclaveRevoked(address enclaveID)
func (_NetworkEnclaveRegistry *NetworkEnclaveRegistryFilterer) ParseSequencerEnclaveRevoked(log types.Log) (*NetworkEnclaveRegistrySequencerEnclaveRevoked, error) {
	event := new(NetworkEnclaveRegistrySequencerEnclaveRevoked)
	if err := _NetworkEnclaveRegistry.contract.UnpackLog(event, "SequencerEnclaveRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
