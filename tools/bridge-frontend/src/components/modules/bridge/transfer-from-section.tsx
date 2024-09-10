import { IChain, IToken } from "@/src/types";
import { Separator } from "../../ui/separator";
import { Skeleton } from "../../ui/skeleton";
import { AmountInput } from "./amount-input";
import { ChainSelect } from "./chain-select";
import { PercentageButtons } from "./percentage-buttons";
import { TokenSelect } from "./token-select";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";
import { RefreshCcwDotIcon } from "lucide-react";
import Spinner from "../../ui/spinner";

export const TransferFromSection = ({
  form,
  fromChains,
  tokens,
  tokenBalance,
  balanceLoading,
  balanceFetching,
  setAmount,
  walletConnected,
  refreshBalance,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  fromChains: IChain[];
  tokens: IToken[];
  tokenBalance: number;
  balanceLoading: boolean;
  balanceFetching: boolean;
  setAmount: (value: number) => void;
  walletConnected: boolean;
  refreshBalance: () => void;
}) => {
  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <strong>Transfer from</strong>
        <ChainSelect form={form} chains={fromChains} name="fromChain" />
      </div>
      <div className="bg-muted dark:bg-[#15171D] rounded-lg border">
        <div className="flex items-center justify-between p-2">
          <TokenSelect form={form} tokens={tokens} />
          <div className="pl-2">
            <p className="text-sm text-muted-foreground">Balance:</p>
            <strong className="text-lg float-right word-wrap">
              {balanceLoading ? <Skeleton /> : tokenBalance || "0.00"}{" "}
              {balanceLoading || balanceFetching ? (
                <Spinner
                  size="sm"
                  className="h-4 w-4 inline-block cursor-pointer"
                />
              ) : (
                <RefreshCcwDotIcon
                  className="h-4 w-4 inline-block cursor-pointer"
                  onClick={refreshBalance}
                />
              )}
            </strong>
          </div>
        </div>
        <Separator />
        <div className="flex items-center justify-between flex-wrap p-2">
          <AmountInput form={form} walletConnected={walletConnected} />
          <PercentageButtons setAmount={setAmount} />
        </div>
      </div>
    </div>
  );
};
