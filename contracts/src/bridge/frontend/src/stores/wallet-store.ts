import { create } from "zustand";
import { toast } from "@/src/components/ui/use-toast";
import { ToastType, WalletNetwork } from "@/src/types";
import { requestMethods } from "@/src/routes";
import {
  getEthereumProvider,
  handleStorage,
} from "@/src/lib/utils/walletUtils";
import {
  initializeSigner,
  setupEventListeners,
} from "@/src/lib/utils/walletEvents";

interface WalletState {
  provider: any;
  signer: any;
  address: string;
  walletConnected: boolean;
  isL1ToL2: boolean;
  initializeProvider: () => void;
  connectWallet: () => void;
  disconnectWallet: () => void;
  switchNetwork: () => void;
  restoreWalletState: () => void;
}

const useWalletStore = create<WalletState>((set, get) => ({
  provider: null,
  signer: null,
  address: "",
  walletConnected: false,
  isL1ToL2: true,

  initializeProvider: async () => {
    const detectedProvider = await getEthereumProvider();
    console.log(
      "🚀 ~ initializeProvider: ~ detectedProvider:",
      detectedProvider
    );

    const newSigner = initializeSigner(detectedProvider);

    const chainId = await detectedProvider?.request({
      method: requestMethods.getChainId,
      params: [],
    });
    const isL1 = chainId === WalletNetwork.L1_SEPOLIA;

    set({
      provider: detectedProvider,
      signer: newSigner,
      isL1ToL2: isL1,
    });

    const cleanup = setupEventListeners(detectedProvider, (address: string) => {
      set({ address });
    });

    return cleanup;
  },

  connectWallet: async () => {
    try {
      const detectedProvider = await getEthereumProvider();
      const accounts = await detectedProvider?.request({
        method: requestMethods.connectAccounts,
        params: [],
      });
      const newSigner = initializeSigner(detectedProvider);

      set({
        provider: detectedProvider,
        signer: newSigner,
        address: accounts[0],
        walletConnected: true,
      });

      handleStorage.save("walletAddress", accounts[0]);
      handleStorage.save("isL1ToL2", get().isL1ToL2.toString());

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
  },

  disconnectWallet: () => {
    try {
      const { provider } = get();
      if (provider) {
        provider.removeAllListeners();
        set({
          provider: null,
          signer: null,
          address: "",
          walletConnected: false,
        });

        handleStorage.remove("walletAddress");
        handleStorage.remove("isL1ToL2");
        handleStorage.remove("tenBridgeReceiver");
        window.location.reload();

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
  },

  switchNetwork: async () => {
    const { provider, isL1ToL2 } = get();
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

      set({ isL1ToL2: isL1 });
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
  },

  restoreWalletState: () => {
    const storedAddress = handleStorage.get("walletAddress");
    const storedIsL1ToL2 = handleStorage.get("isL1ToL2") === "true";

    if (storedAddress) {
      set({
        address: storedAddress,
        walletConnected: true,
        isL1ToL2: storedIsL1ToL2,
      });
    }

    get().initializeProvider();
  },
}));

export default useWalletStore;
