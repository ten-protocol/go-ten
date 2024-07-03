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
import {ethMethods, tenCustomQueryMethods} from "../routes";

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

  const personalTransactions = async () => {
    try {
      setPersonalTxnsLoading(true);
      if (provider) {
        const requestPayload = {
          address: walletAddress,
          pagination: {
            ...options,
          },
        };
        const personalTxResp = await provider.send(ethMethods.getStorageAt, [
          tenCustomQueryMethods.listPersonalTransactions,
          JSON.stringify(requestPayload),
          null,
        ]);
        const personalTxData = jsonHexToObj(personalTxResp);
        setPersonalTxns(personalTxData);
      }
    } catch (error) {
      console.error("Error fetching personal transactions:", error);
      setPersonalTxns(undefined);
      showToast(ToastType.DESTRUCTIVE, "Error fetching personal transactions");
      throw error;
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

function jsonHexToObj(hex: string) {
  return JSON.parse(Buffer.from(hex.slice(2), "hex").toString());
}