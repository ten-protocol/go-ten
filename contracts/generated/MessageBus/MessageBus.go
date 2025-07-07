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
	_ = abi.ConvertType
)

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint64
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// MessageBusMetaData contains all meta data concerning the MessageBus contract.
var MessageBusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublishFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"withdrawal\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feesAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b50610018610025565b610020610025565b610104565b5f61002e6100c5565b805490915068010000000000000000900460ff16156100605760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100c25780546001600160401b0319166001600160401b0390811782556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916100b9916100ef565b60405180910390a15b50565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6001600160401b0382168152602081016100e9565b6117d7806101115f395ff3fe6080604052600436106100ba575f3560e01c80638da5cb5b11610073578063e30c39781161004d578063e30c3978146101e6578063ea4cf970146101fa578063f2fde38b14610219576100ba565b80638da5cb5b1461017a57806391643fdd1461019b578063c0c53b8b146101c7576100ba565b80632e1a0b8e116100a45780632e1a0b8e14610131578063715018a61461015057806379ba509714610166576100ba565b8062a1b815146100e75780630d3fd67c14610111575b3480156100c5575f5ffd5b5060405162461bcd60e51b81526004016100de90610b36565b60405180910390fd5b3480156100f2575f5ffd5b506100fb610238565b6040516101089190610b4e565b60405180910390f35b61012461011f366004610bf6565b6102c1565b6040516101089190610c7e565b34801561013c575f5ffd5b506100fb61014b366004610ca5565b6103c9565b34801561015b575f5ffd5b50610164610426565b005b348015610171575f5ffd5b50610164610446565b348015610185575f5ffd5b5061018e610485565b6040516101089190610cf6565b3480156101a6575f5ffd5b506101ba6101b5366004610ca5565b6104b9565b6040516101089190610d0c565b3480156101d2575f5ffd5b506101646101e1366004610d2e565b610509565b3480156101f1575f5ffd5b5061018e6106aa565b348015610205575f5ffd5b50610164610214366004610d85565b6106d2565b348015610224575f5ffd5b50610164610233366004610dd7565b610817565b600354604080517f1a90a21900000000000000000000000000000000000000000000000000000000815290515f926001600160a01b031691631a90a2199160048083019260209291908290030181865afa158015610298573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102bc9190610dff565b905090565b6003545f906001600160a01b031615610372575f6102dd610238565b9050803410156102ff5760405162461bcd60e51b81526004016100de90610e76565b6003546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610349576040519150601f19603f3d011682016040523d82523d5f602084013e61034e565b606091505b505090508061036f5760405162461bcd60e51b81526004016100de90610ede565b50505b61037b336108a9565b90507fd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061338288888888886040516103b89796959493929190610f2e565b60405180910390a195945050505050565b5f5f826040516020016103dc91906110cb565b60408051601f1981840301815291815281516020928301205f818152928390529120549091508061041f5760405162461bcd60e51b81526004016100de90611134565b9392505050565b61042e610906565b60405162461bcd60e51b81526004016100de9061119c565b33806104506106aa565b6001600160a01b031614610479578060405163118cdaa760e01b81526004016100de9190610cf6565b6104828161093a565b50565b5f807f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c1993005b546001600160a01b031692915050565b5f5f826040516020016104cc91906110cb565b60408051601f1981840301815291815281516020928301205f8181529283905291205490915080158015906105015750428111155b949350505050565b5f610512610983565b805490915060ff68010000000000000000820416159067ffffffffffffffff165f8115801561053e5750825b90505f8267ffffffffffffffff16600114801561055a5750303b155b905081158015610568575080155b1561059f576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156105d357845468ff00000000000000001916680100000000000000001785555b6001600160a01b0386166105f95760405162461bcd60e51b81526004016100de906111de565b6001600160a01b03881661061f5760405162461bcd60e51b81526004016100de90611220565b610628886109ad565b6003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03881617905583156106a057845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061069790600190611253565b60405180910390a15b5050505050505050565b5f807f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c006104a9565b5f6106de600130611275565b90506106e8610485565b6001600160a01b0316336001600160a01b0316148061070f5750336001600160a01b038216145b61072b5760405162461bcd60e51b81526004016100de906112ca565b5f61073683426112da565b90505f8460405160200161074a91906110cb565b60408051601f1981840301815291815281516020928301205f818152928390529120549091501561078d5760405162461bcd60e51b81526004016100de90611345565b5f818152602081815260408220849055600191906107ad90880188610dd7565b6001600160a01b0316815260208101919091526040015f908120906107d86080880160608901611355565b63ffffffff1681526020808201929092526040015f908120805460018101825590825291902086916004020161080e8282611773565b50505050505050565b61081f610906565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383169081178255610870610485565b6001600160a01b03167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a35050565b6001600160a01b0381165f908152600260205260408120805467ffffffffffffffff1691600191906108db838561177d565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b3361090f610485565b6001600160a01b031614610938573360405163118cdaa760e01b81526004016100de9190610cf6565b565b7f237e158222e3e6968b72b9db0d8043aacf074ad9f650f0d1606b4d82ee432c00805473ffffffffffffffffffffffffffffffffffffffff1916815561097f826109c6565b5050565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005b92915050565b6109b5610a43565b6109be81610a81565b610482610a92565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b610a4b610a9a565b610938576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a89610a43565b61048281610ab8565b610938610a43565b5f610aa3610983565b5468010000000000000000900460ff16919050565b610ac0610a43565b6001600160a01b038116610479575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100de9190610cf6565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b602080825281016109a781610b02565b805b82525050565b602081016109a78284610b46565b67ffffffffffffffff81165b8114610482575f5ffd5b80356109a781610b5c565b63ffffffff8116610b68565b80356109a781610b7d565b5f5f83601f840112610ba757610ba75f5ffd5b50813567ffffffffffffffff811115610bc157610bc15f5ffd5b602083019150836001820283011115610bdb57610bdb5f5ffd5b9250929050565b60ff8116610b68565b80356109a781610be2565b5f5f5f5f5f60808688031215610c0d57610c0d5f5ffd5b610c178787610b72565b9450610c268760208801610b89565b9350604086013567ffffffffffffffff811115610c4457610c445f5ffd5b610c5088828901610b94565b9350935050610c628760608801610beb565b90509295509295909350565b67ffffffffffffffff8116610b48565b602081016109a78284610c6e565b5f60c08284031215610c9f57610c9f5f5ffd5b50919050565b5f60208284031215610cb857610cb85f5ffd5b813567ffffffffffffffff811115610cd157610cd15f5ffd5b61050184828501610c8c565b5f6001600160a01b0382166109a7565b610b4881610cdd565b602081016109a78284610ced565b801515610b48565b602081016109a78284610d04565b610b6881610cdd565b80356109a781610d1a565b5f5f5f60608486031215610d4357610d435f5ffd5b610d4d8585610d23565b9250610d5c8560208601610d23565b9150610d6b8560408601610d23565b90509250925092565b80610b68565b80356109a781610d74565b5f5f60408385031215610d9957610d995f5ffd5b823567ffffffffffffffff811115610db257610db25f5ffd5b610dbe85828601610c8c565b925050610dce8460208501610d7a565b90509250929050565b5f60208284031215610dea57610dea5f5ffd5b61041f8383610d23565b80516109a781610d74565b5f60208284031215610e1257610e125f5ffd5b61041f8383610df4565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f7373616765000000000000000000000000000000000000000000000000000000602082015290505b60400190565b602080825281016109a781610e1c565b60248152602081017f4661696c656420746f2073656e64206665657320746f206665657320636f6e7481527f726163740000000000000000000000000000000000000000000000000000000060208201529050610e70565b602080825281016109a781610e86565b63ffffffff8116610b48565b82818337505f910152565b818352602083019250610f19828483610efa565b50601f01601f19160190565b60ff8116610b48565b60c08101610f3c828a610ced565b610f496020830189610c6e565b610f566040830188610c6e565b610f636060830187610eee565b8181036080830152610f76818587610f05565b9050610f8560a0830184610f25565b98975050505050505050565b505f6109a76020830183610d23565b505f6109a76020830183610b72565b505f6109a76020830183610b89565b5f808335601e1936859003018112610fd757610fd75f5ffd5b830160208101925035905067ffffffffffffffff811115610ff957610ff95f5ffd5b36819003821315610bdb57610bdb5f5ffd5b505f6109a76020830183610beb565b5f60c083016110298380610f91565b6110338582610ced565b506110416020840184610fa0565b61104e6020860182610c6e565b5061105c6040840184610fa0565b6110696040860182610c6e565b506110776060840184610faf565b6110846060860182610eee565b506110926080840184610fbe565b85830360808701526110a5838284610f05565b925050506110b660a084018461100b565b6110c360a0860182610f25565b509392505050565b6020808252810161041f818461101a565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d697474656481527f2e0000000000000000000000000000000000000000000000000000000000000060208201529050610e70565b602080825281016109a7816110dc565b60348152602081017f556e72656e6f756e6361626c654f776e61626c6532537465703a2063616e6e6f81527f742072656e6f756e6365206f776e65727368697000000000000000000000000060208201529050610e70565b602080825281016109a781611144565b601a8152602081017f4665657320616464726573732063616e6e6f742062652030783000000000000081529050610b30565b602080825281016109a7816111ac565b60148152602081017f43616c6c65722063616e6e6f742062652030783000000000000000000000000081529050610b30565b602080825281016109a7816111ee565b5f6109a78261123d565b90565b67ffffffffffffffff1690565b610b4881611230565b602081016109a7828461124a565b634e487b7160e01b5f52601160045260245ffd5b6001600160a01b039182169190811690828203908111156109a7576109a7611261565b60118152602081017f4e6f74206f776e6572206f722073656c6600000000000000000000000000000081529050610b30565b602080825281016109a781611298565b808201808211156109a7576109a7611261565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f210000000000000000000000000000000000000000000000000000000000000060208201529050610e70565b602080825281016109a7816112ed565b5f60208284031215611368576113685f5ffd5b61041f8383610b89565b5f81356109a781610d1a565b5f6001600160a01b03835b81169019929092169190911792915050565b5f6109a782610cdd565b5f6109a78261139b565b6113b8826113a5565b6113c381835461137e565b8255505050565b5f81356109a781610b5c565b5f7bffffffffffffffff00000000000000000000000000000000000000006113898460a01b90565b5f6109a767ffffffffffffffff831661123d565b61141b826113fe565b6113c38183546113d6565b5f67ffffffffffffffff83611389565b61143f826113fe565b6113c3818354611426565b5f81356109a781610b7d565b5f6bffffffff00000000000000006113898460401b90565b5f63ffffffff82166109a7565b6114848261146e565b6113c3818354611456565b5f808335601e19368590030181126114a8576114a85f5ffd5b8301915050803567ffffffffffffffff8111156114c6576114c65f5ffd5b602082019150600181023603821315610bdb57610bdb5f5ffd5b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061151c57607f821691505b602082108103610c9f57610c9f6114f4565b5f6109a761123a8381565b6115428361152e565b81545f1960089490940293841b1916921b91909117905550565b5f611568818484611539565b505050565b8181101561097f5761157f5f8261155c565b60010161156d565b601f821115611568575f818152602090206020601f850104810160208510156115ad5750805b6115bf6020601f86010483018261156d565b5050505050565b8267ffffffffffffffff8111156115df576115df6114e0565b6115e98254611508565b6115f4828285611587565b505f601f821160018114611626575f831561160f5750848201355b5f19600885021c198116600285021785555061167d565b5f84815260208120601f198516915b828110156116555787850135825560209485019460019092019101611635565b5084821015611671575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b6115688383836115c6565b5f81356109a781610be2565b5f60ff82166109a7565b6116af8261169c565b815460ff191660ff8216176113c3565b8082806116cb81611372565b90506116d781846113af565b505060208301806116e7826113ca565b90506116f38184611412565b505050600181016040830180611708826113ca565b90506117148184611436565b505060608301806117248261144a565b9050611730818461147b565b505050611740608083018361148f565b61174e818360028601611685565b505060a082018061175e82611690565b905061176d81600385016116a6565b50505050565b61097f82826116bf565b67ffffffffffffffff9182169190811690828201908111156109a7576109a761126156fea2646970667358221220eda2193f3930c02738adcab1acddff01bbebd562f51ed0aa7da1290e65655bd864736f6c634300081c0033",
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
	parsed, err := MessageBusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x2e1a0b8e.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MessageBus.Contract.GetMessageTimeOfFinality(&_MessageBus.CallOpts, crossChainMessage)
}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MessageBus *MessageBusCaller) GetPublishFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "getPublishFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MessageBus *MessageBusSession) GetPublishFee() (*big.Int, error) {
	return _MessageBus.Contract.GetPublishFee(&_MessageBus.CallOpts)
}

