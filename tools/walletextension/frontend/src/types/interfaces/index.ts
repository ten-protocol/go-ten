import React from "react";

export interface SeoProps {
  title: string;
  description: string;
  canonicalUrl: string;
  ogTwitterImage: string;
  ogImageUrl: string;
  ogType: string;
  children: React.ReactNode;
}

export interface ErrorType {
  statusCode?: number;
  showRedirectText?: boolean;
  heading?: string;
  statusText?: string;
  message?: string;
  redirectText?: string;
  customPageTitle?: string;
  isFullWidth?: boolean;
  style?: React.CSSProperties;
  hasGetInitialPropsRun?: boolean;
  err?: Error;
  showMessage?: boolean;
  showStatusText?: boolean;
  isModal?: boolean;
  redirectLink?: string;
  children?: React.ReactNode;
  [key: string]: any;
}

export interface IconProps {
  width?: string;
  height?: string;
  fill?: string;
  stroke?: string;
  strokeWidth?: string;
  className?: string;
  isActive?: boolean;
  onClick?: () => void;
}

export interface GetInfinitePagesInterface<T> {
  nextId?: number;
  previousId?: number;
  data: T;
  count: number;
}

export interface PaginationInterface {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
}

export interface ResponseDataInterface<T> {
  result: T;
  item: T;
  message: string;
  pagination?: PaginationInterface;
  success: boolean;
}

export type NavLink = {
  label: string;
  href?: string;
  isDropdown?: boolean;
  isExternal?: boolean;
  subNavLinks?: NavLink[];
};

export enum ToastType {
  INFO = "info",
  SUCCESS = "success",
  WARNING = "warning",
  DESTRUCTIVE = "destructive",
  DEFAULT = "default",
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

export type Environment =
  | "uat-testnet"
  | "sepolia-testnet"
  | "dev-testnet"
  | "local-testnet";
