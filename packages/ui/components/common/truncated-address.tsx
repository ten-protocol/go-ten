import Link from "next/link";
import { Popover, PopoverContent, PopoverTrigger } from "../shared/popover";
import Copy from "./copy";
import { useMediaQuery } from "../../hooks/useMediaQuery";
import React from "react";

// Note: This assumes EncryptedTextAnimation is available in the consuming app
// We'll need to pass it as a prop or import it from the consuming app's context

const TruncatedAddress = ({
  address,
  prefixLength,
  suffixLength,
  showCopy = true,
  link,
  showFullLength = false,
  animate = false,
  AnimationComponent,
  showPopover = true,
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
  animate?: boolean;
  AnimationComponent?: React.ComponentType<{
    text: string;
    hover: boolean;
    active: boolean;
    onView: boolean;
  }>;
  showPopover?: boolean;
}) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const isDesktop = useMediaQuery("(min-width: 1024px)");

  if (!address) {
    return <span>-</span>;
  }

  // should show full address only on desktop if showFullLength is true; otherwise show truncated version
  const shouldShowFullAddress = isDesktop && showFullLength;

  const truncatedAddress = `${address?.substring(
    0,
    prefixLength || 6,
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
          <Popover open={isOpen} onOpenChange={setIsOpen}>
            <PopoverTrigger asChild>
              <span
                className="cursor-pointer hover:text-primary transition-colors"
                onMouseEnter={() => showPopover && setIsOpen(true)}
                onMouseLeave={() => showPopover && setIsOpen(false)}
              >
                {animate && AnimationComponent ? (
                  <AnimationComponent
                    text={truncatedAddress}
                    hover={false}
                    active={animate}
                    onView={false}
                  />
                ) : (
                  truncatedAddress
                )}
              </span>
            </PopoverTrigger>
            <PopoverContent
              className="w-auto p-2"
              onMouseEnter={() => showPopover && setIsOpen(true)}
              onMouseLeave={() => showPopover && setIsOpen(false)}
            >
              <p className="text-primary font-mono text-sm break-all">
                {animate && AnimationComponent ? (
                  <AnimationComponent
                    text={address}
                    hover={false}
                    active={animate}
                    onView={false}
                  />
                ) : (
                  address
                )}
              </p>
            </PopoverContent>
          </Popover>
        )}
        {showCopy && <Copy value={address} />}
      </div>
    </>
  );
};

export default TruncatedAddress;
