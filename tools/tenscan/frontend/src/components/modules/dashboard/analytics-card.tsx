import { DashboardAnalyticsData } from "@/src/types/interfaces";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from "@repo/ui/components/shared/card";
import { Skeleton } from "@repo/ui/components/shared/skeleton";
import React from "react";
import { cn } from "@repo/ui/lib/utils";

export default function AnalyticsCard({
  item,
}: {
  item: DashboardAnalyticsData;
}) {
  return (
    <div className="relative">
      <div className="stat-shape absolute inset-[1px] pointer-events-none">
        <div className="absolute inset-0 animate-scan-overlay" />
      </div>
      <svg
        width="100%"
        height="100%"
        viewBox="0 0 330 120"
        preserveAspectRatio="none"
        version="1.1"
        xmlns="http://www.w3.org/2000/svg"
        className="z-10 absolute inset-0 pointer-events-none"
      >
        <g>
          {/*  TODO: Decide whether to remove stroke*/}
          <g transform="matrix(0.99697,0,0,0.995721,0.5,0.124897)">
            <path
              d="M10,0L0,10L0,110L10,120L320,120L330,110L330,10L320,0L260,0L254,6L76,6L70,0L10,0Z"
              style={{
                fill: "none",
                stroke: "rgba(255,255,255,.1)",
                strokeWidth: "1px",
              }}
            />
          </g>
          <g transform="matrix(1,0,0,0.752747,0,0.236264)">
            <path
              fill="rgba(255,255,255,0.2)"
              d="M73,-0.314L77.5,5.664L253,5.664L257,-0.314L73,-0.314Z"
            />
          </g>
        </g>
      </svg>

      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-lg font-medium opacity-80">{item.title}</CardTitle>
        {React.createElement(item.icon)}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold truncate mb-1">
          {item.loading ? (
            <Skeleton className="w-[100px] h-[20px] rounded-full" />
          ) : (
            <h6 className="text-4xl font-bold truncate">{item.value}</h6>
          )}
        </div>
        {item?.change && (
          <p className="text-xs text-muted-foreground">{item.change}</p>
        )}
      </CardContent>
    </div>
  );
}
