import { CheckIcon, Cross2Icon } from "@repo/ui/components/shared/react-icons";

export const statuses = [
  {
    value: "0x1",
    label: "Success",
    icon: CheckIcon,
    variant: "success",
  },
  {
    value: "0x0",
    label: "Failure",
    icon: Cross2Icon,
    variant: "destructive",
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
