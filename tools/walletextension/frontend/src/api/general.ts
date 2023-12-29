import { ResponseDataInterface } from "@/types/interfaces";
import { httpRequest } from ".";
import { pathToUrl } from "@/routes/router";
import { apiRoutes } from "@/routes";

export const fetchTestnetStatus = async (): Promise<
  ResponseDataInterface<any>
> => {
  return await httpRequest<ResponseDataInterface<any>>({
    method: "post",
    url: pathToUrl(apiRoutes.getHealthStatus),
    data: { jsonrpc: "2.0", method: "obscuro_health", params: [], id: 1 },
  });
};
