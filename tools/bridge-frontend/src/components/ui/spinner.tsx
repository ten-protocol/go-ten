import { cn } from "@/src/lib/utils";
import React from "react";

const Spinner = ({
  size,
  className,
}: {
  size?: "sm" | "md" | "lg" | "xl";
  className?: string;
}) => {
  const sizeMap = {
    sm: "h-3 w-3",
    md: "h-8 w-8",
    lg: "h-16 w-16",
    xl: "h-32 w-32",
  };

  return (
    <div className={cn("flex justify-center items-center", className)}>
      <div
        className={cn(
          "animate-spin rounded-full border-t-2 border-b-2 border-primary",
          sizeMap[size || "xl"]
        )}
      />
    </div>
  );
};

export default Spinner;
