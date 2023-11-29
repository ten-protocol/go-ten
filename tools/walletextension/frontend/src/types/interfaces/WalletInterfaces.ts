import { ethers } from "ethers";

export interface WalletConnectionContextType {
  accounts: Account[] | null;
  walletConnected: boolean;
  connectAccount: (account: string) => Promise<void>;
  userID: string | null;
  setUserID: (userID: string) => void;
  provider: ethers.providers.Web3Provider | null;
  version: string | null;
  revokeAccounts: () => void;
  getAccounts: (
    provider: ethers.providers.Web3Provider,
    userID: string
  ) => Promise<void>;
  loading: boolean;
}

export interface Props {
  children: React.ReactNode;
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
