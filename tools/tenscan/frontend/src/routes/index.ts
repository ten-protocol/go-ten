import { NavLink } from "../types/interfaces";

export const apiRoutes = {
  // **** BATCHES ****
  getLatestBatch: "/items/batch/latest/",
  getBatches: "/items/v2/batches/",
  getBatchByHash: "/items/batch/:hash",
  getBatchByHeight: "/items/batch/height/:height",
  getBatchTransactions: "/items/batch/:fullHash/transactions",
  getBatchesInRollup: "/items/rollup/:hash/batches",

  // **** BLOCKS ****
  getBlocks: "/items/blocks/",

  // **** CONTRACTS ****
  getContractCount: "/count/contracts/",
  getVerifiedContracts: "/info/obscuro/",

  // **** TRANSACTIONS ****
  getTransactions: "/items/transactions/",
  getTransactionCount: "/count/transactions/",
  getTransactionByHash: "/items/transaction/:hash",

  getEtherPrice:
    "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd",

  // **** ROLLUPS ****
  getLatestRollup: "/items/rollup/latest/",
  decryptEncryptedRollup: "/actions/decryptTxBlob/",
  getRollups: "/items/rollups/",
  getRollupByHash: "/items/rollup/:hash",
  getRollupByBatchSequence: "/items/rollup/batch/:seq",

  // **** INFO ****
  getHealthStatus: "/info/health/",
};

export const ethMethods = {
  getStorageAt: "eth_getStorageAt",
  getTransactionReceipt: "eth_getTransactionReceipt",
};
// to send TEN Custom Queries (CQ) through the provider we call eth_getStorageAt and use these addresses to identify the TEN CQ method
export const tenCustomQueryMethods = {
  getUserID: "0x0000000000000000000000000000000000000001",
  listPersonalTransactions: "0x0000000000000000000000000000000000000002",
};

export const NavLinks: NavLink[] = [
  {
    href: "/",
    label: "Home",
    isExternal: false,
    isDropdown: false,
  },
  {
    href: "/personal",
    label: "Personal",
    isExternal: false,
    isDropdown: false,
  },
  {
    label: "Blockchain",
    isExternal: false,
    isDropdown: true,
    subNavLinks: [
      {
        href: "/transactions",
        label: "Transactions",
        isExternal: false,
      },
      {
        href: "/batches",
        label: "Batches",
        isExternal: false,
      },
      {
        href: "/rollups",
        label: "Rollups",
        isExternal: false,
      },
    ],
  },
  {
    label: "Resources",
    isExternal: false,
    isDropdown: true,
    subNavLinks: [
      {
        href: "/resources/decrypt",
        label: "Decrypt",
        isExternal: false,
      },
      {
        href: "/resources/verified-data",
        label: "Verified Data",
        isExternal: false,
      },
    ],
  },
];

export const externalLinks = {
  // Dev and UAT environments don't have etherscan, hence, we're hardcoding this to just Sepolia
  etherscanBlock: "https://sepolia.etherscan.io/block/",
};
