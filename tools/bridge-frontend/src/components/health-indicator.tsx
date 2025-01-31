import React from "react";
import { Badge, badgeVariants } from "./ui/badge";
import { useGeneralService } from "../services/useGeneralService";
import { Skeleton } from "./ui/skeleton";

const HealthIndicator = () => {
  const [status, setStatus] = React.useState<boolean>(false);
  const { testnetStatus, isStatusLoading } = useGeneralService();

  // if testnetStatus.result is true, status is set to true
  // if testnetStatus.result is false but testnetStatus.error includes [p2p], status is set to true

  React.useEffect(() => {
    if (testnetStatus?.result) {
      setStatus(true);
      //@ts-ignore
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
          (isStatusLoading
            ? "secondary"
            : status
            ? "success"
            : "destructive") as keyof typeof badgeVariants
        }
      >
        {isStatusLoading ? (
          <Skeleton className="w-5 h-5" />
        ) : status ? (
          "Live"
        ) : (
          "Down"
        )}
      </Badge>
    </div>
  );
};

export default HealthIndicator;
