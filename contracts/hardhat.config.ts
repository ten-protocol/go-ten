import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import 'hardhat-deploy';
import * as abigen from './scripts/abigen';
import './scripts/run-obscuro-node.ts';

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
    simGeth: {
      url: "http://127.0.0.1:32000",
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l1/' ],
      accounts: [ 'f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb']
    },
    simObscuro: {
      url: "http://127.0.0.1:3000",
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l2/' ],
      accounts: [ 'f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb']
    },
    localGeth: {
      url: "http://127.0.0.1:8025",
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l1/' ],
      accounts: [ 'f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb']
    },
    localObscuro: {
      url: "http://127.0.0.1:3000",
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l2/' ],
      companionNetworks: {
        layer1: 'localGeth'
      },
      accounts: [ 
        '6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682', 
        '4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8'
      ]
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
