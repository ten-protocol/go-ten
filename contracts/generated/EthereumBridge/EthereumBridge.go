// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EthereumBridge

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

// EthereumBridgeMetaData contains all meta data concerning the EthereumBridge contract.
var EthereumBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"remoteAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"localAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"CreatedWrappedToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messengerAddress\",\"type\":\"address\"}],\"name\":\"configure\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wrappedToken\",\"type\":\"address\"}],\"name\":\"hasTokenMapping\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"messenger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"remoteBridge\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"localToRemoteToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"crossChainAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"onCreateTokenCommand\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"receiveAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"remoteToLocalToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"sendNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"wrappedTokens\",\"outputs\":[{\"internalType\":\"contractWrappedERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040526001805463ffffffff60a01b1916905534801561002057600080fd5b50612681806100306000396000f3fe608060405260043610620000c65760003560e01c806383bece4d11620000735780639e405b7111620000555780639e405b7114620002dc578063a381c8e21462000316578063d5c6b504146200033b576200013f565b806383bece4d14620002675780639813c7b2146200028c576200013f565b8063458ffd6311620000a9578063458ffd6314620001f8578063485cc955146200021d57806375cb26721462000242576200013f565b80628d48e314620001885780631888d71214620001df576200013f565b366200013f5760405162461bcd60e51b815260206004820152602360248201527f436f6e747261637420646f6573206e6f7420737570706f72742072656365697660448201527f652829000000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b60405162461bcd60e51b815260206004820152601d60248201527f66616c6c6261636b2829206d6574686f6420756e737570706f72746564000000604482015260640162000136565b3480156200019557600080fd5b50620001c2620001a736600462000f7f565b6004602052600090815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b620001f6620001f036600462000f7f565b62000375565b005b3480156200020557600080fd5b50620001f66200021736600462000ff2565b62000526565b3480156200022a57600080fd5b50620001f66200023c3660046200107e565b62000739565b3480156200024f57600080fd5b50620001f66200026136600462000f7f565b620008a7565b3480156200027457600080fd5b50620001f662000286366004620010bc565b62000987565b3480156200029957600080fd5b50620002cb620002ab36600462000f7f565b6001600160a01b0390811660009081526002602052604090205416151590565b6040519015158152602001620001d6565b348015620002e957600080fd5b50620001c2620002fb36600462000f7f565b6003602052600090815260409020546001600160a01b031681565b3480156200032357600080fd5b50620001f662000335366004620010bc565b62000bca565b3480156200034857600080fd5b50620001c26200035a36600462000f7f565b6002602052600090815260409020546001600160a01b031681565b60003411620003c75760405162461bcd60e51b815260206004820152600d60248201527f4e6f7468696e672073656e742e00000000000000000000000000000000000000604482015260640162000136565b6000805260026020527fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b546001600160a01b0316620004495760405162461bcd60e51b815260206004820152601560248201527f4e6f206d617070696e6720666f7220746f6b656e2e0000000000000000000000604482015260640162000136565b600080805260036020527f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff546040516001600160a01b03918216602482015234604482015290831660648201526383bece4d60e01b9060840160408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915260055490915062000522906001600160a01b03168260005b60008062000d59565b5050565b6005546000546001600160a01b0391821691163314620005af5760405162461bcd60e51b815260206004820152603060248201527f436f6e74726163742063616c6c6572206973206e6f742074686520726567697360448201527f7465726564206d657373656e6765722100000000000000000000000000000000606482015260840162000136565b806001600160a01b0316620005c362000e6e565b6001600160a01b031614620006415760405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f72726563742073656e64657221000000000000000000000000000000606482015260840162000136565b600085858585604051620006559062000f58565b6200066494939291906200112c565b604051809103906000f08015801562000681573d6000803e3d6000fd5b506001600160a01b038082166000818152600260209081526040808320805473ffffffffffffffffffffffffffffffffffffffff199081168617909155600383528184208054968f169682168717905594835260049091529081902080549093169091179091555190915081907f30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca9062000727908a9084908b908b908b908b9062001162565b60405180910390a15050505050505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff16600081158015620007855750825b905060008267ffffffffffffffff166001148015620007a35750303b155b905081158015620007b2575080155b15620007ea576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156200081f57845468ff00000000000000001916680100000000000000001785555b6200082a87620008a7565b6005805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03881617905583156200089e57845468ff000000000000000019168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b620008b162000eee565b6000805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117909155604080517fa1a227fa000000000000000000000000000000000000000000000000000000008152905163a1a227fa916004808201926020929091908290030181865afa15801562000931573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620009579190620011b3565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039290921691909117905550565b6005546000546001600160a01b039182169116331462000a105760405162461bcd60e51b815260206004820152603060248201527f436f6e74726163742063616c6c6572206973206e6f742074686520726567697360448201527f7465726564206d657373656e6765722100000000000000000000000000000000606482015260840162000136565b806001600160a01b031662000a2462000e6e565b6001600160a01b03161462000aa25760405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f72726563742073656e64657221000000000000000000000000000000606482015260840162000136565b6001600160a01b0380851660009081526004602090815260408083205484168084526002909252909120549091168062000b455760405162461bcd60e51b815260206004820152602b60248201527f526563656976696e672061737365747320666f7220756e6b6e6f776e2077726160448201527f7070656420746f6b656e21000000000000000000000000000000000000000000606482015260840162000136565b6040517f979005ad0000000000000000000000000000000000000000000000000000000081526001600160a01b0385811660048301526024820187905282169063979005ad90604401600060405180830381600087803b15801562000ba957600080fd5b505af115801562000bbe573d6000803e3d6000fd5b50505050505050505050565b6001600160a01b038084166000908152600260205260409020541662000c335760405162461bcd60e51b815260206004820152601560248201527f4e6f206d617070696e6720666f7220746f6b656e2e0000000000000000000000604482015260640162000136565b6001600160a01b03838116600090815260026020526040908190205490517f1dd319cb000000000000000000000000000000000000000000000000000000008152336004820152602481018590529116908190631dd319cb90604401600060405180830381600087803b15801562000caa57600080fd5b505af115801562000cbf573d6000803e3d6000fd5b505050506001600160a01b03848116600090815260036020908152604080832054815190851660248201526044810188905286851660648083019190915282518083039091018152608490910190915290810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff166383bece4d60e01b179052600554909262000d5292911690839062000519565b5050505050565b60006040518060600160405280876001600160a01b031681526020018681526020018481525060405160200162000d9191906200121b565b60408051808303601f19018152919052600180549192506001600160a01b0382169163b1454caa917401000000000000000000000000000000000000000090910463ffffffff1690601462000de68362001262565b91906101000a81548163ffffffff021916908363ffffffff1602179055508684866040518563ffffffff1660e01b815260040162000e289493929190620012ad565b6020604051808303816000875af115801562000e48573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200089e9190620012ec565b60008060009054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000ec3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000ee99190620011b3565b905090565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff1662000f56576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b611333806200131983390190565b6001600160a01b038116811462000f7c57600080fd5b50565b60006020828403121562000f9257600080fd5b813562000f9f8162000f66565b9392505050565b60008083601f84011262000fb957600080fd5b50813567ffffffffffffffff81111562000fd257600080fd5b60208301915083602082850101111562000feb57600080fd5b9250929050565b6000806000806000606086880312156200100b57600080fd5b8535620010188162000f66565b9450602086013567ffffffffffffffff808211156200103657600080fd5b6200104489838a0162000fa6565b909650945060408801359150808211156200105e57600080fd5b506200106d8882890162000fa6565b969995985093965092949392505050565b600080604083850312156200109257600080fd5b82356200109f8162000f66565b91506020830135620010b18162000f66565b809150509250929050565b600080600060608486031215620010d257600080fd5b8335620010df8162000f66565b9250602084013591506040840135620010f88162000f66565b809150509250925092565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6040815260006200114260408301868862001103565b82810360208401526200115781858762001103565b979650505050505050565b60006001600160a01b038089168352808816602084015250608060408301526200119160808301868862001103565b8281036060840152620011a681858762001103565b9998505050505050505050565b600060208284031215620011c657600080fd5b815162000f9f8162000f66565b6000815180845260005b81811015620011fb57602081850181015186830182015201620011dd565b506000602082860101526020601f19601f83011685010191505092915050565b602081526001600160a01b03825116602082015260006020830151606060408401526200124c6080840182620011d3565b9050604084015160608401528091505092915050565b600063ffffffff808316818103620012a3577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6001019392505050565b600063ffffffff808716835280861660208401525060806040830152620012d86080830185620011d3565b905060ff8316606083015295945050505050565b600060208284031215620012ff57600080fd5b815167ffffffffffffffff8116811462000f9f57600080fdfe6080604052600580546001600160a01b03191673deb34a740eca1ec42c8b8204cbec0ba34fdd27f31790553480156200003757600080fd5b5060405162001333380380620013338339810160408190526200005a9162000233565b8181818160036200006c83826200032c565b5060046200007b82826200032c565b5050505050620000b27fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177533620000bb60201b60201c565b505050620003f8565b60008281526007602090815260408083206001600160a01b038516845290915281205460ff16620001645760008381526007602090815260408083206001600160a01b03861684529091529020805460ff191660011790556200011b3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600162000168565b5060005b92915050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200019657600080fd5b81516001600160401b0380821115620001b357620001b36200016e565b604051601f8301601f19908116603f01168101908282118183101715620001de57620001de6200016e565b81604052838152602092508683858801011115620001fb57600080fd5b600091505b838210156200021f578582018301518183018401529082019062000200565b600093810190920192909252949350505050565b600080604083850312156200024757600080fd5b82516001600160401b03808211156200025f57600080fd5b6200026d8683870162000184565b935060208501519150808211156200028457600080fd5b50620002938582860162000184565b9150509250929050565b600181811c90821680620002b257607f821691505b602082108103620002d357634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200032757600081815260208120601f850160051c81016020861015620003025750805b601f850160051c820191505b8181101562000323578281556001016200030e565b5050505b505050565b81516001600160401b038111156200034857620003486200016e565b62000360816200035984546200029d565b84620002d9565b602080601f8311600181146200039857600084156200037f5750858301515b600019600386901b1c1916600185901b17855562000323565b600085815260208120601f198616915b82811015620003c957888601518255948401946001909101908401620003a8565b5085821015620003e85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610f2b80620004086000396000f3fe608060405234801561001057600080fd5b50600436106101515760003560e01c806336568abe116100cd578063979005ad11610081578063a9059cbb11610066578063a9059cbb146102ce578063d547741f146102e1578063dd62ed3e146102f457600080fd5b8063979005ad146102b3578063a217fddf146102c657600080fd5b806375b238fc116100b257806375b238fc1461024b57806391d148541461027257806395d89b41146102ab57600080fd5b806336568abe1461022557806370a082311461023857600080fd5b80631dd319cb11610124578063248a9ca311610109578063248a9ca3146101e05780632f2ff15d14610203578063313ce5671461021657600080fd5b80631dd319cb146101b857806323b872dd146101cd57600080fd5b806301ffc9a71461015657806306fdde031461017e578063095ea7b31461019357806318160ddd146101a6575b600080fd5b610169610164366004610cf7565b610307565b60405190151581526020015b60405180910390f35b6101866103a0565b6040516101759190610d40565b6101696101a1366004610daa565b610432565b6002545b604051908152602001610175565b6101cb6101c6366004610daa565b61044a565b005b6101696101db366004610dd4565b6104e0565b6101aa6101ee366004610e10565b60009081526007602052604090206001015490565b6101cb610211366004610e29565b610504565b60405160128152602001610175565b6101cb610233366004610e29565b61052f565b6101aa610246366004610e55565b61057b565b6101aa7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b610169610280366004610e29565b60009182526007602090815260408084206001600160a01b0393909316845291905290205460ff1690565b610186610621565b6101cb6102c1366004610daa565b610630565b6101aa600081565b6101696102dc366004610daa565b610664565b6101cb6102ef366004610e29565b610672565b6101aa610302366004610e70565b610697565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061039a57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6060600380546103af90610e9a565b80601f01602080910402602001604051908101604052809291908181526020018280546103db90610e9a565b80156104285780601f106103fd57610100808354040283529160200191610428565b820191906000526020600020905b81548152906001019060200180831161040b57829003601f168201915b5050505050905090565b6000336104408185856107a8565b5060019392505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c21775610474816107b5565b8161047e8461057b565b10156104d15760405162461bcd60e51b815260206004820152601560248201527f496e73756666696369656e742062616c616e63652e000000000000000000000060448201526064015b60405180910390fd5b6104db83836107c2565b505050565b6000336104ee8582856107fc565b6104f9858585610875565b506001949350505050565b60008281526007602052604090206001015461051f816107b5565b61052983836108d4565b50505050565b6001600160a01b0381163314610571576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104db8282610982565b60006001600160a01b03821632036105ab576001600160a01b03821660009081526020819052604090205461039a565b6001600160a01b03821633036105d9576001600160a01b03821660009081526020819052604090205461039a565b60405162461bcd60e51b815260206004820152601f60248201527f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e63650060448201526064016104c8565b6060600480546103af90610e9a565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561065a816107b5565b6104db8383610a09565b600033610440818585610875565b60008281526007602052604090206001015461068d816107b5565b6105298383610982565b6000326001600160a01b03841614806106b85750326001600160a01b038316145b156106eb576001600160a01b038084166000908152600160209081526040808320938616835292905220545b905061039a565b336001600160a01b038416148061070a5750336001600160a01b038316145b1561073a576001600160a01b038084166000908152600160209081526040808320938616835292905220546106e4565b60405162461bcd60e51b815260206004820152602160248201527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360448201527f650000000000000000000000000000000000000000000000000000000000000060648201526084016104c8565b6104db8383836001610a3f565b6107bf8133610b46565b50565b6001600160a01b0382166107ec57604051634b637e8f60e11b8152600060048201526024016104c8565b6107f882600083610bb4565b5050565b60006108088484610697565b905060001981146105295781811015610866576040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526001600160a01b038416600482015260248101829052604481018390526064016104c8565b61052984848484036000610a3f565b6001600160a01b03831661089f57604051634b637e8f60e11b8152600060048201526024016104c8565b6001600160a01b0382166108c95760405163ec442f0560e01b8152600060048201526024016104c8565b6104db838383610bb4565b60008281526007602090815260408083206001600160a01b038516845290915281205460ff1661097a5760008381526007602090815260408083206001600160a01b03861684529091529020805460ff191660011790556109323390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161039a565b50600061039a565b60008281526007602090815260408083206001600160a01b038516845290915281205460ff161561097a5760008381526007602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161039a565b6001600160a01b038216610a335760405163ec442f0560e01b8152600060048201526024016104c8565b6107f860008383610bb4565b6001600160a01b038416610a82576040517fe602df05000000000000000000000000000000000000000000000000000000008152600060048201526024016104c8565b6001600160a01b038316610ac5576040517f94280d62000000000000000000000000000000000000000000000000000000008152600060048201526024016104c8565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561052957826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610b3891815260200190565b60405180910390a350505050565b60008281526007602090815260408083206001600160a01b038516845290915290205460ff166107f8576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602481018390526044016104c8565b6001600160a01b038316610bdf578060026000828254610bd49190610ed4565b90915550610c6a9050565b6001600160a01b03831660009081526020819052604090205481811015610c4b576040517fe450d38c0000000000000000000000000000000000000000000000000000000081526001600160a01b038516600482015260248101829052604481018390526064016104c8565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b038216610c8657600280548290039055610ca5565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610cea91815260200190565b60405180910390a3505050565b600060208284031215610d0957600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610d3957600080fd5b9392505050565b600060208083528351808285015260005b81811015610d6d57858101830151858201604001528201610d51565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610da557600080fd5b919050565b60008060408385031215610dbd57600080fd5b610dc683610d8e565b946020939093013593505050565b600080600060608486031215610de957600080fd5b610df284610d8e565b9250610e0060208501610d8e565b9150604084013590509250925092565b600060208284031215610e2257600080fd5b5035919050565b60008060408385031215610e3c57600080fd5b82359150610e4c60208401610d8e565b90509250929050565b600060208284031215610e6757600080fd5b610d3982610d8e565b60008060408385031215610e8357600080fd5b610e8c83610d8e565b9150610e4c60208401610d8e565b600181811c90821680610eae57607f821691505b602082108103610ece57634e487b7160e01b600052602260045260246000fd5b50919050565b8082018082111561039a57634e487b7160e01b600052601160045260246000fdfea2646970667358221220b77e0982f6ea01ef5acc9548caa6305649f5edd8efcbda658dff65d6882af5d364736f6c63430008140033a26469706673582212209e0aa1c97b5d192d64f989aec4b86d5ce2bfbe5d5f345cb45245d90bad25709264736f6c63430008140033",
}

// EthereumBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use EthereumBridgeMetaData.ABI instead.
var EthereumBridgeABI = EthereumBridgeMetaData.ABI

// EthereumBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EthereumBridgeMetaData.Bin instead.
var EthereumBridgeBin = EthereumBridgeMetaData.Bin

// DeployEthereumBridge deploys a new Ethereum contract, binding an instance of EthereumBridge to it.
func DeployEthereumBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EthereumBridge, error) {
	parsed, err := EthereumBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EthereumBridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EthereumBridge{EthereumBridgeCaller: EthereumBridgeCaller{contract: contract}, EthereumBridgeTransactor: EthereumBridgeTransactor{contract: contract}, EthereumBridgeFilterer: EthereumBridgeFilterer{contract: contract}}, nil
}

// EthereumBridge is an auto generated Go binding around an Ethereum contract.
type EthereumBridge struct {
	EthereumBridgeCaller     // Read-only binding to the contract
	EthereumBridgeTransactor // Write-only binding to the contract
	EthereumBridgeFilterer   // Log filterer for contract events
}

// EthereumBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthereumBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthereumBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthereumBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthereumBridgeSession struct {
	Contract     *EthereumBridge   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthereumBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthereumBridgeCallerSession struct {
	Contract *EthereumBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// EthereumBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthereumBridgeTransactorSession struct {
	Contract     *EthereumBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// EthereumBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthereumBridgeRaw struct {
	Contract *EthereumBridge // Generic contract binding to access the raw methods on
}

// EthereumBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthereumBridgeCallerRaw struct {
	Contract *EthereumBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// EthereumBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthereumBridgeTransactorRaw struct {
	Contract *EthereumBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthereumBridge creates a new instance of EthereumBridge, bound to a specific deployed contract.
func NewEthereumBridge(address common.Address, backend bind.ContractBackend) (*EthereumBridge, error) {
	contract, err := bindEthereumBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthereumBridge{EthereumBridgeCaller: EthereumBridgeCaller{contract: contract}, EthereumBridgeTransactor: EthereumBridgeTransactor{contract: contract}, EthereumBridgeFilterer: EthereumBridgeFilterer{contract: contract}}, nil
}

// NewEthereumBridgeCaller creates a new read-only instance of EthereumBridge, bound to a specific deployed contract.
func NewEthereumBridgeCaller(address common.Address, caller bind.ContractCaller) (*EthereumBridgeCaller, error) {
	contract, err := bindEthereumBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumBridgeCaller{contract: contract}, nil
}

// NewEthereumBridgeTransactor creates a new write-only instance of EthereumBridge, bound to a specific deployed contract.
func NewEthereumBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*EthereumBridgeTransactor, error) {
	contract, err := bindEthereumBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumBridgeTransactor{contract: contract}, nil
}

