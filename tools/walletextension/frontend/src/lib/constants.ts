export const tenGatewayAddress = process.env.NEXT_PUBLIC_GATEWAY_URL;

export const tenNetworkName = process.env.NEXT_PUBLIC_NETWORK_NAME || 'TEN Testnet';

export const tenscanAddress = process.env.NEXT_PUBLIC_TENSCAN_URL || 'https://tenscan.io';

export const socialLinks = {
    faucet: 'https://faucet.ten.xyz',
    github: 'https://github.com/ten-protocol',
    discord: 'https://discord.gg/tenprotocol',
    twitter: 'https://twitter.com/tenprotocol',
    twitterHandle: '@tenprotocol',
};

export const GOOGLE_ANALYTICS_ID = process.env.NEXT_PUBLIC_GOOGLE_ANALYTICS_ID;

export const SWITCHED_CODE = 4902;
export const tokenHexLength = 42;

export const tenGatewayVersion = 'v1';
export const tenChainIDDecimal = 443;

export const environment = process.env.NEXT_PUBLIC_ENVIRONMENT;

export const tenChainIDHex = '0x' + tenChainIDDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'

export const METAMASK_CONNECTION_TIMEOUT = 3000;

export const userStorageAddress = '0x0000000000000000000000000000000000000001';

// Server-side only - this will NOT be exposed to the client
export const HMAC_SECRET = process.env.HMAC_SECRET;

// GetUserID custom query method address
export const getUserIDMethodAddress = '0x0000000000000000000000000000000000000009';

export const nativeCurrency = {
    name: 'Sepolia Ether',
    symbol: 'ETH',
    decimals: 18,
};

export const CONNECTION_STEPS = [
    'Hit Connect to TEN and start your journey',
    'Allow MetaMask to switch networks to the TEN Testnet',
    'Sign the <b>Signature Request</b> (this is not a transaction)',
];
