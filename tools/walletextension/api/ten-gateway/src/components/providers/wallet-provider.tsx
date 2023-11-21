import { createContext, useContext, useEffect, useState } from "react";
import { ethers } from "ethers";
import {
  WalletConnectionContextType,
  WalletConnectionProviderProps,
} from "@/types/interfaces/WalletInterfaces";
import { useToast } from "../ui/use-toast";
import {
  getNetworkName,
  getRPCFromUrl,
  getRandomIntAsString,
  isValidUserIDFormat,
  metamaskPersonalSign,
  pathAuthenticate,
  pathJoin,
  pathQuery,
  pathRevoke,
  pathVersion,
  tenChainIDHex,
  tenGatewayVersion,
} from "@/lib/utils";
import { set } from "date-fns";

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
  const [accounts, setAccounts] = useState<string[] | null>(null);
  const [walletAddress, setWalletAddress] = useState<string | null>(null);
  const [provider, setProvider] =
    useState<ethers.providers.Web3Provider | null>(null);

  const connectWallet = async () => {
    if ((window as any).ethereum) {
      const ethProvider = new ethers.providers.Web3Provider(
        (window as any).ethereum
      );
      setProvider(ethProvider);

      try {
        await ethProvider.send("eth_requestAccounts", []);
        const signer = ethProvider.getSigner();
        const address = await signer.getAddress();
        setWalletAddress(address);
        setWalletConnected(true);
      } catch (error: any) {
        toast({ description: "Error connecting to wallet:" + error.message });
      }
    } else {
      toast({ description: "No ethereum object found." });
    }
  };

  const disconnectWallet = () => {
    if (provider) {
      provider.removeAllListeners();
      setWalletConnected(false);
      setWalletAddress(null);
      setProvider(null);
    }
  };

  async function getUserID() {
    if (!provider) {
      return null;
    }
    try {
      if (await isTenChain()) {
        return await provider.send("eth_getStorageAt", [
          "getUserID",
          getRandomIntAsString(0, 1000),
          null,
        ]);
      } else {
        return null;
      }
    } catch (e) {
      console.log(e);
      return null;
    }
  }

  async function isTenChain() {
    let currentChain = await (window as any).ethereum.request({
      method: "eth_chainId",
    });
    return currentChain === tenChainIDHex;
  }

  function checkIfMetamaskIsLoaded() {
    if (window && (window as any).ethereum) {
      handleEthereum();
    } else {
      toast({ description: "Connecting to Metamask..." });
      window.addEventListener("ethereum#initialized", handleEthereum, {
        once: true,
      });

      // If the event is not dispatched by the end of the timeout,
      // the user probably doesn't have MetaMask installed.
      setTimeout(handleEthereum, 3000); // 3 seconds
    }
  }

  function handleEthereum() {
    const { ethereum } = window as any;
    if (ethereum && ethereum.isMetaMask) {
      initialize();
    } else {
      toast({ description: "Please install MetaMask to use Obscuro Gateway." });
    }
  }

  const initialize = async () => {
    // getUserID from the gateway with getStorageAt method
    let userID = await getUserID();

    await displayCorrectScreenBasedOnMetamaskAndUserID();
    await fetchAndDisplayVersion();
    await displayCorrectScreenBasedOnMetamaskAndUserID();
  };

  const revokeUserID = async () => {
    const queryAccountUserID = pathRevoke + "?u=" + userID;
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
        toast({
          description: "Failed to fetch the version. Please try again later.",
        });
        throw new Error("Failed to fetch the version");
      }

      let response = await versionResp.text();
      setVersion(response);
    } catch (error) {
      console.error("Error fetching the version:", error);
      toast({ description: `Error fetching the version: ${error}` });
    }
  }

  async function displayCorrectScreenBasedOnMetamaskAndUserID() {
    console.log("displayCorrectScreenBasedOnMetamaskAndUserID");
    // check if we are on Obscuro Chain
    if (await isTenChain()) {
      // check if we have valid userID in rpcURL
      if (!provider) {
        return;
      }
      if (userID && isValidUserIDFormat(userID)) {
        const accounts = await provider.listAccounts();
        setAccounts(accounts);
      }
      setWalletConnected(true);
      return;
    }
    setWalletConnected(false);
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
    return -1;
  }

  async function addNetworkToMetaMask() {
    // add network to MetaMask
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
            rpcUrls: [
              getRPCFromUrl() + "/" + tenGatewayVersion + "/?u=" + userID,
            ],
            blockExplorerUrls: ["https://testnet.obscuroscan.io"],
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
      // TODO: Display warning to user to allow it and refresh page...
      console.error("User denied account access:", error);
      toast({ description: `User denied account access: ${error}` });
      return null;
    }
  }

  async function isMetamaskConnected() {
    let accounts;
    if (!provider) {
      return false;
    }
    try {
      accounts = await provider.listAccounts();
      return accounts.length > 0;
    } catch (error) {
      console.log("Unable to get accounts");
    }
    return false;
  }

  async function accountIsAuthenticated(account: string) {
    const queryAccountUserID = pathQuery + "?u=" + userID + "&a=" + account;
    const isAuthenticatedResponse = await fetch(queryAccountUserID, {
      method: "get",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });
    let response = await isAuthenticatedResponse.text();
    let jsonResponseObject = JSON.parse(response);
    return jsonResponseObject.status;
  }

  async function authenticateAccountWithTenGateway(account: string) {
    const isAuthenticated = await accountIsAuthenticated(account);
    if (isAuthenticated) {
      return "Account is already authenticated";
    }

    const textToSign = "Register " + userID + " for " + account.toLowerCase();
    const signature = await (window as any).ethereum
      .request({
        method: metamaskPersonalSign,
        params: [textToSign, account],
      })
      .catch((error: any) => {
        return -1;
      });
    if (signature === -1) {
      return "Signing failed";
    }

    const authenticateUserURL = pathAuthenticate + "?u=" + userID;
    const authenticateFields = { signature: signature, message: textToSign };
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

  const connectToTenTestnet = async () => {
    // check if we are on an Ten chain
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
      // we are not on a Ten network - try to switch
      let switched = await switchToTenNetwork();
      // error 4902 means that the chain does not exist
      if (switched === 4902 || !isValidUserIDFormat(await getUserID())) {
        // join the network
        const joinResp = await fetch(pathJoin, {
          method: "get",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        });
        if (!joinResp.ok) {
          console.log("Error joining Ten Gateway");
          toast({
            description: "Error joining Ten Gateway. Please try again later.",
          });
          return;
        }
        const user = await joinResp.text();
        setUserID(user);

        // add Ten network
        await addNetworkToMetaMask();
      }

      // we have to check if user has accounts connected with metamask - and prompt to connect if not
      if (!(await isMetamaskConnected())) {
        await connectAccounts();
      }

      if (!provider) {
        return;
      }

      // connect all accounts
      // Get an accounts and prompt user to sign joining with a selected account
      const accounts = await provider.listAccounts();
      if (accounts.length === 0) {
        toast({ description: "No MetaMask accounts found." });
        return;
      }

      for (const account of accounts) {
        await authenticateAccountWithTenGateway(account);
      }
    }
  };

  useEffect(() => {
    const ethereum = (window as any).ethereum;
    const handleAccountsChanged = async (accounts: string[]) => {
      if (accounts.length === 0) {
        toast({ description: "Please connect to MetaMask." });
      } else {
        if (userID && isValidUserIDFormat(userID)) {
          for (const account of accounts) {
            await authenticateAccountWithTenGateway(account);
          }
        }
      }
    };
    ethereum.on("accountsChanged", handleAccountsChanged);
    return () => {
      ethereum.removeListener("accountsChanged", handleAccountsChanged);
    };
  });

  const user: any = async () => await getUserID();
  useEffect(() => {
    if (window) {
      const ethProvider = new ethers.providers.Web3Provider(
        (window as any).ethereum
      );
      setProvider(ethProvider);
      setUserID(user);
      fetchAndDisplayVersion();
      displayCorrectScreenBasedOnMetamaskAndUserID();
    }
  }, []);

  const walletConnectionContextValue: WalletConnectionContextType = {
    walletConnected,
    walletAddress,
    connectWallet,
    disconnectWallet,
    connectToTenTestnet,
    accounts,
  };

  return (
    <WalletConnectionContext.Provider value={walletConnectionContextValue}>
      {children}
    </WalletConnectionContext.Provider>
  );
};
