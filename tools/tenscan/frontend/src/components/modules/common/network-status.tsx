import { cn } from "@/src/lib/utils";
import React from "react";

const MessageContent = (
  <p className="text-sm">
    <b>You seem to be offline</b>
    <br /> Please check your internet connection.
  </p>
);

export const NetworkStatus = ({ message = MessageContent }) => {
  const [isOnline, setIsOnline] = React.useState(true);

  React.useEffect(() => {
    const setOnlineStatus = () => {
      setIsOnline(navigator.onLine);
    };

    window.addEventListener("online", setOnlineStatus);
    window.addEventListener("offline", setOnlineStatus);

    return () => {
      window.removeEventListener("online", setOnlineStatus);
      window.removeEventListener("offline", setOnlineStatus);
    };
  }, []);

  return (
    <div
      className={cn(
        "fixed z-50 right-0 bottom-0 p-4 transform transition-transform",
        isOnline ? "translate-x-full" : "translate-x-0"
      )}
    >
      <div className="bg-red-500 text-white p-4 rounded-md shadow-lg ring-1 ring-gray-800 backdrop-blur transition dark:bg-gray-800/90 dark:ring-white/10 dark:hover:ring-white/20 cursor-pointer">
        {message}
      </div>
    </div>
  );
};
