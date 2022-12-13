// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObscuroL2Bridge

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

// ObscuroL2BridgeMetaData contains all meta data concerning the ObscuroL2Bridge contract.
var ObscuroL2BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteBridge\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"localAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"CreatedWrappedToken\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"crossChainAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"createWrappedToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wrappedToken\",\"type\":\"address\"}],\"name\":\"hasTokenMapping\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"localToRemoteToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"remoteToLocalToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"wrappedTokens\",\"outputs\":[{\"internalType\":\"contractObscuroERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040526001805463ffffffff60a01b1916905534801561002057600080fd5b5060405161261938038061261983398101604081905261003f91610121565b600080546001600160a01b0319166001600160a01b038416908117909155604080516350d113fd60e11b8152905184929163a1a227fa916004808301926020929190829003018186803b15801561009557600080fd5b505afa1580156100a9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100cd9190610154565b600180546001600160a01b039283166001600160a01b0319918216179091556005805494909216931692909217909155506101769050565b80516001600160a01b038116811461011c57600080fd5b919050565b6000806040838503121561013457600080fd5b61013d83610105565b915061014b60208401610105565b90509250929050565b60006020828403121561016657600080fd5b61016f82610105565b9392505050565b612494806101856000396000f3fe608060405260043610620000845760003560e01c80639813c7b211620000545780639813c7b214620001435780639e405b711462000193578063c432a46f14620001cd578063d5c6b50414620001f157600080fd5b80628d48e314620000895780631888d71214620000e057806318bcac5014620000f957806383bece4d146200011e575b600080fd5b3480156200009657600080fd5b50620000c3620000a83660046200089d565b6004602052600090815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b620000f7620000f13660046200089d565b6200022b565b005b3480156200010657600080fd5b50620000f76200011836600462000910565b620003f7565b3480156200012b57600080fd5b50620000f76200013d3660046200099c565b62000584565b3480156200015057600080fd5b5062000182620001623660046200089d565b6001600160a01b0390811660009081526002602052604090205416151590565b6040519015158152602001620000d7565b348015620001a057600080fd5b50620000c3620001b23660046200089d565b6003602052600090815260409020546001600160a01b031681565b348015620001da57600080fd5b50620000f7620001ec3660046200099c565b505050565b348015620001fe57600080fd5b50620000c3620002103660046200089d565b6002602052600090815260409020546001600160a01b031681565b60003411620002815760405162461bcd60e51b815260206004820152600d60248201527f4e6f7468696e672073656e742e0000000000000000000000000000000000000060448201526064015b60405180910390fd5b6000805260026020527fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b546001600160a01b0316620003035760405162461bcd60e51b815260206004820152601560248201527f4e6f206d617070696e6720666f7220746f6b656e2e0000000000000000000000604482015260640162000278565b600080805260036020527f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff546040516001600160a01b03918216602482015234604482015290831660648201527f83bece4d000000000000000000000000000000000000000000000000000000009060840160408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152600554909150620003f3906001600160a01b03168260008080620006b8565b5050565b6005546000546001600160a01b03918216911633148015620004335750806001600160a01b031662000428620007e6565b6001600160a01b0316145b620004815760405162461bcd60e51b815260206004820152601b60248201527f4d657373616765206973206e6f742063726f737320636861696e2e0000000000604482015260640162000278565b600085858585604051620004959062000876565b620004a4949392919062000a0c565b604051809103906000f080158015620004c1573d6000803e3d6000fd5b506001600160a01b03808216600081815260026020908152604080832080547fffffffffffffffffffffffff00000000000000000000000000000000000000009081168617909155600383528184208054968f169682168717905594835260049091529081902080549093169091179091555190915081907f30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca9062000572908a9084908b908b908b908b9062000a42565b60405180910390a15050505050505050565b6005546000546001600160a01b03918216911633148015620005c05750806001600160a01b0316620005b5620007e6565b6001600160a01b0316145b6200060e5760405162461bcd60e51b815260206004820152601b60248201527f4d657373616765206973206e6f742063726f737320636861696e2e0000000000604482015260640162000278565b6001600160a01b0380851660009081526002602052604090205416806200063457600080fd5b6040517f979005ad0000000000000000000000000000000000000000000000000000000081526001600160a01b0384811660048301526024820186905282169063979005ad90604401600060405180830381600087803b1580156200069857600080fd5b505af1158015620006ad573d6000803e3d6000fd5b505050505050505050565b60006040518060600160405280876001600160a01b0316815260200186815260200184815250604051602001620006f0919062000ae3565b60408051808303601f19018152919052600180549192506001600160a01b0382169163b1454caa917401000000000000000000000000000000000000000090910463ffffffff16906014620007458362000b2a565b91906101000a81548163ffffffff021916908363ffffffff1602179055508684866040518563ffffffff1660e01b815260040162000787949392919062000b76565b602060405180830381600087803b158015620007a257600080fd5b505af1158015620007b7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007dd919062000bb5565b50505050505050565b60008060009054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b815260040160206040518083038186803b1580156200083657600080fd5b505afa1580156200084b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000871919062000be1565b905090565b61185d8062000c0283390190565b6001600160a01b03811681146200089a57600080fd5b50565b600060208284031215620008b057600080fd5b8135620008bd8162000884565b9392505050565b60008083601f840112620008d757600080fd5b50813567ffffffffffffffff811115620008f057600080fd5b6020830191508360208285010111156200090957600080fd5b9250929050565b6000806000806000606086880312156200092957600080fd5b8535620009368162000884565b9450602086013567ffffffffffffffff808211156200095457600080fd5b6200096289838a01620008c4565b909650945060408801359150808211156200097c57600080fd5b506200098b88828901620008c4565b969995985093965092949392505050565b600080600060608486031215620009b257600080fd5b8335620009bf8162000884565b9250602084013591506040840135620009d88162000884565b809150509250925092565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b60408152600062000a22604083018688620009e3565b828103602084015262000a37818587620009e3565b979650505050505050565b60006001600160a01b0380891683528088166020840152506080604083015262000a71608083018688620009e3565b828103606084015262000a86818587620009e3565b9998505050505050505050565b6000815180845260005b8181101562000abb5760208185018101518683018201520162000a9d565b8181111562000ace576000602083870101525b50601f01601f19169290920160200192915050565b602081526001600160a01b038251166020820152600060208301516060604084015262000b14608084018262000a93565b9050604084015160608401528091505092915050565b600063ffffffff8083168181141562000b6c577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b600063ffffffff80871683528086166020840152506080604083015262000ba1608083018562000a93565b905060ff8316606083015295945050505050565b60006020828403121562000bc857600080fd5b815167ffffffffffffffff81168114620008bd57600080fd5b60006020828403121562000bf457600080fd5b8151620008bd816200088456fe60806040523480156200001157600080fd5b506040516200185d3803806200185d8339810160408190526200003491620002b8565b8151829082906200004d90600390602085019062000145565b5080516200006390600490602084019062000145565b505050620000987fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177533620000a060201b60201c565b50506200035f565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff16620001415760008281526005602090815260408083206001600160a01b03851684529091529020805460ff19166001179055620001003390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b828054620001539062000322565b90600052602060002090601f016020900481019282620001775760008555620001c2565b82601f106200019257805160ff1916838001178555620001c2565b82800160010185558215620001c2579182015b82811115620001c2578251825591602001919060010190620001a5565b50620001d0929150620001d4565b5090565b5b80821115620001d05760008155600101620001d5565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200021357600080fd5b81516001600160401b0380821115620002305762000230620001eb565b604051601f8301601f19908116603f011681019082821181831017156200025b576200025b620001eb565b816040528381526020925086838588010111156200027857600080fd5b600091505b838210156200029c57858201830151818301840152908201906200027d565b83821115620002ae5760008385830101525b9695505050505050565b60008060408385031215620002cc57600080fd5b82516001600160401b0380821115620002e457600080fd5b620002f28683870162000201565b935060208501519150808211156200030957600080fd5b50620003188582860162000201565b9150509250929050565b600181811c908216806200033757607f821691505b602082108114156200035957634e487b7160e01b600052602260045260246000fd5b50919050565b6114ee806200036f6000396000f3fe608060405234801561001057600080fd5b50600436106101775760003560e01c806339509351116100d8578063979005ad1161008c578063a9059cbb11610066578063a9059cbb14610330578063d547741f14610343578063dd62ed3e1461035657600080fd5b8063979005ad14610302578063a217fddf14610315578063a457c2d71461031d57600080fd5b806375b238fc116100bd57806375b238fc1461029a57806391d14854146102c157806395d89b41146102fa57600080fd5b8063395093511461025e57806370a082311461027157600080fd5b806323b872dd1161012f5780632f2ff15d116101145780632f2ff15d14610229578063313ce5671461023c57806336568abe1461024b57600080fd5b806323b872dd146101f3578063248a9ca31461020657600080fd5b8063095ea7b311610160578063095ea7b3146101b957806318160ddd146101cc5780631dd319cb146101de57600080fd5b806301ffc9a71461017c57806306fdde03146101a4575b600080fd5b61018f61018a3660046111a8565b61038f565b60405190151581526020015b60405180910390f35b6101ac610428565b60405161019b9190611216565b61018f6101c7366004611265565b6104ba565b6002545b60405190815260200161019b565b6101f16101ec366004611265565b6104d2565b005b61018f61020136600461128f565b61057f565b6101d06102143660046112cb565b60009081526005602052604090206001015490565b6101f16102373660046112e4565b6105a3565b6040516012815260200161019b565b6101f16102593660046112e4565b6105c9565b61018f61026c366004611265565b610655565b6101d061027f366004611310565b6001600160a01b031660009081526020819052604090205490565b6101d07fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b61018f6102cf3660046112e4565b60009182526005602090815260408084206001600160a01b0393909316845291905290205460ff1690565b6101ac610694565b6101f1610310366004611265565b6106a3565b6101d0600081565b61018f61032b366004611265565b6106d8565b61018f61033e366004611265565b610782565b6101f16103513660046112e4565b610790565b6101d061036436600461132b565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061042257507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b60606003805461043790611355565b80601f016020809104026020016040519081016040528092919081815260200182805461046390611355565b80156104b05780601f10610485576101008083540402835291602001916104b0565b820191906000526020600020905b81548152906001019060200180831161049357829003601f168201915b5050505050905090565b6000336104c88185856107b6565b5060019392505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756104fd813361090e565b8161051d846001600160a01b031660009081526020819052604090205490565b10156105705760405162461bcd60e51b815260206004820152601560248201527f496e73756666696369656e742062616c616e63652e000000000000000000000060448201526064015b60405180910390fd5b61057a838361098e565b505050565b60003361058d858285610b13565b610598858585610ba5565b506001949350505050565b6000828152600560205260409020600101546105bf813361090e565b61057a8383610dbc565b6001600160a01b03811633146106475760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201527f20726f6c657320666f722073656c6600000000000000000000000000000000006064820152608401610567565b6106518282610e5e565b5050565b3360008181526001602090815260408083206001600160a01b03871684529091528120549091906104c8908290869061068f9087906113a6565b6107b6565b60606004805461043790611355565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756106ce813361090e565b61057a8383610ee1565b3360008181526001602090815260408083206001600160a01b0387168452909152812054909190838110156107755760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f0000000000000000000000000000000000000000000000000000006064820152608401610567565b61059882868684036107b6565b6000336104c8818585610ba5565b6000828152600560205260409020600101546107ac813361090e565b61057a8383610e5e565b6001600160a01b0383166108315760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f72657373000000000000000000000000000000000000000000000000000000006064820152608401610567565b6001600160a01b0382166108ad5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f73730000000000000000000000000000000000000000000000000000000000006064820152608401610567565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff166106515761094c816001600160a01b03166014610fc0565b610957836020610fc0565b6040516020016109689291906113be565b60408051601f198184030181529082905262461bcd60e51b825261056791600401611216565b6001600160a01b038216610a0a5760405162461bcd60e51b815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f73000000000000000000000000000000000000000000000000000000000000006064820152608401610567565b6001600160a01b03821660009081526020819052604090205481811015610a995760405162461bcd60e51b815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f63650000000000000000000000000000000000000000000000000000000000006064820152608401610567565b6001600160a01b0383166000908152602081905260408120838303905560028054849290610ac890849061143f565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b038381166000908152600160209081526040808320938616835292905220546000198114610b9f5781811015610b925760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e63650000006044820152606401610567565b610b9f84848484036107b6565b50505050565b6001600160a01b038316610c215760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f64726573730000000000000000000000000000000000000000000000000000006064820152608401610567565b6001600160a01b038216610c9d5760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f65737300000000000000000000000000000000000000000000000000000000006064820152608401610567565b6001600160a01b03831660009081526020819052604090205481811015610d2c5760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e636500000000000000000000000000000000000000000000000000006064820152608401610567565b6001600160a01b03808516600090815260208190526040808220858503905591851681529081208054849290610d639084906113a6565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610daf91815260200190565b60405180910390a3610b9f565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff166106515760008281526005602090815260408083206001600160a01b03851684529091529020805460ff19166001179055610e1a3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff16156106515760008281526005602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b6001600160a01b038216610f375760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f2061646472657373006044820152606401610567565b8060026000828254610f4991906113a6565b90915550506001600160a01b03821660009081526020819052604081208054839290610f769084906113a6565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b60606000610fcf836002611456565b610fda9060026113a6565b67ffffffffffffffff811115610ff257610ff2611475565b6040519080825280601f01601f19166020018201604052801561101c576020820181803683370190505b5090507f3000000000000000000000000000000000000000000000000000000000000000816000815181106110535761105361148b565b60200101906001600160f81b031916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061109e5761109e61148b565b60200101906001600160f81b031916908160001a90535060006110c2846002611456565b6110cd9060016113a6565b90505b6001811115611152577f303132333435363738396162636465660000000000000000000000000000000085600f166010811061110e5761110e61148b565b1a60f81b8282815181106111245761112461148b565b60200101906001600160f81b031916908160001a90535060049490941c9361114b816114a1565b90506110d0565b5083156111a15760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152606401610567565b9392505050565b6000602082840312156111ba57600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146111a157600080fd5b60005b838110156112055781810151838201526020016111ed565b83811115610b9f5750506000910152565b60208152600082518060208401526112358160408501602087016111ea565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461126057600080fd5b919050565b6000806040838503121561127857600080fd5b61128183611249565b946020939093013593505050565b6000806000606084860312156112a457600080fd5b6112ad84611249565b92506112bb60208501611249565b9150604084013590509250925092565b6000602082840312156112dd57600080fd5b5035919050565b600080604083850312156112f757600080fd5b8235915061130760208401611249565b90509250929050565b60006020828403121561132257600080fd5b6111a182611249565b6000806040838503121561133e57600080fd5b61134783611249565b915061130760208401611249565b600181811c9082168061136957607f821691505b6020821081141561138a57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b600082198211156113b9576113b9611390565b500190565b7f416363657373436f6e74726f6c3a206163636f756e74200000000000000000008152600083516113f68160178501602088016111ea565b7f206973206d697373696e6720726f6c652000000000000000000000000000000060179184019182015283516114338160288401602088016111ea565b01602801949350505050565b60008282101561145157611451611390565b500390565b600081600019048311821515161561147057611470611390565b500290565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000816114b0576114b0611390565b50600019019056fea2646970667358221220d28b4fc0e891bd0f3bc48a9e65bd48145e89755c3ef32d3a232ac6b1f585595064736f6c63430008090033a2646970667358221220579b58ec18b42bffb10e5e8d90eb48bebf6dea2059b819496501776eb7712dae64736f6c63430008090033",
}

