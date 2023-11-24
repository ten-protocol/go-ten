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
    setVersion(await fetchVersion());

    if (await isTenChain()) {
      const userID = await getUserID();
      setUserID(userID);
      if (provider && userID && isValidUserIDFormat(userID)) {
        await getAccounts();
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

    const revokeResponse = await revokeAccountsApi(userID);
    if (revokeResponse === "success") {
      toast({
        variant: "success",
        description: "Successfully revoked all accounts.",
      });
      setAccounts(null);
      setWalletConnected(false);
    }
  };

  const getAccounts = async () => {
    if (!provider) {
      toast({
        variant: "destructive",
        description: "No provider found. Please try again later.",
      });
      return;
    }

    toast({ variant: "info", description: "Getting accounts..." });
    if (!(await isTenChain())) {
      toast({
        variant: "warning",
        description: "Please connect to the Ten chain.",
      });
      return;
    }
    const accounts = await provider.listAccounts();

    if (accounts.length === 0) {
      toast({
        variant: "destructive",
        description: "No MetaMask accounts found.",
      });
      return;
    }

    const user = await getUserID();
    setUserID(user);

    setAccounts(
      await Promise.all(
        accounts.map(async (account) => ({
          name: account,
          connected: await accountIsAuthenticated(user, account),
        }))
      )
    );
    setWalletConnected(true);
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
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
