"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Badge, badgeVariants } from "@/src/components/ui/badge";

import { statuses } from "@/src/components/modules/transactions/constants";
import { DataTableColumnHeader } from "@/src/components/modules/common/data-table/data-table-column-header";
import TruncatedAddress from "../common/truncated-address";
import { IPendingTx } from "@/src/types";
import { formatDistanceToNow } from "date-fns";

export const columns: ColumnDef<IPendingTx>[] = [
  {
    id: "status",
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => {
      const status = statuses.find((s) => s.value === "Pending");
      return (
        <Badge variant={status?.variant as keyof typeof badgeVariants}>
          {status?.icon && <status.icon className="h-5 w-5 mr-2" />}
          {status?.label}
        </Badge>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "timestamp",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Timestamp" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            {formatDistanceToNow(new Date(row.getValue("timestamp")), {
              addSuffix: true,
            })}
          </span>
        </div>
      );
    },
  },
  {
    accessorKey: "txHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Txn Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("txHash")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "resumeStep",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Txn Step" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <code className="max-w-[500px] truncate">
            {row.getValue("resumeStep")}
          </code>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "receiver",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Receiver" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("receiver")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "value",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Amount" />
    ),
    cell: ({ row }) => {
      return <Badge variant={"outline"}>{row.getValue("value")}</Badge>;
    },
    enableSorting: false,
    enableHiding: false,
  },
];
