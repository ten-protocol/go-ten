import { Alert, AlertDescription } from "../../ui/alert";
import useGatewayService from "../../../services/useGatewayService";
import { Terminal, Badge } from "lucide-react";
import React from "react";

import { Button } from "../../ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../../ui/dialog";
import Copy from "../common/copy";
import {tenChainIDDecimal, tenGatewayAddress} from "../../../lib/constants";
import { downloadMetaMask, ethereum } from "@/lib/utils";

const CONNECTION_STEPS = [
  "Hit Connect to Ten and start your journey",
  "Allow MetaMask to switch networks to the Ten Testnet",
  "Sign the <b>Signature Request</b> (this is not a transaction)",
];

const Disconnected = () => {
  const { connectToTenTestnet } = useGatewayService();

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

      <Dialog>
        <DialogTrigger asChild>
          <Button
            variant={"clear"}
            className="text-primary underline flex justify-end"
          >
            How does this work?
          </Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>How does the Ten Gateway work?</DialogTitle>
            <DialogDescription></DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <p>
              By connecting your wallet to Ten and signing the signature request
              you will get a unique token, which is also your <b>viewing key</b>
              . It is contained in the RPC link and unique for each user.
            </p>
            <Alert variant={"warning"} className="flex items-center space-x-2">
              <Terminal className="h-4 w-4" />
              <AlertDescription>
                Do not share your viewing key unless you want others to see the
                details of your transactions.
              </AlertDescription>
            </Alert>
            <p>
              Signing the Signature Request is completely secure. It’s not a
              transaction so cannot spend any of your assets and it doesn’t give
              Ten control over your account.
            </p>
            <div className="flex items-center space-x-2">
              <Badge className="h-4 w-4" />
              <p className="text-sm">RPC URL: {tenGatewayAddress}</p>
            </div>
            <Alert variant={"warning"} className="flex items-center space-x-2">
              <Terminal className="h-4 w-4" />
              <AlertDescription>
                You can only get your RPC link when you begin the connection
                process.
              </AlertDescription>
            </Alert>
            <div className="flex items-center space-x-2">
              <Badge className="h-4 w-4" />
              <p className="text-sm">Chain ID: {tenChainIDDecimal}</p>
              <Copy value={tenChainIDDecimal} />
            </div>
          </div>
          <DialogFooter className="sm:justify-start">
            <DialogClose asChild>
              <Button type="button" variant="secondary">
                Close
              </Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Button
        className="mt-4"
        onClick={ethereum ? connectToTenTestnet : downloadMetaMask}
      >
        <Terminal />
        <span className="ml-2">
          {ethereum ? "Connect to Ten Testnet" : "Install MetaMask to continue"}
        </span>
      </Button>
    </div>
  );
};

export default Disconnected;
