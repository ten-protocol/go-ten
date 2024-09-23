import { Card, CardHeader, CardTitle, CardContent } from "@repo/ui/shared/card";
import { Skeleton } from "@repo/ui/shared/skeleton";
import React from "react";

export default function AnalyticsCard({
  item,
}: {
  item: { title: string; value: string; change: string; icon: any };
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{item.title}</CardTitle>
        {React.createElement(item.icon)}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold truncate mb-1">
          {item.value ? (
            item.value
          ) : (
            <Skeleton className="w-[100px] h-[20px] rounded-full" />
          )}
        </div>
        {item.change && (
          <p className="text-xs text-muted-foreground">{item.change}</p>
        )}
      </CardContent>
    </Card>
  );
}
