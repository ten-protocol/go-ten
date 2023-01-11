import { task } from "hardhat/config";
import { TASK_COMPILE } from "hardhat/builtin-tasks/task-names";
import * as path from "path";
import { HardhatPluginError } from 'hardhat/plugins';
import { spawnSync } from "node:child_process";
import * as fs from "fs";

export const abiExportPath = "./exported/";
export const bytecodeExporterPath = "./exported/"

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