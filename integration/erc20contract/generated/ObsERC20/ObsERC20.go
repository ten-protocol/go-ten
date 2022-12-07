// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ObsERC20

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

// ObsERC20MetaData contains all meta data concerning the ObsERC20 contract.
var ObsERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"busAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405273deb34a740eca1ec42c8b8204cbec0ba34fdd27f3600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503480156200006657600080fd5b50604051620027213803806200272183398181016040528101906200008c9190620006a7565b838381600390816200009f919062000998565b508060049081620000b1919062000998565b505050620000c633836200011160201b60201c565b80600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050505062000de9565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000183576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200017a9062000ae0565b60405180910390fd5b62000197600083836200028960201b60201c565b8060026000828254620001ab919062000b31565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825462000202919062000b31565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405162000269919062000b7d565b60405180910390a362000285600083836200046f60201b60201c565b5050565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603156200046a57600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036200046957600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115620003d957620003d862000b9a565b5b85604051602001620003ec919062000c33565b60405160208183030381529060405260006040518563ffffffff1660e01b81526004016200041e949392919062000d1e565b6020604051808303816000875af11580156200043e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000464919062000db7565b905050505b5b505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620004dd8262000492565b810181811067ffffffffffffffff82111715620004ff57620004fe620004a3565b5b80604052505050565b60006200051462000474565b9050620005228282620004d2565b919050565b600067ffffffffffffffff821115620005455762000544620004a3565b5b620005508262000492565b9050602081019050919050565b60005b838110156200057d57808201518184015260208101905062000560565b60008484015250505050565b6000620005a06200059a8462000527565b62000508565b905082815260208101848484011115620005bf57620005be6200048d565b5b620005cc8482856200055d565b509392505050565b600082601f830112620005ec57620005eb62000488565b5b8151620005fe84826020860162000589565b91505092915050565b6000819050919050565b6200061c8162000607565b81146200062857600080fd5b50565b6000815190506200063c8162000611565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200066f8262000642565b9050919050565b620006818162000662565b81146200068d57600080fd5b50565b600081519050620006a18162000676565b92915050565b60008060008060808587031215620006c457620006c36200047e565b5b600085015167ffffffffffffffff811115620006e557620006e462000483565b5b620006f387828801620005d4565b945050602085015167ffffffffffffffff81111562000717576200071662000483565b5b6200072587828801620005d4565b935050604062000738878288016200062b565b92505060606200074b8782880162000690565b91505092959194509250565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620007aa57607f821691505b602082108103620007c057620007bf62000762565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026200082a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620007eb565b620008368683620007eb565b95508019841693508086168417925050509392505050565b6000819050919050565b600062000879620008736200086d8462000607565b6200084e565b62000607565b9050919050565b6000819050919050565b620008958362000858565b620008ad620008a48262000880565b848454620007f8565b825550505050565b600090565b620008c4620008b5565b620008d18184846200088a565b505050565b5b81811015620008f957620008ed600082620008ba565b600181019050620008d7565b5050565b601f82111562000948576200091281620007c6565b6200091d84620007db565b810160208510156200092d578190505b620009456200093c85620007db565b830182620008d6565b50505b505050565b600082821c905092915050565b60006200096d600019846008026200094d565b1980831691505092915050565b60006200098883836200095a565b9150826002028217905092915050565b620009a38262000757565b67ffffffffffffffff811115620009bf57620009be620004a3565b5b620009cb825462000791565b620009d8828285620008fd565b600060209050601f83116001811462000a105760008415620009fb578287015190505b62000a0785826200097a565b86555062000a77565b601f19841662000a2086620007c6565b60005b8281101562000a4a5784890151825560018201915060208501945060208101905062000a23565b8683101562000a6a578489015162000a66601f8916826200095a565b8355505b6001600288020188555050505b505050505050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000ac8601f8362000a7f565b915062000ad58262000a90565b602082019050919050565b6000602082019050818103600083015262000afb8162000ab9565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600062000b3e8262000607565b915062000b4b8362000607565b925082820190508082111562000b665762000b6562000b02565b5b92915050565b62000b778162000607565b82525050565b600060208201905062000b94600083018462000b6c565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b62000bd48162000662565b82525050565b62000be58162000607565b82525050565b60608201600082015162000c03600085018262000bc9565b50602082015162000c18602085018262000bc9565b50604082015162000c2d604085018262000bda565b50505050565b600060608201905062000c4a600083018462000beb565b92915050565b600063ffffffff82169050919050565b62000c6b8162000c50565b82525050565b600081519050919050565b600082825260208201905092915050565b600062000c9a8262000c71565b62000ca6818562000c7c565b935062000cb88185602086016200055d565b62000cc38162000492565b840191505092915050565b6000819050919050565b600060ff82169050919050565b600062000d0662000d0062000cfa8462000cce565b6200084e565b62000cd8565b9050919050565b62000d188162000ce5565b82525050565b600060808201905062000d35600083018762000c60565b62000d44602083018662000c60565b818103604083015262000d58818562000c8d565b905062000d69606083018462000d0d565b95945050505050565b600067ffffffffffffffff82169050919050565b62000d918162000d72565b811462000d9d57600080fd5b50565b60008151905062000db18162000d86565b92915050565b60006020828403121562000dd05762000dcf6200047e565b5b600062000de08482850162000da0565b91505092915050565b6119288062000df96000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610eeb565b60405180910390f35b6100e660048036038101906100e19190610fa6565b610308565b6040516100f39190611001565b60405180910390f35b61010461032b565b604051610111919061102b565b60405180910390f35b610134600480360381019061012f9190611046565b610335565b6040516101419190611001565b60405180910390f35b610152610364565b60405161015f91906110b5565b60405180910390f35b610182600480360381019061017d9190610fa6565b61036d565b60405161018f9190611001565b60405180910390f35b6101b260048036038101906101ad91906110d0565b6103a4565b6040516101bf919061102b565b60405180910390f35b6101d061046c565b6040516101dd9190610eeb565b60405180910390f35b61020060048036038101906101fb9190610fa6565b6104fe565b60405161020d9190611001565b60405180910390f35b610230600480360381019061022b9190610fa6565b610575565b60405161023d9190611001565b60405180910390f35b610260600480360381019061025b91906110fd565b610598565b60405161026d919061102b565b60405180910390f35b6060600380546102859061116c565b80601f01602080910402602001604051908101604052809291908181526020018280546102b19061116c565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b6000806103136106d1565b90506103208185856106d9565b600191505092915050565b6000600254905090565b6000806103406106d1565b905061034d8582856108a2565b61035885858561092e565b60019150509392505050565b60006012905090565b6000806103786106d1565b905061039981858561038a8589610598565b61039491906111cc565b6106d9565b600191505092915050565b60008173ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff16036103e9576103e282610bad565b9050610467565b8173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff160361042c5761042582610bad565b9050610467565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161045e9061124c565b60405180910390fd5b919050565b60606004805461047b9061116c565b80601f01602080910402602001604051908101604052809291908181526020018280546104a79061116c565b80156104f45780601f106104c9576101008083540402835291602001916104f4565b820191906000526020600020905b8154815290600101906020018083116104d757829003601f168201915b5050505050905090565b6000806105096106d1565b905060006105178286610598565b90508381101561055c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610553906112de565b60405180910390fd5b61056982868684036106d9565b60019250505092915050565b6000806105806106d1565b905061058d81858561092e565b600191505092915050565b60008273ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff1614806105ff57508173ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff16145b156106155761060e8383610bf5565b90506106cb565b8273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061067a57508173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b15610690576106898383610bf5565b90506106cb565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106c290611370565b60405180910390fd5b92915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610748576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073f90611402565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036107b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107ae90611494565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610895919061102b565b60405180910390a3505050565b60006108ae8484610598565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610928578181101561091a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161091190611500565b60405180910390fd5b61092784848484036106d9565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361099d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161099490611592565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610a0c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a0390611624565b60405180910390fd5b610a17838383610c7c565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610a9d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a94906116b6565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610b3091906111cc565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610b94919061102b565b60405180910390a3610ba7848484610e56565b50505050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160315610e5157600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610e5057600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115610dc757610dc66116d6565b5b85604051602001610dd89190611765565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401610e089493929190611839565b6020604051808303816000875af1158015610e27573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e4b91906118c5565b905050505b5b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610e95578082015181840152602081019050610e7a565b60008484015250505050565b6000601f19601f8301169050919050565b6000610ebd82610e5b565b610ec78185610e66565b9350610ed7818560208601610e77565b610ee081610ea1565b840191505092915050565b60006020820190508181036000830152610f058184610eb2565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610f3d82610f12565b9050919050565b610f4d81610f32565b8114610f5857600080fd5b50565b600081359050610f6a81610f44565b92915050565b6000819050919050565b610f8381610f70565b8114610f8e57600080fd5b50565b600081359050610fa081610f7a565b92915050565b60008060408385031215610fbd57610fbc610f0d565b5b6000610fcb85828601610f5b565b9250506020610fdc85828601610f91565b9150509250929050565b60008115159050919050565b610ffb81610fe6565b82525050565b60006020820190506110166000830184610ff2565b92915050565b61102581610f70565b82525050565b6000602082019050611040600083018461101c565b92915050565b60008060006060848603121561105f5761105e610f0d565b5b600061106d86828701610f5b565b935050602061107e86828701610f5b565b925050604061108f86828701610f91565b9150509250925092565b600060ff82169050919050565b6110af81611099565b82525050565b60006020820190506110ca60008301846110a6565b92915050565b6000602082840312156110e6576110e5610f0d565b5b60006110f484828501610f5b565b91505092915050565b6000806040838503121561111457611113610f0d565b5b600061112285828601610f5b565b925050602061113385828601610f5b565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061118457607f821691505b6020821081036111975761119661113d565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006111d782610f70565b91506111e283610f70565b92508282019050808211156111fa576111f961119d565b5b92915050565b7f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e636500600082015250565b6000611236601f83610e66565b915061124182611200565b602082019050919050565b6000602082019050818103600083015261126581611229565b9050919050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b60006112c8602583610e66565b91506112d38261126c565b604082019050919050565b600060208201905081810360008301526112f7816112bb565b9050919050565b7f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360008201527f6500000000000000000000000000000000000000000000000000000000000000602082015250565b600061135a602183610e66565b9150611365826112fe565b604082019050919050565b600060208201905081810360008301526113898161134d565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006113ec602483610e66565b91506113f782611390565b604082019050919050565b6000602082019050818103600083015261141b816113df565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b600061147e602283610e66565b915061148982611422565b604082019050919050565b600060208201905081810360008301526114ad81611471565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b60006114ea601d83610e66565b91506114f5826114b4565b602082019050919050565b60006020820190508181036000830152611519816114dd565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b600061157c602583610e66565b915061158782611520565b604082019050919050565b600060208201905081810360008301526115ab8161156f565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b600061160e602383610e66565b9150611619826115b2565b604082019050919050565b6000602082019050818103600083015261163d81611601565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b60006116a0602683610e66565b91506116ab82611644565b604082019050919050565b600060208201905081810360008301526116cf81611693565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b61170e81610f32565b82525050565b61171d81610f70565b82525050565b6060820160008201516117396000850182611705565b50602082015161174c6020850182611705565b50604082015161175f6040850182611714565b50505050565b600060608201905061177a6000830184611723565b92915050565b600063ffffffff82169050919050565b61179981611780565b82525050565b600081519050919050565b600082825260208201905092915050565b60006117c68261179f565b6117d081856117aa565b93506117e0818560208601610e77565b6117e981610ea1565b840191505092915050565b6000819050919050565b6000819050919050565b600061182361181e611819846117f4565b6117fe565b611099565b9050919050565b61183381611808565b82525050565b600060808201905061184e6000830187611790565b61185b6020830186611790565b818103604083015261186d81856117bb565b905061187c606083018461182a565b95945050505050565b600067ffffffffffffffff82169050919050565b6118a281611885565b81146118ad57600080fd5b50565b6000815190506118bf81611899565b92915050565b6000602082840312156118db576118da610f0d565b5b60006118e9848285016118b0565b9150509291505056fea264697066735822122070ad4a48d579576c33a29850a9848b40de3c6338ef1a92a842bff6fbaacae06764736f6c63430008110033",
}

// ObsERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use ObsERC20MetaData.ABI instead.
var ObsERC20ABI = ObsERC20MetaData.ABI

// ObsERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ObsERC20MetaData.Bin instead.
var ObsERC20Bin = ObsERC20MetaData.Bin

// DeployObsERC20 deploys a new Ethereum contract, binding an instance of ObsERC20 to it.
func DeployObsERC20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int, busAddress common.Address) (common.Address, *types.Transaction, *ObsERC20, error) {
	parsed, err := ObsERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ObsERC20Bin), backend, name, symbol, initialSupply, busAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ObsERC20{ObsERC20Caller: ObsERC20Caller{contract: contract}, ObsERC20Transactor: ObsERC20Transactor{contract: contract}, ObsERC20Filterer: ObsERC20Filterer{contract: contract}}, nil
}

// ObsERC20 is an auto generated Go binding around an Ethereum contract.
type ObsERC20 struct {
	ObsERC20Caller     // Read-only binding to the contract
	ObsERC20Transactor // Write-only binding to the contract
	ObsERC20Filterer   // Log filterer for contract events
}

// ObsERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ObsERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ObsERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ObsERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ObsERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ObsERC20Session struct {
	Contract     *ObsERC20         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ObsERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ObsERC20CallerSession struct {
	Contract *ObsERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ObsERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ObsERC20TransactorSession struct {
	Contract     *ObsERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ObsERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ObsERC20Raw struct {
	Contract *ObsERC20 // Generic contract binding to access the raw methods on
}

// ObsERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ObsERC20CallerRaw struct {
	Contract *ObsERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ObsERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ObsERC20TransactorRaw struct {
	Contract *ObsERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewObsERC20 creates a new instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20(address common.Address, backend bind.ContractBackend) (*ObsERC20, error) {
	contract, err := bindObsERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ObsERC20{ObsERC20Caller: ObsERC20Caller{contract: contract}, ObsERC20Transactor: ObsERC20Transactor{contract: contract}, ObsERC20Filterer: ObsERC20Filterer{contract: contract}}, nil
}

// NewObsERC20Caller creates a new read-only instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Caller(address common.Address, caller bind.ContractCaller) (*ObsERC20Caller, error) {
	contract, err := bindObsERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Caller{contract: contract}, nil
}

// NewObsERC20Transactor creates a new write-only instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ObsERC20Transactor, error) {
	contract, err := bindObsERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Transactor{contract: contract}, nil
}

// NewObsERC20Filterer creates a new log filterer instance of ObsERC20, bound to a specific deployed contract.
func NewObsERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ObsERC20Filterer, error) {
	contract, err := bindObsERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ObsERC20Filterer{contract: contract}, nil
}

