export type ButtonVariants =
  | "outline"
  | "link"
  | "default"
  | "destructive"
  | "secondary"
  | "ghost"
  | "clear";

export type Environment = "uat-testnet" | "sepolia-testnet" | "dev-testnet";

export type NavLink = {
  label: string;
  href?: string;
  isDropdown?: boolean;
  isExternal?: boolean;
  subNavLinks?: NavLink[];
  icon?: React.ElementType;
};
