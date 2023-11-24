import {
  fetchEtherPrice,
  fetchTransactions,
  fetchTransactionCount,
} from "@/api/transactions";
import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { pollingInterval, pricePollingInterval } from "../lib/constants";
import { PersonalTransactionsResponse } from "../types/interfaces/TransactionInterfaces";
import { useToast } from "@/src/components/ui/use-toast";
import { useRouter } from "next/router";

export const useTransactionsService = () => {
  const router = useRouter();
  const { toast } = useToast();
  const { walletAddress, provider } = useWalletConnection();

  const [personalTxnsLoading, setPersonalTxnsLoading] = useState(false);
  const [personalTxns, setPersonalTxns] =
    useState<PersonalTransactionsResponse>();

  useEffect(() => {
    personalTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [walletAddress]);

  const [options, setOptions] = useState({
    offset: 0,
    limit: 10,
  });

  const [noPolling, setNoPolling] = useState(false);

  const updateQueryParams = (query: any) => {
    console.log(
      "🚀 ~ file: useBatchesService.ts:15 ~ updateQueryParams ~ query:",
      query
    );
    // setOptions({
    //   offset: newOffset,
    //   limit: newSize,
    // });
    // refetchBatches;
  };

  const {
    data: transactions,
    isLoading: isTransactionsLoading,
    refetch,
  } = useQuery({
    queryKey: ["transactions"],
    queryFn: () => fetchTransactions(),
    refetchInterval: pollingInterval,
  });

  const { data: transactionCount, isLoading: isTransactionCountLoading } =
    useQuery({
      queryKey: ["transactionCount"],
      queryFn: () => fetchTransactionCount(),
      refetchInterval: pollingInterval,
    });

  const personalTransactions = async (payload?: {
    pagination: { offset: number; size: number };
  }) => {
    try {
      setPersonalTxnsLoading(true);
      if (provider) {
        const requestPayload = {
          address: walletAddress,
          ...payload,
        };
        const personalTxData = await provider.send("eth_getStorageAt", [
          "listPersonalTransactions",
          requestPayload,
          null,
        ]);
        setPersonalTxns(personalTxData);
      }
    } catch (error) {
      toast({
        description: "Error fetching personal transactions",
      });
    } finally {
      setPersonalTxnsLoading(false);
    }
  };

  const { data: price, isLoading: isPriceLoading } = useQuery({
    queryKey: ["price"],
    queryFn: () => fetchEtherPrice(),
    refetchInterval: pricePollingInterval,
  });

  return {
    transactions,
    updateTransactionQueryParams: updateQueryParams,
    isTransactionsLoading,
    transactionCount,
    isTransactionCountLoading,
    personalTxns,
    personalTxnsLoading,
    price,
  };
};
