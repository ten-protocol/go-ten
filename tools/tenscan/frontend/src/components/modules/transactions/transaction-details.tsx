import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import KeyValueItem, {
  KeyValueList,
} from "@repo/ui/components/shared/key-value";
import { formatTimeAgo } from "@repo/ui/lib/utils";
import { BadgeType } from "@repo/ui/lib/enums/badge";
import { Badge } from "@repo/ui/components/shared/badge";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import Link from "next/link";

export function TransactionDetailsComponent({
  transactionDetails,
}: {
  transactionDetails: Transaction;
}) {
  return (
    <div className="space-y-8">
      <KeyValueList>
        <KeyValueItem
          label="Batch Height"
          value={
            <Link
              href={`/batch/height/${transactionDetails?.BatchHeight}`}
              className="text-primary"
            >
              {"#" + Number(transactionDetails?.BatchHeight)}
            </Link>
          }
        />
        <KeyValueItem
          label="Transaction Hash"
          value={
            <TruncatedAddress
              address={transactionDetails?.TransactionHash}
              link={`/tx/${transactionDetails?.TransactionHash}`}
            />
          }
        />
        <KeyValueItem
          label="Timestamp"
          value={formatTimeAgo(transactionDetails?.BatchTimestamp)}
        />
        <KeyValueItem
          label="Finality"
          value={
            <Badge
              variant={
                transactionDetails?.Finality === "Final"
                  ? BadgeType.SUCCESS
                  : BadgeType.DESTRUCTIVE
              }
            >
              {transactionDetails?.Finality}
            </Badge>
          }
          isLastItem
        />
      </KeyValueList>
    </div>
  );
}
