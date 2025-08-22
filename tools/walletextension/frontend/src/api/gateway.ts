import { apiRoutes } from '../routes';
import { httpRequest } from './index';
import { pathToUrl } from '../routes/router';
import { AuthenticationResponse } from '@/types/interfaces/GatewayInterfaces';
import { tenGatewayAddress } from '../lib/constants';

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

export async function setTokenToCookie(token: string): Promise<void> {
    await httpRequest({
        method: 'post',
        url: tenGatewayAddress + pathToUrl(apiRoutes.setToken),
        data: { token },
        withCredentials: true,
    });
}

export async function getTokenFromCookie(): Promise<string> {
    try {
        const token = await httpRequest<string>({
            method: 'get',
            url: tenGatewayAddress + pathToUrl(apiRoutes.getToken),
            withCredentials: true,
        });

        // Validate token size and content
        if (!token || token.length === 0) {
            return '';
        }
        
        // Check if response is an error message instead of a token
        if (token.includes('cookie not found') || 
            token.includes('not found') || 
            token.includes('error') ||
            token.toLowerCase().includes('invalid')) {
            return '';
        }
        
        // Validate token format - should be hex (with or without 0x prefix)
        const cleanToken = token.startsWith('0x') ? token.slice(2) : token;
        if (!/^[0-9a-fA-F]+$/.test(cleanToken)) {
            return '';
        }
        
        return token;
    } catch (error) {
        return '';
    }
}

