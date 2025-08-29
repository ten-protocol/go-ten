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
import { Batch } from "@/src/types/interfaces/BatchInterfaces";
import Link from "next/link";
import {
  EyeClosedIcon,
  EyeOpenIcon,
} from "@repo/ui/components/shared/react-icons";
import { Button } from "@repo/ui/components/shared/button";
import { useRollupsService } from "@/src/services/useRollupsService";
import JSONPretty from "react-json-pretty";
import React, {useMemo, useState} from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@repo/ui/components/shared/tooltip";
import { pageLinks } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";

export function BatchDetailsComponent({
  batchDetails,
}: {
  batchDetails: Batch;
}) {
  const { decryptedRollup, decryptEncryptedData } = useRollupsService();
  const [showDecryptedData, setShowDecryptedData] = useState(false);
  return (
    <div className="space-y-8">
      <KeyValueList>
        <KeyValueItem
          label="Height"
          value={"#" + Number(batchDetails?.height)}
        />
        <KeyValueItem
          label="Sequence"
          value={"#" + Number(batchDetails?.sequence)}
        />
        <KeyValueItem
          label="Hash"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.hash}
              link={pathToUrl(pageLinks.batchByHash, {
                hash: batchDetails?.header?.hash,
              })}
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
            <Badge variant={"secondary"}>
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
            <Badge variant={"outline"}>
              {formatNumber(batchDetails?.header?.gasLimit)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Gas Used"
          value={
            <Badge variant={"outline"}>
              {formatNumber(batchDetails?.header?.gasUsed)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Base Fee Per Gas"
          value={batchDetails?.header?.baseFeePerGas || "-"}
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
          label="Base Fee Per Gas"
          value={formatNumber(batchDetails?.header?.baseFeePerGas)}
          isLastItem
        />
      </KeyValueList>
      <Separator />
      <KeyValueList>
        <KeyValueItem
          label="No. of Transactions"
          value={
            batchDetails?.txHashes && batchDetails.txHashes.length > 0 ? (
              <span>
                {batchDetails.txHashes.length}{" "}
                <Link
                  href={pathToUrl(pageLinks.batchTransactions, {
                    hash: batchDetails?.header?.hash,
                  })}
                  className="underline text-primary"
                >
                  View
                </Link>
              </span>
            ) : (
              "0"
            )
          }
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
                <TruncatedAddress address={batchDetails?.encryptedTxBlob} />{" "}
                <Button
                  className="text-sm font-bold leading-none hover:text-primary hover:bg-transparent"
                  variant="ghost"
                  onClick={() => {
                    decryptEncryptedData({
                      StrData: batchDetails?.encryptedTxBlob,
                    });
                    setShowDecryptedData(!showDecryptedData);
                  }}
                >
                  {showDecryptedData && decryptedRollup ? (
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger>
                          <EyeClosedIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer ml-2" />
                        </TooltipTrigger>
                        <TooltipContent>Hide Encrypted Data</TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  ) : (
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger>
                          <EyeOpenIcon className="h-5 w-5 text-muted-foreground hover:text-primary transition-colors cursor-pointer ml-2" />
                        </TooltipTrigger>
                        <TooltipContent>Show Encrypted Data</TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  )}
                </Button>
              </div>
              {decryptedRollup && showDecryptedData ? (
                <>
                  <Separator className="my-4" />
                  <JSONPretty
                    id="json-pretty"
                    data={decryptedRollup}
                  ></JSONPretty>
                </>
              ) : null}
            </>
          }
          isLastItem
        />
      </KeyValueList>
    </div>
  );
}
