// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MerkleTreeMessageBus

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

// StructsValueTransferMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsValueTransferMessage struct {
	Sender   common.Address
	Receiver common.Address
	Amount   *big.Int
	Sequence uint64
}

// MerkleTreeMessageBusMetaData contains all meta data concerning the MerkleTreeMessageBus contract.
var MerkleTreeMessageBusMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"LogMessagePublished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"name\":\"ValueTransfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"activationTime\",\"type\":\"uint256\"}],\"name\":\"addStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"name\":\"disableStateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"getMessageTimeOfFinality\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"name\":\"publishMessage\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"receiveValueFromL2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"retrieveAllFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sendValueToL2\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"finalAfterTimestamp\",\"type\":\"uint256\"}],\"name\":\"storeCrossChainMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"crossChainMessage\",\"type\":\"tuple\"}],\"name\":\"verifyMessageFinalized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyMessageInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"internalType\":\"structStructs.ValueTransferMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyValueTransferInclusion\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052346200002657620000146200002b565b604051611ae5620001c58239611ae590f35b600080fd5b6200003562000037565b565b62000035336200008d565b6200005562000052620000529290565b90565b6001600160a01b031690565b620000529062000042565b620000779062000055565b9052565b6020810192916200003591906200006c565b62000099600062000061565b620000a48162000055565b620000af8362000055565b14620000c15750620000359062000162565b620000e790620000d060405190565b631e4fbdf760e01b8152918291600483016200007b565b0390fd5b620000529062000055565b620000529054620000eb565b620000529062000055906001600160a01b031682565b620000529062000102565b620000529062000118565b9062000142620000526200015e9262000123565b82546001600160a01b0319166001600160a01b03919091161790565b9055565b6200016e6000620000f6565b906200017c8160006200012e565b620001b3620001ac7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09362000123565b9162000123565b91620001be60405190565b600090a356fe6080604052600436101561001d575b366119cf5761001b611a3d565b005b60003560e01c80630fcfbd11146100fd5780630fe9188e146100f857806333a88c72146100f3578063346633fb146100ee57806336d2da90146100e9578063715018a6146100e45780638da5cb5b146100df5780639730886d146100da57806399a3ad21146100d5578063b1454caa146100d0578063b201246f146100cb578063b6aed0cb146100c6578063e138a8d2146100c15763f2fde38b0361000e57610560565b610544565b6104d4565b61049e565b6103e2565b610316565b6102fd565b610291565b610279565b610256565b61022e565b6101b6565b610199565b610140565b908160c09103126101105790565b600080fd5b9060208282031261011057813567ffffffffffffffff81116101105761013b9201610102565b90565b565b346101105761016b61015b610156366004610115565b61142d565b6040519182918290815260200190565b0390f35b805b0361011057565b9050359061013e8261016f565b906020828203126101105761013b91610178565b34610110576101b16101ac366004610185565b61091b565b604051005b346101105761016b6101d16101cc366004610115565b61138a565b60405191829182901515815260200190565b6001600160a01b031690565b6001600160a01b038116610171565b9050359061013e826101ef565b91906040838203126101105761013b9061022581856101fe565b93602001610178565b6101b161023c36600461020b565b90611084565b906020828203126101105761013b916101fe565b34610110576101b1610269366004610242565b6119c6565b600091031261011057565b346101105761028936600461026e565b6101b16105b0565b34610110576102a136600461026e565b61016b6102b66000546001600160a01b031690565b604051918291826001600160a01b03909116815260200190565b91906040838203126101105782359067ffffffffffffffff8211610110576102258161013b938601610102565b34610110576101b16103103660046102d0565b90611996565b34610110576101b161032936600461020b565b906112c0565b63ffffffff8116610171565b9050359061013e8261032f565b909182601f830112156101105781359167ffffffffffffffff831161011057602001926001830284011161011057565b60ff8116610171565b9050359061013e82610378565b919091608081840312610110576103a5838261033b565b926103b3816020840161033b565b92604083013567ffffffffffffffff8111610110576103d78361013b928601610348565b939094606001610381565b346101105761016b6104016103f836600461038e565b93929092611331565b6040519182918267ffffffffffffffff909116815260200190565b908160809103126101105790565b909182601f830112156101105781359167ffffffffffffffff831161011057602001926020830284011161011057565b9060c0828203126101105761046f818361041c565b92608083013567ffffffffffffffff8111610110576104938361013b92860161042a565b93909460a001610178565b34610110576101b16104b136600461045a565b92919091610f8a565b91906040838203126101105761013b906102258185610178565b34610110576101b16104e73660046104ba565b90610876565b9060608282031261011057813567ffffffffffffffff81116101105781610515918401610102565b92602083013567ffffffffffffffff8111610110576105398361013b92860161042a565b939094604001610178565b34610110576101b16105573660046104ed565b92919091610c62565b34610110576101b1610573366004610242565b6106a1565b6105806105b8565b61013e61059e565b6101e361013b61013b9290565b61013b90610588565b61013e6105ab6000610595565b610702565b61013e610578565b60005433906001600160a01b03168190036105d05750565b61061a906105dd60405190565b9182917f118cdaa7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b0390fd5b61013e9061062a6105b8565b6106346000610595565b6001600160a01b0381166001600160a01b03831614610657575061013e90610702565b61061a9061066460405190565b9182917f1e4fbdf7000000000000000000000000000000000000000000000000000000008352600483016001600160a01b03909116815260200190565b61013e9061061e565b906001600160a01b03905b9181191691161790565b6101e361013b61013b926001600160a01b031690565b61013b906106bf565b61013b906106d5565b906106f761013b6106fe926106de565b82546106aa565b9055565b6000546001600160a01b03169061071a8160006106e7565b61074d6107477f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0936106de565b916106de565b9161075760405190565b80805b0390a3565b9061013e9161076c6105b8565b610836565b905b600052602052604060002090565b61013b9081565b61013b9054610781565b61013b61013b61013b9290565b0190565b156107aa57565b60405162461bcd60e51b815260206004820152602560248201527f526f6f7420616c726561647920616464656420746f20746865206d657373616760448201527f65206275730000000000000000000000000000000000000000000000000000006064820152608490fd5b90600019906106b5565b9061082f61013b6106fe92610792565b8254610815565b9061087161013e92600461086c6108556108508484610771565b610788565b6108666108626000610792565b9190565b146107a3565b610771565b61081f565b9061013e9161075f565b61013e9061088c6105b8565b6108dd565b1561089857565b60405162461bcd60e51b815260206004820152601a60248201527f537461746520726f6f7420646f6573206e6f742065786973742e0000000000006044820152606490fd5b61013e906108716004916109156108f76108508386610771565b61091060009161090961086284610792565b1415610891565b610792565b92610771565b61013e90610880565b1561092b57565b60405162461bcd60e51b815260206004820152602a60248201527f526f6f74206973206e6f74207075626c6973686564206f6e2074686973206d6560448201527f7373616765206275732e000000000000000000000000000000000000000000006064820152608490fd5b1561099d57565b60405162461bcd60e51b815260206004820152602160248201527f526f6f74206973206e6f7420636f6e736964657265642066696e616c207965746044820152601760f91b6064820152608490fd5b5061013b9060208101906101fe565b67ffffffffffffffff8116610171565b9050359061013e826109fb565b5061013b906020810190610a0b565b5061013b90602081019061033b565b9035601e19368390030181121561011057016020813591019167ffffffffffffffff82116101105736829003831361011057565b90826000939282370152565b9190610a9481610a8d8161079f9560209181520190565b8095610a6a565b601f01601f191690565b5061013b906020810190610381565b9061013b9060a0610b58610b4e60c08401610ad8610acb88806109ec565b6001600160a01b03168652565b610af9610ae86020890189610a18565b67ffffffffffffffff166020870152565b610b16610b096040890189610a27565b63ffffffff166040870152565b610b33610b266060890189610a27565b63ffffffff166060870152565b610b406080880188610a36565b908683036080880152610a76565b9482810190610a9e565b60ff16910152565b602080825261013b92910190610aad565b634e487b7160e01b600052604160045260246000fd5b90601f01601f1916810190811067ffffffffffffffff821117610ba957604052565b610b71565b60408082526001908201527f6d000000000000000000000000000000000000000000000000000000000000006060820152919061013e90602060808501940152565b15610bf757565b60405162461bcd60e51b815260206004820152603060248201527f496e76616c696420696e636c7573696f6e2070726f6f6620666f722063726f7360448201527f7320636861696e206d6573736167652e000000000000000000000000000000006064820152608490fd5b92610d469261013e94610caf6004610c94610c806108508784610771565b610c8d6108626000610792565b1415610924565b610ca861086261013b610850884295610771565b1015610996565b610cd7610cbb60405190565b8092610ccb602083019182610b60565b90810382520382610b87565b610ce9610ce2825190565b9160200190565b20610d06610cf660405190565b8092610ccb602083019182610bae565b610d11610ce2825190565b20610d34610d1e60405190565b8092610ccb60208301918261079f816020939052565b610d3f610ce2825190565b2092610d4b565b610bf0565b610d6a929361086292610d6692610d60600090565b50610dcd565b9290565b1490565b634e487b7160e01b600052601160045260246000fd5b6000198114610d935760010190565b610d6e565b634e487b7160e01b600052603260045260246000fd5b9190811015610dbe576020020190565b610d98565b3561013b8161016f565b610dd76000610792565b925b82841015610e0d57610e01610e0791610dfb610df6878787610dae565b610dc3565b90610e14565b93610d84565b92610dd9565b9250505090565b81811015610e2f579061013b91600052602052604060002090565b61013b91600052602052604060002090565b5061013b906020810190610178565b906060610eb761013e93610e77610e6a60008301836109ec565b6001600160a01b03168552565b610e97610e8760208301836109ec565b6001600160a01b03166020860152565b610eae610ea76040830183610e41565b6040860152565b82810190610a18565b67ffffffffffffffff16910152565b60808101929161013e9190610e50565b60408082526001908201527f76000000000000000000000000000000000000000000000000000000000000006060820152919061013e90602060808501940152565b15610f1f57565b60405162461bcd60e51b815260206004820152603360248201527f496e76616c696420696e636c7573696f6e2070726f6f6620666f722076616c7560448201527f65207472616e73666572206d6573736167652e000000000000000000000000006064820152608490fd5b92610fec9261013e94610fa86004610c94610c806108508784610771565b610fc4610fb460405190565b8092610ccb602083019182610ec6565b610fcf610ce2825190565b20610d06610fdc60405190565b8092610ccb602083019182610ed6565b610f18565b15610ff857565b60405162461bcd60e51b815260206004820152603060248201527f417474656d7074696e6720746f2073656e642076616c756520776974686f757460448201527f2070726f766964696e67204574686572000000000000000000000000000000006064820152608490fd5b90815260408101929161013e916020905b019067ffffffffffffffff169052565b6110a234926110936000610792565b841190816110f5575b50610ff1565b336110ac81611191565b6110df6110d97f50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6936106de565b936106de565b9361075a6110ec60405190565b92839283611063565b905083143861109c565b90610773906106de565b61013b905b67ffffffffffffffff1690565b61013b9054611109565b61110e61013b61013b9290565b67ffffffffffffffff9081169116019067ffffffffffffffff8211610d9357565b9067ffffffffffffffff906106b5565b61110e61013b61013b9267ffffffffffffffff1690565b9061118a61013b6106fe92611163565b8254611153565b9061013e60036111cf6111c06111af6111aa87856110ff565b61111b565b956111ba6001611125565b936110ff565b916111ca8361111b565b611132565b9061117a565b9061013e916111e26105b8565b61129a565b9061013e6111f460405190565b9283610b87565b67ffffffffffffffff8111610ba957602090601f01601f19160190565b9061122a611225836111fb565b6111e7565b918252565b3d156112495761123e3d611218565b903d6000602084013e565b606090565b1561125557565b60405162461bcd60e51b815260206004820152601460248201527f6661696c65642073656e64696e672076616c75650000000000000000000000006044820152606490fd5b6000809161013e936112ab60405190565b90818003925af16112ba61122f565b5061124e565b9061013e916111d5565b97969461131a60a09561130d6113279661130061013e9b969760208f61107460c082019b60008301906001600160a01b03169052565b63ffffffff1660408d0152565b63ffffffff1660608b0152565b88830360808a0152610a76565b94019060ff169052565b9290919361138490611341600090565b50339261134d84611191565b9687957fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef7759379761137b60405190565b978897886112ca565b0390a190565b6108506113b991611399600090565b506113a6610cbb60405190565b6113b1610ce2825190565b206001610771565b6113c36000610792565b811190816113cf575090565b429150111590565b156113de57565b60405162461bcd60e51b815260206004820152602160248201527f54686973206d65737361676520776173206e65766572207375626d69747465646044820152601760f91b6064820152608490fd5b61085061143c91611399600090565b61013b6114496000610792565b82116113d7565b9061013e9161145d6105b8565b6118ff565b91908201809211610d9357565b1561147657565b60405162461bcd60e51b815260206004820152602160248201527f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636560448201527f21000000000000000000000000000000000000000000000000000000000000006064820152608490fd5b3561013b816101ef565b3561013b8161032f565b61150861013b61013b9263ffffffff1690565b63ffffffff1690565b90610773906114f5565b8054821015610dbe57611535600491600052602060002090565b91020190600090565b634e487b7160e01b600052600060045260246000fd5b3561013b816109fb565b907bffffffffffffffff00000000000000000000000000000000000000009060a01b6106b5565b9061159561013b6106fe92611163565b825461155e565b907fffffffff000000000000000000000000000000000000000000000000000000009060e01b6106b5565b906115d761013b6106fe926114f5565b825461159c565b9063ffffffff906106b5565b906115fa61013b6106fe926114f5565b82546115de565b903590601e193682900301821215610110570180359067ffffffffffffffff8211610110576020019136829003831361011057565b634e487b7160e01b600052602260045260246000fd5b906001600283049216801561166c575b602083101461166757565b611636565b91607f169161165c565b9160001960089290920291821b911b6106b5565b919061169b61013b6106fe93610792565b908354611676565b61013e9160009161168a565b8181106116ba575050565b806116c860006001936116a3565b016116af565b9190601f81116116dd57505050565b6116ef61013e93600052602060002090565b906020601f840181900483019310611711575b6020601f9091010401906116af565b9091508190611702565b919067ffffffffffffffff8211610ba9576117408261173a855461164c565b856116ce565b600090601f831160011461177b576106fe929160009183611770575b5050600019600883021c1916906002021790565b01359050388061175c565b601f1983169161179085600052602060002090565b92815b8181106117cd575091600293918560019694106117b4575b50505002019055565b0135600019601f84166008021c191690553880806117ab565b91936020600181928787013581550195019201611793565b9061013e929161171b565b3561013b81610378565b9060ff906106b5565b61181361013b61013b9260ff1690565b60ff1690565b9061182961013b6106fe92611803565b82546117fa565b61013e91906118b79060a09060039061187c8161185561184f886114e1565b826106e7565b61186a61186460208901611554565b82611585565b611876604088016114eb565b906115c7565b6118956001820161188f606088016114eb565b906115ea565b6118af600282016118a96080880188611601565b916117e5565b0192016117f0565b90611819565b91906118cc5761013e91611830565b61153e565b9081549168010000000000000000831015610ba957826118f991600161013e9501815561151b565b906118bd565b61199161013b61197f61191561013e9542611462565b61196c61192160405190565b6119338160208101610ccb8a82610b60565b61193e610ce2825190565b20916108716001936119536108508287610771565b9461086c60009661196661086289610792565b1461146f565b61197960029186016114e1565b906110ff565b61198b606085016114eb565b90611511565b6118d1565b9061013e91611450565b61013e906119ac6105b8565b60008061013e926119bc306106de565b316112ab60405190565b61013e906119a0565b60405162461bcd60e51b815260206004820152600b60248201527f756e737570706f727465640000000000000000000000000000000000000000006044820152606490fd5b6001600160a01b03909116815260408101929161013e9160200152565b6040513d6000823e3d90fd5b611a46306106de565b63346633fb34913390803b1561011057611a7684600094611a81611a6960405190565b9788968795869460e01b90565b845260048401611a14565b03925af18015611aaa57611a925750565b61013e906000611aa28183610b87565b81019061026e565b611a3156fea2646970667358221220f1d785fdbb2d651555a3e150a57286c72acedb6ae4e42ca19d2777f4fbfc83b664736f6c63430008140033",
}

// MerkleTreeMessageBusABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleTreeMessageBusMetaData.ABI instead.
var MerkleTreeMessageBusABI = MerkleTreeMessageBusMetaData.ABI

// MerkleTreeMessageBusBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MerkleTreeMessageBusMetaData.Bin instead.
var MerkleTreeMessageBusBin = MerkleTreeMessageBusMetaData.Bin

// DeployMerkleTreeMessageBus deploys a new Ethereum contract, binding an instance of MerkleTreeMessageBus to it.
func DeployMerkleTreeMessageBus(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MerkleTreeMessageBus, error) {
	parsed, err := MerkleTreeMessageBusMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MerkleTreeMessageBusBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MerkleTreeMessageBus{MerkleTreeMessageBusCaller: MerkleTreeMessageBusCaller{contract: contract}, MerkleTreeMessageBusTransactor: MerkleTreeMessageBusTransactor{contract: contract}, MerkleTreeMessageBusFilterer: MerkleTreeMessageBusFilterer{contract: contract}}, nil
}

// MerkleTreeMessageBus is an auto generated Go binding around an Ethereum contract.
type MerkleTreeMessageBus struct {
	MerkleTreeMessageBusCaller     // Read-only binding to the contract
	MerkleTreeMessageBusTransactor // Write-only binding to the contract
	MerkleTreeMessageBusFilterer   // Log filterer for contract events
}

// MerkleTreeMessageBusCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeMessageBusTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeMessageBusFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleTreeMessageBusFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeMessageBusSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleTreeMessageBusSession struct {
	Contract     *MerkleTreeMessageBus // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// MerkleTreeMessageBusCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleTreeMessageBusCallerSession struct {
	Contract *MerkleTreeMessageBusCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// MerkleTreeMessageBusTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleTreeMessageBusTransactorSession struct {
	Contract     *MerkleTreeMessageBusTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// MerkleTreeMessageBusRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleTreeMessageBusRaw struct {
	Contract *MerkleTreeMessageBus // Generic contract binding to access the raw methods on
}

// MerkleTreeMessageBusCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusCallerRaw struct {
	Contract *MerkleTreeMessageBusCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleTreeMessageBusTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleTreeMessageBusTransactorRaw struct {
	Contract *MerkleTreeMessageBusTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleTreeMessageBus creates a new instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBus(address common.Address, backend bind.ContractBackend) (*MerkleTreeMessageBus, error) {
	contract, err := bindMerkleTreeMessageBus(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBus{MerkleTreeMessageBusCaller: MerkleTreeMessageBusCaller{contract: contract}, MerkleTreeMessageBusTransactor: MerkleTreeMessageBusTransactor{contract: contract}, MerkleTreeMessageBusFilterer: MerkleTreeMessageBusFilterer{contract: contract}}, nil
}

// NewMerkleTreeMessageBusCaller creates a new read-only instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBusCaller(address common.Address, caller bind.ContractCaller) (*MerkleTreeMessageBusCaller, error) {
	contract, err := bindMerkleTreeMessageBus(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusCaller{contract: contract}, nil
}

// NewMerkleTreeMessageBusTransactor creates a new write-only instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBusTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleTreeMessageBusTransactor, error) {
	contract, err := bindMerkleTreeMessageBus(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusTransactor{contract: contract}, nil
}

// NewMerkleTreeMessageBusFilterer creates a new log filterer instance of MerkleTreeMessageBus, bound to a specific deployed contract.
func NewMerkleTreeMessageBusFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleTreeMessageBusFilterer, error) {
	contract, err := bindMerkleTreeMessageBus(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusFilterer{contract: contract}, nil
}

// bindMerkleTreeMessageBus binds a generic wrapper to an already deployed contract.
func bindMerkleTreeMessageBus(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MerkleTreeMessageBusMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleTreeMessageBus.Contract.MerkleTreeMessageBusCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.MerkleTreeMessageBusTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.MerkleTreeMessageBusTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleTreeMessageBus.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.contract.Transact(opts, method, params...)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) GetMessageTimeOfFinality(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "getMessageTimeOfFinality", crossChainMessage)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetMessageTimeOfFinality(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// GetMessageTimeOfFinality is a free data retrieval call binding the contract method 0x0fcfbd11.
//
// Solidity: function getMessageTimeOfFinality((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(uint256)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) GetMessageTimeOfFinality(crossChainMessage StructsCrossChainMessage) (*big.Int, error) {
	return _MerkleTreeMessageBus.Contract.GetMessageTimeOfFinality(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Owner() (common.Address, error) {
	return _MerkleTreeMessageBus.Contract.Owner(&_MerkleTreeMessageBus.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) Owner() (common.Address, error) {
	return _MerkleTreeMessageBus.Contract.Owner(&_MerkleTreeMessageBus.CallOpts)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyMessageFinalized(opts *bind.CallOpts, crossChainMessage StructsCrossChainMessage) (bool, error) {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyMessageFinalized", crossChainMessage)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MerkleTreeMessageBus.Contract.VerifyMessageFinalized(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageFinalized is a free data retrieval call binding the contract method 0x33a88c72.
//
// Solidity: function verifyMessageFinalized((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage) view returns(bool)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) VerifyMessageFinalized(crossChainMessage StructsCrossChainMessage) (bool, error) {
	return _MerkleTreeMessageBus.Contract.VerifyMessageFinalized(&_MerkleTreeMessageBus.CallOpts, crossChainMessage)
}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xe138a8d2.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint32,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyMessageInclusion(opts *bind.CallOpts, message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyMessageInclusion", message, proof, root)

	if err != nil {
		return err
	}

	return err

}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xe138a8d2.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint32,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyMessageInclusion(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyMessageInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// VerifyMessageInclusion is a free data retrieval call binding the contract method 0xe138a8d2.
//
// Solidity: function verifyMessageInclusion((address,uint64,uint32,uint32,bytes,uint8) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) VerifyMessageInclusion(message StructsCrossChainMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyMessageInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// VerifyValueTransferInclusion is a free data retrieval call binding the contract method 0xb201246f.
//
// Solidity: function verifyValueTransferInclusion((address,address,uint256,uint64) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCaller) VerifyValueTransferInclusion(opts *bind.CallOpts, message StructsValueTransferMessage, proof [][32]byte, root [32]byte) error {
	var out []interface{}
	err := _MerkleTreeMessageBus.contract.Call(opts, &out, "verifyValueTransferInclusion", message, proof, root)

	if err != nil {
		return err
	}

	return err

}

// VerifyValueTransferInclusion is a free data retrieval call binding the contract method 0xb201246f.
//
// Solidity: function verifyValueTransferInclusion((address,address,uint256,uint64) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) VerifyValueTransferInclusion(message StructsValueTransferMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyValueTransferInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// VerifyValueTransferInclusion is a free data retrieval call binding the contract method 0xb201246f.
//
// Solidity: function verifyValueTransferInclusion((address,address,uint256,uint64) message, bytes32[] proof, bytes32 root) view returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusCallerSession) VerifyValueTransferInclusion(message StructsValueTransferMessage, proof [][32]byte, root [32]byte) error {
	return _MerkleTreeMessageBus.Contract.VerifyValueTransferInclusion(&_MerkleTreeMessageBus.CallOpts, message, proof, root)
}

// AddStateRoot is a paid mutator transaction binding the contract method 0xb6aed0cb.
//
// Solidity: function addStateRoot(bytes32 stateRoot, uint256 activationTime) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) AddStateRoot(opts *bind.TransactOpts, stateRoot [32]byte, activationTime *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "addStateRoot", stateRoot, activationTime)
}

// AddStateRoot is a paid mutator transaction binding the contract method 0xb6aed0cb.
//
// Solidity: function addStateRoot(bytes32 stateRoot, uint256 activationTime) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) AddStateRoot(stateRoot [32]byte, activationTime *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot, activationTime)
}

// AddStateRoot is a paid mutator transaction binding the contract method 0xb6aed0cb.
//
// Solidity: function addStateRoot(bytes32 stateRoot, uint256 activationTime) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) AddStateRoot(stateRoot [32]byte, activationTime *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.AddStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot, activationTime)
}

// DisableStateRoot is a paid mutator transaction binding the contract method 0x0fe9188e.
//
// Solidity: function disableStateRoot(bytes32 stateRoot) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) DisableStateRoot(opts *bind.TransactOpts, stateRoot [32]byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "disableStateRoot", stateRoot)
}

// DisableStateRoot is a paid mutator transaction binding the contract method 0x0fe9188e.
//
// Solidity: function disableStateRoot(bytes32 stateRoot) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) DisableStateRoot(stateRoot [32]byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.DisableStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot)
}

