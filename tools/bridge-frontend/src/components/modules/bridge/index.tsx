import React from "react";
import {
  CardHeader,
  CardTitle,
  CardContent,
  Card,
  CardDescription,
} from "@/src/components/ui/card";
import { Separator } from "../../ui/separator";
import { Form } from "@/src/components/ui/form";
import { toast } from "@/src/components/ui/use-toast";
import { DrawerDialog } from "../common/drawer-dialog";
import {
  balancePollingInterval,
  L1CHAINS,
  L1TOKENS,
  L2CHAINS,
  L2TOKENS,
} from "@/src/lib/constants";
import { useWatch } from "react-hook-form";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";
import { ToastType, IToken } from "@/src/types";
import { useContractsService } from "@/src/services/useContractsService";
import { TransferFromSection } from "./transfer-from-section";
import { SubmitButton } from "./submit-button";
import { SwitchNetworkButton } from "./switch-network-button";
import { TransferToSection } from "./transfer-to-section";
import { bridgeSchema } from "@/src/schemas/bridge";
import { handleStorage } from "@/src/lib/utils";
import useWalletStore from "@/src/stores/wallet-store";
import { useQuery } from "@tanstack/react-query";

export default function Dashboard() {
  const { address, walletConnected, switchNetwork, isL1ToL2, loading } =
    useWalletStore();
  const { getNativeBalance, getTokenBalance, sendERC20, sendNative } =
    useContractsService();

  const tokens = isL1ToL2 ? L1TOKENS : L2TOKENS;
  const fromChains = isL1ToL2 ? L1CHAINS : L2CHAINS;
  const toChains = isL1ToL2 ? L2CHAINS : L1CHAINS;

  const defaultValues = {
    fromChain: isL1ToL2 ? L1CHAINS[0].value : L2CHAINS[0].value,
    toChain: isL1ToL2 ? L2CHAINS[0].value : L1CHAINS[0].value,
    token: isL1ToL2 ? L1TOKENS[0].value : L2TOKENS[0].value,
    receiver: address,
    amount: "",
  };

  const form = useCustomHookForm(bridgeSchema, { defaultValues });

  const { handleSubmit, control, setValue, reset, setError, formState } = form;

  const textValues = useWatch({
    control,
    name: ["fromChain", "toChain", "token", "receiver", "amount"],
  });

  const [fromChain, toChain, token, receiver, amount] = textValues;

  async function fetchTokenBalance(token: string, address: string) {
    if (!token || !address) return null;

    const selectedToken = tokens.find((t) => t.value === token);
    if (!selectedToken) return null;

    return selectedToken.isNative
      ? await getNativeBalance(address)
      : await getTokenBalance(selectedToken.address, address);
  }

  const {
    data,
    isLoading: isBalanceLoading,
    isFetching: isBalanceFetching,
    refetch: refreshBalance,
  } = useQuery({
    queryKey: ["tokenBalance", token, address, fromChain],
    queryFn: () => fetchTokenBalance(token, address),
    enabled: !!token && !!address && !!fromChain && walletConnected,
    refetchInterval: balancePollingInterval,
  });

  const tokenBalance = (data || 0.0) as number;

  const [open, setOpen] = React.useState(false);

  const onSubmit = React.useCallback(
    async (data: any) => {
      if (amount > tokenBalance) {
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

        let res;
        if (selectedToken.isNative) {
          res = await sendNative({
            receiver: transactionData.receiver,
            value: transactionData.amount,
          });
        } else {
          res = await sendERC20(
            transactionData.receiver,
            transactionData.amount,
            selectedToken.address
          );
        }

        toast({
          title: "Bridge Transaction",
          description: `Completed: ${res.transactionHash}`,
          variant: ToastType.SUCCESS,
        });
        reset();
      } catch (error) {
        console.error(error);
        toast({
          title: "Bridge Transaction",
          description: `Error: ${
            error instanceof Error
              ? error.message
              : "initiating bridge transaction"
          }`,
          variant: ToastType.DESTRUCTIVE,
        });
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [address, token, tokens, receiver, tokenBalance]
  );

  const setAmount = React.useCallback(
    (value: number) => {
      if (!token) {
        setError("token", {
          type: "manual",
          message: "Please select a token to get the balance",
        });
        return;
      }
      const amount = Math.floor(((tokenBalance * value) / 100) * 100) / 100;
      setValue("amount", amount.toString());
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [tokenBalance, token]
  );

  const handleSwitchNetwork = React.useCallback(
    async (event: React.MouseEvent<HTMLButtonElement>) => {
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
    },
    [switchNetwork]
  );

  React.useEffect(() => {
    const storedReceiver = handleStorage.get("tenBridgeReceiver");
    setValue("receiver", storedReceiver ? storedReceiver : address);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [address]);

  React.useEffect(() => {
    setValue("fromChain", isL1ToL2 ? L1CHAINS[0].value : L2CHAINS[0].value);
    setValue("toChain", isL1ToL2 ? L2CHAINS[0].value : L1CHAINS[0].value);
    setValue("token", isL1ToL2 ? L1TOKENS[0].value : L2TOKENS[0].value);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isL1ToL2]);

  return (
    <div className="h-full flex flex-col space-y-4 justify-center items-center">
      <Card className="max-w-[600px] p-0">
        <CardHeader>
          <CardTitle>Bridge</CardTitle>
          <CardDescription>
            You are currently bridging from {isL1ToL2 ? "L1" : "L2"} to{" "}
            {isL1ToL2 ? "L2" : "L1"}.
          </CardDescription>
        </CardHeader>
        <Separator />
        <CardContent className="p-4">
          <Form {...form}>
            <form onSubmit={handleSubmit(onSubmit)}>
              <TransferFromSection
                form={form}
                fromChains={fromChains}
                tokens={tokens}
                tokenBalance={tokenBalance}
                balanceLoading={isBalanceLoading}
                balanceFetching={isBalanceFetching}
                setAmount={setAmount}
                walletConnected={walletConnected}
                refreshBalance={refreshBalance}
              />
              <SwitchNetworkButton
                handleSwitchNetwork={handleSwitchNetwork}
                loading={loading || formState.isSubmitting}
              />
              <TransferToSection
                form={form}
                toChains={toChains}
                loading={loading || formState.isSubmitting}
                receiver={receiver}
                address={address}
                setOpen={setOpen}
              />
              <SubmitButton
                walletConnected={walletConnected}
                loading={loading || formState.isSubmitting}
                tokenBalance={tokenBalance}
              />
            </form>
          </Form>
          <DrawerDialog open={open} setOpen={setOpen} form={form} />
        </CardContent>
      </Card>
    </div>
  );
}
