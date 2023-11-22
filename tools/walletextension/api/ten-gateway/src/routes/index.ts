import { tenGatewayVersion } from "@/lib/constants";
import { NavLink } from "../types/interfaces";

export const NavLinks: NavLink[] = [];

export const apiRoutes = {
  join: `/${tenGatewayVersion}/join/`,
  authenticate: `/${tenGatewayVersion}/authenticate/`,
  queryAccountUserID: `/${tenGatewayVersion}/query/`,
  revoke: `/${tenGatewayVersion}/revoke/`,
  version: `/version/`,
};
