import 'ten-hardhat-plugin';
import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

// Hardhat-deploy plugin - https://www.npmjs.com/package/hardhat-deploy
import 'hardhat-deploy';
// Hardhat upgrade plugin - https://www.npmjs.com/package/hardhat-upgrades
// import "@openzeppelin/hardhat-upgrades";
// Hardhat ignore warnings plugin - https://www.npmjs.com/package/hardhat-ignore-warnings
import 'hardhat-ignore-warnings';
import '@openzeppelin/hardhat-upgrades';

import * as abigen from './tasks/abigen';
import './tasks/obscuro-deploy';

import * as fs from "fs";

const config: HardhatUserConfig = {
  paths: {
    sources: "src"
  },
  solidity: {
    version: "0.8.28",
    settings: {
      viaIR: true,
      optimizer: {
        enabled: true,
        runs: 1000,
        details: {
          yulDetails: {
            optimizerSteps: "u",
          },
        },  
      },
      evmVersion: "cancun",
      outputSelection: { "*": { "*": [ "*" ], "": [ "*" ] } }
    },
  },
  abiExporter : {
    path: abigen.abiExportPath,
    runOnCompile: true,
    clear: true,
    format: "json",
  },
  bytecodeExporter : {
    path: abigen.bytecodeExporterPath,
    runOnCompile: true,
    clear: true,
  },
  namedAccounts: {
    deployer: { // Addressed used for deploying.
        default: 0,
    },
  },
  // For help configuring - https://www.npmjs.com/package/hardhat-ignore-warnings
  warnings : {
    '*' : {
      default: 'warn'
    },
    'src/testing/**/*': {
      default: 'off'
    }
  },
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY || "",
  }
};

try {
  config.networks = JSON.parse(fs.readFileSync('config/networks.json', 'utf8'));
} catch (e) {
  console.log(`Failed parsing "config/networks.json" with reason - ${e}`);
}

try {
  if (process.env.NETWORK_JSON != null) {
    config.networks = JSON.parse(process.env.NETWORK_JSON!!);
  }
} catch (e) {
  console.log(`Unable to parse networks configuration from environment variable. Reason is ${e}`);
}

export default config;
