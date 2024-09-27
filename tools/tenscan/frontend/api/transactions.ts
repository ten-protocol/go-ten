import { jsonHexToObj } from "@repo/ui/lib/utils";
import { httpRequest } from ".";
import { apiRoutes, ethMethods, tenCustomQueryMethods } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import {
  TransactionCount,
  Price,
  TransactionResponse,
  Transaction,
} from "@/src/types/interfaces/TransactionInterfaces";
import { showToast } from "@repo/ui/shared/use-toast";
import { ResponseDataInterface } from "@repo/ui/lib/types/common";
import { ToastType } from "@repo/ui/lib/enums/toast";

export const fetchTransactions = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<TransactionResponse>> => {
  return await httpRequest<ResponseDataInterface<TransactionResponse>>({
    method: "get",
    url: pathToUrl(apiRoutes.getTransactions),
    searchParams: payload,
  });
};

export const fetchTransactionCount = async (): Promise<TransactionCount> => {
  return await httpRequest<TransactionCount>({
    method: "get",
    url: pathToUrl(apiRoutes.getTransactionCount),
  });
};

export const fetchEtherPrice = async (): Promise<Price> => {
  return await httpRequest<Price>({
    method: "get",
    url: apiRoutes.getEtherPrice,
  });
};

export const fetchTransactionByHash = async (
  hash: string
): Promise<ResponseDataInterface<Transaction>> => {
  return await httpRequest<ResponseDataInterface<Transaction>>({
    method: "get",
    url: pathToUrl(apiRoutes.getTransactionByHash, { hash }),
  });
};

export const personalTransactionsData = async (
  provider: any,
  walletAddress: string | null,
  options: Record<string, any>
) => {
  try {
    if (provider && walletAddress) {
      const requestPayload = {
        address: walletAddress,
        pagination: {
          ...options,
        },
      };
      const personalTxResp = await provider.send(ethMethods.getStorageAt, [
        tenCustomQueryMethods.listPersonalTransactions,
        JSON.stringify(requestPayload),
        null,
      ]);
      const personalTxData = jsonHexToObj(personalTxResp);
      return personalTxData;
    }

    return null;
  } catch (error) {
    console.error("Error fetching personal transactions:", error);
    showToast(ToastType.DESTRUCTIVE, "Error fetching personal transactions");
    throw error;
  }
};

export const fetchPersonalTxnByHash = async (
  provider: any,
  hash: string
): Promise<any> => {
  try {
    if (provider) {
      const personalTxnResp = await provider.send(
        ethMethods.getTransactionReceipt,
        [hash]
      );
      return personalTxnResp;
    }
  } catch (error) {
    console.error("Error fetching personal transaction:", error);
    showToast(ToastType.DESTRUCTIVE, "Error fetching personal transaction");
    throw error;
  }
};
