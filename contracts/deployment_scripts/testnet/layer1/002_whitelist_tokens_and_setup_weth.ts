import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/* 
    This script whitelists USDC and USDT tokens on the TenBridge contract
    and sets up the WETH address for mainnet compatibility.
    
    Environment variables required:
    - USDC_ADDRESS: Mainnet USDC token address
    - USDT_ADDRESS: Mainnet USDT token address  
    - WETH_ADDRESS: Mainnet WETH token address
*/

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    // Get environment variables
    const usdcAddress = process.env.USDC_ADDRESS;
    const usdtAddress = process.env.USDT_ADDRESS;
    const wethAddress = process.env.WETH_ADDRESS;

    if (!usdcAddress) {
        console.error("USDC_ADDRESS environment variable is not set.");
        process.exit(1);
    }

    if (!usdtAddress) {
        console.error("USDT_ADDRESS environment variable is not set.");
        process.exit(1);
    }

    if (!wethAddress) {
        console.error("WETH_ADDRESS environment variable is not set.");
        process.exit(1);
    }

    console.log(`Whitelisting USDC: ${usdcAddress}`);
    console.log(`Whitelisting USDT: ${usdtAddress}`);
    console.log(`Setting WETH address: ${wethAddress}`);

    // Whitelist USDC token
    await deployments.execute('TenBridge', {
        from: deployer,
        log: true
    }, 'whitelistToken', usdcAddress, 'USD Coin', 'USDC');

    // Whitelist USDT token  
    await deployments.execute('TenBridge', {
        from: deployer,
        log: true
    }, 'whitelistToken', usdtAddress, 'Tether USD', 'USDT');

    // Set WETH address for WETH unwrapping functionality
    await deployments.execute('TenBridge', {
        from: deployer,
        log: true
    }, 'setWeth', wethAddress);

    console.log('✅ Successfully whitelisted USDC, USDT and configured WETH address');
};

export default func;
func.tags = ['TOKENWHITELIST', 'TOKENWHITELIST_deploy'];
func.dependencies = ['TenBridge'];
