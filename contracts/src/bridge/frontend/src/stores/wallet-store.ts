import { create } from "zustand";
import { toast } from "@/src/components/ui/use-toast";
import { IWalletState, ToastType } from "@/src/types";
import { requestMethods } from "@/src/routes";
import {
  getEthereumProvider,
  handleStorage,
} from "@/src/lib/utils/walletUtils";
import {
  initializeSigner,
  setupEventListeners,
} from "@/src/lib/utils/walletEvents";
import { currentNetwork } from "../lib/utils";

const useWalletStore = create<IWalletState>((set, get) => ({
  provider: null,
  signer: null,
  address: "",
  walletConnected: false,
  isL1ToL2: true,
  isWrongNetwork: false,

  initializeProvider: async () => {
    const detectedProvider = await getEthereumProvider();
    const newSigner = initializeSigner(detectedProvider);

    //@ts-ignore
    const chainId = await detectedProvider?.request({
      method: requestMethods.getChainId,
      params: [],
    });

    const isL1 = chainId === currentNetwork.l1;
    const expectedChainId = isL1 ? currentNetwork.l1 : currentNetwork.l2;

    set({
      provider: detectedProvider,
      signer: newSigner,
      isL1ToL2: !isL1,
      isWrongNetwork: chainId !== expectedChainId,
    });

    const cleanup = setupEventListeners(detectedProvider, (address: string) => {
      set({ address });
    });

    return cleanup;
  },

  connectWallet: async () => {
    try {
      const detectedProvider = await getEthereumProvider();
      //@ts-ignore
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

    const desiredNetwork = isL1ToL2 ? currentNetwork.l2 : currentNetwork.l1;

    try {
      await provider?.request({
        method: requestMethods.switchNetwork,
        params: [{ chainId: desiredNetwork }],
      });

      set({ isL1ToL2: !isL1ToL2 });
      handleStorage.save("isL1ToL2", (!isL1ToL2).toString());

      toast({
        title: "Network Switched",
        description: `Switched to ${isL1ToL2 ? "L1" : "L2 TEN Testnet"}`,
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

  updateNetworkStatus: (isWrongNetwork: boolean) => {
    set({ isWrongNetwork });
  },
}));

export default useWalletStore;
