import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import 'hardhat-deploy';
import * as hdnode from '@ethersproject/hdnode';
import * as abigen from './scripts/abigen';

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.9",
    settings: {
      optimizer: {
        enabled: true,
        runs: 1000,
      },
      remappings : ["@openzeppelin=node_modules/@openzeppelin"],
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
  networks: {
    localGeth: {
      url: "http://127.0.0.1:8025",
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l1/' ],
      companionNetworks: {
        layer1: 'localObscuro'
      },
    },
    localObscuro: {
      url: "http://127.0.0.1:13001",
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l2/' ],
      companionNetworks: {
        layer1: 'localGeth'
      },
    },
    hardhat : {
      deploy : [ 'deploy_l1/' ],
      live: false,
      saveDeployments: false,
    }
  },
  namedAccounts: {
    deployer: { // Addressed used for deploying.
        default: 0,
    },
    sequncer:{ // For management contract.
        default: 1,
    },
    relayer: { // Address that will relay messages to enable boostrapping deposits.
        default: 2,
    }
  }
};

export default config;
