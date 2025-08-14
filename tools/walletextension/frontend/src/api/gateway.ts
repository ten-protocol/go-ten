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

export async function getTokenFromCookie(): Promise<string> {
    console.log('[getTokenFromCookie] Attempting to retrieve token from cookie');
    try {
        const token = await httpRequest<string>({
            method: 'get',
            url: tenGatewayAddress + pathToUrl(apiRoutes.getToken),
        });
        
        // Validate token size and content
        if (!token || token.length === 0) {
            console.log('[getTokenFromCookie] Token is empty or null - treating as not set');
            return '';
        }
        
        // Check if token is valid hex and correct size (40 chars without 0x, 42 with 0x)
        const cleanToken = token.startsWith('0x') ? token.slice(2) : token;
        if (cleanToken.length !== 40 || !/^[0-9a-fA-F]+$/.test(cleanToken)) {
            console.log('[getTokenFromCookie] Token has invalid format or size:', token.length, '- treating as not set');
            return '';
        }
        
        console.log('[getTokenFromCookie] Successfully retrieved valid token from cookie');
        return token;
    } catch (error) {
        console.log('[getTokenFromCookie] Failed to retrieve token from cookie:', error);
        return '';
    }
}

export async function setTokenInCookie(token: string): Promise<boolean> {
    console.log('[setTokenInCookie] Attempting to set token in cookie');
    try {
        if (!token || token.length === 0) {
            console.log('[setTokenInCookie] Token is empty - cannot set');
            return false;
        }
        
        // Validate token format
        const cleanToken = token.startsWith('0x') ? token.slice(2) : token;
        if (cleanToken.length !== 40 || !/^[0-9a-fA-F]+$/.test(cleanToken)) {
            console.log('[setTokenInCookie] Invalid token format:', token);
            return false;
        }
        
        await httpRequest({
            method: 'post',
            url: tenGatewayAddress + pathToUrl(apiRoutes.setToken),
            data: { token },
        });
        
        console.log('[setTokenInCookie] Successfully set token in cookie');
        return true;
    } catch (error) {
        console.log('[setTokenInCookie] Failed to set token in cookie:', error);
        return false;
    }
}
