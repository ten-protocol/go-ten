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
	_ = abi.ConvertType
)

// ContractTransparencyConfigEventLogConfig is an auto generated low-level Go binding around an user-defined struct.
type ContractTransparencyConfigEventLogConfig struct {
	EventSignature  []byte
	IsPublic        bool
	Topic1CanView   bool
	Topic2CanView   bool
	Topic3CanView   bool
	VisibleToSender bool
}

// ContractTransparencyConfigVisibilityConfig is an auto generated low-level Go binding around an user-defined struct.
type ContractTransparencyConfigVisibilityConfig struct {
	IsTransparent   bool
	EventLogConfigs []ContractTransparencyConfigEventLogConfig
}

// ObsERC20MetaData contains all meta data concerning the ObsERC20 contract.
var ObsERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"busAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"visibilityRules\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isTransparent\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"eventSignature\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"isPublic\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"topic1CanView\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"topic2CanView\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"topic3CanView\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"visibleToSender\",\"type\":\"bool\"}],\"internalType\":\"structContractTransparencyConfig.EventLogConfig[]\",\"name\":\"eventLogConfigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structContractTransparencyConfig.VisibilityConfig\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405273deb34a740eca1ec42c8b8204cbec0ba34fdd27f3600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503480156200006657600080fd5b5060405162002b6d38038062002b6d83398181016040528101906200008c9190620006b1565b838381600390816200009f9190620009a2565b508060049081620000b19190620009a2565b505050620000c633836200011160201b60201c565b80600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050505062000e15565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000183576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200017a9062000aea565b60405180910390fd5b62000197600083836200028960201b60201c565b8060026000828254620001ab919062000b3b565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825462000202919062000b3b565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405162000269919062000ba9565b60405180910390a362000285600083836200046f60201b60201c565b5050565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603156200046a57600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036200046957600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115620003d957620003d862000bc6565b5b85604051602001620003ec919062000c5f565b60405160208183030381529060405260006040518563ffffffff1660e01b81526004016200041e949392919062000d4a565b6020604051808303816000875af11580156200043e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000464919062000de3565b905050505b5b505050565b505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620004dd8262000492565b810181811067ffffffffffffffff82111715620004ff57620004fe620004a3565b5b80604052505050565b60006200051462000474565b9050620005228282620004d2565b919050565b600067ffffffffffffffff821115620005455762000544620004a3565b5b620005508262000492565b9050602081019050919050565b60005b838110156200057d57808201518184015260208101905062000560565b838111156200058d576000848401525b50505050565b6000620005aa620005a48462000527565b62000508565b905082815260208101848484011115620005c957620005c86200048d565b5b620005d68482856200055d565b509392505050565b600082601f830112620005f657620005f562000488565b5b81516200060884826020860162000593565b91505092915050565b6000819050919050565b620006268162000611565b81146200063257600080fd5b50565b60008151905062000646816200061b565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600062000679826200064c565b9050919050565b6200068b816200066c565b81146200069757600080fd5b50565b600081519050620006ab8162000680565b92915050565b60008060008060808587031215620006ce57620006cd6200047e565b5b600085015167ffffffffffffffff811115620006ef57620006ee62000483565b5b620006fd87828801620005de565b945050602085015167ffffffffffffffff81111562000721576200072062000483565b5b6200072f87828801620005de565b9350506040620007428782880162000635565b925050606062000755878288016200069a565b91505092959194509250565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620007b457607f821691505b602082108103620007ca57620007c96200076c565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620008347fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620007f5565b620008408683620007f5565b95508019841693508086168417925050509392505050565b6000819050919050565b6000620008836200087d620008778462000611565b62000858565b62000611565b9050919050565b6000819050919050565b6200089f8362000862565b620008b7620008ae826200088a565b84845462000802565b825550505050565b600090565b620008ce620008bf565b620008db81848462000894565b505050565b5b818110156200090357620008f7600082620008c4565b600181019050620008e1565b5050565b601f82111562000952576200091c81620007d0565b6200092784620007e5565b8101602085101562000937578190505b6200094f6200094685620007e5565b830182620008e0565b50505b505050565b600082821c905092915050565b6000620009776000198460080262000957565b1980831691505092915050565b600062000992838362000964565b9150826002028217905092915050565b620009ad8262000761565b67ffffffffffffffff811115620009c957620009c8620004a3565b5b620009d582546200079b565b620009e282828562000907565b600060209050601f83116001811462000a1a576000841562000a05578287015190505b62000a11858262000984565b86555062000a81565b601f19841662000a2a86620007d0565b60005b8281101562000a545784890151825560018201915060208501945060208101905062000a2d565b8683101562000a74578489015162000a70601f89168262000964565b8355505b6001600288020188555050505b505050505050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b600062000ad2601f8362000a89565b915062000adf8262000a9a565b602082019050919050565b6000602082019050818103600083015262000b058162000ac3565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600062000b488262000611565b915062000b558362000611565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111562000b8d5762000b8c62000b0c565b5b828201905092915050565b62000ba38162000611565b82525050565b600060208201905062000bc0600083018462000b98565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b62000c00816200066c565b82525050565b62000c118162000611565b82525050565b60608201600082015162000c2f600085018262000bf5565b50602082015162000c44602085018262000bf5565b50604082015162000c59604085018262000c06565b50505050565b600060608201905062000c76600083018462000c17565b92915050565b600063ffffffff82169050919050565b62000c978162000c7c565b82525050565b600081519050919050565b600082825260208201905092915050565b600062000cc68262000c9d565b62000cd2818562000ca8565b935062000ce48185602086016200055d565b62000cef8162000492565b840191505092915050565b6000819050919050565b600060ff82169050919050565b600062000d3262000d2c62000d268462000cfa565b62000858565b62000d04565b9050919050565b62000d448162000d11565b82525050565b600060808201905062000d61600083018762000c8c565b62000d70602083018662000c8c565b818103604083015262000d84818562000cb9565b905062000d95606083018462000d39565b95945050505050565b600067ffffffffffffffff82169050919050565b62000dbd8162000d9e565b811462000dc957600080fd5b50565b60008151905062000ddd8162000db2565b92915050565b60006020828403121562000dfc5762000dfb6200047e565b5b600062000e0c8482850162000dcc565b91505092915050565b611d488062000e256000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c80633950935111610071578063395093511461019157806370a08231146101c157806395d89b41146101f1578063a457c2d71461020f578063a9059cbb1461023f578063dd62ed3e1461026f576100b4565b806306fdde03146100b9578063095ea7b3146100d757806318160ddd1461010757806323b872dd1461012557806330173dd114610155578063313ce56714610173575b600080fd5b6100c161029f565b6040516100ce9190611088565b60405180910390f35b6100f160048036038101906100ec9190611143565b610331565b6040516100fe919061119e565b60405180910390f35b61010f610354565b60405161011c91906111c8565b60405180910390f35b61013f600480360381019061013a91906111e3565b61035e565b60405161014c919061119e565b60405180910390f35b61015d61038d565b60405161016a9190611422565b60405180910390f35b61017b61049c565b6040516101889190611460565b60405180910390f35b6101ab60048036038101906101a69190611143565b6104a5565b6040516101b8919061119e565b60405180910390f35b6101db60048036038101906101d6919061147b565b6104dc565b6040516101e891906111c8565b60405180910390f35b6101f96105a4565b6040516102069190611088565b60405180910390f35b61022960048036038101906102249190611143565b610636565b604051610236919061119e565b60405180910390f35b61025960048036038101906102549190611143565b6106ad565b604051610266919061119e565b60405180910390f35b610289600480360381019061028491906114a8565b6106d0565b60405161029691906111c8565b60405180910390f35b6060600380546102ae90611517565b80601f01602080910402602001604051908101604052809291908181526020018280546102da90611517565b80156103275780601f106102fc57610100808354040283529160200191610327565b820191906000526020600020905b81548152906001019060200180831161030a57829003601f168201915b5050505050905090565b60008061033c610809565b9050610349818585610811565b600191505092915050565b6000600254905090565b600080610369610809565b90506103768582856109da565b610381858585610a66565b60019150509392505050565b610395610f93565b6000600167ffffffffffffffff8111156103b2576103b1611548565b5b6040519080825280602002602001820160405280156103eb57816020015b6103d8610faf565b8152602001906001900390816103d05790505b5090506040518060c001604052806040518060400160405280602081526020017fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8152508152602001600015158152602001600115158152602001600115158152602001600015158152602001600015158152508160008151811061047357610472611577565b5b602002602001018190525060405180604001604052806000151581526020018281525091505090565b60006012905090565b6000806104b0610809565b90506104d18185856104c285896106d0565b6104cc91906115d5565b610811565b600191505092915050565b60008173ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff16036105215761051a82610ce5565b905061059f565b8173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16036105645761055d82610ce5565b905061059f565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059690611677565b60405180910390fd5b919050565b6060600480546105b390611517565b80601f01602080910402602001604051908101604052809291908181526020018280546105df90611517565b801561062c5780601f106106015761010080835404028352916020019161062c565b820191906000526020600020905b81548152906001019060200180831161060f57829003601f168201915b5050505050905090565b600080610641610809565b9050600061064f82866106d0565b905083811015610694576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161068b90611709565b60405180910390fd5b6106a18286868403610811565b60019250505092915050565b6000806106b8610809565b90506106c5818585610a66565b600191505092915050565b60008273ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff16148061073757508173ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff16145b1561074d576107468383610d2d565b9050610803565b8273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806107b257508173ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b156107c8576107c18383610d2d565b9050610803565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107fa9061179b565b60405180910390fd5b92915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610880576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108779061182d565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036108ef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108e6906118bf565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516109cd91906111c8565b60405180910390a3505050565b60006109e684846106d0565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610a605781811015610a52576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a499061192b565b60405180910390fd5b610a5f8484848403610811565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610ad5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610acc906119bd565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610b44576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b3b90611a4f565b60405180910390fd5b610b4f838383610db4565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610bd5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bcc90611ae1565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610c6891906115d5565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610ccc91906111c8565b60405180910390a3610cdf848484610f8e565b50505050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160315610f8957600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610f8857600060405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381525090506000600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b1454caa43600180811115610eff57610efe611b01565b5b85604051602001610f109190611b90565b60405160208183030381529060405260006040518563ffffffff1660e01b8152600401610f409493929190611c59565b6020604051808303816000875af1158015610f5f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f839190611ce5565b905050505b5b505050565b505050565b6040518060400160405280600015158152602001606081525090565b6040518060c00160405280606081526020016000151581526020016000151581526020016000151581526020016000151581526020016000151581525090565b600081519050919050565b600082825260208201905092915050565b60005b8381101561102957808201518184015260208101905061100e565b83811115611038576000848401525b50505050565b6000601f19601f8301169050919050565b600061105a82610fef565b6110648185610ffa565b935061107481856020860161100b565b61107d8161103e565b840191505092915050565b600060208201905081810360008301526110a2818461104f565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006110da826110af565b9050919050565b6110ea816110cf565b81146110f557600080fd5b50565b600081359050611107816110e1565b92915050565b6000819050919050565b6111208161110d565b811461112b57600080fd5b50565b60008135905061113d81611117565b92915050565b6000806040838503121561115a576111596110aa565b5b6000611168858286016110f8565b92505060206111798582860161112e565b9150509250929050565b60008115159050919050565b61119881611183565b82525050565b60006020820190506111b3600083018461118f565b92915050565b6111c28161110d565b82525050565b60006020820190506111dd60008301846111b9565b92915050565b6000806000606084860312156111fc576111fb6110aa565b5b600061120a868287016110f8565b935050602061121b868287016110f8565b925050604061122c8682870161112e565b9150509250925092565b61123f81611183565b82525050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b600061129882611271565b6112a2818561127c565b93506112b281856020860161100b565b6112bb8161103e565b840191505092915050565b600060c08301600083015184820360008601526112e3828261128d565b91505060208301516112f86020860182611236565b50604083015161130b6040860182611236565b50606083015161131e6060860182611236565b5060808301516113316080860182611236565b5060a083015161134460a0860182611236565b508091505092915050565b600061135b83836112c6565b905092915050565b6000602082019050919050565b600061137b82611245565b6113858185611250565b93508360208202850161139785611261565b8060005b858110156113d357848403895281516113b4858261134f565b94506113bf83611363565b925060208a0199505060018101905061139b565b50829750879550505050505092915050565b60006040830160008301516113fd6000860182611236565b50602083015184820360208601526114158282611370565b9150508091505092915050565b6000602082019050818103600083015261143c81846113e5565b905092915050565b600060ff82169050919050565b61145a81611444565b82525050565b60006020820190506114756000830184611451565b92915050565b600060208284031215611491576114906110aa565b5b600061149f848285016110f8565b91505092915050565b600080604083850312156114bf576114be6110aa565b5b60006114cd858286016110f8565b92505060206114de858286016110f8565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061152f57607f821691505b602082108103611542576115416114e8565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006115e08261110d565b91506115eb8361110d565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156116205761161f6115a6565b5b828201905092915050565b7f4e6f7420616c6c6f77656420746f2072656164207468652062616c616e636500600082015250565b6000611661601f83610ffa565b915061166c8261162b565b602082019050919050565b6000602082019050818103600083015261169081611654565b9050919050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b60006116f3602583610ffa565b91506116fe82611697565b604082019050919050565b60006020820190508181036000830152611722816116e6565b9050919050565b7f4e6f7420616c6c6f77656420746f20726561642074686520616c6c6f77616e6360008201527f6500000000000000000000000000000000000000000000000000000000000000602082015250565b6000611785602183610ffa565b915061179082611729565b604082019050919050565b600060208201905081810360008301526117b481611778565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b6000611817602483610ffa565b9150611822826117bb565b604082019050919050565b600060208201905081810360008301526118468161180a565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b60006118a9602283610ffa565b91506118b48261184d565b604082019050919050565b600060208201905081810360008301526118d88161189c565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b6000611915601d83610ffa565b9150611920826118df565b602082019050919050565b6000602082019050818103600083015261194481611908565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b60006119a7602583610ffa565b91506119b28261194b565b604082019050919050565b600060208201905081810360008301526119d68161199a565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b6000611a39602383610ffa565b9150611a44826119dd565b604082019050919050565b60006020820190508181036000830152611a6881611a2c565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b6000611acb602683610ffa565b9150611ad682611a6f565b604082019050919050565b60006020820190508181036000830152611afa81611abe565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b611b39816110cf565b82525050565b611b488161110d565b82525050565b606082016000820151611b646000850182611b30565b506020820151611b776020850182611b30565b506040820151611b8a6040850182611b3f565b50505050565b6000606082019050611ba56000830184611b4e565b92915050565b600063ffffffff82169050919050565b611bc481611bab565b82525050565b600082825260208201905092915050565b6000611be682611271565b611bf08185611bca565b9350611c0081856020860161100b565b611c098161103e565b840191505092915050565b6000819050919050565b6000819050919050565b6000611c43611c3e611c3984611c14565b611c1e565b611444565b9050919050565b611c5381611c28565b82525050565b6000608082019050611c6e6000830187611bbb565b611c7b6020830186611bbb565b8181036040830152611c8d8185611bdb565b9050611c9c6060830184611c4a565b95945050505050565b600067ffffffffffffffff82169050919050565b611cc281611ca5565b8114611ccd57600080fd5b50565b600081519050611cdf81611cb9565b92915050565b600060208284031215611cfb57611cfa6110aa565b5b6000611d0984828501611cd0565b9150509291505056fea26469706673582212203981d1ae057eaa4af2285ae435eff9df86a18e86fb08cdf7753394717c9a5d6964736f6c634300080f0033",
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
	parsed, err := ObsERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// VisibilityRules is a free data retrieval call binding the contract method 0x30173dd1.
