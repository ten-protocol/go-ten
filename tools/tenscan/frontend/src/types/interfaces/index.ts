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
  errors?: string[] | string;
  item: T;
  message: string;
  pagination?: PaginationInterface;
  success: string;
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

export enum BadgeType {
  SUCCESS = "success",
  SECONDARY = "secondary",
  DESTRUCTIVE = "destructive",
  DEFAULT = "default",
  OUTLINE = "outline",
}

export interface DashboardAnalyticsData {
  title: string;
  value: string | number | JSX.Element;
  change?: string;
  icon: any;
}

export enum ItemPosition {
  FIRST = "first",
  LAST = "last",
}
