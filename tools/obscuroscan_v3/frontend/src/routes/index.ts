export const apiRoutes = {
  // **** BATCHES ****
  getLatestBatch: "/items/batch/latest/",
  getBatches: "/items/batches/",
  getBatchByHash: "/items/batch/:hash",

  // **** BLOCKS ****
  getBlocks: "/items/blocks/",

  // **** CONTRACTS ****
  getContractCount: "/count/contracts/",
  getVerifiedContracts: "/info/obscuro/",

  // **** TRANSACTIONS ****
  getTransactions: "/items/transactions/",
  getTransactionCount: "/count/transactions/",

  getPrice:
    "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd",

  // **** ROLLUPS ****
  getRollups: "/items/rollup/latest/",
  decryptEncryptedRollup: "/actions/decryptTxBlob/",
};
