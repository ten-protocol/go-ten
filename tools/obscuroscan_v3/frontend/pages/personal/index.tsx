import React from "react";
import { columns } from "@/components/modules/blockchain/transactions/columns";
import { DataTable } from "@/components/modules/common/data-table/data-table";
import Layout from "@/components/layouts/default-layout";
import { useTransactions } from "@/src/hooks/useTransactions";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Personal Transactions",
  description: "ObscuroScan Personal Transactions",
};

export default function PersonalTransactions() {
  const { personalTxns } = useTransactions();

  return (
    <>
      <Layout>
        <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
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
          {personalTxns ? (
            <DataTable columns={columns} data={personalTxns} />
          ) : (
            <p>Loading...</p>
          )}
        </div>
      </Layout>
    </>
  );
}
