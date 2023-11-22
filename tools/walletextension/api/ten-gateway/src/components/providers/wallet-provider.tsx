import { createContext, useContext, useEffect, useState } from "react";
import { ethers } from "ethers";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  Account,
} from "@/types/interfaces/WalletInterfaces";
import { useToast } from "../ui/use-toast";
import {
  getRandomIntAsString,
  isTenChain,
  isValidUserIDFormat,
} from "@/lib/utils";
import {
  accountIsAuthenticated,
  authenticateAccountWithTenGateway,
  fetchVersion,
  revokeAccountsApi,
} from "@/api/gateway";
import { METAMASK_CONNECTION_TIMEOUT } from "@/lib/constants";

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
  const { toast } = useToast();

  const [walletConnected, setWalletConnected] = useState(false);
  const [userID, setUserID] = useState<string | null>(null);
  const [version, setVersion] = useState<string | null>(null);
  const [accounts, setAccounts] = useState<Account[] | null>(null);
  const [provider, setProvider] =
    useState<ethers.providers.Web3Provider | null>(null);

  useEffect(() => {
    const { ethereum } = window as any;
    const handleAccountsChanged = async () => {
      if (userID && isValidUserIDFormat(userID)) {
        await displayCorrectScreenBasedOnMetamaskAndUserID();
      }
    };

    ethereum.on("accountsChanged", handleAccountsChanged);

    return () => {
      ethereum.removeListener("accountsChanged", handleAccountsChanged);
    };
  }, []);

  useEffect(() => {
    const initialize = async () => {
      await checkIfMetamaskIsLoaded();
    };
    initialize();
  }, []);

  const checkIfMetamaskIsLoaded = async () => {
    if (window && (window as any).ethereum) {
      await handleEthereum();
    } else {
      toast({ description: "Connecting to MetaMask..." });
      window.addEventListener("ethereum#initialized", handleEthereum, {
        once: true,
      });
      setTimeout(handleEthereum, METAMASK_CONNECTION_TIMEOUT);
    }
  };

  const handleEthereum = async () => {
    const { ethereum } = window as any;
    if (ethereum && ethereum.isMetaMask) {
      const provider = new ethers.providers.Web3Provider(
        (window as any).ethereum
      );
      setProvider(provider);
      await displayCorrectScreenBasedOnMetamaskAndUserID();
    } else {
      toast({ description: "Please install MetaMask to use Ten Gateway." });
    }
  };

  const getUserID = async () => {
    if (!provider) {
      console.log("no provider");
      return null;
    }

    try {
      if (await isTenChain()) {
        const id = await provider.send("eth_getStorageAt", [
          "getUserID",
          getRandomIntAsString(0, 1000),
          null,
        ]);
        return id;
      } else {
        return null;
      }
    } catch (e: any) {
      toast({
        description:
          `${e.message} ${e.data?.message}` ||
          "Error: Could not fetch your userID. Please try again later.",
      });
      console.error(e);
      return null;
    }
  };

  const displayCorrectScreenBasedOnMetamaskAndUserID = async () => {
    const userID = await getUserID();
    setUserID(userID);
    setVersion(await fetchVersion());

    if (await isTenChain()) {
      if (provider && userID && isValidUserIDFormat(userID)) {
        const accounts = await provider.listAccounts();
        const formattedAccounts = await Promise.all(
          accounts.map(async (account) => ({
            name: account,
            connected: await accountIsAuthenticated(userID, account),
          }))
        );
        setAccounts(formattedAccounts);
        setWalletConnected(true); // Set walletConnected after all async operations
      } else {
        setWalletConnected(false);
      }
    } else {
      setWalletConnected(false);
    }
  };

  const connectAccount = async (account: string) => {
    if (!userID) {
      return;
    }
    await authenticateAccountWithTenGateway(userID, account);
  };

  const revokeAccounts = async () => {
    if (!userID) {
      return;
    }

    await revokeAccountsApi(userID);
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
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
