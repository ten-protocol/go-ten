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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"genesisAttestation\",\"type\":\"string\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"verifyAttester\",\"type\":\"bool\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061269f806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063d4c806641161005b578063d4c8066414610113578063e0643dfc14610143578063e0fd84bd14610176578063e34fbfc81461019257610088565b8063324ff8661461008d57806359a90071146100ab5780638ef74f89146100c7578063bbd79e15146100f7575b600080fd5b6100956101ae565b6040516100a291906110c9565b60405180910390f35b6100c560048036038101906100c09190611348565b610287565b005b6100e160048036038101906100dc919061140b565b610351565b6040516100ee9190611482565b60405180910390f35b610111600480360381019061010c919061157d565b6103f1565b005b61012d6004803603810190610128919061140b565b6105ce565b60405161013a919061166d565b60405180910390f35b61015d600480360381019061015891906116be565b6105ee565b60405161016d9493929190611735565b60405180910390f35b610190600480360381019061018b91906117a6565b61065b565b005b6101ac60048036038101906101a79190611840565b610798565b005b60606003805480602002602001604051908101604052809291908181526020016000905b8282101561027e5783829060005260206000200180546101f1906118bc565b80601f016020809104026020016040519081016040528092919081815260200182805461021d906118bc565b801561026a5780601f1061023f5761010080835404028352916020019161026a565b820191906000526020600020905b81548152906001019060200180831161024d57829003601f168201915b5050505050815260200190600101906101d2565b50505050905090565b600460009054906101000a900460ff16156102a157600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506003839080600181540180825580915050600190039060005260206000200160009091909190915090816103489190611a99565b50505050505050565b60016020528060005260406000206000915090508054610370906118bc565b80601f016020809104026020016040519081016040528092919081815260200182805461039c906118bc565b80156103e95780601f106103be576101008083540402835291602001916103e9565b820191906000526020600020905b8154815290600101906020018083116103cc57829003601f168201915b505050505081565b6000600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061044c57600080fd5b81156105385760006104828888868860405160200161046e9493929190611c36565b6040516020818303038152906040526107eb565b905060006104908288610826565b90508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146104ca8a61084d565b6104d38361084d565b6040516020016104e4929190611dd2565b60405160208183030381529060405290610534576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161052b9190611482565b60405180910390fd5b5050505b6001600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506003839080600181540180825580915050600190039060005260206000200160009091909190915090816105c49190611a99565b5050505050505050565b60026020528060005260406000206000915054906101000a900460ff1681565b6000602052816000526040600020818154811061060a57600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166106b157600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826107e6929190611e22565b505050565b60006107f78251610a10565b82604051602001610809929190611f3e565b604051602081830303815290604052805190602001209050919050565b60008060006108358585610b70565b9150915061084281610bf1565b819250505092915050565b60606000602867ffffffffffffffff81111561086c5761086b6111c7565b5b6040519080825280601f01601f19166020018201604052801561089e5781602001600182028036833780820191505090505b50905060005b6014811015610a065760008160136108bc9190611f9c565b60086108c89190611fd0565b60026108d4919061215d565b8573ffffffffffffffffffffffffffffffffffffffff166108f591906121d7565b60f81b9050600060108260f81c61090c9190612215565b60f81b905060008160f81c60106109239190612246565b8360f81c6109319190612281565b60f81b905061093f82610dbd565b8585600261094d9190611fd0565b8151811061095e5761095d6122b5565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061099681610dbd565b8560018660026109a69190611fd0565b6109b091906122e4565b815181106109c1576109c06122b5565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050505080806109fe9061233a565b9150506108a4565b5080915050919050565b606060008203610a57576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610b6b565b600082905060005b60008214610a89578080610a729061233a565b915050600a82610a8291906121d7565b9150610a5f565b60008167ffffffffffffffff811115610aa557610aa46111c7565b5b6040519080825280601f01601f191660200182016040528015610ad75781602001600182028036833780820191505090505b5090505b60008514610b6457600182610af09190611f9c565b9150600a85610aff9190612382565b6030610b0b91906122e4565b60f81b818381518110610b2157610b206122b5565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610b5d91906121d7565b9450610adb565b8093505050505b919050565b6000806041835103610bb15760008060006020860151925060408601519150606086015160001a9050610ba587828585610e03565b94509450505050610bea565b6040835103610be1576000806020850151915060408501519050610bd6868383610f0f565b935093505050610bea565b60006002915091505b9250929050565b60006004811115610c0557610c046123b3565b5b816004811115610c1857610c176123b3565b5b0315610dba5760016004811115610c3257610c316123b3565b5b816004811115610c4557610c446123b3565b5b03610c85576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7c9061242e565b60405180910390fd5b60026004811115610c9957610c986123b3565b5b816004811115610cac57610cab6123b3565b5b03610cec576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ce39061249a565b60405180910390fd5b60036004811115610d0057610cff6123b3565b5b816004811115610d1357610d126123b3565b5b03610d53576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d4a9061252c565b60405180910390fd5b600480811115610d6657610d656123b3565b5b816004811115610d7957610d786123b3565b5b03610db9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610db0906125be565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610de85760308260f81c610dde91906125de565b60f81b9050610dfe565b60578260f81c610df891906125de565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610e3e576000600391509150610f06565b601b8560ff1614158015610e565750601c8560ff1614155b15610e68576000600491509150610f06565b600060018787878760405160008152602001604052604051610e8d9493929190612624565b6020604051602081039080840390855afa158015610eaf573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610efd57600060019250925050610f06565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c610f5291906122e4565b9050610f6087828885610e03565b935093505050935093915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610fd4578082015181840152602081019050610fb9565b83811115610fe3576000848401525b50505050565b6000601f19601f8301169050919050565b600061100582610f9a565b61100f8185610fa5565b935061101f818560208601610fb6565b61102881610fe9565b840191505092915050565b600061103f8383610ffa565b905092915050565b6000602082019050919050565b600061105f82610f6e565b6110698185610f79565b93508360208202850161107b85610f8a565b8060005b858110156110b757848403895281516110988582611033565b94506110a383611047565b925060208a0199505060018101905061107f565b50829750879550505050505092915050565b600060208201905081810360008301526110e38184611054565b905092915050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061112a826110ff565b9050919050565b61113a8161111f565b811461114557600080fd5b50565b60008135905061115781611131565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126111825761118161115d565b5b8235905067ffffffffffffffff81111561119f5761119e611162565b5b6020830191508360018202830111156111bb576111ba611167565b5b9250929050565b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6111ff82610fe9565b810181811067ffffffffffffffff8211171561121e5761121d6111c7565b5b80604052505050565b60006112316110eb565b905061123d82826111f6565b919050565b600067ffffffffffffffff82111561125d5761125c6111c7565b5b61126682610fe9565b9050602081019050919050565b82818337600083830152505050565b600061129561129084611242565b611227565b9050828152602081018484840111156112b1576112b06111c2565b5b6112bc848285611273565b509392505050565b600082601f8301126112d9576112d861115d565b5b81356112e9848260208601611282565b91505092915050565b60008083601f8401126113085761130761115d565b5b8235905067ffffffffffffffff81111561132557611324611162565b5b60208301915083600182028301111561134157611340611167565b5b9250929050565b60008060008060008060808789031215611365576113646110f5565b5b600061137389828a01611148565b965050602087013567ffffffffffffffff811115611394576113936110fa565b5b6113a089828a0161116c565b9550955050604087013567ffffffffffffffff8111156113c3576113c26110fa565b5b6113cf89828a016112c4565b935050606087013567ffffffffffffffff8111156113f0576113ef6110fa565b5b6113fc89828a016112f2565b92509250509295509295509295565b600060208284031215611421576114206110f5565b5b600061142f84828501611148565b91505092915050565b600082825260208201905092915050565b600061145482610f9a565b61145e8185611438565b935061146e818560208601610fb6565b61147781610fe9565b840191505092915050565b6000602082019050818103600083015261149c8184611449565b905092915050565b600067ffffffffffffffff8211156114bf576114be6111c7565b5b6114c882610fe9565b9050602081019050919050565b60006114e86114e3846114a4565b611227565b905082815260208101848484011115611504576115036111c2565b5b61150f848285611273565b509392505050565b600082601f83011261152c5761152b61115d565b5b813561153c8482602086016114d5565b91505092915050565b60008115159050919050565b61155a81611545565b811461156557600080fd5b50565b60008135905061157781611551565b92915050565b60008060008060008060c0878903121561159a576115996110f5565b5b60006115a889828a01611148565b96505060206115b989828a01611148565b955050604087013567ffffffffffffffff8111156115da576115d96110fa565b5b6115e689828a01611517565b945050606087013567ffffffffffffffff811115611607576116066110fa565b5b61161389828a01611517565b935050608087013567ffffffffffffffff811115611634576116336110fa565b5b61164089828a016112c4565b92505060a061165189828a01611568565b9150509295509295509295565b61166781611545565b82525050565b6000602082019050611682600083018461165e565b92915050565b6000819050919050565b61169b81611688565b81146116a657600080fd5b50565b6000813590506116b881611692565b92915050565b600080604083850312156116d5576116d46110f5565b5b60006116e3858286016116a9565b92505060206116f4858286016116a9565b9150509250929050565b6000819050919050565b611711816116fe565b82525050565b6117208161111f565b82525050565b61172f81611688565b82525050565b600060808201905061174a6000830187611708565b6117576020830186611717565b6117646040830185611708565b6117716060830184611726565b95945050505050565b611783816116fe565b811461178e57600080fd5b50565b6000813590506117a08161177a565b92915050565b60008060008060008060a087890312156117c3576117c26110f5565b5b60006117d189828a01611791565b96505060206117e289828a01611148565b95505060406117f389828a01611791565b945050606061180489828a016116a9565b935050608087013567ffffffffffffffff811115611825576118246110fa565b5b61183189828a016112f2565b92509250509295509295509295565b60008060208385031215611857576118566110f5565b5b600083013567ffffffffffffffff811115611875576118746110fa565b5b611881858286016112f2565b92509250509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806118d457607f821691505b6020821081036118e7576118e661188d565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261194f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611912565b6119598683611912565b95508019841693508086168417925050509392505050565b6000819050919050565b600061199661199161198c84611688565b611971565b611688565b9050919050565b6000819050919050565b6119b08361197b565b6119c46119bc8261199d565b84845461191f565b825550505050565b600090565b6119d96119cc565b6119e48184846119a7565b505050565b5b81811015611a08576119fd6000826119d1565b6001810190506119ea565b5050565b601f821115611a4d57611a1e816118ed565b611a2784611902565b81016020851015611a36578190505b611a4a611a4285611902565b8301826119e9565b50505b505050565b600082821c905092915050565b6000611a7060001984600802611a52565b1980831691505092915050565b6000611a898383611a5f565b9150826002028217905092915050565b611aa282610f9a565b67ffffffffffffffff811115611abb57611aba6111c7565b5b611ac582546118bc565b611ad0828285611a0c565b600060209050601f831160018114611b035760008415611af1578287015190505b611afb8582611a7d565b865550611b63565b601f198416611b11866118ed565b60005b82811015611b3957848901518255600182019150602085019450602081019050611b14565b86831015611b565784890151611b52601f891682611a5f565b8355505b6001600288020188555050505b505050505050565b60008160601b9050919050565b6000611b8382611b6b565b9050919050565b6000611b9582611b78565b9050919050565b611bad611ba88261111f565b611b8a565b82525050565b600081905092915050565b6000611bc982610f9a565b611bd38185611bb3565b9350611be3818560208601610fb6565b80840191505092915050565b600081519050919050565b600081905092915050565b6000611c1082611bef565b611c1a8185611bfa565b9350611c2a818560208601610fb6565b80840191505092915050565b6000611c428287611b9c565b601482019150611c528286611b9c565b601482019150611c628285611bbe565b9150611c6e8284611c05565b915081905095945050505050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e2774206d617463682000000000000000000000000000000000000000602082015250565b6000611cd8602d83611bb3565b9150611ce382611c7c565b602d82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b6000611d4a602483611bb3565b9150611d5582611cee565b602482019050919050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b6000611dbc602483611bb3565b9150611dc782611d60565b602482019050919050565b6000611ddd82611ccb565b9150611de882611d3d565b9150611df48285611bbe565b9150611dff82611daf565b9150611e0b8284611bbe565b91508190509392505050565b600082905092915050565b611e2c8383611e17565b67ffffffffffffffff811115611e4557611e446111c7565b5b611e4f82546118bc565b611e5a828285611a0c565b6000601f831160018114611e895760008415611e77578287013590505b611e818582611a7d565b865550611ee9565b601f198416611e97866118ed565b60005b82811015611ebf57848901358255600182019150602085019450602081019050611e9a565b86831015611edc5784890135611ed8601f891682611a5f565b8355505b6001600288020188555050505b50505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611f28601a83611bb3565b9150611f3382611ef2565b601a82019050919050565b6000611f4982611f1b565b9150611f558285611bbe565b9150611f618284611c05565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611fa782611688565b9150611fb283611688565b925082821015611fc557611fc4611f6d565b5b828203905092915050565b6000611fdb82611688565b9150611fe683611688565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561201f5761201e611f6d565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b60018511156120815780860481111561205d5761205c611f6d565b5b600185161561206c5780820291505b808102905061207a8561202a565b9450612041565b94509492505050565b60008261209a5760019050612156565b816120a85760009050612156565b81600181146120be57600281146120c8576120f7565b6001915050612156565b60ff8411156120da576120d9611f6d565b5b8360020a9150848211156120f1576120f0611f6d565b5b50612156565b5060208310610133831016604e8410600b841016171561212c5782820a90508381111561212757612126611f6d565b5b612156565b6121398484846001612037565b925090508184048111156121505761214f611f6d565b5b81810290505b9392505050565b600061216882611688565b915061217383611688565b92506121a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff848461208a565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006121e282611688565b91506121ed83611688565b9250826121fd576121fc6121a8565b5b828204905092915050565b600060ff82169050919050565b600061222082612208565b915061222b83612208565b92508261223b5761223a6121a8565b5b828204905092915050565b600061225182612208565b915061225c83612208565b92508160ff048311821515161561227657612275611f6d565b5b828202905092915050565b600061228c82612208565b915061229783612208565b9250828210156122aa576122a9611f6d565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60006122ef82611688565b91506122fa83611688565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561232f5761232e611f6d565b5b828201905092915050565b600061234582611688565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361237757612376611f6d565b5b600182019050919050565b600061238d82611688565b915061239883611688565b9250826123a8576123a76121a8565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b6000612418601883611438565b9150612423826123e2565b602082019050919050565b600060208201905081810360008301526124478161240b565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b6000612484601f83611438565b915061248f8261244e565b602082019050919050565b600060208201905081810360008301526124b381612477565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b6000612516602283611438565b9150612521826124ba565b604082019050919050565b6000602082019050818103600083015261254581612509565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006125a8602283611438565b91506125b38261254c565b604082019050919050565b600060208201905081810360008301526125d78161259b565b9050919050565b60006125e982612208565b91506125f483612208565b92508260ff0382111561260a57612609611f6d565b5b828201905092915050565b61261e81612208565b82525050565b60006080820190506126396000830187611708565b6126466020830186612615565b6126536040830185611708565b6126606060830184611708565b9594505050505056fea2646970667358221220f3503be95e96aca3d97acf8a5ef7a2beaeeacf21e6c4e5299929bedf66d7b91b64736f6c634300080f0033",
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

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x59a90071.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress, string genesisAttestation) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, aggregatorID common.Address, initSecret []byte, hostAddress string, genesisAttestation string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "InitializeNetworkSecret", aggregatorID, initSecret, hostAddress, genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x59a90071.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress, string genesisAttestation) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte, hostAddress string, genesisAttestation string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret, hostAddress, genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x59a90071.
//
// Solidity: function InitializeNetworkSecret(address aggregatorID, bytes initSecret, string hostAddress, string genesisAttestation) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) InitializeNetworkSecret(aggregatorID common.Address, initSecret []byte, hostAddress string, genesisAttestation string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.InitializeNetworkSecret(&_GeneratedManagementContract.TransactOpts, aggregatorID, initSecret, hostAddress, genesisAttestation)
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
