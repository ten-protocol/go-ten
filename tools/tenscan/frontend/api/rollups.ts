import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";

export const fetchRollups = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<any>> => {
  return await httpRequest<ResponseDataInterface<any>>({
    method: "get",
    url: pathToUrl(apiRoutes.getLatestRollup),
    searchParams: payload,
  });
};

export const decryptEncryptedRollup = async ({
  StrData,
}: {
  StrData: string;
}): Promise<ResponseDataInterface<any>> => {
  return await httpRequest<ResponseDataInterface<any>>({
    method: "post",
    url: pathToUrl(apiRoutes.decryptEncryptedRollup),
    data: { StrData },
  });
};
