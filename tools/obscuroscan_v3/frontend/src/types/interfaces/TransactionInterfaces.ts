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
  TransactionData: Transaction[];
  Total: number;
};
