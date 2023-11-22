import { type ClassValue, clsx } from "clsx";
import { formatDistanceToNow } from "date-fns";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatTimeAgo(unixTimestampSeconds: string) {
  const date = new Date(Number(unixTimestampSeconds) * 1000);
  return formatDistanceToNow(date, { addSuffix: true });
}

const gatewayaddress = "https://testnet.obscu.ro";
export const tenscanLink = "https://testnet.tenscan.com";
export const tenGatewayVersion = "v1";
export const pathJoin = gatewayaddress + "/" + tenGatewayVersion + "/join/";
export const pathAuthenticate =
  gatewayaddress + "/" + tenGatewayVersion + "/authenticate/";
export const pathQuery = gatewayaddress + "/" + tenGatewayVersion + "/query/";
export const pathRevoke = gatewayaddress + "/" + tenGatewayVersion + "/revoke/";
export const pathVersion = gatewayaddress + "/" + "version/";
export const tenChainIDDecimal = 443;

export const metamaskPersonalSign = "personal_sign";
export const tenChainIDHex = "0x" + tenChainIDDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'
export const METAMASK_CONNECTION_TIMEOUT = 3000;

export function isValidUserIDFormat(value: string) {
  return typeof value === "string" && value.length === 64;
}

export let tenGatewayAddress = gatewayaddress;

export function getRandomIntAsString(min: number, max: number) {
  min = Math.ceil(min);
  max = Math.floor(max);
  const randomInt = Math.floor(Math.random() * (max - min + 1)) + min;
  return randomInt.toString();
}
export function getNetworkName() {
  switch (tenGatewayAddress) {
    case "https://uat-testnet.obscu.ro":
      return "Ten UAT-Testnet";
    case "https://dev-testnet.obscu.ro":
      return "Ten Dev-Testnet";
    default:
      return "Ten Testnet";
  }
}

export function getRPCFromUrl() {
  // get the correct RPC endpoint for each network
  switch (tenGatewayAddress) {
    // case 'https://testnet.obscu.ro':
    //     return 'https://rpc.sepolia-testnet.obscu.ro'
    case "https://sepolia-testnet.obscu.ro":
      return "https://rpc.sepolia-testnet.obscu.ro";
    case "https://uat-testnet.obscu.ro":
      return "https://rpc.uat-testnet.obscu.ro";
    case "https://dev-testnet.obscu.ro":
      return "https://rpc.dev-testnet.obscu.ro";
    default:
      return tenGatewayAddress;
  }
}

export async function isTenChain() {
  let currentChain = await (window as any).ethereum.request({
    method: "eth_chainId",
  });
  return currentChain === tenChainIDHex;
}
