import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/* 
    This script whitelists USDC and USDT tokens on the TenBridge contract
    and sets up the WETH address for bridge functionality.
    
    Environment variables:
    - USDC_ADDRESS: USDC token address (optional)
    - USDT_ADDRESS: USDT token address (optional)
    - WETH_ADDRESS: WETH token address (optional - defaults to genesis WETH address)
    
    WETH is pre-deployed at genesis on both L1 and L2 at address 0x1000000000000000000000000000000000000042.
    If WETH_ADDRESS is not set, this address will be used automatically.
*/

// WETH9 is pre-deployed at genesis at this address (same on L1 and L2)
const GENESIS_WETH_ADDRESS = '0x1000000000000000000000000000000000000042';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const { 
        deployments, 
        getNamedAccounts
    } = hre;

    const {deployer} = await getNamedAccounts();

    // Get environment variables
    const usdcAddress = process.env.USDC_ADDRESS;
    const usdtAddress = process.env.USDT_ADDRESS;
    // Use provided WETH address or fall back to genesis WETH address
    const wethAddress = process.env.WETH_ADDRESS || GENESIS_WETH_ADDRESS;

    // Whitelist USDC token if address is provided
    if (usdcAddress) {
        console.log(`Whitelisting USDC: ${usdcAddress}`);
        await deployments.execute('TenBridge', {
            from: deployer,
            log: true
        }, 'whitelistToken', usdcAddress, 'USD Coin', 'USDC');
    } else {
        console.log('Skipping USDC whitelist - USDC_ADDRESS not set');
    }

    // Whitelist USDT token if address is provided
    if (usdtAddress) {
        console.log(`Whitelisting USDT: ${usdtAddress}`);
        await deployments.execute('TenBridge', {
            from: deployer,
            log: true
        }, 'whitelistToken', usdtAddress, 'Tether USD', 'USDT');
    } else {
        console.log('Skipping USDT whitelist - USDT_ADDRESS not set');
    }

    // Set WETH address for WETH unwrapping functionality
    // This also grants ERC20_TOKEN_ROLE to WETH so it can be bridged via sendERC20
    console.log(`Setting WETH address: ${wethAddress}`);
    await deployments.execute('TenBridge', {
        from: deployer,
        log: true
    }, 'setWeth', wethAddress);

    console.log('Token whitelist and WETH configuration completed');
};

export default func;
func.tags = ['TOKENWHITELIST', 'TOKENWHITELIST_deploy'];
func.dependencies = ['TenBridge'];
