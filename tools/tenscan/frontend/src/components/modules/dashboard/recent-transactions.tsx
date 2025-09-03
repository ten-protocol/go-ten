import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import { Badge } from "@repo/ui/components/shared/badge";
import { formatTimeAgo } from "@repo/ui/lib/utils";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";
import GlitchTextAnimation from "@/src/components/GlitchTextAnimation";
import { RecentItemsList } from "./recent-items-list";

export function RecentTransactions({ transactions }: { transactions: any }) {
  const renderTransactionItem = (transaction: Transaction, isNewItem: boolean) => (
    <>
      <div className="ml-4 space-y-1 relative z-10">
        <p className="text-sm font-medium leading-none">
          <GlitchTextAnimation text={`#${Number(transaction?.BatchHeight)}`} hover={false} active={isNewItem} onView={false} />
        </p>
        <p className="text-sm text-muted-foreground word-break-all">
          <GlitchTextAnimation text={formatTimeAgo(transaction?.BatchTimestamp)} hover={false} active={isNewItem} onView={false} />
        </p>
      </div>
      
      <div className="ml-auto font-medium min-w-[140px] relative z-10" onClick={(e) => e.stopPropagation()}>
        <TruncatedAddress
          address={transaction?.TransactionHash}
          animate={isNewItem}
          AnimationComponent={GlitchTextAnimation}
          showPopover={false}
        />
      </div>
      <div className="ml-auto relative z-10">
        <Badge variant={"static-default"}>
          <GlitchTextAnimation text={transaction?.Finality} hover={false} active={isNewItem} onView={false} />
        </Badge>
      </div>
    </>
  );

  const headers = (
    <div className="flex items-center text-xs font-medium text-muted-foreground uppercase tracking-wide">
      <div className="ml-2 flex-1">
        <span>Batch</span>
      </div>
      <div className="min-w-[140px] text-center">
        <span>Hash</span>
      </div>
      <div className="min-w-[80px] text-right">
        <span>Status</span>
      </div>
    </div>
  );

  return (
    <RecentItemsList
      items={transactions?.result?.TransactionsData || []}
      getItemId={(transaction: Transaction) => transaction.TransactionHash}
      getItemLink={(transaction: Transaction) => pathToUrl(pageLinks.txByHash, { hash: transaction?.TransactionHash })}
      renderItem={renderTransactionItem}
      headers={headers}
    />
  );
}