// NewEthereumBridgeFilterer creates a new log filterer instance of EthereumBridge, bound to a specific deployed contract.
func NewEthereumBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*EthereumBridgeFilterer, error) {
	contract, err := bindEthereumBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthereumBridgeFilterer{contract: contract}, nil
}

// bindEthereumBridge binds a generic wrapper to an already deployed contract.
func bindEthereumBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EthereumBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumBridge *EthereumBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthereumBridge.Contract.EthereumBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumBridge *EthereumBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumBridge.Contract.EthereumBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumBridge *EthereumBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumBridge.Contract.EthereumBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumBridge *EthereumBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthereumBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumBridge *EthereumBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumBridge *EthereumBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumBridge.Contract.contract.Transact(opts, method, params...)
}

// HasTokenMapping is a free data retrieval call binding the contract method 0x9813c7b2.
//
// Solidity: function hasTokenMapping(address wrappedToken) view returns(bool)
func (_EthereumBridge *EthereumBridgeCaller) HasTokenMapping(opts *bind.CallOpts, wrappedToken common.Address) (bool, error) {
	var out []interface{}
	err := _EthereumBridge.contract.Call(opts, &out, "hasTokenMapping", wrappedToken)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasTokenMapping is a free data retrieval call binding the contract method 0x9813c7b2.
//
// Solidity: function hasTokenMapping(address wrappedToken) view returns(bool)
func (_EthereumBridge *EthereumBridgeSession) HasTokenMapping(wrappedToken common.Address) (bool, error) {
	return _EthereumBridge.Contract.HasTokenMapping(&_EthereumBridge.CallOpts, wrappedToken)
}

// HasTokenMapping is a free data retrieval call binding the contract method 0x9813c7b2.
//
// Solidity: function hasTokenMapping(address wrappedToken) view returns(bool)
func (_EthereumBridge *EthereumBridgeCallerSession) HasTokenMapping(wrappedToken common.Address) (bool, error) {
	return _EthereumBridge.Contract.HasTokenMapping(&_EthereumBridge.CallOpts, wrappedToken)
}

// LocalToRemoteToken is a free data retrieval call binding the contract method 0x9e405b71.
//
// Solidity: function localToRemoteToken(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeCaller) LocalToRemoteToken(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _EthereumBridge.contract.Call(opts, &out, "localToRemoteToken", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LocalToRemoteToken is a free data retrieval call binding the contract method 0x9e405b71.
//
// Solidity: function localToRemoteToken(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeSession) LocalToRemoteToken(arg0 common.Address) (common.Address, error) {
	return _EthereumBridge.Contract.LocalToRemoteToken(&_EthereumBridge.CallOpts, arg0)
}

