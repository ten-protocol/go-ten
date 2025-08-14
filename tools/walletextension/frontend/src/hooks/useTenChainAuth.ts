import { useEffect, useState } from 'react';
import { useSignTypedData } from 'wagmi';
import { accountIsAuthenticated, authenticateUser, revokeAccountsApi } from '@/api/gateway';
import { Address, getAddress } from 'viem';
import { useMutation, useQuery } from '@tanstack/react-query';
import { generateEIP712 } from '@/lib/eip712';
import { useUiStore } from '@/stores/ui.store';
import { useTenToken } from '@/hooks/useTenToken';

export function useTenChainAuth(walletAddress?: Address) {
    const authEvents = useUiStore((state) => state.authEvents);
    const [address, setAddress] = useState<Address | undefined>(walletAddress);
    const { token: tenToken, loading: tenTokenLoading } = useTenToken();
    const {
        data: signature,
        signTypedData,
        error: signError,
        isSuccess: signSuccess,
    } = useSignTypedData();
    const [isLoading, setIsLoading] = useState<boolean>(false);

    const {
        data: beAuthCheck,
        error: beAuthCheckError,
        refetch: beAuthCheckRefetch,
    } = useQuery({
        queryKey: [`beAuthCheck`, address, signature, authEvents, tenToken],
        queryFn: () => {
            console.log('[useTenChainAuth] beAuthCheck queryFn - tenToken:', tenToken, 'address:', address);
            return accountIsAuthenticated(tenToken || '', address ?? '');
        },
        enabled: !!address && !!tenToken && !tenTokenLoading,
    });

    const authenticationMutation = useMutation({
        mutationFn: (signature: string) => {
            return authenticateWalletWithBE(signature);
        },
    });

    const getSignature = async () => {
        console.log('[useTenChainAuth] getSignature called - tenToken value:', tenToken, 'tenTokenLoading:', tenTokenLoading, 'type:', typeof tenToken);
        
        if (!tenToken) {
            console.log('[useTenChainAuth] No tenToken available for signing - tenToken is:', tenToken);
            return null;
        }

        const tokenValue = tenToken.startsWith('0x') ? tenToken : '0x' + tenToken;
        console.log('[useTenChainAuth] Getting signature for token:', tokenValue);

        try {
            const checksummedAddress = getAddress(tokenValue);
            signTypedData(generateEIP712(checksummedAddress));
        } catch (err) {
            console.log('[useTenChainAuth] Error getting signature:', err);
            throw err;
        }
    };

    const authenticateWalletWithBE = async (signature: string) => {
        if (!tenToken || !address) {
            console.log('[useTenChainAuth] Missing tenToken or address for authentication');
            return null;
        }

        console.log('[useTenChainAuth] Authenticating wallet with backend');
        const response = await authenticateUser(tenToken, {
            signature,
            address,
        });

        if (response === 'internal error') {
            console.log('[useTenChainAuth] Authentication failed - internal error');
            setIsLoading(false);
            throw Error('Unable to authenticate wallet with the network.');
        } else if (response === 'success') {
            console.log('[useTenChainAuth] Authentication successful');
            await beAuthCheckRefetch();
            setIsLoading(false);
        }
    };

    const authenticateAccount = async (walletAddress?: Address) => {
        setIsLoading(true);

        if (walletAddress) {
            setAddress(walletAddress);
        }
        await getSignature();
    };

    useEffect(() => {
        if (signSuccess) {
            authenticationMutation.mutate(signature);
        }
    }, [signTypedData, signSuccess]);

    useEffect(() => {
        if (authenticationMutation.isSuccess || authenticationMutation.isError) {
            setIsLoading(false);
        }
    }, [authenticationMutation.status]);

    const revokeAccount = async () => {
        if (!tenToken) {
            console.log('[useTenChainAuth] Ten token not found for revocation');
            throw Error('Ten token not found.');
        }

        console.log('[useTenChainAuth] Revoking account');
        await revokeAccountsApi(tenToken);
    };

    return {
        isAuthenticated: !!beAuthCheck?.status,
        isAuthenticatedLoading: isLoading || tenTokenLoading,
        authenticateAccount,
        revokeAccount,
        beAuthCheckError,
        authenticationError: authenticationMutation.error || signError,
    };
}
