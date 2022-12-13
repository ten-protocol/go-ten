import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import 'hardhat-deploy';
import * as hdnode from '@ethersproject/hdnode';
import './scripts/abigen';


const abiExportPath = "./artifacts/abi/";
const bytecodeExporterPath = "./artifacts/bin/"

const config: HardhatUserConfig = {
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
    path: abiExportPath,
    runOnCompile: true,
    clear: true,
    format: "json",
  },
  bytecodeExporter : {
    path: bytecodeExporterPath,
    runOnCompile: true,
    clear: true,
  },
};

export default config;
