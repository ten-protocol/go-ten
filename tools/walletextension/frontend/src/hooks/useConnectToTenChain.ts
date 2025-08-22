import { Connector, useAccount, useConnectors } from 'wagmi';
import { useEffect, useState } from 'react';
import {
    nativeCurrency,
    tenChainIDDecimal,
    tenChainIDHex,
    tenGatewayAddress,
    tenNetworkName,
} from '@/lib/constants';
import { joinTestnet, setTokenToCookie } from '@/api/gateway';
import { useTenToken } from '@/contexts/TenTokenContext';
import { useTenChainAuth } from '@/hooks/useTenChainAuth';
import { useUiStore } from '@/stores/ui.store';
import sleep from '@/utils/sleep';
import { shallow } from 'zustand/shallow';

export default function useConnectToTenChain() {
    const incrementAuthEvents = useUiStore((state) => state.incrementAuthEvents, shallow);
    const { address, connector, isConnected, chainId } = useAccount();
    const connectors = useConnectors();
    const [step, setStep] = useState<number>(0);
    const [selectedConnector, setSelectedConnector] = useState<Connector | null>(null);
    const { token: tenToken, loading: tokenLoading, refreshToken } = useTenToken();
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);
    const { isAuthenticated, isAuthenticatedLoading, authenticateAccount, authenticationError } =
        useTenChainAuth();

    const uniqueConnectors = connectors;

    const connectToTen = async (connector: Connector) => {
        console.log('[useConnectToTenChain] Starting connectToTen with connector:', connector.name);
        setStep(1);
        setLoading(true);
        setError(null);
        setSelectedConnector(connector);
        console.log('[useConnectToTenChain] Attempting to connect to wallet...');
        await connector.connect().catch((error) => {
            console.log('[useConnectToTenChain] Wallet connection failed:', error);
            setError({
                name: 'Unable to connect to wallet.',
                message: error.message,
                cause: error.cause,
            });
            setLoading(false);
            throw Error(error);
        });
        console.log('[useConnectToTenChain] Wallet connection initiated successfully');
    };

    useEffect(() => {
        console.log('[useConnectToTenChain] useEffect triggered - step:', step, 'tokenLoading:', tokenLoading);
        if (step !== 1) return;
        if (tokenLoading) {
            console.log('[useConnectToTenChain] Token still loading, deferring switchToTen');
            return;
        }

        async function switchToTen() {
            console.log('[useConnectToTenChain] Starting switchToTen, tenToken:', tenToken ? 'present' : 'empty', 'connector:', connector?.name);
            
            if (!tenToken && connector) {
                console.log('[useConnectToTenChain] No token found, checking if chain exists');
                const chainExists = await chainExistsCheck(connector);

                if (chainExists) {
                    console.log('[useConnectToTenChain] Existing chain found, showing error');
                    setError({
                        name: 'Existing chain found',
                        message:
                            'Found an existing TEN Protocol chain in your wallet. Please delete it and try again.',
                        cause: 'existing-chain-found',
                    });
                    return;
                }
            }

            console.log('[useConnectToTenChain] Getting token - current tenToken:', tenToken);
            let newTenToken: string;
            if (tenToken === '') {
                console.log('[useConnectToTenChain] No existing token found, calling /join endpoint');
                console.log('[useConnectToTenChain] /join call context - step:', step, 'connector:', connector?.name, 'address:', address);
                const startTime = Date.now();
                
                try {
                    newTenToken = await joinTestnet();
                    const endTime = Date.now();
                    const duration = endTime - startTime;
                    console.log('[useConnectToTenChain] /join endpoint SUCCESS - duration:', duration + 'ms', 'token received:', newTenToken ? 'yes' : 'no');
                    console.log('[useConnectToTenChain] /join response full token:', newTenToken);
                    console.log('[useConnectToTenChain] /join response token length:', newTenToken?.length || 0, 'first 10 chars:', newTenToken?.substring(0, 10) || 'N/A');
                    
                    // Set the token to cookie using /set-token endpoint
                    console.log('[useConnectToTenChain] Setting token to cookie via /set-token...');
                    await setTokenToCookie(newTenToken);
                    console.log('[useConnectToTenChain] Token set to cookie successfully');
                    
                    // Refresh the token context to sync with the cookie
                    console.log('[useConnectToTenChain] Refreshing token context...');
                    await refreshToken();
                    console.log('[useConnectToTenChain] Token context refreshed');
                    
                    console.log('[useConnectToTenChain] Using token from /join response for RPC URL');
                } catch (error: any) {
                    const endTime = Date.now();
                    const duration = endTime - startTime;
                    console.log('[useConnectToTenChain] /join endpoint FAILED - duration:', duration + 'ms');
                    console.log('[useConnectToTenChain] /join error details:', {
                        message: error?.message,
                        cause: error?.cause,
                        stack: error?.stack,
                        name: error?.name
                    });
                    setError({
                        name: 'Unable to retrieve TEN token',
                        message: error?.message || 'Unknown error',
                        cause: error?.cause,
                    });
                    setLoading(false);
                    return; // Exit early on error
                }
            } else {
                console.log('[useConnectToTenChain] Using existing token, skipping /join endpoint call');
                newTenToken = tenToken;
            }

            console.log('[useConnectToTenChain] New token obtained:', newTenToken ? 'success' : 'failed');
            if (!newTenToken) {
                throw Error('No tenToken found');
            }

            console.log('[useConnectToTenChain] Moving to step 2');
            setStep(2);

            // Ensure token has 0x prefix for processing
            if (!newTenToken.startsWith('0x')) {
                newTenToken = `0x${newTenToken}`;
                console.log('[useConnectToTenChain] Added 0x prefix to token:', newTenToken.substring(0, 10) + '...');
            }

            if (chainId === tenChainIDDecimal) {
                console.log('[useConnectToTenChain] Already on correct chain, moving to step 3');
                setStep(3);
                return;
            }
            
            console.log('[useConnectToTenChain] Wrong chain detected, current chainId:', chainId, 'expected:', tenChainIDDecimal);
            let switchSuccess = true;

            if (!connector) {
                console.log('[useConnectToTenChain] ERROR: Connector is undefined!');
                throw 'Connector is undefined!';
                return;
            }

            // Remove 0x prefix from token for RPC URL
            const cleanTokenForRpc = newTenToken.startsWith('0x') ? newTenToken.slice(2) : newTenToken;
            console.log('[useConnectToTenChain] Token for RPC URL - original:', newTenToken.substring(0, 10) + '...', 'clean:', cleanTokenForRpc.substring(0, 10) + '...');
            const rpcUrl = `${tenGatewayAddress}/v1/?token=${cleanTokenForRpc}`;
            console.log('[useConnectToTenChain] Attempting to switch chain with RPC URL:', rpcUrl.replace(/token=[^&]*/, 'token=***HIDDEN***'));
            
            //@ts-expect-error Revisit later
            await connector
                .switchChain({
                    chainId: tenChainIDDecimal,
                    addEthereumChainParameter: {
                        rpcUrls: [rpcUrl],
                        chainName: tenNetworkName,
                        nativeCurrency: nativeCurrency,
                    },
                })
                .catch((error: Error) => {
                    console.log('[useConnectToTenChain] Chain switch error:', error);
                    if (error?.message.includes('is not a function')) {
                        console.log('[useConnectToTenChain] Ignoring "is not a function" error');
                    } else {
                        console.log('[useConnectToTenChain] Setting switchSuccess to false due to error');
                        switchSuccess = false;
                        setError({
                            name: 'Error switching chains',
                            message: error.message,
                            cause: error.cause,
                        });
                    }
                });

            if (switchSuccess) {
                console.log('[useConnectToTenChain] Chain switch successful, waiting 500ms then moving to step 3');
                await sleep(500);
                setStep(3);
            } else {
                console.log('[useConnectToTenChain] Chain switch failed, staying in current step');
            }
        }

        if (isConnected && selectedConnector?.uid === connector?.uid) {
            console.log('[useConnectToTenChain] Conditions met for switchToTen - isConnected:', isConnected, 'selectedConnector.uid:', selectedConnector?.uid, 'connector.uid:', connector?.uid, 'tokenLoading:', tokenLoading);
            if (!tokenLoading) {
                switchToTen();
            } else {
                console.log('[useConnectToTenChain] Token still loading, waiting...');
            }
        }
    }, [connector, isConnected, selectedConnector, step, tokenLoading, tenToken, chainId]);

    useEffect(() => {
        if (step !== 3) {
            return;
        }

        console.log('[useConnectToTenChain] Step 3 - Authentication phase. isAuthenticated:', isAuthenticated, 'isAuthenticatedLoading:', isAuthenticatedLoading, 'authenticationError:', authenticationError, 'address:', address);

        if (authenticationError) {
            console.log('[useConnectToTenChain] Authentication error detected:', authenticationError);
            setError({
                name: 'Error authenticating token',
                message: authenticationError.message,
                cause: authenticationError.cause,
            });
        }

        if (!isAuthenticated && !isAuthenticatedLoading && !authenticationError) {
            console.log('[useConnectToTenChain] Starting authentication for address:', address, 'with tenToken:', tenToken);
            if (!tenToken) {
                console.log('[useConnectToTenChain] ERROR: Attempting to authenticate but no tenToken available!');
                return;
            }
            authenticateAccount(address);
        } else if (isAuthenticated && !isAuthenticatedLoading) {
            console.log('[useConnectToTenChain] Authentication successful, moving to step 4');
            setStep(4);
            incrementAuthEvents();
        }
    }, [isAuthenticated, isAuthenticatedLoading, authenticationError, step, address, authenticateAccount, incrementAuthEvents, tenToken]);

    const reset = () => {
        setStep(0);
        setError(null);
        setLoading(false);
    };

    const chainExistsCheck = async (connector: Connector) => {
        const provider = await connector.getProvider();
        let chainExists = true;

        if (!provider) {
            throw 'Provider not found!';
        }

        await provider
            //@ts-expect-error Revisit later
            .request({
                method: 'wallet_switchEthereumChain',
                params: [{ chainId: tenChainIDHex }],
            })
            .then(() => true)
            .catch((error: { code: number }) => {
                if (error.code === 4902 || error.code === -32603) {
                    chainExists = false;
                }
            });

        return chainExists;
    };

    return {
        step,
        error,
        errors: [],
        connectors: uniqueConnectors,
        connectToTen,
        reset,
        loading: loading || tokenLoading,
    };
}
