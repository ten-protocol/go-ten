// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generatedManagementContract

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

// GeneratedManagementContractMetaData contains all meta data concerning the GeneratedManagementContract contract.
var GeneratedManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"verifyAttester\",\"type\":\"bool\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061266a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063d4c806641161005b578063d4c8066414610113578063e0643dfc14610143578063e0fd84bd14610176578063e34fbfc81461019257610088565b8063324ff8661461008d57806368e10383146100ab5780638ef74f89146100c7578063bbd79e15146100f7575b600080fd5b6100956101ae565b6040516100a291906110c7565b60405180910390f35b6100c560048036038101906100c091906112f0565b610287565b005b6100e160048036038101906100dc9190611380565b61034f565b6040516100ee91906113f7565b60405180910390f35b610111600480360381019061010c91906114f2565b6103ef565b005b61012d60048036038101906101289190611380565b6105cc565b60405161013a91906115e2565b60405180910390f35b61015d60048036038101906101589190611633565b6105ec565b60405161016d94939291906116aa565b60405180910390f35b610190600480360381019061018b9190611771565b610659565b005b6101ac60048036038101906101a7919061180b565b610796565b005b60606003805480602002602001604051908101604052809291908181526020016000905b8282101561027e5783829060005260206000200180546101f190611887565b80601f016020809104026020016040519081016040528092919081815260200182805461021d90611887565b801561026a5780601f1061023f5761010080835404028352916020019161026a565b820191906000526020600020905b81548152906001019060200180831161024d57829003601f168201915b5050505050815260200190600101906101d2565b50505050905090565b600460009054906101000a900460ff16156102a157600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506003819080600181540180825580915050600190039060005260206000200160009091909190915090816103489190611a64565b5050505050565b6001602052806000526040600020600091509050805461036e90611887565b80601f016020809104026020016040519081016040528092919081815260200182805461039a90611887565b80156103e75780601f106103bc576101008083540402835291602001916103e7565b820191906000526020600020905b8154815290600101906020018083116103ca57829003601f168201915b505050505081565b6000600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061044a57600080fd5b81156105365760006104808888878760405160200161046c9493929190611c01565b6040516020818303038152906040526107e9565b9050600061048e8288610824565b90508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146104c88a61084b565b6104d18361084b565b6040516020016104e2929190611d9d565b60405160208183030381529060405290610532576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161052991906113f7565b60405180910390fd5b5050505b6001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506003839080600181540180825580915050600190039060005260206000200160009091909190915090816105c29190611a64565b5050505050505050565b60026020528060005260406000206000915054906101000a900460ff1681565b6000602052816000526040600020818154811061060857600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166106af57600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826107e4929190611ded565b505050565b60006107f58251610a0e565b82604051602001610807929190611f09565b604051602081830303815290604052805190602001209050919050565b60008060006108338585610b6e565b9150915061084081610bef565b819250505092915050565b60606000602867ffffffffffffffff81111561086a576108696111c5565b5b6040519080825280601f01601f19166020018201604052801561089c5781602001600182028036833780820191505090505b50905060005b6014811015610a045760008160136108ba9190611f67565b60086108c69190611f9b565b60026108d29190612128565b8573ffffffffffffffffffffffffffffffffffffffff166108f391906121a2565b60f81b9050600060108260f81c61090a91906121e0565b60f81b905060008160f81c60106109219190612211565b8360f81c61092f919061224c565b60f81b905061093d82610dbb565b8585600261094b9190611f9b565b8151811061095c5761095b612280565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061099481610dbb565b8560018660026109a49190611f9b565b6109ae91906122af565b815181106109bf576109be612280565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050505080806109fc90612305565b9150506108a2565b5080915050919050565b606060008203610a55576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610b69565b600082905060005b60008214610a87578080610a7090612305565b915050600a82610a8091906121a2565b9150610a5d565b60008167ffffffffffffffff811115610aa357610aa26111c5565b5b6040519080825280601f01601f191660200182016040528015610ad55781602001600182028036833780820191505090505b5090505b60008514610b6257600182610aee9190611f67565b9150600a85610afd919061234d565b6030610b0991906122af565b60f81b818381518110610b1f57610b1e612280565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610b5b91906121a2565b9450610ad9565b8093505050505b919050565b6000806041835103610baf5760008060006020860151925060408601519150606086015160001a9050610ba387828585610e01565b94509450505050610be8565b6040835103610bdf576000806020850151915060408501519050610bd4868383610f0d565b935093505050610be8565b60006002915091505b9250929050565b60006004811115610c0357610c0261237e565b5b816004811115610c1657610c1561237e565b5b0315610db85760016004811115610c3057610c2f61237e565b5b816004811115610c4357610c4261237e565b5b03610c83576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7a906123f9565b60405180910390fd5b60026004811115610c9757610c9661237e565b5b816004811115610caa57610ca961237e565b5b03610cea576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ce190612465565b60405180910390fd5b60036004811115610cfe57610cfd61237e565b5b816004811115610d1157610d1061237e565b5b03610d51576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d48906124f7565b60405180910390fd5b600480811115610d6457610d6361237e565b5b816004811115610d7757610d7661237e565b5b03610db7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dae90612589565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610de65760308260f81c610ddc91906125a9565b60f81b9050610dfc565b60578260f81c610df691906125a9565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610e3c576000600391509150610f04565b601b8560ff1614158015610e545750601c8560ff1614155b15610e66576000600491509150610f04565b600060018787878760405160008152602001604052604051610e8b94939291906125ef565b6020604051602081039080840390855afa158015610ead573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610efb57600060019250925050610f04565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c610f5091906122af565b9050610f5e87828885610e01565b935093505050935093915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610fd2578082015181840152602081019050610fb7565b83811115610fe1576000848401525b50505050565b6000601f19601f8301169050919050565b600061100382610f98565b61100d8185610fa3565b935061101d818560208601610fb4565b61102681610fe7565b840191505092915050565b600061103d8383610ff8565b905092915050565b6000602082019050919050565b600061105d82610f6c565b6110678185610f77565b93508360208202850161107985610f88565b8060005b858110156110b557848403895281516110968582611031565b94506110a183611045565b925060208a0199505060018101905061107d565b50829750879550505050505092915050565b600060208201905081810360008301526110e18184611052565b905092915050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000611128826110fd565b9050919050565b6111388161111d565b811461114357600080fd5b50565b6000813590506111558161112f565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126111805761117f61115b565b5b8235905067ffffffffffffffff81111561119d5761119c611160565b5b6020830191508360018202830111156111b9576111b8611165565b5b9250929050565b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6111fd82610fe7565b810181811067ffffffffffffffff8211171561121c5761121b6111c5565b5b80604052505050565b600061122f6110e9565b905061123b82826111f4565b919050565b600067ffffffffffffffff82111561125b5761125a6111c5565b5b61126482610fe7565b9050602081019050919050565b82818337600083830152505050565b600061129361128e84611240565b611225565b9050828152602081018484840111156112af576112ae6111c0565b5b6112ba848285611271565b509392505050565b600082601f8301126112d7576112d661115b565b5b81356112e7848260208601611280565b91505092915050565b6000806000806060858703121561130a576113096110f3565b5b600061131887828801611146565b945050602085013567ffffffffffffffff811115611339576113386110f8565b5b6113458782880161116a565b9350935050604085013567ffffffffffffffff811115611368576113676110f8565b5b611374878288016112c2565b91505092959194509250565b600060208284031215611396576113956110f3565b5b60006113a484828501611146565b91505092915050565b600082825260208201905092915050565b60006113c982610f98565b6113d381856113ad565b93506113e3818560208601610fb4565b6113ec81610fe7565b840191505092915050565b6000602082019050818103600083015261141181846113be565b905092915050565b600067ffffffffffffffff821115611434576114336111c5565b5b61143d82610fe7565b9050602081019050919050565b600061145d61145884611419565b611225565b905082815260208101848484011115611479576114786111c0565b5b611484848285611271565b509392505050565b600082601f8301126114a1576114a061115b565b5b81356114b184826020860161144a565b91505092915050565b60008115159050919050565b6114cf816114ba565b81146114da57600080fd5b50565b6000813590506114ec816114c6565b92915050565b60008060008060008060c0878903121561150f5761150e6110f3565b5b600061151d89828a01611146565b965050602061152e89828a01611146565b955050604087013567ffffffffffffffff81111561154f5761154e6110f8565b5b61155b89828a0161148c565b945050606087013567ffffffffffffffff81111561157c5761157b6110f8565b5b61158889828a0161148c565b935050608087013567ffffffffffffffff8111156115a9576115a86110f8565b5b6115b589828a016112c2565b92505060a06115c689828a016114dd565b9150509295509295509295565b6115dc816114ba565b82525050565b60006020820190506115f760008301846115d3565b92915050565b6000819050919050565b611610816115fd565b811461161b57600080fd5b50565b60008135905061162d81611607565b92915050565b6000806040838503121561164a576116496110f3565b5b60006116588582860161161e565b92505060206116698582860161161e565b9150509250929050565b6000819050919050565b61168681611673565b82525050565b6116958161111d565b82525050565b6116a4816115fd565b82525050565b60006080820190506116bf600083018761167d565b6116cc602083018661168c565b6116d9604083018561167d565b6116e6606083018461169b565b95945050505050565b6116f881611673565b811461170357600080fd5b50565b600081359050611715816116ef565b92915050565b60008083601f8401126117315761173061115b565b5b8235905067ffffffffffffffff81111561174e5761174d611160565b5b60208301915083600182028301111561176a57611769611165565b5b9250929050565b60008060008060008060a0878903121561178e5761178d6110f3565b5b600061179c89828a01611706565b96505060206117ad89828a01611146565b95505060406117be89828a01611706565b94505060606117cf89828a0161161e565b935050608087013567ffffffffffffffff8111156117f0576117ef6110f8565b5b6117fc89828a0161171b565b92509250509295509295509295565b60008060208385031215611822576118216110f3565b5b600083013567ffffffffffffffff8111156118405761183f6110f8565b5b61184c8582860161171b565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061189f57607f821691505b6020821081036118b2576118b1611858565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261191a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826118dd565b61192486836118dd565b95508019841693508086168417925050509392505050565b6000819050919050565b600061196161195c611957846115fd565b61193c565b6115fd565b9050919050565b6000819050919050565b61197b83611946565b61198f61198782611968565b8484546118ea565b825550505050565b600090565b6119a4611997565b6119af818484611972565b505050565b5b818110156119d3576119c860008261199c565b6001810190506119b5565b5050565b601f821115611a18576119e9816118b8565b6119f2846118cd565b81016020851015611a01578190505b611a15611a0d856118cd565b8301826119b4565b50505b505050565b600082821c905092915050565b6000611a3b60001984600802611a1d565b1980831691505092915050565b6000611a548383611a2a565b9150826002028217905092915050565b611a6d82610f98565b67ffffffffffffffff811115611a8657611a856111c5565b5b611a908254611887565b611a9b8282856119d7565b600060209050601f831160018114611ace5760008415611abc578287015190505b611ac68582611a48565b865550611b2e565b601f198416611adc866118b8565b60005b82811015611b0457848901518255600182019150602085019450602081019050611adf565b86831015611b215784890151611b1d601f891682611a2a565b8355505b6001600288020188555050505b505050505050565b60008160601b9050919050565b6000611b4e82611b36565b9050919050565b6000611b6082611b43565b9050919050565b611b78611b738261111d565b611b55565b82525050565b600081519050919050565b600081905092915050565b6000611b9f82611b7e565b611ba98185611b89565b9350611bb9818560208601610fb4565b80840191505092915050565b600081905092915050565b6000611bdb82610f98565b611be58185611bc5565b9350611bf5818560208601610fb4565b80840191505092915050565b6000611c0d8287611b67565b601482019150611c1d8286611b67565b601482019150611c2d8285611b94565b9150611c398284611bd0565b915081905095945050505050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e2774206d617463682000000000000000000000000000000000000000602082015250565b6000611ca3602d83611bc5565b9150611cae82611c47565b602d82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b6000611d15602483611bc5565b9150611d2082611cb9565b602482019050919050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b6000611d87602483611bc5565b9150611d9282611d2b565b602482019050919050565b6000611da882611c96565b9150611db382611d08565b9150611dbf8285611bd0565b9150611dca82611d7a565b9150611dd68284611bd0565b91508190509392505050565b600082905092915050565b611df78383611de2565b67ffffffffffffffff811115611e1057611e0f6111c5565b5b611e1a8254611887565b611e258282856119d7565b6000601f831160018114611e545760008415611e42578287013590505b611e4c8582611a48565b865550611eb4565b601f198416611e62866118b8565b60005b82811015611e8a57848901358255600182019150602085019450602081019050611e65565b86831015611ea75784890135611ea3601f891682611a2a565b8355505b6001600288020188555050505b50505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611ef3601a83611bc5565b9150611efe82611ebd565b601a82019050919050565b6000611f1482611ee6565b9150611f208285611bd0565b9150611f2c8284611b94565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611f72826115fd565b9150611f7d836115fd565b925082821015611f9057611f8f611f38565b5b828203905092915050565b6000611fa6826115fd565b9150611fb1836115fd565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611fea57611fe9611f38565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b600185111561204c5780860481111561202857612027611f38565b5b60018516156120375780820291505b808102905061204585611ff5565b945061200c565b94509492505050565b6000826120655760019050612121565b816120735760009050612121565b81600181146120895760028114612093576120c2565b6001915050612121565b60ff8411156120a5576120a4611f38565b5b8360020a9150848211156120bc576120bb611f38565b5b50612121565b5060208310610133831016604e8410600b84101617156120f75782820a9050838111156120f2576120f1611f38565b5b612121565b6121048484846001612002565b9250905081840481111561211b5761211a611f38565b5b81810290505b9392505050565b6000612133826115fd565b915061213e836115fd565b925061216b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8484612055565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006121ad826115fd565b91506121b8836115fd565b9250826121c8576121c7612173565b5b828204905092915050565b600060ff82169050919050565b60006121eb826121d3565b91506121f6836121d3565b92508261220657612205612173565b5b828204905092915050565b600061221c826121d3565b9150612227836121d3565b92508160ff048311821515161561224157612240611f38565b5b828202905092915050565b6000612257826121d3565b9150612262836121d3565b92508282101561227557612274611f38565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006122ba826115fd565b91506122c5836115fd565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156122fa576122f9611f38565b5b828201905092915050565b6000612310826115fd565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361234257612341611f38565b5b600182019050919050565b6000612358826115fd565b9150612363836115fd565b92508261237357612372612173565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b60006123e36018836113ad565b91506123ee826123ad565b602082019050919050565b60006020820190508181036000830152612412816123d6565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b600061244f601f836113ad565b915061245a82612419565b602082019050919050565b6000602082019050818103600083015261247e81612442565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006124e16022836113ad565b91506124ec82612485565b604082019050919050565b60006020820190508181036000830152612510816124d4565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006125736022836113ad565b915061257e82612517565b604082019050919050565b600060208201905081810360008301526125a281612566565b9050919050565b60006125b4826121d3565b91506125bf836121d3565b92508260ff038211156125d5576125d4611f38565b5b828201905092915050565b6125e9816121d3565b82525050565b6000608082019050612604600083018761167d565b61261160208301866125e0565b61261e604083018561167d565b61262b606083018461167d565b9594505050505056fea2646970667358221220f8145af01251715d41b9b506fc93c23274e7e0ee9333ed8893db743bf42ecde564736f6c634300080f0033",
}

// GeneratedManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use GeneratedManagementContractMetaData.ABI instead.
var GeneratedManagementContractABI = GeneratedManagementContractMetaData.ABI

// GeneratedManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GeneratedManagementContractMetaData.Bin instead.
var GeneratedManagementContractBin = GeneratedManagementContractMetaData.Bin

// DeployGeneratedManagementContract deploys a new Ethereum contract, binding an instance of GeneratedManagementContract to it.
func DeployGeneratedManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GeneratedManagementContract, error) {
	parsed, err := GeneratedManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GeneratedManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// GeneratedManagementContract is an auto generated Go binding around an Ethereum contract.
type GeneratedManagementContract struct {
	GeneratedManagementContractCaller     // Read-only binding to the contract
	GeneratedManagementContractTransactor // Write-only binding to the contract
	GeneratedManagementContractFilterer   // Log filterer for contract events
}

// GeneratedManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GeneratedManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeneratedManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GeneratedManagementContractSession struct {
	Contract     *GeneratedManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// GeneratedManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GeneratedManagementContractCallerSession struct {
	Contract *GeneratedManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// GeneratedManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GeneratedManagementContractTransactorSession struct {
	Contract     *GeneratedManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// GeneratedManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GeneratedManagementContractRaw struct {
	Contract *GeneratedManagementContract // Generic contract binding to access the raw methods on
}

// GeneratedManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GeneratedManagementContractCallerRaw struct {
	Contract *GeneratedManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// GeneratedManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GeneratedManagementContractTransactorRaw struct {
	Contract *GeneratedManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGeneratedManagementContract creates a new instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContract(address common.Address, backend bind.ContractBackend) (*GeneratedManagementContract, error) {
	contract, err := bindGeneratedManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContract{GeneratedManagementContractCaller: GeneratedManagementContractCaller{contract: contract}, GeneratedManagementContractTransactor: GeneratedManagementContractTransactor{contract: contract}, GeneratedManagementContractFilterer: GeneratedManagementContractFilterer{contract: contract}}, nil
}

// NewGeneratedManagementContractCaller creates a new read-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractCaller(address common.Address, caller bind.ContractCaller) (*GeneratedManagementContractCaller, error) {
	contract, err := bindGeneratedManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractCaller{contract: contract}, nil
}

// NewGeneratedManagementContractTransactor creates a new write-only instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*GeneratedManagementContractTransactor, error) {
	contract, err := bindGeneratedManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractTransactor{contract: contract}, nil
}

// NewGeneratedManagementContractFilterer creates a new log filterer instance of GeneratedManagementContract, bound to a specific deployed contract.
func NewGeneratedManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*GeneratedManagementContractFilterer, error) {
	contract, err := bindGeneratedManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GeneratedManagementContractFilterer{contract: contract}, nil
}

