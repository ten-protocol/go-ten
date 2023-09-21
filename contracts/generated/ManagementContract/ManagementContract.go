// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ManagementContract

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

// StructsCrossChainMessage is an auto generated low-level Go binding around an user-defined struct.
type StructsCrossChainMessage struct {
	Sender           common.Address
	Sequence         uint64
	Nonce            uint32
	Topic            uint32
	Payload          []byte
	ConsistencyLevel uint8
}

// StructsHeaderCrossChainData is an auto generated low-level Go binding around an user-defined struct.
type StructsHeaderCrossChainData struct {
	Messages []StructsCrossChainMessage
}

// StructsMetaRollup is an auto generated low-level Go binding around an user-defined struct.
type StructsMetaRollup struct {
	Hash               [32]byte
	AggregatorID       common.Address
	LastSequenceNumber *big.Int
}

// ManagementContractMetaData contains all meta data concerning the ManagementContract contract.
var ManagementContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"messageBusAddress\",\"type\":\"address\"}],\"name\":\"LogManagementContractCreated\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"r\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"_rollupData\",\"type\":\"string\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"nonce\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"topic\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"consistencyLevel\",\"type\":\"uint8\"}],\"internalType\":\"structStructs.CrossChainMessage[]\",\"name\":\"messages\",\"type\":\"tuple[]\"}],\"internalType\":\"structStructs.HeaderCrossChainData\",\"name\":\"crossChainData\",\"type\":\"tuple\"}],\"name\":\"AddRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"Attested\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHostAddresses\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rollupHash\",\"type\":\"bytes32\"}],\"name\":\"GetRollupByHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"Hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"AggregatorID\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"LastSequenceNumber\",\"type\":\"uint256\"}],\"internalType\":\"structStructs.MetaRollup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregatorID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_hostAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_genesisAttestation\",\"type\":\"string\"}],\"name\":\"InitializeNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"IsWithdrawalAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestReport\",\"type\":\"string\"}],\"name\":\"RequestNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"attesterID\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"requesterID\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"attesterSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"responseSecret\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"hostAddress\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"verifyAttester\",\"type\":\"bool\"}],\"name\":\"RespondNetworkSecret\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSeqNo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"messageBus\",\"outputs\":[{\"internalType\":\"contractIMessageBus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060045534801561001557600080fd5b5060405161002290610096565b604051809103906000f08015801561003e573d6000803e3d6000fd5b50600680546001600160a01b0319166001600160a01b039290921691821790556040519081527fbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf9060200160405180910390a16100a3565b6110e58061159f83390190565b6114ed806100b26000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638fa0d05311610076578063a52f433c1161005b578063a52f433c14610222578063bbd79e1514610232578063e34fbfc81461024557600080fd5b80638fa0d053146101e4578063a1a227fa146101f757600080fd5b8063440c953b116100a7578063440c953b1461011d57806359a90071146101345780638236a7ba1461014957600080fd5b8063324ff866146100c357806343348b2f146100e1575b600080fd5b6100cb610258565b6040516100d89190610d03565b60405180910390f35b61010d6100ef366004610d92565b6001600160a01b031660009081526001602052604090205460ff1690565b60405190151581526020016100d8565b61012660045481565b6040519081526020016100d8565b610147610142366004610e9b565b610331565b005b6101b1610157366004610f42565b6040805160608082018352600080835260208084018290529284018190528481526005835283902083519182018452805480835260018201546001600160a01b031693830193909352600201549281019290925290911491565b60408051921515835281516020808501919091528201516001600160a01b031683820152015160608201526080016100d8565b6101476101f2366004610f5b565b6103b9565b60065461020a906001600160a01b031681565b6040516001600160a01b0390911681526020016100d8565b600354610100900460ff1661010d565b610147610240366004610fe2565b610453565b6101476102533660046110a8565b6105b6565b60606002805480602002602001604051908101604052809291908181526020016000905b8282101561032857838290600052602060002001805461029b906110ea565b80601f01602080910402602001604051908101604052809291908181526020018280546102c7906110ea565b80156103145780601f106102e957610100808354040283529160200191610314565b820191906000526020600020905b8154815290600101906020018083116102f757829003601f168201915b50505050508152602001906001019061027c565b50505050905090565b60035460ff161561034157600080fd5b60038054600160ff1991821681179092556001600160a01b038816600090815260208381526040822080549093168417909255600280549384018155905284516103b0927f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace0191860190610bca565b50505050505050565b600160006103cd6040870160208801610d92565b6001600160a01b0316815260208101919091526040016000205460ff1661043b5760405162461bcd60e51b815260206004820152601760248201527f61676772656761746f72206e6f7420617474657374656400000000000000000060448201526064015b60405180910390fd5b610444846105d5565b61044d8161060d565b50505050565b6001600160a01b03861660009081526001602052604090205460ff168061047957600080fd5b81156105495760006104af8888868860405160200161049b9493929190611125565b6040516020818303038152906040526106c7565b905060006104bd8288610702565b9050886001600160a01b0316816001600160a01b0316146105465760405162461bcd60e51b815260206004820152602c60248201527f63616c63756c61746564206164647265737320616e642061747465737465724960448201527f4420646f6e74206d6174636800000000000000000000000000000000000000006064820152608401610432565b50505b6001600160a01b03861660009081526001602081815260408320805460ff19168317905560028054928301815590925284516105ac927f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace90920191860190610bca565b5050505050505050565b3360009081526020819052604090206105d0908383610c4e565b505050565b8035600090815260056020526040902081906105f18282611181565b50506004546040820135111561060a5760408101356004555b50565b600061061982806111d8565b9050905060005b818110156105d0576006546001600160a01b0316639730886d61064385806111d8565b8481811061065357610653611222565b90506020028101906106659190611238565b60016040518363ffffffff1660e01b81526004016106849291906112f1565b600060405180830381600087803b15801561069e57600080fd5b505af11580156106b2573d6000803e3d6000fd5b50505050806106c0906113be565b9050610620565b60006106d38251610726565b826040516020016106e59291906113d9565b604051602081830303815290604052805190602001209050919050565b60008060006107118585610860565b9150915061071e816108d0565b509392505050565b60608161076657505060408051808201909152600181527f3000000000000000000000000000000000000000000000000000000000000000602082015290565b8160005b8115610790578061077a816113be565b91506107899050600a8361144a565b915061076a565b60008167ffffffffffffffff8111156107ab576107ab610df8565b6040519080825280601f01601f1916602001820160405280156107d5576020820181803683370190505b5090505b8415610858576107ea60018361145e565b91506107f7600a86611475565b610802906030611489565b60f81b81838151811061081757610817611222565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350610851600a8661144a565b94506107d9565b949350505050565b6000808251604114156108975760208301516040840151606085015160001a61088b87828585610a8b565b945094505050506108c9565b8251604014156108c157602083015160408401516108b6868383610b78565b9350935050506108c9565b506000905060025b9250929050565b60008160048111156108e4576108e46114a1565b14156108ed5750565b6001816004811115610901576109016114a1565b141561094f5760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e617475726500000000000000006044820152606401610432565b6002816004811115610963576109636114a1565b14156109b15760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e677468006044820152606401610432565b60038160048111156109c5576109c56114a1565b1415610a1e5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b6064820152608401610432565b6004816004811115610a3257610a326114a1565b141561060a5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b6064820152608401610432565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0831115610ac25750600090506003610b6f565b8460ff16601b14158015610ada57508460ff16601c14155b15610aeb5750600090506004610b6f565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015610b3f573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610b6857600060019250925050610b6f565b9150600090505b94509492505050565b6000807f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831681610bae60ff86901c601b611489565b9050610bbc87828885610a8b565b935093505050935093915050565b828054610bd6906110ea565b90600052602060002090601f016020900481019282610bf85760008555610c3e565b82601f10610c1157805160ff1916838001178555610c3e565b82800160010185558215610c3e579182015b82811115610c3e578251825591602001919060010190610c23565b50610c4a929150610cc2565b5090565b828054610c5a906110ea565b90600052602060002090601f016020900481019282610c7c5760008555610c3e565b82601f10610c955782800160ff19823516178555610c3e565b82800160010185558215610c3e579182015b82811115610c3e578235825591602001919060010190610ca7565b5b80821115610c4a5760008155600101610cc3565b60005b83811015610cf2578181015183820152602001610cda565b8381111561044d5750506000910152565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b82811015610d7057878503603f1901845281518051808752610d51818989018a8501610cd7565b601f01601f191695909501860194509285019290850190600101610d2a565b5092979650505050505050565b6001600160a01b038116811461060a57600080fd5b600060208284031215610da457600080fd5b8135610daf81610d7d565b9392505050565b60008083601f840112610dc857600080fd5b50813567ffffffffffffffff811115610de057600080fd5b6020830191508360208285010111156108c957600080fd5b634e487b7160e01b600052604160045260246000fd5b600082601f830112610e1f57600080fd5b813567ffffffffffffffff80821115610e3a57610e3a610df8565b604051601f8301601f19908116603f01168101908282118183101715610e6257610e62610df8565b81604052838152866020858801011115610e7b57600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008060008060008060808789031215610eb457600080fd5b8635610ebf81610d7d565b9550602087013567ffffffffffffffff80821115610edc57600080fd5b610ee88a838b01610db6565b90975095506040890135915080821115610f0157600080fd5b610f0d8a838b01610e0e565b94506060890135915080821115610f2357600080fd5b50610f3089828a01610db6565b979a9699509497509295939492505050565b600060208284031215610f5457600080fd5b5035919050565b60008060008084860360a0811215610f7257600080fd5b6060811215610f8057600080fd5b50849350606085013567ffffffffffffffff80821115610f9f57600080fd5b610fab88838901610db6565b90955093506080870135915080821115610fc457600080fd5b50850160208188031215610fd757600080fd5b939692955090935050565b60008060008060008060c08789031215610ffb57600080fd5b863561100681610d7d565b9550602087013561101681610d7d565b9450604087013567ffffffffffffffff8082111561103357600080fd5b61103f8a838b01610e0e565b9550606089013591508082111561105557600080fd5b6110618a838b01610e0e565b9450608089013591508082111561107757600080fd5b5061108489828a01610e0e565b92505060a0870135801515811461109a57600080fd5b809150509295509295509295565b600080602083850312156110bb57600080fd5b823567ffffffffffffffff8111156110d257600080fd5b6110de85828601610db6565b90969095509350505050565b600181811c908216806110fe57607f821691505b6020821081141561111f57634e487b7160e01b600052602260045260246000fd5b50919050565b60006bffffffffffffffffffffffff19808760601b168352808660601b16601484015250835161115c816028850160208801610cd7565b835190830190611173816028840160208801610cd7565b016028019695505050505050565b8135815560018101602083013561119781610d7d565b6001600160a01b0381167fffffffffffffffffffffffff00000000000000000000000000000000000000008354161782555050604082013560028201555050565b6000808335601e198436030181126111ef57600080fd5b83018035915067ffffffffffffffff82111561120a57600080fd5b6020019150600581901b36038213156108c957600080fd5b634e487b7160e01b600052603260045260246000fd5b6000823560be1983360301811261124e57600080fd5b9190910192915050565b803563ffffffff8116811461126c57600080fd5b919050565b6000808335601e1984360301811261128857600080fd5b830160208101925035905067ffffffffffffffff8111156112a857600080fd5b8036038313156108c957600080fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b803560ff8116811461126c57600080fd5b604081526000833561130281610d7d565b6001600160a01b03166040830152602084013567ffffffffffffffff811680821461132c57600080fd5b60608401525061133e60408501611258565b63ffffffff16608083015261135560608501611258565b63ffffffff1660a083015261136d6080850185611271565b60c080850152611382610100850182846112b7565b91505061139160a086016112e0565b60ff1660e084015260209092019290925292915050565b634e487b7160e01b600052601160045260246000fd5b60006000198214156113d2576113d26113a8565b5060010190565b7f19457468657265756d205369676e6564204d6573736167653a0a00000000000081526000835161141181601a850160208801610cd7565b83519083019061142881601a840160208801610cd7565b01601a01949350505050565b634e487b7160e01b600052601260045260246000fd5b60008261145957611459611434565b500490565b600082821015611470576114706113a8565b500390565b60008261148457611484611434565b500690565b6000821982111561149c5761149c6113a8565b500190565b634e487b7160e01b600052602160045260246000fdfea26469706673582212205288dbcbf0c23f4cd66f125292c721cab0d30ab2a175d6c4984aa96b0d35fdd864736f6c63430008090033608060405234801561001057600080fd5b5061001a3361001f565b61006f565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6110678061007e6000396000f3fe60806040526004361061009a5760003560e01c80638da5cb5b1161006957806399a3ad211161004e57806399a3ad2114610235578063b1454caa14610255578063f2fde38b1461028e57610113565b80638da5cb5b146101ed5780639730886d1461021557610113565b80630fcfbd111461016057806333a88c7214610193578063346633fb146101c3578063715018a6146101d857610113565b36610113576040517f346633fb0000000000000000000000000000000000000000000000000000000081523360048201523460248201819052309163346633fb91906044016000604051808303818588803b1580156100f857600080fd5b505af115801561010c573d6000803e3d6000fd5b5050505050005b60405162461bcd60e51b815260206004820152600b60248201527f756e737570706f7274656400000000000000000000000000000000000000000060448201526064015b60405180910390fd5b34801561016c57600080fd5b5061018061017b36600461099d565b6102ae565b6040519081526020015b60405180910390f35b34801561019f57600080fd5b506101b36101ae36600461099d565b610364565b604051901515815260200161018a565b6101d66101d13660046109e7565b6103b7565b005b3480156101e457600080fd5b506101d6610483565b3480156101f957600080fd5b506000546040516001600160a01b03909116815260200161018a565b34801561022157600080fd5b506101d6610230366004610a13565b6104e9565b34801561024157600080fd5b506101d66102503660046109e7565b61068d565b34801561026157600080fd5b50610275610270366004610a89565b61078f565b60405167ffffffffffffffff909116815260200161018a565b34801561029a57600080fd5b506101d66102a9366004610b36565b6107e8565b600080826040516020016102c29190610b92565b60408051601f198184030181529181528151602092830120600081815260019093529120549091508061035d5760405162461bcd60e51b815260206004820152602160248201527f54686973206d65737361676520776173206e65766572207375626d697474656460448201527f2e000000000000000000000000000000000000000000000000000000000000006064820152608401610157565b9392505050565b600080826040516020016103789190610b92565b60408051601f1981840301815291815281516020928301206000818152600190935291205490915080158015906103af5750428111155b949350505050565b6000341180156103c657508034145b6104385760405162461bcd60e51b815260206004820152603060248201527f417474656d7074696e6720746f2073656e642076616c756520776974686f757460448201527f2070726f766964696e67204574686572000000000000000000000000000000006064820152608401610157565b604080513381526001600160a01b0384166020820152348183015290517ff1365f826a788d6c1a955db0eed5ba8642674219c4771f8c65918617511a15609181900360600190a15050565b6000546001600160a01b031633146104dd5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610157565b6104e760006108ca565b565b6000546001600160a01b031633146105435760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610157565b600061054f8242610c92565b90506000836040516020016105649190610b92565b60408051601f19818403018152918152815160209283012060008181526001909352912054909150156105ff5760405162461bcd60e51b815260206004820152602160248201527f4d657373616765207375626d6974746564206d6f7265207468616e206f6e636560448201527f21000000000000000000000000000000000000000000000000000000000000006064820152608401610157565b600081815260016020908152604082208490556002919061062290870187610b36565b6001600160a01b03168152602081019190915260400160009081209061064e6080870160608801610caa565b63ffffffff168152602080820192909252604001600090812080546001810182559082529190208591600402016106858282610e87565b505050505050565b6000546001600160a01b031633146106e75760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610157565b6000826001600160a01b03168260405160006040518083038185875af1925050503d8060008114610734576040519150601f19603f3d011682016040523d82523d6000602084013e610739565b606091505b505090508061078a5760405162461bcd60e51b815260206004820152601460248201527f6661696c65642073656e64696e672076616c75650000000000000000000000006044820152606401610157565b505050565b600061079a33610927565b90507fb93c37389233beb85a3a726c3f15c2d15533ee74cb602f20f490dfffef775937338288888888886040516107d79796959493929190610fa5565b60405180910390a195945050505050565b6000546001600160a01b031633146108425760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610157565b6001600160a01b0381166108be5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610157565b6108c7816108ca565b50565b600080546001600160a01b0383811673ffffffffffffffffffffffffffffffffffffffff19831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6001600160a01b0381166000908152600360205260408120805467ffffffffffffffff16916001919061095a8385611005565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550919050565b600060c0828403121561099757600080fd5b50919050565b6000602082840312156109af57600080fd5b813567ffffffffffffffff8111156109c657600080fd5b6103af84828501610985565b6001600160a01b03811681146108c757600080fd5b600080604083850312156109fa57600080fd5b8235610a05816109d2565b946020939093013593505050565b60008060408385031215610a2657600080fd5b823567ffffffffffffffff811115610a3d57600080fd5b610a4985828601610985565b95602094909401359450505050565b63ffffffff811681146108c757600080fd5b60ff811681146108c757600080fd5b8035610a8481610a6a565b919050565b600080600080600060808688031215610aa157600080fd5b8535610aac81610a58565b94506020860135610abc81610a58565b9350604086013567ffffffffffffffff80821115610ad957600080fd5b818801915088601f830112610aed57600080fd5b813581811115610afc57600080fd5b896020828501011115610b0e57600080fd5b6020830195508094505050506060860135610b2881610a6a565b809150509295509295909350565b600060208284031215610b4857600080fd5b813561035d816109d2565b67ffffffffffffffff811681146108c757600080fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6020815260008235610ba3816109d2565b6001600160a01b0381166020840152506020830135610bc181610b53565b67ffffffffffffffff808216604085015260408501359150610be282610a58565b63ffffffff808316606086015260608601359250610bff83610a58565b80831660808601525060808501359150601e19853603018212610c2157600080fd5b90840190813581811115610c3457600080fd5b803603861315610c4357600080fd5b60c060a0860152610c5b60e086018260208601610b69565b92505050610c6b60a08501610a79565b60ff811660c0850152509392505050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610ca557610ca5610c7c565b500190565b600060208284031215610cbc57600080fd5b813561035d81610a58565b60008135610cd481610a58565b92915050565b6000808335601e19843603018112610cf157600080fd5b83018035915067ffffffffffffffff821115610d0c57600080fd5b602001915036819003821315610d2157600080fd5b9250929050565b634e487b7160e01b600052604160045260246000fd5b600181811c90821680610d5257607f821691505b6020821081141561099757634e487b7160e01b600052602260045260246000fd5b601f82111561078a57600081815260208120601f850160051c81016020861015610d9a5750805b601f850160051c820191505b8181101561068557828155600101610da6565b67ffffffffffffffff831115610dd157610dd1610d28565b610de583610ddf8354610d3e565b83610d73565b6000601f841160018114610e195760008515610e015750838201355b600019600387901b1c1916600186901b178355610e73565b600083815260209020601f19861690835b82811015610e4a5786850135825560209485019460019092019101610e2a565b5086821015610e675760001960f88860031b161c19848701351681555b505060018560011b0183555b5050505050565b60008135610cd481610a6a565b8135610e92816109d2565b6001600160a01b038116905081548173ffffffffffffffffffffffffffffffffffffffff1982161783556020840135610eca81610b53565b7bffffffffffffffff00000000000000000000000000000000000000008160a01b1690507fffffffff0000000000000000000000000000000000000000000000000000000081848285161717855560408601359250610f2883610a58565b921760e09190911b909116178155610f60610f4560608401610cc7565b6001830163ffffffff821663ffffffff198254161781555050565b610f6d6080830183610cda565b610f7b818360028601610db9565b5050610fa1610f8c60a08401610e7a565b6003830160ff821660ff198254161781555050565b5050565b6001600160a01b038816815267ffffffffffffffff87166020820152600063ffffffff808816604084015280871660608401525060c06080830152610fee60c083018587610b69565b905060ff831660a083015298975050505050505050565b600067ffffffffffffffff80831681851680830382111561102857611028610c7c565b0194935050505056fea2646970667358221220fe09a0d141463ff7d6782fdab79699e88422f37be6e4a724c8b6c27cf5a760d764736f6c63430008090033",
}

// ManagementContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ManagementContractMetaData.ABI instead.
var ManagementContractABI = ManagementContractMetaData.ABI

// ManagementContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ManagementContractMetaData.Bin instead.
var ManagementContractBin = ManagementContractMetaData.Bin

// DeployManagementContract deploys a new Ethereum contract, binding an instance of ManagementContract to it.
func DeployManagementContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ManagementContract, error) {
	parsed, err := ManagementContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ManagementContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ManagementContract{ManagementContractCaller: ManagementContractCaller{contract: contract}, ManagementContractTransactor: ManagementContractTransactor{contract: contract}, ManagementContractFilterer: ManagementContractFilterer{contract: contract}}, nil
}

// ManagementContract is an auto generated Go binding around an Ethereum contract.
type ManagementContract struct {
	ManagementContractCaller     // Read-only binding to the contract
	ManagementContractTransactor // Write-only binding to the contract
	ManagementContractFilterer   // Log filterer for contract events
}

// ManagementContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ManagementContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ManagementContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ManagementContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ManagementContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ManagementContractSession struct {
	Contract     *ManagementContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ManagementContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ManagementContractCallerSession struct {
	Contract *ManagementContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ManagementContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ManagementContractTransactorSession struct {
	Contract     *ManagementContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ManagementContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ManagementContractRaw struct {
	Contract *ManagementContract // Generic contract binding to access the raw methods on
}

// ManagementContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ManagementContractCallerRaw struct {
	Contract *ManagementContractCaller // Generic read-only contract binding to access the raw methods on
}

// ManagementContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ManagementContractTransactorRaw struct {
	Contract *ManagementContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewManagementContract creates a new instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContract(address common.Address, backend bind.ContractBackend) (*ManagementContract, error) {
	contract, err := bindManagementContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ManagementContract{ManagementContractCaller: ManagementContractCaller{contract: contract}, ManagementContractTransactor: ManagementContractTransactor{contract: contract}, ManagementContractFilterer: ManagementContractFilterer{contract: contract}}, nil
}

// NewManagementContractCaller creates a new read-only instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractCaller(address common.Address, caller bind.ContractCaller) (*ManagementContractCaller, error) {
	contract, err := bindManagementContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ManagementContractCaller{contract: contract}, nil
}

