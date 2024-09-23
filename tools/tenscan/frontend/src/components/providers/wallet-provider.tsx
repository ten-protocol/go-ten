import { createContext, useContext, useEffect, useState } from "react";
import { ethers } from "ethers";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
} from "@/src/types/interfaces/WalletInterfaces";
import { showToast } from "@repo/ui/shared/use-toast";
import { ToastType } from "@repo/ui/lib/types/interfaces";
import { ethereum, currentNetwork } from "@repo/ui/lib/utils";
import { handleError } from "@repo/ui/lib/walletUtils";
import { ethMethods } from "@/src/routes";

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
  const [walletConnected, setWalletConnected] = useState(false);
  const [isWrongNetwork, setIsWrongNetwork] = useState(false);
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [provider, setProvider] =
    useState<ethers.providers.Web3Provider | null>(null);

  const connectWallet = async () => {
    if (ethereum) {
      const ethProvider = new ethers.providers.Web3Provider(ethereum);
      setProvider(ethProvider);

      try {
        await ethProvider.send("eth_requestAccounts", []);
        const signer = ethProvider.getSigner();
        const address = await signer.getAddress();
        const network = await ethProvider?.getNetwork();
        const chainId = network?.chainId;
        const expectedChainId = currentNetwork.l2;
        setIsWrongNetwork(chainId !== expectedChainId);
        setWalletAddress(address);
        setWalletConnected(true);
      } catch (error: any) {
        showToast(
          ToastType.DESTRUCTIVE,
          "Error connecting to wallet:" + error?.message
        );
      }
    } else {
      showToast(
        ToastType.DESTRUCTIVE,
        "No ethereum object found. Please install MetaMask!"
      );
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
    if (!ethereum) {
      return;
    }

    const handleAccountsChanged = (accounts: string[]) => {
      if (accounts.length === 0) {
        showToast(ToastType.DESTRUCTIVE, "Please connect to MetaMask.");
      } else if (accounts[0] !== walletAddress) {
        setWalletAddress(accounts[0]);
      }
    };

    ethereum.on("accountsChanged", handleAccountsChanged);
    return () => {
      if (!ethereum) return;
      ethereum.removeListener("accountsChanged", handleAccountsChanged);
    };
  });

  const switchNetwork = async () => {
    if (!provider) {
      showToast(ToastType.DESTRUCTIVE, "Please connect to wallet first");
      return;
    }

    const desiredNetwork = ethers.utils.hexValue(currentNetwork.l2);

    try {
      await provider?.send(ethMethods.switchNetwork, [
        { chainId: desiredNetwork },
      ]);

      showToast("Switched to TEN Testnet", ToastType.SUCCESS);
    } catch (error: any) {
      if (error.code === 4902) {
        showToast(
          error.message || "Network not found in wallet",
          ToastType.INFO
        );
        showToast("Redirecting to TEN Gateway...", ToastType.INFO);
        return window.open(currentNetwork.l2Gateway, "_blank");
      } else {
        showToast(
          ToastType.DESTRUCTIVE,
          error.message || "Error switching network"
        );
      }
      handleError(error, "Error switching network");
    }
  };

  const walletConnectionContextValue: WalletConnectionContextType = {
    provider,
    walletConnected,
    walletAddress,
    connectWallet,
    disconnectWallet,
    switchNetwork,
    isWrongNetwork,
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