// GetPublishFee is a free data retrieval call binding the contract method 0x00a1b815.
//
// Solidity: function getPublishFee() view returns(uint256)
func (_MessageBus *MessageBusCallerSession) GetPublishFee() (*big.Int, error) {
	return _MessageBus.Contract.GetPublishFee(&_MessageBus.CallOpts)
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

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MessageBus *MessageBusCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MessageBus *MessageBusSession) PendingOwner() (common.Address, error) {
	return _MessageBus.Contract.PendingOwner(&_MessageBus.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_MessageBus *MessageBusCallerSession) PendingOwner() (common.Address, error) {
	return _MessageBus.Contract.PendingOwner(&_MessageBus.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MessageBus *MessageBusCaller) RenounceOwnership(opts *bind.CallOpts) error {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "renounceOwnership")

	if err != nil {
		return err
	}

	return err

}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MessageBus *MessageBusSession) RenounceOwnership() error {
	return _MessageBus.Contract.RenounceOwnership(&_MessageBus.CallOpts)
}

// RenounceOwnership is a free data retrieval call binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() view returns()
func (_MessageBus *MessageBusCallerSession) RenounceOwnership() error {
	return _MessageBus.Contract.RenounceOwnership(&_MessageBus.CallOpts)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x91643fdd.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MessageBus *MessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MessageBus.Contract.VerifyMessageFinalized(&_MessageBus.CallOpts, crossChainMessage)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MessageBus *MessageBusTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MessageBus *MessageBusSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageBus.Contract.AcceptOwnership(&_MessageBus.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_MessageBus *MessageBusTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _MessageBus.Contract.AcceptOwnership(&_MessageBus.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address caller, address withdrawal, address feesAddress) returns()
func (_MessageBus *MessageBusTransactor) Initialize(opts *bind.TransactOpts, caller common.Address, withdrawal common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "initialize", caller, withdrawal, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address caller, address withdrawal, address feesAddress) returns()
func (_MessageBus *MessageBusSession) Initialize(caller common.Address, withdrawal common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.Initialize(&_MessageBus.TransactOpts, caller, withdrawal, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address caller, address withdrawal, address feesAddress) returns()
func (_MessageBus *MessageBusTransactorSession) Initialize(caller common.Address, withdrawal common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.Initialize(&_MessageBus.TransactOpts, caller, withdrawal, feesAddress)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusSession) PublishMessage(nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0x0d3fd67c.
//
// Solidity: function publishMessage(uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactorSession) PublishMessage(nonce uint64, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusTransactor) StoreCrossChainMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "storeCrossChainMessage", crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MessageBus *MessageBusSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.StoreCrossChainMessage(&_MessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0xea4cf970.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint64,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
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
// Solidity: fallback() returns()
func (_MessageBus *MessageBusTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _MessageBus.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_MessageBus *MessageBusSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_MessageBus *MessageBusTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MessageBus.Contract.Fallback(&_MessageBus.TransactOpts, calldata)
}

// MessageBusInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the MessageBus contract.
type MessageBusInitializedIterator struct {
	Event *MessageBusInitialized // Event containing the contract specifics and raw log

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
func (it *MessageBusInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusInitialized)
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
		it.Event = new(MessageBusInitialized)
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
func (it *MessageBusInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusInitialized represents a Initialized event raised by the MessageBus contract.
type MessageBusInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MessageBus *MessageBusFilterer) FilterInitialized(opts *bind.FilterOpts) (*MessageBusInitializedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MessageBusInitializedIterator{contract: _MessageBus.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_MessageBus *MessageBusFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MessageBusInitialized) (event.Subscription, error) {

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusInitialized)
				if err := _MessageBus.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_MessageBus *MessageBusFilterer) ParseInitialized(log types.Log) (*MessageBusInitialized, error) {
	event := new(MessageBusInitialized)
	if err := _MessageBus.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	Nonce            uint64
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogMessagePublished is a free log retrieval operation binding the contract event 0xd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts) (*MessageBusLogMessagePublishedIterator, error) {

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return &MessageBusLogMessagePublishedIterator{contract: _MessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
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

// ParseLogMessagePublished is a log parse operation binding the contract event 0xd3cfd274dfddb1195699ee44f0ba7aaabf97e75965b5191c2c0bd56776ff2061.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint64 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MessageBus *MessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MessageBusLogMessagePublished, error) {
	event := new(MessageBusLogMessagePublished)
	if err := _MessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MessageBusOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the MessageBus contract.
type MessageBusOwnershipTransferStartedIterator struct {
	Event *MessageBusOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *MessageBusOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusOwnershipTransferStarted)
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
		it.Event = new(MessageBusOwnershipTransferStarted)
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
func (it *MessageBusOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MessageBus contract.
type MessageBusOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MessageBusOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusOwnershipTransferStartedIterator{contract: _MessageBus.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *MessageBusOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusOwnershipTransferStarted)
				if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_MessageBus *MessageBusFilterer) ParseOwnershipTransferStarted(log types.Log) (*MessageBusOwnershipTransferStarted, error) {
	event := new(MessageBusOwnershipTransferStarted)
	if err := _MessageBus.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
