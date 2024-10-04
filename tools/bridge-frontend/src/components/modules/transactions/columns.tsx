"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Badge, badgeVariants } from "@/src/components/ui/badge";

import { statuses } from "@/src/components/modules/transactions/constants";
import { DataTableColumnHeader } from "@/src/components/modules/common/data-table/data-table-column-header";
import TruncatedAddress from "../common/truncated-address";
import { Transactions } from "@/src/types";

export const columns: ColumnDef<Transactions>[] = [
  {
    id: "status",
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => {
      const status = statuses.find((s) => s.value === row.original.status);
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
    accessorKey: "blockNumber",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Ten Batch" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            #{Number(row.getValue("blockNumber"))}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "blockHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Ten Batch Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("blockHash")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "transactionHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Transaction Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("transactionHash")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "address",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Address" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("address")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
];
