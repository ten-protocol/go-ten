import { Button } from "@/components/ui/button";
import TruncatedAddress from "../common/truncated-address";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { EyeOpenIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";

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
                Height: #{transaction?.BatchHeight}
              </p>
            </div>
            <div className="ml-auto font-medium">
              <TruncatedAddress address={transaction?.TransactionHash} />
            </div>
            <div className="ml-auto font-medium">
              <Link
                href={{
                  pathname: `/transactions/${transaction?.TransactionHash}`,
                }}
              >
                <Button variant="link" size="sm">
                  <EyeOpenIcon />
                </Button>
              </Link>
            </div>
          </div>
        )
      )}
    </div>
  );
}
