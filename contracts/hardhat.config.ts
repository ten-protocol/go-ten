import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import 'hardhat-deploy';
import * as abigen from './scripts/abigen';

import * as fs from 'fs';

import './scripts/obscuro-deploy';
import './scripts/run-wallet-extension';

import * as dotenv from "dotenv";
dotenv.config({ path: __dirname+'/.env' });


let config: HardhatUserConfig = {
  paths: {
    sources: "src"
  },
  solidity: {
    version: "0.8.9",
    settings: {
      optimizer: {
        enabled: true,
        runs: 1000,
      }
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
    sequncer:{ // For management contract.
        default: 1,
    },
  }
};


try {
  config.networks = JSON.parse(fs.readFileSync('config/network.json', 'utf8'));
} catch (e) {
  console.log(e);
}

try {
  config.networks = JSON.parse(process.env.NETWORK_JSON!!);
} catch (e) {
  console.log(e);
}

export default config;
