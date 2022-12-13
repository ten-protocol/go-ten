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
.addFlag('wait')
.addParam('dockerImage', 
    'The docker image to use for wallet extension', 
    'testnetobscuronet.azurecr.io/obscuronet/walletextension')
.addParam('rpcURL', "Which network to pick the node connection info from?")
.setAction(async function(args, hre) {
    const docker = new dockerApi.Docker({ socketPath: '/var/run/docker.sock' });

    const parsedUrl = url.parse(args.rpcURL)

    console.log("Starting container...")
    const container = await docker.container.create({
        Image: args.dockerImage,
        Cmd: [
            "--port=3000",
            "--portWS=3001",
            `--nodeHost=${parsedUrl.hostname}`,
            `--nodePortHTTP=${parsedUrl.port}`
        ],
        NetworkMode: 'host',
        ExposedPorts: { "3000/tcp": {} },
        PortBindings:  { "3000/tcp": [{ "HostPort": "3000" }] }
    })


    process.on('SIGINT', ()=>{
        container.stop();
    })
    
    await container.start();

    const stream: any = await container.logs({
        stdout: true,
        stderr: true,
    })
    
    await new Promise((resolveInner)=> {
        stream.on('data', (msg: any)=>console.log(msg.toString()))
        stream.on('error', (msg: any)=>console.log(msg.toString()))
        setTimeout(resolveInner, 30_000);
    })


    if (args.wait) {   

        await container.wait();
    }
});

task("stop-wallet-extension", "Starts up the wallet extension docker container.")
.addParam('dockerImage', 
    'The docker image to use for wallet extension', 
    'testnetobscuronet.azurecr.io/obscuronet/walletextension')
.setAction(async function(args, hre) {
    const docker = new dockerApi.Docker({ socketPath: '/var/run/docker.sock' });
    const containers = await docker.container.list();

    const container = containers.find((c)=> { 
       const data : any = c.data; 
       return data.Image == 'testnetobscuronet.azurecr.io/obscuronet/walletextension'
    })

    await container?.stop()
});