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
  prefixLength = 6,
  suffixLength = 4,
  showCopy = true,
}: {
  address: string;
  prefixLength?: number;
  suffixLength?: number;
  showCopy?: boolean;
}) => {
  const truncatedAddress =
    address &&
    `${address.substring(0, prefixLength)}...${address.substring(
      address.length - suffixLength
    )}`;

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
          {showCopy && <Copy value={address} />}
        </div>
      ) : (
        <div>N/A</div>
      )}
    </>
  );
};

export default TruncatedAddress;
