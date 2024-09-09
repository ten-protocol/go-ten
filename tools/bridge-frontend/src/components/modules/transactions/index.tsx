import React from "react";
import { columns } from "@/src/components/modules/transactions/columns";
import { columns as pending } from "@/src/components/modules/transactions/pending-columns";
import { DataTable } from "@/src/components/modules/common/data-table/data-table";
import { getItem } from "@/src/lib/utils";
import { useContractsService } from "@/src/services/useContractsService";
import { ItemPosition } from "@/src/types";
import useWalletStore from "@/src/stores/wallet-store";
import { useQuery } from "@tanstack/react-query";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../../ui/tabs";
import { getPendingBridgeTransactions } from "@/src/lib/utils/txnUtils";
import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import { Button } from "../../ui/button";

export default function TransactionsComponent() {
  const { isL1ToL2 } = useWalletStore();
  const { getBridgeTransactions, finaliseTransaction, finalisingTxHashes } =
    useContractsService();

  const pendingColumns = [
    ...pending,
    {
      accessorKey: "actions",
      header: ({ column }: { column: any }) => (
        <DataTableColumnHeader column={column} title="Actions" />
      ),
      cell: ({ row }: { row: any }) => {
        const txHash = row.getValue("txHash");
        const isDisabled = finalisingTxHashes.has(txHash as string);
        return (
          <Button
            variant="secondary"
            size="sm"
            className="dark:bg-[#292929]"
            onClick={() => {
              const tx = row.original;
              finaliseTransaction(tx);
            }}
            disabled={isDisabled}
          >
            {isDisabled ? "Finalising..." : "Finalise"}
          </Button>
        );
      },
      enableSorting: false,
      enableHiding: false,
    },
  ];

  const {
    data: transactions = [],
    isLoading: isTransactionsLoading,
    refetch,
  } = useQuery({
    queryKey: ["bridgeTransactions", isL1ToL2 ? "l1" : "l2"],
    queryFn: () => getBridgeTransactions(),
    refetchInterval: 10000,
    refetchOnMount: true,
  });

  const {
    data: pendingTransactions = [],
    isLoading: isPendingTransactionsLoading,
    refetch: refetchPending,
  } = useQuery({
    queryKey: ["bridgePendingTransactions", isL1ToL2 ? "l1" : "l2"],
    queryFn: () => getPendingBridgeTransactions(),
    refetchInterval: 10000,
    refetchOnMount: true,
  });

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
      <Tabs defaultValue="all">
        <TabsList className="flex justify-start bg-background border-b">
          <TabsTrigger value="all">All Transactions</TabsTrigger>
          <TabsTrigger value="pending">Pending</TabsTrigger>
        </TabsList>
        <TabsContent value="all">
          <DataTable
            columns={columns}
            data={transactions}
            refetch={refetch}
            total={transactions?.length}
            isLoading={isTransactionsLoading}
            noResultsText="transactions"
            noPagination={true}
          />
        </TabsContent>
        <TabsContent value="pending">
          <DataTable
            columns={pendingColumns}
            data={pendingTransactions}
            refetch={refetchPending}
            total={pendingTransactions?.length}
            isLoading={isPendingTransactionsLoading}
            noResultsText="Pending Transactions"
            noPagination={true}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
