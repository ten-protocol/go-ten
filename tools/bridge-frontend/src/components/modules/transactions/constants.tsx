import { TransactionStatus } from "@/src/types";
import { CheckIcon, Cross2Icon } from "@radix-ui/react-icons";
import { HourglassIcon } from "lucide-react";

export const statuses = [
  {
    value: TransactionStatus.Success,
    label: "Success",
    icon: CheckIcon,
    variant: "success",
  },
  {
    value: TransactionStatus.Failure,
    label: "Failure",
    icon: Cross2Icon,
    variant: "destructive",
  },
  {
    value: TransactionStatus.Pending,
    label: "Pending",
    icon: HourglassIcon,
    variant: "warning",
  },
];

export const types = [
  {
    value: "0x0",
    label: "Legacy",
    variant: "primary",
  },
  {
    value: "0x1",
    label: "Access List",
    variant: "secondary",
  },
  {
    value: "0x2",
    label: "Dynamic Fee",
    variant: "outline",
  },
  {
    value: "0x3",
    label: "Blob",
    variant: "outline",
  },
];
