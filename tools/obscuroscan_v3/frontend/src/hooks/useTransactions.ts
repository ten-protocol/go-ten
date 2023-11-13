import {
  getTransactions,
  getTransactionCount,
  getPrice,
} from "@/api/transactions";
import { useWalletConnection } from "@/components/providers/wallet-provider";
import { useQuery } from "@tanstack/react-query";
import { ethers } from "ethers";
import { useEffect, useState } from "react";

export const useTransactions = () => {
  const { walletAddress, provider } = useWalletConnection();
  const [personalTxns, setPersonalTxns] = useState<Uint8Array>();

  useEffect(() => {
    personalTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [walletAddress]);

  const { data: transactions, isLoading: isTransactionsLoading } = useQuery({
    queryKey: ["transactions"],
    queryFn: () => getTransactions(),
  });

  const { data: transactionCount, isLoading: isTransactionCountLoading } =
    useQuery({
      queryKey: ["transactionCount"],
      queryFn: () => getTransactionCount(),
    });

  const personalTransactions = async () => {
    if (provider) {
      const personalTxData = await provider.send("eth_getStorageAt", [
        walletAddress,
        "0x0",
        null,
      ]);
      const personalTx = ethers.utils.arrayify(personalTxData);
      setPersonalTxns(personalTx);
    }
  };

  const { data: price, isLoading: isPriceLoading } = useQuery({
    queryKey: ["price"],
    queryFn: () => getPrice(),
  });

  return {
    transactions,
    isTransactionsLoading,
    transactionCount,
    isTransactionCountLoading,
    personalTxns,
    price,
  };
};