// DisableStateRoot is a paid mutator transaction binding the contract method 0x0fe9188e.
//
// Solidity: function disableStateRoot(bytes32 stateRoot) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) DisableStateRoot(stateRoot [32]byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.DisableStateRoot(&_MerkleTreeMessageBus.TransactOpts, stateRoot)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) PublishMessage(opts *bind.TransactOpts, nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "publishMessage", nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.PublishMessage(&_MerkleTreeMessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// PublishMessage is a paid mutator transaction binding the contract method 0xb1454caa.
//
// Solidity: function publishMessage(uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel) returns(uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) PublishMessage(nonce uint32, topic uint32, payload []byte, consistencyLevel uint8) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.PublishMessage(&_MerkleTreeMessageBus.TransactOpts, nonce, topic, payload, consistencyLevel)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) ReceiveValueFromL2(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "receiveValueFromL2", receiver, amount)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) ReceiveValueFromL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.ReceiveValueFromL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// ReceiveValueFromL2 is a paid mutator transaction binding the contract method 0x99a3ad21.
//
// Solidity: function receiveValueFromL2(address receiver, uint256 amount) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) ReceiveValueFromL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.ReceiveValueFromL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RenounceOwnership() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RenounceOwnership(&_MerkleTreeMessageBus.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RenounceOwnership(&_MerkleTreeMessageBus.TransactOpts)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) RetrieveAllFunds(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "retrieveAllFunds", receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RetrieveAllFunds(&_MerkleTreeMessageBus.TransactOpts, receiver)
}

