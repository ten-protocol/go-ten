import { getBlocks } from "@/api/blocks";
import { useQuery } from "@tanstack/react-query";

export const useBlocks = () => {
  const { data: blocks, isLoading: isBlocksLoading } = useQuery({
    queryKey: ["blocks"],
    queryFn: () => getBlocks(),
  });

  return { blocks, isBlocksLoading };
};
