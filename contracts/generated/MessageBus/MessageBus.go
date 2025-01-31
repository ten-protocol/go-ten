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
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// MessageBusMetaData contains all meta data concerning the MessageBus contract.
var MessageBusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NativeDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"name\":\"ValueTransfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublishFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feesAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"notifyDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"receiveValueFromL2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"retrieveAllFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendValueToL2\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b5061001933610026565b610021610096565b610148565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100e65760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146101455780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b611add806101555f395ff3fe6080604052600436106100d4575f3560e01c8063715018a61161007e57806399a3ad211161005857806399a3ad21146102a2578063ab53bddc146102c1578063b1454caa146102e0578063f2fde38b1461030057610144565b8063715018a61461022b5780638da5cb5b1461023f5780639730886d1461028357610144565b8063346633fb116100af578063346633fb146101da57806336d2da90146101ed578063485cc9551461020c57610144565b8062a1b815146101655780630fcfbd111461018f57806333a88c72146101ae57610144565b36610144576040517f346633fb000000000000000000000000000000000000000000000000000000008152309063346633fb9034906101199033908390600401610db1565b5f604051808303818588803b158015610130575f5ffd5b505af1158015610142573d5f5f3e3d5ffd5b005b60405162461bcd60e51b815260040161015c90610e00565b60405180910390fd5b348015610170575f5ffd5b5061017961031f565b6040516101869190610e10565b60405180910390f35b34801561019a575f5ffd5b506101796101a9366004610e37565b6103a8565b3480156101b9575f5ffd5b506101cd6101c8366004610e37565b610405565b6040516101869190610e77565b6101426101e8366004610eb4565b610455565b3480156101f8575f5ffd5b50610142610207366004610eea565b61059e565b348015610217575f5ffd5b50610142610226366004610f07565b610619565b348015610236575f5ffd5b50610142610782565b34801561024a575f5ffd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166040516101869190610f34565b34801561028e575f5ffd5b5061014261029d366004610f42565b610795565b3480156102ad575f5ffd5b506101426102bc366004610eb4565b6108fa565b3480156102cc575f5ffd5b506101426102db366004610eb4565b610996565b6102f36102ee366004611004565b610a5e565b604051610186919061108c565b34801561030b575f5ffd5b5061014261031a366004610eea565b610b66565b600354604080517f1a90a21900000000000000000000000000000000000000000000000000000000815290515f926001600160a01b031691631a90a2199160048083019260209291908290030181865afa15801561037f573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103a391906110a5565b905090565b5f5f826040516020016103bb9190611257565b60408051601f1981840301815291815281516020928301205f81815292839052912054909150806103fe5760405162461bcd60e51b815260040161015c906112c2565b9392505050565b5f5f826040516020016104189190611257565b60408051601f1981840301815291815281516020928301205f81815292839052912054909150801580159061044d5750428111155b949350505050565b5f3411801561046357508034145b61047f5760405162461bcd60e51b815260040161015c9061132a565b60035434906001600160a01b03161561053f575f61049b61031f565b9050803410156104bd5760405162461bcd60e51b815260040161015c9061136a565b6104c7813461138e565b6003546040519193505f916001600160a01b039091169083908381818185875af1925050503d805f8114610516576040519150601f19603f3d011682016040523d82523d5f602084013e61051b565b606091505b505090508061053c5760405162461bcd60e51b815260040161015c906113f9565b50505b5f61054933610bbc565b9050836001600160a01b0316336001600160a01b03167f50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b68484604051610590929190611409565b60405180910390a350505050565b6105a6610c19565b5f816001600160a01b0316476040515f6040518083038185875af1925050503d805f81146105ef576040519150601f19603f3d011682016040523d82523d5f602084013e6105f4565b606091505b50509050806106155760405162461bcd60e51b815260040161015c90611456565b5050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156106635750825b90505f8267ffffffffffffffff16600114801561067f5750303b155b90508115801561068d575080155b156106c4576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156106f857845468ff00000000000000001916680100000000000000001785555b61070187610c8d565b6003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038816179055831561077957845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29061077090600190611489565b60405180910390a15b50505050505050565b61078a610c19565b6107935f610c9e565b565b5f6107a1600130611497565b90506107d47f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b031614806107fb5750336001600160a01b038216145b6108175760405162461bcd60e51b815260040161015c906114ec565b5f61082283426114fc565b90505f846040516020016108369190611257565b60408051601f1981840301815291815281516020928301205f81815292839052912054909150156108795760405162461bcd60e51b815260040161015c90611567565b5f8181526020818152604082208490556001919061089990880188610eea565b6001600160a01b0316815260208101919091526040015f908120906108c46080880160608901611577565b63ffffffff1681526020808201929092526040015f908120805460018101825590825291902086916004020161077982826119ae565b610902610c19565b804710156109225760405162461bcd60e51b815260040161015c9061136a565b5f826001600160a01b0316826040515f6040518083038185875af1925050503d805f811461096b576040519150601f19603f3d011682016040523d82523d5f602084013e610970565b606091505b50509050806109915760405162461bcd60e51b815260040161015c90611456565b505050565b5f6109a2600130611497565b90506109d57f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b031614806109fc5750336001600160a01b038216145b610a185760405162461bcd60e51b815260040161015c906114ec565b826001600160a01b03167fcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d383604051610a519190610e10565b60405180910390a2505050565b6003545f906001600160a01b031615610b0f575f610a7a61031f565b905080341015610a9c5760405162461bcd60e51b815260040161015c90611a10565b6003546040515f916001600160a01b03169083908381818185875af1925050503d805f8114610ae6576040519150601f19603f3d011682016040523d82523d5f602084013e610aeb565b606091505b5050905080610b0c5760405162461bcd60e51b815260040161015c906113f9565b50505b610b1833610bbc565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef77593733828888888888604051610b559796959493929190611a20565b60405180910390a195945050505050565b610b6e610c19565b6001600160a01b038116610bb0575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161015c9190610f34565b610bb981610c9e565b50565b6001600160a01b0381165f908152600260205260408120805467ffffffffffffffff169160019190610bee8385611a83565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b33610c4b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b03161461079357336040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161015c9190610f34565b610c95610d1b565b610bb981610d82565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610793576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b6e610d1b565b5f6001600160a01b0382165b92915050565b610da581610d8a565b82525050565b80610da5565b60408101610dbf8285610d9c565b6103fe6020830184610dab565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b60208082528101610d9681610dcc565b60208101610d968284610dab565b5f60c08284031215610e3157610e315f5ffd5b50919050565b5f60208284031215610e4a57610e4a5f5ffd5b813567ffffffffffffffff811115610e6357610e635f5ffd5b61044d84828501610e1e565b801515610da5565b60208101610d968284610e6f565b610e8e81610d8a565b8114610bb9575f5ffd5b8035610d9681610e85565b80610e8e565b8035610d9681610ea3565b5f5f60408385031215610ec857610ec85f5ffd5b610ed28484610e98565b9150610ee18460208501610ea9565b90509250929050565b5f60208284031215610efd57610efd5f5ffd5b6103fe8383610e98565b5f5f60408385031215610f1b57610f1b5f5ffd5b610f258484610e98565b9150610ee18460208501610e98565b60208101610d968284610d9c565b5f5f60408385031215610f5657610f565f5ffd5b823567ffffffffffffffff811115610f6f57610f6f5f5ffd5b610f7b85828601610e1e565b925050610ee18460208501610ea9565b63ffffffff8116610e8e565b8035610d9681610f8b565b5f5f83601f840112610fb557610fb55f5ffd5b50813567ffffffffffffffff811115610fcf57610fcf5f5ffd5b602083019150836001820283011115610fe957610fe95f5ffd5b9250929050565b60ff8116610e8e565b8035610d9681610ff0565b5f5f5f5f5f6080868803121561101b5761101b5f5ffd5b6110258787610f97565b94506110348760208801610f97565b9350604086013567ffffffffffffffff811115611052576110525f5ffd5b61105e88828901610fa2565b93509350506110708760608801610ff9565b90509295509295909350565b67ffffffffffffffff8116610da5565b60208101610d96828461107c565b8051610d9681610ea3565b5f602082840312156110b8576110b85f5ffd5b6103fe838361109a565b505f610d966020830183610e98565b67ffffffffffffffff8116610e8e565b8035610d96816110d1565b505f610d9660208301836110e1565b505f610d966020830183610f97565b63ffffffff8116610da5565b5f808335601e193685900301811261112f5761112f5f5ffd5b830160208101925035905067ffffffffffffffff811115611151576111515f5ffd5b36819003821315610fe957610fe95f5ffd5b82818337505f910152565b818352602083019250611182828483611163565b50601f01601f19160190565b505f610d966020830183610ff9565b60ff8116610da5565b5f60c083016111b583806110c2565b6111bf8582610d9c565b506111cd60208401846110ec565b6111da602086018261107c565b506111e860408401846110fb565b6111f5604086018261110a565b5061120360608401846110fb565b611210606086018261110a565b5061121e6080840184611116565b858303608087015261123183828461116e565b9250505061124260a084018461118e565b61124f60a086018261119d565b509392505050565b602080825281016103fe81846111a6565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d697474656481527f2e00000000000000000000000000000000000000000000000000000000000000602082015290505b60400190565b60208082528101610d9681611268565b60308152602081017f417474656d7074696e6720746f2073656e642076616c756520776974686f757481527f2070726f766964696e6720457468657200000000000000000000000000000000602082015290506112bc565b60208082528101610d96816112d2565b60208082527f496e73756666696369656e742066756e647320746f2073656e642076616c75659101908152610dfa565b60208082528101610d968161133a565b634e487b7160e01b5f52601160045260245ffd5b81810381811115610d9657610d9661137a565b60248152602081017f4661696c656420746f2073656e64206665657320746f206665657320636f6e7481527f7261637400000000000000000000000000000000000000000000000000000000602082015290506112bc565b60208082528101610d96816113a1565b604081016114178285610dab565b6103fe602083018461107c565b60148152602081017f6661696c65642073656e64696e672076616c756500000000000000000000000081529050610dfa565b60208082528101610d9681611424565b5f610d9682611473565b90565b67ffffffffffffffff1690565b610da581611466565b60208101610d968284611480565b6001600160a01b03918216919081169082820390811115610d9657610d9661137a565b60118152602081017f4e6f74206f776e6572206f722073656c6600000000000000000000000000000081529050610dfa565b60208082528101610d96816114ba565b80820180821115610d9657610d9661137a565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f2100000000000000000000000000000000000000000000000000000000000000602082015290506112bc565b60208082528101610d968161150f565b5f6020828403121561158a5761158a5f5ffd5b6103fe8383610f97565b5f8135610d9681610e85565b5f6001600160a01b03835b81169019929092169190911792915050565b5f6001600160a01b038216610d96565b5f610d96826115bd565b5f610d96826115cd565b6115ea826115d7565b6115f58183546115a0565b8255505050565b5f8135610d96816110d1565b5f7bffffffffffffffff00000000000000000000000000000000000000006115ab8460a01b90565b5f610d9667ffffffffffffffff8316611473565b61164d82611630565b6115f5818354611608565b5f8135610d9681610f8b565b5f7fffffffff000000000000000000000000000000000000000000000000000000006115ab8460e01b90565b5f63ffffffff8216610d96565b6116a682611690565b6115f5818354611664565b5f63ffffffff836115ab565b6116c682611690565b6115f58183546116b1565b5f808335601e19368590030181126116ea576116ea5f5ffd5b8301915050803567ffffffffffffffff811115611708576117085f5ffd5b602082019150600181023603821315610fe957610fe95f5ffd5b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061175e57607f821691505b602082108103610e3157610e31611736565b5f610d966114708381565b61178483611770565b81545f1960089490940293841b1916921b91909117905550565b5f61099181848461177b565b81811015610615576117bc5f8261179e565b6001016117aa565b601f821115610991575f818152602090206020601f850104810160208510156117ea5750805b6117fc6020601f8601048301826117aa565b5050505050565b8267ffffffffffffffff81111561181c5761181c611722565b611826825461174a565b6118318282856117c4565b505f601f821160018114611863575f831561184c5750848201355b5f19600885021c19811660028502178555506118ba565b5f84815260208120601f198516915b828110156118925787850135825560209485019460019092019101611872565b50848210156118ae575f196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b610991838383611803565b5f8135610d9681610ff0565b5f60ff8216610d96565b6118ec826118d9565b815460ff191660ff8216176115f5565b80828061190881611594565b905061191481846115e1565b50506020830180611924826115fc565b90506119308184611644565b5050604083018061194082611658565b905061194c818461169d565b505050606082018061195d82611658565b905061196c81600185016116bd565b505061197b60808301836116d1565b6119898183600286016118c2565b505060a0820180611999826118cd565b90506119a881600385016118e3565b50505050565b61061582826118fc565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f7373616765000000000000000000000000000000000000000000000000000000602082015290506112bc565b60208082528101610d96816119b8565b60c08101611a2e828a610d9c565b611a3b602083018961107c565b611a48604083018861110a565b611a55606083018761110a565b8181036080830152611a6881858761116e565b9050611a7760a083018461119d565b98975050505050505050565b67ffffffffffffffff918216919081169082820190811115610d9657610d9661137a56fea26469706673582212202bd7fd3401fcdb5c821c2216f05c222109591814004f071ca37762f35a1a4faa64736f6c634300081c0033",
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

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address caller, address feesAddress) returns()
func (_MessageBus *MessageBusTransactor) Initialize(opts *bind.TransactOpts, caller common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "initialize", caller, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address caller, address feesAddress) returns()
func (_MessageBus *MessageBusSession) Initialize(caller common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.Initialize(&_MessageBus.TransactOpts, caller, feesAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address caller, address feesAddress) returns()
func (_MessageBus *MessageBusTransactorSession) Initialize(caller common.Address, feesAddress common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.Initialize(&_MessageBus.TransactOpts, caller, feesAddress)
}

// NotifyDeposit is a paid mutator transaction binding the contract method 0xab53bddc.
//
// Solidity: function notifyDeposit(address receiver, uint256 amount) returns()
func (_MessageBus *MessageBusTransactor) NotifyDeposit(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "notifyDeposit", receiver, amount)
}

