import { createContext, useContext, useEffect, useState } from "react";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  Account,
} from "../../types/interfaces/WalletInterfaces";
import { showToast } from "../ui/use-toast";
import { ethereum } from "../../lib/utils";
import {
  accountIsAuthenticated,
  fetchVersion,
  revokeAccountsApi,
} from "../../api/gateway";
import { ToastType } from "@/types/interfaces";
import {
  authenticateAccountWithTenGatewayEIP712,
  getUserID,
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
  const [userID, setUserID] = useState<string>("");
  const [version, setVersion] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [accounts, setAccounts] = useState<Account[] | null>(null);
  const [provider, setProvider] = useState({} as ethers.providers.Web3Provider);

  useEffect(() => {
    initialize();
    ethereum.on("accountsChanged", handleAccountsChanged);

    return () => {
      ethereum.removeListener("accountsChanged", handleAccountsChanged);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleAccountsChanged = async () => {
    if (!ethereum) {
      return;
    }
    const status = await ethService.isUserConnectedToTenChain(userID);
    setWalletConnected(status);
  };

  const initialize = async () => {
    const providerInstance = new ethers.providers.Web3Provider(ethereum);
    setProvider(providerInstance);
    await ethService.checkIfMetamaskIsLoaded(providerInstance);
    const id = await getUserID(providerInstance);
    setUserID(id);
    handleAccountsChanged();
    const accounts = await ethService.getAccounts(providerInstance);
    setAccounts(accounts || null);
    setVersion(await fetchVersion());
    setLoading(false);
  };

  const connectAccount = async (account: string) => {
    if (!userID) {
      showToast(ToastType.INFO, "User ID is required to connect an account.");
      return;
    }
    await authenticateAccountWithTenGatewayEIP712(userID, account);
    const { status } = await accountIsAuthenticated(userID, account);
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
  };

  const revokeAccounts = async () => {
    if (!userID) {
      showToast(ToastType.INFO, "User ID is required to revoke accounts");
      return;
    }
    const revokeResponse = await revokeAccountsApi(userID);
    if (revokeResponse === ToastType.SUCCESS) {
      showToast(ToastType.DESTRUCTIVE, "Accounts revoked!");
      setAccounts(null);
      setWalletConnected(false);
      setUserID("");
    }
  };

  const fetchUserAccounts = async () => {
    const accounts = await ethService.getAccounts(provider);
    setAccounts(accounts || null);
    setWalletConnected(true);
  };

  const walletConnectionContextValue: WalletConnectionContextType = {
    walletConnected,
    accounts,
    userID,
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
