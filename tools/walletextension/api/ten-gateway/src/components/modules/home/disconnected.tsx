import { useWalletConnection } from "@/components/providers/wallet-provider";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "@radix-ui/react-tooltip";
import { Terminal, Badge, CopyIcon } from "lucide-react";
import React from "react";

const CONNECTION_STEPS = [
  "Hit Connect to Ten and start your journey",
  "Allow MetaMask to switch networks to the Ten Testnet",
  "Sign the <b>Signature Request</b> (this is not a transaction)",
];

const Disconnected = () => {
  const { connectToTenTestnet } = useWalletConnection();
  return (
    <div className="flex flex-col items-center justify-center space-y-4">
      <h1 className="text-4xl font-bold">Welcome to the Ten Gateway!</h1>
      <h3 className="text-sm text-muted-foreground my-4">
        Three clicks to setup encrypted communication between MetaMask and TEN.
      </h3>
      <ol className="list-decimal list-inside">
        {CONNECTION_STEPS.map((step, index) => (
          <li key={index}>
            <span dangerouslySetInnerHTML={{ __html: step }} />
          </li>
        ))}
      </ol>
      <Button className="mt-4" onClick={connectToTenTestnet}>
        <Terminal />
        <span className="ml-2">Connect to Ten Testnet</span>
      </Button>
    </div>
  );
};

export default Disconnected;
