import { Separator } from "@/src/components/ui/separator";
import TruncatedAddress from "../common/truncated-address";
import KeyValueItem, { KeyValueList } from "@/src/components/ui/key-value";
import { formatNumber, formatTimeAgo } from "@/src/lib/utils";
import { Badge } from "@/src/components/ui/badge";
import { BatchDetails } from "@/src/types/interfaces/BatchInterfaces";

export function BatchDetailsComponent({
  batchDetails,
}: {
  batchDetails: BatchDetails;
}) {
  return (
    <div className="space-y-8">
      <KeyValueList>
        <KeyValueItem
          label="Batch Number"
          value={"#" + Number(batchDetails?.Header?.number)}
        />
        <KeyValueItem
          label="Hash"
          value={<TruncatedAddress address={batchDetails?.Header?.hash} />}
        />
        <KeyValueItem
          label="Parent Hash"
          value={
            <TruncatedAddress address={batchDetails?.Header?.parentHash} />
          }
        />
        <KeyValueItem
          label="State Root"
          value={<TruncatedAddress address={batchDetails?.Header?.stateRoot} />}
        />
        <KeyValueItem
          label="Transactions Root"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.transactionsRoot}
            />
          }
        />
        <KeyValueItem
          label="Receipts Root"
          value={
            <TruncatedAddress address={batchDetails?.Header?.receiptsRoot} />
          }
        />
        <KeyValueItem
          label="Timestamp"
          value={
            <Badge variant={"secondary"}>
              {formatTimeAgo(batchDetails?.Header?.timestamp)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Coinbase"
          value={<TruncatedAddress address={batchDetails?.Header?.coinbase} />}
        />
        <KeyValueItem
          label="Gas Limit"
          value={formatNumber(batchDetails?.Header?.gasLimit)}
        />
        <KeyValueItem
          label="Gas Used"
          value={formatNumber(batchDetails?.Header?.gasUsed)}
        />
        <KeyValueItem
          label="Base Fee"
          value={batchDetails?.Header?.baseFee || "-"}
        />
        <KeyValueItem
          label="Inbound Cross Chain Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.inboundCrossChainHash}
            />
          }
        />
        <KeyValueItem
          label="Inbound Cross Chain Height"
          value={Number(batchDetails?.Header?.inboundCrossChainHeight)}
        />
        <KeyValueItem
          label="Transfers Tree"
          value={
            <TruncatedAddress address={batchDetails?.Header?.transfersTree} />
          }
        />
        <KeyValueItem
          label="Miner"
          value={<TruncatedAddress address={batchDetails?.Header?.miner} />}
        />
        <KeyValueItem
          label="Base Fee Per Gas"
          value={formatNumber(batchDetails?.Header?.baseFeePerGas)}
          isLastItem
        />
      </KeyValueList>
      <Separator />
      <KeyValueList>
        <KeyValueItem
          label="No. of Transactions"
          value={batchDetails?.TxHashes?.length || "-"}
          isLastItem
        />
      </KeyValueList>
      <Separator />
      <KeyValueList>
        <KeyValueItem
          label="Encrypted Tx Blob"
          value={<TruncatedAddress address={batchDetails?.EncryptedTxBlob} />}
          isLastItem
        />
      </KeyValueList>
    </div>
  );
}
