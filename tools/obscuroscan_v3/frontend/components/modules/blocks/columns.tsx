"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Checkbox } from "@/components/ui/checkbox";

import { labels, statuses } from "./data";
import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import { DataTableRowActions } from "../common/data-table/data-table-row-actions";
import { Block, BlockHeader } from "@/src/types/interfaces/BlockInterfaces";
import TruncatedAddress from "../common/truncated-address";
import { formatTimeAgo } from "@/src/lib/utils";

export const columns: ColumnDef<Block>[] = [
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
      <DataTableColumnHeader column={column} title="Block" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {+blockHeader?.number}
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
          <span className="max-w-[500px] truncate font-medium">
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
    accessorKey: "blockHeader.gasUsed",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Gas Used" />
    ),
    cell: ({ row }) => {
      const blockHeader = row.original.blockHeader as BlockHeader;
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {+blockHeader?.gasUsed}
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
          <span className="max-w-[500px] truncate font-medium">
            {+blockHeader?.gasLimit}
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
    accessorKey: "rollupHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Rollup Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.original.rollupHash} />;
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
          <span>{finality.label}</span>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
  {
    id: "actions",
    cell: ({ row }) => <DataTableRowActions row={row} labels={labels} />,
  },
];
