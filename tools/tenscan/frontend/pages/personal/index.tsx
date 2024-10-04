import React from "react";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import PersonalTransactions from "@/src/components/modules/personal";
import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import ConnectWalletButton from "@/src/components/modules/common/connect-wallet";
import EmptyState from "@/src/components/modules/common/empty-state";
import { ethereum } from "@/src/lib/utils";
import HeadSeo from "@/src/components/head-seo";
import { siteMetadata } from "@/src/lib/siteMetadata";

export const metadata: Metadata = {
  title: "Personal Transactions",
  description: "Tenscan Personal Transactions",
};

export default function PersonalPage() {
  const { walletConnected } = useWalletConnection();

  return (
    <>
      <HeadSeo
        title={`${siteMetadata.personal.title} `}
        description={siteMetadata.personal.description}
        canonicalUrl={`${siteMetadata.personal.canonicalUrl}`}
        ogImageUrl={siteMetadata.personal.ogImageUrl}
        ogTwitterImage={siteMetadata.personal.ogTwitterImage}
        ogType={siteMetadata.personal.ogType}
      ></HeadSeo>
      <Layout>
        {walletConnected ? (
          <PersonalTransactions />
        ) : (
          <EmptyState
            title="Connect Wallet"
            description="Connect your wallet to view your personal transactions."
            action={
              <ConnectWalletButton
                text={
                  ethereum
                    ? "Connect Wallet to continue"
                    : "Install MetaMask to continue"
                }
              />
            }
          />
        )}
      </Layout>
    </>
  );
}
