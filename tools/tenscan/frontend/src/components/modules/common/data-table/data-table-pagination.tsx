import {
  ChevronLeftIcon,
  ChevronRightIcon,
  DoubleArrowLeftIcon,
} from "@radix-ui/react-icons";
import { PaginationState, Table } from "@tanstack/react-table";
import { Button } from "@/src/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/src/components/ui/select";
import { Input } from "@/src/components/ui/input";
import { useState } from "react";

interface DataTablePaginationProps<TData> {
  table: Table<TData>;
  refetch?: () => void;
  setPagination: (pagination: PaginationState) => void;
}

export function DataTablePagination<TData>({
  table,
  refetch,
  setPagination,
}: DataTablePaginationProps<TData>) {
  const [page, setPage] = useState(table.getState().pagination.pageIndex);

  const handlePageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPage(Number(e.target.value));
  };

  const handleKey = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (
      e.key === "Enter" &&
      page > 0 &&
      page !== table.getState().pagination.pageIndex
    ) {
      table.setPageIndex(page);
      refetch?.();
    }
  };

  return (
    <div className="flex items-center flex-wrap justify-between space-x-2">
      <div className="flex-1 text-sm text-muted-foreground mb-2">
        Showing {table.getFilteredRowModel().rows.length} row(s)
      </div>
      <div className="flex flex-2 gap-1 items-center justify-between flex-wrap space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium">Rows per page</p>
          <Select
            value={`${table.getState().pagination.pageSize}`}
            onValueChange={(value: string) => {
              table.setPageSize(Number(value));
              setPagination({ pageIndex: 1, pageSize: Number(value) });
            }}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={table.getState().pagination.pageSize} />
            </SelectTrigger>
            <SelectContent side="top">
              {[5, 10, 20, 30, 40, 50, 100].map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          <span>Page</span>
          <Input
            className="w-[70px] h-8 text-center mx-2 text-ellipsis"
            type="number"
            value={page}
            onChange={handlePageChange}
            onKeyDown={handleKey}
            min={1}
            onFocus={(e) => e.target.select()}
            onBlur={() => setPage(table.getState().pagination.pageIndex)}
          />
          {/* uncomment the following line when total count feature is implemented */}
          {/* of {formatNumber(table.getPageCount())} */}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => {
              setPage(1);
              table.setPageIndex(1);
            }}
            disabled={table.getState().pagination.pageIndex === 1}
          >
            <span className="sr-only">Go to first page</span>
            <DoubleArrowLeftIcon className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => {
              setPage(table.getState().pagination.pageIndex - 1);
              table.previousPage();
            }}
            disabled={table.getState().pagination.pageIndex === 1}
          >
            <span className="sr-only">Go to previous page</span>
            <ChevronLeftIcon className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => {
              setPage(table.getState().pagination.pageIndex + 1);
              table.nextPage();
            }}
            // uncomment the following line when total count feature is implemented
            // disabled={!table.getCanNextPage()}
          >
            <span className="sr-only">Go to next page</span>
            <ChevronRightIcon className="h-4 w-4" />
          </Button>
          {/* uncomment the following line when total count feature is implemented */}
          {/* <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => table.setPageIndex(table.getPageCount() - 1)}
            // disabled={!table.getCanNextPage()}
          >
            <span className="sr-only">Go to last page</span>
            <DoubleArrowRightIcon className="h-4 w-4" />
          </Button> */}
        </div>
      </div>
    </div>
  );
}
