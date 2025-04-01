import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import { network } from 'hardhat';


/* 
    This deployment script deploys the Obscuro Bridge smart contracts on both L1 and L2
    and links them together using the 'setRemoteBridge' call.
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {
        deployments, 
        getNamedAccounts
    } = hre;
        const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
        const bridgeAddress = networkConfig.L1BridgeAddress;

    const bridgeContract = (await hre.ethers.getContractFactory('TenBridge')).attach(bridgeAddress);
    const tx = await bridgeContract.getFunction("setRemoteBridge").populateTransaction(networkConfig.L2BridgeAddress);
    const receipt =await hre.companionNetworks.layer1.deployments.rawTx({
        from: deployer,
        to: bridgeAddress,
        data: tx.data,
        log: true,
        waitConfirmations: 1,
    });
    if (receipt.events?.length === 0) {
        console.log(`Failed to set L2BridgeAddress=${networkConfig.L2BridgeAddress} on bridge contract.`);
    } else {
        console.log(`L2BridgeAddress=${networkConfig.L2BridgeAddress}`);
    }
};

export default func;
func.tags = ['EthereumBridge', 'EthereumBridge_deploy'];