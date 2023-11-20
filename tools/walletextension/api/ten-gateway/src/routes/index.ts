import { NavLink } from "../types/interfaces";

export const NavLinks: NavLink[] = [
  {
    href: "/",
    label: "Home",
    isExternal: false,
    isDropdown: false,
  },
  {
    href: "/personal",
    label: "Personal",
    isExternal: false,
    isDropdown: false,
  },
  {
    label: "Blockchain",
    isExternal: false,
    isDropdown: true,
    subNavLinks: [
      {
        href: "/transactions",
        label: "Transactions",
        isExternal: false,
      },
      {
        href: "/blocks",
        label: "Blocks",
        isExternal: false,
      },
      {
        href: "/batches",
        label: "Batches",
        isExternal: false,
      },
    ],
  },
  {
    label: "Resources",
    isExternal: false,
    isDropdown: true,
    subNavLinks: [
      {
        href: "/resources/decrypt",
        label: "Decrypt",
        isExternal: false,
      },
      {
        href: "/resources/verified-data",
        label: "Verified Data",
        isExternal: false,
      },
    ],
  },
];
