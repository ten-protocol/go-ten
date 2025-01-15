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
	Bin: "0x608060405234801561001057600080fd5b5061001a33610027565b610022610098565b61014a565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100e85760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146101475780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b611b60806101596000396000f3fe6080604052600436106100d55760003560e01c8063715018a61161007f57806399a3ad211161005957806399a3ad21146102af578063ab53bddc146102cf578063b1454caa146102ef578063f2fde38b1461030f57610149565b8063715018a6146102355780638da5cb5b1461024a5780639730886d1461028f57610149565b8063346633fb116100b0578063346633fb146101e257806336d2da90146101f5578063485cc9551461021557610149565b8062a1b8151461016a5780630fcfbd111461019557806333a88c72146101b557610149565b36610149576040517f346633fb000000000000000000000000000000000000000000000000000000008152309063346633fb90349061011a9033908390600401610dea565b6000604051808303818588803b15801561013357600080fd5b505af1158015610147573d6000803e3d6000fd5b005b60405162461bcd60e51b815260040161016190610e39565b60405180910390fd5b34801561017657600080fd5b5061017f61032f565b60405161018c9190610e49565b60405180910390f35b3480156101a157600080fd5b5061017f6101b0366004610e72565b6103bb565b3480156101c157600080fd5b506101d56101d0366004610e72565b61041a565b60405161018c9190610eb5565b6101476101f0366004610ef3565b61046c565b34801561020157600080fd5b50610147610210366004610f2b565b6105bb565b34801561022157600080fd5b50610147610230366004610f4a565b61063a565b34801561024157600080fd5b506101476107a5565b34801561025657600080fd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031660405161018c9190610f79565b34801561029b57600080fd5b506101476102aa366004610f87565b6107b9565b3480156102bb57600080fd5b506101476102ca366004610ef3565b610925565b3480156102db57600080fd5b506101476102ea366004610ef3565b6109c5565b6103026102fd366004611050565b610a8e565b60405161018c91906110dd565b34801561031b57600080fd5b5061014761032a366004610f2b565b610b9b565b600354604080517f1a90a21900000000000000000000000000000000000000000000000000000000815290516000926001600160a01b031691631a90a2199160048083019260209291908290030181865afa158015610392573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103b691906110f6565b905090565b600080826040516020016103cf91906112b4565b60408051601f198184030181529181528151602092830120600081815292839052912054909150806104135760405162461bcd60e51b81526004016101619061131f565b9392505050565b6000808260405160200161042e91906112b4565b60408051601f19818403018152918152815160209283012060008181529283905291205490915080158015906104645750428111155b949350505050565b60003411801561047b57508034145b6104975760405162461bcd60e51b815260040161016190611387565b60035434906001600160a01b03161561055b5760006104b461032f565b9050803410156104d65760405162461bcd60e51b8152600401610161906113c7565b6104e081346113ed565b6003546040519193506000916001600160a01b039091169083908381818185875af1925050503d8060008114610532576040519150601f19603f3d011682016040523d82523d6000602084013e610537565b606091505b50509050806105585760405162461bcd60e51b815260040161016190611458565b50505b600061056633610bf2565b9050836001600160a01b0316336001600160a01b03167f50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b684846040516105ad929190611468565b60405180910390a350505050565b6105c3610c50565b6000816001600160a01b03164760405160006040518083038185875af1925050503d8060008114610610576040519150601f19603f3d011682016040523d82523d6000602084013e610615565b606091505b50509050806106365760405162461bcd60e51b8152600401610161906114b5565b5050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff166000811580156106855750825b905060008267ffffffffffffffff1660011480156106a25750303b155b9050811580156106b0575080155b156106e7576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561071b57845468ff00000000000000001916680100000000000000001785555b61072487610cc4565b6003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038816179055831561079c57845468ff0000000000000000191685556040517fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290610793906001906114e9565b60405180910390a15b50505050505050565b6107ad610c50565b6107b76000610cd5565b565b60006107c66001306114f7565b90506107f97f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b031614806108205750336001600160a01b038216145b61083c5760405162461bcd60e51b81526004016101619061154c565b6000610848834261155c565b905060008460405160200161085d91906112b4565b60408051601f198184030181529181528151602092830120600081815292839052912054909150156108a15760405162461bcd60e51b8152600401610161906115c7565b6000818152602081815260408220849055600191906108c290880188610f2b565b6001600160a01b0316815260208101919091526040016000908120906108ee60808801606089016115d7565b63ffffffff1681526020808201929092526040016000908120805460018101825590825291902086916004020161079c8282611a31565b61092d610c50565b8047101561094d5760405162461bcd60e51b8152600401610161906113c7565b6000826001600160a01b03168260405160006040518083038185875af1925050503d806000811461099a576040519150601f19603f3d011682016040523d82523d6000602084013e61099f565b606091505b50509050806109c05760405162461bcd60e51b8152600401610161906114b5565b505050565b60006109d26001306114f7565b9050610a057f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316336001600160a01b03161480610a2c5750336001600160a01b038216145b610a485760405162461bcd60e51b81526004016101619061154c565b826001600160a01b03167fcd9850463422a7449c406a036e35e5edb6fbe35a64c9f12a2354be98a750c0d383604051610a819190610e49565b60405180910390a2505050565b6003546000906001600160a01b031615610b44576000610aac61032f565b905080341015610ace5760405162461bcd60e51b815260040161016190611a93565b6003546040516000916001600160a01b03169083908381818185875af1925050503d8060008114610b1b576040519150601f19603f3d011682016040523d82523d6000602084013e610b20565b606091505b5050905080610b415760405162461bcd60e51b815260040161016190611458565b50505b610b4d33610bf2565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef77593733828888888888604051610b8a9796959493929190611aa3565b60405180910390a195945050505050565b610ba3610c50565b6001600160a01b038116610be65760006040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016101619190610f79565b610bef81610cd5565b50565b6001600160a01b0381166000908152600260205260408120805467ffffffffffffffff169160019190610c258385611b06565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b33610c827f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146107b757336040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016101619190610f79565b610ccc610d53565b610bef81610dba565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff166107b7576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ba3610d53565b60006001600160a01b0382165b92915050565b610dde81610dc2565b82525050565b80610dde565b60408101610df88285610dd5565b6104136020830184610de4565b600b8152602081017f756e737570706f72746564000000000000000000000000000000000000000000815290505b60200190565b60208082528101610dcf81610e05565b60208101610dcf8284610de4565b600060c08284031215610e6c57610e6c600080fd5b50919050565b600060208284031215610e8757610e87600080fd5b813567ffffffffffffffff811115610ea157610ea1600080fd5b61046484828501610e57565b801515610dde565b60208101610dcf8284610ead565b610ecc81610dc2565b8114610bef57600080fd5b8035610dcf81610ec3565b80610ecc565b8035610dcf81610ee2565b60008060408385031215610f0957610f09600080fd5b610f138484610ed7565b9150610f228460208501610ee8565b90509250929050565b600060208284031215610f4057610f40600080fd5b6104138383610ed7565b60008060408385031215610f6057610f60600080fd5b610f6a8484610ed7565b9150610f228460208501610ed7565b60208101610dcf8284610dd5565b60008060408385031215610f9d57610f9d600080fd5b823567ffffffffffffffff811115610fb757610fb7600080fd5b610fc385828601610e57565b925050610f228460208501610ee8565b63ffffffff8116610ecc565b8035610dcf81610fd3565b60008083601f840112610fff57610fff600080fd5b50813567ffffffffffffffff81111561101a5761101a600080fd5b60208301915083600182028301111561103557611035600080fd5b9250929050565b60ff8116610ecc565b8035610dcf8161103c565b60008060008060006080868803121561106b5761106b600080fd5b6110758787610fdf565b94506110848760208801610fdf565b9350604086013567ffffffffffffffff8111156110a3576110a3600080fd5b6110af88828901610fea565b93509350506110c18760608801611045565b90509295509295909350565b67ffffffffffffffff8116610dde565b60208101610dcf82846110cd565b8051610dcf81610ee2565b60006020828403121561110b5761110b600080fd5b61041383836110eb565b506000610dcf6020830183610ed7565b67ffffffffffffffff8116610ecc565b8035610dcf81611125565b506000610dcf6020830183611135565b506000610dcf6020830183610fdf565b63ffffffff8116610dde565b6000808335601e193685900301811261118757611187600080fd5b830160208101925035905067ffffffffffffffff8111156111aa576111aa600080fd5b3681900382131561103557611035600080fd5b82818337506000910152565b8183526020830192506111dd8284836111bd565b50601f01601f19160190565b506000610dcf6020830183611045565b60ff8116610dde565b600060c083016112128380611115565b61121c8582610dd5565b5061122a6020840184611140565b61123760208601826110cd565b506112456040840184611150565b6112526040860182611160565b506112606060840184611150565b61126d6060860182611160565b5061127b608084018461116c565b858303608087015261128e8382846111c9565b9250505061129f60a08401846111e9565b6112ac60a08601826111f9565b509392505050565b602080825281016104138184611202565b60218152602081017f54686973206d65737361676520776173206e65766572207375626d697474656481527f2e00000000000000000000000000000000000000000000000000000000000000602082015290505b60400190565b60208082528101610dcf816112c5565b60308152602081017f417474656d7074696e6720746f2073656e642076616c756520776974686f757481527f2070726f766964696e672045746865720000000000000000000000000000000060208201529050611319565b60208082528101610dcf8161132f565b60208082527f496e73756666696369656e742066756e647320746f2073656e642076616c75659101908152610e33565b60208082528101610dcf81611397565b634e487b7160e01b600052601160045260246000fd5b81810381811115610dcf57610dcf6113d7565b60248152602081017f4661696c656420746f2073656e64206665657320746f206665657320636f6e7481527f726163740000000000000000000000000000000000000000000000000000000060208201529050611319565b60208082528101610dcf81611400565b604081016114768285610de4565b61041360208301846110cd565b60148152602081017f6661696c65642073656e64696e672076616c756500000000000000000000000081529050610e33565b60208082528101610dcf81611483565b6000610dcf826114d3565b90565b67ffffffffffffffff1690565b610dde816114c5565b60208101610dcf82846114e0565b6001600160a01b03918216919081169082820390811115610dcf57610dcf6113d7565b60118152602081017f4e6f74206f776e6572206f722073656c6600000000000000000000000000000081529050610e33565b60208082528101610dcf8161151a565b80820180821115610dcf57610dcf6113d7565b60218152602081017f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636581527f210000000000000000000000000000000000000000000000000000000000000060208201529050611319565b60208082528101610dcf8161156f565b6000602082840312156115ec576115ec600080fd5b6104138383610fdf565b60008135610dcf81610ec3565b60006001600160a01b03835b81169019929092169190911792915050565b60006001600160a01b038216610dcf565b6000610dcf82611621565b6000610dcf82611632565b6116518261163d565b61165c818354611603565b8255505050565b60008135610dcf81611125565b60007bffffffffffffffff000000000000000000000000000000000000000061160f8460a01b90565b6000610dcf67ffffffffffffffff83166114d3565b6116b782611699565b61165c818354611670565b60008135610dcf81610fd3565b60007fffffffff0000000000000000000000000000000000000000000000000000000061160f8460e01b90565b600063ffffffff8216610dcf565b611713826116fc565b61165c8183546116cf565b600063ffffffff8361160f565b611734826116fc565b61165c81835461171e565b6000808335601e193685900301811261175a5761175a600080fd5b8301915050803567ffffffffffffffff81111561177957611779600080fd5b60208201915060018102360382131561103557611035600080fd5b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052602260045260246000fd5b6002810460018216806117d457607f821691505b602082108103610e6c57610e6c6117aa565b6000610dcf6114d08381565b6117fb836117e6565b815460001960089490940293841b1916921b91909117905550565b60006109c08184846117f2565b8181101561063657611836600082611816565b600101611823565b601f8211156109c0576000818152602090206020601f850104810160208510156118655750805b6118776020601f860104830182611823565b5050505050565b8267ffffffffffffffff81111561189757611897611794565b6118a182546117c0565b6118ac82828561183e565b506000601f8211600181146118e157600083156118c95750848201355b600019600885021c198116600285021785555061193b565b600084815260209020601f19841690835b8281101561191257878501358255602094850194600190920191016118f2565b508482101561192f576000196008601f8716021c19878501351681555b50506001600284020184555b505050505050565b6109c083838361187e565b60008135610dcf8161103c565b600060ff8216610dcf565b61196f8261195b565b815460ff191660ff82161761165c565b80828061198b816115f6565b90506119978184611648565b505060208301806119a782611663565b90506119b381846116ae565b505060408301806119c3826116c2565b90506119cf818461170a565b50505060608201806119e0826116c2565b90506119ef816001850161172b565b50506119fe608083018361173f565b611a0c818360028601611943565b505060a0820180611a1c8261194e565b9050611a2b8160038501611966565b50505050565b610636828261197f565b60258152602081017f496e73756666696369656e742066756e647320746f207075626c697368206d6581527f737361676500000000000000000000000000000000000000000000000000000060208201529050611319565b60208082528101610dcf81611a3b565b60c08101611ab1828a610dd5565b611abe60208301896110cd565b611acb6040830188611160565b611ad86060830187611160565b8181036080830152611aeb8185876111c9565b9050611afa60a08301846111f9565b98975050505050505050565b67ffffffffffffffff918216919081169082820190811115610dcf57610dcf6113d756fea2646970667358221220f95f53578abddab4dbfcba4e38ed51470d4d4cbc31b23b532e46b0005f46de4b64736f6c634300081c0033",
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
