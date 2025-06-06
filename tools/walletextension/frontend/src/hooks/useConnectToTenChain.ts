import { Connector, useAccount, useConnectors } from 'wagmi';
import { useEffect, useState } from 'react';
import {
    nativeCurrency,
    tenChainIDDecimal,
    tenChainIDHex,
    tenGatewayAddress,
    tenNetworkName,
} from '@/lib/constants';
import { useLocalStorage } from 'usehooks-ts';
import { joinTestnet } from '@/api/gateway';
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
    const [tenToken, setTenToken] = useLocalStorage<string>('ten_token', '');
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);
    const { isAuthenticated, isAuthenticatedLoading, authenticateAccount, authenticationError } =
        useTenChainAuth();

    const uniqueConnectors = connectors;

    const connectToTen = async (connector: Connector) => {
        setStep(1);
        setLoading(true);
        setError(null);
        setSelectedConnector(connector);
        await connector.connect().catch((error) => {
            setError({
                name: 'Unable to connect to wallet.',
                message: error.message,
                cause: error.cause,
            });
            setLoading(false);
            throw Error(error);
        });
    };

    useEffect(() => {
        if (step !== 1) return;

        async function switchToTen() {
            if (!tenToken && connector) {
                const chainExists = await chainExistsCheck(connector);

                if (chainExists) {
                    setError({
                        name: 'Existing chain found',
                        message:
                            'Found an existing TEN Protocol chain in your wallet. Please delete it and try again.',
                        cause: 'existing-chain-found',
                    });
                    return;
                }
            }

            const newTenToken =
                tenToken === ''
                    ? await joinTestnet().catch((error) => {
                          setError({
                              name: 'Unable to retrieve TEN token',
                              message: error.message,
                              cause: error.cause,
                          });
                          setLoading(false);
                      })
                    : tenToken;

            if (!newTenToken) {
                throw Error('No tenToken found');
            }

            setStep(2);

            setTenToken(`0x${newTenToken}`);
            setStoreTenToken(`0x${newTenToken}`);

            if (chainId === tenChainIDDecimal) {
                setStep(3);
                return;
            }
            let switchSuccess = true;

            if (!connector) {
                throw 'Connector is undefined!';
                return;
            }

            //@ts-expect-error Revisit later
            await connector
                .switchChain({
                    chainId: tenChainIDDecimal,
                    addEthereumChainParameter: {
                        rpcUrls: [`${tenGatewayAddress}/v1/?token=${newTenToken}`],
                        chainName: tenNetworkName,
                        nativeCurrency: nativeCurrency,
                    },
                })
                .catch((error: Error) => {
                    console.log(error);
                    if (error?.message.includes('is not a function')) {
                        console.log('IGNORE THIS ERROR');
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
                await sleep(500);
                setStep(3);
            }
        }

        if (isConnected && selectedConnector?.uid === connector?.uid) {
            switchToTen();
        }
    }, [connector, isConnected, selectedConnector, step]);

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
    }, [isAuthenticated, isAuthenticatedLoading, authenticationError, step]);

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
