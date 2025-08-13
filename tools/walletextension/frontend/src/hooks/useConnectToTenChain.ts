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
import { useTokenFromCookie } from './useTokenFromCookie';
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
    const [tenToken, setTenTokenToCookie, isTokenLoading] = useTokenFromCookie();
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);
    const { isAuthenticated, isAuthenticatedLoading, authenticateAccount, authenticationError } =
        useTenChainAuth();

    const uniqueConnectors = connectors;

    const connectToTen = async (connector: Connector) => {
        console.log('🚀 connectToTen: Starting connection process');
        console.log('🔌 connectToTen: Connector =', connector.name, connector.id);
        console.log('📊 connectToTen: Current step =', step);
        
        setStep(1);
        setLoading(true);
        setError(null);
        setSelectedConnector(connector);
        
        console.log('📈 connectToTen: Set step to 1, loading to true');
        console.log('🔗 connectToTen: About to call connector.connect()');
        
        await connector.connect().catch((error) => {
            console.error('❌ connectToTen: connector.connect() failed:', error);
            setError({
                name: 'Unable to connect to wallet.',
                message: error.message,
                cause: error.cause,
            });
            setLoading(false);
            throw Error(error);
        });
        
        console.log('✅ connectToTen: connector.connect() successful');
    };

    useEffect(() => {
        console.log('🔄 useEffect[step=1]: Triggered');
        console.log('📊 useEffect[step=1]: step =', step);
        console.log('🔌 useEffect[step=1]: isConnected =', isConnected);
        console.log('🔌 useEffect[step=1]: selectedConnector =', selectedConnector?.name, selectedConnector?.uid);
        console.log('🔌 useEffect[step=1]: connector =', connector?.name, connector?.uid);
        
        if (step !== 1) {
            console.log('⏭️ useEffect[step=1]: step !== 1, returning');
            return;
        }

        async function switchToTen() {
            console.log('🔄 switchToTen: Starting network addition flow');
            console.log('📝 switchToTen: Current tenToken from hook =', tenToken);
            console.log('📝 switchToTen: Is token loading =', isTokenLoading);
            
            // Check if TEN chain already exists in wallet
            if (connector) {
                console.log('🔍 switchToTen: Checking if TEN chain already exists...');
                const chainExists = await chainExistsCheck(connector);
                console.log('🔍 switchToTen: Chain exists check result =', chainExists);

                if (chainExists) {
                    console.log('⚠️ switchToTen: Existing chain found, showing error');
                    setError({
                        name: 'Existing chain found',
                        message:
                            'Found an existing TEN Protocol chain in your wallet. Please delete it and try again.',
                        cause: 'existing-chain-found',
                    });
                    return;
                }
            }

            // Always call /join to get a fresh token when adding new network
            console.log('🎯 switchToTen: About to call joinTestnet()...');
            const newTenToken = await joinTestnet().catch((error) => {
                console.error('❌ switchToTen: joinTestnet() failed:', error);
                setError({
                    name: 'Unable to retrieve TEN token',
                    message: error.message,
                    cause: error.cause,
                });
                setLoading(false);
                return null;
            });

            console.log('🎯 switchToTen: joinTestnet() returned =', newTenToken);

            if (!newTenToken) {
                console.error('❌ switchToTen: No token received from joinTestnet()');
                throw Error('No tenToken found');
            }

            setStep(2);
            console.log('📈 switchToTen: Set step to 2');

            // Store the new token in cookie
            console.log('🍪 switchToTen: About to store token in cookie...');
            await setTenTokenToCookie(newTenToken);
            console.log('🍪 switchToTen: Token stored in cookie');
            
            console.log('🏪 switchToTen: About to store token in UI store...');
            setStoreTenToken(newTenToken);
            console.log('🏪 switchToTen: Token stored in UI store');

            if (chainId === tenChainIDDecimal) {
                console.log('✅ switchToTen: Already on TEN chain, skipping to step 3');
                setStep(3);
                return;
            }
            
            let switchSuccess = true;

            if (!connector) {
                console.error('❌ switchToTen: No connector available');
                throw 'Connector is undefined!';
                return;
            }

            const rpcUrl = `${tenGatewayAddress}/v1/?token=${newTenToken}`;
            console.log('🌐 switchToTen: About to add network with RPC URL =', rpcUrl);
            console.log('🌐 switchToTen: Chain ID =', tenChainIDDecimal);
            console.log('🌐 switchToTen: Chain Name =', tenNetworkName);

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
                    console.error('❌ switchToTen: switchChain failed:', error);
                    if (error?.message.includes('is not a function')) {
                        console.log('ℹ️ switchToTen: Ignoring "is not a function" error');
                    } else {
                        switchSuccess = false;
                        setError({
                            name: 'Error switching chains',
                            message: error.message,
                            cause: error.cause,
                        });
                    }
                });

            if (switchSuccess) {
                console.log('✅ switchToTen: Network addition successful, proceeding to step 3');
                await sleep(500);
                setStep(3);
            } else {
                console.error('❌ switchToTen: Network addition failed, staying on current step');
            }
        }

        if (isConnected && selectedConnector?.uid === connector?.uid) {
            console.log('✅ useEffect[step=1]: Conditions met, calling switchToTen()');
            switchToTen();
        } else {
            console.log('❌ useEffect[step=1]: Conditions not met');
            console.log('   isConnected =', isConnected);
            console.log('   selectedConnector?.uid =', selectedConnector?.uid);
            console.log('   connector?.uid =', connector?.uid);
            console.log('   UIDs match =', selectedConnector?.uid === connector?.uid);
        }
    }, [connector, isConnected, selectedConnector, step, chainId, tenToken, setStoreTenToken, setTenTokenToCookie, isTokenLoading]);

    useEffect(() => {
        if (step !== 3) {
            return;
        }

        if (authenticationError) {
            setError({
                name: 'Error authenticating token',
                message: authenticationError.message,
                cause: authenticationError.cause,
            });
        }

        if (!isAuthenticated && !isAuthenticatedLoading && !authenticationError) {
            authenticateAccount(address);
        } else if (isAuthenticated && !isAuthenticatedLoading) {
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
        loading,
    };
}
