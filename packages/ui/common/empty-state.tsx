import { cn } from "@/src/lib/utils";
import Image from "next/image";
import React from "react";

const EmptyState = ({
  title,
  description,
  icon,
  imageSrc,
  imageAlt,
  action,
  className,
}: {
  title?: string;
  description?: string;
  icon?: React.ReactNode;
  imageSrc?: string;
  imageAlt?: string;
  action?: React.ReactNode;
  className?: string;
}) => {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center space-y-4",
        className
      )}
    >
      <div className="flex flex-col items-center justify-center space-y-4">
        {icon && <div className="w-24 h-24">{icon}</div>}
        {imageSrc && (
          <Image
            src={imageSrc}
            alt={imageAlt || "Empty state"}
            className="w-24 h-24 rounded-full"
            width={96}
            height={96}
          />
        )}
        {title && (
          <h3 className="text-2xl font-semibold leading-none tracking-tight">
            {title}
          </h3>
        )}
        {description && (
          <p className="text-sm text-muted-foreground">{description}</p>
        )}
        {action && <div className="flex items-center">{action}</div>}
      </div>
    </div>
  );
};

export default EmptyState;
