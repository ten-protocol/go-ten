import React, { useEffect } from "react";
import { Badge, badgeVariants } from "./ui/badge";
import useGatewayService from "@/services/useGatewayService";
import { Skeleton } from "./ui/skeleton";

const HealthIndicator = () => {
  const { getTestnetStatus } = useGatewayService();
  const [loading, setLoading] = React.useState(false);
  const [status, setStatus] = React.useState<boolean>();

  const testnetStatus = async () => {
    setLoading(true);
    try {
      const status = await getTestnetStatus();
      return status;
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let isMounted = true;

    testnetStatus().then((res) => {
      if (isMounted) {
        setStatus(res?.result?.OverallHealth);
      }
    });

    return () => {
      isMounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="flex items-center space-x-2 border rounded-lg p-2">
      <h3 className="text-sm">Testnet Status: </h3>
      {loading ? (
        <Skeleton className="h-4 w-10" />
      ) : (
        <Badge
          variant={
            (status ? "success" : "destructive") as keyof typeof badgeVariants
          }
        >
          {status ? "Live" : "Down"}
        </Badge>
      )}
    </div>
  );
};

export default HealthIndicator;
