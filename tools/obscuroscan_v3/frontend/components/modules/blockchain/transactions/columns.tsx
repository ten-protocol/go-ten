"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";

import { labels, statuses } from "./data";
import { DataTableColumnHeader } from "../../common/data-table/data-table-column-header";
import { DataTableRowActions } from "../../common/data-table/data-table-row-actions";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import TruncatedAddress from "../../common/truncated-address";

export const columns: ColumnDef<Transaction>[] = [
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
    accessorKey: "BatchHeight",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Batch" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate font-medium">
            {row.getValue("BatchHeight")}
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
      return <TruncatedAddress address={row.getValue("TransactionHash")} />;
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
