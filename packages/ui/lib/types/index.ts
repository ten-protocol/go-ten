export interface ResponseDataInterface<T> {
  result: T;
  errors?: string[] | string;
  item: T;
  message: string;
  pagination?: PaginationInterface;
  success: string;
}

export interface PaginationInterface {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
}

export type ButtonVariants =
  | "outline"
  | "link"
  | "default"
  | "destructive"
  | "secondary"
  | "ghost"
  | "clear";

export enum ToastType {
  INFO = "info",
  SUCCESS = "success",
  WARNING = "warning",
  DESTRUCTIVE = "destructive",
  DEFAULT = "default",
}

export enum BadgeType {
  SUCCESS = "success",
  SECONDARY = "secondary",
  DESTRUCTIVE = "destructive",
  DEFAULT = "default",
  OUTLINE = "outline",
}

export enum IErrorMessages {
  UnknownAccount = "unknown account",
  InsufficientFunds = "insufficient funds",
  UserDeniedTransactionSignature = "User denied transaction signature",
  UserRejectedTheRequest = "User rejected the request",
  ExecutionReverted = "execution reverted",
  RateLimitExceeded = "rate limit exceeded",
  WithdrawalSpent = "withdrawal already spent",
}

export enum ItemPosition {
  FIRST = "first",
  LAST = "last",
}

export enum L1Network {
  MAINNET = 1,
  SEPOLIA = 11155111,
  UAT = 1337,
  DEV = 1337,
  LOCAL = 1337,
}

export enum L2Network {
  SEPOLIA = 443,
  UAT = 443,
  DEV = 443,
  LOCAL = 443,
}

export type Environment = "uat-testnet" | "sepolia-testnet" | "dev-testnet";
