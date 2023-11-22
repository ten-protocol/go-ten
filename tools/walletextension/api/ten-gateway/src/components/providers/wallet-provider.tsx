import { createContext, useContext, useEffect, useState } from "react";
import { ethers } from "ethers";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
  Account,
} from "@/types/interfaces/WalletInterfaces";
import { useToast } from "../ui/use-toast";
import {
  METAMASK_CONNECTION_TIMEOUT,
  getNetworkName,
  getRPCFromUrl,
  getRandomIntAsString,
  isTenChain,
  isValidUserIDFormat,
  metamaskPersonalSign,
  pathAuthenticate,
  pathJoin,
  pathQuery,
  pathRevoke,
  pathVersion,
  tenChainIDHex,
  tenGatewayVersion,
  tenscanLink,
} from "@/lib/utils";

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
    const ethereum = (window as any).ethereum;

    const handleAccountsChanged = async (accounts: string[]) => {
      if (accounts.length === 0) {
        toast({ description: "Please connect to MetaMask." });
      } else if (userID && isValidUserIDFormat(userID)) {
        await Promise.all(accounts.map(authenticateAccountWithTenGateway));
      }
    };

    ethereum.on("accountsChanged", handleAccountsChanged);

    return () => {
      ethereum.removeListener("accountsChanged", handleAccountsChanged);
    };
  }, [userID]);

  useEffect(() => {
    checkIfMetamaskIsLoaded();
  }, []);

  async function checkIfMetamaskIsLoaded() {
    if (window && (window as any).ethereum) {
      const provider = new ethers.providers.Web3Provider(
        (window as any).ethereum
      );
      setProvider(provider);
      handleEthereum();
    } else {
      toast({ description: "Connecting to MetaMask..." });

      window.addEventListener("ethereum#initialized", handleEthereum, {
        once: true,
      });
      setTimeout(handleEthereum, METAMASK_CONNECTION_TIMEOUT);
    }
  }

  function handleEthereum() {
    const { ethereum } = window as any;

    if (ethereum && ethereum.isMetaMask) {
      initialize();
    } else {
      toast({ description: "Please install MetaMask to use Ten Gateway." });
    }
  }

  async function getUserID() {
    if (!provider || !(await isTenChain())) {
      return null;
    }

    try {
      return await provider.send("eth_getStorageAt", [
        "getUserID",
        getRandomIntAsString(0, 1000),
        null,
      ]);
    } catch (e) {
      console.error(e);
      return null;
    }
  }

  const initialize = async () => {
    const userID = await getUserID();
    setUserID(userID);
    await fetchAndDisplayVersion();
    await displayCorrectScreenBasedOnMetamaskAndUserID();
  };

  async function fetchAndDisplayVersion() {
    try {
      const versionResp = await fetch(pathVersion, {
        method: "get",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      });

      if (!versionResp.ok) {
        handleFetchError("Failed to fetch the version");
        return;
      }

      const response = await versionResp.text();
      setVersion(response);
    } catch (error) {
      handleFetchError(`Error fetching the version: ${error}`);
    }
  }

  async function displayCorrectScreenBasedOnMetamaskAndUserID() {
    if (await isTenChain()) {
      if (provider && userID && isValidUserIDFormat(userID)) {
        const accounts = await provider.listAccounts();
        const formattedAccounts = await Promise.all(
          accounts.map(async (account) => ({
            name: account,
            connected: await accountIsAuthenticated(account),
          }))
        );

        setAccounts(formattedAccounts);
      }

      setWalletConnected(true);
    } else {
      setWalletConnected(false);
    }
  }

  async function switchToTenNetwork() {
    try {
      await (window as any).ethereum.request({
        method: "wallet_switchEthereumChain",
        params: [{ chainId: tenChainIDHex }],
      });

      return 0;
    } catch (switchError: any) {
      return switchError.code;
    }
  }

  async function addNetworkToMetaMask() {
    try {
      await (window as any).ethereum.request({
        method: "wallet_addEthereumChain",
        params: [
          {
            chainId: tenChainIDHex,
            chainName: getNetworkName(),
            nativeCurrency: {
              name: "Sepolia Ether",
              symbol: "ETH",
              decimals: 18,
            },
            rpcUrls: [`${getRPCFromUrl()}/${tenGatewayVersion}/?u=${userID}`],
            blockExplorerUrls: [tenscanLink],
          },
        ],
      });
    } catch (error) {
      console.error(error);
      return false;
    }

    return true;
  }

  async function connectAccounts() {
    try {
      return await (window as any).ethereum.request({
        method: "eth_requestAccounts",
      });
    } catch (error) {
      console.error("User denied account access:", error);
      toast({ description: `User denied account access: ${error}` });
      return null;
    }
  }

  async function isMetamaskConnected() {
    if (!provider) {
      return false;
    }

    try {
      const accounts = await provider.listAccounts();
      return accounts.length > 0;
    } catch (error) {
      console.log("Unable to get accounts");
    }

    return false;
  }

  async function accountIsAuthenticated(account: string) {
    const queryAccountUserID = `${pathQuery}?u=${userID}&a=${account}`;
    const isAuthenticatedResponse = await fetch(queryAccountUserID, {
      method: "get",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });

    const response = await isAuthenticatedResponse.text();
    const jsonResponseObject = JSON.parse(response);
    return jsonResponseObject.status;
  }

  async function authenticateAccountWithTenGateway(account: string) {
    const isAuthenticated = await accountIsAuthenticated(account);

    if (isAuthenticated) {
      return "Account is already authenticated";
    }

    const textToSign = `Register ${userID} for ${account.toLowerCase()}`;
    const signature = await (window as any).ethereum
      .request({
        method: metamaskPersonalSign,
        params: [textToSign, account],
      })
      .catch((error: any) => -1);

    if (signature === -1) {
      return "Signing failed";
    }

    const authenticateUserURL = `${pathAuthenticate}?u=${userID}`;
    const authenticateFields = { signature, message: textToSign };
    const authenticateResp = await fetch(authenticateUserURL, {
      method: "post",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(authenticateFields),
    });

    return await authenticateResp.text();
  }

  const revokeAccounts = async () => {
    const queryAccountUserID = `${pathRevoke}?u=${userID}`;
    const revokeResponse = await fetch(queryAccountUserID, {
      method: "get",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });

    if (revokeResponse.ok) {
      setWalletConnected(false);
    } else {
      toast({ description: "Revoking UserID failed" });
    }
  };

  const connectToTenTestnet = async () => {
    if (await isTenChain()) {
      const user = await getUserID();
      setUserID(user);

      if (!isValidUserIDFormat(user)) {
        toast({
          description:
            "Existing Ten network detected in MetaMask. Please remove before hitting begin",
        });
      }
    } else {
      const switched = await switchToTenNetwork();

      if (switched === 4902 || !isValidUserIDFormat(await getUserID())) {
        const joinResp = await fetch(pathJoin, {
          method: "get",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        });

        if (!joinResp.ok) {
          handleFetchError("Error joining Ten Gateway");
          return;
        }

        const user = await joinResp.text();
        setUserID(user);

        await addNetworkToMetaMask();
      }

      if (!(await isMetamaskConnected())) {
        await connectAccounts();
      }

      if (!provider) {
        return;
      }

      const accounts = await provider.listAccounts();

      if (accounts.length === 0) {
        toast({ description: "No MetaMask accounts found." });
        return;
      }
    }
  };

  const connectAccount = async (account: string) => {
    await authenticateAccountWithTenGateway(account);
  };

  const disconnectAccount = async (account: string) => {
    const revokeAccountURL = `${pathRevoke}?u=${userID}&a=${account}`;
    const revokeAccountResp = await fetch(revokeAccountURL, {
      method: "get",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });

    if (revokeAccountResp.ok) {
      const formattedAccounts = await Promise.all(
        accounts!.map(async (acc) => ({
          name: acc.name,
          connected: await accountIsAuthenticated(acc.name),
        }))
      );

      setAccounts(formattedAccounts);
    } else {
      toast({ description: "Revoking account failed" });
    }
  };

  const walletConnectionContextValue: WalletConnectionContextType = {
    walletConnected,
    connectToTenTestnet,
    accounts,
    revokeAccounts,
    connectAccount,
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );

  function handleFetchError(errorMessage: string) {
    console.error(`Error: ${errorMessage}`);
    toast({ description: `${errorMessage}. Please try again later.` });
  }
};
