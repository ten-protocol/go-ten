import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

/* 
    This script sets up the WETH address for bridge functionality.
    
    Note: Token whitelisting (USDC, USDT, etc.) is intentionally not done here.
    Whitelisting must be done via the manual-whitelist-bridge-token workflow, which
    handles the full e2e process: 
    1. L1 whitelist
    2. Wait for L2 message
    3. Finalisation + relaying the message to create the L2 wrapped token.
    
    Whitelisting here without the relay step leaves the bridge in a broken half-state.
    
    Environment variables:
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
    // Use provided WETH address or fall back to genesis WETH address
    const wethAddress = process.env.WETH_ADDRESS || GENESIS_WETH_ADDRESS;

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
