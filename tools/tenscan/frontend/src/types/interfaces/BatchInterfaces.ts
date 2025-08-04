import { CrossChainMessage } from "./RollupInterfaces";

export type Batch = {
  sequence: number;
  hash: string;
  fullHash: string;
  height: number;
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
  txHashes: string[];
};

export interface LatestBatch {
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
  crossChainMessages: CrossChainMessage[];
  inboundCrossChainHash: string;
  inboundCrossChainHeight: string;
  TransfersTree: string;
  crossChainTree: string;
}

export type BatchResponse = {
  BatchesData: Batch[];
  Total: string;
};
