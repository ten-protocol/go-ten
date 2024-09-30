import { IWalletState } from "@repo/ui/lib/interfaces/wallet";
import { Account } from "@/types/interfaces/WalletInterfaces";
import { ethService } from "@/services/ethService";
import { create } from "zustand";

interface IGatewayWalletState extends IWalletState {
  token: string;
  version: string | null;
  accounts: Account[] | null;

  initializeGateway: () => void;
  connectAccount: (account: string) => void;
  revokeAccounts: () => void;
}

const useWalletStore = create<IGatewayWalletState>((set, get) => ({
  ...get(),

  token: "",
  version: null,
  accounts: null,

  initializeGateway: () => ethService.initializeGateway(set, get),
  connectAccount: (account: string) =>
    ethService.connectAccount(set, get, account),
  revokeAccounts: () => ethService.revokeAccounts(set, get),
}));

export default useWalletStore;
