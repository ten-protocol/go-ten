export type Batch = {
  sequence: number;
  hash: string;
  fullHash: string;
  height: number;
  txCount: number;
  header: {
    hash: string;
    parentHash: string;
    stateRoot: string;
    transactionsRoot: string;
    receiptsRoot: string;
    number: string;
    sequencerOrderNo: string;
    gasLimit: string;
    gasUsed: string;
    timestamp: string;
    extraData: string;
    baseFeePerGas: string;
    miner: string;
    l1Proof: string;
    signature: string;
    crossChainMessages: [];
    inboundCrossChainHash: string;
    inboundCrossChainHeight: string;
    TransfersTree: string;
    crossChainTree: string;
  };
  encryptedTxBlob: string;
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
  TxHashes: string[];
  EncryptedTxBlob: string;
};

export type BatchResponse = {
  BatchesData: Batch[];
  Total: string;
};
