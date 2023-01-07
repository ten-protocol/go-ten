import { task } from "hardhat/config";
import 'hardhat/types/config';

import axios from 'axios';

import {on, exit} from 'process';
import http from 'http';
import { string } from "hardhat/internal/core/params/argumentTypes";
import { HardhatNetworkUserConfig } from "hardhat/types/config";

async function viewingKeyForAddress(address: string) : Promise<string> {
    return new Promise((resolve, fail)=> {

        const data = {"address": address}

        const req = http.request({
            host: 'localhost',
            port: 3000,
            path: '/generateviewingkey/',
            method: 'post',
            headers: {
                'Content-Type': 'application/json'
            }
        }, (response)=>{
            if (response.statusCode != 200) {
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
            host: 'localhost',
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
                fail(response.statusCode);
            }
        });

        req.write(JSON.stringify(signedData));
        req.end()
    });
}

task("deploy", "Prepares for deploying.")
.setAction(async function(args, hre, runSuper) {

    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;
    
    if (!rpcURL) {
        await runSuper();
        return;
    } 

    process.on('exit', ()=>{
        hre.run("stop-wallet-extension")
    });
    
    await hre.run("run-wallet-extension", {rpcUrl : rpcURL});


    

    process.on('SIGINT', ()=>exit(1));

    try {
        const { deployer } = await hre.getNamedAccounts();
        const key = await viewingKeyForAddress(deployer);

        console.log(`Generated viewing key for ${deployer} - ${key}`);

        const signaturePromise = (await hre.ethers.getSigner(deployer)).signMessage(`vk${key}`);
        const signedData = { 'signature': await signaturePromise, 'address': deployer };
        await submitKey(signedData);

        await runSuper();
    } 
    finally 
    {
        await hre.run("stop-wallet-extension");
    }

});


task("obscuro:deploy", "Prepares for deploying.")
.setAction(async function(args, hre, runSuper) {

    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;
    
    if (!rpcURL) {
        return;
    } 
    
    await hre.run("start:local:wallet-extension", {
        rpcUrl : rpcURL,
        withStdOut: true
    });    

    process.on('SIGINT', ()=>exit(1));

    const { deployer } = await hre.getNamedAccounts();
    const key = await viewingKeyForAddress(deployer);

    console.log(`Generated viewing key for ${deployer} - ${key}`);

    const signaturePromise = (await hre.ethers.getSigner(deployer)).signMessage(`vk${key}`);
    const signedData = { 'signature': await signaturePromise, 'address': deployer };
    await submitKey(signedData);

    await hre.run('deploy');
});

declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}