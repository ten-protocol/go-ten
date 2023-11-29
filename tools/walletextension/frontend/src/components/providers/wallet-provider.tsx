import { createContext, useContext, useEffect, useState } from "react";
import { ethers } from "ethers";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  Account,
} from "../../types/interfaces/WalletInterfaces";
import { showToast } from "../ui/use-toast";
import {
  getRandomIntAsString,
  isTenChain,
  isValidUserIDFormat,
} from "../../lib/utils";
import {
  accountIsAuthenticated,
  authenticateAccountWithTenGatewayEIP712,
  fetchVersion,
  revokeAccountsApi,
} from "../../api/gateway";
import { METAMASK_CONNECTION_TIMEOUT } from "../../lib/constants";
import { requestMethods } from "@/routes";
import { ToastType } from "@/types/interfaces";

const { ethereum } = typeof window !== "undefined" ? window : ({} as any);

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
    const handleAccountsChanged = async () => {
      if (!ethereum) {
        return;
      }
      if (userID && isValidUserIDFormat(userID)) {
        await displayCorrectScreenBasedOnMetamaskAndUserID(userID, provider);
      }
    };
    ethereum.on("accountsChanged", handleAccountsChanged);

    return () => {
      ethereum.removeListener("accountsChanged", handleAccountsChanged);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    const initialize = async () => {
      await checkIfMetamaskIsLoaded();
    };
    initialize();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const checkIfMetamaskIsLoaded = async () => {
    if (ethereum) {
      await handleEthereum();
    } else {
      showToast(ToastType.INFO, "Connecting to MetaMask...");
      window.addEventListener("ethereum#initialized", handleEthereum, {
        once: true,
      });
      setTimeout(handleEthereum, METAMASK_CONNECTION_TIMEOUT);
    }
  };

  const handleEthereum = async () => {
    if (ethereum && ethereum.isMetaMask) {
      const provider = new ethers.providers.Web3Provider(ethereum);
      setProvider(provider);
      const fetchedUserID = await getUserID(provider);
      await displayCorrectScreenBasedOnMetamaskAndUserID(
        fetchedUserID,
        provider
      );
    } else {
      showToast(
        ToastType.WARNING,
        "Please install MetaMask to use Ten Gateway."
      );
    }
  };

  const getUserID = async (provider: ethers.providers.Web3Provider) => {
    if (!provider) {
      return null;
    }

    try {
      if (await isTenChain()) {
        const id = await provider.send(requestMethods.getStorageAt, [
          "getUserID",
          getRandomIntAsString(0, 1000),
          null,
        ]);
        setUserID(id);
        return id;
      } else {
        return null;
      }
    } catch (e: any) {
      showToast(
        ToastType.DESTRUCTIVE,
        `${e.message} ${e.data?.message}` ||
          "Error: Could not fetch your user ID. Please try again later."
      );
      console.error(e);
      return null;
    }
  };

  const displayCorrectScreenBasedOnMetamaskAndUserID = async (
    userID: string,
    provider: ethers.providers.Web3Provider
  ) => {
    setVersion(await fetchVersion());
    if (await isTenChain()) {
      if (userID) {
        await getAccounts(provider, userID);
      } else {
        setWalletConnected(false);
      }
    } else {
      setWalletConnected(false);
    }

    setLoading(false);
  };

  const connectAccount = async (account: string) => {
    if (!userID) {
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
      return;
    }
    const revokeResponse = await revokeAccountsApi(userID);
    if (revokeResponse === ToastType.SUCCESS) {
      showToast(ToastType.DESTRUCTIVE, "Accounts revoked!");
      setAccounts(null);
      setWalletConnected(false);
    }
  };

  const getAccounts = async (
    provider: ethers.providers.Web3Provider,
    id: string
  ) => {
    try {
      if (!provider) {
        showToast(
          ToastType.DESTRUCTIVE,
          "No provider found. Please try again later."
        );
        return;
      }

      showToast(ToastType.INFO, "Getting accounts...");

      if (!(await isTenChain())) {
        showToast(ToastType.DESTRUCTIVE, "Please connect to the Ten chain.");
        return;
      }

      const accounts = await provider.listAccounts();

      if (accounts.length === 0) {
        showToast(ToastType.DESTRUCTIVE, "No MetaMask accounts found.");
        return;
      }

      let updatedAccounts: Account[] = [];

      for (let i = 0; i < accounts.length; i++) {
        const account = accounts[i];
        authenticateAccountWithTenGatewayEIP712(id, account);
        const { status } = await accountIsAuthenticated(id, account);
        updatedAccounts.push({
          name: account,
          connected: status,
        });
      }

      setAccounts(updatedAccounts);
      setWalletConnected(true);

      showToast(ToastType.SUCCESS, "Accounts authenticated");
    } catch (error) {
      console.error(error);
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
    }
  };

  const walletConnectionContextValue: WalletConnectionContextType = {
    walletConnected,
    accounts,
    userID,
    setUserID,
    connectAccount,
    provider,
    version,
    revokeAccounts,
    getAccounts,
    loading,
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
