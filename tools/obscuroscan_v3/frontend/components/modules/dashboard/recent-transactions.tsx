import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import TruncatedAddress from "../common/truncated-address";

export function RecentTransactions({ transactions }: any) {
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
          {transactions?.result?.TransactionsData.map(
            (transaction: any, i: number) => (
              <TableRow key={i}>
                <TableCell className="font-medium">
                  {transaction?.BatchHeight}
                </TableCell>
                <TableCell>{transaction?.Finality}</TableCell>
                <TableCell>
                  <TruncatedAddress address={transaction?.TransactionHash} />
                </TableCell>
              </TableRow>
            )
          )}
        </TableBody>
      </Table>
    </div>
  );
}
