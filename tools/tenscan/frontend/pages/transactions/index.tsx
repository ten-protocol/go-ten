import React from "react";
import { columns } from "@/src/components/modules/transactions/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { useTransactionsService } from "@/src/services/useTransactionsService";
import { Metadata } from "next";
import { firstItem, lastItem } from "@/src/lib/utils";

export const metadata: Metadata = {
  title: "Transactions",
  description: "A table of transactions.",
};

export default function Transactions() {
  const {
    transactions,
    refetchTransactions,
    setNoPolling,
    isTransactionsLoading,
  } = useTransactionsService();
  const { TransactionsData, Total } = transactions?.result || {
    TransactionsData: [],
    Total: 0,
  };

  React.useEffect(() => {
    setNoPolling(true);

    return () => {
      setNoPolling(false);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const firstBatchHeight = firstItem(TransactionsData, "BatchHeight");
  const lastBatchHeight = lastItem(TransactionsData, "BatchHeight");

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Transactions</h2>
            {TransactionsData?.length > 0 && (
              <p className="text-sm text-muted-foreground">
                Showing transactions in batch
                {firstBatchHeight !== lastBatchHeight && "es"} #
                {firstBatchHeight}{" "}
                {firstBatchHeight !== lastBatchHeight &&
                  "to #" + lastBatchHeight}
                {/* uncomment the following line when total count feature is implemented */}
                {/* of {formatNumber(Total)} transactions. */}
              </p>
            )}
          </div>
        </div>
        {TransactionsData ? (
          <DataTable
            columns={columns}
            data={TransactionsData}
            refetch={refetchTransactions}
            total={+Total}
            isLoading={isTransactionsLoading}
          />
        ) : (
          <div>No rollups found.</div>
        )}
      </div>
    </Layout>
  );
}