// ObscuroL2BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ObscuroL2BridgeMetaData.ABI instead.
var ObscuroL2BridgeABI = ObscuroL2BridgeMetaData.ABI

// ObscuroL2BridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObscuroL2BridgeMetaData.Bin instead.
var ObscuroL2BridgeBin = ObscuroL2BridgeMetaData.Bin

// DeployObscuroL2Bridge deploys a new Ethereum contract, binding an instance of ObscuroL2Bridge to it.
func DeployObscuroL2Bridge(auth *bind.TransactOpts, backend bind.ContractBackend, messenger common.Address, remoteBridge common.Address) (common.Address, *types.Transaction, *ObscuroL2Bridge, error) {
	parsed, err := ObscuroL2BridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObscuroL2BridgeBin), backend, messenger, remoteBridge)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObscuroL2Bridge{ObscuroL2BridgeCaller: ObscuroL2BridgeCaller{contract: contract}, ObscuroL2BridgeTransactor: ObscuroL2BridgeTransactor{contract: contract}, ObscuroL2BridgeFilterer: ObscuroL2BridgeFilterer{contract: contract}}, nil
}

// ObscuroL2Bridge is an auto generated Go binding around an Ethereum contract.
type ObscuroL2Bridge struct {
	ObscuroL2BridgeCaller     // Read-only binding to the contract
	ObscuroL2BridgeTransactor // Write-only binding to the contract
	ObscuroL2BridgeFilterer   // Log filterer for contract events
}

