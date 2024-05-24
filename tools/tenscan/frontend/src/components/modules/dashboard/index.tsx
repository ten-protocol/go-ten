import React from "react";
import { CardHeader, CardTitle, CardContent, Card } from "@repo/ui/shared/card";
import {
  LayersIcon,
  FileTextIcon,
  ReaderIcon,
  CubeIcon,
  RocketIcon,
} from "@repo/ui/shared/react-icons";

import { RecentBatches } from "./recent-batches";
import { RecentTransactions } from "./recent-transactions";
import { Button } from "@repo/ui/shared/button";
import { useTransactionsService } from "@/src/services/useTransactionsService";
import { useBatchesService } from "@/src/services/useBatchesService";
import TruncatedAddress from "../common/truncated-address";
import { useContractsService } from "@/src/services/useContractsService";
import { Skeleton } from "@repo/ui/shared/skeleton";
import { RecentBlocks } from "./recent-blocks";
import { useBlocksService } from "@/src/services/useBlocksService";
import AnalyticsCard from "./analytics-card";
import Link from "next/link";
import { cn, formatNumber } from "@/src/lib/utils";
import { Badge } from "@repo/ui/shared/badge";
import { BlocksIcon } from "@repo/ui/shared/react-icons";

export default function Dashboard() {
  const { price, transactions, transactionCount } = useTransactionsService();
  const { contractCount } = useContractsService();
  const { batches, latestBatch } = useBatchesService();
  const { blocks } = useBlocksService();

  const DASHBOARD_DATA = [
    {
      title: "Ether Price",
      value: price?.ethereum?.usd
        ? `$${formatNumber(price.ethereum.usd)}`
        : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: RocketIcon,
    },
    {
      title: "Latest L2 Batch",
      value: latestBatch?.item?.height
        ? Number(latestBatch.item.height)
        : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: LayersIcon,
    },
    {
      title: "Latest L1 Rollup",
      value: latestBatch?.item?.header?.l1Proof ? (
        <TruncatedAddress
          address={latestBatch?.item?.header?.l1Proof}
          prefixLength={6}
          suffixLength={4}
        />
      ) : (
        "N/A"
      ),
      // TODO: add change
      // change: "+20.1%",
      icon: CubeIcon,
    },
    {
      title: "Transactions",
      value: transactionCount?.count
        ? formatNumber(transactionCount.count)
        : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: ReaderIcon,
    },
    {
      title: "Contracts",
      value: contractCount?.count ? formatNumber(contractCount.count) : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: FileTextIcon,
    },
    {
      title: "Nodes",
      value: <Badge>Coming Soon</Badge>,
      icon: BlocksIcon,
    },
  ];

  const RECENT_DATA = [
    {
      title: "Recent Blocks",
      data: blocks,
      component: <RecentBlocks blocks={blocks} />,
      goTo: "/blocks",
      className: "col-span-1 md:col-span-2 lg:col-span-3",
    },
    {
      title: "Recent Batches",
      data: batches,
      component: <RecentBatches batches={batches} />,
      goTo: "/batches",
      className: "col-span-1 md:col-span-2 lg:col-span-3",
    },
    {
      title: "Recent Transactions",
      data: transactions,
      component: <RecentTransactions transactions={transactions} />,
      goTo: "/transactions",
      className: "col-span-1 md:col-span-2 lg:col-span-3",
    },
  ];

  return (
    <div className="h-full flex-1 flex-col space-y-8 md:flex">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Tenscan</h2>
      </div>
      <div className="grid gap-4 grid-cols-1 sm:grid-cols-3 md:grid-cols-2 lg:grid-cols-4">
        {DASHBOARD_DATA.map((item: any, index) => (
          <AnalyticsCard key={index} item={item} />
        ))}
      </div>
      <div className="grid gap-4 grid-cols-1 md:grid-cols-6 lg:grid-cols-9">
        {RECENT_DATA.map((item: any, index) => (
          <Card
            key={index}
            className={cn(item.className, "h-[450px] overflow-y-auto")}
          >
            <CardHeader className="flex flex-row items-center justify-between space-y-0 p-3">
              <CardTitle>{item.title}</CardTitle>
              <Link
                href={{
                  pathname: item.goTo,
                }}
              >
                <Button variant="outline" size="sm">
                  View All
                </Button>
              </Link>
            </CardHeader>
            <CardContent className="p-3">
              {item.data ? (
                item.component
              ) : (
                <Skeleton className="w-full h-[200px] rounded-lg" />
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
