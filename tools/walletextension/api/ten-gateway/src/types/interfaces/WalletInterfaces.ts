import { ethers } from "ethers";

export interface WalletConnectionContextType {
  connectToTenTestnet: () => Promise<void>;
  accounts: Account[] | null;
  walletConnected: boolean;
  connectAccount: (account: string) => Promise<void>;
  revokeAccounts: () => Promise<void>;
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