// ObscuroL2BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ObscuroL2BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroL2BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ObscuroL2BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroL2BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObscuroL2BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroL2BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObscuroL2BridgeSession struct {
	Contract     *ObscuroL2Bridge  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObscuroL2BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObscuroL2BridgeCallerSession struct {
	Contract *ObscuroL2BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ObscuroL2BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObscuroL2BridgeTransactorSession struct {
	Contract     *ObscuroL2BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ObscuroL2BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ObscuroL2BridgeRaw struct {
	Contract *ObscuroL2Bridge // Generic contract binding to access the raw methods on
}

// ObscuroL2BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObscuroL2BridgeCallerRaw struct {
	Contract *ObscuroL2BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ObscuroL2BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObscuroL2BridgeTransactorRaw struct {
	Contract *ObscuroL2BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewObscuroL2Bridge creates a new instance of ObscuroL2Bridge, bound to a specific deployed contract.
func NewObscuroL2Bridge(address common.Address, backend bind.ContractBackend) (*ObscuroL2Bridge, error) {
	contract, err := bindObscuroL2Bridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObscuroL2Bridge{ObscuroL2BridgeCaller: ObscuroL2BridgeCaller{contract: contract}, ObscuroL2BridgeTransactor: ObscuroL2BridgeTransactor{contract: contract}, ObscuroL2BridgeFilterer: ObscuroL2BridgeFilterer{contract: contract}}, nil
}

// NewObscuroL2BridgeCaller creates a new read-only instance of ObscuroL2Bridge, bound to a specific deployed contract.
func NewObscuroL2BridgeCaller(address common.Address, caller bind.ContractCaller) (*ObscuroL2BridgeCaller, error) {
	contract, err := bindObscuroL2Bridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroL2BridgeCaller{contract: contract}, nil
}

// NewObscuroL2BridgeTransactor creates a new write-only instance of ObscuroL2Bridge, bound to a specific deployed contract.
func NewObscuroL2BridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ObscuroL2BridgeTransactor, error) {
	contract, err := bindObscuroL2Bridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroL2BridgeTransactor{contract: contract}, nil
}

// NewObscuroL2BridgeFilterer creates a new log filterer instance of ObscuroL2Bridge, bound to a specific deployed contract.
func NewObscuroL2BridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ObscuroL2BridgeFilterer, error) {
	contract, err := bindObscuroL2Bridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObscuroL2BridgeFilterer{contract: contract}, nil
}

// bindObscuroL2Bridge binds a generic wrapper to an already deployed contract.
func bindObscuroL2Bridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ObscuroL2BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroL2Bridge *ObscuroL2BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroL2Bridge.Contract.ObscuroL2BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroL2Bridge *ObscuroL2BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.ObscuroL2BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroL2Bridge *ObscuroL2BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.ObscuroL2BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroL2Bridge *ObscuroL2BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroL2Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.contract.Transact(opts, method, params...)
}

