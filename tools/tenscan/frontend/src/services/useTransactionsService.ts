import {
  fetchEtherPrice,
  fetchTransactions,
  fetchTransactionCount,
  personalTransactionsData,
} from "@/api/transactions";
import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import {
  getOptions,
  pollingInterval,
  pricePollingInterval,
} from "../lib/constants";
import { useRouter } from "next/router";

export const useTransactionsService = () => {
  const { query } = useRouter();
  const { walletAddress, provider } = useWalletConnection();

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

  const { data: personalTxns, isLoading: personalTxnsLoading } = useQuery({
    queryKey: ["personalTxns", options],
    queryFn: () => personalTransactionsData(provider, walletAddress, options),
    enabled: !!walletAddress && !!provider,
  });

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
    isPriceLoading,
  };
};
