import { ethers } from 'hardhat';
import hre from 'hardhat';

/**
 * Script to deploy an ERC20 token and register it in the NetworkConfig contract
 * 
 * Environment variables:
 * - TOKEN_NAME: Name of the token (e.g., "USD Coin")
 * - TOKEN_SYMBOL: Symbol of the token (e.g., "USDC")
 * - TOKEN_DECIMALS: Number of decimals (e.g., 6 for USDC/USDT, 18 for most tokens)
 * - TOKEN_SUPPLY: Initial supply (will be multiplied by 10^decimals), defaults to 1 billion
 * - NETWORK_CONFIG_ADDR: Address of the NetworkConfig contract
 */

const deployAndRegisterToken = async function (): Promise<void> {
    console.log('=== Starting ERC20 token deployment and registration ===');
    
    await hre.run('compile');
    
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    console.log(`Using deployer: ${deployer.address}`);
    
    const tokenName = process.env.TOKEN_NAME;
    const tokenSymbol = process.env.TOKEN_SYMBOL;
    const tokenDecimalsStr = process.env.TOKEN_DECIMALS;
    const tokenSupplyStr = process.env.TOKEN_SUPPLY || '10000'; // 1 billion default
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    
    if (!tokenName) {
        throw new Error('TOKEN_NAME environment variable is required');
    }
    if (!tokenSymbol) {
        throw new Error('TOKEN_SYMBOL environment variable is required');
    }
    if (!tokenDecimalsStr) {
        throw new Error('TOKEN_DECIMALS environment variable is required');
    }
    if (!networkConfigAddr) {
        throw new Error('NETWORK_CONFIG_ADDR environment variable is required');
    }
    
    const tokenDecimals = parseInt(tokenDecimalsStr);
    if (isNaN(tokenDecimals) || tokenDecimals < 0 || tokenDecimals > 18) {
        throw new Error('TOKEN_DECIMALS must be a number between 0 and 18');
    }
    
    const tokenSupply = BigInt(tokenSupplyStr);
    const initialSupply = tokenSupply * (10n ** BigInt(tokenDecimals));
    
    console.log('\nToken parameters:');
    console.table({
        Name: tokenName,
        Symbol: tokenSymbol,
        Decimals: tokenDecimals,
        Supply: tokenSupply.toString(),
        'Initial Supply (with decimals)': initialSupply.toString()
    });
    
    // Deploy the ERC20 token with random salt for unique address
    console.log('\nDeploying ERC20 token...');
    const salt = BigInt(Date.now()); // Use timestamp as salt for unique address
    console.log(`Using salt: ${salt}`);
    
    const TokenFactory = await ethers.getContractFactory('ConfigurableERC20');
    const token = await TokenFactory.deploy(tokenName, tokenSymbol, tokenDecimals, initialSupply, salt);
    await token.waitForDeployment();
    const tokenAddress = await token.getAddress();

    console.log('\n=== Deployment and registration completed successfully ===');
    console.log('Summary:');
    console.table({
        'Token Address': tokenAddress,
        'Token Name': tokenName,
        'Token Symbol': tokenSymbol,
        'Registered in NetworkConfig': 'No',
        'NetworkConfig Address': networkConfigAddr
    });
}

deployAndRegisterToken()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error('Error during deployment:', error);
        process.exit(1);
    });

export default deployAndRegisterToken;