// NewManagementContractTransactor creates a new write-only instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ManagementContractTransactor, error) {
	contract, err := bindManagementContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ManagementContractTransactor{contract: contract}, nil
}

// NewManagementContractFilterer creates a new log filterer instance of ManagementContract, bound to a specific deployed contract.
func NewManagementContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ManagementContractFilterer, error) {
	contract, err := bindManagementContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ManagementContractFilterer{contract: contract}, nil
}

// bindManagementContract binds a generic wrapper to an already deployed contract.
func bindManagementContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ManagementContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagementContract *ManagementContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ManagementContract.Contract.ManagementContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagementContract *ManagementContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.Contract.ManagementContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagementContract *ManagementContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagementContract.Contract.ManagementContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ManagementContract *ManagementContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ManagementContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ManagementContract *ManagementContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ManagementContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ManagementContract *ManagementContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ManagementContract.Contract.contract.Transact(opts, method, params...)
}

// Attested is a free data retrieval call binding the contract method 0x43348b2f.
//
// Solidity: function Attested(address _addr) view returns(bool)
func (_ManagementContract *ManagementContractCaller) Attested(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "Attested", _addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Attested is a free data retrieval call binding the contract method 0x43348b2f.
//
// Solidity: function Attested(address _addr) view returns(bool)
func (_ManagementContract *ManagementContractSession) Attested(_addr common.Address) (bool, error) {
	return _ManagementContract.Contract.Attested(&_ManagementContract.CallOpts, _addr)
}

// Attested is a free data retrieval call binding the contract method 0x43348b2f.
//
// Solidity: function Attested(address _addr) view returns(bool)
func (_ManagementContract *ManagementContractCallerSession) Attested(_addr common.Address) (bool, error) {
	return _ManagementContract.Contract.Attested(&_ManagementContract.CallOpts, _addr)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_ManagementContract *ManagementContractCaller) GetHostAddresses(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "GetHostAddresses")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_ManagementContract *ManagementContractSession) GetHostAddresses() ([]string, error) {
	return _ManagementContract.Contract.GetHostAddresses(&_ManagementContract.CallOpts)
}

