import { useWalletConnection } from "../../providers/wallet-provider";
import { Button } from "../../ui/button";
import useGatewayService from "../../../services/useGatewayService";
import { Link2Icon, LinkBreak2Icon } from "@radix-ui/react-icons";
import React from "react";
import { ethereum } from "@/src/lib/utils";
const ConnectWalletButton = () => {
  return (
    <Button
      className="text-sm font-medium leading-none"
      variant={"outline"}
      // onClick={
      //   ethereum
      //     ? walletConnected
      //       ? revokeAccounts
      //       : connectToTenTestnet
      //     : downloadMetaMask
      // }
      // suppressHydrationWarning
    >
      {/* {walletConnected ? (
        <>
          <LinkBreak2Icon className="h-4 w-4 mr-1" />
          Disconnect
        </>
      ) : ( */}
      <>
        <Link2Icon className="h-4 w-4 mr-1" />
        {/* {ethereum ? "Connect" : "Install"} */}
      </>
      {/* )} */}
    </Button>
  );
};

export default ConnectWalletButton;
