import { tenGatewayVersion } from "../lib/constants";
import { NavLink } from "../types/interfaces";

export const NavLinks: NavLink[] = [];

export const apiRoutes = {
  join: `/${tenGatewayVersion}/join/`,
  authenticate: `/${tenGatewayVersion}/authenticate/`,
  queryAccountUserID: `/${tenGatewayVersion}/query/`,
  revoke: `/${tenGatewayVersion}/revoke/`,
  version: `/version/`,
};

export const requestMethods = {
  connectAccounts: "eth_requestAccounts",
  switchNetwork: "wallet_switchEthereumChain",
  addNetwork: "wallet_addEthereumChain",
};
