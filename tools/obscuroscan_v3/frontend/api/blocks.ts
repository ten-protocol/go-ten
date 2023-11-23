import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";

export const fetchBlocks = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<any>> => {
  return await httpRequest<ResponseDataInterface<any>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBlocks),
    searchParams: payload,
  });
};
