import {
  addNetworkToMetaMask,
  joinTestnet,
  switchToTenNetwork,
} from "@/api/gateway";
import { useWalletConnection } from "@/components/providers/wallet-provider";
import { useToast } from "@/components/ui/use-toast";
import { SWITCHED_CODE, tenGatewayVersion } from "@/lib/constants";
import { getRPCFromUrl, isTenChain, isValidUserIDFormat } from "@/lib/utils";
import { requestMethods } from "@/routes";

const useGatewayService = () => {
  const { toast } = useToast();
  const { provider } = useWalletConnection();
  const { userID, setUserID, getAccounts } = useWalletConnection();

  const connectAccounts = async () => {
    try {
      await (window as any).ethereum.request({
        method: requestMethods.connectAccounts,
      });
      toast({ variant: "success", description: "Connected to Ten Network" });
    } catch (error) {
      toast({
        variant: "destructive",
        description: "Unable to connect to Ten Network",
      });
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
      toast({ variant: "destructive", description: "Unable to get accounts" });
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
        const rpcUrls = [`${getRPCFromUrl()}/${tenGatewayVersion}/?u=${user}`];
        await addNetworkToMetaMask(rpcUrls);
      }

      if (!(await isMetamaskConnected())) {
        toast({
          variant: "info",
          description: "No accounts found, connecting...",
        });
        await connectAccounts();
      }

      if (!provider) {
        return;
      }
      await getAccounts();
    } catch (error: any) {
      console.error("Error:", error.message);
      toast({
        variant: "destructive",
        description: `${error.message}`,
      });
    }
  };

  return {
    connectToTenTestnet,
  };
};

export default useGatewayService;
