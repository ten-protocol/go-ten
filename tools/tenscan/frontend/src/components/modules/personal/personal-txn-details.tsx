import React from "react";
import KeyValueItem, {
  KeyValueList,
} from "@repo/ui/components/shared/key-value";
import {
  PersonalTransactionType,
  TransactionReceipt,
  TransactionType,
} from "../../../types/interfaces/TransactionInterfaces";
import Link from "next/link";
import TruncatedAddress from "@repo/ui/components/common/truncated-address";
import { Badge } from "@repo/ui/components/shared/badge";
import { BadgeType } from "@repo/ui/lib/enums/badge";

export function PersonalTxnDetailsComponent({
  transactionDetails,
}: {
  transactionDetails: TransactionReceipt;
}) {
  const getTransactionType = (type: TransactionType) => {
    switch (type) {
      case PersonalTransactionType.Legacy:
        return "Legacy";
      case PersonalTransactionType.AccessList:
        return "Access List";
      case PersonalTransactionType.DynamicFee:
        return "Dynamic Fee";
      case PersonalTransactionType.Blob:
        return "Blob";
      default:
        return "Unknown";
    }
  };

  return (
    <div className="space-y-8">
      <KeyValueList>
        <KeyValueItem
          label="From"
          value={
            <TruncatedAddress
              address={transactionDetails?.from}
              link={`/address/${transactionDetails?.from}`}
            />
          }
        />
        <KeyValueItem
          label="To"
          value={
            <TruncatedAddress
              address={transactionDetails?.to}
              link={`/address/${transactionDetails?.to}`}
            />
          }
        />
        <KeyValueItem
          label="Transaction Index"
          value={
            <Badge variant={"outline"}>
              {Number(transactionDetails?.transactionIndex)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Transaction Type"
          value={
            <Badge
              variant={
                transactionDetails?.type
                  ? BadgeType.SUCCESS
                  : BadgeType.DESTRUCTIVE
              }
            >
              {getTransactionType(transactionDetails?.type)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Status"
          value={
            <Badge
              variant={
                transactionDetails?.status
                  ? BadgeType.SUCCESS
                  : BadgeType.DESTRUCTIVE
              }
            >
              {transactionDetails?.status ? "Success" : "Failed"}
            </Badge>
          }
        />

        <KeyValueItem
          label="Block Number"
          value={
            <Link href={`/block/${transactionDetails?.blockNumber}`}>
              {Number(transactionDetails?.blockNumber)}
            </Link>
          }
        />
        <KeyValueItem
          label="Gas Used"
          value={
            <Badge variant={"outline"}>
              {Number(transactionDetails?.gasUsed)}{" "}
            </Badge>
          }
        />
        <KeyValueItem
          label="Cumulative Gas Used"
          value={
            <Badge variant={"outline"}>
              {Number(transactionDetails?.cumulativeGasUsed)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Effective Gas Price"
          value={
            <Badge variant={"outline"}>
              {Number(transactionDetails?.effectiveGasPrice)}
            </Badge>
          }
        />
        <KeyValueItem
          label="Block Hash"
          value={<TruncatedAddress address={transactionDetails?.blockHash} />}
        />
        <KeyValueItem
          label="Logs Bloom"
          value={<TruncatedAddress address={transactionDetails?.logsBloom} />}
        />
        <KeyValueItem
          label="Contract Address"
          value={
            <TruncatedAddress
              address={transactionDetails?.contractAddress}
              link={`/address/${transactionDetails?.contractAddress}`}
            />
          }
        />
        <KeyValueItem
          label="Transaction Hash"
          value={
            <TruncatedAddress
              address={transactionDetails?.transactionHash}
              link={`/tx/${transactionDetails?.transactionHash}`}
            />
          }
        />

        <KeyValueItem
          label="Logs"
          value={
            transactionDetails?.logs.length > 0 ? (
              <div className="space-y-4">
                {transactionDetails?.logs.map((log, index) => (
                  <div key={index}>
                    <KeyValueList>
                      <KeyValueItem
                        label="Address"
                        value={<TruncatedAddress address={log.address} />}
                      />
                      <KeyValueItem
                        label="Block Hash"
                        value={<TruncatedAddress address={log.blockHash} />}
                      />
                      <KeyValueItem
                        label="Block Number"
                        value={Number(log.blockNumber)}
                      />
                      <KeyValueItem
                        label="Data"
                        value={<TruncatedAddress address={log.data} />}
                      />
                      <KeyValueItem label="Log Index" value={log.logIndex} />
                      <KeyValueItem
                        label="Removed"
                        value={
                          <Badge
                            variant={
                              log.removed
                                ? BadgeType.DESTRUCTIVE
                                : BadgeType.SUCCESS
                            }
                          >
                            {log.removed ? "Yes" : "No"}
                          </Badge>
                        }
                      />
                      <KeyValueItem
                        label="Topics"
                        value={
                          <div className="space-y-4">
                            {log.topics.map((topic, index) => (
                              <div key={index}>
                                <KeyValueItem
                                  value={<TruncatedAddress address={topic} />}
                                />
                              </div>
                            ))}
                          </div>
                        }
                      />
                      <KeyValueItem
                        label="Transaction Hash"
                        value={
                          <TruncatedAddress
                            address={log.transactionHash}
                            link={`/tx/${log.transactionHash}`}
                          />
                        }
                      />
                      <KeyValueItem
                        label="Transaction Index"
                        value={
                          <Badge variant={"outline"}>
                            {Number(transactionDetails?.transactionIndex)}
                          </Badge>
                        }
                        isLastItem
                      />
                    </KeyValueList>
                  </div>
                ))}
              </div>
            ) : (
              "No logs found"
            )
          }
          isLastItem
        />
      </KeyValueList>
    </div>
  );
}
