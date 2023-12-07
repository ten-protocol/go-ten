import {
  fetchEtherPrice,
  fetchTransactions,
  fetchTransactionCount,
} from "@/api/transactions";
import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import {
  getOptions,
  pollingInterval,
  pricePollingInterval,
} from "../lib/constants";
import { PersonalTransactionsResponse } from "../types/interfaces/TransactionInterfaces";
import { useRouter } from "next/router";
import { showToast } from "../components/ui/use-toast";
import { ToastType } from "../types/interfaces";
import { ethMethods } from "../routes";

export const useTransactionsService = () => {
  const { query } = useRouter();
  const { walletAddress, provider } = useWalletConnection();

  const [personalTxnsLoading, setPersonalTxnsLoading] = useState(false);
  const [personalTxns, setPersonalTxns] =
    useState<PersonalTransactionsResponse>();

  useEffect(() => {
    personalTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [walletAddress]);

  const [noPolling, setNoPolling] = useState(false);

  const options = getOptions(query);

  const {
    data: transactions,
    isLoading: isTransactionsLoading,
    refetch: refetchTransactions,
  } = useQuery({
    queryKey: ["transactions", options],
    queryFn: () => fetchTransactions(options),
    refetchInterval: noPolling ? false : pollingInterval,
  });

  const { data: transactionCount, isLoading: isTransactionCountLoading } =
    useQuery({
      queryKey: ["transactionCount"],
      queryFn: () => fetchTransactionCount(),
      refetchInterval: noPolling ? false : pollingInterval,
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
        const personalTxData = await provider.send(ethMethods.getStorageAt, [
          "listPersonalTransactions",
          requestPayload,
          null,
        ]);
        console.log(
          "🚀 ~ file: useTransactionsService.ts:68 ~ useTransactionsService ~ personalTxData:",
          personalTxData
        );
        setPersonalTxns(personalTxData);
      }
    } catch (error) {
      showToast(ToastType.DESTRUCTIVE, "Error fetching personal transactions");
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
    isTransactionsLoading,
    refetchTransactions,
    setNoPolling,
    transactionCount,
    isTransactionCountLoading,
    personalTxns,
    personalTxnsLoading,
    price,
  };
};
