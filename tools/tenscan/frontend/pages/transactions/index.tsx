import React from "react";
import { columns } from "@/src/components/modules/transactions/columns";
import { DataTable } from "@repo/ui/components/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { useTransactionsService } from "@/src/services/useTransactionsService";
import { Metadata } from "next";
import { getItem } from "@repo/ui/lib/utils";
import HeadSeo from "@/src/components/head-seo";
import { siteMetadata } from "@/src/lib/siteMetadata";
import { ItemPosition } from "@repo/ui/lib/enums/ui";

export const metadata: Metadata = {
  title: "Transactions",
  description: "A table of transactions.",
};

export default function Transactions() {
  const { transactions, refetchTransactions, isTransactionsLoading } =
    useTransactionsService();
  const { TransactionsData, Total } = transactions?.result || {
    TransactionsData: [],
    Total: 0,
  };

 const firstBatchHeight = TransactionsData?.[0]?.BatchHeight
 const lastBatchHeight =
   TransactionsData?.[TransactionsData.length - 1]?.BatchHeight

  return (
    <>
      <HeadSeo
        title={`${siteMetadata.transactions.title} `}
        description={siteMetadata.transactions.description}
        canonicalUrl={`${siteMetadata.transactions.canonicalUrl}`}
        ogImageUrl={siteMetadata.transactions.ogImageUrl}
        ogTwitterImage={siteMetadata.transactions.ogTwitterImage}
        ogType={siteMetadata.transactions.ogType}
      ></HeadSeo>
      <Layout>
        <div className="h-full flex-1 flex-col space-y-8 md:flex">
          <div className="flex items-center justify-between space-y-2">
            <div>
              <h2 className="text-2xl font-bold tracking-tight">
                Transactions
              </h2>
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
          <DataTable
            columns={columns}
            data={TransactionsData}
            refetch={refetchTransactions}
            total={+Total}
            isLoading={isTransactionsLoading}
            noResultsText="transactions"
          />
        </div>
      </Layout>
    </>
  );
}
