export type Batch = {
  parentHash: string;
  stateRoot: string;
  transactionsRoot: string;
  receiptsRoot: string;
  number: number;
  sequencerOrderNo: number;
  gasLimit: number;
  gasUsed: number;
  timestamp: string;
  extraData: string;
  baseFee: number;
  coinbase: string;
  l1Proof: string;
  R: number;
  S: number;
  crossChainMessages: any[];
  inboundCrossChainHash: string;
  inboundCrossChainHeight: number;
  transfersTree: string;
  hash: string;
  sha3Uncles: string;
  miner: string;
  logsBloom: string;
  difficulty: string;
  nonce: string;
  baseFeePerGas: number;
  EncryptedTxBlob: string;
};

export type BatchDetails = {
  Header: {
    parentHash: string;
    stateRoot: string;
    transactionsRoot: string;
    receiptsRoot: string;
    number: number;
    sequencerOrderNo: number;
    gasLimit: number;
    gasUsed: number;
    timestamp: string;
    extraData: string;
    baseFee: number;
    coinbase: string;
    l1Proof: string;
    R: number;
    S: number;
    crossChainMessages: any[];
    inboundCrossChainHash: string;
    inboundCrossChainHeight: number;
    transfersTree: string;
    hash: string;
    sha3Uncles: string;
    miner: string;
    logsBloom: string;
    difficulty: string;
    nonce: string;
    baseFeePerGas: number;
  };
  TxHashes: [];
  EncryptedTxBlob: string;
};

export type BatchResponse = {
  BatchesData: Batch[];
  Total: string;
};
