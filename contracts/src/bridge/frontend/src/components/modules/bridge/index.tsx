import React from "react";
import {
  CardHeader,
  CardTitle,
  CardContent,
  Card,
  CardDescription,
} from "@/src/components/ui/card";
import {
  LayersIcon,
  FileTextIcon,
  ReaderIcon,
  CubeIcon,
  RocketIcon,
  ArrowRightIcon,
} from "@radix-ui/react-icons";

import { Button } from "@/src/components/ui/button";
import { Skeleton } from "@/src/components/ui/skeleton";
import Link from "next/link";
import { cn } from "@/src/lib/utils";
import { Badge } from "../../ui/badge";
import { ArrowDownUpIcon, BlocksIcon, PlusIcon } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../ui/select";
import { Separator } from "../../ui/separator";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

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

interface Token {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
}

const L1CHAINS = [
  {
    name: "Ethereum",
    value: "ETH",
    isNative: true,
    isEnabled: true,
  },
];

const L2CHAINS = [
  {
    name: "TEN",
    value: "TEN",
    isNative: false,
    isEnabled: true,
  },
];

const L1TOKENS = [
  {
    name: "ETH",
    value: "ETH",
    isNative: true,
    isEnabled: true,
  },
];

const L2TOKENS = [
  // {
  //   name: "ETH",
  //   value: "ETH",
  //   isNative: true,
  //   isEnabled: true,
  // },
  {
    name: "USDC",
    value: "USDC",
    isNative: false,
    isEnabled: true,
  },
  {
    name: "USDT",
    value: "USDT",
    isNative: false,
    isEnabled: true,
  },
  {
    name: "TEN",
    value: "TEN",
    isNative: false,
    isEnabled: false,
  },
];

const PERCENTAGES = [
  {
    name: "25%",
    value: 25,
  },
  {
    name: "50%",
    value: 50,
  },
  {
    name: "MAX",
    value: 100,
  },
];

export default function Dashboard() {
  const [loading, setLoading] = React.useState(false);
  const [isL1ToL2, setIsL1ToL2] = React.useState(true);

  const fromChains = isL1ToL2 ? L1CHAINS : L2CHAINS;
  const toChains = isL1ToL2 ? L2CHAINS : L1CHAINS;

  const fromTokens = isL1ToL2 ? L1TOKENS : L2TOKENS;
  const toTokens = isL1ToL2 ? L2TOKENS : L1TOKENS;

  const swapTokens = () => {
    setIsL1ToL2(!isL1ToL2);
  };

  const FormSchema = z.object({
    amount: z.string().nonempty({
      message: "Amount is required.",
    }),
    fromChain: z.string().nonempty({
      message: "From Chain is required.",
    }),
    toChain: z.string().nonempty({
      message: "To Chain is required.",
    }),
    fromToken: z.string().nonempty({
      message: "From Token is required.",
    }),
    toToken: z.string().nonempty({
      message: "To Token is required.",
    }),
  });

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      amount: "",
      fromChain: "",
      toChain: "",
      fromToken: "",
      toToken: "",
    },
  });

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
    if (!form.getValues("fromToken")) {
      form.setError("fromToken", {
        type: "manual",
        message: "Please select a token first.",
      });
      return;
    }
    const balance = 101;
    const amount = Math.floor((balance * value) / 100);
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
                    {/* From Token Select */}
                    <FormField
                      control={form.control}
                      name="fromToken"
                      render={({ field }) => (
                        <FormItem>
                          <Select
                            defaultValue={field.value}
                            onValueChange={field.onChange}
                          >
                            <FormControl>
                              <SelectTrigger className="h-8 w-[80px] bg-[#292929]">
                                <SelectValue placeholder="From" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent side="top">
                              {fromTokens.map((token) => (
                                <SelectItem
                                  key={token.value}
                                  value={token.value}
                                >
                                  {token.name}
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
                        100.00
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
                        <Button type="submit">Submit</Button>
                      </form>
                    )}
                  />
                </div>
                <div className="bg-[#15171D]">
                  <div className="flex items-center justify-between p-2">
                    {/* To Token Select */}
                    <FormField
                      control={form.control}
                      name="toToken"
                      render={({ field }) => (
                        <FormItem>
                          <Select
                            defaultValue={field.value}
                            onValueChange={field.onChange}
                          >
                            <FormControl>
                              <SelectTrigger className="h-8 w-[80px] bg-[#292929]">
                                <SelectValue placeholder="To" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent side="top">
                              {toTokens.map((token) => (
                                <SelectItem
                                  key={token.value}
                                  value={token.value}
                                  disabled={!token.isEnabled}
                                >
                                  {token.name}
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
                        100.00
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
