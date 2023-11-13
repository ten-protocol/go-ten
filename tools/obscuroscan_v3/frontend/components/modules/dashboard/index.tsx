import React from "react";
import { CalendarDateRangePicker } from "@/components/date-range-picker";
import { CardHeader, CardTitle, CardContent, Card } from "@/components/ui/card";
import { TabsList, TabsTrigger, TabsContent, Tabs } from "@/components/ui/tabs";
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
import { useRollups } from "@/src/hooks/useRollups";
import { useBatches } from "@/src/hooks/useBatches";
import TruncatedAddress from "../common/truncated-address";

export default function Dashboard() {
  const { transactions } = useTransactions();
  const { batches, latestBatch } = useBatches();

  const DASHBOARD_DATA = [
    {
      title: "Ether Price",
      value: "$1967.89",
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
      value:
        (
          <TruncatedAddress
            address={latestBatch?.item?.l1Proof}
            prefixLength={6}
            suffixLength={4}
          />
        ) || "N/A",
      change: "+20.1% from last month",
      icon: CubeIcon,
    },
    {
      title: "Transactions",
      value: "5",
      change: "+20.1% from last month",
      icon: ReaderIcon,
    },
    {
      title: "Contracts",
      value: "3",
      change: "+20.1% from last month",
      icon: FileTextIcon,
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
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="analytics" disabled>
            Analytics
          </TabsTrigger>
          <TabsTrigger value="reports" disabled>
            Reports
          </TabsTrigger>
          <TabsTrigger value="notifications" disabled>
            Notifications
          </TabsTrigger>
        </TabsList>
        <TabsContent value="overview" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {DASHBOARD_DATA.map((item: any, index) => (
              <Card key={index}>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    {item.title}
                  </CardTitle>
                  {React.createElement(item.icon)}
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold truncate">
                    {item.value}
                  </div>
                  <p className="text-xs text-muted-foreground">{item.change}</p>
                </CardContent>
              </Card>
            ))}
          </div>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-9">
            <Card className="col-span-6">
              <CardHeader>
                <CardTitle>Recent Batches</CardTitle>
              </CardHeader>
              <CardContent>
                <RecentBatches batches={batches} />
              </CardContent>
            </Card>
            {/* <Card className="col-span-3">
              <CardHeader>
                <CardTitle>Recent Rollups</CardTitle>
              </CardHeader>
              <CardContent>
                <RecentRollups rollups={rollups} />
              </CardContent>
            </Card> */}
            <Card className="col-span-3">
              <CardHeader>
                <CardTitle>Recent Transactions</CardTitle>
              </CardHeader>
              <CardContent>
                <RecentTransactions transactions={transactions} />
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
