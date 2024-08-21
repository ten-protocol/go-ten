// import { IWalletState } from "@/types/interfaces/WalletInterfaces";
import { Account } from "@/types/interfaces/WalletInterfaces";
import { ethers } from "ethers";
import create from "zustand";

 interface IWalletState {
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

export const useWalletStore = create<IWalletState>((set) => ({
  walletConnected: false,
  token: "",
  version: null,
  accounts: null,
  provider: null,
  loading: true,

  initialize: async (providerInstance) => {
    try {
      // Initialize wallet connection and update state
      const fetchedToken = await getToken(providerInstance);
      const accounts = await ethService.getAccounts(providerInstance);
      set({
        token: fetchedToken,
        accounts,
        walletConnected: true,
        provider: providerInstance,
        loading: false,
      });
    } catch (error) {
      set({ loading: false });
    }
  },

  const connectAccount = async (account: string) => {
    try {
      if (!token) {
        showToast(
          ToastType.INFO,
          "Encryption token is required to connect an account."
        );
        return;
      }
      await authenticateAccountWithTenGatewayEIP712(token, account);
      const { status } = await accountIsAuthenticated(token, account);
      if (status) {
        showToast(ToastType.SUCCESS, "Account authenticated!");
        setAccounts((accounts) => {
          if (!accounts) {
            return null;
          }
          return accounts.map((acc) => {
            if (acc.name === account) {
              return {
                ...acc,
                connected: status,
              };
            }
            return acc;
          });
        });
      } else {
        showToast(ToastType.DESTRUCTIVE, "Account authentication failed.");
      }
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, "Account authentication failed.");
    }
  };

  const revokeAccounts = async () => {
    if (!token) {
      showToast(
        ToastType.INFO,
        "Encryption token is required to revoke accounts"
      );
      return;
    }
    const revokeResponse = await revokeAccountsApi(token);
    if (revokeResponse === ToastType.SUCCESS) {
      showToast(ToastType.DESTRUCTIVE, "Accounts revoked!");
      setAccounts(null);
      setWalletConnected(false);
      setToken("");
    }
  };

  const fetchUserAccounts = async () => {
    if (!provider) {
      showToast(
        ToastType.INFO,
        "Provider is required to fetch user accounts. Please connect your wallet."
      );
      return;
    }
    const token = await getToken(provider);

    if (!isValidTokenFormat(token)) {
      showToast(
        ToastType.INFO,
        "Invalid token format. Please refresh the page."
      );
      setAccounts([]);
      setWalletConnected(false);
      return;
    }

    setToken(token);

    try {
      const accounts = await ethService.getAccounts(provider);
      let updatedAccounts: Account[] = [];

      if (!accounts || accounts.length === 0) {
        setAccounts([]);
      } else {
        updatedAccounts = await Promise.all(
          accounts!.map(async (account) => {
            await ethService.authenticateWithGateway(token, account.name);
            const { status } = await accountIsAuthenticated(
              token,
              account.name
            );
            return {
              ...account,
              connected: status,
            };
          })
        );
        showToast(ToastType.INFO, "Accounts authenticated with gateway!");
        setAccounts(updatedAccounts);
      }
    } catch (error: any) {
      showToast(
        ToastType.DESTRUCTIVE,
        `Error fetching user accounts: ${error?.message}`
      );
    } finally {
      setWalletConnected(true);
      setLoading(false);
    }
    },
}));