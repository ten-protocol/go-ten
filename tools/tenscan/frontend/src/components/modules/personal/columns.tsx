"use client";
import React from "react";
import { ColumnDef } from "@tanstack/react-table";
import { Badge, badgeVariants } from "@repo/ui/components/shared/badge";

import { statuses, types } from "./data";
import { DataTableColumnHeader } from "@repo/ui/components/common/data-table/data-table-column-header";
import { PersonalTransactions } from "../../..//types/interfaces/TransactionInterfaces";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { formatNumber } from "@repo/ui/lib/utils";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";
import Link from "next/link";

export const columns: ColumnDef<PersonalTransactions>[] = [
  {
      accessorKey: "transactionHash",
      header: ({ column }) => (
          <DataTableColumnHeader column={column} title="Transaction Hash" />
      ),
      cell: ({ row }) => {
          return (
              <TruncatedAddress
                  address={row.getValue("transactionHash")}
                  link={pathToUrl(pageLinks.personalTxByHash, {
                      hash: row.original.transactionHash,
                  })}
              />
          );
      },
      enableSorting: false,
      enableHiding: false,
  },
  {
    accessorKey: "blockNumber",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="TEN Batch" />
    ),
    cell: ({ row }) => {
      return (
        <Link
          href={pathToUrl(pageLinks.batchByHeight, {
            height: Number(row.getValue("blockNumber")),
          })}
          className="text-primary"
        >
          <span className="max-w-[500px] truncate">
            #{Number(row.getValue("blockNumber"))}
          </span>
        </Link>
      );
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "blockHash",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="TEN Batch Hash" />
    ),
    cell: ({ row }) => {
      return (
        <TruncatedAddress
          address={row.getValue("blockHash")}
          link={pathToUrl(pageLinks.batchByHash, {
            hash: row.original.blockHash,
          })}
        />
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
];
