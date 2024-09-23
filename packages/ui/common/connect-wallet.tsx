import { Button } from "../shared/button";
import { Link2Icon, LinkBreak2Icon } from "../shared/react-icons";
import TruncatedAddress from "./truncated-address";
import { downloadMetaMask, ethereum } from "../lib/utils";
const ConnectWalletButton = ({
  text,
  walletConnected,
  walletAddress,
  connectWallet,
  disconnectWallet,
}: {
  text?: string;
  walletConnected: boolean;
  walletAddress: string | null;
  connectWallet: () => void;
  disconnectWallet: () => void;
}) => {
  return (
    <Button
      className="text-sm font-medium leading-none"
      variant={"outline"}
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
          {walletAddress ? (
            <TruncatedAddress address={walletAddress} showCopy={false} />
          ) : (
            "Wallet"
          )}
        </>
      ) : (
        <>
          <Link2Icon className="h-4 w-4 mr-1" />
          {ethereum ? "Connect Wallet" : text || "Install MetaMask"}
        </>
      )}
    </Button>
  );
};

export default ConnectWalletButton;
