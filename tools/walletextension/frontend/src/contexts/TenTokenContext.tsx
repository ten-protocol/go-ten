import React, { createContext, useContext, useEffect, useState } from 'react';
import { getTokenFromCookie } from '@/api/gateway';

interface TenTokenContextType {
    token: string;
    loading: boolean;
    error: string | null;
}

const TenTokenContext = createContext<TenTokenContextType | undefined>(undefined);

interface TenTokenProviderProps {
    children: React.ReactNode;
}

export function TenTokenProvider({ children }: TenTokenProviderProps) {
    const [token, setTokenState] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    // Load token from cookie on mount
    useEffect(() => {
        const loadToken = async () => {
            console.log('[TenTokenProvider] Loading token on mount');
            setLoading(true);
            setError(null);
            
            try {
                const retrievedToken = await getTokenFromCookie();
                setTokenState(retrievedToken);
                console.log('[TenTokenProvider] Token loaded successfully:', retrievedToken ? 'present' : 'empty');
            } catch (err) {
                const errorMsg = err instanceof Error ? err.message : 'Unknown error loading token';
                console.log('[TenTokenProvider] Error loading token:', errorMsg);
                setError(errorMsg);
                setTokenState('');
            } finally {
                setLoading(false);
            }
        };

        loadToken();
    }, []);


    const value = {
        token,
        loading,
        error,
    };

    console.log('[TenTokenProvider] Context value updated - token:', token ? 'present' : 'empty', 'loading:', loading);

    return (
        <TenTokenContext.Provider value={value}>
            {children}
        </TenTokenContext.Provider>
    );
}

export function useTenToken(): TenTokenContextType {
    const context = useContext(TenTokenContext);
    if (context === undefined) {
        throw new Error('useTenToken must be used within a TenTokenProvider');
    }
    return context;
}