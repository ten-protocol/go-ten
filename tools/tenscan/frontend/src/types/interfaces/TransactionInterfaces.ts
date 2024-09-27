export type Transaction = {
  Id: string;
  label: string;
  BatchHeight: number;
  BatchTimestamp: number;
  Finality: string;
  TransactionHash: string;
};

export type TransactionCount = {
  count: number;
};

export type Price = {
  ethereum: {
    usd: number;
  };
};

export type TransactionResponse = {
  TransactionsData: Transaction[];
  Total: number;
};

export type PersonalTransactionsResponse = {
  Receipts: PersonalTransactions[];
  Total: number;
};

export type TransactionType = "0x0" | "0x1" | "0x2" | "0x3";

export enum PersonalTransactionType {
  Legacy = "0x0",
  AccessList = "0x1",
  DynamicFee = "0x2",
  Blob = "0x3",
}

export type PersonalTransactions = {
  id: number;
  blockNumber: string;
  transactionHash: string;
  status: string;
  gasUsed: string;
  blockHash: string;
  contractAddress: string;
  cumulativeGasUsed: string;
  effectiveGasPrice: string;
  logs: any[];
  logsBloom: string;
  root: string;
  transactionIndex: string;
  type: TransactionType;
};

export type TransactionReceipt = {
  blockHash: string;
  blockNumber: string;
  contractAddress: string;
  cumulativeGasUsed: string;
  effectiveGasPrice: string;
  from: string;
  gasUsed: string;
  logs: Log[];
  logsBloom: string;
  status: string;
  to: string;
  transactionHash: string;
  transactionIndex: string;
  type: TransactionType;
};

export type Log = {
  address: string;
  blockHash: string;
  blockNumber: string;
  data: string;
  logIndex: string;
  removed: boolean;
  topics: string[];
  transactionHash: string;
  transactionIndex: string;
};
