import React from "react";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import PersonalTransactions from "@/src/components/modules/personal";
import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import ConnectWalletButton from "@/src/components/modules/common/connect-wallet";
import EmptyState from "@/src/components/modules/common/empty-state";

export const metadata: Metadata = {
  title: "Personal Transactions",
  description: "Tenscan Personal Transactions",
};

export default function PersonalPage() {
  const { walletConnected } = useWalletConnection();

  return (
    <Layout>
      {walletConnected ? (
        <PersonalTransactions />
      ) : (
        <EmptyState
          title="Connect your wallet"
          description="Connect your wallet to view your personal transactions."
          action={<ConnectWalletButton />}
        />
      )}
    </Layout>
  );
}
