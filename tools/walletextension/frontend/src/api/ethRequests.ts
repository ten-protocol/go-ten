import {
  nativeCurrency,
  tenChainIDDecimal,
  tenChainIDHex,
  tenscanLink,
  userStorageAddress,
} from "@/lib/constants";
import {
  getNetworkName,
  getRandomIntAsString,
  isTenChain,
  ethereum,
} from "@/lib/utils";
import { requestMethods } from "@/routes";
import { ethers } from "ethers";
import { accountIsAuthenticated, authenticateUser } from "./gateway";
import { showToast } from "@/components/ui/use-toast";
import { ToastType } from "@/types/interfaces";

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
    throw "No ethereum object found";
  }
  try {
    await ethereum.request({
      method: requestMethods.switchNetwork,
      params: [{ chainId: tenChainIDHex }],
    });

    return 0;
  } catch (switchError: any) {
    showToast(ToastType.DESTRUCTIVE, `switchToTenNetwork: ${switchError.code}`);
    return switchError.code;
  }
};

export const connectAccounts = async () => {
  if (!ethereum) {
    throw "No ethereum object found";
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
    throw "No ethereum object found";
  }
  try {
    return await ethereum.request({
      method: requestMethods.signTypedData,
      params: [account, JSON.stringify(data)],
    });
  } catch (error) {
    console.error(error);
    throw "Failed to get signature";
  }
};

export const getToken = async (provider: ethers.providers.Web3Provider) => {
  if (!provider.send) {
    return null;
  }
  try {
    if (await isTenChain()) {
      const token = await provider.send(requestMethods.getStorageAt, [
        userStorageAddress,
        getRandomIntAsString(0, 1000),
        null,
      ]);
      return token;
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
    throw "No ethereum object found";
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
  token: string,
  account: string
): Promise<any> {
  if (!token) {
    return showToast(
      ToastType.INFO,
      "Encryption token not found. Please try again later."
    );
  }

  try {
    const isAuthenticated = await accountIsAuthenticated(token, account);
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
        "Encryption Token": token,
      },
    };
    const signature = await getSignature(account, data);

    return await authenticateUser(token, {
      signature,
      address: account,
    });
  } catch (error: any) {
    throw error;
  }
}
