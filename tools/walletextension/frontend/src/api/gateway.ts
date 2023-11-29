import { apiRoutes, requestMethods } from "../routes";
import { httpRequest } from ".";
import { pathToUrl } from "../routes/router";
import { getNetworkName } from "../lib/utils";
import {
  tenChainIDHex,
  tenscanLink,
  nativeCurrency,
  tenChainIDDecimal,
} from "../lib/constants";
import { AuthenticationResponse } from "@/types/interfaces/GatewayInterfaces";

const typedData = {
  types: {
    EIP712Domain: [
      { name: "name", type: "string" },
      { name: "version", type: "string" },
      { name: "chainId", type: "uint256" },
    ],
    Authentication: [{ name: "Encryption Token", type: "address" }],
  },
  primaryType: "Authentication",
  domain: {
    name: "Ten",
    version: "1.0",
    chainId: tenChainIDDecimal,
  },
  message: {
    "Encryption Token": "0x",
  },
};

const { ethereum } = typeof window !== "undefined" ? window : ({} as any);

export async function switchToTenNetwork() {
  if (!ethereum) {
    return;
  }
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
): Promise<AuthenticationResponse> {
  return await httpRequest<AuthenticationResponse>({
    method: "get",
    url: pathToUrl(apiRoutes.queryAccountUserID),
    searchParams: {
      token: userID,
      a: account,
    },
  });
}

const getSignature = async (account: string, data: any) => {
  if (!ethereum) {
    return;
  }
  return await ethereum.request({
    method: requestMethods.signTypedData,
    params: [account, JSON.stringify(data)],
  });
};

export async function authenticateAccountWithTenGatewayEIP712(
  userID: string,
  account: string
): Promise<any> {
  if (!userID) {
    return;
  }
  try {
    const isAuthenticated = await accountIsAuthenticated(userID, account);
    if (isAuthenticated.status) {
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
  return await httpRequest({
    method: "post",
    url: pathToUrl(apiRoutes.authenticate),
    data: authenticateFields,
    searchParams: {
      token: userID,
    },
  });
};

export async function revokeAccountsApi(userID: string): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.revoke),
    searchParams: {
      token: userID,
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
  if (!ethereum) {
    return;
  }
  try {
    await ethereum.request({
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
