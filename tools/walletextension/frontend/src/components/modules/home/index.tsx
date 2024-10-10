import React from "react";
import { useWalletConnection } from "../../providers/wallet-provider";
import Connected from "./connected";
import Disconnected from "./disconnected";
import { Skeleton } from "@repo/ui/components/shared/skeleton";

const Home = () => {
  const { walletConnected, loading } = useWalletConnection();

  return (
    <div className="w-[800px] mx-auto">
      {loading ? (
        <Skeleton className="h-[400px]" />
      ) : walletConnected ? (
        <Connected />
      ) : (
        <Disconnected />
      )}
    </div>
  );
};

export default Home;