// LocalToRemoteToken is a free data retrieval call binding the contract method 0x9e405b71.
//
// Solidity: function localToRemoteToken(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeCallerSession) LocalToRemoteToken(arg0 common.Address) (common.Address, error) {
	return _EthereumBridge.Contract.LocalToRemoteToken(&_EthereumBridge.CallOpts, arg0)
}

// RemoteToLocalToken is a free data retrieval call binding the contract method 0x008d48e3.
//
// Solidity: function remoteToLocalToken(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeCaller) RemoteToLocalToken(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _EthereumBridge.contract.Call(opts, &out, "remoteToLocalToken", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RemoteToLocalToken is a free data retrieval call binding the contract method 0x008d48e3.
//
// Solidity: function remoteToLocalToken(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeSession) RemoteToLocalToken(arg0 common.Address) (common.Address, error) {
	return _EthereumBridge.Contract.RemoteToLocalToken(&_EthereumBridge.CallOpts, arg0)
}

// RemoteToLocalToken is a free data retrieval call binding the contract method 0x008d48e3.
//
// Solidity: function remoteToLocalToken(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeCallerSession) RemoteToLocalToken(arg0 common.Address) (common.Address, error) {
	return _EthereumBridge.Contract.RemoteToLocalToken(&_EthereumBridge.CallOpts, arg0)
}

