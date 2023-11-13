import { getBatches, getLatestBatch } from "@/api/batches";
import { useQuery } from "@tanstack/react-query";

export const useBatches = () => {
  const { data: batches, isLoading: isBatchesLoading } = useQuery({
    queryKey: ["batches"],
    queryFn: () => getBatches(),
  });

  const { data: latestBatch, isLoading: isLatestBatchLoading } = useQuery({
    queryKey: ["latestBatch"],
    queryFn: () => getLatestBatch(),
  });

  return { batches, isBatchesLoading, latestBatch, isLatestBatchLoading };
};