// RetrieveAllFunds is a paid mutator transaction binding the contract method 0x36d2da90.
//
// Solidity: function retrieveAllFunds(address receiver) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) RetrieveAllFunds(receiver common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.RetrieveAllFunds(&_MerkleTreeMessageBus.TransactOpts, receiver)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) SendValueToL2(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "sendValueToL2", receiver, amount)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) SendValueToL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.SendValueToL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// SendValueToL2 is a paid mutator transaction binding the contract method 0x346633fb.
//
// Solidity: function sendValueToL2(address receiver, uint256 amount) payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) SendValueToL2(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.SendValueToL2(&_MerkleTreeMessageBus.TransactOpts, receiver, amount)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) StoreCrossChainMessage(opts *bind.TransactOpts, crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "storeCrossChainMessage", crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.StoreCrossChainMessage(&_MerkleTreeMessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// StoreCrossChainMessage is a paid mutator transaction binding the contract method 0x9730886d.
//
// Solidity: function storeCrossChainMessage((address,uint64,uint32,uint32,bytes,uint8) crossChainMessage, uint256 finalAfterTimestamp) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) StoreCrossChainMessage(crossChainMessage StructsCrossChainMessage, finalAfterTimestamp *big.Int) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.StoreCrossChainMessage(&_MerkleTreeMessageBus.TransactOpts, crossChainMessage, finalAfterTimestamp)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.TransferOwnership(&_MerkleTreeMessageBus.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.TransferOwnership(&_MerkleTreeMessageBus.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Fallback(&_MerkleTreeMessageBus.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Fallback(&_MerkleTreeMessageBus.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeMessageBus.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusSession) Receive() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Receive(&_MerkleTreeMessageBus.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MerkleTreeMessageBus *MerkleTreeMessageBusTransactorSession) Receive() (*types.Transaction, error) {
	return _MerkleTreeMessageBus.Contract.Receive(&_MerkleTreeMessageBus.TransactOpts)
}