// WrappedTokens is a free data retrieval call binding the contract method 0xd5c6b504.
//
// Solidity: function wrappedTokens(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeCaller) WrappedTokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _EthereumBridge.contract.Call(opts, &out, "wrappedTokens", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedTokens is a free data retrieval call binding the contract method 0xd5c6b504.
//
// Solidity: function wrappedTokens(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeSession) WrappedTokens(arg0 common.Address) (common.Address, error) {
	return _EthereumBridge.Contract.WrappedTokens(&_EthereumBridge.CallOpts, arg0)
}

// WrappedTokens is a free data retrieval call binding the contract method 0xd5c6b504.
//
// Solidity: function wrappedTokens(address ) view returns(address)
func (_EthereumBridge *EthereumBridgeCallerSession) WrappedTokens(arg0 common.Address) (common.Address, error) {
	return _EthereumBridge.Contract.WrappedTokens(&_EthereumBridge.CallOpts, arg0)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_EthereumBridge *EthereumBridgeTransactor) Configure(opts *bind.TransactOpts, messengerAddress common.Address) (*types.Transaction, error) {
	return _EthereumBridge.contract.Transact(opts, "configure", messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_EthereumBridge *EthereumBridgeSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.Configure(&_EthereumBridge.TransactOpts, messengerAddress)
}

// Configure is a paid mutator transaction binding the contract method 0x75cb2672.
//
// Solidity: function configure(address messengerAddress) returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) Configure(messengerAddress common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.Configure(&_EthereumBridge.TransactOpts, messengerAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address remoteBridge) returns()
func (_EthereumBridge *EthereumBridgeTransactor) Initialize(opts *bind.TransactOpts, messenger common.Address, remoteBridge common.Address) (*types.Transaction, error) {
	return _EthereumBridge.contract.Transact(opts, "initialize", messenger, remoteBridge)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address remoteBridge) returns()
func (_EthereumBridge *EthereumBridgeSession) Initialize(messenger common.Address, remoteBridge common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.Initialize(&_EthereumBridge.TransactOpts, messenger, remoteBridge)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address messenger, address remoteBridge) returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) Initialize(messenger common.Address, remoteBridge common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.Initialize(&_EthereumBridge.TransactOpts, messenger, remoteBridge)
}

