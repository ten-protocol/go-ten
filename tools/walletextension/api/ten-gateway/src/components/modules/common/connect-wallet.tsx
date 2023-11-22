import { useWalletConnection } from "@/components/providers/wallet-provider";
import { Button } from "@/components/ui/button";
import useGatewayService from "@/services/useGatewayService";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";
const ConnectWalletButton = () => {
  const { walletConnected, revokeAccounts } = useWalletConnection();
  const { connectToTenTestnet } = useGatewayService();

  return (
    <Button
      className="text-sm font-medium leading-none"
      variant={"outline"}
      onClick={walletConnected ? revokeAccounts : connectToTenTestnet}
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
          <span className="hidden sm:inline">&nbsp;Account(s)</span>
        </>
      )}
    </Button>
  );
};

export default ConnectWalletButton;
