import {
  Link2Icon,
  LinkBreak2Icon,
  ExclamationTriangleIcon,
} from "@radix-ui/react-icons";
import { cn, downloadMetaMask, ethereum } from "../lib/utils";
import { ButtonVariants } from "../lib/types";
import TruncatedAddress from "./truncated-address";
import { Button } from "../shared/button";

interface ConnectWalletButtonProps {
  className?: string;
  variant?: ButtonVariants;
  text?: string;
  walletConnected: boolean;
  walletAddress: string | null;
  connectWallet: () => void;
  disconnectWallet: () => void;
  switchNetwork: () => void;
  isWrongNetwork: boolean;
}

const ConnectWalletButton = ({
  className,
  text = "Connect Wallet",
  variant = "outline",
  walletConnected,
  connectWallet,
  disconnectWallet,
  switchNetwork,
  isWrongNetwork,
  walletAddress,
}: ConnectWalletButtonProps) => {
  const handleClick = () => {
    if (!ethereum) {
      downloadMetaMask();
      return;
    }

    if (isWrongNetwork) {
      switchNetwork();
      return;
    }

    if (walletConnected) {
      disconnectWallet();
    } else {
      connectWallet();
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
        {
          <TruncatedAddress
            address={walletAddress as string}
            showCopy={false}
          />
        }
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
