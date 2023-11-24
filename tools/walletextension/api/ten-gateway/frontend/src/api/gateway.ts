import { apiRoutes, requestMethods } from "../routes";
import { httpRequest } from ".";
import { pathToUrl } from "../routes/router";
import { getNetworkName } from "../lib/utils";
import {
  metamaskPersonalSign,
  tenChainIDHex,
  tenscanLink,
  nativeCurrency,
} from "../lib/constants";

export async function switchToTenNetwork() {
  try {
    await (window as any).ethereum.request({
      method: requestMethods.switchNetwork,
      params: [{ chainId: tenChainIDHex }],
    });

    return 0;
  } catch (switchError: any) {
    return switchError.code;
  }
}

export async function fetchVersion(): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.version),
  });
}

export async function accountIsAuthenticated(
  userID: string,
  account: string
): Promise<boolean> {
  return await httpRequest<boolean>({
    method: "get",
    url: pathToUrl(apiRoutes.queryAccountUserID),
    searchParams: {
      u: userID,
      a: account,
    },
  });
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

  return await httpRequest<string>({
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
}

export async function revokeAccountsApi(userID: string): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.revoke),
    searchParams: {
      u: userID,
    },
  });
}

export async function joinTestnet(): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.join),
  });
}

export async function addNetworkToMetaMask(rpcUrls: string[]) {
  try {
    await (window as any).ethereum.request({
      method: requestMethods.addNetwork,
      params: [
        {
          chainId: tenChainIDHex,
          chainName: getNetworkName(),
          nativeCurrency,
          rpcUrls,
          blockExplorerUrls: [tenscanLink],
        },
      ],
    });
  } catch (error) {
    console.error(error);
    return error;
  }

  return true;
}
