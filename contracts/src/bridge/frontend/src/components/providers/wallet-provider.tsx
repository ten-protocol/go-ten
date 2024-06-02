import React, { createContext, useState, useContext, useEffect } from "react";
import { ethers } from "ethers";
import detectEthereumProvider from "@metamask/detect-provider";
import { toast } from "../ui/use-toast";
import {
  ToastType,
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  WalletNetwork,
} from "@/src/types";
import { L1CHAINS, L2CHAINS } from "@/src/lib/constants";
import { requestMethods } from "@/src/routes";

const WalletContext = createContext<WalletConnectionContextType | null>(null);

const WalletProvider = ({ children }: WalletConnectionProviderProps) => {
  const [provider, setProvider] = useState<any>(null);
  const [isWalletConnected, setIsWalletConnected] = useState<boolean>(false);
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
        method: requestMethods.getChainId,
      });

      if (isL1ToL2 && chainId !== WalletNetwork.L1_SEPOLIA) {
        switchNetwork();
      }

      if (!isL1ToL2 && chainId !== WalletNetwork.L2_TEN_TESTNET) {
        switchNetwork();
      }

      //@ts-ignore
      const accounts = await detectedProvider!.request({
        method: requestMethods.connectAccounts,
      });
      setIsWalletConnected(true);
      setAddress(accounts[0]);
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

  // Function to disconnect wallet
  const disconnectWallet = () => {
    try {
      if (provider) {
        provider.removeAllListeners();
        setProvider(null);
        setSigner(null);
        setAddress("");
        setIsWalletConnected(false);
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

  // Function to switch network
  const switchNetwork = async () => {
    try {
      if (isL1ToL2) {
        await provider.request({
          method: requestMethods.switchNetwork,
          params: [{ chainId: WalletNetwork.L2_TEN_TESTNET }],
        });
      } else {
        await provider.request({
          method: requestMethods.switchNetwork,
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
            variant: ToastType.INFO,
          });
        } else {
          toast({
            title: "Wrong Network",
            description: "Please install the network to continue.",
            variant: ToastType.INFO,
          });
        }
      } else {
        // generic error message
        toast({
          title: "Error",
          description: "Error switching network",
          variant: ToastType.DESTRUCTIVE,
        });
      }
    }
  };

  // Effect to set signer and handle events
  useEffect(() => {
    if (provider) {
      const newSigner = new ethers.providers.Web3Provider(provider).getSigner();
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

  // useEffect(() => {
  //   connectWallet();
  // }, []);

  // Context value
  const value = {
    provider,
    signer,
    address,
    walletConnected: isWalletConnected,
    isL1ToL2,
    fromChains,
    toChains,
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
