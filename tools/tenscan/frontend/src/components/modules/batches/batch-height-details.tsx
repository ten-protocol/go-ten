import { Separator } from "@/src/components/ui/separator";
import TruncatedAddress from "../common/truncated-address";
import KeyValueItem, { KeyValueList } from "@/src/components/ui/key-value";
import { formatNumber, formatTimeAgo } from "@/src/lib/utils";
import { Badge } from "@/src/components/ui/badge";
import { Batch, BatchDetails } from "@/src/types/interfaces/BatchInterfaces";
import Link from "next/link";
import { EyeClosedIcon, EyeOpenIcon } from "@radix-ui/react-icons";
import { Button } from "../../ui/button";
import { useRollupsService } from "@/src/services/useRollupsService";
import JSONPretty from "react-json-pretty";
import { useState } from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "../../ui/tooltip";

export function BatchHeightDetailsComponent({
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
          label="Batch Height"
          value={"#" + Number(batchDetails?.height)}
        />
        <KeyValueItem
          label="Sequence"
          value={"#" + Number(batchDetails?.sequence)}
        />
        <KeyValueItem
          label="Hash"
          value={
            <Link
              href={`/batch/${batchDetails?.header?.hash}`}
              className="text-primary"
            >
              <TruncatedAddress address={batchDetails?.header?.hash} />
            </Link>
          }
        />
        <KeyValueItem
          label="Parent Hash"
          value={
            <TruncatedAddress address={batchDetails?.header?.parentHash} />
          }
        />
        <KeyValueItem
          label="State Root"
          value={<TruncatedAddress address={batchDetails?.header?.stateRoot} />}
        />
        <KeyValueItem
          label="Transactions Root"
          value={
            <TruncatedAddress
              address={batchDetails?.header?.transactionsRoot}
            />
          }
        />
        <KeyValueItem
          label="Receipts Root"
          value={
            <TruncatedAddress address={batchDetails?.header?.receiptsRoot} />
          }
        />
        <KeyValueItem
          label="Timestamp"
          value={
            <Badge variant={"secondary"}>
              {formatTimeAgo(batchDetails?.header?.timestamp)}
            </Badge>
          }
        />
        <KeyValueItem
          label="L1 Proof"
          value={<TruncatedAddress address={batchDetails?.header?.l1Proof} />}
        />
        <KeyValueItem
          label="Gas Limit"
          value={formatNumber(batchDetails?.header?.gasLimit)}
        />
        <KeyValueItem
          label="Gas Used"
          value={formatNumber(batchDetails?.header?.gasUsed)}
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
            />
          }
        />
        <KeyValueItem
          label="Inbound Cross Chain Height"
          value={Number(batchDetails?.header?.inboundCrossChainHeight)}
        />
        <KeyValueItem
          label="Transfers Tree"
          value={
            <TruncatedAddress address={batchDetails?.header?.TransfersTree} />
          }
        />
        <KeyValueItem
          label="Miner"
          value={<TruncatedAddress address={batchDetails?.header?.miner} />}
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
            <Link
              href={`/batch/txs/${batchDetails?.header?.hash}`}
              className="text-primary"
            >
              {batchDetails?.txCount || "-"}
            </Link>
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
