import { IPendingTx } from "@/src/types";
import { PENDING_TRANSACTIONS_KEY } from "../constants";
import { handleStorage } from "../utils";

// get all pending transactions
const getPendingBridgeTransactions = () => {
  const transactions = handleStorage.get(PENDING_TRANSACTIONS_KEY);
  return transactions ? JSON.parse(transactions) : [];
};

// save pending txns list to local storage
const savePendingBridgeTransactions = (transactions: IPendingTx[]) => {
  handleStorage.save(PENDING_TRANSACTIONS_KEY, JSON.stringify(transactions));
};

// add a new pending txn
const addPendingBridgeTransaction = (transaction: IPendingTx) => {
  const transactions = getPendingBridgeTransactions();
  // adding the new txn to the beginning of the list
  transactions.unshift(transaction);
  savePendingBridgeTransactions(transactions);
};

// update a specific pending txn by its hash
const updatePendingBridgeTransaction = (
  txHash: string,
  txnUpdates: IPendingTx
) => {
  const transactions = getPendingBridgeTransactions();
  const index = transactions.findIndex(
    (tx: IPendingTx) => tx.txHash === txHash
  );

  if (index !== -1) {
    transactions[index] = { ...transactions[index], ...txnUpdates };
    savePendingBridgeTransactions(transactions);
  }
};

// rm a completed/cancelled txn
const removePendingBridgeTransaction = (txHash: string) => {
  const transactions = getPendingBridgeTransactions();
  const updatedTransactions = transactions.filter(
    (tx: any) => tx.txHash !== txHash
  );
  savePendingBridgeTransactions(updatedTransactions);
};

export {
  getPendingBridgeTransactions,
  savePendingBridgeTransactions,
  addPendingBridgeTransaction,
  updatePendingBridgeTransaction,
  removePendingBridgeTransaction,
};
