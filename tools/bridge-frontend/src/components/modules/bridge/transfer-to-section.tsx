import { PlusIcon } from "lucide-react";
import { Button } from "../../ui/button";
import { Skeleton } from "../../ui/skeleton";
import TruncatedAddress from "../common/truncated-address";
import { ChainSelect } from "./chain-select";
import { IChain } from "@/src/types";
import { isAddress } from "ethers/lib/utils";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";
import { ethers } from "ethers";
import useContractStore from "@/src/stores/contract-store";
import { estimateGas } from "@/src/lib/utils/contractUtils";
import { useQuery } from "@tanstack/react-query";

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

  const fetchGasEstimate = async () => {
    const amount = form.watch("amount");

    if (!receiver) {
      throw new Error("Enter receiver address.");
    }

    if (!amount) {
      throw new Error("Enter amount to estimate gas fee.");
    }

    if (!bridgeContract) {
      throw new Error("Bridge contract is not available.");
    }

    return await estimateGas(receiver, amount, bridgeContract);
  };

  const { data: gasEstimate, isLoading: isEstimating } = useQuery({
    queryKey: ["gasEstimate", receiver, form.watch("amount")],
    queryFn: fetchGasEstimate,
    enabled: !!receiver && !!form.watch("amount"),
  });

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
            ) : (
              <span className="text-lg font-bold">
                {gasEstimate
                  ? `${ethers.utils.formatEther(gasEstimate)} ETH`
                  : "0.00 ETH"}
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
