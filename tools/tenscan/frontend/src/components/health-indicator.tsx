import React from "react";
import { Badge, badgeVariants } from "./ui/badge";
import { useGeneralService } from "../services/useGeneralService";
import { BarChart } from "lucide-react";
import {
  TooltipProvider,
  TooltipTrigger,
  TooltipContent,
  Tooltip,
} from "./ui/tooltip";

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
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger>
          <div className="flex items-center space-x-1 border rounded-lg p-2">
            <span className="text-sm hidden lg:block">Testnet Status</span>
            <BarChart
              className={`h-4 w-4 ${status ? "text-success" : "text-destructive"}`}
            />
            <Badge
              variant={
                (status
                  ? "success"
                  : "destructive") as keyof typeof badgeVariants
              }
            >
              {status ? "Live" : "Down"}
            </Badge>
          </div>
        </TooltipTrigger>
        <TooltipContent>
          {status ? "Testnet status: Live" : "Testnet status: Down"}
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default HealthIndicator;