// HasTokenMapping is a free data retrieval call binding the contract method 0x9813c7b2.
//
// Solidity: function hasTokenMapping(address wrappedToken) view returns(bool)
func (_ObscuroL2Bridge *ObscuroL2BridgeCaller) HasTokenMapping(opts *bind.CallOpts, wrappedToken common.Address) (bool, error) {
	var out []interface{}
	err := _ObscuroL2Bridge.contract.Call(opts, &out, "hasTokenMapping", wrappedToken)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasTokenMapping is a free data retrieval call binding the contract method 0x9813c7b2.
//
// Solidity: function hasTokenMapping(address wrappedToken) view returns(bool)
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) HasTokenMapping(wrappedToken common.Address) (bool, error) {
	return _ObscuroL2Bridge.Contract.HasTokenMapping(&_ObscuroL2Bridge.CallOpts, wrappedToken)
}

// HasTokenMapping is a free data retrieval call binding the contract method 0x9813c7b2.
//
// Solidity: function hasTokenMapping(address wrappedToken) view returns(bool)
func (_ObscuroL2Bridge *ObscuroL2BridgeCallerSession) HasTokenMapping(wrappedToken common.Address) (bool, error) {
	return _ObscuroL2Bridge.Contract.HasTokenMapping(&_ObscuroL2Bridge.CallOpts, wrappedToken)
}

