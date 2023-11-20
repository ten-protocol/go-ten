import { ThemeProvider } from "@/components/providers/theme-provider";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Toaster } from "@/components/ui/toaster";
import { WalletConnectionProvider } from "@/components/providers/wallet-provider";
import { NetworkStatus } from "@/components/modules/common/network-status";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <WalletConnectionProvider>
        <Component {...pageProps} />
        <Toaster />
        <NetworkStatus />
        <ReactQueryDevtools initialIsOpen={false} />
      </WalletConnectionProvider>
    </ThemeProvider>
  );
}
