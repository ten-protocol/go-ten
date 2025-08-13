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
    return await httpRequest<void>({
        method: 'post',
        url: tenGatewayAddress + pathToUrl(apiRoutes.setToken),
        data: { token },
    });
}
