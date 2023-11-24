import { fetchBatches, fetchLatestBatch } from "@/api/batches";
import { useQuery } from "@tanstack/react-query";
import { pollingInterval } from "../lib/constants";
import { useState } from "react";

export const useBatchesService = () => {
  const [options, setOptions] = useState({
    offset: 0,
    limit: 10,
  });

  const [noPolling, setNoPolling] = useState(false);

  const updateQueryParams = (query: any) => {
    console.log(
      "🚀 ~ file: useBatchesService.ts:15 ~ updateQueryParams ~ query:",
      query
    );
    // setOptions({
    //   offset: newOffset,
    //   limit: newSize,
    // });
    // refetchBatches;
  };

  const {
    data: batches,
    isLoading: isBatchesLoading,
    refetch: refetchBatches,
  } = useQuery({
    queryKey: ["batches"],
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
    updateQueryParams,
    setNoPolling,
    refetchBatches,
  };
};
