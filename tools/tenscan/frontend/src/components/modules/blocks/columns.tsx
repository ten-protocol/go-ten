"use client";

import Link from "next/link";
import { ColumnDef } from "@tanstack/react-table";

import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import { Block, BlockHeader } from "@/src/types/interfaces/BlockInterfaces";
import TruncatedAddress from "../common/truncated-address";
import { formatNumber, formatTimeAgo } from "@/src/lib/utils";
import { Badge } from "@repo/ui/shared/badge";
import ExternalLink from "@repo/ui/shared/external-link";
import { externalLinks } from "@/src/routes";
import { EyeOpenIcon } from "@repo/ui/shared/react-icons";

export const columns: ColumnDef<Block>[] = [
  {
    accessorKey: "number",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Block" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            {Number(blockHeader?.number)}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "blockHeader.timestamp",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Age" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            {blockHeader?.timestamp
              ? formatTimeAgo(blockHeader?.timestamp)
              : "N/A"}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },

  {
    accessorKey: "rollupHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Rollup Hash" />
    ),
    cell: ({ row }) => {
      return Number(row.original.rollupHash) === 0 ? (
        <Badge>No rollup</Badge>
      ) : (
        <Link
          href={`/rollup/${row.original.rollupHash}`}
          className="text-primary"
        >
          <TruncatedAddress address={row.original.rollupHash} />
        </Link>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "blockHeader.gasUsed",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Rollup Gas Used" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            <Badge variant={"outline"}>
              {formatNumber(blockHeader?.gasUsed)}
            </Badge>
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "blockHeader.gasLimit",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Gas Limit" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            {formatNumber(blockHeader?.gasLimit)}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "blockHeader.hash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Hash" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return <TruncatedAddress address={blockHeader?.hash} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "blockHeader.parentHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Parent Hash" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return <TruncatedAddress address={blockHeader?.parentHash} />;
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
      const blockHeader = row.original.blockHeader as BlockHeader;
      return <TruncatedAddress address={blockHeader?.miner} />;
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    id: "actions",
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return (
        <ExternalLink
          href={`${externalLinks.etherscanBlock}${blockHeader?.hash}`}
        >
          <EyeOpenIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer mr-2" />
        </ExternalLink>
      );
    },
  },
];
