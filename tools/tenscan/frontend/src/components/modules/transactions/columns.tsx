"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Badge } from "@repo/ui/components/shared/badge";

import { statuses } from "./constants";
import { DataTableColumnHeader } from "@repo/ui/components/common/data-table/data-table-column-header";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { formatTimeAgo } from "@repo/ui/lib/utils";
import Link from "next/link";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";

export const columns: ColumnDef<Transaction>[] = [
  {
    accessorKey: "BatchHeight",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Batch" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            #{row.getValue("BatchHeight")}
          </span>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
    enableHiding: false,
  },

  {
    accessorKey: "BatchTimestamp",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Batch Age" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            {formatTimeAgo(row.getValue("BatchTimestamp"))}
          </span>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
    enableHiding: false,
  },

  {
    accessorKey: "TransactionHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Transaction Hash" />
    ),
    cell: ({ row }) => {
      return (
        <TruncatedAddress
          address={row.getValue("TransactionHash")}
          link={pathToUrl(pageLinks.txByHash, {
            hash: row.original.TransactionHash,
          })}
        />
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "Finality",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Finality" />
    ),
    cell: ({ row }) => {
      const finality = statuses.find(
        (finality) => finality.value === row.getValue("Finality")
      );

      if (!finality) {
        return null;
      }

      return (
        <div className="flex items-center">
          {finality.icon && (
            <finality.icon className="mr-2 h-4 w-4 text-muted-foreground" />
          )}
          <Badge>{finality.label}</Badge>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
];
