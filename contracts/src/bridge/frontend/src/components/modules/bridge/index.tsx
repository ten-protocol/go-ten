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
import { BlocksIcon } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../ui/select";
import { table } from "console";
import { Separator } from "../../ui/separator";

const L1TOKENS = [
  {
    name: "ETH",
    address: "0x",
  },
];

const L2TOKENS = [
  {
    name: "ETH",
    address: "0x",
  },
  {
    name: "USDC",
    address: "0x",
  },
  {
    name: "USDT",
    address: "0x",
  },
  {
    name: "DAI",
    address: "0x",
  },
  {
    name: "WBTC",
    address: "0x",
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

export default function Dashboard() {
  const [loading, setLoading] = React.useState(false);
  const [selectedToken, setSelectedToken] = React.useState(L1TOKENS[0]);
  const [selectedBridge, setSelectedBridge] = React.useState(BRIDGE[0]);

  return (
    <div className="h-full flex flex-col space-y-4 justify-center items-center">
      <Card className="w-[650px]">
        <CardHeader>
          <CardTitle>Bridge</CardTitle>
        </CardHeader>
        <CardContent>
          <div>
            <small>Transfer from</small>

            <div className="bg-[#15171D] p-2 rounded-lg space-x-2 border">
              <div className="flex items-center justify-between">
                <Select
                  value={selectedToken.name}
                  onValueChange={(value) => {
                    setSelectedToken(
                      L1TOKENS.find((token) => token.name === value)
                    );
                  }}
                >
                  <SelectTrigger className="h-8 w-[70px] bg-[#292929]">
                    <SelectValue placeholder={selectedToken.name} />
                  </SelectTrigger>
                  <SelectContent side="top">
                    {L1TOKENS.map((token) => (
                      <SelectItem key={token.address} value={token.name}>
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
              <div className="flex items-center justify-between">
                <div>
                  <strong className="text-lg text-white float-right">0</strong>
                </div>
                <Select
                  value={selectedToken.name}
                  onValueChange={(value) => {
                    setSelectedToken(
                      L1TOKENS.find((token) => token.name === value)
                    );
                  }}
                >
                  <SelectTrigger className="h-8 w-[70px] bg-[#292929]">
                    <SelectValue placeholder={selectedToken.name} />
                  </SelectTrigger>
                  <SelectContent side="top">
                    {L1TOKENS.map((token) => (
                      <SelectItem key={token.address} value={token.name}>
                        {token.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
