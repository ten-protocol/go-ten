import { getContractCount, getVerifiedContracts } from "@/api/contracts";
import { useQuery } from "@tanstack/react-query";

export const useContracts = () => {
  const { data: contracts, isLoading: isContractsLoading } = useQuery({
    queryKey: ["contracts"],
    queryFn: () => getVerifiedContracts(),
  });

  const { data: contractCount, isLoading: isContractCountLoading } = useQuery({
    queryKey: ["contractCount"],
    queryFn: () => getContractCount(),
  });

  return {
    contracts,
    isContractsLoading,
    contractCount,
    isContractCountLoading,
  };
};
