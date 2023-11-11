export const apiRoutes = {
  // **** BATCHES ****
  getLatestBatches: '/items/batch/latest/',
  getBatches: '/items/batches/',
  getBatchByHash: '/items/batch/:hash',

  // **** BLOCKS ****
  getBlocks: '/items/blocks/',

  // **** CONTRACTS ****
  getContractCount: '/count/contracts/',
  getVerifiedContracts: '/info/obscuro/',

  // **** TRANSACTIONS ****
  getTransactions: '/items/transactions/',
  getTransactionCount: '/count/transactions/',

  // **** ROLLUPS ****
  getRollups: '/items/rollups/latest/'
}
