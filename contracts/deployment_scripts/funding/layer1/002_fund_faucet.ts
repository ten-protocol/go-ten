import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';
import {HardhatEthersProvider} from "@nomicfoundation/hardhat-ethers/internal/hardhat-ethers-provider";


const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const layer1 =  hre.companionNetworks.layer1;
    
    const {deployer} = await hre.getNamedAccounts();
    const l1Accs = await layer1.getNamedAccounts();
    
    const networkConfig : any = await hre.network.provider.request({method: 'net_config'});
    const bridgeAddress = networkConfig.L1Bridge;
    console.log(`TenBridge address = ${bridgeAddress}`);

    const prefundAmountStr = process.env.PREFUND_FAUCET_AMOUNT!! || "500"

    if (prefundAmountStr == "0") {
        return;
    }

    const tenBridge = (await hre.ethers.getContractFactory('TenBridge')).attach(bridgeAddress)
    const prefundAmount = hre.ethers.parseEther(prefundAmountStr);

    // this block is here to prevent underpriced tx failures on testnet startup
    const provider = new HardhatEthersProvider(layer1.provider, "layer1");
    const nonce = await provider.getTransactionCount(l1Accs.deployer, 'latest');
    const feeData = await provider.getFeeData();
    const gasPrice = (feeData.gasPrice! * BigInt(120)) / BigInt(100);

    const tx = await tenBridge.getFunction("sendNative").populateTransaction("0xA58C60cc047592DE97BF1E8d2f225Fc5D959De77");
    console.log(`Sending ${prefundAmount} to ${deployer} through TenBridge ${bridgeAddress}`);

    const receipt = await layer1.deployments.rawTx({
        from: l1Accs.deployer,
        to: bridgeAddress,
        value: prefundAmount.toString(),
        data: tx.data,
        nonce: nonce,
        gasPrice: gasPrice,
        log: true,
        waitConfirmations: 2, // Increase confirmations to ensure tx is mined
    });
    if (receipt.events?.length === 0) {
        console.log(`Account prefunding status = FAILURE NO BRIDGE EVENT`);
    } else {
        console.log(`Account prefunding status = ${receipt.status}`);
    }
};

export default func;
func.tags = ['GasPrefunding', 'GasPrefunding_deploy'];
// No dependencies