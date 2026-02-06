import { ethers } from 'hardhat';

/**
 * Script to whitelist existing ERC20 tokens on the TenBridge and register them in NetworkConfig
 * 
 * This script is for adding well-known tokens like USDC/USDT that are already deployed on L1
 * 
 * Environment variables:
 * - TOKEN_ADDRESS: Address of the existing token contract
 * - TOKEN_NAME: Name of the token (e.g., "USD Coin")
 * - TOKEN_SYMBOL: Symbol of the token (e.g., "USDC")
 * - NETWORK_CONFIG_ADDR: Address of the NetworkConfig contract
 * 
 * Example for Sepolia:
 * - USDC: 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
 * - USDT: 0x7169D38820dfd117C3FA1f22a697dBA58d90BA06
 */

const whitelistAndRegisterToken = async function (): Promise<void> {
    console.log('=== Starting token whitelisting and registration ===');
    
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    console.log(`Using signer: ${deployer.address}`);
    
    const tokenAddress = process.env.TOKEN_ADDRESS;
    const tokenName = process.env.TOKEN_NAME;
    const tokenSymbol = process.env.TOKEN_SYMBOL;
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    
    if (!tokenAddress) {
        throw new Error('TOKEN_ADDRESS environment variable is required');
    }
    if (!tokenName) {
        throw new Error('TOKEN_NAME environment variable is required');
    }
    if (!tokenSymbol) {
        throw new Error('TOKEN_SYMBOL environment variable is required');
    }
    if (!networkConfigAddr) {
        throw new Error('NETWORK_CONFIG_ADDR environment variable is required');
    }
    
    if (!ethers.isAddress(tokenAddress)) {
        throw new Error(`Invalid token address: ${tokenAddress}`);
    }
    
    console.log('\nToken parameters:');
    console.table({
        Address: tokenAddress,
        Name: tokenName,
        Symbol: tokenSymbol,
        NetworkConfig: networkConfigAddr
    });
    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
    const addresses = await networkConfig.addresses();
    const bridgeAddress = addresses.l1Bridge;
    
    if (bridgeAddress === ethers.ZeroAddress) {
        throw new Error('Bridge address not set in NetworkConfig');
    }
    console.log(`Bridge address: ${bridgeAddress}`);
    
    console.log('\n Whitelisting token on TenBridge...');
    const bridge = await ethers.getContractAt('TenBridge', bridgeAddress);
    
    // check if token is already whitelisted
    try {
        const hasRole = await bridge.hasRole(
            await bridge.ERC20_TOKEN_ROLE(),
            tokenAddress
        );
        
        if (hasRole) {
            console.log(`Token already whitelisted on bridge`);
        } else {
            console.log(`Whitelisting token: ${tokenName} (${tokenSymbol}) at ${tokenAddress}`);
            const whitelistTx = await bridge.whitelistToken(tokenAddress, tokenName, tokenSymbol);
            console.log(`Transaction hash: ${whitelistTx.hash}`);
            
            const whitelistReceipt = await whitelistTx.wait();
            console.log(`Whitelisted in block: ${whitelistReceipt?.blockNumber}`);
        }
    } catch (error) {
        console.error('Error checking/whitelisting token:', error);
        throw error;
    }
    
    console.log('\n Registering token in NetworkConfig...');
    
    // check if token is already registered
    try {
        const existingAddr = await networkConfig.additionalAddresses(tokenSymbol);
        if (existingAddr !== ethers.ZeroAddress) {
            console.log(`Token "${tokenSymbol}" already registered at ${existingAddr}`);
            if (existingAddr.toLowerCase() !== tokenAddress.toLowerCase()) {
                console.log(`ERROR: Registered address doesn't match! Expected ${tokenAddress}`);
                throw new Error('Address mismatch - token registered with different address');
            } else {
                console.log('Token already registered with correct address');
            }
        } else {
            console.log(`Adding token to NetworkConfig with key: "${tokenSymbol}"...`);
            const registerTx = await networkConfig.addAdditionalAddress(tokenSymbol, tokenAddress);
            console.log(`Transaction hash: ${registerTx.hash}`);
            
            const registerReceipt = await registerTx.wait();
            console.log(`Registered in block: ${registerReceipt?.blockNumber}`);
            
            const registeredAddr = await networkConfig.additionalAddresses(tokenSymbol);
            if (registeredAddr.toLowerCase() !== tokenAddress.toLowerCase()) {
                throw new Error(`Registration failed: expected ${tokenAddress}, got ${registeredAddr}`);
            }
        }
    } catch (error) {
        console.error('Error registering token:', error);
        throw error;
    }
    
    console.log('\n=== Whitelisting and registration completed successfully ===');
    console.table({
        'Token Address': tokenAddress,
        'Token Name': tokenName,
        'Token Symbol': tokenSymbol,
        'Bridge Address': bridgeAddress,
        'NetworkConfig Address': networkConfigAddr,
        'Whitelisted on Bridge': 'Yes',
        'Registered in NetworkConfig': 'Yes'
    });
    console.log('2. Users can query ten_config() to see the registered token');
    console.log('3. Users can now bridge this token using bridge.sendERC20()');
}

whitelistAndRegisterToken()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error('Error during whitelisting/registration:', error);
        process.exit(1);
    });

export default whitelistAndRegisterToken;
