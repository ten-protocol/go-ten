import React, { createContext, useState, useContext, useEffect } from "react";
import { toast } from "../ui/use-toast";
import {
  ToastType,
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  WalletNetwork,
} from "@/src/types";
import { requestMethods } from "@/src/routes";
import {
  getEthereumProvider,
  handleStorage,
} from "@/src/lib/utils/walletUtils";
import {
  setupEventListeners,
  initializeSigner,
} from "@/src/lib/utils/walletEvents";

const WalletContext = createContext<WalletConnectionContextType | null>(null);

const WalletProvider = ({ children }: WalletConnectionProviderProps) => {
  const [provider, setProvider] = useState<any>(null);
  const [isWalletConnected, setIsWalletConnected] = useState<boolean>(false);
  const [signer, setSigner] = useState<any>(null);
  const [address, setAddress] = useState<string>("");
  const [isL1ToL2, setIsL1ToL2] = useState<boolean>(true);

  useEffect(() => {
    const storedAddress = handleStorage.get("walletAddress");
    const storedIsL1ToL2 = handleStorage.get("isL1ToL2") === "true";

    if (storedAddress) {
      setAddress(storedAddress);
      setIsWalletConnected(true);
      setIsL1ToL2(storedIsL1ToL2);
    }

    const initializeProvider = async () => {
      const detectedProvider = await getEthereumProvider();
      setProvider(detectedProvider);

      const chainId = await detectedProvider.request({
        method: requestMethods.getChainId,
      });

      const isL1 = chainId === WalletNetwork.L1_SEPOLIA;
      setIsL1ToL2(isL1);
    };

    initializeProvider();
  }, []);

  const connectWallet = async () => {
    try {
      const detectedProvider = await getEthereumProvider();
      setProvider(detectedProvider);

      const chainId = await detectedProvider.request({
        method: requestMethods.getChainId,
      });

      const isL1 = chainId === WalletNetwork.L1_SEPOLIA;
      setIsL1ToL2(isL1);

      const accounts = await detectedProvider.request({
        method: requestMethods.connectAccounts,
      });
      setIsWalletConnected(true);
      setAddress(accounts[0]);
      handleStorage.save("walletAddress", accounts[0]);
      handleStorage.save("isL1ToL2", isL1.toString());
      toast({
        title: "Connected",
        description: "Connected to wallet! Account: " + accounts[0],
        variant: ToastType.SUCCESS,
      });
    } catch (error) {
      console.error("Error connecting to wallet:", error);
      toast({
        title: "Error",
        description: "Error connecting to wallet",
        variant: ToastType.DESTRUCTIVE,
      });
    }
  };

  const disconnectWallet = () => {
    try {
      if (provider) {
        provider.removeAllListeners();
        setProvider(null);
        setSigner(null);
        setAddress("");
        setIsWalletConnected(false);
        handleStorage.remove("walletAddress");
        handleStorage.remove("isL1ToL2");
        toast({
          title: "Disconnected",
          description: "Disconnected from wallet",
          variant: ToastType.INFO,
        });
      }
    } catch (error) {
      console.error("Error disconnecting from wallet:", error);
      toast({
        title: "Error",
        description: "Error disconnecting from wallet",
        variant: ToastType.DESTRUCTIVE,
      });
    }
  };

  useEffect(() => {
    if (provider) {
      const newSigner = initializeSigner(provider);
      setSigner(newSigner);

      const cleanup = setupEventListeners(provider, setAddress);
      return cleanup;
    }
  }, [provider]);

  const switchNetwork = async () => {
    if (!provider) {
      toast({
        title: "Error",
        description: "Please connect to wallet first",
        variant: ToastType.DESTRUCTIVE,
      });
      return;
    }
    try {
      const desiredNetwork = isL1ToL2
        ? WalletNetwork.L2_TEN_TESTNET
        : WalletNetwork.L1_SEPOLIA;
      await provider.request({
        method: requestMethods.switchNetwork,
        params: [{ chainId: desiredNetwork }],
      });
      const isL1 = desiredNetwork === WalletNetwork.L1_SEPOLIA;
      setIsL1ToL2(isL1);
      handleStorage.save("isL1ToL2", isL1.toString());
      toast({
        title: "Network Switched",
        description: `Switched to ${
          desiredNetwork === WalletNetwork.L2_TEN_TESTNET
            ? "L2 TEN Testnet"
            : "L1 Sepolia"
        }`,
        variant: ToastType.SUCCESS,
      });
    } catch (error: any) {
      console.error("Error switching network:", error);
      if (error.code === 4902) {
        toast({
          title: "Network Not Found",
          description: "Network not found in wallet",
          variant: ToastType.INFO,
        });
      } else {
        toast({
          title: "Error",
          description: error.message || "Error switching network",
          variant: ToastType.DESTRUCTIVE,
        });
      }
    }
  };

  const value = {
    provider,
    signer,
    address,
    walletConnected: isWalletConnected,
    isL1ToL2,
    connectWallet,
    disconnectWallet,
    switchNetwork,
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
