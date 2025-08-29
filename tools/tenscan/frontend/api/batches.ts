import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import {
  Batch,
  BatchResponse,
  LatestBatch,
} from "@/src/types/interfaces/BatchInterfaces";
import { ResponseDataInterface } from "@repo/ui/lib/types/common";

export const fetchBatches = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<BatchResponse>> => {
  return await httpRequest<ResponseDataInterface<BatchResponse>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatches),
    searchParams: payload,
  });
};

export const fetchLatestBatch = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<LatestBatch>> => {
  return await httpRequest<ResponseDataInterface<LatestBatch>>({
    method: "get",
    url: pathToUrl(apiRoutes.getLatestBatch),
    searchParams: payload,
  });
};

export const fetchBatchByHash = async (
  hash: string
): Promise<ResponseDataInterface<Batch>> => {
  return await httpRequest<ResponseDataInterface<Batch>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatchByHash, { hash }),
  });
};

export const fetchBatchByHeight = async (
  height: string
): Promise<ResponseDataInterface<Batch>> => {
  return await httpRequest<ResponseDataInterface<Batch>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatchByHeight, { height }),
  });
};

export const fetchBatchTransactions = async (
  fullHash: string,
  options?: Record<string, any>
): Promise<ResponseDataInterface<any>> => {
  return await httpRequest<ResponseDataInterface<any>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatchTransactions, { fullHash }),
    searchParams: options,
  });
};

export const fetchBatchBySequence = async (
  sequence: string
): Promise<ResponseDataInterface<Batch>> => {
  return await httpRequest<ResponseDataInterface<Batch>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatchBySequence, { seq: sequence }),
  });
};
