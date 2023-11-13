import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";
import { ContractCount } from "@/src/types/interfaces/ContractInterface";

export const getContractCount = async (
  payload?: Record<string, any>
): Promise<ContractCount> => {
  const data = await httpRequest<ContractCount>({
    method: "get",
    url: pathToUrl(apiRoutes.getContractCount),
    searchParams: payload,
  });
  return data;
};

export const getVerifiedContracts = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<any>> => {
  const data = await httpRequest<ResponseDataInterface<any>>({
    method: "get",
    url: pathToUrl(apiRoutes.getVerifiedContracts),
    searchParams: payload,
  });
  return data;
};
