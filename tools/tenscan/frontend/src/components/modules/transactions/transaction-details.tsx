import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import KeyValueItem, {
  KeyValueList,
} from "@repo/ui/components/shared/key-value";
import { formatTimeAgo, formatTimestampToDate } from "@repo/ui/lib/utils";
import { BadgeType } from "@repo/ui/lib/enums/badge";
import { Badge } from "@repo/ui/components/shared/badge";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import Link from "next/link";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";

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
              href={pathToUrl(pageLinks.batchByHeight, {
                height: transactionDetails?.BatchHeight,
              })}
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
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Timestamp"
          value={
            <Badge variant="outline">
              {formatTimeAgo(transactionDetails?.BatchTimestamp) +
                " - " +
                formatTimestampToDate(transactionDetails?.BatchTimestamp)}
            </Badge>
          }
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
