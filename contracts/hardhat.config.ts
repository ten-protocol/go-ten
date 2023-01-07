import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import 'hardhat-deploy';
import * as abigen from './scripts/abigen';

import './scripts/obscuro-deploy';
import './scripts/run-wallet-extension';

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
      accounts: [ '8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b']
    },
    localGeth: {
      url: "http://127.0.0.1:8025",
      chainId: 1337,
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l1/' ],
      accounts: [ 'f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb']
    },
    localObscuro: {
      chainId: 777,
      url: "http://127.0.0.1:3000",
      obscuroEncRpcUrl: "ws://host.docker.internal:13001",
      gasPrice: 2000000000,
      live: false,
      saveDeployments: true,
      tags: ["local"],
      deploy: [ 'deploy_l2/' ],
      companionNetworks: {
        layer1: 'localGeth'
      },
      accounts: [
        `8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b`
      ]
    },
    hardhat : {
      deploy : [ 'deploy_l1/' ],
      chainId: 1337,
      live: false,
      saveDeployments: false,
      accounts:[ {
        privateKey: 'f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb',
        balance: "174165200000000000",
      }]
    }
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

export default config;
