export type Transaction = {
  Id: string;
  label: string;
  BatchHeight: number;
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
  Result: PersonalTransactions[];
  Total: number;
};

export type PersonalTransactions = {
  id: number;
  blockNumber: string;
  transactionHash: string;
  status: "Success" | "Failed";
  gasUsed: string;
  blockHash: string;
};
