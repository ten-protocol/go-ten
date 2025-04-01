import { ToastType } from "@/types/interfaces";
import { joinTestnet } from "../api/gateway";
import { useWalletConnection } from "../components/providers/wallet-provider";
import { showToast } from "../components/ui/use-toast";
import {
  SWITCHED_CODE,
  tenGatewayAddress,
  tenGatewayVersion,
} from "../lib/constants";
import { isTenChain, isValidTokenFormat } from "../lib/utils";
import {
  addNetworkToMetaMask,
  connectAccounts,
  getToken,
  switchToTenNetwork,
} from "@/api/ethRequests";
import { fetchTestnetStatus } from "@/api/general";

const useGatewayService = () => {
  const { token, provider, fetchUserAccounts, setLoading } =
    useWalletConnection();

  const isMetamaskConnected = async () => {
    if (!provider) {
      return false;
    }
    try {
      const accounts = await provider.listAccounts();
      return accounts.length > 0;
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to get accounts");
      throw error;
    }
  };

  const finaliseConnectAccounts = async () => {
    const metamaskConnected = await isMetamaskConnected()
    try {
      if (!metamaskConnected) {
        showToast(ToastType.INFO, "No accounts found, connecting...");
        await connectAccounts();
        showToast(ToastType.SUCCESS, "Connected to TEN Testnet");
      }

      await fetchUserAccounts();
    } catch (error: Error | any) {
      showToast(ToastType.DESTRUCTIVE, `${error?.message}`);
      throw error;
    }
  }

  const connectToTenTestnet = async () => {
    showToast(ToastType.INFO, "Connecting to TEN Testnet...");
    setLoading(true);
    try {
      if (await isTenChain()) {
        if (!token || !isValidTokenFormat(token)) {
          showToast(
            ToastType.DESTRUCTIVE,
            "Existing TEN Testnet detected in MetaMask. Please remove before hitting begin"
          );
          return;
        }
      }
      showToast(ToastType.INFO, "Switching to TEN Testnet...");
      const switched = await switchToTenNetwork();
      showToast(ToastType.SUCCESS, `Switched to TEN Testnet`);
      // SWITCHED_CODE=4902; error 4902 means that the chain does not exist
      if (
        switched === SWITCHED_CODE ||
        !isValidTokenFormat(await getToken())
      ) {
        showToast(ToastType.INFO, "Adding TEN Testnet...");
        const user = await joinTestnet();

        // Store the token in localStorage
        localStorage.setItem("ten_token", "0x" + user);

        const rpcUrls = [
          `${tenGatewayAddress}/${tenGatewayVersion}/?token=${user}`,
        ];
        await addNetworkToMetaMask(rpcUrls);
        showToast(ToastType.SUCCESS, "Added TEN Testnet");
      }

      await finaliseConnectAccounts()

    } catch (error: Error | any) {
      // MetaMask v12.14.* is throwing an unexpected error while adding a network to.
      // Despite the error the network is successfully added.
      // This is a temporary workaround.
      if (error.message === 'l is not a function') {
        await finaliseConnectAccounts()
      } else {
        showToast(ToastType.DESTRUCTIVE, `${error?.message}`);
        throw error;
      }

    } finally {
      setLoading(false);
    }
  };

  const getTestnetStatus = async () => {
    try {
      return await fetchTestnetStatus();
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to connect to TEN Testnet");
      throw error;
    }
  };

  return {
    connectToTenTestnet,
    getTestnetStatus,
  };
};

export default useGatewayService;
