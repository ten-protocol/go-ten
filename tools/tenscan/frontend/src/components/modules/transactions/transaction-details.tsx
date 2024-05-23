import TruncatedAddress from "../common/truncated-address";
import KeyValueItem, { KeyValueList } from "@/src/components/ui/key-value";
import { formatTimeAgo } from "@/src/lib/utils";
import { Badge } from "@/src/components/ui/badge";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import { BadgeType } from "@/src/types/interfaces";
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
            <Link
              href={`/tx/${transactionDetails?.BatchHeight}`}
              className="text-primary"
            >
              <TruncatedAddress address={transactionDetails?.TransactionHash} />
            </Link>
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
