// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TransactionsAnalyzer

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

// TransactionsAnalyzerBlockTransactions is an auto generated low-level Go binding around an user-defined struct.
type TransactionsAnalyzerBlockTransactions struct {
	Transactions [][]byte
}

// TransactionsAnalyzerMetaData contains all meta data concerning the TransactionsAnalyzer contract.
var TransactionsAnalyzerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EOA_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HOOK_CALLER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"}],\"name\":\"addOnBlockEndCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"eoaAdmin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"authorizedCaller\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"transactions\",\"type\":\"bytes[]\"}],\"internalType\":\"structTransactionsAnalyzer.BlockTransactions\",\"name\":\"_block\",\"type\":\"tuple\"}],\"name\":\"onBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60806040523461001a57604051611812610020823961181290f35b600080fdfe6080604052600436101561001257600080fd5b60003560e01c806301ffc9a7146100d2578063248a9ca3146100cd5780632f2ff15d146100c857806336568abe146100c3578063485cc955146100be578063508a50f4146100b95780635f03a661146100b457806391d14854146100af5780639dbbcf8e146100aa578063a217fddf146100a5578063d547741f146100a05763ee546fd803610102576103b5565b610388565b61036d565b610336565b6102e6565b6102ad565b610274565b610250565b61021d565b6101ff565b610189565b610131565b7fffffffff0000000000000000000000000000000000000000000000000000000081165b0361010257565b600080fd5b90503590610114826100d7565b565b906020828203126101025761012a91610107565b90565b9052565b346101025761015e61014c610147366004610116565b6103cd565b60405191829182901515815260200190565b0390f35b806100fb565b9050359061011482610162565b906020828203126101025761012a91610168565b346101025761015e6101a461019f366004610175565b6104f0565b6040519182918290815260200190565b6001600160a01b031690565b6001600160a01b0381166100fb565b90503590610114826101c0565b91906040838203126101025761012a906101f68185610168565b936020016101cf565b34610102576102186102123660046101dc565b90610532565b604051005b34610102576102186102303660046101dc565b906105e2565b91906040838203126101025761012a906101f681856101cf565b3461010257610218610263366004610236565b906109d3565b600091031261010257565b3461010257610284366004610269565b61015e7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a5986101a4565b34610102576102bd366004610269565b61015e7f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db6101a4565b346101025761015e61014c6102fc3660046101dc565b906104c0565b908160209103126101025790565b9060208282031261010257813567ffffffffffffffff81116101025761012a9201610302565b3461010257610218610349366004610310565b6110ca565b61012a61012a61012a9290565b61012a600061034e565b61012a61035b565b346101025761037d366004610269565b61015e6101a4610365565b346101025761021861039b3660046101dc565b906105d8565b906020828203126101025761012a916101cf565b34610102576102186103c83660046103a1565b610ae3565b7f7965db0b000000000000000000000000000000000000000000000000000000007fffffffff0000000000000000000000000000000000000000000000000000000082161490811561041d575090565b61012a91507fffffffff00000000000000000000000000000000000000000000000000000000167f01ffc9a7000000000000000000000000000000000000000000000000000000001490565b905b600052602052604060002090565b6101b461012a61012a926001600160a01b031690565b61012a90610479565b61012a9061048f565b9061046b90610498565b61012a905b60ff1690565b61012a90546104ab565b61012a916104da916104d460009182610469565b016104a1565b6104b6565b61012a9081565b61012a90546104df565b600161050961012a92610501600090565b506000610469565b016104e6565b9061011491610525610520826104f0565b61053c565b9061052f91610659565b50565b906101149161050f565b610114903390610564565b6001600160a01b0390911681526040810192916101149160200152565b9061057661057282846104c0565b1590565b61057e575050565b6105b961058a60405190565b9283927fe2517d3f00000000000000000000000000000000000000000000000000000000845260048401610547565b0390fd5b90610114916105ce610520826104f0565b9061052f916106d5565b90610114916105bd565b906105ec336101b4565b6001600160a01b038216036106045761052f916106d5565b6040517f6697b232000000000000000000000000000000000000000000000000000000008152600490fd5b9060ff905b9181191691161790565b9061064e61012a61065592151590565b825461062f565b9055565b61066661057283836104c0565b156106ce57600191610687836106828360006104d48782610469565b61063e565b33906106bd6106b76106b77f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9590565b92610498565b926106c760405190565b600090a490565b5050600090565b906106e081836104c0565b156106ce576106f9600061068283826104d48782610469565b33906107296106b76106b77ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9590565b9261073360405190565b600090a4600190565b61012a9060401c6104b0565b61012a905461073c565b61012a905b67ffffffffffffffff1690565b61012a9054610752565b61075761012a61012a9290565b9067ffffffffffffffff90610634565b61075761012a61012a9267ffffffffffffffff1690565b906107b261012a6106559261078b565b825461077b565b9068ff00000000000000009060401b610634565b906107dd61012a61065592151590565b82546107b9565b61012d9061076e565b60208101929161011491906107e4565b907ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00908161083661083061057283610748565b91610764565b936000926108438461076e565b67ffffffffffffffff87161480610965575b6001966108726108648961076e565b9167ffffffffffffffff1690565b14908161093d575b155b9081610934575b50610909576108ac91836108a38661089a8a61076e565b980197886107a2565b6108fa5761096c565b6108b557505050565b6108be916107cd565b6108f57fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2916108ec60405190565b918291826107ed565b0390a1565b61090487876107cd565b61096c565b6040517ff92ee8a9000000000000000000000000000000000000000000000000000000008152600490fd5b15905038610883565b905061087c61094b30610498565b3b61095c6109588861034e565b9190565b1491905061087a565b5082610855565b906109ac61052f926109858161098061035b565b610659565b507ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a598610659565b507f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db610659565b90610114916107fd565b61011490610a0a7ff16bb8781ef1311f8fe06747bcbe481e695502acdcb0cb8c03aa03899e39a59861053c565b610ace565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b8054821015610a5e57610a55600191600052602060002090565b91020190600090565b610a25565b91906008610634910291610a7d6001600160a01b03841b90565b921b90565b9190610a9361012a61065593610498565b908354610a63565b9081549168010000000000000000831015610ac95782610ac391600161011495018155610a3b565b90610a82565b610a0f565b61011490610add600191610498565b90610a9b565b610114906109dd565b61011490610b197f33dd54660937884a707404066945db647918933f71cc471efc6d6d0c3665d8db61053c565b610f59565b903590601e193682900301821215610102570180359067ffffffffffffffff82116101025760200191602082023603831361010257565b90601f01601f1916810190811067ffffffffffffffff821117610ac957604052565b90610114610b8460405190565b9283610b55565b67ffffffffffffffff8111610ac95760208091020190565b90610bb5610bb083610b8b565b610b77565b918252565b61012a6101c0610b77565b60209081808080808080808080808080610bdd610bba565b9e8f60008152016000815201600081520160008152016000815201600081520160008152016060815201600081520160008152016000815201600081520160008152016060905250565b61012a610bc5565b60005b828110610c3e57505050565b602090610c49610c27565b8184015201610c32565b90610114610c69610c6384610ba3565b93610b8b565b601f190160208401610c2f565b634e487b7160e01b600052601160045260246000fd5b6000198114610c9b5760010190565b610c76565b903590601e193682900301821215610102570180359067ffffffffffffffff8211610102576020019136829003831361010257565b90821015610a5e576020610cec9202810190610ca0565b9091565b90610cf9825190565b811015610a5e576020809102010190565b61012a916008021c6101b4565b9061012a9154610d0a565b60005b838110610d355750506000910152565b8181015183820152602001610d25565b610d66610d6f602093610d7993610d5a815190565b80835293849260200190565b95869101610d22565b601f01601f191690565b0190565b90610d9d610d96610d8c845190565b8084529260200190565b9260200190565b9060005b818110610dae5750505090565b909192610dd4610dcd60019286516001600160a01b0316815260200190565b9460200190565b929101610da1565b805160ff16825261012a916101a0610e626101c08301610e0160208601516020860152565b610e1060408601516040860152565b610e1f60608601516060860152565b610e2e60808601516080860152565b60a0858101516001600160a01b031690850152610e5060c086015160c0860152565b60e085015184820360e0860152610d45565b6101008085015160ff169084015292610e82610120820151610120850152565b610e93610140820151610140850152565b610ea4610160820151610160850152565b610eb5610180820151610180850152565b0151906101a0818403910152610d7d565b9061012a91610ddc565b90610ee6610edc835190565b8083529160200190565b9081610ef86020830284019460200190565b926000915b838310610f0c57505050505090565b90919293946020610f2f610f2883856001950387528951610ec6565b9760200190565b9301930191939290610efd565b602080825261012a92910190610ed0565b6040513d6000823e3d90fd5b610f6b610f668280610b1e565b905090565b90600091610f7b6109588461034e565b146110c65781810190610f99610f94610f668484610b1e565b610c53565b93610fa38461034e565b610fb361012a610f668686610b1e565b811015610ffb5780610fdb610fd5610ff693610fcf8888610b1e565b90610cd5565b90611155565b610fe58289610cf0565b52610ff08188610cf0565b50610c8c565b610fa3565b50929150506110098161034e565b600161101661012a825490565b8210156110bf5761103361102d8361103893610a3b565b90610d17565b610498565b9063630ac52c91803b156101025761105e92849161105560405190565b94859260e01b90565b82528183816110708b60048301610f3c565b03925af19182156110ba576110899261108e5750610c8c565b611009565b6110ad90843d86116110b3575b6110a58183610b55565b810190610269565b38610ff0565b503d61109b565b610f4d565b5050509050565b5050565b61011490610aec565b156110da57565b60405162461bcd60e51b815260206004820152601660248201527f456d707479207472616e73616374696f6e2064617461000000000000000000006044820152606490fd5b9190811015610a5e570190565b6104b061012a61012a9260ff1690565b61012a9060f81c61112c565b6104b061012a61012a9290565b9061115e610c27565b506000611177826111716109588461034e565b116110d3565b6111c46111be6111b961119361118c8561034e565b868861111f565b357fff000000000000000000000000000000000000000000000000000000000000001690565b61113c565b91611148565b60ff8216036111d7575061012a91611378565b6111e16001611148565b60ff8216036111f4575061012a9161156e565b6112086112016002611148565b9160ff1690565b036112165761012a9161171b565b60405162461bcd60e51b815260206004820152601c60248201527f556e737570706f72746564207472616e73616374696f6e2074797065000000006044820152606490fd5b67ffffffffffffffff8111610ac957602090601f01601f19160190565b90826000939282370152565b90929192611294610bb08261125b565b938185526020850190828401116101025761011492611278565b9080601f830112156101025781602061012a93359101611284565b60ff81166100fb565b90503590610114826112c9565b909161012082840312610102576112f68383610168565b926113048160208501610168565b926113128260408301610168565b9261132083606084016101cf565b9261132e8160808501610168565b9260a081013567ffffffffffffffff8111610102578261134f9183016112ae565b9261012a6113608460c085016112d2565b9361136e8160e08601610168565b9361010001610168565b9061141761140861140861103361012a94611391610c27565b506114086113c26113a0610c27565b9889936113ba60006113b181611148565b60ff1690870152565b8101906112df565b94929c95979d939060408b9c93989c019d60608c019b611413608082019a60a083019961140b60c085019660e08601956114086101008201946101406101208401930152565b52565b9060ff169052565b5252565b6001600160a01b03169052565b90939293848311610102578411610102578101920390565b9092919261144c610bb082610b8b565b938185526020808601920283019281841161010257915b8383106114705750505050565b6020809161147e84866101cf565b815201920191611463565b9080601f830112156101025781602061012a9335910161143c565b919061016083820312610102576114bb8184610168565b926114c98260208301610168565b926114d78360408401610168565b926114e58160608501610168565b926114f382608083016101cf565b926115018360a08401610168565b9260c083013567ffffffffffffffff811161010257816115229185016112ae565b9260e081013567ffffffffffffffff81116101025782611543918301611489565b9261012a6115558461010085016112d2565b93611564816101208601610168565b9361014001610168565b611576610c27565b5061157f610c27565b9180600161158c81611148565b60ff16855261159a9061034e565b906115a493611424565b81016115af916114a4565b9160208c9b9a989497929599969b019a60408d01998d6060810199608082019860a083019760c084019660e08501946101a081019361010082019361012083019261014001906115fc9152565b6116039152565b60ff16905252526116119152565b61161a90610498565b6001600160a01b031690525b61162d9152565b6116349152565b61163b9152565b61012a9152565b909161018082840312610102576116598383610168565b926116678160208501610168565b926116758260408301610168565b926116838360608401610168565b926116918160808501610168565b9261169f8260a083016101cf565b926116ad8360c08401610168565b9260e083013567ffffffffffffffff811161010257816116ce9185016112ae565b9261010081013567ffffffffffffffff811161010257826116f0918301611489565b9261012a6117028461012085016112d2565b93611711816101408601610168565b9361016001610168565b90611724610c27565b5061172d610c27565b91829161173a6002611148565b60ff1683528061174a600161034e565b9061175493611424565b810161175f91611642565b6101408d015260208c019b604081019a99610160820199909861018083019891976080840197919660a0850196939560c0860195929460e08401936101a08101926101008201929161012001906117b39152565b60ff16905252526117c19152565b6117ca90610498565b6001600160a01b03169052611626915256fea2646970667358221220b8ff3cf3247a12d054de35f2950617cc72e1a0c6e90f9eebcd928df1a78d874b64736f6c63430008140033",
}

