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

  const connectToTenTestnet = async () => {
    showToast(ToastType.INFO, "Connecting to Obscuro Testnet...");
    setLoading(true);
    try {
      if (await isTenChain()) {
        if (!token || !isValidTokenFormat(token)) {
          showToast(
            ToastType.DESTRUCTIVE,
            "Existing Obscuro Testnet detected in MetaMask. Please remove before hitting begin"
          );
          return;
        }
      }
      showToast(ToastType.INFO, "Switching to Obscuro Testnet...");
      const switched = await switchToTenNetwork();
      showToast(ToastType.SUCCESS, `Switched to Obscuro Testnet: ${switched}`);
      // SWITCHED_CODE=4902; error 4902 means that the chain does not exist
      if (
        switched === SWITCHED_CODE ||
        !isValidTokenFormat(await getToken(provider))
      ) {
        showToast(ToastType.INFO, "Adding Obscuro Testnet...");
        const user = await joinTestnet();
        const rpcUrls = [
          `${tenGatewayAddress}/${tenGatewayVersion}/?token=${user}`,
        ];
        await addNetworkToMetaMask(rpcUrls);
        showToast(ToastType.SUCCESS, "Added Obscuro Testnet");
      }

      if (!(await isMetamaskConnected())) {
        showToast(ToastType.INFO, "No accounts found, connecting...");
        await connectAccounts();
        showToast(ToastType.SUCCESS, "Connected to Obscuro Testnet");
      }
      await fetchUserAccounts();
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, `${error?.message}`);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const getTestnetStatus = async () => {
    try {
      return await fetchTestnetStatus();
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to connect to Obscuro Testnet");
      throw error;
    }
  };

  return {
    connectToTenTestnet,
    getTestnetStatus,
  };
};

export default useGatewayService;
