import React from "react";
import {
  CardHeader,
  CardTitle,
  CardContent,
  Card,
} from "@repo/ui/components/shared/card";
import {
  LayersIcon,
  FileTextIcon,
  ReaderIcon,
  CubeIcon,
  RocketIcon,
  BlocksIcon,
} from "@repo/ui/components/shared/react-icons";

import { RecentBatches } from "./recent-batches";
import { RecentTransactions } from "./recent-transactions";
import { Button } from "@repo/ui/components/shared/button";
import { useTransactionsService } from "@/src/services/useTransactionsService";
import { useBatchesService } from "@/src/services/useBatchesService";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { useContractsService } from "@/src/services/useContractsService";
import { Skeleton } from "@repo/ui/components/shared/skeleton";
import AnalyticsCard from "./analytics-card";
import Link from "next/link";
import { cn, formatNumber } from "@repo/ui/lib/utils";
import { Badge } from "@repo/ui/components/shared/badge";

import { useRollupsService } from "@/src/services/useRollupsService";
import { RecentRollups } from "./recent-rollups";
import { DashboardAnalyticsData } from "@/src/types/interfaces";
import { pageLinks } from "@/src/routes";

interface RecentData {
  title: string;
  data: any;
  component: JSX.Element;
  goTo: string;
  className: string;
}

export default function Dashboard() {
  const {
    price,
    isPriceLoading,
    transactions,
    transactionCount,
    isTransactionCountLoading,
    setNoPolling: setNoPollingTransactions,
  } = useTransactionsService();
  const { contractCount, isContractCountLoading } = useContractsService();
  const {
    batches,
    latestBatch,
    isLatestBatchLoading,
    setNoPolling: setNoPollingBatches,
  } = useBatchesService();
  const { rollups, setNoPolling: setNoPollingRollups } = useRollupsService();

  React.useEffect(() => {
    setNoPollingTransactions(false);
    setNoPollingBatches(false);
    setNoPollingRollups(false);

    return () => {
      setNoPollingTransactions(true);
      setNoPollingBatches(true);
      setNoPollingRollups(true);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const DASHBOARD_DATA = [
    {
      title: "Ether Price",
      value: price?.ethereum?.usd
        ? `$${formatNumber(price.ethereum.usd)}`
        : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: RocketIcon,
      loading: isPriceLoading,
    },
    {
      title: "Latest L2 Batch",
      value: latestBatch?.item?.number
        ? Number(latestBatch.item.number)
        : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: LayersIcon,
      loading: isLatestBatchLoading,
    },
    {
      title: "Latest L1 Rollup",
      value: latestBatch?.item?.l1Proof ? (
        <TruncatedAddress
          address={latestBatch?.item?.l1Proof}
          prefixLength={6}
          suffixLength={4}
        />
      ) : (
        "N/A"
      ),
      // TODO: add change
      // change: "+20.1%",
      icon: CubeIcon,
      loading: isLatestBatchLoading,
    },
    {
      title: "Transactions",
      value: transactionCount?.count
        ? formatNumber(transactionCount.count)
        : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: ReaderIcon,
      loading: isTransactionCountLoading,
    },
    {
      title: "Contracts",
      value: contractCount?.count ? formatNumber(contractCount.count) : "N/A",
      // TODO: add change
      // change: "+20.1%",
      icon: FileTextIcon,
      loading: isContractCountLoading,
    },
    {
      title: "Nodes",
      value: <Badge variant={"static-default"}>Coming Soon</Badge>,
      icon: BlocksIcon,
    },
  ];

  const RECENT_DATA = [
    {
      title: "Recent Rollups",
      data: rollups,
      component: <RecentRollups rollups={rollups} />,
      goTo: pageLinks.rollups,
      className: "col-span-1 md:col-span-2 lg:col-span-3",
    },
    {
      title: "Recent Batches",
      data: batches,
      component: <RecentBatches batches={batches} />,
      goTo: pageLinks.batches,
      className: "col-span-1 md:col-span-2 lg:col-span-3",
    },
    {
      title: "Recent Transactions",
      data: transactions,
      component: <RecentTransactions transactions={transactions} />,
      goTo: pageLinks.transactions,
      className: "col-span-1 md:col-span-2 lg:col-span-3",
    },
  ];

  return (
    <div className="h-full flex-1 flex-col space-y-8 md:flex">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Tenscan</h2>
      </div>
      <div className="grid gap-4 grid-cols-1 sm:grid-cols-3 md:grid-cols-2 lg:grid-cols-4">
        {DASHBOARD_DATA.map((item: DashboardAnalyticsData, index: number) => (
          <AnalyticsCard key={index} item={item} />
        ))}
      </div>
      <div className="grid gap-4 grid-cols-1 md:grid-cols-6 lg:grid-cols-9">
        {RECENT_DATA.map((item: RecentData, index) => (
          <Card
            key={index}
            className={cn(item.className, "h-[450px] overflow-y-auto relative")}
          >
            <CardHeader className="flex flex-row items-center justify-between space-y-0 p-3 sticky top-0 left-0 right-0 bg-background z-10">
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
