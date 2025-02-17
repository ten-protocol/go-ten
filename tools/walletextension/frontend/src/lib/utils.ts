import { type ClassValue, clsx } from "clsx";
import { formatDistanceToNow } from "date-fns";
import { twMerge } from "tailwind-merge";
import { environment, tenChainIDHex, tokenHexLength } from "./constants";
import { L1Network, L2Network, Environment } from "@/types/interfaces";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatTimeAgo(unixTimestampSeconds: string) {
  const date = new Date(Number(unixTimestampSeconds) * 1000);
  return formatDistanceToNow(date, { addSuffix: true });
}

export function isValidTokenFormat(value: string) {
  return typeof value === "string" && value.length === tokenHexLength;
}

export function getRandomIntAsString(min: number, max: number) {
  min = Math.ceil(min);
  max = Math.floor(max);
  const randomInt = Math.floor(Math.random() * (max - min + 1)) + min;
  return randomInt.toString();
}

export async function isTenChain() {
  if (!ethereum) {
    return false;
  }
  let currentChain = await ethereum.request({
    method: "eth_chainId",
  });
  return currentChain === tenChainIDHex;
}

export const { ethereum } =
  typeof window !== "undefined" ? window : ({} as any);

export const downloadMetaMask = () => {
  window ? window.open("https://metamask.io/download", "_blank") : null;
};

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
    l1Rpc: "https://mainnet.infura.io/v3/",
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
};

export const currentNetwork = networkMappings[environment as Environment];
