import { getRollups } from "@/api/rollups";
import { useQuery } from "@tanstack/react-query";

export const useRollups = () => {
  const { data: rollups, isLoading: isRollupsLoading } = useQuery({
    queryKey: ["rollups"],
    queryFn: () => getRollups(),
  });

  return { rollups, isRollupsLoading };
};