// TransactionsAnalyzerABI is the input ABI used to generate the binding from.
// Deprecated: Use TransactionsAnalyzerMetaData.ABI instead.
var TransactionsAnalyzerABI = TransactionsAnalyzerMetaData.ABI

// TransactionsAnalyzerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransactionsAnalyzerMetaData.Bin instead.
var TransactionsAnalyzerBin = TransactionsAnalyzerMetaData.Bin

// DeployTransactionsAnalyzer deploys a new Ethereum contract, binding an instance of TransactionsAnalyzer to it.
func DeployTransactionsAnalyzer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TransactionsAnalyzer, error) {
	parsed, err := TransactionsAnalyzerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransactionsAnalyzerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransactionsAnalyzer{TransactionsAnalyzerCaller: TransactionsAnalyzerCaller{contract: contract}, TransactionsAnalyzerTransactor: TransactionsAnalyzerTransactor{contract: contract}, TransactionsAnalyzerFilterer: TransactionsAnalyzerFilterer{contract: contract}}, nil
}

// TransactionsAnalyzer is an auto generated Go binding around an Ethereum contract.
type TransactionsAnalyzer struct {
	TransactionsAnalyzerCaller     // Read-only binding to the contract
	TransactionsAnalyzerTransactor // Write-only binding to the contract
	TransactionsAnalyzerFilterer   // Log filterer for contract events
}

// TransactionsAnalyzerCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransactionsAnalyzerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionsAnalyzerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransactionsAnalyzerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionsAnalyzerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransactionsAnalyzerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransactionsAnalyzerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransactionsAnalyzerSession struct {
	Contract     *TransactionsAnalyzer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TransactionsAnalyzerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransactionsAnalyzerCallerSession struct {
	Contract *TransactionsAnalyzerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// TransactionsAnalyzerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransactionsAnalyzerTransactorSession struct {
	Contract     *TransactionsAnalyzerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// TransactionsAnalyzerRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransactionsAnalyzerRaw struct {
	Contract *TransactionsAnalyzer // Generic contract binding to access the raw methods on
}

// TransactionsAnalyzerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransactionsAnalyzerCallerRaw struct {
	Contract *TransactionsAnalyzerCaller // Generic read-only contract binding to access the raw methods on
}

// TransactionsAnalyzerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransactionsAnalyzerTransactorRaw struct {
	Contract *TransactionsAnalyzerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransactionsAnalyzer creates a new instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzer(address common.Address, backend bind.ContractBackend) (*TransactionsAnalyzer, error) {
	contract, err := bindTransactionsAnalyzer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzer{TransactionsAnalyzerCaller: TransactionsAnalyzerCaller{contract: contract}, TransactionsAnalyzerTransactor: TransactionsAnalyzerTransactor{contract: contract}, TransactionsAnalyzerFilterer: TransactionsAnalyzerFilterer{contract: contract}}, nil
}

