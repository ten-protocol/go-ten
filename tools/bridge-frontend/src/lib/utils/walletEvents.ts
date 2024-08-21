import { ethers } from "ethers";

export const setupEventListeners = (
  provider: any,
  setAddress: (address: string) => void
) => {
  const handleAccountsChange = (accounts: string[]) => {
    setAddress(accounts[0]);
    localStorage.setItem("walletAddress", accounts[0]);
  };

  const handleChainChange = () => {
    window.location.reload();
  };

  provider.on("accountsChanged", handleAccountsChange);
  provider.on("chainChanged", handleChainChange);

  return () => {
    provider.removeListener("accountsChanged", handleAccountsChange);
    provider.removeListener("chainChanged", handleChainChange);
  };
};

export const initializeSigner = (provider: any) => {
  if (!provider) {
    return null;
  }

  return new ethers.providers.Web3Provider(provider).getSigner();
};
