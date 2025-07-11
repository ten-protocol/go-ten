import { useQuery } from "@tanstack/react-query";
import { useContractsService } from "@/src/services/useContractsService";
import { balancePollingInterval, CHAINS, TOKENS } from "@/src/lib/constants";
import { IToken, ToastType } from "../types";
import { showToast, toast } from "../components/ui/use-toast";
import React from "react";
import useWalletStore from "../stores/wallet-store";

export const useBridgeUtils = () => {
  const { getNativeBalance, getTokenBalance, sendERC20, sendNative } =
    useContractsService();

  const { switchNetwork } = useWalletStore();

  const useTokenBalance = (
    tokens: IToken[],
    token: string,
    address: string,
    fromChain: string,
    walletConnected: boolean,
    loading: boolean,
  ) => {
    const fetchTokenBalance = async (token: string, address: string) => {
      if (!token || !address) return null;

      const selectedToken = tokens.find((t) => t.value === token);
      if (!selectedToken) return null;

      return selectedToken.isNative
        ? await getNativeBalance(address)
        : await getTokenBalance(selectedToken.address, address);
    };

    const {
      data,
      isLoading: isBalanceLoading,
      isFetching: isBalanceFetching,
      refetch: refreshBalance,
    } = useQuery({
      queryKey: ["tokenBalance", token, address, fromChain],
      queryFn: () => fetchTokenBalance(token, address),
      enabled:
        !!token && !!address && !!fromChain && walletConnected && !loading,
      refetchInterval: balancePollingInterval,
    });

    const tokenBalance = (data || 0.0) as number;

    return {
      tokenBalance,
      isBalanceLoading,
      isBalanceFetching,
      refreshBalance,
    };
  };

  const useBridgeTransaction = (
    address: string,
    token: string,
    tokens: IToken[],
    receiver: string,
    tokenBalance: number,
    setError: (name: string, error: { type: string; message: string }) => void,
  ) => {
    const initiateBridgeTransaction = React.useCallback(
      async (data: any) => {
        const amount = Number(data.amount);
        if (isNaN(amount) || amount <= 0) {
          setError("amount", {
            type: "manual",
            message: "Please enter a valid amount.",
          });
          return;
        }
        if (amount > tokenBalance) {
          setError("amount", {
            type: "manual",
            message: "Amount must be less than or equal to your balance.",
          });
          return;
        }
        try {
          const transactionData = { ...data, receiver: receiver || address };
          if (!transactionData.receiver) {
            throw new Error("Receiver address is required.");
          }

          const selectedToken = tokens.find((t: IToken) => t.value === token);

          if (!selectedToken) {
            setError("token", {
              type: "manual",
              message: "Selected token is invalid.",
            });
            return;
          }

          toast({
            description: "Bridge transaction initiated.",
            variant: ToastType.INFO,
          });

          let res;
          if (selectedToken.isNative) {
            res = await sendNative({
              receiver: transactionData.receiver,
              value: amount.toString(),
            });
          } else {
            res = await sendERC20(
              transactionData.receiver,
              amount,
              selectedToken.address,
            );
          }

          if (res?.transactionHash) {
            showToast(ToastType.SUCCESS, "Transaction completed successfully");
          }
        } catch (error: any) {
          console.error(error);
          setError("submit", {
            type: "manual",
            message: error.message || "An unexpected error occurred.",
          });
          toast({
            title: "Transaction Failed",
            description:
              error.message || "An error occurred during the transaction.",
            variant: ToastType.DESTRUCTIVE,
          });
        }
      },
      [address, token, tokens, receiver, tokenBalance],
    );

    return { initiateBridgeTransaction };
  };

  const getDefaultValues = (isL1ToL2: boolean, address: string) => {
    return {
      fromChain: TOKENS[isL1ToL2 ? "L1" : "L2"][0].value,
      toChain: CHAINS[isL1ToL2 ? "L2" : "L1"][0].value,
      token: TOKENS[isL1ToL2 ? "L1" : "L2"][0].value,
      receiver: address,
      amount: "",
    };
  };

  const handleSwitchNetwork = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    try {
      switchNetwork();
    } catch (error) {
      console.error("Network switch failed", error);
      toast({
        title: "Network Switch",
        description: `Error: ${
          error instanceof Error ? error.message : "switching network"
        }`,
        variant: ToastType.DESTRUCTIVE,
      });
    }
  };

  return {
    useTokenBalance,
    useBridgeTransaction,
    getDefaultValues,
    handleSwitchNetwork,
  };
};