// GetHostAddresses is a free data retrieval call binding the contract method 0x324ff866.
//
// Solidity: function GetHostAddresses() view returns(string[])
func (_ManagementContract *ManagementContractCallerSession) GetHostAddresses() ([]string, error) {
	return _ManagementContract.Contract.GetHostAddresses(&_ManagementContract.CallOpts)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x8236a7ba.
//
// Solidity: function GetRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,address,uint256))
func (_ManagementContract *ManagementContractCaller) GetRollupByHash(opts *bind.CallOpts, rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "GetRollupByHash", rollupHash)

	if err != nil {
		return *new(bool), *new(StructsMetaRollup), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(StructsMetaRollup)).(*StructsMetaRollup)

	return out0, out1, err

}

// GetRollupByHash is a free data retrieval call binding the contract method 0x8236a7ba.
//
// Solidity: function GetRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,address,uint256))
func (_ManagementContract *ManagementContractSession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _ManagementContract.Contract.GetRollupByHash(&_ManagementContract.CallOpts, rollupHash)
}

// GetRollupByHash is a free data retrieval call binding the contract method 0x8236a7ba.
//
// Solidity: function GetRollupByHash(bytes32 rollupHash) view returns(bool, (bytes32,address,uint256))
func (_ManagementContract *ManagementContractCallerSession) GetRollupByHash(rollupHash [32]byte) (bool, StructsMetaRollup, error) {
	return _ManagementContract.Contract.GetRollupByHash(&_ManagementContract.CallOpts, rollupHash)
}

