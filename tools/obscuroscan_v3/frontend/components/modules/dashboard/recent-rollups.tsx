import { Button } from "@/components/ui/button";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";

export function RecentRollups({ rollups }: any) {
  return (
    <div className="space-y-8">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Batch Height</TableHead>
            <TableHead>Finality</TableHead>
            <TableHead className="text-right">Tx Hash</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rollups?.result?.RollupsData.map((rollup: any, i: number) => (
            <TableRow key={i}>
              <TableCell className="font-medium">
                {rollup?.BatchHeight}
              </TableCell>
              <TableCell>{rollup?.Finality}</TableCell>
              <TableCell>{rollup?.TransactionHash}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Button>View All Rollups</Button>
    </div>
  );
}