// NewTransactionsAnalyzerCaller creates a new read-only instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzerCaller(address common.Address, caller bind.ContractCaller) (*TransactionsAnalyzerCaller, error) {
	contract, err := bindTransactionsAnalyzer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerCaller{contract: contract}, nil
}

// NewTransactionsAnalyzerTransactor creates a new write-only instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzerTransactor(address common.Address, transactor bind.ContractTransactor) (*TransactionsAnalyzerTransactor, error) {
	contract, err := bindTransactionsAnalyzer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerTransactor{contract: contract}, nil
}

// NewTransactionsAnalyzerFilterer creates a new log filterer instance of TransactionsAnalyzer, bound to a specific deployed contract.
func NewTransactionsAnalyzerFilterer(address common.Address, filterer bind.ContractFilterer) (*TransactionsAnalyzerFilterer, error) {
	contract, err := bindTransactionsAnalyzer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerFilterer{contract: contract}, nil
}

// bindTransactionsAnalyzer binds a generic wrapper to an already deployed contract.
func bindTransactionsAnalyzer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransactionsAnalyzerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransactionsAnalyzer *TransactionsAnalyzerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransactionsAnalyzer.Contract.TransactionsAnalyzerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransactionsAnalyzer *TransactionsAnalyzerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.TransactionsAnalyzerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransactionsAnalyzer *TransactionsAnalyzerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.TransactionsAnalyzerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransactionsAnalyzer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.DEFAULTADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.DEFAULTADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) EOAADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "EOA_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) EOAADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.EOAADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// EOAADMINROLE is a free data retrieval call binding the contract method 0x508a50f4.
//
// Solidity: function EOA_ADMIN_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) EOAADMINROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.EOAADMINROLE(&_TransactionsAnalyzer.CallOpts)
}

