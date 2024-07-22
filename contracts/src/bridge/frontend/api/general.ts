import { ObscuroConfig, ResponseDataInterface } from "@/src/types";
import { httpRequest } from "@/api";
import { pathToUrl } from "@/src/routes/router";
import { apiRoutes } from "@/src/routes";

export const fetchTestnetStatus = async (): Promise<
  ResponseDataInterface<boolean>
> => {
  return await httpRequest<ResponseDataInterface<boolean>>({
    method: "get",
    url: pathToUrl(apiRoutes.getHealthStatus),
  });
};

export const fetchObscuroConfig = async (): Promise<
  ResponseDataInterface<ObscuroConfig>
> => {
  const res = await httpRequest<ResponseDataInterface<ObscuroConfig>>({
    method: "post",
    url: apiRoutes.getObscuroConfig,
    data: {
      jsonrpc: "2.0",
      method: "obscuro_config",
      params: [],
      id: 1,
    },
  });
  console.log("🚀 ~ res:", res);
  return res;
};
