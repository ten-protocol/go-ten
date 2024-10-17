import { create } from "zustand";
import ethService from "@/services/ethService";
import { IGatewayWalletState } from "@/types/interfaces";

const useWalletStore = create<IGatewayWalletState>((set, get) => ({
  ...get(),
  token: "",
  version: null,
  accounts: null,
  loading: true,
  walletConnected: false,
  provider: null,

  initializeGateway: async () => {
    await ethService.initializeGateway(set, get);
  },

  connectAccount: async (account: string) => {
    await ethService.connectAccount(set, get, account);
  },

  revokeAccounts: async () => {
    await ethService.revokeAccounts(set, get);
  },

  setLoading: (loading: boolean) => {
    set({ loading });
  },

  fetchUserAccounts: async () => {
    const { provider } = get();
    if (!provider) {
      return;
    }
    const accounts = await ethService.getAccounts(provider);
    set({ accounts });
  },
}));

export default useWalletStore;
