"use client";

import * as React from "react";
import { useRouter } from "next/router";
import {
  ColumnDef,
  ColumnFiltersState,
  OnChangeFn,
  PaginationState,
  SortingState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  getFacetedRowModel,
  getFacetedUniqueValues,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../shared/table";
import { DataTablePagination } from "./data-table-pagination";
import { DataTableToolbar } from "./data-table-toolbar";
import { Skeleton } from "../../shared/skeleton";
import { Button } from "../../shared/button";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  toolbar?: {
    column: string;
    title: string;
    options: { label: string; value: string }[];
  }[];
  updateQueryParams?: (query: any) => void;
  refetch?: () => void;
  total: number;
  isLoading?: boolean;
  noPagination?: boolean;
  noResultsText?: string;
  noResultsMessage?: string;
}

export function DataTable<TData, TValue>({
  columns,
  data,
  toolbar,
  refetch,
  total,
  isLoading,
  noPagination,
  noResultsText,
  noResultsMessage,
}: DataTableProps<TData, TValue>) {
  const { query, push, pathname } = useRouter();
  const [rowSelection, setRowSelection] = React.useState({});
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});
  const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
    []
  );
  const [sorting, setSorting] = React.useState<SortingState>([]);

  const pagination = React.useMemo(() => {
    return {
      pageIndex: Number(query.page) || 1,
      pageSize: Number(query.size) || 20,
    };
  }, [query.page, query.size]);

  const setPagination: OnChangeFn<PaginationState> = (func) => {
    const { pageIndex, pageSize } =
      typeof func === "function" ? func(pagination) : func;
    const newPageIndex = pagination.pageSize !== pageSize ? 1 : pageIndex;
    const params = {
      ...query,
      page: newPageIndex > 0 ? newPageIndex : 1,
      size: pageSize <= 100 ? pageSize : 100,
    };
    push({ pathname, query: params });
  };

  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      columnVisibility,
      rowSelection,
      columnFilters,
      pagination,
    },
    onPaginationChange: setPagination,
    manualPagination: true,
    // pageCount: Math.ceil(total / pagination.pageSize),
    enableRowSelection: true,
    onRowSelectionChange: setRowSelection,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFacetedRowModel: getFacetedRowModel(),
    getFacetedUniqueValues: getFacetedUniqueValues(),
  });

  return (
    <div className="space-y-4">
      {data && (
        <DataTableToolbar table={table} toolbar={toolbar} refetch={refetch} />
      )}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
                  return (
                    <TableHead key={header.id}>
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <>
                <TableRow>
                  <TableCell
                    colSpan={columns.length}
                    className="h-24 text-center"
                  >
                    <Skeleton className="w-full h-full" />
                  </TableCell>
                </TableRow>
              </>
            ) : data && table?.getRowModel()?.rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  {pagination.pageIndex > 1 ? (
                    <p>
                      No {noResultsText || "results"} found for the selected
                      filters.
                      <Button
                        variant={"link"}
                        onClick={() => {
                          setPagination({ pageIndex: 1, pageSize: 20 });
                          refetch?.();
                        }}
                      >
                        Clear Filters
                      </Button>
                    </p>
                  ) : (
                    <p>
                      {noResultsMessage ||
                        `No ${noResultsText || "results"} found.`}
                    </p>
                  )}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      {data && !isLoading && !noPagination && (
        <DataTablePagination
          table={table}
          refetch={refetch}
          setPagination={setPagination}
        />
      )}
    </div>
  );
}
