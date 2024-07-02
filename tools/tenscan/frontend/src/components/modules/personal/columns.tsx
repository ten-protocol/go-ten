"use client";

import { ColumnDef } from "@tanstack/react-table";
import { Badge, badgeVariants } from "@/src/components/ui/badge";

import { statuses, types } from "./data";
import { DataTableColumnHeader } from "../common/data-table/data-table-column-header";
import { PersonalTransactions } from "@/src/types/interfaces/TransactionInterfaces";
import TruncatedAddress from "../common/truncated-address";
import { formatNumber } from "@/src/lib/utils";
import Link from "next/link";
import { EyeOpenIcon } from "@radix-ui/react-icons";

export const columns: ColumnDef<PersonalTransactions>[] = [
  {
    accessorKey: "blockNumber",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Ten Batch" />
    ),
    cell: ({ row }) => {
      return (
        <div className="flex space-x-2">
          <span className="max-w-[500px] truncate">
            #{Number(row.getValue("blockNumber"))}
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
      <DataTableColumnHeader column={column} title="Ten Batch Hash" />
    ),
    cell: ({ row }) => {
      return <TruncatedAddress address={row.getValue("blockHash")} />;
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
          <span className="max-w-[500px] truncate">
            <Badge variant={"secondary"}>
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
    accessorKey: "type",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Type" />
    ),
    cell: ({ row }) => {
      const type = types.find((type) => {
        return type.value === row.getValue("type");
      });

      if (!type) {
        return null;
      }

      return (
        <div className="flex items-center">
          <Badge variant={type.variant as keyof typeof badgeVariants}>
            {type.label}
          </Badge>
        </div>
      );
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id));
    },
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => {
      const status = statuses.find(
        (status) => status.value === row.getValue("status")
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
    cell: ({ row }) => {
      return (
        <Link href={`/tx/${row.original.transactionHash}`}>
          <EyeOpenIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer" />
        </Link>
      );
    },
  },
];
