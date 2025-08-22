import { tenGatewayVersion } from '../lib/constants';
import { NavLink } from '../types/interfaces';

export const NavLinks: NavLink[] = [];

export const apiRoutes = {
    join: `/${tenGatewayVersion}/join/`,
    authenticate: `/${tenGatewayVersion}/authenticate/`,
    queryAccountToken: `/${tenGatewayVersion}/query/`,
    revoke: `/${tenGatewayVersion}/revoke/`,
    version: `/${tenGatewayVersion}/version/`,
    getToken: `/${tenGatewayVersion}/get-token/`,
    setToken: `/${tenGatewayVersion}/set-token/`,

    // **** INFO ****
    getHealthStatus: `/${tenGatewayVersion}/network-health/`,
};

export const requestMethods = {
    connectAccounts: 'eth_requestAccounts',
    switchNetwork: 'wallet_switchEthereumChain',
    addNetwork: 'wallet_addEthereumChain',
    getStorageAt: 'eth_getStorageAt',
    signTypedData: 'eth_signTypedData_v4',
};
