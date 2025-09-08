"use client";

import { ColumnDef } from "@tanstack/react-table";

import { DataTableColumnHeader } from "@repo/ui/components/common/data-table/data-table-column-header";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { formatNumber, formatTimeAgo } from "@repo/ui/lib/utils";
import { Batch } from "@/src/types/interfaces/BatchInterfaces";
import Link from "next/link";
import { Badge } from "@repo/ui/components/shared/badge";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";

export const columns: ColumnDef<Batch>[] = [
  {
    accessorKey: "number",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Batch" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <Link
            href={pathToUrl(pageLinks.batchByHeight, {
              height: row.original.height,
            })}
            className="text-primary"
          >
            <span className="max-w-[500px] truncate">
              #{Number(row.original.height)}
            </span>
          </Link>
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
            {row.original.header.timestamp
              ? formatTimeAgo(row.original.header.timestamp)
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
              {formatNumber(row.original?.header?.gasUsed) || "N/A"}
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
            <Badge variant={"outline"}>
              {formatNumber(row.original?.header?.gasUsed) || "N/A"}
            </Badge>
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
      return (
        <TruncatedAddress
          address={row.original.header.hash}
          link={pathToUrl(pageLinks.batchByHash, {
            hash: row.original.header.hash,
          })}
        />
      );
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
      return <TruncatedAddress address={row.original.header.parentHash} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "sequence",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Sequence" />
    ),
    cell: ({ row }) => {
      return (
        <Link
          href={pathToUrl(pageLinks.batchBySequence, {
            sequence: row.original.sequence,
          })}
          className="text-primary"
        >
          {row.original.sequence}
        </Link>
      )
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "txHashes",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Tx Count" />
    ),
    cell: ({ row }) => {
      return row.original.txHashes?.length > 0 ? (
        <Link
          href={pathToUrl(pageLinks.batchTransactions, {
            hash: row.original.fullHash,
          })}
          className="text-primary"
        >
          {row.original.txHashes?.length ?? 0}
        </Link>
      ) : (
        <span>{row.original.txHashes?.length ?? 0}</span>
      )
    },
    enableSorting: false,
    enableHiding: false,
  },
];
