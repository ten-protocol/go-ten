import { useWalletConnection } from "@/components/providers/wallet-provider";
import { Button } from "@/components/ui/button";
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
          <LinkBreak2Icon className="h-4 w-4 mr-1" />
          Disconnect
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
