import React, { createContext, useContext, useEffect } from "react";
import useWalletStore from "@/src/stores/wallet-store";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
} from "@/src/types";

const WalletContext = createContext<WalletConnectionContextType | null>(null);

const WalletProvider = ({ children }: WalletConnectionProviderProps) => {
  const { restoreWalletState } = useWalletStore();

  useEffect(() => {
    restoreWalletState();

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const value = {};

  return (
    <WalletContext.Provider value={value}>{children}</WalletContext.Provider>
  );
};

const useWallet = () => {
  const context = useContext(WalletContext);
  if (!context) {
    throw new Error("useWallet must be used within a WalletProvider");
  }
  return context;
};

export { WalletProvider, useWallet };
