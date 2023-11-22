import React from "react";

import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/src/components/ui/tooltip";

import { useCopy } from "@/src/hooks/useCopy";
import { CopyIcon } from "@radix-ui/react-icons";
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
