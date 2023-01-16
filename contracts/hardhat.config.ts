import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import 'hardhat-deploy';

import './tasks/wallet-extension';
import * as abigen from './tasks/abigen';
import './tasks/obscuro-deploy';

import * as fs from "fs";

const config: HardhatUserConfig = {
  paths: {
    sources: "src"
  },
  solidity: {
    version: "0.8.9",
    settings: {
      optimizer: {
        enabled: true,
        runs: 1000,
      },
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
    sequencer:{ // For management contract.
        default: 1,
    },
    hocowner: {
        default: 2,
    },
    pocowner: {
        default: 3,
    },
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
