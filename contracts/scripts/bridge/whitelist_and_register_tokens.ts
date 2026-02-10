import { ethers } from 'hardhat';

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * Script to whitelist existing ERC20 tokens on the TenBridge and register the L2 wrapped token in NetworkConfig
 * 
 * This script follows the L1→L2 bridge whitelisting flow:
 * 1. Call whitelistToken() on L1 TenBridge
 * 2. Extract cross-chain message from logs[0] LogMessagePublished event
 * 3. Poll L2 MessageBus.verifyMessageFinalized() until it returns true
 * 4. Call L2 CrossChainMessenger.relayMessage() to trigger L2 wrapped token creation
 * 5. Get L2 wrapped token address from logs[0] CreatedWrappedToken event
 * 6. Register L2 token address in NetworkConfig with "L2_" prefix
 * 
 * Environment variables:
 * - TOKEN_ADDRESS: Address of the existing L1 token contract
 * - TOKEN_NAME: Name of the token (e.g., "USD Coin")
 * - TOKEN_SYMBOL: Symbol of the token (e.g., "USDC")
 * - NETWORK_CONFIG_ADDR: Address of the NetworkConfig contract (L1)
 * - L2_RPC_URL: RPC URL for L2 (for verifying message finalization and relaying)
 * 
 * Example for Sepolia:
 * - USDC: 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
 * - USDT: 0x7169D38820dfd117C3FA1f22a697dBA58d90BA06
 */

interface CrossChainMessage {
    sender: string;
    sequence: bigint;
    nonce: number;
    topic: number;
    payload: string;
    consistencyLevel: number;
}

async function waitForMessageFinalized(
    l2MessageBus: any,
    crossChainMessage: CrossChainMessage,
    maxRetries = 300
): Promise<void> {
    console.log(`Polling for message finalization (max ${maxRetries} attempts)...`);
    console.log('Cross-chain message:', {
        sender: crossChainMessage.sender,
        sequence: crossChainMessage.sequence.toString(),
        nonce: crossChainMessage.nonce.toString(),
        topic: crossChainMessage.topic.toString(),
        payload: crossChainMessage.payload,
        consistencyLevel: crossChainMessage.consistencyLevel.toString()
    });
    
    let lastError: any = null;
    
    for (let i = 0; i < maxRetries; i++) {
        try {
            const isFinalized = await l2MessageBus.verifyMessageFinalized(crossChainMessage);
            console.log(`Attempt ${i + 1}/${maxRetries}: isFinalized = ${isFinalized}`);
            if (isFinalized) {
                console.log('Message is finalized on L2');
                return;
            }
        } catch (error) {
            lastError = error;
            console.log(`Attempt ${i + 1}/${maxRetries}: Error checking finalization - ${error}`);
        }
        
        await sleep(2000); 
    }
    
    console.error('Last error:', lastError);
    throw new Error('Timeout waiting for message finalization');
}

