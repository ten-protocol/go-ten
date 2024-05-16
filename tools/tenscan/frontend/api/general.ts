import { ResponseDataInterface } from "@/src/types/interfaces";
import { httpRequest } from ".";
import { pathToUrl } from "@/src/routes/router";
import { apiRoutes } from "@/src/routes";

export const fetchTestnetStatus = async (): Promise<
  ResponseDataInterface<boolean>
> => {
  return await httpRequest<ResponseDataInterface<boolean>>({
    method: "get",
    url: pathToUrl(apiRoutes.getHealthStatus),
    data: { jsonrpc: "2.0", method: "obscuro_health", id: 1 },
  });
};
