import { ethers } from "ethers";
import { ethereum } from ".";

export const getEthereumProvider = async () => {
  const provider = new ethers.providers.Web3Provider(ethereum);
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