// LocalToRemoteToken is a free data retrieval call binding the contract method 0x9e405b71.
//
// Solidity: function localToRemoteToken(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeCaller) LocalToRemoteToken(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ObscuroL2Bridge.contract.Call(opts, &out, "localToRemoteToken", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LocalToRemoteToken is a free data retrieval call binding the contract method 0x9e405b71.
//
// Solidity: function localToRemoteToken(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) LocalToRemoteToken(arg0 common.Address) (common.Address, error) {
	return _ObscuroL2Bridge.Contract.LocalToRemoteToken(&_ObscuroL2Bridge.CallOpts, arg0)
}

// LocalToRemoteToken is a free data retrieval call binding the contract method 0x9e405b71.
//
// Solidity: function localToRemoteToken(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeCallerSession) LocalToRemoteToken(arg0 common.Address) (common.Address, error) {
	return _ObscuroL2Bridge.Contract.LocalToRemoteToken(&_ObscuroL2Bridge.CallOpts, arg0)
}

// RemoteToLocalToken is a free data retrieval call binding the contract method 0x008d48e3.
//
// Solidity: function remoteToLocalToken(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeCaller) RemoteToLocalToken(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ObscuroL2Bridge.contract.Call(opts, &out, "remoteToLocalToken", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RemoteToLocalToken is a free data retrieval call binding the contract method 0x008d48e3.
//
// Solidity: function remoteToLocalToken(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) RemoteToLocalToken(arg0 common.Address) (common.Address, error) {
	return _ObscuroL2Bridge.Contract.RemoteToLocalToken(&_ObscuroL2Bridge.CallOpts, arg0)
}

