import { cn } from "../../lib/utils";

function Skeleton({
  className,
  width,
  height,
  style,
  ...props
}: {
  className?: string;
  width?: number | string;
  height?: number | string;
  style?: React.CSSProperties;
}) {
  return (
    <span
      className={cn(
        "animate-pulse rounded-md bg-muted",
        width && `w-${width}`,
        height && `h-${height}`,
        className
      )}
      {...props}
    />
  );
}

export { Skeleton };
