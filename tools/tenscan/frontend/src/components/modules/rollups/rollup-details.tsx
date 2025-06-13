import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import KeyValueItem, {
  KeyValueList,
} from "@repo/ui/components/shared/key-value";
import { formatTimeAgo, formatTimestampToDate } from "@repo/ui/lib/utils";
import Link from "next/link";
import { Rollup } from "@/src/types/interfaces/RollupInterfaces";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";
import { Badge } from "@repo/ui/components/shared/badge";

export function RollupDetailsComponent({
  rollupDetails,
}: {
  rollupDetails: Rollup;
}) {
  return (
    <div className="space-y-8">
      <KeyValueList>
        <KeyValueItem label="ID" value={"#" + Number(rollupDetails?.ID)} />
        <KeyValueItem
          label="Timestamp"
          value={
            <Badge variant="outline">
              {formatTimeAgo(rollupDetails?.Timestamp) +
                " - " +
                formatTimestampToDate(rollupDetails?.Timestamp)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Full Hash"
          value={
            <TruncatedAddress
              address={rollupDetails?.Hash}
              link={pathToUrl(pageLinks.rollupByHash, {
                hash: rollupDetails?.Hash,
              })}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Rollup Header Hash"
          value={
            <TruncatedAddress
              address={rollupDetails?.Header?.hash}
              link={pathToUrl(pageLinks.rollupByHash, {
                hash: rollupDetails?.Header?.hash,
              })}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="L1 Hash"
          value={
            <TruncatedAddress address={rollupDetails?.L1Hash} showFullLength />
          }
        />
        <KeyValueItem
          label="First Batch Seq No."
          value={
            <Link
              href={pathToUrl(pageLinks.rollupByBatchSequence, {
                sequence: rollupDetails?.FirstSeq,
              })}
              className="text-primary"
            >
              {"#" + rollupDetails?.FirstSeq}
            </Link>
          }
        />
        <KeyValueItem
          label="Last Batch Seq No."
          value={
            <Link
              href={pathToUrl(pageLinks.rollupByBatchSequence, {
                sequence: rollupDetails?.LastSeq,
              })}
              className="text-primary"
            >
              {"#" + rollupDetails?.LastSeq}
            </Link>
          }
        />
        <KeyValueItem
          label="Compression L1 Head"
          value={
            <TruncatedAddress
              address={rollupDetails?.Header?.CompressionL1Head}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Payload Hash"
          value={
            <TruncatedAddress
              address={rollupDetails?.Header?.PayloadHash}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Signature"
          value={
            <TruncatedAddress address={rollupDetails?.Header?.Signature} />
          }
        />
        <KeyValueItem
          label="Cross Chain Messages"
          value={
            rollupDetails?.Header?.crossChainMessages?.length > 0
              ? rollupDetails?.Header?.crossChainMessages?.map((msg, index) => (
                  <div key={index} className="space-y-4">
                    <KeyValueList>
                      <KeyValueItem label="Sender" value={msg.Sender} />
                      <KeyValueItem
                        label="Sequence"
                        value={
                          <Link
                            href={pathToUrl(pageLinks.rollupByBatchSequence, {
                              sequence: msg?.Sequence,
                            })}
                            className="text-primary"
                          >
                            {"#" + msg.Sequence}
                          </Link>
                        }
                      />
                      <KeyValueItem label="Nonce" value={msg.Nonce} />
                      <KeyValueItem label="Topic" value={msg.Topic} />
                      <KeyValueItem
                        label="Payload"
                        value={msg.Payload.map((payload, index) => (
                          <div key={index}>{payload}</div>
                        ))}
                      />
                      <KeyValueItem
                        label="Consistency Level"
                        value={msg.ConsistencyLevel}
                      />
                    </KeyValueList>
                  </div>
                ))
              : "No cross chain messages found."
          }
          isLastItem
        />
      </KeyValueList>
    </div>
  );
}
