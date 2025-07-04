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

/**
 * Custom hook for connecting to the Ten Protocol blockchain
 * Manages the complete flow from wallet connection to chain switching and authentication
 */
export default function useConnectToTenChain() {
    // UI store actions for tracking authentication events
    const incrementAuthEvents = useUiStore((state) => state.incrementAuthEvents, shallow);
    
    // Wagmi hooks for wallet connection state
    const { address, connector, isConnected, chainId } = useAccount();
    const connectors = useConnectors();
    
    // UI store for storing Ten token globally
    const setStoreTenToken = useUiStore((state) => state.setTenToken);
    
    // Local state management
    const [step, setStep] = useState<number>(0); // Tracks connection progress (0-4)
    const [selectedConnector, setSelectedConnector] = useState<Connector | null>(null); // Currently selected wallet connector
    const [tenToken, setTenToken] = useLocalStorage<string>('ten_token', ''); // Persisted Ten authentication token
    const [loading, setLoading] = useState<boolean>(false); // Loading state for UI feedback
    const [error, setError] = useState<Error | null>(null); // Error state for handling failures
    
    // Ten Chain authentication hook
    const { isAuthenticated, isAuthenticatedLoading, authenticateAccount, authenticationError } =
        useTenChainAuth();

    // Available wallet connectors (could be filtered for specific wallets)
    const uniqueConnectors = connectors;

    /**
     * Initiates connection to a specific wallet connector
     * @param connector - The wallet connector to connect to
     */
    const connectToTen = async (connector: Connector) => {
        setStep(1); // Move to wallet connection step
        setLoading(true);
        setError(null);
        setSelectedConnector(connector);
        
        // Attempt to connect to the wallet
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

    /**
     * Effect: Handles chain switching to Ten Protocol after wallet connection
     * Triggered when step === 1 and wallet is connected
     */
    useEffect(() => {
        if (step !== 1) return;

        async function switchToTen() {
            // Check for existing Ten chain if no token exists
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

            // Get or use existing Ten token
            let newTenToken =
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

            setStep(2); // Move to chain switching step

            // Format and store new token if it's fresh
            if (tenToken === '') {
                newTenToken = `0x${newTenToken}`;
                setTenToken(newTenToken);
                setStoreTenToken(`0x${newTenToken}`);
            }

            // Skip chain switching if already on Ten chain
            if (chainId === tenChainIDDecimal) {
                setStep(3);
                return;
            }
            
            let switchSuccess = true;

            if (!connector) {
                throw 'Connector is undefined!';
                return;
            }

            // Attempt to switch to Ten chain or add it if it doesn't exist
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
                    // Ignore specific function not found errors (likely from wallet implementation)
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

            // Proceed to authentication if chain switch was successful
            if (switchSuccess) {
                await sleep(500); // Brief delay for chain switch to complete
                setStep(3);
            }
        }

        // Only proceed if wallet is connected and we have the correct connector
        if (isConnected && selectedConnector?.uid === connector?.uid) {
            switchToTen();
        }
    }, [connector, isConnected, selectedConnector, step]);

    /**
     * Effect: Handles Ten Chain authentication after successful chain switch
     * Triggered when step === 3
     */
    useEffect(() => {
        if (step !== 3) {
            return;
        }

        // Handle authentication errors
        if (authenticationError) {
            setError({
                name: 'Error authenticating token',
                message: authenticationError.message,
                cause: authenticationError.cause,
            });
        }

        // Start authentication if not already authenticated and no errors
        if (!isAuthenticated && !isAuthenticatedLoading && !authenticationError) {
            authenticateAccount(address);
        } else if (isAuthenticated && !isAuthenticatedLoading) {
            // Authentication successful - complete the flow
            setStep(4);
            incrementAuthEvents(); // Track successful authentication
        }
    }, [isAuthenticated, isAuthenticatedLoading, authenticationError, step]);

    /**
     * Resets the connection state to initial values
     * Useful for retrying connections or starting fresh
     */
    const reset = () => {
        setStep(0);
        setError(null);
        setLoading(false);
    };

    /**
     * Checks if Ten chain already exists in the user's wallet
     * @param connector - The wallet connector to check
     * @returns Promise<boolean> - True if chain exists, false otherwise
     */
    const chainExistsCheck = async (connector: Connector) => {
        const provider = await connector.getProvider();
        let chainExists = true;

        if (!provider) {
            throw 'Provider not found!';
        }

        // Try to switch to Ten chain - if it fails with specific error codes, chain doesn't exist
        await provider
            //@ts-expect-error Revisit later
            .request({
                method: 'wallet_switchEthereumChain',
                params: [{ chainId: tenChainIDHex }],
            })
            .then(() => true)
            .catch((error: { code: number }) => {
                // Error codes: 4902 = chain not added, -32603 = chain not found
                if (error.code === 4902 || error.code === -32603) {
                    chainExists = false;
                }
            });

        return chainExists;
    };

    // Return hook interface
    return {
        step, // Current step in connection process (0-4)
        error, // Current error state
        errors: [], // Array of errors (currently unused)
        connectors: uniqueConnectors, // Available wallet connectors
        connectToTen, // Function to initiate connection
        reset, // Function to reset state
        loading, // Loading state for UI
    };
}
