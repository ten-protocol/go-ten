import { apiRoutes, requestMethods } from "../routes";
import { httpRequest } from ".";
import { pathToUrl } from "../routes/router";
import { getNetworkName } from "../lib/utils";
import {
  metamaskPersonalSign,
  tenChainIDHex,
  tenscanLink,
  nativeCurrency,
  typedData,
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

const getSignature = async (account: string, data: any) => {
  const { ethereum } = window as any;
  const signature = await ethereum.request({
    method: metamaskPersonalSign,
    params: [account, JSON.stringify(data)],
  });

  return signature;
};

export async function authenticateAccountWithTenGatewayEIP712(
  userID: string,
  account: string
): Promise<any> {
  try {
    const isAuthenticated = await accountIsAuthenticated(userID, account);
    if (isAuthenticated) {
      return "Account is already authenticated";
    }
    const data = {
      ...typedData,
      message: {
        ...typedData.message,
        "Encryption Token": "0x" + userID,
      },
    };
    const signature = await getSignature(account, data);

    const auth = await authenticateUser(userID, {
      signature,
      address: account,
    });
    return auth;
  } catch (error) {
    throw error;
  }
}

const authenticateUser = async (
  userID: string,
  authenticateFields: {
    signature: string;
    address: string;
  }
) => {
  const authenticateResp = await httpRequest({
    method: "post",
    url: pathToUrl(apiRoutes.authenticate),
    data: authenticateFields,
    searchParams: {
      u: userID,
    },
  });
  return authenticateResp;
};

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
    return true;
  } catch (error) {
    console.error(error);
    return error;
  }
}
