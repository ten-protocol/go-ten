// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MessageBus

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
)

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// MessageBusMetaData contains all meta data concerning the MessageBus contract.
var MessageBusMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061002d61002261003260201b60201c565b61003a60201b60201c565b6100fe565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b611a678061010d6000396000f3fe6080604052600436106100745760003560e01c80638da5cb5b1161004e5780638da5cb5b146101805780639730886d146101ab578063b1454caa146101d4578063f2fde38b14610211576100b4565b80630fcfbd11146100ef57806333a88c721461012c578063715018a614610169576100b4565b366100b4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100ab90610851565b60405180910390fd5b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100e6906108bd565b60405180910390fd5b3480156100fb57600080fd5b506101166004803603810190610111919061090b565b61023a565b604051610123919061096d565b60405180910390f35b34801561013857600080fd5b50610153600480360381019061014e919061090b565b6102cb565b60405161016091906109a3565b60405180910390f35b34801561017557600080fd5b5061017e610316565b005b34801561018c57600080fd5b5061019561032a565b6040516101a291906109ff565b60405180910390f35b3480156101b757600080fd5b506101d260048036038101906101cd9190610a46565b610353565b005b3480156101e057600080fd5b506101fb60048036038101906101f69190610b7c565b6104c2565b6040516102089190610c27565b60405180910390f35b34801561021d57600080fd5b5061023860048036038101906102339190610c6e565b61051b565b005b6000808260405160200161024e9190610eea565b604051602081830303815290604052805190602001209050600060016000838152602001908152602001600020549050600081116102c1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102b890610f7e565b60405180910390fd5b8092505050919050565b600080826040516020016102df9190610eea565b6040516020818303038152906040528051906020012090504260016000838152602001908152602001600020541015915050919050565b61031e61059e565b610328600061061c565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b61035b61059e565b600081426103699190610fcd565b905060008360405160200161037e9190610eea565b60405160208183030381529060405280519060200120905060006001600083815260200190815260200160002054146103ec576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103e390611073565b60405180910390fd5b8160016000838152602001908152602001600020819055506002600085600001602081019061041b9190610c6e565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600085606001602081019061046a9190611093565b63ffffffff1663ffffffff16815260200190815260200160002084908060018154018082558091505060019003906000526020600020906004020160009091909190915081816104ba9190611823565b505050505050565b60006104cd336106e0565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef7759373382888888888860405161050a979695949392919061188d565b60405180910390a195945050505050565b61052361059e565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610592576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058990611969565b60405180910390fd5b61059b8161061c565b50565b6105a66107c6565b73ffffffffffffffffffffffffffffffffffffffff166105c461032a565b73ffffffffffffffffffffffffffffffffffffffff161461061a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610611906119d5565b60405180910390fd5b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900467ffffffffffffffff1690506001600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282829054906101000a900467ffffffffffffffff1661079b91906119f5565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b600033905090565b600082825260208201905092915050565b7f74686520576f726d686f6c6520636f6e747261637420646f6573206e6f74206160008201527f6363657074206173736574730000000000000000000000000000000000000000602082015250565b600061083b602c836107ce565b9150610846826107df565b604082019050919050565b6000602082019050818103600083015261086a8161082e565b9050919050565b7f756e737570706f72746564000000000000000000000000000000000000000000600082015250565b60006108a7600b836107ce565b91506108b282610871565b602082019050919050565b600060208201905081810360008301526108d68161089a565b9050919050565b600080fd5b600080fd5b600080fd5b600060c08284031215610902576109016108e7565b5b81905092915050565b600060208284031215610921576109206108dd565b5b600082013567ffffffffffffffff81111561093f5761093e6108e2565b5b61094b848285016108ec565b91505092915050565b6000819050919050565b61096781610954565b82525050565b6000602082019050610982600083018461095e565b92915050565b60008115159050919050565b61099d81610988565b82525050565b60006020820190506109b86000830184610994565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006109e9826109be565b9050919050565b6109f9816109de565b82525050565b6000602082019050610a1460008301846109f0565b92915050565b610a2381610954565b8114610a2e57600080fd5b50565b600081359050610a4081610a1a565b92915050565b60008060408385031215610a5d57610a5c6108dd565b5b600083013567ffffffffffffffff811115610a7b57610a7a6108e2565b5b610a87858286016108ec565b9250506020610a9885828601610a31565b9150509250929050565b600063ffffffff82169050919050565b610abb81610aa2565b8114610ac657600080fd5b50565b600081359050610ad881610ab2565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610b0357610b02610ade565b5b8235905067ffffffffffffffff811115610b2057610b1f610ae3565b5b602083019150836001820283011115610b3c57610b3b610ae8565b5b9250929050565b600060ff82169050919050565b610b5981610b43565b8114610b6457600080fd5b50565b600081359050610b7681610b50565b92915050565b600080600080600060808688031215610b9857610b976108dd565b5b6000610ba688828901610ac9565b9550506020610bb788828901610ac9565b945050604086013567ffffffffffffffff811115610bd857610bd76108e2565b5b610be488828901610aed565b93509350506060610bf788828901610b67565b9150509295509295909350565b600067ffffffffffffffff82169050919050565b610c2181610c04565b82525050565b6000602082019050610c3c6000830184610c18565b92915050565b610c4b816109de565b8114610c5657600080fd5b50565b600081359050610c6881610c42565b92915050565b600060208284031215610c8457610c836108dd565b5b6000610c9284828501610c59565b91505092915050565b6000610caa6020840184610c59565b905092915050565b610cbb816109de565b82525050565b610cca81610c04565b8114610cd557600080fd5b50565b600081359050610ce781610cc1565b92915050565b6000610cfc6020840184610cd8565b905092915050565b610d0d81610c04565b82525050565b6000610d226020840184610ac9565b905092915050565b610d3381610aa2565b82525050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610d6557610d64610d43565b5b83810192508235915060208301925067ffffffffffffffff821115610d8d57610d8c610d39565b5b600182023603831315610da357610da2610d3e565b5b509250929050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b6000610de88385610dab565b9350610df5838584610dbc565b610dfe83610dcb565b840190509392505050565b6000610e186020840184610b67565b905092915050565b610e2981610b43565b82525050565b600060c08301610e426000840184610c9b565b610e4f6000860182610cb2565b50610e5d6020840184610ced565b610e6a6020860182610d04565b50610e786040840184610d13565b610e856040860182610d2a565b50610e936060840184610d13565b610ea06060860182610d2a565b50610eae6080840184610d48565b8583036080870152610ec1838284610ddc565b92505050610ed260a0840184610e09565b610edf60a0860182610e20565b508091505092915050565b60006020820190508181036000830152610f048184610e2f565b905092915050565b7f54686973206d65737361676520776173206e65766572207375626d697474656460008201527f2e00000000000000000000000000000000000000000000000000000000000000602082015250565b6000610f686021836107ce565b9150610f7382610f0c565b604082019050919050565b60006020820190508181036000830152610f9781610f5b565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610fd882610954565b9150610fe383610954565b9250828201905080821115610ffb57610ffa610f9e565b5b92915050565b7f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636560008201527f2100000000000000000000000000000000000000000000000000000000000000602082015250565b600061105d6021836107ce565b915061106882611001565b604082019050919050565b6000602082019050818103600083015261108c81611050565b9050919050565b6000602082840312156110a9576110a86108dd565b5b60006110b784828501610ac9565b91505092915050565b600081356110cd81610c42565b80915050919050565b60008160001b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff611103846110d6565b9350801983169250808416831791505092915050565b6000819050919050565b600061113e611139611134846109be565b611119565b6109be565b9050919050565b600061115082611123565b9050919050565b600061116282611145565b9050919050565b6000819050919050565b61117c82611157565b61118f61118882611169565b83546110e3565b8255505050565b600081356111a381610cc1565b80915050919050565b60008160a01b9050919050565b60007bffffffffffffffff00000000000000000000000000000000000000006111e1846111ac565b9350801983169250808416831791505092915050565b600061121261120d61120884610c04565b611119565b610c04565b9050919050565b6000819050919050565b61122c826111f7565b61123f61123882611219565b83546111b9565b8255505050565b6000813561125381610ab2565b80915050919050565b60008160e01b9050919050565b60007fffffffff000000000000000000000000000000000000000000000000000000006112958461125c565b9350801983169250808416831791505092915050565b60006112c66112c16112bc84610aa2565b611119565b610aa2565b9050919050565b6000819050919050565b6112e0826112ab565b6112f36112ec826112cd565b8354611269565b8255505050565b600063ffffffff61130a846110d6565b9350801983169250808416831791505092915050565b611329826112ab565b61133c611335826112cd565b83546112fa565b8255505050565b600080fd5b600080fd5b600080fd5b6000808335600160200384360303811261136f5761136e611343565b5b80840192508235915067ffffffffffffffff82111561139157611390611348565b5b6020830192506001820236038313156113ad576113ac61134d565b5b509250929050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061143657607f821691505b602082108103611449576114486113ef565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026114b17fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611474565b6114bb8683611474565b95508019841693508086168417925050509392505050565b60006114ee6114e96114e484610954565b611119565b610954565b9050919050565b6000819050919050565b611508836114d3565b61151c611514826114f5565b848454611481565b825550505050565b600090565b611531611524565b61153c8184846114ff565b505050565b5b8181101561156057611555600082611529565b600181019050611542565b5050565b601f8211156115a5576115768161144f565b61157f84611464565b8101602085101561158e578190505b6115a261159a85611464565b830182611541565b50505b505050565b600082821c905092915050565b60006115c8600019846008026115aa565b1980831691505092915050565b60006115e183836115b7565b9150826002028217905092915050565b6115fb83836113b5565b67ffffffffffffffff811115611614576116136113c0565b5b61161e825461141e565b611629828285611564565b6000601f8311600181146116585760008415611646578287013590505b61165085826115d5565b8655506116b8565b601f1984166116668661144f565b60005b8281101561168e57848901358255600182019150602085019450602081019050611669565b868310156116ab57848901356116a7601f8916826115b7565b8355505b6001600288020188555050505b50505050505050565b6116cc8383836115f1565b505050565b600081356116de81610b50565b80915050919050565b600060ff6116f4846110d6565b9350801983169250808416831791505092915050565b600061172561172061171b84610b43565b611119565b610b43565b9050919050565b6000819050919050565b61173f8261170a565b61175261174b8261172c565b83546116e7565b8255505050565b60008101600083018061176b816110c0565b90506117778184611173565b50505060008101602083018061178c81611196565b90506117988184611223565b5050506000810160408301806117ad81611246565b90506117b981846112d7565b5050506001810160608301806117ce81611246565b90506117da8184611320565b50505060028101608083016117ef8185611352565b6117fa8183866116c1565b505050506003810160a0830180611810816116d1565b905061181c8184611736565b5050505050565b61182d8282611759565b5050565b61183a81610aa2565b82525050565b600082825260208201905092915050565b600061185d8385611840565b935061186a838584610dbc565b61187383610dcb565b840190509392505050565b61188781610b43565b82525050565b600060c0820190506118a2600083018a6109f0565b6118af6020830189610c18565b6118bc6040830188611831565b6118c96060830187611831565b81810360808301526118dc818587611851565b90506118eb60a083018461187e565b98975050505050505050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b60006119536026836107ce565b915061195e826118f7565b604082019050919050565b6000602082019050818103600083015261198281611946565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006119bf6020836107ce565b91506119ca82611989565b602082019050919050565b600060208201905081810360008301526119ee816119b2565b9050919050565b6000611a0082610c04565b9150611a0b83610c04565b9250828201905067ffffffffffffffff811115611a2b57611a2a610f9e565b5b9291505056fea26469706673582212202d3aa8bcdf198912723dcd799dd5e260ffbc01024084482688441199db8b103564736f6c63430008110033",
}

