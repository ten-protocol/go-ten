import { task } from "hardhat/config";
import 'hardhat/types/config';

import {on, exit} from 'process';
import { HardhatNetworkUserConfig } from "hardhat/types/config";

task("obscuro:deploy", "Prepares for deploying.")
.setAction(async function(args, hre, runSuper) {    
    // Trigger shutdown on CTRL + C
    process.on('SIGINT', ()=>exit(1));
    
    const accounts = await hre.getUnnamedAccounts()
    console.log(`Found ${accounts.length} accounts.`);
    accounts.forEach((acc)=>console.log(`Account: ${acc}`));

    // Execute the deploy task provided by the HH deploy plugin.
    await hre.run('deploy');
});

// This extends the hardhat config object for networks to have a key for
// the encrypted rpc endpoint of obscuro nodes.
declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}