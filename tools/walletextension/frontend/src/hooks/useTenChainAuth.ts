import {useLocalStorage} from "@/hooks/useLocalStorage";
import {useEffect, useState} from "react";
import { useSignTypedData} from "wagmi";
import {accountIsAuthenticated, authenticateUser} from "@/api/gateway";
import {getAddress} from "viem";
import {useQuery} from "@tanstack/react-query";
import {generateEIP712} from "@/lib/eip712";
import {toast} from "sonner";




export function useTenChainAuth(address: `0x${string}`) {
    const [tenToken] = useLocalStorage<string|null>('ten_token', null)
    const [signature, setSignature] = useState<`0x${string}` | null>(null);
    const { signTypedDataAsync } = useSignTypedData()
    const [isLoading, setIsLoading] = useState<boolean>(false);

    const {
        data:beAuthCheck,
        isPending: beAuthCheckLoading,
        error: beAuthCheckError,
        isSuccess: beAuthCheckSuccess,
        refetch: beAuthCheckRefetch,
    } = useQuery({queryKey:[`beAuthCheck`, address, signature], queryFn: () => accountIsAuthenticated(tenToken ?? '', address)})


    useEffect(() => {
        toast.error("Authentication error", {
            description: "Error: "+beAuthCheckError?.message,
        })

        if (beAuthCheckSuccess) {
            toast("Address successfully authenticated", {
                description: "You're now ready to explore TEN."
            })
        }
    }, [beAuthCheckError, beAuthCheckSuccess])


    const getSignature = async () => {
        if (!tenToken) return null;

        const tokenValue = tenToken.startsWith('0x') ? tenToken : '0x' + tenToken;
        let signature = null

        try {
            const checksummedAddress = getAddress(tokenValue);

            signature  = await signTypedDataAsync(generateEIP712(checksummedAddress));

            setSignature(signature);
        } catch (err) {
            console.error('Invalid address format:', err);
        }

        return signature
    }

    const authenticateWalletWithBE = async (signature: string) => {
        if (!tenToken) return null;

        await authenticateUser(tenToken, {
            signature,
            address
        });

        await beAuthCheckRefetch()
    }

    const authenticateAccount = async () => {
        setIsLoading(true);
        const newSignature = await getSignature()
        if (newSignature) {
            await authenticateWalletWithBE(newSignature)
        }
        setIsLoading(false);
    }


    return {
        isAuthenticated: !!beAuthCheck?.status,
        isAuthenticatedLoading: isLoading,
        authenticateAccount
    }
}