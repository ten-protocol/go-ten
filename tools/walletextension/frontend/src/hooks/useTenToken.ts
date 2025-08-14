import { useState, useEffect, useCallback } from 'react';
import { getTokenFromCookie, setTokenInCookie } from '@/api/gateway';

interface UseTenTokenReturn {
    token: string;
    setToken: (newToken: string) => Promise<boolean>;
    loading: boolean;
    error: string | null;
}

export function useTenToken(): UseTenTokenReturn {
    const [token, setTokenState] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    // Load token from cookie on mount
    useEffect(() => {
        const loadToken = async () => {
            console.log('[useTenToken] Loading token on mount');
            setLoading(true);
            setError(null);
            
            try {
                const retrievedToken = await getTokenFromCookie();
                setTokenState(retrievedToken);
                console.log('[useTenToken] Token loaded successfully:', retrievedToken ? 'present' : 'empty');
            } catch (err) {
                const errorMsg = err instanceof Error ? err.message : 'Unknown error loading token';
                console.log('[useTenToken] Error loading token:', errorMsg);
                setError(errorMsg);
                setTokenState('');
            } finally {
                setLoading(false);
            }
        };

        loadToken();
    }, []);

    // Set token function
    const setToken = useCallback(async (newToken: string): Promise<boolean> => {
        console.log('[useTenToken] Setting new token');
        setError(null);
        
        try {
            const success = await setTokenInCookie(newToken);
            if (success) {
                setTokenState(newToken);
                console.log('[useTenToken] Token set successfully');
                return true;
            } else {
                const errorMsg = 'Failed to set token in cookie';
                console.log('[useTenToken] ' + errorMsg);
                setError(errorMsg);
                return false;
            }
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : 'Unknown error setting token';
            console.log('[useTenToken] Error setting token:', errorMsg);
            setError(errorMsg);
            return false;
        }
    }, []);

    return {
        token,
        setToken,
        loading,
        error,
    };
}