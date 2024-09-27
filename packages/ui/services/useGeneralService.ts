import { useQuery } from "@tanstack/react-query";
import { fetchTestnetStatus } from "../api/general";

export const useGeneralService = () => {
  const { data: testnetStatus, isLoading: isStatusLoading } = useQuery({
    queryKey: ["testnetStatus"],
    queryFn: () => fetchTestnetStatus(),
    refetchInterval: 10000,
  });

  return {
    testnetStatus,
    isStatusLoading,
  };
};
