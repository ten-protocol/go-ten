import React from "react";

import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "@radix-ui/react-tooltip";
import Copy from "./copy";

const TruncatedAddress = ({
  address,
  prefixLength,
  suffixLength,
}: {
  address: string;
  prefixLength?: number;
  suffixLength?: number;
}) => {
  const truncatedAddress = `${address?.substring(
    0,
    prefixLength || 6
  )}...${address?.substring(address.length - (suffixLength || 4))}`;

  return (
    <>
      {address ? (
        <div className="flex items-center">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger>{truncatedAddress}</TooltipTrigger>
              <TooltipContent>
                <p>{address}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <Copy value={address} />
        </div>
      ) : (
        <div>N/A</div>
      )}
    </>
  );
};

export default TruncatedAddress;
