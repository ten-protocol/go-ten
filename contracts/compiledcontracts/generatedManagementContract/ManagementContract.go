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
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"AddHostAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"hostAddresses\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50612705806100206000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c8063c719bf5011610066578063c719bf5014610159578063d4c8066414610175578063e0643dfc146101a5578063e0fd84bd146101d8578063e34fbfc8146101f45761009e565b8063324ff866146100a3578063597a9723146100c157806365a293c2146100dd5780638ef74f891461010d578063981214ba1461013d575b600080fd5b6100ab610210565b6040516100b89190611197565b60405180910390f35b6100db60048036038101906100d69190611302565b6102e9565b005b6100f760048036038101906100f29190611381565b610321565b60405161010491906113f8565b60405180910390f35b61012760048036038101906101229190611478565b6103cd565b60405161013491906113f8565b60405180910390f35b61015760048036038101906101529190611546565b61046d565b005b610173600480360381019061016e9190611645565b61060a565b005b61018f600480360381019061018a9190611478565b61069c565b60405161019c91906116c0565b60405180910390f35b6101bf60048036038101906101ba91906116db565b6106bc565b6040516101cf9493929190611752565b60405180910390f35b6101f260048036038101906101ed9190611819565b610729565b005b61020e600480360381019061020991906118b3565b610866565b005b60606003805480602002602001604051908101604052809291908181526020016000905b828210156102e05783829060005260206000200180546102539061192f565b80601f016020809104026020016040519081016040528092919081815260200182805461027f9061192f565b80156102cc5780601f106102a1576101008083540402835291602001916102cc565b820191906000526020600020905b8154815290600101906020018083116102af57829003601f168201915b505050505081526020019060010190610234565b50505050905090565b60038190806001815401808255809150506001900390600052602060002001600090919091909150908161031d9190611b0c565b5050565b6003818154811061033157600080fd5b90600052602060002001600091509050805461034c9061192f565b80601f01602080910402602001604051908101604052809291908181526020018280546103789061192f565b80156103c55780601f1061039a576101008083540402835291602001916103c5565b820191906000526020600020905b8154815290600101906020018083116103a857829003601f168201915b505050505081565b600160205280600052604060002060009150905080546103ec9061192f565b80601f01602080910402602001604051908101604052809291908181526020018280546104189061192f565b80156104655780601f1061043a57610100808354040283529160200191610465565b820191906000526020600020905b81548152906001019060200180831161044857829003601f168201915b505050505081565b6000600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050806104c857600080fd5b60006104f68686856040516020016104e293929190611c6d565b6040516020818303038152906040526108b9565b9050600061050482866108f4565b90508673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461053e8861091b565b6105478361091b565b604051602001610558929190611e38565b604051602081830303815290604052906105a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059f91906113f8565b60405180910390fd5b506001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050505050505050565b600460009054906101000a900460ff161561062457600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550505050565b60026020528060005260406000206000915054906101000a900460ff1681565b600060205281600052604060002081815481106106d857600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1661077f57600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826108b4929190611e88565b505050565b60006108c58251610ade565b826040516020016108d7929190611fa4565b604051602081830303815290604052805190602001209050919050565b60008060006109038585610c3e565b9150915061091081610cbf565b819250505092915050565b60606000602867ffffffffffffffff81111561093a576109396111d7565b5b6040519080825280601f01601f19166020018201604052801561096c5781602001600182028036833780820191505090505b50905060005b6014811015610ad457600081601361098a9190612002565b60086109969190612036565b60026109a291906121c3565b8573ffffffffffffffffffffffffffffffffffffffff166109c3919061223d565b60f81b9050600060108260f81c6109da919061227b565b60f81b905060008160f81c60106109f191906122ac565b8360f81c6109ff91906122e7565b60f81b9050610a0d82610e8b565b85856002610a1b9190612036565b81518110610a2c57610a2b61231b565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350610a6481610e8b565b856001866002610a749190612036565b610a7e919061234a565b81518110610a8f57610a8e61231b565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505050508080610acc906123a0565b915050610972565b5080915050919050565b606060008203610b25576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610c39565b600082905060005b60008214610b57578080610b40906123a0565b915050600a82610b50919061223d565b9150610b2d565b60008167ffffffffffffffff811115610b7357610b726111d7565b5b6040519080825280601f01601f191660200182016040528015610ba55781602001600182028036833780820191505090505b5090505b60008514610c3257600182610bbe9190612002565b9150600a85610bcd91906123e8565b6030610bd9919061234a565b60f81b818381518110610bef57610bee61231b565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610c2b919061223d565b9450610ba9565b8093505050505b919050565b6000806041835103610c7f5760008060006020860151925060408601519150606086015160001a9050610c7387828585610ed1565b94509450505050610cb8565b6040835103610caf576000806020850151915060408501519050610ca4868383610fdd565b935093505050610cb8565b60006002915091505b9250929050565b60006004811115610cd357610cd2612419565b5b816004811115610ce657610ce5612419565b5b0315610e885760016004811115610d0057610cff612419565b5b816004811115610d1357610d12612419565b5b03610d53576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d4a90612494565b60405180910390fd5b60026004811115610d6757610d66612419565b5b816004811115610d7a57610d79612419565b5b03610dba576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610db190612500565b60405180910390fd5b60036004811115610dce57610dcd612419565b5b816004811115610de157610de0612419565b5b03610e21576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e1890612592565b60405180910390fd5b600480811115610e3457610e33612419565b5b816004811115610e4757610e46612419565b5b03610e87576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e7e90612624565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610eb65760308260f81c610eac9190612644565b60f81b9050610ecc565b60578260f81c610ec69190612644565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610f0c576000600391509150610fd4565b601b8560ff1614158015610f245750601c8560ff1614155b15610f36576000600491509150610fd4565b600060018787878760405160008152602001604052604051610f5b949392919061268a565b6020604051602081039080840390855afa158015610f7d573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610fcb57600060019250925050610fd4565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c611020919061234a565b905061102e87828885610ed1565b935093505050935093915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b838110156110a2578082015181840152602081019050611087565b838111156110b1576000848401525b50505050565b6000601f19601f8301169050919050565b60006110d382611068565b6110dd8185611073565b93506110ed818560208601611084565b6110f6816110b7565b840191505092915050565b600061110d83836110c8565b905092915050565b6000602082019050919050565b600061112d8261103c565b6111378185611047565b93508360208202850161114985611058565b8060005b8581101561118557848403895281516111668582611101565b945061117183611115565b925060208a0199505060018101905061114d565b50829750879550505050505092915050565b600060208201905081810360008301526111b18184611122565b905092915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61120f826110b7565b810181811067ffffffffffffffff8211171561122e5761122d6111d7565b5b80604052505050565b60006112416111b9565b905061124d8282611206565b919050565b600067ffffffffffffffff82111561126d5761126c6111d7565b5b611276826110b7565b9050602081019050919050565b82818337600083830152505050565b60006112a56112a084611252565b611237565b9050828152602081018484840111156112c1576112c06111d2565b5b6112cc848285611283565b509392505050565b600082601f8301126112e9576112e86111cd565b5b81356112f9848260208601611292565b91505092915050565b600060208284031215611318576113176111c3565b5b600082013567ffffffffffffffff811115611336576113356111c8565b5b611342848285016112d4565b91505092915050565b6000819050919050565b61135e8161134b565b811461136957600080fd5b50565b60008135905061137b81611355565b92915050565b600060208284031215611397576113966111c3565b5b60006113a58482850161136c565b91505092915050565b600082825260208201905092915050565b60006113ca82611068565b6113d481856113ae565b93506113e4818560208601611084565b6113ed816110b7565b840191505092915050565b6000602082019050818103600083015261141281846113bf565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006114458261141a565b9050919050565b6114558161143a565b811461146057600080fd5b50565b6000813590506114728161144c565b92915050565b60006020828403121561148e5761148d6111c3565b5b600061149c84828501611463565b91505092915050565b600067ffffffffffffffff8211156114c0576114bf6111d7565b5b6114c9826110b7565b9050602081019050919050565b60006114e96114e4846114a5565b611237565b905082815260208101848484011115611505576115046111d2565b5b611510848285611283565b509392505050565b600082601f83011261152d5761152c6111cd565b5b813561153d8482602086016114d6565b91505092915050565b600080600080608085870312156115605761155f6111c3565b5b600061156e87828801611463565b945050602061157f87828801611463565b935050604085013567ffffffffffffffff8111156115a05761159f6111c8565b5b6115ac87828801611518565b925050606085013567ffffffffffffffff8111156115cd576115cc6111c8565b5b6115d987828801611518565b91505092959194509250565b600080fd5b600080fd5b60008083601f840112611605576116046111cd565b5b8235905067ffffffffffffffff811115611622576116216115e5565b5b60208301915083600182028301111561163e5761163d6115ea565b5b9250929050565b60008060006040848603121561165e5761165d6111c3565b5b600061166c86828701611463565b935050602084013567ffffffffffffffff81111561168d5761168c6111c8565b5b611699868287016115ef565b92509250509250925092565b60008115159050919050565b6116ba816116a5565b82525050565b60006020820190506116d560008301846116b1565b92915050565b600080604083850312156116f2576116f16111c3565b5b60006117008582860161136c565b92505060206117118582860161136c565b9150509250929050565b6000819050919050565b61172e8161171b565b82525050565b61173d8161143a565b82525050565b61174c8161134b565b82525050565b60006080820190506117676000830187611725565b6117746020830186611734565b6117816040830185611725565b61178e6060830184611743565b95945050505050565b6117a08161171b565b81146117ab57600080fd5b50565b6000813590506117bd81611797565b92915050565b60008083601f8401126117d9576117d86111cd565b5b8235905067ffffffffffffffff8111156117f6576117f56115e5565b5b602083019150836001820283011115611812576118116115ea565b5b9250929050565b60008060008060008060a08789031215611836576118356111c3565b5b600061184489828a016117ae565b965050602061185589828a01611463565b955050604061186689828a016117ae565b945050606061187789828a0161136c565b935050608087013567ffffffffffffffff811115611898576118976111c8565b5b6118a489828a016117c3565b92509250509295509295509295565b600080602083850312156118ca576118c96111c3565b5b600083013567ffffffffffffffff8111156118e8576118e76111c8565b5b6118f4858286016117c3565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061194757607f821691505b60208210810361195a57611959611900565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026119c27fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611985565b6119cc8683611985565b95508019841693508086168417925050509392505050565b6000819050919050565b6000611a09611a046119ff8461134b565b6119e4565b61134b565b9050919050565b6000819050919050565b611a23836119ee565b611a37611a2f82611a10565b848454611992565b825550505050565b600090565b611a4c611a3f565b611a57818484611a1a565b505050565b5b81811015611a7b57611a70600082611a44565b600181019050611a5d565b5050565b601f821115611ac057611a9181611960565b611a9a84611975565b81016020851015611aa9578190505b611abd611ab585611975565b830182611a5c565b50505b505050565b600082821c905092915050565b6000611ae360001984600802611ac5565b1980831691505092915050565b6000611afc8383611ad2565b9150826002028217905092915050565b611b1582611068565b67ffffffffffffffff811115611b2e57611b2d6111d7565b5b611b38825461192f565b611b43828285611a7f565b600060209050601f831160018114611b765760008415611b64578287015190505b611b6e8582611af0565b865550611bd6565b601f198416611b8486611960565b60005b82811015611bac57848901518255600182019150602085019450602081019050611b87565b86831015611bc95784890151611bc5601f891682611ad2565b8355505b6001600288020188555050505b505050505050565b60008160601b9050919050565b6000611bf682611bde565b9050919050565b6000611c0882611beb565b9050919050565b611c20611c1b8261143a565b611bfd565b82525050565b600081519050919050565b600081905092915050565b6000611c4782611c26565b611c518185611c31565b9350611c61818560208601611084565b80840191505092915050565b6000611c798286611c0f565b601482019150611c898285611c0f565b601482019150611c998284611c3c565b9150819050949350505050565b600081905092915050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e2774206d617463682000000000000000000000000000000000000000602082015250565b6000611d0d602d83611ca6565b9150611d1882611cb1565b602d82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b6000611d7f602483611ca6565b9150611d8a82611d23565b602482019050919050565b6000611da082611068565b611daa8185611ca6565b9350611dba818560208601611084565b80840191505092915050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b6000611e22602483611ca6565b9150611e2d82611dc6565b602482019050919050565b6000611e4382611d00565b9150611e4e82611d72565b9150611e5a8285611d95565b9150611e6582611e15565b9150611e718284611d95565b91508190509392505050565b600082905092915050565b611e928383611e7d565b67ffffffffffffffff811115611eab57611eaa6111d7565b5b611eb5825461192f565b611ec0828285611a7f565b6000601f831160018114611eef5760008415611edd578287013590505b611ee78582611af0565b865550611f4f565b601f198416611efd86611960565b60005b82811015611f2557848901358255600182019150602085019450602081019050611f00565b86831015611f425784890135611f3e601f891682611ad2565b8355505b6001600288020188555050505b50505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611f8e601a83611ca6565b9150611f9982611f58565b601a82019050919050565b6000611faf82611f81565b9150611fbb8285611d95565b9150611fc78284611c3c565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061200d8261134b565b91506120188361134b565b92508282101561202b5761202a611fd3565b5b828203905092915050565b60006120418261134b565b915061204c8361134b565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561208557612084611fd3565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b60018511156120e7578086048111156120c3576120c2611fd3565b5b60018516156120d25780820291505b80810290506120e085612090565b94506120a7565b94509492505050565b60008261210057600190506121bc565b8161210e57600090506121bc565b8160018114612124576002811461212e5761215d565b60019150506121bc565b60ff8411156121405761213f611fd3565b5b8360020a91508482111561215757612156611fd3565b5b506121bc565b5060208310610133831016604e8410600b84101617156121925782820a90508381111561218d5761218c611fd3565b5b6121bc565b61219f848484600161209d565b925090508184048111156121b6576121b5611fd3565b5b81810290505b9392505050565b60006121ce8261134b565b91506121d98361134b565b92506122067fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846120f0565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006122488261134b565b91506122538361134b565b9250826122635761226261220e565b5b828204905092915050565b600060ff82169050919050565b60006122868261226e565b91506122918361226e565b9250826122a1576122a061220e565b5b828204905092915050565b60006122b78261226e565b91506122c28361226e565b92508160ff04831182151516156122dc576122db611fd3565b5b828202905092915050565b60006122f28261226e565b91506122fd8361226e565b9250828210156123105761230f611fd3565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006123558261134b565b91506123608361134b565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561239557612394611fd3565b5b828201905092915050565b60006123ab8261134b565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036123dd576123dc611fd3565b5b600182019050919050565b60006123f38261134b565b91506123fe8361134b565b92508261240e5761240d61220e565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b600061247e6018836113ae565b915061248982612448565b602082019050919050565b600060208201905081810360008301526124ad81612471565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b60006124ea601f836113ae565b91506124f5826124b4565b602082019050919050565b60006020820190508181036000830152612519816124dd565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061257c6022836113ae565b915061258782612520565b604082019050919050565b600060208201905081810360008301526125ab8161256f565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061260e6022836113ae565b9150612619826125b2565b604082019050919050565b6000602082019050818103600083015261263d81612601565b9050919050565b600061264f8261226e565b915061265a8361226e565b92508260ff038211156126705761266f611fd3565b5b828201905092915050565b6126848161226e565b82525050565b600060808201905061269f6000830187611725565b6126ac602083018661267b565b6126b96040830185611725565b6126c66060830184611725565b9594505050505056fea2646970667358221220a9e04e697fad6411a4351cd0bbf7942158238fb47e9f312e9da0b27493cefa6764736f6c634300080f0033",
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

