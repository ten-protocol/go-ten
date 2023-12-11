import { ToastType } from "@/types/interfaces";
import { joinTestnet } from "../api/gateway";
import { useWalletConnection } from "../components/providers/wallet-provider";
import { showToast } from "../components/ui/use-toast";
import {SWITCHED_CODE, tenGatewayAddress, tenGatewayVersion} from "../lib/constants";
import { isTenChain, isValidTokenFormat } from "../lib/utils";
import {
  addNetworkToMetaMask,
  connectAccounts,
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

      const switched = await switchToTenNetwork();

      if (switched === SWITCHED_CODE || (token && !isValidTokenFormat(token))) {
        const user = await joinTestnet();
        const rpcUrls = [
          `${tenGatewayAddress}/${tenGatewayVersion}/?token=${user}`,
        ];
        await addNetworkToMetaMask(rpcUrls);
      }

      if (!(await isMetamaskConnected())) {
        showToast(ToastType.INFO, "No accounts found, connecting...");
        await connectAccounts();
        showToast(ToastType.SUCCESS, "Connected to Ten Network");
      }

      await fetchUserAccounts();
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, `${error.message}`);
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
