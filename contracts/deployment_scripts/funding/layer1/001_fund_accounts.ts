import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const layer1 =  hre.companionNetworks.layer1;
    
    const {deployer} = await hre.getNamedAccounts();
    const l1Accs = await layer1.getNamedAccounts();
    
    const messageBusAddress = process.env.MESSAGE_BUS_ADDRESS!!

    const messageBus = (await hre.ethers.getContractFactory('MessageBus')).attach(messageBusAddress)
    const prefundAmount = hre.ethers.utils.parseEther("0.5");
    const tx = await messageBus.populateTransaction.sendValueToL2(deployer, prefundAmount, {
        value: prefundAmount
    });


    console.log(`Sending ${prefundAmount} to ${deployer}`);

    const receipt = await layer1.deployments.rawTx({
        from: l1Accs.deployer,
        to: messageBusAddress,
        value: prefundAmount,
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