// NotifyDeposit is a paid mutator transaction binding the contract method 0xab53bddc.
//
// Solidity: function notifyDeposit(address receiver, uint256 amount) returns()
func (_MessageBus *MessageBusSession) NotifyDeposit(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.NotifyDeposit(&_MessageBus.TransactOpts, receiver, amount)
}

// NotifyDeposit is a paid mutator transaction binding the contract method 0xab53bddc.
//
// Solidity: function notifyDeposit(address receiver, uint256 amount) returns()
func (_MessageBus *MessageBusTransactorSession) NotifyDeposit(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.NotifyDeposit(&_MessageBus.TransactOpts, receiver, amount)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) payable returns(uint64 sequence)
func (_MessageBus *MessageBusTransactorSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MessageBus.Contract.PublishMessage(&_MessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MessageBus *MessageBusTransactor) ReceiveValueFromL2(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "receiveValueFromL2", receiver, amount)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MessageBus *MessageBusSession) ReceiveValueFromL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.ReceiveValueFromL2(&_MessageBus.TransactOpts, receiver, amount)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MessageBus *MessageBusTransactorSession) ReceiveValueFromL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.ReceiveValueFromL2(&_MessageBus.TransactOpts, receiver, amount)
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

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MessageBus *MessageBusTransactor) RetrieveAllFunds(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "retrieveAllFunds", receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MessageBus *MessageBusSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RetrieveAllFunds(&_MessageBus.TransactOpts, receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MessageBus *MessageBusTransactorSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _MessageBus.Contract.RetrieveAllFunds(&_MessageBus.TransactOpts, receiver)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MessageBus *MessageBusTransactor) SendValueToL2(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.contract.Transact(opts, "sendValueToL2", receiver, amount)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MessageBus *MessageBusSession) SendValueToL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.SendValueToL2(&_MessageBus.TransactOpts, receiver, amount)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MessageBus *MessageBusTransactorSession) SendValueToL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MessageBus.Contract.SendValueToL2(&_MessageBus.TransactOpts, receiver, amount)
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

