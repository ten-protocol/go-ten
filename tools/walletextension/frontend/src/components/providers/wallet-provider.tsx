import { createContext, useContext, useEffect, useState } from "react";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  Account,
} from "../../types/interfaces/WalletInterfaces";
import { showToast } from "../ui/use-toast";
import { ethereum, isValidTokenFormat } from "../../lib/utils";
import {
  accountIsAuthenticated,
  fetchVersion,
  revokeAccountsApi,
} from "../../api/gateway";
import { ToastType } from "@/types/interfaces";
import {
  authenticateAccountWithTenGatewayEIP712,
  getToken,
} from "@/api/ethRequests";
import { ethers } from "ethers";
import ethService from "@/services/ethService";

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

  const initialize = async () => {
    if (!provider) {
      return showToast(
        ToastType.INFO,
        "Provider is required to initialize wallet connection."
      );
    }

    try {
      await ethService.checkIfMetamaskIsLoaded(provider);

      const fetchedToken = await getToken(provider);
      setToken(fetchedToken);

      const status = await ethService.isUserConnectedToTenChain(fetchedToken);
      setWalletConnected(status);

      const accounts = await ethService.getAccounts(provider);
      setAccounts(accounts || null);
      setVersion(await fetchVersion());
    } catch (error) {
      showToast(
        ToastType.DESTRUCTIVE,
        "Error initializing wallet connection. Please refresh the page."
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

    showToast(ToastType.SUCCESS, "Token fetched!");
    setToken(token);

    try {
      const accounts = await ethService.getAccounts(provider);
      showToast(ToastType.SUCCESS, "Accounts fetched!");
      let updatedAccounts: Account[] = [];

      // updatedAccounts = await Promise.all(
      //   accounts!.map(async (account) => {
      //     await ethService.authenticateWithGateway(token, account.name);
      //     showToast(ToastType.SUCCESS, "Account authenticated with gateway!");
      //     const { status } = await accountIsAuthenticated(token, account.name);
      //     showToast(ToastType.SUCCESS, "Account authenticated!");
      //     return {
      //       ...account,
      //       connected: status,
      //     };
      //   })
      // );
      if (!accounts || accounts.length === 0) {
        setAccounts([]);
      } else {
        for (const account of accounts) {
          await ethService.authenticateWithGateway(token, account.name);
          showToast(ToastType.SUCCESS, "Account authenticated with gateway!");
          const { status } = await accountIsAuthenticated(token, account.name);
          showToast(ToastType.SUCCESS, "Account authenticated!");
          updatedAccounts.push({
            ...account,
            connected: status,
          });
        }
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
      initialize();

      ethereum.on("accountsChanged", () => {
        fetchUserAccounts();
      });
    }

    return () => {
      if (ethereum && ethereum.removeListener) {
        ethereum.removeListener("accountsChanged", () => {
          fetchUserAccounts();
        });
      }
    };
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
