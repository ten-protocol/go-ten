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
import TruncatedAddress from "../common/truncated-address";
import { formatTimeAgo } from "@/src/lib/utils";
import { Batch } from "@/src/types/interfaces/BatchInterfaces";

export function RecentBatches({ batches }: any) {
  return (
    <div className="space-y-8">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Batch Hash</TableHead>
            <TableHead>L1 Block</TableHead>
            <TableHead>Txn Count</TableHead>
            <TableHead>Timestamp</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {batches?.result?.BatchesData.map((batch: Batch, i: number) => (
            <TableRow key={i}>
              <TableCell>{batch?.number}</TableCell>
              <TableCell className="font-medium">
                <TruncatedAddress address={batch.hash} />
              </TableCell>
              <TableCell className="font-medium">
                <TruncatedAddress address={batch.l1Proof} />
              </TableCell>
              <TableCell>{formatTimeAgo(batch?.timestamp)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
