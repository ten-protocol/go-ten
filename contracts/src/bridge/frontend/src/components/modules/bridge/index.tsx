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
import { ToastType, Token } from "@/src/types";
import { useContract } from "@/src/hooks/useContract";
import { TransferFromSection } from "./transfer-from-section";
import { SubmitButton } from "./submit-button";
import { SwitchNetworkButton } from "./switch-network-button";
import { TransferToSection } from "./transfer-to-section";
import { bridgeSchema } from "@/src/schemas/bridge";
import { handleStorage } from "@/src/lib/utils/walletUtils";
import useWalletStore from "@/src/stores/wallet-store";

export default function Dashboard() {
  const { provider, address, walletConnected, switchNetwork, isL1ToL2 } =
    useWalletStore();
  const { getNativeBalance, getTokenBalance, sendERC20, sendNative } =
    useContract();
  const intervalId = React.useRef<any>(null);

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

  const [fromTokenBalance, setFromTokenBalance] = React.useState<any>(0);
  const [loading, setLoading] = React.useState(false);
  const [open, setOpen] = React.useState(false);

  const onSubmit = React.useCallback(
    async (data: any) => {
      if (amount > fromTokenBalance) {
        setError("amount", {
          type: "manual",
          message: "Amount must be less than balance",
        });
        return;
      }
      try {
        setLoading(true);
        const transactionData = { ...data, receiver: receiver || address };
        toast({
          title: "Bridge Transaction",
          description: "Bridge transaction initiated",
          variant: ToastType.INFO,
        });

        const selectedToken = token
          ? tokens.find((t: Token) => t.value === token)
          : null;

        if (!selectedToken) throw new Error("Invalid token");

        const sendTransaction = selectedToken.isNative ? sendNative : sendERC20;
        const res = await sendTransaction(
          transactionData.receiver,
          transactionData.amount,
          selectedToken.address
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
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [address, token, tokens, receiver, fromTokenBalance]
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
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [fromTokenBalance, token]
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

  React.useEffect(() => {
    const storedReceiver = handleStorage.get("tenBridgeReceiver");
    setValue("receiver", storedReceiver ? storedReceiver : address);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [address]);

  React.useEffect(() => {
    const fetchTokenBalance = async (token: Token) => {
      setLoading(true);
      try {
        const balance = token.isNative
          ? await getNativeBalance(address)
          : await getTokenBalance(token.address, address);
        setFromTokenBalance(balance);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    intervalId.current = setInterval(() => {
      if (token) {
        const selectedToken = tokens.find((t: Token) => t.value === token);
        if (selectedToken) {
          fetchTokenBalance(selectedToken);
        }
      }
    }, balancePollingInterval);

    return () => {
      clearInterval(intervalId.current);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    fromChain,
    token,
    amount,
    receiver,
    provider,
    isL1ToL2,
    walletConnected,
    address,
  ]);

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
