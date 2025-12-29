import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/* 
    This script whitelists USDC and USDT tokens on the TenBridge contract
    and sets up the WETH address for mainnet compatibility.
    
    Environment variables (optional for local testnet):
    - USDC_ADDRESS: Mainnet USDC token address
    - USDT_ADDRESS: Mainnet USDT token address  
    - WETH_ADDRESS: Mainnet WETH token address
    
    If these are not set, the script will skip the corresponding operations.
    This allows local testnets to deploy without requiring external token addresses.
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

    // Check if any token addresses are configured
    const hasTokenConfig = usdcAddress || usdtAddress || wethAddress;
    
    if (!hasTokenConfig) {
        console.log('⏭️  Skipping token whitelist and WETH setup - no token addresses configured');
        console.log('   Set USDC_ADDRESS, USDT_ADDRESS, and WETH_ADDRESS env vars to enable');
        return;
    }

    // Whitelist USDC token if address is provided
    if (usdcAddress) {
        console.log(`Whitelisting USDC: ${usdcAddress}`);
        await deployments.execute('TenBridge', {
            from: deployer,
            log: true
        }, 'whitelistToken', usdcAddress, 'USD Coin', 'USDC');
    } else {
        console.log('⏭️  Skipping USDC whitelist - USDC_ADDRESS not set');
    }

    // Whitelist USDT token if address is provided
    if (usdtAddress) {
        console.log(`Whitelisting USDT: ${usdtAddress}`);
        await deployments.execute('TenBridge', {
            from: deployer,
            log: true
        }, 'whitelistToken', usdtAddress, 'Tether USD', 'USDT');
    } else {
        console.log('⏭️  Skipping USDT whitelist - USDT_ADDRESS not set');
    }

    // Set WETH address for WETH unwrapping functionality if provided
    if (wethAddress) {
        console.log(`Setting WETH address: ${wethAddress}`);
        await deployments.execute('TenBridge', {
            from: deployer,
            log: true
        }, 'setWeth', wethAddress);
    } else {
        console.log('⏭️  Skipping WETH setup - WETH_ADDRESS not set');
    }

    console.log('✅ Token whitelist and WETH configuration completed');
};

export default func;
func.tags = ['TOKENWHITELIST', 'TOKENWHITELIST_deploy'];
func.dependencies = ['TenBridge'];
