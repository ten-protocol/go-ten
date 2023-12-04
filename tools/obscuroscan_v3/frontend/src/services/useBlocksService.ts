import { fetchBlocks } from "@/api/blocks";
import { useQuery } from "@tanstack/react-query";
import { pollingInterval } from "../lib/constants";
import { useRouter } from "next/router";
import { useState } from "react";

export const useBlocksService = () => {
  const { query } = useRouter();

  const [noPolling, setNoPolling] = useState(false);

  const options = {
    offset: query.page ? parseInt(query.page as string) : 1,
    size: query.size ? parseInt(query.size as string) : 10,
  };

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
