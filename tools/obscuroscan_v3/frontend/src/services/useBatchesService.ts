import { fetchBatches, fetchLatestBatch } from "@/api/batches";
import { useQuery } from "@tanstack/react-query";
import { getOptions, pollingInterval } from "../lib/constants";
import { useState } from "react";
import { useRouter } from "next/router";

export const useBatchesService = () => {
  const { query } = useRouter();

  const [noPolling, setNoPolling] = useState(false);

  const options = getOptions(query);

  const {
    data: batches,
    isLoading: isBatchesLoading,
    refetch: refetchBatches,
  } = useQuery({
    queryKey: ["batches", options],
    queryFn: () => fetchBatches(options),
    refetchInterval: noPolling ? false : pollingInterval,
  });

  const { data: latestBatch, isLoading: isLatestBatchLoading } = useQuery({
    queryKey: ["latestBatch"],
    queryFn: () => fetchLatestBatch(),
    refetchInterval: noPolling ? false : pollingInterval,
  });

  return {
    batches,
    isBatchesLoading,
    latestBatch,
    isLatestBatchLoading,
    setNoPolling,
    refetchBatches,
  };
};