//
// Solidity: function visibilityRules() pure returns((bool,(bytes,bool,bool,bool,bool,bool)[]))
func (_ObsERC20 *ObsERC20Caller) VisibilityRules(opts *bind.CallOpts) (ContractTransparencyConfigVisibilityConfig, error) {
	var out []interface{}
	err := _ObsERC20.contract.Call(opts, &out, "visibilityRules")

	if err != nil {
		return *new(ContractTransparencyConfigVisibilityConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(ContractTransparencyConfigVisibilityConfig)).(*ContractTransparencyConfigVisibilityConfig)

	return out0, err

}

// VisibilityRules is a free data retrieval call binding the contract method 0x30173dd1.
//
// Solidity: function visibilityRules() pure returns((bool,(bytes,bool,bool,bool,bool,bool)[]))
func (_ObsERC20 *ObsERC20Session) VisibilityRules() (ContractTransparencyConfigVisibilityConfig, error) {
	return _ObsERC20.Contract.VisibilityRules(&_ObsERC20.CallOpts)
}

// VisibilityRules is a free data retrieval call binding the contract method 0x30173dd1.
//
// Solidity: function visibilityRules() pure returns((bool,(bytes,bool,bool,bool,bool,bool)[]))
func (_ObsERC20 *ObsERC20CallerSession) VisibilityRules() (ContractTransparencyConfigVisibilityConfig, error) {
	return _ObsERC20.Contract.VisibilityRules(&_ObsERC20.CallOpts)
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
