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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"rollupData\",\"type\":\"string\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attestationRequests\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"hostAddresses\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollups\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"ParentHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"L1Block\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"Number\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50612738806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063d4c8066411610066578063d4c8066414610132578063e0643dfc14610162578063e0fd84bd14610195578063e34fbfc8146101b1578063f1846d0c146101cd57610093565b8063324ff8661461009857806365a293c2146100b657806368e10383146100e65780638ef74f8914610102575b600080fd5b6100a06101e9565b6040516100ad91906111a6565b60405180910390f35b6100d060048036038101906100cb9190611212565b6102c2565b6040516100dd9190611289565b60405180910390f35b61010060048036038101906100fb919061149e565b61036e565b005b61011c6004803603810190610117919061152e565b610436565b6040516101299190611289565b60405180910390f35b61014c6004803603810190610147919061152e565b6104d6565b6040516101599190611576565b60405180910390f35b61017c60048036038101906101779190611591565b6104f6565b60405161018c9493929190611608565b60405180910390f35b6101af60048036038101906101aa91906116cf565b610563565b005b6101cb60048036038101906101c69190611769565b6106a0565b005b6101e760048036038101906101e29190611857565b6106f3565b005b60606003805480602002602001604051908101604052809291908181526020016000905b828210156102b957838290600052602060002001805461022c90611955565b80601f016020809104026020016040519081016040528092919081815260200182805461025890611955565b80156102a55780601f1061027a576101008083540402835291602001916102a5565b820191906000526020600020905b81548152906001019060200180831161028857829003601f168201915b50505050508152602001906001019061020d565b50505050905090565b600381815481106102d257600080fd5b9060005260206000200160009150905080546102ed90611955565b80601f016020809104026020016040519081016040528092919081815260200182805461031990611955565b80156103665780601f1061033b57610100808354040283529160200191610366565b820191906000526020600020905b81548152906001019060200180831161034957829003601f168201915b505050505081565b600460009054906101000a900460ff161561038857600080fd5b6001600460006101000a81548160ff0219169083151502179055506001600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555060038190806001815401808255809150506001900390600052602060002001600090919091909150908161042f9190611b32565b5050505050565b6001602052806000526040600020600091509050805461045590611955565b80601f016020809104026020016040519081016040528092919081815260200182805461048190611955565b80156104ce5780601f106104a3576101008083540402835291602001916104ce565b820191906000526020600020905b8154815290600101906020018083116104b157829003601f168201915b505050505081565b60026020528060005260406000206000915054906101000a900460ff1681565b6000602052816000526040600020818154811061051257600080fd5b9060005260206000209060040201600091509150508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154905084565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166105b957600080fd5b600060405180608001604052808881526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581525090506000804381526020019081526020016000208190806001815401808255809150506001900390600052602060002090600402016000909190919091506000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030155505050505050505050565b8181600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002091826106ee929190611c0f565b505050565b6000600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508061074e57600080fd5b600061077e8787868660405160200161076a9493929190611daa565b6040516020818303038152906040526108c8565b9050600061078c8287610903565b90508773ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146107c68961092a565b6107cf8361092a565b6040516020016107e0929190611f46565b60405160208183030381529060405290610830576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108279190611289565b60405180910390fd5b506001600260008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506003849080600181540180825580915050600190039060005260206000200160009091909190915090816108bd9190611b32565b505050505050505050565b60006108d48251610aed565b826040516020016108e6929190611fd7565b604051602081830303815290604052805190602001209050919050565b60008060006109128585610c4d565b9150915061091f81610cce565b819250505092915050565b60606000602867ffffffffffffffff81111561094957610948611373565b5b6040519080825280601f01601f19166020018201604052801561097b5781602001600182028036833780820191505090505b50905060005b6014811015610ae35760008160136109999190612035565b60086109a59190612069565b60026109b191906121f6565b8573ffffffffffffffffffffffffffffffffffffffff166109d29190612270565b60f81b9050600060108260f81c6109e991906122ae565b60f81b905060008160f81c6010610a0091906122df565b8360f81c610a0e919061231a565b60f81b9050610a1c82610e9a565b85856002610a2a9190612069565b81518110610a3b57610a3a61234e565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350610a7381610e9a565b856001866002610a839190612069565b610a8d919061237d565b81518110610a9e57610a9d61234e565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053505050508080610adb906123d3565b915050610981565b5080915050919050565b606060008203610b34576040518060400160405280600181526020017f30000000000000000000000000000000000000000000000000000000000000008152509050610c48565b600082905060005b60008214610b66578080610b4f906123d3565b915050600a82610b5f9190612270565b9150610b3c565b60008167ffffffffffffffff811115610b8257610b81611373565b5b6040519080825280601f01601f191660200182016040528015610bb45781602001600182028036833780820191505090505b5090505b60008514610c4157600182610bcd9190612035565b9150600a85610bdc919061241b565b6030610be8919061237d565b60f81b818381518110610bfe57610bfd61234e565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600a85610c3a9190612270565b9450610bb8565b8093505050505b919050565b6000806041835103610c8e5760008060006020860151925060408601519150606086015160001a9050610c8287828585610ee0565b94509450505050610cc7565b6040835103610cbe576000806020850151915060408501519050610cb3868383610fec565b935093505050610cc7565b60006002915091505b9250929050565b60006004811115610ce257610ce161244c565b5b816004811115610cf557610cf461244c565b5b0315610e975760016004811115610d0f57610d0e61244c565b5b816004811115610d2257610d2161244c565b5b03610d62576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d59906124c7565b60405180910390fd5b60026004811115610d7657610d7561244c565b5b816004811115610d8957610d8861244c565b5b03610dc9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dc090612533565b60405180910390fd5b60036004811115610ddd57610ddc61244c565b5b816004811115610df057610def61244c565b5b03610e30576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e27906125c5565b60405180910390fd5b600480811115610e4357610e4261244c565b5b816004811115610e5657610e5561244c565b5b03610e96576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e8d90612657565b60405180910390fd5b5b50565b6000600a8260f81c60ff161015610ec55760308260f81c610ebb9190612677565b60f81b9050610edb565b60578260f81c610ed59190612677565b60f81b90505b919050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c1115610f1b576000600391509150610fe3565b601b8560ff1614158015610f335750601c8560ff1614155b15610f45576000600491509150610fe3565b600060018787878760405160008152602001604052604051610f6a94939291906126bd565b6020604051602081039080840390855afa158015610f8c573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610fda57600060019250925050610fe3565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c61102f919061237d565b905061103d87828885610ee0565b935093505050935093915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b838110156110b1578082015181840152602081019050611096565b838111156110c0576000848401525b50505050565b6000601f19601f8301169050919050565b60006110e282611077565b6110ec8185611082565b93506110fc818560208601611093565b611105816110c6565b840191505092915050565b600061111c83836110d7565b905092915050565b6000602082019050919050565b600061113c8261104b565b6111468185611056565b93508360208202850161115885611067565b8060005b8581101561119457848403895281516111758582611110565b945061118083611124565b925060208a0199505060018101905061115c565b50829750879550505050505092915050565b600060208201905081810360008301526111c08184611131565b905092915050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b6111ef816111dc565b81146111fa57600080fd5b50565b60008135905061120c816111e6565b92915050565b600060208284031215611228576112276111d2565b5b6000611236848285016111fd565b91505092915050565b600082825260208201905092915050565b600061125b82611077565b611265818561123f565b9350611275818560208601611093565b61127e816110c6565b840191505092915050565b600060208201905081810360008301526112a38184611250565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006112d6826112ab565b9050919050565b6112e6816112cb565b81146112f157600080fd5b50565b600081359050611303816112dd565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261132e5761132d611309565b5b8235905067ffffffffffffffff81111561134b5761134a61130e565b5b60208301915083600182028301111561136757611366611313565b5b9250929050565b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6113ab826110c6565b810181811067ffffffffffffffff821117156113ca576113c9611373565b5b80604052505050565b60006113dd6111c8565b90506113e982826113a2565b919050565b600067ffffffffffffffff82111561140957611408611373565b5b611412826110c6565b9050602081019050919050565b82818337600083830152505050565b600061144161143c846113ee565b6113d3565b90508281526020810184848401111561145d5761145c61136e565b5b61146884828561141f565b509392505050565b600082601f83011261148557611484611309565b5b813561149584826020860161142e565b91505092915050565b600080600080606085870312156114b8576114b76111d2565b5b60006114c6878288016112f4565b945050602085013567ffffffffffffffff8111156114e7576114e66111d7565b5b6114f387828801611318565b9350935050604085013567ffffffffffffffff811115611516576115156111d7565b5b61152287828801611470565b91505092959194509250565b600060208284031215611544576115436111d2565b5b6000611552848285016112f4565b91505092915050565b60008115159050919050565b6115708161155b565b82525050565b600060208201905061158b6000830184611567565b92915050565b600080604083850312156115a8576115a76111d2565b5b60006115b6858286016111fd565b92505060206115c7858286016111fd565b9150509250929050565b6000819050919050565b6115e4816115d1565b82525050565b6115f3816112cb565b82525050565b611602816111dc565b82525050565b600060808201905061161d60008301876115db565b61162a60208301866115ea565b61163760408301856115db565b61164460608301846115f9565b95945050505050565b611656816115d1565b811461166157600080fd5b50565b6000813590506116738161164d565b92915050565b60008083601f84011261168f5761168e611309565b5b8235905067ffffffffffffffff8111156116ac576116ab61130e565b5b6020830191508360018202830111156116c8576116c7611313565b5b9250929050565b60008060008060008060a087890312156116ec576116eb6111d2565b5b60006116fa89828a01611664565b965050602061170b89828a016112f4565b955050604061171c89828a01611664565b945050606061172d89828a016111fd565b935050608087013567ffffffffffffffff81111561174e5761174d6111d7565b5b61175a89828a01611679565b92509250509295509295509295565b600080602083850312156117805761177f6111d2565b5b600083013567ffffffffffffffff81111561179e5761179d6111d7565b5b6117aa85828601611679565b92509250509250929050565b600067ffffffffffffffff8211156117d1576117d0611373565b5b6117da826110c6565b9050602081019050919050565b60006117fa6117f5846117b6565b6113d3565b9050828152602081018484840111156118165761181561136e565b5b61182184828561141f565b509392505050565b600082601f83011261183e5761183d611309565b5b813561184e8482602086016117e7565b91505092915050565b600080600080600060a08688031215611873576118726111d2565b5b6000611881888289016112f4565b9550506020611892888289016112f4565b945050604086013567ffffffffffffffff8111156118b3576118b26111d7565b5b6118bf88828901611829565b935050606086013567ffffffffffffffff8111156118e0576118df6111d7565b5b6118ec88828901611829565b925050608086013567ffffffffffffffff81111561190d5761190c6111d7565b5b61191988828901611470565b9150509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061196d57607f821691505b6020821081036119805761197f611926565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026119e87fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826119ab565b6119f286836119ab565b95508019841693508086168417925050509392505050565b6000819050919050565b6000611a2f611a2a611a25846111dc565b611a0a565b6111dc565b9050919050565b6000819050919050565b611a4983611a14565b611a5d611a5582611a36565b8484546119b8565b825550505050565b600090565b611a72611a65565b611a7d818484611a40565b505050565b5b81811015611aa157611a96600082611a6a565b600181019050611a83565b5050565b601f821115611ae657611ab781611986565b611ac08461199b565b81016020851015611acf578190505b611ae3611adb8561199b565b830182611a82565b50505b505050565b600082821c905092915050565b6000611b0960001984600802611aeb565b1980831691505092915050565b6000611b228383611af8565b9150826002028217905092915050565b611b3b82611077565b67ffffffffffffffff811115611b5457611b53611373565b5b611b5e8254611955565b611b69828285611aa5565b600060209050601f831160018114611b9c5760008415611b8a578287015190505b611b948582611b16565b865550611bfc565b601f198416611baa86611986565b60005b82811015611bd257848901518255600182019150602085019450602081019050611bad565b86831015611bef5784890151611beb601f891682611af8565b8355505b6001600288020188555050505b505050505050565b600082905092915050565b611c198383611c04565b67ffffffffffffffff811115611c3257611c31611373565b5b611c3c8254611955565b611c47828285611aa5565b6000601f831160018114611c765760008415611c64578287013590505b611c6e8582611b16565b865550611cd6565b601f198416611c8486611986565b60005b82811015611cac57848901358255600182019150602085019450602081019050611c87565b86831015611cc95784890135611cc5601f891682611af8565b8355505b6001600288020188555050505b50505050505050565b60008160601b9050919050565b6000611cf782611cdf565b9050919050565b6000611d0982611cec565b9050919050565b611d21611d1c826112cb565b611cfe565b82525050565b600081519050919050565b600081905092915050565b6000611d4882611d27565b611d528185611d32565b9350611d62818560208601611093565b80840191505092915050565b600081905092915050565b6000611d8482611077565b611d8e8185611d6e565b9350611d9e818560208601611093565b80840191505092915050565b6000611db68287611d10565b601482019150611dc68286611d10565b601482019150611dd68285611d3d565b9150611de28284611d79565b915081905095945050505050565b7f7265636f7665726564206164647265737320616e64206174746573746572494460008201527f20646f6e2774206d617463682000000000000000000000000000000000000000602082015250565b6000611e4c602d83611d6e565b9150611e5782611df0565b602d82019050919050565b7f0a2045787065637465643a20202020202020202020202020202020202020202060008201527f2020202000000000000000000000000000000000000000000000000000000000602082015250565b6000611ebe602483611d6e565b9150611ec982611e62565b602482019050919050565b7f0a202f207265636f7665726564416464725369676e656443616c63756c61746560008201527f643a202000000000000000000000000000000000000000000000000000000000602082015250565b6000611f30602483611d6e565b9150611f3b82611ed4565b602482019050919050565b6000611f5182611e3f565b9150611f5c82611eb1565b9150611f688285611d79565b9150611f7382611f23565b9150611f7f8284611d79565b91508190509392505050565b7f19457468657265756d205369676e6564204d6573736167653a0a000000000000600082015250565b6000611fc1601a83611d6e565b9150611fcc82611f8b565b601a82019050919050565b6000611fe282611fb4565b9150611fee8285611d79565b9150611ffa8284611d3d565b91508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000612040826111dc565b915061204b836111dc565b92508282101561205e5761205d612006565b5b828203905092915050565b6000612074826111dc565b915061207f836111dc565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156120b8576120b7612006565b5b828202905092915050565b60008160011c9050919050565b6000808291508390505b600185111561211a578086048111156120f6576120f5612006565b5b60018516156121055780820291505b8081029050612113856120c3565b94506120da565b94509492505050565b60008261213357600190506121ef565b8161214157600090506121ef565b8160018114612157576002811461216157612190565b60019150506121ef565b60ff84111561217357612172612006565b5b8360020a91508482111561218a57612189612006565b5b506121ef565b5060208310610133831016604e8410600b84101617156121c55782820a9050838111156121c0576121bf612006565b5b6121ef565b6121d284848460016120d0565b925090508184048111156121e9576121e8612006565b5b81810290505b9392505050565b6000612201826111dc565b915061220c836111dc565b92506122397fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8484612123565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061227b826111dc565b9150612286836111dc565b92508261229657612295612241565b5b828204905092915050565b600060ff82169050919050565b60006122b9826122a1565b91506122c4836122a1565b9250826122d4576122d3612241565b5b828204905092915050565b60006122ea826122a1565b91506122f5836122a1565b92508160ff048311821515161561230f5761230e612006565b5b828202905092915050565b6000612325826122a1565b9150612330836122a1565b92508282101561234357612342612006565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000612388826111dc565b9150612393836111dc565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156123c8576123c7612006565b5b828201905092915050565b60006123de826111dc565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036124105761240f612006565b5b600182019050919050565b6000612426826111dc565b9150612431836111dc565b92508261244157612440612241565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b60006124b160188361123f565b91506124bc8261247b565b602082019050919050565b600060208201905081810360008301526124e0816124a4565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b600061251d601f8361123f565b9150612528826124e7565b602082019050919050565b6000602082019050818103600083015261254c81612510565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b60006125af60228361123f565b91506125ba82612553565b604082019050919050565b600060208201905081810360008301526125de816125a2565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061264160228361123f565b915061264c826125e5565b604082019050919050565b6000602082019050818103600083015261267081612634565b9050919050565b6000612682826122a1565b915061268d836122a1565b92508260ff038211156126a3576126a2612006565b5b828201905092915050565b6126b7816122a1565b82525050565b60006080820190506126d260008301876115db565b6126df60208301866126ae565b6126ec60408301856115db565b6126f960608301846115db565b9594505050505056fea26469706673582212205ef63f23517d65809967b4298383908341339b44fb803df396acbc13dea621de64736f6c634300080f0033",
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

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xf1846d0c.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret, hostAddress)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xf1846d0c.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xf1846d0c.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress) returns()
func (_GeneratedManagementContract *GeneratedManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string) (*types.Transaction, error) {
	return _GeneratedManagementContract.Contract.RespondNetworkSecret(&_GeneratedManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress)
}
