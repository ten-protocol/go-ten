import detectEthereumProvider from "@metamask/detect-provider";
import { ToastType, WalletNetwork } from "../../types";
import { requestMethods } from "../../routes";
import { toast } from "@/src/components/ui/use-toast";

export const getEthereumProvider = async () => {
  const provider = await detectEthereumProvider();
  console.log("🚀 ~ getEthereumProvider ~ provider:", provider);
  if (!provider) {
    throw new Error("No Ethereum provider detected");
  }
  return provider;
};

export const switchNetwork = async (
  provider: any,
  desiredNetwork: WalletNetwork
) => {
  if (!provider) {
    toast({
      title: "Error",
      description: "Please connect to wallet first",
      variant: ToastType.DESTRUCTIVE,
    });
    return;
  }
  await provider.request({
    method: requestMethods.switchNetwork,
    params: [{ chainId: desiredNetwork }],
  });
};

export const handleStorage = {
  save: (key: string, value: string) => localStorage.setItem(key, value),
  get: (key: string) => localStorage.getItem(key),
  remove: (key: string) => localStorage.removeItem(key),
};
