import { task } from "hardhat/config";
import 'hardhat/types/config';


import { on } from 'process';


task("deploy", "Prepares for deploying.")
.setAction(async function(args, hre, runSuper) {

    const rpcURL = hre.network.config.obscuroEncRpcUrl;
    if (rpcURL) {
        await hre.run("run-wallet-extension", {rpcUrl : rpcURL});
    }

    await new Promise((resolve)=>{
        process.on('SIGINT', ()=>{
            resolve(true);
        })
    })

    try {

        await runSuper();

    } finally {
        if (rpcURL) {
            await new Promise((resolve)=>{
                process.on('SIGINT', ()=>{
                    resolve(true);
                })
            })

            await hre.run("stop-wallet-extension");
        }
    }

});


declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}