"use client";

import { Cross2Icon, ReloadIcon } from "../../shared/react-icons";
import { Table } from "@tanstack/react-table";

import { Button } from "../../shared/button";
import { DataTableViewOptions } from "./data-table-view-options";

import { DataTableFacetedFilter } from "./data-table-faceted-filter";

interface DataTableToolbarProps<TData> {
  table: Table<TData>;
  refetch?: () => void;
  toolbar?: {
    column: string;
    title: string;
    options: { label: string; value: string }[];
  }[];
}
export function DataTableToolbar<TData>({
  table,
  toolbar,
  refetch,
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0;

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-1 items-center space-x-2">
        {toolbar?.map(
          (item, index) =>
            table.getColumn(item.column) && (
              <DataTableFacetedFilter
                key={index}
                column={table.getColumn(item.column)}
                title={item.title}
                options={item.options}
              />
            )
        )}
        {isFiltered && (
          <Button
            variant="ghost"
            onClick={() => table.resetColumnFilters()}
            className="h-8 px-2 lg:px-3"
          >
            Reset
            <Cross2Icon className="ml-2 h-4 w-4" />
          </Button>
        )}
      </div>
      <div className="flex items-center space-x-2">
        {refetch && (
          <Button
            variant="ghost"
            onClick={refetch}
            className="h-8 px-2 lg:px-3"
          >
            <ReloadIcon className="mr-2 h-4 w-4" />
            Refresh
          </Button>
        )}
        <DataTableViewOptions table={table} />
      </div>
    </div>
  );
}
