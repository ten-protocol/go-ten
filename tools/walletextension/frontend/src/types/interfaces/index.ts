import { IWalletState } from "@repo/ui/lib/interfaces/wallet";
import { Account } from "./WalletInterfaces";

export type NavLink = {
  label: string;
  href?: string;
  isDropdown?: boolean;
  isExternal?: boolean;
  subNavLinks?: NavLink[];
};

export type StoreSet = (
  partial:
    | IGatewayWalletState
    | Partial<IGatewayWalletState>
    | ((
        state: IGatewayWalletState
      ) => IGatewayWalletState | Partial<IGatewayWalletState>),
  replace?: boolean | undefined
) => void;

export type StoreGet = () => IGatewayWalletState;

interface IGatewayWalletState extends IWalletState {
  token: string;
  version: string | null;
  accounts: Account[] | null;

  initializeGateway: () => void;
  connectAccount: (account: string) => void;
  revokeAccounts: () => void;
}
