import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";
import {
  Batch,
  BatchDetails,
  BatchResponse,
} from "@/src/types/interfaces/BatchInterfaces";

export const getBatches = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<BatchResponse>> => {
  const data = await httpRequest<ResponseDataInterface<BatchResponse>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatches),
    searchParams: payload,
  });
  return data;
};

export const getLatestBatch = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<Batch>> => {
  const data = await httpRequest<ResponseDataInterface<Batch>>({
    method: "get",
    url: pathToUrl(apiRoutes.getLatestBatch),
    searchParams: payload,
  });
  return data;
};

export const getBatchByHash = async (
  hash: string
): Promise<ResponseDataInterface<BatchDetails>> => {
  const data = await httpRequest<ResponseDataInterface<BatchDetails>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatchByHash, { hash }),
  });
  return data;
};
