import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const layer1 =  hre.companionNetworks.layer1;
    
    const {deployer} = await hre.getNamedAccounts();
    const l1Accs = await layer1.getNamedAccounts();
    
    var messageBusAddress = process.env.MESSAGE_BUS_ADDRESS!!
    if (messageBusAddress === undefined) {
        const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
        messageBusAddress = networkConfig.MessageBusAddress;
        console.log(`Fallback read of message bus address = ${messageBusAddress}`);
    }
    
    const prefundAmountStr = process.env.PREFUND_FAUCET_AMOUNT!! || "1"

    if (prefundAmountStr == "0") {
        return;
    }

    const messageBus = (await hre.ethers.getContractFactory('MessageBus')).attach(messageBusAddress)
    const prefundAmount = hre.ethers.parseEther(prefundAmountStr);
    const tx = await messageBus.getFunction("sendValueToL2").populateTransaction("0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77", prefundAmount, {
        value: prefundAmount
    });

    console.log(`Sending ${prefundAmount} to ${deployer} through ${messageBusAddress}`);

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
func.tags = ['FaucetFunding', 'FaucetFunding_deploy'];
func.dependencies = ['InitialFunding'];
// No dependencies