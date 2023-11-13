import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";

export const getBlocks = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<any>> => {
  const data = await httpRequest<ResponseDataInterface<any>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBlocks),
    searchParams: payload,
  });
  return data;
};