// RemoteToLocalToken is a free data retrieval call binding the contract method 0x008d48e3.
//
// Solidity: function remoteToLocalToken(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeCallerSession) RemoteToLocalToken(arg0 common.Address) (common.Address, error) {
	return _ObscuroL2Bridge.Contract.RemoteToLocalToken(&_ObscuroL2Bridge.CallOpts, arg0)
}

// WrappedTokens is a free data retrieval call binding the contract method 0xd5c6b504.
//
// Solidity: function wrappedTokens(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeCaller) WrappedTokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ObscuroL2Bridge.contract.Call(opts, &out, "wrappedTokens", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedTokens is a free data retrieval call binding the contract method 0xd5c6b504.
//
// Solidity: function wrappedTokens(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) WrappedTokens(arg0 common.Address) (common.Address, error) {
	return _ObscuroL2Bridge.Contract.WrappedTokens(&_ObscuroL2Bridge.CallOpts, arg0)
}

// WrappedTokens is a free data retrieval call binding the contract method 0xd5c6b504.
//
// Solidity: function wrappedTokens(address ) view returns(address)
func (_ObscuroL2Bridge *ObscuroL2BridgeCallerSession) WrappedTokens(arg0 common.Address) (common.Address, error) {
	return _ObscuroL2Bridge.Contract.WrappedTokens(&_ObscuroL2Bridge.CallOpts, arg0)
}

// CreateWrappedToken is a paid mutator transaction binding the contract method 0x18bcac50.
//
// Solidity: function createWrappedToken(address crossChainAddress, string name, string symbol) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactor) CreateWrappedToken(opts *bind.TransactOpts, crossChainAddress common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroL2Bridge.contract.Transact(opts, "createWrappedToken", crossChainAddress, name, symbol)
}

// CreateWrappedToken is a paid mutator transaction binding the contract method 0x18bcac50.
//
// Solidity: function createWrappedToken(address crossChainAddress, string name, string symbol) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) CreateWrappedToken(crossChainAddress common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.CreateWrappedToken(&_ObscuroL2Bridge.TransactOpts, crossChainAddress, name, symbol)
}

// CreateWrappedToken is a paid mutator transaction binding the contract method 0x18bcac50.
//
// Solidity: function createWrappedToken(address crossChainAddress, string name, string symbol) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactorSession) CreateWrappedToken(crossChainAddress common.Address, name string, symbol string) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.CreateWrappedToken(&_ObscuroL2Bridge.TransactOpts, crossChainAddress, name, symbol)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.ReceiveAssets(&_ObscuroL2Bridge.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.ReceiveAssets(&_ObscuroL2Bridge.TransactOpts, asset, amount, receiver)
}

