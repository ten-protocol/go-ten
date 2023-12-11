import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import { Button } from "@/src/components/ui/button";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";
import TruncatedAddress from "./truncated-address";
const ConnectWalletButton = () => {
  const { walletConnected, walletAddress, connectWallet, disconnectWallet } =
    useWalletConnection();

  return (
    <Button
      className="text-sm font-medium leading-none"
      variant={"outline"}
      onClick={walletConnected ? disconnectWallet : connectWallet}
    >
      {walletConnected ? (
        <>
          <LinkBreak2Icon className="h-4 w-4 mr-1" />
          {walletAddress ? (
            <TruncatedAddress address={walletAddress} />
          ) : (
            "Wallet"
          )}
        </>
      ) : (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          Connect
          <span className="hidden sm:inline">&nbsp;Wallet</span>
        </>
      )}
    </Button>
  );
};

export default ConnectWalletButton;
