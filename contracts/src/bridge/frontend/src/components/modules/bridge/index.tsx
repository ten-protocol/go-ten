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
import { L1TOKENS, L2TOKENS } from "@/src/lib/constants";
import { useWatch } from "react-hook-form";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";
import { useWalletStore } from "../../providers/wallet-provider";
import { ToastType, Token } from "@/src/types";
import { useContract } from "@/src/hooks/useContract";
import { TransferFromSection } from "./transfer-from-section";
import { SubmitButton } from "./submit-button";
import { SwitchNetworkButton } from "./switch-network-button";
import { TransferToSection } from "./transfer-to-section";
import { bridgeSchema } from "@/src/schemas/bridge";
import { isAddress } from "ethers/lib/utils";

export default function Dashboard() {
  const {
    provider,
    address,
    walletConnected,
    switchNetwork,
    isL1ToL2,
    fromChains,
    toChains,
  } = useWalletStore();
  const { getNativeBalance, getTokenBalance, sendERC20, sendNative } =
    useContract();

  const defaultValues = {
    amount: "",
    fromChain: fromChains[0].value,
    toChain: toChains[0].value,
    token: isL1ToL2 ? L1TOKENS[0].value : L2TOKENS[0].value,
    receiver: address,
  };

  const form = useCustomHookForm(bridgeSchema, { defaultValues });

  const {
    handleSubmit,
    control,
    setValue,
    reset,
    setError,
    formState: { errors },
  } = form;

  const textValues = useWatch({
    control,
    name: ["fromChain", "toChain", "token", "receiver", "amount"],
  });
  const [fromChain, toChain, token, receiver, amount] = textValues;

  const [loading, setLoading] = React.useState(false);
  const [fromTokenBalance, setFromTokenBalance] = React.useState<any>(0);

  const tokens = isL1ToL2 ? L1TOKENS : L2TOKENS;

  const [open, setOpen] = React.useState(false);

  React.useEffect(() => {
    const fetchTokenBalance = async (token: Token) => {
      setLoading(true);
      try {
        const balance = token.isNative
          ? await getNativeBalance(
              provider,
              isAddress(receiver) ? receiver : address
            )
          : await getTokenBalance(
              token.address,
              isAddress(receiver) ? receiver : address,
              provider
            );
        setFromTokenBalance(balance);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    if (token) {
      const selectedToken = tokens.find((t) => t.value === token);
      if (selectedToken) {
        fetchTokenBalance(selectedToken);
      }
    }
  }, [fromChain, token, amount, receiver, provider, tokens, isL1ToL2]);

  const onSubmit = React.useCallback(
    async (data: any) => {
      try {
        setLoading(true);
        const transactionData = { ...data, receiver: receiver || address };
        toast({
          title: "Bridge Transaction",
          description: "Bridge transaction initiated",
          variant: ToastType.INFO,
        });
        const token = tokens.find((t) => t.value === transactionData.token);
        if (!token) throw new Error("Invalid token");

        const sendTransaction = token.isNative ? sendNative : sendERC20;
        const res = await sendTransaction(
          transactionData.receiver,
          transactionData.amount,
          token.address
        );

        toast({
          title: "Bridge Transaction",
          description: `Bridge transaction completed: ${res.transactionHash}`,
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
      } finally {
        setLoading(false);
      }
    },
    [address, tokens, fromChain]
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
      const amount = Math.floor(((fromTokenBalance * value) / 100) * 100) / 100;
      setValue("amount", amount.toString());
    },
    [form, fromTokenBalance]
  );

  const handleSwitchNetwork = React.useCallback(
    async (event: any) => {
      event.preventDefault();
      try {
        setLoading(true);
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
      } finally {
        setLoading(false);
      }
    },
    [switchNetwork]
  );

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
                fromTokenBalance={fromTokenBalance}
                loading={loading}
                setAmount={setAmount}
                walletConnected={walletConnected}
              />
              <SwitchNetworkButton
                handleSwitchNetwork={handleSwitchNetwork}
                loading={loading}
              />
              <TransferToSection
                form={form}
                toChains={toChains}
                loading={loading}
                receiver={receiver}
                address={address}
                setOpen={setOpen}
              />
              <SubmitButton
                walletConnected={walletConnected}
                loading={loading}
                fromTokenBalance={fromTokenBalance}
              />
            </form>
          </Form>
          <DrawerDialog open={open} setOpen={setOpen} form={form} />
        </CardContent>
      </Card>
    </div>
  );
}
