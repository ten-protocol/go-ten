// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObscuroERC20

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

// ObscuroERC20MetaData contains all meta data concerning the ObscuroERC20 contract.
var ObscuroERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"giver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issueFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001a1b38038062001a1b8339810160408190526200003491620002b8565b8151829082906200004d90600390602085019062000145565b5080516200006390600490602084019062000145565b505050620000987fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177533620000a060201b60201c565b50506200035f565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff16620001415760008281526005602090815260408083206001600160a01b03851684529091529020805460ff19166001179055620001003390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45b5050565b828054620001539062000322565b90600052602060002090601f016020900481019282620001775760008555620001c2565b82601f106200019257805160ff1916838001178555620001c2565b82800160010185558215620001c2579182015b82811115620001c2578251825591602001919060010190620001a5565b50620001d0929150620001d4565b5090565b5b80821115620001d05760008155600101620001d5565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200021357600080fd5b81516001600160401b0380821115620002305762000230620001eb565b604051601f8301601f19908116603f011681019082821181831017156200025b576200025b620001eb565b816040528381526020925086838588010111156200027857600080fd5b600091505b838210156200029c57858201830151818301840152908201906200027d565b83821115620002ae5760008385830101525b9695505050505050565b60008060408385031215620002cc57600080fd5b82516001600160401b0380821115620002e457600080fd5b620002f28683870162000201565b935060208501519150808211156200030957600080fd5b50620003188582860162000201565b9150509250929050565b600181811c908216806200033757607f821691505b602082108114156200035957634e487b7160e01b600052602260045260246000fd5b50919050565b6116ac806200036f6000396000f3fe60806040526004361061016e5760003560e01c806339509351116100cb578063979005ad1161007f578063a9059cbb11610059578063a9059cbb146104cc578063d547741f146104ec578063dd62ed3e1461050c576101e6565b8063979005ad14610477578063a217fddf14610497578063a457c2d7146104ac576101e6565b806375b238fc116100b057806375b238fc146103e857806391d148541461041c57806395d89b4114610462576101e6565b8063395093511461039257806370a08231146103b2576101e6565b806323b872dd116101225780632f2ff15d116101075780632f2ff15d14610336578063313ce5671461035657806336568abe14610372576101e6565b806323b872dd146102e6578063248a9ca314610306576101e6565b8063095ea7b311610153578063095ea7b31461028557806318160ddd146102a55780631dd319cb146102c4576101e6565b806301ffc9a71461022e57806306fdde0314610263576101e6565b366101e65760405162461bcd60e51b815260206004820152602360248201527f436f6e747261637420646f6573206e6f7420737570706f72742072656365697660448201527f652829000000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b60405162461bcd60e51b815260206004820152601d60248201527f66616c6c6261636b2829206d6574686f6420756e737570706f7274656400000060448201526064016101dd565b34801561023a57600080fd5b5061024e610249366004611366565b610552565b60405190151581526020015b60405180910390f35b34801561026f57600080fd5b506102786105eb565b60405161025a91906113d4565b34801561029157600080fd5b5061024e6102a0366004611423565b61067d565b3480156102b157600080fd5b506002545b60405190815260200161025a565b3480156102d057600080fd5b506102e46102df366004611423565b610695565b005b3480156102f257600080fd5b5061024e61030136600461144d565b61073d565b34801561031257600080fd5b506102b6610321366004611489565b60009081526005602052604090206001015490565b34801561034257600080fd5b506102e46103513660046114a2565b610761565b34801561036257600080fd5b506040516012815260200161025a565b34801561037e57600080fd5b506102e461038d3660046114a2565b610787565b34801561039e57600080fd5b5061024e6103ad366004611423565b610813565b3480156103be57600080fd5b506102b66103cd3660046114ce565b6001600160a01b031660009081526020819052604090205490565b3480156103f457600080fd5b506102b67fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177581565b34801561042857600080fd5b5061024e6104373660046114a2565b60009182526005602090815260408084206001600160a01b0393909316845291905290205460ff1690565b34801561046e57600080fd5b50610278610852565b34801561048357600080fd5b506102e4610492366004611423565b610861565b3480156104a357600080fd5b506102b6600081565b3480156104b857600080fd5b5061024e6104c7366004611423565b610896565b3480156104d857600080fd5b5061024e6104e7366004611423565b610940565b3480156104f857600080fd5b506102e46105073660046114a2565b61094e565b34801561051857600080fd5b506102b66105273660046114e9565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806105e557507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6060600380546105fa90611513565b80601f016020809104026020016040519081016040528092919081815260200182805461062690611513565b80156106735780601f1061064857610100808354040283529160200191610673565b820191906000526020600020905b81548152906001019060200180831161065657829003601f168201915b5050505050905090565b60003361068b818585610974565b5060019392505050565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c217756106c08133610acc565b816106e0846001600160a01b031660009081526020819052604090205490565b101561072e5760405162461bcd60e51b815260206004820152601560248201527f496e73756666696369656e742062616c616e63652e000000000000000000000060448201526064016101dd565b6107388383610b4c565b505050565b60003361074b858285610cd1565b610756858585610d63565b506001949350505050565b60008281526005602052604090206001015461077d8133610acc565b6107388383610f7a565b6001600160a01b03811633146108055760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201527f20726f6c657320666f722073656c66000000000000000000000000000000000060648201526084016101dd565b61080f828261101c565b5050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490919061068b908290869061084d908790611564565b610974565b6060600480546105fa90611513565b7fa49807205ce4d355092ef5a8a18f56e8913cf4a201fbe287825b095693c2177561088c8133610acc565b610738838361109f565b3360008181526001602090815260408083206001600160a01b0387168452909152812054909190838110156109335760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f00000000000000000000000000000000000000000000000000000060648201526084016101dd565b6107568286868403610974565b60003361068b818585610d63565b60008281526005602052604090206001015461096a8133610acc565b610738838361101c565b6001600160a01b0383166109ef5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f726573730000000000000000000000000000000000000000000000000000000060648201526084016101dd565b6001600160a01b038216610a6b5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f737300000000000000000000000000000000000000000000000000000000000060648201526084016101dd565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff1661080f57610b0a816001600160a01b0316601461117e565b610b1583602061117e565b604051602001610b2692919061157c565b60408051601f198184030181529082905262461bcd60e51b82526101dd916004016113d4565b6001600160a01b038216610bc85760405162461bcd60e51b815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f730000000000000000000000000000000000000000000000000000000000000060648201526084016101dd565b6001600160a01b03821660009081526020819052604090205481811015610c575760405162461bcd60e51b815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f636500000000000000000000000000000000000000000000000000000000000060648201526084016101dd565b6001600160a01b0383166000908152602081905260408120838303905560028054849290610c869084906115fd565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b038381166000908152600160209081526040808320938616835292905220546000198114610d5d5781811015610d505760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e636500000060448201526064016101dd565b610d5d8484848403610974565b50505050565b6001600160a01b038316610ddf5760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f647265737300000000000000000000000000000000000000000000000000000060648201526084016101dd565b6001600160a01b038216610e5b5760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f657373000000000000000000000000000000000000000000000000000000000060648201526084016101dd565b6001600160a01b03831660009081526020819052604090205481811015610eea5760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e6365000000000000000000000000000000000000000000000000000060648201526084016101dd565b6001600160a01b03808516600090815260208190526040808220858503905591851681529081208054849290610f21908490611564565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610f6d91815260200190565b60405180910390a3610d5d565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff1661080f5760008281526005602090815260408083206001600160a01b03851684529091529020805460ff19166001179055610fd83390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff161561080f5760008281526005602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b6001600160a01b0382166110f55760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f20616464726573730060448201526064016101dd565b80600260008282546111079190611564565b90915550506001600160a01b03821660009081526020819052604081208054839290611134908490611564565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b6060600061118d836002611614565b611198906002611564565b67ffffffffffffffff8111156111b0576111b0611633565b6040519080825280601f01601f1916602001820160405280156111da576020820181803683370190505b5090507f30000000000000000000000000000000000000000000000000000000000000008160008151811061121157611211611649565b60200101906001600160f81b031916908160001a9053507f78000000000000000000000000000000000000000000000000000000000000008160018151811061125c5761125c611649565b60200101906001600160f81b031916908160001a9053506000611280846002611614565b61128b906001611564565b90505b6001811115611310577f303132333435363738396162636465660000000000000000000000000000000085600f16601081106112cc576112cc611649565b1a60f81b8282815181106112e2576112e2611649565b60200101906001600160f81b031916908160001a90535060049490941c936113098161165f565b905061128e565b50831561135f5760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016101dd565b9392505050565b60006020828403121561137857600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461135f57600080fd5b60005b838110156113c35781810151838201526020016113ab565b83811115610d5d5750506000910152565b60208152600082518060208401526113f38160408501602087016113a8565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461141e57600080fd5b919050565b6000806040838503121561143657600080fd5b61143f83611407565b946020939093013593505050565b60008060006060848603121561146257600080fd5b61146b84611407565b925061147960208501611407565b9150604084013590509250925092565b60006020828403121561149b57600080fd5b5035919050565b600080604083850312156114b557600080fd5b823591506114c560208401611407565b90509250929050565b6000602082840312156114e057600080fd5b61135f82611407565b600080604083850312156114fc57600080fd5b61150583611407565b91506114c560208401611407565b600181811c9082168061152757607f821691505b6020821081141561154857634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b600082198211156115775761157761154e565b500190565b7f416363657373436f6e74726f6c3a206163636f756e74200000000000000000008152600083516115b48160178501602088016113a8565b7f206973206d697373696e6720726f6c652000000000000000000000000000000060179184019182015283516115f18160288401602088016113a8565b01602801949350505050565b60008282101561160f5761160f61154e565b500390565b600081600019048311821515161561162e5761162e61154e565b500290565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b60008161166e5761166e61154e565b50600019019056fea2646970667358221220fdfb0e3eb91e955b3f1793bbc9277dbf6c071da0a046d2e22c55fd8ca54bccf864736f6c63430008090033",
}

// ObscuroERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ObscuroERC20MetaData.ABI instead.
var ObscuroERC20ABI = ObscuroERC20MetaData.ABI

// ObscuroERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObscuroERC20MetaData.Bin instead.
var ObscuroERC20Bin = ObscuroERC20MetaData.Bin

// DeployObscuroERC20 deploys a new Ethereum contract, binding an instance of ObscuroERC20 to it.
func DeployObscuroERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *ObscuroERC20, error) {
	parsed, err := ObscuroERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObscuroERC20Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObscuroERC20{ObscuroERC20Caller: ObscuroERC20Caller{contract: contract}, ObscuroERC20Transactor: ObscuroERC20Transactor{contract: contract}, ObscuroERC20Filterer: ObscuroERC20Filterer{contract: contract}}, nil
}

// ObscuroERC20 is an auto generated Go binding around an Ethereum contract.
type ObscuroERC20 struct {
	ObscuroERC20Caller     // Read-only binding to the contract
	ObscuroERC20Transactor // Write-only binding to the contract
	ObscuroERC20Filterer   // Log filterer for contract events
}

// ObscuroERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ObscuroERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ObscuroERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObscuroERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObscuroERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObscuroERC20Session struct {
	Contract     *ObscuroERC20     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObscuroERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObscuroERC20CallerSession struct {
	Contract *ObscuroERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ObscuroERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObscuroERC20TransactorSession struct {
	Contract     *ObscuroERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ObscuroERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ObscuroERC20Raw struct {
	Contract *ObscuroERC20 // Generic contract binding to access the raw methods on
}

// ObscuroERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObscuroERC20CallerRaw struct {
	Contract *ObscuroERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ObscuroERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObscuroERC20TransactorRaw struct {
	Contract *ObscuroERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewObscuroERC20 creates a new instance of ObscuroERC20, bound to a specific deployed contract.
func NewObscuroERC20(address common.Address, backend bind.ContractBackend) (*ObscuroERC20, error) {
	contract, err := bindObscuroERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20{ObscuroERC20Caller: ObscuroERC20Caller{contract: contract}, ObscuroERC20Transactor: ObscuroERC20Transactor{contract: contract}, ObscuroERC20Filterer: ObscuroERC20Filterer{contract: contract}}, nil
}

// NewObscuroERC20Caller creates a new read-only instance of ObscuroERC20, bound to a specific deployed contract.
func NewObscuroERC20Caller(address common.Address, caller bind.ContractCaller) (*ObscuroERC20Caller, error) {
	contract, err := bindObscuroERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20Caller{contract: contract}, nil
}

// NewObscuroERC20Transactor creates a new write-only instance of ObscuroERC20, bound to a specific deployed contract.
func NewObscuroERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ObscuroERC20Transactor, error) {
	contract, err := bindObscuroERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20Transactor{contract: contract}, nil
}

// NewObscuroERC20Filterer creates a new log filterer instance of ObscuroERC20, bound to a specific deployed contract.
func NewObscuroERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ObscuroERC20Filterer, error) {
	contract, err := bindObscuroERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20Filterer{contract: contract}, nil
}

// bindObscuroERC20 binds a generic wrapper to an already deployed contract.
func bindObscuroERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ObscuroERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroERC20 *ObscuroERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroERC20.Contract.ObscuroERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroERC20 *ObscuroERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.ObscuroERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroERC20 *ObscuroERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.ObscuroERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObscuroERC20 *ObscuroERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObscuroERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObscuroERC20 *ObscuroERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObscuroERC20 *ObscuroERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20Caller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20Session) ADMINROLE() ([32]byte, error) {
	return _ObscuroERC20.Contract.ADMINROLE(&_ObscuroERC20.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20CallerSession) ADMINROLE() ([32]byte, error) {
	return _ObscuroERC20.Contract.ADMINROLE(&_ObscuroERC20.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _ObscuroERC20.Contract.DEFAULTADMINROLE(&_ObscuroERC20.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ObscuroERC20.Contract.DEFAULTADMINROLE(&_ObscuroERC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObscuroERC20.Contract.Allowance(&_ObscuroERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObscuroERC20.Contract.Allowance(&_ObscuroERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObscuroERC20.Contract.BalanceOf(&_ObscuroERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObscuroERC20.Contract.BalanceOf(&_ObscuroERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObscuroERC20 *ObscuroERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObscuroERC20 *ObscuroERC20Session) Decimals() (uint8, error) {
	return _ObscuroERC20.Contract.Decimals(&_ObscuroERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObscuroERC20 *ObscuroERC20CallerSession) Decimals() (uint8, error) {
	return _ObscuroERC20.Contract.Decimals(&_ObscuroERC20.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ObscuroERC20.Contract.GetRoleAdmin(&_ObscuroERC20.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ObscuroERC20 *ObscuroERC20CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ObscuroERC20.Contract.GetRoleAdmin(&_ObscuroERC20.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroERC20 *ObscuroERC20Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroERC20 *ObscuroERC20Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ObscuroERC20.Contract.HasRole(&_ObscuroERC20.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ObscuroERC20 *ObscuroERC20CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ObscuroERC20.Contract.HasRole(&_ObscuroERC20.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObscuroERC20 *ObscuroERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObscuroERC20 *ObscuroERC20Session) Name() (string, error) {
	return _ObscuroERC20.Contract.Name(&_ObscuroERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObscuroERC20 *ObscuroERC20CallerSession) Name() (string, error) {
	return _ObscuroERC20.Contract.Name(&_ObscuroERC20.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroERC20 *ObscuroERC20Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroERC20 *ObscuroERC20Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ObscuroERC20.Contract.SupportsInterface(&_ObscuroERC20.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ObscuroERC20 *ObscuroERC20CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ObscuroERC20.Contract.SupportsInterface(&_ObscuroERC20.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObscuroERC20 *ObscuroERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObscuroERC20 *ObscuroERC20Session) Symbol() (string, error) {
	return _ObscuroERC20.Contract.Symbol(&_ObscuroERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObscuroERC20 *ObscuroERC20CallerSession) Symbol() (string, error) {
	return _ObscuroERC20.Contract.Symbol(&_ObscuroERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ObscuroERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20Session) TotalSupply() (*big.Int, error) {
	return _ObscuroERC20.Contract.TotalSupply(&_ObscuroERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObscuroERC20 *ObscuroERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ObscuroERC20.Contract.TotalSupply(&_ObscuroERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Approve(&_ObscuroERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Approve(&_ObscuroERC20.TransactOpts, spender, amount)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_ObscuroERC20 *ObscuroERC20Transactor) BurnFor(opts *bind.TransactOpts, giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "burnFor", giver, amount)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_ObscuroERC20 *ObscuroERC20Session) BurnFor(giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.BurnFor(&_ObscuroERC20.TransactOpts, giver, amount)
}

// BurnFor is a paid mutator transaction binding the contract method 0x1dd319cb.
//
// Solidity: function burnFor(address giver, uint256 amount) returns()
func (_ObscuroERC20 *ObscuroERC20TransactorSession) BurnFor(giver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.BurnFor(&_ObscuroERC20.TransactOpts, giver, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.DecreaseAllowance(&_ObscuroERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObscuroERC20 *ObscuroERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.DecreaseAllowance(&_ObscuroERC20.TransactOpts, spender, subtractedValue)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.GrantRole(&_ObscuroERC20.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.GrantRole(&_ObscuroERC20.TransactOpts, role, account)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.IncreaseAllowance(&_ObscuroERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObscuroERC20 *ObscuroERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.IncreaseAllowance(&_ObscuroERC20.TransactOpts, spender, addedValue)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_ObscuroERC20 *ObscuroERC20Transactor) IssueFor(opts *bind.TransactOpts, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "issueFor", receiver, amount)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_ObscuroERC20 *ObscuroERC20Session) IssueFor(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.IssueFor(&_ObscuroERC20.TransactOpts, receiver, amount)
}

// IssueFor is a paid mutator transaction binding the contract method 0x979005ad.
//
// Solidity: function issueFor(address receiver, uint256 amount) returns()
func (_ObscuroERC20 *ObscuroERC20TransactorSession) IssueFor(receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.IssueFor(&_ObscuroERC20.TransactOpts, receiver, amount)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20Session) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.RenounceRole(&_ObscuroERC20.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20TransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.RenounceRole(&_ObscuroERC20.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.RevokeRole(&_ObscuroERC20.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ObscuroERC20 *ObscuroERC20TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.RevokeRole(&_ObscuroERC20.TransactOpts, role, account)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Transfer(&_ObscuroERC20.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Transfer(&_ObscuroERC20.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.TransferFrom(&_ObscuroERC20.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObscuroERC20 *ObscuroERC20TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.TransferFrom(&_ObscuroERC20.TransactOpts, from, to, amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ObscuroERC20 *ObscuroERC20Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ObscuroERC20.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ObscuroERC20 *ObscuroERC20Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Fallback(&_ObscuroERC20.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ObscuroERC20 *ObscuroERC20TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Fallback(&_ObscuroERC20.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ObscuroERC20 *ObscuroERC20Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObscuroERC20.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ObscuroERC20 *ObscuroERC20Session) Receive() (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Receive(&_ObscuroERC20.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ObscuroERC20 *ObscuroERC20TransactorSession) Receive() (*types.Transaction, error) {
	return _ObscuroERC20.Contract.Receive(&_ObscuroERC20.TransactOpts)
}

// ObscuroERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ObscuroERC20 contract.
type ObscuroERC20ApprovalIterator struct {
	Event *ObscuroERC20Approval // Event containing the contract specifics and raw log

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
func (it *ObscuroERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroERC20Approval)
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
		it.Event = new(ObscuroERC20Approval)
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
func (it *ObscuroERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroERC20Approval represents a Approval event raised by the ObscuroERC20 contract.
type ObscuroERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObscuroERC20 *ObscuroERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ObscuroERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObscuroERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20ApprovalIterator{contract: _ObscuroERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObscuroERC20 *ObscuroERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ObscuroERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObscuroERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroERC20Approval)
				if err := _ObscuroERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObscuroERC20 *ObscuroERC20Filterer) ParseApproval(log types.Log) (*ObscuroERC20Approval, error) {
	event := new(ObscuroERC20Approval)
	if err := _ObscuroERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroERC20RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ObscuroERC20 contract.
type ObscuroERC20RoleAdminChangedIterator struct {
	Event *ObscuroERC20RoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ObscuroERC20RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroERC20RoleAdminChanged)
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
		it.Event = new(ObscuroERC20RoleAdminChanged)
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
func (it *ObscuroERC20RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroERC20RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroERC20RoleAdminChanged represents a RoleAdminChanged event raised by the ObscuroERC20 contract.
type ObscuroERC20RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ObscuroERC20 *ObscuroERC20Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ObscuroERC20RoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ObscuroERC20.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20RoleAdminChangedIterator{contract: _ObscuroERC20.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ObscuroERC20 *ObscuroERC20Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ObscuroERC20RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ObscuroERC20.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroERC20RoleAdminChanged)
				if err := _ObscuroERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ObscuroERC20 *ObscuroERC20Filterer) ParseRoleAdminChanged(log types.Log) (*ObscuroERC20RoleAdminChanged, error) {
	event := new(ObscuroERC20RoleAdminChanged)
	if err := _ObscuroERC20.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroERC20RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ObscuroERC20 contract.
type ObscuroERC20RoleGrantedIterator struct {
	Event *ObscuroERC20RoleGranted // Event containing the contract specifics and raw log

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
func (it *ObscuroERC20RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroERC20RoleGranted)
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
		it.Event = new(ObscuroERC20RoleGranted)
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
func (it *ObscuroERC20RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroERC20RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroERC20RoleGranted represents a RoleGranted event raised by the ObscuroERC20 contract.
type ObscuroERC20RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroERC20 *ObscuroERC20Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ObscuroERC20RoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ObscuroERC20.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20RoleGrantedIterator{contract: _ObscuroERC20.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroERC20 *ObscuroERC20Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ObscuroERC20RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ObscuroERC20.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroERC20RoleGranted)
				if err := _ObscuroERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroERC20 *ObscuroERC20Filterer) ParseRoleGranted(log types.Log) (*ObscuroERC20RoleGranted, error) {
	event := new(ObscuroERC20RoleGranted)
	if err := _ObscuroERC20.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroERC20RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ObscuroERC20 contract.
type ObscuroERC20RoleRevokedIterator struct {
	Event *ObscuroERC20RoleRevoked // Event containing the contract specifics and raw log

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
func (it *ObscuroERC20RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroERC20RoleRevoked)
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
		it.Event = new(ObscuroERC20RoleRevoked)
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
func (it *ObscuroERC20RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroERC20RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroERC20RoleRevoked represents a RoleRevoked event raised by the ObscuroERC20 contract.
type ObscuroERC20RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroERC20 *ObscuroERC20Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ObscuroERC20RoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ObscuroERC20.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20RoleRevokedIterator{contract: _ObscuroERC20.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroERC20 *ObscuroERC20Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ObscuroERC20RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ObscuroERC20.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroERC20RoleRevoked)
				if err := _ObscuroERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ObscuroERC20 *ObscuroERC20Filterer) ParseRoleRevoked(log types.Log) (*ObscuroERC20RoleRevoked, error) {
	event := new(ObscuroERC20RoleRevoked)
	if err := _ObscuroERC20.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObscuroERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ObscuroERC20 contract.
type ObscuroERC20TransferIterator struct {
	Event *ObscuroERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ObscuroERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObscuroERC20Transfer)
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
		it.Event = new(ObscuroERC20Transfer)
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
func (it *ObscuroERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObscuroERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObscuroERC20Transfer represents a Transfer event raised by the ObscuroERC20 contract.
type ObscuroERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObscuroERC20 *ObscuroERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ObscuroERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObscuroERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ObscuroERC20TransferIterator{contract: _ObscuroERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObscuroERC20 *ObscuroERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ObscuroERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObscuroERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObscuroERC20Transfer)
				if err := _ObscuroERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObscuroERC20 *ObscuroERC20Filterer) ParseTransfer(log types.Log) (*ObscuroERC20Transfer, error) {
	event := new(ObscuroERC20Transfer)
	if err := _ObscuroERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
