import { Loader } from "lucide-react";
import { Button } from "../../ui/button";
import ConnectWalletButton from "../common/connect-wallet";
import { ethers } from "ethers";

export const SubmitButton = ({
  walletConnected,
  loading,
  tokenBalance,
  provider,
}: {
  walletConnected: boolean;
  loading: boolean;
  tokenBalance: number;
  provider: ethers.providers.Web3Provider | null;
}) => {
  return (
    <div className="flex items-center justify-center mt-4">
      {walletConnected ? (
        <Button
          type="submit"
          className="text-sm font-bold leading-none w-full"
          size="lg"
          disabled={loading || tokenBalance <= 0 || !provider}
        >
          {loading ? <Loader /> : "Initiate Bridge Transaction"}
        </Button>
      ) : (
        <ConnectWalletButton
          className="text-sm font-bold leading-none w-full"
          variant="default"
        />
      )}
    </div>
  );
};
