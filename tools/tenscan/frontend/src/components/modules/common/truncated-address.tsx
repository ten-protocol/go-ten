import React from "react";

import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/src/components/ui/tooltip";

import Copy from "./copy";
import Link from "next/link";

const TruncatedAddress = ({
  address,
  prefixLength,
  suffixLength,
  showCopy = true,
  link,
}: {
  address: string;
  prefixLength?: number;
  suffixLength?: number;
  showCopy?: boolean;
  link?:
    | string
    | {
        pathname: string;
        query: { [key: string]: string | number };
      };
}) => {
  const truncatedAddress = `${address?.substring(
    0,
    prefixLength || 6
  )}...${address?.substring(address.length - (suffixLength || 4))}`;

  return (
    <>
      {address ? (
        <div className="flex items-center">
          {link ? (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger>
                  <Link href={link} className="text-primary hover:underline">
                    {truncatedAddress}
                  </Link>
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-primary">{address}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          ) : (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger>{truncatedAddress}</TooltipTrigger>
                <TooltipContent>
                  <p className="text-primary">{address}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          )}
          {showCopy && <Copy value={address} />}
        </div>
      ) : (
        <div>N/A</div>
      )}
    </>
  );
};

export default TruncatedAddress;
