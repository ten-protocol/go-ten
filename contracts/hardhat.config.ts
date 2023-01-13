import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import './tasks/wallet-extension';
import * as abigen from './tasks/abigen';

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
      remappings : ["@openzeppelin=node_modules/@openzeppelin"],
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
  }
};

export default config;
