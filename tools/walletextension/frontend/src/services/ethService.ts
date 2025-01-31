import { ethers } from "ethers";
import {
  authenticateAccountWithTenGatewayEIP712,
  getToken,
} from "@/api/ethRequests";
import { accountIsAuthenticated } from "@/api/gateway";
import { showToast } from "@/components/ui/use-toast";
import { METAMASK_CONNECTION_TIMEOUT } from "@/lib/constants";
import { isTenChain, isValidTokenFormat, ethereum } from "@/lib/utils";
import { ToastType } from "@/types/interfaces";
import { Account } from "@/types/interfaces/WalletInterfaces";

const ethService = {
  checkIfMetamaskIsLoaded: async (provider: ethers.providers.Web3Provider) => {
    try {
      if (ethereum) {
        // There are some wallets that are conflicting with MetaMask - we want to check that and throw an error if they are connected
        const conflictingWalletMap = {
          "Exodus Wallet": ethereum.isExodus,
          "Nest Wallet": ethereum.isNestWallet,
          // Add other wallets here as needed
        };

        // Iterate over the wallet map and handle conflicts
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
          {
            once: true,
          }
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
      if (token && isValidTokenFormat(token)) {
        return true;
      } else {
        return false;
      }
    } else {
      return false;
    }
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
    let updatedAccounts: Account[] = [];
    showToast(ToastType.INFO, "Checking account authentication status...");
    const authenticationPromise = accounts.map((account) =>
      accountIsAuthenticated(token, account).then(({ status }) => {
        return {
          name: account,
          connected: status,
        };
      })
    );
    updatedAccounts = await Promise.all(authenticationPromise);
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

    const token = await getToken();

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
      return await authenticateAccountWithTenGatewayEIP712(token, account);
    } catch (error) {
      showToast(
        ToastType.DESTRUCTIVE,
        `Error authenticating account: ${account}`
      );
    }
  },
};

export default ethService;
