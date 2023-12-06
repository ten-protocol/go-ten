import { apiRoutes } from "../routes";
import { httpRequest } from ".";
import { pathToUrl } from "../routes/router";
import { AuthenticationResponse } from "@/types/interfaces/GatewayInterfaces";

export async function fetchVersion(): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.version),
  });
}

export async function accountIsAuthenticated(
  token: string,
  account: string
): Promise<AuthenticationResponse> {
  return await httpRequest<AuthenticationResponse>({
    method: "get",
    url: pathToUrl(apiRoutes.queryAccountToken),
    searchParams: {
      token,
      a: account,
    },
  });
}

export const authenticateUser = async (
  token: string,
  authenticateFields: {
    signature: string;
    address: string;
  }
) => {
  return await httpRequest({
    method: "post",
    url: pathToUrl(apiRoutes.authenticate),
    data: authenticateFields,
    searchParams: {
      token,
    },
  });
};

export async function revokeAccountsApi(token: string): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.revoke),
    searchParams: {
      token,
    },
  });
}

export async function joinTestnet(): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.join),
  });
}
