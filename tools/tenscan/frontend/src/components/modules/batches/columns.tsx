"use client";

import { ColumnDef } from "@tanstack/react-table";

import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import TruncatedAddress from "../common/truncated-address";
import { formatNumber, formatTimeAgo } from "@/src/lib/utils";
import { Batch } from "@/src/types/interfaces/BatchInterfaces";
import { EyeOpenIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { Badge } from "../../ui/badge";

export const columns: ColumnDef<Batch>[] = [
  {
    accessorKey: "number",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Batch" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            #{Number(row.getValue("number"))}
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
          <span className="max-w-[500px] truncate">
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
          <span className="max-w-[500px] truncate">
            <Badge variant={"outline"}>
              {formatNumber(row.getValue("gasUsed"))}
            </Badge>
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
          <span className="max-w-[500px] truncate">
            {formatNumber(row.getValue("gasLimit"))}
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
    accessorKey: "miner",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Miner" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.original.miner} />;
    },
    enableSorting: false,
    enableHiding: false,
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
