import { ethers } from 'hardhat';
import type { MessageBus, Structs } from '../../typechain-types/src/cross_chain_messaging/common/MessageBus';
import { MessageBus__factory } from '../../typechain-types/factories/src/cross_chain_messaging/common/MessageBus__factory';

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * Script to whitelist existing ERC20 tokens on the TenBridge and register the L2 wrapped token in NetworkConfig
 * 
 * This script follows the L1→L2 bridge whitelisting flow:
 * 1. Query NetworkConfig on L1 to get all contract addresses (L1 Bridge, L1 MessageBus, L2 Bridge, L2 CrossChainMessenger)
 * 2. Call whitelistToken() on L1 TenBridge
 * 3. Extract cross-chain message from logs[0] LogMessagePublished event
 * 4. Poll L2 MessageBus.verifyMessageFinalized() (via gateway) until it returns true
 * 5. Call L2 CrossChainMessenger.relayMessage() to trigger L2 wrapped token creation
 * 6. Get L2 wrapped token address from CreatedWrappedToken event
 * 7. Register L2 token address in L1 NetworkConfig with "L2_" prefix
 * 
 * Environment variables:
 * - TOKEN_ADDRESS: Address of the existing L1 token contract
 * - TOKEN_NAME: Name of the token (e.g., "USD Coin")
 * - TOKEN_SYMBOL: Symbol of the token (e.g., "USDC")
 * - NETWORK_CONFIG_ADDR: Address of the NetworkConfig contract (L1)
 * - NETWORK_JSON: JSON config with layer1 network info { url: L1_RPC_URL, accounts: [private_key] }
 * - NETWORK_ENV: The environment to get the chain ID
 * - L2_GATEWAY_URL (optional): Gateway URL for finalized-message polling (defaults to NETWORK_JSON.layer2.url)
 *   Supports either ten_config or net_config on the target RPC.
 * 
 * Note: For local testing, L2_RPC_URL should point to the L2 node directly.
 * For production, it may point to a wallet extension gateway that provides full RPC support.
 * 
 * Example for Sepolia:
 * - USDC: 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
 * - USDT: 0x7169D38820dfd117C3FA1f22a697dBA58d90BA06
 */

type CrossChainMessage = Structs.CrossChainMessageStruct;
const PREFUND_L2_AMOUNT = ethers.parseEther('0.5');
const ENV_CHAIN_IDS: Record<string, number> = {
    sepolia: 8443,
    uat: 7443,
    dev: 6443,
    local: 5443,
    mainnet: 443
};

function normalizeNetworkEnv(value: string): string {
    return value.toLowerCase().trim().replace(/-testnet$/, '');
}

function resolveL2ChainId(networkEnvRaw: string, overrideChainId?: string): number {
    if (overrideChainId) {
        const parsed = Number(overrideChainId);
        if (!Number.isFinite(parsed) || parsed <= 0) {
            throw new Error(`Invalid NETWORK_CHAINID override: ${overrideChainId}`);
        }
        return parsed;
    }

    const networkEnv = normalizeNetworkEnv(networkEnvRaw);
    const mapped = ENV_CHAIN_IDS[networkEnv];
    if (!mapped) {
        throw new Error(
            `Unsupported NETWORK_ENV "${networkEnvRaw}". Supported values: ${Object.keys(ENV_CHAIN_IDS).join(', ')}`
        );
    }
    return mapped;
}

async function fetchRpcConfig(rpcUrl: string): Promise<any> {
    const methods = ['ten_config', 'net_config'];
    let lastError: unknown = null;

    for (const method of methods) {
        try {
            const response = await fetch(rpcUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    jsonrpc: '2.0',
                    method,
                    params: [],
                    id: 1
                })
            });

            if (!response.ok) {
                throw new Error(`${method} RPC failed with status ${response.status}`);
            }

            const data = await response.json();
            if (data.error || !data.result) {
                throw new Error(`${method} returned invalid response: ${JSON.stringify(data.error ?? data)}`);
            }

            return data.result;
        } catch (error) {
            lastError = error;
        }
    }

    throw new Error(`Unable to fetch RPC config from gateway. Last error: ${String(lastError)}`);
}

