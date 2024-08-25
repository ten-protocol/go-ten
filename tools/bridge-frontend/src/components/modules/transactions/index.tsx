import React from "react";
import { columns } from "@/src/components/modules/transactions/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import { getItem } from "@/src/lib/utils";
import { useContractService } from "@/src/services/useContractService";
import { ItemPosition } from "@/src/types";
import useWalletStore from "@/src/stores/wallet-store";

export default function TransactionsComponent() {
  const { isL1ToL2 } = useWalletStore();
  const { getBridgeTransactions } = useContractService();
  const [transactions, setTransactions] = React.useState<any[]>([]);
  const [isTransactionsLoading, setIsTransactionsLoading] =
    React.useState(true);

  const getTransactions = async () => {
    try {
      const txns = await getBridgeTransactions();
      setTransactions(txns);
    } catch (error) {
      console.error(error);
    }
  };

  const interval = React.useRef<any>();

  const fetchTransactions = () => {
    setIsTransactionsLoading(true);
    getTransactions();
    setIsTransactionsLoading(false);
  };

  React.useEffect(() => {
    fetchTransactions();

    interval.current = setInterval(() => {
      getTransactions();
    }, 10000);

    return () => {
      clearInterval(interval.current);
    };
  }, []);

  const firstBatchHeight = getItem(transactions, "blockNumber");
  const lastBatchHeight = getItem(
    transactions,
    "blockNumber",
    ItemPosition.LAST
  );

  return (
    <div className="h-full flex-1 flex-col space-y-8 md:flex">
      <div className="flex items-center justify-between space-y-2">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">
            Latest {isL1ToL2 ? "L1-L2" : "L2-L1"} Transactions
          </h2>
          {transactions?.length > 0 && (
            <p className="text-sm text-muted-foreground">
              Showing transactions in batch
              {firstBatchHeight !== lastBatchHeight && "es"} #{firstBatchHeight}{" "}
              {firstBatchHeight !== lastBatchHeight && "to #" + lastBatchHeight}
            </p>
          )}
        </div>
      </div>
      <DataTable
        columns={columns}
        data={transactions}
        refetch={fetchTransactions}
        total={transactions?.length}
        isLoading={isTransactionsLoading}
        noResultsText="transactions"
      />
    </div>
  );
}