// HostAddresses is a free data retrieval call binding the contract method 0x65a293c2.
//
// Solidity: function hostAddresses(uint256 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCaller) HostAddresses(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _GeneratedManagementContract.contract.Call(opts, &out, "hostAddresses", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// HostAddresses is a free data retrieval call binding the contract method 0x65a293c2.
//
// Solidity: function hostAddresses(uint256 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractSession) HostAddresses(arg0 *big.Int) (string, error) {
	return _GeneratedManagementContract.Contract.HostAddresses(&_GeneratedManagementContract.CallOpts, arg0)
}

// HostAddresses is a free data retrieval call binding the contract method 0x65a293c2.
//
// Solidity: function hostAddresses(uint256 ) view returns(string)
func (_GeneratedManagementContract *GeneratedManagementContractCallerSession) HostAddresses(arg0 *big.Int) (string, error) {
	return _GeneratedManagementContract.Contract.HostAddresses(&_GeneratedManagementContract.CallOpts, arg0)
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

// AddHostAddress is a paid mutator transaction binding the contract method 0x597a9723.
//
// Solidity: function AddHostAddress(string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) AddHostAddress(opts *bind.TransactOpts, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "AddHostAddress", hostAddress)
}

// AddHostAddress is a paid mutator transaction binding the contract method 0x597a9723.
//
// Solidity: function AddHostAddress(string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) AddHostAddress(hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddHostAddress(&_GeneratedManagementContract.TransactOpts, hostAddress)
}

// AddHostAddress is a paid mutator transaction binding the contract method 0x597a9723.
//
// Solidity: function AddHostAddress(string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) AddHostAddress(hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.AddHostAddress(&_GeneratedManagementContract.TransactOpts, hostAddress)
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
