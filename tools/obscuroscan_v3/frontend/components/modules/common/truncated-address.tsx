import { useCopy } from "@/src/hooks/useCopy";
import { CopyIcon } from "@radix-ui/react-icons";
import React from "react";

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
          <div>{truncatedAddress}</div>
          <button
            className="text-muted-foreground hover:text-primary transition-colors"
            onClick={() => copyToClipboard(address)}
          >
            <CopyIcon />
          </button>
        </>
      ) : (
        <div>N/A</div>
      )}
    </div>
  );
};

export default TruncatedAddress;
