import { HardhatUserConfig, task } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

import "hardhat-abi-exporter";
import "@solidstate/hardhat-bytecode-exporter";

import 'hardhat-deploy';

const abiExportPath = "./exported/";
const bytecodeExporterPath = "./exported/"

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
  namedAccounts: {
    deployer: { // Addressed used for deploying.
        default: 0,
    },
    sequncer:{ // For management contract.
        default: 1,
    },
  }
};

import { TASK_COMPILE } from "hardhat/builtin-tasks/task-names";
import * as path from "path";
import { HardhatPluginError } from 'hardhat/plugins';
import { spawn, spawnSync } from "node:child_process";
import * as fs from "fs";


task("generate-abi-bindings", "Using the evm bytecode and exported abi's of the contract export go bindings.")
.addFlag('noCompile', 'Don\'t compile before running this task')
.addParam('outputDir', 'Location to dump bindings')
.setAction(async function(args, hre) {
  if (!args.noCompile) {
    await hre.run(TASK_COMPILE);
  }

  const outputDirectory = path.resolve(hre.config.paths.root, args.outputDir);
  const bytecodeDirectory = path.resolve(hre.config.paths.root, bytecodeExporterPath)
  const abiDirectory = path.resolve(hre.config.paths.root, abiExportPath)


  if (outputDirectory === hre.config.paths.root) {
    throw new HardhatPluginError("AbiGen", 'resolved path must not be root directory');
  }
  const { bytecodeGroupConfig: config } = args;

  const fullNames = await hre.artifacts.getAllFullyQualifiedNames();

  await Promise.all(fullNames.map(async function (fullName) {
    
    let { bytecode, sourceName, contractName } = await hre.artifacts.readArtifact(fullName);

    // Some contracts like interfaces have only ABI.
    // This is enough to generate bindings, but since we haven't needed them so far we will skip those for now.
    bytecode = bytecode.replace(/^0x/, '');
    if (!bytecode.length) return;


    const contractBinDir = path.resolve(bytecodeDirectory, sourceName);
    const contractAbiDir = path.resolve(abiDirectory, sourceName);

    const binFilePath = path.resolve(contractBinDir, contractName + '.bin');
    const abiFilePath = path.resolve(contractAbiDir, contractName + '.json');
    
    const outputFileDir = path.resolve(outputDirectory, contractName + "/");
    const outputFilePath = path.resolve(outputFileDir, contractName + ".go");

    if (!fs.existsSync(outputDirectory)) {
      fs.mkdirSync(outputDirectory);
    }

    if (!fs.existsSync(abiFilePath)) {
      console.log(`No artifact for ${sourceName}`)
    } else {
      
      if (!fs.existsSync(outputFileDir)) {
        fs.mkdirSync(outputFileDir);
      }

      const abigenSpawn = spawnSync("abigen", [
        `--abi=${abiFilePath}`, 
        `--bin=${binFilePath}`, 
        `--pkg=${contractName}`, 
        `--out=${outputFilePath}`]
      );
    
      if (abigenSpawn.status == 0) {
        console.log(`Successfully generated go binding for ${sourceName}`);
      } else {
        console.log(`Error[${abigenSpawn.status}] generating go binding for ${sourceName};\n   Output: ${abigenSpawn.stderr}`);
      }
    }
  }));
})

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
