import React from "react";
import { Button } from "../../ui/button";
import {
  Link2Icon,
  LinkBreak2Icon,
  ExclamationTriangleIcon,
} from "@radix-ui/react-icons";
import { useWalletStore } from "@/stores/wallet-store";
import { cn, downloadMetaMask, ethereum } from "@/lib/utils";
import useGatewayService from "@/services/useGatewayService";

interface ConnectWalletButtonProps {
  className?: string;
  text?: string;
  variant?: any;
}

const ConnectWalletButton = ({
  className,
  text = "Connect Wallet",
  variant = "outline",
}: ConnectWalletButtonProps) => {
  const { connectToTenTestnet } = useGatewayService();
  const { walletConnected, revokeAccounts, isWrongNetwork } = useWalletStore();

  const handleClick = () => {
    if (!ethereum) {
      downloadMetaMask();
      return;
    }

    if (isWrongNetwork) {
      connectToTenTestnet();
      return;
    }

    if (walletConnected) {
      revokeAccounts();
    } else {
      connectToTenTestnet();
    }
  };

  const renderButtonContent = () => {
    if (!ethereum) {
      return (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          Download MetaMask
        </>
      );
    }

    if (isWrongNetwork) {
      return (
        <>
          <ExclamationTriangleIcon className="h-4 w-4 mr-1 text-yellow-500" />
          Unsupported network
        </>
      );
    }

    return walletConnected ? (
      <>
        <LinkBreak2Icon className="h-4 w-4 mr-1" />
        Disconnect
      </>
    ) : (
      <>
        <Link2Icon className="h-4 w-4 mr-1" />
        {text}
      </>
    );
  };

  return (
    <Button
      className={cn("text-sm font-medium leading-none", className)}
      variant={variant}
      onClick={handleClick}
      suppressHydrationWarning
    >
      {renderButtonContent()}
    </Button>
  );
};

export default ConnectWalletButton;
