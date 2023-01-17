import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const l1Network = hre.companionNetworks.layer1;
    const l2Network = hre; 

    const l1Accounts = await l1Network.getNamedAccounts();
    
    const layer2BridgeDeployment = await l2Network.deployments.get("EthereumBridge");
    
    const setResult = await l1Network.deployments.execute("ObscuroBridge", {
        from: l1Accounts.deployer, 
        log: true,
    }, "setRemoteBridge", layer2BridgeDeployment.address);

    if (setResult.status != 1) {
        console.error("Unable to link L1 and L2 bridges!");
        throw Error("Unable to link L1 and L2 bridges!");
    }
};

export default func;
func.tags = ['SetBridge', 'SetBridge_deploy'];
func.dependencies = ['EthereumBridge'];
