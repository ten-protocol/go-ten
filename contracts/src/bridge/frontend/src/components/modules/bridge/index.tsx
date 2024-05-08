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
import { cn } from "@/src/lib/utils";
import { ArrowDownUpIcon, Terminal } from "lucide-react";
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
import { Label } from "../../ui/label";
import { L1TOKENS, L2TOKENS, PERCENTAGES } from "@/src/lib/constants";
import { z } from "zod";
import { useFormHook } from "@/src/hooks/useForm";
import Web3Service from "@/src/services/web3service";
import { useWalletStore } from "../../providers/wallet-provider";
import { Token } from "@/src/types";
import { Alert, AlertDescription } from "../../ui/alert";

export default function Dashboard() {
  const {
    signer,
    provider,
    address,
    switchNetwork,
    isL1ToL2,
    fromChains,
    toChains,
  } = useWalletStore();
  const web3Service = new Web3Service(signer);

  const { form, FormSchema } = useFormHook();
  const [loading, setLoading] = React.useState(false);
  const [fromTokenBalance, setFromTokenBalance] = React.useState<any>(0);
  const [toReceive, setToReceive] = React.useState(0);

  const tokens = isL1ToL2 ? L1TOKENS : L2TOKENS;

  const swapTokens = () => {
    switchNetwork(isL1ToL2 ? "L2" : "L1");
  };

  const watchTokenChange = form.watch("token");
  React.useEffect(() => {
    const getTokenBalance = async (value: string, token: Token) => {
      setLoading(true);
      try {
        tokens.find((t) => t.value === value);
        let balance;
        if (token.isNative) {
          balance = await web3Service.getNativeBalance(provider, address);
        } else {
          balance = await web3Service.getTokenBalance(
            token.address,
            address,
            provider
          );
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
        getTokenBalance(watchTokenChange, token);
      }
    }
  }, [watchTokenChange, address, provider]);

  function onSubmit(data: z.infer<typeof FormSchema>) {
    toast({
      title: "You submitted the following values:",
      description: (
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(data, null, 2)}</code>
        </pre>
      ),
    });
  }

  const setAmount = (value: number) => {
    const vals = form.getValues();
    console.log("token", vals);
    if (!form.getValues("token")) {
      form.setError("token", {
        type: "manual",
        message: "Please select a token to get the balance",
      });
      return;
    }
    const amount = Math.floor((fromTokenBalance * value) / 100);
    form.setValue("amount", amount.toString());
  };

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
                              <SelectItem key={chain.value} value={chain.value}>
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

                <div className="bg-[#15171D] rounded-lg border">
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
                              <SelectTrigger className="h-8 bg-[#292929]">
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

                    <div>
                      <p className="text-sm text-muted-foreground">Balance:</p>
                      <strong className="text-lg text-white float-right">
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
                              placeholder="0"
                              className="text-2xl font-bold w-full bg-[#292929] border-none outline-none overflow-ellipsis"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <div className="flex items-center p-3">
                      {/* Percentage Buttons */}
                      <div className="flex items-center space-x-2">
                        {PERCENTAGES.map((percentage) => (
                          <Button
                            type="button"
                            key={percentage.name}
                            variant="outline"
                            size={"sm"}
                            className="bg-[#292929]"
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
                              <SelectItem key={chain.value} value={chain.value}>
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
                  <DrawerDialog
                    FormComponent={({ className }) => (
                      <form className={cn("grid items-start gap-4", className)}>
                        <div className="grid gap-2">
                          <Label htmlFor="address">Address</Label>
                          <Input type="address" id="address" defaultValue="" />
                        </div>
                        <Alert
                          variant={"warning"}
                          className="flex items-center space-x-2"
                        >
                          <Terminal className="h-4 w-4" />
                          <AlertDescription>
                            Make sure the address is correct before submitting.
                          </AlertDescription>
                        </Alert>
                        <Button type="submit">Add destination address</Button>
                      </form>
                    )}
                  />
                </div>
                <div className="bg-[#15171D]">
                  <div className="flex items-center justify-between p-2">
                    <strong className="text-lg text-white">
                      {form.getValues().token}
                    </strong>

                    <div>
                      <p className="text-sm text-muted-foreground">
                        You will receive:
                      </p>
                      <strong className="text-lg text-white float-right">
                        {loading ? <Skeleton /> : toReceive || 0}
                      </strong>
                    </div>
                  </div>
                </div>

                <div className="bg-[#15171D] rounded-lg border flex items-center justify-between mt-2 p-2 h-14">
                  <strong>Refuel gas</strong>
                  <div className="flex items-center">Not supported</div>
                </div>
              </div>
              <div className="flex items-center justify-center mt-4">
                <Button
                  type="submit"
                  className="text-sm font-bold leading-none w-full"
                  size={"lg"}
                  onClick={() => {
                    console.log("Transfer");
                  }}
                >
                  Initiate Bridge Transaction
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
