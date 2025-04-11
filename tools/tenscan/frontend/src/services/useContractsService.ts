import { fetchContractCount, fetchVerifiedContracts } from "@/api/contracts";
import { useQuery } from "@tanstack/react-query";

export const useContractsService = () => {
  const { data: contracts, isLoading: isContractsLoading } = useQuery({
    queryKey: ["contracts"],
    queryFn: () => fetchVerifiedContracts(),
  });

  const { data: contractCount, isLoading: isContractCountLoading } = useQuery({
    queryKey: ["contractCount"],
    queryFn: () => fetchContractCount(),
  });

  const formattedContracts = [
    {
      name: "Network Config",
      address: contracts?.item.NetworkConfig,
      confirmed: true,
    },
    {
      name: "Enclave Registry",
      address: contracts?.item.EnclaveRegistry,
      confirmed: true,
    },
    {
      name: "Cross Chain",
      address: contracts?.item.CrossChain,
      confirmed: true,
    },
    {
      name: "Data Availability Registry",
      address: contracts?.item.DataAvailabilityRegistry,
      confirmed: true,
    },
    {
      name: "L1 Message Bus",
      address: contracts?.item.L1MessageBus,
      confirmed: true,
    },
    {
      name: "L2 Message Bus",
      address: contracts?.item.L2MessageBus,
      confirmed: true,
    },
    {
      name: "L1 Bridge",
      address: contracts?.item.L1Bridge,
      confirmed: true,
    },
    {
      name: "L2 Bridge",
      address: contracts?.item.L2Bridge,
      confirmed: true,
    },
    {
      name: "L1 Cross Chain Messenger",
      address: contracts?.item.L1CrossChainMessenger,
      confirmed: true,
    },
    {
      name: "L2 Cross Chain Messenger",
      address: contracts?.item.L2CrossChainMessenger,
      confirmed: true,
    },
    {
      name: "Transactions Post Processor",
      address: contracts?.item.TransactionsPostProcessor,
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
