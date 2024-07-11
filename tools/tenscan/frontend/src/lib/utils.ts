import { type ClassValue, clsx } from "clsx";
import { formatDistanceToNow } from "date-fns";
import { twMerge } from "tailwind-merge";
import { ItemPosition } from "../types/interfaces";

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

export const getItem = <T>(
  arr: T[],
  key: string,
  position: ItemPosition = ItemPosition.FIRST
) => {
  if (!arr || !arr.length) {
    return null;
  }

  const keys = key.split(".");
  const item = position === ItemPosition.FIRST ? arr[0] : arr[arr.length - 1];
  let value: any = item;

  for (const k of keys) {
    if (value[k] === undefined) {
      return null;
    }
    value = value[k];
  }

  return value;
};
