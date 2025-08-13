import { useState, useEffect, useCallback } from 'react';
import { getTokenFromCookie, setTokenToCookie, joinTestnet } from '../api/gateway';

export function useTokenFromCookie(): [string, (token: string) => Promise<void>, boolean] {
    const [token, setToken] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);

    // TODO: Add caching mechanism to avoid calling get-token API every time
    const fetchToken = useCallback(async () => {
        try {
            setIsLoading(true);
            const tokenFromCookie = await getTokenFromCookie();
            setToken(tokenFromCookie || '');
        } catch (error: any) {
            console.error('Failed to fetch token from cookie:', error);
            
            // Check if error message indicates cookie not found
            if (error?.message?.includes('gateway_token cookie not found')) {
                console.log('No cookie found, calling /join to create new token');
                try {
                    const newToken = await joinTestnet();
                    if (newToken) {
                        await setTokenToCookie(newToken);
                        setToken(newToken);
                        return;
                    }
                } catch (joinError) {
                    console.error('Failed to join testnet:', joinError);
                }
            }
            
            setToken(''); // Empty token means not authenticated
        } finally {
            setIsLoading(false);
        }
    }, []);

    const updateToken = useCallback(async (newToken: string) => {
        try {
            await setTokenToCookie(newToken);
            setToken(newToken);
        } catch (error) {
            console.error('Failed to set token to cookie:', error);
            throw error; // Re-throw so components can handle the error
        }
    }, []);

    useEffect(() => {
        fetchToken();
    }, [fetchToken]);

    return [token, updateToken, isLoading];
}