// SendAssets is a paid mutator transaction binding the contract method 0xc432a46f.
//
// Solidity: function sendAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactor) SendAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.contract.Transact(opts, "sendAssets", asset, amount, receiver)
}

// SendAssets is a paid mutator transaction binding the contract method 0xc432a46f.
//
// Solidity: function sendAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) SendAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.SendAssets(&_ObscuroL2Bridge.TransactOpts, asset, amount, receiver)
}

// SendAssets is a paid mutator transaction binding the contract method 0xc432a46f.
//
// Solidity: function sendAssets(address asset, uint256 amount, address receiver) returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactorSession) SendAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.SendAssets(&_ObscuroL2Bridge.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactor) SendNative(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.contract.Transact(opts, "sendNative", receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.SendNative(&_ObscuroL2Bridge.TransactOpts, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_ObscuroL2Bridge *ObscuroL2BridgeTransactorSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _ObscuroL2Bridge.Contract.SendNative(&_ObscuroL2Bridge.TransactOpts, receiver)
}

// ObscuroL2BridgeCreatedWrappedTokenIterator is returned from FilterCreatedWrappedToken and is used to iterate over the raw logs and unpacked data for CreatedWrappedToken events raised by the ObscuroL2Bridge contract.
type ObscuroL2BridgeCreatedWrappedTokenIterator struct {
	Event *ObscuroL2BridgeCreatedWrappedToken // Event containing the contract specifics and raw log

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
func (it *ObscuroL2BridgeCreatedWrappedTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroL2BridgeCreatedWrappedToken)
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
		it.Event = new(ObscuroL2BridgeCreatedWrappedToken)
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
func (it *ObscuroL2BridgeCreatedWrappedTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroL2BridgeCreatedWrappedTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroL2BridgeCreatedWrappedToken represents a CreatedWrappedToken event raised by the ObscuroL2Bridge contract.
type ObscuroL2BridgeCreatedWrappedToken struct {
	RemoteAddress common.Address
	LocalAddress  common.Address
	Name          string
	Symbol        string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCreatedWrappedToken is a free log retrieval operation binding the contract event 0x30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca.
//
// Solidity: event CreatedWrappedToken(address remoteAddress, address localAddress, string name, string symbol)
func (_ObscuroL2Bridge *ObscuroL2BridgeFilterer) FilterCreatedWrappedToken(opts *bind.FilterOpts) (*ObscuroL2BridgeCreatedWrappedTokenIterator, error) {

	logs, sub, err := _ObscuroL2Bridge.contract.FilterLogs(opts, "CreatedWrappedToken")
	if err != nil {
		return nil, err
	}
	return &ObscuroL2BridgeCreatedWrappedTokenIterator{contract: _ObscuroL2Bridge.contract, event: "CreatedWrappedToken", logs: logs, sub: sub}, nil
}

// WatchCreatedWrappedToken is a free log subscription operation binding the contract event 0x30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca.
//
// Solidity: event CreatedWrappedToken(address remoteAddress, address localAddress, string name, string symbol)
func (_ObscuroL2Bridge *ObscuroL2BridgeFilterer) WatchCreatedWrappedToken(opts *bind.WatchOpts, sink chan<- *ObscuroL2BridgeCreatedWrappedToken) (event.Subscription, error) {

	logs, sub, err := _ObscuroL2Bridge.contract.WatchLogs(opts, "CreatedWrappedToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroL2BridgeCreatedWrappedToken)
				if err := _ObscuroL2Bridge.contract.UnpackLog(event, "CreatedWrappedToken", log); err != nil {
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

// ParseCreatedWrappedToken is a log parse operation binding the contract event 0x30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca.
//
// Solidity: event CreatedWrappedToken(address remoteAddress, address localAddress, string name, string symbol)
func (_ObscuroL2Bridge *ObscuroL2BridgeFilterer) ParseCreatedWrappedToken(log types.Log) (*ObscuroL2BridgeCreatedWrappedToken, error) {
	event := new(ObscuroL2BridgeCreatedWrappedToken)
	if err := _ObscuroL2Bridge.contract.UnpackLog(event, "CreatedWrappedToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
