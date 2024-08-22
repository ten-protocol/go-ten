import React from "react";
import { columns } from "@/src/components/modules/transactions/columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import Layout from "@/src/components/layouts/default-layout";
import { Metadata } from "next";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/src/components/ui/tabs";
import { useContractService } from "@/src/services/useContractService";

export const metadata: Metadata = {
  title: "Transactions",
  description: "A table of transactions.",
};

export default function Transactions() {
  const { getBridgeTransactions } = useContractService();
  const { transactions, refetchTransactions } = {
    transactions: {
      result: {
        TransactionsData: [],
        Total: 0,
      },
    },
    refetchTransactions: () => {},
  };
  const { TransactionsData, Total } = transactions?.result || {
    TransactionsData: [],
    Total: 0,
  };

  const getTransactions = async () => {
    const transactions = await getBridgeTransactions();
    console.log(transactions);
  };

  React.useEffect(() => {
    getTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <Layout>
      <div className="h-full flex-1 flex-col space-y-8 md:flex">
        <h1 className="text-2xl font-bold">Transaction History</h1>
        <p>
          View all transactions that have been made on the bridge. You can also
          filter transactions by status.
        </p>
        <Tabs defaultValue="all">
          <TabsList className="flex justify-start bg-background border-b">
            <TabsTrigger value="all">All Transactions</TabsTrigger>
            <TabsTrigger value="pending">Pending</TabsTrigger>
          </TabsList>
          <TabsContent value="all">
            {TransactionsData ? (
              <DataTable
                columns={columns}
                data={TransactionsData}
                refetch={refetchTransactions}
                total={+Total}
              />
            ) : (
              <p>Loading...</p>
            )}
          </TabsContent>
          <TabsContent value="pending">
            {TransactionsData ? (
              <DataTable
                columns={columns}
                data={TransactionsData}
                refetch={getTransactions}
                total={+Total}
              />
            ) : (
              <p>Loading...</p>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </Layout>
  );
}
