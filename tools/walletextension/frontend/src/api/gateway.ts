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
    console.log('🚀 joinTestnet: Starting /join API call');
    console.log('📍 joinTestnet: URL =', tenGatewayAddress + pathToUrl(apiRoutes.join));
    
    try {
        const result = await httpRequest<string>(
            {
                method: 'get',
                url: tenGatewayAddress + pathToUrl(apiRoutes.join),
            },
            {
                withCredentials: false, // Don't send existing cookies for fresh token generation
            }
        );
        
        console.log('✅ joinTestnet: Success! Received token =', result);
        console.log('📊 joinTestnet: Token length =', result?.length);
        console.log('🔤 joinTestnet: Token type =', typeof result);
        
        return result;
    } catch (error) {
        console.error('❌ joinTestnet: Error occurred:', error);
        throw error;
    }
}

export async function getTokenFromCookie(): Promise<string> {
    const response = await httpRequest<string>({
        method: 'get',
        url: tenGatewayAddress + pathToUrl(apiRoutes.getToken),
    });
    
    // Check if response indicates cookie not found
    if (response === 'gateway_token cookie not found') {
        throw new Error('gateway_token cookie not found');
    }
    
    return response;
}

export async function setTokenToCookie(token: string): Promise<void> {
    console.log('🍪 setTokenToCookie: About to store token in cookie =', token);
    console.log('📍 setTokenToCookie: URL =', tenGatewayAddress + pathToUrl(apiRoutes.setToken));
    
    try {
        const result = await httpRequest<void>({
            method: 'post',
            url: tenGatewayAddress + pathToUrl(apiRoutes.setToken),
            data: { token },
        });
        
        console.log('✅ setTokenToCookie: Successfully stored token in cookie');
        return result;
    } catch (error) {
        console.error('❌ setTokenToCookie: Error storing token in cookie:', error);
        throw error;
    }
}
