import { fetchTestnetStatus, fetchNetworkConfig } from "@/api/general";
import { useQuery } from "@tanstack/react-query";

export const useGeneralService = () => {
  const { data: testnetStatus, isLoading: isStatusLoading } = useQuery({
    queryKey: ["testnetStatus"],
    queryFn: () => fetchTestnetStatus(),
    refetchInterval: 10000,
  });

  const { data: networkConfig, isLoading: isNetworkConfigLoading } = useQuery({
    queryKey: ["networkConfig"],
    queryFn: () => fetchNetworkConfig(),
    // refetchInterval: 10000, // TODO: confirm if this is needed
  });

  return {
    testnetStatus,
    isStatusLoading,
    networkConfig,
    isNetworkConfigLoading,
  };
};
