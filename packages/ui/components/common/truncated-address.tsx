import Link from "next/link";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "../shared/tooltip";
import Copy from "./copy";
import { useMediaQuery } from "../../hooks/useMediaQuery";

const TruncatedAddress = ({
  address,
  prefixLength,
  suffixLength,
  showCopy = true,
  link,
  showFullLength = false,
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
  showFullLength?: boolean;
}) => {
  if (!address) {
    return <span>-</span>;
  }

  const isDesktop = useMediaQuery("(min-width: 1024px)");

  // should show full address only on desktop if showFullLength is true; otherwise show truncated version
  const shouldShowFullAddress = isDesktop && showFullLength;

  const truncatedAddress = `${address?.substring(
    0,
    prefixLength || 6
  )}...${address?.substring(address.length - (suffixLength || 4))}`;

  return (
    <>
      <div className="flex items-center">
        {link ? (
          <Link href={link} className="text-primary hover:underline">
            {shouldShowFullAddress ? address : truncatedAddress}
          </Link>
        ) : shouldShowFullAddress ? (
          <span>{address}</span>
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
    </>
  );
};

export default TruncatedAddress;
