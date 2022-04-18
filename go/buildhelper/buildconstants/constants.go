package buildconstants

import "github.com/ethereum/go-ethereum/common"

const ADDR1PK = "f6b8ec4aed6edb5fd6087cbef005d4d2d89707228a4633499ad04fbeb99d3701"
const ADDR2PK = "15c1b11f3545f8f6e482b9656d2e163b7383c6fddb089a98ba3090733c5f1f3d"
const ADDR3PK = "4ce7bb228134edd313fd1041b1cd8d20bbc9c1a109d0e43e934d5cafda4d603a"

var (
	CONTRACT_ADDRESS = common.HexToAddress("0x7e6B032d74AA8afE44441987b6b3691adAee98FD")
)

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
