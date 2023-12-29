import React, { useEffect } from "react";
import { Badge, badgeVariants } from "./ui/badge";
import useGatewayService from "@/services/useGatewayService";

const HealthIndicator = () => {
  const { getTestnetStatus } = useGatewayService();
  // const testnetStatus = async () => {
  //   const status = await getTestnetStatus();
  //   console.log(
  //     "🚀 ~ file: health-indicator.tsx:9 ~ testnetStatus ~ status:",
  //     status
  //   );
  //   return status;
  // };

  // useEffect(() => {
  //   testnetStatus();
  // }, []);

  return (
    <div className="flex items-center space-x-2 border rounded-lg p-2">
      <h3 className="text-sm">Testnet Status: </h3>
      {/* <Badge
        variant={
          (testnetStatus?.result
            ? "success"
            : "destructive") as keyof typeof badgeVariants
        }
      >
        {testnetStatus?.result ? "Live" : "Down"}
      </Badge> */}
    </div>
  );
};

export default HealthIndicator;
