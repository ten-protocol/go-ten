import { createContext, useContext, useEffect } from "react";
import { ethers } from "ethers";
import { ethereum } from "@repo/ui/lib/utils";
import { showToast } from "@repo/ui/components/shared/use-toast";
import { ToastType } from "@repo/ui/lib/enums/toast";
import useWalletStore from "@/stores/wallet-store";

const WalletConnectionContext = createContext<void | ReturnType<
  typeof useWalletStore
>>(undefined);

export const useWalletConnection = () => {
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
}: {
  children: React.ReactNode;
}) => {
  const walletStore = useWalletStore();
  const { initializeGateway, fetchUserAccounts, setLoading } = walletStore;

  useEffect(() => {
    const initializeWallet = async () => {
      setLoading(true);
      if (ethereum && ethereum.isMetaMask) {
        try {
          const providerInstance = new ethers.providers.Web3Provider(ethereum);
          useWalletStore.setState({ provider: providerInstance });
          await initializeGateway();

          ethereum.on("accountsChanged", fetchUserAccounts);
        } catch (error) {
          console.error("Failed to initialize wallet:", error);
          showToast(
            ToastType.DESTRUCTIVE,
            "Failed to initialize wallet. Please refresh and try again."
          );
        }
      } else {
        showToast(
          ToastType.WARNING,
          "MetaMask not detected. Some features may be unavailable."
        );
      }
      setLoading(false);
    };

    initializeWallet();

    return () => {
      if (ethereum && ethereum.removeListener) {
        ethereum.removeListener("accountsChanged", fetchUserAccounts);
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ethereum]);

  return (
    <WalletConnectionContext.Provider value={walletStore}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
