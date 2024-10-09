import { destr } from "destr";
import { is } from "ramda";

type AuthCredentials = {
  id: string;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  tokenType: string;
};

export function getAuthData(): AuthCredentials | undefined {
  const auth_data = destr(localStorage.getItem("authentication"));

  if (!is(Object, auth_data)) return undefined;

  return auth_data;
}
