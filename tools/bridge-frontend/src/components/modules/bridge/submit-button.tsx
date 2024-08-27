import { Loader } from "lucide-react";
import { Button } from "../../ui/button";
import ConnectWalletButton from "../common/connect-wallet";

export const SubmitButton = ({
  walletConnected,
  loading,
  tokenBalance,
}: {
  walletConnected: boolean;
  loading: boolean;
  tokenBalance: number;
}) => {
  return (
    <div className="flex items-center justify-center mt-4">
      {walletConnected ? (
        <Button
          type="submit"
          className="text-sm font-bold leading-none w-full"
          size="lg"
          disabled={loading || tokenBalance <= 0}
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
