import TruncatedAddress from "../common/truncated-address";
import { Avatar, AvatarFallback } from "@/src/components/ui/avatar";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import { Badge } from "../../ui/badge";
import { formatTimeAgo } from "@/src/lib/utils";
import Link from "next/link";

export function RecentTransactions({ transactions }: { transactions: any }) {
  return (
    <div className="space-y-8">
      {transactions?.result?.TransactionsData.map(
        (transaction: Transaction, i: number) => (
          <div className="flex items-center" key={i}>
            <Avatar className="h-9 w-9">
              <AvatarFallback>TX</AvatarFallback>
            </Avatar>
            <div className="ml-4 space-y-1">
              <p className="text-sm font-medium leading-none">
                <Link
                  href={`/batch/height/${transaction?.BatchHeight}`}
                  className="text-primary"
                >
                  #{Number(transaction?.BatchHeight)}
                </Link>
              </p>
              <p className="text-sm text-muted-foreground word-break-all">
                {formatTimeAgo(transaction?.BatchTimestamp)}
              </p>
            </div>
            <div className="ml-auto font-medium">
              <Link
                href={`/tx/${transaction?.TransactionHash}`}
                className="text-primary"
              >
                {" "}
                <TruncatedAddress address={transaction?.TransactionHash} />
              </Link>
            </div>
            <div className="ml-auto">
              <Badge>{transaction?.Finality}</Badge>
            </div>
          </div>
        )
      )}
    </div>
  );
}
