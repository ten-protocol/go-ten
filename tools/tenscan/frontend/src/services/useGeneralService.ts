import { fetchTestnetStatus } from "@/api/general";
import { useQuery } from "@tanstack/react-query";

export const useGeneralService = () => {
  const { data: testnetStatus, isLoading: isStatusLoading } = useQuery({
    queryKey: ["testnetStatus"],
    queryFn: () => fetchTestnetStatus(),
    refetchInterval: 10000,
  });

  return { testnetStatus, isStatusLoading };
};
