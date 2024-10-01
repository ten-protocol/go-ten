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
import { DrawerDialog } from "../common/drawer-dialog";
import { useWatch } from "react-hook-form";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";
import { bridgeSchema } from "@/src/schemas/bridge";
import useWalletStore from "@/src/stores/wallet-store";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import { useBridgeUtils } from "@/src/hooks/useBridgeUtils";
import { SubmitButton } from "./submit-button";
import { TransferToSection } from "./transfer-to-section";
import { SwitchNetworkButton } from "./switch-network-button";
import { TransferFromSection } from "./transfer-from-section";
import { CHAINS, TOKENS } from "@/src/lib/constants";
import { handleStorage } from "@/src/lib/utils";

export default function Dashboard() {
  const queryClient = useQueryClient();
  const {
    useBridgeTransaction,
    useTokenBalance,
    handleSwitchNetwork,
    getDefaultValues,
  } = useBridgeUtils();
  const { address, walletConnected, isL1ToL2, loading, provider } =
    useWalletStore();

  const tokens = TOKENS[isL1ToL2 ? "L1" : "L2"];
  const fromChains = CHAINS[isL1ToL2 ? "L1" : "L2"];
  const toChains = CHAINS[isL1ToL2 ? "L2" : "L1"];

  const defaultValues = getDefaultValues(isL1ToL2, address);

  const form = useCustomHookForm(bridgeSchema, { defaultValues });
  const { handleSubmit, control, setValue, reset, setError, formState } = form;

  const [fromChain, toChain, token, receiver, amount] = useWatch({
    control,
    name: ["fromChain", "toChain", "token", "receiver", "amount"],
  });

  const { tokenBalance, isBalanceLoading, isBalanceFetching, refreshBalance } =
    useTokenBalance(
      tokens,
      token,
      address,
      fromChain,
      walletConnected,
      loading
    );

  const [open, setOpen] = React.useState(false);

  const { initiateBridgeTransaction } = useBridgeTransaction(
    address,
    token,
    tokens,
    receiver,
    tokenBalance,
    setError
  );

  const { mutate, isPending } = useMutation({
    mutationFn: initiateBridgeTransaction,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["bridgeTransactions", isL1ToL2 ? "l1" : "l2"],
      });
      reset();
    },
  });

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
    [tokenBalance, token, setError, setValue]
  );

  React.useEffect(() => {
    const storedReceiver = handleStorage.get("tenBridgeReceiver");
    setValue("receiver", storedReceiver ? storedReceiver : address);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [address]);

  React.useEffect(() => {
    setValue("fromChain", isL1ToL2 ? CHAINS.L1[0].value : CHAINS.L2[0].value);
    setValue("toChain", isL1ToL2 ? CHAINS.L2[0].value : CHAINS.L1[0].value);
    setValue("token", isL1ToL2 ? TOKENS.L1[0].value : TOKENS.L2[0].value);
  }, [isL1ToL2, setValue]);

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
            <form onSubmit={handleSubmit((data) => mutate(data))}>
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
                loading={loading || formState.isSubmitting || isPending}
              />
              <TransferToSection
                form={form}
                toChains={toChains}
                loading={loading || formState.isSubmitting || isPending}
                receiver={receiver}
                address={address}
                setOpen={setOpen}
              />
              <SubmitButton
                walletConnected={walletConnected}
                loading={loading}
                isSubmitting={formState.isSubmitting || isPending}
                tokenBalance={tokenBalance}
                provider={provider}
                hasValue={!!amount}
              />
            </form>
          </Form>
          <DrawerDialog open={open} setOpen={setOpen} form={form} />
        </CardContent>
      </Card>
    </div>
  );
}
