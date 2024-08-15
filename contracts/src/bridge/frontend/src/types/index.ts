import { ethers } from "ethers";
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
  icon?: any;
};

export enum ToastType {
  INFO = "info",
  SUCCESS = "success",
  WARNING = "warning",
  DESTRUCTIVE = "destructive",
  DEFAULT = "default",
}

export interface SeoProps {
  title: string;
  description: string;
  canonicalUrl: string;
  ogTwitterImage: string;
  ogImageUrl: string;
  ogType: string;
  children: React.ReactNode;
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

export interface WalletConnectionContextType {}

export interface Props {
  children: React.ReactNode;
}

export interface State {
  hasError: boolean;
}

export interface WalletConnectionProviderProps {
  children: React.ReactNode;
}

export enum WalletNetwork {
  L1_MAINNET = "0x1",
  L1_SEPOLIA = "0xaa36a7",
  L2_TEN_TESTNET = "0x1bb",
}

export interface Chain {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
  chainId?: string;
}

export interface Token {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
  address: string;
}

export interface INetworkConfig {
  ManagementContractAddress: string;
  L1StartHash: string;
  MessageBusAddress: string;
  L2MessageBusAddress: string;
  ImportantContracts: {
    L1Bridge: string;
    L1CrossChainMessenger: string;
    L2Bridge: string;
    L2CrossChainMessenger: string;
  };
}

export interface ContractState {
  bridgeContract?: ethers.Contract;
  managementContract?: ethers.Contract;
  messageBusContract?: ethers.Contract;
  wallet?: ethers.Wallet;
  messageBusAddress: string;
}
