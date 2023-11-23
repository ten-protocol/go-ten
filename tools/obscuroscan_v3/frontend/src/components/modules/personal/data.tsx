import { CheckIcon, Cross2Icon } from "@radix-ui/react-icons";

export const statuses = [
  {
    value: "Success",
    label: "Success",
    icon: CheckIcon,
    variant: "success",
  },
  {
    value: "Failure",
    label: "Failure",
    icon: Cross2Icon,
    variant: "destructive",
  },
];

export const toolbar = [
  {
    column: "status",
    title: "Status",
    options: statuses,
  },
];
