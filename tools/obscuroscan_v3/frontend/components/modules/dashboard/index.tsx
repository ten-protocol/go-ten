import React from "react";
import { CalendarDateRangePicker } from "@/components/date-range-picker";
import { CardHeader, CardTitle, CardContent, Card } from "@/components/ui/card";
import {
  LayersIcon,
  FileTextIcon,
  ReaderIcon,
  CubeIcon,
  RocketIcon,
} from "@radix-ui/react-icons";

import { RecentBatches } from "./recent-batches";
import { RecentTransactions } from "./recent-transactions";
import { Button } from "@/components/ui/button";
import { useTransactions } from "@/src/hooks/useTransactions";
import { useBatches } from "@/src/hooks/useBatches";
import TruncatedAddress from "../common/truncated-address";
import { useContracts } from "@/src/hooks/useContracts";
import { Skeleton } from "@/components/ui/skeleton";
import { RecentBlocks } from "./recent-blocks";
import { useBlocks } from "@/src/hooks/useBlocks";
import AnalyticsCard from "./analytics-card";
import Link from "next/link";

export default function Dashboard() {
  const { price, transactions, transactionCount } = useTransactions();
  const { contractCount } = useContracts();
  const { batches, latestBatch } = useBatches();
  const { blocks } = useBlocks();

  const DASHBOARD_DATA = [
    {
      title: "Ether Price",
      value: price?.ethereum?.usd,
      change: "+20.1% from last month",
      icon: RocketIcon,
    },
    {
      title: "Latest Batch",
      value: batches?.result?.Total,
      change: "+20.1% from last month",
      icon: LayersIcon,
    },
    {
      title: "Latest Rollup",
      value: latestBatch?.item?.l1Proof ? (
        <TruncatedAddress
          address={latestBatch?.item?.l1Proof}
          prefixLength={6}
          suffixLength={4}
        />
      ) : (
        "N/A"
      ),
      change: "+20.1% from last month",
      icon: CubeIcon,
    },
    {
      title: "Transactions",
      value: transactionCount?.count,
      change: "+20.1% from last month",
      icon: ReaderIcon,
    },
    {
      title: "Contracts",
      value: contractCount?.count,
      change: "+20.1% from last month",
      icon: FileTextIcon,
    },
  ];

  const RECENT_DATA = [
    {
      title: "Recent Batches",
      data: batches,
      component: <RecentBatches batches={batches} />,
      goTo: "/blockchain/batches",
    },
    {
      title: "Recent Blocks",
      data: blocks,
      component: <RecentBlocks blocks={blocks} />,
      goTo: "/blockchain/blocks",
    },
    {
      title: "Recent Transactions",
      data: transactions,
      component: <RecentTransactions transactions={transactions} />,
      goTo: "/blockchain/transactions",
    },
  ];

  return (
    <div className="flex-1 space-y-4 p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Obscuroscan</h2>
        <div className="flex items-center space-x-2">
          <CalendarDateRangePicker />
          <Button>Download</Button>
        </div>
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {DASHBOARD_DATA.map((item: any, index) => (
          <AnalyticsCard key={index} item={item} />
        ))}
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-9">
        {RECENT_DATA.map((item: any, index) => (
          <Card key={index} className="col-span-4">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
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
            <CardContent>
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
