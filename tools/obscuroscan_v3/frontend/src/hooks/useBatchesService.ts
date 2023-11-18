import { fetchBatches, fetchLatestBatch } from "@/api/batches";
import { useQuery } from "@tanstack/react-query";
import { pollingInterval } from "../lib/constants";

export const useBatchesService = () => {
  const { data: batches, isLoading: isBatchesLoading } = useQuery({
    queryKey: ["batches"],
    queryFn: () => fetchBatches(),
    refetchInterval: pollingInterval,
  });

  const { data: latestBatch, isLoading: isLatestBatchLoading } = useQuery({
    queryKey: ["latestBatch"],
    queryFn: () => fetchLatestBatch(),
    refetchInterval: pollingInterval,
  });

  return { batches, isBatchesLoading, latestBatch, isLatestBatchLoading };
};
