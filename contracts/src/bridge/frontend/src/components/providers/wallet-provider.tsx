import React, { createContext, useState, useContext, useEffect } from "react";
import { ethers } from "ethers";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
} from "@/src/types";

const WalletContext = createContext<WalletConnectionContextType | null>(null);

const WalletProvider = ({ children }: WalletConnectionProviderProps) => {
  const [provider, setProvider] = useState<any>(null);
  const [signer, setSigner] = useState<any>(null);
  const [address, setAddress] = useState<string | null>(null);

  useEffect(() => {
    if (provider) {
      const newSigner = new ethers.providers.Web3Provider(provider).getSigner();
      setSigner(newSigner);
    }
  }, [provider]);

  const handleSetProvider = (newProvider: any) => {
    setProvider(newProvider);
  };

  const handleSetAddress = (newAddress: string) => {
    setAddress(newAddress);
  };

  const value = {
    provider,
    signer,
    address,
    setProvider: handleSetProvider,
    setAddress: handleSetAddress,
  };

  return (
    <WalletContext.Provider value={value}>{children}</WalletContext.Provider>
  );
};

const useWalletStore = () => {
  const context = useContext(WalletContext);
  if (!context) {
    throw new Error("useWalletStore must be used within a WalletProvider");
  }
  return context;
};

export { WalletProvider, useWalletStore };
