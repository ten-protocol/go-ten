import { RouteIcon } from "lucide-react";
import { NavLink } from "@/src/types";
import { ReaderIcon } from "@radix-ui/react-icons";

export const NavLinks: NavLink[] = [
  {
    href: "/",
    label: "Bridge",
    isExternal: false,
    isDropdown: false,
    icon: RouteIcon,
  },
  {
    href: "/transactions",
    label: "Transactions",
    isExternal: false,
    isDropdown: false,
    icon: ReaderIcon,
  },
];

export const apiRoutes = {
  getHealthStatus: `/network-health/`,
  getNetworkConfig: `/network-config/`,
};

export const requestMethods = {
  requestAccounts: "eth_requestAccounts",
  switchNetwork: "wallet_switchEthereumChain",
  addNetwork: "wallet_addEthereumChain",
  getStorageAt: "eth_getStorageAt",
  signTypedData: "eth_signTypedData_v4",
  getChainId: "eth_chainId",
};
