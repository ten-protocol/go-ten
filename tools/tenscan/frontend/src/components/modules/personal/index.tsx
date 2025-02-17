import React from "react";
import { columns } from "@/src/components/modules/personal/columns";
import { DataTable } from "@repo/ui/components/common/data-table/data-table";
import { useTransactionsService } from "@/src/services/useTransactionsService";

export default function PersonalTransactions() {
  const { personalTxns, personalTxnsLoading } = useTransactionsService();
  const { Receipts, Total } = personalTxns || {
    Receipts: [],
    Total: 0,
  };

  return (
    <>
      <div className="flex items-center justify-between space-y-2">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">
            Personal Transactions
          </h2>
          {/* uncomment the following line when total count feature is implemented */}
          {/* <p className="text-muted-foreground">
            {formatNumber(Total)} personal transaction(s).
          </p> */}
        </div>
      </div>
      <DataTable
        columns={columns}
        data={Receipts}
        total={Total}
        isLoading={personalTxnsLoading}
        noResultsText="personal transactions"
      />
    </>
  );
}
