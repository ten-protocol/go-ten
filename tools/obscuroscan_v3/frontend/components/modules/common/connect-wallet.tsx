import { useWalletConnection } from "@/components/providers/wallet-provider";
import { Button } from "@/components/ui/button";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";

const ConnectWalletButton = ({ children }: { children?: React.ReactNode }) => {
  const { walletConnected, connectWallet, disconnectWallet } =
    useWalletConnection();

  return (
    <>
      {children}
      <Button
        variant={"outline"}
        onClick={walletConnected ? disconnectWallet : connectWallet}
      >
        {walletConnected ? (
          <>
            <LinkBreak2Icon className="w-4 h-4 mr-2" />
            <span className="hidden sm:inline">Disconnect Wallet</span>
          </>
        ) : (
          <>
            <Link2Icon className="w-4 h-4 mr-2" />
            <span className="sm:hidden">Connect to Metamask</span>
          </>
        )}
      </Button>
    </>
  );
};

export default ConnectWalletButton;
