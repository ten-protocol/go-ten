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

export const useTransactionsService = () => {
  const { toast } = useToast();
  const { walletAddress, provider } = useWalletConnection();

  const [personalTxnsLoading, setPersonalTxnsLoading] = useState(false);
  const [personalTxns, setPersonalTxns] =
    useState<PersonalTransactionsResponse>();

  useEffect(() => {
    personalTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [walletAddress]);

  const { data: transactions, isLoading: isTransactionsLoading } = useQuery({
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
    isTransactionsLoading,
    transactionCount,
    isTransactionCountLoading,
    personalTxns,
    personalTxnsLoading,
    price,
  };
};
