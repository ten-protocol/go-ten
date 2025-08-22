import React, { createContext, useContext, useEffect, useState } from 'react';
import { getTokenFromCookie } from '@/api/gateway';

interface TenTokenContextType {
    token: string;
    loading: boolean;
    error: string | null;
    refreshToken: () => Promise<void>;
}

const TenTokenContext = createContext<TenTokenContextType | undefined>(undefined);

interface TenTokenProviderProps {
    children: React.ReactNode;
}

export function TenTokenProvider({ children }: TenTokenProviderProps) {
    const [token, setTokenState] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    const loadToken = async () => {
        setLoading(true);
        setError(null);
        
        try {
            const retrievedToken = await getTokenFromCookie();
            setTokenState(retrievedToken);
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : 'Unknown error loading token';
            setError(errorMsg);
            setTokenState('');
        } finally {
            setLoading(false);
        }
    };

    const refreshToken = async () => {
        await loadToken();
    };

    // Load token from cookie on mount
    useEffect(() => {
        loadToken();
    }, []);

    const value = {
        token,
        loading,
        error,
        refreshToken,
    };


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