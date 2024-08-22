import {
  authenticateAccountWithTenGatewayEIP712,
  getToken,
} from "@/api/ethRequests";
import { accountIsAuthenticated, revokeAccountsApi } from "@/api/gateway";
import { showToast } from "@/components/ui/use-toast";
import { tenChainIDDecimal } from "@/lib/constants";
import { handleError, validateToken } from "@/lib/utils/walletUtils";
import ethService from "@/services/ethService";
import { ToastType } from "@/types/interfaces";
import { Account, IWalletState } from "@/types/interfaces/WalletInterfaces";
import { ethers } from "ethers";
import { create } from "zustand";

export const useWalletStore = create<IWalletState>((set, get) => ({
  walletConnected: false,
  token: "",
  version: null,
  accounts: null,
  provider: null,
  loading: true,
  isWrongNetwork: false,
  setLoading: (loading: boolean) => set({ loading }),

  initialize: async (providerInstance: ethers.providers.Web3Provider) => {
    try {
      const fetchedToken = await getToken(providerInstance);
      const accounts = await ethService.getAccounts(providerInstance);

      const network = await providerInstance.getNetwork();

      set({
        token: fetchedToken,
        accounts,
        walletConnected: true,
        provider: providerInstance,
        loading: false,
        isWrongNetwork: network.chainId !== tenChainIDDecimal,
      });
    } catch (error: any) {
      handleError(error, "Failed to initialize wallet connection");
      set({ loading: false });
    }
  },

  connectAccount: async (account: string) => {
    const { token } = get();

    if (!token) {
      return showToast(
        ToastType.INFO,
        "Encryption token is required to connect an account."
      );
    }

    try {
      await authenticateAccountWithTenGatewayEIP712(token, account);
      const { status } = await accountIsAuthenticated(token, account);

      const updateAccountStatus = (account: string, status: boolean) => {
        set((state: IWalletState) => ({
          accounts: state.accounts?.map((acc) =>
            acc.name === account ? { ...acc, connected: status } : acc
          ),
        }));
      };

      if (status) {
        updateAccountStatus(account, true);
        showToast(ToastType.SUCCESS, "Account authenticated!");
      } else {
        throw new Error("Account authentication failed.");
      }
    } catch (error: any) {
      handleError(error, "Account authentication failed");
    }
  },

  revokeAccounts: async () => {
    const { token } = get();

    if (!token) {
      showToast(
        ToastType.INFO,
        "Encryption token is required to revoke accounts"
      );
    }

    try {
      const revokeResponse = await revokeAccountsApi(token);

      if (revokeResponse === ToastType.SUCCESS) {
        showToast(ToastType.INFO, "Accounts revoked!");
        set({
          accounts: null,
          walletConnected: false,
          token: "",
        });
      }
    } catch (error: any) {
      handleError(error, "Failed to revoke accounts");
    }
  },

  fetchUserAccounts: async () => {
    const { provider } = get();

    if (!provider) {
      return showToast(
        ToastType.INFO,
        "Provider is required to fetch user accounts. Please connect your wallet."
      );
    }

    try {
      const token = await getToken(provider);
      validateToken(token);

      const accounts = await ethService.getAccounts(provider);

      if (!accounts || accounts.length === 0) {
        set({ accounts: [] });
        return;
      }

      const authenticateAccountsWithGateway = async (
        accounts: Account[],
        token: string
      ) => {
        return await Promise.all(
          accounts.map(async (account) => {
            await ethService.authenticateWithGateway(token, account.name);
            const { status } = await accountIsAuthenticated(
              token,
              account.name
            );
            return { ...account, connected: status };
          })
        );
      };

      const updatedAccounts = await authenticateAccountsWithGateway(
        accounts,
        token
      );
      console.log(
        "🚀 ~ fetchUserAccounts: ~ updatedAccounts:",
        updatedAccounts
      );

      showToast(ToastType.INFO, "Accounts authenticated with gateway!");
      set({ accounts: updatedAccounts });
    } catch (error: any) {
      handleError(error, "Error fetching user accounts");
    } finally {
      set({ walletConnected: true, loading: false });
    }
  },
}));
