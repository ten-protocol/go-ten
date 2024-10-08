import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { BlockResponse } from "@/src/types/interfaces/BlockInterfaces";
import { ResponseDataInterface } from "@repo/ui/lib/types/common";

export const fetchBlocks = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<BlockResponse>> => {
  return await httpRequest<ResponseDataInterface<BlockResponse>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBlocks),
    searchParams: payload,
  });
};
