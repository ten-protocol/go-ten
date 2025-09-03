import React, {useEffect} from "react";
import { motion, AnimatePresence } from "framer-motion";
import Link from "next/link";

interface RecentItemsListProps<T> {
  items: T[];
  getItemId: (item: T) => string;
  renderItem: (item: T, isNewItem: boolean) => React.ReactNode;
  getItemLink: (item: T) => string;
  className?: string;
  headers?: React.ReactNode;
}

export function RecentItemsList<T>({
  items,
  getItemId,
  renderItem,
  getItemLink,
  className = "",
  headers
}: RecentItemsListProps<T>) {
  const previousItemsRef = React.useRef<Set<string>>(new Set());
  
  return (
    <div className="space-y-0">
      {headers && (
        <div className="sticky top-0 z-20 bg-background/60 backdrop-blur-sm border-b border-border/50 px-2 py-2 mb-2">
          {headers}
        </div>
      )}
      <div className="space-y-2 overflow-y-scroll">
        <AnimatePresence mode="popLayout">
        {items?.map((item: T, i: number) => {
          // Check if this item is new (wasn't in previous render)
          const itemId = getItemId(item);
          const isNewItem = !previousItemsRef.current.has(itemId);
          
          // Update previous items ref
          useEffect(() => {
            previousItemsRef.current.add(itemId);
          }, [itemId]);
          
          return (
            <motion.div 
              className={`flex items-center group ${className}`}
              key={itemId}
              layout
              initial={isNewItem ? { opacity: 0, y: -20, scale: 0.95 } : false}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: 20, scale: 0.95 }}
              transition={{ 
                duration: 0.4,
                ease: [0.25, 0.46, 0.45, 0.94],
                layout: { duration: 0.3 }
              }}
            >
              <Link
                href={getItemLink(item)}
                className={`flex items-center flex-1 p-2 -mx-2 gradient-flash-hover gradient-flash-entry`}
              >
                {renderItem(item, isNewItem)}
              </Link>
            </motion.div>
          );
        })}
        </AnimatePresence>
      </div>
    </div>
  );
}
