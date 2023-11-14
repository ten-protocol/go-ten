"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Checkbox } from "@/components/ui/checkbox";

import { labels, statuses } from "./data";
import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import TruncatedAddress from "../common/truncated-address";
import { formatTimeAgo } from "@/src/lib/utils";
import { Batch } from "@/src/types/interfaces/BatchInterfaces";
import { EyeOpenIcon } from "@radix-ui/react-icons";
import Link from "next/link";

export const columns: ColumnDef<Batch>[] = [
  {
    id: "select",
    header: ({ table }) => (
      <Checkbox
        checked={table.getIsAllPageRowsSelected()}
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
        className="translate-y-[2px]"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
        className="translate-y-[2px]"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "number",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Batch" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {row.getValue("number")}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "timestamp",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Age" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {row.getValue("timestamp")
              ? formatTimeAgo(row.getValue("timestamp"))
              : "N/A"}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "gasUsed",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Gas Used" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {row.getValue("gasUsed")}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "gasLimit",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Gas Limit" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {row.getValue("gasLimit")}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "hash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("hash")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "parentHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Parent Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("parentHash")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "l1Proof",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="L1 Proof" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.original.l1Proof} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "difficulty",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Difficulty" />
    ),
    cell: ({ row }) => {
      const finality = statuses.find(
        (finality) => finality.value === row.getValue("difficulty")
      );

      if (!finality) {
        return null;
      }

      return (
        <div className="flex items-center">
          {finality.icon && (
            <finality.icon className="mr-2 h-4 w-4 text-muted-foreground" />
          )}
          <span>{finality.label || "N/A"}</span>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return (
        <Link href={`/batches/${row.original.hash}`}>
          <EyeOpenIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer" />
        </Link>
      );
    },
  },
];
