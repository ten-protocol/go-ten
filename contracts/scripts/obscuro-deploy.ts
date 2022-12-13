import { task } from "hardhat/config";
import 'hardhat/types/config';

import axios from 'axios';

import {on, exit} from 'process';
import http from 'http';
import { string } from "hardhat/internal/core/params/argumentTypes";

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

    const rpcURL = hre.network.config.obscuroEncRpcUrl;
    
    if (!rpcURL) {
        await runSuper();
        return;
    } 
    
    await hre.run("run-wallet-extension", {rpcUrl : rpcURL});


    process.on('SIGINT', ()=>{
        exit(1);
    });

    try {
        const { deployer } = await hre.getNamedAccounts()
        const key = await viewingKeyForAddress(deployer);

        const signature = (await hre.ethers.getSigner(deployer)).signMessage(`vk${key}`);
        const signedData = { 'signature': await signature, 'address': deployer }

        await submitKey(signedData);

        setTimeout(()=>exit(1), 60_000);

        await runSuper();
    } 
    finally 
    {
        await hre.run("stop-wallet-extension");
    }

});


declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}