// OnCreateTokenCommand is a paid mutator transaction binding the contract method 0x458ffd63.
//
// Solidity: function onCreateTokenCommand(address crossChainAddress, string name, string symbol) returns()
func (_EthereumBridge *EthereumBridgeTransactor) OnCreateTokenCommand(opts *bind.TransactOpts, crossChainAddress common.Address, name string, symbol string) (*types.Transaction, error) {
	return _EthereumBridge.contract.Transact(opts, "onCreateTokenCommand", crossChainAddress, name, symbol)
}

// OnCreateTokenCommand is a paid mutator transaction binding the contract method 0x458ffd63.
//
// Solidity: function onCreateTokenCommand(address crossChainAddress, string name, string symbol) returns()
func (_EthereumBridge *EthereumBridgeSession) OnCreateTokenCommand(crossChainAddress common.Address, name string, symbol string) (*types.Transaction, error) {
	return _EthereumBridge.Contract.OnCreateTokenCommand(&_EthereumBridge.TransactOpts, crossChainAddress, name, symbol)
}

// OnCreateTokenCommand is a paid mutator transaction binding the contract method 0x458ffd63.
//
// Solidity: function onCreateTokenCommand(address crossChainAddress, string name, string symbol) returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) OnCreateTokenCommand(crossChainAddress common.Address, name string, symbol string) (*types.Transaction, error) {
	return _EthereumBridge.Contract.OnCreateTokenCommand(&_EthereumBridge.TransactOpts, crossChainAddress, name, symbol)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_EthereumBridge *EthereumBridgeTransactor) ReceiveAssets(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.contract.Transact(opts, "receiveAssets", asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_EthereumBridge *EthereumBridgeSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.ReceiveAssets(&_EthereumBridge.TransactOpts, asset, amount, receiver)
}

// ReceiveAssets is a paid mutator transaction binding the contract method 0x83bece4d.
//
// Solidity: function receiveAssets(address asset, uint256 amount, address receiver) returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) ReceiveAssets(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.ReceiveAssets(&_EthereumBridge.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_EthereumBridge *EthereumBridgeTransactor) SendERC20(opts *bind.TransactOpts, asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.contract.Transact(opts, "sendERC20", asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_EthereumBridge *EthereumBridgeSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.SendERC20(&_EthereumBridge.TransactOpts, asset, amount, receiver)
}

// SendERC20 is a paid mutator transaction binding the contract method 0xa381c8e2.
//
// Solidity: function sendERC20(address asset, uint256 amount, address receiver) returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) SendERC20(asset common.Address, amount *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.SendERC20(&_EthereumBridge.TransactOpts, asset, amount, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_EthereumBridge *EthereumBridgeTransactor) SendNative(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.contract.Transact(opts, "sendNative", receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_EthereumBridge *EthereumBridgeSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.SendNative(&_EthereumBridge.TransactOpts, receiver)
}

