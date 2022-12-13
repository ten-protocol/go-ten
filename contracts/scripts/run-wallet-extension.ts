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
.addParam('rpcUrl', "Which network to pick the node connection info from?")
.setAction(async function(args, hre) {
    const docker = new dockerApi.Docker({ socketPath: '/var/run/docker.sock' });

    const parsedUrl = url.parse(args.rpcUrl)

    const container = await docker.container.create({
        Image: args.dockerImage,
        Cmd: [
            "--port=3000",
            "--portWS=3001",
            `--nodeHost=${parsedUrl.hostname}`,
            `--nodePortWS=${parsedUrl.port}`
        ],
        ExposedPorts: { "3000/tcp": {} },
        PortBindings:  { "3000/tcp": [{ "HostPort": "3000" }] }
    })


    process.on('SIGINT', ()=>{
        container.stop();
    })
    
    await container.start();

    const stream: any = await container.logs({
        follow: true,
        stdout: true,
        stderr: true
    })

    console.log(`\nWallet Extension{ ${container.id.slice(0, 5)} } %>\n`);
    const startupPromise = new Promise((resolveInner)=> {    
        stream.on('data', (msg: any)=> {
            const message = msg.toString();

            console.log(message);    
            if(message.includes("Wallet extension started")) {
                console.log(`Wallet - success!`);
                resolveInner(true);
            }
        });

        setTimeout(resolveInner, 20_000);
    });

    await startupPromise;
    console.log("\n[ . . . ]\n");


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