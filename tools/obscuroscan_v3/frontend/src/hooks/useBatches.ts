import { getBatches, getLatestBatch } from "@/api/batches";
import { useQuery } from "@tanstack/react-query";
import { pollingInterval } from "../lib/constants";

export const useBatches = () => {
  const { data: batches, isLoading: isBatchesLoading } = useQuery({
    queryKey: ["batches"],
    queryFn: () => getBatches(),
    refetchInterval: pollingInterval,
  });

  const { data: latestBatch, isLoading: isLatestBatchLoading } = useQuery({
    queryKey: ["latestBatch"],
    queryFn: () => getLatestBatch(),
    refetchInterval: pollingInterval,
  });

  return { batches, isBatchesLoading, latestBatch, isLatestBatchLoading };
};
