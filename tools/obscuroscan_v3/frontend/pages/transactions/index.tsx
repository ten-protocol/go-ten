import React from "react";
import { columns } from "@/src/components/modules/transactions/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { useTransactionsService } from "@/src/hooks/useTransactionsService";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Transactions",
  description: "A table of transactions.",
};

export default function Transactions() {
  const { transactions } = useTransactionsService();

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Transactions</h2>
            <p className="text-muted-foreground">A table of transactions.</p>
          </div>
        </div>
        {transactions?.result?.TransactionsData ? (
          <DataTable
            columns={columns}
            data={transactions?.result?.TransactionsData}
          />
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </Layout>
  );
}
