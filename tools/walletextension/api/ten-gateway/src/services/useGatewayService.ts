import {
  addNetworkToMetaMask,
  joinTestnet,
  switchToTenNetwork,
} from "@/api/gateway";
import { useWalletConnection } from "@/components/providers/wallet-provider";
import { useToast } from "@/components/ui/use-toast";
import { SWITCHED_CODE, tenGatewayVersion } from "@/lib/constants";
import { getRPCFromUrl, isTenChain, isValidUserIDFormat } from "@/lib/utils";

const useGatewayService = () => {
  const { toast } = useToast();
  const { provider } = useWalletConnection();
  const { userID, setUserID } = useWalletConnection();

  const connectAccounts = async () => {
    try {
      return await (window as any).ethereum.request({
        method: "eth_requestAccounts",
      });
    } catch (error) {
      toast({ description: `User denied account access: ${error}` });
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
      toast({ description: "Unable to get accounts" });
    }
    return false;
  };

  const connectToTenTestnet = async () => {
    console.log(
      "🚀 ~ file: useGatewayService.ts:43 ~ connectToTenTestnet ~ userID:",
      userID
    );
    if (await isTenChain()) {
      if (!userID || !isValidUserIDFormat(userID)) {
        return toast({
          description:
            "Existing Ten network detected in MetaMask. Please remove before hitting begin",
        });
      }
    }

    const switched = await switchToTenNetwork();
    const rpcUrls = [`${getRPCFromUrl()}/${tenGatewayVersion}/?u=${userID}`];

    if (
      switched === SWITCHED_CODE ||
      (userID && !isValidUserIDFormat(userID))
    ) {
      const user = await joinTestnet();
      setUserID(user);
      await addNetworkToMetaMask(rpcUrls);
    }

    if (!(await isMetamaskConnected())) {
      toast({ description: "No accounts found, connecting..." });
      await connectAccounts();
    }
    toast({ description: "Connected to Ten Network" });
    if (!provider) {
      return;
    }
    toast({ description: "Getting accounts..." });
    const accounts = await provider.listAccounts();
    if (accounts.length === 0) {
      toast({ description: "No accounts found" });
      toast({ description: "No MetaMask accounts found." });
      return;
    }
  };

  return {
    connectToTenTestnet,
  };
};

export default useGatewayService;
