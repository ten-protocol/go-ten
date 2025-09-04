import { useMemo, useState } from "react";
import { Separator } from "@repo/ui/components/shared/separator";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import KeyValueItem, {
  KeyValueList,
} from "@repo/ui/components/shared/key-value";
import {
  formatNumber,
  formatTimeAgo,
  formatTimestampToDate,
} from "@repo/ui/lib/utils";
import { Badge } from "@repo/ui/components/shared/badge";
import Link from "next/link";
import Copy from "@repo/ui/components/common/copy";
import { useRollupsService } from "@/src/services/useRollupsService";
import JSONPretty from "react-json-pretty";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@repo/ui/components/shared/tooltip";
import { pathToUrl } from "@/src/routes/router";
import { pageLinks } from "@/src/routes";
import {Batch} from "@/src/types/interfaces/BatchInterfaces";

export function BatchHashDetailsComponent({
  batchDetails,
}: {
  batchDetails: Batch;
}) {
  const { decryptedRollup, decryptEncryptedData } = useRollupsService();
  const [showDecryptedData, setShowDecryptedData] = useState(false);

  const transactionHashes = useMemo(
    () =>
      batchDetails?.txHashes.length > 0 ? (
        <ul>
          {batchDetails.txHashes.map((txHash, index) => (
            <li key={index} className="text-sm">
              <TruncatedAddress
                address={txHash}
                link={pathToUrl(pageLinks.txByHash, { hash: txHash })}
                showFullLength
              />
            </li>
          ))}
        </ul>
      ) : (
        "-"
      ),
    [batchDetails?.txHashes]
  );

  const handleDecryptToggle = () => {
    decryptEncryptedData({ StrData: batchDetails?.encryptedTxBlob });
    setShowDecryptedData(!showDecryptedData);
  };

  return (
    <div className="space-y-8">
      <KeyValueList>
        <KeyValueItem
          label="Height"
          value={
            <Link
              href={pathToUrl(pageLinks.batchByHeight, {
                height: +batchDetails?.header?.number,
              })}
              className="text-primary"
            >
              {"#" + Number(batchDetails?.header?.number)}
            </Link>
          }
        />
        <KeyValueItem
          label="Sequence"
          value={
            <Link
              href={pathToUrl(pageLinks.batchBySequence, {
                sequence: +batchDetails?.header?.sequencerOrderNo,
              })}
              className="text-primary"
            >
              {"#" + Number(batchDetails?.header?.sequencerOrderNo.toString())}
            </Link>
          }
        />
        <KeyValueItem
          label="Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.hash}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Parent Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.parentHash}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="State Root"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.stateRoot}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Transactions Root"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.transactionsRoot}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Receipts Root"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.receiptsRoot}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Timestamp"
          value={
            <Badge variant="secondary">
              {formatTimeAgo(batchDetails?.header?.timestamp) +
                " - " +
                formatTimestampToDate(batchDetails?.header?.timestamp)}
            </Badge>
          }
        />
        <KeyValueItem
          label="L1 Proof"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.l1Proof}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Gas Limit"
          value={
            <Badge variant="secondary">
              {formatNumber(batchDetails?.header?.gasLimit)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Gas Used"
          value={formatNumber(batchDetails?.header?.gasUsed)}
        />
        <KeyValueItem
          label="Base Fee Per Gas"
          value={
            <Badge variant="secondary">
              {formatNumber(batchDetails?.header?.baseFeePerGas)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Inbound Cross Chain Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.inboundCrossChainHash}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Inbound Cross Chain Height"
          value={Number(batchDetails?.header?.inboundCrossChainHeight)}
        />
        <KeyValueItem
          label="Miner"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.miner}
              showFullLength
            />
          }
        />

        <KeyValueItem
          label="Transaction Hashes"
          value={transactionHashes}
          isLastItem
        />
      </KeyValueList>

      <Separator />

      <KeyValueList>
        <KeyValueItem
          label="Encrypted Tx Blob"
          value={
            <>
              <div className="flex items-center space-x-2">
                <TruncatedAddress address={batchDetails?.encryptedTxBlob} />
                {showDecryptedData && decryptedRollup && (
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger>
                        <Copy value={decryptedRollup} />
                      </TooltipTrigger>
                      <TooltipContent>
                        Copy Decrypted Data to Clipboard
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                )}
              </div>
              {showDecryptedData && decryptedRollup && (
                <>
                  <Separator className="my-4" />
                  <JSONPretty id="json-pretty" data={decryptedRollup} />
                </>
              )}
            </>
          }
          isLastItem
        />
      </KeyValueList>
    </div>
  )
}
