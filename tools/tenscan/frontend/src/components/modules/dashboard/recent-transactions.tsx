import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { Avatar, AvatarFallback } from "@repo/ui/components/shared/avatar";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import { Badge } from "@repo/ui/components/shared/badge";
import { formatTimeAgo } from "@repo/ui/lib/utils";
import Link from "next/link";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";

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
                <span className="text-muted-foreground">Batch </span>
                <Link
                  href={pathToUrl(pageLinks.batchByHeight, {
                    height: transaction?.BatchHeight,
                  })}
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
              <TruncatedAddress
                address={transaction?.TransactionHash}
                link={pathToUrl(pageLinks.txByHash, {
                  hash: transaction?.TransactionHash,
                })}
              />
            </div>
            <div className="ml-auto">
              <Badge variant={"static-default"}>{transaction?.Finality}</Badge>
            </div>
          </div>
        )
      )}
    </div>
  );
}
