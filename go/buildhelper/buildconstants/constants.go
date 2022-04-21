package buildconstants

import "github.com/ethereum/go-ethereum/common"

// TODO these will be loaded from a keyset
const (
	ADDR1PK = "5dbbff1b5ff19f1ad6ea656433be35f6846e890b3f3ec6ef2b2e2137a8cab4ae"
	ADDR2PK = "b728cd9a9f54cede03a82fc189eab4830a612703d48b7ef43ceed2cbad1a06c7"
	ADDR3PK = "1e1e76d5c0ea1382b6acf76e873977fd223c7fa2a6dc57db2b94e93eb303ba85"
)

// TODO this will be dynamic once the deployment of multiple enviroments is a feature
var CONTRACT_ADDRESS = common.HexToAddress("0x63eA4F258791234606F9d11C0C66D83c99bfC7CA")

const CONTRACT_ABI = `[
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "rollupData",
				"type": "string"
			}
		],
		"name": "AddRollup",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "dest",
				"type": "address"
			}
		],
		"name": "Deposit",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "inputSecret",
				"type": "string"
			}
		],
		"name": "StoreSecret",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "withdrawAmount",
				"type": "uint256"
			},
			{
				"internalType": "address payable",
				"name": "destination",
				"type": "address"
			}
		],
		"name": "Withdraw",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "deposits",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "RequestSecret",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "Rollup",
		"outputs": [
			{
				"internalType": "string[]",
				"name": "",
				"type": "string[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "rollups",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`
