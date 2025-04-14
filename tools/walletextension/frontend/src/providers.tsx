'use client';

import {
  RainbowKitProvider,
  darkTheme,
  lightTheme,
  connectorsForWallets
} from '@rainbow-me/rainbowkit';
import { 
  WagmiProvider,
} from 'wagmi';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import '@rainbow-me/rainbowkit/styles.css';
import { createConfig } from 'wagmi';
import { http } from 'viem';
import { useEffect, useState } from 'react';
import { type Chain } from 'wagmi/chains';
import {braveWallet, metaMaskWallet, rabbyWallet} from "@rainbow-me/rainbowkit/wallets";
import {joinTestnet} from "@/api/gateway";
import {tenGatewayAddress} from "@/lib/constants";
import {useLocalStorage} from "@/hooks/useLocalStorage";
import {ThemeProvider} from "@/components/ThemeProvider/ThemeProvider";
import {TooltipProvider} from "@/components/ui/tooltip";

const createCustomChain = (rpcUrl: string): Chain => ({
  id: 443,
  name: 'Ten Testnet',
  "nativeCurrency": {
    "name":"Sepolia Ether",
    "symbol":"ETH",
    "decimals":18
  },
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
      url: 'https://tenscan.io',
    },
  },
} as const satisfies Chain);

export const connectors = connectorsForWallets(
  [
    {
      groupName: 'Recommended',
      wallets: [metaMaskWallet, braveWallet, rabbyWallet],
    },
  ],
  {
    appName: 'TEN Testnet Gateway',
    projectId: '443',
  }
);

const queryClient = new QueryClient();

export function Providers({ children }: { children: React.ReactNode }) {
  const [mounted, setMounted] = useState(false);
  const [loading, setLoading] = useState(true);
  const [config, setConfig] = useState<ReturnType<typeof createConfig> | null>(null);
  const [customChain, setCustomChain] = useState<Chain | null>(null);
  const [tenToken, setTenToken] = useLocalStorage<string|null>('ten_token', null)

  useEffect(() => {
    const initializeChain = async () => {
      try {
        setLoading(true);

        try {
          if (!tenToken) {
            const newTenToken = await joinTestnet();
            setTenToken(newTenToken);
          }
        } catch (error) {
          console.error('Failed to fetch custom RPC, using default:', error);
        }

        if (tenToken) {
          const chain = createCustomChain(`${tenGatewayAddress}/v1/?token=${tenToken}`);
          setCustomChain(chain);

          const newConfig = createConfig({
            chains: [chain],
            transports: {
              [chain.id]: http(chain.rpcUrls.default.http[0]),
            },
            connectors,
          });

          setConfig(newConfig);
          setLoading(false);
        }

      } catch (error) {
        console.error('Failed to initialize chain:', error);

        setLoading(false);
      }
    };
    
    initializeChain();
    
    setMounted(true);
  }, [tenToken]);
  
  // Show loading or nothing while initializing
  if (loading || !mounted || !config || !customChain) {
    return null;
  }
  
  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider 
          theme={{
            lightMode: lightTheme(),
            darkMode: darkTheme({
              accentColor: 'white',
              accentColorForeground: 'black',
              borderRadius: 'small',
              fontStack: 'system',
              overlayBlur: 'small',
            })
          }}
          initialChain={customChain.id}
        >
          <TooltipProvider>
          <ThemeProvider
              attribute="class"
              defaultTheme="dark"
              enableSystem
              disableTransitionOnChange
          >
            {children}
          </ThemeProvider>
          </TooltipProvider>
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
} 