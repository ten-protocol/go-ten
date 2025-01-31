import React from "react";

export type NavLink = {
  label: string;
  href?: string;
  isDropdown?: boolean;
  isExternal?: boolean;
  subNavLinks?: NavLink[];
};

export interface DashboardAnalyticsData {
  title: string;
  value: string | number | JSX.Element;
  change?: string;
  icon: any;
  loading?: boolean;
}
