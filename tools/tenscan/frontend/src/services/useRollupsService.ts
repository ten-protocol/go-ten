import {
  decryptEncryptedRollup,
  fetchLatestRollups,
  fetchRollups,
} from "@/api/rollups";
import { toast } from "@/src/components/ui/use-toast";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { pollingInterval } from "../lib/constants";

export const useRollupsService = () => {
  const [noPolling, setNoPolling] = useState(false);
  const [decryptedRollup, setDecryptedRollup] = useState<any>();

  const { data: latestRollups } = useQuery({
    queryKey: ["latestRollups"],
    queryFn: () => fetchLatestRollups(),
  });

  const {
    data: rollups,
    isLoading: isRollupsLoading,
    refetch: refetchRollups,
  } = useQuery({
    queryKey: ["rollups"],
    queryFn: () => fetchRollups(),
    refetchInterval: noPolling ? false : pollingInterval,
  });

  const { mutate: decryptEncryptedData } = useMutation({
    mutationFn: decryptEncryptedRollup,
    onSuccess: (data: any) => {
      setDecryptedRollup(data);
      toast({ description: "Decryption successful!" });
    },
    onError: (error: any) => {
      toast({ description: error.message });
    },
  });

  return {
    rollups,
    latestRollups,
    refetchRollups,
    isRollupsLoading,
    decryptEncryptedData,
    decryptedRollup,
    setNoPolling,
  };
};
