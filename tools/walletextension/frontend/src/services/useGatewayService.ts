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
    showToast(ToastType.INFO, "Connecting to Ten Network...");
    setLoading(true);
    try {
      if (await isTenChain()) {
        if (!token || !isValidTokenFormat(token)) {
          showToast(
            ToastType.DESTRUCTIVE,
            "Existing Ten network detected in MetaMask. Please remove before hitting begin"
          );
          return;
        }
      }
      showToast(ToastType.INFO, "Switching to Ten Network...");
      const switched = await switchToTenNetwork();
      showToast(ToastType.SUCCESS, `Switched to Ten Network: ${switched}`);
      // SWITCHED_CODE=4902; error 4902 means that the chain does not exist
      if (switched === 4902 || !isValidTokenFormat(await getToken(provider))) {
        showToast(ToastType.INFO, "Joining Ten Network...");
        // if (switched === SWITCHED_CODE || (token && !isValidTokenFormat(token))) {
        const user = await joinTestnet();
        showToast(ToastType.SUCCESS, `Joined Ten Network: ${user}`);
        const rpcUrls = [
          `${tenGatewayAddress}/${tenGatewayVersion}/?token=${user}`,
        ];
        await addNetworkToMetaMask(rpcUrls);
        showToast(ToastType.SUCCESS, "Added Ten Network to MetaMask");
      }

      if (!(await isMetamaskConnected())) {
        showToast(ToastType.INFO, "No accounts found, connecting...");
        await connectAccounts();
        showToast(ToastType.SUCCESS, "Connected to Ten Network");
      }

      await fetchUserAccounts();
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, `${error?.message}`);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  return {
    connectToTenTestnet,
  };
};

export default useGatewayService;