// SendNative is a paid mutator transaction binding the contract method 0x1888d712.
//
// Solidity: function sendNative(address receiver) payable returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) SendNative(receiver common.Address) (*types.Transaction, error) {
	return _EthereumBridge.Contract.SendNative(&_EthereumBridge.TransactOpts, receiver)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_EthereumBridge *EthereumBridgeTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _EthereumBridge.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_EthereumBridge *EthereumBridgeSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _EthereumBridge.Contract.Fallback(&_EthereumBridge.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _EthereumBridge.Contract.Fallback(&_EthereumBridge.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EthereumBridge *EthereumBridgeTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumBridge.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EthereumBridge *EthereumBridgeSession) Receive() (*types.Transaction, error) {
	return _EthereumBridge.Contract.Receive(&_EthereumBridge.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_EthereumBridge *EthereumBridgeTransactorSession) Receive() (*types.Transaction, error) {
	return _EthereumBridge.Contract.Receive(&_EthereumBridge.TransactOpts)
}

// EthereumBridgeCreatedWrappedTokenIterator is returned from FilterCreatedWrappedToken and is used to iterate over the raw logs and unpacked data for CreatedWrappedToken events raised by the EthereumBridge contract.
type EthereumBridgeCreatedWrappedTokenIterator struct {
	Event *EthereumBridgeCreatedWrappedToken // Event containing the contract specifics and raw log

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
func (it *EthereumBridgeCreatedWrappedTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumBridgeCreatedWrappedToken)
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
		it.Event = new(EthereumBridgeCreatedWrappedToken)
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
func (it *EthereumBridgeCreatedWrappedTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumBridgeCreatedWrappedTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumBridgeCreatedWrappedToken represents a CreatedWrappedToken event raised by the EthereumBridge contract.
type EthereumBridgeCreatedWrappedToken struct {
	RemoteAddress common.Address
	LocalAddress  common.Address
	Name          string
	Symbol        string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCreatedWrappedToken is a free log retrieval operation binding the contract event 0x30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca.
//
// Solidity: event CreatedWrappedToken(address remoteAddress, address localAddress, string name, string symbol)
func (_EthereumBridge *EthereumBridgeFilterer) FilterCreatedWrappedToken(opts *bind.FilterOpts) (*EthereumBridgeCreatedWrappedTokenIterator, error) {

	logs, sub, err := _EthereumBridge.contract.FilterLogs(opts, "CreatedWrappedToken")
	if err != nil {
		return nil, err
	}
	return &EthereumBridgeCreatedWrappedTokenIterator{contract: _EthereumBridge.contract, event: "CreatedWrappedToken", logs: logs, sub: sub}, nil
}

// WatchCreatedWrappedToken is a free log subscription operation binding the contract event 0x30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca.
//
// Solidity: event CreatedWrappedToken(address remoteAddress, address localAddress, string name, string symbol)
func (_EthereumBridge *EthereumBridgeFilterer) WatchCreatedWrappedToken(opts *bind.WatchOpts, sink chan<- *EthereumBridgeCreatedWrappedToken) (event.Subscription, error) {

	logs, sub, err := _EthereumBridge.contract.WatchLogs(opts, "CreatedWrappedToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumBridgeCreatedWrappedToken)
				if err := _EthereumBridge.contract.UnpackLog(event, "CreatedWrappedToken", log); err != nil {
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
func (_EthereumBridge *EthereumBridgeFilterer) ParseCreatedWrappedToken(log types.Log) (*EthereumBridgeCreatedWrappedToken, error) {
	event := new(EthereumBridgeCreatedWrappedToken)
	if err := _EthereumBridge.contract.UnpackLog(event, "CreatedWrappedToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthereumBridgeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the EthereumBridge contract.
type EthereumBridgeInitializedIterator struct {
	Event *EthereumBridgeInitialized // Event containing the contract specifics and raw log

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
func (it *EthereumBridgeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumBridgeInitialized)
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
		it.Event = new(EthereumBridgeInitialized)
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
func (it *EthereumBridgeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumBridgeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumBridgeInitialized represents a Initialized event raised by the EthereumBridge contract.
type EthereumBridgeInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_EthereumBridge *EthereumBridgeFilterer) FilterInitialized(opts *bind.FilterOpts) (*EthereumBridgeInitializedIterator, error) {

	logs, sub, err := _EthereumBridge.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &EthereumBridgeInitializedIterator{contract: _EthereumBridge.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_EthereumBridge *EthereumBridgeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *EthereumBridgeInitialized) (event.Subscription, error) {

	logs, sub, err := _EthereumBridge.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumBridgeInitialized)
				if err := _EthereumBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_EthereumBridge *EthereumBridgeFilterer) ParseInitialized(log types.Log) (*EthereumBridgeInitialized, error) {
	event := new(EthereumBridgeInitialized)
	if err := _EthereumBridge.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
