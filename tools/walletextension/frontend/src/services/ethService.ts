import { ethers } from "ethers";
import { ToastType } from "@repo/ui/lib/enums/toast";
import { showToast } from "@repo/ui/components/shared/use-toast";
import { ethereum } from "@repo/ui/lib/utils";
import {
  nativeCurrency,
  tenChainIDDecimal,
  tenChainIDHex,
  tenNetworkName,
  tenscanAddress,
  userStorageAddress,
  METAMASK_CONNECTION_TIMEOUT,
  SWITCHED_CODE,
} from "@/lib/constants";
import {
  getRandomIntAsString,
  isTenChain,
  isValidTokenFormat,
} from "@/lib/utils";
import { requestMethods } from "@/routes";
import {
  accountIsAuthenticated,
  authenticateUser,
  fetchVersion,
  revokeAccountsApi,
} from "@/api/gateway";
import { Account } from "@/types/interfaces/WalletInterfaces";

const typedData = {
  types: {
    EIP712Domain: [
      { name: "name", type: "string" },
      { name: "version", type: "string" },
      { name: "chainId", type: "uint256" },
    ],
    Authentication: [{ name: "Encryption Token", type: "address" }],
  },
  primaryType: "Authentication",
  domain: {
    name: "Ten",
    version: "1.0",
    chainId: tenChainIDDecimal,
  },
  message: {
    "Encryption Token": "0x",
  },
};