// bindObsERC20 binds a generic wrapper to an already deployed contract.
func bindObsERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ObsERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObsERC20 *ObsERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObsERC20.Contract.ObsERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObsERC20 *ObsERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObsERC20.Contract.ObsERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObsERC20 *ObsERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObsERC20.Contract.ObsERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ObsERC20 *ObsERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ObsERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ObsERC20 *ObsERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ObsERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ObsERC20 *ObsERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ObsERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.Allowance(&_ObsERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.Allowance(&_ObsERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.BalanceOf(&_ObsERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ObsERC20.Contract.BalanceOf(&_ObsERC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20Session) Decimals() (uint8, error) {
	return _ObsERC20.Contract.Decimals(&_ObsERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ObsERC20 *ObsERC20CallerSession) Decimals() (uint8, error) {
	return _ObsERC20.Contract.Decimals(&_ObsERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20Session) Name() (string, error) {
	return _ObsERC20.Contract.Name(&_ObsERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ObsERC20 *ObsERC20CallerSession) Name() (string, error) {
	return _ObsERC20.Contract.Name(&_ObsERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20Session) Symbol() (string, error) {
	return _ObsERC20.Contract.Symbol(&_ObsERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ObsERC20 *ObsERC20CallerSession) Symbol() (string, error) {
	return _ObsERC20.Contract.Symbol(&_ObsERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20Session) TotalSupply() (*big.Int, error) {
	return _ObsERC20.Contract.TotalSupply(&_ObsERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ObsERC20 *ObsERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ObsERC20.Contract.TotalSupply(&_ObsERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Approve(&_ObsERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Approve(&_ObsERC20.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObsERC20 *ObsERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.DecreaseAllowance(&_ObsERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.DecreaseAllowance(&_ObsERC20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObsERC20 *ObsERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.IncreaseAllowance(&_ObsERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.IncreaseAllowance(&_ObsERC20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Session) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Transfer(&_ObsERC20.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.Transfer(&_ObsERC20.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20Session) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.TransferFrom(&_ObsERC20.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ObsERC20 *ObsERC20TransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ObsERC20.Contract.TransferFrom(&_ObsERC20.TransactOpts, from, to, amount)
}

// ObsERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ObsERC20 contract.
type ObsERC20ApprovalIterator struct {
	Event *ObsERC20Approval // Event containing the contract specifics and raw log

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
func (it *ObsERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObsERC20Approval)
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
		it.Event = new(ObsERC20Approval)
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
func (it *ObsERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObsERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObsERC20Approval represents a Approval event raised by the ObsERC20 contract.
type ObsERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ObsERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObsERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ObsERC20ApprovalIterator{contract: _ObsERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ObsERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ObsERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObsERC20Approval)
				if err := _ObsERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ObsERC20 *ObsERC20Filterer) ParseApproval(log types.Log) (*ObsERC20Approval, error) {
	event := new(ObsERC20Approval)
	if err := _ObsERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ObsERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ObsERC20 contract.
type ObsERC20TransferIterator struct {
	Event *ObsERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ObsERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ObsERC20Transfer)
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
		it.Event = new(ObsERC20Transfer)
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
func (it *ObsERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ObsERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ObsERC20Transfer represents a Transfer event raised by the ObsERC20 contract.
type ObsERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ObsERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObsERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ObsERC20TransferIterator{contract: _ObsERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ObsERC20 *ObsERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ObsERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ObsERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ObsERC20Transfer)
				if err := _ObsERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ObsERC20 *ObsERC20Filterer) ParseTransfer(log types.Log) (*ObsERC20Transfer, error) {
	event := new(ObsERC20Transfer)
	if err := _ObsERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
