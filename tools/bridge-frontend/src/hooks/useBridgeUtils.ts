import { useQuery } from "@tanstack/react-query";
import { useContractsService } from "@/src/services/useContractsService";
import { balancePollingInterval, CHAINS, TOKENS } from "@/src/lib/constants";
import { IToken, ToastType } from "../types";
import { toast } from "../components/ui/use-toast";
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
    loading: boolean
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
    setError: (name: string, error: { type: string; message: string }) => void
  ) => {
    const initiateBridgeTransaction = React.useCallback(
      async (data: any) => {
        if (data.amount > tokenBalance) {
          setError("amount", {
            type: "manual",
            message: "Amount must be less than balance",
          });
          return;
        }
        try {
          const transactionData = { ...data, receiver: receiver || address };
          toast({
            title: "Bridge Transaction",
            description: "Bridge transaction initiated",
            variant: ToastType.INFO,
          });

          const selectedToken = token
            ? tokens.find((t: IToken) => t.value === token)
            : null;

          if (!selectedToken) throw new Error("Invalid token");

          if (selectedToken.isNative) {
            await sendNative({
              receiver: transactionData.receiver,
              value: transactionData.amount,
            });
          } else {
            await sendERC20(
              transactionData.receiver,
              transactionData.amount,
              selectedToken.address
            );
          }
        } catch (error) {
          console.error(error);
          throw error;
        }
      },
      [
        address,
        token,
        tokens,
        receiver,
        tokenBalance,
        sendERC20,
        sendNative,
        setError,
      ]
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
