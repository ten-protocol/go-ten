import { Button } from "@/components/ui/button";
import {
  Table,
  TableCaption,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";
import { format } from "date-fns";
import TruncatedAddress from "../common/truncated-address";

export function RecentBatches({ batches }: any) {
  return (
    <div className="space-y-8">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Batch Hash</TableHead>
            <TableHead>L1 Block</TableHead>
            <TableHead>Txn Count</TableHead>
            <TableHead className="text-right">Timestamp</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {batches?.result?.BatchesData.map((batch: any, i: number) => (
            <TableRow key={i}>
              <TableCell className="font-medium">
                <TruncatedAddress address={batch.hash} />
              </TableCell>
              <TableCell className="font-medium">
                <TruncatedAddress address={batch.l1Proof} />
              </TableCell>
              <TableCell>{batch?.txHashes?.length || 0}</TableCell>
              <TableCell>{format(batch?.timestamp, "LLL dd, y")}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Button>View All Rollups</Button>
    </div>
  );
}
