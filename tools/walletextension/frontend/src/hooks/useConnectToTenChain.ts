import { Connector, useAccount, useConnectors } from 'wagmi';
import { useEffect, useState } from 'react';
import {
    nativeCurrency,
    tenChainIDDecimal,
    tenChainIDHex,
    tenGatewayAddress,
    tenNetworkName,
} from '@/lib/constants';
import { joinTestnet } from '@/api/gateway';
import { useTenToken } from '@/hooks/useTenToken';
import { useTenChainAuth } from '@/hooks/useTenChainAuth';
import { useUiStore } from '@/stores/ui.store';
import sleep from '@/utils/sleep';
import { shallow } from 'zustand/shallow';

export default function useConnectToTenChain() {
    const incrementAuthEvents = useUiStore((state) => state.incrementAuthEvents, shallow);
    const { address, connector, isConnected, chainId } = useAccount();
    const connectors = useConnectors();
    const setStoreTenToken = useUiStore((state) => state.setTenToken);
    const [step, setStep] = useState<number>(0);
    const [selectedConnector, setSelectedConnector] = useState<Connector | null>(null);
    const { token: tenToken, setToken: setTenToken, loading: tokenLoading } = useTenToken();
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
            let newTenToken =
                tenToken === ''
                    ? await joinTestnet().catch((error) => {
                          console.log('[useConnectToTenChain] Error joining testnet:', error);
                          setError({
                              name: 'Unable to retrieve TEN token',
                              message: error.message,
                              cause: error.cause,
                          });
                          setLoading(false);
                      })
                    : tenToken;

            console.log('[useConnectToTenChain] New token obtained:', newTenToken ? 'success' : 'failed');
            if (!newTenToken) {
                throw Error('No tenToken found');
            }

            console.log('[useConnectToTenChain] Moving to step 2');
            setStep(2);

            if (tenToken === '') {
                newTenToken = `0x${newTenToken}`;
                console.log('[useConnectToTenChain] Setting new token in cookie:', newTenToken);
                await setTenToken(newTenToken);
                setStoreTenToken(newTenToken);
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
            const rpcUrl = `${tenGatewayAddress}/v1/?token=${cleanTokenForRpc}`;
            console.log('[useConnectToTenChain] Attempting to switch chain with RPC URL:', rpcUrl);
            
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
    }, [connector, isConnected, selectedConnector, step, tokenLoading, tenToken, setTenToken, setStoreTenToken, chainId]);

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
            console.log('[useConnectToTenChain] Starting authentication for address:', address);
            authenticateAccount(address);
        } else if (isAuthenticated && !isAuthenticatedLoading) {
            console.log('[useConnectToTenChain] Authentication successful, moving to step 4');
            setStep(4);
            incrementAuthEvents();
        }
    }, [isAuthenticated, isAuthenticatedLoading, authenticationError, step, address, authenticateAccount, incrementAuthEvents]);

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
