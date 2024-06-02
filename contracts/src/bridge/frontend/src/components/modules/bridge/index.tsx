import React from "react";
import {
  CardHeader,
  CardTitle,
  CardContent,
  Card,
  CardDescription,
} from "@/src/components/ui/card";
import { Button } from "@/src/components/ui/button";
import { Skeleton } from "@/src/components/ui/skeleton";
import { ArrowDownUpIcon, PlusIcon, Terminal } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../ui/select";
import { Separator } from "../../ui/separator";

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/src/components/ui/form";
import { Input } from "@/src/components/ui/input";
import { toast } from "@/src/components/ui/use-toast";
import { DrawerDialog } from "../common/drawer-dialog";
import { L1TOKENS, L2TOKENS, PERCENTAGES } from "@/src/lib/constants";
import { z } from "zod";
import { useFormHook } from "@/src/hooks/useForm";
import { useWalletStore } from "../../providers/wallet-provider";
import { ToastType, Token } from "@/src/types";
import { Alert, AlertDescription } from "../../ui/alert";
import ConnectWalletButton from "../common/connect-wallet";
import TruncatedAddress from "../common/truncated-address";
import { useContract } from "@/src/hooks/useContract";

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

  const { form, FormSchema } = useFormHook();
  const [loading, setLoading] = React.useState(false);
  const [fromTokenBalance, setFromTokenBalance] = React.useState<any>(0);

  const tokens = isL1ToL2 ? L1TOKENS : L2TOKENS;
  const receiver = form.watch("receiver");

  const swapTokens = () => {
    switchNetwork(isL1ToL2 ? "L2" : "L1");
  };

  const [open, setOpen] = React.useState(false);
  const watchTokenChange = form.watch("token");
  React.useEffect(() => {
    const tokenBalance = async (value: string, token: Token) => {
      setLoading(true);
      try {
        tokens.find((t) => t.value === value);
        let balance;
        if (token.isNative) {
          balance = await getNativeBalance(provider, address);
        } else {
          balance = await getTokenBalance(token.address, address, provider);
        }
        setFromTokenBalance(balance);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    if (watchTokenChange) {
      const token = tokens.find((t) => t.value === watchTokenChange);
      if (token) {
        tokenBalance(watchTokenChange, token);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [watchTokenChange, address, provider]);

  async function onSubmit(data: z.infer<typeof FormSchema>) {
    try {
      setLoading(true);
      const d = {
        ...data,
        receiver: receiver ? receiver : address,
      };
      toast({
        title: "Bridge Transaction",
        description: "Bridge transaction initiated",
        variant: ToastType.INFO,
      });
      const token = d.token;
      const t = tokens.find((t) => t.value === token);
      if (!t) {
        throw new Error("Invalid token");
      }

      const sendTransaction = t.isNative ? sendNative : sendERC20;
      const res = await sendTransaction(
        d.receiver ? d.receiver : address,
        d.amount,
        t.address
      );
      toast({
        title: "Bridge Transaction",
        description: "Bridge transaction completed",
        variant: ToastType.SUCCESS,
      });
      form.reset();
    } catch (error) {
      setLoading(false);
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
  }

  const setAmount = (value: number) => {
    if (!form.getValues("token")) {
      form.setError("token", {
        type: "manual",
        message: "Please select a token to get the balance",
      });
      return;
    }
    const roundDown = (num: number) => Math.floor(num * 100) / 100;
    const amount = roundDown((fromTokenBalance * value) / 100);
    form.setValue("amount", amount.toString());
  };

  return (
    <>
      {!isL1ToL2 && (
        <Alert variant={"destructive"} className="flex items-center space-x-2">
          <Terminal className="h-4 w-4" />
          <AlertDescription>
            Sorry! 😔 L2 to L1 bridge is not supported yet.
          </AlertDescription>
        </Alert>
      )}
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
              <form onSubmit={form.handleSubmit(onSubmit)}>
                <div>
                  <div className="flex items-center justify-between mb-4">
                    <strong>Transfer from</strong>
                    {/* From Chain Select */}
                    <FormField
                      control={form.control}
                      name="fromChain"
                      render={({ field }) => (
                        <FormItem>
                          <Select
                            defaultValue={field.value}
                            onValueChange={field.onChange}
                          >
                            <FormControl>
                              <SelectTrigger className="h-8 bg-muted">
                                <SelectValue
                                  placeholder={field.value || "Select Chain"}
                                />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              {fromChains.map((chain) => (
                                <SelectItem
                                  key={chain.value}
                                  value={chain.value}
                                >
                                  {chain.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <div className="bg-muted dark:bg-[#15171D] rounded-lg border">
                    <div className="flex items-center justify-between p-2">
                      {/* Token Select */}
                      <FormField
                        control={form.control}
                        name="token"
                        render={({ field }) => (
                          <FormItem>
                            <Select
                              defaultValue={field.value}
                              onValueChange={field.onChange}
                            >
                              <FormControl>
                                <SelectTrigger className="h-8 dark:bg-[#292929]">
                                  <SelectValue
                                    placeholder={field.value || "Select Token"}
                                  />
                                </SelectTrigger>
                              </FormControl>
                              <SelectContent side="top">
                                {tokens.map((token) => (
                                  <SelectItem
                                    key={token.value}
                                    value={token.value}
                                    disabled={!token.isEnabled}
                                  >
                                    {token.value}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      {/* Balance */}
                      <div className="pl-2">
                        <p className="text-sm text-muted-foreground">
                          Balance:
                        </p>
                        <strong className="text-lg float-right word-wrap">
                          {loading ? <Skeleton /> : fromTokenBalance || 0}
                        </strong>
                      </div>
                    </div>
                    <Separator />
                    <div className="flex items-center justify-between flex-wrap p-2">
                      {/* Amount Input */}
                      <FormField
                        control={form.control}
                        name="amount"
                        render={({ field }) => (
                          <FormItem>
                            <FormControl>
                              <Input
                                type="number"
                                placeholder="0"
                                className="text-2xl font-bold w-full dark:bg-[#292929] overflow-ellipsis"
                                disabled={!walletConnected}
                                {...field}
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <div className="flex items-center p-2">
                        {/* Percentage Buttons */}
                        <div className="flex items-center space-x-2">
                          {PERCENTAGES.map((percentage) => (
                            <Button
                              type="button"
                              key={percentage.name}
                              variant="outline"
                              size={"sm"}
                              className="dark:bg-[#292929]"
                              onClick={() => {
                                setAmount(percentage.value);
                              }}
                            >
                              {percentage.name}
                            </Button>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Swap Bridge Tokens */}
                <div className="flex items-center justify-center">
                  <Button
                    type="button"
                    className="mt-4"
                    variant="outline"
                    size={"sm"}
                    onClick={swapTokens}
                  >
                    <ArrowDownUpIcon className="h-4 w-4" />
                  </Button>
                </div>

                {/* Transfer to */}
                <div>
                  <div className="flex items-center justify-between">
                    <strong>Transfer to</strong>
                    {/* To Chain Select */}
                    <FormField
                      control={form.control}
                      name="toChain"
                      render={({ field }) => (
                        <FormItem>
                          <Select
                            defaultValue={field.value}
                            onValueChange={field.onChange}
                          >
                            <FormControl>
                              <SelectTrigger className="h-8 bg-muted">
                                <SelectValue
                                  placeholder={field.value || "Select Chain"}
                                />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent side="top">
                              {toChains.map((chain) => (
                                <SelectItem
                                  key={chain.value}
                                  value={chain.value}
                                >
                                  {chain.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <div className="flex items-center justify-end">
                    {/* Destination Address Input */}
                    <Button
                      variant="ghost"
                      className="text-sm font-bold leading-none"
                      onClick={() => setOpen(true)}
                    >
                      <PlusIcon className="h-3 w-3 mr-1" />
                      <small>Edit destination address</small>
                    </Button>
                  </div>
                  <div className="bg-muted dark:bg-[#15171D]">
                    <div className="flex items-center justify-between p-2">
                      <strong className="text-lg">
                        {form.getValues().token}
                      </strong>

                      <div className="flex flex-col items-end">
                        <p className="text-sm text-muted-foreground">
                          You will receive:
                        </p>
                        <strong className="text-lg float-right">
                          {loading ? <Skeleton /> : form.watch("amount") || 0}
                        </strong>
                      </div>
                    </div>
                  </div>

                  <div className="bg-muted dark:bg-[#15171D] rounded-lg border flex items-center justify-between mt-2 p-2 h-14">
                    <strong className="text-lg">Receiver Address</strong>
                    <div className="flex items-center">
                      {receiver || address ? (
                        <TruncatedAddress address={receiver || address} />
                      ) : null}
                    </div>
                  </div>
                </div>
                <div className="flex items-center justify-center mt-4">
                  {walletConnected ? (
                    <Button
                      type="submit"
                      className="text-sm font-bold leading-none w-full"
                      size={"lg"}
                      disabled={!isL1ToL2}
                    >
                      {loading ? (
                        <Skeleton className="w-20" />
                      ) : (
                        "Initiate Bridge Transaction"
                      )}
                    </Button>
                  ) : (
                    <ConnectWalletButton
                      className="text-sm font-bold leading-none w-full"
                      variant="default"
                    />
                  )}
                </div>
              </form>
            </Form>
            <DrawerDialog open={open} setOpen={setOpen} />
          </CardContent>
        </Card>
      </div>
    </>
  );
}
