import React, { createContext, useState, useContext, useEffect } from "react";
import { ethers } from "ethers";
import detectEthereumProvider from "@metamask/detect-provider";
import Web3Service from "@/src/services/web3service";
import { toast } from "../ui/use-toast";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  WalletNetwork,
} from "@/src/types";
import { L1CHAINS, L2CHAINS } from "@/src/lib/constants";

const WalletContext = createContext<WalletConnectionContextType | null>(null);

const WalletProvider = ({ children }: WalletConnectionProviderProps) => {
  const [provider, setProvider] = useState<any>(null);
  const [signer, setSigner] = useState<any>(null);
  const [address, setAddress] = useState<string>("");
  const [isL1ToL2, setIsL1ToL2] = React.useState(true);
  const fromChains = isL1ToL2 ? L1CHAINS : L2CHAINS;
  const toChains = isL1ToL2 ? L2CHAINS : L1CHAINS;

  // Function to connect wallet
  const connectWallet = async () => {
    try {
      const detectedProvider = await detectEthereumProvider();
      setProvider(detectedProvider);
      //@ts-ignore
      const chainId = await detectedProvider!.request({
        method: "eth_chainId",
      });

      if (isL1ToL2 && chainId !== WalletNetwork.L1_SEPOLIA) {
        switchNetwork();
      }

      if (!isL1ToL2 && chainId !== WalletNetwork.L2_TEN_TESTNET) {
        switchNetwork();
      }

      //@ts-ignore
      const accounts = await detectedProvider!.request({
        method: "eth_requestAccounts",
      });
      setAddress(accounts[0]);
      toast({
        title: "Connected",
        description: "Connected to wallet! Account: " + accounts[0],
        variant: "success",
      });
    } catch (error) {
      console.error("Error connecting to wallet:", error);
      toast({
        title: "Error",
        description: "Error connecting to wallet",
        variant: "destructive",
      });
    }
  };

  // Function to disconnect wallet
  const disconnectWallet = () => {
    try {
      if (provider) {
        provider.removeAllListeners();
        setProvider(null);
        setSigner(null);
        setAddress("");
        toast({
          title: "Disconnected",
          description: "Disconnected from wallet",
          variant: "info",
        });
      }
    } catch (error) {
      console.error("Error disconnecting from wallet:", error);
      toast({
        title: "Error",
        description: "Error disconnecting from wallet",
        variant: "destructive",
      });
    }
  };

  // Function to switch network
  const switchNetwork = async () => {
    try {
      if (isL1ToL2) {
        await provider.request({
          method: "wallet_switchEthereumChain",
          params: [{ chainId: WalletNetwork.L2_TEN_TESTNET }],
        });
      } else {
        await provider.request({
          method: "wallet_switchEthereumChain",
          params: [{ chainId: WalletNetwork.L1_SEPOLIA }],
        });
      }
      setIsL1ToL2(!isL1ToL2);
    } catch (error: any) {
      console.error("Error switching network:", error);
      if (error.code === 4902) {
        // if the network is not installed
        if (isL1ToL2) {
          toast({
            title: "Wrong Network",
            description: (
              <>
                <span> You are not connected to Ten! Connect at: </span>
                <pre className="mt-2 w-[500px] rounded-md bg-slate-950 p-4">
                  <a
                    href="https://testnet.ten.xyz/"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    https://testnet.ten.xyz/
                  </a>
                </pre>
              </>
            ),
            variant: "info",
          });
        } else {
          toast({
            title: "Wrong Network",
            description: "Please install the network to continue.",
            variant: "info",
          });
        }
      } else {
        // generic error message
        toast({
          title: "Error",
          description: "Error switching network",
          variant: "destructive",
        });
      }
    }
  };

  // Effect to set signer and handle events
  useEffect(() => {
    if (provider) {
      const newSigner = new ethers.providers.Web3Provider(provider).getSigner();
      new Web3Service(newSigner);
      setSigner(newSigner);

      const handleAccountsChange = (accounts: string[]) => {
        setAddress(accounts[0]);
      };
      const handleChainChange = () => {
        window.location.reload();
      };

      provider.on("accountsChanged", handleAccountsChange);
      provider.on("chainChanged", handleChainChange);

      return () => {
        provider.removeListener("accountsChanged", handleAccountsChange);
        provider.removeListener("chainChanged", handleChainChange);
      };
    }
  }, [provider]);

  useEffect(() => {
    //connect wallet
    connectWallet();
  }, []);

  // Context value
  const value = {
    provider,
    signer,
    address,
    walletConnected: !!provider,
    isL1ToL2,
    fromChains,
    toChains,
    connectWallet,
    disconnectWallet,
    switchNetwork,
  };

  // Render the provider with the context value
  return (
    <WalletContext.Provider value={value}>{children}</WalletContext.Provider>
  );
};

// Custom hook to use the wallet context
const useWalletStore = () => {
  const context = useContext(WalletContext);
  if (!context) {
    throw new Error("useWalletStore must be used within a WalletProvider");
  }
  return context;
};

// Export provider and custom hook
export { WalletProvider, useWalletStore };
