import { fetchBlocks } from "@/api/blocks";
import { useQuery } from "@tanstack/react-query";
import { getOptions, pollingInterval } from "../lib/constants";
import { useRouter } from "next/router";
import { useState } from "react";

export const useBlocksService = () => {
  const { query } = useRouter();

  const [noPolling, setNoPolling] = useState(true);

  const options = getOptions(query);

  const {
    data: blocks,
    isLoading: isBlocksLoading,
    refetch: refetchBlocks,
  } = useQuery({
    queryKey: ["blocks", options],
    queryFn: () => fetchBlocks(options),
    refetchInterval: noPolling ? false : pollingInterval,
  });

  return { blocks, isBlocksLoading, setNoPolling, refetchBlocks };
};
