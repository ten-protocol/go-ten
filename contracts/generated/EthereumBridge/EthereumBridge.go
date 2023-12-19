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
	Bin: "0x60806040526001805463ffffffff60a01b1916905534801561001f575f80fd5b506125ca8061002d5f395ff3fe608060405260043610620000c5575f3560e01c806383bece4d11620000725780639e405b7111620000545780639e405b7114620002d3578063a381c8e2146200030b578063d5c6b504146200032f576200013e565b806383bece4d14620002615780639813c7b21462000285576200013e565b8063458ffd6311620000a8578063458ffd6314620001f5578063485cc955146200021957806375cb2672146200023d576200013e565b80628d48e314620001875780631888d71214620001dc576200013e565b366200013e5760405162461bcd60e51b815260206004820152602360248201527f436f6e747261637420646f6573206e6f7420737570706f72742072656365697660448201527f652829000000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b60405162461bcd60e51b815260206004820152601d60248201527f66616c6c6261636b2829206d6574686f6420756e737570706f72746564000000604482015260640162000135565b34801562000193575f80fd5b50620001bf620001a536600462000f49565b60046020525f90815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b620001f3620001ed36600462000f49565b62000367565b005b34801562000201575f80fd5b50620001f36200021336600462000fb6565b62000513565b34801562000225575f80fd5b50620001f3620002373660046200103c565b62000720565b34801562000249575f80fd5b50620001f36200025b36600462000f49565b6200088c565b3480156200026d575f80fd5b50620001f36200027f36600462001078565b62000969565b34801562000291575f80fd5b50620002c2620002a336600462000f49565b6001600160a01b039081165f9081526002602052604090205416151590565b6040519015158152602001620001d3565b348015620002df575f80fd5b50620001bf620002f136600462000f49565b60036020525f90815260409020546001600160a01b031681565b34801562000317575f80fd5b50620001f36200032936600462001078565b62000ba5565b3480156200033b575f80fd5b50620001bf6200034d36600462000f49565b60026020525f90815260409020546001600160a01b031681565b5f3411620003b85760405162461bcd60e51b815260206004820152600d60248201527f4e6f7468696e672073656e742e00000000000000000000000000000000000000604482015260640162000135565b5f805260026020527fac33ff75c19e70fe83507db0d683fd3465c996598dc972688b7ace676c89077b546001600160a01b0316620004395760405162461bcd60e51b815260206004820152601560248201527f4e6f206d617070696e6720666f7220746f6b656e2e0000000000000000000000604482015260640162000135565b5f80805260036020527f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff546040516001600160a01b03918216602482015234604482015290831660648201526383bece4d60e01b9060840160408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526005549091506200050f906001600160a01b0316825f5b5f8062000d2c565b5050565b6005545f546001600160a01b03918216911633146200059b5760405162461bcd60e51b815260206004820152603060248201527f436f6e74726163742063616c6c6572206973206e6f742074686520726567697360448201527f7465726564206d657373656e6765722100000000000000000000000000000000606482015260840162000135565b806001600160a01b0316620005af62000e3d565b6001600160a01b0316146200062d5760405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f72726563742073656e64657221000000000000000000000000000000606482015260840162000135565b5f85858585604051620006409062000f23565b6200064f9493929190620010e4565b604051809103905ff08015801562000669573d5f803e3d5ffd5b506001600160a01b038082165f818152600260209081526040808320805473ffffffffffffffffffffffffffffffffffffffff199081168617909155600383528184208054968f169682168717905594835260049091529081902080549093169091179091555190915081907f30c05779f384e0ae9d43bbf7ec4417f28bdc53d02a35551b6eb270a9c4c71dca906200070e908a9084908b908b908b908b9062001119565b60405180910390a15050505050505050565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff165f811580156200076b5750825b90505f8267ffffffffffffffff166001148015620007885750303b155b90508115801562000797575080155b15620007cf576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156200080457845468ff00000000000000001916680100000000000000001785555b6200080f876200088c565b6005805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03881617905583156200088357845468ff000000000000000019168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b6200089662000eb9565b5f805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b038316908117909155604080517fa1a227fa000000000000000000000000000000000000000000000000000000008152905163a1a227fa916004808201926020929091908290030181865afa15801562000913573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062000939919062001169565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b039290921691909117905550565b6005545f546001600160a01b0391821691163314620009f15760405162461bcd60e51b815260206004820152603060248201527f436f6e74726163742063616c6c6572206973206e6f742074686520726567697360448201527f7465726564206d657373656e6765722100000000000000000000000000000000606482015260840162000135565b806001600160a01b031662000a0562000e3d565b6001600160a01b03161462000a835760405162461bcd60e51b815260206004820152603160248201527f43726f737320636861696e206d65737361676520636f6d696e672066726f6d2060448201527f696e636f72726563742073656e64657221000000000000000000000000000000606482015260840162000135565b6001600160a01b038085165f9081526004602090815260408083205484168084526002909252909120549091168062000b255760405162461bcd60e51b815260206004820152602b60248201527f526563656976696e672061737365747320666f7220756e6b6e6f776e2077726160448201527f7070656420746f6b656e21000000000000000000000000000000000000000000606482015260840162000135565b6040517f979005ad0000000000000000000000000000000000000000000000000000000081526001600160a01b0385811660048301526024820187905282169063979005ad906044015f604051808303815f87803b15801562000b86575f80fd5b505af115801562000b99573d5f803e3d5ffd5b50505050505050505050565b6001600160a01b038084165f908152600260205260409020541662000c0d5760405162461bcd60e51b815260206004820152601560248201527f4e6f206d617070696e6720666f7220746f6b656e2e0000000000000000000000604482015260640162000135565b6001600160a01b038381165f90815260026020526040908190205490517f1dd319cb000000000000000000000000000000000000000000000000000000008152336004820152602481018590529116908190631dd319cb906044015f604051808303815f87803b15801562000c80575f80fd5b505af115801562000c93573d5f803e3d5ffd5b505050506001600160a01b038481165f90815260036020908152604080832054815190851660248201526044810188905286851660648083019190915282518083039091018152608490910190915290810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff166383bece4d60e01b179052600554909262000d2592911690839062000507565b5050505050565b5f6040518060600160405280876001600160a01b031681526020018681526020018481525060405160200162000d639190620011cc565b60408051808303601f19018152919052600180549192506001600160a01b0382169163b1454caa917401000000000000000000000000000000000000000090910463ffffffff1690601462000db88362001212565b91906101000a81548163ffffffff021916908363ffffffff1602179055508684866040518563ffffffff1660e01b815260040162000dfa94939291906200125a565b6020604051808303815f875af115801562000e17573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062000883919062001298565b5f805f9054906101000a90046001600160a01b03166001600160a01b03166363012de56040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000e8e573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062000eb4919062001169565b905090565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff1662000f21576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6112d380620012c283390190565b6001600160a01b038116811462000f46575f80fd5b50565b5f6020828403121562000f5a575f80fd5b813562000f678162000f31565b9392505050565b5f8083601f84011262000f7f575f80fd5b50813567ffffffffffffffff81111562000f97575f80fd5b60208301915083602082850101111562000faf575f80fd5b9250929050565b5f805f805f6060868803121562000fcb575f80fd5b853562000fd88162000f31565b9450602086013567ffffffffffffffff8082111562000ff5575f80fd5b6200100389838a0162000f6e565b909650945060408801359150808211156200101c575f80fd5b506200102b8882890162000f6e565b969995985093965092949392505050565b5f80604083850312156200104e575f80fd5b82356200105b8162000f31565b915060208301356200106d8162000f31565b809150509250929050565b5f805f606084860312156200108b575f80fd5b8335620010988162000f31565b9250602084013591506040840135620010b18162000f31565b809150509250925092565b81835281816020850137505f828201602090810191909152601f909101601f19169091010190565b604081525f620010f9604083018688620010bc565b82810360208401526200110e818587620010bc565b979650505050505050565b5f6001600160a01b0380891683528088166020840152506080604083015262001147608083018688620010bc565b82810360608401526200115c818587620010bc565b9998505050505050505050565b5f602082840312156200117a575f80fd5b815162000f678162000f31565b5f81518084525f5b81811015620011ad576020818501810151868301820152016200118f565b505f602082860101526020601f19601f83011685010191505092915050565b602081526001600160a01b0382511660208201525f602083015160606040840152620011fc608084018262001187565b9050604084015160608401528091505092915050565b5f63ffffffff80831681810362001250577f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b6001019392505050565b5f63ffffffff80871683528086166020840152506080604083015262001284608083018562001187565b905060ff8316606083015295945050505050565b5f60208284031215620012a9575f80fd5b815167ffffffffffffffff8116811462000f67575f80fdfe6080604052600580546001600160a01b03191673deb34a740eca1ec42c8b8204cbec0ba34fdd27f317905534801562000036575f80fd5b50604051620012d3380380620012d3833981016040819052620000599162000228565b8181818160036200006b83826200031a565b5060046200007a82826200031a565b5050505050620000b17fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177533620000ba60201b60201c565b505050620003e2565b5f8281526007602090815260408083206001600160a01b038516845290915281205460ff1662000161575f8381526007602090815260408083206001600160a01b03861684529091529020805460ff19166001179055620001183390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600162000164565b505f5b92915050565b634e487b7160e01b5f52604160045260245ffd5b5f82601f8301126200018e575f80fd5b81516001600160401b0380821115620001ab57620001ab6200016a565b604051601f8301601f19908116603f01168101908282118183101715620001d657620001d66200016a565b81604052838152602092508683858801011115620001f2575f80fd5b5f91505b83821015620002155785820183015181830184015290820190620001f6565b5f93810190920192909252949350505050565b5f80604083850312156200023a575f80fd5b82516001600160401b038082111562000251575f80fd5b6200025f868387016200017e565b9350602085015191508082111562000275575f80fd5b5062000284858286016200017e565b9150509250929050565b600181811c90821680620002a357607f821691505b602082108103620002c257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111562000315575f81815260208120601f850160051c81016020861015620002f05750805b601f850160051c820191505b818110156200031157828155600101620002fc565b5050505b505050565b81516001600160401b038111156200033657620003366200016a565b6200034e816200034784546200028e565b84620002c8565b602080601f83116001811462000384575f84156200036c5750858301515b5f19600386901b1c1916600185901b17855562000311565b5f85815260208120601f198616915b82811015620003b45788860151825594840194600190910190840162000393565b5085821015620003d257878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b610ee380620003f05f395ff3fe608060405234801561000f575f80fd5b5060043610610149575f3560e01c806336568abe116100c7578063979005ad1161007d578063a9059cbb11610063578063a9059cbb146102c2578063d547741f146102d5578063dd62ed3e146102e8575f80fd5b8063979005ad146102a8578063a217fddf146102bb575f80fd5b806375b238fc116100ad57806375b238fc1461024157806391d148541461026857806395d89b41146102a0575f80fd5b806336568abe1461021b57806370a082311461022e575f80fd5b80631dd319cb1161011c578063248a9ca311610102578063248a9ca3146101d75780632f2ff15d146101f9578063313ce5671461020c575f80fd5b80631dd319cb146101af57806323b872dd146101c4575f80fd5b806301ffc9a71461014d57806306fdde0314610175578063095ea7b31461018a57806318160ddd1461019d575b5f80fd5b61016061015b366004610cc7565b6102fb565b60405190151581526020015b60405180910390f35b61017d610393565b60405161016c9190610d0d565b610160610198366004610d73565b610423565b6002545b60405190815260200161016c565b6101c26101bd366004610d73565b61043a565b005b6101606101d2366004610d9b565b6104d0565b6101a16101e5366004610dd4565b5f9081526007602052604090206001015490565b6101c2610207366004610deb565b6104f3565b6040516012815260200161016c565b6101c2610229366004610deb565b61051d565b6101a161023c366004610e15565b610569565b6101a17fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b610160610276366004610deb565b5f9182526007602090815260408084206001600160a01b0393909316845291905290205460ff1690565b61017d61060c565b6101c26102b6366004610d73565b61061b565b6101a15f81565b6101606102d0366004610d73565b61064f565b6101c26102e3366004610deb565b61065c565b6101a16102f6366004610e2e565b610680565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061038d57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6060600380546103a290610e56565b80601f01602080910402602001604051908101604052809291908181526020018280546103ce90610e56565b80156104195780601f106103f057610100808354040283529160200191610419565b820191905f5260205f20905b8154815290600101906020018083116103fc57829003601f168201915b5050505050905090565b5f3361043081858561078e565b5060019392505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756104648161079b565b8161046e84610569565b10156104c15760405162461bcd60e51b815260206004820152601560248201527f496e73756666696369656e742062616c616e63652e000000000000000000000060448201526064015b60405180910390fd5b6104cb83836107a8565b505050565b5f336104dd8582856107e0565b6104e8858585610856565b506001949350505050565b5f8281526007602052604090206001015461050d8161079b565b61051783836108b3565b50505050565b6001600160a01b038116331461055f576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104cb828261095e565b5f6001600160a01b0382163203610597576001600160a01b0382165f9081526020819052604090205461038d565b6001600160a01b03821633036105c4576001600160a01b0382165f9081526020819052604090205461038d565b60405162461bcd60e51b815260206004820152601f60248201527f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e63650060448201526064016104b8565b6060600480546103a290610e56565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756106458161079b565b6104cb83836109e3565b5f33610430818585610856565b5f828152600760205260409020600101546106768161079b565b610517838361095e565b5f326001600160a01b03841614806106a05750326001600160a01b038316145b156106d2576001600160a01b038084165f908152600160209081526040808320938616835292905220545b905061038d565b336001600160a01b03841614806106f15750336001600160a01b038316145b15610720576001600160a01b038084165f908152600160209081526040808320938616835292905220546106cb565b60405162461bcd60e51b815260206004820152602160248201527f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360448201527f650000000000000000000000000000000000000000000000000000000000000060648201526084016104b8565b6104cb8383836001610a17565b6107a58133610b1b565b50565b6001600160a01b0382166107d157604051634b637e8f60e11b81525f60048201526024016104b8565b6107dc825f83610b88565b5050565b5f6107eb8484610680565b90505f1981146105175781811015610848576040517ffb8f41b20000000000000000000000000000000000000000000000000000000081526001600160a01b038416600482015260248101829052604481018390526064016104b8565b61051784848484035f610a17565b6001600160a01b03831661087f57604051634b637e8f60e11b81525f60048201526024016104b8565b6001600160a01b0382166108a85760405163ec442f0560e01b81525f60048201526024016104b8565b6104cb838383610b88565b5f8281526007602090815260408083206001600160a01b038516845290915281205460ff16610957575f8381526007602090815260408083206001600160a01b03861684529091529020805460ff1916600117905561090f3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161038d565b505f61038d565b5f8281526007602090815260408083206001600160a01b038516845290915281205460ff1615610957575f8381526007602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161038d565b6001600160a01b038216610a0c5760405163ec442f0560e01b81525f60048201526024016104b8565b6107dc5f8383610b88565b6001600160a01b038416610a59576040517fe602df050000000000000000000000000000000000000000000000000000000081525f60048201526024016104b8565b6001600160a01b038316610a9b576040517f94280d620000000000000000000000000000000000000000000000000000000081525f60048201526024016104b8565b6001600160a01b038085165f908152600160209081526040808320938716835292905220829055801561051757826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610b0d91815260200190565b60405180910390a350505050565b5f8281526007602090815260408083206001600160a01b038516845290915290205460ff166107dc576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602481018390526044016104b8565b6001600160a01b038316610bb2578060025f828254610ba79190610e8e565b90915550610c3b9050565b6001600160a01b0383165f9081526020819052604090205481811015610c1d576040517fe450d38c0000000000000000000000000000000000000000000000000000000081526001600160a01b038516600482015260248101829052604481018390526064016104b8565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b038216610c5757600280548290039055610c75565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610cba91815260200190565b60405180910390a3505050565b5f60208284031215610cd7575f80fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610d06575f80fd5b9392505050565b5f6020808352835180828501525f5b81811015610d3857858101830151858201604001528201610d1c565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610d6e575f80fd5b919050565b5f8060408385031215610d84575f80fd5b610d8d83610d58565b946020939093013593505050565b5f805f60608486031215610dad575f80fd5b610db684610d58565b9250610dc460208501610d58565b9150604084013590509250925092565b5f60208284031215610de4575f80fd5b5035919050565b5f8060408385031215610dfc575f80fd5b82359150610e0c60208401610d58565b90509250929050565b5f60208284031215610e25575f80fd5b610d0682610d58565b5f8060408385031215610e3f575f80fd5b610e4883610d58565b9150610e0c60208401610d58565b600181811c90821680610e6a57607f821691505b602082108103610e8857634e487b7160e01b5f52602260045260245ffd5b50919050565b8082018082111561038d57634e487b7160e01b5f52601160045260245ffdfea2646970667358221220d78cf2dd9ecf4da03679034ab54dbc38a0748a2a50f13ee501925e8df4738b1264736f6c63430008140033a2646970667358221220910ed4d806fc5d201f041052722e240d19224d7a3957d6d54e271e5fa89dde2664736f6c63430008140033",
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
