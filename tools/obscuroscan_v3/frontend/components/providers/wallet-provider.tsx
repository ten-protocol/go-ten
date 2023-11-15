import { createContext, useContext, useEffect, useState } from "react";
import { ethers } from "ethers";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
} from "@/src/types/interfaces/WalletInterfaces";
import { useToast } from "../ui/use-toast";

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
  const { toast } = useToast();

  const [walletConnected, setWalletConnected] = useState(false);
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [provider, setProvider] =
    useState<ethers.providers.Web3Provider | null>(null);

  const connectWallet = async () => {
    if ((window as any).ethereum) {
      const ethProvider = new ethers.providers.Web3Provider(
        (window as any).ethereum
      );
      setProvider(ethProvider);

      try {
        await ethProvider.send("eth_requestAccounts", []);
        const signer = ethProvider.getSigner();
        const address = await signer.getAddress();
        setWalletAddress(address);
        setWalletConnected(true);
      } catch (error: any) {
        toast({ description: "Error connecting to wallet:" + error.message });
      }
    } else {
      toast({ description: "No ethereum object found." });
    }
  };

  const disconnectWallet = () => {
    if (provider) {
      provider.removeAllListeners();
      setWalletConnected(false);
      setWalletAddress(null);
      setProvider(null);
    }
  };

  useEffect(() => {
    const ethereum = (window as any).ethereum;
    const handleAccountsChanged = (accounts: string[]) => {
      if (accounts.length === 0) {
        toast({ description: "Please connect to MetaMask." });
      } else if (accounts[0] !== walletAddress) {
        setWalletAddress(accounts[0]);
      }
    };
    ethereum.on("accountsChanged", handleAccountsChanged);
    return () => {
      ethereum.removeListener("accountsChanged", handleAccountsChanged);
    };
  });

  const walletConnectionContextValue: WalletConnectionContextType = {
    provider,
    walletConnected,
    walletAddress,
    connectWallet,
    disconnectWallet,
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
