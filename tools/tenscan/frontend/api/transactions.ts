import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";
import {
  TransactionCount,
  Price,
  TransactionResponse,
  Transaction,
} from "@/src/types/interfaces/TransactionInterfaces";

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
