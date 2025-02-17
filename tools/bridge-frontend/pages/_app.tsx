import { useState } from "react";
import { ThemeProvider } from "@/src/components/providers/theme-provider";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { QueryClient, MutationCache } from "@tanstack/react-query";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Toaster } from "@/src/components/ui/toaster";
import { NetworkStatus } from "@/src/components/modules/common/network-status";
import HeadSeo from "@/src/components/head-seo";
import { siteMetadata } from "@/src/lib/siteMetadata";
import Script from "next/script";
import { GOOGLE_ANALYTICS_ID } from "@/src/lib/constants";
import { showToast } from "@/src/components/ui/use-toast";
import { ToastType } from "@/src/types";
import { WalletProvider } from "@/src/components/providers/wallet-provider";
import { PersistQueryClientProvider } from "@tanstack/react-query-persist-client";
import { createSyncStoragePersister } from "@tanstack/query-sync-storage-persister";

export default function App({ Component, pageProps }: AppProps) {
  const mutationCache = new MutationCache({
    onSuccess: (mutation: any) => {
      if (mutation?.message) {
        showToast(ToastType.SUCCESS, mutation?.message);
      }
    },
    onError: (error: any, mutation: any) => {
      if (error?.response?.data?.message) {
        showToast(ToastType.DESTRUCTIVE, error?.response?.data?.message);
      }
    },
  });

  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            refetchOnWindowFocus: false,
            staleTime: 300000,
            gcTime: 1000 * 60 * 60 * 24, // 24 hours
          },
        },
        mutationCache,
      })
  );

  const localStoragePersister = createSyncStoragePersister({
    storage: typeof window !== "undefined" ? window.localStorage : null,
  });

  return (
    <>
      <Script
        strategy="lazyOnload"
        src={`https://www.googletagmanager.com/gtag/js?id='${GOOGLE_ANALYTICS_ID}'`}
      />

      <Script strategy="lazyOnload" id="google-analytics">
        {`
        window.dataLayer = window.dataLayer || [];
        function gtag(){dataLayer.push(arguments);}
        gtag('js', new Date());
        gtag('config', '${GOOGLE_ANALYTICS_ID}');
        `}
      </Script>

      <HeadSeo
        title={`${siteMetadata.companyName} `}
        description={siteMetadata.description}
        canonicalUrl={`${siteMetadata.siteUrl}`}
        ogImageUrl={siteMetadata.siteLogo}
        ogTwitterImage={siteMetadata.siteLogo}
        ogType={"website"}
      >
        <link rel="icon" href="/favicon/favicon.ico" />
        <link
          rel="apple-touch-icon"
          sizes="180x180"
          href="/favicon/apple-touch-icon.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="32x32"
          href="/favicon/favicon-32x32.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="16x16"
          href="/favicon/favicon-16x16.png"
        />
        <link rel="manifest" href="/favicon/site.webmanifest" />
      </HeadSeo>
      <PersistQueryClientProvider
        client={queryClient}
        persistOptions={{
          persister: localStoragePersister,
        }}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <WalletProvider>
            <Component {...pageProps} />
            <Toaster />
            <NetworkStatus />
            <ReactQueryDevtools initialIsOpen={false} />
          </WalletProvider>
        </ThemeProvider>
      </PersistQueryClientProvider>
    </>
  );
}
