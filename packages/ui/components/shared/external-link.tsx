import { ExternalLinkIcon } from "@radix-ui/react-icons";
import React from "react";

const ExternalLink = ({
  href,
  children,
  className,
}: {
  href: string;
  children: React.ReactNode;
  className?: string;
}) => {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className={`flex items-center hover:underline ${className}`}
    >
      {children} <ExternalLinkIcon className="inline-block h-4 w-4" />
    </a>
  );
};

export default ExternalLink;
