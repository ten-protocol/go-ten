import { fetchBlocks } from "@/api/blocks";
import { useQuery } from "@tanstack/react-query";
import { pollingInterval } from "../lib/constants";

export const useBlocksService = () => {
  const { data: blocks, isLoading: isBlocksLoading } = useQuery({
    queryKey: ["blocks"],
    queryFn: () => fetchBlocks(),
    refetchInterval: pollingInterval,
  });

  return { blocks, isBlocksLoading };
};