// MessageBusNativeDepositIterator is returned from FilterNativeDeposit and is used to iterate over the raw logs and unpacked data for NativeDeposit events raised by the MessageBus contract.
type MessageBusNativeDepositIterator struct {
	Event *MessageBusNativeDeposit // Event containing the contract specifics and raw log

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
func (it *MessageBusNativeDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusNativeDeposit)
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
		it.Event = new(MessageBusNativeDeposit)
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
func (it *MessageBusNativeDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusNativeDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusNativeDeposit represents a NativeDeposit event raised by the MessageBus contract.
type MessageBusNativeDeposit struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNativeDeposit is a free log retrieval operation binding the contract event 0xcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d3.
//
// Solidity: event NativeDeposit(address indexed receiver, uint256 amount)
func (_MessageBus *MessageBusFilterer) FilterNativeDeposit(opts *bind.FilterOpts, receiver []common.Address) (*MessageBusNativeDepositIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "NativeDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusNativeDepositIterator{contract: _MessageBus.contract, event: "NativeDeposit", logs: logs, sub: sub}, nil
}

// WatchNativeDeposit is a free log subscription operation binding the contract event 0xcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d3.
//
// Solidity: event NativeDeposit(address indexed receiver, uint256 amount)
func (_MessageBus *MessageBusFilterer) WatchNativeDeposit(opts *bind.WatchOpts, sink chan<- *MessageBusNativeDeposit, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "NativeDeposit", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusNativeDeposit)
				if err := _MessageBus.contract.UnpackLog(event, "NativeDeposit", log); err != nil {
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

// ParseNativeDeposit is a log parse operation binding the contract event 0xcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d3.
//
// Solidity: event NativeDeposit(address indexed receiver, uint256 amount)
func (_MessageBus *MessageBusFilterer) ParseNativeDeposit(log types.Log) (*MessageBusNativeDeposit, error) {
	event := new(MessageBusNativeDeposit)
	if err := _MessageBus.contract.UnpackLog(event, "NativeDeposit", log); err != nil {
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

// MessageBusValueTransferIterator is returned from FilterValueTransfer and is used to iterate over the raw logs and unpacked data for ValueTransfer events raised by the MessageBus contract.
type MessageBusValueTransferIterator struct {
	Event *MessageBusValueTransfer // Event containing the contract specifics and raw log

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
func (it *MessageBusValueTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MessageBusValueTransfer)
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
		it.Event = new(MessageBusValueTransfer)
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
func (it *MessageBusValueTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MessageBusValueTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MessageBusValueTransfer represents a ValueTransfer event raised by the MessageBus contract.
type MessageBusValueTransfer struct {
	Sender   common.Address
	Receiver common.Address
	Amount   *big.Int
	Sequence uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValueTransfer is a free log retrieval operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MessageBus *MessageBusFilterer) FilterValueTransfer(opts *bind.FilterOpts, sender []common.Address, receiver []common.Address) (*MessageBusValueTransferIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MessageBus.contract.FilterLogs(opts, "ValueTransfer", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &MessageBusValueTransferIterator{contract: _MessageBus.contract, event: "ValueTransfer", logs: logs, sub: sub}, nil
}

// WatchValueTransfer is a free log subscription operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MessageBus *MessageBusFilterer) WatchValueTransfer(opts *bind.WatchOpts, sink chan<- *MessageBusValueTransfer, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MessageBus.contract.WatchLogs(opts, "ValueTransfer", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MessageBusValueTransfer)
				if err := _MessageBus.contract.UnpackLog(event, "ValueTransfer", log); err != nil {
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

// ParseValueTransfer is a log parse operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MessageBus *MessageBusFilterer) ParseValueTransfer(log types.Log) (*MessageBusValueTransfer, error) {
	event := new(MessageBusValueTransfer)
	if err := _MessageBus.contract.UnpackLog(event, "ValueTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
