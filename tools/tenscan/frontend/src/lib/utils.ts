import { type ClassValue, clsx } from "clsx";
import { formatDistanceToNow } from "date-fns";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatTimeAgo(unixTimestampSeconds: string | number) {
  if (!unixTimestampSeconds) {
    return "Unknown";
  }
  const date = new Date(Number(unixTimestampSeconds) * 1000);
  return formatDistanceToNow(date, { addSuffix: true });
}

export const { ethereum } =
  typeof window !== "undefined" ? window : ({} as any);

export const downloadMetaMask = () => {
  window ? window.open("https://metamask.io/download", "_blank") : null;
};

export const formatNumber = (number: string | number) => {
  const num = Number(number);
  return num.toLocaleString();
};

export const firstItem = <T>(arr: T[], key: keyof T) => {
  if (!arr.length) return null;
  if (!arr[0][key]) return null;
  return arr[0][key];
};

export const lastItem = <T>(arr: T[], key: keyof T) => {
  if (!arr.length) return null;
  if (!arr[0][key]) return null;
  return arr[arr.length - 1][key];
};
