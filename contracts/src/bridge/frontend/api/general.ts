import { ResponseDataInterface } from "@/src/types";
import { httpRequest } from "@/api";
import { pathToUrl } from "@/src/routes/router";
import { apiRoutes } from "@/src/routes";

export const fetchTestnetStatus = async (): Promise<
  ResponseDataInterface<boolean>
> => {
  return await httpRequest<ResponseDataInterface<boolean>>({
    method: "post",
    url: pathToUrl(apiRoutes.getHealthStatus),
    data: { jsonrpc: "2.0", method: "obscuro_health", params: [], id: 1 },
  });
};
