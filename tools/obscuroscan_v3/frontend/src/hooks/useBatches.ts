import { getBatches, getLatestBatches } from "@/api/batches";
import { useQuery } from "@tanstack/react-query";

export const useBatches = () => {
  const { data: batches, isLoading: isBatchesLoading } = useQuery({
    queryKey: ["batches"],
    queryFn: () => getBatches(),
  });

  const { data: batchCount, isLoading: isBatchCountLoading } = useQuery({
    queryKey: ["batchCount"],
    queryFn: () => getLatestBatches(),
  });

  return { batches, isBatchesLoading, batchCount, isBatchCountLoading };
};
