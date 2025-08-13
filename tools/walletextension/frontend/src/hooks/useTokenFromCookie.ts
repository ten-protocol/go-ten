import { useState, useEffect, useCallback } from 'react';
import { getTokenFromCookie, setTokenToCookie } from '../api/gateway';

export function useTokenFromCookie(): [string, (token: string) => Promise<void>, boolean, () => Promise<void>] {
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
            
            // If no cookie found, just set empty token - don't auto-join
            // The join flow should only happen during wallet connection, not authentication
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

    const refreshToken = useCallback(async () => {
        await fetchToken();
    }, [fetchToken]);

    useEffect(() => {
        fetchToken();
    }, [fetchToken]);

    return [token, updateToken, isLoading, refreshToken];
}