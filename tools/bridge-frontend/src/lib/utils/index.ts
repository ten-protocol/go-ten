import { Environment, ItemPosition, L1Network, L2Network } from "@/src/types";
import { type ClassValue, clsx } from "clsx";
import { formatDistanceToNow } from "date-fns";
import { twMerge } from "tailwind-merge";
import { environment } from "../constants";

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

export function trackEvent(eventName: string, eventData: object) {
  // @ts-ignore
  if (!window.gtag) {
    return;
  }
  // @ts-ignore
  window.gtag("event", eventName, eventData);
}

export const networkMappings = {
  "uat-testnet": {
    l1: L1Network.UAT,
    l2: L2Network.UAT,
    l1Rpc: "https://uat-testnet-eth2network.ten.xyz",
    l1Explorer: "https://uat-testnet-tenscan.io",
    l2Gateway: "https://uat-testnet.ten.xyz",
  },
  "sepolia-testnet": {
    l1: L1Network.SEPOLIA,
    l2: L2Network.SEPOLIA,
    l1Rpc: "https://rpc.sepolia.org",
    l1Explorer: "https://sepolia.etherscan.io",
    l2Gateway: "https://sepolia-testnet.ten.xyz",
  },
  "dev-testnet": {
    l1: L1Network.DEV,
    l2: L2Network.DEV,
    l1Rpc: "https://dev-testnet-eth2network.ten.xyz",
    l1Explorer: "https://dev-testnet-tenscan.io",
    l2Gateway: "https://dev-testnet.ten.xyz",
  },
  "local-testnet": {
    l1: L1Network.LOCAL,
    l2: L2Network.LOCAL,
    l1Rpc: `${process.env.NEXT_PUBLIC_L1NodeHostHTTP}:${process.env.NEXT_PUBLIC_L1NodePortHTTP}`,
    l1Explorer: `${process.env.NEXT_PUBLIC_TENSCAN_URL}`,
    l2Gateway: `${process.env.NEXT_PUBLIC_GATEWAY_URL}`,
  },
  "mainnet": {
    l1: L1Network.MAINNET,
    l2: L2Network.MAINNET,
    l1Rpc: `${process.env.NEXT_PUBLIC_L1NodeHostHTTP}:${process.env.NEXT_PUBLIC_L1NodePortHTTP}`,
    l1Explorer: `${process.env.NEXT_PUBLIC_TENSCAN_URL}`,
    l2Gateway: `${process.env.NEXT_PUBLIC_GATEWAY_URL}`,
  },
};

export const currentNetwork = networkMappings[environment as Environment];

export const formatNumber = (number: string | number) => {
  const num = Number(number);
  return num.toLocaleString();
};

export const getItem = <T>(
  arr: T[],
  key: string,
  position: ItemPosition = ItemPosition.FIRST,
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

export const handleStorage = {
  save: (key: string, value: string) => localStorage.setItem(key, value),
  get: (key: string) => localStorage.getItem(key),
  remove: (key: string) => localStorage.removeItem(key),
};
