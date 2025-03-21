import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const layer1 =  hre.companionNetworks.layer1;
    
    const {deployer} = await hre.getNamedAccounts();
    const l1Accs = await layer1.getNamedAccounts();
    
    var messageBusAddress = process.env.MESSAGE_BUS_ADDRESS!!
    if (messageBusAddress === undefined) {
        const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
        console.log(`Network config = ${JSON.stringify(networkConfig)}`);
        messageBusAddress = networkConfig.L1MessageBus;
        console.log(`Fallback read of message bus address = ${messageBusAddress}`);
    }

    const messageBus = (await hre.ethers.getContractFactory('MessageBus')).attach(messageBusAddress)
    const prefundAmount = hre.ethers.parseEther("0.5");
    console.log(`Prefund amount ${prefundAmount}; MB = ${messageBus}`);

    console.log(`Deployer = ${messageBusAddress}`);
    const tx = await messageBus.getFunction("sendValueToL2").populateTransaction(deployer, prefundAmount, {
        value: prefundAmount.toString()
    });

    console.log(`Sending ${prefundAmount} to ${deployer} through MessageBus ${messageBusAddress}`);

    const receipt = await layer1.deployments.rawTx({
        from: l1Accs.deployer,
        to: messageBusAddress,
        value: prefundAmount.toString(),
        data: tx.data,
        log: true,
        waitConfirmations: 1,
    });
    if (receipt.events?.length === 0) {
        console.log(`Account prefunding status = FAILURE NO CROSS CHAIN EVENT`);
    } else {
        console.log(`Account prefunding status = ${receipt.status}`);
    }
};

export default func;
func.tags = ['GasPrefunding', 'GasPrefunding_deploy'];
// No dependencies