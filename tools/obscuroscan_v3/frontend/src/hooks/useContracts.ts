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

  const formattedContracts = [
    {
      name: "Management Contract",
      address: contracts?.item.ManagementContractAddress,
      confirmed: true,
    },
    {
      name: "Message Bus Contract",
      address: contracts?.item.MessageBusAddress,
      confirmed: true,
    },
  ];
  const sequencerData = [
    {
      name: "Sequencer ID",
      address: contracts?.item.SequencerID,
      confirmed: true,
    },
    {
      name: "L1 Start Hash",
      address: contracts?.item.L1StartHash,
      confirmed: true,
    },
  ];

  return {
    formattedContracts,
    sequencerData,
    isContractsLoading,
    contractCount,
    isContractCountLoading,
  };
};
