import { createContext, useContext, useEffect, useState } from "react";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  Account,
} from "../../types/interfaces/WalletInterfaces";
import { isValidTokenFormat } from "../../lib/utils";
import {
  accountIsAuthenticated,
  fetchVersion,
  revokeAccountsApi,
} from "../../api/gateway";
import { ToastType } from "@repo/ui/lib/enums/toast";
import {
  authenticateAccountWithTenGatewayEIP712,
  getToken,
} from "@/api/ethRequests";
import { ethers } from "ethers";
import ethService from "@/services/ethService";
import { showToast } from "@repo/ui/components/shared/use-toast";
import { ethereum } from "@repo/ui/lib/utils";

const WalletConnectionContext =
  createContext<WalletConnectionContextType | null>(null);

export const useWalletConnection = (): WalletConnectionContextType => {
  const context = useContext(WalletConnectionContext);

  if (!context) {
    throw new Error(
      "useWalletConnection must be used within a WalletConnectionProvider"
    );
  }
  return context;
};

export const WalletConnectionProvider = ({
  children,
}: WalletConnectionProviderProps) => {
  const [walletConnected, setWalletConnected] = useState(false);
  const [token, setToken] = useState<string>("");
  const [version, setVersion] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [accounts, setAccounts] = useState<Account[] | null>(null);
  const [provider, setProvider] = useState({} as ethers.providers.Web3Provider);

  const initialize = async (
    providerInstance: ethers.providers.Web3Provider
  ) => {
    if (!providerInstance) {
      return showToast(
        ToastType.INFO,
        "Provider is required to initialize wallet connection."
      );
    }

    try {
      await ethService.checkIfMetamaskIsLoaded(providerInstance);

      const fetchedToken = await getToken(providerInstance);
      setToken(fetchedToken);

      const status = await ethService.isUserConnectedToTenChain(fetchedToken);
      setWalletConnected(status);

      const accounts = await ethService.getAccounts(providerInstance);
      setAccounts(accounts || null);
      setVersion(await fetchVersion());
    } catch (error) {
      showToast(
        ToastType.DESTRUCTIVE,
        error instanceof Error
          ? error.message
          : "Error initializing wallet connection. Please refresh the page."
      );
    } finally {
      setLoading(false);
    }
  };

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
  };

  useEffect(() => {
    if (ethereum && ethereum.isMetaMask) {
      const providerInstance = new ethers.providers.Web3Provider(ethereum);
      setProvider(providerInstance);
      initialize(providerInstance);

      ethereum.on("accountsChanged", fetchUserAccounts);
    } else {
      setLoading(false);
    }

    return () => {
      if (ethereum && ethereum.removeListener) {
        ethereum.removeListener("accountsChanged", fetchUserAccounts);
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const walletConnectionContextValue: WalletConnectionContextType = {
    walletConnected,
    accounts,
    token,
    connectAccount,
    version,
    revokeAccounts,
    loading,
    provider,
    fetchUserAccounts,
    setLoading,
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
