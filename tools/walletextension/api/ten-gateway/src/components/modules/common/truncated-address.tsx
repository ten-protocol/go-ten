import React from "react";

import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import { useCopy } from "@/hooks/useCopy";
import { CopyIcon } from "@radix-ui/react-icons";
import { Button } from "@/components/ui/button";
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

  const { copyToClipboard } = useCopy();

  return (
    <div className="flex items-center space-x-2">
      {address ? (
        <>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger>{truncatedAddress}</TooltipTrigger>
              <TooltipContent>
                <p>{address}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
          <Copy value={address} />
        </>
      ) : (
        <div>N/A</div>
      )}
    </div>
  );
};

export default TruncatedAddress;
