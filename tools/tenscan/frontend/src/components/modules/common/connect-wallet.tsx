import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import { Button } from "@/src/components/ui/button";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";
import TruncatedAddress from "./truncated-address";
import { downloadMetaMask, ethereum } from "@/src/lib/utils";
const ConnectWalletButton = ({ text }: { text?: string }) => {
  const { walletConnected, walletAddress, connectWallet, disconnectWallet } =
    useWalletConnection();

  return (
    <Button
      className="font-medium leading-none"
      variant={"outline"}
      onClick={
        ethereum
          ? walletConnected
            ? disconnectWallet
            : connectWallet
          : downloadMetaMask
      }
      suppressHydrationWarning
    >
      {walletConnected ? (
        <>
          <LinkBreak2Icon className="h-4 w-4 mr-1" />
          {walletAddress ? (
            <TruncatedAddress address={walletAddress} showCopy={false} />
          ) : (
            "Wallet"
          )}
        </>
      ) : (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          {ethereum ? "Connect" : text || "Install MetaMask"}
        </>
      )}
    </Button>
  );
};

export default ConnectWalletButton;