// MerkleTreeMessageBusLogMessagePublishedIterator is returned from FilterLogMessagePublished and is used to iterate over the raw logs and unpacked data for LogMessagePublished events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusLogMessagePublishedIterator struct {
	Event *MerkleTreeMessageBusLogMessagePublished // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusLogMessagePublishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusLogMessagePublished)
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
		it.Event = new(MerkleTreeMessageBusLogMessagePublished)
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
func (it *MerkleTreeMessageBusLogMessagePublishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusLogMessagePublishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusLogMessagePublished represents a LogMessagePublished event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusLogMessagePublished struct {
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
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterLogMessagePublished(opts *bind.FilterOpts) (*MerkleTreeMessageBusLogMessagePublishedIterator, error) {

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusLogMessagePublishedIterator{contract: _MerkleTreeMessageBus.contract, event: "LogMessagePublished", logs: logs, sub: sub}, nil
}

// WatchLogMessagePublished is a free log subscription operation binding the contract event 0xb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937.
//
// Solidity: event LogMessagePublished(address sender, uint64 sequence, uint32 nonce, uint32 topic, bytes payload, uint8 consistencyLevel)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchLogMessagePublished(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusLogMessagePublished) (event.Subscription, error) {

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "LogMessagePublished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusLogMessagePublished)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
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
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseLogMessagePublished(log types.Log) (*MerkleTreeMessageBusLogMessagePublished, error) {
	event := new(MerkleTreeMessageBusLogMessagePublished)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "LogMessagePublished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusOwnershipTransferredIterator struct {
	Event *MerkleTreeMessageBusOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusOwnershipTransferred)
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
		it.Event = new(MerkleTreeMessageBusOwnershipTransferred)
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
func (it *MerkleTreeMessageBusOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusOwnershipTransferred represents a OwnershipTransferred event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MerkleTreeMessageBusOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusOwnershipTransferredIterator{contract: _MerkleTreeMessageBus.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusOwnershipTransferred)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseOwnershipTransferred(log types.Log) (*MerkleTreeMessageBusOwnershipTransferred, error) {
	event := new(MerkleTreeMessageBusOwnershipTransferred)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MerkleTreeMessageBusValueTransferIterator is returned from FilterValueTransfer and is used to iterate over the raw logs and unpacked data for ValueTransfer events raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusValueTransferIterator struct {
	Event *MerkleTreeMessageBusValueTransfer // Event containing the contract specifics and raw log

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
func (it *MerkleTreeMessageBusValueTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MerkleTreeMessageBusValueTransfer)
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
		it.Event = new(MerkleTreeMessageBusValueTransfer)
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
func (it *MerkleTreeMessageBusValueTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MerkleTreeMessageBusValueTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MerkleTreeMessageBusValueTransfer represents a ValueTransfer event raised by the MerkleTreeMessageBus contract.
type MerkleTreeMessageBusValueTransfer struct {
	Sender   common.Address
	Receiver common.Address
	Amount   *big.Int
	Sequence uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterValueTransfer is a free log retrieval operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) FilterValueTransfer(opts *bind.FilterOpts, sender []common.Address, receiver []common.Address) (*MerkleTreeMessageBusValueTransferIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.FilterLogs(opts, "ValueTransfer", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeMessageBusValueTransferIterator{contract: _MerkleTreeMessageBus.contract, event: "ValueTransfer", logs: logs, sub: sub}, nil
}

// WatchValueTransfer is a free log subscription operation binding the contract event 0x50c536ac33a920f00755865b831d17bf4cff0b2e0345f65b16d52bfc004068b6.
//
// Solidity: event ValueTransfer(address indexed sender, address indexed receiver, uint256 amount, uint64 sequence)
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) WatchValueTransfer(opts *bind.WatchOpts, sink chan<- *MerkleTreeMessageBusValueTransfer, sender []common.Address, receiver []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _MerkleTreeMessageBus.contract.WatchLogs(opts, "ValueTransfer", senderRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MerkleTreeMessageBusValueTransfer)
				if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "ValueTransfer", log); err != nil {
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
func (_MerkleTreeMessageBus *MerkleTreeMessageBusFilterer) ParseValueTransfer(log types.Log) (*MerkleTreeMessageBusValueTransfer, error) {
	event := new(MerkleTreeMessageBusValueTransfer)
	if err := _MerkleTreeMessageBus.contract.UnpackLog(event, "ValueTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
