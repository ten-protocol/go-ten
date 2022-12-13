import { subtask, task } from "hardhat/config";
import { TASK_COMPILE } from "hardhat/builtin-tasks/task-names";
import * as path from "path";
import { HardhatPluginError } from 'hardhat/plugins';

import * as dockerApi from 'node-docker-api';
import YAML from 'yaml'
import * as fs from 'fs';
import { on } from 'process';
import { NetworksUserConfig, NetworkUserConfig } from "hardhat/types";

import * as url from 'node:url';


task("run-wallet-extension", "Starts up the wallet extension docker container.")
.addParam('dockerImage', 
    'The docker image to use for wallet extension', 
    'testnetobscuronet.azurecr.io/obscuronet/walletextension')
.addParam('environment', "Which network to pick the node connection info from?", "default")
.setAction(async function(args, hre) {
    const docker = new dockerApi.Docker({ socketPath: '/var/run/docker.sock' });

    const defaultConfig = {
        port: 3000,
        nodeHost: "http://127.0.0.1",
        nodeHttpPort: "8025"
    }


    if (args.environment != "default") {
        const network : any = hre.userConfig.networks![args.environment]

        console.log(`Url for environment - ${network.url}`);

        const parsedUrl = url.parse(network.url)
        defaultConfig.nodeHost = parsedUrl.hostname!
        defaultConfig.nodeHttpPort = parsedUrl.port!
    }

    console.log("Starting container...")
    const container = await docker.container.create({
        Image: args.dockerImage,
        Cmd: [
            "--port=3000",
            "--portWS=3001",
            `--nodeHost=${defaultConfig.nodeHost}`,
            `--nodePortHTTP=${defaultConfig.nodeHttpPort}`
        ]
    })


    process.on('SIGINT', ()=>{
        container.stop();
    })
    
    await container.start();

    const stream:any = await container.logs({
        follow: true,
        stdout: true,
        stderr: true
    })

    stream.on('data', (info: any) => console.log(info.toString()))
    stream.on('error', (err: any) => console.log(err.toString()))
    
    await container.wait();

});