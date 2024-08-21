import { ToastType } from "@/types/interfaces";
import { joinTestnet } from "../api/gateway";
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
  switchToTenNetwork,
} from "@/api/ethRequests";
import { fetchTestnetStatus } from "@/api/general";
import { useWalletStore } from "@/stores/wallet-store";

const useGatewayService = () => {
  const { token, provider, fetchUserAccounts, setLoading } = useWalletStore();

  const addTenTestnet = async (userToken: string) => {
    const rpcUrls = [
      `${tenGatewayAddress}/${tenGatewayVersion}/?token=${userToken}`,
    ];
    await addNetworkToMetaMask(rpcUrls);
    showToast(ToastType.SUCCESS, "Ten Testnet added to MetaMask.");
  };

  const isMetamaskConnected = async () => {
    try {
      if (!provider) throw new Error("Ethereum provider not found.");
      const accounts = await provider.listAccounts();
      return accounts.length > 0;
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to retrieve MetaMask accounts.");
      throw error;
    }
  };

  const connectToTenTestnet = async () => {
    showToast(ToastType.INFO, "Connecting to Ten Testnet...");
    setLoading(true);

    try {
      if (await isTenChain()) {
        if (!token || !isValidTokenFormat(token)) {
          return showToast(
            ToastType.DESTRUCTIVE,
            "Existing Ten Testnet detected. Please remove and reconnect."
          );
        }
      }

      showToast(ToastType.INFO, "Switching to Ten Testnet...");
      const switched = await switchToTenNetwork();

      if (switched === SWITCHED_CODE) {
        const userToken = await joinTestnet();
        await addTenTestnet(userToken);
      }

      if (!(await isMetamaskConnected())) {
        showToast(ToastType.INFO, "No accounts found. Connecting...");
        await connectAccounts();
        showToast(ToastType.SUCCESS, "Connected to Ten Testnet.");
      }

      await fetchUserAccounts();
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, `Connection failed: ${error.message}`);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const getTestnetStatus = async () => {
    try {
      return await fetchTestnetStatus();
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Failed to fetch Ten Testnet status.");
      throw error;
    }
  };

  return {
    connectToTenTestnet,
    getTestnetStatus,
  };
};

export default useGatewayService;
