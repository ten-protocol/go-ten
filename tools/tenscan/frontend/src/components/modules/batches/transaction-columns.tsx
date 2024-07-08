"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Badge } from "@/src/components/ui/badge";

import { statuses } from "../transactions/constants";
import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import TruncatedAddress from "../common/truncated-address";
import { formatNumber, formatTimeAgo } from "@/src/lib/utils";
import Link from "next/link";
import { EyeOpenIcon } from "@radix-ui/react-icons";

export const columns: ColumnDef<Transaction>[] = [
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
    enableSorting: false,
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
          link={`/tx/${row.original.TransactionHash}`}
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