// MessageBusABI is the input ABI used to generate the binding from.
// Deprecated: Use MessageBusMetaData.ABI instead.
var MessageBusABI = MessageBusMetaData.ABI

// MessageBusBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MessageBusMetaData.Bin instead.
var MessageBusBin = MessageBusMetaData.Bin

// DeployMessageBus deploys a new Ethereum contract, binding an instance of MessageBus to it.
func DeployMessageBus(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MessageBus, error) {
	parsed, err := MessageBusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MessageBusBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MessageBus{MessageBusCaller: MessageBusCaller{contract: contract}, MessageBusTransactor: MessageBusTransactor{contract: contract}, MessageBusFilterer: MessageBusFilterer{contract: contract}}, nil
}

// MessageBus is an auto generated Go binding around an Ethereum contract.
type MessageBus struct {
	MessageBusCaller     // Read-only binding to the contract
	MessageBusTransactor // Write-only binding to the contract
	MessageBusFilterer   // Log filterer for contract events
}

// MessageBusCaller is an auto generated read-only Go binding around an Ethereum contract.
type MessageBusCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MessageBusTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MessageBusFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MessageBusSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MessageBusSession struct {
	Contract     *MessageBus       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MessageBusCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MessageBusCallerSession struct {
	Contract *MessageBusCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// MessageBusTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MessageBusTransactorSession struct {
	Contract     *MessageBusTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MessageBusRaw is an auto generated low-level Go binding around an Ethereum contract.
type MessageBusRaw struct {
	Contract *MessageBus // Generic contract binding to access the raw methods on
}

// MessageBusCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MessageBusCallerRaw struct {
	Contract *MessageBusCaller // Generic read-only contract binding to access the raw methods on
}

// MessageBusTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MessageBusTransactorRaw struct {
	Contract *MessageBusTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMessageBus creates a new instance of MessageBus, bound to a specific deployed contract.
func NewMessageBus(address common.Address, backend bind.ContractBackend) (*MessageBus, error) {
	contract, err := bindMessageBus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MessageBus{MessageBusCaller: MessageBusCaller{contract: contract}, MessageBusTransactor: MessageBusTransactor{contract: contract}, MessageBusFilterer: MessageBusFilterer{contract: contract}}, nil
}

// NewMessageBusCaller creates a new read-only instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusCaller(address common.Address, caller bind.ContractCaller) (*MessageBusCaller, error) {
	contract, err := bindMessageBus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MessageBusCaller{contract: contract}, nil
}

