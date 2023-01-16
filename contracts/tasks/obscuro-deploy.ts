import { task } from "hardhat/config";
import 'hardhat/types/config';

import {on, exit} from 'process';
import { HardhatNetworkUserConfig } from "hardhat/types/config";

task("obscuro:deploy", "Prepares for deploying.")
.setAction(async function(args, hre, runSuper) {

    const rpcURL = (hre.network.config as HardhatNetworkUserConfig).obscuroEncRpcUrl;
    
    if (!rpcURL) {
        return;
    } 
    
    await hre.run("obscuro:wallet-extension:start:local", {
        rpcUrl : rpcURL,
        withStdOut: true
    });    

    process.on('SIGINT', ()=>exit(1));

    const { deployer } = await hre.getNamedAccounts();
    await hre.run('obscuro:wallet-extension:add-key', {address: deployer});

    await hre.run('deploy');
});

declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}