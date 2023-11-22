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

  async function connectAccounts() {
    try {
      return await (window as any).ethereum.request({
        method: "eth_requestAccounts",
      });
    } catch (error) {
      console.error("User denied account access:", error);
      toast({ description: `User denied account access: ${error}` });
      return null;
    }
  }

  async function isMetamaskConnected() {
    if (!provider) {
      return false;
    }

    try {
      const accounts = await provider.listAccounts();
      return accounts.length > 0;
    } catch (error) {
      console.log("Unable to get accounts");
    }

    return false;
  }

  const connectToTenTestnet = async () => {
    if (await isTenChain()) {
      if (!userID || !isValidUserIDFormat(userID)) {
        toast({
          description:
            "Existing Ten network detected in MetaMask. Please remove before hitting begin",
        });
      }
    } else {
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
        console.log("No accounts found, connecting...");
        await connectAccounts();
      }
      console.log("Connected to Ten Network");

      if (!provider) {
        return;
      }
      console.log("Getting accounts...");
      const accounts = await provider.listAccounts();
      if (accounts.length === 0) {
        console.log("No accounts found");
        toast({ description: "No MetaMask accounts found." });
        return;
      }
    }
  };

  const handleFetchError = (message: string) => {
    toast({ description: message });
  };

  return {
    connectToTenTestnet,
    handleFetchError,
  };
};

export default useGatewayService;
