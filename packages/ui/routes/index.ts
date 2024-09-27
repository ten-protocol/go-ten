export const apiRoutes = {
  getHealthStatus: `/network-health/`,
};

export const requestMethods = {
  requestAccounts: "eth_requestAccounts",
  switchNetwork: "wallet_switchEthereumChain",
  addNetwork: "wallet_addEthereumChain",
  getStorageAt: "eth_getStorageAt",
  signTypedData: "eth_signTypedData_v4",
  getChainId: "eth_chainId",
};