async function fetchChainIdFromGateway(gatewayUrl: string): Promise<number> {
    const response = await fetch(gatewayUrl, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            jsonrpc: '2.0',
            method: 'eth_chainId',
            params: [],
            id: 1
        })
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch chain ID from gateway: HTTP ${response.status}`);
    }

    const data = await response.json();
    if (data?.error || !data?.result || typeof data.result !== 'string') {
        throw new Error(`Invalid eth_chainId response from gateway: ${JSON.stringify(data?.error ?? data)}`);
    }

    return Number(BigInt(data.result));
}

async function ensureGatewayAccountRegistered(gatewayBaseUrl: string, privateKey: string): Promise<string> {
    const baseUrl = gatewayBaseUrl.replace(/\/+$/, '');
    const wallet = new ethers.Wallet(privateKey);

    const joinResponse = await fetch(`${baseUrl}/join`, { method: 'GET' });
    if (!joinResponse.ok) {
        throw new Error(`Gateway join failed with status ${joinResponse.status}`);
    }
    const token = (await joinResponse.text()).trim();
    if (!token) {
        throw new Error('Gateway join returned an empty token');
    }

    const queryUrl = new URL(`${baseUrl}/query/address`);
    queryUrl.searchParams.set('token', token);
    queryUrl.searchParams.set('a', wallet.address);
    const queryResponse = await fetch(queryUrl.toString(), { method: 'GET' });

    let isRegistered = false;
    if (queryResponse.ok) {
        const queryData = await queryResponse.json();
        isRegistered = queryData?.status === true;
    }

    if (!isRegistered) {
        const chainId = await fetchChainIdFromGateway(gatewayBaseUrl);
        const domain = {
            name: 'Ten',
            version: '1.0',
            chainId,
            verifyingContract: ethers.ZeroAddress
        };
        const types = {
            Authentication: [{ name: 'Encryption Token', type: 'address' }]
        };
        const message = {
            'Encryption Token': `0x${token}`
        };
        const signature = await wallet.signTypedData(domain, types, message);

        const authenticateUrl = new URL(`${baseUrl}/authenticate/`);
        authenticateUrl.searchParams.set('token', token);
        const authResponse = await fetch(authenticateUrl.toString(), {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                address: wallet.address,
                signature
            })
        });

        if (!authResponse.ok) {
            throw new Error(`Gateway account authentication failed with status ${authResponse.status}`);
        }
        const verifyResponse = await fetch(queryUrl.toString(), { method: 'GET' });
        if (!verifyResponse.ok) {
            throw new Error(`Gateway registration verification failed with status ${verifyResponse.status}`);
        }
        const verifyData = await verifyResponse.json();
        if (verifyData?.status !== true) {
            throw new Error(`Gateway registration failed for ${wallet.address}`);
        }

        console.log(`Registered ${wallet.address} with gateway for token ${token}`);
    }

    const rpcUrl = new URL(gatewayBaseUrl);
    rpcUrl.searchParams.set('token', token);
    return rpcUrl.toString();
}

async function waitForMessageFinalized(
    l2messageBus: MessageBus,
    crossChainMessage: CrossChainMessage,
    maxRetries = 100,
    networkLabel = 'target'
): Promise<void> {
    let lastError: unknown = null;

    for (let i = 0; i < maxRetries; i++) {
        try {
            if (await l2messageBus.verifyMessageFinalized(crossChainMessage)) {
                return;
            }
        } catch (error) {
            lastError = error;
        }

        if ((i + 1) % 10 === 0) {
            console.log(`Still waiting for ${networkLabel} message finalization (${i + 1}/${maxRetries})...`);
        }

        await sleep(2000);
    }

    const suffix = lastError ? ` Last error: ${String(lastError)}` : '';
    throw new Error(`Timeout waiting for message finalization.${suffix}`);
}

function extractCrossChainMessageFromReceipt(
    receipt: { logs: Array<{ topics: readonly string[]; data: string }> },
    l1MessageBus: MessageBus
): CrossChainMessage {
    for (const log of receipt.logs) {
        let parsed: ReturnType<typeof l1MessageBus.interface.parseLog> | null = null;
        try {
            parsed = l1MessageBus.interface.parseLog({ topics: [...log.topics], data: log.data });
        } catch {
            continue;
        }

        if (parsed && parsed.name === 'LogMessagePublished') {
            return {
                sender: parsed.args.sender as string,
                sequence: BigInt(parsed.args.sequence.toString()),
                nonce: Number(parsed.args.nonce),
                topic: Number(parsed.args.topic),
                payload: parsed.args.payload as string,
                consistencyLevel: Number(parsed.args.consistencyLevel)
            };
        }
    }

    throw new Error('LogMessagePublished event not found in transaction receipt');
}

async function waitForWrappedTokenAddress(l2Bridge: any, l1TokenAddress: string): Promise<string> {
    const filter = l2Bridge.filters.CreatedWrappedToken();
    const maxAttempts = 12;

    for (let attempt = 1; attempt <= maxAttempts; attempt++) {
        const waitSeconds = Math.min(attempt * 5, 30);
        await sleep(waitSeconds * 1000);

        const events = await l2Bridge.queryFilter(filter, -300, 'latest');
        const match = [...events]
            .reverse()
            .find((event: any) => event.args?.remoteAddress?.toLowerCase() === l1TokenAddress.toLowerCase());

        if (match?.args?.localAddress) {
            return match.args.localAddress;
        }

        console.log(`Waiting for wrapped token creation (${attempt}/${maxAttempts})...`);
    }

    throw new Error(
        `CreatedWrappedToken event not found after ${maxAttempts} attempts.\n` +
        'Check enclave logs for wrapped token creation.'
    );
}

type ScriptConfig = {
    tokenAddress: string;
    tokenName: string;
    tokenSymbol: string;
    networkConfigAddr: string;
    l2RpcUrl: string;
    l1PrivateKey: string;
    l2PrivateKey: string;
    l1RpcUrl: string;
    l2ChainId: number;
    networkEnv: string;
};

type ContractAddresses = {
    l1BridgeAddress: string;
    l1MessageBusAddress: string;
    l2BridgeAddress: string;
    l2CrossChainMessengerAddress: string;
};

function getScriptConfig(): ScriptConfig {
    const tokenAddress = process.env.TOKEN_ADDRESS;
    const tokenName = process.env.TOKEN_NAME;
    const tokenSymbol = process.env.TOKEN_SYMBOL;
    const networkConfigAddr = process.env.NETWORK_CONFIG_ADDR;
    const networkEnvRaw = process.env.NETWORK_ENV ?? 'mainnet';
    const networkChainIdOverride = process.env.NETWORK_CHAINID;

    const networkJson = JSON.parse(process.env.NETWORK_JSON || '{}');
    const l1PrivateKey = networkJson.layer1?.accounts?.[0];
    const l2PrivateKey = networkJson.layer2?.accounts?.[0] ?? l1PrivateKey;
    const l1RpcUrl = networkJson.layer1?.url;
    const l2RpcUrl = networkJson.layer2?.url;
    const l2ChainId = resolveL2ChainId(networkEnvRaw, networkChainIdOverride);

    if (!l1PrivateKey) {
        throw new Error('Private key not found in NETWORK_JSON.layer1.accounts');
    }
    if (!l2PrivateKey) {
        throw new Error('Private key not found in NETWORK_JSON.layer2.accounts');
    }
    if (!l1RpcUrl) {
        throw new Error('L1 RPC URL not found in NETWORK_JSON.layer1.url');
    }
    if (!l2RpcUrl) {
        throw new Error('L2 RPC URL not found in NETWORK_JSON.layer2.url');
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
    
    if (!ethers.isAddress(tokenAddress)) {
        throw new Error(`Invalid token address: ${tokenAddress}`);
    }

    return {
        tokenAddress,
        tokenName,
        tokenSymbol,
        networkConfigAddr,
        l2RpcUrl,
        l1PrivateKey,
        l2PrivateKey,
        l1RpcUrl,
        l2ChainId,
        networkEnv: normalizeNetworkEnv(networkEnvRaw)
    };
}

async function step1QueryNetworkAddresses(networkConfigAddr: string): Promise<{ networkConfig: any; addresses: ContractAddresses }> {
    const networkConfig = await ethers.getContractAt('NetworkConfig', networkConfigAddr);
    const addresses = await networkConfig.addresses();
    const l1BridgeAddress = addresses.l1Bridge;
    const l1MessageBusAddress = addresses.messageBus;
    const l2BridgeAddress = addresses.l2Bridge;
    const l2CrossChainMessengerAddress = addresses.l2CrossChainMessenger;

    const parsedAddresses: ContractAddresses = {
        l1BridgeAddress,
        l1MessageBusAddress,
        l2BridgeAddress,
        l2CrossChainMessengerAddress
    };

    if (l1BridgeAddress === ethers.ZeroAddress) {
        throw new Error('L1 Bridge address not set in NetworkConfig');
    }
    if (l1MessageBusAddress === ethers.ZeroAddress) {
        throw new Error('L1 MessageBus address not set in NetworkConfig');
    }

    if (l2BridgeAddress === ethers.ZeroAddress) {
        throw new Error('L2 Bridge address not set in NetworkConfig');
    }
    if (l2CrossChainMessengerAddress === ethers.ZeroAddress) {
        throw new Error('L2 CrossChainMessenger address not set in NetworkConfig');
    }

    return { networkConfig, addresses: parsedAddresses };
}

async function step2WhitelistToken(
    l1BridgeAddress: string,
    l1MessageBusAddress: string,
    tokenAddress: string,
    tokenName: string,
    tokenSymbol: string
): Promise<{ whitelistReceipt: any; l1MessageBus: MessageBus }> {
    const l1Bridge = await ethers.getContractAt('TenBridge', l1BridgeAddress);
    const l1MessageBus = await ethers.getContractAt('MessageBus', l1MessageBusAddress);

    const whitelistTx = await l1Bridge.whitelistToken(tokenAddress, tokenName, tokenSymbol);
    console.log(`Transaction hash: ${whitelistTx.hash}`);

    const whitelistReceipt = await whitelistTx.wait();
    if (!whitelistReceipt) {
        throw new Error('Failed to fetch whitelist transaction receipt');
    }
    if (whitelistReceipt.status !== 1) {
        throw new Error(`Whitelist transaction failed with status: ${whitelistReceipt.status}`);
    }

    return { whitelistReceipt, l1MessageBus };
}

function step3ExtractCrossChainMessage(whitelistReceipt: any, l1MessageBus: MessageBus): CrossChainMessage {
    return extractCrossChainMessageFromReceipt(whitelistReceipt, l1MessageBus);
}

async function step4WaitForMessageFinalization(
    l2RpcUrl: string,
    l2PrivateKey: string,
    crossChainMessage: CrossChainMessage
): Promise<string> {
    const authenticatedGatewayUrl = await ensureGatewayAccountRegistered(l2RpcUrl, l2PrivateKey);
    const tenConfig = await fetchRpcConfig(authenticatedGatewayUrl);
    const l2MessageBusAddress = tenConfig.L2MessageBus;

    if (!l2MessageBusAddress || !ethers.isAddress(l2MessageBusAddress) || l2MessageBusAddress === ethers.ZeroAddress) {
        throw new Error(`Invalid L2MessageBus address from ten_config: ${String(l2MessageBusAddress)}`);
    }

    const l2GatewayProvider = new ethers.JsonRpcProvider(authenticatedGatewayUrl);
    const l2MessageBus = MessageBus__factory.connect(l2MessageBusAddress, l2GatewayProvider);
    await waitForMessageFinalized(l2MessageBus, crossChainMessage, 120, 'L2');
    return authenticatedGatewayUrl;
}

async function step5PrefundL2SignerIfNeeded(
    l1BridgeAddress: string,
    l2Signer: any,
    l2Provider: any
): Promise<void> {
    const currentBalance = await l2Provider.getBalance(l2Signer.address);
    if (currentBalance > 0n) {
        console.log(`L2 signer already funded: ${ethers.formatEther(currentBalance)} ETH`);
        return;
    }

    const l1Bridge = await ethers.getContractAt('TenBridge', l1BridgeAddress);
    const prefundTx = await l1Bridge.sendNative(l2Signer.address, { value: PREFUND_L2_AMOUNT });
    await prefundTx.wait();
    console.log(`Prefund tx submitted (${ethers.formatEther(PREFUND_L2_AMOUNT)} ETH): ${prefundTx.hash}`);

    const maxAttempts = 30;
    for (let i = 0; i < maxAttempts; i++) {
        const balance = await l2Provider.getBalance(l2Signer.address);
        if (balance > 0n) {
            console.log(`L2 signer funded: ${ethers.formatEther(balance)} ETH`);
            return;
        }
        await sleep(2000);
    }

    throw new Error('Timed out waiting for L2 prefund to arrive');
}

async function step5RelayMessageOnL2(
    l2CrossChainMessengerAddress: string,
    l2Signer: any,
    crossChainMessage: CrossChainMessage
): Promise<void> {
    const l2CrossChainMessenger = await ethers.getContractAt('CrossChainMessenger', l2CrossChainMessengerAddress, l2Signer);
    const relayTx = await l2CrossChainMessenger.relayMessage(crossChainMessage);
    console.log(`Relay tx submitted: ${relayTx.hash}`);
}

async function step6GetWrappedTokenAddress(
    l2BridgeAddress: string,
    l2Signer: any,
    tokenAddress: string
): Promise<string> {
    const l2Bridge = await ethers.getContractAt('EthereumBridge', l2BridgeAddress, l2Signer);
    return waitForWrappedTokenAddress(l2Bridge, tokenAddress);
}

async function step7RegisterWrappedToken(
    networkConfig: any,
    l2TokenKey: string,
    l2TokenAddress: string
): Promise<void> {
    const existingAddr = await networkConfig.additionalAddresses(l2TokenKey);
    if (existingAddr === ethers.ZeroAddress) {
        const registerTx = await networkConfig.addAdditionalAddress(l2TokenKey, l2TokenAddress);
        console.log(`Register tx hash: ${registerTx.hash}`);
        await registerTx.wait();
    } else if (existingAddr.toLowerCase() !== l2TokenAddress.toLowerCase()) {
        console.warn(`"${l2TokenKey}" already exists with different address: ${existingAddr}`);
    }

    const registeredAddr = await networkConfig.additionalAddresses(l2TokenKey);
    if (registeredAddr.toLowerCase() !== l2TokenAddress.toLowerCase()) {
        throw new Error(`Registration failed: expected ${l2TokenAddress}, got ${registeredAddr}`);
    }
}

async function verifyViaTenConfig(l2RpcUrl: string, l2TokenKey: string, l2TokenAddress: string): Promise<void> {
    try {
        const tenConfig = await fetchRpcConfig(l2RpcUrl);
        const tenConfigAddr = tenConfig.AdditionalContracts?.[l2TokenKey];

        if (tenConfigAddr?.toLowerCase() === l2TokenAddress.toLowerCase()) {
            console.log(`Verified in ten_config: ${l2TokenKey} = ${tenConfigAddr}`);
        } else if (tenConfigAddr) {
            console.warn(`ten_config mismatch for ${l2TokenKey}: ${tenConfigAddr}`);
        } else {
            console.warn(`${l2TokenKey} not found in ten_config yet`);
        }
    } catch (error) {
        console.warn(`Could not verify via ten_config: ${String(error)}`);
    }
}

const whitelistAndRegisterToken = async function (): Promise<void> {
    console.log('Starting token whitelisting and L2 registration...');
    const config = getScriptConfig();

    console.table({
        'L1 Token Address': config.tokenAddress,
        'Token Name': config.tokenName,
        'Token Symbol': config.tokenSymbol,
        'NetworkConfig': config.networkConfigAddr,
        'L1 RPC URL': config.l1RpcUrl,
        'L2 RPC URL': config.l2RpcUrl,
        'Network Env': config.networkEnv,
        'L2 Chain ID': config.l2ChainId
    });

    console.log('[1/7] Querying network addresses...');
    const { networkConfig, addresses } = await step1QueryNetworkAddresses(config.networkConfigAddr);

    console.log('[2/7] Whitelisting token on L1 bridge...');
    const { whitelistReceipt, l1MessageBus } = await step2WhitelistToken(
        addresses.l1BridgeAddress,
        addresses.l1MessageBusAddress,
        config.tokenAddress,
        config.tokenName,
        config.tokenSymbol
    );

    console.log('[3/7] Extracting cross-chain message...');
    const crossChainMessage = step3ExtractCrossChainMessage(whitelistReceipt, l1MessageBus);

    console.log('[4/7] Waiting for L2 message finalization (via gateway)...');
    const authenticatedL2RpcUrl = await step4WaitForMessageFinalization(config.l2RpcUrl, config.l2PrivateKey, crossChainMessage);

    const l2Network = ethers.Network.from({ name: `ten-${config.networkEnv}`, chainId: config.l2ChainId });
    const l2Provider = new ethers.JsonRpcProvider(authenticatedL2RpcUrl, l2Network, { staticNetwork: l2Network });
    const l2Signer = new ethers.Wallet(config.l2PrivateKey, l2Provider);

    console.log('[5/7] Prefunding L2 signer...');
    await step5PrefundL2SignerIfNeeded(addresses.l1BridgeAddress, l2Signer, l2Provider);

    console.log('[6/7] Relaying message on L2...');
    await step5RelayMessageOnL2(addresses.l2CrossChainMessengerAddress, l2Signer, crossChainMessage);

    console.log('[7/7] Waiting for wrapped token creation...');
    const l2TokenAddress = await step6GetWrappedTokenAddress(addresses.l2BridgeAddress, l2Signer, config.tokenAddress);

    console.log('[8/8] Registering wrapped token in NetworkConfig...');
    const l2TokenKey = `L2_${config.tokenSymbol}`;
    await step7RegisterWrappedToken(networkConfig, l2TokenKey, l2TokenAddress);

    console.log('Verifying registration via ten_config...');
    await verifyViaTenConfig(authenticatedL2RpcUrl, l2TokenKey, l2TokenAddress);

    console.log('Token whitelisting and registration completed successfully.');
    console.table({
        'L1 Token Address': config.tokenAddress,
        'L2 Token Address': l2TokenAddress,
        'Token Name': config.tokenName,
        'Token Symbol': config.tokenSymbol,
        'NetworkConfig Key': l2TokenKey,
        'L1 Bridge Address': addresses.l1BridgeAddress,
        'NetworkConfig Address': config.networkConfigAddr
    });
}

whitelistAndRegisterToken()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error('Error during whitelisting/registration:', error);
        process.exit(1);
    });

export default whitelistAndRegisterToken;
