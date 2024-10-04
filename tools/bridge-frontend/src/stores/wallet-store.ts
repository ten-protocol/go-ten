import { create } from "zustand";
import { walletService } from "../services/walletService";
import { handleStorage } from "../lib/utils";
import { IWalletState } from "../types";

const useWalletStore = create<IWalletState>((set, get) => ({
  provider: null,
  signer: null,
  address: "",
  walletConnected: false,
  isL1ToL2: true,
  isWrongNetwork: false,
  loading: true,

  initializeProvider: () => walletService.initializeProvider(set, get),
  connectWallet: () => walletService.connectWallet(set, get),
  disconnectWallet: () => walletService.disconnectWallet(set, get),
  switchNetwork: () => walletService.switchNetwork(set, get),
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

    walletService.initializeProvider(set, get);
  },
}));

export default useWalletStore;
