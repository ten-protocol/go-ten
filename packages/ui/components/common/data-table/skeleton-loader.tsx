import { Skeleton } from "../../shared/skeleton";
import { TableHeader, TableRow, TableBody } from "../../shared/table";
import { Table } from "lucide-react";

const DataTableSkeleton = ({ columns }: { columns: number }) => {
  const renderSkeletonColumns = () => {
    return Array.from({ length: columns }).map((_, index) => (
      <Skeleton key={index} className="w-[100px] h-[20px] rounded-full" />
    ));
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>{renderSkeletonColumns()}</TableRow>
      </TableHeader>
      <TableBody>{renderSkeletonColumns()}</TableBody>
    </Table>
  );
};

export default DataTableSkeleton;
