import React from "react";
import {
  CardHeader,
  CardTitle,
  CardContent,
  Card,
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

interface Token {
  name: string;
  value: string;
  isNative: boolean;
  isEnabled: boolean;
}

const L1TOKENS = [
  {
    name: "ETH",
    value: "ETH",
    isNative: true,
    isEnabled: true,
  },
];

const L2TOKENS = [
  {
    name: "ETH",
    value: "ETH",
    isNative: true,
    isEnabled: true,
  },
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

const BRIDGE = [
  {
    name: "L1 to L2",
    address: "0x",
  },
  {
    name: "L2 to L1",
    address: "0x",
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
  const [selectedFromToken, setSelectedFromToken] = React.useState<string>(
    L1TOKENS[0].value
  );
  const [selectedToToken, setSelectedToToken] = React.useState<string>(
    L2TOKENS[0].value
  );
  const [isL1ToL2, setIsL1ToL2] = React.useState(true);

  const swapTokens = () => {
    setIsL1ToL2(!isL1ToL2);
  };

  const setFromToken = (token: string) => {
    setSelectedFromToken(token);
  };

  const setToToken = (token: string) => {
    setSelectedToToken(token);
  };

  return (
    <div className="h-full flex flex-col space-y-4 justify-center items-center">
      <Card className="w-[600px] p-0">
        <CardHeader>
          <CardTitle>Bridge</CardTitle>
        </CardHeader>
        <Separator />
        <CardContent className="p-4">
          <div>
            <strong>Transfer from</strong>

            <div className="bg-[#15171D] rounded-lg border">
              <div className="flex items-center justify-between p-2">
                <Select value={selectedFromToken} onValueChange={setFromToken}>
                  <SelectTrigger className="h-8 w-[70px] bg-[#292929]">
                    <SelectValue placeholder={selectedFromToken} />
                  </SelectTrigger>
                  <SelectContent side="top">
                    {L1TOKENS.map((token) => (
                      <SelectItem key={token.value} value={token.value}>
                        {token.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <div>
                  <p className="text-sm text-muted-foreground">Balance:</p>
                  <strong className="text-lg text-white float-right">
                    100.00
                  </strong>
                </div>
              </div>
              <Separator />
              <div className="flex items-center justify-between p-2">
                <div>
                  <strong className="text-2xl text-white float-right">0</strong>
                </div>
                <div className="flex items-center">
                  <div className="space-x-2">
                    {PERCENTAGES.map((percentage) => (
                      <Button
                        key={percentage.name}
                        variant="outline"
                        size={"sm"}
                        className="bg-[#292929]"
                        onClick={() => {
                          console.log(percentage.value);
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

          {/* Swap Tokens */}
          <div className="flex items-center justify-center">
            <Button
              className="mt-4"
              variant="outline"
              size={"sm"}
              onClick={swapTokens}
            >
              <ArrowDownUpIcon className="h-4 w-4" />
            </Button>
          </div>

          <div>
            <div className="flex items-center justify-between">
              <strong>Transfer to</strong>
              <Button
                className="text-sm font-bold leading-none hover:text-primary hover:bg-transparent"
                variant="ghost"
                onClick={() => {
                  console.log("Switch transfer direction");
                }}
              >
                <PlusIcon className="h-3 w-3 mr-1" />
                <small>Transfer to a different address</small>
              </Button>
            </div>

            <div className="bg-[#15171D]">
              <div className="flex items-center justify-between p-2">
                <Select value={selectedToToken} onValueChange={setToToken}>
                  <SelectTrigger className="h-8 w-[70px] bg-[#292929]">
                    <SelectValue placeholder={selectedToToken} />
                  </SelectTrigger>
                  <SelectContent side="top">
                    {L1TOKENS.map((token) => (
                      <SelectItem key={token.value} value={token.value}>
                        {token.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
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
              className="text-sm font-bold leading-none w-full"
              size={"lg"}
              onClick={() => {
                console.log("Transfer");
              }}
            >
              Initiate Bridge Transaction
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