// HOOKCALLERROLE is a free data retrieval call binding the contract method 0x5f03a661.
//
// Solidity: function HOOK_CALLER_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) HOOKCALLERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "HOOK_CALLER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HOOKCALLERROLE is a free data retrieval call binding the contract method 0x5f03a661.
//
// Solidity: function HOOK_CALLER_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) HOOKCALLERROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.HOOKCALLERROLE(&_TransactionsAnalyzer.CallOpts)
}

// HOOKCALLERROLE is a free data retrieval call binding the contract method 0x5f03a661.
//
// Solidity: function HOOK_CALLER_ROLE() view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) HOOKCALLERROLE() ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.HOOKCALLERROLE(&_TransactionsAnalyzer.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.GetRoleAdmin(&_TransactionsAnalyzer.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TransactionsAnalyzer.Contract.GetRoleAdmin(&_TransactionsAnalyzer.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TransactionsAnalyzer.Contract.HasRole(&_TransactionsAnalyzer.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TransactionsAnalyzer.Contract.HasRole(&_TransactionsAnalyzer.CallOpts, role, account)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TransactionsAnalyzer.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TransactionsAnalyzer.Contract.SupportsInterface(&_TransactionsAnalyzer.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TransactionsAnalyzer *TransactionsAnalyzerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TransactionsAnalyzer.Contract.SupportsInterface(&_TransactionsAnalyzer.CallOpts, interfaceId)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) AddOnBlockEndCallback(opts *bind.TransactOpts, callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "addOnBlockEndCallback", callbackAddress)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) AddOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.AddOnBlockEndCallback(&_TransactionsAnalyzer.TransactOpts, callbackAddress)
}

// AddOnBlockEndCallback is a paid mutator transaction binding the contract method 0xee546fd8.
//
// Solidity: function addOnBlockEndCallback(address callbackAddress) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) AddOnBlockEndCallback(callbackAddress common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.AddOnBlockEndCallback(&_TransactionsAnalyzer.TransactOpts, callbackAddress)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.GrantRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.GrantRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address eoaAdmin, address authorizedCaller) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) Initialize(opts *bind.TransactOpts, eoaAdmin common.Address, authorizedCaller common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "initialize", eoaAdmin, authorizedCaller)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address eoaAdmin, address authorizedCaller) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) Initialize(eoaAdmin common.Address, authorizedCaller common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.Initialize(&_TransactionsAnalyzer.TransactOpts, eoaAdmin, authorizedCaller)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address eoaAdmin, address authorizedCaller) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) Initialize(eoaAdmin common.Address, authorizedCaller common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.Initialize(&_TransactionsAnalyzer.TransactOpts, eoaAdmin, authorizedCaller)
}

