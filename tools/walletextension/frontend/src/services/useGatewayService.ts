import { ToastType } from "@/types/interfaces";
import {
  addNetworkToMetaMask,
  joinTestnet,
  switchToTenNetwork,
} from "../api/gateway";
import { useWalletConnection } from "../components/providers/wallet-provider";
import { showToast } from "../components/ui/use-toast";
import { SWITCHED_CODE, tenGatewayVersion } from "../lib/constants";
import { getRPCFromUrl, isTenChain, isValidUserIDFormat } from "../lib/utils";
import { requestMethods } from "../routes";

const { ethereum } = typeof window !== "undefined" ? window : ({} as any);

const useGatewayService = () => {
  const { provider } = useWalletConnection();
  const { userID, setUserID, getAccounts } = useWalletConnection();

  const connectAccounts = async () => {
    if (!ethereum) {
      return null;
    }
    try {
      await ethereum.request({
        method: requestMethods.connectAccounts,
      });
      showToast(ToastType.SUCCESS, "Connected to Ten Network");
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to connect to Ten Network");
      return null;
    }
  };

  const isMetamaskConnected = async () => {
    if (!provider) {
      return false;
    }
    try {
      const accounts = await provider.listAccounts();
      return accounts.length > 0;
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Unable to get accounts");
    }
    return false;
  };

  const connectToTenTestnet = async () => {
    try {
      if (await isTenChain()) {
        if (!userID || !isValidUserIDFormat(userID)) {
          throw new Error(
            "Existing Ten network detected in MetaMask. Please remove before hitting begin"
          );
        }
      }

      const switched = await switchToTenNetwork();

      if (
        switched === SWITCHED_CODE ||
        (userID && !isValidUserIDFormat(userID))
      ) {
        const user = await joinTestnet();
        setUserID(user);
        const rpcUrls = [
          `${getRPCFromUrl()}/${tenGatewayVersion}/?token=${user}`,
        ];
        await addNetworkToMetaMask(rpcUrls);
      }

      if (!(await isMetamaskConnected())) {
        showToast(ToastType.INFO, "No accounts found, connecting...");
        await connectAccounts();
      }

      if (!provider || !userID) {
        return;
      }
      await getAccounts(provider, userID);
    } catch (error: any) {
      showToast(ToastType.DESTRUCTIVE, `${error.message}`);
    }
  };

  return {
    connectToTenTestnet,
  };
};

export default useGatewayService;
