import { Button } from "@/components/ui/button";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";
import { formatTimeAgo } from "@/src/lib/utils";
import { Block } from "@/src/types/interfaces/BlockInterfaces";
import TruncatedAddress from "../common/truncated-address";

export function RecentBlocks({ blocks }: any) {
  return (
    <div className="space-y-8">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[100px]">Block</TableHead>
            <TableHead>Age</TableHead>
            <TableHead>Tx Hash</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {blocks?.result?.BlocksData.map((block: Block, i: number) => (
            <TableRow key={i}>
              <TableCell>{block?.blockHeader?.number}</TableCell>
              <TableCell className="font-medium">
                {formatTimeAgo(block?.blockHeader?.timestamp)}
              </TableCell>
              <TableCell>
                <TruncatedAddress address={block?.blockHeader?.hash} />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
