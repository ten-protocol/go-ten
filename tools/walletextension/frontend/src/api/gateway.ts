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
    _token: string, // Not used - we get token from cookie instead
    account: string
): Promise<AuthenticationResponse> {
    console.log('🔍 accountIsAuthenticated: Starting cookie-based authentication check');
    console.log('📝 accountIsAuthenticated: account =', account);
    
    try {
        // First, get the token from the cookie via /get-token endpoint
        console.log('🍪 accountIsAuthenticated: Getting token from cookie...');
        const tokenFromCookie = await getTokenFromCookie();
        console.log('🍪 accountIsAuthenticated: Token from cookie =', tokenFromCookie);
        
        // Then use that token to query authentication status
        console.log('🔍 accountIsAuthenticated: Calling /query with token from cookie');
        console.log('📍 accountIsAuthenticated: URL =', tenGatewayAddress + pathToUrl(apiRoutes.queryAccountToken));
        
        const result = await httpRequest<AuthenticationResponse>({
            method: 'get',
            url: tenGatewayAddress + pathToUrl(apiRoutes.queryAccountToken),
            searchParams: {
                token: tokenFromCookie, // Use token from cookie
                a: account,
            },
        });
        
        console.log('✅ accountIsAuthenticated: Success! Response =', result);
        console.log('📊 accountIsAuthenticated: Status =', result?.status);
        
        return result;
    } catch (error) {
        console.error('❌ accountIsAuthenticated: Error occurred:', error);
        throw error;
    }
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
        
        // Verify the token was actually stored by immediately reading it back
        console.log('🔍 setTokenToCookie: Verifying token was stored correctly...');
        const verifyToken = await getTokenFromCookie();
        console.log('🔍 setTokenToCookie: Token verification - stored:', token);
        console.log('🔍 setTokenToCookie: Token verification - retrieved:', verifyToken);
        console.log('🔍 setTokenToCookie: Token verification - match:', token === verifyToken);
        
        return result;
    } catch (error) {
        console.error('❌ setTokenToCookie: Error storing token in cookie:', error);
        throw error;
    }
}
