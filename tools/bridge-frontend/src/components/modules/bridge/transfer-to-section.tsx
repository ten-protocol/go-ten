import { useState, useEffect, useCallback } from "react";
import { PlusIcon } from "lucide-react";
import { Button } from "../../ui/button";
import { Skeleton } from "../../ui/skeleton";
import TruncatedAddress from "../common/truncated-address";
import { ChainSelect } from "./chain-select";
import { IChain } from "@/src/types";
import { isAddress } from "ethers/lib/utils";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";
import { useContractService } from "@/src/services/useContractService";
import { ethers } from "ethers";
import useContractStore from "@/src/stores/contract-store";
import { cn } from "@/src/lib/utils";
import { estimateGas } from "@/src/lib/utils/contractUtils";

export const TransferToSection = ({
  form,
  toChains,
  loading,
  receiver,
  address,
  setOpen,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  toChains: IChain[];
  loading: boolean;
  receiver?: string;
  address: string;
  setOpen: (open: boolean) => void;
}) => {
  const { bridgeContract } = useContractStore();
  const [gas, setGas] = useState<string | null>(null);
  const [isEstimating, setIsEstimating] = useState(false);
  const [message, setMessage] = useState<{
    type: "error" | "info";
    message: string;
  }>({
    type: "info",
    message: "",
  });

  const estimateGasFee = useCallback(async () => {
    if (!bridgeContract || !estimateGas) return;
    const amount = form.watch("amount");

    if (!receiver) {
      return setMessage({
        type: "info",
        message: "Enter receiver address.",
      });
    }

    if (!amount) {
      return setMessage({
        type: "info",
        message: "Enter amount to estimate gas fee.",
      });
    }

    try {
      const gasEstimate = await estimateGas(receiver, amount, bridgeContract);
      setGas(ethers.utils.formatEther(gasEstimate));
      setMessage({
        type: "info",
        message: "",
      });
    } catch (err) {
      console.error("Error estimating gas:", err);
      setMessage({
        type: "error",
        message: "Error estimating gas.",
      });
    } finally {
      setIsEstimating(false);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [receiver, form, estimateGas]);

  useEffect(() => {
    estimateGasFee();
  }, [estimateGasFee]);

  return (
    <div>
      <div className="flex items-center justify-between">
        <strong>Transfer to</strong>
        <ChainSelect form={form} chains={toChains} name="toChain" />
      </div>
      <div className="flex items-center justify-end">
        <Button
          variant="ghost"
          className="text-sm font-bold leading-none"
          onClick={() => setOpen(true)}
        >
          <PlusIcon className="h-3 w-3 mr-1" />
          <small>Edit destination address</small>
        </Button>
      </div>
      <div className="bg-muted dark:bg-[#15171D]">
        <div className="flex items-center justify-between p-2">
          <div className="flex flex-col">
            <p className="text-sm text-muted-foreground">Amount to Receive</p>
            <strong className="text-lg">
              {form.watch("amount") || 0} {form.getValues().token}
            </strong>
          </div>
          <div className="flex flex-col items-end">
            <p className="text-sm text-muted-foreground">Est. Gas Fee</p>
            {isEstimating ? (
              <Skeleton />
            ) : message?.message ? (
              <span className={cn("text-sm", "text" + message.type)}>
                {message.message}
              </span>
            ) : (
              <span className="text-lg font-bold">
                {gas ? `${gas}` : "0.00"} ETH
              </span>
            )}
          </div>
        </div>
      </div>
      <div className="bg-muted dark:bg-[#15171D] rounded-lg border flex items-center justify-between mt-4 p-2">
        <span>To:</span>
        {receiver && isAddress(receiver as string) ? (
          <TruncatedAddress address={receiver} />
        ) : address ? (
          <TruncatedAddress address={address} />
        ) : null}
      </div>
    </div>
  );
};
