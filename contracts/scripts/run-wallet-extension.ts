import { task } from "hardhat/config";

import * as dockerApi from 'node-docker-api';
import YAML from 'yaml'
import * as fs from 'fs';
import { on } from 'process';
import { NetworksUserConfig, NetworkUserConfig } from "hardhat/types";

import * as url from 'node:url';

import { spawn } from 'node:child_process';
import * as path from "path";



task("start:local:wallet-extension")
.addFlag('withStdOut')
.addParam('rpcUrl', "Node rpc endpoint where the wallet extension should connect to.")
.addOptionalParam('port', "Port that the wallet extension will open for incoming requests.", "3001")
.setAction(async function(args, hre) {
    const nodeUrl = url.parse(args.rpcUrl)

    if (args.withStdOut) {
        console.log(`Node url = ${JSON.stringify(nodeUrl, null, "  ")}`);
    }

    const walletExtensionPath = path.resolve(hre.config.paths.root, "../tools/walletextension/bin/wallet_extension_linux");
    const weProcess = spawn(walletExtensionPath, [
        `-portWS`, `${args.port}`,
        `-nodeHost`, `${nodeUrl.hostname}`,
        `-nodePortWS`, `${nodeUrl.port}`
    ]);

    await new Promise((resolve, fail)=>{
        const timeoutSchedule = setTimeout(fail, 40_000);
        weProcess.stdout.on('data', (data: string) => {
            if (args.withStdOut) {
                console.log(data.toString());
            }

            if (data.includes("Wallet extension started")) {
                clearTimeout(timeoutSchedule);
                resolve(true)
            }
        });

        weProcess.stderr.on('data', (data: string) => {
            console.log(data.toString());
        });
    });

    return weProcess;
});

task("start:docker:wallet-extension", "Starts up the wallet extension docker container.")
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
        ExposedPorts: { "3000/tcp": {}, "3001/tcp": {}, "3000/udp": {}, "3001/udp": {} },
        PortBindings:  { "3000/tcp": [{ "HostPort": "3000" }], "3001/tcp": [{ "HostPort": "3001" }] }
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