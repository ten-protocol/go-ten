import React from "react";
import { columns } from "@/src/components/modules/personal/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import { useTransactionsService } from "@/src/services/useTransactionsService";
import { toolbar } from "./data";
import { Skeleton } from "@/src/components/ui/skeleton";

export default function PersonalTransactions() {
  const { personalTxns, personalTxnsLoading } = useTransactionsService();

  return (
    <>
      <div className="flex items-center justify-between space-y-2">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">
            Personal Transactions
          </h2>
          <p className="text-muted-foreground">
            A table of personal transactions.
          </p>
        </div>
      </div>
      {personalTxnsLoading ? (
        <Skeleton className="h-96" />
      ) : personalTxns?.Results ? (
        <DataTable
          columns={columns}
          data={personalTxns?.Results}
          toolbar={toolbar}
        />
      ) : (
        <p>No transactions found.</p>
      )}
    </>
  );
}
