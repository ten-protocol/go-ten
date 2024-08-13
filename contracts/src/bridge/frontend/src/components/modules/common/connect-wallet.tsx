import { Button } from "../../ui/button";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";
import { cn, downloadMetaMask, ethereum } from "@/src/lib/utils";
import useWalletStore from "@/src/stores/wallet-store";

const ConnectWalletButton = ({ className, text, variant }: any) => {
  const { walletConnected, connectWallet, disconnectWallet } = useWalletStore();

  return (
    <Button
      className={cn("text-sm font-medium leading-none", className)}
      variant={variant ? variant : "outline"}
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
          Disconnect
        </>
      ) : (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          {ethereum ? text ?? "Connect Wallet" : "Download MetaMask"}
        </>
      )}
    </Button>
  );
};

export default ConnectWalletButton;
