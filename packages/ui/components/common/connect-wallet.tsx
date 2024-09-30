import {
  Link2Icon,
  LinkBreak2Icon,
  ExclamationTriangleIcon,
} from "@radix-ui/react-icons";
import { cn, downloadMetaMask, ethereum } from "../../lib/utils";
import useWalletStore from "../../stores/wallet-store";
import TruncatedAddress from "./truncated-address";
import { Button } from "../shared/button";
import { ConnectWalletButtonProps } from "../../lib/interfaces/ui";

const ConnectWalletButton = ({
  className,
  text = "Connect Wallet",
  variant = "outline",
}: ConnectWalletButtonProps) => {
  const {
    walletConnected,
    connectWallet,
    disconnectWallet,
    isWrongNetwork,
    switchNetwork,
    address,
  } = useWalletStore();

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
        {<TruncatedAddress address={address as string} showCopy={false} />}
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
