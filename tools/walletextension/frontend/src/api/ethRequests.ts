import {
  nativeCurrency,
  tenChainIDDecimal,
  tenChainIDHex,
  tenscanLink,
} from "@/lib/constants";
import { getNetworkName, getRandomIntAsString, isTenChain } from "@/lib/utils";
import { requestMethods } from "@/routes";
import { ethers } from "ethers";
import { accountIsAuthenticated, authenticateUser } from "./gateway";

const { ethereum } = typeof window !== "undefined" ? window : ({} as any);

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

export const switchToTenNetwork = async () => {
  if (!ethereum) {
    throw new Error("No ethereum object found");
  }
  try {
    await ethereum.request({
      method: requestMethods.switchNetwork,
      params: [{ chainId: tenChainIDHex }],
    });

    return 0;
  } catch (error: any) {
    return error.code;
  }
};

export const connectAccounts = async () => {
  if (!ethereum) {
    throw new Error("No ethereum object found");
  }
  try {
    return await ethereum.request({
      method: requestMethods.connectAccounts,
    });
  } catch (error) {
    console.error(error);
    throw error;
  }
};

export const getSignature = async (account: string, data: any) => {
  if (!ethereum) {
    throw new Error("No ethereum object found");
  }
  return await ethereum.request({
    method: requestMethods.signTypedData,
    params: [account, JSON.stringify(data)],
  });
};

export const getUserID = async (provider: ethers.providers.Web3Provider) => {
  if (!provider) {
    return null;
  }

  try {
    if (await isTenChain()) {
      const id = await provider.send(requestMethods.getStorageAt, [
        "getUserID",
        getRandomIntAsString(0, 1000),
        null,
      ]);
      return id;
    } else {
      return null;
    }
  } catch (e: any) {
    console.error(e);
    throw e;
  }
};

export async function addNetworkToMetaMask(rpcUrls: string[]) {
  if (!ethereum) {
    throw new Error("No ethereum object found");
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
    throw error;
  }
}

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
      return {
        status: true,
        message: "Account already authenticated",
      };
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
