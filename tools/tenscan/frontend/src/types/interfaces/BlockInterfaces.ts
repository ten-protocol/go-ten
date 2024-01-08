export type Block = {
  blockHeader: BlockHeader;
  rollupHash: string;
};

export type BlockHeader = {
  parentHash: string;
  sha3Uncles: string;
  miner: string;
  stateRoot: string;
  transactionsRoot: string;
  receiptsRoot: string;
  logsBloom: string;
  difficulty: string;
  number: string;
  gasLimit: string;
  gasUsed: string;
  timestamp: string;
  extraData: string;
  mixHash: string;
  nonce: string;
  baseFeePerGas: string;
  withdrawalsRoot: null;
  blobGasUsed: null;
  excessBlobGas: null;
  hash: string;
};
