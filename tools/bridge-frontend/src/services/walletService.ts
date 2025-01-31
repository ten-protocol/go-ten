import { ethers } from "ethers";
import { toast } from "@/src/components/ui/use-toast";
import { getEthereumProvider, handleError } from "@/src/lib/utils/walletUtils";
import {
  initializeSigner,
  setupEventListeners,
} from "@/src/lib/utils/walletEvents";
import { StoreGet, StoreSet, ToastType } from "../types";
import { currentNetwork, handleStorage } from "../lib/utils";
import { requestMethods } from "../routes";
import { environment } from "../lib/constants";

export const walletService = {
  initializeProvider: async (set: StoreSet, get: StoreGet) => {
    try {
      let detectedProvider = await getEthereumProvider();
      if (!detectedProvider) {
        let attempts = 0;
        const maxAttempts = 5;
        const intervalId = setInterval(async () => {
          attempts++;
          try {
            detectedProvider = await getEthereumProvider();
            if (detectedProvider) {
              clearInterval(intervalId);
              set({ provider: detectedProvider, loading: false });
              await walletService.proceedWithProviderInitialization(set, get);
            } else if (attempts >= maxAttempts) {
              clearInterval(intervalId);
              toast({
                title: "Unable to connect",
                description:
                  "Please check if your web3 provider is installed and unlocked.",
                variant: ToastType.INFO,
              });
            }
          } catch (error) {
            console.error("Error detecting provider:", error);
          }
        }, 3000);
        return;
      }
      set({ provider: detectedProvider, loading: false });
      await walletService.proceedWithProviderInitialization(set, get);
    } catch (error) {
      console.error("Error initializing provider:", error);
    }
  },

  proceedWithProviderInitialization: async (set: StoreSet, get: StoreGet) => {
    try {
      const { provider: detectedProvider } = get();
      await detectedProvider?.send("eth_requestAccounts", []);
      const signer = await initializeSigner(detectedProvider!);

      const network = await detectedProvider?.getNetwork();
      const chainId = network?.chainId;
      const isL1 = chainId === currentNetwork.l1;
      const expectedChainId = isL1 ? currentNetwork.l1 : currentNetwork.l2;

      set({
        signer: signer,
        isL1ToL2: isL1,
        isWrongNetwork: chainId !== expectedChainId,
      });

      const cleanup = setupEventListeners((address) => set({ address }));
      return cleanup;
    } catch (error) {
      console.error("Error during provider initialization:", error);
    }
  },

  connectWallet: async (set: StoreSet, get: StoreGet) => {
    const { isL1ToL2 } = get() as { isL1ToL2: boolean };
    try {
      const detectedProvider = await getEthereumProvider();
      const accounts = await detectedProvider?.send(
        requestMethods.requestAccounts,
        []
      );
      const newSigner = await initializeSigner(detectedProvider);

      set({
        provider: detectedProvider,
        signer: newSigner,
        address: accounts[0],
        walletConnected: true,
      });

      handleStorage.save("walletAddress", accounts[0]);
      handleStorage.save("isL1ToL2", isL1ToL2.toString());

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

  disconnectWallet: async (set: StoreSet, get: StoreGet) => {
    const { provider } = get() as { provider: any };
    try {
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

  switchNetwork: async (set: StoreSet, get: StoreGet) => {
    const { provider, isL1ToL2 } = get() as {
      provider: any;
      isL1ToL2: boolean;
    };
    if (!provider) {
      toast({
        title: "Error",
        description: "Please connect to wallet first",
        variant: ToastType.DESTRUCTIVE,
      });
      return;
    }

    set({ loading: true });

    const desiredNetwork = ethers.utils.hexValue(
      isL1ToL2 ? currentNetwork.l2 : currentNetwork.l1
    );

    try {
      await provider?.send(requestMethods.switchNetwork, [
        { chainId: desiredNetwork },
      ]);

      set({ isL1ToL2: !isL1ToL2 });
      handleStorage.save("isL1ToL2", (!isL1ToL2).toString());

      toast({
        title: "Network Switched",
        description: `Switched to ${isL1ToL2 ? "L1" : "L2 TEN Testnet"}`,
        variant: ToastType.SUCCESS,
      });
    } catch (error: any) {
      if (error.code === 4902) {
        toast({
          title: "Network not found",
          description: error.message || "Network not found in wallet",
          variant: ToastType.INFO,
        });
        if (isL1ToL2) {
          toast({
            title: "Network not found",
            description: "Redirecting to TEN Gateway...",
            variant: ToastType.INFO,
          });
          return window.open(currentNetwork.l2Gateway, "_blank");
        }

        const networkConfig = {
          chainId: desiredNetwork,
          rpcUrls: [currentNetwork.l1Rpc],
          chainName: environment,
          nativeCurrency: {
            name: "ETH",
            symbol: "ETH",
            decimals: 18,
          },
          blockExplorerUrls: [currentNetwork.l1Explorer],
        };

        await provider?.send(requestMethods.addNetwork, [networkConfig]);
      } else {
        toast({
          title: "Error",
          description: error.message || "Error switching network",
          variant: ToastType.DESTRUCTIVE,
        });
      }
      handleError(error, "Error switching network");
    } finally {
      set({ loading: false });
    }
  },
};
