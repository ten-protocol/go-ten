import { useQuery } from "@tanstack/react-query";
import { fetchDocument, fetchTestnetStatus } from "../api/general";

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

export const useDocumentService = (id: string | undefined) => {
  return useQuery({
    queryKey: ["document", id],
    queryFn: () => fetchDocument(id || ""),
    enabled: !!id,
  });
};
