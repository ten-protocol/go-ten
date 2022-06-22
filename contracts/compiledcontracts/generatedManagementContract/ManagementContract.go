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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50612256806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063d4c806641161005b578063d4c80664146100ea578063e0643dfc1461011a578063e0fd84bd1461014d578063e34fbfc8146101695761007d565b80638ef74f8914610082578063981214ba146100b2578063c719bf50146100ce575b600080fd5b61009c60048036038101906100979190610eaa565b610185565b6040516100a99190610f70565b60405180910390f35b6100cc60048036038101906100c791906110c7565b610225565b005b6100e860048036038101906100e391906111c6565b610406565b005b61010460048036038101906100ff9190610eaa565b610498565b6040516101119190611241565b60405180910390f35b610134600480360381019061012f9190611292565b6104b8565b6040516101449493929190611309565b60405180910390f35b610167600480360381019061016291906113d0565b610525565b005b610183600480360381019061017e919061146a565b610662565b005b600160205280600052604060002060009150905080546101a4906114e6565b80601f01602080910402602001604051908101604052809291908181526020018280546101d0906114e6565b801561021d5780601f106101f25761010080835404028352916020019161021d565b820191906000526020600020905b81548152906001019060200180831161020057829003601f168201915b505050505081565b6002600114610269576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161026090611563565b60405180910390fd5b6000600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050806102c457600080fd5b60006102f28686856040516020016102de93929190611612565b6040516020818303038152906040526106b5565b9050600061030082866106f0565b90508673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461033a88610717565b61034383610717565b6040516020016103549291906117dd565b604051602081830303815290604052906103a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161039b9190610f70565b60405180910390fd5b506001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050505050505050565b600360009054906101000a900460ff161561042057600080fd5b6001600360006101000a81548160ff0219169083151502179055506001600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550505050565b60026020528060005260406000206000915054906101000a900460ff1681565b600060205281600052604060002081815481106104d457600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1661057b57600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826106b09291906119d9565b505050565b60006106c182516108da565b826040516020016106d3929190611af5565b604051602081830303815290604052805190602001209050919050565b60008060006106ff8585610a3a565b9150915061070c81610abb565b819250505092915050565b60606000602867ffffffffffffffff81111561073657610735610f9c565b5b6040519080825280601f01601f1916602001820160405280156107685781602001600182028036833780820191505090505b50905060005b60148110156108d05760008160136107869190611b53565b60086107929190611b87565b600261079e9190611d14565b8573ffffffffffffffffffffffffffffffffffffffff166107bf9190611d8e565b60f81b9050600060108260f81c6107d69190611dcc565b60f81b905060008160f81c60106107ed9190611dfd565b8360f81c6107fb9190611e38565b60f81b905061080982610c87565b858560026108179190611b87565b8151811061082857610827611e6c565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061086081610c87565b8560018660026108709190611b87565b61087a9190611e9b565b8151811061088b5761088a611e6c565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050505080806108c890611ef1565b91505061076e565b5080915050919050565b606060008203610921576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610a35565b600082905060005b6000821461095357808061093c90611ef1565b915050600a8261094c9190611d8e565b9150610929565b60008167ffffffffffffffff81111561096f5761096e610f9c565b5b6040519080825280601f01601f1916602001820160405280156109a15781602001600182028036833780820191505090505b5090505b60008514610a2e576001826109ba9190611b53565b9150600a856109c99190611f39565b60306109d59190611e9b565b60f81b8183815181106109eb576109ea611e6c565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610a279190611d8e565b94506109a5565b8093505050505b919050565b6000806041835103610a7b5760008060006020860151925060408601519150606086015160001a9050610a6f87828585610ccd565b94509450505050610ab4565b6040835103610aab576000806020850151915060408501519050610aa0868383610dd9565b935093505050610ab4565b60006002915091505b9250929050565b60006004811115610acf57610ace611f6a565b5b816004811115610ae257610ae1611f6a565b5b0315610c845760016004811115610afc57610afb611f6a565b5b816004811115610b0f57610b0e611f6a565b5b03610b4f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b4690611fe5565b60405180910390fd5b60026004811115610b6357610b62611f6a565b5b816004811115610b7657610b75611f6a565b5b03610bb6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bad90612051565b60405180910390fd5b60036004811115610bca57610bc9611f6a565b5b816004811115610bdd57610bdc611f6a565b5b03610c1d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c14906120e3565b60405180910390fd5b600480811115610c3057610c2f611f6a565b5b816004811115610c4357610c42611f6a565b5b03610c83576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7a90612175565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610cb25760308260f81c610ca89190612195565b60f81b9050610cc8565b60578260f81c610cc29190612195565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610d08576000600391509150610dd0565b601b8560ff1614158015610d205750601c8560ff1614155b15610d32576000600491509150610dd0565b600060018787878760405160008152602001604052604051610d5794939291906121db565b6020604051602081039080840390855afa158015610d79573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610dc757600060019250925050610dd0565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c610e1c9190611e9b565b9050610e2a87828885610ccd565b935093505050935093915050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610e7782610e4c565b9050919050565b610e8781610e6c565b8114610e9257600080fd5b50565b600081359050610ea481610e7e565b92915050565b600060208284031215610ec057610ebf610e42565b5b6000610ece84828501610e95565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610f11578082015181840152602081019050610ef6565b83811115610f20576000848401525b50505050565b6000601f19601f8301169050919050565b6000610f4282610ed7565b610f4c8185610ee2565b9350610f5c818560208601610ef3565b610f6581610f26565b840191505092915050565b60006020820190508181036000830152610f8a8184610f37565b905092915050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610fd482610f26565b810181811067ffffffffffffffff82111715610ff357610ff2610f9c565b5b80604052505050565b6000611006610e38565b90506110128282610fcb565b919050565b600067ffffffffffffffff82111561103257611031610f9c565b5b61103b82610f26565b9050602081019050919050565b82818337600083830152505050565b600061106a61106584611017565b610ffc565b90508281526020810184848401111561108657611085610f97565b5b611091848285611048565b509392505050565b600082601f8301126110ae576110ad610f92565b5b81356110be848260208601611057565b91505092915050565b600080600080608085870312156110e1576110e0610e42565b5b60006110ef87828801610e95565b945050602061110087828801610e95565b935050604085013567ffffffffffffffff81111561112157611120610e47565b5b61112d87828801611099565b925050606085013567ffffffffffffffff81111561114e5761114d610e47565b5b61115a87828801611099565b91505092959194509250565b600080fd5b600080fd5b60008083601f84011261118657611185610f92565b5b8235905067ffffffffffffffff8111156111a3576111a2611166565b5b6020830191508360018202830111156111bf576111be61116b565b5b9250929050565b6000806000604084860312156111df576111de610e42565b5b60006111ed86828701610e95565b935050602084013567ffffffffffffffff81111561120e5761120d610e47565b5b61121a86828701611170565b92509250509250925092565b60008115159050919050565b61123b81611226565b82525050565b60006020820190506112566000830184611232565b92915050565b6000819050919050565b61126f8161125c565b811461127a57600080fd5b50565b60008135905061128c81611266565b92915050565b600080604083850312156112a9576112a8610e42565b5b60006112b78582860161127d565b92505060206112c88582860161127d565b9150509250929050565b6000819050919050565b6112e5816112d2565b82525050565b6112f481610e6c565b82525050565b6113038161125c565b82525050565b600060808201905061131e60008301876112dc565b61132b60208301866112eb565b61133860408301856112dc565b61134560608301846112fa565b95945050505050565b611357816112d2565b811461136257600080fd5b50565b6000813590506113748161134e565b92915050565b60008083601f8401126113905761138f610f92565b5b8235905067ffffffffffffffff8111156113ad576113ac611166565b5b6020830191508360018202830111156113c9576113c861116b565b5b9250929050565b60008060008060008060a087890312156113ed576113ec610e42565b5b60006113fb89828a01611365565b965050602061140c89828a01610e95565b955050604061141d89828a01611365565b945050606061142e89828a0161127d565b935050608087013567ffffffffffffffff81111561144f5761144e610e47565b5b61145b89828a0161137a565b92509250509295509295509295565b6000806020838503121561148157611480610e42565b5b600083013567ffffffffffffffff81111561149f5761149e610e47565b5b6114ab8582860161137a565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806114fe57607f821691505b602082108103611511576115106114b7565b5b50919050565b7f6a6a6a2073686f756c64206661696c0000000000000000000000000000000000600082015250565b600061154d600f83610ee2565b915061155882611517565b602082019050919050565b6000602082019050818103600083015261157c81611540565b9050919050565b60008160601b9050919050565b600061159b82611583565b9050919050565b60006115ad82611590565b9050919050565b6115c56115c082610e6c565b6115a2565b82525050565b600081519050919050565b600081905092915050565b60006115ec826115cb565b6115f681856115d6565b9350611606818560208601610ef3565b80840191505092915050565b600061161e82866115b4565b60148201915061162e82856115b4565b60148201915061163e82846115e1565b9150819050949350505050565b600081905092915050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e74206d61746368200000000000000000000000000000000000000000602082015250565b60006116b2602c8361164b565b91506116bd82611656565b602c82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b600061172460248361164b565b915061172f826116c8565b602482019050919050565b600061174582610ed7565b61174f818561164b565b935061175f818560208601610ef3565b80840191505092915050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b60006117c760248361164b565b91506117d28261176b565b602482019050919050565b60006117e8826116a5565b91506117f382611717565b91506117ff828561173a565b915061180a826117ba565b9150611816828461173a565b91508190509392505050565b600082905092915050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261188f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611852565b6118998683611852565b95508019841693508086168417925050509392505050565b6000819050919050565b60006118d66118d16118cc8461125c565b6118b1565b61125c565b9050919050565b6000819050919050565b6118f0836118bb565b6119046118fc826118dd565b84845461185f565b825550505050565b600090565b61191961190c565b6119248184846118e7565b505050565b5b818110156119485761193d600082611911565b60018101905061192a565b5050565b601f82111561198d5761195e8161182d565b61196784611842565b81016020851015611976578190505b61198a61198285611842565b830182611929565b50505b505050565b600082821c905092915050565b60006119b060001984600802611992565b1980831691505092915050565b60006119c9838361199f565b9150826002028217905092915050565b6119e38383611822565b67ffffffffffffffff8111156119fc576119fb610f9c565b5b611a0682546114e6565b611a1182828561194c565b6000601f831160018114611a405760008415611a2e578287013590505b611a3885826119bd565b865550611aa0565b601f198416611a4e8661182d565b60005b82811015611a7657848901358255600182019150602085019450602081019050611a51565b86831015611a935784890135611a8f601f89168261199f565b8355505b6001600288020188555050505b50505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611adf601a8361164b565b9150611aea82611aa9565b601a82019050919050565b6000611b0082611ad2565b9150611b0c828561173a565b9150611b1882846115e1565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611b5e8261125c565b9150611b698361125c565b925082821015611b7c57611b7b611b24565b5b828203905092915050565b6000611b928261125c565b9150611b9d8361125c565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611bd657611bd5611b24565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b6001851115611c3857808604811115611c1457611c13611b24565b5b6001851615611c235780820291505b8081029050611c3185611be1565b9450611bf8565b94509492505050565b600082611c515760019050611d0d565b81611c5f5760009050611d0d565b8160018114611c755760028114611c7f57611cae565b6001915050611d0d565b60ff841115611c9157611c90611b24565b5b8360020a915084821115611ca857611ca7611b24565b5b50611d0d565b5060208310610133831016604e8410600b8410161715611ce35782820a905083811115611cde57611cdd611b24565b5b611d0d565b611cf08484846001611bee565b92509050818404811115611d0757611d06611b24565b5b81810290505b9392505050565b6000611d1f8261125c565b9150611d2a8361125c565b9250611d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8484611c41565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611d998261125c565b9150611da48361125c565b925082611db457611db3611d5f565b5b828204905092915050565b600060ff82169050919050565b6000611dd782611dbf565b9150611de283611dbf565b925082611df257611df1611d5f565b5b828204905092915050565b6000611e0882611dbf565b9150611e1383611dbf565b92508160ff0483118215151615611e2d57611e2c611b24565b5b828202905092915050565b6000611e4382611dbf565b9150611e4e83611dbf565b925082821015611e6157611e60611b24565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000611ea68261125c565b9150611eb18361125c565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611ee657611ee5611b24565b5b828201905092915050565b6000611efc8261125c565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611f2e57611f2d611b24565b5b600182019050919050565b6000611f448261125c565b9150611f4f8361125c565b925082611f5f57611f5e611d5f565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b6000611fcf601883610ee2565b9150611fda82611f99565b602082019050919050565b60006020820190508181036000830152611ffe81611fc2565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b600061203b601f83610ee2565b915061204682612005565b602082019050919050565b6000602082019050818103600083015261206a8161202e565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006120cd602283610ee2565b91506120d882612071565b604082019050919050565b600060208201905081810360008301526120fc816120c0565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061215f602283610ee2565b915061216a82612103565b604082019050919050565b6000602082019050818103600083015261218e81612152565b9050919050565b60006121a082611dbf565b91506121ab83611dbf565b92508260ff038211156121c1576121c0611b24565b5b828201905092915050565b6121d581611dbf565b82525050565b60006080820190506121f060008301876112dc565b6121fd60208301866121cc565b61220a60408301856112dc565b61221760608301846112dc565b9594505050505056fea2646970667358221220b58823c4a86c70dd793872f473056cbcaa7ecc1ad907c2f1e2212fa494c430d564736f6c634300080f0033",
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

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0xc719bf50.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret)
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

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0x981214ba.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret)
}
