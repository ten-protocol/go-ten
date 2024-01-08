import React from "react";
import { Badge, badgeVariants } from "./ui/badge";
import { useGeneralService } from "../services/useGeneralService";

const HealthIndicator = () => {
  const { testnetStatus } = useGeneralService();

  return (
    <div className="flex items-center space-x-2 border rounded-lg p-2">
      <h3 className="text-sm">Testnet Status: </h3>
      <Badge
        variant={
          (testnetStatus?.result
            ? "success"
            : "destructive") as keyof typeof badgeVariants
        }
      >
        {testnetStatus?.result ? "Live" : "Down"}
      </Badge>
    </div>
  );
};

export default HealthIndicator;
