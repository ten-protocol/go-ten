import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "../../lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary text-primary-foreground hover:bg-primary/80 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        secondary:
          "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        destructive:
          "border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        success:
          "border-transparent bg-success text-success-foreground hover:bg-success/80 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        outline:
          "text-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        "static-default":
          "border-transparent bg-primary text-primary-foreground",
        "static-secondary":
          "border-transparent bg-secondary text-secondary-foreground",
        "static-destructive":
          "border-transparent bg-destructive text-destructive-foreground",
        "static-success":
          "border-transparent bg-success text-success-foreground",
        "static-outline": "border-current text-foreground",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