// IsWithdrawalAvailable is a free data retrieval call binding the contract method 0xa52f433c.
//
// Solidity: function IsWithdrawalAvailable() view returns(bool)
func (_ManagementContract *ManagementContractCaller) IsWithdrawalAvailable(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "IsWithdrawalAvailable")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWithdrawalAvailable is a free data retrieval call binding the contract method 0xa52f433c.
//
// Solidity: function IsWithdrawalAvailable() view returns(bool)
func (_ManagementContract *ManagementContractSession) IsWithdrawalAvailable() (bool, error) {
	return _ManagementContract.Contract.IsWithdrawalAvailable(&_ManagementContract.CallOpts)
}

// IsWithdrawalAvailable is a free data retrieval call binding the contract method 0xa52f433c.
//
// Solidity: function IsWithdrawalAvailable() view returns(bool)
func (_ManagementContract *ManagementContractCallerSession) IsWithdrawalAvailable() (bool, error) {
	return _ManagementContract.Contract.IsWithdrawalAvailable(&_ManagementContract.CallOpts)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_ManagementContract *ManagementContractCaller) LastBatchSeqNo(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "lastBatchSeqNo")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_ManagementContract *ManagementContractSession) LastBatchSeqNo() (*big.Int, error) {
	return _ManagementContract.Contract.LastBatchSeqNo(&_ManagementContract.CallOpts)
}