// bindGeneratedManagementContract binds a generic wrapper to an already deployed contract.
func bindGeneratedManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GeneratedManagementContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.GeneratedManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GeneratedManagementContract *GeneratedManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GeneratedManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GeneratedManagementContract *GeneratedManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.contract.Transact(opts, method, params...)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractCaller) GetHostAddresses(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "GetHostAddresses")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractSession) GetHostAddresses() ([]string, error) {
	return _GeneratedManagementContract.Contract.GetHostAddresses(&_GeneratedManagementContract.CallOpts)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) GetHostAddresses() ([]string, error) {
	return _GeneratedManagementContract.Contract.GetHostAddresses(&_GeneratedManagementContract.CallOpts)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) AttestationRequests(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attestationRequests", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// AttestationRequests is a free data retrieval call binding the contract method 0x8ef74f89.
//
// Solidity: function attestationRequests(address ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) AttestationRequests(arg0 common.Address) (string, error) {
	return _GeneratedManagementContract.Contract.AttestationRequests(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Attested(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "attested", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Attested(arg0 common.Address) (bool, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// Attested is a free data retrieval call binding the contract method 0xd4c80664.
//
// Solidity: function attested(address ) view returns(bool)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Attested(arg0 common.Address) (bool, error) {
	return _GeneratedManagementContract.Contract.Attested(&_GeneratedManagementContract.CallOpts, arg0)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) Rollups(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "rollups", arg0, arg1)

	outstruct := new(struct {
		ParentHash   [32]byte
		AggregatorID common.Address
		L1Block      [32]byte
		Number       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ParentHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.AggregatorID = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.L1Block = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Number = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// Rollups is a free data retrieval call binding the contract method 0xe0643dfc.
//
// Solidity: function rollups(uint256 , uint256 ) view returns(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) Rollups(arg0 *big.Int, arg1 *big.Int) (struct {
	ParentHash   [32]byte
	AggregatorID common.Address
	L1Block      [32]byte
	Number       *big.Int
}, error) {
	return _GeneratedManagementContract.Contract.Rollups(&_GeneratedManagementContract.CallOpts, arg0, arg1)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) AddRollup(opts *bind.TransactOpts, ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "AddRollup", ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// AddRollup is a paid mutator transaction binding the contract method 0xe0fd84bd.
//
// Solidity: function AddRollup(bytes32 ParentHash, address AggregatorID, bytes32 L1Block, uint256 Number, string rollupData) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) AddRollup(ParentHash [32]byte, AggregatorID common.Address, L1Block [32]byte, Number *big.Int, rollupData string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddRollup(&_GeneratedManagementContract.TransactOpts, ParentHash, AggregatorID, L1Block, Number, rollupData)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret, hostAddress)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret, hostAddress)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x68e10383.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret, hostAddress)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RequestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RequestNetworkSecret(&_GeneratedManagementContract.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}
