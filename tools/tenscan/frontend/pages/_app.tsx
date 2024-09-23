import { useState } from "react";
import { ThemeProvider } from "@/src/components/providers/theme-provider";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import {
  QueryClient,
  MutationCache,
  QueryClientProvider,
} from "@tanstack/react-query";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Toaster } from "@repo/ui/shared/toaster";
import { WalletConnectionProvider } from "@/src/components/providers/wallet-provider";
import { NetworkStatus } from "@repo/ui/common/network-status";
import HeadSeo from "@/src/components/head-seo";
import { siteMetadata } from "@/src/lib/siteMetadata";
import Script from "next/script";
import { GOOGLE_ANALYTICS_ID } from "@/src/lib/constants";
import { showToast } from "@repo/ui/shared/use-toast";
import { ToastType } from "@/src/types/interfaces";

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
          },
        },
        mutationCache,
      })
  );

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
      ></HeadSeo>
      <QueryClientProvider client={queryClient}>
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
      </QueryClientProvider>
    </>
  );
}