// LastBatchSeqNo is a free data retrieval call binding the contract method 0x440c953b.
//
// Solidity: function lastBatchSeqNo() view returns(uint256)
func (_ManagementContract *ManagementContractCallerSession) LastBatchSeqNo() (*big.Int, error) {
	return _ManagementContract.Contract.LastBatchSeqNo(&_ManagementContract.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_ManagementContract *ManagementContractCaller) MessageBus(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ManagementContract.contract.Call(opts, &out, "messageBus")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_ManagementContract *ManagementContractSession) MessageBus() (common.Address, error) {
	return _ManagementContract.Contract.MessageBus(&_ManagementContract.CallOpts)
}

// MessageBus is a free data retrieval call binding the contract method 0xa1a227fa.
//
// Solidity: function messageBus() view returns(address)
func (_ManagementContract *ManagementContractCallerSession) MessageBus() (common.Address, error) {
	return _ManagementContract.Contract.MessageBus(&_ManagementContract.CallOpts)
}

// AddRollup is a paid mutator transaction binding the contract method 0x8fa0d053.
//
// Solidity: function AddRollup((bytes32,address,uint256) r, string _rollupData, ((address,uint64,uint32,uint32,bytes,uint8)[]) crossChainData) returns()
func (_ManagementContract *ManagementContractTransactor) AddRollup(opts *bind.TransactOpts, r StructsMetaRollup, _rollupData string, crossChainData StructsHeaderCrossChainData) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "AddRollup", r, _rollupData, crossChainData)
}

