'use client';

import { fallback, injected, unstable_connector, WagmiProvider } from 'wagmi';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { createConfig } from 'wagmi';
import { http } from 'viem';
import { type Chain } from 'wagmi/chains';
import {
    nativeCurrency,
    tenChainIDDecimal,
    tenGatewayAddress,
    tenNetworkName,
    tenscanAddress,
} from '@/lib/constants';
import { ThemeProvider } from '@/components/ThemeProvider/ThemeProvider';
import { TooltipProvider } from '@/components/ui/tooltip';
import { createStorage } from 'wagmi';
import { SidebarProvider } from '@/components/ui/sidebar';
import { AppSidebar } from '@/components/AppSidebar/AppSidebar';

const createCustomChain = (rpcUrl: string): Chain =>
    ({
        id: tenChainIDDecimal,
        name: tenNetworkName,
        nativeCurrency: nativeCurrency,
        rpcUrls: {
            default: {
                http: [rpcUrl],
            },
            public: {
                http: [rpcUrl],
            },
        },
        blockExplorers: {
            default: {
                name: 'TenScan',
                url: tenscanAddress,
            },
        },
    }) as const satisfies Chain;

const queryClient = new QueryClient();

export function Providers({ children }: { children: React.ReactNode }) {
    const chain = createCustomChain(`${tenGatewayAddress}/v1/`);

    const config = createConfig({
        chains: [chain],
        connectors: [injected()],
        transports: {
            [tenChainIDDecimal]: fallback([unstable_connector(injected), http()]),
        },
        ssr: true,
        storage: createStorage({
            storage: typeof window !== 'undefined' ? window.localStorage : undefined,
        }),
    });

    return (
        <WagmiProvider config={config}>
            <QueryClientProvider client={queryClient}>
                <TooltipProvider>
                    <ThemeProvider
                        attribute="class"
                        defaultTheme="dark"
                        enableSystem
                        disableTransitionOnChange
                    >
                        <SidebarProvider>
                            <AppSidebar />
                            {children}
                        </SidebarProvider>
                    </ThemeProvider>
                </TooltipProvider>
            </QueryClientProvider>
        </WagmiProvider>
    );
}
