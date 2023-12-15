import React from "react";

const EmptyState = ({
  title,
  description,
  icon,
  image,
  action,
}: {
  title?: string;
  description?: string;
  icon?: React.ReactNode;
  image?: string;
  action?: React.ReactNode;
}) => {
  return (
    <div className="flex flex-col items-center justify-center h-full">
      <div className="flex flex-col items-center justify-center space-y-4">
        {icon && <div className="w-24 h-24">{icon}</div>}
        {image && (
          <img
            src={image}
            alt="Empty state"
            className="w-24 h-24 rounded-full"
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
