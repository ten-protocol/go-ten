import React from "react";
import { Badge, badgeVariants } from "@repo/ui/shared/badge";
import { useGeneralService } from "../services/useGeneralService";

const HealthIndicator = () => {
  const [status, setStatus] = React.useState<boolean>(false);
  const { testnetStatus } = useGeneralService();

  // if testnetStatus.result is true, status is set to true
  // if testnetStatus.result is false but testnetStatus.error includes [p2p], status is set to true

  React.useEffect(() => {
    if (testnetStatus?.result) {
      setStatus(true);
    } else if (testnetStatus?.errors?.includes("[p2p]")) {
      setStatus(true);
    } else {
      setStatus(false);
    }
  }, [testnetStatus]);

  return (
    <div className="flex items-center space-x-2 border rounded-lg p-2">
      <h3 className="text-sm">Testnet Status: </h3>
      <Badge
        variant={
          (status ? "success" : "destructive") as keyof typeof badgeVariants
        }
      >
        {status ? "Live" : "Down"}
      </Badge>
    </div>
  );
};

export default HealthIndicator;
