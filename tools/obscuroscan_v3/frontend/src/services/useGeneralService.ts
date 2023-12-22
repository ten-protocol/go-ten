import { fetchTestnetStatus } from "@/api/general";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

export const useGeneralService = () => {
  const [noPolling, setNoPolling] = useState(false);

  const {
    data: testnetStatus,
    isLoading: isStatusLoading,
    refetch: refetchTestnetStatus,
  } = useQuery({
    queryKey: ["testnetStatus"],
    queryFn: () => fetchTestnetStatus(),
    // refetchInterval: noPolling ? false : pollingInterval,
  });

  return { testnetStatus, isStatusLoading, setNoPolling, refetchTestnetStatus };
};
