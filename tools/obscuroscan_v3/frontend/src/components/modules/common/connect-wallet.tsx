import { useWalletConnection } from "@/src/components/providers/wallet-provider";
import { Button } from "@/src/components/ui/button";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";
const ConnectWalletButton = () => {
  const { walletConnected, connectWallet, disconnectWallet } =
    useWalletConnection();

  return (
    <Button
      className="text-sm font-medium leading-none"
      variant={"outline"}
      onClick={walletConnected ? disconnectWallet : connectWallet}
    >
      {walletConnected ? (
        <>
          <LinkBreak2Icon className="w-4 h-4" />
          <span className="hidden md:inline ml-2">Disconnect Wallet</span>
        </>
      ) : (
        <>
          <Link2Icon className="w-4 h-4" />
          <span className="hidden md:inline ml-2">Connect to Metamask</span>
        </>
      )}
    </Button>
  );
};

export default ConnectWalletButton;
