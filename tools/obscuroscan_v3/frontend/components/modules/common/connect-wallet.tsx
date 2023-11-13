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
            Disconnect
          </>
        ) : (
          <>
            <Link2Icon className="w-4 h-4 mr-2" />
            Connect
          </>
        )}{" "}
        Wallet
      </Button>
    </>
  );
};

export default ConnectWalletButton;
