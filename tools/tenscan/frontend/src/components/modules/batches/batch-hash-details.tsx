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
import {
  EyeClosedIcon,
  EyeOpenIcon,
} from "@repo/ui/components/shared/react-icons";
import { Button } from "@repo/ui/components/shared/button";
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
import { BatchDetails } from "@/src/types/interfaces/BatchInterfaces";

export function BatchHashDetailsComponent({
  batchDetails,
}: {
  batchDetails: BatchDetails;
}) {
  const { decryptedRollup, decryptEncryptedData } = useRollupsService();
  const [showDecryptedData, setShowDecryptedData] = useState(false);

  const transactionHashes = useMemo(
    () =>
      batchDetails?.TxHashes.length > 0 ? (
        <ul>
          {batchDetails.TxHashes.map((txHash, index) => (
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
    [batchDetails?.TxHashes]
  );

  const handleDecryptToggle = () => {
    decryptEncryptedData({ StrData: batchDetails?.EncryptedTxBlob });
    setShowDecryptedData(!showDecryptedData);
  };

  return (
    <div className="space-y-8">
      <KeyValueList>
        <KeyValueItem
          label="Batch Height"
          value={
            <Link
              href={pathToUrl(pageLinks.batchByHeight, {
                height: +batchDetails?.Header?.number,
              })}
              className="text-primary"
            >
              {"#" + Number(batchDetails?.Header?.number)}
            </Link>
          }
        />
        <KeyValueItem
          label="Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.hash}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Parent Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.parentHash}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="State Root"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.stateRoot}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Transactions Root"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.transactionsRoot}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Receipts Root"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.receiptsRoot}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Timestamp"
          value={
            <Badge variant="secondary">
              {formatTimeAgo(batchDetails?.Header?.timestamp) +
                " - " +
                formatTimestampToDate(batchDetails?.Header?.timestamp)}
            </Badge>
          }
        />
        <KeyValueItem
          label="L1 Proof"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.l1Proof}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Gas Limit"
          value={
            <Badge variant="secondary">
              {formatNumber(batchDetails?.Header?.gasLimit)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Gas Used"
          value={formatNumber(batchDetails?.Header?.gasUsed)}
        />
        <KeyValueItem
          label="Base Fee Per Gas"
          value={
            <Badge variant="secondary">
              {formatNumber(batchDetails?.Header?.baseFeePerGas)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Inbound Cross Chain Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.inboundCrossChainHash}
              showFullLength
            />
          }
        />
        <KeyValueItem
          label="Inbound Cross Chain Height"
          value={Number(batchDetails?.Header?.inboundCrossChainHeight)}
        />
        <KeyValueItem
          label="Miner"
          value={
            <TruncatedAddress
              address={batchDetails?.Header?.miner}
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
                <TruncatedAddress address={batchDetails?.EncryptedTxBlob} />
                <Button
                  className="text-sm font-bold leading-none hover:text-primary hover:bg-transparent"
                  variant="ghost"
                  onClick={handleDecryptToggle}
                >
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger>
                        {showDecryptedData ? (
                          <EyeClosedIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer ml-2" />
                        ) : (
                          <EyeOpenIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer ml-2" />
                        )}
                      </TooltipTrigger>
                      <TooltipContent>
                        {showDecryptedData
                          ? "Hide Encrypted Data"
                          : "Show Encrypted Data"}
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </Button>

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
  );
}