// AddRollup is a paid mutator transaction binding the contract method 0x8fa0d053.
//
// Solidity: function AddRollup((bytes32,address,uint256) r, string _rollupData, ((address,uint64,uint32,uint32,bytes,uint8)[]) crossChainData) returns()
func (_ManagementContract *ManagementContractSession) AddRollup(r StructsMetaRollup, _rollupData string, crossChainData StructsHeaderCrossChainData) (*types.Transaction, error) {
	return _ManagementContract.Contract.AddRollup(&_ManagementContract.TransactOpts, r, _rollupData, crossChainData)
}

// AddRollup is a paid mutator transaction binding the contract method 0x8fa0d053.
//
// Solidity: function AddRollup((bytes32,address,uint256) r, string _rollupData, ((address,uint64,uint32,uint32,bytes,uint8)[]) crossChainData) returns()
func (_ManagementContract *ManagementContractTransactorSession) AddRollup(r StructsMetaRollup, _rollupData string, crossChainData StructsHeaderCrossChainData) (*types.Transaction, error) {
	return _ManagementContract.Contract.AddRollup(&_ManagementContract.TransactOpts, r, _rollupData, crossChainData)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x59a90071.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret, string _hostAddress, string _genesisAttestation) returns()
func (_ManagementContract *ManagementContractTransactor) InitializeNetworkSecret(opts *bind.TransactOpts, _aggregatorID common.Address, _initSecret []byte, _hostAddress string, _genesisAttestation string) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "InitializeNetworkSecret", _aggregatorID, _initSecret, _hostAddress, _genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x59a90071.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret, string _hostAddress, string _genesisAttestation) returns()
func (_ManagementContract *ManagementContractSession) InitializeNetworkSecret(_aggregatorID common.Address, _initSecret []byte, _hostAddress string, _genesisAttestation string) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, _aggregatorID, _initSecret, _hostAddress, _genesisAttestation)
}

