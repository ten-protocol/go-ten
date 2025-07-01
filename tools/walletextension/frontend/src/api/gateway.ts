import { apiRoutes } from '../routes';
import { httpRequest } from './index';
import { pathToUrl } from '../routes/router';
import { AuthenticationResponse } from '@/types/interfaces/GatewayInterfaces';
import { tenGatewayAddress, getUserIDMethodAddress } from '../lib/constants';

export async function fetchVersion(): Promise<string> {
    return await httpRequest<string>({
        method: 'get',
        url: tenGatewayAddress + pathToUrl(apiRoutes.version),
    });
}

export async function accountIsAuthenticated(
    token: string,
    account: string
): Promise<AuthenticationResponse> {
    return await httpRequest<AuthenticationResponse>({
        method: 'get',
        url: tenGatewayAddress + pathToUrl(apiRoutes.queryAccountToken),
        searchParams: {
            token,
            a: account,
        },
    });
}

export const authenticateUser = async (
    token: string,
    authenticateFields: {
        signature: string;
        address: string;
    }
) => {
    return await httpRequest({
        method: 'post',
        url: tenGatewayAddress + pathToUrl(apiRoutes.authenticate),
        data: authenticateFields,
        searchParams: {
            token,
        },
    });
};

export async function revokeAccountsApi(token: string): Promise<string> {
    return await httpRequest<string>({
        method: 'get',
        url: tenGatewayAddress + pathToUrl(apiRoutes.revoke),
        searchParams: {
            token,
        },
    });
}

export async function joinTestnet(): Promise<string> {
    return await httpRequest<string>({
        method: 'get',
        url: tenGatewayAddress + pathToUrl(apiRoutes.join),
    });
}

// Generate HMAC signature server-side
async function generateHMAC(timestamp: string): Promise<string> {
    const response = await httpRequest<{ signature: string }>({
        method: 'post',
        url: '/api/generate-hmac',
        data: { timestamp },
    });
    return response.signature;
}

// Convert string to hex (matches the HTML example)
function stringToHex(str: string): string {
    const encoder = new TextEncoder();
    const bytes = encoder.encode(str);
    return '0x' + Array.from(bytes)
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');
}

// Convert hex to string (for decoding the response)
function hexToString(hex: string): string {
    if (hex.startsWith('0x')) {
        hex = hex.slice(2);
    }
    const bytes: number[] = [];
    for (let i = 0; i < hex.length; i += 2) {
        bytes.push(parseInt(hex.substr(i, 2), 16));
    }
    const decoder = new TextDecoder();
    return decoder.decode(new Uint8Array(bytes));
}

export async function getUserID(token: string): Promise<string> {
    try {
        // Generate timestamp
        const timestamp = Date.now().toString();
        
        // Generate HMAC signature server-side
        const signature = await generateHMAC(timestamp);
        
        // Create authentication data (matches the HTML example)
        const authData = {
            timestamp: timestamp,
            signature: signature
        };
        
        const authDataJSON = JSON.stringify(authData);
        const storageKeyHex = stringToHex(authDataJSON);
        
        // Create RPC payload (matches the HTML example)
        const rpcPayload = {
            jsonrpc: "2.0",
            method: "eth_getStorageAt",
            params: [
                getUserIDMethodAddress,
                storageKeyHex,
                "latest"
            ],
            id: 1
        };
        
        // Send RPC request to backend
        const rpcResponse = await httpRequest<{ result: string; error?: any }>({
            method: 'post',
            url: `${tenGatewayAddress}/v1/?token=${token}`,
            data: rpcPayload,
        });
        
        if (rpcResponse.error) {
            throw new Error(`RPC error: ${rpcResponse.error.message}`);
        }
        
        const result = rpcResponse.result;
        
        if (!result || result === '0x' || result === '0x0') {
            throw new Error('Empty result from backend');
        }
        
        // Decode the hex response back to user ID string
        const userID = hexToString(result);
        return userID;
        
    } catch (error) {
        console.error('Failed to get user ID:', error);
        throw error;
    }
}
