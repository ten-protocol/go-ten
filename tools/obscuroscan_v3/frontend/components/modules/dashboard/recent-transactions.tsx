import { useTransactions } from "@/src/hooks/useTransactions";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export function RecentTransactions() {
  const { transactions } = useTransactions();

  return (
    <div className="space-y-8">
      <Table>
        <TableCaption>A list of your recent invoices.</TableCaption>
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
                <TableCell>{transaction?.TransactionHash}</TableCell>
              </TableRow>
            )
          )}
        </TableBody>
      </Table>

      <Button>View All Transactions</Button>
    </div>
  );
}