// InitializeNetworkSecret is a paid mutator transaction binding the contract method 0x59a90071.
//
// Solidity: function InitializeNetworkSecret(address _aggregatorID, bytes _initSecret, string _hostAddress, string _genesisAttestation) returns()
func (_ManagementContract *ManagementContractTransactorSession) InitializeNetworkSecret(_aggregatorID common.Address, _initSecret []byte, _hostAddress string, _genesisAttestation string) (*types.Transaction, error) {
	return _ManagementContract.Contract.InitializeNetworkSecret(&_ManagementContract.TransactOpts, _aggregatorID, _initSecret, _hostAddress, _genesisAttestation)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_ManagementContract *ManagementContractTransactor) RequestNetworkSecret(opts *bind.TransactOpts, requestReport string) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "RequestNetworkSecret", requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_ManagementContract *ManagementContractSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _ManagementContract.Contract.RequestNetworkSecret(&_ManagementContract.TransactOpts, requestReport)
}

// RequestNetworkSecret is a paid mutator transaction binding the contract method 0xe34fbfc8.
//
// Solidity: function RequestNetworkSecret(string requestReport) returns()
func (_ManagementContract *ManagementContractTransactorSession) RequestNetworkSecret(requestReport string) (*types.Transaction, error) {
	return _ManagementContract.Contract.RequestNetworkSecret(&_ManagementContract.TransactOpts, requestReport)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_ManagementContract *ManagementContractTransactor) RespondNetworkSecret(opts *bind.TransactOpts, attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _ManagementContract.contract.Transact(opts, "RespondNetworkSecret", attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_ManagementContract *ManagementContractSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _ManagementContract.Contract.RespondNetworkSecret(&_ManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}

// RespondNetworkSecret is a paid mutator transaction binding the contract method 0xbbd79e15.
//
// Solidity: function RespondNetworkSecret(address attesterID, address requesterID, bytes attesterSig, bytes responseSecret, string hostAddress, bool verifyAttester) returns()
func (_ManagementContract *ManagementContractTransactorSession) RespondNetworkSecret(attesterID common.Address, requesterID common.Address, attesterSig []byte, responseSecret []byte, hostAddress string, verifyAttester bool) (*types.Transaction, error) {
	return _ManagementContract.Contract.RespondNetworkSecret(&_ManagementContract.TransactOpts, attesterID, requesterID, attesterSig, responseSecret, hostAddress, verifyAttester)
}

// ManagementContractLogManagementContractCreatedIterator is returned from FilterLogManagementContractCreated and is used to iterate over the raw logs and unpacked data for LogManagementContractCreated events raised by the ManagementContract contract.
type ManagementContractLogManagementContractCreatedIterator struct {
	Event *ManagementContractLogManagementContractCreated // Event containing the contract specifics and raw log

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
func (it *ManagementContractLogManagementContractCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ManagementContractLogManagementContractCreated)
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
		it.Event = new(ManagementContractLogManagementContractCreated)
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
func (it *ManagementContractLogManagementContractCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ManagementContractLogManagementContractCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ManagementContractLogManagementContractCreated represents a LogManagementContractCreated event raised by the ManagementContract contract.
type ManagementContractLogManagementContractCreated struct {
	MessageBusAddress common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterLogManagementContractCreated is a free log retrieval operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_ManagementContract *ManagementContractFilterer) FilterLogManagementContractCreated(opts *bind.FilterOpts) (*ManagementContractLogManagementContractCreatedIterator, error) {

	logs, sub, err := _ManagementContract.contract.FilterLogs(opts, "LogManagementContractCreated")
	if err != nil {
		return nil, err
	}
	return &ManagementContractLogManagementContractCreatedIterator{contract: _ManagementContract.contract, event: "LogManagementContractCreated", logs: logs, sub: sub}, nil
}

// WatchLogManagementContractCreated is a free log subscription operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_ManagementContract *ManagementContractFilterer) WatchLogManagementContractCreated(opts *bind.WatchOpts, sink chan<- *ManagementContractLogManagementContractCreated) (event.Subscription, error) {

	logs, sub, err := _ManagementContract.contract.WatchLogs(opts, "LogManagementContractCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ManagementContractLogManagementContractCreated)
				if err := _ManagementContract.contract.UnpackLog(event, "LogManagementContractCreated", log); err != nil {
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

// ParseLogManagementContractCreated is a log parse operation binding the contract event 0xbd726cf82ac9c3260b1495107182e336e0654b25c10915648c0cc15b2bb72cbf.
//
// Solidity: event LogManagementContractCreated(address messageBusAddress)
func (_ManagementContract *ManagementContractFilterer) ParseLogManagementContractCreated(log types.Log) (*ManagementContractLogManagementContractCreated, error) {
	event := new(ManagementContractLogManagementContractCreated)
	if err := _ManagementContract.contract.UnpackLog(event, "LogManagementContractCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
