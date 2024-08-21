import { ethers } from "ethers";

export interface WalletConnectionContextType {
  accounts: Account[] | null;
  walletConnected: boolean;
  loading: boolean;
  provider: ethers.providers.Web3Provider | null;
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

export interface IWalletState {
  walletConnected: boolean;
  token: string;
  version: string | null;
  accounts: Account[] | null;
  provider: ethers.providers.Web3Provider | null;
  loading: boolean;
  initialize: (
    providerInstance: ethers.providers.Web3Provider
  ) => Promise<void>;
  connectAccount: (account: string) => Promise<void>;
  revokeAccounts: () => Promise<void>;
  fetchUserAccounts: () => Promise<void>;
}
