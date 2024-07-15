import React from "react";
import { Separator } from "./separator";

export const KeyValueList = ({ children }: { children: React.ReactNode }) => (
  <ul className="divide-y divide-gray-200">{children}</ul>
);

export const KeyValueItem = ({
  label,
  value,
  isLastItem,
}: {
  label?: string;
  value: string | number | React.ReactNode;
  isLastItem?: boolean;
}) => (
  <li
    className={`border-none
   ${isLastItem ? "" : "mb-2"}`}
  >
    <div className="py-2 grid grid-cols-9">
      {label && <span className="font-bold col-span-3">{label}</span>}
      <span className="ml-2 col-span-6">{value}</span>
    </div>
    {!isLastItem && <Separator />}
  </li>
);

export default KeyValueItem;
