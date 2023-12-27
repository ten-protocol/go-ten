import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";
import { ContractCount } from "@/src/types/interfaces/ContractInterface";

export const fetchContractCount = async (
  payload?: Record<string, any>
): Promise<ContractCount> => {
  return await httpRequest<ContractCount>({
    method: "get",
    url: pathToUrl(apiRoutes.getContractCount),
    searchParams: payload,
  });
};

export const fetchVerifiedContracts = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<any>> => {
  return await httpRequest<ResponseDataInterface<any>>({
    method: "get",
    url: pathToUrl(apiRoutes.getVerifiedContracts),
    searchParams: payload,
  });
};