// OnBlock is a paid mutator transaction binding the contract method 0x9dbbcf8e.
//
// Solidity: function onBlock((bytes[]) _block) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) OnBlock(opts *bind.TransactOpts, _block TransactionsAnalyzerBlockTransactions) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "onBlock", _block)
}

// OnBlock is a paid mutator transaction binding the contract method 0x9dbbcf8e.
//
// Solidity: function onBlock((bytes[]) _block) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) OnBlock(_block TransactionsAnalyzerBlockTransactions) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.OnBlock(&_TransactionsAnalyzer.TransactOpts, _block)
}

// OnBlock is a paid mutator transaction binding the contract method 0x9dbbcf8e.
//
// Solidity: function onBlock((bytes[]) _block) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) OnBlock(_block TransactionsAnalyzerBlockTransactions) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.OnBlock(&_TransactionsAnalyzer.TransactOpts, _block)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RenounceRole(&_TransactionsAnalyzer.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RenounceRole(&_TransactionsAnalyzer.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RevokeRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TransactionsAnalyzer *TransactionsAnalyzerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TransactionsAnalyzer.Contract.RevokeRole(&_TransactionsAnalyzer.TransactOpts, role, account)
}

// TransactionsAnalyzerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerInitializedIterator struct {
	Event *TransactionsAnalyzerInitialized // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerInitialized)
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
		it.Event = new(TransactionsAnalyzerInitialized)
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
func (it *TransactionsAnalyzerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerInitialized represents a Initialized event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterInitialized(opts *bind.FilterOpts) (*TransactionsAnalyzerInitializedIterator, error) {

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerInitializedIterator{contract: _TransactionsAnalyzer.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerInitialized) (event.Subscription, error) {

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerInitialized)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseInitialized(log types.Log) (*TransactionsAnalyzerInitialized, error) {
	event := new(TransactionsAnalyzerInitialized)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionsAnalyzerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleAdminChangedIterator struct {
	Event *TransactionsAnalyzerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerRoleAdminChanged)
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
		it.Event = new(TransactionsAnalyzerRoleAdminChanged)
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
func (it *TransactionsAnalyzerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerRoleAdminChanged represents a RoleAdminChanged event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TransactionsAnalyzerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerRoleAdminChangedIterator{contract: _TransactionsAnalyzer.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerRoleAdminChanged)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseRoleAdminChanged(log types.Log) (*TransactionsAnalyzerRoleAdminChanged, error) {
	event := new(TransactionsAnalyzerRoleAdminChanged)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionsAnalyzerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleGrantedIterator struct {
	Event *TransactionsAnalyzerRoleGranted // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerRoleGranted)
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
		it.Event = new(TransactionsAnalyzerRoleGranted)
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
func (it *TransactionsAnalyzerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerRoleGranted represents a RoleGranted event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TransactionsAnalyzerRoleGrantedIterator, error) {

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

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerRoleGrantedIterator{contract: _TransactionsAnalyzer.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerRoleGranted)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseRoleGranted(log types.Log) (*TransactionsAnalyzerRoleGranted, error) {
	event := new(TransactionsAnalyzerRoleGranted)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TransactionsAnalyzerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleRevokedIterator struct {
	Event *TransactionsAnalyzerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TransactionsAnalyzerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransactionsAnalyzerRoleRevoked)
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
		it.Event = new(TransactionsAnalyzerRoleRevoked)
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
func (it *TransactionsAnalyzerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransactionsAnalyzerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransactionsAnalyzerRoleRevoked represents a RoleRevoked event raised by the TransactionsAnalyzer contract.
type TransactionsAnalyzerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TransactionsAnalyzerRoleRevokedIterator, error) {

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

	logs, sub, err := _TransactionsAnalyzer.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TransactionsAnalyzerRoleRevokedIterator{contract: _TransactionsAnalyzer.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TransactionsAnalyzerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _TransactionsAnalyzer.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransactionsAnalyzerRoleRevoked)
				if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_TransactionsAnalyzer *TransactionsAnalyzerFilterer) ParseRoleRevoked(log types.Log) (*TransactionsAnalyzerRoleRevoked, error) {
	event := new(TransactionsAnalyzerRoleRevoked)
	if err := _TransactionsAnalyzer.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
