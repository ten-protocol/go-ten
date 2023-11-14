"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge, badgeVariants } from "@/components/ui/badge";

import { statuses } from "./data";
import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import { DataTableRowActions } from "../common/data-table/data-table-row-actions";
import { PersonalTransactions } from "@/src/types/interfaces/TransactionInterfaces";
import TruncatedAddress from "../common/truncated-address";

export const columns: ColumnDef<PersonalTransactions>[] = [
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
    accessorKey: "blockNumber",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Block Number" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {row.getValue("blockNumber")}
          </span>
        </div>
      );
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
    accessorKey: "blockHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Block Hash" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {row.getValue("blockHash")}
          </span>
        </div>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => {
      const status = statuses.find(
        (status) => status.value === row.getValue("Status")
      );

      if (!status) {
        return null;
      }

      return (
        <div className="flex items-center">
          {status.icon && (
            <status.icon className="mr-2 h-4 w-4 text-muted-foreground" />
          )}
          <Badge variant={status.variant as keyof typeof badgeVariants}>
            {status.label}
          </Badge>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
  {
    id: "actions",
    cell: ({ row }) => <DataTableRowActions row={row} labels={statuses} />,
  },
];
