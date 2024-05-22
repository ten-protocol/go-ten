export interface RollupsResponse {
  RollupsData: Rollup[];
  Total: number;
}

export interface Rollup {
  ID: number;
  Hash: string;
  FirstSeq: number;
  LastSeq: number;
  Timestamp: number;
  Header: Header;
  L1Hash: string;
}

export interface Header {
  CompressionL1Head: string;
  crossChainMessages: CrossChainMessage[];
  PayloadHash: string;
  Signature: string;
  LastBatchSeqNo: number;
  hash: string;
}

type CrossChainMessage = {
  Sender: string;
  Sequence: number;
  Nonce: number;
  Topic: number;
  Payload: string[];
  ConsistencyLevel: number;
};
