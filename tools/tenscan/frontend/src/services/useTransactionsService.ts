import {
  fetchEtherPrice,
  fetchTransactions,
  fetchTransactionCount,
  personalTransactionsData,
} from "@/api/transactions";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import {
  getOptions,
  pollingInterval,
  pricePollingInterval,
} from "../lib/constants";
import { useRouter } from "next/router";
import useWalletStore from "@repo/ui/stores/wallet-store";

export const useTransactionsService = () => {
  const { query } = useRouter();
  const { address, provider } = useWalletStore();

  const [noPolling, setNoPolling] = useState(true);

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
    queryFn: () => personalTransactionsData(provider, address, options),
    enabled: !!address && !!provider,
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
