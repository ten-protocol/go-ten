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
  userID: string,
  account: string
): Promise<AuthenticationResponse> {
  return await httpRequest<AuthenticationResponse>({
    method: "get",
    url: pathToUrl(apiRoutes.queryAccountUserID),
    searchParams: {
      token: userID,
      a: account,
    },
  });
}

export const authenticateUser = async (
  userID: string,
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
      token: userID,
    },
  });
};

export async function revokeAccountsApi(userID: string): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.revoke),
    searchParams: {
      token: userID,
    },
  });
}

export async function joinTestnet(): Promise<string> {
  return await httpRequest<string>({
    method: "get",
    url: pathToUrl(apiRoutes.join),
  });
}
