import { PlusIcon } from "lucide-react";
import { Button } from "../../ui/button";
import { Skeleton } from "../../ui/skeleton";
import TruncatedAddress from "../common/truncated-address";
import { ChainSelectTo } from "./chain-select";
import { Chain } from "@/src/types";
import { isAddress } from "ethers/lib/utils";
import useCustomHookForm from "@/src/hooks/useCustomHookForm";

export const TransferToSection = ({
  form,
  toChains,
  loading,
  receiver,
  address,
  setOpen,
}: {
  form: ReturnType<typeof useCustomHookForm>;
  toChains: Chain[];
  loading: boolean;
  receiver?: string;
  address: string;
  setOpen: (open: boolean) => void;
}) => {
  return (
    <div>
      <div className="flex items-center justify-between">
        <strong>Transfer to</strong>
        <ChainSelectTo form={form} chains={toChains} />
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
          <strong className="text-lg">{form.getValues().token}</strong>
          <div className="flex flex-col items-end">
            <p className="text-sm text-muted-foreground">You will receive:</p>
            <strong className="text-lg float-right">
              {loading ? <Skeleton /> : form.watch("amount") || 0}
            </strong>
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
