import { ethers } from "ethers";

export interface WalletConnectionContextType {
  accounts: Account[] | null;
  walletConnected: boolean;
  connectAccount: (account: string) => Promise<void>;
  token: string | null;
  version: string | null;
  revokeAccounts: () => void;
  loading: boolean;
  provider: ethers.providers.Web3Provider;
  fetchUserAccounts: () => Promise<void>;
  setLoading: (loading: boolean) => void;
}

export interface State {
  hasError: boolean;
}

export interface WalletConnectionProviderProps {
  children: React.ReactNode;
}

export type Account = {
  name: string;
  connected: boolean;
};