// NewMessageBusTransactor creates a new write-only instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusTransactor(address common.Address, transactor bind.ContractTransactor) (*MessageBusTransactor, error) {
	contract, err := bindMessageBus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MessageBusTransactor{contract: contract}, nil
}

// NewMessageBusFilterer creates a new log filterer instance of MessageBus, bound to a specific deployed contract.
func NewMessageBusFilterer(address common.Address, filterer bind.ContractFilterer) (*MessageBusFilterer, error) {
	contract, err := bindMessageBus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MessageBusFilterer{contract: contract}, nil
}

// bindMessageBus binds a generic wrapper to an already deployed contract.
func bindMessageBus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MessageBusABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageBus *MessageBusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageBus.Contract.MessageBusCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageBus *MessageBusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.Contract.MessageBusTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageBus *MessageBusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageBus.Contract.MessageBusTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MessageBus *MessageBusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MessageBus.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MessageBus *MessageBusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MessageBus *MessageBusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MessageBus.Contract.contract.Transact(opts, method, params...)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MessageBus *MessageBusCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MessageBus *MessageBusSession) Owner() (common.Address, error) {
	return _MessageBus.Contract.Owner(&_MessageBus.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MessageBus *MessageBusCallerSession) Owner() (common.Address, error) {
	return _MessageBus.Contract.Owner(&_MessageBus.CallOpts)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MessageBus *MessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MessageBus *MessageBusSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MessageBus *MessageBusTransactorSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MessageBus *MessageBusTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MessageBus *MessageBusSession) RenounceOwnership() (*types.Transaction, error) {
	return _MessageBus.Contract.RenounceOwnership(&_MessageBus.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MessageBus *MessageBusTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MessageBus.Contract.RenounceOwnership(&_MessageBus.TransactOpts)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactor) StoreCrossChainMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "storeCrossChainMessage", crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.StoreCrossChainMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactorSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.StoreCrossChainMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MessageBus *MessageBusTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MessageBus *MessageBusSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.TransferOwnership(&_MessageBus.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MessageBus *MessageBusTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.TransferOwnership(&_MessageBus.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MessageBus *MessageBusTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _MessageBus.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MessageBus *MessageBusSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MessageBus *MessageBusTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MessageBus *MessageBusTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MessageBus *MessageBusSession) Receive() (*types.Transaction, error) {
	return _MessageBus.Contract.Receive(&_MessageBus.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MessageBus *MessageBusTransactorSession) Receive() (*types.Transaction, error) {
	return _MessageBus.Contract.Receive(&_MessageBus.TransactOpts)
}

// MessageBusLogMessagePublishedIterator is returned from FilterLogMessagePublished and is used to iterate over the raw logs and unpacked data for LogMessagePublished events raised by the MessageBus contract.
type MessageBusLogMessagePublishedIterator struct {
	Event *MessageBusLogMessagePublished // Event containing the contract specifics and raw log

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
func (it *MessageBusLogMessagePublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusLogMessagePublished)
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
		it.Event = new(MessageBusLogMessagePublished)
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
func (it *MessageBusLogMessagePublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusLogMessagePublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusLogMessagePublished represents a LogMessagePublished event raised by the MessageBus contract.
type MessageBusLogMessagePublished struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogMessagePublished is a free log retrieval operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts) (*MessageBusLogMessagePublishedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return &MessageBusLogMessagePublishedIterator{contract: _MessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) WatchLogMessagePublished(opts *bind.WatchOpts, sink chan<- *MessageBusLogMessagePublished) (event.Subscription, error) {

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusLogMessagePublished)
				if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
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

// ParseLogMessagePublished is a log parse operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MessageBusLogMessagePublished, error) {
	event := new(MessageBusLogMessagePublished)
	if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MessageBus contract.
type MessageBusOwnershipTransferredIterator struct {
	Event *MessageBusOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MessageBusOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusOwnershipTransferred)
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
		it.Event = new(MessageBusOwnershipTransferred)
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
func (it *MessageBusOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusOwnershipTransferred represents a OwnershipTransferred event raised by the MessageBus contract.
type MessageBusOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MessageBusOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusOwnershipTransferredIterator{contract: _MessageBus.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MessageBusOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusOwnershipTransferred)
				if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseOwnershipTransferred(log types.Log) (*MessageBusOwnershipTransferred, error) {
	event := new(MessageBusOwnershipTransferred)
	if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
