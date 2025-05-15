import { useEffect, useState } from 'react';
import { useAccount, useSignTypedData } from 'wagmi';
import { accountIsAuthenticated, authenticateUser, revokeAccountsApi } from '@/api/gateway';
import { Address, getAddress } from 'viem';
import { useMutation, useQuery } from '@tanstack/react-query';
import { generateEIP712 } from '@/lib/eip712';
import { useUiStore } from '@/stores/ui.store';
import { useLocalStorage } from 'usehooks-ts';

export function useTenChainAuth(walletAddress?: Address) {
    const { connector } = useAccount();
    const authEvents = useUiStore((state) => state.authEvents);
    const [address, setAddress] = useState<Address | undefined>(walletAddress);
    const [tenToken] = useLocalStorage<string | null>('ten_token', null);
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
        queryKey: [`beAuthCheck`, address, signature, authEvents],
        queryFn: () => accountIsAuthenticated(tenToken ?? '', address ?? ''),
        enabled: !!address,
    });

    const authenticationMutation = useMutation({
        mutationFn: (signature) => {
            return authenticateWalletWithBE(signature);
        },
    });

    const getSignature = async () => {
        if (!tenToken) return null;

        const tokenValue = tenToken.startsWith('0x') ? tenToken : '0x' + tenToken;

        try {
            const checksummedAddress = getAddress(tokenValue);
            signTypedData(generateEIP712(checksummedAddress));
        } catch (err) {
            throw err;
        }
    };

    const authenticateWalletWithBE = async (signature: string) => {
        if (!tenToken) return null;

        const response = await authenticateUser(tenToken, {
            signature,
            address,
        });

        if (response === 'internal error') {
            setIsLoading(false);
            throw Error('Unable to authenticate wallet with the network.');
        } else if (response === 'success') {
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
            throw Error('Ten token not found.');
        }

        await revokeAccountsApi(tenToken);
    };

    return {
        isAuthenticated: !!beAuthCheck?.status,
        isAuthenticatedLoading: isLoading,
        authenticateAccount,
        revokeAccount,
        beAuthCheckError,
        authenticationError: authenticationMutation.error || signError,
    };
}
