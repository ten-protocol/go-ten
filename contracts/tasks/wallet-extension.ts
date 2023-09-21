import { task } from "hardhat/config";

import * as dockerApi from 'node-docker-api';
import { on } from 'process';

import * as url from 'node:url';

import { spawn } from 'node:child_process';
import * as path from "path";
import http from 'http';




task("obscuro:wallet-extension:start:local")
.addFlag('withStdOut')
.addParam('rpcUrl', "Node rpc endpoint where the wallet extension should connect to.")
.addOptionalParam('port', "Port that the wallet extension will open for incoming requests.", "3001")
.setAction(async function(args, hre) {
    const nodeUrl = url.parse(args.rpcUrl)

/*    if (args.withStdOut) {
        console.log(`Node url = ${JSON.stringify(nodeUrl, null, "  ")}`);
    }*/

    const walletExtensionPath = path.resolve(hre.config.paths.root, "../tools/walletextension/bin/wallet_extension_linux");
    const weProcess = spawn(walletExtensionPath, [
        `-portWS`, `${args.port}`,
        `-nodeHost`, `${nodeUrl.hostname}`,
        `-nodePortWS`, `${nodeUrl.port}`
    ]);

    console.log("Waiting for Wallet Extension to start");
    await new Promise((resolve, fail)=>{
        const timeoutSchedule = setTimeout(fail, 60_000);
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

    console.log("Wallet Exension started successfully");
    return weProcess;
});

// This is not to be used for internal development. It is targeted at external devs when the obscuro hh plugin is finished!
task("obscuro:wallet-extension:start:docker", "Starts up the wallet extension docker container.")
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

// This is not to be used for internal development. It is targeted at external devs when the obscuro hh plugin is finished!
task("obscuro:wallet-extension:stop:docker", "Stops the docker container with matching image name.")
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

task("obscuro:wallet-extension:add-key", "Creates a viewing key for a specifiec address")
.addParam("address", "The address for which to add key")
.setAction(async function(args, hre) {
    async function viewingKeyForAddress(address: string) : Promise<string> {
        return new Promise((resolve, fail)=> {
    
            const data = {"address": address}
    
            const req = http.request({
                host: '127.0.0.1',
                port: 3000,
                path: '/generateviewingkey/',
                method: 'post',
                headers: {
                    'Content-Type': 'application/json'
                }
            }, (response)=>{
                if (response.statusCode != 200) {
                    console.error(response);
                    fail(response.statusCode);
                    return;
                }
    
                let chunks : string[] = []
                response.on('data', (chunk)=>{
                    chunks.push(chunk);
                });
    
                response.on('end', ()=> { 
                    resolve(chunks.join('')); 
                });
            });
            req.write(JSON.stringify(data));
            req.end()
            setTimeout(resolve, 15_000);
        });
    }
    
    interface SignedData { signature: string, address: string }
    
    async function submitKey(signedData: SignedData) : Promise<number> {
        return await new Promise(async (resolve, fail)=>{ 
            const req = http.request({
                host: '127.0.0.1',
                port: 3000,
                path: '/submitviewingkey/',
                method: 'post',
                headers: {
                    'Content-Type': 'application/json'
                }
            }, (response)=>{
                if (response.statusCode == 200) { 
                    resolve(response.statusCode);
                } else {
                    console.log(response.statusMessage)
                    fail(response.statusCode);
                }
            });
    
            req.write(JSON.stringify(signedData));
            req.end()
        });
    }

    const key = await viewingKeyForAddress(args.address);

    const signaturePromise = (await hre.ethers.getSigner(args.address)).signMessage(`vk${key}`);
    const signedData = { 'signature': await signaturePromise, 'address': args.address };
    await submitKey(signedData)
});