import { apiRoutes } from "@/routes";
import { httpRequest } from ".";
import { pathToUrl } from "@/routes/router";
import { getNetworkName } from "@/lib/utils";
import {
  metamaskPersonalSign,
  tenChainIDHex,
  tenscanLink,
} from "@/lib/constants";

export async function switchToTenNetwork() {
  try {
    await (window as any).ethereum.request({
      method: "wallet_switchEthereumChain",
      params: [{ chainId: tenChainIDHex }],
    });

    return 0;
  } catch (switchError: any) {
    return switchError.code;
  }
}

export async function fetchVersion(): Promise<string> {
  const data = await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.version),
  });
  return data;
}

export async function accountIsAuthenticated(
  userID: string,
  account: string
): Promise<boolean> {
  const data = await httpRequest<boolean>({
    method: "get",
    url: pathToUrl(apiRoutes.queryAccountUserID),
    searchParams: {
      u: userID,
      a: account,
    },
  });
  return data;
}

export async function authenticateAccountWithTenGateway(
  userID: string,
  account: string
): Promise<string> {
  const textToSign = `Register ${userID} for ${account.toLowerCase()}`;
  const signature = await (window as any).ethereum
    .request({
      method: metamaskPersonalSign,
      params: [textToSign, account],
    })
    .catch((error: any) => -1);

  if (signature === -1) {
    return "Signing failed";
  }

  const data = await httpRequest<string>({
    method: "post",
    url: pathToUrl(apiRoutes.authenticate),
    data: {
      signature,
      message: textToSign,
    },
    searchParams: {
      u: userID,
    },
  });
  return data;
}

export async function revokeAccountsApi(userID: string): Promise<void> {
  const data = await httpRequest<void>({
    method: "get",
    url: pathToUrl(apiRoutes.revoke),
    searchParams: {
      u: userID,
    },
  });
}

export async function joinTestnet(): Promise<string> {
  const data = await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.join),
  });
  return data;
}

export async function addNetworkToMetaMask(rpcUrls: string[]) {
  try {
    await (window as any).ethereum.request({
      method: "wallet_addEthereumChain",
      params: [
        {
          chainId: tenChainIDHex,
          chainName: getNetworkName(),
          nativeCurrency: {
            name: "Sepolia Ether",
            symbol: "ETH",
            decimals: 18,
          },
          rpcUrls,
          blockExplorerUrls: [tenscanLink],
        },
      ],
    });
  } catch (error) {
    console.error(error);
    return false;
  }

  return true;
}
