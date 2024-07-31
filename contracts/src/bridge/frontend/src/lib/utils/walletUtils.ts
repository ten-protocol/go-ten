import detectEthereumProvider from "@metamask/detect-provider";

export const getEthereumProvider = async () => {
  const provider = await detectEthereumProvider();
  if (!provider) {
    throw new Error("No Ethereum provider detected");
  }
  return provider;
};

export const handleStorage = {
  save: (key: string, value: string) => localStorage.setItem(key, value),
  get: (key: string) => localStorage.getItem(key),
  remove: (key: string) => localStorage.removeItem(key),
};
