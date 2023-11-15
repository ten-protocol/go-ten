import { decryptEncryptedRollup, getRollups } from "@/api/rollups";
import { toast } from "@/src/components/ui/use-toast";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";

export const useRollups = () => {
  const [decryptedRollup, setDecryptedRollup] = useState<any>();

  const { data: rollups, isLoading: isRollupsLoading } = useQuery({
    queryKey: ["rollups"],
    queryFn: () => getRollups(),
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

  return { rollups, isRollupsLoading, decryptEncryptedData, decryptedRollup };
};
