import { ethers } from "ethers";

export interface WalletConnectionContextType {
  connectToTenTestnet: () => Promise<void>;
  accounts: string[] | null;
  walletConnected: boolean;
  walletAddress: string | null;
  connectWallet: () => Promise<void>;
  disconnectWallet: () => void;
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
