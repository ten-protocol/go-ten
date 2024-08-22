import { ethereum } from "@/lib/utils";
import { useWalletStore } from "@/stores/wallet-store";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
} from "@/types/interfaces/WalletInterfaces";
import { ethers } from "ethers";
import { createContext, useContext, useEffect } from "react";

const WalletConnectionContext =
  createContext<WalletConnectionContextType | null>(null);

export const useWalletConnection = (): WalletConnectionContextType => {
  const context = useContext(WalletConnectionContext);

  if (!context) {
    throw new Error(
      "useWalletConnection must be used within a WalletConnectionProvider"
    );
  }
  return context;
};
export const WalletConnectionProvider = ({
  children,
}: WalletConnectionProviderProps) => {
  const {
    initialize,
    walletConnected,
    loading,
    accounts,
    provider,
    fetchUserAccounts,
  } = useWalletStore();

  useEffect(() => {
    if (ethereum && ethereum.isMetaMask) {
      const providerInstance = new ethers.providers.Web3Provider(ethereum);
      initialize(providerInstance);
      ethereum.on("accountsChanged", fetchUserAccounts);
      ethereum.on("chainChanged", window.location.reload);
    }

    return () => {
      if (ethereum && ethereum.removeListener) {
        ethereum.removeListener("accountsChanged", fetchUserAccounts);
        ethereum.removeListener("chainChanged", window.location.reload);
      }
    };
  }, [initialize, fetchUserAccounts]);

  const walletConnectionContextValue: WalletConnectionContextType = {
    walletConnected,
    loading,
    accounts,
    provider,
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
