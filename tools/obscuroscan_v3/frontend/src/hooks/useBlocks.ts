import { getBlocks } from "@/api/blocks";
import { useQuery } from "@tanstack/react-query";
import { pollingInterval } from "../lib/constants";

export const useBlocks = () => {
  const { data: blocks, isLoading: isBlocksLoading } = useQuery({
    queryKey: ["blocks"],
    queryFn: () => getBlocks(),
    refetchInterval: pollingInterval,
  });

  return { blocks, isBlocksLoading };
};