const ethService = {
  checkIfMetamaskIsLoaded: async (provider: ethers.providers.Web3Provider) => {
    try {
      if (ethereum) {
        const conflictingWalletMap = {
          "Exodus Wallet": ethereum.isExodus,
          "Nest Wallet": ethereum.isNestWallet,
        };

        for (const [walletName, isWalletConnected] of Object.entries(
          conflictingWalletMap
        )) {
          if (isWalletConnected) {
            const message = `${walletName} is connected and is conflicting with MetaMask. Please disable ${walletName} and try again.`;
            showToast(ToastType.DESTRUCTIVE, message);
            throw new Error(message);
          }
        }

        return await ethService.handleEthereum(provider);
      } else {
        showToast(ToastType.INFO, "Connecting to MetaMask...");

        let timeoutId: ReturnType<typeof setTimeout>;

        const handleEthereumOnce = async () => {
          await ethService.handleEthereum(provider);
        };

        window.addEventListener(
          "ethereum#initialized",
          () => {
            clearTimeout(timeoutId);
            handleEthereumOnce();
          },
          { once: true }
        );

        timeoutId = setTimeout(() => {
          handleEthereumOnce();
        }, METAMASK_CONNECTION_TIMEOUT);
      }
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
      throw error;
    }
  },

  handleEthereum: async (provider: ethers.providers.Web3Provider) => {
    try {
      if (ethereum && ethereum.isMetaMask) {
        return;
      } else {
        return showToast(
          ToastType.WARNING,
          "Please install MetaMask to use TEN Gateway."
        );
      }
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
      throw error;
    }
  },

  isUserConnectedToTenChain: async (token: string) => {
    if (await isTenChain()) {
      return !!(token && isValidTokenFormat(token));
    }
    return false;
  },

  formatAccounts: async (
    accounts: string[],
    provider: ethers.providers.Web3Provider,
    token: string
  ) => {
    if (!provider) {
      showToast(
        ToastType.DESTRUCTIVE,
        "No provider found. Please try again later."
      );
      return;
    }
    showToast(ToastType.INFO, "Checking account authentication status...");
    const updatedAccounts = await Promise.all(
      accounts.map((account) =>
        accountIsAuthenticated(token, account).then(({ status }) => ({
          name: account,
          connected: status,
        }))
      )
    );
    showToast(ToastType.INFO, "Account authentication status updated!");
    return updatedAccounts;
  },

  getAccounts: async (provider: ethers.providers.Web3Provider) => {
    if (!provider) {
      showToast(
        ToastType.DESTRUCTIVE,
        "No provider found. Please try again later."
      );
      return;
    }

    const token = await ethService.getToken(provider);

    if (!token || !isValidTokenFormat(token)) {
      return;
    }

    try {
      showToast(ToastType.INFO, "Getting accounts...");

      if (!(await isTenChain())) {
        showToast(ToastType.DESTRUCTIVE, "Please connect to the TEN chain.");
        return;
      }

      const accounts = await provider.listAccounts();

      if (accounts.length === 0) {
        showToast(ToastType.DESTRUCTIVE, "No MetaMask accounts found.");
        return [];
      }
      showToast(ToastType.SUCCESS, "Accounts found!");

      return ethService.formatAccounts(accounts, provider, token);
    } catch (error) {
      console.error(error);
      showToast(ToastType.DESTRUCTIVE, "An error occurred. Please try again.");
      throw error;
    }
  },

  authenticateWithGateway: async (token: string, account: string) => {
    try {
      return await ethService.authenticateAccountWithTenGatewayEIP712(
        token,
        account
      );
    } catch (error) {
      showToast(
        ToastType.DESTRUCTIVE,
        `Error authenticating account: ${account}`
      );
    }
  },

  switchToTenNetwork: async () => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      await ethereum.request({
        method: requestMethods.switchNetwork,
        params: [{ chainId: tenChainIDHex }],
      });
      return 0;
    } catch (switchError: any) {
      showToast(
        ToastType.DESTRUCTIVE,
        `switchToTenNetwork: ${switchError.code}`
      );
      return switchError.code;
    }
  },

  connectAccounts: async () => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      return await ethereum.request({
        method: requestMethods.connectAccounts,
      });
    } catch (error) {
      console.error(error);
      throw error;
    }
  },

  getSignature: async (account: string, data: any) => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      return await ethereum.request({
        method: requestMethods.signTypedData,
        params: [account, JSON.stringify(data)],
      });
    } catch (error) {
      console.error(error);
      throw "Failed to get signature";
    }
  },

  getToken: async (provider: ethers.providers.Web3Provider) => {
    if (!provider.send) {
      return null;
    }
    try {
      if (await isTenChain()) {
        const token = await provider.send(requestMethods.getStorageAt, [
          userStorageAddress,
          getRandomIntAsString(0, 1000),
          null,
        ]);
        return token;
      } else {
        return null;
      }
    } catch (e: any) {
      console.error(e);
      throw e;
    }
  },

  addNetworkToMetaMask: async (rpcUrls: string[]) => {
    if (!ethereum) {
      throw "No ethereum object found";
    }
    try {
      await ethereum.request({
        method: requestMethods.addNetwork,
        params: [
          {
            chainId: tenChainIDHex,
            chainName: tenNetworkName,
            nativeCurrency,
            rpcUrls,
            blockExplorerUrls: [tenscanAddress],
          },
        ],
      });
      return true;
    } catch (error) {
      console.error(error);
      throw error;
    }
  },

  authenticateAccountWithTenGatewayEIP712: async (
    token: string,
    account: string
  ): Promise<any> => {
    if (!token) {
      return showToast(
        ToastType.INFO,
        "Encryption token not found. Please try again later."
      );
    }

    try {
      const isAuthenticated = await accountIsAuthenticated(token, account);
      if (isAuthenticated.status) {
        return {
          status: true,
          message: "Account already authenticated",
        };
      }
      const data = {
        ...typedData,
        message: {
          ...typedData.message,
          "Encryption Token": token,
        },
      };
      const signature = await ethService.getSignature(account, data);

      return await authenticateUser(token, {
        signature,
        address: account,
      });
    } catch (error: any) {
      throw error;
    }
  },

  initializeGateway: async (set: any, get: any) => {
    try {
      const providerInstance = new ethers.providers.Web3Provider(ethereum);
      set({ provider: providerInstance });

      await ethService.checkIfMetamaskIsLoaded(providerInstance);

      const fetchedToken = await ethService.getToken(providerInstance);
      set({ token: fetchedToken });

      const status = await ethService.isUserConnectedToTenChain(fetchedToken);
      set({ walletConnected: status });

      const accounts = await ethService.getAccounts(providerInstance);
      set({ accounts: accounts || null });

      const version = await fetchVersion();
      set({ version });
    } catch (error) {
      showToast(
        ToastType.DESTRUCTIVE,
        error instanceof Error
          ? error.message
          : "Error initializing wallet connection. Please refresh the page."
      );
    } finally {
      set({ loading: false });
    }
  },

  connectAccount: async (set: any, get: any, account: string) => {
    const { token } = get();
    try {
      if (!token) {
        showToast(
          ToastType.INFO,
          "Encryption token is required to connect an account."
        );
        return;
      }
      await ethService.authenticateAccountWithTenGatewayEIP712(token, account);
      const { status } = await accountIsAuthenticated(token, account);
      if (status) {
        showToast(ToastType.SUCCESS, "Account authenticated!");
        set((state: any) => ({
          accounts:
            state.accounts?.map((acc: Account) =>
              acc.name === account ? { ...acc, connected: status } : acc
            ) || null,
        }));
      } else {
        showToast(ToastType.DESTRUCTIVE, "Account authentication failed.");
      }
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, "Account authentication failed.");
    }
  },

  revokeAccounts: async (set: any, get: any) => {
    const { token } = get();
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
      set({ accounts: null, walletConnected: false, token: "" });
    }
  },
};

export default ethService;