const whitelistAndRegisterToken = async function (): Promise<void> {
    console.log('=== Starting token whitelisting and L2 registration ===');
    
    const [deployer] = await ethers.getSigners();
    if (!deployer) {
        throw new Error('No deployer signer found');
    }
    console.log(`Using signer: ${deployer.address}`);
    
    const tokenAddress = process.env.TOKEN_ADDRESS;
    const tokenName = process.env.TOKEN_NAME;
    const tokenSymbol = process.env.TOKEN_SYMBOL;
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    const l2RpcUrl = process.env.L2_RPC_URL;
    
    // Get private key from hardhat config (passed via NETWORK_JSON env var)
    const networkJson = JSON.parse(process.env.NETWORK_JSON || '{}');
    const privateKey = networkJson.layer1?.accounts?.[0];
    if (!privateKey) {
        throw new Error('Private key not found in NETWORK_JSON');
    }
    
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
    if (!l2RpcUrl) {
        throw new Error('L2_RPC_URL environment variable is required');
    }
    
    if (!ethers.isAddress(tokenAddress)) {
        throw new Error(`Invalid token address: ${tokenAddress}`);
    }
    
    console.log('\nToken parameters:');
    console.table({
        'L1 Token Address': tokenAddress,
        'Token Name': tokenName,
        'Token Symbol': tokenSymbol,
        'NetworkConfig': networkConfigAddr,
        'L2 RPC URL': l2RpcUrl
    });
    
    // Get L1 addresses from NetworkConfig contract
    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
    const addresses = await networkConfig.addresses();
    const l1BridgeAddress = addresses.l1Bridge;
    const l1MessageBusAddress = addresses.messageBus;
    
    if (l1BridgeAddress === ethers.ZeroAddress) {
        throw new Error('L1 Bridge address not set in NetworkConfig');
    }
    if (l1MessageBusAddress === ethers.ZeroAddress) {
        throw new Error('L1 MessageBus address not set in NetworkConfig');
    }
    
    console.log(`\nL1 Contract Addresses:`);
    console.log(`  L1 Bridge: ${l1BridgeAddress}`);
    console.log(`  L1 MessageBus: ${l1MessageBusAddress}`);
    
    // Get L2 addresses from ten_config RPC
    console.log(`\nQuerying L2 contract addresses from ten_config...`);
    const l2ConfigResponse = await fetch(l2RpcUrl, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            jsonrpc: '2.0',
            method: 'ten_config',
            params: [],
            id: 1
        })
    });
    
    const l2ConfigData = await l2ConfigResponse.json();
    if (!l2ConfigData.result) {
        throw new Error('Failed to get L2 config from ten_config RPC');
    }
    
    const l2BridgeAddress = l2ConfigData.result.L2Bridge;
    const l2CrossChainMessengerAddress = l2ConfigData.result.L2CrossChainMessenger;
    const l2MessageBusAddress = l2ConfigData.result.L2MessageBus;
    
    if (!l2BridgeAddress || l2BridgeAddress === ethers.ZeroAddress) {
        throw new Error('L2 Bridge address not found in ten_config');
    }
    if (!l2CrossChainMessengerAddress || l2CrossChainMessengerAddress === ethers.ZeroAddress) {
        throw new Error('L2 CrossChainMessenger address not found in ten_config');
    }
    if (!l2MessageBusAddress || l2MessageBusAddress === ethers.ZeroAddress) {
        throw new Error('L2 MessageBus address not found in ten_config');
    }
    
    console.log(`\nL2 Contract Addresses (from ten_config):`);
    console.log(`  L2 Bridge: ${l2BridgeAddress}`);
    console.log(`  L2 CrossChainMessenger: ${l2CrossChainMessengerAddress}`);
    console.log(`  L2 MessageBus: ${l2MessageBusAddress}`);
    
    // Connect to L2 using the same deployer wallet
    // Use staticNetwork to skip network detection (Ten doesn't support all eth_ methods)
    const l2Network = ethers.Network.from(443); // Ten testnet chainId
    const l2Provider = new ethers.JsonRpcProvider(l2RpcUrl, l2Network, { staticNetwork: l2Network });
    const l2Signer = new ethers.Wallet(privateKey, l2Provider);
    
    // 1. whitelist token on L1 Bridge
    console.log('\n[STEP 1] Whitelisting token on L1 TenBridge...');
    const l1Bridge = await ethers.getContractAt('TenBridge', l1BridgeAddress);
    const l1MessageBus = await ethers.getContractAt('MessageBus', l1MessageBusAddress);
    
    console.log(`Whitelisting token: ${tokenName} (${tokenSymbol}) at ${tokenAddress}`);
    const whitelistTx = await l1Bridge.whitelistToken(tokenAddress, tokenName, tokenSymbol);
    console.log(`Transaction hash: ${whitelistTx.hash}`);
    
    const whitelistReceipt = await whitelistTx.wait();
    if (!whitelistReceipt) {
        throw new Error('Failed to get whitelist transaction receipt');
    }
    
    if (whitelistReceipt.status !== 1) {
        console.error('Transaction failed! Receipt:', whitelistReceipt);
        throw new Error(`Whitelist transaction failed with status: ${whitelistReceipt.status}`);
    }
    
    console.log(`Whitelisted successfully in block: ${whitelistReceipt.blockNumber}`);
    
    // Extract LogMessagePublished event - use MessageBus contract to parse logs
    const logMessagePublishedEvents = whitelistReceipt.logs
        .map(log => {
            try {
                return l1MessageBus.interface.parseLog({ topics: log.topics as string[], data: log.data });
            } catch {
                return null;
            }
        })
        .filter(event => event && event.name === 'LogMessagePublished');
    
    if (logMessagePublishedEvents.length === 0) {
        throw new Error('LogMessagePublished event not found in transaction receipt');
    }
    
    // Get the first LogMessagePublished event (logs[0])
    const logEvent = logMessagePublishedEvents[0]!;
    const crossChainMessage: CrossChainMessage = {
        sender: logEvent.args.sender,
        sequence: logEvent.args.sequence,
        nonce: logEvent.args.nonce,
        topic: logEvent.args.topic,
        payload: logEvent.args.payload,
        consistencyLevel: logEvent.args.consistencyLevel
    };
    
    console.log('Cross-chain message extracted from logs[0]');
    
    // 2. Wait for L2 to process cross-chain message (happens via synthetic tx)
    console.log('\n[STEP 2] Waiting for L2 to process cross-chain message...');
    console.log('The Ten enclave creates a synthetic transaction to store the message on L2 MessageBus');
    console.log('Waiting 60 seconds for this to complete...');
    await sleep(60000);
    
    // 3. relay message on L2 to create wrapped token
    console.log('\n[STEP 3] Relaying message on L2 to create wrapped token...');
    
    // ISSUE IS HERE WITH THE SIGNER
    // ISSUE IS HERE WITH THE SIGNER
    // ISSUE IS HERE WITH THE SIGNER
    // ISSUE IS HERE WITH THE SIGNER
    // ISSUE IS HERE WITH THE SIGNER
    // ISSUE IS HERE WITH THE SIGNER
    // ISSUE IS HERE WITH THE SIGNER
    // ISSUE IS HERE WITH THE SIGNER
    const l2CrossChainMessenger = await ethers.getContractAt('CrossChainMessenger', l2CrossChainMessengerAddress, l2Signer);
    
    // Manually set nonce to 0 to avoid eth_getTransactionCount call
    // This assumes it's the first transaction from this address on L2
    
    const relayTx = await l2CrossChainMessenger.relayMessage(crossChainMessage, { nonce: 0 });
    console.log(`Relay transaction hash: ${relayTx.hash}`);
    
    const relayReceipt = await relayTx.wait();
    if (!relayReceipt) {
        throw new Error('Failed to get relay transaction receipt');
    }
    console.log(`Relayed in block: ${relayReceipt.blockNumber}`);
    
    // Extract CreatedWrappedToken event from relay receipt
    const l2Bridge = await ethers.getContractAt('EthereumBridge', l2BridgeAddress, l2Signer);
    
    const createdWrappedTokenEvents = relayReceipt.logs
        .map((log: any) => {
            try {
                return l2Bridge.interface.parseLog({ topics: log.topics as string[], data: log.data });
            } catch {
                return null;
            }
        })
        .filter((event: any) => event && event.name === 'CreatedWrappedToken');
    
    if (createdWrappedTokenEvents.length === 0) {
        throw new Error('CreatedWrappedToken event not found in transaction logs');
    }
    
    const l2TokenAddress = createdWrappedTokenEvents[0]!.args.localAddress;
    console.log(`L2 wrapped token created at: ${l2TokenAddress}`);
    
    // 4. register L2 token address in NetworkConfig
    console.log('\n[STEP 4] Registering L2 token in NetworkConfig...');
    
    const l2TokenKey = `L2_${tokenSymbol}`;
    
    try {
        const existingAddr = await networkConfig.additionalAddresses(l2TokenKey);
        if (existingAddr !== ethers.ZeroAddress) {
            console.log(`Token "${l2TokenKey}" already registered at ${existingAddr}`);
            if (existingAddr.toLowerCase() !== l2TokenAddress.toLowerCase()) {
                console.warn(`WARNING: Registered address (${existingAddr}) doesn't match created L2 token (${l2TokenAddress})`);
                console.warn('This may indicate a previous failed run. Consider manual verification.');
            } else {
                console.log(' L2 token already registered with correct address');
            }
        } else {
            console.log(`Adding L2 token to NetworkConfig with key: "${l2TokenKey}"...`);
            const registerTx = await networkConfig.addAdditionalAddress(l2TokenKey, l2TokenAddress);
            console.log(`Transaction hash: ${registerTx.hash}`);
            
            const registerReceipt = await registerTx.wait();
            console.log(` Registered in block: ${registerReceipt?.blockNumber}`);
            
            const registeredAddr = await networkConfig.additionalAddresses(l2TokenKey);
            if (registeredAddr.toLowerCase() !== l2TokenAddress.toLowerCase()) {
                throw new Error(`Registration failed: expected ${l2TokenAddress}, got ${registeredAddr}`);
            }
        }
    } catch (error) {
        console.error('Error registering L2 token:', error);
        throw error;
    }
    
    // 5. verify registration via ten_config RPC
    console.log('\n[STEP 5] Verifying registration via ten_config...');
    
    try {
        const response = await fetch(l2RpcUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                jsonrpc: '2.0',
                method: 'ten_config',
                params: [],
                id: 1
            })
        });
        
        const data = await response.json();
        if (data.result && data.result.AdditionalContracts) {
            const registeredAddress = data.result.AdditionalContracts[l2TokenKey];
            if (registeredAddress && registeredAddress.toLowerCase() === l2TokenAddress.toLowerCase()) {
                console.log(`Verified: ${l2TokenKey} = ${registeredAddress}`);
            } else if (registeredAddress) {
                console.warn(`Warning: Address mismatch in ten_config!`);
                console.warn(`Expected: ${l2TokenAddress}`);
                console.warn(`Got: ${registeredAddress}`);
            } else {
                console.warn(`Warning: ${l2TokenKey} not found in ten_config (may need time to sync)`);
            }
        } else {
            console.warn('Could not verify via ten_config (RPC may not be ready)');
        }
    } catch (error) {
        console.warn('ould not verify via ten_config:', error);
    }
    
    console.log('\n=== Token whitelisting and registration completed successfully ===');
    console.table({
        'L1 Token Address': tokenAddress,
        'L2 Token Address': l2TokenAddress,
        'Token Name': tokenName,
        'Token Symbol': tokenSymbol,
        'NetworkConfig Key': l2TokenKey,
        'L1 Bridge Address': l1BridgeAddress,
        'NetworkConfig Address': networkConfigAddr
    });
}

whitelistAndRegisterToken()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error('Error during whitelisting/registration:', error);
        process.exit(1);
    });

export default whitelistAndRegisterToken;
