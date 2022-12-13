import { task } from "hardhat/config";
import 'hardhat/types/config';


task("deploy", "Prepares for deploying.")
.setAction(async function(args, hre, runSuper) {

    const rpcURL = hre.network.config.obscuroEncRpcUrl;
    if (rpcURL) {
        await hre.run("run-wallet-extension", {rpcURL : rpcURL});
    }

    await runSuper();

    if (rpcURL) {
        await hre.run("stop-wallet-extension");
    }
});


declare module 'hardhat/types/config' {
    interface HardhatNetworkUserConfig {
        obscuroEncRpcUrl?: string
      }    
}