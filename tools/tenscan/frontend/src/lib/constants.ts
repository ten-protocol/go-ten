export const socialLinks = {
  github: "https://github.com/ten-protocol",
  discord: "https://discord.gg/tenprotocol",
  twitter: "https://twitter.com/tenprotocol",
  twitterHandle: "@tenprotocol",
};

export const pollingInterval = 5000;
export const maxRetries = 3;
export const pricePollingInterval = 60 * 1000;

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

export const version = process.env.NEXT_PUBLIC_FE_VERSION;
export const apiHost = process.env.NEXT_PUBLIC_API_HOST;

export const currentEncryptedKey =
  "bddbc0d46a0666ce57a466168d99c1830b0c65e052d77188f2cbfc3f6486588c";

export const GOOGLE_ANALYTICS_ID = "G-M82QX9RT4L";
