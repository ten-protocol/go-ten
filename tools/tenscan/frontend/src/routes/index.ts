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

export const externalPageLinks = {
  // Dev and UAT environments don't have etherscan, hence, we're hardcoding this to just Sepolia
  etherscanBlock: "https://sepolia.etherscan.io/block/:hash",
};

export const pageLinks = {
  home: "/",
  address: "/address/:address",

  // **** TRANSACTIONS ****
  personalTransactions: "/personal",
  transactions: "/transactions",
  txByHash: "/tx/:hash",
  personalTxByHash: "/personal/tx/:hash",

  // **** BATCHES ****
  batches: "/batches",
  batchByHash: "/batch/:hash",
  batchByHeight: "/batch/height/:height",
  batchTransactions: "/batch/txs/:hash",

  // **** ROLLUPS ****
  rollups: "/rollups",
  rollupByHash: "/rollup/:hash",
  rollupByBatchSequence: "/rollup/batch/sequence/:sequence",
  rollupBatches: "/rollup/:hash/batches",

  // **** BLOCKS ****
  blocks: "/blocks",

  // **** CONTRACTS ****
  verifiedData: "/resources/verified-data",
  decrypt: "/resources/decrypt",

  // **** INFO ****
  privacyPolicy: "/docs/privacy-policy",
  termsOfService: "/docs/terms",
};

export const NavLinks: NavLink[] = [
  {
    href: pageLinks.home,
    label: "Home",
    isExternal: false,
    isDropdown: false,
  },
  {
    href: pageLinks.personalTransactions,
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
        href: pageLinks.transactions,
        label: "Transactions",
        isExternal: false,
      },
      {
        href: pageLinks.batches,
        label: "Batches",
        isExternal: false,
      },
      {
        href: pageLinks.rollups,
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
        href: pageLinks.verifiedData,
        label: "Verified Data",
        isExternal: false,
      },
    ],
  },
];
