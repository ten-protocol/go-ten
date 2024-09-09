export const socialLinks = {
  github: "https://github.com/ten-protocol",
  discord: "https://discord.gg/tenprotocol",
  twitter: "https://twitter.com/tenprotocol",
  twitterHandle: "@tenprotocol",
};

export const pollingInterval = 5000;
export const maxRetries = 3;
export const pricePollingInterval = 60 * 1000; // 1 minute in milliseconds

export const RESET_COPIED_TIMEOUT = 2000;

export const getOptions = (query: {
  page?: string | string[];
  size?: string | string[];
}) => {
  const offset =
    query.page && query.size
      ? (parseInt(query.page as string, 10) - 1) *
        parseInt(query.size as string, 10)
      : 0;
  const options = {
    offset: Number.isNaN(offset) ? 0 : offset,
    size: Number.isNaN(parseInt(query.size as string, 10))
      ? 10
      : parseInt(query.size as string, 10),
    // sort: query.sort ? (query.sort as string) : "blockNumber",
    // order: query.order ? (query.order as string) : "desc",
    // filter: query.filter ? (query.filter as string) : "",
  };
  return options;
};

export const apiHost = process.env.NEXT_PUBLIC_BRIDGE_API_HOST;
export const environment = process.env.NEXT_PUBLIC_BRIDGE_API_HOST_ENVIRONMENT;

export const l1Bridge = process.env.NEXT_PUBLIC_BRIDGE_L1;
export const l2Bridge = process.env.NEXT_PUBLIC_BRIDGE_L2;
export const messageBusAddress = process.env.NEXT_PUBLIC_BRIDGE_MESSAGE_BUS;

export const GOOGLE_ANALYTICS_ID =
  process.env.NEXT_PUBLIC_BRIDGE_GOOGLE_ANALYTICS_ID;

export const L1CHAINS = [
  {
    name: "Ethereum",
    value: "ETH",
    isNative: true,
    isEnabled: true,
    chainId: "0x1",
  },
];

export const L2CHAINS = [
  {
    name: "TEN",
    value: "TEN",
    isNative: false,
    isEnabled: true,
    chainId: "0x1bb",
  },
];

export const L2TOKENS = [
  {
    name: "Ether",
    value: "ETH",
    isNative: true,
    isEnabled: true,
    address: "",
  },
  {
    name: "USD Coin",
    value: "USDC",
    isNative: false,
    isEnabled: false,
    address: "0xb0E09857675Dc4c23ce90D4Ba62aC66fAb8b8155",
  },
  {
    name: "Tether USD",
    value: "USDT",
    isNative: false,
    isEnabled: false,
    address: "0x41ef84feDff3cE53d4C39097A81a74DD9A71280c",
  },
  {
    name: "TEN",
    value: "TEN",
    isNative: false,
    isEnabled: false,
    address: "",
  },
];

export const L1TOKENS = [
  {
    name: "Ether",
    value: "ETH",
    isNative: true,
    isEnabled: true,
    address: "",
  },
  {
    name: "USD Coin",
    value: "USDC",
    isNative: false,
    isEnabled: false,
    address: "0x718b239FFBB2dff8054ef424545A074d4EAbF220",
  },
  {
    name: "Tether USD",
    value: "USDT",
    isNative: false,
    isEnabled: false,
    address: "0x9Fa2813Fecc4706b3CA488EF21c0c73c7aD52c1F",
  },
  {
    name: "TEN",
    value: "TEN",
    isNative: false,
    isEnabled: false,
    address: "",
  },
];

export const PERCENTAGES = [
  {
    name: "25%",
    value: 25,
  },
  {
    name: "50%",
    value: 50,
  },
  {
    name: "MAX",
    value: 100,
  },
];
