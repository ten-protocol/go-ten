import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const layer1 =  hre.companionNetworks.layer1;

    const {deployer} = await hre.getNamedAccounts();
    const l1Accs = await layer1.getNamedAccounts();

    const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
    console.log(`Network config = ${JSON.stringify(networkConfig)}`);

    const bridgeAddress = networkConfig.L1Bridge;
    console.log(`TenBridge address = ${bridgeAddress}`);

    const tenBridge = (await hre.ethers.getContractFactory('TenBridge')).attach(bridgeAddress)
    const prefundAmount = hre.ethers.parseEther("0.5");
    console.log(`Prefund amount ${prefundAmount}; MB = ${tenBridge}`);

    const tx = await tenBridge.getFunction("sendNative").populateTransaction(deployer);

    console.log(`Sending ${prefundAmount} to ${deployer} through TenBridge ${bridgeAddress}`);

    const receipt = await layer1.deployments.rawTx({
        from: l1Accs.deployer,
        to: bridgeAddress,
        value: prefundAmount.toString(),
        data: tx.data,
        log: true,
        waitConfirmations: 1,
    });
    if (receipt.events?.length === 0) {
        console.log(`Account prefunding status = FAILURE BRIDGE EVENT`);
    } else {
        console.log(`Account prefunding status = ${receipt.status}`);
    }
};

export default func;
func.tags = ['GasPrefunding', 'GasPrefunding_deploy'];
// No dependencies