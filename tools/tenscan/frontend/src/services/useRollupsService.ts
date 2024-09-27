import {
  decryptEncryptedRollup,
  fetchBatchesInRollups,
  fetchLatestRollups,
  fetchRollups,
} from "@/api/rollups";
import { toast } from "@/src/components/ui/use-toast";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { getOptions, pollingInterval } from "../lib/constants";
import { useRouter } from "next/router";

export const useRollupsService = () => {
  const { query } = useRouter();

  const [noPolling, setNoPolling] = useState(false);
  const [decryptedRollup, setDecryptedRollup] = useState<any>();

  const options = getOptions(query);

  const { data: latestRollups } = useQuery({
    queryKey: ["latestRollups"],
    queryFn: () => fetchLatestRollups(),
  });

  const {
    data: rollups,
    isLoading: isRollupsLoading,
    refetch: refetchRollups,
  } = useQuery({
    queryKey: ["rollups", options],
    queryFn: () => fetchRollups(options),
    refetchInterval: noPolling ? false : pollingInterval,
  });

  const { mutate: decryptEncryptedData } = useMutation({
    mutationFn: decryptEncryptedRollup,
    onSuccess: (data: any) => {
      setDecryptedRollup(data);
    },
    onError: (error: any) => {
      toast({ description: error.message });
    },
  });

  const {
    data: rollupBatches,
    isLoading: isRollupBatchesLoading,
    refetch: refetchRollupBatches,
  } = useQuery({
    queryKey: ["rollupBatches", { hash: query.hash, options }],
    queryFn: () => fetchBatchesInRollups(query.hash as string, options),
  });

  return {
    rollups,
    latestRollups,
    refetchRollups,
    isRollupsLoading,
    decryptEncryptedData,
    decryptedRollup,
    setNoPolling,
    rollupBatches,
    isRollupBatchesLoading,
    refetchRollupBatches,
  };
};
