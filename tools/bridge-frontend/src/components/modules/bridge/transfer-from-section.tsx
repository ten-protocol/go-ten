import { IChain, IToken } from "@/src/types";
import { Separator } from "../../ui/separator";
import { Skeleton } from "../../ui/skeleton";
import { AmountInput } from "./amount-input";
import { ChainSelect } from "./chain-select";
import { PercentageButtons } from "./percentage-buttons";
import { TokenSelect } from "./token-select";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";

export const TransferFromSection = ({
  form,
  fromChains,
  tokens,
  tokenBalance,
  loading,
  setAmount,
  walletConnected,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  fromChains: IChain[];
  tokens: IToken[];
  tokenBalance: number;
  loading: boolean;
  setAmount: (value: number) => void;
  walletConnected: boolean;
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
              {loading ? <Skeleton /> : tokenBalance || "0.00"}{" "}
              {form.watch("token")}
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
