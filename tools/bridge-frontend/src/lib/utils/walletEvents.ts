import { ethers } from "ethers";
import { ethereum } from ".";

export const setupEventListeners = (setAddress: (address: string) => void) => {
  const handleAccountsChange = (accounts: string[]) => {
    setAddress(accounts[0]);
    localStorage.setItem("walletAddress", accounts[0]);
  };

  const handleChainChange = () => {
    console.log("Chain changed; reloading the page");
    window.location.reload();
  };

  ethereum.on("accountsChanged", handleAccountsChange);
  ethereum.on("chainChanged", handleChainChange);

  return () => {
    ethereum.removeListener("accountsChanged", handleAccountsChange);
    ethereum.removeListener("chainChanged", handleChainChange);
  };
};

export const initializeSigner = (provider: ethers.providers.Web3Provider) => {
  if (!provider) {
    return null;
  }

  return provider.getSigner();
};
