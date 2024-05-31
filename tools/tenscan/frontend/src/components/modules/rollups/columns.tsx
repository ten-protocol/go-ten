"use client";

import { ColumnDef } from "@tanstack/react-table";
import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import TruncatedAddress from "../common/truncated-address";
import { formatTimeAgo } from "@/src/lib/utils";
import Link from "next/link";
import { EyeOpenIcon } from "@radix-ui/react-icons";
import { Rollup } from "@/src/types/interfaces/RollupInterfaces";

export const columns: ColumnDef<Rollup>[] = [
  {
    accessorKey: "ID",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="ID" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">{row.getValue("ID")}</span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "Hash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Hash" />
    ),
    cell: ({ row }) => {
      return (
        <Link
          href={{
            pathname: "/rollup/[hash]/batches",
            query: { hash: row.getValue("Hash") },
          }}
          className="text-primary"
        >
          <TruncatedAddress address={row.getValue("Hash")} />
        </Link>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "Timestamp",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Timestamp" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            {formatTimeAgo(row.getValue("Timestamp"))}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "L1Hash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="L1 Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("L1Hash")} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "FirstSeq",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="First Batch Seq. No." />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <Link
            href={`/rollup/batch/sequence/${row.original.FirstSeq}`}
            className="text-primary"
          >
            <span className="max-w-[500px] truncate">
              {row.getValue("FirstSeq")}
            </span>
          </Link>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "LastSeq",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Last Batch Seq. No." />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <Link
            href={`/rollup/batch/sequence/${row.original.LastSeq}`}
            className="text-primary"
          >
            <span className="max-w-[500px] truncate">
              {row.getValue("LastSeq")}
            </span>
          </Link>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return (
        <Link href={`/rollup/${row.original.Hash}`}>
          <EyeOpenIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer" />
        </Link>
      );
    },
  },
];
