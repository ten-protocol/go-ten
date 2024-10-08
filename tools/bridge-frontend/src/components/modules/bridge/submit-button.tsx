import { Button } from "../../ui/button";
import ConnectWalletButton from "../common/connect-wallet";
import { ethers } from "ethers";

export const SubmitButton = ({
  walletConnected,
  loading,
  tokenBalance,
  provider,
  hasValue,
  isSubmitting,
}: {
  walletConnected: boolean;
  loading: boolean;
  tokenBalance: number;
  provider: ethers.providers.Web3Provider | null;
  hasValue: boolean;
  isSubmitting: boolean;
}) => {
  return (
    <div className="flex items-center justify-center mt-4">
      {walletConnected ? (
        <Button
          type="submit"
          className="text-sm font-bold leading-none w-full"
          size="lg"
          disabled={
            loading ||
            tokenBalance <= 0 ||
            !provider ||
            !hasValue ||
            isSubmitting
          }
          loading={isSubmitting}
          loadingText="Initiating..."
        >
          Initiate Bridge Transaction
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
